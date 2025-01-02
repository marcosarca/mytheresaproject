package response_test

import (
	"errors"
	"mytheresa/internal/apierror"
	"mytheresa/internal/response"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRespondWithData_OK(t *testing.T) {
	w := httptest.NewRecorder()
	data := struct {
		Data string `json:"data"`
	}{
		Data: "something to test",
	}
	err := response.RespondWithData(w, http.StatusOK, data)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{\"data\":\"something to test\"}\n", w.Body.String())
}

func TestRespondWithError_ApierrorSent(t *testing.T) {
	w := httptest.NewRecorder()
	apierr := apierror.BadRequest("test message")

	err := response.RespondWithError(w, apierr)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "{\"message\":\"test message\"}\n", w.Body.String())
}

func TestRespondWithError_GenericErrorSent(t *testing.T) {
	w := httptest.NewRecorder()
	genericErr := errors.New("generic error")

	err := response.RespondWithError(w, genericErr)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "{\"message\":\"Internal Server Error\"}\n", w.Body.String())
}
