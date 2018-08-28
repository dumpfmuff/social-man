package webutils

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func ParseJson(w http.ResponseWriter, r io.Reader, v interface{}) error {
	decoder := json.NewDecoder(r)
	err := decoder.Decode(v)
	if err != nil {
		logrus.Error("failed to decode")
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	return err
}
