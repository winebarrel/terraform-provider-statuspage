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
	_ datasource.DataSource              = &pageAccessUserDataSource{}
	_ datasource.DataSourceWithConfigure = &pageAccessUserDataSource{}
)

type pageAccessUserDataSource struct {
	client *apiclient.Client
}

type pageAccessUserDataSourceModel struct {
	ID            types.String `tfsdk:"id"`
	PageID        types.String `tfsdk:"page_id"`
	ExternalLogin types.String `tfsdk:"external_login"`
	ExternalEmail types.String `tfsdk:"external_email"`
	ComponentIDs  types.List   `tfsdk:"component_ids"`
	MetricIDs     types.List   `tfsdk:"metric_ids"`
}

func NewPageAccessUserDataSource() datasource.DataSource {
	return &pageAccessUserDataSource{}
}

func (d *pageAccessUserDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_page_access_user"
}

func (d *pageAccessUserDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a Statuspage page access user.",
		Attributes: map[string]schema.Attribute{
			"id":             schema.StringAttribute{Required: true},
			"page_id":        schema.StringAttribute{Required: true},
			"external_login": schema.StringAttribute{Computed: true},
			"external_email": schema.StringAttribute{Computed: true},
			"component_ids":  schema.ListAttribute{Computed: true, ElementType: types.StringType},
			"metric_ids":     schema.ListAttribute{Computed: true, ElementType: types.StringType},
		},
	}
}

func (d *pageAccessUserDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *pageAccessUserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config pageAccessUserDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user, err := d.client.GetPageAccessUser(ctx, config.PageID.ValueString(), config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading page access user", err.Error())
		return
	}

	d.mapToState(ctx, &config, user)
	resp.Diagnostics.Append(resp.State.Set(ctx, config)...)
}

func (d *pageAccessUserDataSource) mapToState(ctx context.Context, state *pageAccessUserDataSourceModel, u *apiclient.PageAccessUser) {
	state.ID = types.StringValue(u.ID)
	state.PageID = types.StringValue(u.PageID)
	state.ExternalLogin = types.StringValue(u.ExternalLogin)
	state.ExternalEmail = types.StringValue(u.ExternalEmail)
	if len(u.ComponentIDs) > 0 {
		componentIDs, _ := types.ListValueFrom(ctx, types.StringType, u.ComponentIDs)
		state.ComponentIDs = componentIDs
	} else {
		state.ComponentIDs = types.ListNull(types.StringType)
	}
	if len(u.MetricIDs) > 0 {
		metricIDs, _ := types.ListValueFrom(ctx, types.StringType, u.MetricIDs)
		state.MetricIDs = metricIDs
	} else {
		state.MetricIDs = types.ListNull(types.StringType)
	}
}
