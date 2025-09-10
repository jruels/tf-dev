package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure implementation satisfies Terraform framework interfaces
var _ resource.Resource = &customerAccountResource{}

type customerAccountResource struct {
	client *BankingDBClient
}

// ✅ Define customerAccountPlan struct to match Terraform state
type customerAccountPlan struct {
	ID          types.String  `tfsdk:"id"`
	FirstName   types.String  `tfsdk:"first_name"`
	LastName    types.String  `tfsdk:"last_name"`
	Email       types.String  `tfsdk:"email"`
	AccountType types.String  `tfsdk:"account_type"`
	Balance     types.Float64 `tfsdk:"balance"`
}

// **Metadata: Registers the resource with Terraform**
func (r *customerAccountResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "banking_customer_account"
}

// **Schema: Defines resource attributes**
func (r *customerAccountResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a customer account in the banking system.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier for the customer account.",
			},
			"first_name": schema.StringAttribute{
				Required:    true,
				Description: "Customer's first name.",
			},
			"last_name": schema.StringAttribute{
				Required:    true,
				Description: "Customer's last name.",
			},
			"email": schema.StringAttribute{
				Required:    true,
				Description: "Customer's unique email address.",
			},
			"account_type": schema.StringAttribute{
				Required:    true,
				Description: "Type of bank account (e.g., savings, checking).",
			},
			"balance": schema.Float64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Current balance of the customer account.",
			},
		},
	}
}

// **Create: Inserts new customer account into the database**
func (r *customerAccountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan customerAccountPlan

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	accountID, err := r.client.CreateCustomerAccount(
		plan.FirstName.ValueString(),
		plan.LastName.ValueString(),
		plan.Email.ValueString(),
		plan.AccountType.ValueString(),
		plan.Balance.ValueFloat64(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Database Error", fmt.Sprintf("Failed to create account: %s", err.Error()))
		return
	}

	// ✅ Ensure `id` is properly stored in Terraform State
	plan.ID = types.StringValue(fmt.Sprintf("%d", accountID))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// **Read: Fetches customer account details from database**
func (r *customerAccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state customerAccountPlan

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	account, err := r.client.GetCustomerAccount(state.Email.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Database Error", fmt.Sprintf("Could not retrieve account: %s", err.Error()))
		return
	}

	// ✅ Ensure `id` is correctly set
	state.ID = types.StringValue(fmt.Sprintf("%d", account.ID))
	state.FirstName = types.StringValue(account.FirstName)
	state.LastName = types.StringValue(account.LastName)
	state.Email = types.StringValue(account.Email)
	state.AccountType = types.StringValue(account.AccountType)
	state.Balance = types.Float64Value(account.Balance)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// **Update: Modifies an existing customer account**
func (r *customerAccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan customerAccountPlan

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.UpdateCustomerAccount(
		plan.ID.ValueString(),
		plan.FirstName.ValueString(),
		plan.LastName.ValueString(),
		plan.Email.ValueString(),
		plan.AccountType.ValueString(),
		plan.Balance.ValueFloat64(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Database Error", fmt.Sprintf("Could not update account: %s", err.Error()))
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r *customerAccountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state customerAccountPlan

	// ✅ Read full state before deleting
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// ✅ Delete the customer account from the database
	err := r.client.DeleteCustomerAccount(state.Email.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Database Error", fmt.Sprintf("Could not delete account: %s", err.Error()))
		return
	}

	// ✅ Remove state after successful deletion
	resp.State.RemoveResource(ctx)
}
