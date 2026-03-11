package resources

import (
	"context"
	"fmt"

	"github.com/akmatori/terraform-provider-akmatori/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &settingsSlackResource{}
	_ resource.ResourceWithConfigure   = &settingsSlackResource{}
	_ resource.ResourceWithImportState = &settingsSlackResource{}
)

func NewSettingsSlackResource() resource.Resource {
	return &settingsSlackResource{}
}

type settingsSlackResource struct {
	client *client.Client
}

type settingsSlackResourceModel struct {
	BotToken      types.String `tfsdk:"bot_token"`
	SigningSecret types.String `tfsdk:"signing_secret"`
	AppToken      types.String `tfsdk:"app_token"`
	AlertsChannel types.String `tfsdk:"alerts_channel"`
	Enabled       types.Bool   `tfsdk:"enabled"`
}

func (r *settingsSlackResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_settings_slack"
}

func (r *settingsSlackResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages Slack integration settings. Singleton resource — create adopts existing settings, delete resets to defaults.",
		Attributes: map[string]schema.Attribute{
			"bot_token": schema.StringAttribute{
				Description: "Slack bot token.",
				Optional:    true,
				Sensitive:   true,
			},
			"signing_secret": schema.StringAttribute{
				Description: "Slack signing secret.",
				Optional:    true,
				Sensitive:   true,
			},
			"app_token": schema.StringAttribute{
				Description: "Slack app-level token.",
				Optional:    true,
				Sensitive:   true,
			},
			"alerts_channel": schema.StringAttribute{
				Description: "Slack channel for alerts.",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether Slack integration is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
		},
	}
}

func (r *settingsSlackResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *settingsSlackResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan settingsSlackResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := buildSlackUpdateRequest(&plan)
	_, err := r.client.UpdateSlackSettings(updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error Creating Slack Settings", err.Error())
		return
	}

	// State preserves the raw values (not masked API response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *settingsSlackResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state settingsSlackResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	settings, err := r.client.GetSlackSettings()
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Slack Settings", err.Error())
		return
	}

	// Only update non-sensitive fields from API; sensitive fields preserved from state
	state.AlertsChannel = types.StringValue(settings.AlertsChannel)
	state.Enabled = types.BoolValue(settings.Enabled)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *settingsSlackResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan settingsSlackResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := buildSlackUpdateRequest(&plan)
	_, err := r.client.UpdateSlackSettings(updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error Updating Slack Settings", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *settingsSlackResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Reset to defaults
	empty := ""
	disabled := false
	resetReq := client.UpdateSlackSettingsRequest{
		BotToken:      &empty,
		SigningSecret: &empty,
		AppToken:      &empty,
		AlertsChannel: &empty,
		Enabled:       &disabled,
	}
	if _, err := r.client.UpdateSlackSettings(resetReq); err != nil {
		resp.Diagnostics.AddError("Error Resetting Slack Settings", err.Error())
	}
}

func (r *settingsSlackResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	settings, err := r.client.GetSlackSettings()
	if err != nil {
		resp.Diagnostics.AddError("Error Importing Slack Settings", err.Error())
		return
	}

	state := settingsSlackResourceModel{
		BotToken:      types.StringValue(settings.BotToken),
		SigningSecret: types.StringValue(settings.SigningSecret),
		AppToken:      types.StringValue(settings.AppToken),
		AlertsChannel: types.StringValue(settings.AlertsChannel),
		Enabled:       types.BoolValue(settings.Enabled),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func buildSlackUpdateRequest(m *settingsSlackResourceModel) client.UpdateSlackSettingsRequest {
	req := client.UpdateSlackSettingsRequest{}
	if !m.BotToken.IsNull() {
		v := m.BotToken.ValueString()
		req.BotToken = &v
	}
	if !m.SigningSecret.IsNull() {
		v := m.SigningSecret.ValueString()
		req.SigningSecret = &v
	}
	if !m.AppToken.IsNull() {
		v := m.AppToken.ValueString()
		req.AppToken = &v
	}
	if !m.AlertsChannel.IsNull() {
		v := m.AlertsChannel.ValueString()
		req.AlertsChannel = &v
	}
	if !m.Enabled.IsNull() {
		v := m.Enabled.ValueBool()
		req.Enabled = &v
	}
	return req
}
