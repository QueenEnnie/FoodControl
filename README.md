# FoodControl

FoodControl is a Go backend service for tracking food inventory and checking which recipes can be cooked from available products.

## Features

- PostgreSQL schema for products, recipes, and recipe ingredients.
- Database access through `pgx/v5` and `pgxpool`.
- REST API built on the Go standard library `net/http`.
- JSON request and response bodies.
- Recipe availability checks that ignore expired products.
- Docker and Docker Compose setup for local deployment.

## Run With Docker

```bash
docker compose up --build
```

The API listens on `http://localhost:8080`.
PostgreSQL is exposed on `localhost:5433`.

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
POST /products
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

### Recipes

```http
GET /recipes
POST /recipes
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
