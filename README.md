# ed

`ed` is a clone of the UNIX command-line tool by the same name `ed` a line
editor that was nortorious for being and most unfriendly editor ever created.

This is a modern take on that editor writtein in the [Go](https://golang.org)
language for portability with all the basic ed commands, a modern readline
line editor with vi bindings and friendly error messages when things go wrong.

Some future ideas may include:

-   Syntax Highlighting
-   Color Prompt
-   A pony?!

## Quick Start

```\#!sh
$ go get github.com/prologic/ed
$ ed
```

For help on how to use `ed` in general please refer to this excellent guide:

-   [Actually using ed](https://sanctum.geek.nz/arabesque/actually-using-ed/)

Or the [GNU ed manual](chrome-extension://klbibkeccnjlkjkiokjodocebajanakg/suspended.html#ttl=GNU%20'ed'%20Manual&pos=7563&uri=https://www.gnu.org/software/ed/manual/ed_manual.html)

**NOTE:** This port/clone does not implement all ed commands, only a subset.

## Example

Here is an example `ed` session:

```\#!sh
$ ed
> a
Start typing a few lines.
Anything you want, really.
Just go right ahead.
When you're done, just type a period by itself.
.
> p
When you're done, just type a period by itself.
> n
4    When you're done, just type a period by itself.
> .n
4    When you're done, just type a period by itself.
> 3
Just go right ahead.
> .n
3    Just go right ahead.
> a
Entering another line.
.
> n
4    Entering another line.
> i
A line before the current line.
.
> n
4    A line before the current line.
> c
I decided I like this line better.
.
> n
4    I decided I like this line better.
>
```

## Commands

This implementation supports the following commands:
(_not all commands from the original `ed` are supported nor the GNU `ed`_)

| Command   | Description  | Notes                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         |
| --------- | ------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| ``        | newline      | Null command. An address alone prints the addressed line. A <newline> alone is equivalent to '+1p'. The current address is set to the address of the printed line. |
| `=`       | print index  | Prints the line number of the addressed line. The current address is unchanged. |
| `a`       | append text  | Appends text to the buffer after the addressed line. Text is entered in input mode. The current address is set to the address of the last line entered or, if there were none, to the addressed line.                                                                                                                                                                                                                                                                                                         |
| `c`       | change lines | Changes lines in the buffer. The addressed lines are deleted from the buffer, and text is inserted in their place. Text is entered in input mode. The current address is set to the address of the last line entered or, if there were none, to the new address of the line after the last line deleted; if the lines deleted were originally at the end of the buffer, the current address is set to the address of the new last line; if no lines remain in the buffer, the current address is set to zero. |
| `d`       | delete lines | Deletes the addressed lines from the buffer. The current address is set to the new address of the line after the last line deleted; if the lines deleted were originally at the end of the buffer, the current address is set to the address of the new last line; if no lines remain in the buffer, the current address is set to zero.                                                                                                                                                                      |
| `f file`  | set filename | Sets the default filename to file. If file is not specified, then the default unescaped filename is printed.                                                                                                                                                                                                                                                                                                                                                                                                  |
| `i`       | insert text  | Inserts text in the buffer before the addressed line. The address '0' (zero) is valid for this command; it places the entered text at the beginning of the buffer. Text is entered in input mode. The current address is set to the address of the last line entered or, if there were none, to the addressed line.                                                                                                                                                                                           |
| `j` .     | join lines . | Joins the addressed lines, replacing them by a single line containing their joined text. If only one address is given, this command does nothing. If lines are joined, the current address is set to the address of the joined line. Else, the current address is unchanged.                                                                                                                                                                                                                                  |
| `n`       | numbered     | Number command. Prints the addressed lines, preceding each line by its line number and a <tab>. The current address is set to the address of the last line printed.                                                                                                                                                                                                                                                                                                                                           |
| `p`       | print lines  | Prints the addressed lines. The current address is set to the address of the last line printed.                                                                                                                                                                                                                                                                                                                                                                                                               |
| `q`       | quit         | Quits ed. A warning is printed if any changes have been made in the buffer since the last 'w' command that wrote the entire buffer to a file.                                                                                                                                                                                                                                                                                                                                                                 |
| `r file`  | read file    | Reads file and appends it after the addressed line. If file is not specified, then the default filename is used. If there is no default filename prior to the command, then the default filename is set to file. Otherwise, the default filename is unchanged. The address '0' (zero) is valid for this command; it reads the file at the beginning of the buffer. The current address is set to the address of the last line read or, if there were none, to the addressed line.                             |
| `w file`  | write file   | Writes the addressed lines to file. Any previous contents of file is lost without warning. If there is no default filename, then the default filename is set to file, otherwise it is unchanged. If no filename is specified, then the default filename is used. The current address is unchanged.                                                                                                                                                                                                            |
| `wq file` | write & quit | Writes the addressed lines to file, and then executes a 'q' command.                                                                                                                                                                                                                                                                                                                                                                                                                                          |
| `x`       | put text     | Copies (puts) the contents of the cut buffer to after the addressed line. The current address is set to the address of the last line copied. |
| `y`       | yank text    | Copies (yanks) the addressed lines to the cut buffer. The cut buffer is overwritten by subsequent 'c', 'd', 'j', 's', or 'y' commands. The current address is unchanged. |

## License

    `ed` is licensed under the terms of the [MIT License](https://opensource.org/licenses/MIT)
