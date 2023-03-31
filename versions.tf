terraform {
  required_providers {
    azuredevops = {
      source  = "microsoft/azuredevops"
      version = "= 0.4.0"
    }
    bblnazuredevops = {
      source  = "babylonhealth.com/babylonhealth/bblnazuredevops"
      version = "~> 0.0.4"
    }
  }
  required_version = ">= 0.13"
}