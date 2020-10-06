package client

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/checks"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestClient_DeleteCheck(t *testing.T) {
	type fields struct {
		baseUrl                 string
		client                  *http.Client
		authorization           string
		suppressFedAuthRedirect bool
		forceMsaPassThrough     bool
		userAgent               string
	}
	type args struct {
		projectID string
		checkID   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Delete test",
			args: args{
				projectID: "4f7f5d92-0e11-4311-ac85-9972864acbc2",
				checkID:   "46",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			personalAccessToken := getAuthString()

			duration := 60 * time.Second

			ts := getTestServer(struct{}{})
			defer ts.Close()

			c := NewClient(ts.URL, personalAccessToken, &duration)

			if err := c.DeleteCheck(tt.args.projectID, tt.args.checkID); (err != nil) != tt.wantErr {
				t.Errorf("DeleteCheck() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_AddCheck(t *testing.T) {
	type fields struct {
		baseUrl                 string
		client                  *http.Client
		authorization           string
		suppressFedAuthRedirect bool
		forceMsaPassThrough     bool
		userAgent               string
	}
	type args struct {
		projectID  string
		resourceID string
		check      checks.InvokeRESTAPIValues
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "Add test",
			args: args{
				projectID:  "4f7f5d92-0e11-4311-ac85-9972864acbc2",
				resourceID: "02c325bc-f8ec-47cd-a466-374b2f8cd835",
				check: checks.InvokeRESTAPIValues{
					ServiceConnectionId: "02c325bc-f8ec-47cd-a466-374b2f8cd835",
					LinkedVariableGroup: "",
					Timeout:             43200,
					RetryInterval:       5,
					DisplayName:         "Terraform test",
					Method:              "POST",
					UseCallback:         false,
					Body:                "{}",
					UrlSuffix:           "",
					SuccessCriteria:     "",
					Headers: map[string]string{
						"k": "v",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			personalAccessToken := getAuthString()

			duration := 60 * time.Second
			ts := getTestServer(populateInvokeRestAPIPayload(tt.args.resourceID, tt.args.check))
			defer ts.Close()

			c := NewClient(ts.URL, personalAccessToken, &duration)

			got, err := c.AddCheck(tt.args.projectID, tt.args.resourceID, tt.args.check)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddCheck() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			expectedResp := populateCheckConf(tt.args.check)

			if !reflect.DeepEqual(got, expectedResp) {
				t.Errorf("AddCheck() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_UpdateCheck(t *testing.T) {
	type fields struct {
		baseUrl                 string
		client                  *http.Client
		authorization           string
		suppressFedAuthRedirect bool
		forceMsaPassThrough     bool
		userAgent               string
	}
	type args struct {
		projectID  string
		resourceID string
		checkID    string
		check      checks.InvokeRESTAPIValues
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantResp checks.CheckConfiguration
		wantErr  bool
	}{
		{
			name: "Update test",
			args: args{
				projectID:  "4f7f5d92-0e11-4311-ac85-9972864acbc2",
				resourceID: "02c325bc-f8ec-47cd-a466-374b2f8cd835",
				checkID:    "57",
				check: checks.InvokeRESTAPIValues{
					ServiceConnectionId: "02c325bc-f8ec-47cd-a466-374b2f8cd835",
					LinkedVariableGroup: "",
					Timeout:             43200,
					RetryInterval:       5,
					DisplayName:         "Update test",
					Method:              "POST",
					UseCallback:         false,
					Body:                "{}",
					UrlSuffix:           "",
					SuccessCriteria:     "",
					Headers: map[string]string{
						"k": "v",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := getTestServer(populateInvokeRestAPIPayload(tt.args.resourceID, tt.args.check))
			defer ts.Close()

			duration := 60 * time.Second
			c := NewClient(ts.URL, "", &duration)
			gotResp, err := c.UpdateCheck(tt.args.projectID, tt.args.resourceID, tt.args.checkID, tt.args.check)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateCheck() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			expectedResp := populateCheckConf(tt.args.check)

			if diff := cmp.Diff(expectedResp, gotResp); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestClient_GetCheckByID(t *testing.T) {
	type fields struct {
		baseUrl                 string
		client                  *http.Client
		authorization           string
		suppressFedAuthRedirect bool
		forceMsaPassThrough     bool
		userAgent               string
	}
	type args struct {
		projectID  string
		resourceID string
		checkID    int64
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      checks.CheckConfigurationData
		wantErr   bool
		wantFound bool
	}{
		{
			name: "client test",
			args: args{
				checkID: 50,
			},
			want: checks.CheckConfigurationData{
				CheckConfiguration: checks.CheckConfiguration{
					ID:  50,
					URL: "test",
				},
			},
			wantFound: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hr := checks.HeirarchyResp{}
			configData := []checks.CheckConfigurationData{tt.want}
			hr.DataProviders.MsVssPipelinechecksChecksDataProvider.CheckConfigurationDataList = configData

			ts := getTestServer(hr)
			defer ts.Close()

			duration := 60 * time.Second
			c := NewClient(ts.URL, "", &duration)

			got, found, err := c.GetCheckByID(tt.args.projectID, tt.args.resourceID, tt.args.checkID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetCheckByID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCheckByID() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(found, tt.wantFound) {
				t.Errorf("GetCheckByID() found = %v, want %v", found, tt.wantFound)
			}
		})
	}
}

func getAuthString() string {
	auth := ":" + os.Getenv("TEST_TOKEN")
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

func getTestServer(wantedResponse interface{}) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResp, err := json.Marshal(wantedResponse)
		if err != nil {
			logrus.Fatalf("error setting up test server: %v", err)
		}

		fmt.Fprintf(w, string(jsonResp))
	}))

	return ts
}

func populateCheckConf(values checks.InvokeRESTAPIValues) checks.CheckConfiguration {
	conf := checks.CheckConfiguration{}
	conf.Resource.Type = "endpoint"
	conf.Resource.ID = "02c325bc-f8ec-47cd-a466-374b2f8cd835"

	conf.Type.ID = "fe1de3ee-a436-41b4-bb20-f6eb4cb879a7"
	conf.Type.Name = "Task Check"

	conf.Settings.DefinitionRef.ID = "9c3e8943-130d-4c78-ac63-8af81df62dfb"
	conf.Settings.DefinitionRef.Name = "InvokeRESTAPI"
	conf.Settings.DefinitionRef.Version = "1.152.3"

	conf.Timeout = values.Timeout
	conf.Settings.RetryInterval = values.RetryInterval
	conf.Settings.DisplayName = values.DisplayName

	headersBytes, err := json.Marshal(values.Headers)
	if err != nil {
		logrus.Fatal(err)
	}

	conf.Settings.Inputs.Headers = string(headersBytes)
	conf.Settings.Inputs.SuccessCriteria = values.SuccessCriteria
	conf.Settings.Inputs.URLSuffix = values.UrlSuffix
	conf.Settings.Inputs.Body = values.Body
	conf.Settings.Inputs.WaitForCompletion = strconv.FormatBool(values.UseCallback)
	conf.Settings.Inputs.Method = values.Method
	conf.Settings.Inputs.ConnectedServiceName = values.ServiceConnectionId
	conf.Settings.Inputs.ConnectedServiceNameSelector = "connectedServiceName"

	return conf
}
