# S3 upload & read with imgIX

Upload a image to on S3 bucket and render image with imgIX

## Install
```shell
cp .env.example .env
```

## Run Serve

```shell
$ go run *.go
```

## Usage

Use this api with [Insomnia](https://insomnia.rest/download/) for example

## Upload image

### Request

`POST: /api/profiles`

Structured

`Multipart Form`

Headers

`Content-Type: multipart/form-data`

Field

`file: <input file>`

### Response
```json
{
  "data": "picture-5ec409f39ada6f0a4a6dcc48.jpg"
}
```

### Read image with imgIX

Copy last response and use as image param

### Request

`GET: /api/profiles/{my_image}?h=200&w=200`

`GET: /api/profiles/{my_image}`

### Response
```json
{
  "data": {
    "custom": "https://serviceprofile.imgix.net/picture-5ec409f39ada6f0a4a6dcc48.jpg?h=200&w=200",
    "original": "https://serviceprofile.imgix.net/picture-5ec409f39ada6f0a4a6dcc48.jpg"
  }
}
```
