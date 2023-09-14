/**
 * The purpose of this exercise is to demonstrate how a programmer's basic
 * understanding of their own program can give them a significant advantage
 * over an optimizing compiler.  It will also provide an opportunity to use
 * a dissasembler and basic measurement tools.
 *
 * We've provided an incredibly simple program in this file.  Modern
 * operating systems divide memory notionally into units called "pages".
 * If each page is 4Kb, say, and we have 4Gb of pages, say, that means we
 * have approximately 1 million pages of memory.  The pagecount function in
 * this file does this calculation simply by dividing the available memory
 * (in bytes) by the page size (also in bytes).
 *
 * As a first step, read the code for this function then consider:
 *
 * * Which instructions would you expect your compiler to generate for this
 *   function?
 *   
 *   read memory to register
 *   read page to register
 *   divide memory and page
 *
 * * What does it in fact generate?
 *   
 *   pagecount:
 *   	push	rbp 				; callee-save
 *   	mov	rbp, rsp			; use stack pointer
 *   	mov	qword ptr[rbp - 8], rdi		; mov mem_size onto stack
 *   	mov	qword ptr[rbp - 16], rsi	; mov page_size onto stack
 *   	mov	rax, qword ptr[rbp -8]		; move mem_size to return register
 *   	xor	ecx, ecx			; clear ecx
 *   	mov	edx, ecx			; set edx=0
 *   	div	qword ptr[rbp - 16]		; divide mem_size by page_size	
 *   	pop	rbp				; restore callee-save
 *   	ret
 *   	nop	word ptr [rax + rax]		; unsure
 *
 *
 * * If you change the optimization level, is the function substantially
 *   different?
 *
 *   Using -O2 optimization we see the compiler try to avoid 64-bit division
 *   in favor of 32-bit division when possible.  Presumably this is quicker.
 *
 *   pagecount:
 *   0:		push	rbp			; callee-save
 *   1:		mov	rbp, rsp		; move stack pointer to rbp
 *   4:		mov	rax, rdi		; move mem_size to return register
 *   7:		mov	rcx, rdi		; move mem_size to rcx
 *   a:		or	rcx, rsi		; rcx |= page_size
 *   d:		shr	rcx, 32			; rcx >> 32
 *   11:	je	0x1a 			; jumps to  if the 32nd bit is 1
 *   13:	xor	edx, edx		; clear edx
 *   15:	div	rsi			; mem_size / page_size
 *   18:	pop	rbp			; restore callee-save
 *   19:	ret				
 *   1a:	xor	edx, edx		; clear edx
 *   1c:	div	esi			; rax / 32 bits of page_size (rax was shifted 32 bits)
 *   1e:	pop	rbp			; restore callee-size
 *   1f:	ret
 *
 * * Use godbolt.org to explore a few different compilers.  Do any of them
 *   generate substantially different instructions?
 *
 *   I'm seeing pretty similar approaches to my gcc across the x86-64
 *   compilers
 *
 * * By using Agner Fog's instruction tables or reviewing CS:APP ch5.7, can
 *   you determine which of the generated instructions may be slow?
 *
 *   division is slow, it takes 3-30 cycles per CS:APP Fig. 5.12.
 *
 * Next, let's improve performance!
 *
 * * Noting that a page size is always a power of 2, and that the size of
 *   memory will always be divisible by the page size, can you think of a
 *   performance optimization we could employ?  You are welcome to change
 *   the function signature and test runner code.
 *
 *   I reckon we can employ shifting instead of division which would
 *   be quicker.  Using a builtin function we can find the operand
 *   more quickly than would be achievable by manual shifting.
 *
 * * How much of an improvement would you expect to see?
 * 
 *   We should reduce the pagecount time by at least 3 cycles per loop
 *   if our divisions were quick, and 30 cyles per loop if they were slow.
 *   My processor is 2.7 GHz so this would be 1-10 ns per loop
 *
 * * Go ahead and make the improvement, and measure the speed up.  Did this
 *   match your expectations?
 *
 *   I'm seeing a drop of about 5.6 ns per test, so this is in line with my
 *   expectations.
 *
 * * Consider, what is stopping the compiler from making the same
 *   optimization that you did?
 *
 *   The compiler isn't going to assume anything about the nature of the
 *   input
 */
#include <stdint.h>
#include <stdio.h>
#include <time.h>

#define TEST_LOOPS 10000000

uint64_t pagecount(uint64_t memory_size, uint64_t page_size) {
  int page_power = __builtin_ctzl(page_size) + 1;
  return memory_size >> page_power;
}

int main (int argc, char** argv) {
  clock_t baseline_start, baseline_end, test_start, test_end;
  uint64_t memory_size, page_size;
  double clocks_elapsed, time_elapsed;
  int i, ignore = 0;

  uint64_t msizes[] = {1L << 32, 1L << 40, 1L << 52};
  uint64_t psizes[] = {1L << 12, 1L << 16, 1L << 32};

  baseline_start = clock();
  for (i = 0; i < TEST_LOOPS; i++) {
    memory_size = msizes[i % 3];
    page_size = psizes[i % 3];
    ignore += 1 + memory_size +
              page_size; // so that this loop isn't just optimized away
  }
  baseline_end = clock();

  test_start = clock();
  for (i = 0; i < TEST_LOOPS; i++) {
    memory_size = msizes[i % 3];
    page_size = psizes[i % 3];
    ignore += pagecount(memory_size, page_size) + memory_size + page_size;
  }
  test_end = clock();

  clocks_elapsed = test_end - test_start - (baseline_end - baseline_start);
  time_elapsed = clocks_elapsed / CLOCKS_PER_SEC;

  printf("%.2fs to run %d tests (%.2fns per test)\n", time_elapsed, TEST_LOOPS,
         time_elapsed * 1e9 / TEST_LOOPS);
  return ignore;
}

