# gin-dump

Gin middleware/handler to dump header/body of request and response .

Very helpful for debugging your applications.

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
	"time"
	"github.com/gin-gonic/gin"
	"github.com/tpkeeper/gin-dump"
)

func main() {
	router := gin.Default()
	
	router.Use(gindump.Dump())
	//...
	router.Run()
}
```

### Output is as follows

```sh
[GIN-dump]:
Request-Header:
        {
                        Content-Length : [66]
                        Content-Type : [application/json;charset=utf-8]
                        Accept-Encoding : [gzip]
                        User-Agent : [Go-http-client/1.1]
        }

Request-Body:
        {
                        sms_code : 1111
                        telephone : 18322889845
                        password : lfajkdffsefadfare
        }

Response-Header:
        {
                        Content-Type : [application/json; charset=utf-8]
        }

Response-Body:
        {
                        data : sms_code error
                        ok : %!s(bool=false)
        }

```
