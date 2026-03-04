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
	_ datasource.DataSource              = &pageDataSource{}
	_ datasource.DataSourceWithConfigure = &pageDataSource{}
)

type pageDataSource struct {
	client *apiclient.Client
}

type pageDataSourceModel struct {
	ID                       types.String `tfsdk:"id"`
	Name                     types.String `tfsdk:"name"`
	PageDescription          types.String `tfsdk:"page_description"`
	Headline                 types.String `tfsdk:"headline"`
	Branding                 types.String `tfsdk:"branding"`
	Subdomain                types.String `tfsdk:"subdomain"`
	Domain                   types.String `tfsdk:"domain"`
	URL                      types.String `tfsdk:"url"`
	SupportURL               types.String `tfsdk:"support_url"`
	HiddenFromSearch         types.Bool   `tfsdk:"hidden_from_search"`
	AllowPageSubscribers     types.Bool   `tfsdk:"allow_page_subscribers"`
	AllowIncidentSubscribers types.Bool   `tfsdk:"allow_incident_subscribers"`
	AllowEmailSubscribers    types.Bool   `tfsdk:"allow_email_subscribers"`
	AllowSmsSubscribers      types.Bool   `tfsdk:"allow_sms_subscribers"`
	AllowRssAtomFeeds        types.Bool   `tfsdk:"allow_rss_atom_feeds"`
	AllowWebhookSubscribers  types.Bool   `tfsdk:"allow_webhook_subscribers"`
	NotificationsFromEmail   types.String `tfsdk:"notifications_from_email"`
	NotificationsEmailFooter types.String `tfsdk:"notifications_email_footer"`
	TimeZone                 types.String `tfsdk:"time_zone"`
	City                     types.String `tfsdk:"city"`
	State                    types.String `tfsdk:"state"`
	Country                  types.String `tfsdk:"country"`
}

func NewPageDataSource() datasource.DataSource {
	return &pageDataSource{}
}

func (d *pageDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_page"
}

func (d *pageDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a Statuspage page.",
		Attributes: map[string]schema.Attribute{
			"id":                         schema.StringAttribute{Required: true, Description: "Page identifier."},
			"name":                       schema.StringAttribute{Computed: true},
			"page_description":           schema.StringAttribute{Computed: true},
			"headline":                   schema.StringAttribute{Computed: true},
			"branding":                   schema.StringAttribute{Computed: true},
			"subdomain":                  schema.StringAttribute{Computed: true},
			"domain":                     schema.StringAttribute{Computed: true},
			"url":                        schema.StringAttribute{Computed: true},
			"support_url":                schema.StringAttribute{Computed: true},
			"hidden_from_search":         schema.BoolAttribute{Computed: true},
			"allow_page_subscribers":     schema.BoolAttribute{Computed: true},
			"allow_incident_subscribers": schema.BoolAttribute{Computed: true},
			"allow_email_subscribers":    schema.BoolAttribute{Computed: true},
			"allow_sms_subscribers":      schema.BoolAttribute{Computed: true},
			"allow_rss_atom_feeds":       schema.BoolAttribute{Computed: true},
			"allow_webhook_subscribers":  schema.BoolAttribute{Computed: true},
			"notifications_from_email":   schema.StringAttribute{Computed: true},
			"notifications_email_footer": schema.StringAttribute{Computed: true},
			"time_zone":                  schema.StringAttribute{Computed: true},
			"city":                       schema.StringAttribute{Computed: true},
			"state":                      schema.StringAttribute{Computed: true},
			"country":                    schema.StringAttribute{Computed: true},
		},
	}
}

func (d *pageDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *pageDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config pageDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	page, err := d.client.GetPage(ctx, config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading page", err.Error())
		return
	}

	d.mapToState(&config, page)
	resp.Diagnostics.Append(resp.State.Set(ctx, config)...)
}

func (d *pageDataSource) mapToState(state *pageDataSourceModel, p *apiclient.Page) {
	state.ID = types.StringValue(p.ID)
	state.Name = types.StringValue(p.Name)
	state.PageDescription = types.StringValue(p.PageDescription)
	state.Headline = types.StringValue(p.Headline)
	state.Branding = types.StringValue(p.Branding)
	state.Subdomain = types.StringValue(p.Subdomain)
	state.Domain = types.StringValue(p.Domain)
	state.URL = types.StringValue(p.URL)
	state.SupportURL = types.StringValue(p.SupportURL)
	state.HiddenFromSearch = types.BoolValue(p.HiddenFromSearch)
	state.AllowPageSubscribers = types.BoolValue(p.AllowPageSubscribers)
	state.AllowIncidentSubscribers = types.BoolValue(p.AllowIncidentSubscribers)
	state.AllowEmailSubscribers = types.BoolValue(p.AllowEmailSubscribers)
	state.AllowSmsSubscribers = types.BoolValue(p.AllowSmsSubscribers)
	state.AllowRssAtomFeeds = types.BoolValue(p.AllowRssAtomFeeds)
	state.AllowWebhookSubscribers = types.BoolValue(p.AllowWebhookSubscribers)
	state.NotificationsFromEmail = types.StringValue(p.NotificationsFromEmail)
	state.NotificationsEmailFooter = types.StringValue(p.NotificationsEmailFooter)
	state.TimeZone = types.StringValue(p.TimeZone)
	state.City = types.StringValue(p.City)
	state.State = types.StringValue(p.State)
	state.Country = types.StringValue(p.Country)
}
