# Task 2: HTTP service

This project is an HTTP service that lets users manage their products. User management is provided by [AIexMoran/httpCRUD](https://github.com/AIexMoran/httpCRUD) service.

## How to run

`Docker` and `docker-compose` must be installed on the system to be able to run this project.

- `make up` creates and runs Docker containers of this project and `AIexMoran/httpCRUD` a path to which must be specified in the `docker-compose.yml` file.

- `make down` stops and removes containers, networks, images, and volumes.

- `make build` builds or rebuilds services.

- `make start` starts services.

- `make stop` stops services.

## Usage

`GET /products/` lists all products.

`GET /products/{id}` shows product details by the specified id.

`POST /products/` adds a product to the product list. Authorization required.

`PUT /products/{id}` replaces the product with a new one by the specified id. Authorization required.

`DELETE /products/{id}` removes the product by the specified id. Authorization required.

### Authorization

The request is expected to have an `Authorization` header with the token issued by `AIexMoran/httpCRUD`. The usage may be found [here](https://github.com/AIexMoran/httpCRUD).

Example:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
```