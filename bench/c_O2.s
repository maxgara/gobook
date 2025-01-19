	.section	__TEXT,__text,regular,pure_instructions
	.build_version macos, 15, 0	sdk_version 15, 2
	.section	__TEXT,__literal16,16byte_literals
	.p2align	4, 0x0                          ; -- Begin function main
lCPI0_0:
	.quad	2                               ; 0x2
	.quad	3                               ; 0x3
lCPI0_1:
	.quad	0                               ; 0x0
	.quad	1                               ; 0x1
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
Lloh7:
	adrp	x8, lCPI0_1@PAGE
Lloh8:
	ldr	q1, [x8, lCPI0_1@PAGEOFF]
	mov	w8, #4                          ; =0x4
	dup.2d	v2, x8
	mov	w8, #8                          ; =0x8
	dup.2d	v3, x8
	mov	w8, #12                         ; =0xc
	dup.2d	v4, x8
	mov	w8, #10                         ; =0xa
	dup.2d	v5, x8
	mov	x8, #-40000                     ; =0xffffffffffff63c0
	movi.4s	v6, #1
	add	x9, sp, #24
	mov	w10, #16                        ; =0x10
	dup.2d	v7, x10
LBB0_1:                                 ; =>This Inner Loop Header: Depth=1
	add.2d	v16, v1, v2
	add.2d	v17, v0, v2
	add.2d	v18, v1, v3
	add.2d	v19, v0, v3
	add.2d	v20, v1, v4
	add.2d	v21, v0, v4
	cmhi.2d	v22, v5, v0
	cmhi.2d	v23, v5, v1
	uzp1.4s	v22, v23, v22
	cmhi.2d	v17, v5, v17
	cmhi.2d	v16, v5, v16
	uzp1.4s	v16, v16, v17
	cmhi.2d	v17, v5, v19
	cmhi.2d	v18, v5, v18
	uzp1.4s	v17, v18, v17
	cmhi.2d	v18, v5, v21
	cmhi.2d	v19, v5, v20
	uzp1.4s	v18, v19, v18
	uzp1.4s	v19, v1, v0
	eor.16b	v20, v22, v19
	eor.16b	v16, v16, v19
	eor.16b	v17, v17, v19
	eor.16b	v18, v18, v19
	and.16b	v19, v20, v6
	and.16b	v16, v16, v6
	and.16b	v17, v17, v6
	and.16b	v18, v18, v6
	add	x10, x9, x8
	str	q19, [x10, #40000]
	str	q16, [x10, #40016]
	str	q17, [x10, #40032]
	str	q18, [x10, #40048]
	add.2d	v0, v0, v7
	add.2d	v1, v1, v7
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
Lloh9:
	adrp	x0, l_.str@PAGE
Lloh10:
	add	x0, x0, l_.str@PAGEOFF
	bl	_printf
	sub	x8, x20, x19
	str	x8, [sp]
Lloh11:
	adrp	x0, l_.str.1@PAGE
Lloh12:
	add	x0, x0, l_.str.1@PAGEOFF
	bl	_printf
	ldur	x8, [x29, #-40]
Lloh13:
	adrp	x9, ___stack_chk_guard@GOTPAGE
Lloh14:
	ldr	x9, [x9, ___stack_chk_guard@GOTPAGEOFF]
Lloh15:
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
	.loh AdrpLdr	Lloh7, Lloh8
	.loh AdrpAdrp	Lloh5, Lloh7
	.loh AdrpLdr	Lloh5, Lloh6
	.loh AdrpLdrGotLdr	Lloh2, Lloh3, Lloh4
	.loh AdrpLdrGot	Lloh0, Lloh1
	.loh AdrpLdrGotLdr	Lloh13, Lloh14, Lloh15
	.loh AdrpAdd	Lloh11, Lloh12
	.loh AdrpAdd	Lloh9, Lloh10
	.cfi_endproc
                                        ; -- End function
	.globl	_iseven                         ; -- Begin function iseven
	.p2align	2
_iseven:                                ; @iseven
	.cfi_startproc
; %bb.0:
	add	w8, w0, #5
	cmp	w0, #10
	csel	w8, w8, w0, lt
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
