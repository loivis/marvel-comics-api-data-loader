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
