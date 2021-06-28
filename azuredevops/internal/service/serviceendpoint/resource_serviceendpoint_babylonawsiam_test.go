package serviceendpoint

import (
	"fmt"
	"testing"

	"github.com/go-test/deep"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

func TestResourceServiceEndpointBabylonAwsIam(t *testing.T) {
	tests := []struct {
		name           string
		expectedSchema map[string]*schema.Schema
	}{
		{
			name: "test",
			expectedSchema: map[string]*schema.Schema{
				"username": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "AWS Access Key ID of the IAM user",
				},
				"password": {
					Type:             schema.TypeString,
					Required:         true,
					Description:      "AWS Secret Access Key of the IAM user",
					Sensitive:        true,
					DiffSuppressFunc: tfhelper.DiffFuncSuppressSecretChanged,
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
				"password_hash": {
					Type:        schema.TypeString,
					Computed:    true,
					Default:     nil,
					Description: fmt.Sprintf("A bcrypted hash of the attribute '%s'", "password"),
					Sensitive:   true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			resource := ResourceServiceEndpointBabylonAwsIam()
			resourceSchema := resource.Schema

			if diff := deep.Equal(resourceSchema, tt.expectedSchema); len(diff) > 0 {
				t.Errorf("ResourceServiceEndpointBabylonAwsIam() mismatch:\n%s", diff)
			}
		})
	}
}

func Test_expandServiceEndpointBabylonAwsIam(t *testing.T) {
	type args struct {
		username string
		password string
		project  string
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
				username: "user",
				password: "password",
				project:  "project",
			},
			want: &serviceendpoint.ServiceEndpoint{
				Authorization: &serviceendpoint.EndpointAuthorization{
					Parameters: &map[string]string{
						"username": "user",
						"password": "password",
					},
					Scheme: converter.String("UsernamePassword"),
				},
				Data:        &map[string]string{},
				Description: converter.String("Managed by Terraform"),
				Owner:       converter.String("library"),
				Type:        converter.String("babylon-service-endpoint-aws-iam"),
				Name:        converter.String(""),
				Url:         converter.String("https://s3.amazonaws.com/"),
			},
			wantProject: converter.String("project"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := ResourceServiceEndpointBabylonAwsIam()
			resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)

			multiErr := &multierror.Error{}

			err := resourceData.Set("username", tt.args.username)
			if err != nil {
				multiErr = multierror.Append(err, multiErr.Errors...)
			}

			err = resourceData.Set("password", tt.args.password)
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

			got, got1, err := expandServiceEndpointBabylonAwsIam(resourceData)
			if (err != nil) != tt.wantErr {
				t.Errorf("expandServiceEndpointBabylonAwsIam() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := deep.Equal(got, tt.want); len(diff) > 0 {
				t.Errorf("expandServiceEndpointBabylonAwsIam() mismatch:\n%s", diff)
			}

			if diff := deep.Equal(got1, tt.wantProject); len(diff) > 0 {
				t.Errorf("expandServiceEndpointBabylonAwsIam() got1 = %v, want %v", got1, tt.wantProject)
			}
		})
	}
}

func Test_flattenServiceEndpointBabylonAwsIam(t *testing.T) {
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
					Id:  converter.UUID("1ceae7ff-565c-4cdf-9214-6e2246cba764"),
					Url: converter.String("https://s3.amazonaws.com/"),
					Authorization: &serviceendpoint.EndpointAuthorization{
						Parameters: &map[string]string{
							"username": "user1",
							"password": "password1",
						},
						Scheme: converter.String("UsernamePassword"),
					},
				},
				projectID: converter.String("project"),
			},
			expected: map[string]string{
				"id":                    "1ceae7ff-565c-4cdf-9214-6e2246cba764",
				"authorization.%":       "1",
				"authorization.scheme":  "UsernamePassword",
				"description":           "",
				"password":              "password1",
				"project_id":            "project",
				"service_endpoint_name": "",
				"username":              "user1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := ResourceServiceEndpointBabylonAwsIam()

			resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)

			flattenServiceEndpointBabylonAwsIam(resourceData, tt.args.serviceEndpoint, tt.args.projectID)
			state := resourceData.State()

			if diff := deep.Equal(tt.expected, state.Attributes); len(diff) > 0 {
				t.Errorf("mismatch:\n%s", diff)
			}
		})
	}
}
