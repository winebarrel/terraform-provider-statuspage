package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/winebarrel/terraform-provider-statuspage/internal/apiclient"
)

var (
	_ resource.Resource                = &incidentTemplateResource{}
	_ resource.ResourceWithConfigure   = &incidentTemplateResource{}
	_ resource.ResourceWithImportState = &incidentTemplateResource{}
)

type incidentTemplateResource struct {
	client *apiclient.Client
}

type incidentTemplateResourceModel struct {
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

func NewIncidentTemplateResource() resource.Resource {
	return &incidentTemplateResource{}
}

func (r *incidentTemplateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_incident_template"
}

func (r *incidentTemplateResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Statuspage incident template.",
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
				Description: "Template name.",
			},
			"title": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The incident title to use when applying this template.",
			},
			"body": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The incident body to use when applying this template.",
			},
			"group_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"update_status": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("investigating", "identified", "monitoring", "resolved", "scheduled", "in_progress", "verifying", "completed"),
				},
			},
			"should_tweet": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"should_send_notifications": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"component_ids": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (r *incidentTemplateResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *incidentTemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan incidentTemplateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := r.buildBody(ctx, &plan)

	tmpl, err := r.client.CreateIncidentTemplate(ctx, plan.PageID.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError("Error creating incident template", err.Error())
		return
	}

	r.mapToState(ctx, &plan, tmpl)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *incidentTemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state incidentTemplateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tmpl, err := r.client.GetIncidentTemplate(ctx, state.PageID.ValueString(), state.ID.ValueString())
	if err != nil {
		var apiErr *apiclient.APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading incident template", err.Error())
		return
	}

	r.mapToState(ctx, &state, tmpl)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *incidentTemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan incidentTemplateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := r.buildBody(ctx, &plan)

	tmpl, err := r.client.UpdateIncidentTemplate(ctx, plan.PageID.ValueString(), plan.ID.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError("Error updating incident template", err.Error())
		return
	}

	r.mapToState(ctx, &plan, tmpl)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *incidentTemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state incidentTemplateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteIncidentTemplate(ctx, state.PageID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting incident template", err.Error())
		return
	}
}

func (r *incidentTemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := splitImportID(req.ID, 2)
	if idParts == nil {
		resp.Diagnostics.AddError("Invalid import ID", "Expected format: page_id/template_id")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("page_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)
}

func (r *incidentTemplateResource) buildBody(ctx context.Context, plan *incidentTemplateResourceModel) apiclient.IncidentTemplateBody {
	body := apiclient.IncidentTemplateBody{
		Name:         plan.Name.ValueString(),
		Title:        plan.Title.ValueString(),
		Body:         plan.Body.ValueString(),
		GroupID:      plan.GroupID.ValueString(),
		UpdateStatus: plan.UpdateStatus.ValueString(),
	}
	setBoolPtr(&body.ShouldTweet, plan.ShouldTweet)
	setBoolPtr(&body.ShouldSendNotifications, plan.ShouldSendNotifications)
	if !plan.ComponentIDs.IsNull() && !plan.ComponentIDs.IsUnknown() {
		var ids []string
		plan.ComponentIDs.ElementsAs(ctx, &ids, false)
		body.ComponentIDs = ids
	}
	return body
}

func (r *incidentTemplateResource) mapToState(ctx context.Context, state *incidentTemplateResourceModel, t *apiclient.IncidentTemplate) {
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
	}
}
