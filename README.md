# Raster Media Example API

An example RESTful API cataloging the releases of German experimental music label, [Raster Media](https://raster-media.net/).

Built with Go, gorilla/mux, net/http, and postgresql.

## Hosting

Example currently hosted using Cloud Run [here](https://rasterapi-tz76zkxxqq-uw.a.run.app/)

Postgresql database hosted at ElephantSQL

## API

`/albums`

- `GET` : Get all albums

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
