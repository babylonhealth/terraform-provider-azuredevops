module github.com/babylonhealth/terraform-provider-bblnazuredevops

go 1.16

require (
	//	patch versus microsoft-origin above
	github.com/ahmetb/go-linq v3.0.0+incompatible
	//	patch versus microsoft-origin below
	github.com/go-test/deep v1.0.3
	github.com/golang/mock v1.4.1
	github.com/google/go-cmp v0.5.2
	github.com/google/uuid v1.1.1
	github.com/hashicorp/go-multierror v1.0.0
	github.com/hashicorp/go-uuid v1.0.1
	github.com/hashicorp/terraform v0.12.23
	github.com/hashicorp/terraform-plugin-sdk v1.13.1
	github.com/microsoft/azure-devops-go-api/azuredevops v1.0.0-b3
	github.com/sirupsen/logrus v1.2.0
	github.com/stretchr/testify v1.4.0
	golang.org/x/crypto v0.0.0-20200427165652-729f1e841bcc
	gopkg.in/yaml.v2 v2.2.4
)
