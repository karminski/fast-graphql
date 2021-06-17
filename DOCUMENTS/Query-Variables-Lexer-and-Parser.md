Query-Variables-Lexer-and-Parser.md
-----------------------------------


由于参数代换后, 所有的变量全都集中在 Query Varibales 中, 因此 Query Varibales 的解析性能至关重要.  
然而 Go 内置的 encoding/json 库并不是为了这种 variables dictionary 场景而准备的, 因此性能并不理想. 
大概消耗了整体请求过程耗时的 10%, 作为对比, 执行缓存流程的 fast-graphql 耗时占比为 30% 左右, 其余是 net/http 库和 runtime, I/O 耗时).


本质上应该是一个 json 的子集.

# QueryVariables EBNF 定义

QueryVariables ::= Ignored "{" Ignored Variables+ Ignored "}" Ignored
Variables      ::= Ignored Name Ignored ":" Ignored Value Ignored
Name           ::= '"' [_A-Za-z][_0-9A-Za-z]* '"'
Value          ::= IntValue | FloatValue | StringValue | BooleanValue | NullValue 