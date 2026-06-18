# FoodControl

FoodControl is a Go backend service for tracking food inventory and checking which recipes can be cooked from available products.

## Features

- PostgreSQL schema for products, recipes, and recipe ingredients.
- Database access through `pgx/v5` and `pgxpool`.
- REST API built on the Go standard library `net/http`.
- JSON request and response bodies.
- Recipe availability checks that ignore expired products.
- Atomic ingredient consumption with serializable transactions and row-level locking.
- Expired and soon-to-expire inventory filtering with a configurable time window.
- Soft deletion for products, so removed inventory is preserved in the database but hidden from API results and recipe availability checks.
- Structured JSON logs with request IDs, request logging middleware, and panic recovery middleware.
- Docker and Docker Compose setup for local deployment.

## Run With Docker

```bash
docker compose up --build
```

The API listens on `http://localhost:8080`.
PostgreSQL is exposed on `localhost:5433`.
The service writes structured JSON logs to stdout. In Docker, view them with
`docker compose logs app`.

The local database is seeded with demo products and recipes, so you can call
`GET /suggestions` right after startup and see which recipes can be cooked from
the available inventory.

If you already started the database before adding seed data, recreate the local
volume:

```bash
docker compose down -v
docker compose up --build
```

## Environment

```text
HTTP_ADDR=:8080
DATABASE_URL=postgres://fcuser:fcpass@localhost:5433/fooddb
```

Inside Docker Compose, the application uses:

```text
DATABASE_URL=postgres://fcuser:fcpass@db:5432/fooddb
```

## API

### Health

```http
GET /health
```

### Products

```http
GET /products
GET /products/{id}
GET /products/expiring?days=3
POST /products
PUT /products/{id}
DELETE /products/{id}
```

Example product:

```json
{
  "name": "milk",
  "description": "2.5%",
  "quantity": 1,
  "unit": "l",
  "expiration_date": "2026-05-25T00:00:00Z"
}
```

Get one product:

```http
GET /products/1
```

Successful product lookups return `200 OK` with the product.
If the product does not exist, the API returns `404 Not Found`.

List expired products and products expiring within the next three days:

```http
GET /products/expiring?days=3
```

The `days` parameter is optional and defaults to `3`. It accepts values from
`0` to `365`; `0` returns products expiring today or already expired.

Update an existing product:

```http
PUT /products/1
Content-Type: application/json
```

```json
{
  "name": "milk",
  "description": "2.5%",
  "quantity": 0.5,
  "unit": "l",
  "expiration_date": "2026-05-25T00:00:00Z"
}
```

Successful product updates return `200 OK` with the updated product.
If the product does not exist, the API returns `404 Not Found`.

Delete uses soft deletion:

```http
DELETE /products/1
```

Successful deletes return `204 No Content`. Soft-deleted products are preserved
in the database with `deleted_at`, but they are hidden from product endpoints and
ignored by recipe availability checks.

### Recipes

```http
GET /recipes
POST /recipes
POST /recipes/{id}/cook
GET /recipes/{id}/availability
GET /suggestions
```

Example recipe:

```json
{
  "name": "Pancakes",
  "description": "Simple breakfast recipe",
  "ingredients": [
    {
      "product_name": "milk",
      "required_quantity": 0.3,
      "unit": "l"
    },
    {
      "product_name": "flour",
      "required_quantity": 200,
      "unit": "g"
    }
  ]
}
```

Cook a recipe and atomically consume its ingredients:

```http
POST /recipes/1/cook
```

```json
{
  "recipe_id": 1,
  "recipe_name": "Pancakes",
  "status": "cooked",
  "consumed": [
    {
      "product_name": "milk",
      "quantity": 0.5,
      "unit": "l"
    }
  ]
}
```

Cooking runs in a serializable PostgreSQL transaction and locks matching
inventory rows. If ingredients are missing, expired, soft-deleted, or
insufficient, the entire transaction is rolled back and the API returns
`409 Conflict`.
