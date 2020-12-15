package permissions

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	securityhelper "github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/permissions/utils"
	"log"
)

// resourcePipelinePermissions schema and implementation for project permission resource
func ResourcePipelinePermissions() *schema.Resource {
	return &schema.Resource{
		Create: resourcePipelinePermissionsCreateOrUpdate,
		Read:   resourcePipelinePermissionsRead,
		Update: resourcePipelinePermissionsCreateOrUpdate,
		Delete: resourcePipelinePermissionsDelete,
		Schema: securityhelper.CreatePermissionResourceSchema(map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				ValidateFunc: validation.IsUUID,
				Required:     true,
				ForceNew:     true,
			},
			"build_id": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"project_level": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		}),
	}
}

func resourcePipelinePermissionsCreateOrUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	sn, err := securityhelper.NewSecurityNamespace(d, clients, securityhelper.SecurityNamespaceIDValues.Build, createBuildToken)
	if err != nil {
		return err
	}

	if err := securityhelper.SetPrincipalPermissions(d, sn, nil, false); err != nil {
		return err
	}

	return resourcePipelinePermissionsRead(d, m)
}

func resourcePipelinePermissionsRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	sn, err := securityhelper.NewSecurityNamespace(d, clients, securityhelper.SecurityNamespaceIDValues.Build, createBuildToken)
	if err != nil {
		return err
	}

	principalPermissions, err := securityhelper.GetPrincipalPermissions(d, sn)
	if err != nil {
		return err
	}
	if principalPermissions == nil {
		d.SetId("")
		log.Printf("[INFO] Permissions for ACL token %q not found. Removing from state", sn.GetToken())
		return nil
	}

	d.Set("permissions", principalPermissions.Permissions)
	return nil
}

func resourcePipelinePermissionsDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	sn, err := securityhelper.NewSecurityNamespace(d, clients, securityhelper.SecurityNamespaceIDValues.Build, createBuildToken)
	if err != nil {
		return err
	}

	if err := securityhelper.SetPrincipalPermissions(d, sn, &securityhelper.PermissionTypeValues.NotSet, true); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func createBuildToken(d *schema.ResourceData, clients *client.AggregatedClient) (string, error) {
	projectID, ok := d.GetOk("project_id")
	if !ok {
		return "", fmt.Errorf("failed to get 'project_id' from schema")
	}

	projectLevelBool := d.Get("project_level").(bool)

	buildID, buildIdOK := d.GetOk("build_id")
	if !projectLevelBool && !buildIdOK {
		return "", fmt.Errorf("build_id required when project_level is not true")
	}

	if projectLevelBool && buildIdOK {
		return "", fmt.Errorf("build_id cannot be set when project_level is true")
	}

	buildString := buildID.(string)
	if !projectLevelBool && buildString == "" {
		return "", fmt.Errorf("build_id cannot be empty when project_level is not true")
	}

	if buildString != "" {
		buildString = "/" + buildString
	}

	aclToken := fmt.Sprintf("%s%s", projectID.(string), buildString)
	return aclToken, nil
}
