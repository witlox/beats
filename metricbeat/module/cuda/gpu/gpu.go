package gpu

import (
	"github.com/pkg/errors"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/metricbeat/mb"
	"github.com/elastic/beats/metricbeat/mb/parse"

	"os/exec"
	"encoding/xml"
)

func init() {
	if err := mb.Registry.AddMetricSet("cuda", "gpu", New, parse.EmptyHostParser); err != nil {
		panic(err)
	}
}

// nvidia smi binary
const (
	BinaryName = "nvidia-smi"
)

// MetricSet for fetching GPU metrics.
type MetricSet struct {
	mb.BaseMetricSet
	NvidiaSmi
}

// New is a mb.MetricSetFactory that returns a MetricSet (containing NvidiaSmi).
func New(base mb.BaseMetricSet) (mb.MetricSet, error) {
	exe, err := exec.LookPath(BinaryName)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot find %v in path", BinaryName)
	}
	cmd := exec.Command(exe, "-x", "-q", "-a")
	bts, err := cmd.CombinedOutput()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to run %v -x -q -a", exe)
	}
	res := new(NvidiaSmi)
	err = xml.Unmarshal(bts, res)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal %v -x -q -a", exe)
	}
	return &MetricSet{
		BaseMetricSet: base,
		res,
	}, nil

}

// Fetch fetches GPU metrics from the NvidiaSmi.
func (m *MetricSet) Fetch() (common.MapStr, error) {
	if m.NvidiaSmi.HasGPU() == false {
		return nil, "no nvidia gpu detected"
	}

	gpus := m.NvidiaSmi.GPUS

	event := common.MapStr{"gpus": len(gpus)}

	for i := 0; i < len(gpus); i++ {
		gpu := gpus[i]
		event.Put("gpu"+i+".type", gpu.ProductName)
		event.Put("gpu"+i+".temp", gpu.GpuTemp)
		event.Put("gpu"+i+".fan", gpu.FanSpeed)
		event.Put("gpu"+i+".pwrdrw", gpu.PowerDraw)
		event.Put("gpu"+i+".pwrlim", gpu.PowerLimit)
		event.Put("gpu"+i+".perfst", gpu.PerformanceState)
		event.Put("gpu"+i+".memtot", gpu.TotalFbMemoryUsageGpu)
		event.Put("gpu"+i+".memfree", gpu.FreeFbMemoryUsageGpu)
		event.Put("gpu"+i+".memused", gpu.UsedFbMemoryUsageGpu)
		event.Put("gpu"+i+".gpupct", gpu.GpuUtil)
		event.Put("gpu"+i+".mempct", gpu.MemoryUtil)
		event.Put("gpu"+i+".gclock", gpu.GraphicsClock)
		event.Put("gpu"+i+".mclock", gpu.MemClock)
		event.Put("gpu"+i+".smclock", gpu.SmClock)
	}

	return event, nil
}
