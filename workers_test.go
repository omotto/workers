package workers

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestPoolError(t *testing.T) {
	pool := New()

	if _, err := pool.AddWorker(true, func () { fmt.Println("Hello, world") }, 10); err == nil {
		t.Error("This AddFunc should return Error, wrong number of args")
	}

	if _, err := pool.AddWorker(true, nil); err == nil {
		t.Error("This AddFunc should return Error, fn is nil")
	}

	if _, err := pool.AddWorker(true, 0); err == nil {
		t.Error("This AddFunc should return Error, fn is not func kind")
	}

	if _, err := pool.AddWorker(true, func (s string, n int) {	fmt.Printf("We have params here, string `%s` and nymber %d\n", s, n) }, "s", 10, 12); err == nil {
		t.Error("This AddFunc should return Error, wrong number of args")
	}

	if _, err := pool.AddWorker(true, func (s string, n int) {	fmt.Printf("We have params here, string `%s` and nymber %d\n", s, n) }, "s", "s2"); err == nil {
		t.Error("This AddFunc should return Error, args are not the correct type")
	}

	if _, err := pool.AddWorker(true, func (s string, n int) {	fmt.Printf("We have params here, string `%s` and nymber %d\n", s, n) }, "s", "s2"); err == nil {
		t.Error("This AddFunc should return Error, syntax error")
	}

	// custom types and interfaces as function params
	type user struct {
		ID   int
		Name string
	}
	var u user
	if _, err := pool.AddWorker(true, func (u user) { fmt.Println("Custom type as param") }, u); err != nil {
		t.Error(err)
	}

	type Foo interface {
		Bar() string
	}
	if _, err := pool.AddWorker(true, func (i Foo) { i.Bar() }, u); err == nil {
		t.Error("This should return error, type that don't implements interface assigned as param")
	}

	if _, err := pool.AddWorker(true, func (value int) int { fmt.Println(value); time.Sleep(time.Second*6); return value }, 10); err == nil {
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(4*time.Second))
		defer cancel()
		if err = pool.Run(ctx); err == nil {
			t.Error("This should return error, exceeded context timeout")
		}
	} else {
		t.Error(err)
	}

	pool = New()
	if id, err := pool.AddWorker(true, func (u user) user { fmt.Println(u.Name); return u }, user{ID:10, Name: "pepe"}); err == nil {
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(4*time.Second))
		defer cancel()
		if err = pool.Run(ctx); err == nil {
			if _, err = pool.GetResults("id"); err == nil {
				t.Error("This should return error, id do not exist")
			} else {
				if results, err := pool.GetResults(id); err == nil {
					r := results[0].Interface()
					switch r.(type) {
					case string:
						t.Error("result value is not string type")
					case int, int8, int16, int32, int64, uint, uint16, uint32, uint64:
						t.Error("result value is not int type")
					case user:
						t.Log(r.(user).Name)
						t.Log(r.(user).ID)
					default:
						t.Error("result should be solved before this point")
					}
				} else {
					t.Error(err)
				}
			}
		} else {
			t.Error(err)
		}
	} else {
		t.Error(err)
	}
}