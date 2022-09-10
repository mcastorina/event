# event

A Go package providing convenience struct and functions for goroutine
synchronization.


## Timeout

Execute a task and return an error on timeout.

```go
package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/mcastorina/event"
)

func main() {
	defer func(start time.Time) {
		if r := recover(); r != nil {
			fmt.Println(r)
			return
		}
		fmt.Printf("finished in %s\n", time.Since(start))
	}(time.Now())

	if err := event.Timeout(5*time.Second, func(ctx context.Context) {
		_, _ = http.Get("https://api.github.com/repos/mcastorina/event")
	}); err != nil {
		panic("timed out")
	}
}
```

```
finished in 277.432939ms
```

## Trigger

Have multiple goroutines wait for a single event to trigger. Once a `Trigger`
has fired, all calls to `Wait()` will return immediately.

```go
package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/mcastorina/event"
)

func main() {
	var t event.Trigger
	var wg sync.WaitGroup

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			fmt.Println("waiting for trigger")
			t.Wait()
			time.Sleep(time.Duration(i) * time.Second)
			fmt.Println(i)
		}(i)
	}

	time.Sleep(3 * time.Second)
	fmt.Println("triggering")
	t.Trigger()

	wg.Wait()
	fmt.Println("exiting")
}
```

```
waiting for trigger
waiting for trigger
waiting for trigger
waiting for trigger
triggering
0
1
2
3
exiting
```
