// set TF_LOG=DEBUG 
// 
package main
import (
	"context"
	"log"

	"github.com/donis/terraform-provider-filemanager/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/donis/filemanager",
	}

	err := providerserver.Serve(context.Background(), provider.New, opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
