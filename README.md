fast-graphql
------------

An extremely fast GraphQL implementation with Go and GOASM.


Steps
-----

- 利用 c2goasm, 验证 C -> ASM -> GOASM + go-sourcecode -> bin 的构建流程. [已完成]
- 调研是否存在 GraphQL 的 C 实现, 如果没有, 则需要自己实现.
- 实现 C 版本的 GraphQL 最小单元 Parser Demo.
- 验证 Parser Demo 生成 GOASM 嵌入 Go 并调用测试.
- 实现 C 版本的 GraphQL.
- 编写 GraphQL 测试用例.
- 接入 C 版本的 GraphQL 到 fast-graphql.
- 测试 fast-graphql.


Reference
---------

- [c2goasm 简单教程](./DOCUMENTS/c2goasm-usage.md)