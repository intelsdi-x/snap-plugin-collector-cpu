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

//Package cpu provides implementation of snap-plugin-collector-cpu plugin
package cpu

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
)

const (
	//vendor namespace part
	vendor = "intel"

	//fs namespace part
	fs = "procfs"

	//pluginName namespace part
	Name = "cpu"

	// version of cpu plugin
	Version = 7

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

	//utilizationProcStat "utilization" snap metric
	utilizationProcStat = "utilization"

	//jiffiesRepresentationType jiffies representation type
	jiffiesRepresentationType = "jiffies"

	//percentageRepresentationType percentage representation type
	percentageRepresentationType = "percentage"

	//maxNamespaceSize max size of namespace for metrics
	maxNamespaceSize = 5

	//allCPU string indentifier for aggregation metrics (for all CPUs)
	allCPU = "all"

	//cpuStr string indentifier for /proc/stat line which have desired CPU metrics
	cpuStr = "cpu"
)

// CPUCollector plugin struct which gathers plugin specific data

/* stats - metrics per cpu read from file /proc/stat:
map ["all": map["user_jiffies": x
		"user_percentage" x
		... ]
     "0": map["user_jiffies": x
	      "user_percentage" x
	      ... ]
     "1": ... ]
*/
type CPUCollector struct {
	initialized          bool
	proc_path            string
	cpuMetricsNumber     int // number of cpu + "all" metric
	stats                map[string]map[string]interface{}
	prevMetricsSum       map[string]float64
	procStatMetricsNames []string
	snapMetricsNames     []string
}

// defaultProcPath source of data for metrics
var defaultProcPath = "/proc"

// New creates instance of interface info plugin
func New() *CPUCollector {
	return &CPUCollector{
		proc_path: defaultProcPath + "/stat",
	}
}

// GetConfigPolicy returns config policy
// It returns error in case retrieval was not successful
func (p *CPUCollector) GetConfigPolicy() (plugin.ConfigPolicy, error) {
	policy := plugin.NewConfigPolicy()
	policy.AddNewStringRule([]string{vendor, fs, Name}, "proc_path", false, plugin.SetDefaultString(defaultProcPath))

	return *policy, nil
}

// GetMetricTypes returns list of available metric types
// It returns error in case retrieval was not successful
func (p *CPUCollector) GetMetricTypes(cfg plugin.Config) ([]plugin.Metric, error) {
	if !p.initialized {
		if err := p.init(cfg); err != nil {
			return nil, err
		}
	}
	if err := getStats(p.proc_path, p.stats, p.prevMetricsSum, p.cpuMetricsNumber, p.snapMetricsNames, p.procStatMetricsNames); err != nil {
		return nil, err
	}
	mts := []plugin.Metric{}

	namespaces := []string{}

	prefix := filepath.Join(vendor, fs, Name)
	for cpu, stats := range p.stats {
		for metric, _ := range stats {
			namespaces = append(namespaces, prefix+"/"+cpu+"/"+metric)
		}
	}

	// List of terminal metric names
	mList := make(map[string]bool)
	for _, namespace := range namespaces {
		namespace = strings.TrimRight(namespace, string(os.PathSeparator))
		mt := plugin.Metric{
			Namespace: plugin.NewNamespace(strings.Split(namespace, string(os.PathSeparator))...)}
		ns := mt.Namespace
		// CPU metric (aka last element in namespace)
		mItem := ns[len(ns)-1]
		// Keep it if not already seen before
		if !mList[mItem.Value] {
			mList[mItem.Value] = true
			mts = append(mts, plugin.Metric{
				Namespace: plugin.NewNamespace(strings.Split(prefix, string(os.PathSeparator))...).
					AddDynamicElement("cpuID", "ID of CPU ('all' for aggregate)").
					AddStaticElement(mItem.Value),
				Description: "dynamic CPU metric: " + mItem.Value,
			})
		}
	}

	return mts, nil
}

// CollectMetrics returns list of requested metric values
// It returns error in case retrieval was not successful
func (p *CPUCollector) CollectMetrics(mts []plugin.Metric) ([]plugin.Metric, error) {
	metrics := []plugin.Metric{}
	if !p.initialized {
		if err := p.init(mts[0].Config); err != nil {
			return nil, err
		}
	}
	if err := getStats(p.proc_path, p.stats, p.prevMetricsSum, p.cpuMetricsNumber, p.snapMetricsNames, p.procStatMetricsNames); err != nil {
		return nil, err
	}
	ts := time.Now()
	for _, mt := range mts {
		ns := mt.Namespace
		if len(ns) != maxNamespaceSize {
			return nil, fmt.Errorf("Incorrect namespace length (len = %d)", len(ns))
		}
		if ns[len(ns)-2].Value == "*" {
			for cpuId, cpuStats := range p.stats {
				for k, v := range cpuStats {
					if ns[len(ns)-1].Value == k && v != nil {
						ns1 := make([]plugin.NamespaceElement, len(ns))
						copy(ns1, ns)
						ns1[len(ns)-2].Value = cpuId
						metric := plugin.Metric{
							Namespace: ns1,
							Data:      v,
							Timestamp: ts,
							Version:   Version,
						}
						metrics = append(metrics, metric)
					}
				}
			}
		} else {
			val, err := getMapValueByNamespace(p.stats[ns.Strings()[3]], ns.Strings()[4:])
			if err != nil {
				return metrics, err
			}
			metric := plugin.Metric{
				Namespace: ns,
				Data:      val,
				Timestamp: ts,
				Version:   Version,
			}
			metrics = append(metrics, metric)
		}
	}
	return metrics, nil
}

func (p *CPUCollector) init(cfg plugin.Config) error {
	procPath, err := cfg.GetString("proc_path")
	if err == nil {
		// change default if proc_path supplied
		p.proc_path = procPath + "/stat"
	}

	fh, err := os.Open(p.proc_path)
	if err != nil {
		return err
	}
	defer fh.Close()

	var procStatMetricsNumber int
	p.cpuMetricsNumber, procStatMetricsNumber, err = getInitialProcStatData(p.proc_path)
	if err != nil {
		return err
	}

	// initialize metric names arrays
	p.procStatMetricsNames = []string{userProcStat, niceProcStat, systemProcStat, idleProcStat,
		iowaitProcStat, irqProcStat, softirqProcStat, stealProcStat, guestProcStat, guestNiceProcStat}[0:procStatMetricsNumber]
	snapSpecificMetricsNames := []string{activeProcStat, utilizationProcStat}

	// build snapMetricsNames to support different kernels
	// var snapMetricsNames []string
	p.snapMetricsNames = append(p.snapMetricsNames, p.procStatMetricsNames...)
	p.snapMetricsNames = append(p.snapMetricsNames, snapSpecificMetricsNames...)
	p.stats = make(map[string]map[string]interface{})
	p.prevMetricsSum = make(map[string]float64)
	p.initialized = true
	return nil
}

// getStats gets metrics from /proc/stat output and calculates snap specific metrics
func getStats(path string, stats map[string]map[string]interface{}, prevMetricsSum map[string]float64, cpuMetricsNumber int, snapMetricsNames []string, procStatMetricsNames []string) (err error) {
	fh, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fh.Close()

	scanner := bufio.NewScanner(fh)
	for i := 0; i < cpuMetricsNumber; i++ {
		scanErr := scanner.Scan()
		if !scanErr {
			return fmt.Errorf("Wrong %s format", path)
		}
		fields := strings.Fields(scanner.Text())

		if len(fields) < 2 {
			return fmt.Errorf("Wrong %s format", path)
		}

		cpuID := strings.TrimSpace(fields[0])
		if cpuID == cpuStr {
			if err != nil {
				return err
			}
			cpuID = allCPU //change CPU identifier for aggregation metrics
		} else {
			cpuID = strings.TrimPrefix(cpuID, cpuStr) //get number from CPU indentifier, for example if CPU identifier is cpu42 then 42 is get
		}
		metrics := fields[1:]

		if len(metrics) != len(procStatMetricsNames) {
			return fmt.Errorf("Wrong data length. Expected {%d} is {%d}",
				len(procStatMetricsNames), len(metrics))
		}

		//sum of new data in line
		currDataSum, err := strTabSum(metrics)
		if err != nil {
			return err
		}

		metricStats := make(map[string]interface{})
		for j := range snapMetricsNames {

			metricName := snapMetricsNames[j]
			var currVal float64
			//data collecting, there is an assumption that firstly metrics from /proc/stat/
			//are gathered then snap specific metrics (e.g. active and utilization are calculated)
			if metricName == activeProcStat {
				idleVal, err := getMapFloatValueByNamespace(metricStats,
					[]string{getNamespaceMetricPart(idleProcStat, jiffiesRepresentationType)})
				if err != nil {
					return err
				}
				currVal = currDataSum - idleVal
			} else if metricName == utilizationProcStat {
				nonActiveVal, err := getMapFloatValueByNamespace(metricStats,
					[]string{getNamespaceMetricPart(idleProcStat, jiffiesRepresentationType)})
				if err != nil {
					return err
				}

				currVal = currDataSum - nonActiveVal

				nonActiveVal, err = getMapFloatValueByNamespace(metricStats,
					[]string{getNamespaceMetricPart(iowaitProcStat, jiffiesRepresentationType)})
				if err != nil {
					return err
				}

				currVal = currVal - nonActiveVal
			} else {
				currVal, err = strconv.ParseFloat(metrics[j], 64)
				if err != nil {
					return err
				}
			}

			metricStats[getNamespaceMetricPart(metricName, percentageRepresentationType)] = nil

			if mapKeyExists(cpuID, prevMetricsSum) {
				diffSum := currDataSum - prevMetricsSum[cpuID]
				if diffSum > 0 {
					prevVal, err := getMapFloatValueByNamespace(stats[cpuID],
						[]string{getNamespaceMetricPart(metricName, jiffiesRepresentationType)})
					if err != nil {
						return err
					}

					if percVal := float64(100 * (currVal - prevVal) / diffSum); percVal < 0 {
						fmt.Fprintf(os.Stderr, "Percentage value of %v could not be calculated due to invalid data reported by /proc/stat\n", getNamespaceMetricPart(metricName, percentageRepresentationType))
					} else {
						metricStats[getNamespaceMetricPart(metricName, percentageRepresentationType)] = percVal
					}
				} else {
					fmt.Fprintf(os.Stderr, "Percentage value of %v could not be calculated due to invalid data reported by /proc/stat\n", getNamespaceMetricPart(metricName, percentageRepresentationType))
				}
			}
			metricStats[getNamespaceMetricPart(metricName, jiffiesRepresentationType)] = currVal
		}
		stats[cpuID] = metricStats
		prevMetricsSum[cpuID] = currDataSum
	}
	return nil
}

// getNamespaceMetricPart builds part of namespace specific for metric and representation type
func getNamespaceMetricPart(metricName string, representationType string) (s string) {
	s = metricName + "_" + representationType
	return s
}

// mapKeyExists checks if element with given key exists in map
func mapKeyExists(key string, m map[string]float64) bool {
	var ret = false
	if _, ok := m[key]; ok {
		ret = true
	}
	return ret
}

// strTabSum adds string data as float
func strTabSum(metrics []string) (sum float64, err error) {
	sum = 0
	for i := range metrics {
		val, err := strconv.ParseFloat(metrics[i], 64)
		if err != nil {
			return sum, err
		}
		sum += val
	}
	return sum, err
}

// getMapFloatValueByNamespace gets value as float from map by namespace given in array of strings
func getMapFloatValueByNamespace(m map[string]interface{}, ns []string) (val float64, err error) {
	var interfaceVal interface{}
	interfaceVal, err = getMapValueByNamespace(m, ns)
	if err != nil {
		return val, err
	}
	var errBool bool
	val, errBool = interfaceVal.(float64)
	if !errBool {
		return val, fmt.Errorf("Parsing error")
	}
	return val, err
}

// getMapValueByNamespace gets value from map by namespace given in array of strings
func getMapValueByNamespace(m map[string]interface{}, ns []string) (val interface{}, err error) {
	if m == nil {
		return nil, fmt.Errorf("Invalid (nil) value of argument m")
	}

	if len(ns) == 0 {
		return nil, fmt.Errorf("Namespace length equal to zero")
	}

	current := ns[0]

	if len(ns) == 1 {
		if val, ok := m[current]; ok {
			return val, err
		}
		return val, fmt.Errorf("Key does not exist in map {key %s}", current)
	}

	if v, ok := m[current].(map[string]interface{}); ok {
		val, err = getMapValueByNamespace(v, ns[1:])
		return val, err
	}
	return val, err
}

// getInitialProcStatData gets number of CPUs and number of metrics available in /proc/stat output
func getInitialProcStatData(path string) (cpuMetricsNumber int, procStatMetricNumber int, err error) {
	fh, err := os.Open(path)
	if err != nil {
		return cpuMetricsNumber, procStatMetricNumber, err
	}
	defer fh.Close()

	scanner := bufio.NewScanner(fh)
	// read the first line to start the loop
	scanErr := scanner.Scan()
	if scanErr != true {
		return cpuMetricsNumber, procStatMetricNumber, fmt.Errorf("Cannot read from %s", path)
	}

	procStatLine := strings.Fields(scanner.Text())
	procStatMetricNumber = len(procStatLine) - 1 //without cpu identifier (e.g. cpu or cpu0)

	cpuMetricsNumber = 0
	for strings.Contains(procStatLine[0], cpuStr) {
		// check line length, compare current length with length of the first line
		// length of the first line = len(procStatMetricNumber) + CPU identifier
		if len(procStatLine) != (procStatMetricNumber + 1) {
			return cpuMetricsNumber, procStatMetricNumber, fmt.Errorf("Incorrect %s output", path)
		}
		cpuMetricsNumber++
		// read the next line to be able to check loop condition
		scanErr = scanner.Scan()
		if scanErr != true {
			break // no more lines to read
		}
		procStatLine = strings.Fields(scanner.Text())
	}
	return cpuMetricsNumber, procStatMetricNumber, err
}
