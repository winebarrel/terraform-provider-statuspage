package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/winebarrel/terraform-provider-statuspage/internal/apiclient"
)

var _ provider.Provider = &statuspageProvider{}

type statuspageProvider struct {
	version     string
	clientOpts  []apiclient.ClientOption
}

type statuspageProviderModel struct {
	APIKey types.String `tfsdk:"api_key"`
}

type providerData struct {
	Client *apiclient.Client
}

func New(version string, clientOpts ...apiclient.ClientOption) func() provider.Provider {
	return func() provider.Provider {
		return &statuspageProvider{version: version, clientOpts: clientOpts}
	}
}

func (p *statuspageProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "statuspage"
	resp.Version = p.version
}

func (p *statuspageProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Terraform provider for managing Statuspage.io resources.",
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Description: "The Statuspage.io API key. Can also be set via STATUSPAGE_API_KEY environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *statuspageProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config statuspageProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiKey := config.APIKey.ValueString()
	if apiKey == "" {
		apiKey = os.Getenv("STATUSPAGE_API_KEY")
	}
	if apiKey == "" {
		resp.Diagnostics.AddError(
			"Missing API Key",
			"api_key must be set in provider configuration or via STATUSPAGE_API_KEY environment variable.",
		)
		return
	}

	client := apiclient.NewClient(apiKey, p.clientOpts...)
	pd := &providerData{
		Client: client,
	}

	resp.DataSourceData = pd
	resp.ResourceData = pd
}

func (p *statuspageProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewComponentResource,
		NewComponentGroupResource,
		NewIncidentResource,
		NewIncidentTemplateResource,
		NewIncidentPostmortemResource,
		NewMetricResource,
		NewMetricsProviderResource,
		NewPageResource,
		NewPageAccessUserResource,
		NewPageAccessGroupResource,
		NewStatusEmbedConfigResource,
		NewSubscriberResource,
		NewUserResource,
		NewUserPermissionsResource,
	}
}

func (p *statuspageProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
