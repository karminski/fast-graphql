2-Int-Type-Assert-in-Query-Variables-Issue.md
---------------------------------------------

# Query Variables 的 Int 类型解析问题

# Description
在客户端传入 Query Variables 的时候, 需要进行 json.Unmarshal.    
这时 json 数字类型会默认转换为 float64 类型, 是需要根据用户定义的类型Name进行类型转换还是保留当前类型?   
因为看到了有些实现例如 (graphql-go) 并没有做转换, 而是推迟到了用户构建服务端 Schema 的时候让用户自己断言类型(float64)并手动转换.

# Solution

<del> 添加了 correctJsonUnmarshalIntValue() 方法用来转换. <del>

使用了 RequestParser, 在反序列化时直接解析为 Int 类型. 详见: [Request-Parser](../Request-Parser-CN.md)