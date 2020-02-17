// Copyright 2020 Sven Rebhan.
// Licensed under the MIT license which can be found in the LICENSE file.

package dpt

import (
	"fmt"
	"reflect"
	"strings"
)

type Registry struct {
	lut map[string]DatapointValue
}

// Fill the LUT with the given datatype
func (r *Registry) add(d DatapointValue) error {
	if r.lut == nil {
		return fmt.Errorf("registry is not initialized!")
	}

	// Determine the name of the datatype
	d_type := reflect.TypeOf(d)
	name := d_type.Name()
	if d_type.Kind() == reflect.Ptr {
		name = d_type.Elem().Name()
	}

	// Make sure we only handle DPT types^
	if !strings.HasPrefix(name, "DPT_") {
		return fmt.Errorf("invalid type \"%v\" for registry!", name)
	}

	// Convert the name into KNX yy.xxx (e.g. DPT_1001 --> 1.001) format
	name = strings.TrimPrefix(name, "DPT_")
	name = name[:len(name)-3] + "." + name[len(name)-3:]

	// Register the type
	r.lut[name] = d

	return nil
}

// Init function used to add all types
func (r *Registry) Init() error {
	r.lut = make(map[string]DatapointValue)

	// Create a list of all known datapoint-types
	dpts := make([]DatapointValue, 0)

	// 1.xxx
	dpts = append(dpts, new(DPT_1001))
	dpts = append(dpts, new(DPT_1002))
	dpts = append(dpts, new(DPT_1003))
	dpts = append(dpts, new(DPT_1009))
	dpts = append(dpts, new(DPT_1010))

	// 5.xxx
	dpts = append(dpts, new(DPT_5001))
	dpts = append(dpts, new(DPT_5003))
	dpts = append(dpts, new(DPT_5004))

	// 9.xxx
	dpts = append(dpts, new(DPT_9001))
	dpts = append(dpts, new(DPT_9004))
	dpts = append(dpts, new(DPT_9005))
	dpts = append(dpts, new(DPT_9007))

	// 12.xxx
	dpts = append(dpts, new(DPT_12001))

	// 13.xxx
	dpts = append(dpts, new(DPT_13001))
	dpts = append(dpts, new(DPT_13002))
	dpts = append(dpts, new(DPT_13010))
	dpts = append(dpts, new(DPT_13011))
	dpts = append(dpts, new(DPT_13012))
	dpts = append(dpts, new(DPT_13013))
	dpts = append(dpts, new(DPT_13014))
	dpts = append(dpts, new(DPT_13015))

	// Register the types
	for _, d := range dpts {
		err := r.add(d)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r Registry) List() []string {
	// Initialize the key-list
	keys := make([]string, len(r.lut))

	// Fill the key-list
	i := 0
	for k := range r.lut {
		keys[i] = k
		i++
	}

	return keys
}

func (r Registry) Lookup(name string) (t DatapointValue, err error) {
	t, ok := r.lut[name]
	if ok {
		err = nil
	} else {
		err = fmt.Errorf("datapoint-type \"%v\" not found!", name)
	}

	return
}
