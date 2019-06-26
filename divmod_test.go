package divmod

import (
	"testing"
)

func TestDivmod(t *testing.T) {
	q, r := Divmod(0, 10, 3)

	if got, want := q, uint(3); got != want {
		t.Errorf("GOT: %#v; WANT: %#v", got, want)
	}
	if got, want := r, uint(1); got != want {
		t.Errorf("GOT: %#v; WANT: %#v", got, want)
	}
}

func BenchmarkOperators(b *testing.B) {
	var q, r uint
	for i := 0; i < b.N; i++ {
		q = 10 / 3
		r = 10 % 3
	}
	_, _ = q, r
}

func BenchmarkOperatorFunction(b *testing.B) {
	var q, r uint
	for i := 0; i < b.N; i++ {
		q, r = divmod(10, 3)
	}
	_, _ = q, r
}

func divmod_operator_function(n, d uint) (q, r uint) {
	q = n / d
	r = n % d
	return
}

func BenchmarkOperatorFunctionNoInline(b *testing.B) {
	var q, r uint
	for i := 0; i < b.N; i++ {
		q, r = divmod_operator_function_noinline(10, 3)
	}
	_, _ = q, r
}

//go:noinline
func divmod_operator_function_noinline(n, d uint) (q, r uint) {
	q = n / d
	r = n % d
	return
}

func BenchmarkDivmod(b *testing.B) {
	var q, r uint
	for i := 0; i < b.N; i++ {
		q, r = Divmod(0, 10, 3)
	}
	_, _ = q, r
}

func BenchmarkDivmodPure(b *testing.B) {
	var q, r uint
	for i := 0; i < b.N; i++ {
		q, r = divmod_go(0, 10, 3)
	}
	_, _ = q, r
}
