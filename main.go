package main

import (
	"github.com/babylonhealth/terraform-provider-bblnazuredevops/bblnazuredevops"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderAddr: "babylonhealth.com/babylonhealth/bblnazuredevops",
		ProviderFunc: func() *schema.Provider {
			return bblnazuredevops.Provider()
		},
	})
}
