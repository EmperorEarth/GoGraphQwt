package main

import (
	"encoding/json"
	"net/http"
	"log"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/SermoDigital/jose/jws"
	"github.com/SermoDigital/jose/jwt"
)

func customHandler(schema *graphql.Schema) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {

		jwt, err := jws.ParseJWTFromRequest(r)

		opts := handler.NewRequestOptions(r)

		rootValue := map[string]interface{}{
			"response": rw,
			"request":  r,
			"claims":   jwt.Claims(),
		}

		params := graphql.Params{
			Schema:         *schema,
			RequestString:  opts.Query,
			VariableValues: opts.Variables,
			OperationName:  opts.OperationName,
			RootObject:     rootValue,
		}

		result := graphql.Do(params)

		jsonStr, err := json.Marshal(result)

		if err != nil {
			panic(err)
		}

		rw.Header().Set("Content-Type", "application/json")
		
		// I've testing some GraphiQL, this is needed
		rw.Header().Set("Access-Control-Allow-Credentials", "true")
		rw.Header().Set("Access-Control-Allow-Methods", "POST")
		rw.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		rw.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")

		rw.Write(jsonStr)
	}
}

type Ping struct {
	Pong string `json:"pong"`
	User string `json:"user"`
}

func main() {
	pingType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Ping",
		Description: "Ping to pong...",

		Fields: graphql.Fields{
			"pong": &graphql.Field{
				Type: graphql.String,
			},
			"user": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

	queryType := graphql.ObjectConfig{
		Name: "RootQuery",
		Description: "Big Bang!",

		Fields: graphql.Fields{
			"ping": &graphql.Field{
				Type: pingType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					rv := p.Info.RootValue.(map[string]interface{})
					c := rv["claims"].(jwt.Claims)

					log.Println(c)

					r := Ping{
						Pong: "It works!",
						User: c.Get("name").(string),
					}

					return r, nil
				},
			},
		},
	}

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(queryType),
	})

	if err != nil {
		panic(err)
	}

	http.HandleFunc("/graphql", customHandler(&schema))
	http.ListenAndServe(":1337", nil)
}
