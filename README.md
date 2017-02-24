# rate-limit
Very simple rate limiter for HTTP requests

### How to use ?

```go
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/ahmdrz/rate-limit"
)

func main() {
	r := ratelimit.InitRateLimit(5, 5*time.Second, ratelimit.DefaultHandler)
	http.HandleFunc("/", r.RateLimit(mainHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}
```

Or you can use it as a handler 

```go
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ahmdrz/rate-limit"
)

func main() {
	mux := http.DefaultServeMux
	mux.HandleFunc("/", mainHandler)
	limiter := ratelimit.NewHandler(mux, 10, 1*time.Minute, ratelimit.DefaultHandler)
	server := &http.Server{
		Addr:    fmt.Sprintf("localhost:%s", "8080"),
		Handler: limiter,
	}
	server.ListenAndServe()
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}
```

*NOTE* : If you are using Proxy , set the `IsUsingProxy` to `True` for getting correct IP address in ratelimit function.

You can use `WhiteList` :

```go
   ratelimit.WhiteList.Add("/index")
   ratelimit.WhiteList.HasPrefix("/images")
```

And ratelimit will skip whitelist URIs.

**NOTE** WhiteList operation such as Add and etc is not thread safe.