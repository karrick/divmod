package divmod_test

import (
	"testing"

	"github.com/karrick/divmod"
)

func TestDivmod(t *testing.T) {
	q, r := divmod.Divmod(0, 10, 3)

	if got, want := q, uint(3); got != want {
		t.Errorf("GOT: %#v; WANT: %#v", got, want)
	}
	if got, want := r, uint(1); got != want {
		t.Errorf("GOT: %#v; WANT: %#v", got, want)
	}
}
