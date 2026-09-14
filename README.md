# standardgh

Type-safe HTTP handlers for Go's `net/http` using generics. Bind requests, encode responses, validate input, and handle errors with zero boilerplate.

## What is a Generic Handler?

A generic handler is a typed wrapper that removes boilerplate from your handlers. But the real benefit goes beyond convenience.

**Zero framework dependencies.** Your handler function receives a typed request struct and returns a typed response struct — no framework imports, no context casting, no framework-specific code. This means your business logic is decoupled from the router.

```go
// This handler has ZERO framework imports
func CreateUserHandler(ctx context.Context, req *CreateUserReq) (*CreateUserResp, int, error) {
    // pure business logic
}
```

**Router-agnostic.** Since your handlers don't depend on any framework, you could switch between chi, echo, Fiber, or any other router without rewriting handler code. Router changes are rare in real projects, but this flexibility means your core logic is never held hostage by a framework choice.

**Readable request/response handling.** The most common boilerplate in web apps is reading requests and writing responses. Struct tags make this declarative and self-documenting — you see exactly what comes in and what goes out at a glance, without digging through handler code.

```go
// You can immediately see: name from body, age from query, token from cookie
type Req struct {
    Name  string `json:"name"`
    Age   int    `query:"age"`
    Token string `cookie:"token"`
}
```

## Install

```bash
go get github.com/ksckaan1/standardgh
```

## Usage

### `GH` — Generic HTTP Handler

`GH[Req, Resp]` wraps a typed handler function into an `http.HandlerFunc`. Your handler returns `(response, statusCode, error)` — the `int` is the HTTP status code sent to the client.

It automatically:

- Binds request body, headers, query params, URI params, and cookies into your `Req` struct
- Validates the struct using `go-playground/validator/v10`
- Encodes response headers from `header` struct tags
- Sets/clears response cookies from `cookie` struct tags
- Returns JSON responses with proper status codes
- Converts errors into `{"error": "..."}` JSON responses

```go
package main

import (
	"context"
	"log"
	"net/http"

	"github.com/ksckaan1/standardgh"
)

type CreateUserReq struct {
	Name  string `json:"name"`
	Email string `json:"email" validate:"required,email"`
}

type CreateUserResp struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Token string `cookie:"token,24h" json:"-"`
}

func CreateUserHandler(ctx context.Context, req *CreateUserReq) (*CreateUserResp, int, error) {
	// Your business logic here
	return &CreateUserResp{
		ID:    "abc-123",
		Name:  req.Name,
		Token: "secret-token",
	}, http.StatusCreated, nil
}

func main() {
	http.HandleFunc("/users", standardgh.GH(CreateUserHandler))

	log.Fatal(http.ListenAndServe(":3000", nil))
}
```

### `GH` with Multipart Form

Use `multipart.FileHeader` for file uploads with `multipart/form-data` content type:

```go
package main

import (
	"context"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/ksckaan1/standardgh"
)

type UploadReq struct {
	Title string                `json:"title"`
	File  *multipart.FileHeader `json:"file"`
}

type UploadResp struct {
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
}

func UploadHandler(ctx context.Context, req *UploadReq) (*UploadResp, int, error) {
	return &UploadResp{
		Filename: req.File.Filename,
		Size:     req.File.Size,
	}, http.StatusOK, nil
}

func main() {
	http.HandleFunc("/upload", standardgh.GH(UploadHandler))

	log.Fatal(http.ListenAndServe(":3000", nil))
}
```

```bash
curl -X POST http://localhost:3000/upload \
  -F "title=My Document" \
  -F "file=@report.pdf"
```

### Struct Tags

**Request binding** uses struct tags to map incoming data to your struct fields:

```go
type Req struct {
	Name      string `json:"name"`         // body (JSON)
	Age       int    `query:"age"`         // query param
	UserID    string `header:"X-User-ID"`  // request header (case-insensitive)
	PostID    string `uri:"id"`            // path param (Go 1.22+ r.PathValue)
	SessionID string `cookie:"session_id"` // cookie (case-sensitive)
}
```

**Response encoding** uses `header` and `cookie` struct tags on the response struct:

```go
type Resp struct {
	Name      string `json:"name"`
	AuthToken string `cookie:"token,24h" json:"-"`     // sets cookie, omits from JSON
	Session   string `cookie:"session,clear" json:"-"`  // clears cookie
	CustomID  string `header:"X-Custom-ID" json:"-"`    // sets response header
}
```

### Binding Order

The binding order is: **URI → Body → Header → Query → Cookie → Validate**.

This means body fields can override URI params if they share the same name, and validation runs after all binding is complete.

### Cookie Tag Options

The `cookie` tag format is `cookie:"<name>,<duration>"` or `cookie:"<name>,clear"`.

Additional options are set via separate struct tags:

| Tag                 | Description        | Default |
| ------------------- | ------------------ | ------- |
| `cookiePath`        | Cookie path        | `/`     |
| `cookieSameSite`    | SameSite attribute | `Lax`   |
| `cookieSecure`      | Secure flag        | `false` |
| `cookieHTTPOnly`    | HttpOnly flag      | `false` |
| `cookieDomain`      | Domain attribute   | (empty) |
| `cookieMaxAge`      | Max-Age in seconds | `0`     |
| `cookiePartitioned` | Partitioned flag   | `false` |
| `cookieSessionOnly` | Session-only flag  | `false` |

```go
type Resp struct {
	Token string `cookie:"token,24h" cookiePath:"/" cookieSecure:"true" cookieHTTPOnly:"true" json:"-"`
}
```

### Default Values

Set default values using the `default` option in the struct tag. Supported types: `bool`, `int`, `uint`, `float`, `string`, slices (use `|` as separator), and pointers to these types.

```go
type Req struct {
	Name     string   `query:"name,default:john"`
	Age      int      `query:"age,default:18"`
	Products []string `query:"products,default:shoe|hat"`
}
```

### Validation

`standardgh` includes built-in validation using `go-playground/validator/v10`. Validation runs automatically after all binding completes.

```go
type Req struct {
	Name  string `json:"name" validate:"required"`
	Age   int    `json:"age" validate:"gte=0,lte=150"`
	Email string `json:"email" validate:"required,email"`
}
```

Validation errors return **422 Unprocessable Entity** with a formatted message:

```json
{
	"error": "validation failed: Name is required, Email must be a valid email"
}
```

### Error Handling

Errors from your handler are converted to JSON responses:

```go
func Handler(ctx context.Context, req *Req) (*Resp, int, error) {
	return nil, http.StatusBadRequest, fmt.Errorf("invalid input")
}
```

Response:

```json
{
	"error": "invalid input"
}
```

### Binding Errors

Binding errors (type conversion, missing required fields, unsupported content types) return **400 Bad Request**:

```json
{
	"error": "invalid query parameter: age must be an integer"
}
```

### `GHforSSE` — Server-Sent Events

`GHforSSE[Req, Data]` provides type-safe SSE with request binding and named events. It accepts a retry duration and a handler that receives the parsed request and a `send` function.

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ksckaan1/standardgh"
)

type ChatReq struct {
	RoomID string `query:"room"`
}

type ChatEvent struct {
	Message string `json:"message"`
}

func chatHandler(ctx context.Context, req *ChatReq, send func(name string, data ChatEvent) error) error {
	for i := 0; i < 10; i++ {
		if err := send("message", ChatEvent{Message: fmt.Sprintf("hello %d", i)}); err != nil {
			return err
		}
		time.Sleep(time.Second)
	}
	return nil
}

func main() {
	http.HandleFunc("/stream", standardgh.GHforSSE(5*time.Second, chatHandler))

	log.Fatal(http.ListenAndServe(":3000", nil))
}
```

The handler returns an error to signal abnormal termination (e.g., context cancellation, connection lost). Successful completion returns `nil`.

### SSE Event Format

Each event includes:
- `event:` — optional event name (omitted if empty)
- `id:` — UUID v7 (time-ordered, globally unique)
- `data:` — JSON-serialized payload

```
event: message
id: 0192e4c8-8b6a-7c5d-9e0f-1a2b3c4d5e6f
data: {"message":"hello 0"}

event: message
id: 0192e4c8-8b6a-7c5d-9e0f-1a2b3c4d5e70
data: {"message":"hello 1"}
```

## License

[MIT](LICENSE)
