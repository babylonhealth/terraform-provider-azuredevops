package model

import (
	"encoding/json"
	"github.com/microsoft/terraform-provider-azuredevops/bblnazuredevops/internal/service/checks/common/model"
	"github.com/sirupsen/logrus"
)

type CheckConfigurationData struct {
	DefinitionRefID    string                    `json:"definitionRefId"`
	CheckConfiguration ManualApprovalCheckConfig `json:"checkConfiguration"`
}

type HeirarchyResp struct {
	DataProviders struct {
		MsVssPipelinechecksChecksDataProvider struct {
			CheckConfigurationDataList []CheckConfigurationData `json:"checkConfigurationDataList"`
		} `json:"ms.vss-pipelinechecks.checks-data-provider"`
	} `json:"dataProviders"`
}

type ManualApprovalValues struct {
	Approvers         []string
	Instructions      string
	AllowSelfApproval bool
	Timeout           int64
	ApproveInOrder    bool
	MinimumApprovers  int64
}

type ManualApprovalCheckConfig struct {
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
type Approvers struct {
	DisplayName interface{} `json:"displayName"`
	ID          string      `json:"id"`
}
type Settings struct {
	RequesterCannotBeApprover bool          `json:"requesterCannotBeApprover"`
	Approvers                 []Approvers   `json:"approvers"`
	ExecutionOrder            int64         `json:"executionOrder"`
	MinRequiredApprovers      int64         `json:"minRequiredApprovers"`
	Instructions              string        `json:"instructions"`
	BlockedApprovers          []interface{} `json:"blockedApprovers"`
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

type Approver struct {
	ID string `json:"id"`
}

type ManualApprovalCheckPayload struct {
	Type     model.CheckPayloadType `json:"type"`
	Settings struct {
		Approvers                 []Approver `json:"approvers"`
		ExecutionOrder            int64      `json:"executionOrder"`
		Instructions              string     `json:"instructions"`
		BlockedApprovers          []string   `json:"blockedApprovers"`
		MinRequiredApprovers      int64      `json:"minRequiredApprovers"`
		RequesterCannotBeApprover bool       `json:"requesterCannotBeApprover"`
	} `json:"settings"`
	Resource model.CheckResource `json:"resource"`
	Timeout  int64               `json:"timeout"`
	ID       string              `json:"id,omitempty"`
}

func NewManualApprovalCheckPayload() ManualApprovalCheckPayload {
	jsonPayload := `{
    "type": {
        "id": "8C6F20A7-A545-4486-9777-F762FAFE0D4D",
        "name": "Approval"
    },
    "settings": {
        "approvers": [],
        "executionOrder": 1,
        "instructions": "",
        "blockedApprovers": [],
        "minRequiredApprovers": 0,
        "requesterCannotBeApprover": false
    },
    "resource": {
        "type": "endpoint",
        "id": ""
    },
    "timeout": 43200
}`

	checkPayload := ManualApprovalCheckPayload{}

	// should not error has payload is unchanging, caught via test
	err := json.Unmarshal([]byte(jsonPayload), &checkPayload)

	if err != nil {
		logrus.Fatal(err)
	}

	return checkPayload
}
