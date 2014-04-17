# Hawk http.Handler

```Go
import "github.com/kelseyhightower/hawkhandler"
```

## Example

```Go
package main

import (
	"crypto/sha1"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/kelseyhightower/hawkhandler"
	"github.com/tent/hawk-go"
)

var creds map[string]string

func CredentialsLookupFunc(c *hawk.Credentials) error {
	key, ok := creds[c.ID]
	if !ok {
		return errors.New("client ID does not exist: " + c.ID)
	}
	c.Key = key
	c.Hash = sha1.New
	return nil
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello\n")
}

func init() {
	creds = make(map[string]string)
	creds["1325003"] = "7CCFC015-3D8C-45F5-AF41-0505F29AA04A"
}

func main() {
	h := hawkhandler.HawkHandler(http.DefaultServeMux, nil, CredentialsLookupFunc)
	http.HandleFunc("/", HelloServer)
	err := http.ListenAndServe(":8080", h)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
```
