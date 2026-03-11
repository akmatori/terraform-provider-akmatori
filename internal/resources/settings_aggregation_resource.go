package resources

import (
	"context"
	"fmt"

	"github.com/akmatori/terraform-provider-akmatori/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &settingsAggregationResource{}
	_ resource.ResourceWithConfigure   = &settingsAggregationResource{}
	_ resource.ResourceWithImportState = &settingsAggregationResource{}
)

func NewSettingsAggregationResource() resource.Resource {
	return &settingsAggregationResource{}
}

type settingsAggregationResource struct {
	client *client.Client
}

type settingsAggregationResourceModel struct {
	Enabled                        types.Bool    `tfsdk:"enabled"`
	CorrelationConfidenceThreshold types.Float64 `tfsdk:"correlation_confidence_threshold"`
	MergeConfidenceThreshold       types.Float64 `tfsdk:"merge_confidence_threshold"`
	RecorrelationEnabled           types.Bool    `tfsdk:"recorrelation_enabled"`
	RecorrelationIntervalMinutes   types.Int64   `tfsdk:"recorrelation_interval_minutes"`
	MaxIncidentsToAnalyze          types.Int64   `tfsdk:"max_incidents_to_analyze"`
	ObservingDurationMinutes       types.Int64   `tfsdk:"observing_duration_minutes"`
	CorrelatorTimeoutSeconds       types.Int64   `tfsdk:"correlator_timeout_seconds"`
	MergeAnalyzerTimeoutSeconds    types.Int64   `tfsdk:"merge_analyzer_timeout_seconds"`
}

func (r *settingsAggregationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_settings_aggregation"
}

func (r *settingsAggregationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages alert aggregation settings. Singleton resource — create adopts, delete resets to defaults.",
		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				Description: "Whether alert aggregation is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"correlation_confidence_threshold": schema.Float64Attribute{
				Description: "Minimum confidence threshold for alert correlation.",
				Optional:    true,
				Computed:    true,
				Default:     float64default.StaticFloat64(0.70),
			},
			"merge_confidence_threshold": schema.Float64Attribute{
				Description: "Minimum confidence threshold for incident merge.",
				Optional:    true,
				Computed:    true,
				Default:     float64default.StaticFloat64(0.75),
			},
			"recorrelation_enabled": schema.BoolAttribute{
				Description: "Whether periodic recorrelation is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"recorrelation_interval_minutes": schema.Int64Attribute{
				Description: "Interval in minutes between recorrelation runs.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(3),
			},
			"max_incidents_to_analyze": schema.Int64Attribute{
				Description: "Maximum number of incidents to analyze for correlation.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(20),
			},
			"observing_duration_minutes": schema.Int64Attribute{
				Description: "Duration in minutes to observe an incident before closing.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(30),
			},
			"correlator_timeout_seconds": schema.Int64Attribute{
				Description: "Timeout in seconds for the correlator.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(5),
			},
			"merge_analyzer_timeout_seconds": schema.Int64Attribute{
				Description: "Timeout in seconds for the merge analyzer.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(30),
			},
		},
	}
}

func (r *settingsAggregationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *settingsAggregationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan settingsAggregationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := toAggregationAPI(&plan)
	settings, err := r.client.UpdateAggregationSettings(apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error Creating Aggregation Settings", err.Error())
		return
	}

	fromAggregationAPI(settings, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *settingsAggregationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state settingsAggregationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	settings, err := r.client.GetAggregationSettings()
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Aggregation Settings", err.Error())
		return
	}

	fromAggregationAPI(settings, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *settingsAggregationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan settingsAggregationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := toAggregationAPI(&plan)
	settings, err := r.client.UpdateAggregationSettings(apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error Updating Aggregation Settings", err.Error())
		return
	}

	fromAggregationAPI(settings, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *settingsAggregationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Reset to defaults
	defaults := client.AggregationSettings{
		Enabled:                        true,
		CorrelationConfidenceThreshold: 0.70,
		MergeConfidenceThreshold:       0.75,
		RecorrelationEnabled:           true,
		RecorrelationIntervalMinutes:   3,
		MaxIncidentsToAnalyze:          20,
		ObservingDurationMinutes:       30,
		CorrelatorTimeoutSeconds:       5,
		MergeAnalyzerTimeoutSeconds:    30,
	}
	if _, err := r.client.UpdateAggregationSettings(defaults); err != nil {
		resp.Diagnostics.AddError("Error Resetting Aggregation Settings", err.Error())
	}
}

func (r *settingsAggregationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	settings, err := r.client.GetAggregationSettings()
	if err != nil {
		resp.Diagnostics.AddError("Error Importing Aggregation Settings", err.Error())
		return
	}

	var state settingsAggregationResourceModel
	fromAggregationAPI(settings, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func toAggregationAPI(m *settingsAggregationResourceModel) client.AggregationSettings {
	return client.AggregationSettings{
		Enabled:                        m.Enabled.ValueBool(),
		CorrelationConfidenceThreshold: m.CorrelationConfidenceThreshold.ValueFloat64(),
		MergeConfidenceThreshold:       m.MergeConfidenceThreshold.ValueFloat64(),
		RecorrelationEnabled:           m.RecorrelationEnabled.ValueBool(),
		RecorrelationIntervalMinutes:   int(m.RecorrelationIntervalMinutes.ValueInt64()),
		MaxIncidentsToAnalyze:          int(m.MaxIncidentsToAnalyze.ValueInt64()),
		ObservingDurationMinutes:       int(m.ObservingDurationMinutes.ValueInt64()),
		CorrelatorTimeoutSeconds:       int(m.CorrelatorTimeoutSeconds.ValueInt64()),
		MergeAnalyzerTimeoutSeconds:    int(m.MergeAnalyzerTimeoutSeconds.ValueInt64()),
	}
}

func fromAggregationAPI(s *client.AggregationSettings, m *settingsAggregationResourceModel) {
	m.Enabled = types.BoolValue(s.Enabled)
	m.CorrelationConfidenceThreshold = types.Float64Value(s.CorrelationConfidenceThreshold)
	m.MergeConfidenceThreshold = types.Float64Value(s.MergeConfidenceThreshold)
	m.RecorrelationEnabled = types.BoolValue(s.RecorrelationEnabled)
	m.RecorrelationIntervalMinutes = types.Int64Value(int64(s.RecorrelationIntervalMinutes))
	m.MaxIncidentsToAnalyze = types.Int64Value(int64(s.MaxIncidentsToAnalyze))
	m.ObservingDurationMinutes = types.Int64Value(int64(s.ObservingDurationMinutes))
	m.CorrelatorTimeoutSeconds = types.Int64Value(int64(s.CorrelatorTimeoutSeconds))
	m.MergeAnalyzerTimeoutSeconds = types.Int64Value(int64(s.MergeAnalyzerTimeoutSeconds))
}
