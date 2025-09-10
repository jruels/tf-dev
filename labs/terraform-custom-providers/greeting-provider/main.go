// C:\Users\donis\AppData\Roaming\terraform.rc
// This is an executable and the entrypoint for
//		the custom provider.
//	 go mod init {name}
package main

import (
	// context:controls the lifetime of a
	//    request, such as cancellations or timeouts.
	// flag: read flags from the command line
	// log: diagnostic messages

	"context"
	"flag"
	"log"

	// Identifies the provider code
	"registry.terraform.io/donis/greeting/internal/provider"
	// References the "plumbing" to connect
	//	 a custom provider with Terraform.
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	// debug: hold the result of the --debug flag
	var debug bool

	// Create the --debug command line flag.
	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	// Update the debug variable with the flag status
	flag.Parse()

	// Create options for the provider
	//		Address: names the provider
	//		Debug: sets the debug mode
	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/donis/greeting",
		Debug:   debug,
	}

	// Registers the provider factory with
	//        Terraform using gRPC.
	//	contextBackground: no context
	//  provider.New: factory method
	//  opts:  provider options set previosly
	err := providerserver.Serve(context.Background(), provider.New, opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
