# PCap Exercises

## Global `pcap-savefile` header
```shell
xxd net.cap | head
```

```
00000000: d4c3 b2a1 0200 0400 0000 0000 0000 0000  ................
00000010: ea05 0000 0100 0000 4098 d057 0a1f 0300  ........@..W....
00000020: 4e00 0000 4e00 0000 c4e9 8487 6028 a45e  N...N.......`(.^
00000030: 60df 2e1b 0800 4500 0040 d003 0000 4006  `.....E..@....@.
00000040: 2cee c0a8 0065 c01e fc9a e79f 0050 5eab  ,....e.......P^.
00000050: 2265 0000 0000 b002 ffff 5823 0000 0204  "e........X#....
00000060: 05b4 0103 0305 0101 080a 3a4d bdc5 0000  ..........:M....
00000070: 0000 0402 0000 4098 d057 97ab 0400 4c00  ......@..W....L.
00000080: 0000 4c00 0000 a45e 60df 2e1b c4e9 8487  ..L....^`.......
00000090: 6028 0800 4548 003c 0000 4000 2906 d3ad  `(..EH.<..@.)...
```

*What's the magic number?  What does it tell you about the byte ordering in the
pcap-specific aspects of the file?*

Magic number is `0xd4c3b2a1` which means I'm reading from a computer that has
the opposite byte ordering as the machine that wrote the file. (If the machine
had the same byte ordering, the magic number would be `0xa1b2c3d4`)

*What are the major and minor versions?  Don't forget about byte ordering!*

Major version is 2, minor version is 4.

*Are the values zero that ought to be zero in fact zero?*

Yes, I see 8 bytes of zeroes after the major and minor version 2-byte sections.

*What is the snapshot length?*

Snapshot length is the 4 bytes following the zeroes.  Reversing bytes, I see 
`0x000005ea`, which is 1514 bytes.

*What is the link layer header type?*

The 4 bytes of link layer, when reversed, are: `0x0000_0001`. This corresponds
to ethernet (https://www.tcpdump.org/linktypes.html).

## Per-packet Headers

*The bytes immediately following the global header will be the first per-packet
header data.  Parse these values manually as well:*

Without reversing bytes:
```
 4098 d057 0a1f 0300 4e00 0000 4e00 0000
| ts (s)  | ts (ms) | len (tr)| len     |
```

*What is the size of the first packet?*

First packet is `0x0000004e` or 78 bytes

*Was any data truncated?*

No, the bytes representing untruncated length match those representing truncated
length.

