package serviceendpoint

import (
	"fmt"
	"github.com/babylonhealth/terraform-provider-bblnazuredevops/bblnazuredevops/internal/utils/converter"
	"github.com/babylonhealth/terraform-provider-bblnazuredevops/bblnazuredevops/internal/utils/tfhelper"
	"github.com/go-test/deep"
	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/serviceendpoint"
	"golang.org/x/crypto/bcrypt"
	"os"
	"testing"
)

func TestResourceServiceEndpointGenericWebhook(t *testing.T) {
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
					DefaultFunc: schema.EnvDefaultFunc("AZDO_GENERIC_WEBHOOK_URL", nil),
					Description: "The endpoint URL",
				},
				"username": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("AZDO_GENERIC_WEBHOOK_USERNAME", nil),
					Description: "The username for the endpoint",
				},
				"password": {
					Type:             schema.TypeString,
					Optional:         true,
					DefaultFunc:      schema.EnvDefaultFunc("AZDO_GENERIC_WEBHOOK_PASSWORD", nil),
					Description:      "The Password for the endpoint",
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

			resource := ResourceServiceEndpointGenericWebhook()
			resourceSchema := resource.Schema

			if diff := deep.Equal(resourceSchema, tt.expectedSchema); len(diff) > 0 {
				t.Errorf("ResourceServiceEndpointGenericWebhook() mismatch:\n%s", diff)
			}
		})
	}
}

func Test_expandServiceEndpointGenericWebhook(t *testing.T) {
	type args struct {
		username string
		password string
		url      string
		project  string
	}
	tests := []struct {
		name        string
		args        args
		want        *serviceendpoint.ServiceEndpoint
		wantProject *uuid.UUID
		wantErr     bool
	}{
		{
			name: "test expandServiceEndpoint",
			args: args{
				username: "user",
				password: "password",
				url:      "http://http.cat",
				project:  "3c49c3b6-a06d-424d-a6b6-0cd375ee9261",
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
				Type:        converter.String("generic"),
				Url:         converter.String("http://http.cat"),
				Name:        converter.String(""),
				ServiceEndpointProjectReferences: &[]serviceendpoint.ServiceEndpointProjectReference{
					{
						Name: converter.String(""),
						ProjectReference: &serviceendpoint.ProjectReference{
							Id: converter.UUID("3c49c3b6-a06d-424d-a6b6-0cd375ee9261"),
						},
					},
				},
			},
			wantProject: converter.UUID("3c49c3b6-a06d-424d-a6b6-0cd375ee9261"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := ResourceServiceEndpointGenericWebhook()
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

			err = resourceData.Set("url", tt.args.url)
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

			got, got1, err := expandServiceEndpointGenericWebhook(resourceData)
			if (err != nil) != tt.wantErr {
				t.Errorf("expandServiceEndpointGenericWebhook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := deep.Equal(got, tt.want); len(diff) > 0 {
				t.Errorf("expandServiceEndpointGenericWebhook() mismatch:\n%s", diff)
			}

			if diff := deep.Equal(got1, tt.wantProject); len(diff) > 0 {
				t.Errorf("expandServiceEndpointGenericWebhook() got1 = %v, wantServiceEndpoint %v", got1, tt.wantProject)
			}
		})
	}
}

func Test_flattenServiceEndpointGenericWebhook(t *testing.T) {
	type args struct {
		d               *schema.ResourceData
		serviceEndpoint *serviceendpoint.ServiceEndpoint
		projectID       *uuid.UUID
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
					Url: converter.String("http://http.cat"),
					Authorization: &serviceendpoint.EndpointAuthorization{
						Parameters: &map[string]string{
							"username": "user1",
							"password": "password1",
						},
						Scheme: converter.String("UsernamePassword"),
					},
				},
				projectID: converter.UUID("3c49c3b6-a06d-424d-a6b6-0cd375ee9261"),
			},
			expected: map[string]string{
				"id":                    "1ceae7ff-565c-4cdf-9214-6e2246cba764",
				"authorization.%":       "1",
				"authorization.scheme":  "UsernamePassword",
				"description":           "",
				"password":              "password1",
				"project_id":            "3c49c3b6-a06d-424d-a6b6-0cd375ee9261",
				"service_endpoint_name": "",
				"url":                   "http://http.cat",
				"username":              "user1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := ResourceServiceEndpointGenericWebhook()
			os.Setenv("AZDO_GENERIC_WEBHOOK_PASSWORD", "env_input")
			defer func() {
				os.Unsetenv("AZDO_GENERIC_WEBHOOK_PASSWORD")
			}()

			resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)

			flattenServiceEndpointGenericWebhook(resourceData, tt.args.serviceEndpoint, tt.args.projectID)
			state := resourceData.State()

			err := bcrypt.CompareHashAndPassword([]byte(state.Attributes["password_hash"]), []byte("env_input"))
			if err != nil {
				t.Errorf("password hashing mismatch")
			}

			delete(state.Attributes, "password_hash")

			if diff := deep.Equal(tt.expected, state.Attributes); len(diff) > 0 {
				t.Errorf("mismatch:\n%s", diff)
			}
		})
	}
}
