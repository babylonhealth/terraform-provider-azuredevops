package serviceendpoint

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

func ResourceServiceEndpointBabylonAWSAssumeRole() *schema.Resource {
	r := genBaseServiceEndpointResource(flattenServiceEndpointBabylonAWSAssumeRole, expandServiceEndpointBabylonAWSAssumeRole)
	r.Schema["username"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "AWS Access Key ID of the IAM user",
	}
	r.Schema["password"] = &schema.Schema{
		Type:             schema.TypeString,
		Required:         true,
		Description:      "AWS Secret Access Key of the IAM user",
		Sensitive:        true,
		DiffSuppressFunc: tfhelper.DiffFuncSuppressSecretChanged,
	}
	secretHashKey, secretHashSchema := tfhelper.GenerateSecreteMemoSchema("password")
	r.Schema[secretHashKey] = secretHashSchema
	return r
}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointBabylonAWSAssumeRole(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *string, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)
	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"username": d.Get("username").(string),
			"password": d.Get("password").(string),
		},
		Scheme: converter.String("UsernamePassword"),
	}
	serviceEndpoint.Data = &map[string]string{}
	serviceEndpoint.Type = converter.String("babylon-service-endpoint-aws-iam")
	serviceEndpoint.Url = converter.String("https://s3.amazonaws.com/")
	return serviceEndpoint, projectID, nil
}

// Convert AzDO data structure to internal Terraform data structure
func flattenServiceEndpointBabylonAWSAssumeRole(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *string) {
	doBaseFlattening(d, serviceEndpoint, projectID)
	tfhelper.HelpFlattenSecret(d, "password")

	d.Set("username", (*serviceEndpoint.Authorization.Parameters)["username"])
	d.Set("password", (*serviceEndpoint.Authorization.Parameters)["password"])
}
