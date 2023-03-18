package logicmonitor

import (
	log "github.com/sourcegraph-ce/logrus"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	lmclient "github.com/logicmonitor/lm-sdk-go/client"
	"github.com/logicmonitor/lm-sdk-go/client/lm"
)

func TestAccLogicMonitorDeviceGroup(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLogicMonitorDeviceGroup,
				Check: resource.ComposeTestCheckFunc(
					testDeviceGroupExists("logicmonitor_device_group.group1"),
					resource.TestCheckResourceAttr(
						"logicmonitor_device_group.group1", "name", "BlackMamba"),
					resource.TestCheckResourceAttr(
						"logicmonitor_device_group.group1", "disable_alerting", "true"),
					resource.TestCheckResourceAttr(
						"logicmonitor_device_group.group1", "properties.system.categories", "GOAT,MVP"),
				),
			},
		},
	})
}

func testDeviceGroupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*lmclient.LMSdkGo)
		if err := testDeviceGroupExistsHelper(s, client); err != nil {
			return err
		}
		return nil
	}
}

func testDeviceGroupExistsHelper(s *terraform.State, client *lmclient.LMSdkGo) error {
	for _, r := range s.RootModule().Resources {
		id, e := strconv.Atoi(r.Primary.ID)
		if e != nil {
			return e
		}
		params := lm.NewGetDeviceGroupByIDParams()
		params.SetID(int32(id))

		restDeviceGroupResponse, err := client.LM.GetDeviceGroupByID(params)
		if err != nil {
			return err
		}
		log.Printf("delete collector response %v", restDeviceGroupResponse)
	}
	return nil
}

const testAccCheckLogicMonitorDeviceGroup = `
resource "logicmonitor_device_group" "group1" {
    name = "BlackMamba"
    disable_alerting = true
    description = "testing group"
    applies_to = "system.displayname =~ \"test\""
    properties = {
     "system.categories" = "GOAT,MVP"
    }
}
`
