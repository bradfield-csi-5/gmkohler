	global fib
	section .text
fib:
	; save registers
	push	rbx			; use for n
	push	rbp			; use for fib(n-1)

	mov	rbx, rdi		; save n to prepare args for later fibs (or return it)
	cmp	rbx, 1			; compare n with 1
	jle	.baseCase		; go to base case when le

	lea	rdi, [rbx - 1]		; prepare n-1 argument
	call 	fib			; call fib(n-1)
	mov	rbp, rax		; save fib(n-1) in our callee-saved register

	lea	rdi, [rbx - 2]		; prepare n-2 argument
	call	fib			; call fib(n-2)
	add	rax, rbp		; add fib(n-1) to the result of fib(n-2)
.done:
	pop	rbp			; restore registers
	pop	rbx
	ret
.baseCase:
	mov	rax, rbx		; move callee-saved n to return value
	jmp	.done
