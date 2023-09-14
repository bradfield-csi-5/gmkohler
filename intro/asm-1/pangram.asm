; --------------------------------------------------------------------------
; A phrase is a pangram if it contains every character at least once, for
; instance "the quick brown fox jumped over the lazy dog".  Write a program
; to determine if a given phrase is a pangram.  As usuall, some basic tests
; are provided: from these you can see that your program should be
; insensitive case, and ignore punctutation.
;
; This exercise should be interesting as it will encourage you to think about
; simple data structures at a low level.  How will you keep track of which
; letters have already been encountered?  If you were using a high level
; language, you might immediately reach for something like a hash map.
; Writing in assembly, you might find yourself forced to discover a much more
; efficient and elegant approach.
; --------------------------------------------------------------------------

; --------------------------------------------------------------------------
; I am completing this problem after watching the recording so I am going to
; implement it with a 26-bit bitfield.
; Assuming rdi has the input string.
; --------------------------------------------------------------------------
	global pangram
	section .text
pangram:
	xor	esi, esi			; use esi as our bitfield
	xor	edx, edx			; hold character codes in dl
.work:
	movzx	edx, byte [rdi]			; read byte, set the rest of register to zero
	cmp 	edx, 0				; check whether we have null character
	je	.done				; if so, prepare return value based on what we've read so far.

	or	edx, 0x20			; case insensitivity: ascii has 'A' and 'a' 32 bits apart so this will map both to 'a'
	sub	edx, 0x61			; set 'a' as lowest bit in our byte array
	bts	esi, edx

	inc	rdi				; prepare for next part of loop
	jmp 	.work
.done:
	xor	rax, rax			; clean out return register
	and	esi, 0x4000000 - 1 		; zero out characters above 'z'
	cmp	esi, 0x4000000 - 1 		; is bitfield equal to 2**26 - 1 (all ones)?
	sete	al				; set our result to whether the two are equal
	ret

