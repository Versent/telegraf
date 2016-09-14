package system

import (
	"fmt"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
)

type MemStats struct {
	ps        PS
	UsageOnly bool `toml:"usage_only"`
}

func (_ *MemStats) Description() string {
	return "Read metrics about memory usage"
}

var memorySampleConfig = `
  ## By default, telegraf gather stats for all mountpoints.
  ## Enable only percentage usage.
  usage_only = true
`

func (_ *MemStats) SampleConfig() string {
	return memorySampleConfig
}

func (s *MemStats) Gather(acc telegraf.Accumulator) error {
	vm, err := s.ps.VMStat()
	if err != nil {
		return fmt.Errorf("error getting virtual memory info: %s", err)
	}

	var fields map[string]interface{}

	if s.UsageOnly {
		fields = map[string]interface{}{
			"used_percent": 100 * float64(vm.Used) / float64(vm.Total),
		}
	} else {
		fields = map[string]interface{}{
			"total":             vm.Total,
			"available":         vm.Available,
			"used":              vm.Used,
			"free":              vm.Free,
			"cached":            vm.Cached,
			"buffered":          vm.Buffers,
			"active":            vm.Active,
			"inactive":          vm.Inactive,
			"used_percent":      100 * float64(vm.Used) / float64(vm.Total),
			"available_percent": 100 * float64(vm.Available) / float64(vm.Total),
		}
	}

	acc.AddCounter("mem", fields, nil)

	return nil
}

type SwapStats struct {
	ps PS
}

func (_ *SwapStats) Description() string {
	return "Read metrics about swap memory usage"
}

func (_ *SwapStats) SampleConfig() string { return "" }

func (s *SwapStats) Gather(acc telegraf.Accumulator) error {
	swap, err := s.ps.SwapStat()
	if err != nil {
		return fmt.Errorf("error getting swap memory info: %s", err)
	}

	fieldsG := map[string]interface{}{
		"total":        swap.Total,
		"used":         swap.Used,
		"free":         swap.Free,
		"used_percent": swap.UsedPercent,
	}
	fieldsC := map[string]interface{}{
		"in":  swap.Sin,
		"out": swap.Sout,
	}
	acc.AddGauge("swap", fieldsG, nil)
	acc.AddCounter("swap", fieldsC, nil)

	return nil
}

func init() {
	inputs.Add("mem", func() telegraf.Input {
		return &MemStats{ps: &systemPS{}}
	})

	inputs.Add("swap", func() telegraf.Input {
		return &SwapStats{ps: &systemPS{}}
	})
}
