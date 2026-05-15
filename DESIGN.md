# Reusable Module Design

Dokumen ini adalah desain awal untuk menjadikan `complete_project_module` sebagai module reusable lintas project.

Use case pertama yang ingin didukung adalah aplikasi monitoring income dan outcome dari WhatsApp. Contoh utama: GoPay mengirim notifikasi WhatsApp ketika ada outcome, lalu sistem membaca pesan tersebut, parsing transaksi, menyimpan data, dan menghasilkan report.

Desain ini sengaja dibuat generic agar module tidak terkunci hanya untuk project finance atau WhatsApp.

## Tujuan Utama

- Satu module bisa dipakai ulang di project baru tanpa menulis ulang utility dasar.
- Resource umum seperti config, database, Redis, Kafka, logger, dan lifecycle dikelola konsisten.
- Utility kecil seperti encode/decode, hash SHA-256, pointer helper, parser, validator, dan error wrapper tersedia sebagai package standar.
- Domain spesifik aplikasi tetap berada di aplikasi pemakai, bukan dipaksa masuk ke module.

## Prinsip Batasan

`module` boleh berisi:

- reusable infrastructure
- reusable helper
- reusable abstraction
- reusable parser engine
- reusable repository pattern
- shared model yang benar-benar umum

`module` tidak boleh berisi:

- route HTTP aplikasi tertentu
- usecase bisnis aplikasi tertentu
- template frontend
- hardcoded provider seperti GoPay sebagai satu-satunya format
- logic yang hanya masuk akal untuk satu project

Jika logic hanya berlaku untuk aplikasi monitoring income/outcome, logic tersebut berada di `backend`. Jika logic itu bisa dipakai untuk membaca banyak jenis pesan notifikasi, logic bisa masuk ke `module/parser`.

## Target Struktur Package

```text
module/
|-- bootstrap/
|   `-- resources.go
|-- config/
|   |-- config.go
|   `-- env.go
|-- database/
|   |-- manager.go
|   |-- postgres.go
|   `-- transaction.go
|-- cache/
|   `-- redis.go
|-- messaging/
|   |-- kafka.go
|   |-- publisher.go
|   `-- consumer.go
|-- logger/
|   `-- logger.go
|-- exception/
|   |-- code.go
|   `-- exception.go
|-- response/
|   `-- api_response.go
|-- codec/
|   |-- base64.go
|   |-- json.go
|   `-- url.go
|-- hash/
|   `-- sha256.go
|-- ptr/
|   `-- ptr.go
|-- parser/
|   |-- text.go
|   |-- money.go
|   |-- time.go
|   `-- rule.go
|-- errors/
|   `-- app_error.go
|-- validate/
|   `-- validate.go
|-- id/
|   `-- id.go
|-- models/
|-- repositories/
`-- utils/
```

Package yang sudah ada seperti `models`, `repositories`, dan `utils` tetap bisa dipakai, lalu pelan-pelan dirapikan agar mengikuti struktur di atas.

## Layer Arsitektur

```text
frontend
   |
   v
backend application
   |
   |-- usecase income/outcome monitoring
   |-- handler WhatsApp webhook/import
   |-- scheduler/reporting
   |
   v
module reusable packages
   |-- bootstrap/config/database/cache/messaging
   |-- response/exception/parser/codec/hash/ptr/errors/validate
   |-- repositories/models
   |
   v
external resources
   |-- PostgreSQL
   |-- Redis
   |-- Kafka
   `-- WhatsApp provider/gateway
```

## Flow Monitoring WhatsApp

Target flow aplikasi:

```text
WhatsApp message
   |
   v
backend receive webhook/import
   |
   v
normalize text
   |
   v
parse candidate transaction
   |
   v
classify income/outcome/provider
   |
   v
deduplicate message
   |
   v
store transaction
   |
   v
publish event
   |
   v
report/dashboard/notification
```

Untuk GoPay:

```text
GoPay WA notification
   |
   v
parser extracts amount, direction, merchant, time, reference
   |
   v
backend maps result into transaction domain
```

Parser di `module` sebaiknya hanya menghasilkan data generic seperti `ParsedMessage` atau `ParsedTransactionCandidate`. Keputusan final bahwa itu income/outcome resmi milik aplikasi dilakukan di backend.

## Desain Parser

Parser dibuat rule-based terlebih dahulu agar mudah dites dan tidak bergantung ke AI.

Konsep:

```go
type ParsedTransactionCandidate struct {
	Provider    string
	Direction   string
	Amount      int64
	Currency    string
	Merchant    string
	Reference   string
	OccurredAt  *time.Time
	RawText     string
	Confidence  float64
}
```

Package parser menyediakan helper:

- normalize whitespace
- lowercase safe comparison
- extract money amount
- extract date/time
- extract reference number
- detect keyword
- apply ordered parsing rules

Contoh package:

```text
parser.Text.Normalize(input)
parser.Money.ParseIDR(input)
parser.Rule.Apply(input, rules)
```

Format GoPay, bank, e-wallet lain, atau provider baru bisa dibuat sebagai rule set terpisah.

## Desain Deduplication

Pesan WhatsApp bisa terkirim ulang atau diproses dua kali. Deduplication wajib ada.

Strategi generic:

```text
dedup_key = sha256(provider + normalized_message + occurred_at + amount)
```

Package `hash` menyediakan SHA-256 helper, tetapi backend yang menentukan komposisi key.

Dedup bisa disimpan di:

- PostgreSQL unique index
- Redis key dengan TTL
- keduanya jika butuh lebih kuat

## Desain Config

Config harus typed dan punya default aman.

```go
type Config struct {
	App       AppConfig
	HTTP      HTTPConfig
	Database  map[string]DatabaseConfig
	Redis     RedisConfig
	Kafka     KafkaConfig
	WhatsApp  WhatsAppConfig
}
```

`WhatsAppConfig` tetap generic:

```go
type WhatsAppConfig struct {
	Provider      string
	WebhookSecret string
	VerifyToken   string
}
```

Module boleh menyediakan struct config, tetapi detail provider dan validasi bisnis tetap bisa ditambahkan oleh backend.

## Desain Bootstrap Resource

Bootstrap bertugas membuat dan menutup resource.

```go
type Resources struct {
	Config config.Config
	DB     *database.Manager
	Redis  cache.Client
	Kafka  messaging.Client
	Logger logger.Logger
}
```

Kontrak:

- `Init` membuat resource yang enabled di config.
- Resource optional boleh disabled tanpa error.
- Resource wajib seperti database utama harus gagal cepat jika tidak bisa connect.
- `Close` menutup semua resource dengan urutan aman.

## Desain Database

Multiple database dikelola berdasarkan nama.

```go
type Manager interface {
	DB(name string) (*gorm.DB, error)
	MustDB(name string) *gorm.DB
	Close() error
}
```

Nama koneksi standar:

- `main`
- `audit`
- `reporting`
- `archive`

Repository menerima `*gorm.DB` dari aplikasi, bukan membuat koneksi sendiri.

## Desain Kafka

Kafka abstraction minimal:

```go
type Publisher interface {
	Publish(ctx context.Context, topic string, key string, payload any) error
}
```

Event generic:

```go
type Event struct {
	ID          string
	Type        string
	Source      string
	OccurredAt  time.Time
	Payload     any
	Metadata    map[string]string
}
```

Contoh event aplikasi monitoring:

- `transaction.detected`
- `transaction.created`
- `transaction.duplicated`
- `transaction.parse_failed`

## Desain Redis

Redis abstraction minimal:

```go
type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
}
```

Use case awal:

- dedup WhatsApp message
- rate limit webhook
- temporary parser cache
- lock saat import batch

## Desain API Response

Semua API buatan project sebaiknya memakai response envelope yang sama. Ini membuat frontend, logging, tracing, dan error handling lebih mudah dirawat.

Package yang disarankan:

```text
module/response
```

Struct utama:

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

Field:

- `CorrelationID` - ID request untuk tracing dari request masuk sampai response keluar.
- `Error` - pesan error public yang aman dikirim ke client.
- `Tin` - waktu request mulai diproses.
- `Tout` - waktu response dibuat.
- `Data` - payload response, generic agar bisa dipakai semua endpoint.
- `Success` - status sukses/gagal.
- `StatusCode` - HTTP status code yang juga ditulis di body.

Helper yang wajib disediakan:

```go
func Success[T any](correlationID string, tin time.Time, statusCode int, data T) ApiResponse[T]
func SuccessPtr[T any](correlationID string, tin time.Time, statusCode int, data *T) ApiResponse[T]
func Error[T any](correlationID string, tin time.Time, statusCode int, message string) ApiResponse[T]
func EmptySuccess(correlationID string, tin time.Time, statusCode int) ApiResponse[any]
```

Helper HTTP yang disarankan:

```go
func WriteJSON[T any](w http.ResponseWriter, res ApiResponse[T])
func WriteSuccess[T any](w http.ResponseWriter, correlationID string, tin time.Time, statusCode int, data T)
func WriteError(w http.ResponseWriter, correlationID string, tin time.Time, statusCode int, message string)
```

Aturan response:

- Semua handler harus mengembalikan `ApiResponse`.
- `CorrelationID` diambil dari header request jika ada, jika tidak generate baru.
- Nama header yang disarankan: `X-Correlation-ID`.
- `Tin` dibuat sedekat mungkin dengan awal request.
- `Tout` dibuat saat response helper dipanggil.
- `Error` harus `nil` saat success.
- `Data` boleh `nil` saat error atau response kosong.
- `StatusCode` di body harus sama dengan HTTP status code.
- Error internal tidak boleh langsung dibocorkan ke `Error`.

Contoh success:

```json
{
  "correlationid": "01HX8S7J9Q9Y3E7Z9J4T8V7P5A",
  "error": null,
  "tin": "2026-05-15T10:00:00Z",
  "tout": "2026-05-15T10:00:00.025Z",
  "data": {
    "id": "trx-001",
    "amount": 25000
  },
  "success": true,
  "statuscode": 200
}
```

Contoh error:

```json
{
  "correlationid": "01HX8S7J9Q9Y3E7Z9J4T8V7P5A",
  "error": "transaction not found",
  "tin": "2026-05-15T10:00:00Z",
  "tout": "2026-05-15T10:00:00.018Z",
  "data": null,
  "success": false,
  "statuscode": 404
}
```

Middleware backend yang disarankan:

```text
request masuk
   |
   v
set tin
   |
   v
get/generate correlation id
   |
   v
inject ke request context
   |
   v
handler
   |
   v
response.WriteSuccess / response.WriteError
```

Context helper yang disarankan:

```go
func CorrelationIDFromContext(ctx context.Context) string
func TinFromContext(ctx context.Context) time.Time
```

Dengan pola ini handler tidak perlu membuat envelope manual berulang-ulang.

## Desain Exception

Package yang disarankan:

```text
module/exception
```

Enum error code:

```go
type Code string

const (
	InvalidArgumentCode     Code = "INVALID_ARGUMENT"
	NotFoundCode            Code = "NOT_FOUND"
	AlreadyExistsCode       Code = "ALREADY_EXISTS"
	PermissionDeniedCode    Code = "PERMISSION_DENIED"
	UnauthenticatedCode     Code = "UNAUTHENTICATED"
	InternalErrorCode       Code = "INTERNAL"
	UnprocessableEntityCode Code = "UNPROCESSABLE_ENTITY"
)
```

Mapping default HTTP status:

- `INVALID_ARGUMENT` -> `400`
- `NOT_FOUND` -> `404`
- `ALREADY_EXISTS` -> `409`
- `PERMISSION_DENIED` -> `403`
- `UNAUTHENTICATED` -> `401`
- `INTERNAL` -> `500`
- `UNPROCESSABLE_ENTITY` -> `422`

Object exception:

```go
type Exception struct {
	Code    Code
	Message string
}
```

Implementation boleh menyimpan internal error dan detail public sebagai field private. Response helper harus menyembunyikan detail internal ketika code adalah `INTERNAL`.

## Utility Wajib

Utility ini sangat reusable dan layak masuk `module`.

### codec

Fungsi:

- JSON encode/decode
- Base64 encode/decode
- URL encode/decode

Prinsip:

- return error, jangan panic
- support generic JSON decode
- tidak mencampur encoding dengan crypto

### hash

Fungsi:

- SHA-256 string
- SHA-256 bytes
- SHA-256 hex
- compare hash string jika diperlukan

Prinsip:

- deterministic
- output default hex lowercase
- tidak digunakan untuk password hashing

Password harus memakai bcrypt/argon2, bukan SHA-256 biasa.

### ptr

Fungsi:

- membuat pointer dari value
- mengambil value dari pointer dengan default
- safe dereference
- helper untuk string/int/bool/time pointer

Contoh kebutuhan:

```go
name := ptr.Value(input.Name, "")
limit := ptr.Value(input.Limit, 20)
enabled := ptr.Value(input.Enabled, false)
```

### parser

Fungsi:

- normalize text
- parse money
- parse date/time
- parse keyword
- parse rule-based message

Parser harus deterministic dan mudah dites.

### errors

Fungsi:

- application error
- error code
- public message
- internal cause
- wrapping

Error harus bisa dipakai oleh HTTP handler tanpa membocorkan detail internal.

### validate

Fungsi:

- required string
- min/max length
- enum validation
- numeric range
- simple email/phone validation

Validasi domain kompleks tetap di backend usecase.

## Domain Monitoring Income/Outcome

Domain ini berada di backend atau aplikasi pemakai, bukan di module reusable.

Model aplikasi yang disarankan:

```text
transactions
|-- id
|-- direction
|-- amount
|-- currency
|-- provider
|-- merchant
|-- reference
|-- occurred_at
|-- source
|-- raw_message
|-- dedup_key
|-- created_at
|-- updated_at
```

Direction:

- `income`
- `outcome`
- `unknown`

Source:

- `whatsapp`
- `manual`
- `import`
- `api`

## WhatsApp Integration Boundary

Module tidak perlu tahu cara spesifik menerima WhatsApp. Backend yang menerima input dari:

- WhatsApp Business API webhook
- provider gateway
- manual paste/import
- forwarded message endpoint

Module hanya membantu:

- verify signature jika generic helper tersedia
- normalize text
- parse candidate transaction
- hash dedup key
- encode/decode payload

## Testing Strategy

Setiap package reusable harus punya unit test table-driven.

Prioritas test:

- parser money IDR
- parser text normalization
- API response success/error envelope
- correlation ID propagation
- SHA-256 deterministic output
- pointer helper nil/non-nil
- config default
- DB manager missing connection
- disabled Redis/Kafka behavior
- error mapping

Untuk WhatsApp parser:

- simpan sample pesan sebagai fixture
- test setiap provider/rule
- pastikan pesan gagal parse menghasilkan confidence rendah, bukan panic

## Migration Strategy

Implementasi sebaiknya bertahap:

1. Rapikan utility dasar: `ptr`, `hash`, `codec`, `errors`.
2. Tambah reusable `response.ApiResponse` dan HTTP writer helper.
3. Tambah `config` typed loader.
4. Tambah `database.Manager`.
5. Pindahkan init PostgreSQL dari backend ke `module/database`.
6. Tambah `bootstrap.Resources`.
7. Tambah parser text/money/time.
8. Tambah parser rule untuk WhatsApp transaction candidate.
9. Tambah Redis optional.
10. Tambah Kafka optional.
11. Baru implement domain monitoring income/outcome di backend.

Urutan ini menjaga module tetap reusable dan menghindari backend terlalu cepat bergantung pada desain yang belum stabil.

## Keputusan Desain Awal

- `module` akan menjadi reusable shared library, bukan framework aplikasi penuh.
- Backend tetap menjadi composition root.
- Parser dibuat deterministic dan rule-based terlebih dahulu.
- SHA-256 dipakai untuk hash/dedup, bukan password.
- Redis dan Kafka optional.
- Multiple database dikelola dengan named connection.
- Semua API response memakai `response.ApiResponse[T]` dari module.
- Correlation ID wajib tersedia di setiap response.
- Utility kecil dipisah ke package sendiri agar mudah dipakai lintas project.
- Logic GoPay spesifik tidak hardcoded sebagai inti module; ia menjadi rule set/parser provider yang bisa ditambah atau diganti.
