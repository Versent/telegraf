package printer

import (
	"fmt"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/filters"
)

type Printer struct {
}

var sampleConfig = `
`

func (p *Printer) SampleConfig() string {
	return sampleConfig
}

func (p *Printer) Description() string {
	return "Print all metrics that pass through this filter."
}

func (p *Printer) Apply(in ...telegraf.Metric) []telegraf.Metric {
	for _, metric := range in {
		fmt.Println(metric.String())
	}
	return in
}

func init() {
	filters.Add("printer", func() telegraf.FilterPlugin {
		return &Printer{}
	})
}
