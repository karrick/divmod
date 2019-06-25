package divmod

// A Word represents a single digit of a multi-precision unsigned integer.
type Word uint

func divmod(x1, x0, y uint) (q, r uint)

func Divmod(x1, x0, y uint) (q, r uint) {
	return divmod(x1, x0, y)
}
