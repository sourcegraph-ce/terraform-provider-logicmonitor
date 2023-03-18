package logicmonitor

import (
	"fmt"
	log "github.com/sourcegraph-ce/logrus"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

//Provider for LogicMonitor
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("LM_API_ID", nil),
			},
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("LM_API_KEY", nil),
			},
			"company": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("LM_COMPANY", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"logicmonitor_collector":       resourceCollector(),
			"logicmonitor_collector_group": resourceCollectorGroup(),
			"logicmonitor_dashboard":       resourceDashboard(),
			"logicmonitor_dashboard_group": resourceDashboardGroup(),
			"logicmonitor_device":          resourceDevices(),
			"logicmonitor_device_group":    resourceDeviceGroup(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"logicmonitor_collectors":      dataSourceFindCollectors(),
			"logicmonitor_dashboard":       dataSourceFindDashboards(),
			"logicmonitor_dashboard_group": dataSourceFindDashboardGroups(),
			"logicmonitor_device_group":    dataSourceFindDeviceGroups(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	company := fmt.Sprintf("%s.logicmonitor.com", strings.Replace(d.Get("company").(string), ".logicmonitor.com", "", -1))
	config := Config{
		AccessID:  d.Get("api_id").(string),
		AccessKey: d.Get("api_key").(string),
		Company:   company,
	}
	log.Println("[INFO] Initializing LM client")
	return config.newLMClient()
}
