package jio

import "encoding/json"

func ValidateJSON(dataRaw *[]byte, schema Schema) error {
	jsonRaw := make(map[string]interface{})
	err := json.Unmarshal(*dataRaw, &jsonRaw)
	if err != nil {
		return err
	}
	jsonNew, err := schema.Validate("", jsonRaw)
	if err != nil {
		return err
	}
	dataNew, err := json.Marshal(jsonNew)
	if err != nil {
		return err
	}
	*dataRaw = dataNew
	return nil
}
