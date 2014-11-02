fmtckage main

import (
	"encoding/json"
	"fmt"
)

type Message struct {
	Name string `json:"tagged_name"`
	Body string `json:"tagged_name2"`
	Time int64  `json:"tagged_name3"`
}

func main() {
	m := Message{"Alice", "Hello", 1294706395881547000}
	fmt.Println(m)

	b, err := json.Marshal(m)
	if err != nil {
		fmt.Println("json error")
	}

	fmt.Println("json encoded as []byte:", b)
	fmt.Println("json encoded as []byte:", string(b))

	var m2 Message
	err = json.Unmarshal(b, &m2)
	fmt.Println("json decoded from []byte:", m2)
}

