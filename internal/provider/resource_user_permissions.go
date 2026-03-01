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
	_ resource.Resource                = &userPermissionsResource{}
	_ resource.ResourceWithConfigure   = &userPermissionsResource{}
	_ resource.ResourceWithImportState = &userPermissionsResource{}
)

type userPermissionsResource struct {
	client *apiclient.Client
}

type userPermissionsResourceModel struct {
	OrganizationID types.String `tfsdk:"organization_id"`
	UserID         types.String `tfsdk:"user_id"`
	Pages          types.Map    `tfsdk:"pages"`
}

func NewUserPermissionsResource() resource.Resource {
	return &userPermissionsResource{}
}

func (r *userPermissionsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_permissions"
}

func (r *userPermissionsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages Statuspage user permissions.",
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"user_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"pages": schema.MapAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "Map of page_id to permission level.",
			},
		},
	}
}

func (r *userPermissionsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *userPermissionsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan userPermissionsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := r.buildBody(ctx, &plan)
	perms, err := r.client.UpdatePermissions(ctx, plan.OrganizationID.ValueString(), plan.UserID.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError("Error setting user permissions", err.Error())
		return
	}

	r.mapToState(ctx, &plan, perms)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *userPermissionsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state userPermissionsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	perms, err := r.client.GetPermissions(ctx, state.OrganizationID.ValueString(), state.UserID.ValueString())
	if err != nil {
		var apiErr *apiclient.APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading user permissions", err.Error())
		return
	}

	r.mapToState(ctx, &state, perms)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *userPermissionsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan userPermissionsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := r.buildBody(ctx, &plan)
	perms, err := r.client.UpdatePermissions(ctx, plan.OrganizationID.ValueString(), plan.UserID.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError("Error updating user permissions", err.Error())
		return
	}

	r.mapToState(ctx, &plan, perms)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *userPermissionsResource) Delete(_ context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddWarning(
		"Permissions not deleted",
		"User permissions cannot be deleted via the API. Removed from Terraform state only.",
	)
}

func (r *userPermissionsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := splitImportID(req.ID, 2)
	if idParts == nil {
		resp.Diagnostics.AddError("Invalid import ID", "Expected format: organization_id/user_id")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("organization_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("user_id"), idParts[1])...)
}

func (r *userPermissionsResource) buildBody(ctx context.Context, plan *userPermissionsResourceModel) apiclient.PermissionsRequest {
	body := apiclient.PermissionsRequest{
		Pages: make(map[string][]string),
	}
	if !plan.Pages.IsNull() && !plan.Pages.IsUnknown() {
		m := make(map[string]string)
		plan.Pages.ElementsAs(ctx, &m, false)
		for pageID, perm := range m {
			body.Pages[pageID] = []string{perm}
		}
	}
	return body
}

func (r *userPermissionsResource) mapToState(ctx context.Context, state *userPermissionsResourceModel, p *apiclient.Permissions) {
	state.UserID = types.StringValue(p.UserID)
	if len(p.Pages) > 0 {
		pages, _ := types.MapValueFrom(ctx, types.StringType, p.Pages)
		state.Pages = pages
	}
}
