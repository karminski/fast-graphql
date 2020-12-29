todo-list.md
------------

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
- 完善内部函数错误处理



