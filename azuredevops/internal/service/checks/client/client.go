package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/checks"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type Client struct {
	baseUrl                 string
	client                  *http.Client
	authorization           string
}

type GetChecksPayload struct {
	ContributionIds     []string `json:"contributionIds"`
	DataProviderContext struct {
		Properties struct {
			ResourceID string `json:"resourceId"`
			SourcePage struct {
				RouteValues struct {
					Project string `json:"project"`
				} `json:"routeValues"`
			} `json:"sourcePage"`
		} `json:"properties"`
	} `json:"dataProviderContext"`
}

func (c *Client) GetCheckByID(projectID string, resourceID string, checkID int64) (checks.CheckConfigurationData, bool, error) {
	found := false

	checkList, err := c.GetChecks(projectID, resourceID)

	if err != nil {
		return checks.CheckConfigurationData{}, found, err
	}

	for _, tempCheck := range checkList {
		if tempCheck.CheckConfiguration.ID == checkID {
			found = true
			return tempCheck, found, nil
		}
	}

	return checks.CheckConfigurationData{}, found, fmt.Errorf("no check found with id: %v under resource: %s, in project: %s",
		checkID, resourceID, projectID)
}

func (c *Client) GetChecks(projectID string, resourceID string) ([]checks.CheckConfigurationData, error) {
	payload := GetChecksPayload{}
	payload.ContributionIds = []string{"ms.vss-pipelinechecks.checks-data-provider"}
	payload.DataProviderContext.Properties.ResourceID = resourceID
	payload.DataProviderContext.Properties.SourcePage.RouteValues.Project = projectID

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return []checks.CheckConfigurationData{}, err
	}

	url := "/_apis/Contribution/HierarchyQuery"
	respBytes, err := c.SendRequest("POST", url, string(jsonPayload))

	result := checks.HeirarchyResp{}
	err = json.Unmarshal(respBytes, &result)
	if err != nil {
		return []checks.CheckConfigurationData{}, err
	}

	return result.DataProviders.MsVssPipelinechecksChecksDataProvider.CheckConfigurationDataList, nil
}

func (c *Client) AddCheck(projectID string, resourceID string, check checks.InvokeRESTAPIValues) (checks.CheckConfiguration, error) {
	restAPIPayload := populateInvokeRestAPIPayload(resourceID, check)

	jsonPayload, err := json.Marshal(restAPIPayload)
	if err != nil {
		return checks.CheckConfiguration{}, err
	}

	url := fmt.Sprintf("/%s/_apis/pipelines/checks/configurations", projectID)
	respBytes, err := c.SendRequest("POST", url, string(jsonPayload))
	if err != nil {
		return checks.CheckConfiguration{}, err
	}

	checkConf := checks.CheckConfiguration{}

	err = json.Unmarshal(respBytes, &checkConf)
	if err != nil {
		return checks.CheckConfiguration{}, err
	}

	return checkConf, nil
}

func (c *Client) UpdateCheck(projectID string, resourceID string, checkID string, check checks.InvokeRESTAPIValues) (checks.CheckConfiguration, error) {
	restAPIPayload := populateInvokeRestAPIPayload(resourceID, check)
	restAPIPayload.ID = checkID

	jsonPayload, err := json.Marshal(restAPIPayload)
	if err != nil {
		return checks.CheckConfiguration{}, err
	}

	url := fmt.Sprintf("/%s/_apis/pipelines/checks/configurations/%s", projectID, checkID)
	respBytes, err := c.SendRequest("PATCH", url, string(jsonPayload))
	if err != nil {
		return checks.CheckConfiguration{}, err
	}

	checkConf := checks.CheckConfiguration{}

	err = json.Unmarshal(respBytes, &checkConf)
	if err != nil {
		return checks.CheckConfiguration{}, err
	}

	return checkConf, nil
}

func (c *Client) DeleteCheck(projectID string, checkID string) error {
	url := fmt.Sprintf("/%s/_apis/pipelines/checks/configurations/%s", projectID, checkID)
	_, err := c.SendRequest("DELETE", url, "")

	return err
}

func populateInvokeRestAPIPayload(resourceID string, check checks.InvokeRESTAPIValues) checks.InvokeRestAPICheckPayload {
	checkPayload := checks.NewInvokeRestCheckPayload()

	checkPayload.Settings.DisplayName = check.DisplayName

	checkPayload.Settings.Inputs.Method = check.Method
	checkPayload.Settings.Inputs.WaitForCompletion = strconv.FormatBool(check.UseCallback)
	checkPayload.Settings.Inputs.Body = check.Body
	checkPayload.Settings.Inputs.URLSuffix = check.UrlSuffix
	checkPayload.Settings.Inputs.SuccessCriteria = check.SuccessCriteria

	headersBytes, err := json.Marshal(check.Headers)
	if err != nil {
		logrus.Fatal(err)
	}

	checkPayload.Settings.Inputs.Headers = string(headersBytes)

	// set to linked resource
	checkPayload.Resource.ID = resourceID

	//set by user
	checkPayload.Settings.Inputs.ConnectedServiceName = check.ServiceConnectionId
	checkPayload.Timeout = check.Timeout
	checkPayload.Settings.RetryInterval = check.RetryInterval

	return checkPayload
}

func (c *Client) SendRequest(httpMethod string, url string, jsonPayload string) ([]byte, error) {
	req, err := http.NewRequest(httpMethod,
		c.baseUrl + url,
		bytes.NewBufferString(jsonPayload))
	if err != nil {
		return []byte{}, err
	}

	req.Header.Set("Authorization", c.authorization)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json;api-version=5.1-preview.1;excludeUrls=true;enumsAsNumbers=true;msDateFormat=true;noArrayWrap=true")

	resp, err := c.client.Do(req)
	if err != nil {
		return []byte{}, err
	}

	if resp.StatusCode > 399 {
		return []byte{}, fmt.Errorf("resp status code from azure 400 or above")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}

func NewClient(baseUrl string, auth string, timeout *time.Duration) *Client {
	defaultTime := time.Duration(60 * time.Second)
	if timeout == nil {
		timeout = &defaultTime
	}

	client := &http.Client{
		Timeout: *timeout,
	}

	return &Client{
		baseUrl:                 baseUrl,
		client:                  client,
		authorization:           auth,
	}
}

type ChecksClient interface {
	GetCheckByID(projectID string, resourceID string, checkID int64) (checks.CheckConfigurationData, bool, error)
	GetChecks(projectID string, resourceID string) ([]checks.CheckConfigurationData, error)
	AddCheck(projectID string, resourceID string, check checks.InvokeRESTAPIValues) (checks.CheckConfiguration, error)
	UpdateCheck(projectID string, resourceID string, checkID string, check checks.InvokeRESTAPIValues) (checks.CheckConfiguration, error)
	DeleteCheck(projectID string, checkID string) error
}