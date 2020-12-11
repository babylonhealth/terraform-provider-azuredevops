locals {
  azure_org = "NiteshLall"
}

provider "azuredevops" {
  org_service_url = "https://dev.azure.com/${local.azure_org}/"
}

resource "azuredevops_project" "project" {
  project_name = "nitesh-custom"
}

resource "azuredevops_serviceendpoint_github" "gh" {
  project_id = azuredevops_project.project.id
  service_endpoint_name = "github"

}

resource "azuredevops_gitapp" "app" {
  project_id = azuredevops_project.project.id
  connection_id = azuredevops_serviceendpoint_github.gh.id
  repo = "babylonhealth/pipeline-historian"
}