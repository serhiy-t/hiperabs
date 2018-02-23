package benchmarks

// Point ...
type Point struct {
	X, Y int
}

func (p Point) Hash() int {
	return p.X*31 ^ p.Y
}
