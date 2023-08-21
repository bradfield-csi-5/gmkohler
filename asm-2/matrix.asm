; --------------------------------------------------------------------------
; Arrays are implemented as contiguous typed data, which can thereby be
; indexed in constant time.  If we know that an array of integers starts at
; location x, we can access arr[5] at a low level by accessing
; byte x + 4 * 5 (since integers are 4 bytes wide).
;
; In this exercise, you will be given a pointer to matrix (an array of
; arrays) as well as its dimensions, and row and column indexes.  Your task
; will be simply to read the value at the given index.  For instance, given
; the following matrix:
;
;	[[1, 2, 3],
;	 [4, 5, 6]]
;
; Given the row and column indices 1 and 2, you should return the value 6.
; Note that we are using the computing convention of zero-indexing, rather
; than the math convention of one-indexing.
; --------------------------------------------------------------------------
	global index
	section .text
index:
	; rdi: matrix
	; rsi: rows
	; rdx: cols
	; rcx: rindex
	; r8: cindex

	; Strategy: compute the offset needed, then access the memory and write
	; it to the return register
	; Equation 3.1 in CS:APP (page 258) describes how to access row-major
	; arrays, I will apply it here.
	mov 	r10, rcx		; (rowIndex)
	imul	r10, rdx		; (rowIndex*cols)
	add	r10, r8			; (rowIndex*cols + colIndex)
	mov	eax, [rdi+r10*4]	; read value
	ret
