package divmod

func Divmod(x1, x0, y uint) (q, r uint)

func Inline(n, d uint) (q, r uint) {
	q = n / d
	r = n % d
	return
}
