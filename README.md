# shellfire
Buffer overflow analysis swiss army knife

## Get it
```
wget -O /tmp/shellfire https://github.com/hashibuto/shellfire/releases/download/v0.1.0/shellfire && chmod +x /tmp/shellfire && sudo mv -f /tmp/shellfire /usr/local/bin/shellfire
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