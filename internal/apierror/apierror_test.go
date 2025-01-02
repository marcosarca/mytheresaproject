package apierror_test

import (
	"mytheresa/internal/apierror"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBadRequest(t *testing.T) {
	err := apierror.BadRequest("test message")

	apierr, ok := err.(*apierror.ApiError)

	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, apierr.Code())
	assert.Equal(t, "test message", apierr.Message)
	assert.Equal(t, "test message", apierr.Error())
}

func TestNotFound(t *testing.T) {
	err := apierror.NotFound("test message")
	apierr, ok := err.(*apierror.ApiError)

	assert.True(t, ok)
	assert.Equal(t, http.StatusNotFound, apierr.Code())
	assert.Equal(t, "test message", apierr.Message)
	assert.Equal(t, "test message", apierr.Error())
}

func TestInternalServerError(t *testing.T) {
	err := apierror.InternalServerError("test message")

	apierr, ok := err.(*apierror.ApiError)

	assert.True(t, ok)
	assert.Equal(t, http.StatusInternalServerError, apierr.Code())
	assert.Equal(t, "test message", apierr.Message)
	assert.Equal(t, "test message", apierr.Error())
}
