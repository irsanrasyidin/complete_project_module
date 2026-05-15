# complete_project_module

Submodule ini adalah shared package untuk project `complete_project`.

## Tujuan

- Menyediakan model data yang dipakai bersama.
- Menyediakan repository abstraction untuk akses data.
- Menyediakan utility umum seperti error wrapper.
- Menjadi tempat desain initializer untuk config, multiple database, Redis, dan Kafka.

## Prinsip Desain

`module` harus reusable dan tidak menjalankan aplikasi secara langsung. Package di dalamnya boleh membuat resource, tetapi keputusan resource mana yang dipakai tetap berada di aplikasi pemanggil seperti `backend`.

Yang boleh berada di `module`:

- struct config dan config loader
- database connection factory
- database manager untuk multiple database
- Redis client factory
- Kafka producer/consumer factory
- bootstrap resource umum
- model, repository, dan utility bersama

Yang tidak sebaiknya berada di `module`:

- HTTP route milik backend
- usecase aplikasi spesifik
- template frontend
- logic startup server

## Rancangan Struktur

Desain lengkap module reusable ada di [DESIGN.md](DESIGN.md).

```text
module/
|-- bootstrap/
|   `-- app.go
|-- config/
|   `-- config.go
|-- database/
|   |-- manager.go
|   `-- postgres.go
|-- cache/
|   `-- redis.go
|-- messaging/
|   `-- kafka.go
|-- exception/
|   |-- code.go
|   `-- exception.go
|-- response/
|   `-- api_response.go
|-- models/
|-- repositories/
`-- utils/
```

Struktur yang sudah ada saat ini:

```text
module/
|-- models/
|-- repositories/
`-- utils/
```

## Desain Config

Config sebaiknya menyimpan semua resource eksternal dalam satu object:

```go
type Config struct {
	App       AppConfig
	Database map[string]DatabaseConfig
	Redis    RedisConfig
	Kafka    KafkaConfig
}
```

Contoh environment:

```text
APP_ENV=local
APP_NAME=complete_project_backend

DB_MAIN_DSN=host=localhost user=postgres password=secret dbname=main port=5432 sslmode=disable
DB_AUDIT_DSN=host=localhost user=postgres password=secret dbname=audit port=5432 sslmode=disable

REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

KAFKA_BROKERS=localhost:9092
KAFKA_CLIENT_ID=complete-project-backend
```

## Desain Multiple Database

Database manager menyimpan koneksi berdasarkan nama.

```go
type DBManager struct {
	connections map[string]*gorm.DB
}

func (m *DBManager) DB(name string) (*gorm.DB, error)
func (m *DBManager) Close() error
```

Contoh pemakaian:

```go
db, err := resources.DB("main")
if err != nil {
	return err
}

repo := masterdata.New(db)
```

Nama koneksi yang disarankan:

- `main` - database utama aplikasi
- `audit` - database audit/log activity
- `reporting` - database reporting/read model

## Desain Redis

Redis dibuat optional. Jika config Redis kosong, bootstrap dapat mengembalikan `nil` atau disabled client wrapper.

Kegunaan Redis yang cocok:

- session cache
- token blacklist
- rate limit
- cache response yang mahal
- distributed lock ringan

## Desain Kafka

Kafka dibuat optional dan sebaiknya dipisah antara producer dan consumer.

Kegunaan Kafka yang cocok:

- event audit
- async notification
- integration event antar service
- background processing

Topik awal yang disarankan:

- `wealth.wallet.created`
- `wealth.category.created`
- `wealth.audit.logged`

## Desain Bootstrap

Package `bootstrap` bertugas menggabungkan resource umum:

```go
type Resources struct {
	DB     *database.DBManager
	Redis  *cache.RedisClient
	Kafka  *messaging.KafkaClient
	Config config.Config
}

func Init(ctx context.Context, cfg config.Config) (*Resources, error)
func (r *Resources) Close() error
```

Backend cukup memanggil `bootstrap.Init`, lalu mengambil dependency yang dibutuhkan.

## Package Saat Ini

- `models/` - struktur data seperti `UserModel`, `UserCredential`, wallet, category, dan rule.
- `repositories/` - implementasi repository untuk user dan login.
- `repositories/masterdata/` - repository SQL untuk wallet, income, expense, dan rules.
- `utils/` - utility umum, termasuk `AppError`.

## Dependensi

Module ini menggunakan:

- Go `1.23.4`
- `gorm.io/gorm`

Jika desain Redis dan Kafka mulai diimplementasikan, dependensi yang mungkin ditambahkan:

- Redis client Go
- Kafka client Go
- PostgreSQL driver jika factory database dipindahkan sepenuhnya ke `module`

## Cara Pakai

Import package dari path modul ini:

```go
github.com/irsanrasyidin/complete_project/module
```

Contoh package yang tersedia:

- `module/models`
- `module/repositories`
- `module/repositories/masterdata`
- `module/utils`

Contoh package yang dirancang:

- `module/config`
- `module/database`
- `module/cache`
- `module/messaging`
- `module/bootstrap`
- `module/exception`
- `module/response`

## Desain Response API

Semua API project disarankan memakai response envelope generic dari package `module/response`.

```go
type ApiResponse[T any] struct {
	CorrelationID string    `json:"correlationid"`
	Error         *string   `json:"error"`
	Tin           time.Time `json:"tin"`
	Tout          time.Time `json:"tout"`
	Data          *T        `json:"data"`
	Success       bool      `json:"success"`
	StatusCode    int       `json:"statuscode"`
}
```

Helper yang dirancang:

- `Success`
- `SuccessPtr`
- `Error`
- `EmptySuccess`
- `WriteJSON`
- `WriteSuccess`
- `WriteError`

Dengan pola ini semua service punya format response yang sama, termasuk `correlationid`, timing request, status code, payload, dan error public.

## Desain Exception

Package `module/exception` menyimpan enum error code dan object exception reusable.

Code yang tersedia:

- `INVALID_ARGUMENT`
- `NOT_FOUND`
- `ALREADY_EXISTS`
- `PERMISSION_DENIED`
- `UNAUTHENTICATED`
- `INTERNAL`
- `UNPROCESSABLE_ENTITY`

`response.NewExceptionErrorResponse` akan menyembunyikan detail error internal jika code adalah `INTERNAL`, tetapi tetap mengirim detail public untuk error non-internal.

## Catatan

Karena ini git submodule, pastikan sudah di-init sebelum build project utama:

```powershell
git submodule update --init --recursive
```
