terraform {
  required_providers {
    azuredevops = {
      source  = "babylonhealth.com/babylonhealth/azuredevops"
      version = "~> 0.0.4"
    }
  }
  required_version = ">= 0.13"
}