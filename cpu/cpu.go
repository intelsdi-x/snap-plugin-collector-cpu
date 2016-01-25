// +build linux

/*
http://www.apache.org/licenses/LICENSE-2.0.txt


Copyright 2015 Intel Corporation

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

//package cpu provides implementation of snap-plugin-collector-cpu plugin
package cpu

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
)

const (
	// vendor namespace part
	vendor = "intel"

	// fs namespace part
	fs = "procfs"

	// pluginName namespace part
	pluginName = "cpu"

	// version of cpu plugin
	version = 1

	//pluginType type of plugin
	pluginType = plugin.CollectorPluginType

	//userProcStat "user" metric from /proc/stat
	userProcStat = "user"

	//userProcStat "nice" metric from /proc/stat
	niceProcStat = "nice"

	//systemProcStat "system" metric from /proc/stat
	systemProcStat = "system"

	//idleProcStat "idle" metric from /proc/stat
	idleProcStat = "idle"
	//iowaitProcStat "iowait" metric from /proc/stat
	iowaitProcStat = "iowait"

	//irqProcStat "irq" metric from /proc/stat
	irqProcStat = "irq"

	//softirqProcStat "softirq" metric from /proc/stat
	softirqProcStat = "softirq"

	//stealProcStat "steal" metric from /proc/stat
	stealProcStat = "steal"

	//guestProcStat "guest" metric from /proc/stat
	guestProcStat = "guest"

	//guestNiceProcStat "guest_nice" metric from /proc/stat
	guestNiceProcStat = "guest_nice"

	//activeProcStat "active" snap metric
	activeProcStat = "active"

	//activeProcStat "utilization" snap metric
	utilizationProcStat = "utilization"

	//jiffiesRepresentationType jiffies representation type
	jiffiesRepresentationType = ""

	//percentageRepresentationType percentage representation type
	percentageRepresentationType = "percentage"

	//maxPercentageNamespaceSize max size of namespace for metrics in percentage
	maxPercentageNamespaceSize = 6

	//maxJiffiesNamespaceSize max size of namespace for metrics in jiffies
	maxJiffiesNamespaceSize = 5

	//allCPU string indentifier for aggregation metrics (for all CPUs)
	allCPU = "all"

	//firstCPU string indentifier for the first CPU
	firstCPU = "0"

	//secondCPU string indentifier for the second CPU
	secondCPU = "1"

	//cpuStr string indentifier for /proc/stat line which have desired CPU metrics
	cpuStr = "cpu"
)

//cpuInfo source of data for metrics
var cpuInfo = "/proc/stat"

// make sure that we actually satisify requierd interface
var _ plugin.CollectorPlugin = (*Plugin)(nil)

//Plugin cpu plugin struct which gathers plugin specific data
type Plugin struct {
	host                     string
	cpuMetricsNumber         int // number of cpu + "all" metric
	cpuData                  []cpuDataStruct
	procStatMetricsNames     []string
	snapSpecificMetricsNames []string
	snapMetricsNames         []string
	representationTypes      []string
}

//cpuDataStruct represents one CPU in CPUs data array
type cpuDataStruct struct {
	statVal map[string]float64
	rateVal map[string]float64 //difference between current and previous value
	rateSum float64
	statSum float64
}

//New creates new instance of cpu plugin
func New() *Plugin {
	var p *Plugin
	p = nil

	fh, err := os.Open(cpuInfo)
	if err != nil {
		return p
	}
	defer fh.Close()

	cpuMetricsNumber, procStatMetricsNumber, err := getInitialProcStatData()
	if err != nil {
		return p
	}

	//build CPUs data array
	var cpuData []cpuDataStruct
	for i := 0; i < cpuMetricsNumber; i++ {
		cpuData = append(cpuData, newCPUDataStruct())
	}
	//initialize metric names arrays
	procStatMetricsNames := []string{userProcStat, niceProcStat, systemProcStat, idleProcStat, iowaitProcStat, irqProcStat, softirqProcStat, stealProcStat, guestProcStat, guestNiceProcStat}
	snapSpecificMetricsNames := []string{activeProcStat, utilizationProcStat}
	representationTypes := []string{jiffiesRepresentationType, percentageRepresentationType}

	//build snapMetricsNames to support different kernels
	var snapMetricsNames []string
	snapMetricsNames = append(snapMetricsNames, procStatMetricsNames[0:procStatMetricsNumber]...)
	snapMetricsNames = append(snapMetricsNames, snapSpecificMetricsNames...)

	host, err := os.Hostname()
	if err != nil {
		host = "localhost"
	}

	p = &Plugin{
		host:                     host,
		cpuMetricsNumber:         cpuMetricsNumber,
		cpuData:                  cpuData,
		procStatMetricsNames:     procStatMetricsNames,
		snapSpecificMetricsNames: snapSpecificMetricsNames,
		snapMetricsNames:         snapMetricsNames,
		representationTypes:      representationTypes}
	return p
}

// GetMetricTypes returns list of available metric types
// It returns error in case retrieval was not successful
func (p *Plugin) GetMetricTypes(cfg plugin.PluginConfigType) ([]plugin.PluginMetricType, error) {
	mts := []plugin.PluginMetricType{}

	for i := range p.cpuData {
		for j := range p.snapMetricsNames {
			for k := range p.representationTypes {
				cpuStrID := allCPU //get cpu identifier for aggregation metric (for all cpu)
				if i != 0 {
					cpuStrID = strconv.Itoa(i - 1) //get cpu identifier (one value lower than cpu identifier in cpu's data array)
				}
				ns := []string{vendor, fs, pluginName, cpuStrID, p.snapMetricsNames[j]}
				if p.representationTypes[k] == percentageRepresentationType {
					ns = append(ns, percentageRepresentationType)
				}
				mt := plugin.PluginMetricType{Namespace_: ns}
				mts = append(mts, mt)
			}
		}
	}
	return mts, nil
}

// CollectMetrics returns list of requested metric values
// It returns error in case retrieval was not successful
func (p *Plugin) CollectMetrics(mts []plugin.PluginMetricType) ([]plugin.PluginMetricType, error) {
	metrics := []plugin.PluginMetricType{}
	err := p.getCPUMetrics()
	if err != nil {
		return metrics, err
	}

	for _, mt := range mts {
		ns := mt.Namespace()
		cpuID, metricID, representationType, err := p.getMetricParameters(ns)
		if err != nil {
			return metrics, err
		}
		val, err := p.getMetricValue(cpuID, metricID, representationType)
		if err != nil {
			return metrics, err
		}

		mt := plugin.PluginMetricType{
			Namespace_: ns,
			Data_:      val,
			Source_:    p.host,
			Timestamp_: time.Now(),
		}
		metrics = append(metrics, mt)
	}
	return metrics, err
}

// GetConfigPolicy returns config policy
// It returns error in case retrieval was not successful
func (p *Plugin) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	return cpolicy.New(), nil
}

//Meta returns meta data for plugin
func Meta() *plugin.PluginMeta {
	return plugin.NewPluginMeta(
		pluginName,
		version,
		pluginType,
		[]string{},
		[]string{plugin.SnapGOBContentType},
		plugin.ConcurrencyCount(1))
}

//newCPUDataStruct initializes one element of CPUs data array
func newCPUDataStruct() (cpu cpuDataStruct) {
	cpu.statVal = make(map[string]float64)
	cpu.rateVal = make(map[string]float64)
	cpu.rateSum = 0
	cpu.statSum = 0
	return cpu
}

//getInitialProcStatData gets number of CPUs needed to set CPUs data array and number of metrics available in /proc/stat output
func getInitialProcStatData() (cpuMetricsNumber int, procStatMetricNumber int, err error) {
	fh, err := os.Open(cpuInfo)
	if err != nil {
		return cpuMetricsNumber, procStatMetricNumber, err
	}
	defer fh.Close()

	scanner := bufio.NewScanner(fh)
	//read the first line to start the loop
	scanErr := scanner.Scan()
	if scanErr != true {
		return cpuMetricsNumber, procStatMetricNumber, fmt.Errorf("Cannot read from /proc/stat")
	}

	procStatLine := strings.Fields(scanner.Text())
	procStatMetricNumber = len(procStatLine) - 1 //without cpu identifier (e.g. cpu or cpu0)

	cpuMetricsNumber = 0
	for strings.Contains(procStatLine[0], cpuStr) {
		//check line length, compare current length with length of the first line (length of the firs line = len(procStatMetricNumber) + CPU identifier
		if len(procStatLine) != (procStatMetricNumber + 1) {
			return cpuMetricsNumber, procStatMetricNumber, fmt.Errorf("Incorrect /proc/stat output")
		}
		cpuMetricsNumber++
		//read the next line to be able to check loop condition
		scanErr = scanner.Scan()
		if scanErr != true {
			break //no more lines to read
		}
		procStatLine = strings.Fields(scanner.Text())
	}
	return cpuMetricsNumber, procStatMetricNumber, err
}

//updateCPUMetricsVal updates metric values. collects information needed for aggregation metrics and write data to CPUs data array
func (p *Plugin) updateCPUMetricsVal(cpuNumber int, metricName string, metricVal float64) {
	if mapKeyExists(metricName, p.cpuData[cpuNumber].statVal) {
		//if value of metrics exist then we could calculate rate value
		p.cpuData[cpuNumber].rateVal[metricName] = metricVal - p.cpuData[cpuNumber].statVal[metricName]
		//update sum of rate values
		if metricName != activeProcStat && metricName != utilizationProcStat {
			p.cpuData[cpuNumber].rateSum = p.cpuData[cpuNumber].rateSum + p.cpuData[cpuNumber].rateVal[metricName]
		}
	} else {
		//previous value of metric does not exist - we can not calculate rate value
		p.cpuData[cpuNumber].rateVal[metricName] = 0
	}

	p.cpuData[cpuNumber].statVal[metricName] = metricVal
	if metricName != activeProcStat && metricName != utilizationProcStat {
		p.cpuData[cpuNumber].statSum = p.cpuData[cpuNumber].statSum + metricVal
	}
}

//getActiveStatMetric calculates CPU active state in the same way as it is made in Collectd
func (p *Plugin) getActiveStatMetric(cpuNumber int) (idleStat float64) {
	if mapKeyExists(idleProcStat, p.cpuData[cpuNumber].statVal) {
		idleStat = p.cpuData[cpuNumber].statVal[idleProcStat]
	}

	return (p.cpuData[cpuNumber].statSum - idleStat)
}

//getUtilizationStatMetric calculates CPU utilization
func (p *Plugin) getUtilizationStatMetric(cpuNumber int) (idleStat float64) {
	if mapKeyExists(idleProcStat, p.cpuData[cpuNumber].statVal) {
		idleStat = p.cpuData[cpuNumber].statVal[idleProcStat]
	}

	if mapKeyExists(iowaitProcStat, p.cpuData[cpuNumber].statVal) {
		idleStat += p.cpuData[cpuNumber].statVal[iowaitProcStat]
	}

	return (p.cpuData[cpuNumber].statSum - idleStat)
}

//getCPUMetrics gets information from /proc/stat output and fills CPUs data array
func (p *Plugin) getCPUMetrics() error {
	fh, err := os.Open(cpuInfo)
	if err != nil {
		return err
	}
	defer fh.Close()
	scanner := bufio.NewScanner(fh)

	for i := range p.cpuData {
		scanErr := scanner.Scan()
		if scanErr != true {
			return fmt.Errorf("Cannot read from /proc/stat")
		}
		procStatLine := strings.Fields(scanner.Text())
		procStatLineLength := len(procStatLine)

		//get desired /proc/stat/ line length = number of possible metrics from /proc/stat + CPU identifier (e.g. cpu or cpu0)
		//current /proc/stat/ line length should have the same length as during initialization
		desiredProcStatLineLength := (len(p.snapMetricsNames) - len(p.snapSpecificMetricsNames)) + 1
		if procStatLineLength <= 1 || procStatLineLength != desiredProcStatLineLength {
			err := fmt.Errorf("Incorrect length of /proc/stat/ line output (len: %d)", procStatLineLength)
			return err
		}

		//skip the first column from /proc/stat/ output
		procStatLine = procStatLine[1:procStatLineLength]

		//clean sum
		p.cpuData[i].rateSum = 0
		p.cpuData[i].statSum = 0

		//get metrics from /proc/stat output
		for j := range procStatLine {
			metricName := p.procStatMetricsNames[j]
			metricVal, err := strconv.ParseFloat(procStatLine[j], 64)
			if err != nil {
				return err
			}
			p.updateCPUMetricsVal(i, metricName, metricVal)
		}

		activeStat := p.getActiveStatMetric(i)
		p.updateCPUMetricsVal(i, activeProcStat, activeStat)

		utilizationStat := p.getUtilizationStatMetric(i)
		p.updateCPUMetricsVal(i, utilizationProcStat, utilizationStat)
	}
	return err
}

//getCPUTabID gets CPU identifier in array of CPUs data and checks identifier correctness
func (p *Plugin) getCPUTabID(cpuID string) (cpuTabID int, err error) {
	//get CPU identifier
	if cpuID == allCPU {
		cpuTabID = 0
		return cpuTabID, err
	} else {
		cpuTabID, err = strconv.Atoi(cpuID)
		if err != nil {
			return cpuTabID, err
		}
		cpuTabID = cpuTabID + 1 //CPUs are numerated from 0, identifier in array is one value higher, cpu0 has id equal to 1, cpu1 has id equal to 2
		//check identifier correctness
		if cpuTabID < 1 || cpuTabID >= len(p.cpuData) {
			err = fmt.Errorf("Incorrect CPU identifier (CPU identifier: %s )", cpuID)
			return cpuTabID, err
		}
	}
	return cpuTabID, err
}

//getMetricValue gets value of metric in jiffies or as percentage value
func (p *Plugin) getMetricValue(cpuID string, metricName string, representationType string) (val float64, err error) {
	cpuTabID, err := p.getCPUTabID(cpuID)
	if err != nil {
		return val, err
	}

	if representationType == jiffiesRepresentationType && mapKeyExists(metricName, p.cpuData[cpuTabID].statVal) {
		val = p.cpuData[cpuTabID].statVal[metricName]
	} else if representationType == percentageRepresentationType && mapKeyExists(metricName, p.cpuData[cpuTabID].rateVal) {
		if p.cpuData[cpuTabID].rateSum != 0 {
			val = 100 * p.cpuData[cpuTabID].rateVal[metricName] / p.cpuData[cpuTabID].rateSum
		} else {
			val = 0
		}
	} else {
		err = fmt.Errorf("Incorrect metric name (metric name: %s )", metricName)
	}
	return val, err
}

//getMetricParameters get parameters needed to get value of metric
func (p *Plugin) getMetricParameters(ns []string) (string, string, string, error) {
	var cpuID string
	var metricID string
	var representationType string
	var err error

	if len(ns) == maxJiffiesNamespaceSize || len(ns) == maxPercentageNamespaceSize {
		cpuID = ns[3]
		metricID = ns[4]
		representationType = jiffiesRepresentationType
		if len(ns) == maxPercentageNamespaceSize {
			representationType = ns[maxPercentageNamespaceSize-1]
		}
	} else {
		err = fmt.Errorf("Incorrect namespace length, (len = %d)", len(ns))
	}
	return cpuID, metricID, representationType, err
}

//mapKeyExists checks if element with given key exist in map
func mapKeyExists(key string, m map[string]float64) bool {
	var ret = false
	_, ok := m[key]
	if ok {
		ret = true
	}
	return ret
}
