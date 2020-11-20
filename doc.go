/*

Package workers implements a pool of concurrent workers.

Installation

To download the specific tagged release, run:

	go get github.com/omotto/workers

Import it in your program as:

	import "github.com/omotto/workers"

Usage

	type user struct {
		ID   int
		Name string
	}

	var (
		err 	error
		id		string
		value	int
	)

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
	defer cancel()

	pool := New()

	if id, err = pool.AddWorker(true, func (u user) int { fmt.Println(u.Name); return u.ID }, user{ID:10, Name: "pepe"}); err == nil {
		if err = pool.Run(ctx); err == nil {
			if results, err := pool.GetResults(id); err == nil {
				switch results[0].(type) {
					case int:
						value = results[0].(int)
					default:
						errors.New("invalid type")
				}
			}
		}
	}
*/
package workers

