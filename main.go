package main

import (
	"flag"
	"github.com/babylonhealth/terraform-provider-bblnazuredevops/bblnazuredevops"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{
		Debug:        debug,
		ProviderAddr: "babylonhealth.com/babylonhealth/bblnazuredevops",
		ProviderFunc: func() *schema.Provider {
			return bblnazuredevops.Provider()
		},
	}
	plugin.Serve(opts)
}
