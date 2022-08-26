# errs-go
errs-go is simple error with stack trace and key-value attribute.

## Usage
### Create
new error instance
```go
func someFunction() error {
    err := errs.New("parameter error")
}
```
with value
```go
func someFunction() error {
    err := errs.New("parameter error").With("key", "value")
}
```
wrapping
```go
func someFunction() error {
    err := errs.Wrap(errors.New("original error"))
}
```
### Template
Declare the definition of error
```go
const (
    ErrNoSession    errs.StringError = "ErrNoSession"
    ErrRetryTimeout errs.StringError = "ErrRetryTimeout"
)

func someFunction1() error {
    session := getSession()
    if session != nil {
        return ErrNoSession.New()
    }
}

func someFunction2() error {
    conn, err := connect()
    if err != nil {
        return ErrRetryTimeout.Wrap(err)
    }
}
```
### StackTrace
Source:
```go
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

```
Output: %s
```
somthing error
```
Output: %+v
```
somthing error
values={"key":"hoge"}
main.something
        ./main.go:11
main.main
        ./main.go:15
runtime.main
        runtime/proc.go:250
runtime.goexit
        runtime/asm_amd64.s:1571

```
Output: JSON
```json
{
    "message": "somthing error",
    "values": {
        "key": "hoge"
    },
    "stackTrace": [
        {
            "func": "main.something",
            "file": "./main.go",
            "line": 11
        },
        {
            "func": "main.main",
            "file": "./main.go",
            "line": 15
        },
        {
            "func": "runtime.main",
            "file": "runtime/proc.go",
            "line": 250
        },
        {
            "func": "runtime.goexit",
            "file": "runtime/asm_amd64.s",
            "line": 1571
        }
    ]
}
```
