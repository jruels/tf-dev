package provider

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// ✅ Rename DbClient to BankingDBClient
// type BankingDBClient struct {
// 	DB *sql.DB
// }

// ✅ Implement function to create a new database client
func NewBankingDBClient(host string, port int64, user, password, dbname string) (*BankingDBClient, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return &BankingDBClient{DB: db}, nil
}

// ✅ Define bankingProvider struct
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
	var config struct {
		DBHost     string `tfsdk:"db_host"`
		DBPort     int64  `tfsdk:"db_port"`
		DBUser     string `tfsdk:"db_user"`
		DBPassword string `tfsdk:"db_password"`
		DBName     string `tfsdk:"db_name"`
	}

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := NewBankingDBClient(config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName)
	if err != nil {
		resp.Diagnostics.AddError("Database Connection Error", fmt.Sprintf("Could not connect to database: %s", err.Error()))
		return
	}
	p.dbClient = client
	resp.ResourceData = client
}

// ✅ Implement DataSources function (required by Terraform Plugin Framework)
func (p *bankingProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return nil // No data sources implemented yet
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
