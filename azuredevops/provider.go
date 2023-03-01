package azuredevops

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	exclusivelock "github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/checks/exclusivelock/resource"
	invokerestapi "github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/checks/invokerestapi/resource"
	manualapproval "github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/checks/manualapproval/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/githubapp"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/permissions"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/serviceendpoint"
)

// Provider - The top level Azure DevOps Provider definition.
func Provider() *schema.Provider {
	p := &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"azuredevops_build_permissions":              permissions.ResourcePipelinePermissions(),
			"azuredevops_check_invokerestapi":            invokerestapi.ResourceCheckInvokeRestAPI(),
			"azuredevops_check_manualapproval":           manualapproval.ResourceCheckManualApproval(),
			"azuredevops_check_exclusivelock":            exclusivelock.ResourceCheckExclusiveLock(),
			"azuredevops_serviceendpoint_genericwebhook": serviceendpoint.ResourceServiceEndpointGenericWebhook(),
			"azuredevops_serviceendpoint_babylonawsiam":  serviceendpoint.ResourceServiceEndpointBabylonAwsIam(),
			"azuredevops_serviceendpoint_babylonvault":   serviceendpoint.ResourceServiceEndpointBabylonVault(),
			"azuredevops_serviceendpoint_githubapp":      githubapp.ResourceGithubApp(),
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

	p.ConfigureFunc = providerConfigure(p)

	return p
}

func providerConfigure(p *schema.Provider) schema.ConfigureFunc {
	return func(d *schema.ResourceData) (interface{}, error) {
		terraformVersion := p.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}

		return client.GetAzdoClient(d.Get("personal_access_token").(string), d.Get("org_service_url").(string), terraformVersion)
	}
}
