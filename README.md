fast-graphql
------------

An extremely fast GraphQL implementation with Go and GOASM.


Steps
-----

- 利用 c2goasm, 验证 C -> ASM -> GOASM + go-sourcecode -> bin 的构建流程 [Check]
- GraphQL 的 BNF 或 EBNF 定义 [Check]
- 编写一个只含有一种 scalar 的 query 的前端 [Check]
- 编写与之配套的后端 [Check]
- 实现单一类型的 Arguments 请求并输出结果 [Check]
- 修正分隔符问题 [Check]
- 修正输入Arguments Feild不存在的错误提示 []
- ObjectField.Type 需要与 ResolveFunction 得到的 type 相匹配 [check] implement -> resolvedDataTypeChecker
- variable 的实现
    - [未定义行为] variable 未传入报错
    - 传入 variable
        - int       [check]
        - string    [check]
        - float     [check]
        - list      []
        - boolean   [check]
        - enum      []
        - object    []
    - 从用户输入获取 variable
    - DecodeVariables 函数使用了 json.Unmarshal, 执行完毕后默认数字类型是float64, 如果需要 int 需要手动转换. [Reference: pkg/encoding/json/#Unmarshal](https://golang.org/pkg/encoding/json/#Unmarshal)
    ```
    To unmarshal JSON into an interface value, Unmarshal stores one of these in the interface value:
    
    bool, for JSON booleans
    float64, for JSON numbers
    string, for JSON strings
    []interface{}, for JSON arrays
    map[string]interface{}, for JSON objects
    nil for JSON null
    ```
- 实现所有类型的 Arguments 请求并输出结果 []
    - variable  [check]
    - int       [check]
    - string    [check]
    - float     [check]
    - list      
    - boolean   [check]
    - enum
    - object

- 修正全部 Ignored Definition
- 完善 errors




Issues 
-----------------

- 前端解析成AST后, AST的构造形式应该是什么样子的? 会对后续查询过程有设计/性能影响么?
- 查询是如何在目标的数据结构上进行查找的?
    - 每个数据结构会被```graphql.NewObject```实例化为专用结构, 并提供```Resolve```方法进行查询. ```Resolve```方法需要自己构造.
- AST如何作用查询的?
- variable 为什么没有纳入 GraphQL Schema? 现有实现(graphql-go, graphql-js) 都是通过参数传递进去的.

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
- [GraphQL Specification](http://spec.graphql.org/)
- [c2goasm 简单教程](./DOCUMENTS/c2goasm-usage.md)
- [GraphQL BNF](./DOCUMENTS/graphql.bnf)
- [GraphQL Specification](https://github.com/graphql/graphql-spec)