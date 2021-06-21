Request-Parser-CN.md
------------------

# Desc

GraphQL 请求是包裹在结构化的 json 数据中的 (详见: [QuerySample](#Query-Sample)). 
为了提升缓存命中几率, 我们实施了参数代换 (详见: [Arguments-Substitution](./Arguments-Substitution-CN.md)).
由于参数代换后, 所有的变量全都集中在 Query Varibales 中, 因此 Query Varibales 的解析性能至关重要.  
然而 Go 内置的 encoding/json 库并不是为了这种 variables dictionary 场景而准备的, 因此性能并不理想. 
使用 encoding/json 解析这段数据结构大概消耗了整体请求过程耗时的 10%, 作为对比, 执行缓存流程的 fast-graphql 耗时占比为 30% 左右, 其余是 net/http 库和 runtime, I/O 耗时).
因此手动实现一个固定结构的解析器, 会明显提升性能.





# Request Sample

```json
{
    "query": "query UserInfo{\n User(Id:$Id){\n    Id \n    Name } }",
    "variables":
    {
        "Id": 1
    },
    "operationName": "UserInfo"
}
```

# Request EBNF 定义

```ebnf
# Request
Request        ::= Ignored "{" Ignored RequestField Ignored "}" Ignored

# RequestField

RequestField   ::= QueryString | QueryVariables | OperationName

# QueryString
QueryString    ::= '"query"' Ignored ":" Ignored '"' StringValue '"'

# QueryVariables
QueryVariables ::= Ignored "{" Ignored QueryVariable+ Ignored "}" Ignored
QueryVariable  ::= Ignored VariableName Ignored ":" Ignored VariableValue Ignored
VariableName   ::= StringValue
VariableValue  ::= IntValue | FloatValue | StringValue | BooleanValue | NullValue 

# OperationName
OperationName  ::= '"operationName"' Ignored ":" Ignored '"' StringValue '"'
```