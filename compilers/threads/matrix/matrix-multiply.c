#include <stdio.h>
#include <stdlib.h>
#include <pthread.h>

/*
  A naive implementation of matrix multiplication.

  DO NOT MODIFY THIS FUNCTION, the tests assume it works correctly, which it
  currently does
*/

#define NUMTHREADS 8

typedef struct matrix_args {
	double **A;
	double **B;
	double **C;
	int a_rows;
	int a_cols;
	int b_cols;
	int a_row_start;
	int a_row_end;
} matrix_args;

void *matrix_multiply_sub(void *ptr)
{
	matrix_args *p = (matrix_args *) ptr;
    for (int i = p->a_row_start; i < p->a_row_end; i++) {
        for (int j = 0; j < p->b_cols; j++)
        {
            double sum = 0;
            for (int k = 0; k < p->a_cols; k++)
            {
                sum += p->A[i][k] * p->B[k][j];
            }
            p->C[i][j] = sum;
        }
    }
    return NULL;

}

void matrix_multiply(
    double **C,
    double **A,
    double **B,
    int a_rows,
    int a_cols,
    int b_cols
) {
    for (int i = 0; i < a_rows; i++)
    {
        for (int j = 0; j < b_cols; j++)
        {
            C[i][j] = 0;
            for (int k = 0; k < a_cols; k++)
            {
                C[i][j] += A[i][k] * B[k][j];
            }
        }
    }
}

/**
 * I'm guessing I could do more applying the memory management techniques
 * we learned in module one but am going to spend my time on the other
 * section of the prework
 */
void parallel_matrix_multiply(
	double **c,
	double **a,
	double **b,
	int a_rows,
	int a_cols,
	int b_cols
) {
	pthread_t threads[NUMTHREADS];
	matrix_args params[NUMTHREADS];
	int t;
	for(t=0;t<NUMTHREADS;t++)
	{
        matrix_args args = {
			a,
			b,
			c,
			a_rows,
			a_cols,
			b_cols,
            t * a_rows / NUMTHREADS,
            (t + 1) * a_rows / NUMTHREADS,
	    };
        params[t] = args;
        pthread_create(&threads[t], NULL, matrix_multiply_sub, &params[t]);
	}

    for(t = 0; t < NUMTHREADS; t++)
    {
        pthread_join(threads[t], NULL);
    }
}
