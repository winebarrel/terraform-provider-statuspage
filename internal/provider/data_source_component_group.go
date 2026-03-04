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
	_ datasource.DataSource              = &componentGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &componentGroupDataSource{}
)

type componentGroupDataSource struct {
	client *apiclient.Client
}

type componentGroupDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	PageID      types.String `tfsdk:"page_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Components  types.List   `tfsdk:"components"`
	Position    types.Int64  `tfsdk:"position"`
}

func NewComponentGroupDataSource() datasource.DataSource {
	return &componentGroupDataSource{}
}

func (d *componentGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_component_group"
}

func (d *componentGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a Statuspage component group.",
		Attributes: map[string]schema.Attribute{
			"id":          schema.StringAttribute{Required: true, Description: "Component group identifier."},
			"page_id":     schema.StringAttribute{Required: true, Description: "Page identifier."},
			"name":        schema.StringAttribute{Computed: true, Description: "Display name for the component group."},
			"description": schema.StringAttribute{Computed: true, Description: "Description for the component group."},
			"components": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "List of component IDs in this group.",
			},
			"position": schema.Int64Attribute{Computed: true, Description: "Order the component group will appear on the page."},
		},
	}
}

func (d *componentGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *componentGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config componentGroupDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	group, err := d.client.GetComponentGroup(ctx, config.PageID.ValueString(), config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading component group", err.Error())
		return
	}

	d.mapToState(ctx, &config, group)
	resp.Diagnostics.Append(resp.State.Set(ctx, config)...)
}

func (d *componentGroupDataSource) mapToState(ctx context.Context, state *componentGroupDataSourceModel, g *apiclient.ComponentGroup) {
	state.ID = types.StringValue(g.ID)
	state.PageID = types.StringValue(g.PageID)
	state.Name = types.StringValue(g.Name)
	state.Description = types.StringValue(g.Description)
	state.Position = types.Int64Value(int64(g.Position))
	components, _ := types.ListValueFrom(ctx, types.StringType, g.Components)
	state.Components = components
}
