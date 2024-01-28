# The StaQ Programming Language

StaQ is an interpreted programming language that aims to be simple, and easy to use. It has a C-like syntax, support for first-class and higher-order functions and general features. Its interpreter is built using the [Go Programming Language](https://golang.org/).

## Syntax

Let's see how StaQ looks like. The following is a simple StaQ program that prints the first 10 numbers of the Fibonacci sequence:

### Value bindings

This is how to bind a value to a name in StaQ:

```
let age = 1;
let name = "StaQ";
let result = 10 * (20 / 2);
let someHex = 0xFF; // Numbers can also be written in hexadecimal notation
let someOct = 0o77; // Or in octal notation
```

Besides primitives such as numbers, booleans and strings, the StaQ interpreter also supports arrays and maps:

### Arrays and maps

```
let myArray = [1, 2, 3, 4];
let myMap = {"name": "StaQ", "version": 0.1};
```

Accessing elements in arrays and maps is done using the `[]` operator:

```
let myArray = [1, 2, 3, 4];
let myMap = {"name": "StaQ", "version": 0.1};
myArray[0]; // 1
myMap["name"]; // "StaQ"
```

### Functions

The assignment statements can also be used to bind functions to names:

```
let add = fn(a, b) {
    return a + b;
};
```

But StaQ not only supports `return` statements. Implicit return values are also supported:

```
let add = fn(a, b) {
    a + b;
};
```

So, a more complex function, such as the `fibonacci` function, can be written as follows:

```
let fibonacci = fn(x) {
    if (x == 0) {
        0;
    } else {
        if (x == 1) {
            1;
        } else {
            fibonacci(x - 1) + fibonacci(x - 2);
        }
    }
};
```

**Note the recursive call to `fibonacci` itself!.**

StaQ also supports higher-order functions. For example, the `map` function can be implemented as follows:

```
let twice = fn(f, x) {
    f(f(x));
};

let multiplyByTwo = fn(x) {
    x * 2;
};

twice(multiplyByTwo, 10); // 40
```

Here, `twice` takes a function `f` and a value `x`, and applies `f` to `x` twice. The `multiplyByTwo` function is then passed to `twice` as the first argument, and `10` as the second argument.

Functions in StaQ are just values, just like numbers and strings. That makes them first-class functions.
