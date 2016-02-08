# snap collector plugin - cpu
This plugin collects metrics from /proc/stat kernel information about the amount of time, 
measured in units of USER_HZ (1/100ths of a second on most architectures, use
sysconf(_SC_CLK_TCK) to obtain the right value), that the system spent in various states. 

It's used in the [snap framework](http://github.com:intelsdi-x/snap).

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Installation](#installation)
  * [Configuration and Usage](configuration-and-usage)
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
* [golang 1.4.3+](https://golang.org/dl/)

### Operating systems
All OSs currently supported by snap:
* Linux/amd64

### Installation
#### Download interface plugin binary:
You can get the pre-built binaries for your OS and architecture at snap's [GitHub Releases](https://github.com/intelsdi-x/snap/releases) page.

#### To build the plugin binary:
Fork https://github.com/intelsdi-x/snap-plugin-collector-cpu  
Clone repo into `$GOPATH/src/github.com/intelsdi-x/`:

```
$ git clone https://github.com/<yourGithubID>/snap-plugin-collector-cpu.git
```

Build the plugin by running make within the cloned repo:
```
$ make
```
This builds the plugin in `/build/rootfs/`

### Configuration and Usage
* Set up the [snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started)
* Ensure `$SNAP_PATH` is exported  
`export SNAP_PATH=$GOPATH/src/github.com/intelsdi-x/snap/build`

## Documentation

### Collected Metrics
This plugin has the ability to gather the following metrics:

Namespace | Description (optional)
----------|-----------------------
/intel/procfs/cpu/all/user_jiffies	| The amount of time spent in user mode by all CPUs
/intel/procfs/cpu/all/nice_jiffies	| The amount of time spent in user mode with low priority by all CPUs
/intel/procfs/cpu/all/system_jiffies	| The amount of time spent in system mode by all CPUs
/intel/procfs/cpu/all/idle_jiffies	|	The amount of time spent in the idle task by all CPUs
/intel/procfs/cpu/all/iowait_jiffies	| The amount of time spent waiting for I/O to complete by all CPUs
/intel/procfs/cpu/all/irq_jiffies	|	The amount of time servicing interrupts by all CPUs
/intel/procfs/cpu/all/softirq_jiffies	|	The amount of time servicing softirqs by all CPUs
/intel/procfs/cpu/all/steal_jiffies	|	The amount of stolen time, which is the time spent in other operating systems when running in a virtualized environment by all CPUs
/intel/procfs/cpu/all/guest_jiffies	|	The amount of time spent running a virtual CPU for guest operating systems under the control of the Linux kernel by all CPUs
/intel/procfs/cpu/all/guest_nice_jiffies| The amount of time spent running a niced guest (virtual CPU for guest operating systems under the control of the Linux kernel) by all CPUs
/intel/procfs/cpu/all/active_jiffies	| The amount of time spend in non idle state by all CPUs
/intel/procfs/cpu/all/utilization_jiffies	| The amount of time spend in non idle and non iowait states by all CPUs
/intel/procfs/cpu/all/user_percentage	| The percent of time spent in user mode by all CPUs
/intel/procfs/cpu/all/nice_percentage	| The percent of time spent in user mode with low priority by all CPUs
/intel/procfs/cpu/all/system_percentage	| The percent of time spent in system mode by all CPUs
/intel/procfs/cpu/all/idle_percentage	|	The percent of time spent in the idle task by all CPUs
/intel/procfs/cpu/all/iowait_percentage	| The percent of time spent waiting for I/O to complete by all CPUs
/intel/procfs/cpu/all/irq_percentage	|	The percent of time servicing interrupts by all CPUs
/intel/procfs/cpu/all/softirq_percentage	|	The percent of time servicing softirqs by all CPUs
/intel/procfs/cpu/all/steal_percentage	|	The percent of stolen time, which is the time spent in other operating systems when running in a virtualized environment by all CPUs
/intel/procfs/cpu/all/guest_percentage	|	The percent of time spent running a virtual CPU for guest operating systems under the control of the Linux kernel by all CPUs
/intel/procfs/cpu/all/guest_nice_percentage| The percent of time spent running a niced guest (virtual CPU for guest operating systems under the control of the Linux kernel) by all CPUs
/intel/procfs/cpu/all/active_percentage	| The percent of time spend in non idle state by all CPUs
/intel/procfs/cpu/all/utilization_percentage	| The percent of time spend in non idle and non iowait states by all CPUs
/intel/procfs/cpu/\<CPU_number\>/user_jiffies	| The amount of time spent in user mode by CPU with given identifier
/intel/procfs/cpu/\<CPU_number\>/nice_jiffies	| The amount of time spent in user mode with low priority by CPU with given identifier
/intel/procfs/cpu/\<CPU_number\>/system_jiffies	| The amount of time spent in system mode by CPU with given identifier
/intel/procfs/cpu/\<CPU_number\>/idle_jiffies	|	The amount of time spent in the idle task by CPU with given identifier
/intel/procfs/cpu/\<CPU_number\>/iowait_jiffies	| The amount of time spent waiting for I/O to complete by CPU with given identifier
/intel/procfs/cpu/\<CPU_number\>/irq_jiffies	|	The amount of time servicing interrupts by CPU with given identifier
/intel/procfs/cpu/\<CPU_number\>/softirq_jiffies	|	The amount of time servicing softirqs by CPU with given identifier
/intel/procfs/cpu/\<CPU_number\>/steal_jiffies	|	The amount of stolen time, which is the time spent in other operating systems when running in a virtualized environment by CPU with given identifier
/intel/procfs/cpu/\<CPU_number\>/guest_jiffies	|	The amount of time spent running a virtual CPU for guest operating systems under the control of the Linux kernel by CPU with given identifier
/intel/procfs/cpu/\<CPU_number\>/guest_nice_jiffies| The amount of time spent running a niced guest (virtual CPU for guest operating systems under the control of the Linux kernel) by CPU with given identifier
/intel/procfs/cpu/\<CPU_number\>/active_jiffies	| The amount of time spend in non idle state by CPU with given identifier
/intel/procfs/cpu/\<CPU_number\>/utilization_jiffies	| The amount of time spend in non idle and non iowait states by CPU with given identifier
/intel/procfs/cpu/\<CPU_number\>/user_percentage	| The percent of time spent in user mode by CPU with given identifier
/intel/procfs/cpu/\<CPU_number\>/nice_percentage	| The percent of time spent in user mode with low priority by CPU with given identifier
/intel/procfs/cpu/\<CPU_number\>/system_percentage	| The percent of time spent in system mode by CPU with given identifier
/intel/procfs/cpu/\<CPU_number\>/idle_percentage	|	The percent of time spent in the idle task by all CPUs
/intel/procfs/cpu/\<CPU_number\>/iowait_percentage	| The percent of time spent waiting for I/O to complete by CPU with given identifier
/intel/procfs/cpu/\<CPU_number\>/irq_percentage	|	The percent of time servicing interrupts by CPU with given identifier
/intel/procfs/cpu/\<CPU_number\>/softirq_percentage	|	The percent of time servicing softirqs by CPU with given identifier
/intel/procfs/cpu/\<CPU_number\>/steal_percentage	|	The percent of stolen time, which is the time spent in other operating systems when running in a virtualized environment by CPU with given identifier
/intel/procfs/cpu/\<CPU_number\>/guest_percentage	|	The percent of time spent running a virtual CPU for guest operating systems under the control of the Linux kernel by all CPUs
/intel/procfs/cpu/\<CPU_number\>/guest_nice_percentage| The percent of time spent running a niced guest (virtual CPU for guest operating systems under the control of the Linux kernel) by CPU with given identifier
/intel/procfs/cpu/\<CPU_number\>/active_percentage	| The percent of time spend in non idle state by CPU with given identifier
/intel/procfs/cpu/\<CPU_number\>/utilization_percentage	| The percent of time spend in non idle and non iowait states by CPU with given identifier

### Examples
Example running interface, passthru processor, and writing data to a file.

This is done from the snap directory.

In one terminal window, open the snap daemon (in this case with logging set to 1 and trust disabled):
```
$ $SNAP_PATH/bin/snapd -l 1 -t 0
```

In another terminal window:
Load interface plugin
```
$ $SNAP_PATH/bin/snapctl plugin load snap-plugin-collector-cpu
```
See available metrics for your system
```
$ $SNAP_PATH/bin/snapctl metric list
```

Create a task manifest file (e.g. `cpu-file.json`):
    
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
	"/intel/procfs/cpu/all/user_jiffies" : {},
	"/intel/procfs/cpu/all/nice_jiffies" : {},
	"/intel/procfs/cpu/all/system_jiffies" : {},
	"/intel/procfs/cpu/all/idle_jiffies" : {},
	"/intel/procfs/cpu/all/iowait_jiffies" : {},
	"/intel/procfs/cpu/all/irq_jiffies" : {},
	"/intel/procfs/cpu/all/softirq_jiffies" : {},
	"/intel/procfs/cpu/all/steal_jiffies" : {},
	"/intel/procfs/cpu/all/guest_jiffies" : {},
	"/intel/procfs/cpu/all/guest_nice_jiffies" : {},
	"/intel/procfs/cpu/all/active_jiffies" : {},
	"/intel/procfs/cpu/all/utilization_jiffies" : {},
	"/intel/procfs/cpu/all/user_percentage" : {},
	"/intel/procfs/cpu/all/nice_percentage" : {},
	"/intel/procfs/cpu/all/system_percentage" : {},
	"/intel/procfs/cpu/all/idle_percentage" : {},
	"/intel/procfs/cpu/all/iowait_percentage" : {},
	"/intel/procfs/cpu/all/irq_percentage" : {},
	"/intel/procfs/cpu/all/softirq_percentage" : {},
	"/intel/procfs/cpu/all/steal_percentage" : {},
	"/intel/procfs/cpu/all/guest_percentage" : {},
	"/intel/procfs/cpu/all/guest_nice_percentage" : {},
	"/intel/procfs/cpu/all/active_percentage" : {},
	"/intel/procfs/cpu/all/utilization_percentage" : {},
	"/intel/procfs/cpu/0/user_jiffies" : {},
	"/intel/procfs/cpu/0/nice_jiffies" : {},
	"/intel/procfs/cpu/0/system_jiffies" : {},
	"/intel/procfs/cpu/0/idle_jiffies" : {},
	"/intel/procfs/cpu/0/iowait_jiffies" : {},
	"/intel/procfs/cpu/0/irq_jiffies" : {},
	"/intel/procfs/cpu/0/softirq_jiffies" : {},
	"/intel/procfs/cpu/0/steal_jiffies" : {},
	"/intel/procfs/cpu/0/guest_jiffies" : {},
	"/intel/procfs/cpu/0/guest_nice_jiffies" : {},
	"/intel/procfs/cpu/0/active_jiffies" : {},
	"/intel/procfs/cpu/0/utilization_jiffies" : {},
	"/intel/procfs/cpu/0/user_percentage" : {},
	"/intel/procfs/cpu/0/nice_percentage" : {},
	"/intel/procfs/cpu/0/system_percentage" : {},
	"/intel/procfs/cpu/0/idle_percentage" : {},
	"/intel/procfs/cpu/0/iowait_percentage" : {},
	"/intel/procfs/cpu/0/irq_percentage" : {},
	"/intel/procfs/cpu/0/softirq_percentage" : {},
	"/intel/procfs/cpu/0/steal_percentage" : {},
	"/intel/procfs/cpu/0/guest_percentage" : {},
	"/intel/procfs/cpu/0/guest_nice_percentage" : {},
	"/intel/procfs/cpu/0/active_percentage" : {},
	"/intel/procfs/cpu/0/utilization_percentage" : {}
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
$ $SNAP_PATH/bin/snapctl plugin load build/plugin/snap-processor-passthru
Plugin loaded
Name: passthru
Version: 1
Type: processor
Signed: false
Loaded Time: Fri, 20 Nov 2015 11:44:03 PST
```

Load file plugin for publishing:
```
$ $SNAP_PATH/bin/snapctl plugin load build/plugin/snap-publisher-file
Plugin loaded
Name: file
Version: 3
Type: publisher
Signed: false
Loaded Time: Fri, 20 Nov 2015 11:41:39 PST
```

Create task:
```
$ $SNAP_PATH/bin/snapctl task create -t examples/tasks/cpu-file.json
Using task manifest to create task
Task created
ID: 02dd7ff4-8106-47e9-8b86-70067cd0a850
Name: Task-02dd7ff4-8106-47e9-8b86-70067cd0a850
State: Running
```

Stop task:
```
$ $SNAP_PATH/bin/snapctl task stop 02dd7ff4-8106-47e9-8b86-70067cd0a850
Task stopped:
ID: 02dd7ff4-8106-47e9-8b86-70067cd0a850
```

### Roadmap
There isn't a current roadmap for this plugin, but it is in active development. As we launch this plugin, we do not have any outstanding requirements for the next release. If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-plugin-collector-cpu/issues/new) and/or submit a [pull request](https://github.com/intelsdi-x/snap-plugin-collector-cpu/pulls).

## Community Support
This repository is one of **many** plugins in **snap**, a powerful telemetry framework. See the full project at http://github.com/intelsdi-x/snap To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support)

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

## License
[snap](http://github.com:intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
* Author: [@katarzyna-z](https://github.com/katarzyna-z)

And **thank you!** Your contribution, through code and participation, is incredibly important to us.
