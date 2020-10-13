locals {
  azure_org = "NiteshLall"
}

provider "azuredevops" {
  version = "= 0.0.4"
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

resource "azuredevops_check_manualapproval" "manual_approval1" {
  project_id = azuredevops_project.project.id
  resource_id = "02c325bc-f8ec-47cd-a466-374b2f8cd835"
  type = "endpoint"

  instructions = "instructions 4"
  timeout = 4320
  approvers = ["e081ade0-63c7-6dc6-b15a-646ce81036ba", "1d22f8ac-06cc-6c5d-a114-03882932bb56"]
  allow_self_approve = true

  minimum_approvers = 0
  approve_in_order = true
}


resource "azuredevops_check_manualapproval" "manual_approval2" {
  project_id = azuredevops_project.project.id
  resource_id = "02c325bc-f8ec-47cd-a466-374b2f8cd835"
  type = "endpoint"

  instructions = "instructions 5"
  timeout = 4320
  approvers = ["e081ade0-63c7-6dc6-b15a-646ce81036ba", "1d22f8ac-06cc-6c5d-a114-03882932bb56"]
  allow_self_approve = true

  minimum_approvers = 0
  approve_in_order = true
}
