#include "textflag.h"

// func divmod(x1, x0, y Word) (q, r Word)
TEXT ·divmod(SB),NOSPLIT,$0
	MOVQ x1+0(FP), DX
	MOVQ x0+8(FP), AX
	DIVQ y+16(FP)
	MOVQ AX, q+24(FP)
	MOVQ DX, r+32(FP)
	RET
