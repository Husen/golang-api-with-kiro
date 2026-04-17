# go-products-api

REST API untuk manajemen produk dibangun dengan Go, Gin, dan Clean Architecture. Storage in-memory, tidak butuh database eksternal. Dilengkapi JWT authentication dan dokumentasi API via Scalar.

## Tech Stack

- **Go** 1.22
- **Gin** v1.10.0
- **JWT** via `golang-jwt/jwt`
- **API Docs** via Scalar (spec di-generate oleh swaggo)
- **Docker** + Docker Compose

## Menjalankan Project

Tidak perlu install Go secara lokal — semua berjalan di Docker.

```bash
# 1. Clone repo
git clone https://github.com/your-username/go-products-api.git
cd go-products-api

# 2. Setup env
cp .env.example .env
# Edit .env dan isi JWT_SECRET

# 3. Jalankan
docker-compose up --build
```

API tersedia di `http://localhost:8080`
Dokumentasi API di `http://localhost:8080/docs`

## Endpoints

### Auth
| Method | Path         | Deskripsi        | Auth |
|--------|--------------|------------------|------|
| POST   | /auth/login  | Login, dapat JWT | -    |

### Products
| Method | Path            | Deskripsi           | Auth    |
|--------|-----------------|---------------------|---------|
| GET    | /products       | List semua produk   | Bearer  |
| GET    | /products/:id   | Detail produk       | Bearer  |
| POST   | /products       | Buat produk baru    | Bearer  |
| PUT    | /products/:id   | Update produk       | Bearer  |
| DELETE | /products/:id   | Hapus produk        | Bearer  |

## Struktur Project

```
go-products-api/
├── main.go
├── internal/
│   ├── domain/        # Entity structs + interfaces
│   ├── handler/       # HTTP layer (Gin handlers)
│   ├── middleware/    # JWT auth + role middleware
│   ├── repository/    # In-memory data access
│   └── usecase/       # Business logic
├── Dockerfile
├── docker-compose.yml
└── .env.example
```

## Environment Variables

Lihat `.env.example` untuk daftar lengkap konfigurasi yang tersedia.
