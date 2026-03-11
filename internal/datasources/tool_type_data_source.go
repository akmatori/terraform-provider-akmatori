package datasources

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/akmatori/terraform-provider-akmatori/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &toolTypeDataSource{}
	_ datasource.DataSourceWithConfigure = &toolTypeDataSource{}
)

func NewToolTypeDataSource() datasource.DataSource {
	return &toolTypeDataSource{}
}

type toolTypeDataSource struct {
	client *client.Client
}

type toolTypeDataSourceModel struct {
	Name        types.String `tfsdk:"name"`
	ID          types.Int64  `tfsdk:"id"`
	Description types.String `tfsdk:"description"`
	SchemaJSON  types.String `tfsdk:"schema_json"`
}

func (d *toolTypeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tool_type"
}

func (d *toolTypeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches information about a tool type by name.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the tool type to look up.",
				Required:    true,
			},
			"id":          schema.Int64Attribute{Computed: true},
			"description": schema.StringAttribute{Computed: true},
			"schema_json": schema.StringAttribute{Computed: true},
		},
	}
}

func (d *toolTypeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *toolTypeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config toolTypeDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	types_, err := d.client.GetToolTypes()
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Tool Types", err.Error())
		return
	}

	name := config.Name.ValueString()
	for _, t := range types_ {
		if t.Name == name {
			config.ID = types.Int64Value(int64(t.ID))
			config.Description = types.StringValue(t.Description)
			if t.Schema != nil {
				if raw, err := json.Marshal(t.Schema); err == nil {
					config.SchemaJSON = types.StringValue(string(raw))
				}
			} else {
				config.SchemaJSON = types.StringNull()
			}
			resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
			return
		}
	}

	resp.Diagnostics.AddError("Tool Type Not Found", fmt.Sprintf("No tool type found with name %q", name))
}
