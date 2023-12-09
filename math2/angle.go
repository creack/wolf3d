package math2

import "math"

// Angle helper type.
type Angle float64

// NewRadAngle returns a new angle from a radians value.
func NewRadAngle[T ~int | ~float32 | ~float64](radians T) Angle {
	return Angle(radians).Normalize()
}

// NewDegAngle returns a new angle from a degrees values.
func NewDegAngle[T ~int | ~float32 | ~float64](degrees T) Angle {
	return NewRadAngle(float64(degrees) * math.Pi / 180)
}

// Radians returns the radians value of the angle.
func (a Angle) Radians() float64 { return float64(a) }

// Degrees returns the degress value of the angle.
func (a Angle) Degrees() float64 { return float64(a) * 180. / math.Pi }

// Normalize the angle on -pi to +pi.
// i.e. 270 degrees -> -90 degrees.
func (a Angle) Normalize() Angle {
	// Module 2pi (360). Doesn't change the sign.
	a = Angle(math.Mod(a.Radians(), 2*math.Pi))

	if a > math.Pi {
		return -(2*math.Pi - a)
	}
	if a < -math.Pi {
		return -(-2*math.Pi - a)
	}
	return a
}
