package subscription

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jackc/pgx"
)

var handler http.Handler

func TestMain(m *testing.M) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		panic("$DATABASE_URL must be set")
	}

	cfg, err := pgx.ParseURI(dbURL)
	if err != nil {
		panic(err)
	}

	pool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:   cfg,
		AfterConnect: PrepareStatements,
	})
	if err != nil {
		panic(err)
	}

	_, err = pool.Exec("TRUNCATE table subscribers CASCADE;")
	if err != nil {
		panic(err)
	}
	handler = New(pool)

	os.Exit(m.Run())
}

func TestSubscribe(t *testing.T) {
	url := `?access_token=1ixLbltjWkzwqLMXT-8UF-UQeKRma0hOOWFA6o91oXw&subscriber_number=9171234567`
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusAccepted {
		t.Errorf("Code: %d, Body:%s", w.Code, w.Body.String())
	}
}

func TestMissingSubscriberNumber(t *testing.T) {
	url := `?access_token=1ixLbltjWkzwqLMXT-8UF-UQeKRma0hOOWFA6o91oXw`
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestMissingAccesToken(t *testing.T) {
	url := `?subscriber_number=9171234567`
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestNoQuery(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestPutNotAllowed(t *testing.T) {
	r, err := http.NewRequest(http.MethodPut, "", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fail()
	}
}

func TestDeleteNotAllowed(t *testing.T) {
	r, err := http.NewRequest(http.MethodDelete, "", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fail()
	}
}
