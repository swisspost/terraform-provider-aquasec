package aquasec

import (
	"fmt"
	"os"
	"testing"

	"github.com/aquasecurity/terraform-provider-aquasec/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var imageData = client.Image{
	Registry:   "Docker Hub",
	Repository: "elasticsearch",
	Tag:        "6.4.2",
}

func TestDataSourceAquasecImage(t *testing.T) {
	rootRef := imageDataRef("test")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: getImageDataSource(&imageData),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rootRef, "registry", imageData.Registry),
					resource.TestCheckResourceAttr(rootRef, "registry_type", "HUB"),
					resource.TestCheckResourceAttr(rootRef, "repository", imageData.Repository),
					resource.TestCheckResourceAttr(rootRef, "tag", imageData.Tag),
					resource.TestCheckResourceAttr(rootRef, "scan_status", "finished"),
					resource.TestCheckResourceAttr(rootRef, "disallowed", "false"),
					resource.TestCheckResourceAttrSet(rootRef, "created"),
					resource.TestCheckResourceAttrSet(rootRef, "scan_date"),
					resource.TestCheckResourceAttr(rootRef, "scan_error", ""),
					resource.TestCheckResourceAttrSet(rootRef, "critical_vulnerabilities"),
					resource.TestCheckResourceAttrSet(rootRef, "high_vulnerabilities"),
					resource.TestCheckResourceAttrSet(rootRef, "medium_vulnerabilities"),
					resource.TestCheckResourceAttrSet(rootRef, "low_vulnerabilities"),
					resource.TestCheckResourceAttrSet(rootRef, "negligible_vulnerabilities"),
					resource.TestCheckResourceAttrSet(rootRef, "total_vulnerabilities"),
					resource.TestCheckResourceAttr(rootRef, "author", os.Getenv("AQUA_USER")),
					resource.TestCheckResourceAttr(rootRef, "os", "centos"),
					resource.TestCheckResourceAttr(rootRef, "os_version", "7"),
					resource.TestCheckResourceAttrSet(rootRef, "docker_version"),
					resource.TestCheckResourceAttrSet(rootRef, "architecture"),
					resource.TestCheckResourceAttrSet(rootRef, "image_size"),
					resource.TestCheckResourceAttrSet(rootRef, "environment_variables.0"),
					resource.TestCheckResourceAttrSet(rootRef, "vulnerabilities.0.name"),
					resource.TestCheckResourceAttrSet(rootRef, "history.0.created"),
					resource.TestCheckResourceAttrSet(rootRef, "assurance_checks_performed.0.control"),
				),
			},
		},
	})
}

func imageDataRef(name string) string {
	return fmt.Sprintf("data.aquasec_image.%s", name)
}

func getImageDataSource(image *client.Image) string {
	return fmt.Sprintf(`
	resource "aquasec_image" "test" {
		registry = "%s"
		repository = "%s"
		tag = "%s"
	}

	data "aquasec_image" "test" {
		registry = aquasec_image.test.registry
		repository = aquasec_image.test.repository
		tag = aquasec_image.test.tag

		depends_on = [
			aquasec_image.test,
		]
	}
`, image.Registry, image.Repository, image.Tag)
}
