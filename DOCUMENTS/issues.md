issues.md
---------


# Opend Issues

### #2, 前端解析成AST后, AST的构造形式应该是什么样子的? 会对后续查询过程有设计/性能影响么?

### #3, 查询是如何在目标的数据结构上进行查找的?

#### Description
每个数据结构会被```graphql.NewObject```实例化为专用结构, 并提供```Resolve```方法进行查询. ```Resolve```方法需要自己构造.

### #4, AST如何作用查询的?

### #5, variable 为什么没有纳入 GraphQL Schema? 现有实现(graphql-go, graphql-js) 都是通过参数传递进去的.

### #7, VariableDefinition 的 Type 需要进行检查.

### #8, 请求存在多个 Selection 的情况下, 如果有一个 Selection执行失败, 错误该怎样处理?

# Closed Issues

### #1, Parser 解析 "," 有问题.

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


### #6, Query Variables 的 Int 类型解析问题

#### Description
在客户端传入 Query Variables 的时候, 需要进行 json.Unmarshal, 这时 json 数字类型会默认转换为 float64 类型, 是需要根据用户定义的类型Name进行类型转换还是保留当前类型? 因为看到了有些实现例如 (graphql-go) 并没有做转换, 而是推迟到了用户构建服务端 Schema 的时候让用户自己断言类型(float64)并手动转换.

#### Solution

总之我添加了 correctJsonUnmarshalIntValue() 方法用来转换.