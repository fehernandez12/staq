# The StaQ Parser

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

### Statements

In the StaQ programming language, a statement is a single instruction that the computer can execute. It can be a variable declaration, a function call, a loop, a conditional, etc. The parser is responsible for parsing the statements and building the AST for them.

### Expressions

Expressions are a bit more complex than statements. They are basically a combination of values, variables, operators and function calls that are evaluated to produce a single value. The parser is responsible for parsing the expressions and building the AST for them.
