# go-products-api

REST API untuk manajemen produk dengan JWT authentication, dibangun menggunakan Go, Gin, dan Clean Architecture.

## Menjalankan

```bash
cp .env.example .env
# isi JWT_SECRET di .env

docker-compose up --build
```

- API: `http://localhost:8080`
- Docs: `http://localhost:8080/docs`

## Endpoints

### Auth
| Method | Path        | Deskripsi     |
|--------|-------------|---------------|
| POST   | /auth/login | Login, dapat JWT |

### Products (butuh Bearer token)
| Method | Path          | Deskripsi        |
|--------|---------------|------------------|
| GET    | /products     | List semua produk |
| GET    | /products/:id | Detail produk    |
| POST   | /products     | Buat produk baru |
| PUT    | /products/:id | Update produk    |
| DELETE | /products/:id | Hapus produk     |

## Struktur

```
internal/
├── domain/      # Entity + interfaces
├── handler/     # HTTP layer
├── middleware/  # JWT + role middleware
├── repository/  # In-memory storage
└── usecase/     # Business logic
```
