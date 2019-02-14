package jsons

import (
	"encoding/json"
	"os"
)

func DecodeFromJsonFile(jsonFile string, v interface{}) error {
	r, err := os.Open(jsonFile)
	if err != nil {
		return err
	}
	dec := json.NewDecoder(r)
	return dec.Decode(v)
}

func EncodeToJsonFile(jsonFile string, v interface{}) error {
	w, err := os.Create(jsonFile)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(w)
	return enc.Encode(v)
}
