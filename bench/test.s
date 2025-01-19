	.section	__TEXT,__text,regular,pure_instructions
	.build_version macos, 15, 0	sdk_version 15, 2
	.globl	_main                           ; -- Begin function main
	.p2align	2
_main:                                  ; @main
	.cfi_startproc
; %bb.0:
	stp	x28, x27, [sp, #-32]!           ; 16-byte Folded Spill
	stp	x29, x30, [sp, #16]             ; 16-byte Folded Spill
	add	x29, sp, #16
	mov	w9, #40064                      ; =0x9c80
	adrp	x16, ___chkstk_darwin@GOTPAGE
	ldr	x16, [x16, ___chkstk_darwin@GOTPAGEOFF]
	blr	x16
	sub	sp, sp, #9, lsl #12             ; =36864
	sub	sp, sp, #3200
	.cfi_def_cfa w29, 16
	.cfi_offset w30, -8
	.cfi_offset w29, -16
	.cfi_offset w27, -24
	.cfi_offset w28, -32
	adrp	x8, ___stack_chk_guard@GOTPAGE
	ldr	x8, [x8, ___stack_chk_guard@GOTPAGEOFF]
	ldr	x8, [x8]
	stur	x8, [x29, #-24]
	str	wzr, [sp, #52]
	bl	_clock
	str	x0, [sp, #40]
	str	wzr, [sp, #28]
	b	LBB0_1
LBB0_1:                                 ; =>This Inner Loop Header: Depth=1
	ldr	w8, [sp, #28]
	mov	w9, #10000                      ; =0x2710
	subs	w8, w8, w9
	cset	w8, ge
	tbnz	w8, #0, LBB0_4
	b	LBB0_2
LBB0_2:                                 ;   in Loop: Header=BB0_1 Depth=1
	ldr	w0, [sp, #28]
	bl	_iseven
	ldrsw	x9, [sp, #28]
	add	x8, sp, #56
	str	w0, [x8, x9, lsl #2]
	b	LBB0_3
LBB0_3:                                 ;   in Loop: Header=BB0_1 Depth=1
	ldr	w8, [sp, #28]
	add	w8, w8, #1
	str	w8, [sp, #28]
	b	LBB0_1
LBB0_4:
	str	wzr, [sp, #24]
	str	wzr, [sp, #20]
	b	LBB0_5
LBB0_5:                                 ; =>This Inner Loop Header: Depth=1
	ldr	w8, [sp, #20]
	mov	w9, #10000                      ; =0x2710
	subs	w8, w8, w9
	cset	w8, ge
	tbnz	w8, #0, LBB0_8
	b	LBB0_6
LBB0_6:                                 ;   in Loop: Header=BB0_5 Depth=1
	ldrsw	x9, [sp, #20]
	add	x8, sp, #56
	ldr	w9, [x8, x9, lsl #2]
	ldr	w8, [sp, #24]
	add	w8, w8, w9
	str	w8, [sp, #24]
	b	LBB0_7
LBB0_7:                                 ;   in Loop: Header=BB0_5 Depth=1
	ldr	w8, [sp, #20]
	add	w8, w8, #1
	str	w8, [sp, #20]
	b	LBB0_5
LBB0_8:
	bl	_clock
	str	x0, [sp, #32]
	ldr	w9, [sp, #24]
                                        ; implicit-def: $x8
	mov	x8, x9
	mov	x9, sp
	mov	x10, #10000                     ; =0x2710
	str	x10, [x9]
	str	x8, [x9, #8]
	adrp	x0, l_.str@PAGE
	add	x0, x0, l_.str@PAGEOFF
	bl	_printf
	ldr	x8, [sp, #32]
	ldr	x9, [sp, #40]
	subs	x8, x8, x9
	mov	x9, sp
	str	x8, [x9]
	adrp	x0, l_.str.1@PAGE
	add	x0, x0, l_.str.1@PAGEOFF
	bl	_printf
	ldr	w8, [sp, #52]
	str	w8, [sp, #16]                   ; 4-byte Folded Spill
	ldur	x9, [x29, #-24]
	adrp	x8, ___stack_chk_guard@GOTPAGE
	ldr	x8, [x8, ___stack_chk_guard@GOTPAGEOFF]
	ldr	x8, [x8]
	subs	x8, x8, x9
	cset	w8, eq
	tbnz	w8, #0, LBB0_10
	b	LBB0_9
LBB0_9:
	bl	___stack_chk_fail
LBB0_10:
	ldr	w0, [sp, #16]                   ; 4-byte Folded Reload
	add	sp, sp, #9, lsl #12             ; =36864
	add	sp, sp, #3200
	ldp	x29, x30, [sp, #16]             ; 16-byte Folded Reload
	ldp	x28, x27, [sp], #32             ; 16-byte Folded Reload
	ret
	.cfi_endproc
                                        ; -- End function
	.globl	_iseven                         ; -- Begin function iseven
	.p2align	2
_iseven:                                ; @iseven
	.cfi_startproc
; %bb.0:
	sub	sp, sp, #16
	.cfi_def_cfa_offset 16
	str	w0, [sp, #8]
	ldr	w8, [sp, #8]
	add	w8, w8, #5
	str	w8, [sp, #4]
	ldr	w8, [sp, #4]
	mov	w10, #2                         ; =0x2
	sdiv	w9, w8, w10
	mul	w9, w9, w10
	subs	w8, w8, w9
	str	w8, [sp]
	ldr	w8, [sp]
	subs	w8, w8, #1
	cset	w8, ne
	tbnz	w8, #0, LBB1_2
	b	LBB1_1
LBB1_1:
	mov	w8, #1                          ; =0x1
	str	w8, [sp, #12]
	b	LBB1_3
LBB1_2:
	str	wzr, [sp, #12]
	b	LBB1_3
LBB1_3:
	ldr	w0, [sp, #12]
	add	sp, sp, #16
	ret
	.cfi_endproc
                                        ; -- End function
	.section	__TEXT,__cstring,cstring_literals
l_.str:                                 ; @.str
	.asciz	"even counter(%d):%d\n"

l_.str.1:                               ; @.str.1
	.asciz	"time elapsed:%lu\n"

.subsections_via_symbols
