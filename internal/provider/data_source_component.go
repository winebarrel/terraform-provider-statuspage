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
	_ datasource.DataSource              = &componentDataSource{}
	_ datasource.DataSourceWithConfigure = &componentDataSource{}
)

type componentDataSource struct {
	client *apiclient.Client
}

type componentDataSourceModel struct {
	ID                 types.String `tfsdk:"id"`
	PageID             types.String `tfsdk:"page_id"`
	GroupID            types.String `tfsdk:"group_id"`
	Name               types.String `tfsdk:"name"`
	Description        types.String `tfsdk:"description"`
	Position           types.Int64  `tfsdk:"position"`
	Status             types.String `tfsdk:"status"`
	Showcase           types.Bool   `tfsdk:"showcase"`
	OnlyShowIfDegraded types.Bool   `tfsdk:"only_show_if_degraded"`
	AutomationEmail    types.String `tfsdk:"automation_email"`
	StartDate          types.String `tfsdk:"start_date"`
}

func NewComponentDataSource() datasource.DataSource {
	return &componentDataSource{}
}

func (d *componentDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_component"
}

func (d *componentDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a Statuspage component.",
		Attributes: map[string]schema.Attribute{
			"id":                   schema.StringAttribute{Required: true, Description: "Component identifier."},
			"page_id":              schema.StringAttribute{Required: true, Description: "Page identifier."},
			"group_id":             schema.StringAttribute{Computed: true, Description: "Component group identifier."},
			"name":                 schema.StringAttribute{Computed: true, Description: "Display name for the component."},
			"description":          schema.StringAttribute{Computed: true, Description: "More detailed description for the component."},
			"position":             schema.Int64Attribute{Computed: true, Description: "Order the component will appear on the page."},
			"status":               schema.StringAttribute{Computed: true, Description: "Status of the component."},
			"showcase":             schema.BoolAttribute{Computed: true, Description: "Should this component be showcased."},
			"only_show_if_degraded": schema.BoolAttribute{Computed: true, Description: "Requires a special feature flag to be enabled."},
			"automation_email":     schema.StringAttribute{Computed: true, Description: "Automation email address for the component."},
			"start_date":           schema.StringAttribute{Computed: true, Description: "The date this component started being used (YYYY-MM-DD)."},
		},
	}
}

func (d *componentDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *componentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config componentDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	component, err := d.client.GetComponent(ctx, config.PageID.ValueString(), config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading component", err.Error())
		return
	}

	d.mapToState(&config, component)
	resp.Diagnostics.Append(resp.State.Set(ctx, config)...)
}

func (d *componentDataSource) mapToState(state *componentDataSourceModel, c *apiclient.Component) {
	state.ID = types.StringValue(c.ID)
	state.PageID = types.StringValue(c.PageID)
	state.Name = types.StringValue(c.Name)
	state.Description = types.StringValue(c.Description)
	state.Position = types.Int64Value(int64(c.Position))
	state.Status = types.StringValue(c.Status)
	state.Showcase = types.BoolValue(c.Showcase)
	state.OnlyShowIfDegraded = types.BoolValue(c.OnlyShowIfDegraded)
	state.AutomationEmail = types.StringValue(c.AutomationEmail)
	state.StartDate = types.StringValue(c.StartDate)
	if c.GroupID != "" {
		state.GroupID = types.StringValue(c.GroupID)
	} else {
		state.GroupID = types.StringNull()
	}
}
