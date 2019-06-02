# Swagger Generated Client

+ [api-spec-converter](https://www.npmjs.com/package/api-spec-converter)

```
api-spec-converter --from=swagger_1 --to=swagger_2 --syntax=json spec/spec-1.0.json  > spec/spec-2.0.json
```

+ [go-swagger](https://github.com/go-spec/go-swagger)

```
swagger generate client -f spec/spec-2.0.json  --default-scheme https -A marvel -c mclient
```
