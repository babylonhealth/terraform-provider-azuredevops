package githubappclient

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestGitApp_AddGithubApp(t *testing.T) {
	type args struct {
		projectID    string
		repo         string
		connectionID string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			args: args{
				projectID:    "4f7f5d92-0e11-4311-ac85-9972864acbc2",
				repo:         "babylonhealth/cd-workflows",
				connectionID: "627166ad-752b-47f7-a115-7bdcb385931e",
			},
			want: "azure",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := AddGitAppResp{}
			resp.DataProviders.MsVssBuildWebAppServiceconnectionsRecommendationDataProvider.CommonConnectionID = tt.want

			g := getRealClient()
			ts := getTestServer(resp)
			defer ts.Close()

			g.baseUrl = ts.URL

			got, err := g.AddGithubApp(tt.args.projectID, tt.args.repo, tt.args.connectionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddGithubApp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddGithubApp() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGitApp_DeleteGihubApp(t *testing.T) {
	type args struct {
		projectID    string
		connectionID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			args: args{
				projectID:    "4f7f5d92-0e11-4311-ac85-9972864acbc2",
				connectionID: "eeb8dca2-7c49-4bbd-aaf1-b27652864aa1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := getRealClient()
			ts := getTestServer(struct{}{})
			defer ts.Close()
			g.baseUrl = ts.URL

			if err := g.DeleteGithubApp(tt.args.projectID, tt.args.connectionID); (err != nil) != tt.wantErr {
				t.Errorf("DeleteGithubApp() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGitApp_GetGithubAppByID(t *testing.T) {
	type fields struct {
		baseUrl       string
		authorization string
	}
	type args struct {
		projectID    string
		connectionID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    GetGithubAppResponse
		want1   bool
		wantErr bool
	}{
		{
			args: args{
				projectID:    "4f7f5d92-0e11-4311-ac85-9972864acbc2",
				connectionID: "d2f95b0d-2107-4c9e-8b63-a3af05d5eea1",
			},
			wantErr: false,
			want1:   true,
			want: GetGithubAppResponse{
				DataProviders: DataProviders{MsVssServiceEndpointsWebServiceEndpointsDetailsDataProvider: MsVssServiceEndpointsWebServiceEndpointsDetailsDataProvider{
					ServiceEndpoint: ServiceEndpoint{
						Authorization: Authorization{
							Scheme: "InstallationToken",
						},
					},
				}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := getRealClient()
			ts := getTestServer(tt.want)
			defer ts.Close()
			g.baseUrl = ts.URL

			got, got1, err := g.GetGithubAppByID(tt.args.projectID, tt.args.connectionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGithubAppByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetGithubAppByID() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetGithubAppByID() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
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

func getRealClient() *GithubApp {
	auth := "Basic " + base64.StdEncoding.EncodeToString([]byte(":"+os.Getenv("TEST_TOKEN")))
	timeout := time.Minute
	return NewGithubApp(os.Getenv("TEST_URL"), auth, &timeout)
}
