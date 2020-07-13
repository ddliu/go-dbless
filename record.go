// Copyright 2019 Liu Dong <ddliuhb@gmail.com>.
// Licensed under the MIT license.

package dbless

import (
	"encoding/json"
)

type Record map[string]interface{}

func (r Record) Unmarshal(input interface{}) error {
	data, err := json.Marshal(r)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, input)
}
