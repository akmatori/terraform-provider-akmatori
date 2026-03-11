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
	_ datasource.DataSource              = &toolTypesDataSource{}
	_ datasource.DataSourceWithConfigure = &toolTypesDataSource{}
)

func NewToolTypesDataSource() datasource.DataSource {
	return &toolTypesDataSource{}
}

type toolTypesDataSource struct {
	client *client.Client
}

type toolTypesDataSourceModel struct {
	ToolTypes []toolTypeItem `tfsdk:"tool_types"`
}

type toolTypeItem struct {
	ID          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	SchemaJSON  types.String `tfsdk:"schema_json"`
}

func (d *toolTypesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tool_types"
}

func (d *toolTypesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches all available tool types.",
		Attributes: map[string]schema.Attribute{
			"tool_types": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":          schema.Int64Attribute{Computed: true},
						"name":        schema.StringAttribute{Computed: true},
						"description": schema.StringAttribute{Computed: true},
						"schema_json": schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *toolTypesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *toolTypesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	apiTypes, err := d.client.GetToolTypes()
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Tool Types", err.Error())
		return
	}

	var state toolTypesDataSourceModel
	for _, t := range apiTypes {
		item := toolTypeItem{
			ID:          types.Int64Value(int64(t.ID)),
			Name:        types.StringValue(t.Name),
			Description: types.StringValue(t.Description),
		}
		if t.Schema != nil {
			if raw, err := json.Marshal(t.Schema); err == nil {
				item.SchemaJSON = types.StringValue(string(raw))
			}
		} else {
			item.SchemaJSON = types.StringNull()
		}
		state.ToolTypes = append(state.ToolTypes, item)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
