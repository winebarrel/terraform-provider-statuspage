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
	_ resource.Resource                = &metricResource{}
	_ resource.ResourceWithConfigure   = &metricResource{}
	_ resource.ResourceWithImportState = &metricResource{}
)

type metricResource struct {
	client *apiclient.Client
}

type metricResourceModel struct {
	ID                 types.String  `tfsdk:"id"`
	PageID             types.String  `tfsdk:"page_id"`
	MetricsProviderID  types.String  `tfsdk:"metrics_provider_id"`
	MetricIdentifier   types.String  `tfsdk:"metric_identifier"`
	Name               types.String  `tfsdk:"name"`
	Display            types.Bool    `tfsdk:"display"`
	TooltipDescription types.String  `tfsdk:"tooltip_description"`
	YAxisMin           types.Float64 `tfsdk:"y_axis_min"`
	YAxisMax           types.Float64 `tfsdk:"y_axis_max"`
	YAxisHidden        types.Bool    `tfsdk:"y_axis_hidden"`
	Suffix             types.String  `tfsdk:"suffix"`
	DecimalPlaces      types.Int64   `tfsdk:"decimal_places"`
}

func NewMetricResource() resource.Resource {
	return &metricResource{}
}

func (r *metricResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_metric"
}

func (r *metricResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Statuspage metric.",
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
			"metrics_provider_id": schema.StringAttribute{
				Required:    true,
				Description: "The metrics provider this metric belongs to.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"metric_identifier": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"display": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"tooltip_description": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"y_axis_min": schema.Float64Attribute{
				Optional: true,
				Computed: true,
			},
			"y_axis_max": schema.Float64Attribute{
				Optional: true,
				Computed: true,
			},
			"y_axis_hidden": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"suffix": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"decimal_places": schema.Int64Attribute{
				Optional: true,
				Computed: true,
			},
		},
	}
}

func (r *metricResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *metricResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan metricResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := apiclient.MetricBody{
		Name:               plan.Name.ValueString(),
		MetricIdentifier:   plan.MetricIdentifier.ValueString(),
		TooltipDescription: plan.TooltipDescription.ValueString(),
		Suffix:             plan.Suffix.ValueString(),
	}
	if !plan.Display.IsNull() && !plan.Display.IsUnknown() {
		v := plan.Display.ValueBool()
		body.Display = &v
	}
	if !plan.YAxisMin.IsNull() && !plan.YAxisMin.IsUnknown() {
		v := plan.YAxisMin.ValueFloat64()
		body.YAxisMin = &v
	}
	if !plan.YAxisMax.IsNull() && !plan.YAxisMax.IsUnknown() {
		v := plan.YAxisMax.ValueFloat64()
		body.YAxisMax = &v
	}
	if !plan.YAxisHidden.IsNull() && !plan.YAxisHidden.IsUnknown() {
		v := plan.YAxisHidden.ValueBool()
		body.YAxisHidden = &v
	}
	if !plan.DecimalPlaces.IsNull() && !plan.DecimalPlaces.IsUnknown() {
		v := int(plan.DecimalPlaces.ValueInt64())
		body.DecimalPlaces = &v
	}

	metric, err := r.client.CreateMetric(ctx, plan.PageID.ValueString(), plan.MetricsProviderID.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError("Error creating metric", err.Error())
		return
	}

	r.mapToState(&plan, metric)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *metricResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state metricResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	metric, err := r.client.GetMetric(ctx, state.PageID.ValueString(), state.ID.ValueString())
	if err != nil {
		var apiErr *apiclient.APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading metric", err.Error())
		return
	}

	r.mapToState(&state, metric)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *metricResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan metricResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := apiclient.MetricBody{
		Name:               plan.Name.ValueString(),
		MetricIdentifier:   plan.MetricIdentifier.ValueString(),
		TooltipDescription: plan.TooltipDescription.ValueString(),
		Suffix:             plan.Suffix.ValueString(),
	}
	if !plan.Display.IsNull() && !plan.Display.IsUnknown() {
		v := plan.Display.ValueBool()
		body.Display = &v
	}
	if !plan.YAxisMin.IsNull() && !plan.YAxisMin.IsUnknown() {
		v := plan.YAxisMin.ValueFloat64()
		body.YAxisMin = &v
	}
	if !plan.YAxisMax.IsNull() && !plan.YAxisMax.IsUnknown() {
		v := plan.YAxisMax.ValueFloat64()
		body.YAxisMax = &v
	}
	if !plan.YAxisHidden.IsNull() && !plan.YAxisHidden.IsUnknown() {
		v := plan.YAxisHidden.ValueBool()
		body.YAxisHidden = &v
	}
	if !plan.DecimalPlaces.IsNull() && !plan.DecimalPlaces.IsUnknown() {
		v := int(plan.DecimalPlaces.ValueInt64())
		body.DecimalPlaces = &v
	}

	metric, err := r.client.UpdateMetric(ctx, plan.PageID.ValueString(), plan.ID.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError("Error updating metric", err.Error())
		return
	}

	r.mapToState(&plan, metric)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *metricResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state metricResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteMetric(ctx, state.PageID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting metric", err.Error())
		return
	}
}

func (r *metricResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := splitImportID(req.ID, 2)
	if idParts == nil {
		resp.Diagnostics.AddError("Invalid import ID", "Expected format: page_id/metric_id")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("page_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)
}

func (r *metricResource) mapToState(state *metricResourceModel, m *apiclient.Metric) {
	state.ID = types.StringValue(m.ID)
	state.PageID = types.StringValue(m.PageID)
	state.MetricsProviderID = types.StringValue(m.MetricsProviderID)
	state.MetricIdentifier = types.StringValue(m.MetricIdentifier)
	state.Name = types.StringValue(m.Name)
	state.Display = types.BoolValue(m.Display)
	state.TooltipDescription = types.StringValue(m.TooltipDescription)
	state.YAxisMin = types.Float64Value(m.YAxisMin)
	state.YAxisMax = types.Float64Value(m.YAxisMax)
	state.YAxisHidden = types.BoolValue(m.YAxisHidden)
	state.Suffix = types.StringValue(m.Suffix)
	state.DecimalPlaces = types.Int64Value(int64(m.DecimalPlaces))
}
