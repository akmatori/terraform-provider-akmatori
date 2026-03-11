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
	_ resource.Resource                = &settingsProxyResource{}
	_ resource.ResourceWithConfigure   = &settingsProxyResource{}
	_ resource.ResourceWithImportState = &settingsProxyResource{}
)

func NewSettingsProxyResource() resource.Resource {
	return &settingsProxyResource{}
}

type settingsProxyResource struct {
	client *client.Client
}

type settingsProxyResourceModel struct {
	ProxyURL      types.String `tfsdk:"proxy_url"`
	NoProxy       types.String `tfsdk:"no_proxy"`
	OpenAIEnabled types.Bool   `tfsdk:"openai_enabled"`
	SlackEnabled  types.Bool   `tfsdk:"slack_enabled"`
	ZabbixEnabled types.Bool   `tfsdk:"zabbix_enabled"`
}

func (r *settingsProxyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_settings_proxy"
}

func (r *settingsProxyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages proxy settings. Singleton resource — create adopts, delete resets.",
		Attributes: map[string]schema.Attribute{
			"proxy_url": schema.StringAttribute{
				Description: "HTTP/HTTPS proxy URL.",
				Optional:    true,
				Computed:    true,
			},
			"no_proxy": schema.StringAttribute{
				Description: "Comma-separated hosts to bypass proxy.",
				Optional:    true,
				Computed:    true,
			},
			"openai_enabled": schema.BoolAttribute{
				Description: "Use proxy for OpenAI API.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"slack_enabled": schema.BoolAttribute{
				Description: "Use proxy for Slack.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"zabbix_enabled": schema.BoolAttribute{
				Description: "Use proxy for Zabbix API.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
		},
	}
}

func (r *settingsProxyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *settingsProxyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan settingsProxyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := buildProxyUpdateRequest(&plan)
	settings, err := r.client.UpdateProxySettings(updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error Creating Proxy Settings", err.Error())
		return
	}

	flattenProxySettings(settings, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *settingsProxyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state settingsProxyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	settings, err := r.client.GetProxySettings()
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Proxy Settings", err.Error())
		return
	}

	// Preserve proxy_url from state since API masks password
	prevURL := state.ProxyURL
	flattenProxySettings(settings, &state)
	if !prevURL.IsNull() {
		state.ProxyURL = prevURL
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *settingsProxyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan settingsProxyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := buildProxyUpdateRequest(&plan)
	settings, err := r.client.UpdateProxySettings(updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error Updating Proxy Settings", err.Error())
		return
	}

	flattenProxySettings(settings, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *settingsProxyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resetReq := client.UpdateProxySettingsRequest{}
	// All zero-values: empty proxy_url, no_proxy, all services disabled
	if _, err := r.client.UpdateProxySettings(resetReq); err != nil {
		resp.Diagnostics.AddError("Error Resetting Proxy Settings", err.Error())
	}
}

func (r *settingsProxyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	settings, err := r.client.GetProxySettings()
	if err != nil {
		resp.Diagnostics.AddError("Error Importing Proxy Settings", err.Error())
		return
	}

	var state settingsProxyResourceModel
	flattenProxySettings(settings, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func buildProxyUpdateRequest(m *settingsProxyResourceModel) client.UpdateProxySettingsRequest {
	var req client.UpdateProxySettingsRequest
	req.ProxyURL = m.ProxyURL.ValueString()
	req.NoProxy = m.NoProxy.ValueString()
	req.Services.OpenAI.Enabled = m.OpenAIEnabled.ValueBool()
	req.Services.Slack.Enabled = m.SlackEnabled.ValueBool()
	req.Services.Zabbix.Enabled = m.ZabbixEnabled.ValueBool()
	return req
}

func flattenProxySettings(s *client.ProxySettings, m *settingsProxyResourceModel) {
	m.ProxyURL = types.StringValue(s.ProxyURL)
	m.NoProxy = types.StringValue(s.NoProxy)
	m.OpenAIEnabled = types.BoolValue(s.Services.OpenAI.Enabled)
	m.SlackEnabled = types.BoolValue(s.Services.Slack.Enabled)
	m.ZabbixEnabled = types.BoolValue(s.Services.Zabbix.Enabled)
}
