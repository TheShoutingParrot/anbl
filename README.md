# ANBL - A New Basic (programming) Language

A New Basic (programming) Language is a very simplistic programming language intended to be used for learning how basic programming works and/or learning how a basic interpreter works.

This ANBL interpreter is the first one and is written in Go.

**MORE DOCUMENTATION COMING SOON**

## Examples

### Example usage without a .anbl file

Building and running a simple program by writing it in standard input:

```
$ go build
$ ./interpreter
>>> 1 RESERVE VAR AS NUMBER
>>> 2 VAR IS 9
>>> 3 PRINTNUM VAR
>>> 5 IF NOT EQUALS VAR 0 JUMP 3
>>> 4 DECREMENT VAR
>>> 6 END
>>> RUN
987654321>>>
```

By typing `EXIT` you exit the program.

`RUNANDEXIT` will first run the inputed program and afterwards exit.
