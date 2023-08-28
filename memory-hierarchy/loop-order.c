/**
 * Two different ways to loop over an array of arrays.
 * Spotted at: http://stackoverflow.com/questions/9936132/why-does-the-order-of-the-loops-affect-performance-when-iterating-over-a-2d-arra
 *
 * This file contains two options for looping through an array of arrays,
 * one in column orer and the other in row orer.  Prior to this class, you
 * may have thought that both strategies would perform the same—in fact,
 * that's what "Mark" on Stack Overflow thought.
 *
 * Now that you have understanding of memory hierarchy, you should be able
 * to answer Mark's question!
 *
 * The purpose of this exercise is to develop your understanding further,
 * by using profiling tools to probe the program.  To do so, add comments
 * to the program to isolate on option or the other, compile each version,
 * and run them through cachegrind.
 *
 * * Which function takes longer to run, if any?
 *
 *   I would expect option_two to take longer to run because of lesser
 *   spatial locality.
 *
 * * Do they execute the same number of instructions?
 *
 *   Based on the output of objdump I'm seeing the same number of
 *   instructions for each option.
 *
 * * What do you notice about the cache utilization of each program?  How
 *   much better is one than the other?  Does this match your expectations?
 *
 *   option_one has fewer D1 cache misses, less than 1% whereas option_two
 *   has roughly 12% D1 cache misses.  When I decrease the dimension of the
 *   matrix, the option_two cache misses go down which is also expected.
 *
 * * As a stretch goal, try to do a first-principles analysis of the
 *   expected performance of both functions, considering your specific
 *   hardware.  How close was the cachegrind simulation?
 *
 * If you're running an intel processor that was made in the last 10 years,
 * you almost certainly have a sizable, shared L3 cache.  If so, valgrind
 * will ignore your L2 cache and use only your L1 and L3 cache in its
 * simulation.  The cachegrind manual makes a case for why this is a
 * reasonable choice, but if you would like to try simulating your L2 cache
 * as the lowewst level cache, you can do so with the --LL flag.
 *
 * We encourage you to explore this simple example further before we move
 * into more challenging ones.  For instance, have you tried changing the
 * dimensions of the data?  What happens when you change compiler
 * optimization levels?  You might also try running your cachegrind output
 * through cg_annotate or qcachegrind although these tools tend to be more
 * reliable on linux than macOS.
*/
#define N 10000

void option_one() {
  int i, j;
  static int x[N][N];
  for (i = 0; i < N; i++) {
    for (j = 0; j < N; j++) {
      x[i][j] = i + j;
    }
  }
}

void option_two() {
  int i, j;
  static int x[N][N];
  for (i = 0; i < N; i++) {
    for (j = 0; j < N; j++) {
      x[j][i] = i + j;
    }
  }
}

int main() {
  // option_one();
  option_two();
  return 0;
}
