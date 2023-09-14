	.section	__TEXT,__text,regular,pure_instructions
	.build_version macos, 12, 0	sdk_version 13, 1
	.globl	_pagecount                      ## -- Begin function pagecount
	.p2align	4, 0x90
_pagecount:                             ## @pagecount
	.cfi_startproc
## %bb.0:
	pushq	%rbp
	.cfi_def_cfa_offset 16
	.cfi_offset %rbp, -16
	movq	%rsp, %rbp
	.cfi_def_cfa_register %rbp
	movq	%rdi, -8(%rbp)
	movq	%rsi, -16(%rbp)
	movq	-8(%rbp), %rax
	xorl	%ecx, %ecx
	movl	%ecx, %edx
	divq	-16(%rbp)
	popq	%rbp
	retq
	.cfi_endproc
                                        ## -- End function
	.section	__TEXT,__literal8,8byte_literals
	.p2align	3                               ## -- Begin function main
LCPI1_0:
	.quad	0x416312d000000000              ## double 1.0E+7
LCPI1_1:
	.quad	0x41cdcd6500000000              ## double 1.0E+9
LCPI1_2:
	.quad	0x412e848000000000              ## double 1.0E+6
	.section	__TEXT,__literal16,16byte_literals
	.p2align	4
LCPI1_3:
	.long	1127219200                      ## 0x43300000
	.long	1160773632                      ## 0x45300000
	.long	0                               ## 0x0
	.long	0                               ## 0x0
LCPI1_4:
	.quad	0x4330000000000000              ## double 4503599627370496
	.quad	0x4530000000000000              ## double 1.9342813113834067E+25
	.section	__TEXT,__text,regular,pure_instructions
	.globl	_main
	.p2align	4, 0x90
_main:                                  ## @main
	.cfi_startproc
## %bb.0:
	pushq	%rbp
	.cfi_def_cfa_offset 16
	.cfi_offset %rbp, -16
	movq	%rsp, %rbp
	.cfi_def_cfa_register %rbp
	subq	$160, %rsp
	movq	___stack_chk_guard@GOTPCREL(%rip), %rax
	movq	(%rax), %rax
	movq	%rax, -8(%rbp)
	movl	$0, -68(%rbp)
	movl	%edi, -72(%rbp)
	movq	%rsi, -80(%rbp)
	movl	$0, -152(%rbp)
	movq	l___const.main.msizes(%rip), %rax
	movq	%rax, -32(%rbp)
	movq	l___const.main.msizes+8(%rip), %rax
	movq	%rax, -24(%rbp)
	movq	l___const.main.msizes+16(%rip), %rax
	movq	%rax, -16(%rbp)
	movq	l___const.main.psizes(%rip), %rax
	movq	%rax, -64(%rbp)
	movq	l___const.main.psizes+8(%rip), %rax
	movq	%rax, -56(%rbp)
	movq	l___const.main.psizes+16(%rip), %rax
	movq	%rax, -48(%rbp)
	callq	_clock
	movq	%rax, -88(%rbp)
	movl	$0, -148(%rbp)
LBB1_1:                                 ## =>This Inner Loop Header: Depth=1
	cmpl	$10000000, -148(%rbp)           ## imm = 0x989680
	jge	LBB1_4
## %bb.2:                               ##   in Loop: Header=BB1_1 Depth=1
	movl	-148(%rbp), %eax
	movl	$3, %ecx
	cltd
	idivl	%ecx
	movslq	%edx, %rax
	movq	-32(%rbp,%rax,8), %rax
	movq	%rax, -120(%rbp)
	movl	-148(%rbp), %eax
	movl	$3, %ecx
	cltd
	idivl	%ecx
	movslq	%edx, %rax
	movq	-64(%rbp,%rax,8), %rax
	movq	%rax, -128(%rbp)
	movq	-120(%rbp), %rcx
	addq	$1, %rcx
	addq	-128(%rbp), %rcx
	movslq	-152(%rbp), %rax
	addq	%rcx, %rax
                                        ## kill: def $eax killed $eax killed $rax
	movl	%eax, -152(%rbp)
## %bb.3:                               ##   in Loop: Header=BB1_1 Depth=1
	movl	-148(%rbp), %eax
	addl	$1, %eax
	movl	%eax, -148(%rbp)
	jmp	LBB1_1
LBB1_4:
	callq	_clock
	movq	%rax, -96(%rbp)
	callq	_clock
	movq	%rax, -104(%rbp)
	movl	$0, -148(%rbp)
LBB1_5:                                 ## =>This Inner Loop Header: Depth=1
	cmpl	$10000000, -148(%rbp)           ## imm = 0x989680
	jge	LBB1_8
## %bb.6:                               ##   in Loop: Header=BB1_5 Depth=1
	movl	-148(%rbp), %eax
	movl	$3, %ecx
	cltd
	idivl	%ecx
	movslq	%edx, %rax
	movq	-32(%rbp,%rax,8), %rax
	movq	%rax, -120(%rbp)
	movl	-148(%rbp), %eax
	movl	$3, %ecx
	cltd
	idivl	%ecx
	movslq	%edx, %rax
	movq	-64(%rbp,%rax,8), %rax
	movq	%rax, -128(%rbp)
	movq	-120(%rbp), %rdi
	movq	-128(%rbp), %rsi
	callq	_pagecount
	movq	%rax, %rcx
	addq	-120(%rbp), %rcx
	addq	-128(%rbp), %rcx
	movslq	-152(%rbp), %rax
	addq	%rcx, %rax
                                        ## kill: def $eax killed $eax killed $rax
	movl	%eax, -152(%rbp)
## %bb.7:                               ##   in Loop: Header=BB1_5 Depth=1
	movl	-148(%rbp), %eax
	addl	$1, %eax
	movl	%eax, -148(%rbp)
	jmp	LBB1_5
LBB1_8:
	callq	_clock
	movq	%rax, -112(%rbp)
	movq	-112(%rbp), %rax
	movq	-104(%rbp), %rcx
	subq	%rcx, %rax
	movq	-96(%rbp), %rdx
	movq	-88(%rbp), %rcx
	subq	%rdx, %rcx
	addq	%rcx, %rax
	movq	%rax, %xmm0
	movaps	LCPI1_3(%rip), %xmm1            ## xmm1 = [1127219200,1160773632,0,0]
	punpckldq	%xmm1, %xmm0            ## xmm0 = xmm0[0],xmm1[0],xmm0[1],xmm1[1]
	movapd	LCPI1_4(%rip), %xmm1            ## xmm1 = [4.503599627370496E+15,1.9342813113834067E+25]
	subpd	%xmm1, %xmm0
	movaps	%xmm0, %xmm1
	unpckhpd	%xmm0, %xmm0                    ## xmm0 = xmm0[1,1]
	addsd	%xmm1, %xmm0
	movsd	%xmm0, -136(%rbp)
	movsd	-136(%rbp), %xmm0               ## xmm0 = mem[0],zero
	movsd	LCPI1_2(%rip), %xmm1            ## xmm1 = mem[0],zero
	divsd	%xmm1, %xmm0
	movsd	%xmm0, -144(%rbp)
	movsd	-144(%rbp), %xmm0               ## xmm0 = mem[0],zero
	movsd	LCPI1_1(%rip), %xmm1            ## xmm1 = mem[0],zero
	mulsd	-144(%rbp), %xmm1
	movsd	LCPI1_0(%rip), %xmm2            ## xmm2 = mem[0],zero
	divsd	%xmm2, %xmm1
	leaq	L_.str(%rip), %rdi
	movl	$10000000, %esi                 ## imm = 0x989680
	movb	$2, %al
	callq	_printf
	movl	-152(%rbp), %eax
	movl	%eax, -156(%rbp)                ## 4-byte Spill
	movq	___stack_chk_guard@GOTPCREL(%rip), %rax
	movq	(%rax), %rax
	movq	-8(%rbp), %rcx
	cmpq	%rcx, %rax
	jne	LBB1_10
## %bb.9:
	movl	-156(%rbp), %eax                ## 4-byte Reload
	addq	$160, %rsp
	popq	%rbp
	retq
LBB1_10:
	callq	___stack_chk_fail
	ud2
	.cfi_endproc
                                        ## -- End function
	.section	__TEXT,__const
	.p2align	4                               ## @__const.main.msizes
l___const.main.msizes:
	.quad	4294967296                      ## 0x100000000
	.quad	1099511627776                   ## 0x10000000000
	.quad	4503599627370496                ## 0x10000000000000

	.p2align	4                               ## @__const.main.psizes
l___const.main.psizes:
	.quad	4096                            ## 0x1000
	.quad	65536                           ## 0x10000
	.quad	4294967296                      ## 0x100000000

	.section	__TEXT,__cstring,cstring_literals
L_.str:                                 ## @.str
	.asciz	"%.2fs to run %d tests (%.2fns per test)\n"

.subsections_via_symbols
