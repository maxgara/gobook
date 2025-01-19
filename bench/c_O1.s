	.section	__TEXT,__text,regular,pure_instructions
	.build_version macos, 15, 0	sdk_version 15, 2
	.globl	_main                           ; -- Begin function main
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
	mov	x8, #0                          ; =0x0
	add	x9, sp, #24
	mov	w10, #10000                     ; =0x2710
LBB0_1:                                 ; =>This Inner Loop Header: Depth=1
	cmp	x8, #10
	cset	w11, lo
	eor	w11, w11, w8
	and	w11, w11, #0x1
	str	w11, [x9, x8, lsl #2]
	add	x8, x8, #1
	cmp	x8, x10
	b.ne	LBB0_1
; %bb.2:
	mov	x8, #0                          ; =0x0
	mov	w21, #0                         ; =0x0
	add	x9, sp, #24
	mov	w10, #40000                     ; =0x9c40
LBB0_3:                                 ; =>This Inner Loop Header: Depth=1
	ldr	w11, [x9, x8]
	add	w21, w11, w21
	add	x8, x8, #4
	cmp	x8, x10
	b.ne	LBB0_3
; %bb.4:
	bl	_clock
	mov	x20, x0
	mov	w8, #10000                      ; =0x2710
	stp	x8, x21, [sp]
Lloh5:
	adrp	x0, l_.str@PAGE
Lloh6:
	add	x0, x0, l_.str@PAGEOFF
	bl	_printf
	sub	x8, x20, x19
	str	x8, [sp]
Lloh7:
	adrp	x0, l_.str.1@PAGE
Lloh8:
	add	x0, x0, l_.str.1@PAGEOFF
	bl	_printf
	ldur	x8, [x29, #-40]
Lloh9:
	adrp	x9, ___stack_chk_guard@GOTPAGE
Lloh10:
	ldr	x9, [x9, ___stack_chk_guard@GOTPAGEOFF]
Lloh11:
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
	.loh AdrpLdrGotLdr	Lloh2, Lloh3, Lloh4
	.loh AdrpLdrGot	Lloh0, Lloh1
	.loh AdrpLdrGotLdr	Lloh9, Lloh10, Lloh11
	.loh AdrpAdd	Lloh7, Lloh8
	.loh AdrpAdd	Lloh5, Lloh6
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
