# gin-dump

* Gin middleware/handler to dump header/body of request and response .

* Very helpful for debugging your applications.

* More beautiful output than httputil.DumpXXX()

## Content-type support / todo

* [x] application/json
* [x] application/x-www-form-urlencoded
* [ ] text/xml
* [ ] application/xml
* [ ] text/plain

## Usage
### Start using it

Download and install it:

```sh
$ go get github.com/tpkeeper/gin-dump
```

Import it in your code:

```go
import "github.com/tpkeeper/gin-dump"
```

### Canonical example:

```go
package main

import (
	"fmt"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/tpkeeper/gin-dump"
)

func main() {
    router := gin.Default()
    
    showReq := true
    showResp := true
    showBody := true
    showHeaders := false
    showCookies := false
    
    router.Use(gindump.Dump(nil))   // prints on stdout
    // or
	router.Use(gindump.Dump(func(dumpStr string) {
	    fmt.Println(dumpStr)
    }))
    // or
    router.Use(gindump.DumpWithOptions(showReq, showResp, showBody, showHeaders, showCookies, nil)   // prints on stdout
    // or
	router.Use(gindump.DumpWithOptions(showReq, showResp, showBody, showHeaders, showCookies, func(dumpStr string) {
	    fmt.Println(dumpStr)
    }))

	//...
	router.Run()
}
```

### Output is as follows

```sh
[GIN-dump]:
Request-Header:
{
    "Content-Type": [
        "application/x-www-form-urlencoded"
    ]
}
Request-Body:
{
    "bar": [
        "baz"
    ],
    "foo": [
        "bar",
        "bar2"
    ]
}
Response-Header:
{
    "Content-Type": [
        "application/json; charset=utf-8"
    ]
}
Response-Body:
{
    "data": {
        "addr": "tpkeeper@qq.com",
        "name": "jfise"
    },
    "ok": true
}
```
