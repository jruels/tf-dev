// provider: package for the provider
//
//	file(s)
package provider

import (
	// context: control lifetime of request.
	"context"

	// These are imports for Terraform's
	//	     plugin framwork.
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure Provider implementation satisfies required interfaces.
//
//	This is a compile-time check.
var _ provider.Provider = &GreetingProvider{}

// Define GreetingProvider as the struct, with no fields.
type GreetingProvider struct{}

// Factory method that returns a specific
//
//	instance of a provider.
func New() provider.Provider {
	return &GreetingProvider{}
}

// Metadata data that is provided to Terraform
func (p *GreetingProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "greeting"
}

// Schema defines data passed into the provider from main.tf
func (p *GreetingProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Greeting provider that returns 'Hello, world!'",
	}
}

// Configure does nothing as there is no configuration needed
func (p *GreetingProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
}

// DataSources method is required but not needed for this provider
func (p *GreetingProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return nil
}

// Tells Terraform what resource(s) the provider supports.
//
//	In this example, the factory function for the
//	greeting resource is returned.
func (p *GreetingProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewGreetingResource,
	}
}
