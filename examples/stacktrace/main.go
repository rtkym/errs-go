package main

import (
	"fmt"

	"github.com/goccy/go-json"

	"github.com/rtkym/errs-go"
)

func something(value string) error {
	return errs.New("somthing error").With("key", value)
}

func main() {
	if err := something("hoge"); err != nil {
		b, _ := json.Marshal(err)
		fmt.Printf("%s\n", err)
		fmt.Printf("%+v\n", err)
		fmt.Println(string(b))
	}
}
