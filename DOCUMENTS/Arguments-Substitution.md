Arguments-Substitution.md
-------------------------


参数代换


根据 WEB 业务的流量模型, 我们假定业务模式是有限的(实际上的确是有限的), 从而可以推导出 GraphQL Schema 的模式也是有限的, 
那么对使用频次较高的 Schema 进行缓存 (缓存AST/缓存后端类型和锚定的解析方法/JIT), 即可提升整体性能.

比如, 请求用户信息的 GraphQL 的 Schema, 基本都会请求用户信息的固定几个字段, Schema很少变动. 而变动的地方经常是参数等 Variables 信息.
那么, 我们只要缓存此 Schema 的前后端解析结果, 即可提升性能.

为此, 我们需要将 Variables 等经常变动的参数从 GraphQL Query 中提取出来, 作为单独的 Query Variables 参数. 并将 GraphQL 重写为含有Arguments的语法.
这样的参数代换过程即可将 GraphQL 归一到可缓存的固定 Schema, 进而完成我们的目标.


e.g.

```graphql
query UserInfo{
    User(Id:2){
        Id,
        Name,
        Email,
    }
}
```   

代换为

```graphql
query UserInfo{
    User(Id:$id){
        Id,
        Name,
        Email,
    }
}
```  

和

Query Variables:
```json
{
    "id": 2
}
```

