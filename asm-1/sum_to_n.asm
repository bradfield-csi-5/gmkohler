; * Which register contains the argument (n)?
;   I am assuming %edi until I learn otherwise because it's
;   the int-sized first argument register
; * In which register will the return value be placed?
;   I am going to use the int-sized one so %eax
; * What instruction is used to do the addition?
;   addl should be able to be used
; * What instruction(s) can be used to construct the loop?
;   jumping when comparing a temp variable to %rdi should let
;   us do this.
;   This program uses the jump-to-middle aprroach described in CS:APP
	section .text
	global sum_to_n
sum_to_n:
	xor	eax, eax	; clear out eax to use for return value
	xor	esi, esi	; clear out esi to use for iterator
	jmp	test
work:
	inc	esi		; increase iterator
	add	eax, esi	; add new value to sum
test:
	cmp	esi, edi	; is i < n?
	jl	work		; if so, continue working
done:
	ret			; we're done
