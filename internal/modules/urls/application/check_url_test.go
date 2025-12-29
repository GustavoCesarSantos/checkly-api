package application

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCheckUrl_Execute_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	sut := NewCheckUrl()
	result, err := sut.Execute(context.Background(), server.URL)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !result.IsSuccess {
		t.Errorf("expected IsSuccess=true, got false")
	}
	if result.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, result.StatusCode)
	}
}

func TestCheckUrl_Execute_StatusCodeError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()
	sut := NewCheckUrl()
	result, err := sut.Execute(context.Background(), server.URL)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.IsSuccess {
		t.Errorf("expected IsSuccess=false, got true")
	}
	if result.StatusCode != http.StatusInternalServerError {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusInternalServerError,
			result.StatusCode,
		)
	}
}

func TestCheckUrl_Execute_InvalidURL(t *testing.T) {
	sut := NewCheckUrl()
	invalidURL := "http://invalid-url"
	_, err := sut.Execute(context.Background(), invalidURL)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCheckUrl_Execute_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second) // maior que o timeout de 3s
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	sut := NewCheckUrl()
	_, err := sut.Execute(context.Background(), server.URL)
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
	var netErr net.Error
	if !errors.As(err, &netErr) && err.Error() != "context deadline exceeded" {
		t.Errorf("expected timeout error, got %v", err)
	}
}
