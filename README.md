# snap collector plugin - cpu

This plugin collects metrics from /proc/stat kernel information about the amount of time, 
measured in units of USER_HZ (1/100ths of a second on most architectures, use
sysconf(_SC_CLK_TCK) to obtain the right value), that the system spent in various states. 

It's used in the [snap framework](http://github.com:intelsdi-x/snap).

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Operating systems](#operating-systems)
  * [Installation](#installation)
  * [Configuration and Usage](#configuration-and-usage)
2. [Documentation](#documentation)
  * [Collected Metrics](#collected-metrics)
  * [Examples](#examples)
  * [Roadmap](#roadmap)
3. [Community Support](#community-support)
4. [Contributing](#contributing)
5. [License](#license-and-authors)
6. [Acknowledgements](#acknowledgements)

## Getting Started
### System Requirements
* [golang 1.5+](https://golang.org/dl/) - needed only for building

### Operating systems
All OSs currently supported by plugin:
* Linux/amd64

### Installation
You can get the pre-built binaries for your OS and architecture at snap's [GitHub Releases](https://github.com/intelsdi-x/snap/releases) page. Download the plugins package from the latest release, unzip and store in a path you want `snapd` to access.

### To build the plugin binary:
Fork https://github.com/intelsdi-x/snap-plugin-collector-cpu
Clone repo into `$GOPATH/src/github.com/intelsdi-x/`:

```
$ git clone https://github.com/<yourGithubID>/snap-plugin-collector-cpu.git
```

Build the plugin by running make within the cloned repo:
```
$ make
```
This builds the plugin in `/build/$GOOS/GOARCH`

### Configuration and Usage
* Set up the [snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started)
* If /proc resides in a different directory, say for example by mounting host /proc inside a container at /hostproc, a proc_path configuration item can be added to snapd global config or as part of the task manifest for the metrics to be collected.

As part of snapd global config

```yaml
---
control:
  plugins:
    collector:
      cpu:
        all:
          proc_path: /hostproc
```

Or as part of the task manifest

```json
{
...
    "workflow": {
        "collect": {
            "metrics": {
	      "/intel/procfs/cpu/*/user_jiffies" : {}
	    },
	    "config": {
	      "/intel/procfs": {
                "proc_path": "/hostproc"
	      }
	    },
	    ...
       },
    },
...
```

* Load the plugin and create a task, see example in [Examples](https://github.com/intelsdi-x/snap-plugin-collector-cpu/blob/master/README.md#examples).

## Documentation
### Collected Metrics
Collected metrics have namespace in following format: `/intel/procfs/cpu/<cpu_identifier>/<metric_name>`.
List of collected metrics in [METRICS.md](https://github.com/intelsdi-x/snap-plugin-collector-cpu/blob/master/METRICS.md)

### Examples
#### Run the example
```bash
./examples/run-cpu-file.sh
```

#### Run the plugin manually
Example running CPU collector plugin, passthru processor plugin, and writing data to a file using file publisher plugin.

Make sure that your `$SNAP_PATH` is set, if not:
```
$ export SNAP_PATH=<snapDirectoryPath>/build
```
Other paths to files should be set according to your configuration, using a file you should indicate where it is located.

In one terminal window, open the snap daemon (in this case with logging set to 1 and trust disabled):
```
$ $SNAP_PATH/bin/snapd -l 1 -t 0
```
In another terminal window:
Load snap-plugin-collector-cpu plugin:
```
$ $SNAP_PATH/bin/snapctl plugin load snap-plugin-collector-cpu
```
See available metrics for your system:
```
$ $SNAP_PATH/bin/snapctl metric list
```
Create a task manifest file ( see [exemplary files] (https://github.com/intelsdi-x/snap-plugin-collector-cpu/blob/master/examples/task/)):
    
```json
{
    "version": 1,
    "schedule": {
        "type": "simple",
        "interval": "1s"
    },
    "workflow": {
        "collect": {
            "metrics": {
	            "/intel/procfs/cpu/*/user_jiffies" : {},
	            "/intel/procfs/cpu/*/nice_jiffies" : {},
	            "/intel/procfs/cpu/*/system_jiffies" : {},
	            "/intel/procfs/cpu/*/idle_jiffies" : {},
	            "/intel/procfs/cpu/*/iowait_jiffies" : {},
	            "/intel/procfs/cpu/*/irq_jiffies" : {},
	            "/intel/procfs/cpu/*/softirq_jiffies" : {},
	            "/intel/procfs/cpu/*/steal_jiffies" : {},
	            "/intel/procfs/cpu/*/guest_jiffies" : {},
	            "/intel/procfs/cpu/*/guest_nice_jiffies" : {},
	            "/intel/procfs/cpu/*/active_jiffies" : {},
	            "/intel/procfs/cpu/*/utilization_jiffies" : {},
	            "/intel/procfs/cpu/*/user_percentage" : {},
	            "/intel/procfs/cpu/*/nice_percentage" : {},
	            "/intel/procfs/cpu/*/system_percentage" : {},
	            "/intel/procfs/cpu/*/idle_percentage" : {},
	            "/intel/procfs/cpu/*/iowait_percentage" : {},
	            "/intel/procfs/cpu/*/irq_percentage" : {},
	            "/intel/procfs/cpu/*/softirq_percentage" : {},
	            "/intel/procfs/cpu/*/steal_percentage" : {},
	            "/intel/procfs/cpu/*/guest_percentage" : {},
	            "/intel/procfs/cpu/*/guest_nice_percentage" : {},
	            "/intel/procfs/cpu/*/active_percentage" : {},
	            "/intel/procfs/cpu/*/utilization_percentage" : {}
            },
            "process": [
                {
                    "plugin_name": "passthru",
                    "process": null,
                    "publish": [
                        {
                            "plugin_name": "file",
                            "config": {
                                "file": "/tmp/published_cpu.log"
                            }
                        }
                    ],
                    "config": null
                }
            ],
            "publish": null
        }
    }
}
```

Load passthru plugin for processing:
```
$ $SNAP_PATH/bin/snapctl plugin load build/plugin/snap-plugin-processor-passthru
Plugin loaded
Name: passthru
Version: 1
Type: processor
Signed: false
Loaded Time: Fri, 26 Aug 2016 12:13:18 CEST
```

Load file plugin for publishing:
```
$ $SNAP_PATH/bin/snapctl plugin load build/plugin/snap-plugin-publisher-file
Plugin loaded
Name: file
Version: 3
Type: publisher
Signed: false
Loaded Time: Fri, 26 Aug 2016 12:13:43 CEST
```

Create a task:
```
$ $SNAP_PATH/bin/snapctl task create -t examples/tasks/cpu-file.json
Using task manifest to create task
Task created
ID: 02dd7ff4-8106-47e9-8b86-70067cd0a850
Name: Task-02dd7ff4-8106-47e9-8b86-70067cd0a850
State: Running
```

Stop previously created task:
```
$ $SNAP_PATH/bin/snapctl task stop 02dd7ff4-8106-47e9-8b86-70067cd0a850
Task stopped:
ID: 02dd7ff4-8106-47e9-8b86-70067cd0a850
```

### Roadmap
There isn't a current roadmap for this plugin, but it is in active development. As we launch this plugin, we do not have any outstanding requirements for the next release. If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-plugin-collector-cpu/issues/new) and/or submit a [pull request](https://github.com/intelsdi-x/snap-plugin-collector-cpu/pulls).

If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-plugin-collector-cpu/issues).

## Community Support
This repository is one of **many** plugins in **snap**, a powerful telemetry framework. The full project is at http://github.com:intelsdi-x/snap.
To reach out on other use cases, visit:
* [snap Gitter channel](https://gitter.im/intelsdi-x/snap)

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

And **thank you!** Your contribution, through code and participation, is incredibly important to us.

## License
[snap](http://github.com:intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
* Author: [Katarzyna Zabrocka](https://github.com/katarzyna-z)
