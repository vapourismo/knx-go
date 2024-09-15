// Copyright 2020 Sven Rebhan.
// Copyright 2024 Martin Müller.
// Licensed under the MIT license which can be found in the LICENSE file.

package dpt

import (
	"reflect"
)

var dptTypes = map[string]Datapoint{
	// 1.xxx
	"1.001": new(DPT_1001),
	"1.002": new(DPT_1002),
	"1.003": new(DPT_1003),
	"1.004": new(DPT_1004),
	"1.005": new(DPT_1005),
	"1.006": new(DPT_1006),
	"1.007": new(DPT_1007),
	"1.008": new(DPT_1008),
	"1.009": new(DPT_1009),
	"1.010": new(DPT_1010),
	"1.011": new(DPT_1011),
	"1.012": new(DPT_1012),
	"1.013": new(DPT_1013),
	"1.014": new(DPT_1014),
	"1.015": new(DPT_1015),
	"1.016": new(DPT_1016),
	"1.017": new(DPT_1017),
	"1.018": new(DPT_1018),
	"1.019": new(DPT_1019),
	"1.021": new(DPT_1021),
	"1.022": new(DPT_1022),
	"1.023": new(DPT_1023),
	"1.024": new(DPT_1024),
	"1.100": new(DPT_1100),

	// 5.xxx
	"5.001": new(DPT_5001),
	"5.003": new(DPT_5003),
	"5.004": new(DPT_5004),
	"5.005": new(DPT_5005),

	// 6.xxx
	"6.010": new(DPT_6010),

	// 7.xxx
	"7.001": new(DPT_7001),
	"7.002": new(DPT_7002),
	"7.003": new(DPT_7003),
	"7.004": new(DPT_7004),
	"7.005": new(DPT_7005),
	"7.006": new(DPT_7006),
	"7.007": new(DPT_7007),
	"7.010": new(DPT_7010),
	"7.011": new(DPT_7011),
	"7.012": new(DPT_7012),
	"7.013": new(DPT_7013),
	"7.600": new(DPT_7600),

	// 8.xxx
	"8.001": new(DPT_8001),
	"8.002": new(DPT_8002),
	"8.003": new(DPT_8003),
	"8.004": new(DPT_8004),
	"8.005": new(DPT_8005),
	"8.006": new(DPT_8006),
	"8.007": new(DPT_8007),
	"8.010": new(DPT_8010),
	"8.011": new(DPT_8011),

	// 9.xxx
	"9.001": new(DPT_9001),
	"9.002": new(DPT_9002),
	"9.003": new(DPT_9003),
	"9.004": new(DPT_9004),
	"9.005": new(DPT_9005),
	"9.006": new(DPT_9006),
	"9.007": new(DPT_9007),
	"9.008": new(DPT_9008),
	"9.010": new(DPT_9010),
	"9.011": new(DPT_9011),
	"9.020": new(DPT_9020),
	"9.021": new(DPT_9021),
	"9.022": new(DPT_9022),
	"9.023": new(DPT_9023),
	"9.024": new(DPT_9024),
	"9.025": new(DPT_9025),
	"9.026": new(DPT_9026),
	"9.027": new(DPT_9027),
	"9.028": new(DPT_9028),
	"9.029": new(DPT_9029),

	// 10.xxx
	"10.001": new(DPT_10001),

	// 11.xxx
	"11.001": new(DPT_11001),

	// 12.xxx
	"12.001": new(DPT_12001),

	// 13.xxx
	"13.001": new(DPT_13001),
	"13.002": new(DPT_13002),
	"13.010": new(DPT_13010),
	"13.011": new(DPT_13011),
	"13.012": new(DPT_13012),
	"13.013": new(DPT_13013),
	"13.014": new(DPT_13014),
	"13.015": new(DPT_13015),
	"13.016": new(DPT_13016),
	"13.100": new(DPT_13100),

	// 14.xxx
	"14.000":  new(DPT_14000),
	"14.001":  new(DPT_14001),
	"14.002":  new(DPT_14002),
	"14.003":  new(DPT_14003),
	"14.004":  new(DPT_14004),
	"14.005":  new(DPT_14005),
	"14.006":  new(DPT_14006),
	"14.007":  new(DPT_14007),
	"14.008":  new(DPT_14008),
	"14.009":  new(DPT_14009),
	"14.010":  new(DPT_14010),
	"14.011":  new(DPT_14011),
	"14.012":  new(DPT_14012),
	"14.013":  new(DPT_14013),
	"14.014":  new(DPT_14014),
	"14.015":  new(DPT_14015),
	"14.016":  new(DPT_14016),
	"14.017":  new(DPT_14017),
	"14.018":  new(DPT_14018),
	"14.019":  new(DPT_14019),
	"14.020":  new(DPT_14020),
	"14.021":  new(DPT_14021),
	"14.022":  new(DPT_14022),
	"14.023":  new(DPT_14023),
	"14.024":  new(DPT_14024),
	"14.025":  new(DPT_14025),
	"14.026":  new(DPT_14026),
	"14.027":  new(DPT_14027),
	"14.028":  new(DPT_14028),
	"14.029":  new(DPT_14029),
	"14.030":  new(DPT_14030),
	"14.031":  new(DPT_14031),
	"14.032":  new(DPT_14032),
	"14.033":  new(DPT_14033),
	"14.034":  new(DPT_14034),
	"14.035":  new(DPT_14035),
	"14.036":  new(DPT_14036),
	"14.037":  new(DPT_14037),
	"14.038":  new(DPT_14038),
	"14.039":  new(DPT_14039),
	"14.040":  new(DPT_14040),
	"14.041":  new(DPT_14041),
	"14.042":  new(DPT_14042),
	"14.043":  new(DPT_14043),
	"14.044":  new(DPT_14044),
	"14.045":  new(DPT_14045),
	"14.046":  new(DPT_14046),
	"14.047":  new(DPT_14047),
	"14.048":  new(DPT_14048),
	"14.049":  new(DPT_14049),
	"14.050":  new(DPT_14050),
	"14.051":  new(DPT_14051),
	"14.052":  new(DPT_14052),
	"14.053":  new(DPT_14053),
	"14.054":  new(DPT_14054),
	"14.055":  new(DPT_14055),
	"14.056":  new(DPT_14056),
	"14.057":  new(DPT_14057),
	"14.058":  new(DPT_14058),
	"14.059":  new(DPT_14059),
	"14.060":  new(DPT_14060),
	"14.061":  new(DPT_14061),
	"14.062":  new(DPT_14062),
	"14.063":  new(DPT_14063),
	"14.064":  new(DPT_14064),
	"14.065":  new(DPT_14065),
	"14.066":  new(DPT_14066),
	"14.067":  new(DPT_14067),
	"14.068":  new(DPT_14068),
	"14.069":  new(DPT_14069),
	"14.070":  new(DPT_14070),
	"14.071":  new(DPT_14071),
	"14.072":  new(DPT_14072),
	"14.073":  new(DPT_14073),
	"14.074":  new(DPT_14074),
	"14.075":  new(DPT_14075),
	"14.076":  new(DPT_14076),
	"14.077":  new(DPT_14077),
	"14.078":  new(DPT_14078),
	"14.079":  new(DPT_14079),
	"14.1200": new(DPT_141200),

	// 16.xxx
	"16.000": new(DPT_16000),
	"16.001": new(DPT_16001),

	// 17.xxx
	"17.001": new(DPT_17001),

	// 18.xxx
	"18.001": new(DPT_18001),

	// 20.xxx
	"20.102": new(DPT_20102),
	"20.105": new(DPT_20105),

	// 28.xxx
	"28.001": new(DPT_28001),

	// 232.xxx
	"232.600": new(DPT_232600),

	// 242.xxx
	"242.600": new(DPT_242600),

	// 251.xxx
	"251.600": new(DPT_251600),
}

// ListSupportedTypes returns the names of all known datapoint types (DPTs).
func ListSupportedTypes() []string {
	keys := make([]string, 0, len(dptTypes))

	for k := range dptTypes {
		keys = append(keys, k)
	}

	return keys
}

// Produce returns a new instance, given the exact datapoint name.
// It returns a DPT_1001 for the name "1.001".
func Produce(name string) (d Datapoint, ok bool) {
	// Lookup the given type and return a datapoint of that type.
	x, ok := dptTypes[name]

	if ok {
		d_type := reflect.TypeOf(x).Elem()
		d = reflect.New(d_type).Interface().(Datapoint)
	}

	return d, ok
}
