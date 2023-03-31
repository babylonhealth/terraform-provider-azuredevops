package serviceendpoint

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/bblnazuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/bblnazuredevops/internal/utils/tfhelper"
)

const BABYLON_VAULT_SERVICE_CONNECTION_TYPE string = "babylon-service-endpoint-vault"

func ResourceServiceEndpointBabylonVault() *schema.Resource {
	r := genBaseServiceEndpointResource(flattenServiceEndpointBabylonVault, expandServiceEndpointBabylonVault)
	r.Schema["url"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		ValidateFunc: func(i interface{}, key string) (_ []string, errors []error) {
			url, ok := i.(string)
			if !ok {
				errors = append(errors, fmt.Errorf("expected type of %q to be string", key))
				return
			}
			if strings.HasSuffix(url, "/") {
				errors = append(errors, fmt.Errorf("%q should not end with slash, got %q.", key, url))
				return
			}
			return validation.IsURLWithHTTPorHTTPS(url, key)
		},
		Description: "Url for the Vault Server",
	}

	r.Schema["vault_role"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Vault role to log in as",
	}
	return r
}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointBabylonVault(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *string, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)
	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{},
		Scheme:     converter.String("None"),
	}
	serviceEndpoint.Data = &map[string]string{"vaultRole": d.Get("vault_role").(string)}
	serviceEndpoint.Type = converter.String(BABYLON_VAULT_SERVICE_CONNECTION_TYPE)
	serviceEndpoint.Url = converter.String(d.Get("url").(string))
	return serviceEndpoint, projectID, nil
}

// Convert AzDO data structure to internal Terraform data structure
func flattenServiceEndpointBabylonVault(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *string) {
	doBaseFlattening(d, serviceEndpoint, projectID)
	tfhelper.HelpFlattenSecret(d, "password")

	d.Set("url", *serviceEndpoint.Url)
	d.Set("vault_role", (*serviceEndpoint.Data)["vaultRole"])
}
