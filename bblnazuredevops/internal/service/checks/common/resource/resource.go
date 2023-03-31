package resource

import (
	"fmt"
	"github.com/babylonhealth/terraform-provider-bblnazuredevops/bblnazuredevops/internal/client"
	"github.com/babylonhealth/terraform-provider-bblnazuredevops/bblnazuredevops/internal/service/checks/invokerestapi/model"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strconv"
)

// See Resource documentation.
func DeleteCheck(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	checkId := d.Id()

	return clients.InvokeCheckClient.DeleteCheck(projectID, checkId)
}

// See Resource documentation.
func ExistsCheck(d *schema.ResourceData, m interface{}) (bool, error) {
	clients := m.(*client.AggregatedClient)

	projectId := d.Get("project_id").(string)
	resourceId := d.Get("resource_id").(string)
	id := d.Id()

	idInt, err := strconv.ParseInt(id, 10, 0)
	if err != nil {
		return false, err
	}

	// any check  client will work as the ID is specified in the same part of the json response
	// the body is ignored, only if it is found or not is relevant
	_, found, err := clients.InvokeCheckClient.GetInvokeRestAPICheckByID(projectId, resourceId, idInt)

	return found, err
}

func buildInvokeRESTAPIValuesFromSchema(d *schema.ResourceData) model.InvokeRESTAPIValues {
	timeout := d.Get("timeout").(int)
	retryInterval := d.Get("retry_interval").(int)

	headersFromSchema := d.Get("headers").(map[string]interface{})

	headersMap := map[string]string{}

	for k, v := range headersFromSchema {
		s := fmt.Sprintf("%s", v)
		headersMap[k] = s
	}

	check := model.InvokeRESTAPIValues{
		ServiceConnectionId: d.Get("service_connection_id").(string),
		LinkedVariableGroup: d.Get("linked_variable_group").(string),
		Timeout:             int64(timeout),
		RetryInterval:       int64(retryInterval),
		DisplayName:         d.Get("display_name").(string),
		Method:              d.Get("method").(string),
		UseCallback:         d.Get("use_callback").(bool),
		Body:                d.Get("body").(string),
		UrlSuffix:           d.Get("url_suffix").(string),
		SuccessCriteria:     d.Get("success_criteria").(string),
		Headers:             headersMap,
	}

	return check
}
