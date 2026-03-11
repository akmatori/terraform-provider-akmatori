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
	_ datasource.DataSource              = &alertSourceTypesDataSource{}
	_ datasource.DataSourceWithConfigure = &alertSourceTypesDataSource{}
)

func NewAlertSourceTypesDataSource() datasource.DataSource {
	return &alertSourceTypesDataSource{}
}

type alertSourceTypesDataSource struct {
	client *client.Client
}

type alertSourceTypesDataSourceModel struct {
	AlertSourceTypes []alertSourceTypeItem `tfsdk:"alert_source_types"`
}

type alertSourceTypeItem struct {
	ID                   types.Int64  `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	DisplayName          types.String `tfsdk:"display_name"`
	Description          types.String `tfsdk:"description"`
	DefaultFieldMappings types.String `tfsdk:"default_field_mappings_json"`
	WebhookSecretHeader  types.String `tfsdk:"webhook_secret_header"`
}

func (d *alertSourceTypesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alert_source_types"
}

func (d *alertSourceTypesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches all available alert source types.",
		Attributes: map[string]schema.Attribute{
			"alert_source_types": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":                          schema.Int64Attribute{Computed: true},
						"name":                        schema.StringAttribute{Computed: true},
						"display_name":                schema.StringAttribute{Computed: true},
						"description":                 schema.StringAttribute{Computed: true},
						"default_field_mappings_json": schema.StringAttribute{Computed: true},
						"webhook_secret_header":       schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *alertSourceTypesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *alertSourceTypesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	apiTypes, err := d.client.GetAlertSourceTypes()
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Alert Source Types", err.Error())
		return
	}

	var state alertSourceTypesDataSourceModel
	for _, t := range apiTypes {
		item := alertSourceTypeItem{
			ID:                  types.Int64Value(int64(t.ID)),
			Name:                types.StringValue(t.Name),
			DisplayName:         types.StringValue(t.DisplayName),
			Description:         types.StringValue(t.Description),
			WebhookSecretHeader: types.StringValue(t.WebhookSecretHeader),
		}
		if t.DefaultFieldMappings != nil {
			if raw, err := json.Marshal(t.DefaultFieldMappings); err == nil {
				item.DefaultFieldMappings = types.StringValue(string(raw))
			}
		} else {
			item.DefaultFieldMappings = types.StringNull()
		}
		state.AlertSourceTypes = append(state.AlertSourceTypes, item)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
