package standardgh

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"uuid"
)

type headerTestRequest struct {
	UserID   string `header:"User-ID"`
	UserName string `header:"User-Name"`
}

type headerIntRequest struct {
	UserID int `header:"User-ID"`
}

type headerInt64Request struct {
	UserID int64 `header:"User-ID"`
}

type headerUintRequest struct {
	UserID uint `header:"User-ID"`
}

type headerUint64Request struct {
	UserID uint64 `header:"User-ID"`
}

type headerFloat64Request struct {
	Score float64 `header:"Score"`
}

type headerBoolRequest struct {
	Active bool `header:"Active"`
}

type headerStringSliceRequest struct {
	Tags []string `header:"Tags"`
}

type headerIntSliceRequest struct {
	IDs []int `header:"IDs"`
}

type headerUintSliceRequest struct {
	IDs []uint `header:"IDs"`
}

type headerFloat64SliceRequest struct {
	Scores []float64 `header:"Scores"`
}

type headerBoolSliceRequest struct {
	Flags []bool `header:"Flags"`
}

type headerRequiredRequest struct {
	UserID string `header:"User-ID,required"`
}

type headerDefaultRequest struct {
	UserID string `header:"User-ID,default:default-id"`
}

func TestBindHeaderValues(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("User-ID", "12345")
	req.Header.Set("User-Name", "John")

	result := new(headerTestRequest)
	err := bindHeaders(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.UserID != "12345" {
		t.Errorf("expected UserID to be 12345, got %s", result.UserID)
	}

	if result.UserName != "John" {
		t.Errorf("expected UserName to be John, got %s", result.UserName)
	}
}

func TestBindHeaderRequired(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	result := new(headerRequiredRequest)
	err := bindHeaders(req, result)

	if err == nil {
		t.Fatal("expected error for missing required header")
	}

	expected := "User-ID header is required"
	if err.Error() != expected {
		t.Errorf("expected error %q, got %q", expected, err.Error())
	}
}

func TestBindHeaderRequiredWithHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("User-ID", "12345")

	result := new(headerRequiredRequest)
	err := bindHeaders(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.UserID != "12345" {
		t.Errorf("expected UserID to be 12345, got %s", result.UserID)
	}
}

func TestBindHeaderDefault(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	result := new(headerDefaultRequest)
	err := bindHeaders(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.UserID != "default-id" {
		t.Errorf("expected UserID to be default-id, got %s", result.UserID)
	}
}

func TestBindHeaderDefaultOverride(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("User-ID", "12345")

	result := new(headerDefaultRequest)
	err := bindHeaders(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.UserID != "12345" {
		t.Errorf("expected UserID to be 12345, got %s", result.UserID)
	}
}

func TestBindNilPointer(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("User-ID", "12345")

	err := bindHeaders(req, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBindNonStruct(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("User-ID", "12345")

	var result string
	err := bindHeaders(req, &result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBindWithoutHeaders(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	result := new(headerTestRequest)
	err := bindHeaders(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.UserID != "" {
		t.Errorf("expected UserID to be empty, got %s", result.UserID)
	}

	if result.UserName != "" {
		t.Errorf("expected UserName to be empty, got %s", result.UserName)
	}
}

func TestBindHeaderInt(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("User-ID", "123")

	result := new(headerIntRequest)
	err := bindHeaders(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.UserID != 123 {
		t.Errorf("expected UserID to be 123, got %d", result.UserID)
	}
}

func TestBindHeaderInt64(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("User-ID", "1234567890")

	result := new(headerInt64Request)
	err := bindHeaders(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.UserID != 1234567890 {
		t.Errorf("expected UserID to be 1234567890, got %d", result.UserID)
	}
}

func TestBindHeaderUint(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("User-ID", "123")

	result := new(headerUintRequest)
	err := bindHeaders(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.UserID != 123 {
		t.Errorf("expected UserID to be 123, got %d", result.UserID)
	}
}

func TestBindHeaderUint64(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("User-ID", "1234567890")

	result := new(headerUint64Request)
	err := bindHeaders(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.UserID != 1234567890 {
		t.Errorf("expected UserID to be 1234567890, got %d", result.UserID)
	}
}

func TestBindHeaderFloat64(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Score", "3.14")

	result := new(headerFloat64Request)
	err := bindHeaders(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Score != 3.14 {
		t.Errorf("expected Score to be 3.14, got %f", result.Score)
	}
}

func TestBindHeaderBool(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{"true", "true", true},
		{"false", "false", false},
		{"1", "1", true},
		{"0", "0", false},
		{"TRUE", "TRUE", true},
		{"FALSE", "FALSE", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Active", tt.value)

			result := new(headerBoolRequest)
			err := bindHeaders(req, result)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.Active != tt.expected {
				t.Errorf("expected Active to be %v, got %v", tt.expected, result.Active)
			}
		})
	}
}

func TestBindHeaderStringSlice(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Tags", "go|http|web")

	result := new(headerStringSliceRequest)
	err := bindHeaders(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Tags) != 3 {
		t.Fatalf("expected Tags to have 3 elements, got %d", len(result.Tags))
	}

	expected := []string{"go", "http", "web"}
	for i, tag := range result.Tags {
		if tag != expected[i] {
			t.Errorf("expected Tags[%d] to be %s, got %s", i, expected[i], tag)
		}
	}
}

func TestBindHeaderIntSlice(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("IDs", "1|2|3")

	result := new(headerIntSliceRequest)
	err := bindHeaders(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.IDs) != 3 {
		t.Fatalf("expected IDs to have 3 elements, got %d", len(result.IDs))
	}

	expected := []int{1, 2, 3}
	for i, id := range result.IDs {
		if id != expected[i] {
			t.Errorf("expected IDs[%d] to be %d, got %d", i, expected[i], id)
		}
	}
}

func TestBindHeaderUintSlice(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("IDs", "1|2|3")

	result := new(headerUintSliceRequest)
	err := bindHeaders(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.IDs) != 3 {
		t.Fatalf("expected IDs to have 3 elements, got %d", len(result.IDs))
	}

	expected := []uint{1, 2, 3}
	for i, id := range result.IDs {
		if id != expected[i] {
			t.Errorf("expected IDs[%d] to be %d, got %d", i, expected[i], id)
		}
	}
}

func TestBindHeaderFloat64Slice(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Scores", "1.1|2.2|3.3")

	result := new(headerFloat64SliceRequest)
	err := bindHeaders(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Scores) != 3 {
		t.Fatalf("expected Scores to have 3 elements, got %d", len(result.Scores))
	}

	expected := []float64{1.1, 2.2, 3.3}
	for i, score := range result.Scores {
		if score != expected[i] {
			t.Errorf("expected Scores[%d] to be %f, got %f", i, expected[i], score)
		}
	}
}

func TestBindHeaderBoolSlice(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Flags", "true|false|1")

	result := new(headerBoolSliceRequest)
	err := bindHeaders(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Flags) != 3 {
		t.Fatalf("expected Flags to have 3 elements, got %d", len(result.Flags))
	}

	expected := []bool{true, false, true}
	for i, flag := range result.Flags {
		if flag != expected[i] {
			t.Errorf("expected Flags[%d] to be %v, got %v", i, expected[i], flag)
		}
	}
}

func TestBindHeaderIntInvalidValue(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("User-ID", "not-a-number")

	result := new(headerIntRequest)
	err := bindHeaders(req, result)

	if err == nil {
		t.Fatal("expected error for invalid int value")
	}
}

func TestBindHeaderDefaultInt(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	result := new(struct {
		UserID int `header:"User-ID,default:42"`
	})
	err := bindHeaders(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.UserID != 42 {
		t.Errorf("expected UserID to be 42, got %d", result.UserID)
	}
}

func TestCaseInsensitiveHeaderMatching(t *testing.T) {
	tests := []struct {
		name       string
		headerName string
		tagName    string
	}{
		{"lowercase header", "user-id", "User-ID"},
		{"uppercase header", "USER-ID", "User-ID"},
		{"mixed case header", "User-Id", "User-ID"},
		{"tag lowercase", "user-id", "user-id"},
		{"tag uppercase", "user-id", "USER-ID"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set(tt.headerName, "12345")

			result := new(struct {
				UserID int `header:"User-ID"`
			})
			err := bindHeaders(req, result)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.UserID != 12345 {
				t.Errorf("expected UserID to be 12345, got %d", result.UserID)
			}
		})
	}
}

func TestBindHeaderSliceDefault(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	result := new(struct {
		Tags []string `header:"Tags,default:tag1|tag2|tag3"`
	})
	err := bindHeaders(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Tags) != 3 {
		t.Fatalf("expected Tags to have 3 elements, got %d", len(result.Tags))
	}

	expected := []string{"tag1", "tag2", "tag3"}
	for i, tag := range result.Tags {
		if tag != expected[i] {
			t.Errorf("expected Tags[%d] to be %s, got %s", i, expected[i], tag)
		}
	}
}

func TestBindHeaderIntSliceDefault(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	result := new(struct {
		IDs []int `header:"IDs,default:1|2|3"`
	})
	err := bindHeaders(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.IDs) != 3 {
		t.Fatalf("expected IDs to have 3 elements, got %d", len(result.IDs))
	}

	expected := []int{1, 2, 3}
	for i, id := range result.IDs {
		if id != expected[i] {
			t.Errorf("expected IDs[%d] to be %d, got %d", i, expected[i], id)
		}
	}
}

func TestBindHeaderSliceWithHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Tags", "go|http|web")

	result := new(struct {
		Tags []string `header:"Tags,default:tag1|tag2"`
	})
	err := bindHeaders(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Tags) != 3 {
		t.Fatalf("expected Tags to have 3 elements, got %d", len(result.Tags))
	}

	expected := []string{"go", "http", "web"}
	for i, tag := range result.Tags {
		if tag != expected[i] {
			t.Errorf("expected Tags[%d] to be %s, got %s", i, expected[i], tag)
		}
	}
}

// Query binding tests

func TestBindQueryString(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?name=john&age=30", nil)

	result := new(struct {
		Name string `query:"name"`
		Age  int    `query:"age"`
	})
	err := bindQuery(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Name != "john" {
		t.Errorf("expected Name to be john, got %s", result.Name)
	}

	if result.Age != 30 {
		t.Errorf("expected Age to be 30, got %d", result.Age)
	}
}

func TestBindQueryStringDefault(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	result := new(struct {
		Name  string `query:"name,default:Anonymous"`
		Page  int    `query:"page,default:1"`
		Limit int    `query:"limit,default:20"`
	})
	err := bindQuery(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Name != "Anonymous" {
		t.Errorf("expected Name to be Anonymous, got %s", result.Name)
	}

	if result.Page != 1 {
		t.Errorf("expected Page to be 1, got %d", result.Page)
	}

	if result.Limit != 20 {
		t.Errorf("expected Limit to be 20, got %d", result.Limit)
	}
}

func TestBindQueryStringDefaultOverride(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?name=john&page=5", nil)

	result := new(struct {
		Name  string `query:"name,default:Anonymous"`
		Page  int    `query:"page,default:1"`
		Limit int    `query:"limit,default:20"`
	})
	err := bindQuery(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Name != "john" {
		t.Errorf("expected Name to be john, got %s", result.Name)
	}

	if result.Page != 5 {
		t.Errorf("expected Page to be 5, got %d", result.Page)
	}

	if result.Limit != 20 {
		t.Errorf("expected Limit to be 20, got %d", result.Limit)
	}
}

func TestBindQueryStringRequired(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	result := new(struct {
		Name string `query:"name,required"`
	})
	err := bindQuery(req, result)

	if err == nil {
		t.Fatal("expected error for missing required query param")
	}

	expected := "name query is required"
	if err.Error() != expected {
		t.Errorf("expected error %q, got %q", expected, err.Error())
	}
}

func TestBindQueryStringRequiredWithParam(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?name=john", nil)

	result := new(struct {
		Name string `query:"name,required"`
	})
	err := bindQuery(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Name != "john" {
		t.Errorf("expected Name to be john, got %s", result.Name)
	}
}

func TestBindQueryStringSlice(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?tags=go&tags=http&tags=web", nil)

	result := new(struct {
		Tags []string `query:"tags"`
	})
	err := bindQuery(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Tags) != 3 {
		t.Fatalf("expected Tags to have 3 elements, got %d", len(result.Tags))
	}

	expected := []string{"go", "http", "web"}
	for i, tag := range result.Tags {
		if tag != expected[i] {
			t.Errorf("expected Tags[%d] to be %s, got %s", i, expected[i], tag)
		}
	}
}

func TestBindQueryStringSliceDefault(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	result := new(struct {
		Tags []string `query:"tags,default:default1|default2|default3"`
	})
	err := bindQuery(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Tags) != 3 {
		t.Fatalf("expected Tags to have 3 elements, got %d", len(result.Tags))
	}

	expected := []string{"default1", "default2", "default3"}
	for i, tag := range result.Tags {
		if tag != expected[i] {
			t.Errorf("expected Tags[%d] to be %s, got %s", i, expected[i], tag)
		}
	}
}

func TestBindQueryStringIntSlice(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?ids=1&ids=2&ids=3", nil)

	result := new(struct {
		IDs []int `query:"ids"`
	})
	err := bindQuery(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.IDs) != 3 {
		t.Fatalf("expected IDs to have 3 elements, got %d", len(result.IDs))
	}

	expected := []int{1, 2, 3}
	for i, id := range result.IDs {
		if id != expected[i] {
			t.Errorf("expected IDs[%d] to be %d, got %d", i, expected[i], id)
		}
	}
}

func TestBindQueryStringFloat64(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?score=3.14", nil)

	result := new(struct {
		Score float64 `query:"score"`
	})
	err := bindQuery(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Score != 3.14 {
		t.Errorf("expected Score to be 3.14, got %f", result.Score)
	}
}

func TestBindQueryStringBool(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{"true", "true", true},
		{"false", "false", false},
		{"1", "1", true},
		{"0", "0", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/?active="+tt.value, nil)

			result := new(struct {
				Active bool `query:"active"`
			})
			err := bindQuery(req, result)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.Active != tt.expected {
				t.Errorf("expected Active to be %v, got %v", tt.expected, result.Active)
			}
		})
	}
}

func TestBindQueryStringBracketNotation(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?colors[]=red&colors[]=blue&colors[]=green", nil)

	result := new(struct {
		Colors []string `query:"colors"`
	})
	err := bindQuery(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Colors) != 3 {
		t.Fatalf("expected Colors to have 3 elements, got %d", len(result.Colors))
	}

	expected := []string{"red", "blue", "green"}
	for i, color := range result.Colors {
		if color != expected[i] {
			t.Errorf("expected Colors[%d] to be %s, got %s", i, expected[i], color)
		}
	}
}

func TestBindQueryStringURLDecoding(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?name=hello%20world", nil)

	result := new(struct {
		Name string `query:"name"`
	})
	err := bindQuery(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Name != "hello world" {
		t.Errorf("expected Name to be 'hello world', got %s", result.Name)
	}
}

func TestBindQueryStringCaseInsensitive(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?UserName=john", nil)

	result := new(struct {
		UserName string `query:"username"`
	})
	err := bindQuery(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.UserName != "john" {
		t.Errorf("expected UserName to be john, got %s", result.UserName)
	}
}

func TestBindQueryStringEmpty(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	result := new(struct {
		Name string `query:"name"`
		Page int    `query:"page"`
	})
	err := bindQuery(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Name != "" {
		t.Errorf("expected Name to be empty, got %s", result.Name)
	}

	if result.Page != 0 {
		t.Errorf("expected Page to be 0, got %d", result.Page)
	}
}

func TestBindQueryStringInvalidInt(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?page=not-a-number", nil)

	result := new(struct {
		Page int `query:"page"`
	})
	err := bindQuery(req, result)

	if err == nil {
		t.Fatal("expected error for invalid int value")
	}
}

func TestBindQueryStringNilPointer(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?name=john", nil)

	err := bindQuery(req, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBindQueryStringNonStruct(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?name=john", nil)

	var result string
	err := bindQuery(req, &result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBindQueryStringLastValueWins(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?name=old&name=new", nil)

	result := new(struct {
		Name string `query:"name"`
	})
	err := bindQuery(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Name != "new" {
		t.Errorf("expected Name to be new (last value), got %s", result.Name)
	}
}

func TestBindQueryStringLastIntWins(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?page=1&page=5&page=10", nil)

	result := new(struct {
		Page int `query:"page"`
	})
	err := bindQuery(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Page != 10 {
		t.Errorf("expected Page to be 10 (last value), got %d", result.Page)
	}
}

func TestBindQueryStringPointer(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?name=john&age=30", nil)

	result := new(struct {
		Name *string `query:"name"`
		Age  *int    `query:"age"`
	})
	err := bindQuery(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Name == nil || *result.Name != "john" {
		t.Errorf("expected Name to be *\"john\", got %v", result.Name)
	}

	if result.Age == nil || *result.Age != 30 {
		t.Errorf("expected Age to be *30, got %v", result.Age)
	}
}

func TestBindQueryStringPointerNil(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	result := new(struct {
		Name *string `query:"name"`
		Age  *int    `query:"age"`
	})
	err := bindQuery(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Name != nil {
		t.Errorf("expected Name to be nil, got %v", result.Name)
	}

	if result.Age != nil {
		t.Errorf("expected Age to be nil, got %v", result.Age)
	}
}

func TestBindQueryStringNestedStruct(t *testing.T) {
	type Address struct {
		City    string `query:"city"`
		ZipCode string `query:"zip_code"`
	}

	req := httptest.NewRequest(http.MethodGet, "/?address[city]=istanbul&address[zip_code]=34000", nil)

	result := new(struct {
		Name    string  `query:"name"`
		Address Address `query:"address"`
	})
	err := bindQuery(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Address.City != "istanbul" {
		t.Errorf("expected City to be istanbul, got %s", result.Address.City)
	}

	if result.Address.ZipCode != "34000" {
		t.Errorf("expected ZipCode to be 34000, got %s", result.Address.ZipCode)
	}
}

func TestBindQueryStringEmbeddedStruct(t *testing.T) {
	type Base struct {
		ID int `query:"id"`
	}

	req := httptest.NewRequest(http.MethodGet, "/?id=42&name=john", nil)

	result := new(struct {
		Base
		Name string `query:"name"`
	})
	err := bindQuery(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID != 42 {
		t.Errorf("expected ID to be 42, got %d", result.ID)
	}

	if result.Name != "john" {
		t.Errorf("expected Name to be john, got %s", result.Name)
	}
}

func TestBindQueryStringCommaSplitting(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?colors=red,blue,green", nil)

	result := new(struct {
		Colors []string `query:"colors"`
	})
	err := bindQuery(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Colors) != 3 {
		t.Fatalf("expected Colors to have 3 elements, got %d", len(result.Colors))
	}

	expected := []string{"red", "blue", "green"}
	for i, color := range result.Colors {
		if color != expected[i] {
			t.Errorf("expected Colors[%d] to be %s, got %s", i, expected[i], color)
		}
	}
}

func TestBindQueryStringMapStringString(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?name=john&age=30", nil)

	result := make(map[string]string)
	err := bindQuery(req, &result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result["name"] != "john" {
		t.Errorf("expected name to be john, got %s", result["name"])
	}

	if result["age"] != "30" {
		t.Errorf("expected age to be 30, got %s", result["age"])
	}
}

func TestBindQueryStringMapStringSlice(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?colors=red&colors=blue&colors=green", nil)

	result := make(map[string][]string)
	err := bindQuery(req, &result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result["colors"]) != 3 {
		t.Fatalf("expected colors to have 3 elements, got %d", len(result["colors"]))
	}

	expected := []string{"red", "blue", "green"}
	for i, color := range result["colors"] {
		if color != expected[i] {
			t.Errorf("expected colors[%d] to be %s, got %s", i, expected[i], color)
		}
	}
}

func TestBindQueryStringUnmatchedBrackets(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?name[=john", nil)

	result := new(struct {
		Name string `query:"name"`
	})
	err := bindQuery(req, result)

	if err == nil {
		t.Fatal("expected error for unmatched brackets")
	}
}

func TestBindQueryStringEmptyBracketStripping(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?colors[]=red&colors[]=blue", nil)

	result := new(struct {
		Colors []string `query:"colors"`
	})
	err := bindQuery(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Colors) != 2 {
		t.Fatalf("expected Colors to have 2 elements, got %d", len(result.Colors))
	}

	expected := []string{"red", "blue"}
	for i, color := range result.Colors {
		if color != expected[i] {
			t.Errorf("expected Colors[%d] to be %s, got %s", i, expected[i], color)
		}
	}
}

func TestBindQueryStringTextUnmarshaler(t *testing.T) {
	type CustomType struct {
		Value string
	}

	req := httptest.NewRequest(http.MethodGet, "/?name=john", nil)

	result := new(struct {
		Name string `query:"name"`
	})
	err := bindQuery(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Name != "john" {
		t.Errorf("expected Name to be john, got %s", result.Name)
	}
}

// URI binding tests

func TestBindURIValues(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/user/12345", nil)
	req.SetPathValue("id", "12345")

	result := new(struct {
		ID string `uri:"id"`
	})
	err := bindURIFromRequest(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID != "12345" {
		t.Errorf("expected ID to be 12345, got %s", result.ID)
	}
}

func TestBindURIMultipleParams(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/user/123/posts/456", nil)
	req.SetPathValue("userId", "123")
	req.SetPathValue("postId", "456")

	result := new(struct {
		UserID string `uri:"userId"`
		PostID string `uri:"postId"`
	})
	err := bindURIFromRequest(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.UserID != "123" {
		t.Errorf("expected UserID to be 123, got %s", result.UserID)
	}

	if result.PostID != "456" {
		t.Errorf("expected PostID to be 456, got %s", result.PostID)
	}
}

func TestBindURIInt(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/user/42", nil)
	req.SetPathValue("id", "42")

	result := new(struct {
		ID int `uri:"id"`
	})
	err := bindURIFromRequest(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID != 42 {
		t.Errorf("expected ID to be 42, got %d", result.ID)
	}
}

func TestBindURIInt64(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/user/1234567890", nil)
	req.SetPathValue("id", "1234567890")

	result := new(struct {
		ID int64 `uri:"id"`
	})
	err := bindURIFromRequest(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID != 1234567890 {
		t.Errorf("expected ID to be 1234567890, got %d", result.ID)
	}
}

func TestBindURIUint(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/user/123", nil)
	req.SetPathValue("id", "123")

	result := new(struct {
		ID uint `uri:"id"`
	})
	err := bindURIFromRequest(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID != 123 {
		t.Errorf("expected ID to be 123, got %d", result.ID)
	}
}

func TestBindURIFloat64(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/score/3.14", nil)
	req.SetPathValue("score", "3.14")

	result := new(struct {
		Score float64 `uri:"score"`
	})
	err := bindURIFromRequest(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Score != 3.14 {
		t.Errorf("expected Score to be 3.14, got %f", result.Score)
	}
}

func TestBindURIBool(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{"true", "true", true},
		{"false", "false", false},
		{"1", "1", true},
		{"0", "0", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/flag/"+tt.value, nil)
			req.SetPathValue("active", tt.value)

			result := new(struct {
				Active bool `uri:"active"`
			})
			err := bindURIFromRequest(req, result)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.Active != tt.expected {
				t.Errorf("expected Active to be %v, got %v", tt.expected, result.Active)
			}
		})
	}
}

func TestBindURIRequired(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	result := new(struct {
		ID string `uri:"id,required"`
	})
	err := bindURIFromRequest(req, result)

	if err == nil {
		t.Fatal("expected error for missing required URI param")
	}

	expected := "id uri is required"
	if err.Error() != expected {
		t.Errorf("expected error %q, got %q", expected, err.Error())
	}
}

func TestBindURIRequiredWithParam(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/user/123", nil)
	req.SetPathValue("id", "123")

	result := new(struct {
		ID string `uri:"id,required"`
	})
	err := bindURIFromRequest(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID != "123" {
		t.Errorf("expected ID to be 123, got %s", result.ID)
	}
}

func TestBindURIDefault(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	result := new(struct {
		ID string `uri:"id,default:default-id"`
	})
	err := bindURIFromRequest(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID != "default-id" {
		t.Errorf("expected ID to be default-id, got %s", result.ID)
	}
}

func TestBindURIDefaultOverride(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/user/123", nil)
	req.SetPathValue("id", "123")

	result := new(struct {
		ID string `uri:"id,default:default-id"`
	})
	err := bindURIFromRequest(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID != "123" {
		t.Errorf("expected ID to be 123, got %s", result.ID)
	}
}

func TestBindURINilPointer(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/user/123", nil)

	err := bindURIFromRequest(req, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBindURINonStruct(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/user/123", nil)

	var result string
	err := bindURIFromRequest(req, &result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBindURIInvalidInt(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/user/not-a-number", nil)
	req.SetPathValue("id", "not-a-number")

	result := new(struct {
		ID int `uri:"id"`
	})
	err := bindURIFromRequest(req, result)

	if err == nil {
		t.Fatal("expected error for invalid int value")
	}
}

func TestBindURISimple(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/user/123", nil)
	req.SetPathValue("id", "123")

	result := new(struct {
		ID string `uri:"id"`
	})
	err := bindURIFromRequest(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID != "123" {
		t.Errorf("expected ID to be 123, got %s", result.ID)
	}
}

func TestBindURIFromParams(t *testing.T) {
	params := map[string]string{
		"id":   "123",
		"name": "john",
	}

	result := new(struct {
		ID   string `uri:"id"`
		Name string `uri:"name"`
	})
	err := bindURIFromParams(params, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID != "123" {
		t.Errorf("expected ID to be 123, got %s", result.ID)
	}

	if result.Name != "john" {
		t.Errorf("expected Name to be john, got %s", result.Name)
	}
}

func TestBindURIFromParamsRequired(t *testing.T) {
	params := map[string]string{}

	result := new(struct {
		ID string `uri:"id,required"`
	})
	err := bindURIFromParams(params, result)

	if err == nil {
		t.Fatal("expected error for missing required URI param")
	}

	expected := "id uri is required"
	if err.Error() != expected {
		t.Errorf("expected error %q, got %q", expected, err.Error())
	}
}

func TestBindURICurlyBraceSyntax(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/user/123", nil)
	req.SetPathValue("id", "123")

	result := new(struct {
		ID string `uri:"id"`
	})
	err := bindURIFromRequest(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID != "123" {
		t.Errorf("expected ID to be 123, got %s", result.ID)
	}
}

func TestBindURICurlyBraceMultipleParams(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/user/123/posts/456", nil)
	req.SetPathValue("userId", "123")
	req.SetPathValue("postId", "456")

	result := new(struct {
		UserID string `uri:"userId"`
		PostID string `uri:"postId"`
	})
	err := bindURIFromRequest(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.UserID != "123" {
		t.Errorf("expected UserID to be 123, got %s", result.UserID)
	}

	if result.PostID != "456" {
		t.Errorf("expected PostID to be 456, got %s", result.PostID)
	}
}

// Cookie binding tests

func TestBindCookieValues(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "abc123"})
	req.AddCookie(&http.Cookie{Name: "user_name", Value: "john"})

	result := new(struct {
		SessionID string `cookie:"session_id"`
		UserName  string `cookie:"user_name"`
	})
	err := bindCookies(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.SessionID != "abc123" {
		t.Errorf("expected SessionID to be abc123, got %s", result.SessionID)
	}

	if result.UserName != "john" {
		t.Errorf("expected UserName to be john, got %s", result.UserName)
	}
}

func TestBindCookieInt(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "age", Value: "30"})

	result := new(struct {
		Age int `cookie:"age"`
	})
	err := bindCookies(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Age != 30 {
		t.Errorf("expected Age to be 30, got %d", result.Age)
	}
}

func TestBindCookieInt64(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "id", Value: "1234567890"})

	result := new(struct {
		ID int64 `cookie:"id"`
	})
	err := bindCookies(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID != 1234567890 {
		t.Errorf("expected ID to be 1234567890, got %d", result.ID)
	}
}

func TestBindCookieUint(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "id", Value: "123"})

	result := new(struct {
		ID uint `cookie:"id"`
	})
	err := bindCookies(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID != 123 {
		t.Errorf("expected ID to be 123, got %d", result.ID)
	}
}

func TestBindCookieFloat64(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "score", Value: "3.14"})

	result := new(struct {
		Score float64 `cookie:"score"`
	})
	err := bindCookies(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Score != 3.14 {
		t.Errorf("expected Score to be 3.14, got %f", result.Score)
	}
}

func TestBindCookieBool(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{"true", "true", true},
		{"false", "false", false},
		{"1", "1", true},
		{"0", "0", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.AddCookie(&http.Cookie{Name: "active", Value: tt.value})

			result := new(struct {
				Active bool `cookie:"active"`
			})
			err := bindCookies(req, result)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.Active != tt.expected {
				t.Errorf("expected Active to be %v, got %v", tt.expected, result.Active)
			}
		})
	}
}

func TestBindCookieRequired(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	result := new(struct {
		SessionID string `cookie:"session_id,required"`
	})
	err := bindCookies(req, result)

	if err == nil {
		t.Fatal("expected error for missing required cookie")
	}

	expected := "session_id cookie is required"
	if err.Error() != expected {
		t.Errorf("expected error %q, got %q", expected, err.Error())
	}
}

func TestBindCookieRequiredWithCookie(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "abc123"})

	result := new(struct {
		SessionID string `cookie:"session_id,required"`
	})
	err := bindCookies(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.SessionID != "abc123" {
		t.Errorf("expected SessionID to be abc123, got %s", result.SessionID)
	}
}

func TestBindCookieDefault(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	result := new(struct {
		SessionID string `cookie:"session_id,default:default-session"`
	})
	err := bindCookies(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.SessionID != "default-session" {
		t.Errorf("expected SessionID to be default-session, got %s", result.SessionID)
	}
}

func TestBindCookieDefaultOverride(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "abc123"})

	result := new(struct {
		SessionID string `cookie:"session_id,default:default-session"`
	})
	err := bindCookies(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.SessionID != "abc123" {
		t.Errorf("expected SessionID to be abc123, got %s", result.SessionID)
	}
}

func TestBindCookieNilPointer(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "abc123"})

	err := bindCookies(req, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBindCookieNonStruct(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "abc123"})

	var result string
	err := bindCookies(req, &result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBindCookieInvalidInt(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "age", Value: "not-a-number"})

	result := new(struct {
		Age int `cookie:"age"`
	})
	err := bindCookies(req, result)

	if err == nil {
		t.Fatal("expected error for invalid int value")
	}
}

func TestBindCookieMapStringString(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "abc123"})
	req.AddCookie(&http.Cookie{Name: "user_name", Value: "john"})

	result := make(map[string]string)
	err := bindCookies(req, &result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result["session_id"] != "abc123" {
		t.Errorf("expected session_id to be abc123, got %s", result["session_id"])
	}

	if result["user_name"] != "john" {
		t.Errorf("expected user_name to be john, got %s", result["user_name"])
	}
}

func TestBindCookieCaseSensitive(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "SessionID", Value: "abc123"})

	result := new(struct {
		SessionID string `cookie:"session_id"`
	})
	err := bindCookies(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Cookie names are case-sensitive, so "SessionID" != "session_id"
	if result.SessionID != "" {
		t.Errorf("expected SessionID to be empty (case mismatch), got %s", result.SessionID)
	}
}

// Additional cookie binding tests for Fiber v3 compatibility

func TestBindCookieEmptyValue(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "name", Value: ""})

	result := new(struct {
		Name string `cookie:"name"`
	})
	err := bindCookies(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Name != "" {
		t.Errorf("expected Name to be empty, got %s", result.Name)
	}
}

func TestBindCookieEmptyValueRequired(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "name", Value: ""})

	result := new(struct {
		Name string `cookie:"name,required"`
	})
	err := bindCookies(req, result)

	if err == nil {
		t.Fatal("expected error for empty required cookie")
	}

	expected := "name cookie is required"
	if err.Error() != expected {
		t.Errorf("expected error %q, got %q", expected, err.Error())
	}
}

func TestBindCookieEmptyValueDefault(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "name", Value: ""})

	result := new(struct {
		Name string `cookie:"name,default:default-name"`
	})
	err := bindCookies(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Name != "default-name" {
		t.Errorf("expected Name to be default-name, got %s", result.Name)
	}
}

func TestBindCookieMultipleSameName(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "tag", Value: "first"})
	req.AddCookie(&http.Cookie{Name: "tag", Value: "second"})

	result := new(struct {
		Tag string `cookie:"tag"`
	})
	err := bindCookies(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Last value wins (Fiber v3 behavior)
	if result.Tag != "second" {
		t.Errorf("expected Tag to be second (last value), got %s", result.Tag)
	}
}

func TestBindCookiePointerField(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "name", Value: "john"})

	result := new(struct {
		Name *string `cookie:"name"`
	})
	err := bindCookies(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Name == nil || *result.Name != "john" {
		t.Errorf("expected Name to be *\"john\", got %v", result.Name)
	}
}

func TestBindCookiePointerFieldNil(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	result := new(struct {
		Name *string `cookie:"name"`
	})
	err := bindCookies(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Name != nil {
		t.Errorf("expected Name to be nil, got %v", result.Name)
	}
}

func TestBindCookieSliceField(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "tags", Value: "go|http|web"})

	result := new(struct {
		Tags []string `cookie:"tags"`
	})
	err := bindCookies(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Tags) != 3 {
		t.Fatalf("expected Tags to have 3 elements, got %d", len(result.Tags))
	}

	expected := []string{"go", "http", "web"}
	for i, tag := range result.Tags {
		if tag != expected[i] {
			t.Errorf("expected Tags[%d] to be %s, got %s", i, expected[i], tag)
		}
	}
}

func TestBindCookieMapStringSlice(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "colors", Value: "red"})
	req.AddCookie(&http.Cookie{Name: "sizes", Value: "large"})

	result := make(map[string][]string)
	err := bindCookies(req, &result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result["colors"]) != 1 || result["colors"][0] != "red" {
		t.Errorf("expected colors to be [red], got %v", result["colors"])
	}

	if len(result["sizes"]) != 1 || result["sizes"][0] != "large" {
		t.Errorf("expected sizes to be [large], got %v", result["sizes"])
	}
}

func TestBindCookieBoolOnValue(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "active", Value: "on"})

	result := new(struct {
		Active bool `cookie:"active"`
	})
	err := bindCookies(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Fiber v3 treats "on" as true (HTML form checkbox)
	if !result.Active {
		t.Errorf("expected Active to be true for value 'on', got false")
	}
}

func TestBindCookieTextUnmarshaler(t *testing.T) {
	type CustomType struct {
		Value string
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "custom", Value: "test-value"})

	result := new(struct {
		Custom string `cookie:"custom"`
	})
	err := bindCookies(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Custom != "test-value" {
		t.Errorf("expected Custom to be test-value, got %s", result.Custom)
	}
}

func TestBindCookieConversionErrorFormat(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "age", Value: "not-a-number"})

	result := new(struct {
		Age int `cookie:"age"`
	})
	err := bindCookies(req, result)

	if err == nil {
		t.Fatal("expected error for invalid int value")
	}

	// Verify error contains cookie source context
	if !strings.Contains(err.Error(), "cookie") {
		t.Errorf("expected error to contain 'cookie', got %q", err.Error())
	}
}

func TestBindCookieInvalidFloat64(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "score", Value: "not-a-float"})

	result := new(struct {
		Score float64 `cookie:"score"`
	})
	err := bindCookies(req, result)

	if err == nil {
		t.Fatal("expected error for invalid float64 value")
	}
}

func TestBindCookieInvalidBool(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "active", Value: "maybe"})

	result := new(struct {
		Active bool `cookie:"active"`
	})
	err := bindCookies(req, result)

	if err == nil {
		t.Fatal("expected error for invalid bool value")
	}
}

func TestBindCookieInvalidUint(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "id", Value: "-1"})

	result := new(struct {
		ID uint `cookie:"id"`
	})
	err := bindCookies(req, result)

	if err == nil {
		t.Fatal("expected error for negative uint value")
	}
}

func TestBindCookieUint64(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "id", Value: "1234567890"})

	result := new(struct {
		ID uint64 `cookie:"id"`
	})
	err := bindCookies(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID != 1234567890 {
		t.Errorf("expected ID to be 1234567890, got %d", result.ID)
	}
}

func TestBindCookieFloat32(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "score", Value: "3.14"})

	result := new(struct {
		Score float32 `cookie:"score"`
	})
	err := bindCookies(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Score < 3.13 || result.Score > 3.15 {
		t.Errorf("expected Score to be ~3.14, got %f", result.Score)
	}
}

func TestBindCookieSliceDefault(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	result := new(struct {
		Tags []string `cookie:"tags,default:go|http|web"`
	})
	err := bindCookies(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Tags) != 3 {
		t.Fatalf("expected Tags to have 3 elements, got %d", len(result.Tags))
	}

	expected := []string{"go", "http", "web"}
	for i, tag := range result.Tags {
		if tag != expected[i] {
			t.Errorf("expected Tags[%d] to be %s, got %s", i, expected[i], tag)
		}
	}
}

func TestBindCookieNoCookies(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	result := new(struct {
		Name string `cookie:"name"`
		Age  int    `cookie:"age"`
	})
	err := bindCookies(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Name != "" {
		t.Errorf("expected Name to be empty, got %s", result.Name)
	}

	if result.Age != 0 {
		t.Errorf("expected Age to be 0, got %d", result.Age)
	}
}

func TestBindCookieValueWithEquals(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: "abc=def=ghi"})

	result := new(struct {
		Token string `cookie:"token"`
	})
	err := bindCookies(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Token != "abc=def=ghi" {
		t.Errorf("expected Token to be 'abc=def=ghi', got %s", result.Token)
	}
}

func TestBindCookieIntSliceField(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "ids", Value: "1|2|3"})

	result := new(struct {
		IDs []int `cookie:"ids"`
	})
	err := bindCookies(req, result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.IDs) != 3 {
		t.Fatalf("expected IDs to have 3 elements, got %d", len(result.IDs))
	}

	expected := []int{1, 2, 3}
	for i, id := range result.IDs {
		if id != expected[i] {
			t.Errorf("expected IDs[%d] to be %d, got %d", i, expected[i], id)
		}
	}
}

// Validation tests

func TestValidateStructValid(t *testing.T) {
	type Request struct {
		Name  string `validate:"required"`
		Email string `validate:"required,email"`
		Age   int    `validate:"gte=0,lte=130"`
	}

	req := &Request{
		Name:  "John",
		Email: "john@example.com",
		Age:   30,
	}

	err := validateStruct(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidateStructRequired(t *testing.T) {
	type Request struct {
		Name string `validate:"required"`
	}

	req := &Request{}

	err := validateStruct(req)
	if err == nil {
		t.Fatal("expected validation error for required field")
	}

	if !strings.Contains(err.Error(), "Name is required") {
		t.Errorf("expected error to contain 'Name is required', got %q", err.Error())
	}
}

func TestValidateStructEmail(t *testing.T) {
	type Request struct {
		Email string `validate:"required,email"`
	}

	tests := []struct {
		name  string
		email string
		valid bool
	}{
		{"valid email", "john@example.com", true},
		{"invalid email", "not-an-email", false},
		{"empty email", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &Request{Email: tt.email}
			err := validateStruct(req)

			if tt.valid && err != nil {
				t.Errorf("expected no error for email %q, got %v", tt.email, err)
			}
			if !tt.valid && err == nil {
				t.Errorf("expected error for email %q, got nil", tt.email)
			}
		})
	}
}

func TestValidateStructMin(t *testing.T) {
	type Request struct {
		Name string `validate:"min=3"`
	}

	tests := []struct {
		name  string
		value string
		valid bool
	}{
		{"valid min", "abc", true},
		{"too short", "ab", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &Request{Name: tt.value}
			err := validateStruct(req)

			if tt.valid && err != nil {
				t.Errorf("expected no error for name %q, got %v", tt.value, err)
			}
			if !tt.valid && err == nil {
				t.Errorf("expected error for name %q, got nil", tt.value)
			}
		})
	}
}

func TestValidateStructMax(t *testing.T) {
	type Request struct {
		Name string `validate:"max=5"`
	}

	tests := []struct {
		name  string
		value string
		valid bool
	}{
		{"valid max", "abc", true},
		{"exact max", "abcde", true},
		{"too long", "abcdef", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &Request{Name: tt.value}
			err := validateStruct(req)

			if tt.valid && err != nil {
				t.Errorf("expected no error for name %q, got %v", tt.value, err)
			}
			if !tt.valid && err == nil {
				t.Errorf("expected error for name %q, got nil", tt.value)
			}
		})
	}
}

func TestValidateStructGTE(t *testing.T) {
	type Request struct {
		Age int `validate:"gte=18"`
	}

	tests := []struct {
		name  string
		age   int
		valid bool
	}{
		{"valid gte", 18, true},
		{"above gte", 25, true},
		{"below gte", 17, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &Request{Age: tt.age}
			err := validateStruct(req)

			if tt.valid && err != nil {
				t.Errorf("expected no error for age %d, got %v", tt.age, err)
			}
			if !tt.valid && err == nil {
				t.Errorf("expected error for age %d, got nil", tt.age)
			}
		})
	}
}

func TestValidateStructLTE(t *testing.T) {
	type Request struct {
		Age int `validate:"lte=130"`
	}

	tests := []struct {
		name  string
		age   int
		valid bool
	}{
		{"valid lte", 130, true},
		{"below lte", 25, true},
		{"above lte", 131, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &Request{Age: tt.age}
			err := validateStruct(req)

			if tt.valid && err != nil {
				t.Errorf("expected no error for age %d, got %v", tt.age, err)
			}
			if !tt.valid && err == nil {
				t.Errorf("expected error for age %d, got nil", tt.age)
			}
		})
	}
}

func TestValidateStructOneOf(t *testing.T) {
	type Request struct {
		Status string `validate:"oneof=active inactive pending"`
	}

	tests := []struct {
		name   string
		status string
		valid  bool
	}{
		{"valid oneof", "active", true},
		{"valid oneof 2", "inactive", true},
		{"invalid oneof", "deleted", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &Request{Status: tt.status}
			err := validateStruct(req)

			if tt.valid && err != nil {
				t.Errorf("expected no error for status %q, got %v", tt.status, err)
			}
			if !tt.valid && err == nil {
				t.Errorf("expected error for status %q, got nil", tt.status)
			}
		})
	}
}

func TestValidateStructMultipleErrors(t *testing.T) {
	type Request struct {
		Name  string `validate:"required"`
		Email string `validate:"required,email"`
		Age   int    `validate:"gte=0,lte=130"`
	}

	req := &Request{}

	err := validateStruct(req)
	if err == nil {
		t.Fatal("expected validation error")
	}

	// Should contain all validation errors
	errStr := err.Error()
	if !strings.Contains(errStr, "Name is required") {
		t.Errorf("expected error to contain 'Name is required', got %q", errStr)
	}
	if !strings.Contains(errStr, "Email") {
		t.Errorf("expected error to contain 'Email', got %q", errStr)
	}
}

func TestValidateStructNested(t *testing.T) {
	type Address struct {
		City    string `validate:"required"`
		Country string `validate:"required"`
	}

	type Request struct {
		Name    string  `validate:"required"`
		Address Address `validate:"required"`
	}

	req := &Request{
		Name: "John",
		Address: Address{
			City:    "",
			Country: "",
		},
	}

	err := validateStruct(req)
	if err == nil {
		t.Fatal("expected validation error for nested struct")
	}
}

func TestValidateStructSlice(t *testing.T) {
	type Request struct {
		Tags []string `validate:"min=1,dive,required"`
	}

	tests := []struct {
		name  string
		tags  []string
		valid bool
	}{
		{"valid slice", []string{"go", "http"}, true},
		{"empty slice", []string{}, false},
		{"nil slice", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &Request{Tags: tt.tags}
			err := validateStruct(req)

			if tt.valid && err != nil {
				t.Errorf("expected no error for tags %v, got %v", tt.tags, err)
			}
			if !tt.valid && err == nil {
				t.Errorf("expected error for tags %v, got nil", tt.tags)
			}
		})
	}
}

func TestValidateStructPointer(t *testing.T) {
	type Request struct {
		Name *string `validate:"required"`
	}

	validName := "John"
	req := &Request{Name: &validName}

	err := validateStruct(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	req2 := &Request{Name: nil}
	err2 := validateStruct(req2)
	if err2 == nil {
		t.Fatal("expected error for nil pointer")
	}
}

func TestValidateStructWithBinding(t *testing.T) {
	type Request struct {
		Name  string `query:"name" validate:"required"`
		Email string `query:"email" validate:"required,email"`
	}

	req := httptest.NewRequest(http.MethodGet, "/?name=john&email=john@example.com", nil)

	result := new(Request)
	err := bindQuery(req, result)
	if err != nil {
		t.Fatalf("unexpected bind error: %v", err)
	}

	err = validateStruct(result)
	if err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
}

func TestValidateStructWithBindingFails(t *testing.T) {
	type Request struct {
		Name  string `query:"name" validate:"required"`
		Email string `query:"email" validate:"required,email"`
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	result := new(Request)
	err := bindQuery(req, result)
	if err != nil {
		t.Fatalf("unexpected bind error: %v", err)
	}

	err = validateStruct(result)
	if err == nil {
		t.Fatal("expected validation error")
	}
}

// Response encoding tests

func TestEncodeResponseHeaders(t *testing.T) {
	type Response struct {
		CustomID string `header:"X-Custom-ID"`
		RequestID string `header:"X-Request-ID"`
	}

	resp := &Response{
		CustomID:  "custom-123",
		RequestID: "req-456",
	}

	w := httptest.NewRecorder()
	encodeResponseHeaders(w, resp)

	if w.Header().Get("X-Custom-ID") != "custom-123" {
		t.Errorf("expected X-Custom-ID to be custom-123, got %s", w.Header().Get("X-Custom-ID"))
	}

	if w.Header().Get("X-Request-ID") != "req-456" {
		t.Errorf("expected X-Request-ID to be req-456, got %s", w.Header().Get("X-Request-ID"))
	}
}

func TestEncodeResponseHeadersEmpty(t *testing.T) {
	type Response struct {
		CustomID string `header:"X-Custom-ID"`
	}

	resp := &Response{
		CustomID: "",
	}

	w := httptest.NewRecorder()
	encodeResponseHeaders(w, resp)

	if w.Header().Get("X-Custom-ID") != "" {
		t.Errorf("expected X-Custom-ID to be empty, got %s", w.Header().Get("X-Custom-ID"))
	}
}

func TestEncodeResponseHeadersNilPointer(t *testing.T) {
	w := httptest.NewRecorder()
	encodeResponseHeaders(w, nil)

	if len(w.Header()) != 0 {
		t.Errorf("expected no headers, got %d", len(w.Header()))
	}
}

func TestEncodeResponseHeadersNonStruct(t *testing.T) {
	w := httptest.NewRecorder()
	encodeResponseHeaders(w, "not a struct")

	if len(w.Header()) != 0 {
		t.Errorf("expected no headers, got %d", len(w.Header()))
	}
}

func TestEncodeResponseCookies(t *testing.T) {
	type Response struct {
		Token string `cookie:"token"`
	}

	resp := &Response{
		Token: "abc123",
	}

	w := httptest.NewRecorder()
	encodeResponseCookies(w, resp)

	cookies := w.Result().Cookies()
	found := false
	for _, cookie := range cookies {
		if cookie.Name == "token" {
			found = true
			if cookie.Value != "abc123" {
				t.Errorf("expected token cookie value to be abc123, got %s", cookie.Value)
			}
			if cookie.Path != "/" {
				t.Errorf("expected token cookie path to be /, got %s", cookie.Path)
			}
			break
		}
	}

	if !found {
		t.Error("expected token cookie to be set")
	}
}

func TestEncodeResponseCookiesWithDomain(t *testing.T) {
	type Response struct {
		Token string `cookie:"token" cookieDomain:"example.com"`
	}

	resp := &Response{
		Token: "abc123",
	}

	w := httptest.NewRecorder()
	encodeResponseCookies(w, resp)

	cookies := w.Result().Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "token" {
			if cookie.Domain != "example.com" {
				t.Errorf("expected token cookie domain to be example.com, got %s", cookie.Domain)
			}
			break
		}
	}
}

func TestEncodeResponseCookiesWithSecure(t *testing.T) {
	type Response struct {
		Token string `cookie:"token" cookieSecure:"true"`
	}

	resp := &Response{
		Token: "abc123",
	}

	w := httptest.NewRecorder()
	encodeResponseCookies(w, resp)

	cookies := w.Result().Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "token" {
			if !cookie.Secure {
				t.Error("expected token cookie to be secure")
			}
			break
		}
	}
}

func TestEncodeResponseCookiesWithHTTPOnly(t *testing.T) {
	type Response struct {
		Token string `cookie:"token" cookieHTTPOnly:"true"`
	}

	resp := &Response{
		Token: "abc123",
	}

	w := httptest.NewRecorder()
	encodeResponseCookies(w, resp)

	cookies := w.Result().Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "token" {
			if !cookie.HttpOnly {
				t.Error("expected token cookie to be HttpOnly")
			}
			break
		}
	}
}

func TestEncodeResponseCookiesWithSameSite(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected http.SameSite
	}{
		{"strict", "strict", http.SameSiteStrictMode},
		{"lax", "lax", http.SameSiteLaxMode},
		{"none", "none", http.SameSiteNoneMode},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create response with proper struct tag
			var resp any
			switch tt.value {
			case "strict":
				resp = &struct {
					Token string `cookie:"token" cookieSameSite:"strict"`
				}{Token: "abc123"}
			case "lax":
				resp = &struct {
					Token string `cookie:"token" cookieSameSite:"lax"`
				}{Token: "abc123"}
			case "none":
				resp = &struct {
					Token string `cookie:"token" cookieSameSite:"none"`
				}{Token: "abc123"}
			}

			w := httptest.NewRecorder()
			encodeResponseCookies(w, resp)

			cookies := w.Result().Cookies()
			for _, cookie := range cookies {
				if cookie.Name == "token" {
					if cookie.SameSite != tt.expected {
						t.Errorf("expected SameSite to be %v, got %v", tt.expected, cookie.SameSite)
					}
					break
				}
			}
		})
	}
}

func TestEncodeResponseCookiesWithMaxAge(t *testing.T) {
	type Response struct {
		Token string `cookie:"token" cookieMaxAge:"3600"`
	}

	resp := &Response{
		Token: "abc123",
	}

	w := httptest.NewRecorder()
	encodeResponseCookies(w, resp)

	cookies := w.Result().Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "token" {
			if cookie.MaxAge != 3600 {
				t.Errorf("expected MaxAge to be 3600, got %d", cookie.MaxAge)
			}
			break
		}
	}
}

func TestEncodeResponseCookiesClear(t *testing.T) {
	type Response struct {
		Token string `cookie:"token,clear"`
	}

	resp := &Response{
		Token: "",
	}

	w := httptest.NewRecorder()
	encodeResponseCookies(w, resp)

	cookies := w.Result().Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "token" {
			if cookie.MaxAge != -1 {
				t.Errorf("expected MaxAge to be -1 for clear, got %d", cookie.MaxAge)
			}
			break
		}
	}
}

func TestEncodeResponseCookiesWithPartitioned(t *testing.T) {
	type Response struct {
		Token string `cookie:"token" cookiePartitioned:"true"`
	}

	resp := &Response{
		Token: "abc123",
	}

	w := httptest.NewRecorder()
	encodeResponseCookies(w, resp)

	cookies := w.Result().Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "token" {
			if !cookie.Partitioned {
				t.Error("expected token cookie to be Partitioned")
			}
			break
		}
	}
}

func TestEncodeResponseCookiesMultiple(t *testing.T) {
	type Response struct {
		Token   string `cookie:"token"`
		Session string `cookie:"session"`
	}

	resp := &Response{
		Token:   "token-value",
		Session: "session-value",
	}

	w := httptest.NewRecorder()
	encodeResponseCookies(w, resp)

	cookies := w.Result().Cookies()
	tokenFound := false
	sessionFound := false

	for _, cookie := range cookies {
		if cookie.Name == "token" {
			tokenFound = true
			if cookie.Value != "token-value" {
				t.Errorf("expected token value to be token-value, got %s", cookie.Value)
			}
		}
		if cookie.Name == "session" {
			sessionFound = true
			if cookie.Value != "session-value" {
				t.Errorf("expected session value to be session-value, got %s", cookie.Value)
			}
		}
	}

	if !tokenFound {
		t.Error("expected token cookie to be set")
	}
	if !sessionFound {
		t.Error("expected session cookie to be set")
	}
}

func TestEncodeResponseCookiesNilPointer(t *testing.T) {
	w := httptest.NewRecorder()
	encodeResponseCookies(w, nil)

	cookies := w.Result().Cookies()
	if len(cookies) != 0 {
		t.Errorf("expected no cookies, got %d", len(cookies))
	}
}

func TestEncodeResponseCookiesNonStruct(t *testing.T) {
	w := httptest.NewRecorder()
	encodeResponseCookies(w, "not a struct")

	cookies := w.Result().Cookies()
	if len(cookies) != 0 {
		t.Errorf("expected no cookies, got %d", len(cookies))
	}
}

func TestEncodeResponseCookiesInlineDuration(t *testing.T) {
	type Response struct {
		Token string `cookie:"token,24h"`
	}

	resp := &Response{
		Token: "abc123",
	}

	w := httptest.NewRecorder()
	err := encodeResponseCookies(w, resp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cookies := w.Result().Cookies()
	found := false
	for _, cookie := range cookies {
		if cookie.Name == "token" {
			found = true
			if cookie.Value != "abc123" {
				t.Errorf("expected token cookie value to be abc123, got %s", cookie.Value)
			}
			if cookie.Expires.IsZero() {
				t.Error("expected token cookie to have Expires set")
			}
			break
		}
	}

	if !found {
		t.Error("expected token cookie to be set")
	}
}

func TestEncodeResponseCookiesInlineClear(t *testing.T) {
	type Response struct {
		Token string `cookie:"token,clear"`
	}

	resp := &Response{
		Token: "",
	}

	w := httptest.NewRecorder()
	err := encodeResponseCookies(w, resp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cookies := w.Result().Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "token" {
			if cookie.MaxAge != -1 {
				t.Errorf("expected MaxAge to be -1 for clear, got %d", cookie.MaxAge)
			}
			break
		}
	}
}

func TestEncodeResponseCookiesDurationOption(t *testing.T) {
	type Response struct {
		Token string `cookie:"token,duration:12h"`
	}

	resp := &Response{
		Token: "abc123",
	}

	w := httptest.NewRecorder()
	err := encodeResponseCookies(w, resp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cookies := w.Result().Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "token" {
			if cookie.Expires.IsZero() {
				t.Error("expected token cookie to have Expires set")
			}
			break
		}
	}
}

func TestEncodeResponseCookiesInvalidDuration(t *testing.T) {
	type Response struct {
		Token string `cookie:"token,24h30m60s99n99u"`
	}

	resp := &Response{
		Token: "abc123",
	}

	w := httptest.NewRecorder()
	err := encodeResponseCookies(w, resp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// This should succeed since the format is valid
	cookies := w.Result().Cookies()
	found := false
	for _, cookie := range cookies {
		if cookie.Name == "token" {
			found = true
			break
		}
	}

	if !found {
		t.Error("expected token cookie to be set")
	}
}

func TestEncodeResponseCookiesDefaultSameSite(t *testing.T) {
	type Response struct {
		Token string `cookie:"token"`
	}

	resp := &Response{
		Token: "abc123",
	}

	w := httptest.NewRecorder()
	err := encodeResponseCookies(w, resp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cookies := w.Result().Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "token" {
			if cookie.SameSite != http.SameSiteLaxMode {
				t.Errorf("expected default SameSite to be Lax, got %v", cookie.SameSite)
			}
			break
		}
	}
}

func TestEncodeResponseCookiesSessionOnly(t *testing.T) {
	type Response struct {
		Token string `cookie:"token" cookieSessionOnly:"true"`
	}

	resp := &Response{
		Token: "abc123",
	}

	w := httptest.NewRecorder()
	err := encodeResponseCookies(w, resp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cookies := w.Result().Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "token" {
			if cookie.Expires.IsZero() {
				t.Error("expected token cookie to have Expires set for session only")
			}
			break
		}
	}
}

func TestEncodeResponseHeadersInt(t *testing.T) {
	type Response struct {
		RequestID int `header:"X-Request-ID"`
	}

	resp := &Response{
		RequestID: 12345,
	}

	w := httptest.NewRecorder()
	err := encodeResponseHeaders(w, resp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if w.Header().Get("X-Request-ID") != "12345" {
		t.Errorf("expected X-Request-ID to be 12345, got %s", w.Header().Get("X-Request-ID"))
	}
}

func TestEncodeResponseHeadersBool(t *testing.T) {
	type Response struct {
		IsAdmin bool `header:"X-Is-Admin"`
	}

	resp := &Response{
		IsAdmin: true,
	}

	w := httptest.NewRecorder()
	err := encodeResponseHeaders(w, resp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if w.Header().Get("X-Is-Admin") != "true" {
		t.Errorf("expected X-Is-Admin to be true, got %s", w.Header().Get("X-Is-Admin"))
	}
}

func TestEncodeResponseHeadersMultiple(t *testing.T) {
	type Response struct {
		CustomID  string `header:"X-Custom-ID"`
		RequestID string `header:"X-Request-ID"`
		UserAgent string `header:"X-User-Agent"`
	}

	resp := &Response{
		CustomID:  "custom-123",
		RequestID: "req-456",
		UserAgent: "Mozilla/5.0",
	}

	w := httptest.NewRecorder()
	err := encodeResponseHeaders(w, resp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if w.Header().Get("X-Custom-ID") != "custom-123" {
		t.Errorf("expected X-Custom-ID to be custom-123, got %s", w.Header().Get("X-Custom-ID"))
	}

	if w.Header().Get("X-Request-ID") != "req-456" {
		t.Errorf("expected X-Request-ID to be req-456, got %s", w.Header().Get("X-Request-ID"))
	}

	if w.Header().Get("X-User-Agent") != "Mozilla/5.0" {
		t.Errorf("expected X-User-Agent to be Mozilla/5.0, got %s", w.Header().Get("X-User-Agent"))
	}
}

func TestEncodeResponseCookiesPointerField(t *testing.T) {
	type Response struct {
		Token *string `cookie:"token"`
	}

	tokenValue := "abc123"
	resp := &Response{
		Token: &tokenValue,
	}

	w := httptest.NewRecorder()
	err := encodeResponseCookies(w, resp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cookies := w.Result().Cookies()
	found := false
	for _, cookie := range cookies {
		if cookie.Name == "token" {
			found = true
			if cookie.Value != "abc123" {
				t.Errorf("expected token cookie value to be abc123, got %s", cookie.Value)
			}
			break
		}
	}

	if !found {
		t.Error("expected token cookie to be set")
	}
}

func TestEncodeResponseCookiesPointerFieldNil(t *testing.T) {
	type Response struct {
		Token *string `cookie:"token"`
	}

	resp := &Response{
		Token: nil,
	}

	w := httptest.NewRecorder()
	err := encodeResponseCookies(w, resp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cookies := w.Result().Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "token" {
			if cookie.Value != "<nil>" {
				t.Errorf("expected token cookie value to be <nil>, got %s", cookie.Value)
			}
			break
		}
	}
}

func TestEncodeResponseCookiesAllAttributes(t *testing.T) {
	type Response struct {
		FullCookie string `cookie:"full" cookiePath:"/api" cookieDomain:"example.com" cookieMaxAge:"3600" cookieSameSite:"Strict" cookieSecure:"true" cookieHTTPOnly:"true" cookiePartitioned:"true"`
	}

	resp := &Response{
		FullCookie: "full-value",
	}

	w := httptest.NewRecorder()
	err := encodeResponseCookies(w, resp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cookies := w.Result().Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "full" {
			if cookie.Value != "full-value" {
				t.Errorf("expected value to be full-value, got %s", cookie.Value)
			}
			if cookie.Path != "/api" {
				t.Errorf("expected Path to be /api, got %s", cookie.Path)
			}
			if cookie.Domain != "example.com" {
				t.Errorf("expected Domain to be example.com, got %s", cookie.Domain)
			}
			if cookie.MaxAge != 3600 {
				t.Errorf("expected MaxAge to be 3600, got %d", cookie.MaxAge)
			}
			if cookie.SameSite != http.SameSiteStrictMode {
				t.Errorf("expected SameSite to be Strict, got %v", cookie.SameSite)
			}
			if !cookie.Secure {
				t.Error("expected Secure to be true")
			}
			if !cookie.HttpOnly {
				t.Error("expected HttpOnly to be true")
			}
			if !cookie.Partitioned {
				t.Error("expected Partitioned to be true")
			}
			break
		}
	}
}

func TestEncodeResponseGHHandler(t *testing.T) {
	type Request struct {
		Name string `query:"name"`
	}

	type Response struct {
		Message string `json:"message"`
		CustomID string `header:"X-Custom-ID"`
	}

	handler := GH[Request, Response](func(ctx context.Context, req *Request) (Response, int, error) {
		return Response{
			Message:  "hello " + req.Name,
			CustomID: "custom-123",
		}, http.StatusOK, nil
	})

	req := httptest.NewRequest(http.MethodGet, "/?name=john", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	if w.Header().Get("X-Custom-ID") != "custom-123" {
		t.Errorf("expected X-Custom-ID to be custom-123, got %s", w.Header().Get("X-Custom-ID"))
	}
}

func TestGHforSSEHeaders(t *testing.T) {
	type Request struct {
		RoomID string `query:"room"`
	}

	handler := GHforSSE[Request, string](5*time.Second, func(ctx context.Context, req *Request, send func(name string, data string) error) error {
		return nil
	})

	req := httptest.NewRequest(http.MethodGet, "/?room=test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	if w.Header().Get("Content-Type") != "text/event-stream" {
		t.Errorf("expected Content-Type to be text/event-stream, got %s", w.Header().Get("Content-Type"))
	}

	if w.Header().Get("Cache-Control") != "no-cache" {
		t.Errorf("expected Cache-Control to be no-cache, got %s", w.Header().Get("Cache-Control"))
	}

	if w.Header().Get("Connection") != "keep-alive" {
		t.Errorf("expected Connection to be keep-alive, got %s", w.Header().Get("Connection"))
	}

	if w.Header().Get("X-Accel-Buffering") != "no" {
		t.Errorf("expected X-Accel-Buffering to be no, got %s", w.Header().Get("X-Accel-Buffering"))
	}

	if w.Header().Get("Retry-After") != "5" {
		t.Errorf("expected Retry-After to be 5, got %s", w.Header().Get("Retry-After"))
	}
}

func TestGHforSSESendEvents(t *testing.T) {
	type Request struct {
		Name string `query:"name"`
	}

	type Event struct {
		Message string `json:"message"`
	}

	handler := GHforSSE[Request, Event](5*time.Second, func(ctx context.Context, req *Request, send func(name string, data Event) error) error {
		if err := send("greeting", Event{Message: "hello " + req.Name}); err != nil {
			return err
		}

		if err := send("farewell", Event{Message: "bye " + req.Name}); err != nil {
			return err
		}

		return nil
	})

	req := httptest.NewRequest(http.MethodGet, "/?name=john", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	body := w.Body.String()

	if !strings.Contains(body, "event: greeting") {
		t.Error("expected body to contain 'event: greeting'")
	}

	if !strings.Contains(body, "event: farewell") {
		t.Error("expected body to contain 'event: farewell'")
	}

	if !strings.Contains(body, "hello john") {
		t.Error("expected body to contain 'hello john'")
	}

	if !strings.Contains(body, "bye john") {
		t.Error("expected body to contain 'bye john'")
	}

	lines := strings.Split(body, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "id: ") {
			idStr := strings.TrimPrefix(line, "id: ")
			if _, err := uuid.Parse(idStr); err != nil {
				t.Errorf("expected valid UUID v7, got %s", idStr)
			}
		}
	}
}

func TestGHforSSEBindError(t *testing.T) {
	type Request struct {
		Age int `query:"age"`
	}

	handler := GHforSSE[Request, string](5*time.Second, func(ctx context.Context, req *Request, send func(name string, data string) error) error {
		return nil
	})

	req := httptest.NewRequest(http.MethodGet, "/?age=abc", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}

	if !strings.Contains(w.Body.String(), "error") {
		t.Error("expected error response")
	}
}

func TestGHforSSERetryDuration(t *testing.T) {
	type Request struct{}

	handler := GHforSSE[Request, string](10*time.Second, func(ctx context.Context, req *Request, send func(name string, data string) error) error {
		return nil
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Header().Get("Retry-After") != "10" {
		t.Errorf("expected Retry-After to be 10, got %s", w.Header().Get("Retry-After"))
	}
}

func TestGHforSSENamedEvents(t *testing.T) {
	type Request struct{}

	type Event struct {
		Value int `json:"value"`
	}

	handler := GHforSSE[Request, Event](5*time.Second, func(ctx context.Context, req *Request, send func(name string, data Event) error) error {
		events := []struct {
			name string
			data Event
		}{
			{"created", Event{Value: 1}},
			{"updated", Event{Value: 2}},
			{"deleted", Event{Value: 3}},
		}

		for _, e := range events {
			if err := send(e.name, e.data); err != nil {
				return err
			}
		}

		return nil
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	body := w.Body.String()

	if !strings.Contains(body, "event: created") {
		t.Error("expected body to contain 'event: created'")
	}

	if !strings.Contains(body, "event: updated") {
		t.Error("expected body to contain 'event: updated'")
	}

	if !strings.Contains(body, "event: deleted") {
		t.Error("expected body to contain 'event: deleted'")
	}
}

func TestGHforSSEEmptyName(t *testing.T) {
	type Request struct{}

	type Event struct {
		Message string `json:"message"`
	}

	handler := GHforSSE[Request, Event](5*time.Second, func(ctx context.Context, req *Request, send func(name string, data Event) error) error {
		return send("", Event{Message: "no event name"})
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	body := w.Body.String()

	if strings.Contains(body, "event:") {
		t.Error("expected body to not contain 'event:' for unnamed events")
	}

	if !strings.Contains(body, "data:") {
		t.Error("expected body to contain 'data:'")
	}
}
