package aez

import (
	"fmt"
)

// ErrResistanceTooLarge resistance sent to interpolator is too large
type ErrResistanceTooLarge float64

func (err ErrResistanceTooLarge) String() string {
	return fmt.Sprint(float64(err))
}

func (err ErrResistanceTooLarge) Error() string {
	return fmt.Sprintf(
		"Resistance %v is too large for interpolator. Returning value for maximum.",
		float64(err),
	)
}

// ErrResistanceTooSmall resistance sent to interpolator is too small
type ErrResistanceTooSmall float64

func (err ErrResistanceTooSmall) String() string {
	return fmt.Sprint(float64(err))
}

func (err ErrResistanceTooSmall) Error() string {
	return fmt.Sprintf(
		"Resistance %v is too small for interpolator. Returning value for minimum.",
		float64(err),
	)
}

// ErrResistancesTemperaturesTooShort resistance sent to interpolator is too small
type ErrResistancesTemperaturesTooShort int

func (err ErrResistancesTemperaturesTooShort) String() string {
	return fmt.Sprint(int(err))
}

func (err ErrResistancesTemperaturesTooShort) Error() string {
	return fmt.Sprintf(
		"Resistances and temperatures must be at least of length 2 (%d < 0). Returning absolute zero.",
		int(err),
	)
}

// ErrResistancesTemperaturesMismatch resistances and temperatures of different lengths
type ErrResistancesTemperaturesMismatch struct{ lenRes, lenTemp int }

func (err ErrResistancesTemperaturesMismatch) String() string {
	return fmt.Sprintf("%d, %d", err.lenRes, err.lenTemp)
}

func (err ErrResistancesTemperaturesMismatch) Error() string {
	return fmt.Sprintf(
		"Resistances and temperatures not of same length (%d != %d). Returning absolute zero.",
		err.lenRes, err.lenTemp,
	)
}

func getResistanceIndex(res float64, resistances []float64) int {
	var i int
	for i = range resistances {
		if res > resistances[i] {
			break
		}
	}
	return i
}

func interpolate(r [2]float64, t [2]float64, resistance float64) float64 {
	return ((t[1]-t[0])/(r[1]-r[0]))*(resistance-r[0]) + t[0]
}

func interpolateTemperature(
	resistance float64, resistances []float64, temperatures []float64,
) (float64, error) {
	if len(resistances) != len(temperatures) {
		return -273.15, ErrResistancesTemperaturesMismatch{
			len(resistances), len(temperatures),
		}
	}
	if len(resistances) < 2 {
		return -273.15, ErrResistancesTemperaturesTooShort(len(resistances))
	}

	i := getResistanceIndex(resistance, resistances)
	if i == 0 {
		return temperatures[0], ErrResistanceTooLarge(resistance)
	} else if i == len(resistances)-1 {
		return temperatures[len(resistances)-1], ErrResistanceTooSmall(resistance)
	}

	var r, t [2]float64
	copy(r[:], resistances[i-1:i+1])
	copy(t[:], temperatures[i-1:i+1])
	return interpolate(r, t, resistance), nil
}
