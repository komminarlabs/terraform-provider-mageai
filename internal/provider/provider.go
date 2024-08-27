package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/komminarlabs/terraform-provider-mageai/internal/sdk/mageai"
)

// Ensure the implementation satisfies the expected interfaces.
var _ provider.Provider = &MageAIProvider{}

// MageAIProvider defines the provider implementation.
type MageAIProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// MageAIProviderModel maps provider schema data to a Go type.
type MageAIProviderModel struct {
	ApiKey types.String `tfsdk:"api_key"`
	Host   types.String `tfsdk:"host"`
}

type providerData struct {
	client mageai.Client
}

// Metadata returns the provider type name.
func (p *MageAIProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "mageai"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *MageAIProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Mage AI provider to deploy and manage resources supported by Mage AI.",

		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Description: "The API key to authenticate calls",
				Optional:    true,
				Sensitive:   true,
			},
			"host": schema.StringAttribute{
				Description: "The host of the Mage AI server",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

// Configure prepares a Mage AI API client for data sources and resources.
func (p *MageAIProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config MageAIProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.
	if config.ApiKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown Mage AI API Key",
			"The provider cannot create the Mage AI client as there is an unknown configuration value for the API Key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the MAGEAI_API_KEY environment variable.",
		)
	}

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown Mage AI Host",
			"The provider cannot create the Mage AI client as there is an unknown configuration value for the Host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the MAGEAI_HOST environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	apiKey := os.Getenv("MAGEAI_API_KEY")
	host := os.Getenv("MAGEAI_HOST")

	if !config.ApiKey.IsNull() {
		apiKey = config.ApiKey.ValueString()
	}

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("apiKey"),
			"Missing Mage AI API Key",
			"The provider cannot create the Mage AI client as there is an unknown configuration value for the API Key. "+
				"Set the API Key value in the configuration or use the MAGEAI_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing Mage AI Host",
			"The provider cannot create the Mage AI client as there is a missing or empty value for the Host. "+
				"Set the Cluster ID value in the configuration or use the MAGEAI_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// If any of the expected configurations are in wrong format, return
	// errors with provider-specific guidance.

	ctx = tflog.SetField(ctx, "MAGEAI_API_KEY", apiKey)
	ctx = tflog.SetField(ctx, "MAGEAI_HOST", host)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "MAGEAI_API_KEY")

	tflog.Debug(ctx, "Creating Mage AI client")

	// Create a new Mage AI client using the configuration values
	client, err := mageai.New(
		&mageai.ClientConfig{
			Host:   host,
			ApiKey: apiKey,
		},
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Mage AI Client",
			"An unexpected error occurred when creating the Mage AI client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Mage AI Client Error: "+err.Error(),
		)
		return
	}

	// Make the Mage AI client available during DataSource and Resource
	// type Configure methods.
	providerData := &providerData{
		client: client,
	}
	resp.DataSourceData = *providerData
	resp.ResourceData = *providerData
	tflog.Info(ctx, "Configured Mage AI client", map[string]any{"success": true})
}

// Resources defines the resources implemented in the provider.
func (p *MageAIProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

// DataSources defines the data sources implemented in the provider.
func (p *MageAIProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewPipelineDataSource,
		NewPipelinesDataSource,
	}
}

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &MageAIProvider{
			version: version,
		}
	}
}
