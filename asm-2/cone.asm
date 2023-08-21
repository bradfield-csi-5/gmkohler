; --------------------------------------------------------------------------
; Write a function that calculates the volume of a cone, given its radius
; and height as the first two arguments.
; 
; Note that both arguments are given as floating points, and the return
; value is expected to be a floating point.  This is after all an exercise
; on floating points.
;
; For the value of pi, feel free to hardcode a resonable approximation
; somewhere in your program.
; --------------------------------------------------------------------------
; --------------------------------------------------------------------------
; The formula for a cone is pi * r**2 * h / 3.
; We have r in xmm0 and h in xmm1;  We can do some math on floats
; (section 3.11.3) and encode a constant for pi (section 3.11.4) to implement
; this program.
; --------------------------------------------------------------------------

	default rel
	global volume
	section .text
volume:
	mulss		xmm0, xmm0		; r**2
	mulss		xmm0, xmm1		; r**2 * h
	mulss		xmm0, [pi]		; r ** 2 * h * pi
	; build 3 as a float
	mov		esi, 3			; add to integer register
	; I had to search around for this instruction, it's not on félix cloutier's site.
	vcvtsi2ss	xmm2, esi		; move to xmm register
	divss		xmm0, xmm2		; divide running reuslt by three
 	ret

	section .rodata
pi:	dd	3.141592653589793238462

