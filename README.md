# divmod
Test implementation of divmod in Go

# Observation

Assume the following function is defined in Go:

```Go
//go:noinline
func divmod_operator_function(n, d uint) (q, r uint) {
	q = n / d
	r = n % d
	return
}
```

One might think that Go is going to perform two division operations,
in the first throwing away the remainder, and the second throwing away
the quotient, but here is the corresponding output from the compiler
for that function:

    go build -gcflags=-S

```
"".divmod_operator_function STEXT nosplit size=58 args=0x20 locals=0x8
	0x0000 00000 (/home/kmcdermo/go/src/github.com/karrick/divmod/divmod_test.go:18)	TEXT	"".divmod_operator_function(SB), NOSPLIT|ABIInternal, $8-32
	0x0000 00000 (/home/kmcdermo/go/src/github.com/karrick/divmod/divmod_test.go:18)	SUBQ	$8, SP
	0x0004 00004 (/home/kmcdermo/go/src/github.com/karrick/divmod/divmod_test.go:18)	MOVQ	BP, (SP)
	0x0008 00008 (/home/kmcdermo/go/src/github.com/karrick/divmod/divmod_test.go:18)	LEAQ	(SP), BP
	0x000c 00012 (/home/kmcdermo/go/src/github.com/karrick/divmod/divmod_test.go:18)	FUNCDATA	$0, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
	0x000c 00012 (/home/kmcdermo/go/src/github.com/karrick/divmod/divmod_test.go:18)	FUNCDATA	$1, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
	0x000c 00012 (/home/kmcdermo/go/src/github.com/karrick/divmod/divmod_test.go:18)	FUNCDATA	$3, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
	0x000c 00012 (/home/kmcdermo/go/src/github.com/karrick/divmod/divmod_test.go:19)	PCDATA	$2, $0
	0x000c 00012 (/home/kmcdermo/go/src/github.com/karrick/divmod/divmod_test.go:19)	PCDATA	$0, $0
	0x000c 00012 (/home/kmcdermo/go/src/github.com/karrick/divmod/divmod_test.go:19)	MOVQ	"".d+24(SP), CX
	0x0011 00017 (/home/kmcdermo/go/src/github.com/karrick/divmod/divmod_test.go:19)	TESTQ	CX, CX
	0x0014 00020 (/home/kmcdermo/go/src/github.com/karrick/divmod/divmod_test.go:19)	JEQ	51
	0x0016 00022 (/home/kmcdermo/go/src/github.com/karrick/divmod/divmod_test.go:19)	MOVQ	"".n+16(SP), AX
	0x001b 00027 (/home/kmcdermo/go/src/github.com/karrick/divmod/divmod_test.go:19)	XORL	DX, DX
	0x001d 00029 (/home/kmcdermo/go/src/github.com/karrick/divmod/divmod_test.go:19)	DIVQ	CX
	0x0020 00032 (/home/kmcdermo/go/src/github.com/karrick/divmod/divmod_test.go:21)	MOVQ	AX, "".q+32(SP)
	0x0025 00037 (/home/kmcdermo/go/src/github.com/karrick/divmod/divmod_test.go:21)	MOVQ	DX, "".r+40(SP)
	0x002a 00042 (/home/kmcdermo/go/src/github.com/karrick/divmod/divmod_test.go:21)	MOVQ	(SP), BP
	0x002e 00046 (/home/kmcdermo/go/src/github.com/karrick/divmod/divmod_test.go:21)	ADDQ	$8, SP
	0x0032 00050 (/home/kmcdermo/go/src/github.com/karrick/divmod/divmod_test.go:21)	RET
	0x0033 00051 (/home/kmcdermo/go/src/github.com/karrick/divmod/divmod_test.go:19)	CALL	runtime.panicdivide(SB)
	0x0038 00056 (/home/kmcdermo/go/src/github.com/karrick/divmod/divmod_test.go:19)	UNDEF
	0x0000 48 83 ec 08 48 89 2c 24 48 8d 2c 24 48 8b 4c 24  H...H.,$H.,$H.L$
	0x0010 18 48 85 c9 74 1d 48 8b 44 24 10 31 d2 48 f7 f1  .H..t.H.D$.1.H..
	0x0020 48 89 44 24 20 48 89 54 24 28 48 8b 2c 24 48 83  H.D$ H.T$(H.,$H.
	0x0030 c4 08 c3 e8 00 00 00 00 0f 0b                    ..........
	rel 52+4 t=8 runtime.panicdivide+0
```

A few things stand out.  After the function prologue the denominator
is loaded into the CX register at 0x000c.  This is immediately
followed by a test operation to check whether the denominator is 0,
and if so, a jump to the line that calls `runtime.panicdivide`.  So
the first observation is that Go always inserts code to ensure the
denominator is not zero before executing the division instruction.
Impressively, adding our own guard into the Go code, I thought there
would be two comparisons of CX against 0, but the Go compiler is smart
enough to omit the extra check.

The next thing to notice is that there is only a single DIV
instruction the resulting assembly language.  The Go compiler emits
code to copy the quotient and back in the stack as a return value,
then does the same for the remainder.

The third observation is that this Go function can be inlined, and
when done so, Go is essentially the exact same code that an assembly
language programmer would emit.

## Inlining

When the following function exists, when I emit the assembly it not
only does not call `divmod_inline`, but encodes the results as
constants in the emitted program.

```Go
func main() {
	q, r := divmod_inline(11, 3)
	fmt.Println(q, r)

	q, r = divmod_noinline(17, 5)
	fmt.Println(q, r)
}

func divmod_inline(n, d uint) (q, r uint) {
	q = n / d
	r = n % d
	return
}

//go:noinline
func divmod_noinline(n, d uint) (q, r uint) {
	q = n / d
	r = n % d
	return
}
```

Modifying the Go source code to take the numerator and denominator
arguments from the command line we can prevent the Go compiler from
computing the results and adding them as constants in the machine code
text.  Here we finally see the benefit of inlining, coupled with Go's
ability to use both outputs from executing the DIV instruction once
rather than executing it twice.

```
	0x00f6 00246 (/home/kmcdermo/go/src/github.com/karrick/divmod/examples/main.go:19)	MOVQ	AX, ""..autotmp_68+72(SP)
	0x00fb 00251 (/home/kmcdermo/go/src/github.com/karrick/divmod/examples/main.go:19)	MOVQ	AX, CX
	0x00fe 00254 (/home/kmcdermo/go/src/github.com/karrick/divmod/examples/main.go:35)	MOVQ	""..autotmp_67+80(SP), AX
	0x0103 00259 (/home/kmcdermo/go/src/github.com/karrick/divmod/examples/main.go:35)	XORL	DX, DX
	0x0105 00261 (/home/kmcdermo/go/src/github.com/karrick/divmod/examples/main.go:35)	DIVQ	CX
	0x0108 00264 (/home/kmcdermo/go/src/github.com/karrick/divmod/examples/main.go:36)	MOVQ	DX, "".r+64(SP)
	0x010d 00269 (/home/kmcdermo/go/src/github.com/karrick/divmod/examples/main.go:28)	MOVQ	AX, (SP)
```
