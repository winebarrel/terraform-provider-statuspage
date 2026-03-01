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
	_ resource.Resource                = &subscriberResource{}
	_ resource.ResourceWithConfigure   = &subscriberResource{}
	_ resource.ResourceWithImportState = &subscriberResource{}
)

type subscriberResource struct {
	client *apiclient.Client
}

type subscriberResourceModel struct {
	ID                           types.String `tfsdk:"id"`
	PageID                       types.String `tfsdk:"page_id"`
	Email                        types.String `tfsdk:"email"`
	PhoneNumber                  types.String `tfsdk:"phone_number"`
	PhoneCountry                 types.String `tfsdk:"phone_country"`
	Endpoint                     types.String `tfsdk:"endpoint"`
	Mode                         types.String `tfsdk:"mode"`
	ComponentIDs                 types.List   `tfsdk:"component_ids"`
	SkipConfirmationNotification types.Bool   `tfsdk:"skip_confirmation_notification"`
}

func NewSubscriberResource() resource.Resource {
	return &subscriberResource{}
}

func (r *subscriberResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_subscriber"
}

func (r *subscriberResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Statuspage subscriber.",
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
			"email": schema.StringAttribute{
				Optional: true,
			},
			"phone_number": schema.StringAttribute{
				Optional: true,
			},
			"phone_country": schema.StringAttribute{
				Optional: true,
			},
			"endpoint": schema.StringAttribute{
				Optional:    true,
				Description: "Webhook endpoint URL.",
			},
			"mode": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("email", "sms", "slack", "webhook", "microsoft_teams"),
				},
			},
			"component_ids": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"skip_confirmation_notification": schema.BoolAttribute{
				Optional: true,
			},
		},
	}
}

func (r *subscriberResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *subscriberResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan subscriberResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := apiclient.SubscriberBody{
		Email:        plan.Email.ValueString(),
		PhoneNumber:  plan.PhoneNumber.ValueString(),
		PhoneCountry: plan.PhoneCountry.ValueString(),
		Endpoint:     plan.Endpoint.ValueString(),
	}
	if !plan.ComponentIDs.IsNull() && !plan.ComponentIDs.IsUnknown() {
		var ids []string
		plan.ComponentIDs.ElementsAs(ctx, &ids, false)
		body.ComponentIDs = ids
	}
	if !plan.SkipConfirmationNotification.IsNull() && !plan.SkipConfirmationNotification.IsUnknown() {
		v := plan.SkipConfirmationNotification.ValueBool()
		body.SkipConfirmationNotification = &v
	}

	subscriber, err := r.client.CreateSubscriber(ctx, plan.PageID.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError("Error creating subscriber", err.Error())
		return
	}

	r.mapToState(ctx, &plan, subscriber)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *subscriberResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state subscriberResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	subscriber, err := r.client.GetSubscriber(ctx, state.PageID.ValueString(), state.ID.ValueString())
	if err != nil {
		var apiErr *apiclient.APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading subscriber", err.Error())
		return
	}

	r.mapToState(ctx, &state, subscriber)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *subscriberResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan subscriberResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := apiclient.SubscriberBody{
		Email:        plan.Email.ValueString(),
		PhoneNumber:  plan.PhoneNumber.ValueString(),
		PhoneCountry: plan.PhoneCountry.ValueString(),
		Endpoint:     plan.Endpoint.ValueString(),
	}
	if !plan.ComponentIDs.IsNull() && !plan.ComponentIDs.IsUnknown() {
		var ids []string
		plan.ComponentIDs.ElementsAs(ctx, &ids, false)
		body.ComponentIDs = ids
	}

	subscriber, err := r.client.UpdateSubscriber(ctx, plan.PageID.ValueString(), plan.ID.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError("Error updating subscriber", err.Error())
		return
	}

	r.mapToState(ctx, &plan, subscriber)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *subscriberResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state subscriberResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteSubscriber(ctx, state.PageID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting subscriber", err.Error())
		return
	}
}

func (r *subscriberResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := splitImportID(req.ID, 2)
	if idParts == nil {
		resp.Diagnostics.AddError("Invalid import ID", "Expected format: page_id/subscriber_id")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("page_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)
}

func (r *subscriberResource) mapToState(ctx context.Context, state *subscriberResourceModel, s *apiclient.Subscriber) {
	state.ID = types.StringValue(s.ID)
	state.PageID = types.StringValue(s.PageID)
	state.Email = types.StringValue(s.Email)
	state.PhoneNumber = types.StringValue(s.PhoneNumber)
	state.PhoneCountry = types.StringValue(s.PhoneCountry)
	state.Endpoint = types.StringValue(s.Endpoint)
	state.Mode = types.StringValue(s.Mode)
	if len(s.ComponentIDs) > 0 {
		componentIDs, _ := types.ListValueFrom(ctx, types.StringType, s.ComponentIDs)
		state.ComponentIDs = componentIDs
	}
}
