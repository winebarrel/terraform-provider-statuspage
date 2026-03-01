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
	_ resource.Resource                = &metricsProviderResource{}
	_ resource.ResourceWithConfigure   = &metricsProviderResource{}
	_ resource.ResourceWithImportState = &metricsProviderResource{}
)

type metricsProviderResource struct {
	client *apiclient.Client
}

type metricsProviderResourceModel struct {
	ID             types.String `tfsdk:"id"`
	PageID         types.String `tfsdk:"page_id"`
	Type           types.String `tfsdk:"type"`
	Email          types.String `tfsdk:"email"`
	MetricBaseURI  types.String `tfsdk:"metric_base_uri"`
	APIKey         types.String `tfsdk:"api_key"`
	APIToken       types.String `tfsdk:"api_token"`
	ApplicationKey types.String `tfsdk:"application_key"`
}

func NewMetricsProviderResource() resource.Resource {
	return &metricsProviderResource{}
}

func (r *metricsProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_metrics_provider"
}

func (r *metricsProviderResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Statuspage metrics provider.",
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
			"type": schema.StringAttribute{
				Required:    true,
				Description: "One of: Pingdom, NewRelic, Datadog, Self, Librato, CalabashCustom.",
				Validators: []validator.String{
					stringvalidator.OneOf("Pingdom", "NewRelic", "Datadog", "Self", "Librato", "CalabashCustom"),
				},
			},
			"email": schema.StringAttribute{
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},
			"metric_base_uri": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"api_key": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"api_token": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"application_key": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func (r *metricsProviderResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *metricsProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan metricsProviderResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := apiclient.MetricsProviderBody{
		Type:           plan.Type.ValueString(),
		Email:          plan.Email.ValueString(),
		MetricBaseURI:  plan.MetricBaseURI.ValueString(),
		APIKey:         plan.APIKey.ValueString(),
		APIToken:       plan.APIToken.ValueString(),
		ApplicationKey: plan.ApplicationKey.ValueString(),
	}

	provider, err := r.client.CreateMetricsProvider(ctx, plan.PageID.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError("Error creating metrics provider", err.Error())
		return
	}

	r.mapToState(&plan, provider)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *metricsProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state metricsProviderResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	provider, err := r.client.GetMetricsProvider(ctx, state.PageID.ValueString(), state.ID.ValueString())
	if err != nil {
		var apiErr *apiclient.APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading metrics provider", err.Error())
		return
	}

	r.mapToState(&state, provider)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *metricsProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan metricsProviderResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := apiclient.MetricsProviderBody{
		Type:           plan.Type.ValueString(),
		Email:          plan.Email.ValueString(),
		MetricBaseURI:  plan.MetricBaseURI.ValueString(),
		APIKey:         plan.APIKey.ValueString(),
		APIToken:       plan.APIToken.ValueString(),
		ApplicationKey: plan.ApplicationKey.ValueString(),
	}

	provider, err := r.client.UpdateMetricsProvider(ctx, plan.PageID.ValueString(), plan.ID.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError("Error updating metrics provider", err.Error())
		return
	}

	r.mapToState(&plan, provider)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *metricsProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state metricsProviderResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteMetricsProvider(ctx, state.PageID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting metrics provider", err.Error())
		return
	}
}

func (r *metricsProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := splitImportID(req.ID, 2)
	if idParts == nil {
		resp.Diagnostics.AddError("Invalid import ID", "Expected format: page_id/metrics_provider_id")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("page_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)
}

func (r *metricsProviderResource) mapToState(state *metricsProviderResourceModel, p *apiclient.MetricsProvider) {
	state.ID = types.StringValue(p.ID)
	state.PageID = types.StringValue(p.PageID)
	state.Type = types.StringValue(p.Type)
	state.Email = types.StringValue(p.Email)
	state.MetricBaseURI = types.StringValue(p.MetricBaseURI)
}
