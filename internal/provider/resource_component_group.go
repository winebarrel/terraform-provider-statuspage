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
	_ resource.Resource                = &componentGroupResource{}
	_ resource.ResourceWithConfigure   = &componentGroupResource{}
	_ resource.ResourceWithImportState = &componentGroupResource{}
)

type componentGroupResource struct {
	client *apiclient.Client
}

type componentGroupResourceModel struct {
	ID          types.String `tfsdk:"id"`
	PageID      types.String `tfsdk:"page_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Components  types.List   `tfsdk:"components"`
	Position    types.Int64  `tfsdk:"position"`
}

func NewComponentGroupResource() resource.Resource {
	return &componentGroupResource{}
}

func (r *componentGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_component_group"
}

func (r *componentGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Statuspage component group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Component group identifier.",
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
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Display name for the component group.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Description for the component group.",
			},
			"components": schema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "List of component IDs in this group.",
			},
			"position": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Order the component group will appear on the page.",
			},
		},
	}
}

func (r *componentGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *componentGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan componentGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var components []string
	resp.Diagnostics.Append(plan.Components.ElementsAs(ctx, &components, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := apiclient.ComponentGroupBody{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Components:  components,
	}

	group, err := r.client.CreateComponentGroup(ctx, plan.PageID.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError("Error creating component group", err.Error())
		return
	}

	r.mapToState(ctx, &plan, group)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *componentGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state componentGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	group, err := r.client.GetComponentGroup(ctx, state.PageID.ValueString(), state.ID.ValueString())
	if err != nil {
		var apiErr *apiclient.APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading component group", err.Error())
		return
	}

	r.mapToState(ctx, &state, group)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *componentGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan componentGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var components []string
	resp.Diagnostics.Append(plan.Components.ElementsAs(ctx, &components, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := apiclient.ComponentGroupBody{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Components:  components,
	}

	group, err := r.client.UpdateComponentGroup(ctx, plan.PageID.ValueString(), plan.ID.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError("Error updating component group", err.Error())
		return
	}

	r.mapToState(ctx, &plan, group)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *componentGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state componentGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteComponentGroup(ctx, state.PageID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting component group", err.Error())
		return
	}
}

func (r *componentGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := splitImportID(req.ID, 2)
	if idParts == nil {
		resp.Diagnostics.AddError("Invalid import ID", "Expected format: page_id/component_group_id")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("page_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)
}

func (r *componentGroupResource) mapToState(ctx context.Context, state *componentGroupResourceModel, g *apiclient.ComponentGroup) {
	state.ID = types.StringValue(g.ID)
	state.PageID = types.StringValue(g.PageID)
	state.Name = types.StringValue(g.Name)
	state.Description = types.StringValue(g.Description)
	state.Position = types.Int64Value(int64(g.Position))
	components, _ := types.ListValueFrom(ctx, types.StringType, g.Components)
	state.Components = components
}
