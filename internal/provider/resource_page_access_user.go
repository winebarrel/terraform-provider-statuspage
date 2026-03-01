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
	_ resource.Resource                = &pageAccessUserResource{}
	_ resource.ResourceWithConfigure   = &pageAccessUserResource{}
	_ resource.ResourceWithImportState = &pageAccessUserResource{}
)

type pageAccessUserResource struct {
	client *apiclient.Client
}

type pageAccessUserResourceModel struct {
	ID            types.String `tfsdk:"id"`
	PageID        types.String `tfsdk:"page_id"`
	ExternalLogin types.String `tfsdk:"external_login"`
	ExternalEmail types.String `tfsdk:"external_email"`
	ComponentIDs  types.List   `tfsdk:"component_ids"`
	MetricIDs     types.List   `tfsdk:"metric_ids"`
}

func NewPageAccessUserResource() resource.Resource {
	return &pageAccessUserResource{}
}

func (r *pageAccessUserResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_page_access_user"
}

func (r *pageAccessUserResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Statuspage page access user.",
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
			"external_login": schema.StringAttribute{
				Required: true,
			},
			"external_email": schema.StringAttribute{
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
		},
	}
}

func (r *pageAccessUserResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *pageAccessUserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan pageAccessUserResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := apiclient.PageAccessUserBody{
		ExternalLogin: plan.ExternalLogin.ValueString(),
		ExternalEmail: plan.ExternalEmail.ValueString(),
	}
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

	user, err := r.client.CreatePageAccessUser(ctx, plan.PageID.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError("Error creating page access user", err.Error())
		return
	}

	r.mapToState(ctx, &plan, user)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *pageAccessUserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state pageAccessUserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user, err := r.client.GetPageAccessUser(ctx, state.PageID.ValueString(), state.ID.ValueString())
	if err != nil {
		var apiErr *apiclient.APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading page access user", err.Error())
		return
	}

	r.mapToState(ctx, &state, user)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *pageAccessUserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan pageAccessUserResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := apiclient.PageAccessUserBody{
		ExternalLogin: plan.ExternalLogin.ValueString(),
		ExternalEmail: plan.ExternalEmail.ValueString(),
	}
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

	user, err := r.client.UpdatePageAccessUser(ctx, plan.PageID.ValueString(), plan.ID.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError("Error updating page access user", err.Error())
		return
	}

	r.mapToState(ctx, &plan, user)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *pageAccessUserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state pageAccessUserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeletePageAccessUser(ctx, state.PageID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting page access user", err.Error())
		return
	}
}

func (r *pageAccessUserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := splitImportID(req.ID, 2)
	if idParts == nil {
		resp.Diagnostics.AddError("Invalid import ID", "Expected format: page_id/page_access_user_id")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("page_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)
}

func (r *pageAccessUserResource) mapToState(ctx context.Context, state *pageAccessUserResourceModel, u *apiclient.PageAccessUser) {
	state.ID = types.StringValue(u.ID)
	state.PageID = types.StringValue(u.PageID)
	state.ExternalLogin = types.StringValue(u.ExternalLogin)
	state.ExternalEmail = types.StringValue(u.ExternalEmail)
	if len(u.ComponentIDs) > 0 {
		componentIDs, _ := types.ListValueFrom(ctx, types.StringType, u.ComponentIDs)
		state.ComponentIDs = componentIDs
	}
	if len(u.MetricIDs) > 0 {
		metricIDs, _ := types.ListValueFrom(ctx, types.StringType, u.MetricIDs)
		state.MetricIDs = metricIDs
	}
}
