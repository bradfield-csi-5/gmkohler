# Lecture Notes: Binary Representations of Data
## Reviewing floats

Remember the "mantisse" is fractional bits, so representing even 100 as a float won't go in as its int representation

Adding floats up 100M times is going to give you a way less precise number than if you just declared the float as 100M (float arithmetic vs representation).

The peak number when incrementing is 2^24 which is the size of the mantisse

## Review of previous exercise (VM work)

### Languages with VMs vs langs compiled to machine code

VM "simulates" a CPU

When you run Python code, what happens?
* reads source code, "cpython" executable analyses it and turns it into
  bytecode
  * think of bytecode as code that's meant to be executed by a VM
* Python takes the bytecode and will translate to machine code (I think?)

(Do not confuse VirtualBox style VM with Python VM / YARV / JVM)

## Looking through asm code for factorial function

Review this in the recording

Also look at 5pm lecture in case they get more into the encoding content
