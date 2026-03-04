package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/winebarrel/terraform-provider-statuspage/internal/apiclient"
)

var (
	_ datasource.DataSource              = &metricDataSource{}
	_ datasource.DataSourceWithConfigure = &metricDataSource{}
)

type metricDataSource struct {
	client *apiclient.Client
}

type metricDataSourceModel struct {
	ID                 types.String  `tfsdk:"id"`
	PageID             types.String  `tfsdk:"page_id"`
	MetricsProviderID  types.String  `tfsdk:"metrics_provider_id"`
	MetricIdentifier   types.String  `tfsdk:"metric_identifier"`
	Name               types.String  `tfsdk:"name"`
	Display            types.Bool    `tfsdk:"display"`
	TooltipDescription types.String  `tfsdk:"tooltip_description"`
	YAxisMin           types.Float64 `tfsdk:"y_axis_min"`
	YAxisMax           types.Float64 `tfsdk:"y_axis_max"`
	YAxisHidden        types.Bool    `tfsdk:"y_axis_hidden"`
	Suffix             types.String  `tfsdk:"suffix"`
	DecimalPlaces      types.Int64   `tfsdk:"decimal_places"`
}

func NewMetricDataSource() datasource.DataSource {
	return &metricDataSource{}
}

func (d *metricDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_metric"
}

func (d *metricDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a Statuspage metric.",
		Attributes: map[string]schema.Attribute{
			"id":                  schema.StringAttribute{Required: true},
			"page_id":             schema.StringAttribute{Required: true},
			"metrics_provider_id": schema.StringAttribute{Computed: true, Description: "The metrics provider this metric belongs to."},
			"metric_identifier":   schema.StringAttribute{Computed: true},
			"name":                schema.StringAttribute{Computed: true},
			"display":             schema.BoolAttribute{Computed: true},
			"tooltip_description": schema.StringAttribute{Computed: true},
			"y_axis_min":          schema.Float64Attribute{Computed: true},
			"y_axis_max":          schema.Float64Attribute{Computed: true},
			"y_axis_hidden":       schema.BoolAttribute{Computed: true},
			"suffix":              schema.StringAttribute{Computed: true},
			"decimal_places":      schema.Int64Attribute{Computed: true},
		},
	}
}

func (d *metricDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	pd, ok := req.ProviderData.(*providerData)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Provider Data", fmt.Sprintf("Expected *providerData, got: %T", req.ProviderData))
		return
	}
	d.client = pd.Client
}

func (d *metricDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config metricDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	metric, err := d.client.GetMetric(ctx, config.PageID.ValueString(), config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading metric", err.Error())
		return
	}

	d.mapToState(&config, metric)
	resp.Diagnostics.Append(resp.State.Set(ctx, config)...)
}

func (d *metricDataSource) mapToState(state *metricDataSourceModel, m *apiclient.Metric) {
	state.ID = types.StringValue(m.ID)
	state.PageID = types.StringValue(m.PageID)
	state.MetricsProviderID = types.StringValue(m.MetricsProviderID)
	state.MetricIdentifier = types.StringValue(m.MetricIdentifier)
	state.Name = types.StringValue(m.Name)
	state.Display = types.BoolValue(m.Display)
	state.TooltipDescription = types.StringValue(m.TooltipDescription)
	state.YAxisMin = types.Float64Value(m.YAxisMin)
	state.YAxisMax = types.Float64Value(m.YAxisMax)
	state.YAxisHidden = types.BoolValue(m.YAxisHidden)
	state.Suffix = types.StringValue(m.Suffix)
	state.DecimalPlaces = types.Int64Value(int64(m.DecimalPlaces))
}
