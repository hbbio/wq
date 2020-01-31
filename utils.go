package wq

// MinInt computes the minimum of two values.
func MinInt(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// Plural displays 's' if value is more than 1.
func Plural(nb int) string {
	if nb > 1 {
		return "s"
	}
	return ""
}
