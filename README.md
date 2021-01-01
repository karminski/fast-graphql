fast-graphql
------------

([English](./README.md) | [中文](./README-zh-CN.md))


# Description

An experimental GraphQL implementation with Go. This repo focuses on improve GraphQL Parse and Resolve speed to a new level.  
For now, we will try JIT Parser, Raw-ASM Parser, and SIMD methods to optimize.  

WARNING: this repo is under construction, not production-ready.

# Steps & Objects

- Finish a minimal example include query operation. (✔️)
- Finish simple ASM optimized demo. (✔️)
- Finish grammar EBNF definition. (✔️)
- Finish full GraphQL lexer & parser. (✔️)
- Finish GraphQL backend with OperationDefinition feature. ()
- Finish simple JIT prototype. ()
- Finish stringify optimize. ()
- Finish GraphQL backend with Directive feature. ()
- Finish GraphQL backend with full FragmentDefinition feature. ()
- Finish GraphQL backend with full Definition feature. ()
- Finish test case. ()
- Finish GC & memory Optimize. ()
- Finish ASM & SIMD Optimize. ()

# Examples

- [HTTP GET Method Demo](./src/cmd/http-get-example/main.go)
- [HTTP POST Method Demo](./src/cmd/http-post-example/main.go)
- [Lexer Demo](./src/cmd/fast-graphql-frontend/main.go)


# Documents

- [Issues](./DOCUMENTS/issues.md)
- [Todo List](./DOCUMENTS/todo-list.md)
- [Ideas](./DOCUMENTS/ideas.md)

# Dependency

The basic backend logic of this repo is port from [graphql-go](https://github.com/graphql-go/graphql), and the lexer & parser are inspired by Lua.

# Contributors

- [karminski](https://github.com/karminski)

# License

- [MIT](./LICENSE)

# Reference

- [GraphQL Grammar EBNF Definition](https://github.com/karminski/graphql-grammar-ebnf-definition)
- [c2goasm usage](./DOCUMENTS/c2goasm-usage.md)
- [GraphQL Specification](http://spec.graphql.org/)
- [GraphQL Specification on Github](https://github.com/graphql/graphql-spec)