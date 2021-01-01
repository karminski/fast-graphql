fast-graphql
------------

([English](./README.md) | [中文](./README-zh-CN.md))


# 说明

该项目是一个 Go 语言编写的实验性 GraphQL 实现. 目的是将 GraphQL 的 Parse 和 Resolve 性能提升到一个新的水平.  
目前, 该项目在尝试使用 JIT Parser, ASM 优化 Parser, SIMD 优化方法来提升性能.  

注意: 本项目仍然在构建中, 并不可以投入生产使用.

# Steps & Objects

- 完成一个只有 query operation 的最小 Demo. (✔️)
- 完成 ASM 优化 Parser 的 Demo. (✔️)
- 完成 GraphQL 的 EBNF 定义, 方便构建编译器前端. (✔️)
- 完成 GraphQL Lexer & Parser. (✔️)
- 完成 OperationDefinition 的相关功能. ()
- 完成 JIT 前端原型. ()
- 完成字符串序列化优化. ()
- 完成 FragmentDefinition 的相关功能. ()
- 完成所有 GraphQL 解析功能. ()
- 完成测试. ()
- 完成 GC & memory 优化. ()
- 完成 ASM & SIMD 优化. ()

# Examples

- [HTTP GET Method Demo](./src/cmd/http-get-example/main.go)
- [HTTP POST Method Demo](./src/cmd/http-post-example/main.go)
- [Lexer Demo](./src/cmd/fast-graphql-frontend/main.go)


# Documents

- [现存问题](./DOCUMENTS/issues.md)
- [事项列表](./DOCUMENTS/todo-list.md)
- [想法和优化方案](./DOCUMENTS/ideas.md)

# Dependency

本项目的后端逻辑是从 [graphql-go](https://github.com/graphql-go/graphql) 移植而来, Parser & Lexer 部分则受到了 Lua 的启发构建而成.

# Contributors

- [karminski](https://github.com/karminski)

# License

- [MIT](./LICENSE)

# Reference

- [GraphQL Grammar EBNF Definition](https://github.com/karminski/graphql-grammar-ebnf-definition)
- [c2goasm usage](./DOCUMENTS/c2goasm-usage.md)
- [GraphQL Specification](http://spec.graphql.org/)
- [GraphQL Specification on Github](https://github.com/graphql/graphql-spec)