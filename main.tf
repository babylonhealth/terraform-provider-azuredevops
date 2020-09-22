locals {
  azure_org = "NiteshLall"
}

provider "azuredevops" {
  version = "= 0.0.3"
  org_service_url = "https://dev.azure.com/${local.azure_org}/"
}

resource "azuredevops_project" "project" {
  project_name = "nitesh-custom"
}

resource "azuredevops_serviceendpoint_genericwebhook" "endpoint_webhook" {
  project_id = azuredevops_project.project.id
  service_endpoint_name = "webhook1"
  description = "dev-platform service connection to webhook"
  url = "https://http.cat"
  username = "cat"
  password = "dog"
}