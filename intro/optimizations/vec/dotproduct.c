#include "vec.h"

/**
 * This seems like a candidate for unrolling (CS:APP §5.8), especially
 * because the vector is longs so we don't risk issues with combining the
 * two.
 *
 * We are also given get_vec_start in vec.h, which looks like a way to
 * avoid reading memory every time (CS:APP §5.6).
 *
 * We can also avoid reading vec_length in each iteration by reading once
 * and storing it in a local variable (CS:APP §5.6).
 */
#define K 4
data_t dotproduct(vec_ptr u, vec_ptr v) {
   data_t acc0 = 0, acc1 = 0, acc2 = 0, acc3 = 0;

   data_t *u_start = get_vec_start(u);
   data_t *v_start = get_vec_start(v);
   long length = vec_length(u);
   long i;

   for (i = 0; i < length - (K - 1); i+= K) { // we can assume both vectors are same length
	acc0 += u_start[i] * v_start[i];
	acc1 += u_start[i + 1] * v_start[i + 1];
	acc2 += u_start[i + 2] * v_start[i + 2];
	acc3 += u_start[i + 3] * v_start[i + 3];
   }   
   for(; i < length; i++) {
	acc0 += u_start[i] * v_start[i];
   }

   return acc0 + acc1 + acc2 + acc3;
}
