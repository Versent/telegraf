package systemu

import (
	"fmt"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
)

type MemStats struct {
	ps PS
}

func (_ *MemStats) Description() string {
	return "Read utilisation of memory metric"
}

func (_ *MemStats) SampleConfig() string { return "" }

func (s *MemStats) Gather(acc telegraf.Accumulator) error {
	vm, err := s.ps.VMStat()
	if err != nil {
		return fmt.Errorf("error getting virtual memory info: %s", err)
	}

	fields := map[string]interface{}{
		"used_percent": 100 * float64(vm.Used) / float64(vm.Total),
	}
	acc.AddCounter("mem", fields, nil)

	return nil
}

type SwapStats struct {
	ps PS
}

func (_ *SwapStats) Description() string {
	return "Read utilisation of swap memory metric"
}

func (_ *SwapStats) SampleConfig() string { return "" }

func (s *SwapStats) Gather(acc telegraf.Accumulator) error {
	swap, err := s.ps.SwapStat()
	if err != nil {
		return fmt.Errorf("error getting swap memory info: %s", err)
	}

	fieldsG := map[string]interface{}{
		"used_percent": swap.UsedPercent,
	}

	acc.AddGauge("swap", fieldsG, nil)

	return nil
}

func init() {
	inputs.Add("memu", func() telegraf.Input {
		return &MemStats{ps: &systemPS{}}
	})

	inputs.Add("swapu", func() telegraf.Input {
		return &SwapStats{ps: &systemPS{}}
	})
}
