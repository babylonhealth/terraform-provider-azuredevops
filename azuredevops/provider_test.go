package azuredevops

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProvider_SchemaIsValid(t *testing.T) {
	type testParams struct {
		name          string
		required      bool
		defaultEnvVar string
		sensitive     bool
	}

	tests := []testParams{
		{"org_service_url", false, "AZDO_ORG_SERVICE_URL", false},
		{"personal_access_token", false, "AZDO_PERSONAL_ACCESS_TOKEN", true},
	}

	schema := Provider().Schema
	require.Equal(t, len(tests), len(schema), "There are an unexpected number of properties in the schema")

	for _, test := range tests {
		require.Contains(t, schema, test.name, "An expected property was not found in the schema")
		require.NotNil(t, schema[test.name], "A property in the schema cannot have a nil value")
		require.Equal(t, test.sensitive, schema[test.name].Sensitive, "A property in the schema has an incorrect sensitivity value")
		require.Equal(t, test.required, schema[test.name].Required, "A property in the schema has an incorrect required value")

		if test.defaultEnvVar != "" {
			expectedValue := os.Getenv(test.defaultEnvVar)

			actualValue, err := schema[test.name].DefaultFunc()
			if actualValue == nil {
				actualValue = ""
			}

			require.Nil(t, err, "An error occurred when getting the default value from the environment")
			require.Equal(t, expectedValue, actualValue, "The default value pulled from the environment has the wrong value")
		}
	}
}
