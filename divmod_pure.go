package divmod

import "math/bits"

const (
	_S = _W / 8 // word size in bytes

	_W = bits.UintSize // word size in bits
	_B = 1 << _W       // digit base
	_M = _B - 1        // digit mask

	_W2 = _W / 2   // half word size in bits
	_B2 = 1 << _W2 // half digit base
	_M2 = _B2 - 1  // half digit mask
)

func divmod_go(u1, u0, v uint) (q, r uint) {
	if u1 >= v {
		return 1<<_W - 1, 1<<_W - 1
	}

	s := uint(bits.LeadingZeros(uint(v)))
	v <<= s

	vn1 := v >> _W2
	vn0 := v & _M2
	un32 := u1<<s | u0>>(_W-s)
	un10 := u0 << s
	un1 := un10 >> _W2
	un0 := un10 & _M2
	q1 := un32 / vn1
	rhat := un32 - q1*vn1

	for q1 >= _B2 || q1*vn0 > _B2*rhat+un1 {
		q1--
		rhat += vn1
		if rhat >= _B2 {
			break
		}
	}

	un21 := un32*_B2 + un1 - q1*v
	q0 := un21 / vn1
	rhat = un21 - q0*vn1

	for q0 >= _B2 || q0*vn0 > _B2*rhat+un0 {
		q0--
		rhat += vn1
		if rhat >= _B2 {
			break
		}
	}

	return q1*_B2 + q0, (un21*_B2 + un0 - q0*v) >> s
}
