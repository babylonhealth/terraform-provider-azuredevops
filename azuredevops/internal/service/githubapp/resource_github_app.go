package githubapp

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

func ResourceGithubApp() *schema.Resource {
	r := &schema.Resource{
		Create: createGithubApp,
		Read:   getGitHubApp,
		Update: updateApp,
		Delete: deleteApp,
		Exists: gitHubAppExists,
	}
	r.Schema = map[string]*schema.Schema{}
	r.Schema["project_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	r.Schema["connection_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	r.Schema["repo"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	r.Schema["app_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: false,
		Computed: true,
	}

	r.Importer = tfhelper.ImportProjectQualifiedResourceUUID()

	return r
}

// See Resource documentation.
func createGithubApp(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	connectionID := d.Get("connection_id").(string)
	repo := d.Get("repo").(string)

	appId, err := clients.GitAppClient.AddGithubApp(projectID, repo, connectionID)
	if err != nil {
		return fmt.Errorf("error creating Github App in Azure DevOps: %+v", err)
	}

	d.SetId(fmt.Sprintf("%v", appId))
	err = d.Set("app_id", appId)

	return err
}

func deleteApp(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	connectionID := d.Get("app_id").(string)

	err := clients.GitAppClient.DeleteGithubApp(projectID, connectionID)

	return err
}

func gitHubAppExists(d *schema.ResourceData, m interface{}) (bool, error) {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	id := d.Get("app_id").(string)

	_, found, err := clients.GitAppClient.GetGithubAppByID(projectID, id)
	if err != nil {
		return false, err
	}

	return found, err
}

func getGitHubApp(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	connectionID := d.Get("app_id").(string)

	resp, _, err := clients.GitAppClient.GetGithubAppByID(projectID, connectionID)
	if err != nil {
		return err
	}

	id := resp.DataProviders.MsVssServiceEndpointsWebServiceEndpointsDetailsDataProvider.ServiceEndpoint.ID

	d.SetId(fmt.Sprintf("%v", id))
	err = d.Set("app_id", id)

	return err
}

func updateApp(d *schema.ResourceData, m interface{}) error {
	return fmt.Errorf("github apps cannot be updated, delete then re-create")
}
