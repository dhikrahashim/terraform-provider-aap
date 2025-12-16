package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/dhikrahashim/terraform-provider-aap/internal/client"
)

// Ensure AapProvider satisfies various interfaces.
var _ provider.Provider = &AapProvider{}

type AapProvider struct {
	version string
}

type AapProviderModel struct {
	Host     types.String `tfsdk:"host"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	Token    types.String `tfsdk:"token"`
	Insecure types.Bool   `tfsdk:"insecure"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &AapProvider{
			version: version,
		}
	}
}

func (p *AapProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "aap"
	resp.Version = p.version
}

func (p *AapProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The AAP provider allows you to configure Ansible Automation Platform resources.",
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Description: "The URI of the Ansible Automation Platform Controller (e.g. https://aap.example.com)",
				Optional:    true,
			},
			"username": schema.StringAttribute{
				Description: "The username for AAP authentication.",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "The password for AAP authentication.",
				Optional:    true,
				Sensitive:   true,
			},
			"token": schema.StringAttribute{
				Description: "The OAuth2 token for AAP authentication.",
				Optional:    true,
				Sensitive:   true,
			},
			"insecure": schema.BoolAttribute{
				Description: "Whether to skip TLS verification.",
				Optional:    true,
			},
		},
	}
}

func (p *AapProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data AapProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	host := os.Getenv("AAP_HOST")
	username := os.Getenv("AAP_USERNAME")
	password := os.Getenv("AAP_PASSWORD")
	token := os.Getenv("AAP_TOKEN")
	insecure := false

	if !data.Host.IsNull() {
		host = data.Host.ValueString()
	}
	if !data.Username.IsNull() {
		username = data.Username.ValueString()
	}
	if !data.Password.IsNull() {
		password = data.Password.ValueString()
	}
	if !data.Token.IsNull() {
		token = data.Token.ValueString()
	}
	if !data.Insecure.IsNull() {
		insecure = data.Insecure.ValueBool()
	}

	if host == "" {
		resp.Diagnostics.AddError("Missing host", "AAP host must be configured via provider or AAP_HOST env var")
		return
	}

	// Basic client setup (placeholder)
	c := client.NewClient(host, username, password, token, insecure)
	
	resp.DataSourceData = c
	resp.ResourceData = c
}

func (p *AapProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewOrganizationResource,
		NewInventoryResource,
		NewJobTemplateResource,
		NewProjectResource,
		NewCredentialMachineResource,
		NewCredentialScmResource,
		NewCredentialTypeResource,
		NewInventorySourceResource,
	}
}

func (p *AapProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
