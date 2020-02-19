// Copyright 2020 Sven Rebhan.
// Licensed under the MIT license which can be found in the LICENSE file.

package dpt

import (
	"reflect"
	"sync"
)

var (
	types = [...]DatapointValue{
		// 1.xxx
		new(DPT_1001),
		new(DPT_1002),
		new(DPT_1003),
		new(DPT_1009),
		new(DPT_1010),

		// 5.xxx
		new(DPT_5001),
		new(DPT_5003),
		new(DPT_5004),

		// 9.xxx
		new(DPT_9001),
		new(DPT_9004),
		new(DPT_9005),
		new(DPT_9007),

		// 12.xxx
		new(DPT_12001),

		// 13.xxx
		new(DPT_13001),
		new(DPT_13002),
		new(DPT_13010),
		new(DPT_13011),
		new(DPT_13012),
		new(DPT_13013),
		new(DPT_13014),
		new(DPT_13015),
	}
	once     sync.Once
	registry map[string]reflect.Type
)

// Init function used to add all types
func setup() {
	// Register the types
	registry = make(map[string]reflect.Type)
	for _, d := range types {
		// Determine the name of the datatype
		d_type := reflect.TypeOf(d).Elem()
		name := d_type.Name()

		// Convert the name into KNX yy.xxx (e.g. DPT_1001 --> 1.001) format
		name = name[4:len(name)-3] + "." + name[len(name)-3:]

		// Register the type
		registry[name] = d_type
	}
}

func ListSupportedTypes() []string {
	// Singleton, can only run once
	once.Do(setup)

	// Initialize the key-list
	keys := make([]string, len(registry))

	// Fill the key-list
	i := 0
	for k := range registry {
		keys[i] = k
		i++
	}

	return keys
}

func Produce(name string) (d DatapointValue, ok bool) {
	// Singleton, can only run once
	once.Do(setup)

	x, ok := registry[name]

	if ok {
		d = reflect.New(x).Interface().(DatapointValue)
	}
	return
}
