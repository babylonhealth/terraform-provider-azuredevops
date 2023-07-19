package bblnazuredevops

import (
	"context"
	"github.com/babylonhealth/terraform-provider-bblnazuredevops/bblnazuredevops/internal/client"
	exclusivelock "github.com/babylonhealth/terraform-provider-bblnazuredevops/bblnazuredevops/internal/service/checks/exclusivelock/resource"
	invokerestapi "github.com/babylonhealth/terraform-provider-bblnazuredevops/bblnazuredevops/internal/service/checks/invokerestapi/resource"
	manualapproval "github.com/babylonhealth/terraform-provider-bblnazuredevops/bblnazuredevops/internal/service/checks/manualapproval/resource"
	"github.com/babylonhealth/terraform-provider-bblnazuredevops/bblnazuredevops/internal/service/githubapp"
	"github.com/babylonhealth/terraform-provider-bblnazuredevops/bblnazuredevops/internal/service/permissions"
	"github.com/babylonhealth/terraform-provider-bblnazuredevops/bblnazuredevops/internal/service/serviceendpoint"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider - The top level Azure DevOps Provider definition.
func Provider() *schema.Provider {
	p := &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"bblnazuredevops_build_permissions":              permissions.ResourcePipelinePermissions(),
			"bblnazuredevops_check_invokerestapi":            invokerestapi.ResourceCheckInvokeRestAPI(),
			"bblnazuredevops_check_manualapproval":           manualapproval.ResourceCheckManualApproval(),
			"bblnazuredevops_check_exclusivelock":            exclusivelock.ResourceCheckExclusiveLock(),
			"bblnazuredevops_serviceendpoint_genericwebhook": serviceendpoint.ResourceServiceEndpointGenericWebhook(),
			"bblnazuredevops_serviceendpoint_babylonawsiam":  serviceendpoint.ResourceServiceEndpointBabylonAwsIam(),
			"bblnazuredevops_serviceendpoint_babylonvault":   serviceendpoint.ResourceServiceEndpointBabylonVault(),
			"bblnazuredevops_serviceendpoint_githubapp":      githubapp.ResourceGithubApp(),
		},
		DataSourcesMap: map[string]*schema.Resource{},
		Schema: map[string]*schema.Schema{
			"org_service_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AZDO_ORG_SERVICE_URL", nil),
				Description: "The url of the Azure DevOps instance which should be used.",
			},
			"personal_access_token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AZDO_PERSONAL_ACCESS_TOKEN", nil),
				Description: "The personal access token which should be used.",
				Sensitive:   true,
			},
		},
	}

	p.ConfigureContextFunc = providerConfigure(p)

	return p
}

func providerConfigure(p *schema.Provider) schema.ConfigureContextFunc {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		terraformVersion := p.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		azdoPAT := d.Get("personal_access_token").(string)
		organizationURL := d.Get("org_service_url").(string)

		azdoClient, err := client.GetAzdoClient(azdoPAT, organizationURL, terraformVersion)

		return azdoClient, diag.FromErr(err)
	}
}
