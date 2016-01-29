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

package main

import (
	"os"

	"github.com/intelsdi-x/snap-plugin-collector-cpu/cpu"
	"github.com/intelsdi-x/snap/control/plugin"
)

func main() {
	cpuPlugin := cpu.New()
	if cpuPlugin == nil {
		panic("Failed to initialize plugin!\n")
	}
	meta := cpu.Meta()
	plugin.Start(meta, cpuPlugin, os.Args[1])
}
