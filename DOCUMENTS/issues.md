issues.md
---------


# Opend Issues

### 前端解析成AST后, AST的构造形式应该是什么样子的? 会对后续查询过程有设计/性能影响么?

### 查询是如何在目标的数据结构上进行查找的?

#### Description
每个数据结构会被```graphql.NewObject```实例化为专用结构, 并提供```Resolve```方法进行查询. ```Resolve```方法需要自己构造.

### AST如何作用查询的?

### variable 为什么没有纳入 GraphQL Schema? 现有实现(graphql-go, graphql-js) 都是通过参数传递进去的.


# Closed Issues

### Parser 解析 "," 有问题.

#### Description

http://localhost:8080/product?query={list{id},user(id:1){name}}

问题是, 貌似 GraphQL 的 EBNF 并没有规定 Selection 之间的分隔符是什么,

这个可以运行
```
{
    list{id}, 
    user(id:1){id,name}
}
```
这个也可以, 并且是等价的.
```
{
    list{id}
    user(id:1){id name}
}
```
看来只要是 Ignored 分割就可以 
Ignore ::= UnicodeBOM | WhiteSpace | LineTerminator | Comment | Comma

#### Solution

修正了 Ignored Lexer 定义.



