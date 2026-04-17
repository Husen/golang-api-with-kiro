# Product

A simple in-memory CRUD REST API for managing products. Built as a learning/reference project demonstrating Clean Architecture in Go.

## Core Entity

`Product` — has `id`, `name`, `price` (float64), and `stock` (int).

## API Endpoints

| Method | Path            | Description       |
|--------|-----------------|-------------------|
| GET    | /products       | List all products |
| GET    | /products/:id   | Get by ID         |
| POST   | /products       | Create product    |
| PUT    | /products/:id   | Update product    |
| DELETE | /products/:id   | Delete product    |

Swagger UI is available at `/swagger/index.html`.

Data is stored in-memory (no database). State resets on restart.
