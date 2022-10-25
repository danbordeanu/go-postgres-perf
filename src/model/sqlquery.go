package model

import (
	"fmt"
)

type SqlUpload struct {
	Request
	Workers string `json:"workers" example:"4"`
	//File    string `json:"file" example:"@/Users/appleboy/test.zip"`
}

func (r *SqlUpload) Validate() error {
	//if r.File == "" {
	//	return fmt.Errorf("file required parameter")
	//}
	if r.Workers == "" {
		return fmt.Errorf("I need some workers")
	}
	return nil
}
