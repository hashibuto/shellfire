# shellfire
Buffer overflow swiss army knife

## Get it
```
wget -O /tmp/shellfire https://github.com/hashibuto/shellfire/releases/download/v0.1.13/shellfire && chmod +x /tmp/shellfire && sudo mv -f /tmp/shellfire /usr/local/bin/shellfire
```


## Buffer overflow return address pattern tool

Create a pattern of n bytes

`shellfire pattern create <num_bytes>`

```
> shellfire pattern create 128
0000111122223333444455556666777788889999AAAABBBBCCCCDDDDEEEEFFFFGGGGHHHHIIIIJJJJKKKKLLLLMMMMNNNNOOOOPPPPQQQQRRRRSSSSTTTTUUUUVVVV
```

Create a pattern with fixed bytes

```
> shellfire pattern create -f 128 228
================================================================================================================================0000111122223333444455556666777788889999AAAABBBBCCCCDDDDEEEEFFFFGGGGHHHHIIIIJJJJKKKKLLLLMMMMNNNNOOOO
```

Determine the offset of a sub-pattern (4 bytes)

`shellfire pattern offset [-h] [-f=int] <pattern>`

```
> shellfire pattern offset LMMM
91
```

## Produce a shellcode payload

`shellfire payload [-hba] [-n=int] <offset> <returncode> <shellcode>`

```
> shellfire payload 128 \x33\x23\x23\x55 \x33\x23\x23\x55\x33\x23\x23\x55\x33\x23\x23\x55\x33\x23\x23\x55\x33\x23\x23\x55`
```

Because the output is binary, it may not be visible in the terminal.  
Use the `-h` flag to output as a hex encoded string.  
Use `-n` to specify nopsled length - by default it will account for half the remaining buffer space after the shellcode.  
Use `-b` to specify big-endian byte order (this will preserve the return address order instead of reversing it).  
Use `-a` to append shellcode after the return address instead of inserting it into the buffer


## Produce a stripe pattern for buffer inspection
This generates a visually inspectable sequence of bytes which should be easy to follow in a debugger, giving you a rough indication of position.

`shellfire stripe [-h] [-a=int] [-e=int] <length>`
```
> shellfire stripe 64
DDDDUUUUDDDDUUUUDDDDUUUUDDDDUUUUDDDDUUUUDDDDUUUUDDDDUUUUDDDDUUUUwwww
```
Use the `-h` flag to output as a hex encoded string
Use the `-a` flag to specify a number of alignment bytes to prepend to the buffer

# Evaluate simple arithmatic expressions
`shellfire eval [-d] <expression>`
```
> shellfire eval -d "0x1FAD + 222"
\x0000208b
```
Use the `-d` option to switch output to unpadded decimal if that's useful.  Currently only `+` and `-` operators are supported.  Expression can contain hexadecimal and decimal (mixed) but all hexadecimal must be disambiguated by providing `\x` or `0x` as a prefix.