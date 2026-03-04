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
	_ datasource.DataSource              = &pageAccessGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &pageAccessGroupDataSource{}
)

type pageAccessGroupDataSource struct {
	client *apiclient.Client
}

type pageAccessGroupDataSourceModel struct {
	ID                 types.String `tfsdk:"id"`
	PageID             types.String `tfsdk:"page_id"`
	Name               types.String `tfsdk:"name"`
	ExternalIdentifier types.String `tfsdk:"external_identifier"`
	ComponentIDs       types.List   `tfsdk:"component_ids"`
	MetricIDs          types.List   `tfsdk:"metric_ids"`
	PageAccessUserIDs  types.List   `tfsdk:"page_access_user_ids"`
}

func NewPageAccessGroupDataSource() datasource.DataSource {
	return &pageAccessGroupDataSource{}
}

func (d *pageAccessGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_page_access_group"
}

func (d *pageAccessGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a Statuspage page access group.",
		Attributes: map[string]schema.Attribute{
			"id":                  schema.StringAttribute{Required: true},
			"page_id":             schema.StringAttribute{Required: true},
			"name":                schema.StringAttribute{Computed: true},
			"external_identifier": schema.StringAttribute{Computed: true},
			"component_ids":       schema.ListAttribute{Computed: true, ElementType: types.StringType},
			"metric_ids":          schema.ListAttribute{Computed: true, ElementType: types.StringType},
			"page_access_user_ids": schema.ListAttribute{Computed: true, ElementType: types.StringType},
		},
	}
}

func (d *pageAccessGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *pageAccessGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config pageAccessGroupDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	group, err := d.client.GetPageAccessGroup(ctx, config.PageID.ValueString(), config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading page access group", err.Error())
		return
	}

	d.mapToState(ctx, &config, group)
	resp.Diagnostics.Append(resp.State.Set(ctx, config)...)
}

func (d *pageAccessGroupDataSource) mapToState(ctx context.Context, state *pageAccessGroupDataSourceModel, g *apiclient.PageAccessGroup) {
	state.ID = types.StringValue(g.ID)
	state.PageID = types.StringValue(g.PageID)
	state.Name = types.StringValue(g.Name)
	state.ExternalIdentifier = types.StringValue(g.ExternalIdentifier)
	if len(g.ComponentIDs) > 0 {
		componentIDs, _ := types.ListValueFrom(ctx, types.StringType, g.ComponentIDs)
		state.ComponentIDs = componentIDs
	} else {
		state.ComponentIDs = types.ListNull(types.StringType)
	}
	if len(g.MetricIDs) > 0 {
		metricIDs, _ := types.ListValueFrom(ctx, types.StringType, g.MetricIDs)
		state.MetricIDs = metricIDs
	} else {
		state.MetricIDs = types.ListNull(types.StringType)
	}
	if len(g.PageAccessUserIDs) > 0 {
		userIDs, _ := types.ListValueFrom(ctx, types.StringType, g.PageAccessUserIDs)
		state.PageAccessUserIDs = userIDs
	} else {
		state.PageAccessUserIDs = types.ListNull(types.StringType)
	}
}
