// resource_greeting.go is included in the
//
//	provider package.
package provider

import (
	// context: control lifetime of request.
	"context"

	// These are imports for Terraform's
	//	     plugin framwork.
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure Resource implementation satisfies required interfaces.
//
//	This is a compile-time check.
var _ resource.Resource = &Greeting2Resource{}

// Define GreetingResource as a struct, with no fields.
type Greeting2Resource struct{}

// Factory method for creating a new resource
func NewGreeting2Resource() resource.Resource {
	return &Greeting2Resource{}
}

// Metadata: assigns a name to the resource,
//
//	when referenced in a .tf file
func (r *Greeting2Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "greeting_message"
}

// Defines the schema for the resource
func (r *Greeting2Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A simple resource that displays 'Hello, world!'",
		Attributes: map[string]schema.Attribute{
			"message": schema.StringAttribute{
				Description: "The greeting message",
				// Not provided by the user
				//    Defined when the resource is created.
				Computed: true,
			},
		},
	}
}

///*******  Beginning of CRUD ********

// Responsible create the resource
//
//	(i.e., terraform apply).
func (r *Greeting2Resource) Create(ctx context.Context,
	req resource.CreateRequest, resp *resource.CreateResponse) {
	// Create a struct for the state with a
	//		Message field that maps to the message field
	//		in the schema.
	resp.State.Set(ctx, &struct {
		Message types.String `tfsdk:"message"`
	}{ // Immediately initialize the struct
		//    (i.e., Message field)
		Message: types.StringValue("Hello, world12!"),
	})
}

// Read does nothing since the message never changes
func (r *Greeting2Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

// Update does nothing since the message never changes
func (r *Greeting2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete does nothing since there is nothing to clean up
func (r *Greeting2Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
