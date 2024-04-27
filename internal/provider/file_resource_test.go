/*
This file contains the acceptance tests for the filedata_file resource.
*/

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFile(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccFileConfig("file1", "[\"one\",\"two\"]"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("filedata_file.file1", "file_name", "file1"),
				),
			},
			// Update and Read testing
			{
				Config: testAccFileConfig("file1", "[\"two\",\"three\"]"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("filedata_file.file1", "lines.0", "two"),
				),
			},
		},
	})
}

func testAccFileConfig(fileName, lines string) string {
	return fmt.Sprintf(`
resource "filedata_file" "%[1]s" {
  file_name  = "%[1]s"
  lines = %[2]s
}`, fileName, lines)
}
