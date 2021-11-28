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

		// 7.xxx
		new(DPT_7001),
		new(DPT_7002),
		new(DPT_7003),
		new(DPT_7004),
		new(DPT_7005),
		new(DPT_7006),
		new(DPT_7007),
		new(DPT_7010),
		new(DPT_7011),
		new(DPT_7012),
		new(DPT_7013),

		// 9.xxx
		new(DPT_9001),
		new(DPT_9002),
		new(DPT_9003),
		new(DPT_9004),
		new(DPT_9005),
		new(DPT_9006),
		new(DPT_9007),
		new(DPT_9008),
		new(DPT_9010),
		new(DPT_9011),
		new(DPT_9020),
		new(DPT_9021),
		new(DPT_9022),
		new(DPT_9023),
		new(DPT_9024),
		new(DPT_9025),
		new(DPT_9026),
		new(DPT_9027),
		new(DPT_9028),

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
		new(DPT_13016),
		new(DPT_13100),

		// 14.xxx
		new(DPT_14000),
		new(DPT_14001),
		new(DPT_14002),
		new(DPT_14010),
		new(DPT_14011),
		new(DPT_14012),
		new(DPT_14013),
		new(DPT_14014),
		new(DPT_14015),
		new(DPT_14016),
		new(DPT_14017),
		new(DPT_14018),
		new(DPT_14019),
		new(DPT_14020),
		new(DPT_14021),
		new(DPT_14022),
		new(DPT_14023),
		new(DPT_14024),
		new(DPT_14025),
		new(DPT_14026),
		new(DPT_14027),
		new(DPT_14028),
		new(DPT_14029),
		new(DPT_14030),
		new(DPT_14031),
		new(DPT_14032),
		new(DPT_14033),
		new(DPT_14034),
		new(DPT_14035),
		new(DPT_14036),
		new(DPT_14037),
		new(DPT_14038),
		new(DPT_14039),
		new(DPT_14040),
		new(DPT_14041),
		new(DPT_14042),
		new(DPT_14043),
		new(DPT_14044),
		new(DPT_14045),
		new(DPT_14046),
		new(DPT_14047),
		new(DPT_14048),
		new(DPT_14049),
		new(DPT_14050),
		new(DPT_14051),
		new(DPT_14052),
		new(DPT_14053),
		new(DPT_14054),
		new(DPT_14055),
		new(DPT_14056),
		new(DPT_14057),
		new(DPT_14058),
		new(DPT_14059),
		new(DPT_14060),
		new(DPT_14061),
		new(DPT_14062),
		new(DPT_14063),
		new(DPT_14064),
		new(DPT_14065),
		new(DPT_14066),
		new(DPT_14067),
		new(DPT_14068),
		new(DPT_14069),
		new(DPT_14070),
		new(DPT_14071),
		new(DPT_14072),
		new(DPT_14073),
		new(DPT_14074),
		new(DPT_14075),
		new(DPT_14076),
		new(DPT_14077),
		new(DPT_14078),
		new(DPT_14079),

		// 17.xxx
		new(DPT_17001),
		// 18.xxx
		new(DPT_18001),
		// 251.xxx
		new(DPT_251600),
	}
	once     sync.Once
	registry map[string]reflect.Type
)

// Init function used to add all types
func setup() {
	// Singleton, can only run once
	once.Do(func() {
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
	})
}

// ListSupportedTypes returns the name all known datapoint-types (DPTs).
func ListSupportedTypes() []string {
	// Setup the registry
	setup()

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

// Produce creates a new instance of the given datapoint-type name e.g. "1.001".
func Produce(name string) (d DatapointValue, ok bool) {
	// Setup the registry
	setup()

	// Lookup the given type and create a new instance of that type
	x, ok := registry[name]
	if ok {
		d = reflect.New(x).Interface().(DatapointValue)
	}
	return d, ok
}
