package standardgh

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	cookieDurationRegex = regexp.MustCompile(`^(.+),([0-9]+[hmsnu]+)$`)
	cookieClearRegex    = regexp.MustCompile(`^(.+),clear$`)
)

func encodeResponseHeaders(w http.ResponseWriter, resp any) error {
	respVal := reflectValue(resp)
	if !respVal.IsValid() || respVal.Kind() != reflect.Struct {
		return nil
	}

	t := respVal.Type()
	for i := range respVal.NumField() {
		field := t.Field(i)
		fieldValue := respVal.Field(i)

		tag := field.Tag.Get("header")
		if tag == "" {
			continue
		}

		name, _ := parseTag(tag)
		if name == "" {
			continue
		}

		if !fieldValue.IsValid() || !fieldValue.CanInterface() {
			continue
		}

		// Handle pointer fields
		if fieldValue.Kind() == reflect.Ptr {
			if fieldValue.IsNil() {
				continue
			}
			fieldValue = fieldValue.Elem()
		}

		value := fmt.Sprintf("%v", fieldValue.Interface())
		w.Header().Set(name, value)
	}

	return nil
}

func encodeResponseCookies(w http.ResponseWriter, resp any) error {
	respVal := reflectValue(resp)
	if !respVal.IsValid() || respVal.Kind() != reflect.Struct {
		return nil
	}

	t := respVal.Type()
	for i := range respVal.NumField() {
		field := t.Field(i)
		fieldValue := respVal.Field(i)

		tag := field.Tag.Get("cookie")
		if tag == "" {
			continue
		}

		cookie, err := parseCookieTag(tag, field, fieldValue)
		if err != nil {
			return err
		}

		if cookie != nil {
			http.SetCookie(w, cookie)
		}
	}

	return nil
}

func parseCookieTag(tag string, field reflect.StructField, fieldValue reflect.Value) (*http.Cookie, error) {
	if !fieldValue.IsValid() || !fieldValue.CanInterface() {
		return nil, nil
	}

	// Handle pointer fields
	if fieldValue.Kind() == reflect.Ptr {
		if fieldValue.IsNil() {
			return nil, nil
		}
		fieldValue = fieldValue.Elem()
	}

	value := fmt.Sprintf("%v", fieldValue.Interface())

	// Try to match "name,duration" format (e.g., "token,24h")
	if matches := cookieDurationRegex.FindStringSubmatch(tag); len(matches) == 3 {
		name := matches[1]
		durationStr := matches[2]

		d, err := time.ParseDuration(durationStr)
		if err != nil {
			return nil, fmt.Errorf("cookie duration parse error: %w", err)
		}

		cookie := &http.Cookie{
			Name:    name,
			Value:   value,
			Path:    "/",
			Expires: time.Now().Add(d),
		}

		applyCookieAttributes(cookie, field)
		return cookie, nil
	}

	// Try to match "name,clear" format
	if matches := cookieClearRegex.FindStringSubmatch(tag); len(matches) == 2 {
		name := matches[1]
		cookie := &http.Cookie{
			Name:   name,
			Value:  value,
			Path:   "/",
			Expires: time.Unix(0, 0),
			MaxAge: -1,
		}

		applyCookieAttributes(cookie, field)
		return cookie, nil
	}

	// Standard parseTag format: "name"
	name, opts := parseTag(tag)
	if name == "" {
		return nil, nil
	}

	cookie := &http.Cookie{
		Name:  name,
		Value: value,
		Path:  "/",
	}

	applyCookieAttributes(cookie, field)

	// Handle duration option (e.g., "duration:24h")
	if duration, ok := opts["duration"]; ok {
		if d, err := time.ParseDuration(duration); err == nil {
			cookie.Expires = time.Now().Add(d)
		}
	} else if _, ok := opts["clear"]; ok {
		cookie.Expires = time.Unix(0, 0)
		cookie.MaxAge = -1
	}

	return cookie, nil
}

func applyCookieAttributes(cookie *http.Cookie, field reflect.StructField) {
	// Default SameSite to Lax
	cookie.SameSite = http.SameSiteLaxMode

	// Read attributes from separate struct tags
	if path := field.Tag.Get("cookiePath"); path != "" {
		cookie.Path = path
	}

	if sameSite := field.Tag.Get("cookieSameSite"); sameSite != "" {
		switch strings.ToLower(sameSite) {
		case "strict":
			cookie.SameSite = http.SameSiteStrictMode
		case "lax":
			cookie.SameSite = http.SameSiteLaxMode
		case "none":
			cookie.SameSite = http.SameSiteNoneMode
		}
	}

	if secure := field.Tag.Get("cookieSecure"); secure != "" {
		cookie.Secure, _ = strconv.ParseBool(secure)
	}

	if httpOnly := field.Tag.Get("cookieHTTPOnly"); httpOnly != "" {
		cookie.HttpOnly, _ = strconv.ParseBool(httpOnly)
	}

	if domain := field.Tag.Get("cookieDomain"); domain != "" {
		cookie.Domain = domain
	}

	if maxAge := field.Tag.Get("cookieMaxAge"); maxAge != "" {
		cookie.MaxAge, _ = strconv.Atoi(maxAge)
	}

	if partitioned := field.Tag.Get("cookiePartitioned"); partitioned != "" {
		cookie.Partitioned, _ = strconv.ParseBool(partitioned)
	}

	if sessionOnly := field.Tag.Get("cookieSessionOnly"); sessionOnly == "true" {
		cookie.Expires = time.Unix(0, 0)
	}
}
