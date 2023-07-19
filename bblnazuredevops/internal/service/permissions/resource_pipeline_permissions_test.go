package permissions

// The tests in this file use the mock clients in mock_client.go to mock out
// the Azure DevOps client operations.

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"testing"
)

/**
 * Begin unit tests
 */

func TestBuildPermissions_CreateBuildTokenBH(t *testing.T) {
	projectID := "9083e944-8e9e-405e-960a-c80180aa71e6"
	buildID := "29"
	expectedToken := fmt.Sprintf("%s/%s", projectID, buildID)

	d := getBuildPermissionsResource(t, projectID, buildID, false)
	gotToken, err := createBuildTokenBH(d, nil)
	assert.NotEmpty(t, gotToken)
	assert.Nil(t, err)
	assert.Equal(t, expectedToken, gotToken)

	expectedErr := fmt.Errorf("failed to get 'project_id' from schema")
	d = getBuildPermissionsResource(t, "", "", false)
	token, err := createBuildTokenBH(d, nil)
	assert.Empty(t, token)
	assert.Equal(t, expectedErr, err)

	expectedToken = fmt.Sprintf("%s", projectID)
	d = getBuildPermissionsResource(t, projectID, "", true)
	err = d.Set("build_id", nil)
	token, err = createBuildTokenBH(d, nil)
	assert.NotEmpty(t, gotToken)
	assert.Nil(t, err)
	assert.Equal(t, expectedToken, token)

	expectedErr = fmt.Errorf("build_id cannot be set when project_level is true")
	d = getBuildPermissionsResource(t, projectID, "1234", true)
	token, err = createBuildTokenBH(d, nil)
	assert.Empty(t, token)
	assert.Equal(t, expectedErr, err)

	expectedErr = fmt.Errorf("build_id required when project_level is not true")
	d = getBuildPermissionsResource(t, projectID, "", false)
	token, err = createBuildTokenBH(d, nil)
	assert.Empty(t, token)
	assert.Equal(t, expectedErr, err)

	expectedErr = fmt.Errorf("build_id required when project_level is not true")
	d = getBuildPermissionsResource(t, projectID, "1234", false)
	err = d.Set("build_id", nil)
	token, err = createBuildTokenBH(d, nil)
	assert.Empty(t, token)
	assert.Equal(t, expectedErr, err)
}

func getBuildPermissionsResource(t *testing.T, projectID string, buildID string, projectLevel bool) *schema.ResourceData {
	d := schema.TestResourceDataRaw(t, ResourcePipelinePermissions().Schema, nil)
	if projectID != "" {
		d.Set("project_id", projectID)
	}
	if buildID != "" {
		d.Set("build_id", buildID)
	}

	d.Set("project_level", projectLevel)

	d.Set("principal", "vssgp.Uy0xLTktMTU1MTM3NDI0NS0yNDU1NjAwOTc1LTI4NjEzNDU5NS0yODk0NDM2NzIyLTIyNTMwNDg3NzAtMS0zNzc2MDEyNDM3LTI1MjAyNTQ3OTAtMjYxOTIwMDAzOS0yNTg5OTY1NzE4")

	d.Set("permissions", map[string]string{
		"ViewBuilds":                     "allow",
		"EditBuildQuality":               "allow",
		"RetainIndefinitely":             "allow",
		"DeleteBuilds":                   "allow",
		"ManageBuildQualities":           "allow",
		"DestroyBuilds":                  "allow",
		"UpdateBuildInformation":         "allow",
		"QueueBuilds":                    "allow",
		"ManageBuildQueue":               "allow",
		"StopBuilds":                     "allow",
		"ViewBuildDefinition":            "allow",
		"EditBuildDefinition":            "allow",
		"DeleteBuildDefinition":          "allow",
		"OverrideBuildCheckInValidation": "allow",
		"AdministerBuildPermissions":     "allow",
	})

	return d
}
