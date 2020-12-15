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

resource "azuredevops_serviceendpoint_github" "gh" {
  project_id = azuredevops_project.project.id
  service_endpoint_name = "github"

}

resource "azuredevops_serviceendpoint_githubapp" "app" {
  project_id = azuredevops_project.project.id
  connection_id = azuredevops_serviceendpoint_github.gh.id
  repo = "babylonhealth/pipeline-historian"
}

resource "azuredevops_build_permissions" "admin_permission3" {
    project_id      = "4f7f5d92-0e11-4311-ac85-9972864acbc2"
    principal = "vssgp.Uy0xLTktMTU1MTM3NDI0NS0yNDU1NjAwOTc1LTI4NjEzNDU5NS0yODk0NDM2NzIyLTIyNTMwNDg3NzAtMS0zNzc2MDEyNDM3LTI1MjAyNTQ3OTAtMjYxOTIwMDAzOS0yNTg5OTY1NzE4"
    project_level = true
    replace         = true
    permissions = {
        ViewBuilds                          = "allow"
        EditBuildQuality                    = "allow"
        RetainIndefinitely                  = "allow"
        DeleteBuilds                        = "allow"
        ManageBuildQualities                = "allow"
        DestroyBuilds                       = "allow"
        UpdateBuildInformation              = "allow"
        QueueBuilds                         = "allow"
        ManageBuildQueue                    = "allow"
        StopBuilds                          = "allow"
        ViewBuildDefinition                 = "allow"
        EditBuildDefinition                 = "allow"
        DeleteBuildDefinition               = "allow"
        OverrideBuildCheckInValidation      = "allow"
        AdministerBuildPermissions          = "deny"
    }
}

resource "azuredevops_build_permissions" "admin_permission4" {
    project_id      = "4f7f5d92-0e11-4311-ac85-9972864acbc2"
    principal = "vssgp.Uy0xLTktMTU1MTM3NDI0NS0yNDU1NjAwOTc1LTI4NjEzNDU5NS0yODk0NDM2NzIyLTIyNTMwNDg3NzAtMS0zNzc2MDEyNDM3LTI1MjAyNTQ3OTAtMjYxOTIwMDAzOS0yNTg5OTY1NzE4"
    project_level = false
    build_id = 39
    replace         = true
    permissions = {
        ViewBuilds                          = "allow"
        EditBuildQuality                    = "allow"
        RetainIndefinitely                  = "allow"
        DeleteBuilds                        = "allow"
        ManageBuildQualities                = "allow"
        DestroyBuilds                       = "allow"
        UpdateBuildInformation              = "allow"
        QueueBuilds                         = "allow"
        ManageBuildQueue                    = "allow"
        StopBuilds                          = "allow"
        ViewBuildDefinition                 = "allow"
        EditBuildDefinition                 = "allow"
        DeleteBuildDefinition               = "allow"
        OverrideBuildCheckInValidation      = "allow"
        AdministerBuildPermissions          = "deny"
    }
}