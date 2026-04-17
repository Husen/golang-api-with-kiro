# Products API

REST API sederhana untuk manajemen produk, dibangun menggunakan Go dan Gin dengan pendekatan Clean Architecture. Data disimpan secara in-memory sehingga tidak memerlukan database eksternal.

---

## Fitur Utama

- CRUD lengkap untuk entitas produk (list, detail, create, update, delete)
- Arsitektur berlapis (Clean Architecture): domain, repository, usecase, handler
- Thread-safe in-memory storage menggunakan `sync.Mutex`
- Dokumentasi API interaktif via Swagger UI (`/swagger/index.html`)
- Siap dijalankan dengan Docker dan Docker Compose

---

## Instalasi dan Menjalankan Proyek

### Prasyarat

- Go 1.22 atau lebih baru
- [swag CLI](https://github.com/swaggo/swag) (untuk generate docs)
- Docker dan Docker Compose (opsional)

### Menjalankan Secara Lokal

```bash
# 1. Clone repositori
git clone <repository-url>
cd go-products-api

# 2. Install dependensi
go mod download

# 3. Generate Swagger docs (wajib sebelum build jika anotasi berubah)
swag init

# 4. Jalankan server
go run .
```

Server akan berjalan di `http://localhost:8080`.

### Menjalankan dengan Docker Compose

```bash
docker-compose up --build
```

Server akan berjalan di `http://localhost:8080` dengan `GIN_MODE=release`.

---

## Endpoint API

| Method | Path            | Deskripsi              |
|--------|-----------------|------------------------|
| GET    | /products       | Ambil semua produk     |
| GET    | /products/:id   | Ambil produk by ID     |
| POST   | /products       | Buat produk baru       |
| PUT    | /products/:id   | Update produk by ID    |
| DELETE | /products/:id   | Hapus produk by ID     |

Dokumentasi interaktif tersedia di: `http://localhost:8080/swagger/index.html`

### Contoh Request Body (POST / PUT)

```json
{
  "name": "Laptop",
  "price": 15000000,
  "stock": 10
}
```

---

## Struktur Folder

```
go-products-api/
├── main.go                          # Entry point: wiring dependensi, registrasi route
├── go.mod
├── Dockerfile
├── docker-compose.yml
├── docs/                            # Auto-generated oleh swag init, jangan diedit manual
└── internal/
    ├── domain/
    │   └── product.go               # Struct entitas + interface repository & usecase
    ├── handler/
    │   └── product_handler.go       # Layer HTTP: Gin handlers + registrasi route
    ├── repository/
    │   └── product_repository.go    # Akses data: in-memory slice dengan mutex
    └── usecase/
        └── product_usecase.go       # Business logic: orkestrasi pemanggilan repository
```

---

## Teknologi yang Digunakan

| Teknologi         | Versi    | Keterangan                          |
|-------------------|----------|-------------------------------------|
| Go                | 1.22     | Bahasa pemrograman utama            |
| Gin               | v1.10.0  | HTTP web framework                  |
| Swaggo/swag       | v1.16.3  | Generator dokumentasi Swagger       |
| Swaggo/gin-swagger| v1.6.0   | Middleware Swagger UI untuk Gin     |
| Docker            | -        | Containerisasi aplikasi             |
| Docker Compose    | -        | Orkestrasi container lokal          |
