package aez

import (
	"time"

	"github.com/innosat-mats/rac-extract-payload/internal/ccsds"
)

type pm uint32

var pm1Voltages = [...]float64{
	0.002, 0.199, 0.396, 0.593, 0.798,
	0.994, 0.191, 1.405, 1.593, 1.798,
	1.994, 2.199, 2.395, 2.592, 2.797,
	3.002, 3.199, 3.395, 3.592, 3.797,
	3.993, 4.199, 4.396, 4.601, 4.798,
	4.995, 5.102, 5.192,
} // V
var pm1AValues = [...]float64{
	1.1, 162.1, 323.0, 484.0, 652.9,
	813.5, 974.8, 1150.2, 1304.3, 1471.0,
	1632.6, 1801.0, 1961.8, 2122.7, 2290.9,
	2459.0, 2619.0, 2780.5, 2941.9, 3110.4,
	3271.0, 3439.0, 3599.3, 3769.0, 3930.9,
	4091.3, 4095.0, 4095.0,
} // a.u. (thermistor)
var pm1BValues = [...]float64{
	1.1, 162.1, 323.0, 484.0, 652.6,
	813.1, 974.2, 1149.9, 1303.2, 1471.0,
	1631.1, 1800.6, 1961.0, 2121.5, 2289.8,
	2458.3, 2619.0, 2779.0, 2940.8, 3108.9,
	3270.7, 3439.0, 3599.0, 3767.0, 3928.0,
	4090.7, 4095.0, 4095.0,
} // a.u. (thermistor)
var pm1SValues = [...]float64{
	1.1, 162.1, 323.1, 484.0, 652.9,
	813.6, 974.7, 1150.1, 1304.3, 1471.0,
	1632.8, 1801.0, 1961.6, 2122.7, 2290.8,
	2459.0, 2619.0, 2780.3, 2940.7, 3110.2,
	3271.0, 3439.0, 3599.4, 3768.8, 3930.9,
	4091.1, 4095.0, 4095.0,
} // a.u. (photo diode)
var pm2Voltages = [...]float64{
	0.002, 0.199, 0.396, 0.593, 0.798,
	0.994, 0.191, 1.405, 1.593, 1.798,
	1.994, 2.199, 2.395, 2.592, 2.797,
	3.002, 3.199, 3.395, 3.600, 3.797,
	3.993, 4.199, 4.396, 4.601, 4.798,
	4.995, 5.102, 5.192,
} // V
var pm2AValues = [...]float64{
	1.0, 162.2, 323.2, 484.4, 653.0,
	814.4, 975.0, 1151.0, 1305.0, 1473.0,
	1633.7, 1802.1, 1963.0, 2123.0, 2291.1,
	2460.0, 2621.0, 2781.8, 2951.0, 3111.1,
	3272.8, 3442.5, 3603.0, 3771.0, 3932.4,
	4094.4, 4095.0, 4095.0,
} // a.u. (thermistor)
var pm2BValues = [...]float64{
	1.0, 162.5, 323.7, 484.9, 653.3,
	815.0, 975.7, 1151.3, 1306.0, 1474.2,
	1635.0, 1803.0, 1964.2, 2125.0, 2293.0,
	2461.5, 2623.0, 2783.0, 2951.7, 3114.3,
	3275.0, 3443.4, 3605.3, 3774.1, 3935.0,
	4095.0, 4095.0, 4095.0,
} // a.u. (thermistor)
var pm2SValues = [...]float64{
	1.0, 162.5, 323.7, 484.9, 653.4,
	815.0, 976.0, 1151.9, 1306.2, 1474.3,
	1635.0, 1803.0, 1964.7, 2125.0, 2293.2,
	2462.3, 2623.0, 2783.0, 2952.9, 3114.8,
	3275.1, 3444.8, 3606.4, 3774.9, 3935.0,
	4095.0, 4095.0, 4095.0,
} // a.u. (photo diode)

// PMData data from photometers
type PMData struct {
	EXPTS    uint32 // Exposure start time, seconds (CUC time format)
	EXPTSS   uint16 // Exposure start time, subseconds (CUC time format)
	PM1A     pm     // Photometer 1, thermistor input A sum
	PM1ACNTR pm     // Photometer 1, thermistor input A counter
	PM1B     pm     // Photometer 1, thermistor input B sum
	PM1BCNTR pm     // Photometer 1, thermistor input B counter
	PM1S     pm     // Photometer 1, photo diode input SIG sum
	PM1SCNTR pm     // Photometer 1, photo diode input SIG counter
	PM2A     pm     // Photometer 2, thermistor input A sum
	PM2ACNTR pm     // Photometer 2, thermistor input A counter
	PM2B     pm     // Photometer 2, thermistor input B sum
	PM2BCNTR pm     // Photometer 2, thermistor input B counter
	PM2S     pm     // Photometer 2, photo diode input SIG sum
	PM2SCNTR pm     // Photometer 2, photo diode input SIG counter
}

// Time returns the measurement time in UTC
func (pm *PMData) Time(epoch time.Time) time.Time {
	return ccsds.UnsegmentedTimeDate(pm.EXPTS, pm.EXPTSS, epoch)
}

// Nanoseconds returns the measurement time in nanoseconds since epoch
func (pm *PMData) Nanoseconds() int64 {
	return ccsds.UnsegmentedTimeNanoseconds(pm.EXPTS, pm.EXPTSS)
}
