#include <stdio.h>
#include <stdlib.h>

int mult(int n, int m) {
	return n * m;
}

int main() {
	int n = 3;
	int m = 5;
	
	printf("%d * %d = %d\n", n, m, mult(n, m));
	abort();
}
