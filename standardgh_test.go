package standardgh

import (
	"net/http"
	"net/http/httptest"
	"testing"
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
