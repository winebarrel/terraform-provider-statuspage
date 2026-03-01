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
	_ resource.Resource                = &componentResource{}
	_ resource.ResourceWithConfigure   = &componentResource{}
	_ resource.ResourceWithImportState = &componentResource{}
)

type componentResource struct {
	client *apiclient.Client
}

type componentResourceModel struct {
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

func NewComponentResource() resource.Resource {
	return &componentResource{}
}

func (r *componentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_component"
}

func (r *componentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Statuspage component.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Component identifier.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"page_id": schema.StringAttribute{
				Required:    true,
				Description: "Page identifier.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"group_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Component group identifier.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Display name for the component.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "More detailed description for the component.",
			},
			"position": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Order the component will appear on the page.",
			},
			"status": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Status of the component.",
				Validators: []validator.String{
					stringvalidator.OneOf("operational", "under_maintenance", "degraded_performance", "partial_outage", "major_outage"),
				},
			},
			"showcase": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Should this component be showcased.",
			},
			"only_show_if_degraded": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Requires a special feature flag to be enabled.",
			},
			"automation_email": schema.StringAttribute{
				Computed:    true,
				Description: "Automation email address for the component.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"start_date": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The date this component started being used (YYYY-MM-DD).",
			},
		},
	}
}

func (r *componentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *componentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan componentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := apiclient.ComponentBody{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Status:      plan.Status.ValueString(),
		GroupID:     plan.GroupID.ValueString(),
		StartDate:   plan.StartDate.ValueString(),
	}
	if !plan.Showcase.IsNull() && !plan.Showcase.IsUnknown() {
		v := plan.Showcase.ValueBool()
		body.Showcase = &v
	}
	if !plan.OnlyShowIfDegraded.IsNull() && !plan.OnlyShowIfDegraded.IsUnknown() {
		v := plan.OnlyShowIfDegraded.ValueBool()
		body.OnlyShowIfDegraded = &v
	}

	component, err := r.client.CreateComponent(ctx, plan.PageID.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError("Error creating component", err.Error())
		return
	}

	r.mapToState(&plan, component)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *componentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state componentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	component, err := r.client.GetComponent(ctx, state.PageID.ValueString(), state.ID.ValueString())
	if err != nil {
		var apiErr *apiclient.APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading component", err.Error())
		return
	}

	r.mapToState(&state, component)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *componentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan componentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := apiclient.ComponentBody{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Status:      plan.Status.ValueString(),
		GroupID:     plan.GroupID.ValueString(),
		StartDate:   plan.StartDate.ValueString(),
	}
	if !plan.Showcase.IsNull() && !plan.Showcase.IsUnknown() {
		v := plan.Showcase.ValueBool()
		body.Showcase = &v
	}
	if !plan.OnlyShowIfDegraded.IsNull() && !plan.OnlyShowIfDegraded.IsUnknown() {
		v := plan.OnlyShowIfDegraded.ValueBool()
		body.OnlyShowIfDegraded = &v
	}

	component, err := r.client.UpdateComponent(ctx, plan.PageID.ValueString(), plan.ID.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError("Error updating component", err.Error())
		return
	}

	r.mapToState(&plan, component)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *componentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state componentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteComponent(ctx, state.PageID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting component", err.Error())
		return
	}
}

func (r *componentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := splitImportID(req.ID, 2)
	if idParts == nil {
		resp.Diagnostics.AddError("Invalid import ID", "Expected format: page_id/component_id")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("page_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)
}

func (r *componentResource) mapToState(state *componentResourceModel, c *apiclient.Component) {
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
