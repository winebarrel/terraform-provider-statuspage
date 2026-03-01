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
	_ resource.Resource                = &statusEmbedConfigResource{}
	_ resource.ResourceWithConfigure   = &statusEmbedConfigResource{}
	_ resource.ResourceWithImportState = &statusEmbedConfigResource{}
)

type statusEmbedConfigResource struct {
	client *apiclient.Client
}

type statusEmbedConfigResourceModel struct {
	PageID                     types.String `tfsdk:"page_id"`
	Position                   types.String `tfsdk:"position"`
	IncidentBackgroundColor    types.String `tfsdk:"incident_background_color"`
	IncidentTextColor          types.String `tfsdk:"incident_text_color"`
	MaintenanceBackgroundColor types.String `tfsdk:"maintenance_background_color"`
	MaintenanceTextColor       types.String `tfsdk:"maintenance_text_color"`
}

func NewStatusEmbedConfigResource() resource.Resource {
	return &statusEmbedConfigResource{}
}

func (r *statusEmbedConfigResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_status_embed_config"
}

func (r *statusEmbedConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Statuspage status embed configuration. Cannot be created or deleted; use terraform import.",
		Attributes: map[string]schema.Attribute{
			"page_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"position": schema.StringAttribute{
				Optional: true, Computed: true,
			},
			"incident_background_color": schema.StringAttribute{
				Optional: true, Computed: true,
			},
			"incident_text_color": schema.StringAttribute{
				Optional: true, Computed: true,
			},
			"maintenance_background_color": schema.StringAttribute{
				Optional: true, Computed: true,
			},
			"maintenance_text_color": schema.StringAttribute{
				Optional: true, Computed: true,
			},
		},
	}
}

func (r *statusEmbedConfigResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *statusEmbedConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan statusEmbedConfigResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config, err := r.client.GetStatusEmbedConfig(ctx, plan.PageID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading status embed config", err.Error())
		return
	}

	r.mapToState(&plan, config)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *statusEmbedConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state statusEmbedConfigResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config, err := r.client.GetStatusEmbedConfig(ctx, state.PageID.ValueString())
	if err != nil {
		var apiErr *apiclient.APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading status embed config", err.Error())
		return
	}

	r.mapToState(&state, config)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *statusEmbedConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan statusEmbedConfigResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := apiclient.StatusEmbedConfigBody{
		Position:                   plan.Position.ValueString(),
		IncidentBackgroundColor:    plan.IncidentBackgroundColor.ValueString(),
		IncidentTextColor:          plan.IncidentTextColor.ValueString(),
		MaintenanceBackgroundColor: plan.MaintenanceBackgroundColor.ValueString(),
		MaintenanceTextColor:       plan.MaintenanceTextColor.ValueString(),
	}

	config, err := r.client.UpdateStatusEmbedConfig(ctx, plan.PageID.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError("Error updating status embed config", err.Error())
		return
	}

	r.mapToState(&plan, config)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *statusEmbedConfigResource) Delete(_ context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddWarning(
		"Status embed config not deleted",
		"Status embed configs cannot be deleted via the API. Removed from Terraform state only.",
	)
}

func (r *statusEmbedConfigResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("page_id"), req, resp)
}

func (r *statusEmbedConfigResource) mapToState(state *statusEmbedConfigResourceModel, c *apiclient.StatusEmbedConfig) {
	state.PageID = types.StringValue(c.PageID)
	state.Position = types.StringValue(c.Position)
	state.IncidentBackgroundColor = types.StringValue(c.IncidentBackgroundColor)
	state.IncidentTextColor = types.StringValue(c.IncidentTextColor)
	state.MaintenanceBackgroundColor = types.StringValue(c.MaintenanceBackgroundColor)
	state.MaintenanceTextColor = types.StringValue(c.MaintenanceTextColor)
}
