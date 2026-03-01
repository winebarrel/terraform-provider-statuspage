package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/winebarrel/terraform-provider-statuspage/internal/apiclient"
)

var (
	_ resource.Resource                = &pageResource{}
	_ resource.ResourceWithConfigure   = &pageResource{}
	_ resource.ResourceWithImportState = &pageResource{}
)

type pageResource struct {
	client *apiclient.Client
}

type pageResourceModel struct {
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

func NewPageResource() resource.Resource {
	return &pageResource{}
}

func (r *pageResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_page"
}

func (r *pageResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Statuspage page. Pages cannot be created or deleted via the API; use terraform import.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "Page identifier.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Optional: true, Computed: true,
			},
			"page_description": schema.StringAttribute{
				Optional: true, Computed: true,
			},
			"headline": schema.StringAttribute{
				Optional: true, Computed: true,
			},
			"branding": schema.StringAttribute{
				Optional: true, Computed: true,
			},
			"subdomain": schema.StringAttribute{
				Optional: true, Computed: true,
			},
			"domain": schema.StringAttribute{
				Optional: true, Computed: true,
			},
			"url": schema.StringAttribute{
				Optional: true, Computed: true,
			},
			"support_url": schema.StringAttribute{
				Optional: true, Computed: true,
			},
			"hidden_from_search": schema.BoolAttribute{
				Optional: true, Computed: true,
			},
			"allow_page_subscribers": schema.BoolAttribute{
				Optional: true, Computed: true,
			},
			"allow_incident_subscribers": schema.BoolAttribute{
				Optional: true, Computed: true,
			},
			"allow_email_subscribers": schema.BoolAttribute{
				Optional: true, Computed: true,
			},
			"allow_sms_subscribers": schema.BoolAttribute{
				Optional: true, Computed: true,
			},
			"allow_rss_atom_feeds": schema.BoolAttribute{
				Optional: true, Computed: true,
			},
			"allow_webhook_subscribers": schema.BoolAttribute{
				Optional: true, Computed: true,
			},
			"notifications_from_email": schema.StringAttribute{
				Optional: true, Computed: true,
			},
			"notifications_email_footer": schema.StringAttribute{
				Optional: true, Computed: true,
			},
			"time_zone": schema.StringAttribute{
				Optional: true, Computed: true,
			},
			"city": schema.StringAttribute{
				Optional: true, Computed: true,
			},
			"state": schema.StringAttribute{
				Optional: true, Computed: true,
			},
			"country": schema.StringAttribute{
				Optional: true, Computed: true,
			},
		},
	}
}

func (r *pageResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	pd, ok := req.ProviderData.(*providerData)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Provider Data", fmt.Sprintf("Expected *providerData, got: %T", req.ProviderData))
		return
	}
	r.client = pd.Client
}

func (r *pageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Pages cannot be created via API. Read the existing page and set state.
	var plan pageResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	page, err := r.client.GetPage(ctx, plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading page", err.Error())
		return
	}

	r.mapToState(&plan, page)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *pageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state pageResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	page, err := r.client.GetPage(ctx, state.ID.ValueString())
	if err != nil {
		var apiErr *apiclient.APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading page", err.Error())
		return
	}

	r.mapToState(&state, page)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *pageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan pageResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := apiclient.PageBody{
		Name:                     plan.Name.ValueString(),
		PageDescription:          plan.PageDescription.ValueString(),
		Headline:                 plan.Headline.ValueString(),
		Branding:                 plan.Branding.ValueString(),
		Subdomain:                plan.Subdomain.ValueString(),
		Domain:                   plan.Domain.ValueString(),
		URL:                      plan.URL.ValueString(),
		SupportURL:               plan.SupportURL.ValueString(),
		NotificationsFromEmail:   plan.NotificationsFromEmail.ValueString(),
		NotificationsEmailFooter: plan.NotificationsEmailFooter.ValueString(),
		TimeZone:                 plan.TimeZone.ValueString(),
		City:                     plan.City.ValueString(),
		State:                    plan.State.ValueString(),
		Country:                  plan.Country.ValueString(),
	}
	setBoolPtr(&body.HiddenFromSearch, plan.HiddenFromSearch)
	setBoolPtr(&body.AllowPageSubscribers, plan.AllowPageSubscribers)
	setBoolPtr(&body.AllowIncidentSubscribers, plan.AllowIncidentSubscribers)
	setBoolPtr(&body.AllowEmailSubscribers, plan.AllowEmailSubscribers)
	setBoolPtr(&body.AllowSmsSubscribers, plan.AllowSmsSubscribers)
	setBoolPtr(&body.AllowRssAtomFeeds, plan.AllowRssAtomFeeds)
	setBoolPtr(&body.AllowWebhookSubscribers, plan.AllowWebhookSubscribers)

	page, err := r.client.UpdatePage(ctx, plan.ID.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError("Error updating page", err.Error())
		return
	}

	r.mapToState(&plan, page)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *pageResource) Delete(_ context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddWarning(
		"Page not deleted",
		"Statuspage pages cannot be deleted via the API. The page has been removed from Terraform state only.",
	)
}

func (r *pageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *pageResource) mapToState(state *pageResourceModel, p *apiclient.Page) {
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
