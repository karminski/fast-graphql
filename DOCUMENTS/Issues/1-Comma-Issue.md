1-Comma-Issue.md
----------------


# Parser 解析 "," 问题.

# Description

GraphQL 的 EBNF 并没有规定 Selection 之间的分隔符是什么,

这个可以运行  

```graphql
{
    list{id}, 
    user(id:1){id,name}
}
```

这个也可以, 并且是等价的.  

```graphql
{
    list{id}
    user(id:1){id name}
}
```

看来只要是 Ignored 分割就可以 
Ignore ::= UnicodeBOM | WhiteSpace | LineTerminator | Comment | Comma


# Solution

修正了 Ignored Lexer 定义, 使其符合 GraphQL Specification.
