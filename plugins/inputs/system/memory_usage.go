package system

import (
	"fmt"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
)

type MemUsageStats struct {
	ps PS
}

func (_ *MemUsageStats) Description() string {
	return "Read utilisation of memory metric"
}

func (_ *MemUsageStats) SampleConfig() string { return "" }

func (s *MemUsageStats) Gather(acc telegraf.Accumulator) error {
	vm, err := s.ps.VMStat()
	if err != nil {
		return fmt.Errorf("error getting virtual memory info: %s", err)
	}

	fields := map[string]interface{}{
		"used_percent": 100 * float64(vm.Used) / float64(vm.Total),
	}
	acc.AddCounter("memu", fields, nil)

	return nil
}

type SwapUsageStats struct {
	ps PS
}

func (_ *SwapUsageStats) Description() string {
	return "Read utilisation of swap memory metric"
}

func (_ *SwapUsageStats) SampleConfig() string { return "" }

func (s *SwapUsageStats) Gather(acc telegraf.Accumulator) error {
	swap, err := s.ps.SwapStat()
	if err != nil {
		return fmt.Errorf("error getting swap memory info: %s", err)
	}

	fieldsG := map[string]interface{}{
		"used_percent": swap.UsedPercent,
	}

	acc.AddGauge("swapu", fieldsG, nil)

	return nil
}

func init() {
	inputs.Add("memu", func() telegraf.Input {
		return &MemUsageStats{ps: &systemPS{}}
	})

	inputs.Add("swapu", func() telegraf.Input {
		return &SwapUsageStats{ps: &systemPS{}}
	})
}
