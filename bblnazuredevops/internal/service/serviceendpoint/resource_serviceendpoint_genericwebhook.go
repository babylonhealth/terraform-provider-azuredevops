package serviceendpoint

import (
	"github.com/babylonhealth/terraform-provider-bblnazuredevops/bblnazuredevops/internal/utils/converter"
	"github.com/babylonhealth/terraform-provider-bblnazuredevops/bblnazuredevops/internal/utils/tfhelper"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/serviceendpoint"
)

// ResourceServiceEndpointGenericWebhook schema and implementation for docker registry service endpoint resource
func ResourceServiceEndpointGenericWebhook() *schema.Resource {
	r := genBaseServiceEndpointResource(flattenServiceEndpointGenericWebhook, expandServiceEndpointGenericWebhook)
	r.Schema["url"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		DefaultFunc: schema.EnvDefaultFunc("AZDO_GENERIC_WEBHOOK_URL", nil),
		Description: "The endpoint URL",
	}
	r.Schema["username"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		DefaultFunc: schema.EnvDefaultFunc("AZDO_GENERIC_WEBHOOK_USERNAME", nil),
		Description: "The username for the endpoint",
	}
	r.Schema["password"] = &schema.Schema{
		Type:             schema.TypeString,
		Optional:         true,
		DefaultFunc:      schema.EnvDefaultFunc("AZDO_GENERIC_WEBHOOK_PASSWORD", nil),
		Description:      "The Password for the endpoint",
		Sensitive:        true,
		DiffSuppressFunc: tfhelper.DiffFuncSuppressSecretChanged,
	}
	secretHashKey, secretHashSchema := tfhelper.GenerateSecreteMemoSchema("password")
	r.Schema[secretHashKey] = secretHashSchema
	return r
}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointGenericWebhook(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *string, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)
	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"username": d.Get("username").(string),
			"password": d.Get("password").(string),
		},
		Scheme: converter.String("UsernamePassword"),
	}
	serviceEndpoint.Data = &map[string]string{}
	serviceEndpoint.Type = converter.String("generic")
	urlString := d.Get("url").(string)
	serviceEndpoint.Url = &urlString
	return serviceEndpoint, projectID, nil
}

// Convert AzDO data structure to internal Terraform data structure
func flattenServiceEndpointGenericWebhook(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *string) {
	doBaseFlattening(d, serviceEndpoint, projectID)
	tfhelper.HelpFlattenSecret(d, "password")

	// ToDo: test with CLI tool if behavior differs from env var and file input password
	d.Set("username", (*serviceEndpoint.Authorization.Parameters)["username"])
	d.Set("password", (*serviceEndpoint.Authorization.Parameters)["password"])
	d.Set("url", *serviceEndpoint.Url)
}
