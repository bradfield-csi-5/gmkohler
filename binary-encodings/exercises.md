# Binary Represenation of Data

## 1 Hexadecimal

### 1.1 Simple conversion

*What are the numbers 9, 136, and 247 in hexadecimal?*

Interpreted as decimal numbers:

* 9 would be **0x9** because it fits into both bases
* 136 fits 16 into it 8 times (128) then has 8 remainder, so **0x88**
* 247 fits 16 into it 15 (F) times (240), then has 7 remainer, so **0xF7**

### 1.2 CSS colors

*In CSS, two ways to specify colors are hexadecimal and Rgb.  For instance,
pure red would be `0xff0000` or `rgb(255, 0, 0);`.  How many colors can be
represented in each form?*


The same number should be able to be represented in each encoding, and
should be:

```
256^2 - 1 == 16^6 - 1 == (2^4)^6 - 1 == 2^(4*6) - 1 == 2^24 - 1 == 16777215
```

Because that's the number of hex digits we have allotted for representing
the colors.

### 1.3 `hellohex`

*`hellohex` is 17 bytes in size.*
*If you were to view a hexidecimal representation of the file, how many
hexadecimal characters would you expect?*

2 hex characters cover 1 byte, so for 17 bytes, I'd expect
2 * 17 == **34 characters** to represent 17 bytes.  However, the number
could be lower if the highest-order bits are insignificant

*Convert the first 5 bytes in the file by hand to binary.  Write these
down as you'll use them again in a later exercise*

5 bytes is 10 hex digits (`0x68656c6c6f`), we can convert each digit to four bits:

```
0x|6    8    6    5    6    c    6    c    6    f
0b 0110 1000 0110 0101 0110 1100 0110 1100 0110 1111
```

(this matches the output of `xxd -b -l 5 hellohex`)

## 2 Integers

### 2.1 Basic conversion

*Convert the following decimal numbers to binary*

4 is 2^2 so fill in the 2nd bit (index 0)
65 is one more than 2^6 so fill in the 6th bit and the 0th bit
105 is 40 more than 65 so fill in the 5th (32) and the 3rd (8) digits to add 40
255 is obviously all of the bits because it's the largest number representable in uint8

```
  4 | 0000 0100
 65 | 0100 0001
105 | 0110 1001
255 | 1111 1111
```

*Convert the following binary representations of unsigned integers to decimal*

```
0000 0010 | 2
0000 0011 | 3
0110 1100 | 108 (64 + 32 + 8 + 4)
0101 0101 | 85  (64 + 16 + 4 + 1)
```

### 2.2 Unsigned binary addition
*Binary addition works the same as decimal addition, except that we carry
from one place to another after exceeding not a 9, but a 1.*

*Add these two binary numbers and determine the result by doing the addition
"in binary".  Convert each number and check that the result matches your
expectation*

```
---1 1111 111- # carries
0000 1111 1111
0000 0000 1101 +
===============
0001 0000 1100
```

Reading the result by adding powers of two, I see is 256 + 8 + 4 == 268.
Reading the two operands as decimal and adding, I see 255 + 13 == 268.
So they agree.

*If my registers are only 8 bits wide, what is the value returned from that
binary addition?  What is this phenomenon called?*

If the values are only 8 bits wide we'd lose the top value and only have 12
as an answer.  Integer overflow is the name of the phenomenon.

### 2.3 Two's complement conversion

*Given the following decimal values, determine their 8-bit two's complement representations*

Two's complement of 8-bit integers gives the 7th bit a weight of -128 (-1 * 2^7) and the rest positive
weights of 2^i

```
 127 | 0111 1111
-128 | 1000 0000
  -1 | 1111 1111
   1 | 0000 0001
 -14 | 1111 0010 # 13 less than -1 (all 1's), so swap the 8,4,1 bits to add to 13
```
#

### 2.4 Addition of two's complement signed integers

*What is the sum of the following two signed integers?  Does this match
your expectations?*

```
0111 1111
1000 0000 +
===========
1111 1111
```

The sum is -1, which is expected because the first number is -128 and the
second number is 127, so the overflow cases don't come into play.

If we interpret twos-complement addition as converting to unsigned, adding,
then converting back to two's complement (described this way in §2.3.2),
we'd be adding up to 255 and converting back to twos-complement which is
again -1.

*How do you negate a number in two's complement?  How can we compute
subtraction of two's complement numbers?*

§2.3.3 of *CS:APP* defines negation as follows:

* if x is the minimum number, -x == x i.e. the negation is the identity
* else, -x = -x as you would normally think.

Practically speaking you can invert each bit and add 1 to the result.
For the minimum number (e.g. `0b1000_000`), the flip would have all but
the most significant bit filled (`0b0111_0111`), and adding 1 to that
carries all the ones over (`0b1000_0000`), which fits the definition
above.

With this in mind, subtraction of b from a could be defined as adding
the negated value.

*What is the value of the most significant bit in the 8-bit two's
complement? What about 32-bit two's complement?*

§2.2.3 of *CS:APP* provides says that for a byte of size _w_, the most
significant bit of a twos-complement is `-1 * 2^(w-1)`.

For an 8-bit number, this would be `-1 * 2^7` or `-128`.
For a 32-bit number, this would be `-1 * 2^31` or `-2_147_483_648`.

### 2.5 Advanced: Integer overflow detection

After reading the answer I'm uncertain how you'd come up with the values
for `carry_in`/`carry_out` to evaluate `carry_in ^ carry_out`.

## 3 Byte ordering

### 3.1 It's over 9000!

`xxd` tells me this is `0x2329` or `0b0010_0011_0010_1001` which is 9001.
Treating this big-endian gives us `9001` whereas if we were to look at
the  first byte as the least significant digits we'd get `10531`.

So I conclude it's big-endian.

### 3.2 TCP

```
➜ xxd tcpheader
00000000: af00 bc06 441e 7368 eff2 a002 81ff 5600  ....D.sh......V.
```

I am taking this part to be the actual values:

```
af00 bc06 441e 7368 eff2 a002 81ff 5600
```

Source Port (2 bytes, 0 offset): `0xAF00` (`44_800`)
Destination Port (2 bytes, 2 offset): `0xBC06` (`48134`)
Sequence Number (4 bytes, 4 offset): `0x441E7368` (`1142846312`)
Acknowledgement Number (4 bytes, 4 offset): `0xEFF2A002` (`4025655298`)


### 3.3 Bonus: Byte ordering and integer encoding in bitmaps

will come back to the bonus after finishing required work

## 4 IEEE Floating Point


### 4.1 Deconstruction
*Identify the 3 components of this 32-bit IEEE Floating Point Number and their values*

```
0 10000100 010 1010 0000 0000 0000 0000
- -------- ----------------------------
s     e               f
```

For 32-bit IEEE we have 1 sign bit, k=8 exponent bits, and s=23 fractional
bits as shown above.

Given k = 8, our "Bias" is `2^(k-1) - 1 == 2^7 - 1 == 127`

Here, we're "normalized" beacuse e is neither 0 nor all ones.
Thus:
`s == 0` so positive number

`e == 2^2 + 2^7 == 4 + 128 == 132`

`E == e - Bias == 132 - 127 == 5`

`f == 1/4 + 1/16 + 1/64 == 21/64`

`M == 1+f == 85/64`

So `2^E * M == 2^5 * 85/64 == 32 * 85/64 == 42.5`

*For the largest fixed exponent, `1111_1110 == 254 - 127 = 127`, what is
the smallest (magnitude) incremental change that can be made to a number?*

Least significant change will always be flipping the least significant
fractional bit (1/2^23).  In the context of the largest fixed exponent,
this is:

2^127 * 2^-23 == 2^104

*For the smallest (most negative) fixed exponent, what is the smallest
(magnitude) incremental change that can be made to a number?*

The smallest exponent is going to be `0000_0000` which by bias is -126

2^-126 * 2^-23 == 2^-149

*What does this imply about the precision of IEEE Floating Point values?*

They are more precise at smaller exponents (close to zero) than larger exponents

### 4.2 Advanced: Float casting

Going to return to this as time allows

## 5 Character encodings

*Is there additional space needed for encoding ASCII as UTF-8*?

No additional space is needed to encode ASCII (by design).  The
largest ASCII character is `7f` ie `0b0111_1111`, and 1-byte
unicode characters start with `0`, so all of ASCII fits into
the one-byte UTF-8 character set.

From my reading of UTF-32, there's more wasted space for the vast
majority of characters.

### 5.1 Snowman

*How large is the Snowman emoji?*

Code point `U+2603` need 3 bytes according to the encoding table on
Wikipedia.

Breaking this down, `0x2603` is `0b10_0110_0000_0011`, which means
we need 14 bits to encode it.  In theory that could be 2 bytes but given
the need for continuation bits we have to fit it into the 16 available bits
of the three-byte encoding (I presume by padding the front):

```
1110xxxx_10xxxxxx_10xxxxxx
```

Calling `wc` on the `snowman` file provided in the materials:

```
➜ wc snowman
       0       1       3 snowman
```

, I see 3 bytes.  Calling `xxd` on it:

```
➜  xxd -b snowman
00000000: 11100010 10011000 10000011
#         1110xxss 10ssssss 10ssssss
```

, the UTF-8–encoded snowman bits (denoted as `s`) are packed as I suspected.

### 5.2 Hello again hellohex

*The 5 first bytes of `hellohex`—which you previously converted from hex
to binary—are actually characters.  What are they?  It is cheating to use
`xxd` or `hexdump` at this point just to read the ASCII interpretation.*

```
0x68656c6c6f
```

Looking at `man ascii` to break down `68 65 6c 6c 6f` I see:

```
hello
```

We could break it down by subtracting from 60 to see the ordinal in the
alphabet:

```
8 5 c c f (hex)
8 5 12 12 15 (dec)
h e l l o
```

*What character encoding is used in this file?  If we do look at the ASCII
interpretation column in xxd or hexdump, we see some dots signifying that
these are not printable ASCII characters.  Could they still be characters?*

Given then `0x9880` at the end of the file, we've exceeded ASCII range (max 0x72).
I'm going to guess UTF-8.  Of course the extra characters can be characters, and
given ASCII encodings are also UTF encodings, I would have to guess it's UTF
encoded.

*What if we told you there was a multi-byte Unicode character at the end of
the file.  Could you figure out what it is?  What if we told you it's
encoded in UTF-8?*

Expanding the terminal non-ASCII characters to binary via `xxd -b`, I see 5 bytes:

```
11110000 10011111 10011000 10000000 00001010
```

, the first of which matches the prefix for a 4-byte UTF-8 character, and the
fifth of which matches the prefix for a 1-byte UTF-8 character.  Unpacking
those bytes, I see:

```
11110000 10011111 10011000 10000000 00001010
xxxxx000 xx011111 xx011000 xx000000 x0001010
```

The significant bits not used for UTF encoding:

```
0b 0001 1111 0110 0000 0000
0x 1F        60        0
```

and

```
1010
A
```

`U+1F600` corresponds to a grinning emoji (😀), and `U+000A` is a line feed.

### 5.3 Bonus: Ding ding ding!

Going to return to this as time allows
