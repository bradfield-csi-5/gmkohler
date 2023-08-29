As a final optimization exercise, you will improve the performance of a simple,
“typical” Go program that was written without consideration of CPU caches.
Using what you know about CPU caches, you should be able to substantially
improve its performance, without using any fancy algorithms or low level
optimization techniques.

From the metrics directory of the included files, you can run:

```bash
go test -bench=.
goos: darwin
goarch: amd64
BenchmarkMetrics/Average_age-12                  606       1879392 ns/op
BenchmarkMetrics/Average_payment-12               56      28236512 ns/op
BenchmarkMetrics/Payment_stddev-12                22      54153588 ns/op
PASS
```
for a basic benchmark. Your task is of course so reduce the time taken per
execution of the functions under test (ns/op).

You are welcome to change any aspects of the organization of the data in memory,
including types, so long as the tests still pass. You are also welcome to change
the code used to load the data from disk, and in fact may need to do so to reflect
changes you make elsewhere. However, please do not precompute the answers
directly in the LoadData function :)

You are also welcome to make minor changes to the benchmark code, so long as it
doesn’t affect the substance of the test. You should not need to alter the code
used to generate the test data, but this is provided for your information, and
may be useful to know the ranges of the provided data.

As a stretch goal, simply try to make the program as fast as possible, using any
techniques at your disposal!
