# shellfire
Buffer overflow analysis swiss army knife

## Get it
```
wget -O /tmp/shellfire https://github.com/hashibuto/shellfire/releases/download/v0.1.1/shellfire && chmod +x /tmp/shellfire && sudo mv -f /tmp/shellfire /usr/local/bin/shellfire
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

`shellfire pattern offset [-f=int] <pattern>`

```
> shellfire pattern offset LMMM
91
```

Produce a shellcode payload

`shellfire payload [-h] [-n] <offset> <returncode> <shellcode>`

```
> shellfire payload 128 \x33\x23\x23\x55 \x33\x23\x23\x55\x33\x23\x23\x55\x33\x23\x23\x55\x33\x23\x23\x55\x33\x23\x23\x55`
```

Because the output is binary, it may not be visible in the terminal.  Use the `-h` flag to output as a hex encoded string.  Use `-n` to specify nopsled length - by default it will account for half the remaining buffer space after the shellcode.