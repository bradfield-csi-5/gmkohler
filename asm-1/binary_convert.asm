; --------------------------------------------------------------------------
; Convert a binary number, represented as a string (e.g. "101010"), to its
; decimal equivalent using first principles.  Given a null-terminated ASCII
; string of '0' and '1' characters, your program will return the equivalent
; integer as output.  For instance, given the input "110" your program should
; return the integer 6.  You can assume the input is well-formed, i.e. that
; it only contains '1' and '0' characters.
; --------------------------------------------------------------------------
section .text
global binary_convert
binary_convert:
	xor	eax, eax		; prepare return value
	xor	sil, sil		; use sil for holding converted bytes
	xor	ecx, ecx		; hold a zero or one based on sil
work:
	mov	sil, byte[rdi]	; load next byte
	cmp	sil, 0			; are we at the end of the string?
	je	done			; if so, we're done
	
	xor 	ecx, ecx
	cmp	sil, 0x31		; are we looking at a '1'?
	sete	cl			; if so, push 1 into the register we're using for numbers

	shl	eax, 1			; push the last-read bit to the left
	or	eax, ecx		; add the new bit (at the least-significant place in the return value)

	inc	rdi
	jmp	work			; continue until we run into a null

	; need to cmp edi[j] to decide when to terminate
	ret
done:
	ret

