package utils

import (
	"encoding/json"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
)

// Convert strut (mostly) to json values
func JsonMap(val interface{}) map[string]interface{} {
	data, err := json.Marshal(val)
	if err != nil {
		log.Error("Unable to marshall json", err)
		return map[string]interface{}{}
	}
	var res interface{}
	if err = json.Unmarshal(data, &res); err != nil {
		log.Error("Unable to unmarshall json", err)
		return map[string]interface{}{}
	}
	jsonMap, ok := res.(map[string]interface{})
	if !ok {
		log.Error("Unable to unmarshall json")
		return map[string]interface{}{}
	}
	return jsonMap
}

func JsonString(val interface{}) string {
	res, err := json.Marshal(val)
	if err != nil {
		log.Error("Unable to marshall json", err)
		return ""
	}
	return string(res)
}

func JsonStringIndent(val interface{}) string {
	res, err := json.MarshalIndent(val, "", "  ")
	if err != nil {
		log.Error("Unable to marshall json with indent", err)
		return ""
	}
	return string(res)
}

func ReadJSON(filepath string, res interface{}) error {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &res)
}
