package system

import (
	"testing"

	"github.com/influxdata/telegraf/testutil"
	"github.com/shirou/gopsutil/disk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiskUsageStats(t *testing.T) {
	var mps MockPS
	defer mps.AssertExpectations(t)
	var acc testutil.Accumulator
	var err error

	duAll := []*disk.UsageStat{
		{
			Path:        "/",
			Fstype:      "ext4",
			Total:       128,
			Free:        23,
			Used:        100,
			InodesTotal: 1234,
			InodesFree:  234,
			InodesUsed:  1000,
		},
		{
			Path:        "/home",
			Fstype:      "ext4",
			Total:       256,
			Free:        46,
			Used:        200,
			InodesTotal: 2468,
			InodesFree:  468,
			InodesUsed:  2000,
		},
	}
	duFiltered := []*disk.UsageStat{
		{
			Path:        "/",
			Fstype:      "ext4",
			Total:       128,
			Free:        23,
			Used:        100,
			InodesTotal: 1234,
			InodesFree:  234,
			InodesUsed:  1000,
		},
	}

	mps.On("DiskUsage", []string(nil), []string(nil)).Return(duAll, nil)
	mps.On("DiskUsage", []string{"/", "/dev"}, []string(nil)).Return(duFiltered, nil)
	mps.On("DiskUsage", []string{"/", "/home"}, []string(nil)).Return(duAll, nil)

	err = (&DiskUsageStats{ps: &mps}).Gather(&acc)
	require.NoError(t, err)

	numDiskMetrics := acc.NFields()
	expectedAllDiskMetrics := 2
	assert.Equal(t, expectedAllDiskMetrics, numDiskMetrics)

	tags1 := map[string]string{
		"path":   "/",
		"fstype": "ext4",
	}
	tags2 := map[string]string{
		"path":   "/home",
		"fstype": "ext4",
	}

	fields1 := map[string]interface{}{
		"used_percent": float64(81.30081300813008),
	}
	fields2 := map[string]interface{}{
		"used_percent": float64(81.30081300813008),
	}
	acc.AssertContainsTaggedFields(t, "disku", fields1, tags1)
	acc.AssertContainsTaggedFields(t, "disku", fields2, tags2)

	// We expect 6 more DiskMetrics to show up with an explicit match on "/"
	// and /home not matching the /dev in MountPoints
	err = (&DiskUsageStats{ps: &mps, MountPoints: []string{"/", "/dev"}}).Gather(&acc)
	assert.Equal(t, expectedAllDiskMetrics+1, acc.NFields())

	// We should see all the diskpoints as MountPoints includes both
	// / and /home
	err = (&DiskUsageStats{ps: &mps, MountPoints: []string{"/", "/home"}}).Gather(&acc)
	assert.Equal(t, 2*expectedAllDiskMetrics+1, acc.NFields())
}
