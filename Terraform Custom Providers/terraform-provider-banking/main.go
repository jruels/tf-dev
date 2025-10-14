package main

import (
	"context"
	"log"

	bankingprovider "github.com/donis/terraform-provider-banking/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	// ✅ Start the Terraform provider with the correct configuration
	err := providerserver.Serve(context.Background(), bankingprovider.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/example/banking",
	})

	if err != nil {
		log.Fatal(err)
	}
}
