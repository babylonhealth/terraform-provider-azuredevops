package githubappclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type GithubAppClient interface {
	GetGithubAppByID(projectID string, connectionID string) (GetGithubAppResponse, bool, error)
	AddGithubApp(projectID string, repo string, connectionID string) (string, error)
	DeleteGithubApp(projectID string, connectionID string) error
}

// NewGithubApp will return a GithubApp struct. This is used to create service connections based on Github Apps
// Using Github Apps is preferred to PATs as there are rate limits on PAT connections
func NewGithubApp(baseUrl string, auth string, timeout *time.Duration) *GithubApp {
	defaultTime := time.Duration(60 * time.Second)
	if timeout == nil {
		timeout = &defaultTime
	}

	client := &http.Client{
		Timeout: *timeout,
	}

	return &GithubApp{
		baseUrl:       baseUrl,
		client:        client,
		authorization: auth,
	}
}

func NewGitHubAppPayload(projectID string, repo string, connectionID string) GitHubAppPayload {
	return GitHubAppPayload{
		ContributionIds: []string{"ms.vss-build-web.app-serviceconnections-recommendation-data-provider"},
		DataProviderContext: DataProviderContext{
			Properties: Properties{
				SourceProvider: "github",
				RepositoryID:   repo,
				RepositoryName: repo,
				ConnectionID:   connectionID,
				StrongBoxKey:   "useWellKnownStrongBoxLocation",
				SourcePage: SourcePage{
					RouteValues: RouteValues{
						Project: projectID,
					},
				},
			},
		},
	}
}

func NewGetGithubAppPayload(projectID string, connectionID string) GitHubAppPayload {
	return GitHubAppPayload{
		ContributionIds: []string{"ms.vss-serviceEndpoints-web.service-endpoints-details-data-provider"},
		DataProviderContext: DataProviderContext{
			Properties: Properties{
				ProjectID:         projectID,
				ServiceEndpointID: connectionID,
			},
		},
	}
}

func (g *GithubApp) GetGithubAppByID(projectID string, connectionID string) (GetGithubAppResponse, bool, error) {
	payload := NewGetGithubAppPayload(projectID, connectionID)

	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return GetGithubAppResponse{}, false, err
	}

	acceptHeaders := "application/json;api-version=5.1-preview.1;excludeUrls=true;enumsAsNumbers=true;msDateFormat=true;noArrayWrap=true"

	url := "/_apis/Contribution/HierarchyQuery"
	resp, err := g.SendRequest("POST", url, string(payloadJson), acceptHeaders)
	if err != nil {
		return GetGithubAppResponse{}, false, err
	}

	addAppResp := GetGithubAppResponse{}
	err = json.Unmarshal(resp, &addAppResp)
	if err != nil {
		return GetGithubAppResponse{}, false, err
	}

	if addAppResp.DataProviderExceptions.MsVssServiceEndpointsWebServiceEndpointsDetailsDataProvider.Message != "" {
		return GetGithubAppResponse{}, false, fmt.Errorf("an exception occurred fetching the gitapp")
	}

	if addAppResp.DataProviders.MsVssServiceEndpointsWebServiceEndpointsDetailsDataProvider.
		ServiceEndpoint.Authorization.Scheme != "InstallationToken" {
		fmt.Errorf("service connection is not github app")
	}

	return addAppResp, true, err
}

func (g *GithubApp) AddGithubApp(projectID string, repo string, connectionID string) (string, error) {
	payload := NewGitHubAppPayload(projectID, repo, connectionID)

	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	acceptHeaders := "application/json;api-version=5.1-preview.1;excludeUrls=true;enumsAsNumbers=true;msDateFormat=true;noArrayWrap=true"

	url := "/_apis/Contribution/HierarchyQuery"
	resp, err := g.SendRequest("POST", url, string(payloadJson), acceptHeaders)
	if err != nil {
		return "", err
	}

	addAppResp := &AddGitAppResp{}
	err = json.Unmarshal(resp, addAppResp)
	if err != nil {
		return "", err
	}

	return addAppResp.DataProviders.MsVssBuildWebAppServiceconnectionsRecommendationDataProvider.CommonConnectionID, nil
}

func (g *GithubApp) DeleteGithubApp(projectID string, connectionID string) error {
	url := fmt.Sprintf("/_apis/serviceendpoint/endpoints/%s?projectIds=%s", connectionID, projectID)

	acceptHeaders := "application/json;api-version=6.0-preview.4;excludeUrls=true;enumsAsNumbers=true;msDateFormat=true;noArrayWrap=true"
	_, err := g.SendRequest("DELETE", url, "", acceptHeaders)

	return err
}

func (c *GithubApp) SendRequest(httpMethod string, url string, jsonPayload string, acceptHeaders string) ([]byte, error) {
	req, err := http.NewRequest(httpMethod,
		c.baseUrl+url,
		bytes.NewBufferString(jsonPayload))
	if err != nil {
		return []byte{}, err
	}

	req.Header.Set("Authorization", c.authorization)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", acceptHeaders)

	resp, err := c.client.Do(req)
	if err != nil {
		return []byte{}, err
	}

	if resp.StatusCode == 203 {
		return []byte{}, fmt.Errorf("resp status code from azure 203 - need auth")
	}

	if resp.StatusCode > 399 {
		return []byte{}, fmt.Errorf("resp status code from azure 400 or above: %v", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}
