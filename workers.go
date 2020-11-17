package workers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"sync"
	"time"
)

type worker struct {
	uuid			string
	function 		interface{}
	funcParams		[]interface{}
	funcReturns		[]reflect.Value
	wait 			bool
}

// Pool Struct
type Pool struct {
	lastIndex		int64				// Last inserted job index
	running			bool				// Scheduler status running (true) or stopped
	workers			[]worker			// Slice of jobs to be executed
	wg 				sync.WaitGroup		// WaitGroup
}

// New Pool constructor
func New() *Pool {
	return &Pool{
		running: false,
		lastIndex: 0,
	}
}

// AddWorker adds a new worker on pool setting if it's necessary to wait for it
func (c *Pool) AddWorker(wait bool, function interface{}, funcParams ...interface{}) (uuid string, err error) {
	// Check input parameters
	// -- Check if param function contains a function
	if function == nil || reflect.ValueOf(function).Kind() != reflect.Func { return uuid, fmt.Errorf("invalid function parameter") }
	// -- Check number of params
	if len(funcParams) != reflect.TypeOf(function).NumIn() { return uuid, fmt.Errorf("number of function params and number of provided params doesn't match") }
	// -- Check input params type belongs to function params type
	for i := 0; i < reflect.TypeOf(function).NumIn(); i++ {
		functionParam := reflect.TypeOf(function).In(i)
		inputParam := reflect.TypeOf(funcParams[i])
		if functionParam != inputParam {
			if functionParam.Kind() != reflect.Interface { return uuid, fmt.Errorf(fmt.Sprintf("param[%d] must be be `%s` not `%s`", i, functionParam, inputParam)) }
			if !inputParam.Implements(functionParam) { return uuid, fmt.Errorf(fmt.Sprintf("param[%d] of type `%s` doesn't implement interface `%s`", i, functionParam, inputParam)) }
		}
	}
	// Add new job
	uuid = strconv.FormatInt(time.Now().UnixNano(), 16) + strconv.FormatInt(c.lastIndex, 16)
	c.lastIndex++
	newWorker := worker { uuid: uuid, function: function, funcParams: funcParams, wait: wait }
	if c.running == false {
		c.workers = append(c.workers, newWorker)
	}
	return uuid, err
}

// GetResults get worker result by its id
func (c *Pool) GetResults(uuid string) ([]reflect.Value, error) {
	for _, w := range c.workers {
		if w.uuid == uuid {
			return w.funcReturns, nil
		}
	}
	return nil, errors.New("uuid not found")
}

// Run start running pool of workers
func (c *Pool) Run(ctx context.Context) error {
	if c.running == false {
		c.running = true
		ch := make(chan struct{})
		go func() {
			for i := 0; i < len(c.workers); i++ {
				if c.workers[i].wait {
					c.wg.Add(1)
				}
				go c.execWorker(&c.workers[i])
			}
			c.wg.Wait()
			ch <- struct{}{}
			c.running = false
		}()
		select {
		case <-ch:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}

// Private Methods

func (c *Pool) execWorker(w *worker) {
	defer func() {
		if r := recover(); r != nil { log.Println("worker func ", w.uuid, " error ", r) }
	}()
	args := make([]reflect.Value, len(w.funcParams))
	for i, param := range w.funcParams {
		args[i] = reflect.ValueOf(param)
	}
	w.funcReturns = reflect.ValueOf(w.function).Call(args)
	if w.wait {
		c.wg.Done()
	}
}
