package resource

import (
	"context"
	"fmt"
	"github.com/babylonhealth/terraform-provider-bblnazuredevops/bblnazuredevops/internal/client"
	"github.com/babylonhealth/terraform-provider-bblnazuredevops/bblnazuredevops/internal/service/checks/common/resource"
	"github.com/babylonhealth/terraform-provider-bblnazuredevops/bblnazuredevops/internal/service/checks/exclusivelock/model"
	"github.com/babylonhealth/terraform-provider-bblnazuredevops/bblnazuredevops/internal/utils/tfhelper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
)

// ResourceServiceEndpointDockerRegistry schema and implementation for docker registry service endpoint resource
func ResourceCheckExclusiveLock() *schema.Resource {
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

	//ToDo: type will be either "endpoint" or "variablegroup"
	// current implementation defaults to endpoint - service endpoint
	// variablegroup support will be added in further ticket
	r.Schema["type"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}

	r.Schema["timeout"] = &schema.Schema{
		Type:     schema.TypeInt,
		Required: false,
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

	check := buildExclusiveLockValuesFromSchema(d)

	resp, err := clients.ExclusiveLockCheckClient.AddExclusiveLockCheck(ctx, projectID, resourceID, check)
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

	projectID := d.Get("project_id").(string)
	resourceID := d.Get("resource_id").(string)

	checkId := d.Id()

	idInt, err := strconv.ParseInt(checkId, 10, 0)
	if err != nil {
		return diag.FromErr(err)
	}

	checkConfig, found, err := clients.ExclusiveLockCheckClient.GetExclusiveLockCheckByID(ctx, projectID, resourceID, idInt)
	if err != nil {
		return diag.FromErr(err)
	}

	if !found {
		d.SetId("")
		return nil
	}

	d.Set("timeout", checkConfig.Timeout)

	return diag.FromErr(err)
}

// See Resource documentation.
func updateCheck(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	resourceID := d.Get("resource_id").(string)

	check := buildExclusiveLockValuesFromSchema(d)

	_, err := clients.ExclusiveLockCheckClient.UpdateExclusiveLockCheck(ctx, projectID, resourceID, d.Id(), check)
	if err != nil {
		return diag.FromErr(err)
	}

	//update ?
	d.SetId(d.Id())

	return nil
}

func buildExclusiveLockValuesFromSchema(d *schema.ResourceData) model.ExclusiveLockValues {
	timeout := d.Get("timeout").(int)
	check := model.ExclusiveLockValues{
		Timeout: int64(timeout),
	}

	return check
}
