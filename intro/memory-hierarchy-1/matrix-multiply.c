/**
 * Matrix multiplication is a very common operation, so we want it to run
 * as fast as possible, even on very large matrices.  There are many ways
 * to speed up matrix multiplication, including parallelism, clever
 * algorithms, and of course cache utilization, which will be the focus of
 * this exercise.
 *
 * In this file, you'll find a naïve implementation of matrix
 * multiplication, and a function optimistically called
 * fast_matrix_multiply for you to implement.
 *
 * To compile and test your code, run:
 *
 * 	cc -Wall matrix-multiply.c benchmark.c && ./a.out 512
 *
 * .  The 512 command line argument represents the size of the matrices to
 * be multiplied, and you may want to change this.  For simplicity, we use
 * square matrices.  The benchmark program will create two random matrices
 * of the given size, multiply them against one another using both the
 * naïve baseline program and your own implementation, verify that the
 * results are the same, and measure the amount of time taken to execute
 * each.
 *
 * If you run the code without changing anything you should be able to
 * verify that both functions run in approximately the same amount of time.
 *
 * As you work towards a solution, you may wish to utilize both
 * first-principles thinking as well as benchmarking and profiling.  Doing
 * this exercise in C should make it easier to use tools like cachegrind as
 * you did previously.
 *
 * As a stretch goal, you may also wish to run the program at various
 * compiler optimization levels, and with different matrix sizes, to see
 * how these factors affect performance.  You may even enjoy plotting the
 * results for various sizes to see what patterns emerge.  Do you expect
 * linear growth or something else?
 */
#include <stdio.h>
#include <stdlib.h>

/*
  A naive implementation of matrix multiplication.

  DO NOT MODIFY THIS FUNCTION, the tests assume it works correctly, which it
  currently does
*/
void matrix_multiply(double **C, double **A, double **B, int a_rows, int a_cols,
                     int b_cols) {
  for (int i = 0; i < a_rows; i++) {
    for (int j = 0; j < b_cols; j++) {
      C[i][j] = 0;
      for (int k = 0; k < a_cols; k++)
        C[i][j] += A[i][k] * B[k][j];
    }
  }
}

/**
 * * First approach is to change the loop ordering so that we access B in
 *   column-major order.  This improves speedup to 2.2-4.3x.
 *
 * * Second approach is found from googling, "loop tiling".  The idea here
 *   seems to be to reduce the ranges being accessed as to have them all in
 *   cache memory at once.
 *
 *   This means we need to pick a tile size.  With 8-byte doubles, tile size
 *   of 2 would let 8 * (2 * 2) = 32 bytes still fit in cache.  tile size of
 *   3 would not fit in a cache.
 *
 *   However I'm seeing better results with tile_size = 8.  Not sure how to
 *   explain this result.
 *
 *   This approach seems to raise average speedup to closer to 3x
 */
#define TILE_SIZE 8
#define min(a, b) (((a) < (b)) ? (a) : (b))
void fast_matrix_multiply(double **C, double **A, double **B, int a_rows,
                          int a_cols, int b_cols) {
    for (int i = 0; i < a_rows; i+=TILE_SIZE) {
        int t_row_stop = min(i + TILE_SIZE, a_rows);

	    for(int dot = 0; dot < a_cols; dot++) {

            for (int j = 0; j < b_cols; j+=TILE_SIZE) {
                int t_col_stop = min(j + TILE_SIZE, b_cols);

                for(int t_row = i; t_row < t_row_stop; t_row++) {
                    for(int t_col = j; t_col < t_col_stop; t_col++) {
                        C[t_row][t_col] += A[t_row][dot] * B[dot][t_col];
                }
            }
        }
		}
	}
}
