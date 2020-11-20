package workers

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestPoolError(t *testing.T) {
	var err error

	pool := New()

	if _, err = pool.AddWorker(true, func () { fmt.Println("Hello, world") }, 10); err == nil {
		t.Error("This AddFunc should return Error, wrong number of args")
	}

	if _, err = pool.AddWorker(true, nil); err == nil {
		t.Error("This AddFunc should return Error, fn is nil")
	}

	if _, err = pool.AddWorker(true, 0); err == nil {
		t.Error("This AddFunc should return Error, fn is not func kind")
	}

	if _, err = pool.AddWorker(true, func (s string, n int) {	fmt.Printf("We have params here, string `%s` and nymber %d\n", s, n) }, "s", 10, 12); err == nil {
		t.Error("This AddFunc should return Error, wrong number of args")
	}

	if _, err = pool.AddWorker(true, func (s string, n int) {	fmt.Printf("We have params here, string `%s` and nymber %d\n", s, n) }, "s", "s2"); err == nil {
		t.Error("This AddFunc should return Error, args are not the correct type")
	}

	if _, err = pool.AddWorker(true, func (s string, n int) {	fmt.Printf("We have params here, string `%s` and nymber %d\n", s, n) }, "s", "s2"); err == nil {
		t.Error("This AddFunc should return Error, syntax error")
	}

	// custom types and interfaces as function params
	type user struct {
		ID   int
		Name string
	}
	var u user
	if _, err = pool.AddWorker(true, func (u user) { fmt.Println("Custom type as param") }, u); err != nil {
		t.Error(err)
	}

	type Foo interface {
		Bar() string
	}
	if _, err = pool.AddWorker(true, func (i Foo) { i.Bar() }, u); err == nil {
		t.Error("This should return error, type that don't implements interface assigned as param")
	}

	if _, err = pool.AddWorker(true, func (value int) int { fmt.Println(value); time.Sleep(time.Second*6); return value }, 10); err == nil {
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(4*time.Second))
		defer cancel()
		if err = pool.Run(ctx); err == nil {
			t.Error("This should return error, exceeded context timeout")
		}
	} else {
		t.Error(err)
	}

	pool = New()
	var id1, id2 string
	if id1, err = pool.AddWorker(true, func (u user) user { fmt.Println(u.Name); return u }, user{ID:10, Name: "pepe"}); err != nil {
		t.Error(err)
		return
	}
	if id2, err = pool.AddWorker(true, func (name string, age int8) (string, error) { return fmt.Sprintf("%s  %d years old", name, age), errors.New("fake error") }, "pepe", int8(34)); err != nil {
		t.Error(err)
		return
	}
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(4*time.Second))
	defer cancel()
	if err = pool.Run(ctx); err == nil {
		if _, err = pool.GetResults("id"); err == nil {
			t.Error("This should return error, id do not exist")
		} else {
			if results, err := pool.GetResults(id1); err == nil {
				switch results[0].(type) {
				case user:
					t.Log(results[0].(user).Name)
					t.Log(results[0].(user).ID)
				default:
					t.Error("result should be user type")
				}
			} else {
				t.Error(err)
			}
			if results, err := pool.GetResults(id2); err == nil {
				switch results[0].(type) {
				case string:
					t.Log(results[0].(string))
				default:
					t.Error("result should be string type")
				}
				switch results[1].(type) {
				case error:
					t.Log(results[1].(error))
				default:
					t.Error("result should be error type")
				}
			} else {
				t.Error(err)
			}
		}
	}
}