main.main STEXT size=55 args=0x0 locals=0x10 funcid=0x0 align=0x0
	0x0000 00000 (/Users/maxgara/Desktop/go-code/gobook/workspace/code-analysis/asm/test.go:3)	TEXT	main.main(SB), ABIInternal, $16-0
	0x0000 00000 (/Users/maxgara/Desktop/go-code/gobook/workspace/code-analysis/asm/test.go:3)	CMPQ	SP, 16(R14)
	0x0004 00004 (/Users/maxgara/Desktop/go-code/gobook/workspace/code-analysis/asm/test.go:3)	PCDATA	$0, $-2
	0x0004 00004 (/Users/maxgara/Desktop/go-code/gobook/workspace/code-analysis/asm/test.go:3)	JLS	48
	0x0006 00006 (/Users/maxgara/Desktop/go-code/gobook/workspace/code-analysis/asm/test.go:3)	PCDATA	$0, $-1
	0x0006 00006 (/Users/maxgara/Desktop/go-code/gobook/workspace/code-analysis/asm/test.go:3)	PUSHQ	BP
	0x0007 00007 (/Users/maxgara/Desktop/go-code/gobook/workspace/code-analysis/asm/test.go:3)	MOVQ	SP, BP
	0x000a 00010 (/Users/maxgara/Desktop/go-code/gobook/workspace/code-analysis/asm/test.go:3)	SUBQ	$8, SP
	0x000e 00014 (/Users/maxgara/Desktop/go-code/gobook/workspace/code-analysis/asm/test.go:3)	FUNCDATA	$0, gclocals·g2BeySu+wFnoycgXfElmcg==(SB)
	0x000e 00014 (/Users/maxgara/Desktop/go-code/gobook/workspace/code-analysis/asm/test.go:3)	FUNCDATA	$1, gclocals·g2BeySu+wFnoycgXfElmcg==(SB)
	0x000e 00014 (/Users/maxgara/Desktop/go-code/gobook/workspace/code-analysis/asm/test.go:4)	PCDATA	$1, $0
	0x000e 00014 (/Users/maxgara/Desktop/go-code/gobook/workspace/code-analysis/asm/test.go:4)	CALL	runtime.printlock(SB)
	0x0013 00019 (/Users/maxgara/Desktop/go-code/gobook/workspace/code-analysis/asm/test.go:4)	MOVL	$3, AX
	0x0018 00024 (/Users/maxgara/Desktop/go-code/gobook/workspace/code-analysis/asm/test.go:4)	CALL	runtime.printint(SB)
	0x001d 00029 (/Users/maxgara/Desktop/go-code/gobook/workspace/code-analysis/asm/test.go:4)	NOP
	0x0020 00032 (/Users/maxgara/Desktop/go-code/gobook/workspace/code-analysis/asm/test.go:4)	CALL	runtime.printnl(SB)
	0x0025 00037 (/Users/maxgara/Desktop/go-code/gobook/workspace/code-analysis/asm/test.go:4)	CALL	runtime.printunlock(SB)
	0x002a 00042 (/Users/maxgara/Desktop/go-code/gobook/workspace/code-analysis/asm/test.go:5)	ADDQ	$8, SP
	0x002e 00046 (/Users/maxgara/Desktop/go-code/gobook/workspace/code-analysis/asm/test.go:5)	POPQ	BP
	0x002f 00047 (/Users/maxgara/Desktop/go-code/gobook/workspace/code-analysis/asm/test.go:5)	RET
	0x0030 00048 (/Users/maxgara/Desktop/go-code/gobook/workspace/code-analysis/asm/test.go:5)	NOP
	0x0030 00048 (/Users/maxgara/Desktop/go-code/gobook/workspace/code-analysis/asm/test.go:3)	PCDATA	$1, $-1
	0x0030 00048 (/Users/maxgara/Desktop/go-code/gobook/workspace/code-analysis/asm/test.go:3)	PCDATA	$0, $-2
	0x0030 00048 (/Users/maxgara/Desktop/go-code/gobook/workspace/code-analysis/asm/test.go:3)	CALL	runtime.morestack_noctxt(SB)
	0x0035 00053 (/Users/maxgara/Desktop/go-code/gobook/workspace/code-analysis/asm/test.go:3)	PCDATA	$0, $-1
	0x0035 00053 (/Users/maxgara/Desktop/go-code/gobook/workspace/code-analysis/asm/test.go:3)	JMP	0
	0x0000 49 3b 66 10 76 2a 55 48 89 e5 48 83 ec 08 e8 00  I;f.v*UH..H.....
	0x0010 00 00 00 b8 03 00 00 00 e8 00 00 00 00 0f 1f 00  ................
	0x0020 e8 00 00 00 00 e8 00 00 00 00 48 83 c4 08 5d c3  ..........H...].
	0x0030 e8 00 00 00 00 eb c9                             .......
	rel 15+4 t=R_CALL runtime.printlock+0
	rel 25+4 t=R_CALL runtime.printint+0
	rel 33+4 t=R_CALL runtime.printnl+0
	rel 38+4 t=R_CALL runtime.printunlock+0
	rel 49+4 t=R_CALL runtime.morestack_noctxt+0
go:cuinfo.producer.<unlinkable> SDWARFCUINFO dupok size=0
	0x0000 72 65 67 61 62 69                                regabi
go:cuinfo.packagename.main SDWARFCUINFO dupok size=0
	0x0000 6d 61 69 6e                                      main
main..inittask SNOPTRDATA size=8
	0x0000 00 00 00 00 00 00 00 00                          ........
gclocals·g2BeySu+wFnoycgXfElmcg== SRODATA dupok size=8
	0x0000 01 00 00 00 00 00 00 00                          ........
