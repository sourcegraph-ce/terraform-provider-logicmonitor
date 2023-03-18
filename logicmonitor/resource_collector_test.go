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

func TestAccLogicMonitorCollector(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCollectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLogicMonitorConfigCollector,
				Check: resource.ComposeTestCheckFunc(
					testCollectorExists("logicmonitor_collector.collector1"),
					resource.TestCheckResourceAttr(
						"logicmonitor_collector.collector1", "description", "test collector"),
					resource.TestCheckResourceAttr(
						"logicmonitor_collector.collector1", "enable_collector_device_failover", "false"),
					resource.TestCheckResourceAttr(
						"logicmonitor_collector.collector1", "enable_failback", "true"),
					resource.TestCheckResourceAttr(
						"logicmonitor_collector.collector1", "resend_interval", "5"),
					resource.TestCheckResourceAttr(
						"logicmonitor_collector.collector1", "suppress_alert_clear", "false"),
				),
			},
		},
	})
}

func testCollectorDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*lmclient.LMSdkGo)
	if err := testCollectorDestroyHelper(s, client); err != nil {
		return err
	}
	return nil
}

func testCollectorDestroyHelper(s *terraform.State, client *lmclient.LMSdkGo) error {
	for _, r := range s.RootModule().Resources {
		id, e := strconv.Atoi(r.Primary.ID)
		if e != nil {
			return e
		}
		params := lm.NewDeleteCollectorByIDParams()
		params.SetID(int32(id))

		restCollectorResponse, err := client.LM.DeleteCollectorByID(params)
		if err != nil {
			return err
		}
		log.Printf("delete collector response %v", restCollectorResponse.Payload)
	}
	return nil
}

func testCollectorExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*lmclient.LMSdkGo)
		if err := testCollectorExistsHelper(s, client); err != nil {
			return err
		}
		return nil
	}
}

func testCollectorExistsHelper(s *terraform.State, client *lmclient.LMSdkGo) error {
	collectorID := s.RootModule().Resources["logicmonitor_collector.collector1"]
	id, e := strconv.Atoi(collectorID.Primary.ID)
	if e != nil {
		return e
	}

	params := lm.NewGetCollectorByIDParams()
	params.SetID(int32(id))

	restCollectorResponse, err := client.LM.GetCollectorByID(params)
	if err != nil {
		return err
	}
	log.Printf("get collector id response %v", restCollectorResponse.Payload)
	return nil
}

const testAccCheckLogicMonitorConfigCollector = `
resource "logicmonitor_collector" "collector1" {
    description                       = "test collector"
    enable_collector_device_failover  = false
    enable_failback                   = true
    resend_interval                   = 5
    suppress_alert_clear              = false
}
`
