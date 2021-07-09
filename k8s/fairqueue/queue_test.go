package fairqueue

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"k8s.io/client-go/util/workqueue"
)

func TestBasic(t *testing.T) {
	// If something is seriously wrong this test will never complete.
	q := New()

	// Start producers
	const producers = 50
	producerWG := sync.WaitGroup{}
	producerWG.Add(producers)
	for i := 0; i < producers; i++ {
		go func(i int) {
			defer producerWG.Done()
			for j := 0; j < 50; j++ {
				q.Add(fmt.Sprintf("%s/%s", strconv.Itoa(i), strconv.Itoa(i)))
				time.Sleep(time.Millisecond)
			}
		}(i)
	}

	// Start consumers
	const consumers = 10
	consumerWG := sync.WaitGroup{}
	consumerWG.Add(consumers)
	for i := 0; i < consumers; i++ {
		go func(i int) {
			defer consumerWG.Done()
			for {
				item, quit := q.Get()
				if item == "added after shutdown!" {
					t.Errorf("Got an item added after shutdown.")
				}
				if quit {
					return
				}
				t.Logf("Worker %v: begin processing %v", i, item)
				time.Sleep(3 * time.Millisecond)
				t.Logf("Worker %v: done processing %v", i, item)
				q.Done(item)
			}
		}(i)
	}

	producerWG.Wait()
	q.ShutDown()
	q.Add("added after shutdown!")
	consumerWG.Wait()
}

func TestAddWhileProcessing(t *testing.T) {
	q := New()

	// Start producers
	const producers = 10
	producerWG := sync.WaitGroup{}
	producerWG.Add(producers)
	for i := 0; i < producers; i++ {
		go func(i int) {
			defer producerWG.Done()
			q.Add(fmt.Sprintf("%s/%s", strconv.Itoa(i), strconv.Itoa(i)))
		}(i)
	}

	// Start consumers
	const consumers = 2
	consumerWG := sync.WaitGroup{}
	consumerWG.Add(consumers)
	for i := 0; i < consumers; i++ {
		go func(i int) {
			defer consumerWG.Done()
			// Every worker will re-add every item up to two times.
			// This tests the dirty-while-processing case.
			counters := map[interface{}]int{}
			for {
				item, quit := q.Get()
				if quit {
					return
				}
				counters[item]++
				if counters[item] < 2 {
					q.Add(item)
				}
				q.Done(item)
			}
		}(i)
	}

	producerWG.Wait()
	q.ShutDown()
	consumerWG.Wait()
}

func TestLen(t *testing.T) {
	q := New()
	q.Add("foo")
	if e, a := 1, q.Len(); e != a {
		t.Errorf("Expected %v, got %v", e, a)
	}
	q.Add("bar")
	if e, a := 2, q.Len(); e != a {
		t.Errorf("Expected %v, got %v", e, a)
	}
	q.Add("foo") // should not increase the queue length.
	if e, a := 2, q.Len(); e != a {
		t.Errorf("Expected %v, got %v", e, a)
	}
}

func TestReinsert(t *testing.T) {
	q := New()
	q.Add("foo")

	// Start processing
	i, _ := q.Get()
	if i != "foo" {
		t.Errorf("Expected %v, got %v", "foo", i)
	}

	// Add it back while processing
	q.Add(i)

	// Finish it up
	q.Done(i)

	// It should be back on the queue
	i, _ = q.Get()
	if i != "foo" {
		t.Errorf("Expected %v, got %v", "foo", i)
	}

	// Finish that one up
	q.Done(i)

	if a := q.Len(); a != 0 {
		t.Errorf("Expected queue to be empty. Has %v items", a)
	}
}

func benchmarkQueueAddRepeat(b *testing.B, q workqueue.Interface) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		q.Add("foo")
	}
}

func BenchmarkPerNamespaceFairQueue_Add_Repeat(b *testing.B) { benchmarkQueueAddRepeat(b, New()) }
func BenchmarkWorkQueue_Add_Repeat(b *testing.B)             { benchmarkQueueAddRepeat(b, workqueue.New()) }

func benchmarkQueueAddDistinct(b *testing.B, q workqueue.Interface) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		q.Add(strconv.Itoa(i))
	}
}

func BenchmarkPerNamespaceFairQueue_Add_Distinct(b *testing.B) { benchmarkQueueAddDistinct(b, New()) }
func BenchmarkWorkQueue_Add_Distinct(b *testing.B)             { benchmarkQueueAddDistinct(b, workqueue.New()) }

func benchmarkQueueGet(b *testing.B, q workqueue.Interface) {
	b.StopTimer()
	for i := 0; i < b.N; i++ {
		q.Add(strconv.Itoa(i))
	}

	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, _ = q.Get()
	}
}

func BenchmarkPerNamespaceFairQueue_Get(b *testing.B) { benchmarkQueueGet(b, New()) }
func BenchmarkWorkQueue_Get(b *testing.B)             { benchmarkQueueGet(b, workqueue.New()) }
