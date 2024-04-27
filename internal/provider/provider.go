// Copyright (c) Harel Safra

package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure FileDataProvider satisfies various provider interfaces.
var _ provider.Provider = &FileDataProvider{}
var _ provider.ProviderWithFunctions = &FileDataProvider{}

// FileDataProvider defines the provider implementation.
type FileDataProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// FileDataProviderModel describes the provider data model.
type FileDataProviderModel struct {
	Base_path types.String `tfsdk:"base_path"`
}

func (p *FileDataProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "filedata"
	resp.Version = p.version
}

func (p *FileDataProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"base_path": schema.StringAttribute{
				MarkdownDescription: "Base path of the files to manage",
				Required:            true,
			},
		},
	}
}

func (p *FileDataProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data FileDataProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.DataSourceData = data.Base_path.ValueString()
	resp.ResourceData = data.Base_path.ValueString()
}

func (p *FileDataProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewFile,
	}
}

func (p *FileDataProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *FileDataProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &FileDataProvider{
			version: version,
		}
	}
}
