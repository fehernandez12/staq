# The StaQ Parser

This is where the magic happens. The parser is responsible for parsing the source code and building the AST for it. The AST is then used by the interpreter to execute the program.

## Introduction

You've probably heard about parsers already. They are programs that take a string of characters as input and produce a tree-like structure as output. This tree is called an abstract syntax tree (AST) and it represents the syntactic structure of the input string.

The _abstract_ in the AST name is based on the fact that certain details that are visible in the source code are omitted for the AST, such as semicolons, newlines, whitespace (depending on the language), comments, braces and parentheses. Those details are not represented in the AST, but are necessary to guide the parser while constructing it.

Something worth noting is that there is not a universal format to represent an AST, so it's really up to the parser, even though all implementations are all quite similar. It basically depends on the programming language that is being parsed.

## The StaQ Parser

### Why not a parser generator?

Honestly, for the StaQ project a parser generator, such as yacc or ANTLR could have been used. After all, parsing is already a well understood problem in computer science and it has basically been solved already.

Nevertheless, writing a parser from scratch is an immensely valuable learning experience. It's a great way to understand how parsers work and how they are implemented. It's also a great way to understand the language that is being parsed, since it's necessary to understand the grammar of the language in order to write a parser for it. It is only after you write a parser that you understand the advantages and drawbacks of parser generators.

### The parsing strategy

For the StaQ Programming language, the parsing strategy of choice is Recursive Descent Parsing. This is a top-down parsing strategy that consists of a set of recursive procedures, one for each non-terminal symbol of the grammar. The parser starts with the start symbol of the grammar and recursively expands it until it reaches the terminal symbols. It also works very similar to the most basic idea of a tree and, after all, the AST is a tree.

We will also be evaluating expressions followign the [Top Down](https://journal.stuffwithstuff.com/2011/03/19/pratt-parsers-expression-parsing-made-easy/) [Operator Precedence](http://crockford.com/javascript/tdop/tdop.html), also called Pratt Parsing. This is a top-down parsing strategy that is based on the idea that every operator has a precedence and a left binding power. The precedence determines the order in which the operators are evaluated and the left binding power determines how tightly the operator binds to its left operand. This is a very powerful strategy that allows us to parse expressions without the need for a grammar.

So, instead of associating parsing functions (like the `parseLetStatement` receiver function) with grammar rules, we associate them with tokens. This is a very powerful strategy that allows us to parse expressions without the need for a grammar.

### Statements

In the StaQ programming language, a statement is a single instruction that the computer can execute. It can be a variable declaration, a function call, a loop, a conditional, etc. The parser is responsible for parsing the statements and building the AST for them.

Some examples of statements in StaQ are:

```
let myArray = [1, 2, 3, 4];
let myMap = {"name": "StaQ", "version": 0.1};
let name = "StaQ";

...

return true;
return 0;
```

### Expressions

Expressions are a bit more complex than statements. They are basically a combination of values, variables, operators and function calls that are evaluated to produce a single value. The parser is responsible for parsing the expressions and building the AST for them.

Some examples of expressions in StaQ are:

```
1 + 2
((5 * 5) + 10)
add(1, 2)
!true
```

Expressions can involve several types of operators, such as prefix operators:

```
!true
-5
```

Infix operators:

```
1 + 2
5 * 5
```

And postfix operators:

```
x++
```

Besides the basic operators, there are also the assignment and comparation operators:

```
x = 5
x == 5
x ?? 5
```

And also call expressions:

```
add(1, 2)
add(add(1, 2), 3)
```

Identifiers are also expressions.

It's worth noting that functions in StaQ are **first-class citizens**, which means that function literals are also expressions. We can use a let statement to bind a function to a name. The function literal is the value that is bound to the name.

```
let add = fn(a, b) {
    return a + b;
};
```

And we can also use function literals in place of identifiers:

```
fn(x, y) { return x + y; }(1, 2)
```

Instead of the ternary operator, StaQ uses the `if` expression, which is more readable:

```
let result = if (x > y) { true } else { false };
```
