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
	_ datasource.DataSource              = &incidentDataSource{}
	_ datasource.DataSourceWithConfigure = &incidentDataSource{}
)

type incidentDataSource struct {
	client *apiclient.Client
}

type incidentDataSourceModel struct {
	ID                                        types.String `tfsdk:"id"`
	PageID                                    types.String `tfsdk:"page_id"`
	Name                                      types.String `tfsdk:"name"`
	Status                                    types.String `tfsdk:"status"`
	ImpactOverride                            types.String `tfsdk:"impact_override"`
	Body                                      types.String `tfsdk:"body"`
	ComponentIDs                              types.List   `tfsdk:"component_ids"`
	Components                                types.List   `tfsdk:"components"`
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

func NewIncidentDataSource() datasource.DataSource {
	return &incidentDataSource{}
}

func (d *incidentDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_incident"
}

func (d *incidentDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a Statuspage incident.",
		Attributes: map[string]schema.Attribute{
			"id":               schema.StringAttribute{Required: true},
			"page_id":          schema.StringAttribute{Required: true},
			"name":             schema.StringAttribute{Computed: true, Description: "Incident name."},
			"status":           schema.StringAttribute{Computed: true},
			"impact_override":  schema.StringAttribute{Computed: true},
			"body":             schema.StringAttribute{Computed: true, Description: "The initial message, created as the first incident update."},
			"component_ids":    schema.ListAttribute{Computed: true, ElementType: types.StringType},
			"components": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of components affected by this incident.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":                    schema.StringAttribute{Computed: true},
						"page_id":               schema.StringAttribute{Computed: true},
						"group_id":              schema.StringAttribute{Computed: true},
						"created_at":            schema.StringAttribute{Computed: true},
						"updated_at":            schema.StringAttribute{Computed: true},
						"group":                 schema.BoolAttribute{Computed: true},
						"name":                  schema.StringAttribute{Computed: true},
						"description":           schema.StringAttribute{Computed: true},
						"position":              schema.Int64Attribute{Computed: true},
						"status":                schema.StringAttribute{Computed: true},
						"showcase":              schema.BoolAttribute{Computed: true},
						"only_show_if_degraded": schema.BoolAttribute{Computed: true},
						"automation_email":       schema.StringAttribute{Computed: true},
						"start_date":            schema.StringAttribute{Computed: true},
					},
				},
			},
			"scheduled_for":    schema.StringAttribute{Computed: true, Description: "The timestamp the maintenance is scheduled for (ISO 8601)."},
			"scheduled_until":  schema.StringAttribute{Computed: true, Description: "The timestamp the maintenance is scheduled until (ISO 8601)."},
			"scheduled_remind_prior":                        schema.BoolAttribute{Computed: true},
			"scheduled_auto_in_progress":                    schema.BoolAttribute{Computed: true},
			"scheduled_auto_completed":                      schema.BoolAttribute{Computed: true},
			"auto_transition_to_maintenance_state":          schema.BoolAttribute{Computed: true},
			"auto_transition_to_operational_state":          schema.BoolAttribute{Computed: true},
			"auto_transition_deliver_notifications_at_start": schema.BoolAttribute{Computed: true},
			"auto_transition_deliver_notifications_at_end":   schema.BoolAttribute{Computed: true},
			"deliver_notifications":                         schema.BoolAttribute{Computed: true},
			"auto_tweet_at_beginning":                       schema.BoolAttribute{Computed: true},
			"auto_tweet_on_creation":                        schema.BoolAttribute{Computed: true},
			"auto_tweet_on_completion":                      schema.BoolAttribute{Computed: true},
			"auto_tweet_one_hour_before":                    schema.BoolAttribute{Computed: true},
			"shortlink":                                     schema.StringAttribute{Computed: true},
		},
	}
}

func (d *incidentDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *incidentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config incidentDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	incident, err := d.client.GetIncident(ctx, config.PageID.ValueString(), config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading incident", err.Error())
		return
	}

	d.mapToState(ctx, &config, incident)
	resp.Diagnostics.Append(resp.State.Set(ctx, config)...)
}

func (d *incidentDataSource) mapToState(ctx context.Context, state *incidentDataSourceModel, i *apiclient.Incident) {
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
	} else {
		state.ComponentIDs = types.ListNull(types.StringType)
	}
	state.Components = mapComponentsToList(ctx, i.Components)
}
