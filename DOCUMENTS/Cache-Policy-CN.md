Cache-Policy-CN.md
------------------

缓存策略
---------


# Desc

本文档描述前端和后端缓存策略



# Frontend

前端的缓存策略基于**request QueryHash/SubstitutedQueryHash**和**参数代换**.

- 如果用户设定了 Query Variables, 那么我们假设用户**没有使用**内嵌参数, 则**不会进行参数代换**, 直接使用当前 QueryHash 查询缓存. 如果缓存未命中, 则直接对 Query 进行 Parse 并缓存结果.
- 如果用户没有设定 Query Variables, 那么我们假设用户**使用了**内嵌参数, 则**会进行**参数代换**, 然后通过 SubstitutedQueryHash 查询缓存. 如果缓存未命中, 则对 SubstitutedQuery 进行 Parse 并缓存结果.
- 如果用户即设定了 Query Variables, 又在 GraphQL Query 中使用内嵌参数, 为了使结果正确, **不会进行参数代换**, 该情况与情况1 策略相同. 因为我们假定用户使用的内嵌参数是固定的, 不经常变化的, 有缓存意义的. 如果不是, 则会浪费性能和内存. 这不是一种良好的使用方法, 因此请尽量避免. 
- 当然, 在参数替换失败的情况, 会退回使用原始数据进行解析.

# Backend

后端缓存策略的核心目标是降低分支判断, 缓存目标 AST 对应的后端逻辑.

数据结构为:

```go
type cachedField struct {
	Name 		    string
	Type 		    int
	Arguments 		map[string]interface{}
	StringifyFunc   StringifyFunc
	ResolveFunction ResolveFunction
}
```

该结构将在第一次执行目标 AST 时缓存 Field 对应的序列化方法 (StringifyFunc) 和解析方法 (ResolveFunction). 
当接下来执行相同目标 AST 时, 直接执行相关方法, 从而避免了分支判断逻辑.  