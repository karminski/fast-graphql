fast-graphql
------------

An extremely fast GraphQL implementation with Go and GOASM.


Steps
-----

- 利用 c2goasm, 验证 C -> ASM -> GOASM + go-sourcecode -> bin 的构建流程. [已完成]
- GraphQL 的 BNF 或 EBNF 定义 [已完成]
- 编写一个只含有一种 scalar 的 query 的前端
- 编写与之配套的后端



- 调研是否存在 GraphQL 的 C 实现, 如果没有, 则需要自己实现.
- 实现 C 版本的 GraphQL 最小单元 Parser Demo.
- 验证 Parser Demo 生成 GOASM 嵌入 Go 并调用测试.
- 实现 C 版本的 GraphQL.
- 编写 GraphQL 测试用例.
- 接入 C 版本的 GraphQL 到 fast-graphql.
- 测试 fast-graphql.


Issues 
-----------------

- 前端解析成AST后, AST的构造形式应该是什么样子的? 会对后续查询过程有设计/性能影响么?
- 查询是如何在目标的数据结构上进行查找的?
    - 每个数据结构会被```graphql.NewObject```实例化为专用结构, 并提供```Resolve```方法进行查询. ```Resolve```方法需要自己构造.
- AST如何作用查询的?

Thoughts
--------

- 查询过程是否可以利用现有数据库技术进行优化?
- 查询过程是否可以进行JIT优化?
- 查询过程是否可以进行缓存优化?
- FlameGraph 中的GC热点问题是否有优化空间?
- 最后的序列化过程是否可以整体序列化以提升性能?


Improvements
------------
- 请求缓存问题
- GC问题
- 返回数据拼接效率问题
- 并行优化



Reference
---------

- [c2goasm 简单教程](./DOCUMENTS/c2goasm-usage.md)
- [GraphQL BNF](./DOCUMENTS/graphql.bnf)
- [GraphQL Specification](https://github.com/graphql/graphql-spec)