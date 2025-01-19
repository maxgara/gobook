	.section	__TEXT,__text,regular,pure_instructions
	.build_version macos, 15, 0	sdk_version 15, 2
	.section	__TEXT,__literal16,16byte_literals
	.p2align	4, 0x0                          ; -- Begin function main
lCPI0_0:
	.long	0                               ; 0x0
	.long	1                               ; 0x1
	.long	2                               ; 0x2
	.long	3                               ; 0x3
	.section	__TEXT,__text,regular,pure_instructions
	.globl	_main
	.p2align	2
_main:                                  ; @main
	.cfi_startproc
; %bb.0:
	stp	x22, x21, [sp, #-48]!           ; 16-byte Folded Spill
	stp	x20, x19, [sp, #16]             ; 16-byte Folded Spill
	stp	x29, x30, [sp, #32]             ; 16-byte Folded Spill
	add	x29, sp, #32
	mov	w9, #40032                      ; =0x9c60
Lloh0:
	adrp	x16, ___chkstk_darwin@GOTPAGE
Lloh1:
	ldr	x16, [x16, ___chkstk_darwin@GOTPAGEOFF]
	blr	x16
	sub	sp, sp, #9, lsl #12             ; =36864
	sub	sp, sp, #3168
	.cfi_def_cfa w29, 16
	.cfi_offset w30, -8
	.cfi_offset w29, -16
	.cfi_offset w19, -24
	.cfi_offset w20, -32
	.cfi_offset w21, -40
	.cfi_offset w22, -48
Lloh2:
	adrp	x8, ___stack_chk_guard@GOTPAGE
Lloh3:
	ldr	x8, [x8, ___stack_chk_guard@GOTPAGEOFF]
Lloh4:
	ldr	x8, [x8]
	stur	x8, [x29, #-40]
	bl	_clock
	mov	x19, x0
Lloh5:
	adrp	x8, lCPI0_0@PAGE
Lloh6:
	ldr	q0, [x8, lCPI0_0@PAGEOFF]
	mov	x8, #-40000                     ; =0xffffffffffff63c0
	movi.4s	v1, #1
	add	x9, sp, #24
	movi.4s	v2, #16
LBB0_1:                                 ; =>This Inner Loop Header: Depth=1
	bic.16b	v3, v1, v0
	add	x10, x9, x8
	str	q3, [x10, #40000]
	str	q3, [x10, #40016]
	str	q3, [x10, #40032]
	str	q3, [x10, #40048]
	add.4s	v0, v0, v2
	adds	x8, x8, #64
	b.ne	LBB0_1
; %bb.2:
	movi.2d	v0, #0000000000000000
	mov	x8, #-40000                     ; =0xffffffffffff63c0
	add	x9, sp, #24
	movi.2d	v1, #0000000000000000
	movi.2d	v2, #0000000000000000
	movi.2d	v3, #0000000000000000
LBB0_3:                                 ; =>This Inner Loop Header: Depth=1
	add	x10, x9, x8
	ldr	q4, [x10, #40000]
	ldr	q5, [x10, #40016]
	ldr	q6, [x10, #40032]
	ldr	q7, [x10, #40048]
	add.4s	v0, v4, v0
	add.4s	v1, v5, v1
	add.4s	v2, v6, v2
	add.4s	v3, v7, v3
	adds	x8, x8, #64
	b.ne	LBB0_3
; %bb.4:
	add.4s	v0, v1, v0
	add.4s	v0, v2, v0
	add.4s	v0, v3, v0
	addv.4s	s0, v0
	fmov	w21, s0
	bl	_clock
	mov	x20, x0
	mov	w8, #10000                      ; =0x2710
	stp	x8, x21, [sp]
Lloh7:
	adrp	x0, l_.str@PAGE
Lloh8:
	add	x0, x0, l_.str@PAGEOFF
	bl	_printf
	sub	x8, x20, x19
	str	x8, [sp]
Lloh9:
	adrp	x0, l_.str.1@PAGE
Lloh10:
	add	x0, x0, l_.str.1@PAGEOFF
	bl	_printf
	ldur	x8, [x29, #-40]
Lloh11:
	adrp	x9, ___stack_chk_guard@GOTPAGE
Lloh12:
	ldr	x9, [x9, ___stack_chk_guard@GOTPAGEOFF]
Lloh13:
	ldr	x9, [x9]
	cmp	x9, x8
	b.ne	LBB0_6
; %bb.5:
	mov	w0, #0                          ; =0x0
	add	sp, sp, #9, lsl #12             ; =36864
	add	sp, sp, #3168
	ldp	x29, x30, [sp, #32]             ; 16-byte Folded Reload
	ldp	x20, x19, [sp, #16]             ; 16-byte Folded Reload
	ldp	x22, x21, [sp], #48             ; 16-byte Folded Reload
	ret
LBB0_6:
	bl	___stack_chk_fail
	.loh AdrpLdr	Lloh5, Lloh6
	.loh AdrpLdrGotLdr	Lloh2, Lloh3, Lloh4
	.loh AdrpLdrGot	Lloh0, Lloh1
	.loh AdrpLdrGotLdr	Lloh11, Lloh12, Lloh13
	.loh AdrpAdd	Lloh9, Lloh10
	.loh AdrpAdd	Lloh7, Lloh8
	.cfi_endproc
                                        ; -- End function
	.globl	_iseven                         ; -- Begin function iseven
	.p2align	2
_iseven:                                ; @iseven
	.cfi_startproc
; %bb.0:
	add	w8, w0, #5
	and	w8, w8, #0x80000001
	cmp	w8, #1
	cset	w0, eq
	ret
	.cfi_endproc
                                        ; -- End function
	.section	__TEXT,__cstring,cstring_literals
l_.str:                                 ; @.str
	.asciz	"even counter(%d):%d\n"

l_.str.1:                               ; @.str.1
	.asciz	"time elapsed:%lu\n"

.subsections_via_symbols
