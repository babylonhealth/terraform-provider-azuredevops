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
  service_endpoint_name = "webhook2"
  description = "dev-platform service connection to webhook"
  url = "https://http.cat"
  username = "cat"
  password = "dog"
}

resource "azuredevops_serviceendpoint_genericwebhook" "endpoint_webhook2" {
  project_id = azuredevops_project.project.id
  service_endpoint_name = "webhook3"
  description = "dev-platform service connection to webhook"
  url = "https://http.cat"
  username = "cat"
  password = "dog"
}

resource "azuredevops_check_invokerestapi" "invoke_rest1" {
  project_id = azuredevops_project.project.id
  resource_id = azuredevops_serviceendpoint_genericwebhook.endpoint_webhook.id
  service_connection_id = azuredevops_serviceendpoint_genericwebhook.endpoint_webhook.id

  timeout = 43200
  retry_interval = 5
  display_name = "cli1"
  method = "POST"
  use_callback = false
  body = "[]"
  url_suffix = "/batman2"
  success_criteria = "[]"
  headers = {
    Content-Type ="application/json"
    PlanUrl = "$(system.CollectionUri)"
    ProjectId = "$(system.TeamProjectId)"
    HubName = "$(system.HostType)"
    PlanId = "$(system.PlanId)"
    JobId = "$(system.JobId)"
    TimelineId = "$(system.TimelineId)"
    TaskInstanceId = "$(system.TaskInstanceId)"
    AuthToken = "$(system.AccessToken)"
    k = "v"
    cat = "dog"
  }
}

resource "azuredevops_check_invokerestapi" "invoke_rest2" {
  project_id = azuredevops_project.project.id
  resource_id = azuredevops_serviceendpoint_genericwebhook.endpoint_webhook2.id
  service_connection_id = azuredevops_serviceendpoint_genericwebhook.endpoint_webhook.id

  timeout = 43200
  retry_interval = 5
  display_name = "cli2"
  method = "POST"
  use_callback = false
  body = "[]"
  url_suffix = "/batman2"
  success_criteria = "[]"
  headers = {
    Content-Type ="application/json"
    PlanUrl = "$(system.CollectionUri)"
    ProjectId = "$(system.TeamProjectId)"
    HubName = "$(system.HostType)"
    PlanId = "$(system.PlanId)"
    JobId = "$(system.JobId)"
    TimelineId = "$(system.TimelineId)"
    TaskInstanceId = "$(system.TaskInstanceId)"
    AuthToken = "$(system.AccessToken)"
    k = "v1"
    cat = "dog2"
  }
}