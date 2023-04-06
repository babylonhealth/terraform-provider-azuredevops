package main

import (
	"github.com/babylonhealth/terraform-provider-bblnazuredevops/bblnazuredevops"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return bblnazuredevops.Provider()
		},
	})
}
