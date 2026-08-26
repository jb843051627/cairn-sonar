package codec

import "encoding/json"

func Encode(value any) ([]byte, error)     { return json.Marshal(value) }
func Decode(data []byte, target any) error { return json.Unmarshal(data, target) }
