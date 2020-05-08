package aez

import (
	"fmt"
)

// ErrXTooLarge x value sent to interpolator is too large
type ErrXTooLarge float64

func (err ErrXTooLarge) String() string {
	return fmt.Sprint(float64(err))
}

func (err ErrXTooLarge) Error() string {
	return fmt.Sprintf(
		"%v is too large for interpolator. Returning value for maximum.",
		float64(err),
	)
}

// ErrXTooSmall value sent to interpolator is too small
type ErrXTooSmall float64

func (err ErrXTooSmall) String() string {
	return fmt.Sprint(float64(err))
}

func (err ErrXTooSmall) Error() string {
	return fmt.Sprintf(
		"%v is too small for interpolator. Returning value for minimum.",
		float64(err),
	)
}

// ErrXYTooShort x or y slice sent to interpolator is too short
type ErrXYTooShort int

func (err ErrXYTooShort) String() string {
	return fmt.Sprint(int(err))
}

func (err ErrXYTooShort) Error() string {
	return fmt.Sprintf(
		"Slices x and y must be at least of length 2 (%d < 0). Returning -9999.",
		int(err),
	)
}

// ErrXYMismatch x and y slices of different lengths
type ErrXYMismatch struct{ lenRes, lenTemp int }

func (err ErrXYMismatch) String() string {
	return fmt.Sprintf("%d, %d", err.lenRes, err.lenTemp)
}

func (err ErrXYMismatch) Error() string {
	return fmt.Sprintf(
		"Slices x and y not of same length (%d != %d). Returning absolute -9999.",
		err.lenRes, err.lenTemp,
	)
}

// Note! Assumes xSlice contains a monotonically increasing values!
func getXIndex(x float64, xSlice []float64) int {
	var i int
	for i = range xSlice {
		if x > xSlice[i] {
			break
		}
	}
	return i
}

// Interpolate y value corresponding to x value, given
// xSlice (e.g. resistances or voltages), and
// ySlice (e.g. temperatures or photometer values).
func Interpolate(
	x float64, xSlice []float64, ySlice []float64,
) (float64, error) {
	if len(xSlice) != len(ySlice) {
		return -9999, ErrXYMismatch{
			len(xSlice), len(ySlice),
		}
	}
	if len(xSlice) < 2 {
		return -9999, ErrXYTooShort(len(xSlice))
	}

	i := getXIndex(x, xSlice)
	if i == 0 {
		return ySlice[0], ErrXTooLarge(x)
	} else if i == len(xSlice)-1 {
		return ySlice[len(xSlice)-1], ErrXTooSmall(x)
	}

	var xs, ys [2]float64
	copy(xs[:], xSlice[i-1:i+1])
	copy(ys[:], ySlice[i-1:i+1])
	return interpolate(xs, ys, x), nil
}

func interpolate(xs [2]float64, ys [2]float64, x float64) float64 {
	return ((ys[1]-ys[0])/(xs[1]-xs[0]))*(x-xs[0]) + ys[0]
}
