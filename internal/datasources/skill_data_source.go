package datasources

import (
	"context"
	"fmt"

	"github.com/akmatori/terraform-provider-akmatori/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &skillDataSource{}
	_ datasource.DataSourceWithConfigure = &skillDataSource{}
)

func NewSkillDataSource() datasource.DataSource {
	return &skillDataSource{}
}

type skillDataSource struct {
	client *client.Client
}

type skillDataSourceModel struct {
	ID          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Category    types.String `tfsdk:"category"`
	IsSystem    types.Bool   `tfsdk:"is_system"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	Prompt      types.String `tfsdk:"prompt"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
}

func (d *skillDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_skill"
}

func (d *skillDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches information about an existing skill.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the skill to look up.",
				Required:    true,
			},
			"id":          schema.Int64Attribute{Computed: true},
			"description": schema.StringAttribute{Computed: true},
			"category":    schema.StringAttribute{Computed: true},
			"is_system":   schema.BoolAttribute{Computed: true},
			"enabled":     schema.BoolAttribute{Computed: true},
			"prompt":      schema.StringAttribute{Computed: true},
			"created_at":  schema.StringAttribute{Computed: true},
			"updated_at":  schema.StringAttribute{Computed: true},
		},
	}
}

func (d *skillDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Data Source Configure Type", fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData))
		return
	}
	d.client = c
}

func (d *skillDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config skillDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	skill, err := d.client.GetSkill(config.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Skill", fmt.Sprintf("Could not read skill %q: %s", config.Name.ValueString(), err))
		return
	}

	config.ID = types.Int64Value(int64(skill.ID))
	config.Description = types.StringValue(skill.Description)
	config.Category = types.StringValue(skill.Category)
	config.IsSystem = types.BoolValue(skill.IsSystem)
	config.Enabled = types.BoolValue(skill.Enabled)
	config.Prompt = types.StringValue(skill.Prompt)
	config.CreatedAt = types.StringValue(skill.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
	config.UpdatedAt = types.StringValue(skill.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
