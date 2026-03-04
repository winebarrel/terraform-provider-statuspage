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
	_ datasource.DataSource              = &incidentTemplateDataSource{}
	_ datasource.DataSourceWithConfigure = &incidentTemplateDataSource{}
)

type incidentTemplateDataSource struct {
	client *apiclient.Client
}

type incidentTemplateDataSourceModel struct {
	ID                      types.String `tfsdk:"id"`
	PageID                  types.String `tfsdk:"page_id"`
	Name                    types.String `tfsdk:"name"`
	Title                   types.String `tfsdk:"title"`
	Body                    types.String `tfsdk:"body"`
	GroupID                 types.String `tfsdk:"group_id"`
	UpdateStatus            types.String `tfsdk:"update_status"`
	ShouldTweet             types.Bool   `tfsdk:"should_tweet"`
	ShouldSendNotifications types.Bool   `tfsdk:"should_send_notifications"`
	ComponentIDs            types.List   `tfsdk:"component_ids"`
}

func NewIncidentTemplateDataSource() datasource.DataSource {
	return &incidentTemplateDataSource{}
}

func (d *incidentTemplateDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_incident_template"
}

func (d *incidentTemplateDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a Statuspage incident template.",
		Attributes: map[string]schema.Attribute{
			"id":                        schema.StringAttribute{Required: true},
			"page_id":                   schema.StringAttribute{Required: true},
			"name":                      schema.StringAttribute{Computed: true, Description: "Template name."},
			"title":                     schema.StringAttribute{Computed: true, Description: "The incident title to use when applying this template."},
			"body":                      schema.StringAttribute{Computed: true, Description: "The incident body to use when applying this template."},
			"group_id":                  schema.StringAttribute{Computed: true},
			"update_status":             schema.StringAttribute{Computed: true},
			"should_tweet":              schema.BoolAttribute{Computed: true},
			"should_send_notifications": schema.BoolAttribute{Computed: true},
			"component_ids":             schema.ListAttribute{Computed: true, ElementType: types.StringType},
		},
	}
}

func (d *incidentTemplateDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *incidentTemplateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config incidentTemplateDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tmpl, err := d.client.GetIncidentTemplate(ctx, config.PageID.ValueString(), config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading incident template", err.Error())
		return
	}

	d.mapToState(ctx, &config, tmpl)
	resp.Diagnostics.Append(resp.State.Set(ctx, config)...)
}

func (d *incidentTemplateDataSource) mapToState(ctx context.Context, state *incidentTemplateDataSourceModel, t *apiclient.IncidentTemplate) {
	state.ID = types.StringValue(t.ID)
	state.PageID = types.StringValue(t.PageID)
	state.Name = types.StringValue(t.Name)
	state.Title = types.StringValue(t.Title)
	state.Body = types.StringValue(t.Body)
	state.GroupID = types.StringValue(t.GroupID)
	state.UpdateStatus = types.StringValue(t.UpdateStatus)
	state.ShouldTweet = types.BoolValue(t.ShouldTweet)
	state.ShouldSendNotifications = types.BoolValue(t.ShouldSendNotifications)
	if len(t.ComponentIDs) > 0 {
		componentIDs, _ := types.ListValueFrom(ctx, types.StringType, t.ComponentIDs)
		state.ComponentIDs = componentIDs
	} else {
		state.ComponentIDs = types.ListNull(types.StringType)
	}
}
