# Raster Media Example API

![](https://raster-media.net/templates/emotion_raster-noton/frontend/_resources/images/logo.png) An example RESTful API cataloging the releases of German experimental music label, [Raster Media](https://raster-media.net/).

Built with golang ( and love ) using [mux](https://pkg.go.dev/github.com/gorilla/mux?utm_source=godoc), [http](https://pkg.go.dev/net/http), and [postgresql](https://www.postgresql.org/).

## Hosting

Example currently hosted using Cloud Run [here](https://rasterapi-tz76zkxxqq-uw.a.run.app/)

Postgresql database hosted at ElephantSQL

## API

`/albums`

- `GET` : Get all albums; Accepts `count` and `start` parameters for amount of results returned and offset

`/album/:id`

- `GET` : Get an album by `id`
- `PUT` : Update an album by `id`
- `DELETE` : Delete an album by `id`

`/album/`

- `POST` : Create a new album

## TODO

- [x] Support basic REST requests
- [x] Test all basic REST functions
- [ ] Automate database migration every week to refresh tables
- [ ] Authentication required for any mutating requests
- [ ] 100% Coverage on package
- [ ] Finish GoDoc documentation
