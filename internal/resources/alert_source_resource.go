package resources

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/akmatori/terraform-provider-akmatori/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &alertSourceResource{}
	_ resource.ResourceWithConfigure   = &alertSourceResource{}
	_ resource.ResourceWithImportState = &alertSourceResource{}
)

func NewAlertSourceResource() resource.Resource {
	return &alertSourceResource{}
}

type alertSourceResource struct {
	client *client.Client
}

type alertSourceResourceModel struct {
	UUID             types.String         `tfsdk:"uuid"`
	ID               types.Int64          `tfsdk:"id"`
	SourceTypeName   types.String         `tfsdk:"source_type_name"`
	Name             types.String         `tfsdk:"name"`
	Description      types.String         `tfsdk:"description"`
	WebhookSecret    types.String         `tfsdk:"webhook_secret"`
	WebhookURL       types.String         `tfsdk:"webhook_url"`
	FieldMappingsJSON jsontypes.Normalized `tfsdk:"field_mappings_json"`
	SettingsJSON     jsontypes.Normalized `tfsdk:"settings_json"`
	Enabled          types.Bool           `tfsdk:"enabled"`
	CreatedAt        types.String         `tfsdk:"created_at"`
	UpdatedAt        types.String         `tfsdk:"updated_at"`
}

func (r *alertSourceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alert_source"
}

func (r *alertSourceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an alert source instance.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Description: "UUID of the alert source.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.Int64Attribute{
				Description: "Numeric identifier of the alert source.",
				Computed:    true,
			},
			"source_type_name": schema.StringAttribute{
				Description: "The name of the alert source type (e.g., alertmanager, pagerduty).",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "User-friendly name of the alert source.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the alert source.",
				Optional:    true,
				Computed:    true,
			},
			"webhook_secret": schema.StringAttribute{
				Description: "Secret used to authenticate incoming webhooks.",
				Optional:    true,
				Sensitive:   true,
			},
			"webhook_url": schema.StringAttribute{
				Description: "The webhook URL for this alert source (computed by the server).",
				Computed:    true,
			},
			"field_mappings_json": schema.StringAttribute{
				Description: "JSON-encoded field mappings for this alert source.",
				Optional:    true,
				CustomType:  jsontypes.NormalizedType{},
			},
			"settings_json": schema.StringAttribute{
				Description: "JSON-encoded settings for this alert source.",
				Optional:    true,
				CustomType:  jsontypes.NormalizedType{},
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the alert source is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"created_at": schema.StringAttribute{
				Description: "Timestamp when the alert source was created.",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Timestamp when the alert source was last updated.",
				Computed:    true,
			},
		},
	}
}

func (r *alertSourceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Provider Data", fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData))
		return
	}
	r.client = c
}

func (r *alertSourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan alertSourceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := client.CreateAlertSourceRequest{
		SourceTypeName: plan.SourceTypeName.ValueString(),
		Name:           plan.Name.ValueString(),
		Description:    plan.Description.ValueString(),
		WebhookSecret:  plan.WebhookSecret.ValueString(),
	}

	if !plan.FieldMappingsJSON.IsNull() && !plan.FieldMappingsJSON.IsUnknown() {
		var fm any
		if err := json.Unmarshal([]byte(plan.FieldMappingsJSON.ValueString()), &fm); err != nil {
			resp.Diagnostics.AddError("Invalid field_mappings_json", err.Error())
			return
		}
		createReq.FieldMappings = fm
	}

	if !plan.SettingsJSON.IsNull() && !plan.SettingsJSON.IsUnknown() {
		var s any
		if err := json.Unmarshal([]byte(plan.SettingsJSON.ValueString()), &s); err != nil {
			resp.Diagnostics.AddError("Invalid settings_json", err.Error())
			return
		}
		createReq.Settings = s
	}

	source, err := r.client.CreateAlertSource(createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error Creating Alert Source", err.Error())
		return
	}

	flattenAlertSource(source, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *alertSourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state alertSourceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	source, err := r.client.GetAlertSource(state.UUID.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error Reading Alert Source", err.Error())
		return
	}

	prevSecret := state.WebhookSecret
	flattenAlertSource(source, &state)
	// Preserve webhook_secret from state since API returns masked values
	if !prevSecret.IsNull() {
		state.WebhookSecret = prevSecret
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *alertSourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan alertSourceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state alertSourceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := plan.Name.ValueString()
	desc := plan.Description.ValueString()
	enabled := plan.Enabled.ValueBool()

	updateReq := client.UpdateAlertSourceRequest{
		Name:        &name,
		Description: &desc,
		Enabled:     &enabled,
	}

	if !plan.WebhookSecret.IsNull() {
		ws := plan.WebhookSecret.ValueString()
		updateReq.WebhookSecret = &ws
	}

	if !plan.FieldMappingsJSON.IsNull() && !plan.FieldMappingsJSON.IsUnknown() {
		var fm any
		if err := json.Unmarshal([]byte(plan.FieldMappingsJSON.ValueString()), &fm); err != nil {
			resp.Diagnostics.AddError("Invalid field_mappings_json", err.Error())
			return
		}
		updateReq.FieldMappings = fm
	}

	if !plan.SettingsJSON.IsNull() && !plan.SettingsJSON.IsUnknown() {
		var s any
		if err := json.Unmarshal([]byte(plan.SettingsJSON.ValueString()), &s); err != nil {
			resp.Diagnostics.AddError("Invalid settings_json", err.Error())
			return
		}
		updateReq.Settings = s
	}

	source, err := r.client.UpdateAlertSource(state.UUID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error Updating Alert Source", err.Error())
		return
	}

	flattenAlertSource(source, &plan)
	plan.UUID = state.UUID
	// Preserve webhook_secret from plan
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *alertSourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state alertSourceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteAlertSource(state.UUID.ValueString()); err != nil {
		if client.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error Deleting Alert Source", err.Error())
	}
}

func (r *alertSourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	source, err := r.client.GetAlertSource(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Error Importing Alert Source", err.Error())
		return
	}

	var state alertSourceResourceModel
	flattenAlertSource(source, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func flattenAlertSource(s *client.AlertSource, m *alertSourceResourceModel) {
	m.UUID = types.StringValue(s.UUID)
	m.ID = types.Int64Value(int64(s.ID))
	m.SourceTypeName = types.StringValue(s.SourceTypeName)
	m.Name = types.StringValue(s.Name)
	m.Description = types.StringValue(s.Description)
	m.WebhookURL = types.StringValue(s.WebhookURL)
	m.Enabled = types.BoolValue(s.Enabled)
	m.CreatedAt = types.StringValue(s.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
	m.UpdatedAt = types.StringValue(s.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))

	// Webhook secret from API is masked, so we don't overwrite it here
	if s.WebhookSecret != "" {
		m.WebhookSecret = types.StringValue(s.WebhookSecret)
	}

	if s.FieldMappings != nil {
		if raw, err := json.Marshal(s.FieldMappings); err == nil {
			m.FieldMappingsJSON = jsontypes.NewNormalizedValue(string(raw))
		}
	} else {
		m.FieldMappingsJSON = jsontypes.NewNormalizedNull()
	}

	if s.Settings != nil {
		if raw, err := json.Marshal(s.Settings); err == nil {
			m.SettingsJSON = jsontypes.NewNormalizedValue(string(raw))
		}
	} else {
		m.SettingsJSON = jsontypes.NewNormalizedNull()
	}
}
