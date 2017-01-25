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

*NOTE* : If you are using Proxy , set the IsUsingProxy to True for getting correct address in ratelimit function.