package json

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Decode JSON to struct
func Decode(r *http.Request, val interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(val); err != nil {
		fmt.Print(err)
		return err
	}
	return nil
}

// H is a shortcut for map[string]interface{}
type H map[string]interface{}
