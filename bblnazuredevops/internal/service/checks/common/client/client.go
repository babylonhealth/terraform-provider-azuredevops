package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	exclusivelockmodel "github.com/babylonhealth/terraform-provider-bblnazuredevops/bblnazuredevops/internal/service/checks/exclusivelock/model"
	invokerestapimodel "github.com/babylonhealth/terraform-provider-bblnazuredevops/bblnazuredevops/internal/service/checks/invokerestapi/model"
	manualapprovalmodel "github.com/babylonhealth/terraform-provider-bblnazuredevops/bblnazuredevops/internal/service/checks/manualapproval/model"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type Client struct {
	baseUrl       string
	client        *http.Client
	authorization string
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

func (c *Client) GetInvokeRestAPICheckByID(projectID string, resourceID string, checkID int64) (invokerestapimodel.CheckConfigurationData, bool, error) {
	found := false

	checkList, err := c.getInvokeRestAPIChecks(projectID, resourceID)

	if err != nil {
		return invokerestapimodel.CheckConfigurationData{}, found, err
	}

	for _, tempCheck := range checkList {
		if tempCheck.CheckConfiguration.ID == checkID {
			found = true
			return tempCheck, found, nil
		}
	}

	return invokerestapimodel.CheckConfigurationData{}, found, fmt.Errorf("no invokerestapi check found with id: %v under resource: %s, in project: %s",
		checkID, resourceID, projectID)
}

func (c *Client) GetManualApprovalCheckByID(projectID string, resourceID string, checkID int64) (manualapprovalmodel.ManualApprovalCheckConfig, bool, error) {
	found := false

	checkList, err := c.getManualApprovalChecks(projectID, resourceID)
	if err != nil {
		return manualapprovalmodel.ManualApprovalCheckConfig{}, found, err
	}

	for _, tempCheck := range checkList {
		if tempCheck.ID == checkID {
			found = true
			return tempCheck, found, nil
		}
	}

	return manualapprovalmodel.ManualApprovalCheckConfig{}, found, fmt.Errorf("no manualapproval check found with id: %v under resource: %s, in project: %s",
		checkID, resourceID, projectID)
}

func (c *Client) GetExclusiveLockCheckByID(projectID string, resourceID string, checkID int64) (exclusivelockmodel.ExclusiveLockCheckConfig, bool, error) {
	found := false

	checkList, err := c.getExclusiveLockChecks(projectID, resourceID)
	if err != nil {
		return exclusivelockmodel.ExclusiveLockCheckConfig{}, found, err
	}

	for _, tempCheck := range checkList {
		if tempCheck.ID == checkID {
			found = true
			return tempCheck, found, nil
		}
	}

	return exclusivelockmodel.ExclusiveLockCheckConfig{}, found, fmt.Errorf("no exclusivelock check found with id: %v under resource: %s, in project: %s",
		checkID, resourceID, projectID)
}

func (c *Client) getAllChecks(projectID string, resourceID string) ([]byte, error) {
	payload := GetChecksPayload{}
	payload.ContributionIds = []string{"ms.vss-pipelinechecks.checks-data-provider"}
	payload.DataProviderContext.Properties.ResourceID = resourceID
	payload.DataProviderContext.Properties.SourcePage.RouteValues.Project = projectID

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return []byte{}, err
	}

	url := "/_apis/Contribution/HierarchyQuery"
	respBytes, err := c.SendRequest("POST", url, string(jsonPayload))
	if err != nil {
		return []byte{}, err
	}

	return respBytes, nil
}

func (c *Client) getInvokeRestAPIChecks(projectID string, resourceID string) ([]invokerestapimodel.CheckConfigurationData, error) {
	allChecksBytes, err := c.getAllChecks(projectID, resourceID)
	if err != nil {
		return []invokerestapimodel.CheckConfigurationData{}, err
	}

	result := invokerestapimodel.HierarchyResp{}
	err = json.Unmarshal(allChecksBytes, &result)
	if err != nil {
		return []invokerestapimodel.CheckConfigurationData{}, err
	}

	return result.DataProviders.MsVssPipelinechecksChecksDataProvider.CheckConfigurationDataList, nil
}

func (c *Client) getManualApprovalChecks(projectID string, resourceID string) ([]manualapprovalmodel.ManualApprovalCheckConfig, error) {
	allChecksBytes, err := c.getAllChecks(projectID, resourceID)
	if err != nil {
		return []manualapprovalmodel.ManualApprovalCheckConfig{}, err
	}

	result := manualapprovalmodel.HeirarchyResp{}
	err = json.Unmarshal(allChecksBytes, &result)
	if err != nil {
		return []manualapprovalmodel.ManualApprovalCheckConfig{}, err
	}

	configs := []manualapprovalmodel.ManualApprovalCheckConfig{}

	for _, v := range result.DataProviders.MsVssPipelinechecksChecksDataProvider.CheckConfigurationDataList {
		configs = append(configs, v.CheckConfiguration)
	}

	return configs, nil
}

func (c *Client) getExclusiveLockChecks(projectID string, resourceID string) ([]exclusivelockmodel.ExclusiveLockCheckConfig, error) {
	allChecksBytes, err := c.getAllChecks(projectID, resourceID)
	if err != nil {
		return []exclusivelockmodel.ExclusiveLockCheckConfig{}, err
	}

	result := exclusivelockmodel.HeirarchyResp{}
	err = json.Unmarshal(allChecksBytes, &result)
	if err != nil {
		return []exclusivelockmodel.ExclusiveLockCheckConfig{}, err
	}

	configs := []exclusivelockmodel.ExclusiveLockCheckConfig{}

	for _, v := range result.DataProviders.MsVssPipelinechecksChecksDataProvider.CheckConfigurationDataList {
		configs = append(configs, v.CheckConfiguration)
	}

	return configs, nil
}

func (c *Client) AddInvokeRestAPICheck(projectID string, resourceID string, check invokerestapimodel.InvokeRESTAPIValues) (invokerestapimodel.CheckConfiguration, error) {
	restAPIPayload := populateInvokeRestAPIPayload(resourceID, check)

	jsonPayload, err := json.Marshal(restAPIPayload)
	if err != nil {
		return invokerestapimodel.CheckConfiguration{}, err
	}

	url := fmt.Sprintf("/%s/_apis/pipelines/checks/configurations", projectID)
	respBytes, err := c.SendRequest("POST", url, string(jsonPayload))
	if err != nil {
		return invokerestapimodel.CheckConfiguration{}, err
	}

	checkConf := invokerestapimodel.CheckConfiguration{}

	err = json.Unmarshal(respBytes, &checkConf)
	if err != nil {
		return invokerestapimodel.CheckConfiguration{}, err
	}

	return checkConf, nil
}

func (c *Client) AddManualApprovalCheck(projectID string, resourceID string,
	check manualapprovalmodel.ManualApprovalValues) (manualapprovalmodel.ManualApprovalCheckConfig, error) {
	manualApproval := populateManualApprovalPayload(resourceID, check)

	jsonPayload, err := json.Marshal(manualApproval)
	if err != nil {
		return manualapprovalmodel.ManualApprovalCheckConfig{}, err
	}

	url := fmt.Sprintf("/%s/_apis/pipelines/checks/configurations", projectID)
	respBytes, err := c.SendRequest("POST", url, string(jsonPayload))
	if err != nil {
		return manualapprovalmodel.ManualApprovalCheckConfig{}, err
	}

	checkConf := manualapprovalmodel.ManualApprovalCheckConfig{}

	err = json.Unmarshal(respBytes, &checkConf)
	if err != nil {
		return manualapprovalmodel.ManualApprovalCheckConfig{}, err
	}

	return checkConf, nil
}

func (c *Client) AddExclusiveLockCheck(projectID string, resourceID string,
	check exclusivelockmodel.ExclusiveLockValues) (exclusivelockmodel.ExclusiveLockCheckConfig, error) {
	exclusiveLock := populateExclusiveLockPayload(resourceID, check)

	jsonPayload, err := json.Marshal(exclusiveLock)
	if err != nil {
		return exclusivelockmodel.ExclusiveLockCheckConfig{}, err
	}

	url := fmt.Sprintf("/%s/_apis/pipelines/checks/configurations", projectID)
	respBytes, err := c.SendRequest("POST", url, string(jsonPayload))
	if err != nil {
		return exclusivelockmodel.ExclusiveLockCheckConfig{}, err
	}

	checkConf := exclusivelockmodel.ExclusiveLockCheckConfig{}

	err = json.Unmarshal(respBytes, &checkConf)
	if err != nil {
		return exclusivelockmodel.ExclusiveLockCheckConfig{}, err
	}

	return checkConf, nil
}

func (c *Client) UpdateManualApprovalCheck(projectID string, resourceID string, checkID string,
	check manualapprovalmodel.ManualApprovalValues) (manualapprovalmodel.ManualApprovalCheckConfig, error) {
	manualApproval := populateManualApprovalPayload(resourceID, check)
	manualApproval.ID = checkID

	jsonPayload, err := json.Marshal(manualApproval)
	if err != nil {
		return manualapprovalmodel.ManualApprovalCheckConfig{}, err
	}

	url := fmt.Sprintf("/%s/_apis/pipelines/checks/configurations/%s", projectID, checkID)
	respBytes, err := c.SendRequest("PATCH", url, string(jsonPayload))
	if err != nil {
		return manualapprovalmodel.ManualApprovalCheckConfig{}, err
	}

	checkConf := manualapprovalmodel.ManualApprovalCheckConfig{}

	err = json.Unmarshal(respBytes, &checkConf)
	if err != nil {
		return manualapprovalmodel.ManualApprovalCheckConfig{}, err
	}

	return checkConf, nil
}

func (c *Client) UpdateExclusiveLockCheck(projectID string, resourceID string, checkID string,
	check exclusivelockmodel.ExclusiveLockValues) (exclusivelockmodel.ExclusiveLockCheckConfig, error) {
	exclusiveLock := populateExclusiveLockPayload(resourceID, check)
	exclusiveLock.ID = checkID

	jsonPayload, err := json.Marshal(exclusiveLock)
	if err != nil {
		return exclusivelockmodel.ExclusiveLockCheckConfig{}, err
	}

	url := fmt.Sprintf("/%s/_apis/pipelines/checks/configurations/%s", projectID, checkID)
	respBytes, err := c.SendRequest("PATCH", url, string(jsonPayload))
	if err != nil {
		return exclusivelockmodel.ExclusiveLockCheckConfig{}, err
	}

	checkConf := exclusivelockmodel.ExclusiveLockCheckConfig{}

	err = json.Unmarshal(respBytes, &checkConf)
	if err != nil {
		return exclusivelockmodel.ExclusiveLockCheckConfig{}, err
	}

	return checkConf, nil
}

func (c *Client) UpdateCheck(projectID string, resourceID string, checkID string, check invokerestapimodel.InvokeRESTAPIValues) (invokerestapimodel.CheckConfiguration, error) {
	restAPIPayload := populateInvokeRestAPIPayload(resourceID, check)
	restAPIPayload.ID = checkID

	jsonPayload, err := json.Marshal(restAPIPayload)
	if err != nil {
		return invokerestapimodel.CheckConfiguration{}, err
	}

	url := fmt.Sprintf("/%s/_apis/pipelines/checks/configurations/%s", projectID, checkID)
	respBytes, err := c.SendRequest("PATCH", url, string(jsonPayload))
	if err != nil {
		return invokerestapimodel.CheckConfiguration{}, err
	}

	checkConf := invokerestapimodel.CheckConfiguration{}

	err = json.Unmarshal(respBytes, &checkConf)
	if err != nil {
		return invokerestapimodel.CheckConfiguration{}, err
	}

	return checkConf, nil
}

func (c *Client) DeleteCheck(projectID string, checkID string) error {
	url := fmt.Sprintf("/%s/_apis/pipelines/checks/configurations/%s", projectID, checkID)
	_, err := c.SendRequest("DELETE", url, "")

	return err
}

func populateInvokeRestAPIPayload(resourceID string, check invokerestapimodel.InvokeRESTAPIValues) invokerestapimodel.InvokeRestAPICheckPayload {
	checkPayload := invokerestapimodel.NewInvokeRestCheckPayload()

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

func populateManualApprovalPayload(resourceID string,
	check manualapprovalmodel.ManualApprovalValues) manualapprovalmodel.ManualApprovalCheckPayload {
	approval := manualapprovalmodel.NewManualApprovalCheckPayload()
	approval.Resource.ID = resourceID
	approval.Settings.Instructions = check.Instructions
	//Allow self approve is logical opposite of cannot request and approve.
	//UI asks for self approve, API asks for RequesterCannotBeApprover
	approval.Settings.RequesterCannotBeApprover = !check.AllowSelfApproval
	approval.Timeout = check.Timeout

	approvers := []manualapprovalmodel.Approver{}

	for _, v := range check.Approvers {
		approver := manualapprovalmodel.Approver{
			ID: v,
		}
		approvers = append(approvers, approver)
	}

	approval.Settings.Approvers = approvers

	if check.ApproveInOrder {
		approval.Settings.ExecutionOrder = 2
	} else {
		approval.Settings.ExecutionOrder = 1
	}

	approval.Settings.MinRequiredApprovers = check.MinimumApprovers

	return approval
}

func populateExclusiveLockPayload(resourceID string,
	check exclusivelockmodel.ExclusiveLockValues) exclusivelockmodel.ExclusiveLockCheckPayload {
	exclusiveLock := exclusivelockmodel.NewExclusiveLockCheckPayload()
	exclusiveLock.Resource.ID = resourceID
	exclusiveLock.Timeout = check.Timeout

	return exclusiveLock
}

func (c *Client) SendRequest(httpMethod string, url string, jsonPayload string) ([]byte, error) {
	req, err := http.NewRequest(httpMethod,
		c.baseUrl+url,
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

func NewClient(baseUrl string, auth string, timeout *time.Duration) *Client {
	defaultTime := time.Duration(60 * time.Second)
	if timeout == nil {
		timeout = &defaultTime
	}

	client := &http.Client{
		Timeout: *timeout,
	}

	return &Client{
		baseUrl:       baseUrl,
		client:        client,
		authorization: auth,
	}
}

type ManualApprovalClient interface {
	GetManualApprovalCheckByID(projectID string, resourceID string, checkID int64) (manualapprovalmodel.ManualApprovalCheckConfig, bool, error)
	AddManualApprovalCheck(projectID string, resourceID string, check manualapprovalmodel.ManualApprovalValues) (manualapprovalmodel.ManualApprovalCheckConfig, error)
	UpdateManualApprovalCheck(projectID string, resourceID string, checkID string, check manualapprovalmodel.ManualApprovalValues) (manualapprovalmodel.ManualApprovalCheckConfig, error)
	DeleteCheck(projectID string, checkID string) error
}

type ExclusiveLockClient interface {
	GetExclusiveLockCheckByID(projectID string, resourceID string, checkID int64) (exclusivelockmodel.ExclusiveLockCheckConfig, bool, error)
	AddExclusiveLockCheck(projectID string, resourceID string, check exclusivelockmodel.ExclusiveLockValues) (exclusivelockmodel.ExclusiveLockCheckConfig, error)
	UpdateExclusiveLockCheck(projectID string, resourceID string, checkID string, check exclusivelockmodel.ExclusiveLockValues) (exclusivelockmodel.ExclusiveLockCheckConfig, error)
	DeleteCheck(projectID string, checkID string) error
}

type ChecksClient interface {
	GetInvokeRestAPICheckByID(projectID string, resourceID string, checkID int64) (invokerestapimodel.CheckConfigurationData, bool, error)
	AddInvokeRestAPICheck(projectID string, resourceID string, check invokerestapimodel.InvokeRESTAPIValues) (invokerestapimodel.CheckConfiguration, error)
	UpdateCheck(projectID string, resourceID string, checkID string, check invokerestapimodel.InvokeRESTAPIValues) (invokerestapimodel.CheckConfiguration, error)
	DeleteCheck(projectID string, checkID string) error
}
