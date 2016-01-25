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

package cpu

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	str "github.com/intelsdi-x/snap-plugin-utilities/strings"
	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
)

const (
	// VENDOR namespace part
	VENDOR = "intel"
	// FS namespace part
	FS = "procfs"
	// PLUGIN namespace part
	PLUGIN = "cpu"
	// VERSION of load info plugin
	VERSION = 1
	//TYPE of plugin
	TYPE = plugin.CollectorPluginType
)

var cpuInfo = "/proc/stat"
var cpuData []cpu_data_struct
var cpuDataAggregatedThroughMetrics map[string]float64
var cpuDataAggregateThroughStates [] float64
var procStatMetricsNames []string
var snapMetricsNames []string


// make sure that we actually satisify requierd interface
var _ plugin.CollectorPlugin = (*CpuPlugin)(nil)

type CpuPlugin struct {
	stats[]interface{}
	host  string
	cpuMetricsNumber  int // number of cpu + "all" metric
}

func New() *CpuPlugin {
	
	fh, err := os.Open(cpuInfo)
	if err != nil {
		return nil
	}
	defer fh.Close()
	
	cpu, err := getCPUs()
	cpuMetricsNumber := cpu + 1
	
	if err != nil {
		return nil
	}
	
	initialize()
	cpuData = createCPUDataTable(cpuMetricsNumber)
	
	host, err := os.Hostname()
	if err != nil {
		host = "localhost"
	}
	
	p := &CpuPlugin{stats:[]interface{}{}, host:host, cpuMetricsNumber:cpuMetricsNumber}
	return p
}


func initialize() {
	procStatMetricsNames = []string{"user", "nice", "system", "idle", "iowait", "irq", "softirq", "steal", "guest", "guest_nice"}
	snapMetricsNames =  []string{"user", "nice", "system", "idle", "iowait", "irq", "softirq", "steal", "guest", "guest_nice", "active"}
	createCPUDataAggregatedThroughMetrics()
}

func createCPUDataAggregatedThroughMetrics() {
	cpuDataAggregatedThroughMetrics = make(map[string]float64)
	for metric := range snapMetricsNames {
		cpuDataAggregatedThroughMetrics[snapMetricsNames[metric]] = 0
	}
}

func cleanCPUDataAggregatedThroughMetrics() {
	for k, _ := range cpuDataAggregatedThroughMetrics {
		cpuDataAggregatedThroughMetrics[k] = 0
	}
}

func getCPUs() (int, error) {
	out, err := exec.Command("lscpu", "-p").Output()
	if err != nil {
		return -1, err
	}

	lines := strings.Split(string(out), "\n")
	lines = str.Filter(lines, func(s string) bool {
		return s != ""
	})
	last := lines[len(lines)-1]
	cpus, err := strconv.Atoi(strings.Split(last, ",")[0])

	if err != nil {
		return -1, err
	}

	return cpus + 1, nil
}

func mapKeyExist( key string, m map[string]float64) (bool) {
	_, ok := m[key]
    if ok {
	return true
   } else {
	return false
   }
}

func sumMapElements(a map[string]float64) (float64){
	var sum float64 = 0
	for _,val := range a {
		sum = sum + val
	}
	return sum
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

type cpu_data_struct struct{
	statVal map[string]float64
	rateVal map[string]float64
	time int64
	rateSum int64
}

func newCPUDataStruct() (cpu_data_struct) {
    var cpu cpu_data_struct
    cpu.statVal = make(map[string]float64)
	cpu.rateVal = make(map[string]float64)
	cpu.time = 0
	cpu.rateSum = 0
	return cpu
}

func newCPUData() ([]cpu_data_struct) {
	return make([]cpu_data_struct,0)
}

func createCPUDataTable(n int) ([]cpu_data_struct) {
	cpuTable:= make([]cpu_data_struct, 0)
	for i := 0; i < n; i++ {
		cpuTable = append(cpuTable, newCPUDataStruct())
	}
	return cpuTable
}

func updateCPUMetricsVal(cpuNumber int, now int64, metricName string, metricVal float64) {
	if (mapKeyExist(metricName, cpuData[cpuNumber].statVal)) {
		//if value of metrics exist then we could calculate rate value
		cpuData[cpuNumber].rateVal[metricName] =  (metricVal - cpuData[cpuNumber].statVal[metricName]) / float64( now - cpuData[cpuNumber].time )
	} else {
		//previous value of metric does not exist - we can not calculate rate value
		cpuData[cpuNumber].rateVal[metricName] = 0
	}

	cpuData[cpuNumber].statVal[metricName] = metricVal
	//get race value aggreagate through states for all cpu
	if(cpuNumber != 0) {
		cpuDataAggregatedThroughMetrics[metricName] = cpuDataAggregatedThroughMetrics[metricName] + cpuData[cpuNumber].rateVal[metricName]
	}
}

func updateCPUMetricsTime(cpuNumber int, now int64) {
	cpuData[cpuNumber].time = now
}

func getActiveStatMetric(cpuNumber int) (float64) {
	var idleStat float64 = 0

	if mapKeyExist("idle", cpuData[cpuNumber].statVal) {
		idleStat = cpuData[cpuNumber].statVal["idle"]
	}

	return  (sumMapElements( cpuData[cpuNumber].statVal) -  idleStat)
}

//getProcStatMetrics returns /proc/stat output as array of lines,
func getProcStatMetrics() (int, [][]string, error) {
	var procStatMetricsLines = make([][]string,0)
	var procStatMetricsNumber int = 0

	fh, err := os.Open(cpuInfo)
	if err != nil {
		return 0, nil, err
	}
	defer fh.Close()

	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		procStatMetricsLines = append(procStatMetricsLines,  strings.Fields(scanner.Text()))
	}

	if len(procStatMetricsLines) != 0 {
		procStatMetricsNumber =  len(procStatMetricsLines[0]) - 1
	}
	return procStatMetricsNumber, procStatMetricsLines, nil
}

func getCPUMetrics() {
	var metricName string
	var metricVal float64
	var err error
	var now int64 = 0
	var procStatMetricsLines [][]string
	var procStatMetricsNumber int
	
	cleanCPUDataAggregatedThroughMetrics()
	
	now = makeTimestamp()
	procStatMetricsNumber, procStatMetricsLines, err = getProcStatMetrics()

	for i:=0; i< len(cpuData); i++ {
		procStatMetricsLineTab := procStatMetricsLines[i]
		fmt.Println(procStatMetricsLineTab)

		//get metrics from /proc/stat output
		for j := 0; j < procStatMetricsNumber; j++ {
			metricName = procStatMetricsNames[j]
			metricVal, err =  strconv.ParseFloat(procStatMetricsLineTab[j], 64)
			if err != nil {
				continue
			}
			updateCPUMetricsVal(i, now, metricName, metricVal)
		}

		//get "active" metrics, "active" means not idle
		activeStat := getActiveStatMetric(i)
		updateCPUMetricsVal(i, now, "active", activeStat)
		updateCPUMetricsTime(i, now)
	}
}

func getCPUId(i int) (string) {
	if i == 0 {
		return "all"
	} else {
		return strconv.Itoa(i-1)
	}
}

func getCPUMetricId(i int) (string){
	if i >= 0 && i <= len(snapMetricsNames) {
		return snapMetricsNames[i]
	} else {
		return "" //TODO error
	}
}

func getRepresentationTypes(cpu_id int) ([]string) {
	// 0 means all
	if cpu_id == 0 {
		return []string{"percentage"}
	} else {
		return []string{"", "percentage"}
	}
}

func getNamespaceTab(cpu_id string, metrics_id string, representation_type string) ([]string) {
	var namespace_tab []string
	if representation_type == "" {
		namespace_tab = []string{VENDOR, FS, PLUGIN, cpu_id, metrics_id}
	} else {
		namespace_tab = []string{VENDOR, FS, PLUGIN, cpu_id, metrics_id, representation_type}
	}
	return namespace_tab
}

func getCPUTabId(cpuId string) (int) {
	if cpuId == "all" {
		return 0
	} else {
		cpuTabId, err := strconv.Atoi(cpuId)
		if err != nil {
			return -1 //TODO: error handling
		} 
		cpuTabId = cpuTabId + 1
		return cpuTabId
	}
}

func getJiffiesCPUMetric(cpuId string, metricName string,  representationType string)  (float64) {
	//TODO: check if element exist
	cpuTabId := getCPUTabId(cpuId)
	return cpuData[cpuTabId].statVal[metricName]
}

func calculateRatePercentageValue(m map[string]float64, metricName string) (float64) {
	var percentMetricVal float64 = 0
	if(mapKeyExist(metricName,m) && mapKeyExist("active",m) && mapKeyExist("idle",m)) {
		ratesSum := m["active"] + m["idle"]
		if ratesSum != 0 {
			percentMetricVal = 100 * m[metricName] / ratesSum
		}
	}
	//TODO: error handling
	return percentMetricVal
}

func getPercentageCPUMetric(cpuId string, metricName string,  representationType string) (float64) {
	//TODO: check if element exist
	cpuTabId := getCPUTabId(cpuId)
	percentMetricVal := calculateRatePercentageValue(cpuData[cpuTabId].rateVal, metricName)
	return percentMetricVal
}

func getAllPercentageCPUMetric(cpuId string, metricName string,  representationType string) (float64) {
	cpuTabId := getCPUTabId(cpuId)
	percentMetricVal := calculateRatePercentageValue(cpuData[cpuTabId].rateVal, metricName)
	return percentMetricVal
}

func getOneCPUMetric(cpuId string, snapMetricsName string,  representationType string) (float64) {
	var val float64 = 0
	switch(representationType) {
		default:
			val = getJiffiesCPUMetric(cpuId,  snapMetricsName,  representationType)
		case "percentage":
			val = getPercentageCPUMetric(cpuId,  snapMetricsName,  representationType)
	}
	return val
}

func getAllCPUMetric(cpuId string, snapMetricsName string,  representationType string) (float64) {
	var val float64 = 0
	switch(representationType) {
		default:
			val = -1 //TODO:
		case "percentage":
			val = getAllPercentageCPUMetric(cpuId,  snapMetricsName,  representationType)
	}
	return val
}

func getMetricValue(cpuId string, snapMetricsName,  representationType string) (float64){
	var val float64 = 0
	switch(cpuId) {
		case "all":
			val = getAllCPUMetric(cpuId,  snapMetricsName,  representationType)
		default:
			val = getOneCPUMetric(cpuId, snapMetricsName,  representationType)
	}
	return val
}

func getMetricParameters(ns []string) (string, string, string, error) {
	if len(ns) == 5 || len(ns) == 6 {
		cpuId := ns[3]
		metricId := ns[4]
		representationType := ""
		if len(ns) == 6 {
			representationType = ns[5]
		}
		return cpuId, metricId, representationType, nil
	} else {
		return "", "", "", fmt.Errorf("Incorrect namespace length, (len = %d)", len(ns))
	}
}

// GetMetricTypes returns list of available metric types
// It returns error in case retrieval was not successful
func (p *CpuPlugin) GetMetricTypes(cfg plugin.PluginConfigType) ([]plugin.PluginMetricType, error) {
	mts := []plugin.PluginMetricType{}

	for i :=  0; i < p.cpuMetricsNumber; i++ {
		for j:= 0; j<len(snapMetricsNames); j++{

			representation_types := getRepresentationTypes(i)
			representation_types_number := len(representation_types)

			for k:=0; k < representation_types_number; k++ {
				mt :=  plugin.PluginMetricType{Namespace_:  getNamespaceTab(getCPUId(i), getCPUMetricId(j), representation_types[k])}
				mts = append(mts, mt)
			}
		}
	}
	return mts, nil
}

// CollectMetrics returns list of requested metric values
// It returns error in case retrieval was not successful
func (p *CpuPlugin) CollectMetrics(mts []plugin.PluginMetricType) ([]plugin.PluginMetricType, error) {
	metrics := []plugin.PluginMetricType{}
	getCPUMetrics()
	
	for _, mt := range mts {
		ns := mt.Namespace()
		cpuId, metricId, representationType, err := getMetricParameters(ns)
		if err != nil {
			return nil, err
		}
		val := getMetricValue(cpuId, metricId, representationType)

		mt := plugin.PluginMetricType{
				Namespace_: ns,
				Data_: val,
				Source_: p.host,
				Timestamp_: time.Now(),
			}
		metrics = append(metrics, mt)	
	}
	return metrics, nil
}

// GetConfigPolicy returns config policy
// It returns error in case retrieval was not successful
func (p *CpuPlugin) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	return cpolicy.New(), nil
}

//Meta returns meta data for testing
func Meta() *plugin.PluginMeta {
	return plugin.NewPluginMeta(
			PLUGIN,
			VERSION,
			plugin.CollectorPluginType,
			[]string{},
			[]string{plugin.SnapGOBContentType},
			plugin.ConcurrencyCount(1))
}

