# snap collector plugin - cpu

## Collected Metrics

Note that in the following table, the dynamic component of the namespace (*)
is either the \<CPU ID/number\> or 'all' when the metric is aggregated across all CPUs

This plugin has the ability to gather the following metrics:

Namespace 					| Description
------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------
/intel/procfs/cpu/*/user_jiffies		| The amount of time spent in user mode by CPU with given identifier
/intel/procfs/cpu/*/nice_jiffies		| The amount of time spent in user mode with low priority by CPU with given identifier
/intel/procfs/cpu/*/system_jiffies		| The amount of time spent in system mode by CPU with given identifier
/intel/procfs/cpu/*/idle_jiffies		| The amount of time spent in the idle task by CPU with given identifier
/intel/procfs/cpu/*/iowait_jiffies		| The amount of time spent waiting for I/O to complete by CPU with given identifier
/intel/procfs/cpu/*/irq_jiffies			| The amount of time servicing interrupts by CPU with given identifier
/intel/procfs/cpu/*/softirq_jiffies		| The amount of time servicing softirqs by CPU with given identifier
/intel/procfs/cpu/*/steal_jiffies		| The amount of stolen time, which is the time spent in other operating systems when running in a virtualized environment by CPU with given identifier
/intel/procfs/cpu/*/guest_jiffies		| The amount of time spent running a virtual CPU for guest operating systems under the control of the Linux kernel by CPU with given identifier
/intel/procfs/cpu/*/guest_nice_jiffies		| The amount of time spent running a niced guest (virtual CPU for guest operating systems under the control of the Linux kernel) by CPU with given identifier
/intel/procfs/cpu/*/active_jiffies		| The amount of time spend in non idle state by CPU with given identifier
/intel/procfs/cpu/*/utilization_jiffies		| The amount of time spend in non idle and non iowait states by CPU with given identifier
/intel/procfs/cpu/*/user_percentage		| The percent of time spent in user mode by CPU with given identifier
/intel/procfs/cpu/*/nice_percentage		| The percent of time spent in user mode with low priority by CPU with given identifier
/intel/procfs/cpu/*/system_percentage		| The percent of time spent in system mode by CPU with given identifier
/intel/procfs/cpu/*/idle_percentage		| The percent of time spent in the idle task by all CPUs
/intel/procfs/cpu/*/iowait_percentage		| The percent of time spent waiting for I/O to complete by CPU with given identifier
/intel/procfs/cpu/*/irq_percentage		| The percent of time servicing interrupts by CPU with given identifier
/intel/procfs/cpu/*/softirq_percentage		| The percent of time servicing softirqs by CPU with given identifier
/intel/procfs/cpu/*/steal_percentage		| The percent of stolen time, which is the time spent in other operating systems when running in a virtualized environment by CPU with given identifier
/intel/procfs/cpu/*/guest_percentage		| The percent of time spent running a virtual CPU for guest operating systems under the control of the Linux kernel by all CPUs
/intel/procfs/cpu/*/guest_nice_percentage	| The percent of time spent running a niced guest (virtual CPU for guest operating systems under the control of the Linux kernel) by CPU with given identifier
/intel/procfs/cpu/*/active_percentage		| The percent of time spend in non idle state by CPU with given identifier
/intel/procfs/cpu/*/utilization_percentage	| The percent of time spend in non idle and non iowait states by CPU with given identifier

