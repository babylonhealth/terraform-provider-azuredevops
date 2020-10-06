package checks

import (
	"encoding/json"
	"log"
)

type HeirarchyResp struct {
	DataProviders struct {
		MsVssPipelinechecksChecksDataProvider struct {
			CheckConfigurationDataList []CheckConfigurationData `json:"checkConfigurationDataList"`
		} `json:"ms.vss-pipelinechecks.checks-data-provider"`
	} `json:"dataProviders"`
}

type CheckConfiguration struct {
	Settings struct {
		DisplayName   string `json:"displayName"`
		DefinitionRef struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			Version string `json:"version"`
		} `json:"definitionRef"`
		Inputs struct {
			ConnectedServiceNameSelector string `json:"connectedServiceNameSelector"`
			Method                       string `json:"method"`
			WaitForCompletion            string `json:"waitForCompletion"`
			ConnectedServiceName         string `json:"connectedServiceName"`
			Body                         string `json:"body"`
			URLSuffix                    string `json:"urlSuffix"`
			SuccessCriteria              string `json:"successCriteria"`
			Headers                      string `json:"headers"`
		} `json:"inputs"`
		RetryInterval int64 `json:"retryInterval"`
		LinkedVariableGroup string `json:"linkedVariableGroup"`
	} `json:"settings"`
	CreatedBy struct {
		DisplayName string `json:"displayName"`
		ID          string `json:"id"`
		UniqueName  string `json:"uniqueName"`
		Descriptor  string `json:"descriptor"`
	} `json:"createdBy"`
	CreatedOn  string `json:"createdOn"`
	ModifiedBy struct {
		DisplayName string `json:"displayName"`
		ID          string `json:"id"`
		UniqueName  string `json:"uniqueName"`
		Descriptor  string `json:"descriptor"`
	} `json:"modifiedBy"`
	ModifiedOn string `json:"modifiedOn"`
	Timeout    int64    `json:"timeout"`
	Links      struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"_links"`
	ID   int64 `json:"id"`
	Type struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"type"`
	URL      string `json:"url"`
	Resource struct {
		Type string `json:"type"`
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"resource"`
}

type CheckConfigurationData struct {
	DefinitionRefID    string `json:"definitionRefId"`
	CheckConfiguration CheckConfiguration `json:"checkConfiguration"`
}


type InvokeRestAPICheckPayload struct {
	Type struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"type"`
	Settings struct {
		DefinitionRef struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			Version string `json:"version"`
		} `json:"definitionRef"`
		DisplayName string `json:"displayName"`
		Inputs      struct {
			ConnectedServiceNameSelector string `json:"connectedServiceNameSelector"`
			Method                       string `json:"method"`
			WaitForCompletion            string `json:"waitForCompletion"`
			ConnectedServiceName         string `json:"connectedServiceName"`
			Body                         string `json:"body"`
			URLSuffix                    string `json:"urlSuffix"`
			SuccessCriteria              string `json:"successCriteria"`
			Headers                      string `json:"headers"`
		} `json:"inputs"`
		RetryInterval       int64         `json:"retryInterval"`
		LinkedVariableGroup interface{} `json:"linkedVariableGroup"`
	} `json:"settings"`
	Resource struct {
		Type string `json:"type"`
		ID   string `json:"id"`
	} `json:"resource"`
	Timeout int64 `json:"timeout"`
	ID string `json:"id,omitempty"`
}


func NewInvokeRestCheckPayload() InvokeRestAPICheckPayload {
	jsonPayload := `{
    "type": {
        "id": "fe1de3ee-a436-41b4-bb20-f6eb4cb879a7",
        "name": "Task Check"
    },
    "settings": {
        "definitionRef": {
            "id": "9c3e8943-130d-4c78-ac63-8af81df62dfb",
            "name": "InvokeRESTAPI",
            "version": "1.152.3"
        },
        "displayName": "", 
        "inputs": {
            "connectedServiceNameSelector": "connectedServiceName",
            "method": "", 
            "waitForCompletion": "false", 
            "connectedServiceName": "USE ID FROM SERVICE IN TERRAFORM", 
            "body": "", 
            "urlSuffix": "", 
            "successCriteria": "",
			"headers": "{\n\"Content-Type\":\"application/json\", \n\"PlanUrl\": \"$(system.CollectionUri)\", \n\"ProjectId\": \"$(system.TeamProjectId)\", \n\"HubName\": \"$(system.HostType)\", \n\"PlanId\": \"$(system.PlanId)\", \n\"JobId\": \"$(system.JobId)\", \n\"TimelineId\": \"$(system.TimelineId)\", \n\"TaskInstanceId\": \"$(system.TaskInstanceId)\", \n\"AuthToken\": \"$(system.AccessToken)\"\n}"
        },
        "retryInterval": 5, 
        "linkedVariableGroup": null 
    },
    "resource": {
        "type": "endpoint",
        "id": "%s"
    },
    "timeout": 43200
}`

	checkPayload := InvokeRestAPICheckPayload{}

	// should not error has payload is unchanging, caught via test
	err := json.Unmarshal([]byte(jsonPayload), &checkPayload)

	if err != nil {
		log.Fatal(err)
	}

	return checkPayload
}

type InvokeRESTAPIValues struct {
	ServiceConnectionId string
	LinkedVariableGroup string
	Timeout int64
	RetryInterval int64

	DisplayName string
	Method string
	UseCallback bool // True is Callback, false is ApiResponse

	Body string
	UrlSuffix string
	SuccessCriteria string
	Headers map[string]string
}
