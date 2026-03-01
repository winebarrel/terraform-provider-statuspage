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
	_ resource.Resource                = &incidentPostmortemResource{}
	_ resource.ResourceWithConfigure   = &incidentPostmortemResource{}
	_ resource.ResourceWithImportState = &incidentPostmortemResource{}
)

type incidentPostmortemResource struct {
	client *apiclient.Client
}

type incidentPostmortemResourceModel struct {
	PageID            types.String `tfsdk:"page_id"`
	IncidentID        types.String `tfsdk:"incident_id"`
	Body              types.String `tfsdk:"body"`
	BodyDraft         types.String `tfsdk:"body_draft"`
	NotifySubscribers types.Bool   `tfsdk:"notify_subscribers"`
	NotifyTwitter     types.Bool   `tfsdk:"notify_twitter"`
}

func NewIncidentPostmortemResource() resource.Resource {
	return &incidentPostmortemResource{}
}

func (r *incidentPostmortemResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_incident_postmortem"
}

func (r *incidentPostmortemResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Statuspage incident postmortem.",
		Attributes: map[string]schema.Attribute{
			"page_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"incident_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"body": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"body_draft": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"notify_subscribers": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"notify_twitter": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
		},
	}
}

func (r *incidentPostmortemResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *incidentPostmortemResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan incidentPostmortemResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := apiclient.PostmortemBody{
		Body:      plan.Body.ValueString(),
		BodyDraft: plan.BodyDraft.ValueString(),
	}
	setBoolPtr(&body.NotifySubscribers, plan.NotifySubscribers)
	setBoolPtr(&body.NotifyTwitter, plan.NotifyTwitter)

	pm, err := r.client.CreateOrUpdatePostmortem(ctx, plan.PageID.ValueString(), plan.IncidentID.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError("Error creating postmortem", err.Error())
		return
	}

	r.mapToState(&plan, pm)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *incidentPostmortemResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state incidentPostmortemResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pm, err := r.client.GetPostmortem(ctx, state.PageID.ValueString(), state.IncidentID.ValueString())
	if err != nil {
		var apiErr *apiclient.APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading postmortem", err.Error())
		return
	}

	r.mapToState(&state, pm)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *incidentPostmortemResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan incidentPostmortemResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := apiclient.PostmortemBody{
		Body:      plan.Body.ValueString(),
		BodyDraft: plan.BodyDraft.ValueString(),
	}
	setBoolPtr(&body.NotifySubscribers, plan.NotifySubscribers)
	setBoolPtr(&body.NotifyTwitter, plan.NotifyTwitter)

	pm, err := r.client.CreateOrUpdatePostmortem(ctx, plan.PageID.ValueString(), plan.IncidentID.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError("Error updating postmortem", err.Error())
		return
	}

	r.mapToState(&plan, pm)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *incidentPostmortemResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state incidentPostmortemResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeletePostmortem(ctx, state.PageID.ValueString(), state.IncidentID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting postmortem", err.Error())
		return
	}
}

func (r *incidentPostmortemResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := splitImportID(req.ID, 2)
	if idParts == nil {
		resp.Diagnostics.AddError("Invalid import ID", "Expected format: page_id/incident_id")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("page_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("incident_id"), idParts[1])...)
}

func (r *incidentPostmortemResource) mapToState(state *incidentPostmortemResourceModel, pm *apiclient.Postmortem) {
	state.Body = types.StringValue(pm.Body)
	state.BodyDraft = types.StringValue(pm.BodyDraft)
	state.NotifySubscribers = types.BoolValue(pm.NotifySubscribers)
	state.NotifyTwitter = types.BoolValue(pm.NotifyTwitter)
}
