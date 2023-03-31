package model

import (
	"encoding/json"
	"github.com/microsoft/terraform-provider-azuredevops/bblnazuredevops/internal/service/checks/common/model"
	"github.com/sirupsen/logrus"
)

type CheckConfigurationData struct {
	DefinitionRefID    string                   `json:"definitionRefId"`
	CheckConfiguration ExclusiveLockCheckConfig `json:"checkConfiguration"`
}

type HeirarchyResp struct {
	DataProviders struct {
		MsVssPipelinechecksChecksDataProvider struct {
			CheckConfigurationDataList []CheckConfigurationData `json:"checkConfigurationDataList"`
		} `json:"ms.vss-pipelinechecks.checks-data-provider"`
	} `json:"dataProviders"`
}

type ExclusiveLockValues struct {
	Timeout int64
}

type ExclusiveLockCheckConfig struct {
	Settings   Settings   `json:"settings"`
	CreatedBy  CreatedBy  `json:"createdBy"`
	CreatedOn  string     `json:"createdOn"`
	ModifiedBy ModifiedBy `json:"modifiedBy"`
	ModifiedOn string     `json:"modifiedOn"`
	Timeout    int64      `json:"timeout"`
	Links      Links      `json:"_links"`
	ID         int64      `json:"id"`
	Type       Type       `json:"type"`
	URL        string     `json:"url"`
	Resource   Resource   `json:"resource"`
}

type Settings struct {
}

type CreatedBy struct {
	DisplayName string `json:"displayName"`
	ID          string `json:"id"`
	UniqueName  string `json:"uniqueName"`
	Descriptor  string `json:"descriptor"`
}
type ModifiedBy struct {
	DisplayName string `json:"displayName"`
	ID          string `json:"id"`
	UniqueName  string `json:"uniqueName"`
	Descriptor  string `json:"descriptor"`
}
type Self struct {
	Href string `json:"href"`
}
type Links struct {
	Self Self `json:"self"`
}
type Type struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type Resource struct {
	Type string `json:"type"`
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ExclusiveLockCheckPayload struct {
	Type     model.CheckPayloadType `json:"type"`
	Resource model.CheckResource    `json:"resource"`
	Timeout  int64                  `json:"timeout"`
	Settings Settings               `json:"settings"`
	ID       string                 `json:"id,omitempty"`
}

func NewExclusiveLockCheckPayload() ExclusiveLockCheckPayload {
	jsonPayload := `{
    "type": {
        "id": "2EF31AD6-BAA0-403A-8B45-2CBC9B4E5563",
        "name": "ExclusiveLock"
    },
    "resource": {
        "type": "endpoint",
        "id": ""
    },
		"settings": {},
    "timeout": 60
}`

	checkPayload := ExclusiveLockCheckPayload{}

	// should not error has payload is unchanging, caught via test
	err := json.Unmarshal([]byte(jsonPayload), &checkPayload)

	if err != nil {
		logrus.Fatal(err)
	}

	return checkPayload
}
