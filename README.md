[![GoDoc](http://godoc.org/github.com/omotto/workers?status.png)](http://godoc.org/github.com/omotto/workers)
[![Build Status](https://travis-ci.com/omotto/workers.svg?branch=main)](https://travis-ci.com/omotto/workers)
[![Coverage Status](https://coveralls.io/repos/github/omotto/workers/badge.svg)](https://coveralls.io/github/omotto/workers)

# basicCron

Package workers implements concurrent workers

### Installation

To download the specific tagged release, run:

```
go get github.com/omotto/workers
```

Import it in your program as:

```
import "github.com/omotto/workers"
```

### Usage

```
    // New pool of workers
    pool := New()

    // Define global context for running workers
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
	defer cancel()

    // Add workers
	if id, err = pool.AddWorker(true, func (string) int { return 0 }, "pepe"); err != nil {
		panic(err)
	}
    ids = append(ids, id)
    ...

    // Execute pool of workers
    if err = pool.Run(ctx); err == nil {
	    // Check results
        for _, id := range(ids) {
            // Get executed worker result
            if results, err := pool.GetResults(id); err == nil {
				r := results[0].Interface()
                // Check result type
				switch r.(type) {
					case int:
						value = r.(int)
					default:
						errors.New("invalid type")
				}
			}
        }
    }
