package provider

import (
	"context"
	"os"

	"github.com/akmatori/terraform-provider-akmatori/internal/client"
	"github.com/akmatori/terraform-provider-akmatori/internal/datasources"
	"github.com/akmatori/terraform-provider-akmatori/internal/resources"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &akmatoriProvider{}

type akmatoriProvider struct {
	version string
}

type akmatoriProviderModel struct {
	Host     types.String `tfsdk:"host"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	Token    types.String `tfsdk:"token"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &akmatoriProvider{
			version: version,
		}
	}
}

func (p *akmatoriProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "akmatori"
	resp.Version = p.version
}

func (p *akmatoriProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Terraform provider for managing Akmatori AIOps platform resources.",
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Description: "The URL of the Akmatori instance. Can also be set with the AKMATORI_HOST environment variable.",
				Optional:    true,
			},
			"username": schema.StringAttribute{
				Description: "Username for Akmatori authentication. Can also be set with the AKMATORI_USERNAME environment variable.",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "Password for Akmatori authentication. Can also be set with the AKMATORI_PASSWORD environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
			"token": schema.StringAttribute{
				Description: "JWT token for Akmatori authentication. Overrides username/password. Can also be set with the AKMATORI_TOKEN environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *akmatoriProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config akmatoriProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	host := os.Getenv("AKMATORI_HOST")
	username := os.Getenv("AKMATORI_USERNAME")
	password := os.Getenv("AKMATORI_PASSWORD")
	token := os.Getenv("AKMATORI_TOKEN")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}
	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}
	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}
	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	}

	if host == "" {
		resp.Diagnostics.AddError(
			"Missing Host",
			"The Akmatori host must be set in the provider configuration or via the AKMATORI_HOST environment variable.",
		)
		return
	}

	c, err := client.NewClient(host, token, username, password)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Akmatori Client",
			"An error occurred when creating the Akmatori API client: "+err.Error(),
		)
		return
	}

	resp.DataSourceData = c
	resp.ResourceData = c
}

func (p *akmatoriProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resources.NewSkillResource,
		resources.NewSkillToolsResource,
		resources.NewSkillScriptResource,
		resources.NewToolInstanceResource,
		resources.NewSSHKeyResource,
		resources.NewAlertSourceResource,
		resources.NewContextFileResource,
		resources.NewSettingsSlackResource,
		resources.NewSettingsLLMResource,
		resources.NewSettingsProxyResource,
		resources.NewSettingsAggregationResource,
	}
}

func (p *akmatoriProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		datasources.NewSkillDataSource,
		datasources.NewToolTypeDataSource,
		datasources.NewToolTypesDataSource,
		datasources.NewAlertSourceTypesDataSource,
	}
}
