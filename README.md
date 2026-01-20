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

## Broadcast

Have multiple goroutines wait for a specific and unique signal.

```go
package main

import (
	"fmt"

	"github.com/mcastorina/event"
)

func main() {
	b := event.NewBroadcast()

	for i := 1; i <= 3; i++ {
		go func(i int) {
			b.Wait("hello")
			fmt.Println("hello", i)
			b.Broadcast(fmt.Sprintf("first[%d] done", i))
		}(i)
		go func(i int) {
			b.Wait("world")
			fmt.Println("world", i)
			b.Broadcast(fmt.Sprintf("second[%d] done", i))
		}(i)
	}

	b.Broadcast("hello")
	b.Wait("first[1] done")
	b.Wait("first[2] done")
	b.Wait("first[3] done")

	b.Broadcast("world")
	b.Wait("second[1] done")
	b.Wait("second[2] done")
	b.Wait("second[3] done")

	fmt.Println("exiting")
}
```

```
hello 2
hello 3
hello 1
world 2
world 1
world 3
exiting
```
