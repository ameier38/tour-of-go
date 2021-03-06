package onepassword

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

const testDataSourceDocumentConfig = `
provider "onepassword" {
	email = "test@testing.com"
	password = "test-password"
	secret_key = "test-secret-key"
	subdomain = "test"
}

data "onepassword_document" "test" {
	vault = "test-vault"
	document = "test-doc"
}

output "test_doc" {
	value = "${data.onepassword_document.test.result}"
} 
`

func TestDataSourceDocument(t *testing.T) {
	progPath, err := buildMockOnePassword()

	if err != nil {
		t.Errorf("failed to build mock 1Password cli: %s", err)
	}

	os.Setenv("OP_PATH", progPath)

	resource.UnitTest(t, resource.TestCase{
		Providers: createTestProviders(),
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDocumentConfig,
				Check: func(s *terraform.State) error {
					_, ok := s.RootModule().Resources["data.onepassword_document.test"]
					if !ok {
						return fmt.Errorf("missing data.onepassword_document.test data source")
					}

					outputs := s.RootModule().Outputs

					if outputs["test_doc"] == nil {
						return fmt.Errorf("missing 'test_doc' output")
					}

					if outputs["test_doc"].Value != "hello world" {
						return fmt.Errorf("'%s' != 'hello world'", outputs["test_doc"].Value)
					}

					return nil
				},
			},
		},
	})

}
