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
	_ datasource.DataSource              = &subscriberDataSource{}
	_ datasource.DataSourceWithConfigure = &subscriberDataSource{}
)

type subscriberDataSource struct {
	client *apiclient.Client
}

type subscriberDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	PageID       types.String `tfsdk:"page_id"`
	Email        types.String `tfsdk:"email"`
	PhoneNumber  types.String `tfsdk:"phone_number"`
	PhoneCountry types.String `tfsdk:"phone_country"`
	Endpoint     types.String `tfsdk:"endpoint"`
	Mode         types.String `tfsdk:"mode"`
	ComponentIDs types.List   `tfsdk:"component_ids"`
}

func NewSubscriberDataSource() datasource.DataSource {
	return &subscriberDataSource{}
}

func (d *subscriberDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_subscriber"
}

func (d *subscriberDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a Statuspage subscriber.",
		Attributes: map[string]schema.Attribute{
			"id":            schema.StringAttribute{Required: true},
			"page_id":       schema.StringAttribute{Required: true},
			"email":         schema.StringAttribute{Computed: true},
			"phone_number":  schema.StringAttribute{Computed: true},
			"phone_country": schema.StringAttribute{Computed: true},
			"endpoint":      schema.StringAttribute{Computed: true, Description: "Webhook endpoint URL."},
			"mode":          schema.StringAttribute{Computed: true},
			"component_ids": schema.ListAttribute{Computed: true, ElementType: types.StringType},
		},
	}
}

func (d *subscriberDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *subscriberDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config subscriberDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	subscriber, err := d.client.GetSubscriber(ctx, config.PageID.ValueString(), config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading subscriber", err.Error())
		return
	}

	d.mapToState(ctx, &config, subscriber)
	resp.Diagnostics.Append(resp.State.Set(ctx, config)...)
}

func (d *subscriberDataSource) mapToState(ctx context.Context, state *subscriberDataSourceModel, s *apiclient.Subscriber) {
	state.ID = types.StringValue(s.ID)
	state.PageID = types.StringValue(s.PageID)
	state.Email = stringValueOrNull(s.Email)
	state.PhoneNumber = stringValueOrNull(s.PhoneNumber)
	state.PhoneCountry = stringValueOrNull(s.PhoneCountry)
	state.Endpoint = stringValueOrNull(s.Endpoint)
	state.Mode = types.StringValue(s.Mode)
	if len(s.ComponentIDs) > 0 {
		componentIDs, _ := types.ListValueFrom(ctx, types.StringType, s.ComponentIDs)
		state.ComponentIDs = componentIDs
	} else {
		state.ComponentIDs = types.ListNull(types.StringType)
	}
}
