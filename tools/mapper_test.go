package mapper

import (
	"fmt"
	"testing"
)

func TestXxx(t *testing.T) {
	type Inner struct {
		Code string `json:"code"`
	}

	type User struct {
		ID    int     `json:"id"`
		Name  string  `json:"name"`
		Score float64 `json:"score"`
		Meta  Inner   `json:"meta"`
	}

	u := User{ID: 1, Name: "Alice", Score: 98.5, Meta: Inner{Code: "X123"}}
	m := StructToMap(u)
	fmt.Println(m)
	/*
		map[
		  id:1
		  name:Alice
		  score:98.5
		  meta.code:X123
		]
	*/

	var u2 User
	MapToStruct(m, &u2)
	fmt.Println(u2)
}

func TestDeep(t *testing.T) {
	type Meta struct {
		Code string `json:"code"`
	}

	type User struct {
		ID    int     `json:"id"`
		Name  string  `json:"name"`
		Score float64 `json:"score"`
		Meta  Meta    `json:"meta"`
	}

	m := map[string]interface{}{
		"id":        1,
		"name":      "Alice",
		"score":     95.8,
		"meta.code": "X999",
	}

	var u User
	MapToStructNested(m, &u)

	fmt.Printf("%+v\n", u)
}
