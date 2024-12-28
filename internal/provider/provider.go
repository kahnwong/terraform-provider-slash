package provider

import (
	"context"
	"github.com/kahnwong/terraform-provider-slash/slash"
	"os"
	//"github.com/hashicorp-demoapp/hashicups-client-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ provider.Provider              = &slashProvider{}
	_ provider.ProviderWithFunctions = &slashProvider{}
)

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &slashProvider{
			version: version,
		}
	}
}

type slashProviderModel struct {
	Host        types.String `tfsdk:"host"`
	AccessToken types.String `tfsdk:"access_token"`
}

type slashProvider struct {
	version string
}

func (p *slashProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "slash"
	resp.Version = p.version
}

func (p *slashProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with Slash.",
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Description: "URI for Slash API. May also be provided via SLASH_HOST environment variable.",
				Optional:    true,
			},
			"access_token": schema.StringAttribute{
				Description: "Password for Slash API. May also be provided via SLASH_PASSWORD environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *slashProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Slash client")

	var config slashProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown Slash API Host",
			"The provider cannot create the Slash API client as there is an unknown configuration value for the Slash API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SLASH_HOST environment variable.",
		)
	}

	if config.AccessToken.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("access_token"),
			"Unknown Slash API access_token",
			"The provider cannot create the Slash API client as there is an unknown configuration value for the Slash API username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SLASH_ACCESS_TOKEN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	host := os.Getenv("SLASH_HOST")
	accessToken := os.Getenv("SLASH_ACCESS_TOKEN")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.AccessToken.IsNull() {
		accessToken = config.AccessToken.ValueString()
	}

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing Slash API Host",
			"The provider cannot create the Slash API client as there is a missing or empty value for the Slash API host. "+
				"Set the host value in the configuration or use the SLASH_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if accessToken == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing Slash API AccessToken",
			"The provider cannot create the Slash API client as there is a missing or empty value for the Slash API username. "+
				"Set the username value in the configuration or use the SLASH_ACCESS_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "slash_host", host)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "slash_access_token")

	tflog.Debug(ctx, "Creating Slash client")

	client, err := slash.NewClient(&host, &accessToken) // hashicups.NewClient(&host, &accessToken)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Slash API Client",
			"An unexpected error occurred when creating the Slash API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Slash Client Error: "+err.Error(),
		)
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Slash client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *slashProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		//NewCoffeesDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *slashProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		//NewOrderResource,
	}
}

func (p *slashProvider) Functions(_ context.Context) []func() function.Function {
	return []func() function.Function{
		//NewComputeTaxFunction,
	}
}
