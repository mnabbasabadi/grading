package http

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func respond(w http.ResponseWriter, data interface{}, statusCode int) error {
	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}

	if data == nil {
		return fmt.Errorf("data expected for all statuses other than %s", http.StatusText(http.StatusNoContent))
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	if _, err := w.Write(bytes); err != nil {
		return err
	}

	return nil
}

func (s server) respondError(w http.ResponseWriter, err error, code int) {
	if err := respond(w, map[string]string{"error": err.Error()}, code); err != nil {
		s.logger.With(err).Error("error responding to request")

	}
}
