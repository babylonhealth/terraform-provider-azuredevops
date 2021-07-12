package serviceendpoint

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

const SERVICE_CONNECTION_TYPE string = "babylon-service-endpoint-aws-iam"
const DEFAULT_SESSION_NAME string = "azure-pipelines-task"

func ResourceServiceEndpointBabylonAwsIam() *schema.Resource {
	r := genBaseServiceEndpointResource(flattenServiceEndpointBabylonAwsIam, expandServiceEndpointBabylonAwsIam)
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
	r.Schema["global_role_arn"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The Amazon Resource Name (ARN) of the role to assume",
	}
	r.Schema["global_sts_session_name"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Session name to be used when assuming the role. The session name should match the one specified in the trust policies of the regional IAM roles.",
		Default:     DEFAULT_SESSION_NAME,
	}
	secretHashKey, secretHashSchema := tfhelper.GenerateSecreteMemoSchema("password")
	r.Schema[secretHashKey] = secretHashSchema
	return r
}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointBabylonAwsIam(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *string, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)
	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"username":             d.Get("username").(string),
			"password":             d.Get("password").(string),
			"globalRoleArn":        d.Get("global_role_arn").(string),
			"globalStsSessionName": d.Get("global_sts_session_name").(string),
		},
		Scheme: converter.String("UsernamePassword"),
	}
	serviceEndpoint.Data = &map[string]string{}
	serviceEndpoint.Type = converter.String(SERVICE_CONNECTION_TYPE)
	serviceEndpoint.Url = converter.String("https://aws.amazon.com/")
	return serviceEndpoint, projectID, nil
}

// Convert AzDO data structure to internal Terraform data structure
func flattenServiceEndpointBabylonAwsIam(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *string) {
	doBaseFlattening(d, serviceEndpoint, projectID)
	tfhelper.HelpFlattenSecret(d, "password")

	d.Set("username", (*serviceEndpoint.Authorization.Parameters)["username"])
	d.Set("password", (*serviceEndpoint.Authorization.Parameters)["password"])
	d.Set("global_role_arn", (*serviceEndpoint.Authorization.Parameters)["globalRoleArn"])
	d.Set("global_sts_session_name", (*serviceEndpoint.Authorization.Parameters)["globalStsSessionName"])
}
