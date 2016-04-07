# snap collector plugin - cpu

## Collected Metrics

This plugin has the ability to gather the following metrics:

Namespace | Description
--------------------|------------------------------
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