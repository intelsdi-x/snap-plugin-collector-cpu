# Example tasks

[This](tasks/cpu-file.json) example task will publish metrics to file 
from the **cpu** plugin.  

## Running the example

### Requirements
 * `docker` and `docker-compose` are **installed** and **configured** 

Running the sample is as *easy* as running the script `./run-cpu-file.sh`. 

## Files

- [run-cpu-file.sh](run-cpu-file.sh) 
    - The example is launched with this script     
- [tasks/cpu-file.json](tasks/cpu-file.json)
    - Snap task definition
- [docker-compose.yml](docker-compose.yml)
    - A docker compose file which defines two linked containers
        - "runner" is the container where snapteld is run from.  You will be dumped 
        into a shell in this container after running 
        [run-cpu-file.sh](run-cpu-file.sh).  Exiting the shell will 
        trigger cleaning up the containers used in the example.
- [cpu-file.sh](cpu-file.sh)
    - Downloads `snapteld`, `snaptel`, `snap-plugin-publisher-file`,
    `snap-plugin-collector-cpu` and starts the task 
    [tasks/cpu-file.json](tasks/cpu-file.json).
- [.setup.sh](.setup.sh)
    - Verifies dependencies and starts the containers.  It's called 
    by [run-cpu-file.sh](run-cpu-file.sh).