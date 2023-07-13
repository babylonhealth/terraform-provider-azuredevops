package serviceendpoint

import (
	"github.com/google/uuid"
	"testing"

	"github.com/babylonhealth/terraform-provider-bblnazuredevops/bblnazuredevops/internal/utils/converter"
	"github.com/go-test/deep"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/serviceendpoint"
)

func TestResourceServiceEndpointBabylonVault(t *testing.T) {
	tests := []struct {
		name           string
		expectedSchema map[string]*schema.Schema
	}{
		{
			name: "test",
			expectedSchema: map[string]*schema.Schema{
				"url": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Url for the Vault Server",
				},
				"vault_role": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Vault role to log in as",
				},
				"project_id": {
					Type:         schema.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotWhiteSpace,
				},
				"service_endpoint_name": {
					Type:         schema.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotWhiteSpace,
				},
				"description": {
					Type:         schema.TypeString,
					Optional:     true,
					Default:      "Managed by Terraform",
					ValidateFunc: validation.StringIsNotWhiteSpace,
				},
				"authorization": {
					Type:         schema.TypeMap,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.StringIsNotWhiteSpace,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			resource := ResourceServiceEndpointBabylonVault()
			resourceSchema := resource.Schema

			if diff := deep.Equal(resourceSchema, tt.expectedSchema); len(diff) > 0 {
				t.Errorf("ResourceServiceEndpointBabylonVault() mismatch:\n%s", diff)
			}
		})
	}
}

func Test_expandServiceEndpointBabylonVault(t *testing.T) {
	type args struct {
		url       string
		vaultRole string
		project   string
	}
	tests := []struct {
		name        string
		args        args
		want        *serviceendpoint.ServiceEndpoint
		wantProject *string
		wantErr     bool
	}{
		{
			name: "test expandServiceEndpoint",
			args: args{
				url:       "https://vault.babylonhealth.com",
				vaultRole: "devtest",
				project:   "project",
			},
			want: &serviceendpoint.ServiceEndpoint{
				Authorization: &serviceendpoint.EndpointAuthorization{
					Parameters: &map[string]string{},
					Scheme:     converter.String("None"),
				},
				Data: &map[string]string{
					"vaultRole": "devtest",
				},
				Description: converter.String("Managed by Terraform"),
				Owner:       converter.String("library"),
				Type:        converter.String(BABYLON_VAULT_SERVICE_CONNECTION_TYPE),
				Name:        converter.String(""),
				Url:         converter.String("https://vault.babylonhealth.com"),
			},
			wantProject: converter.String("project"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := ResourceServiceEndpointBabylonVault()
			resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)

			multiErr := &multierror.Error{}

			err := resourceData.Set("url", tt.args.url)
			if err != nil {
				multiErr = multierror.Append(err, multiErr.Errors...)
			}

			err = resourceData.Set("vault_role", tt.args.vaultRole)
			if err != nil {
				multiErr = multierror.Append(err, multiErr.Errors...)
			}

			err = resourceData.Set("project_id", tt.args.project)
			if err != nil {
				multiErr = multierror.Append(err, multiErr.Errors...)
			}

			if err != nil {
				t.Error(err)
			}

			got, got1, err := expandServiceEndpointBabylonVault(resourceData)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResourceServiceEndpointBabylonVault() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := deep.Equal(got, tt.want); len(diff) > 0 {
				t.Errorf("ResourceServiceEndpointBabylonVault() mismatch:\n%s", diff)
			}

			if diff := deep.Equal(got1, tt.wantProject); len(diff) > 0 {
				t.Errorf("ResourceServiceEndpointBabylonVault() got1 = %v, want %v", got1, tt.wantProject)
			}
		})
	}
}

func Test_flattenServiceEndpointBabylonVault(t *testing.T) {
	type args struct {
		d               *schema.ResourceData
		serviceEndpoint *serviceendpoint.ServiceEndpoint
		projectID       *string
	}
	tests := []struct {
		name     string
		args     args
		expected map[string]string
	}{
		{
			name: "Test flattenServiceEndpoint - secret updated",
			args: args{
				d: &schema.ResourceData{},
				serviceEndpoint: &serviceendpoint.ServiceEndpoint{
					Id:   converter.UUID("1ceae7ff-565c-4cdf-9214-6e2246cba764"),
					Url:  converter.String("https://vault.babylonhealth.com"),
					Data: &map[string]string{"vaultRole": "devtest"},
					Authorization: &serviceendpoint.EndpointAuthorization{
						Parameters: &map[string]string{},
						Scheme:     converter.String("None"),
					},
				},
				projectID: converter.String("project"),
			},
			expected: map[string]string{
				"id":                    "1ceae7ff-565c-4cdf-9214-6e2246cba764",
				"authorization.%":       "1",
				"authorization.scheme":  "None",
				"description":           "",
				"url":                   "https://vault.babylonhealth.com",
				"vault_role":            "devtest",
				"project_id":            "project",
				"service_endpoint_name": "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := ResourceServiceEndpointBabylonVault()

			resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)

			flattenServiceEndpointBabylonVault(resourceData, tt.args.serviceEndpoint, tt.args.projectID)
			state := resourceData.State()

			if diff := deep.Equal(tt.expected, state.Attributes); len(diff) > 0 {
				t.Errorf("mismatch:\n%s", diff)
			}
		})
	}
}
