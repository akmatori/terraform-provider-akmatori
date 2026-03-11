package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/akmatori/terraform-provider-akmatori/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &toolInstanceResource{}
	_ resource.ResourceWithConfigure   = &toolInstanceResource{}
	_ resource.ResourceWithImportState = &toolInstanceResource{}
)

// NewToolInstanceResource is a helper function to simplify the provider implementation.
func NewToolInstanceResource() resource.Resource {
	return &toolInstanceResource{}
}

// toolInstanceResource is the resource implementation.
type toolInstanceResource struct {
	client *client.Client
}

// toolInstanceResourceModel maps the resource schema data.
type toolInstanceResourceModel struct {
	ID           types.Int64          `tfsdk:"id"`
	ToolTypeID   types.Int64          `tfsdk:"tool_type_id"`
	Name         types.String         `tfsdk:"name"`
	SettingsJSON jsontypes.Normalized `tfsdk:"settings_json"`
	Enabled      types.Bool           `tfsdk:"enabled"`
	ToolTypeName types.String         `tfsdk:"tool_type_name"`
	CreatedAt    types.String         `tfsdk:"created_at"`
	UpdatedAt    types.String         `tfsdk:"updated_at"`
}

// Metadata returns the resource type name.
func (r *toolInstanceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tool_instance"
}

// Schema defines the schema for the resource.
func (r *toolInstanceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a tool instance.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Numeric identifier of the tool instance.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"tool_type_id": schema.Int64Attribute{
				Description: "The ID of the tool type this instance belongs to.",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the tool instance.",
				Required:    true,
			},
			"settings_json": schema.StringAttribute{
				Description: "JSON-encoded settings for the tool instance. Sensitive because settings may contain secrets.",
				Optional:    true,
				Sensitive:   true,
				CustomType:  jsontypes.NormalizedType{},
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the tool instance is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"tool_type_name": schema.StringAttribute{
				Description: "The name of the tool type.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_at": schema.StringAttribute{
				Description: "Timestamp when the tool instance was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.StringAttribute{
				Description: "Timestamp when the tool instance was last updated.",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *toolInstanceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = c
}

// Create creates the resource and sets the initial Terraform state.
func (r *toolInstanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan toolInstanceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := client.CreateToolInstanceRequest{
		ToolTypeID: int(plan.ToolTypeID.ValueInt64()),
		Name:       plan.Name.ValueString(),
		Enabled:    plan.Enabled.ValueBool(),
	}

	if !plan.SettingsJSON.IsNull() && !plan.SettingsJSON.IsUnknown() {
		var settings any
		if err := json.Unmarshal([]byte(plan.SettingsJSON.ValueString()), &settings); err != nil {
			resp.Diagnostics.AddError(
				"Invalid settings_json",
				fmt.Sprintf("Could not parse settings_json as JSON: %s", err),
			)
			return
		}
		createReq.Settings = settings
	}

	instance, err := r.client.CreateToolInstance(createReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Tool Instance",
			fmt.Sprintf("Could not create tool instance: %s", err),
		)
		return
	}

	state := flattenToolInstance(instance, plan.SettingsJSON)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *toolInstanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state toolInstanceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	instance, err := r.client.GetToolInstance(int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Tool Instance",
			fmt.Sprintf("Could not read tool instance ID %d: %s", state.ID.ValueInt64(), err),
		)
		return
	}

	// Preserve the existing settings_json from state to avoid overwriting with masked API response.
	newState := flattenToolInstance(instance, state.SettingsJSON)

	diags = resp.State.Set(ctx, newState)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *toolInstanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan toolInstanceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state toolInstanceResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := client.UpdateToolInstanceRequest{
		Name:    plan.Name.ValueString(),
		Enabled: plan.Enabled.ValueBool(),
	}

	if !plan.SettingsJSON.IsNull() && !plan.SettingsJSON.IsUnknown() {
		var settings any
		if err := json.Unmarshal([]byte(plan.SettingsJSON.ValueString()), &settings); err != nil {
			resp.Diagnostics.AddError(
				"Invalid settings_json",
				fmt.Sprintf("Could not parse settings_json as JSON: %s", err),
			)
			return
		}
		updateReq.Settings = settings
	}

	instance, err := r.client.UpdateToolInstance(int(state.ID.ValueInt64()), updateReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Tool Instance",
			fmt.Sprintf("Could not update tool instance ID %d: %s", state.ID.ValueInt64(), err),
		)
		return
	}

	// Preserve plan settings_json to avoid overwriting with masked API response.
	newState := flattenToolInstance(instance, plan.SettingsJSON)

	diags = resp.State.Set(ctx, newState)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *toolInstanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state toolInstanceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteToolInstance(int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Tool Instance",
			fmt.Sprintf("Could not delete tool instance ID %d: %s", state.ID.ValueInt64(), err),
		)
		return
	}
}

// ImportState imports an existing resource into Terraform state.
func (r *toolInstanceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Could not parse import ID %q as integer: %s", req.ID, err),
		)
		return
	}

	instance, err := r.client.GetToolInstance(id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Importing Tool Instance",
			fmt.Sprintf("Could not read tool instance ID %d: %s", id, err),
		)
		return
	}

	// On import, settings_json is unknown — we use Null to signal no prior state.
	// The user will need to set it manually if needed, or accept the masked values.
	state := flattenToolInstance(instance, jsontypes.NewNormalizedNull())

	diags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// flattenToolInstance converts a client.ToolInstance into the resource model.
// prevSettingsJSON is the previously known settings_json value; it is preserved
// as-is because the API may return masked values for sensitive fields.
func flattenToolInstance(instance *client.ToolInstance, prevSettingsJSON jsontypes.Normalized) toolInstanceResourceModel {
	model := toolInstanceResourceModel{
		ID:           types.Int64Value(int64(instance.ID)),
		ToolTypeID:   types.Int64Value(int64(instance.ToolTypeID)),
		Name:         types.StringValue(instance.Name),
		Enabled:      types.BoolValue(instance.Enabled),
		ToolTypeName: types.StringValue(instance.ToolTypeName),
		CreatedAt:    types.StringValue(instance.CreatedAt.String()),
		UpdatedAt:    types.StringValue(instance.UpdatedAt.String()),
		SettingsJSON: prevSettingsJSON,
	}

	// If there is no prior state for settings_json (e.g. first import), attempt to
	// populate it from the API response so the user sees something meaningful.
	if prevSettingsJSON.IsNull() && instance.Settings != nil {
		if raw, err := json.Marshal(instance.Settings); err == nil {
			model.SettingsJSON = jsontypes.NewNormalizedValue(string(raw))
		}
	}

	return model
}
