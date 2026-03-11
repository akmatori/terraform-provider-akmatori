package resources

import (
	"context"
	"fmt"

	"github.com/akmatori/terraform-provider-akmatori/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &settingsLLMResource{}
	_ resource.ResourceWithConfigure   = &settingsLLMResource{}
	_ resource.ResourceWithImportState = &settingsLLMResource{}
)

func NewSettingsLLMResource() resource.Resource {
	return &settingsLLMResource{}
}

type settingsLLMResource struct {
	client *client.Client
}

type settingsLLMResourceModel struct {
	Provider      types.String `tfsdk:"llm_provider"`
	APIKey        types.String `tfsdk:"api_key"`
	Model         types.String `tfsdk:"model"`
	ThinkingLevel types.String `tfsdk:"thinking_level"`
	BaseURL       types.String `tfsdk:"base_url"`
}

func (r *settingsLLMResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_settings_llm"
}

func (r *settingsLLMResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages LLM provider settings. Singleton per provider — create adopts, delete resets.",
		Attributes: map[string]schema.Attribute{
			"llm_provider": schema.StringAttribute{
				Description: "LLM provider name (openai, anthropic, google, openrouter, custom).",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"api_key": schema.StringAttribute{
				Description: "API key for the LLM provider.",
				Optional:    true,
				Sensitive:   true,
			},
			"model": schema.StringAttribute{
				Description: "Model name to use.",
				Optional:    true,
				Computed:    true,
			},
			"thinking_level": schema.StringAttribute{
				Description: "Thinking/reasoning level (off, minimal, low, medium, high, xhigh).",
				Optional:    true,
				Computed:    true,
			},
			"base_url": schema.StringAttribute{
				Description: "Custom base URL for the LLM API.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *settingsLLMResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *settingsLLMResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan settingsLLMResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := buildLLMUpdateRequest(&plan)
	settings, err := r.client.UpdateLLMSettings(updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error Creating LLM Settings", err.Error())
		return
	}

	flattenLLMSettings(settings, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *settingsLLMResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state settingsLLMResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	settings, err := r.client.GetLLMSettings()
	if err != nil {
		resp.Diagnostics.AddError("Error Reading LLM Settings", err.Error())
		return
	}

	// Update non-sensitive fields; preserve api_key from state
	prevAPIKey := state.APIKey
	state.Provider = types.StringValue(settings.Provider)
	state.Model = types.StringValue(settings.Model)
	state.ThinkingLevel = types.StringValue(settings.ThinkingLevel)
	state.BaseURL = types.StringValue(settings.BaseURL)
	state.APIKey = prevAPIKey

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *settingsLLMResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan settingsLLMResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := buildLLMUpdateRequest(&plan)
	settings, err := r.client.UpdateLLMSettings(updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error Updating LLM Settings", err.Error())
		return
	}

	flattenLLMSettings(settings, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *settingsLLMResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state settingsLLMResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	empty := ""
	provider := state.Provider.ValueString()
	resetReq := client.UpdateLLMSettingsRequest{
		Provider: &provider,
		APIKey:   &empty,
		Model:    &empty,
		BaseURL:  &empty,
	}
	if _, err := r.client.UpdateLLMSettings(resetReq); err != nil {
		resp.Diagnostics.AddError("Error Resetting LLM Settings", err.Error())
	}
}

func (r *settingsLLMResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	settings, err := r.client.GetLLMSettings()
	if err != nil {
		resp.Diagnostics.AddError("Error Importing LLM Settings", err.Error())
		return
	}

	state := settingsLLMResourceModel{
		Provider:      types.StringValue(settings.Provider),
		APIKey:        types.StringValue(settings.APIKey),
		Model:         types.StringValue(settings.Model),
		ThinkingLevel: types.StringValue(settings.ThinkingLevel),
		BaseURL:       types.StringValue(settings.BaseURL),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func buildLLMUpdateRequest(m *settingsLLMResourceModel) client.UpdateLLMSettingsRequest {
	req := client.UpdateLLMSettingsRequest{}
	if !m.Provider.IsNull() {
		v := m.Provider.ValueString()
		req.Provider = &v
	}
	if !m.APIKey.IsNull() {
		v := m.APIKey.ValueString()
		req.APIKey = &v
	}
	if !m.Model.IsNull() {
		v := m.Model.ValueString()
		req.Model = &v
	}
	if !m.ThinkingLevel.IsNull() {
		v := m.ThinkingLevel.ValueString()
		req.ThinkingLevel = &v
	}
	if !m.BaseURL.IsNull() {
		v := m.BaseURL.ValueString()
		req.BaseURL = &v
	}
	return req
}

func flattenLLMSettings(s *client.LLMSettings, m *settingsLLMResourceModel) {
	m.Provider = types.StringValue(s.Provider)
	// Preserve api_key from plan (API returns masked)
	m.Model = types.StringValue(s.Model)
	m.ThinkingLevel = types.StringValue(s.ThinkingLevel)
	m.BaseURL = types.StringValue(s.BaseURL)
}
