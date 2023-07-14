package resource

import (
	"context"
	"fmt"
	"github.com/babylonhealth/terraform-provider-bblnazuredevops/bblnazuredevops/internal/client"
	"github.com/babylonhealth/terraform-provider-bblnazuredevops/bblnazuredevops/internal/service/checks/invokerestapi/model"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DeleteCheckContext(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	checkId := d.Id()

	return diag.FromErr(clients.InvokeCheckClient.DeleteCheck(ctx, projectID, checkId))
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
