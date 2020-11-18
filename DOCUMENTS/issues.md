issues.md
---------


# Issue List

### Parser 解析 "," 有问题.

### Desc
http://localhost:8080/product?query={list{id},user(id:1){name}}

问题是, 貌似 GraphQL 的 EBNF 并没有规定 Selection 之间的分隔符是什么,


这个可以运行

{
    list{id}, 
    user(id:1){id,name}
}

这个也可以, 并且是等价的.

{
    list{id}
    user(id:1){id name}
}

看来只要是 Ignored 分割就可以 
Ignore ::= UnicodeBOM | WhiteSpace | LineTerminator | Comment | Comma

