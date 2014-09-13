package main

import (
	"encoding/json"
	"fmt"
	"talks"
)

func main() {
	s := `{
"order_id": 6,
"unused": "mmm",
"buyer": {
  "name": "John Doe"
},
"products": [
  {"name": "foo", "price": 1.99},
  {"name": "bar", "price": 1.55}
]}`

	var msg talks.Message
	_ = json.Unmarshal([]byte(s), &msg)
	fmt.Printf("%+v\n", msg)
}
