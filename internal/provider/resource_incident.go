package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/winebarrel/terraform-provider-statuspage/internal/apiclient"
)

var (
	_ resource.Resource                = &incidentResource{}
	_ resource.ResourceWithConfigure   = &incidentResource{}
	_ resource.ResourceWithImportState = &incidentResource{}
)

type incidentResource struct {
	client *apiclient.Client
}

type incidentResourceModel struct {
	ID                                        types.String `tfsdk:"id"`
	PageID                                    types.String `tfsdk:"page_id"`
	Name                                      types.String `tfsdk:"name"`
	Status                                    types.String `tfsdk:"status"`
	ImpactOverride                            types.String `tfsdk:"impact_override"`
	Body                                      types.String `tfsdk:"body"`
	ComponentIDs                              types.List   `tfsdk:"component_ids"`
	Components                                types.Map    `tfsdk:"components"`
	ScheduledFor                              types.String `tfsdk:"scheduled_for"`
	ScheduledUntil                            types.String `tfsdk:"scheduled_until"`
	ScheduledRemindPrior                      types.Bool   `tfsdk:"scheduled_remind_prior"`
	ScheduledAutoInProgress                   types.Bool   `tfsdk:"scheduled_auto_in_progress"`
	ScheduledAutoCompleted                    types.Bool   `tfsdk:"scheduled_auto_completed"`
	AutoTransitionToMaintenanceState          types.Bool   `tfsdk:"auto_transition_to_maintenance_state"`
	AutoTransitionToOperationalState          types.Bool   `tfsdk:"auto_transition_to_operational_state"`
	AutoTransitionDeliverNotificationsAtStart types.Bool   `tfsdk:"auto_transition_deliver_notifications_at_start"`
	AutoTransitionDeliverNotificationsAtEnd   types.Bool   `tfsdk:"auto_transition_deliver_notifications_at_end"`
	DeliverNotifications                      types.Bool   `tfsdk:"deliver_notifications"`
	AutoTweetAtBeginning                      types.Bool   `tfsdk:"auto_tweet_at_beginning"`
	AutoTweetOnCreation                       types.Bool   `tfsdk:"auto_tweet_on_creation"`
	AutoTweetOnCompletion                     types.Bool   `tfsdk:"auto_tweet_on_completion"`
	AutoTweetOneHourBefore                    types.Bool   `tfsdk:"auto_tweet_one_hour_before"`
	Shortlink                                 types.String `tfsdk:"shortlink"`
}

func NewIncidentResource() resource.Resource {
	return &incidentResource{}
}

func (r *incidentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_incident"
}

func (r *incidentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Statuspage incident.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"page_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Incident name.",
			},
			"status": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("investigating", "identified", "monitoring", "resolved", "scheduled", "in_progress", "verifying", "completed"),
				},
			},
			"impact_override": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("none", "minor", "major", "critical", "maintenance"),
				},
			},
			"body": schema.StringAttribute{
				Optional:    true,
				Description: "The initial message, created as the first incident update.",
			},
			"component_ids": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"components": schema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Map of component IDs to their status.",
			},
			"scheduled_for": schema.StringAttribute{
				Optional:    true,
				Description: "The timestamp the maintenance is scheduled for (ISO 8601).",
			},
			"scheduled_until": schema.StringAttribute{
				Optional:    true,
				Description: "The timestamp the maintenance is scheduled until (ISO 8601).",
			},
			"scheduled_remind_prior": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"scheduled_auto_in_progress": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"scheduled_auto_completed": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"auto_transition_to_maintenance_state": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"auto_transition_to_operational_state": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"auto_transition_deliver_notifications_at_start": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"auto_transition_deliver_notifications_at_end": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"deliver_notifications": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
			},
			"auto_tweet_at_beginning": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"auto_tweet_on_creation": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"auto_tweet_on_completion": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"auto_tweet_one_hour_before": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"shortlink": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *incidentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *incidentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan incidentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := r.buildBody(ctx, &plan)

	incident, err := r.client.CreateIncident(ctx, plan.PageID.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError("Error creating incident", err.Error())
		return
	}

	r.mapToState(ctx, &plan, incident)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *incidentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state incidentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	incident, err := r.client.GetIncident(ctx, state.PageID.ValueString(), state.ID.ValueString())
	if err != nil {
		var apiErr *apiclient.APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading incident", err.Error())
		return
	}

	r.mapToState(ctx, &state, incident)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *incidentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan incidentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := r.buildBody(ctx, &plan)

	incident, err := r.client.UpdateIncident(ctx, plan.PageID.ValueString(), plan.ID.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError("Error updating incident", err.Error())
		return
	}

	r.mapToState(ctx, &plan, incident)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *incidentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state incidentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteIncident(ctx, state.PageID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting incident", err.Error())
		return
	}
}

func (r *incidentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := splitImportID(req.ID, 2)
	if idParts == nil {
		resp.Diagnostics.AddError("Invalid import ID", "Expected format: page_id/incident_id")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("page_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)
}

func (r *incidentResource) buildBody(ctx context.Context, plan *incidentResourceModel) apiclient.IncidentBody {
	body := apiclient.IncidentBody{
		Name:           plan.Name.ValueString(),
		Status:         plan.Status.ValueString(),
		ImpactOverride: plan.ImpactOverride.ValueString(),
		Body:           plan.Body.ValueString(),
		ScheduledFor:   plan.ScheduledFor.ValueString(),
		ScheduledUntil: plan.ScheduledUntil.ValueString(),
	}

	if !plan.ComponentIDs.IsNull() && !plan.ComponentIDs.IsUnknown() {
		var ids []string
		plan.ComponentIDs.ElementsAs(ctx, &ids, false)
		body.ComponentIDs = ids
	}

	if !plan.Components.IsNull() && !plan.Components.IsUnknown() {
		m := make(map[string]string)
		plan.Components.ElementsAs(ctx, &m, false)
		body.Components = m
	}

	setBoolPtr(&body.ScheduledRemindPrior, plan.ScheduledRemindPrior)
	setBoolPtr(&body.ScheduledAutoInProgress, plan.ScheduledAutoInProgress)
	setBoolPtr(&body.ScheduledAutoCompleted, plan.ScheduledAutoCompleted)
	setBoolPtr(&body.AutoTransitionToMaintenanceState, plan.AutoTransitionToMaintenanceState)
	setBoolPtr(&body.AutoTransitionToOperationalState, plan.AutoTransitionToOperationalState)
	setBoolPtr(&body.AutoTransitionDeliverNotificationsAtStart, plan.AutoTransitionDeliverNotificationsAtStart)
	setBoolPtr(&body.AutoTransitionDeliverNotificationsAtEnd, plan.AutoTransitionDeliverNotificationsAtEnd)
	setBoolPtr(&body.DeliverNotifications, plan.DeliverNotifications)
	setBoolPtr(&body.AutoTweetAtBeginning, plan.AutoTweetAtBeginning)
	setBoolPtr(&body.AutoTweetOnCreation, plan.AutoTweetOnCreation)
	setBoolPtr(&body.AutoTweetOnCompletion, plan.AutoTweetOnCompletion)
	setBoolPtr(&body.AutoTweetOneHourBefore, plan.AutoTweetOneHourBefore)

	return body
}

func (r *incidentResource) mapToState(ctx context.Context, state *incidentResourceModel, i *apiclient.Incident) {
	state.ID = types.StringValue(i.ID)
	state.PageID = types.StringValue(i.PageID)
	state.Name = types.StringValue(i.Name)
	state.Status = types.StringValue(i.Status)
	state.ImpactOverride = stringValueOrNull(i.ImpactOverride)
	state.Body = stringValueOrNull(i.Body)
	state.ScheduledFor = stringValueOrNull(i.ScheduledFor)
	state.ScheduledUntil = stringValueOrNull(i.ScheduledUntil)
	state.ScheduledRemindPrior = types.BoolValue(i.ScheduledRemindPrior)
	state.ScheduledAutoInProgress = types.BoolValue(i.ScheduledAutoInProgress)
	state.ScheduledAutoCompleted = types.BoolValue(i.ScheduledAutoCompleted)
	state.AutoTransitionToMaintenanceState = types.BoolValue(i.AutoTransitionToMaintenanceState)
	state.AutoTransitionToOperationalState = types.BoolValue(i.AutoTransitionToOperationalState)
	state.AutoTransitionDeliverNotificationsAtStart = types.BoolValue(i.AutoTransitionDeliverNotificationsAtStart)
	state.AutoTransitionDeliverNotificationsAtEnd = types.BoolValue(i.AutoTransitionDeliverNotificationsAtEnd)
	state.DeliverNotifications = types.BoolValue(i.DeliverNotifications)
	state.AutoTweetAtBeginning = types.BoolValue(i.AutoTweetAtBeginning)
	state.AutoTweetOnCreation = types.BoolValue(i.AutoTweetOnCreation)
	state.AutoTweetOnCompletion = types.BoolValue(i.AutoTweetOnCompletion)
	state.AutoTweetOneHourBefore = types.BoolValue(i.AutoTweetOneHourBefore)
	state.Shortlink = types.StringValue(i.Shortlink)

	if len(i.ComponentIDs) > 0 {
		componentIDs, _ := types.ListValueFrom(ctx, types.StringType, i.ComponentIDs)
		state.ComponentIDs = componentIDs
	}
	if len(i.Components) > 0 {
		components, _ := types.MapValueFrom(ctx, types.StringType, i.Components)
		state.Components = components
	}
}

func stringValueOrNull(s string) types.String {
	if s == "" {
		return types.StringNull()
	}
	return types.StringValue(s)
}

func setBoolPtr(dst **bool, src types.Bool) {
	if !src.IsNull() && !src.IsUnknown() {
		v := src.ValueBool()
		*dst = &v
	}
}
