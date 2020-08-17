call-map.md
-----------



# Server Side

```golang

schema := // sechema settings

query := `
    query {
        Users(id:10){
                id
                name
                email
        }
    }
`


params := graphql.Params{
    Schema:        schema,
    RequestString: query,
}

// execute 

result := graphql.Do(params)

```