package http

import (
	"MyLog-M/internal/domain"
	"MyLog-M/internal/service/mock"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestHomePage(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll pass 'nil' as the third parameter.
	request, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(Home)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(response, request)

	// Check the status code is what we expect.
	if status := response.Code; status != http.StatusOK {
		t.Errorf("wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `Welcome to MyLog-As-A-Service.`
	if response.Body.String() != expected {
		t.Errorf("unexpected body: got %v want %v", response.Body.String(), expected)
	}
}

func TestHealthCheck(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll pass 'nil' as the third parameter.
	request, err := http.NewRequest(http.MethodGet, "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthCheck)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(response, request)

	// Check the status code is what we expect.
	if status := response.Code; status != http.StatusOK {
		t.Errorf("wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `Ok`
	if response.Body.String() != expected {
		t.Errorf("unexpected body: got %v want %v", response.Body.String(), expected)
	}
}

func TestMyTail(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll pass 'nil' as the third parameter.
	request, err := http.NewRequest(http.MethodGet, "/api/tail", nil)
	if err != nil {
		t.Fatal(err)
	}
	// mock service
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mock.NewMockService(ctrl)
	m.EXPECT().Tail(int64(1)).Return(&[]domain.Data{{RID: 1}}, nil)

	// initialize handler
	h := New(m)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(h.MyTail)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(response, request)

	// Check the status code is what we expect.
	if status := response.Code; status != http.StatusOK {
		t.Errorf("wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `[{"RID":1,"UnixTime":0,"LocalTime":"0001-01-01T00:00:00Z","LogType":"","LogSeverity":0,"LogText":"","Status":{"Code":0,"Text":""}}]`
	if response.Body.String() != expected {
		t.Errorf("unexpected body: got %v want %v", response.Body.String(), expected)
	}
}
