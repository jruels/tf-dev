package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type bankingProvider struct {
	dbClient *BankingDBClient
}

// ✅ Implement Metadata function
func (p *bankingProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "banking"
}

// ✅ Implement Schema function
func (p *bankingProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"db_host":     schema.StringAttribute{Required: true},
			"db_port":     schema.Int64Attribute{Required: true},
			"db_user":     schema.StringAttribute{Required: true},
			"db_password": schema.StringAttribute{Required: true, Sensitive: true},
			"db_name":     schema.StringAttribute{Required: true},
		},
	}
}

// ✅ Implement Configure function
func (p *bankingProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	client, err := NewBankingDBClient("localhost", 5432, "postgres", "Post1260", "bankingdb")
	if err != nil {
		resp.Diagnostics.AddError("Database Connection Error", fmt.Sprintf("Could not connect to database: %s", err.Error()))
		return
	}
	p.dbClient = client
	resp.ResourceData = client
}

// ✅ Implement DataSources function (required by Terraform Plugin Framework)
func (p *bankingProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	// Return an empty list if no data sources are implemented
	return nil
}

// ✅ Register Resources (Terraform Resources like customer accounts)
func (p *bankingProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource { return &customerAccountResource{client: p.dbClient} },
	}
}

// ✅ Fix New() function to return the provider correctly
func New() provider.Provider {
	return &bankingProvider{}
}
