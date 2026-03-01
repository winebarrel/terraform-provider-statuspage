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
	_ resource.Resource                = &pageAccessGroupResource{}
	_ resource.ResourceWithConfigure   = &pageAccessGroupResource{}
	_ resource.ResourceWithImportState = &pageAccessGroupResource{}
)

type pageAccessGroupResource struct {
	client *apiclient.Client
}

type pageAccessGroupResourceModel struct {
	ID                 types.String `tfsdk:"id"`
	PageID             types.String `tfsdk:"page_id"`
	Name               types.String `tfsdk:"name"`
	ExternalIdentifier types.String `tfsdk:"external_identifier"`
	ComponentIDs       types.List   `tfsdk:"component_ids"`
	MetricIDs          types.List   `tfsdk:"metric_ids"`
	PageAccessUserIDs  types.List   `tfsdk:"page_access_user_ids"`
}

func NewPageAccessGroupResource() resource.Resource {
	return &pageAccessGroupResource{}
}

func (r *pageAccessGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_page_access_group"
}

func (r *pageAccessGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Statuspage page access group.",
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
				Required: true,
			},
			"external_identifier": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"component_ids": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"metric_ids": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"page_access_user_ids": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (r *pageAccessGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *pageAccessGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan pageAccessGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := apiclient.PageAccessGroupBody{
		Name:               plan.Name.ValueString(),
		ExternalIdentifier: plan.ExternalIdentifier.ValueString(),
	}
	r.setListFields(ctx, &plan, &body)

	group, err := r.client.CreatePageAccessGroup(ctx, plan.PageID.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError("Error creating page access group", err.Error())
		return
	}

	r.mapToState(ctx, &plan, group)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *pageAccessGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state pageAccessGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	group, err := r.client.GetPageAccessGroup(ctx, state.PageID.ValueString(), state.ID.ValueString())
	if err != nil {
		var apiErr *apiclient.APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading page access group", err.Error())
		return
	}

	r.mapToState(ctx, &state, group)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *pageAccessGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan pageAccessGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := apiclient.PageAccessGroupBody{
		Name:               plan.Name.ValueString(),
		ExternalIdentifier: plan.ExternalIdentifier.ValueString(),
	}
	r.setListFields(ctx, &plan, &body)

	group, err := r.client.UpdatePageAccessGroup(ctx, plan.PageID.ValueString(), plan.ID.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError("Error updating page access group", err.Error())
		return
	}

	r.mapToState(ctx, &plan, group)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *pageAccessGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state pageAccessGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeletePageAccessGroup(ctx, state.PageID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting page access group", err.Error())
		return
	}
}

func (r *pageAccessGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := splitImportID(req.ID, 2)
	if idParts == nil {
		resp.Diagnostics.AddError("Invalid import ID", "Expected format: page_id/page_access_group_id")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("page_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)
}

func (r *pageAccessGroupResource) setListFields(ctx context.Context, plan *pageAccessGroupResourceModel, body *apiclient.PageAccessGroupBody) {
	if !plan.ComponentIDs.IsNull() && !plan.ComponentIDs.IsUnknown() {
		var ids []string
		plan.ComponentIDs.ElementsAs(ctx, &ids, false)
		body.ComponentIDs = ids
	}
	if !plan.MetricIDs.IsNull() && !plan.MetricIDs.IsUnknown() {
		var ids []string
		plan.MetricIDs.ElementsAs(ctx, &ids, false)
		body.MetricIDs = ids
	}
	if !plan.PageAccessUserIDs.IsNull() && !plan.PageAccessUserIDs.IsUnknown() {
		var ids []string
		plan.PageAccessUserIDs.ElementsAs(ctx, &ids, false)
		body.PageAccessUserIDs = ids
	}
}

func (r *pageAccessGroupResource) mapToState(ctx context.Context, state *pageAccessGroupResourceModel, g *apiclient.PageAccessGroup) {
	state.ID = types.StringValue(g.ID)
	state.PageID = types.StringValue(g.PageID)
	state.Name = types.StringValue(g.Name)
	state.ExternalIdentifier = types.StringValue(g.ExternalIdentifier)
	if len(g.ComponentIDs) > 0 {
		componentIDs, _ := types.ListValueFrom(ctx, types.StringType, g.ComponentIDs)
		state.ComponentIDs = componentIDs
	}
	if len(g.MetricIDs) > 0 {
		metricIDs, _ := types.ListValueFrom(ctx, types.StringType, g.MetricIDs)
		state.MetricIDs = metricIDs
	}
	if len(g.PageAccessUserIDs) > 0 {
		userIDs, _ := types.ListValueFrom(ctx, types.StringType, g.PageAccessUserIDs)
		state.PageAccessUserIDs = userIDs
	}
}
