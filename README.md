# mcapi-loader

Load [Marvel Comics API](https://developer.marvel.com/) data into MongoDB

# Instructions

## Signup on [Marvel Developer Portal](https://developer.marvel.com/) for API key

## Get ready MongoDB

```
docker run -d --name mongo -p 27017:27017 mongo
```

## Run it !!!

```
MARVEL_API_PRIVATE_KEY="private_key" MARVE_API_PUBLIC_KEY="public_key" MONGODB_URI="mongodb://localhost:27017/marvel-comics" go run main.go
```

# ISSUE

+ `limit` and `offset` not as expected

    Change on any of them will get a totally different response, at least for `/v1/public/comics`.

+ Count not always identical

    `available` for list of comics/events/series/stories returned from `/v1/public/characters/{characterId}` and `total` returned from `/v1/public/characters/{characterId}/{comics/events/series/stories}`

+ Some field types don't follow api definition

    + comic 39237 with `diamondCode` should be `string` but returns `number`

    + comic 70668 with `diamondCode` should be `string` but returns `number`

    + comic 27399 with `diamondCode` should be `string` but returns `number`

    + comic 1405 with `isbn` should be `string` but returns `number`

    + story 30688 with `title` should be `string` but returns `number`

    + story 44568 with `title` should be `string` but returns `number`

    + creator 9551 with `suffix` should be `string` but returns `number`

    + creator 10669 with `lastName` should be `string` but returns `number`

    + and more ...

# Tools

+ [api-spec-converter](https://www.npmjs.com/package/api-spec-converter)

```
api-spec-converter --from=swagger_1 --to=swagger_2 --syntax=json swagger/spec-1.0.json  > swagger/spec-2.0.json
```

+ [go-swagger](https://github.com/go-swagger/go-swagger)

```
swagger generate client -f swagger/spec-2.0.json  --default-scheme https -A marvel -c mclient
```
