package subscription

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSubscribe(t *testing.T) {
	subscribe := `{"access_token":"1ixLbltjWkzwqLMXT-8UF-UQeKRma0hOOWFA6o91oXw", "subscriber_number":"9171234567"}`
	consentHandler := Handler()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "", strings.NewReader(subscribe))
	consentHandler.ServeHTTP(w, req)
	expected := "9171234567 successfully subscribed.\n"
	actual := w.Body.String()
	if w.Code != http.StatusCreated || expected != actual {
		t.Errorf("Subscription fail. Code: %d\t Body: %s\n", w.Code, actual)
	}
}
