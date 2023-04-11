package messagedb

import (
	"encoding/json"
	"os"
)

func SqlDictionaryFromJson(jsonFile string) (dict map[string]string, err error) {
	if jsonFile == "" {
		return nil, nil
	}
	dict = make(map[string]string)

	if _, err = os.Stat(jsonFile); os.IsNotExist(err) {
		return nil, err
	}

	file, err := os.Open(jsonFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&dict)
	if err != nil {
		return nil, err
	}

	return dict, nil
}
