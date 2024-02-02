package config

import (
	"github.com/gdamore/tcell/v2"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/shared/types"
	"github.com/rivo/tview"
)

// The page wrapper for the metrics config
type MetricsConfigPage struct {
	home                       *settingsHome
	page                       *page
	layout                     *standardLayout
	masterConfig               *client.GlobalConfig
	enableMetricsBox           *parameterizedFormItem
	ecMetricsPortBox           *parameterizedFormItem
	bnMetricsPortBox           *parameterizedFormItem
	daemonMetricsPortBox       *parameterizedFormItem
	exporterMetricsPortBox     *parameterizedFormItem
	grafanaItems               []*parameterizedFormItem
	prometheusItems            []*parameterizedFormItem
	exporterItems              []*parameterizedFormItem
	enableBitflyNodeMetricsBox *parameterizedFormItem
	bitflyNodeMetricsItems     []*parameterizedFormItem
}

// Creates a new page for the metrics / stats settings
func NewMetricsConfigPage(home *settingsHome) *MetricsConfigPage {

	configPage := &MetricsConfigPage{
		home:         home,
		masterConfig: home.md.Config,
	}
	configPage.createContent()

	configPage.page = newPage(
		home.homePage,
		"settings-metrics",
		"Monitoring / Metrics",
		"Select this to configure the monitoring and statistics gathering parts of Hyperdrive, such as Grafana and Prometheus.",
		configPage.layout.grid,
	)

	return configPage

}

// Get the underlying page
func (configPage *MetricsConfigPage) getPage() *page {
	return configPage.page
}

// Creates the content for the monitoring / stats settings page
func (configPage *MetricsConfigPage) createContent() {

	// Create the layout
	configPage.layout = newStandardLayout()
	configPage.layout.createForm(&configPage.masterConfig.Hyperdrive.Network, "Monitoring / Metrics Settings")

	// Return to the home page after pressing Escape
	configPage.layout.form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Return to the home page
		if event.Key() == tcell.KeyEsc {
			// Close all dropdowns and break if one was open
			for _, param := range configPage.layout.parameters {
				dropDown, ok := param.item.(*DropDown)
				if ok && dropDown.open {
					dropDown.CloseList(configPage.home.md.app)
					return nil
				}
			}

			configPage.home.md.setPage(configPage.home.homePage)
			return nil
		}
		return event
	})

	// Set up the form items
	configPage.enableMetricsBox = createParameterizedCheckbox(&configPage.masterConfig.Hyperdrive.Metrics.EnableMetrics)
	configPage.ecMetricsPortBox = createParameterizedUint16Field(&configPage.masterConfig.Hyperdrive.Metrics.EcMetricsPort)
	configPage.bnMetricsPortBox = createParameterizedUint16Field(&configPage.masterConfig.Hyperdrive.Metrics.BnMetricsPort)
	configPage.daemonMetricsPortBox = createParameterizedUint16Field(&configPage.masterConfig.Hyperdrive.Metrics.DaemonMetricsPort)
	configPage.exporterMetricsPortBox = createParameterizedUint16Field(&configPage.masterConfig.Hyperdrive.Metrics.ExporterMetricsPort)
	configPage.grafanaItems = createParameterizedFormItems(configPage.masterConfig.Hyperdrive.Metrics.Grafana.GetParameters(), configPage.layout.descriptionBox)
	configPage.prometheusItems = createParameterizedFormItems(configPage.masterConfig.Hyperdrive.Metrics.Prometheus.GetParameters(), configPage.layout.descriptionBox)
	configPage.exporterItems = createParameterizedFormItems(configPage.masterConfig.Hyperdrive.Metrics.Exporter.GetParameters(), configPage.layout.descriptionBox)
	configPage.enableBitflyNodeMetricsBox = createParameterizedCheckbox(&configPage.masterConfig.Hyperdrive.Metrics.EnableBitflyNodeMetrics)
	configPage.bitflyNodeMetricsItems = createParameterizedFormItems(configPage.masterConfig.Hyperdrive.Metrics.BitflyNodeMetrics.GetParameters(), configPage.layout.descriptionBox)

	// Map the parameters to the form items in the layout
	configPage.layout.mapParameterizedFormItems(configPage.enableMetricsBox, configPage.ecMetricsPortBox, configPage.bnMetricsPortBox, configPage.daemonMetricsPortBox, configPage.exporterMetricsPortBox)
	configPage.layout.mapParameterizedFormItems(configPage.grafanaItems...)
	configPage.layout.mapParameterizedFormItems(configPage.prometheusItems...)
	configPage.layout.mapParameterizedFormItems(configPage.exporterItems...)
	configPage.layout.mapParameterizedFormItems(configPage.enableBitflyNodeMetricsBox)
	configPage.layout.mapParameterizedFormItems(configPage.bitflyNodeMetricsItems...)

	// Set up the setting callbacks
	configPage.enableMetricsBox.item.(*tview.Checkbox).SetChangedFunc(func(checked bool) {
		if configPage.masterConfig.Hyperdrive.Metrics.EnableMetrics.Value == checked {
			return
		}
		configPage.masterConfig.Hyperdrive.Metrics.EnableMetrics.Value = checked
		configPage.handleLayoutChanged()
	})
	configPage.enableBitflyNodeMetricsBox.item.(*tview.Checkbox).SetChangedFunc(func(checked bool) {
		if configPage.masterConfig.Hyperdrive.Metrics.EnableBitflyNodeMetrics.Value == checked {
			return
		}
		configPage.masterConfig.Hyperdrive.Metrics.EnableBitflyNodeMetrics.Value = checked
		configPage.handleLayoutChanged()
	})

	// Do the initial draw
	configPage.handleLayoutChanged()
}

// Handle all of the form changes when the Enable Metrics box has changed
func (configPage *MetricsConfigPage) handleLayoutChanged() {
	configPage.layout.form.Clear(true)
	configPage.layout.form.AddFormItem(configPage.enableMetricsBox.item)

	if configPage.masterConfig.Hyperdrive.Metrics.EnableMetrics.Value == true {
		configPage.layout.addFormItems([]*parameterizedFormItem{configPage.ecMetricsPortBox, configPage.bnMetricsPortBox, configPage.daemonMetricsPortBox, configPage.exporterMetricsPortBox})
		configPage.layout.addFormItems(configPage.grafanaItems)
		configPage.layout.addFormItems(configPage.prometheusItems)
		configPage.layout.addFormItems(configPage.exporterItems)
	}

	if configPage.masterConfig.Hyperdrive.ClientMode.Value == types.ClientMode_Local {
		switch configPage.masterConfig.Hyperdrive.LocalBeaconConfig.BeaconNode.Value {
		case types.BeaconNode_Teku, types.BeaconNode_Lighthouse, types.BeaconNode_Lodestar:
			configPage.layout.form.AddFormItem(configPage.enableBitflyNodeMetricsBox.item)
			if configPage.masterConfig.Hyperdrive.Metrics.EnableBitflyNodeMetrics.Value == true {
				configPage.layout.addFormItems(configPage.bitflyNodeMetricsItems)
			}
		}
	}

	configPage.layout.refresh()
}
