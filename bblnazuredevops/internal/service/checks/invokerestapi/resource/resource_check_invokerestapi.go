package resource

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/babylonhealth/terraform-provider-bblnazuredevops/bblnazuredevops/internal/client"
	"github.com/babylonhealth/terraform-provider-bblnazuredevops/bblnazuredevops/internal/service/checks/common/resource"
	"github.com/babylonhealth/terraform-provider-bblnazuredevops/bblnazuredevops/internal/service/checks/invokerestapi/model"
	"github.com/babylonhealth/terraform-provider-bblnazuredevops/bblnazuredevops/internal/utils/tfhelper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
)

// ResourceServiceEndpointDockerRegistry schema and implementation for docker registry service endpoint resource
func ResourceCheckInvokeRestAPI() *schema.Resource {
	r := &schema.Resource{
		CreateContext: createCheck,
		ReadContext:   readCheck,
		UpdateContext: updateCheck,
		DeleteContext: resource.DeleteCheckContext,
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
	r.Schema["service_connection_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	//ToDo: type will be either "endpoint" or "queue"
	// current implementation defaults to endpoint - service endpoint
	// queue support will be added in further ticket
	r.Schema["type"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}

	r.Schema["linked_variable_group"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	r.Schema["timeout"] = &schema.Schema{
		Type:     schema.TypeInt,
		Required: false,
		Optional: true,
	}
	r.Schema["retry_interval"] = &schema.Schema{
		Type:     schema.TypeInt,
		Optional: true,
	}
	r.Schema["display_name"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}

	r.Schema["method"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	r.Schema["use_callback"] = &schema.Schema{
		Type:     schema.TypeBool,
		Required: true,
	}
	r.Schema["body"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	r.Schema["url_suffix"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}

	r.Schema["success_criteria"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	r.Schema["headers"] = &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
	}

	r.Importer = tfhelper.ImportProjectQualifiedResourceUUID()

	return r
}

// See Resource documentation.
func createCheck(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	resourceID := d.Get("resource_id").(string)

	check := buildInvokeRESTAPIValuesFromSchema(d)

	resp, err := clients.InvokeCheckClient.AddInvokeRestAPICheck(ctx, projectID, resourceID, check)
	if err != nil {
		return diag.FromErr(err)
	}

	id := resp.ID

	d.SetId(fmt.Sprintf("%v", id))

	return nil
}

// See Resource documentation.
func readCheck(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	projectId := d.Get("project_id").(string)
	resourceId := d.Get("resource_id").(string)
	id := d.Id()

	idInt, err := strconv.ParseInt(id, 10, 0)
	if err != nil {
		return diag.FromErr(err)
	}

	checkConfig, found, err := clients.InvokeCheckClient.GetInvokeRestAPICheckByID(ctx, projectId, resourceId, idInt)
	if err != nil {
		return diag.FromErr(err)
	}

	if !found {
		d.SetId("")
		return nil
	}

	useCallback, err := strconv.ParseBool(checkConfig.CheckConfiguration.Settings.Inputs.WaitForCompletion)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("timeout", checkConfig.CheckConfiguration.Timeout)
	d.Set("retry_interval", checkConfig.CheckConfiguration.Settings.RetryInterval)
	d.Set("linked_variable_group", checkConfig.CheckConfiguration.Settings.LinkedVariableGroup)

	d.Set("display_name", checkConfig.CheckConfiguration.Settings.DisplayName)
	d.Set("method", checkConfig.CheckConfiguration.Settings.Inputs.Method)
	d.Set("use_callback", useCallback)
	d.Set("body", checkConfig.CheckConfiguration.Settings.Inputs.Body)
	d.Set("url_suffix", checkConfig.CheckConfiguration.Settings.Inputs.URLSuffix)
	d.Set("success_criteria", checkConfig.CheckConfiguration.Settings.Inputs.SuccessCriteria)

	headersMap := map[string]interface{}{}
	err = json.Unmarshal([]byte(checkConfig.CheckConfiguration.Settings.Inputs.Headers), &headersMap)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("headers", headersMap)

	return diag.FromErr(err)
}

// See Resource documentation.
func updateCheck(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	resourceID := d.Get("resource_id").(string)

	check := buildInvokeRESTAPIValuesFromSchema(d)

	_, err := clients.InvokeCheckClient.UpdateCheck(ctx, projectID, resourceID, d.Id(), check)
	if err != nil {
		return diag.FromErr(err)
	}

	//update ?
	d.SetId(d.Id())

	return nil
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
