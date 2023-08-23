
pagecount_o2.o:	file format mach-o 64-bit x86-64

Disassembly of section __TEXT,__text:

0000000000000000 <_pagecount>:
       0: 55                           	push	rbp
       1: 48 89 e5                     	mov	rbp, rsp
       4: 48 89 f8                     	mov	rax, rdi
       7: 48 89 f9                     	mov	rcx, rdi
       a: 48 09 f1                     	or	rcx, rsi
       d: 48 c1 e9 20                  	shr	rcx, 32
      11: 74 07                        	je	0x1a <_pagecount+0x1a>
      13: 31 d2                        	xor	edx, edx
      15: 48 f7 f6                     	div	rsi
      18: 5d                           	pop	rbp
      19: c3                           	ret
      1a: 31 d2                        	xor	edx, edx
      1c: f7 f6                        	div	esi
      1e: 5d                           	pop	rbp
      1f: c3                           	ret

0000000000000020 <_main>:
      20: 55                           	push	rbp
      21: 48 89 e5                     	mov	rbp, rsp
      24: 41 57                        	push	r15
      26: 41 56                        	push	r14
      28: 41 55                        	push	r13
      2a: 41 54                        	push	r12
      2c: 53                           	push	rbx
      2d: 48 83 ec 18                  	sub	rsp, 24
      31: 45 31 ff                     	xor	r15d, r15d
      34: 41 be 08 00 00 00            	mov	r14d, 8
      3a: bb 01 00 00 00               	mov	ebx, 1
      3f: 41 bc 80 96 98 00            	mov	r12d, 10000000
      45: e8 00 00 00 00               	call	0x4a <_main+0x2a>
      4a: 48 89 45 c8                  	mov	qword ptr [rbp - 56], rax
      4e: 49 b8 ab aa aa aa aa aa aa aa	movabs	r8, -6148914691236517205
      58: 4c 8d 0d 00 00 00 00         	lea	r9, [rip]               ## 0x5f <_main+0x3f>
      5f: 48 8d 3d 00 00 00 00         	lea	rdi, [rip]              ## 0x66 <_main+0x46>
      66: 31 c9                        	xor	ecx, ecx
      68: 45 31 ed                     	xor	r13d, r13d
      6b: 0f 1f 44 00 00               	nop	dword ptr [rax + rax]
      70: 48 89 d8                     	mov	rax, rbx
      73: 49 f7 e0                     	mul	r8
      76: 48 c1 e2 02                  	shl	rdx, 2
      7a: 48 83 e2 f8                  	and	rdx, -8
      7e: 48 8d 04 52                  	lea	rax, [rdx + 2*rdx]
      82: 4c 89 f6                     	mov	rsi, r14
      85: 48 29 c6                     	sub	rsi, rax
      88: 4c 89 f8                     	mov	rax, r15
      8b: 49 f7 e0                     	mul	r8
      8e: 48 c1 e2 02                  	shl	rdx, 2
      92: 48 83 e2 f8                  	and	rdx, -8
      96: 48 8d 04 52                  	lea	rax, [rdx + 2*rdx]
      9a: 48 89 ca                     	mov	rdx, rcx
      9d: 48 29 c2                     	sub	rdx, rax
      a0: 41 8b 04 11                  	mov	eax, dword ptr [r9 + rdx]
      a4: 03 04 17                     	add	eax, dword ptr [rdi + rdx]
      a7: 44 01 e8                     	add	eax, r13d
      aa: 41 8b 14 31                  	mov	edx, dword ptr [r9 + rsi]
      ae: 03 14 37                     	add	edx, dword ptr [rdi + rsi]
      b1: 44 8d 2c 10                  	lea	r13d, [rax + rdx]
      b5: 41 83 c5 02                  	add	r13d, 2
      b9: 49 83 c6 10                  	add	r14, 16
      bd: 48 83 c3 02                  	add	rbx, 2
      c1: 48 83 c1 10                  	add	rcx, 16
      c5: 49 83 c7 02                  	add	r15, 2
      c9: 41 83 c4 fe                  	add	r12d, -2
      cd: 75 a1                        	jne	0x70 <_main+0x50>
      cf: 4d 89 cc                     	mov	r12, r9
      d2: 49 89 fe                     	mov	r14, rdi
      d5: e8 00 00 00 00               	call	0xda <_main+0xba>
      da: 48 89 45 d0                  	mov	qword ptr [rbp - 48], rax
      de: 41 bf 80 96 98 00            	mov	r15d, 10000000
      e4: 31 db                        	xor	ebx, ebx
      e6: e8 00 00 00 00               	call	0xeb <_main+0xcb>
      eb: 4d 89 f2                     	mov	r10, r14
      ee: 4d 89 e1                     	mov	r9, r12
      f1: 49 b8 ab aa aa aa aa aa aa aa	movabs	r8, -6148914691236517205
      fb: 49 89 c6                     	mov	r14, rax
      fe: 31 f6                        	xor	esi, esi
     100: eb 2b                        	jmp	0x12d <_main+0x10d>
     102: 66 2e 0f 1f 84 00 00 00 00 00	nop	word ptr cs:[rax + rax]
     10c: 0f 1f 40 00                  	nop	dword ptr [rax]
     110: 48 89 c8                     	mov	rax, rcx
     113: 31 d2                        	xor	edx, edx
     115: 48 f7 f7                     	div	rdi
     118: 01 cf                        	add	edi, ecx
     11a: 01 f8                        	add	eax, edi
     11c: 41 01 c5                     	add	r13d, eax
     11f: 48 83 c6 08                  	add	rsi, 8
     123: 48 83 c3 01                  	add	rbx, 1
     127: 41 83 c7 ff                  	add	r15d, -1
     12b: 74 34                        	je	0x161 <_main+0x141>
     12d: 48 89 d8                     	mov	rax, rbx
     130: 49 f7 e0                     	mul	r8
     133: 48 c1 e2 02                  	shl	rdx, 2
     137: 48 83 e2 f8                  	and	rdx, -8
     13b: 48 8d 04 52                  	lea	rax, [rdx + 2*rdx]
     13f: 48 89 f2                     	mov	rdx, rsi
     142: 48 29 c2                     	sub	rdx, rax
     145: 4a 8b 0c 0a                  	mov	rcx, qword ptr [rdx + r9]
     149: 4a 8b 3c 12                  	mov	rdi, qword ptr [rdx + r10]
     14d: 48 89 c8                     	mov	rax, rcx
     150: 48 09 f8                     	or	rax, rdi
     153: 48 c1 e8 20                  	shr	rax, 32
     157: 75 b7                        	jne	0x110 <_main+0xf0>
     159: 89 c8                        	mov	eax, ecx
     15b: 31 d2                        	xor	edx, edx
     15d: f7 f7                        	div	edi
     15f: eb b7                        	jmp	0x118 <_main+0xf8>
     161: e8 00 00 00 00               	call	0x166 <_main+0x146>
     166: 48 8b 55 d0                  	mov	rdx, qword ptr [rbp - 48]
     16a: 4c 01 f2                     	add	rdx, r14
     16d: 48 8b 4d c8                  	mov	rcx, qword ptr [rbp - 56]
     171: 48 29 d1                     	sub	rcx, rdx
     174: 48 01 c1                     	add	rcx, rax
     177: 66 48 0f 6e c9               	movq	xmm1, rcx
     17c: 66 0f 62 0d 5c 00 00 00      	punpckldq	xmm1, xmmword ptr [rip + 92] ## xmm1 = xmm1[0],mem[0],xmm1[1],mem[1]
                                                                        ## 0x1e0 <_main+0x1c0>
     184: 66 0f 5c 0d 64 00 00 00      	subpd	xmm1, xmmword ptr [rip + 100] ## 0x1f0 <_main+0x1d0>
     18c: 66 0f 28 c1                  	movapd	xmm0, xmm1
     190: 66 0f 15 c1                  	unpckhpd	xmm0, xmm1              ## xmm0 = xmm0[1],xmm1[1]
     194: f2 0f 58 c1                  	addsd	xmm0, xmm1
     198: f2 0f 5e 05 60 00 00 00      	divsd	xmm0, qword ptr [rip + 96] ## 0x200 <_main+0x1e0>
     1a0: f2 0f 10 0d 60 00 00 00      	movsd	xmm1, qword ptr [rip + 96] ## xmm1 = mem[0],zero
                                                                        ## 0x208 <_main+0x1e8>
     1a8: f2 0f 59 c8                  	mulsd	xmm1, xmm0
     1ac: f2 0f 5e 0d 5c 00 00 00      	divsd	xmm1, qword ptr [rip + 92] ## 0x210 <_main+0x1f0>
     1b4: 48 8d 3d 9d 00 00 00         	lea	rdi, [rip + 157]        ## 0x258 <_main+0x238>
     1bb: be 80 96 98 00               	mov	esi, 10000000
     1c0: b0 02                        	mov	al, 2
     1c2: e8 00 00 00 00               	call	0x1c7 <_main+0x1a7>
     1c7: 44 89 e8                     	mov	eax, r13d
     1ca: 48 83 c4 18                  	add	rsp, 24
     1ce: 5b                           	pop	rbx
     1cf: 41 5c                        	pop	r12
     1d1: 41 5d                        	pop	r13
     1d3: 41 5e                        	pop	r14
     1d5: 41 5f                        	pop	r15
     1d7: 5d                           	pop	rbp
     1d8: c3                           	ret
