package big

import "math/bits"

// A Word represents a single digit of a multi-precision unsigned integer.
type Word uint

const (
	_W = bits.UintSize // word size in bits
	_B = 1 << _W       // digit base
	_M = _B - 1        // digit mask
)

// z1<<_W + z0 = x*y
func mulWW_g(x, y Word) (z1, z0 Word) {
	hi, lo := bits.Mul(uint(x), uint(y))
	return Word(hi), Word(lo)
}

// z1<<_W + z0 = x*y + c
func mulAddWWW_g(x, y, c Word) (z1, z0 Word) {
	hi, lo := bits.Mul(uint(x), uint(y))
	var cc uint
	lo, cc = bits.Add(lo, uint(c), 0)
	return Word(hi + cc), Word(lo)
}

// nlz returns the number of leading zeros in x.
// Wraps bits.LeadingZeros call for convenience.
func nlz(x Word) uint {
	return uint(bits.LeadingZeros(uint(x)))
}

// The resulting carry c is either 0 or 1.
func addVV_g(z, x, y []Word) (c Word) {
	// The comment near the top of this file discusses this for loop condition.
	for i := 0; i < len(z) && i < len(x) && i < len(y); i++ {
		zi, cc := bits.Add(uint(x[i]), uint(y[i]), uint(c))
		z[i] = Word(zi)
		c = Word(cc)
	}
	return
}

// The resulting carry c is either 0 or 1.
func subVV_g(z, x, y []Word) (c Word) {
	// The comment near the top of this file discusses this for loop condition.
	for i := 0; i < len(z) && i < len(x) && i < len(y); i++ {
		zi, cc := bits.Sub(uint(x[i]), uint(y[i]), uint(c))
		z[i] = Word(zi)
		c = Word(cc)
	}
	return
}

// The resulting carry c is either 0 or 1.
func addVW_g(z, x []Word, y Word) (c Word) {
	c = y
	// The comment near the top of this file discusses this for loop condition.
	for i := 0; i < len(z) && i < len(x); i++ {
		zi, cc := bits.Add(uint(x[i]), uint(c), 0)
		z[i] = Word(zi)
		c = Word(cc)
	}
	return
}

// addVWlarge is addVW, but intended for large z.
// The only difference is that we check on every iteration
// whether we are done with carries,
// and if so, switch to a much faster copy instead.
// This is only a good idea for large z,
// because the overhead of the check and the function call
// outweigh the benefits when z is small.
func addVWlarge(z, x []Word, y Word) (c Word) {
	c = y
	// The comment near the top of this file discusses this for loop condition.
	for i := 0; i < len(z) && i < len(x); i++ {
		if c == 0 {
			copy(z[i:], x[i:])
			return
		}
		zi, cc := bits.Add(uint(x[i]), uint(c), 0)
		z[i] = Word(zi)
		c = Word(cc)
	}
	return
}

func subVW_g(z, x []Word, y Word) (c Word) {
	c = y
	// The comment near the top of this file discusses this for loop condition.
	for i := 0; i < len(z) && i < len(x); i++ {
		zi, cc := bits.Sub(uint(x[i]), uint(c), 0)
		z[i] = Word(zi)
		c = Word(cc)
	}
	return
}

// subVWlarge is to subVW as addVWlarge is to addVW.
func subVWlarge(z, x []Word, y Word) (c Word) {
	c = y
	// The comment near the top of this file discusses this for loop condition.
	for i := 0; i < len(z) && i < len(x); i++ {
		if c == 0 {
			copy(z[i:], x[i:])
			return
		}
		zi, cc := bits.Sub(uint(x[i]), uint(c), 0)
		z[i] = Word(zi)
		c = Word(cc)
	}
	return
}

func shlVU_g(z, x []Word, s uint) (c Word) {
	if s == 0 {
		copy(z, x)
		return
	}
	if len(z) == 0 {
		return
	}
	s &= _W - 1 // hint to the compiler that shifts by s don't need guard code
	ŝ := _W - s
	ŝ &= _W - 1 // ditto
	c = x[len(z)-1] >> ŝ
	for i := len(z) - 1; i > 0; i-- {
		z[i] = x[i]<<s | x[i-1]>>ŝ
	}
	z[0] = x[0] << s
	return
}

func shrVU_g(z, x []Word, s uint) (c Word) {
	if s == 0 {
		copy(z, x)
		return
	}
	if len(z) == 0 {
		return
	}
	s &= _W - 1 // hint to the compiler that shifts by s don't need guard code
	ŝ := _W - s
	ŝ &= _W - 1 // ditto
	c = x[0] << ŝ
	for i := 0; i < len(z)-1; i++ {
		z[i] = x[i]>>s | x[i+1]<<ŝ
	}
	z[len(z)-1] = x[len(z)-1] >> s
	return
}

func mulAddVWW_g(z, x []Word, y, r Word) (c Word) {
	c = r
	// The comment near the top of this file discusses this for loop condition.
	for i := 0; i < len(z) && i < len(x); i++ {
		c, z[i] = mulAddWWW_g(x[i], y, c)
	}
	return
}

func addMulVVW_g(z, x []Word, y Word) (c Word) {
	// The comment near the top of this file discusses this for loop condition.
	for i := 0; i < len(z) && i < len(x); i++ {
		z1, z0 := mulAddWWW_g(x[i], y, z[i])
		lo, cc := bits.Add(uint(z0), uint(c), 0)
		c, z[i] = Word(cc), Word(lo)
		c += z1
	}
	return
}

// q = ( x1 << _W + x0 - r)/y. m = floor(( _B^2 - 1 ) / d - _B). Requiring x1<y.
// An approximate reciprocal with a reference to "Improved Division by Invariant Integers
// (IEEE Transactions on Computers, 11 Jun. 2010)"
func divWW(x1, x0, y, m Word) (q, r Word) {
	s := nlz(y)
	if s != 0 {
		x1 = x1<<s | x0>>(_W-s)
		x0 <<= s
		y <<= s
	}
	d := uint(y)

	t1, t0 := bits.Mul(uint(m), uint(x1))
	_, c := bits.Add(t0, uint(x0), 0)
	t1, _ = bits.Add(t1, uint(x1), c)

	qq := t1
	dq1, dq0 := bits.Mul(d, qq)
	r0, b := bits.Sub(uint(x0), dq0, 0)
	r1, _ := bits.Sub(uint(x1), dq1, b)

	if r1 != 0 {
		qq++
		r0 -= d
	}
	// If the remainder is still too large, increment q one more time.
	if r0 >= d {
		qq++
		r0 -= d
	}
	return Word(qq), Word(r0 >> s)
}

func divWVW(z []Word, xn Word, x []Word, y Word) (r Word) {
	r = xn
	if len(x) == 1 {
		qq, rr := bits.Div(uint(r), uint(x[0]), uint(y))
		z[0] = Word(qq)
		return Word(rr)
	}
	rec := reciprocalWord(y)
	for i := len(z) - 1; i >= 0; i-- {
		z[i], r = divWW(r, x[i], y, rec)
	}
	return r
}

// reciprocalWord return the reciprocal of the divisor. rec = floor(( _B^2 - 1 ) / u - _B). u = d1 << nlz(d1).
func reciprocalWord(d1 Word) Word {
	u := uint(d1 << nlz(d1))
	x1 := ^u
	x0 := uint(_M)
	rec, _ := bits.Div(x1, x0, u) // (_B^2-1)/U-_B = (_B*(_M-C)+_M)/U
	return Word(rec)
}
