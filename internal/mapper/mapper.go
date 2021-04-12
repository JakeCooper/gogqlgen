package mapper

import "encoding/json"

func ToMap(req interface{}) map[string]interface{} {
	return req.(map[string]interface{})
}

func ToStruct(mp map[string]interface{}, resp interface{}) error {
	b, err := json.Marshal(mp)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &resp)
}
