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
	_ datasource.DataSource              = &metricsProviderDataSource{}
	_ datasource.DataSourceWithConfigure = &metricsProviderDataSource{}
)

type metricsProviderDataSource struct {
	client *apiclient.Client
}

type metricsProviderDataSourceModel struct {
	ID            types.String `tfsdk:"id"`
	PageID        types.String `tfsdk:"page_id"`
	Type          types.String `tfsdk:"type"`
	Email         types.String `tfsdk:"email"`
	MetricBaseURI types.String `tfsdk:"metric_base_uri"`
}

func NewMetricsProviderDataSource() datasource.DataSource {
	return &metricsProviderDataSource{}
}

func (d *metricsProviderDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_metrics_provider"
}

func (d *metricsProviderDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a Statuspage metrics provider.",
		Attributes: map[string]schema.Attribute{
			"id":             schema.StringAttribute{Required: true},
			"page_id":        schema.StringAttribute{Required: true},
			"type":           schema.StringAttribute{Computed: true, Description: "One of: Pingdom, NewRelic, Datadog, Self, Librato, CalabashCustom."},
			"email":          schema.StringAttribute{Computed: true},
			"metric_base_uri": schema.StringAttribute{Computed: true},
		},
	}
}

func (d *metricsProviderDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *metricsProviderDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config metricsProviderDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	provider, err := d.client.GetMetricsProvider(ctx, config.PageID.ValueString(), config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading metrics provider", err.Error())
		return
	}

	d.mapToState(&config, provider)
	resp.Diagnostics.Append(resp.State.Set(ctx, config)...)
}

func (d *metricsProviderDataSource) mapToState(state *metricsProviderDataSourceModel, p *apiclient.MetricsProvider) {
	state.ID = types.StringValue(p.ID)
	state.PageID = types.StringValue(p.PageID)
	state.Type = types.StringValue(p.Type)
	state.Email = types.StringValue(p.Email)
	state.MetricBaseURI = types.StringValue(p.MetricBaseURI)
}
