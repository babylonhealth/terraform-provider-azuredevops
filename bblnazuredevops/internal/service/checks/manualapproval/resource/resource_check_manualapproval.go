package resource

import (
	"fmt"
	"github.com/babylonhealth/terraform-provider-bblnazuredevops/bblnazuredevops/internal/client"
	"github.com/babylonhealth/terraform-provider-bblnazuredevops/bblnazuredevops/internal/service/checks/common/resource"
	"github.com/babylonhealth/terraform-provider-bblnazuredevops/bblnazuredevops/internal/service/checks/manualapproval/model"
	"github.com/babylonhealth/terraform-provider-bblnazuredevops/bblnazuredevops/internal/utils/tfhelper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"strconv"
)

// ResourceServiceEndpointDockerRegistry schema and implementation for docker registry service endpoint resource
func ResourceCheckManualApproval() *schema.Resource {
	r := &schema.Resource{
		Create: createCheck,
		Read:   readCheck,
		Update: updateCheck,
		Delete: resource.DeleteCheck,
		Exists: resource.ExistsCheck,
	}
	r.Schema = map[string]*schema.Schema{}
	r.Schema["project_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	}
	r.Schema["resource_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	}

	//ToDo: type will be either "endpoint" or "queue"
	// current implementation defaults to endpoint - service endpoint
	// queue support will be added in further ticket
	r.Schema["type"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}

	r.Schema["timeout"] = &schema.Schema{
		Type:     schema.TypeInt,
		Required: false,
		Optional: true,
	}

	r.Schema["approvers"] = &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		Elem: &schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validation.NoZeroValues,
		},
	}
	r.Schema["allow_self_approve"] = &schema.Schema{
		Type:     schema.TypeBool,
		Required: true,
	}
	r.Schema["approve_in_order"] = &schema.Schema{
		Type:     schema.TypeBool,
		Required: false,
		Optional: true,
		Default:  false,
	}
	r.Schema["minimum_approvers"] = &schema.Schema{
		Type:     schema.TypeInt,
		Required: false,
		Optional: true,
		Default:  0,
	}

	r.Schema["instructions"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}

	r.Importer = tfhelper.ImportProjectQualifiedResourceUUID()

	return r
}

// See Resource documentation.
func createCheck(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	resourceID := d.Get("resource_id").(string)

	check := buildManualApprovalValuesFromSchema(d)

	resp, err := clients.ManualApprovalCheckClient.AddManualApprovalCheck(projectID, resourceID, check)
	if err != nil {
		return fmt.Errorf("error creating check in Azure DevOps: %+v", err)
	}

	id := resp.ID

	d.SetId(fmt.Sprintf("%v", id))

	return nil
}

// See Resource documentation.
func readCheck(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	resourceID := d.Get("resource_id").(string)

	checkId := d.Id()

	idInt, err := strconv.ParseInt(checkId, 10, 0)
	if err != nil {
		return err
	}

	checkConfig, found, err := clients.ManualApprovalCheckClient.GetManualApprovalCheckByID(projectID, resourceID, idInt)
	if err != nil {
		return fmt.Errorf("error reading check in Azure DevOps: %+v", err)
	}

	if !found {
		d.SetId("")
		return nil
	}

	d.Set("timeout", checkConfig.Timeout)
	d.Set("allow_self_approve", !checkConfig.Settings.RequesterCannotBeApprover)
	d.Set("instructions", checkConfig.Settings.Instructions)

	approvers := []string{}

	for _, approver := range checkConfig.Settings.Approvers {
		approvers = append(approvers, approver.ID)
	}

	d.Set("approvers", approvers)

	d.Set("minimum_approvers", checkConfig.Settings.MinRequiredApprovers)

	approveInOrder := false
	if checkConfig.Settings.ExecutionOrder == 2 {
		approveInOrder = true
	}

	d.Set("approve_in_order", approveInOrder)

	return err

}

// See Resource documentation.
func updateCheck(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	resourceID := d.Get("resource_id").(string)

	check := buildManualApprovalValuesFromSchema(d)

	_, err := clients.ManualApprovalCheckClient.UpdateManualApprovalCheck(projectID, resourceID, d.Id(), check)
	if err != nil {
		return fmt.Errorf("error creating check in Azure DevOps: %+v", err)
	}

	//update ?
	d.SetId(d.Id())

	return nil
}

func buildManualApprovalValuesFromSchema(d *schema.ResourceData) model.ManualApprovalValues {
	timeout := d.Get("timeout").(int)
	minimumApprovers := d.Get("minimum_approvers").(int)

	approversFromSchema := d.Get("approvers").([]interface{})

	approversList := []string{}

	for _, v := range approversFromSchema {
		s := fmt.Sprintf("%s", v)
		approversList = append(approversList, s)
	}

	check := model.ManualApprovalValues{
		Approvers:         approversList,
		Timeout:           int64(timeout),
		AllowSelfApproval: d.Get("allow_self_approve").(bool),
		Instructions:      d.Get("instructions").(string),
		MinimumApprovers:  int64(minimumApprovers),
		ApproveInOrder:    d.Get("approve_in_order").(bool),
	}

	return check
}
