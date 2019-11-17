# ed

`ed` is a clone of the UNIX command-line tool by the same name `ed` a line
editor that was nortorious for being and most unfriendly editor ever created.

This is a modern take on that editor writtein in the [Go](https://golang.org)
language for portability with all the basic ed commands, a modern readline
line editor with vi bindings and friendly error messages when things go wrong.

Some future ideas may include:

- Syntax Highlighting
- Color Prompt
- A pony?!

## Quick Start

```#!sh
$ go get github.com/prologic/ed
$ ed
```

For help on how to use `ed` in general please refer to this excellent guide:

- [Actually using ed](https://sanctum.geek.nz/arabesque/actually-using-ed/)

Or the [GNU ed manual](chrome-extension://klbibkeccnjlkjkiokjodocebajanakg/suspended.html#ttl=GNU%20'ed'%20Manual&pos=7563&uri=https://www.gnu.org/software/ed/manual/ed_manual.html)

**NOTE:** This port/clone does not implement all ed commands, only a subset.

## License

`ed` is licensed under the terms of the [MIT License](https://opensource.org/licenses/MIT)
