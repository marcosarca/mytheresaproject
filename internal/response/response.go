package response

import (
	"encoding/json"
	"mytheresa/internal/apierror"
	"net/http"
)

func RespondWithData(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}

func RespondWithError(w http.ResponseWriter, err error) error {
	apierr, ok := err.(*apierror.ApiError)
	if ok {
		w.WriteHeader(apierr.Code())
		return json.NewEncoder(w).Encode(apierr)
	}

	w.WriteHeader(http.StatusInternalServerError)
	return json.NewEncoder(w).Encode(apierror.InternalServerError("Internal Server Error"))
}
