// +build linux

/*
http://www.apache.org/licenses/LICENSE-2.0.txt


Copyright 2016 Intel Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cpu

import (
	"os"
	"strings"
	"testing"

	"github.com/intelsdi-x/snap/control/plugin"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/suite"
)

type CPUInfoSuite struct {
	suite.Suite
	MockCPUInfo string
}

func (cis *CPUInfoSuite) SetupSuite() {
	cpuInfo = cis.MockCPUInfo
	createMockCPUInfo(0)
}

func (cis *CPUInfoSuite) TearDownSuite() {
	removeMockCPUInfo()
}

func removeMockCPUInfo() {
	os.Remove(cpuInfo)
}

func TestGetStatsSuite(t *testing.T) {
	suite.Run(t, &CPUInfoSuite{MockCPUInfo: "MockCPUInfo"})
}

func mockNew() *Plugin {
	cpuPlg := New()
	So(cpuPlg, ShouldNotBeNil)
	So(cpuPlg.snapMetricsNames, ShouldNotBeNil)
	cpuPlg.cpuMetricsNumber = 3
	cpuPlg.cpuData = cpuPlg.cpuData[:0]
	for i := 0; i < cpuPlg.cpuMetricsNumber; i++ {
		cpuPlg.cpuData = append(cpuPlg.cpuData, newCPUDataStruct())
	}
	return cpuPlg
}

func createMockCPUInfo(dataSetNumber int) {
	var content string
	if dataSetNumber == 0 {
		content = `cpu  23359837 6006716 1209900 402135131 129307 4 2156 0 0 0
			cpu0 3464284 998669 208226 49355234 57380 3 422 0 0 0
			cpu1 3501681 1012206 189642 49374240 11620 0 278 0 0 0
			intr 33594809 19 2 0 0 0 0 0 9 1 4 0 0 4 0 0 0 31 0 0`
	} else if dataSetNumber == 1 {
		content = `cpu  23472679 6048986 1215282 403105970 129312 4 2158 0 0 0
					cpu0 3480506 1005574 209103 49472588 57381 3 424 0 0 0
					cpu1 3516068 1019269 190413 49493320 11620 0 278 0 0 0
					intr 33594809 19 2 0 0 0 0 0 9 1 4 0 0 4 0 0 0 31 0 0`
	} else if dataSetNumber == 2 {
		content = `cpu  23472679 6048986 1215282 403105970 129312 4 2158 0 0 0
					cpu0 3480506 1005574 209103 49472588 57381 3 424 0 0 0
					cpu1 3516068 1019269 190413 49493320 11620 0 278 0 0 0`
	} else if dataSetNumber == 3 {
		content = `cpu  23472679 6048986 1215282 403105970 129312 4 2158 0 0 0
					cpu0 3480506 1005574 209103 49472588 57381 3 424 0 0
					cpu1 3516068 1019269 190413 49493320 11620 0 278 0 0 0`
	} else if dataSetNumber == 4 {
		content = ``
	}

	cpuInfoContent := []byte(content)
	f, err := os.Create(cpuInfo)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.Write(cpuInfoContent)
}

func (cis *CPUInfoSuite) TestGetMetricTypes() {
	_ = plugin.PluginConfigType{}
	Convey("Given cpu info plugin initialized", cis.T(), func() {
		cpuPlg := mockNew()
		Convey("When one wants to get list of available meterics", func() {
			mts, err := cpuPlg.GetMetricTypes(plugin.PluginConfigType{})

			Convey("Then error should not be reported", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then list of metrics is returned", func() {
				So(len(mts), ShouldEqual, cpuPlg.cpuMetricsNumber*len(cpuPlg.snapMetricsNames)*len(cpuPlg.representationTypes))

				namespaces := []string{}
				for _, m := range mts {
					namespaces = append(namespaces, strings.Join(m.Namespace(), "/"))
				}

				So(namespaces, ShouldContain, "intel/procfs/cpu/all/user/percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/all/nice/percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/all/system/percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/all/idle/percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/all/iowait/percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/all/irq/percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/all/softirq/percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/all/steal/percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/all/guest/percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/all/guest_nice/percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/all/active/percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/all/utilization/percentage")

				So(namespaces, ShouldContain, "intel/procfs/cpu/0/user/percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/0/nice/percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/0/system/percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/0/idle/percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/0/iowait/percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/0/irq/percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/0/softirq/percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/0/steal/percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/0/guest/percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/0/guest_nice/percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/0/active/percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/0/utilization/percentage")

				So(namespaces, ShouldContain, "intel/procfs/cpu/1/user/percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/1/nice/percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/1/system/percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/1/idle/percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/1/iowait/percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/1/irq/percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/1/softirq/percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/1/steal/percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/1/guest/percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/1/guest_nice/percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/1/active/percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/1/utilization/percentage")

				So(namespaces, ShouldContain, "intel/procfs/cpu/all/user")
				So(namespaces, ShouldContain, "intel/procfs/cpu/all/nice")
				So(namespaces, ShouldContain, "intel/procfs/cpu/all/system")
				So(namespaces, ShouldContain, "intel/procfs/cpu/all/idle")
				So(namespaces, ShouldContain, "intel/procfs/cpu/all/iowait")
				So(namespaces, ShouldContain, "intel/procfs/cpu/all/irq")
				So(namespaces, ShouldContain, "intel/procfs/cpu/all/softirq")
				So(namespaces, ShouldContain, "intel/procfs/cpu/all/steal")
				So(namespaces, ShouldContain, "intel/procfs/cpu/all/guest")
				So(namespaces, ShouldContain, "intel/procfs/cpu/all/guest_nice")
				So(namespaces, ShouldContain, "intel/procfs/cpu/all/active")
				So(namespaces, ShouldContain, "intel/procfs/cpu/all/utilization")

				So(namespaces, ShouldContain, "intel/procfs/cpu/0/user")
				So(namespaces, ShouldContain, "intel/procfs/cpu/0/nice")
				So(namespaces, ShouldContain, "intel/procfs/cpu/0/system")
				So(namespaces, ShouldContain, "intel/procfs/cpu/0/idle")
				So(namespaces, ShouldContain, "intel/procfs/cpu/0/iowait")
				So(namespaces, ShouldContain, "intel/procfs/cpu/0/irq")
				So(namespaces, ShouldContain, "intel/procfs/cpu/0/softirq")
				So(namespaces, ShouldContain, "intel/procfs/cpu/0/steal")
				So(namespaces, ShouldContain, "intel/procfs/cpu/0/guest")
				So(namespaces, ShouldContain, "intel/procfs/cpu/0/guest_nice")
				So(namespaces, ShouldContain, "intel/procfs/cpu/0/active")
				So(namespaces, ShouldContain, "intel/procfs/cpu/0/utilization")

				So(namespaces, ShouldContain, "intel/procfs/cpu/1/user")
				So(namespaces, ShouldContain, "intel/procfs/cpu/1/nice")
				So(namespaces, ShouldContain, "intel/procfs/cpu/1/system")
				So(namespaces, ShouldContain, "intel/procfs/cpu/1/idle")
				So(namespaces, ShouldContain, "intel/procfs/cpu/1/iowait")
				So(namespaces, ShouldContain, "intel/procfs/cpu/1/irq")
				So(namespaces, ShouldContain, "intel/procfs/cpu/1/softirq")
				So(namespaces, ShouldContain, "intel/procfs/cpu/1/steal")
				So(namespaces, ShouldContain, "intel/procfs/cpu/1/guest")
				So(namespaces, ShouldContain, "intel/procfs/cpu/1/guest_nice")
				So(namespaces, ShouldContain, "intel/procfs/cpu/1/active")
				So(namespaces, ShouldContain, "intel/procfs/cpu/1/utilization")
			})
		})
	})
}

func (cis *CPUInfoSuite) TestCollectMetrics() {
	Convey("Given cpu plugin initlialized", cis.T(), func() {
		cpuPlg := mockNew()
		So(len(cpuPlg.cpuData), ShouldEqual, 3)

		Convey("When one wants to get values for given metric types", func() {
			mTypes := []plugin.PluginMetricType{
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, allCPU, userProcStat, percentageRepresentationType}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, allCPU, niceProcStat, percentageRepresentationType}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, allCPU, systemProcStat, percentageRepresentationType}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, allCPU, idleProcStat, percentageRepresentationType}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, allCPU, iowaitProcStat, percentageRepresentationType}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, allCPU, irqProcStat, percentageRepresentationType}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, allCPU, softirqProcStat, percentageRepresentationType}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, allCPU, stealProcStat, percentageRepresentationType}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, allCPU, guestProcStat, percentageRepresentationType}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, allCPU, guestNiceProcStat, percentageRepresentationType}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, allCPU, activeProcStat, percentageRepresentationType}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, allCPU, utilizationProcStat, percentageRepresentationType}},

				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, allCPU, userProcStat}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, allCPU, niceProcStat}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, allCPU, systemProcStat}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, allCPU, idleProcStat}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, allCPU, iowaitProcStat}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, allCPU, irqProcStat}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, allCPU, softirqProcStat}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, allCPU, stealProcStat}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, allCPU, guestProcStat}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, allCPU, guestNiceProcStat}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, allCPU, activeProcStat}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, allCPU, utilizationProcStat}},

				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, firstCPU, userProcStat, percentageRepresentationType}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, firstCPU, niceProcStat, percentageRepresentationType}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, firstCPU, systemProcStat, percentageRepresentationType}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, firstCPU, idleProcStat, percentageRepresentationType}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, firstCPU, iowaitProcStat, percentageRepresentationType}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, firstCPU, irqProcStat, percentageRepresentationType}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, firstCPU, softirqProcStat, percentageRepresentationType}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, firstCPU, stealProcStat, percentageRepresentationType}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, firstCPU, guestProcStat, percentageRepresentationType}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, firstCPU, guestNiceProcStat, percentageRepresentationType}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, firstCPU, activeProcStat, percentageRepresentationType}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, firstCPU, utilizationProcStat, percentageRepresentationType}},

				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, firstCPU, userProcStat}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, firstCPU, niceProcStat}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, firstCPU, systemProcStat}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, firstCPU, idleProcStat}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, firstCPU, iowaitProcStat}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, firstCPU, irqProcStat}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, firstCPU, softirqProcStat}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, firstCPU, stealProcStat}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, firstCPU, guestProcStat}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, firstCPU, guestNiceProcStat}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, firstCPU, activeProcStat}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, firstCPU, utilizationProcStat}},

				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, secondCPU, userProcStat, percentageRepresentationType}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, secondCPU, niceProcStat, percentageRepresentationType}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, secondCPU, systemProcStat, percentageRepresentationType}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, secondCPU, idleProcStat, percentageRepresentationType}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, secondCPU, iowaitProcStat, percentageRepresentationType}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, secondCPU, irqProcStat, percentageRepresentationType}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, secondCPU, softirqProcStat, percentageRepresentationType}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, secondCPU, stealProcStat, percentageRepresentationType}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, secondCPU, guestProcStat, percentageRepresentationType}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, secondCPU, guestNiceProcStat, percentageRepresentationType}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, secondCPU, activeProcStat, percentageRepresentationType}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, secondCPU, utilizationProcStat, percentageRepresentationType}},

				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, secondCPU, userProcStat}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, secondCPU, niceProcStat}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, secondCPU, systemProcStat}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, secondCPU, idleProcStat}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, secondCPU, iowaitProcStat}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, secondCPU, irqProcStat}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, secondCPU, softirqProcStat}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, secondCPU, stealProcStat}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, secondCPU, guestProcStat}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, secondCPU, guestNiceProcStat}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, secondCPU, activeProcStat}},
				plugin.PluginMetricType{Namespace_: []string{vendor, fs, pluginName, secondCPU, utilizationProcStat}},
			}

			metrics, err := cpuPlg.CollectMetrics(mTypes)

			Convey("Then no erros should be reported", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then proper metrics values are returned", func() {
				So(len(metrics), ShouldEqual, len(mTypes))
				for _, mt := range metrics {
					So(mt.Data_, ShouldNotBeNil)
				}
			})
		})
	})
}
func (cis *CPUInfoSuite) TestgetCPUMetrics() {
	Convey("Given cpu plugin initlialized", cis.T(), func() {
		cpuPlg := mockNew()
		So(cpuPlg, ShouldNotBeNil)
		So(len(cpuPlg.cpuData), ShouldEqual, 3)

		Convey("We want to check if metrics have proper value", func() {
			//get new data set from /proc/stat
			createMockCPUInfo(0)
			err := cpuPlg.getCPUMetrics()
			So(err, ShouldBeNil)

			var sumStat float64
			//cpu
			So(cpuPlg.cpuData[0].statVal[userProcStat], ShouldEqual, 23359837)
			So(cpuPlg.cpuData[0].statVal[niceProcStat], ShouldEqual, 6006716)
			So(cpuPlg.cpuData[0].statVal[systemProcStat], ShouldEqual, 1209900)
			So(cpuPlg.cpuData[0].statVal[idleProcStat], ShouldEqual, 402135131)
			So(cpuPlg.cpuData[0].statVal[iowaitProcStat], ShouldEqual, 129307)
			So(cpuPlg.cpuData[0].statVal[irqProcStat], ShouldEqual, 4)
			So(cpuPlg.cpuData[0].statVal[softirqProcStat], ShouldEqual, 2156)
			So(cpuPlg.cpuData[0].statVal[stealProcStat], ShouldEqual, 0)
			So(cpuPlg.cpuData[0].statVal[guestProcStat], ShouldEqual, 0)
			So(cpuPlg.cpuData[0].statVal[guestNiceProcStat], ShouldEqual, 0)

			sumStat = cpuPlg.cpuData[0].statVal[userProcStat] + cpuPlg.cpuData[0].statVal[niceProcStat] + cpuPlg.cpuData[0].statVal[systemProcStat]
			sumStat += cpuPlg.cpuData[0].statVal[idleProcStat] + cpuPlg.cpuData[0].statVal[iowaitProcStat] + cpuPlg.cpuData[0].statVal[irqProcStat]
			sumStat += cpuPlg.cpuData[0].statVal[softirqProcStat] + cpuPlg.cpuData[0].statVal[stealProcStat] + cpuPlg.cpuData[0].statVal[guestProcStat] + cpuPlg.cpuData[0].statVal[guestNiceProcStat]

			So(cpuPlg.cpuData[0].statSum, ShouldEqual, sumStat)
			So(cpuPlg.cpuData[0].statVal[activeProcStat], ShouldEqual, sumStat-cpuPlg.cpuData[0].statVal[idleProcStat])
			So(cpuPlg.cpuData[0].statVal[utilizationProcStat], ShouldEqual, sumStat-cpuPlg.cpuData[0].statVal[idleProcStat]-cpuPlg.cpuData[0].statVal[iowaitProcStat])
			prevSumCPU := sumStat

			//cpu0
			So(cpuPlg.cpuData[1].statVal[userProcStat], ShouldEqual, 3464284)
			So(cpuPlg.cpuData[1].statVal[niceProcStat], ShouldEqual, 998669)
			So(cpuPlg.cpuData[1].statVal[systemProcStat], ShouldEqual, 208226)
			So(cpuPlg.cpuData[1].statVal[idleProcStat], ShouldEqual, 49355234)
			So(cpuPlg.cpuData[1].statVal[iowaitProcStat], ShouldEqual, 57380)
			So(cpuPlg.cpuData[1].statVal[irqProcStat], ShouldEqual, 3)
			So(cpuPlg.cpuData[1].statVal[softirqProcStat], ShouldEqual, 422)
			So(cpuPlg.cpuData[1].statVal[stealProcStat], ShouldEqual, 0)
			So(cpuPlg.cpuData[1].statVal[guestProcStat], ShouldEqual, 0)
			So(cpuPlg.cpuData[1].statVal[guestNiceProcStat], ShouldEqual, 0)

			sumStat = cpuPlg.cpuData[1].statVal[userProcStat] + cpuPlg.cpuData[1].statVal[niceProcStat] + cpuPlg.cpuData[1].statVal[systemProcStat]
			sumStat += cpuPlg.cpuData[1].statVal[idleProcStat] + cpuPlg.cpuData[1].statVal[iowaitProcStat] + cpuPlg.cpuData[1].statVal[irqProcStat]
			sumStat += cpuPlg.cpuData[1].statVal[softirqProcStat] + cpuPlg.cpuData[1].statVal[stealProcStat] + cpuPlg.cpuData[1].statVal[guestProcStat] + cpuPlg.cpuData[1].statVal[guestNiceProcStat]

			So(cpuPlg.cpuData[1].statSum, ShouldEqual, sumStat)

			So(cpuPlg.cpuData[1].statVal[activeProcStat], ShouldEqual, sumStat-cpuPlg.cpuData[1].statVal[idleProcStat])
			So(cpuPlg.cpuData[1].statVal[utilizationProcStat], ShouldEqual, sumStat-cpuPlg.cpuData[1].statVal[idleProcStat]-cpuPlg.cpuData[1].statVal[iowaitProcStat])
			prevSumCPU0 := sumStat

			//cpu1
			So(cpuPlg.cpuData[2].statVal[userProcStat], ShouldEqual, 3501681)
			So(cpuPlg.cpuData[2].statVal[niceProcStat], ShouldEqual, 1012206)
			So(cpuPlg.cpuData[2].statVal[systemProcStat], ShouldEqual, 189642)
			So(cpuPlg.cpuData[2].statVal[idleProcStat], ShouldEqual, 49374240)
			So(cpuPlg.cpuData[2].statVal[iowaitProcStat], ShouldEqual, 11620)
			So(cpuPlg.cpuData[2].statVal[irqProcStat], ShouldEqual, 0)
			So(cpuPlg.cpuData[2].statVal[softirqProcStat], ShouldEqual, 278)
			So(cpuPlg.cpuData[2].statVal[stealProcStat], ShouldEqual, 0)
			So(cpuPlg.cpuData[2].statVal[guestProcStat], ShouldEqual, 0)
			So(cpuPlg.cpuData[2].statVal[guestNiceProcStat], ShouldEqual, 0)

			sumStat = cpuPlg.cpuData[2].statVal[userProcStat] + cpuPlg.cpuData[2].statVal[niceProcStat] + cpuPlg.cpuData[2].statVal[systemProcStat]
			sumStat += cpuPlg.cpuData[2].statVal[idleProcStat] + cpuPlg.cpuData[2].statVal[iowaitProcStat] + cpuPlg.cpuData[2].statVal[irqProcStat]
			sumStat += cpuPlg.cpuData[2].statVal[softirqProcStat] + cpuPlg.cpuData[2].statVal[stealProcStat] + cpuPlg.cpuData[2].statVal[guestProcStat] + cpuPlg.cpuData[2].statVal[guestNiceProcStat]

			So(cpuPlg.cpuData[2].statSum, ShouldEqual, sumStat)

			So(cpuPlg.cpuData[2].statVal[activeProcStat], ShouldEqual, sumStat-cpuPlg.cpuData[2].statVal[idleProcStat])
			So(cpuPlg.cpuData[2].statVal[utilizationProcStat], ShouldEqual, sumStat-cpuPlg.cpuData[2].statVal[idleProcStat]-cpuPlg.cpuData[2].statVal[iowaitProcStat])
			prevSumCPU1 := sumStat

			//get new data set from /proc/stat
			createMockCPUInfo(1)
			err = cpuPlg.getCPUMetrics()
			So(err, ShouldBeNil)

			//cpu
			So(cpuPlg.cpuData[0].statVal[userProcStat], ShouldEqual, 23472679)
			So(cpuPlg.cpuData[0].statVal[niceProcStat], ShouldEqual, 6048986)
			So(cpuPlg.cpuData[0].statVal[systemProcStat], ShouldEqual, 1215282)
			So(cpuPlg.cpuData[0].statVal[idleProcStat], ShouldEqual, 403105970)
			So(cpuPlg.cpuData[0].statVal[iowaitProcStat], ShouldEqual, 129312)
			So(cpuPlg.cpuData[0].statVal[irqProcStat], ShouldEqual, 4)
			So(cpuPlg.cpuData[0].statVal[softirqProcStat], ShouldEqual, 2158)
			So(cpuPlg.cpuData[0].statVal[stealProcStat], ShouldEqual, 0)
			So(cpuPlg.cpuData[0].statVal[guestProcStat], ShouldEqual, 0)
			So(cpuPlg.cpuData[0].statVal[guestNiceProcStat], ShouldEqual, 0)

			sumStat = cpuPlg.cpuData[0].statVal[userProcStat] + cpuPlg.cpuData[0].statVal[niceProcStat] + cpuPlg.cpuData[0].statVal[systemProcStat]
			sumStat += cpuPlg.cpuData[0].statVal[idleProcStat] + cpuPlg.cpuData[0].statVal[iowaitProcStat] + cpuPlg.cpuData[0].statVal[irqProcStat]
			sumStat += cpuPlg.cpuData[0].statVal[softirqProcStat] + cpuPlg.cpuData[0].statVal[stealProcStat] + cpuPlg.cpuData[0].statVal[guestProcStat] + cpuPlg.cpuData[0].statVal[guestNiceProcStat]

			So(cpuPlg.cpuData[0].statSum, ShouldEqual, sumStat)
			So(cpuPlg.cpuData[0].statVal[activeProcStat], ShouldEqual, sumStat-cpuPlg.cpuData[0].statVal[idleProcStat])
			So(cpuPlg.cpuData[0].statVal[utilizationProcStat], ShouldEqual, sumStat-cpuPlg.cpuData[0].statVal[idleProcStat]-cpuPlg.cpuData[0].statVal[iowaitProcStat])

			//cpu0
			So(cpuPlg.cpuData[1].statVal[userProcStat], ShouldEqual, 3480506)
			So(cpuPlg.cpuData[1].statVal[niceProcStat], ShouldEqual, 1005574)
			So(cpuPlg.cpuData[1].statVal[systemProcStat], ShouldEqual, 209103)
			So(cpuPlg.cpuData[1].statVal[idleProcStat], ShouldEqual, 49472588)
			So(cpuPlg.cpuData[1].statVal[iowaitProcStat], ShouldEqual, 57381)
			So(cpuPlg.cpuData[1].statVal[irqProcStat], ShouldEqual, 3)
			So(cpuPlg.cpuData[1].statVal[softirqProcStat], ShouldEqual, 424)
			So(cpuPlg.cpuData[1].statVal[stealProcStat], ShouldEqual, 0)
			So(cpuPlg.cpuData[1].statVal[guestProcStat], ShouldEqual, 0)
			So(cpuPlg.cpuData[1].statVal[guestNiceProcStat], ShouldEqual, 0)

			sumStat = cpuPlg.cpuData[1].statVal[userProcStat] + cpuPlg.cpuData[1].statVal[niceProcStat] + cpuPlg.cpuData[1].statVal[systemProcStat]
			sumStat += cpuPlg.cpuData[1].statVal[idleProcStat] + cpuPlg.cpuData[1].statVal[iowaitProcStat] + cpuPlg.cpuData[1].statVal[irqProcStat]
			sumStat += cpuPlg.cpuData[1].statVal[softirqProcStat] + cpuPlg.cpuData[1].statVal[stealProcStat] + cpuPlg.cpuData[1].statVal[guestProcStat] + cpuPlg.cpuData[1].statVal[guestNiceProcStat]

			So(cpuPlg.cpuData[1].statSum, ShouldEqual, sumStat)

			So(cpuPlg.cpuData[1].statVal[activeProcStat], ShouldEqual, sumStat-cpuPlg.cpuData[1].statVal[idleProcStat])
			So(cpuPlg.cpuData[1].statVal[utilizationProcStat], ShouldEqual, sumStat-cpuPlg.cpuData[1].statVal[idleProcStat]-cpuPlg.cpuData[1].statVal[iowaitProcStat])

			//cpu1
			So(cpuPlg.cpuData[2].statVal[userProcStat], ShouldEqual, 3516068)
			So(cpuPlg.cpuData[2].statVal[niceProcStat], ShouldEqual, 1019269)
			So(cpuPlg.cpuData[2].statVal[systemProcStat], ShouldEqual, 190413)
			So(cpuPlg.cpuData[2].statVal[idleProcStat], ShouldEqual, 49493320)
			So(cpuPlg.cpuData[2].statVal[iowaitProcStat], ShouldEqual, 11620)
			So(cpuPlg.cpuData[2].statVal[irqProcStat], ShouldEqual, 0)
			So(cpuPlg.cpuData[2].statVal[softirqProcStat], ShouldEqual, 278)
			So(cpuPlg.cpuData[2].statVal[stealProcStat], ShouldEqual, 0)
			So(cpuPlg.cpuData[2].statVal[guestProcStat], ShouldEqual, 0)
			So(cpuPlg.cpuData[2].statVal[guestNiceProcStat], ShouldEqual, 0)

			sumStat = cpuPlg.cpuData[2].statVal[userProcStat] + cpuPlg.cpuData[2].statVal[niceProcStat] + cpuPlg.cpuData[2].statVal[systemProcStat]
			sumStat += cpuPlg.cpuData[2].statVal[idleProcStat] + cpuPlg.cpuData[2].statVal[iowaitProcStat] + cpuPlg.cpuData[2].statVal[irqProcStat]
			sumStat += cpuPlg.cpuData[2].statVal[softirqProcStat] + cpuPlg.cpuData[2].statVal[stealProcStat] + cpuPlg.cpuData[2].statVal[guestProcStat] + cpuPlg.cpuData[2].statVal[guestNiceProcStat]

			So(cpuPlg.cpuData[2].statSum, ShouldEqual, sumStat)

			So(cpuPlg.cpuData[2].statVal[activeProcStat], ShouldEqual, sumStat-cpuPlg.cpuData[2].statVal[idleProcStat])
			So(cpuPlg.cpuData[2].statVal[utilizationProcStat], ShouldEqual, sumStat-cpuPlg.cpuData[2].statVal[idleProcStat]-cpuPlg.cpuData[2].statVal[iowaitProcStat])

			//ensuring that optimized method of metric calculation is correct
			//cpu percentage
			percentageVal, errPercentage := cpuPlg.getMetricValue(allCPU, userProcStat, percentageRepresentationType)
			So(percentageVal, ShouldEqual, 100*cpuPlg.cpuData[0].rateVal[userProcStat]/(cpuPlg.cpuData[0].statSum-prevSumCPU))

			percentageVal, errPercentage = cpuPlg.getMetricValue(allCPU, niceProcStat, percentageRepresentationType)
			So(percentageVal, ShouldEqual, 100*cpuPlg.cpuData[0].rateVal[niceProcStat]/(cpuPlg.cpuData[0].statSum-prevSumCPU))

			percentageVal, errPercentage = cpuPlg.getMetricValue(allCPU, idleProcStat, percentageRepresentationType)
			So(percentageVal, ShouldEqual, 100*cpuPlg.cpuData[0].rateVal[idleProcStat]/(cpuPlg.cpuData[0].statSum-prevSumCPU))

			percentageVal, errPercentage = cpuPlg.getMetricValue(allCPU, iowaitProcStat, percentageRepresentationType)
			So(percentageVal, ShouldEqual, 100*cpuPlg.cpuData[0].rateVal[iowaitProcStat]/(cpuPlg.cpuData[0].statSum-prevSumCPU))

			percentageVal, errPercentage = cpuPlg.getMetricValue(allCPU, irqProcStat, percentageRepresentationType)
			So(percentageVal, ShouldEqual, 100*cpuPlg.cpuData[0].rateVal[irqProcStat]/(cpuPlg.cpuData[0].statSum-prevSumCPU))

			percentageVal, errPercentage = cpuPlg.getMetricValue(allCPU, softirqProcStat, percentageRepresentationType)
			So(percentageVal, ShouldEqual, 100*cpuPlg.cpuData[0].rateVal[softirqProcStat]/(cpuPlg.cpuData[0].statSum-prevSumCPU))

			percentageVal, errPercentage = cpuPlg.getMetricValue(allCPU, stealProcStat, percentageRepresentationType)
			So(percentageVal, ShouldEqual, 100*cpuPlg.cpuData[0].rateVal[stealProcStat]/(cpuPlg.cpuData[0].statSum-prevSumCPU))

			percentageVal, errPercentage = cpuPlg.getMetricValue(allCPU, guestProcStat, percentageRepresentationType)
			So(percentageVal, ShouldEqual, 100*cpuPlg.cpuData[0].rateVal[guestProcStat]/(cpuPlg.cpuData[0].statSum-prevSumCPU))

			percentageVal, errPercentage = cpuPlg.getMetricValue(allCPU, guestNiceProcStat, percentageRepresentationType)
			So(percentageVal, ShouldEqual, 100*cpuPlg.cpuData[0].rateVal[guestNiceProcStat]/(cpuPlg.cpuData[0].statSum-prevSumCPU))

			rateSum := cpuPlg.cpuData[0].rateVal[userProcStat] + cpuPlg.cpuData[0].rateVal[niceProcStat] + cpuPlg.cpuData[0].rateVal[systemProcStat]
			rateSum += cpuPlg.cpuData[0].rateVal[idleProcStat] + cpuPlg.cpuData[0].rateVal[iowaitProcStat] + cpuPlg.cpuData[0].rateVal[irqProcStat]
			rateSum += cpuPlg.cpuData[0].rateVal[softirqProcStat] + cpuPlg.cpuData[0].rateVal[stealProcStat] + cpuPlg.cpuData[0].rateVal[guestProcStat] + cpuPlg.cpuData[0].rateVal[guestNiceProcStat]

			So(rateSum, ShouldEqual, cpuPlg.cpuData[0].statSum-prevSumCPU)

			percentageVal, errPercentage = cpuPlg.getMetricValue(allCPU, activeProcStat, percentageRepresentationType)
			So(percentageVal, ShouldEqual, 100*(rateSum-cpuPlg.cpuData[0].rateVal[idleProcStat])/rateSum)

			percentageVal, errPercentage = cpuPlg.getMetricValue(allCPU, utilizationProcStat, percentageRepresentationType)
			So(percentageVal, ShouldEqual, 100*(rateSum-cpuPlg.cpuData[0].rateVal[iowaitProcStat]-cpuPlg.cpuData[0].rateVal[idleProcStat])/rateSum)

			//cpu1 percentage
			percentageVal, errPercentage = cpuPlg.getMetricValue(firstCPU, userProcStat, percentageRepresentationType)
			So(percentageVal, ShouldEqual, 100*cpuPlg.cpuData[1].rateVal[userProcStat]/(cpuPlg.cpuData[1].statSum-prevSumCPU0))

			percentageVal, errPercentage = cpuPlg.getMetricValue(firstCPU, niceProcStat, percentageRepresentationType)
			So(percentageVal, ShouldEqual, 100*cpuPlg.cpuData[1].rateVal[niceProcStat]/(cpuPlg.cpuData[1].statSum-prevSumCPU0))

			percentageVal, errPercentage = cpuPlg.getMetricValue(firstCPU, idleProcStat, percentageRepresentationType)
			So(percentageVal, ShouldEqual, 100*cpuPlg.cpuData[1].rateVal[idleProcStat]/(cpuPlg.cpuData[1].statSum-prevSumCPU0))

			percentageVal, errPercentage = cpuPlg.getMetricValue(firstCPU, iowaitProcStat, percentageRepresentationType)
			So(percentageVal, ShouldEqual, 100*cpuPlg.cpuData[1].rateVal[iowaitProcStat]/(cpuPlg.cpuData[1].statSum-prevSumCPU0))

			percentageVal, errPercentage = cpuPlg.getMetricValue(firstCPU, irqProcStat, percentageRepresentationType)
			So(percentageVal, ShouldEqual, 100*cpuPlg.cpuData[1].rateVal[irqProcStat]/(cpuPlg.cpuData[1].statSum-prevSumCPU0))

			percentageVal, errPercentage = cpuPlg.getMetricValue(firstCPU, softirqProcStat, percentageRepresentationType)
			So(percentageVal, ShouldEqual, 100*cpuPlg.cpuData[1].rateVal[softirqProcStat]/(cpuPlg.cpuData[1].statSum-prevSumCPU0))

			percentageVal, errPercentage = cpuPlg.getMetricValue(firstCPU, stealProcStat, percentageRepresentationType)
			So(percentageVal, ShouldEqual, 100*cpuPlg.cpuData[1].rateVal[stealProcStat]/(cpuPlg.cpuData[1].statSum-prevSumCPU0))

			percentageVal, errPercentage = cpuPlg.getMetricValue(firstCPU, guestProcStat, percentageRepresentationType)
			So(percentageVal, ShouldEqual, 100*cpuPlg.cpuData[0].rateVal[guestProcStat]/(cpuPlg.cpuData[1].statSum-prevSumCPU0))

			percentageVal, errPercentage = cpuPlg.getMetricValue(firstCPU, guestNiceProcStat, percentageRepresentationType)
			So(percentageVal, ShouldEqual, 100*cpuPlg.cpuData[1].rateVal[guestNiceProcStat]/(cpuPlg.cpuData[1].statSum-prevSumCPU0))

			rateSum = cpuPlg.cpuData[1].rateVal[userProcStat] + cpuPlg.cpuData[1].rateVal[niceProcStat] + cpuPlg.cpuData[1].rateVal[systemProcStat]
			rateSum += cpuPlg.cpuData[1].rateVal[idleProcStat] + cpuPlg.cpuData[1].rateVal[iowaitProcStat] + cpuPlg.cpuData[1].rateVal[irqProcStat]
			rateSum += cpuPlg.cpuData[1].rateVal[softirqProcStat] + cpuPlg.cpuData[1].rateVal[stealProcStat] + cpuPlg.cpuData[1].rateVal[guestProcStat] + cpuPlg.cpuData[1].rateVal[guestNiceProcStat]

			So(rateSum, ShouldEqual, cpuPlg.cpuData[1].statSum-prevSumCPU0)

			percentageVal, errPercentage = cpuPlg.getMetricValue(firstCPU, activeProcStat, percentageRepresentationType)
			So(percentageVal, ShouldEqual, 100*(rateSum-cpuPlg.cpuData[1].rateVal[idleProcStat])/rateSum)

			percentageVal, errPercentage = cpuPlg.getMetricValue(firstCPU, utilizationProcStat, percentageRepresentationType)
			So(percentageVal, ShouldEqual, 100*(rateSum-cpuPlg.cpuData[1].rateVal[iowaitProcStat]-cpuPlg.cpuData[1].rateVal[idleProcStat])/rateSum)

			//cpu2 percentage
			percentageVal, errPercentage = cpuPlg.getMetricValue(secondCPU, userProcStat, percentageRepresentationType)
			So(percentageVal, ShouldEqual, 100*cpuPlg.cpuData[2].rateVal[userProcStat]/(cpuPlg.cpuData[2].statSum-prevSumCPU1))

			percentageVal, errPercentage = cpuPlg.getMetricValue(secondCPU, userProcStat, percentageRepresentationType)
			So(percentageVal, ShouldEqual, 100*cpuPlg.cpuData[2].rateVal[userProcStat]/(cpuPlg.cpuData[2].statSum-prevSumCPU1))

			percentageVal, errPercentage = cpuPlg.getMetricValue(secondCPU, niceProcStat, percentageRepresentationType)
			So(percentageVal, ShouldEqual, 100*cpuPlg.cpuData[2].rateVal[niceProcStat]/(cpuPlg.cpuData[2].statSum-prevSumCPU1))

			percentageVal, errPercentage = cpuPlg.getMetricValue(secondCPU, idleProcStat, percentageRepresentationType)
			So(percentageVal, ShouldEqual, 100*cpuPlg.cpuData[2].rateVal[idleProcStat]/(cpuPlg.cpuData[2].statSum-prevSumCPU1))

			percentageVal, errPercentage = cpuPlg.getMetricValue(secondCPU, iowaitProcStat, percentageRepresentationType)
			So(percentageVal, ShouldEqual, 100*cpuPlg.cpuData[2].rateVal[iowaitProcStat]/(cpuPlg.cpuData[2].statSum-prevSumCPU1))

			percentageVal, errPercentage = cpuPlg.getMetricValue(secondCPU, irqProcStat, percentageRepresentationType)
			So(percentageVal, ShouldEqual, 100*cpuPlg.cpuData[2].rateVal[irqProcStat]/(cpuPlg.cpuData[2].statSum-prevSumCPU1))

			percentageVal, errPercentage = cpuPlg.getMetricValue(secondCPU, softirqProcStat, percentageRepresentationType)
			So(percentageVal, ShouldEqual, 100*cpuPlg.cpuData[2].rateVal[softirqProcStat]/(cpuPlg.cpuData[2].statSum-prevSumCPU1))

			percentageVal, errPercentage = cpuPlg.getMetricValue(secondCPU, stealProcStat, percentageRepresentationType)
			So(percentageVal, ShouldEqual, 100*cpuPlg.cpuData[2].rateVal[stealProcStat]/(cpuPlg.cpuData[2].statSum-prevSumCPU1))

			percentageVal, errPercentage = cpuPlg.getMetricValue(secondCPU, guestProcStat, percentageRepresentationType)
			So(percentageVal, ShouldEqual, 100*cpuPlg.cpuData[2].rateVal[guestProcStat]/(cpuPlg.cpuData[2].statSum-prevSumCPU1))

			percentageVal, errPercentage = cpuPlg.getMetricValue(secondCPU, guestNiceProcStat, percentageRepresentationType)
			So(percentageVal, ShouldEqual, 100*cpuPlg.cpuData[2].rateVal[guestNiceProcStat]/(cpuPlg.cpuData[2].statSum-prevSumCPU1))

			rateSum = cpuPlg.cpuData[2].rateVal[userProcStat] + cpuPlg.cpuData[2].rateVal[niceProcStat] + cpuPlg.cpuData[2].rateVal[systemProcStat]
			rateSum += cpuPlg.cpuData[2].rateVal[idleProcStat] + cpuPlg.cpuData[2].rateVal[iowaitProcStat] + cpuPlg.cpuData[2].rateVal[irqProcStat]
			rateSum += cpuPlg.cpuData[2].rateVal[softirqProcStat] + cpuPlg.cpuData[2].rateVal[stealProcStat] + cpuPlg.cpuData[2].rateVal[guestProcStat] + cpuPlg.cpuData[2].rateVal[guestNiceProcStat]

			So(rateSum, ShouldEqual, cpuPlg.cpuData[2].statSum-prevSumCPU1)

			percentageVal, errPercentage = cpuPlg.getMetricValue(secondCPU, activeProcStat, percentageRepresentationType)
			So(percentageVal, ShouldEqual, 100*(rateSum-cpuPlg.cpuData[2].rateVal[idleProcStat])/rateSum)

			percentageVal, errPercentage = cpuPlg.getMetricValue(secondCPU, utilizationProcStat, percentageRepresentationType)
			So(percentageVal, ShouldEqual, 100*(rateSum-cpuPlg.cpuData[2].rateVal[iowaitProcStat]-cpuPlg.cpuData[2].rateVal[idleProcStat])/rateSum)

			So(errPercentage, ShouldBeNil)

		})
	})
}

func (cis *CPUInfoSuite) TestgetCPUTabID() {
	Convey("Given cpu plugin initlialized", cis.T(), func() {
		cpuPlg := mockNew()
		So(cpuPlg, ShouldNotBeNil)
		So(len(cpuPlg.cpuData), ShouldEqual, 3)
		Convey("We want to check CPU indefier in array of CPUs data", func() {
			_, err := cpuPlg.getCPUTabID("abc")
			So(err, ShouldNotBeNil)
			_, err = cpuPlg.getCPUTabID("3")
			So(err, ShouldNotBeNil)
			_, err = cpuPlg.getCPUTabID("5")
			So(err, ShouldNotBeNil)
			_, err = cpuPlg.getCPUTabID("-1")
			So(err, ShouldNotBeNil)
			_, err = cpuPlg.getCPUTabID("0")
			So(err, ShouldBeNil)
			_, err = cpuPlg.getCPUTabID(allCPU)
			So(err, ShouldBeNil)
		})
	})
}

func (cis *CPUInfoSuite) TestgetInitialProcStatData() {
	Convey("Given cpu plugin initlialized", cis.T(), func() {
		Convey("We want to check initial reading of /proc/stat", func() {
			createMockCPUInfo(1)
			_, _, err := getInitialProcStatData()
			So(err, ShouldBeNil)
			createMockCPUInfo(2)
			_, _, err = getInitialProcStatData()
			So(err, ShouldBeNil)
			createMockCPUInfo(3)
			_, _, err = getInitialProcStatData()
			So(err, ShouldNotBeNil)
			createMockCPUInfo(4)
			_, _, err = getInitialProcStatData()
			So(err, ShouldNotBeNil)
		})
	})
}
