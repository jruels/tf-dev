package main

import (
	"context"
	"log"

	"github.com/donis/terraform-provider-banking/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	err := providerserver.Serve(context.Background(), provider.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/example/banking",
	})

	if err != nil {
		log.Fatal(err)
	}
}
