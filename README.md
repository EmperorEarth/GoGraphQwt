# Go/Golang + GraphQL + JWT

Just another study-case example.

`docker-compose up`

```bash
curl http://<YOUR_DOCKER_MACHINE_IP>:1337/graphql

-H "Authorization: Bearer
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.
eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.
TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"

-d query="query { ping { pong user } }"
```
This token was borrowed from [JWT](https://jwt.io/) debugger.
```json
{
  "data": {
    "ping": {
      "pong": "It works!",
      "user": "John Doe"
    }
  }
}
```
