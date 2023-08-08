# Introduction Lecture


## Stretch goal to understand binary add/mult

For anyone who has this stretch goal, here’s a concrete thing I might
recommend:

1. Write a C program that takes two arbitrarily long integers (ok you can
limit it to like 1000 digits) in two lines of input, then adds them and
prints the result in a line of output
2. Once you get addition working, try multiplication by any means necessary
(i.e. the “grade school way”)
3. Once you have multiplication working, try the
[Karatsuba algorithm](https://en.wikipedia.org/wiki/Karatsuba_algorithm).

## 32 vs 64 bit

64 bit provides more addressable memory

Think of word-size and pointer-size as the same and that they're the same 
as the "size" of the architecture.

## why doesn't catting executables work?

Interpreted as ASCII or UTF-8 (depends on terminal)—the encoding is not
designed to be human readable.

We can use something like `hexdump` instead.

```sh
hexdump -C a.out | less
```

or `objdump`:

```sh
objdump -d a.out | less
```

## Why registers?

Quicker data access (and they're implemented with D(ynamic)-RAM instead of
S(tatic)-RAM for faster access, flip-flop vs capacitor).

Intel extended their registers from 16 to 32-bit so they started with %ax
(16-bit register) then extended it to $eax.

You can refer to lower-order bits of $eax via $ax even if you've used the
64 bits for storage

## Memory addresses

Virtual addresses are used to give each program in an operating system their
own address space

When you load the file to be executed, the objects in the executable are
loaded into RAM according to instructions in the file.

## debugger

lldb lets you set breakpoints in the compiled C file and read registers
