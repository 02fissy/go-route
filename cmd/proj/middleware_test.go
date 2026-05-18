package main

import (
	"testing"
	"bytes"
	"net/http"
	"net/http/httptest"
	"io"

	"displaybox.fisayoai.net/internal/assert"
)

func TestCommonHeaders(t *testing.T){
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil{
		t.Fatal(err)
	}

	next :=http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		w.Write([]byte("OK"))
	})
		commonHeaders(next).ServeHTTP(rr, r)

		rs := rr.Result()
		defer rs.Body.Close()

		tests := []struct{
			headerKey string
			expectedValue string
		}{
			{"Content-Security-Policy", "default-src 'self'; style-src 'self' https://fonts.googleapis.com; font-src 'self' https://fonts.gstatic.com; img-src 'self' data:;"},
			{"Referrer-Policy", "origin-when-cross-origin"},
			{"X-Content-Type-Options", "nosniff"},
			{"X-Frame-Options", "DENY"},
			{"X-XSS-Protection", "0"},
			{"Server", "Go"},
		}
		for _, tt := range tests{
			t.Run("Header_"+tt.headerKey, func(t *testing.T){
				actualValue := rs.Header.Get(tt.headerKey)
				assert.Equal(t, actualValue, tt.expectedValue)
			})
		}
		assert.Equal(t, rs.StatusCode, http.StatusOK)

		body, err := io.ReadAll(rs.Body)
		if err != nil{
			t.Fatal(err)
		}
		assert.Equal(t, string(bytes.TrimSpace(body)), "OK")
}
