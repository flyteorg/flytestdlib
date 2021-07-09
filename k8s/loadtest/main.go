package loadtest

import (
	"fmt"
	"sync"
	"time"

	"github.com/lyft/flytepropeller/pkg/controller/fairqueue"
	"k8s.io/client-go/util/workqueue"
)

type workCounterMap struct {
	*sync.Map
	cycles int

	l sync.Mutex
}

func (w *workCounterMap) allDone() bool {
	w.l.Lock()
	defer w.l.Unlock()
	allDone := true
	w.Map.Range(func(key, value interface{}) bool {
		v := value.(int)
		if v < w.cycles {
			allDone = false
			return false
		}
		return true
	})
	return allDone
}

func newWorkCounterMap(v map[string]time.Duration, cycles int) *workCounterMap {
	w := &workCounterMap{
		cycles: cycles,
		Map:    &sync.Map{},
	}
	for k := range v {
		w.Store(k, 0)
	}
	return w
}

func loadTest(tenantLoads map[string]time.Duration, q workqueue.Interface, cycles int) {

	w := newWorkCounterMap(tenantLoads, cycles)

	wg := sync.WaitGroup{}

	doWork := func() {
		defer wg.Done()
		for {
			i, ok := q.Get()
			if ok {
				return
			}
			d, _ := tenantLoads[i.(string)]
			time.Sleep(d)
			q.Done(i)
			v, _ := w.Load(i)
			iv := v.(int)
			w.Store(i, iv+1)
			if w.allDone() {
				q.ShutDown()
				return
			}
			if iv < cycles {
				q.Add(i)
			}
		}
	}

	wg.Add(3)
	for i := 0; i < 3; i++ {
		go doWork()
	}

	for k := range tenantLoads {
		q.Add(k)
	}

	wg.Wait()
}

func main() {

	rounds := 10
	cycles := 10
	tenantLoads := map[string]time.Duration{
		"ns1/item1": time.Millisecond * 200,
		"ns1/item2": time.Millisecond * 100,
		"ns1/item3": time.Millisecond * 50,
		"ns2/item1": time.Millisecond * 5,
		"ns3/item1": time.Millisecond * 100,
		"ns3/item2": time.Millisecond * 50,
		"ns4/item1": time.Millisecond * 5,
		"ns5/item1": time.Millisecond * 5,
	}

	avgTime := time.Duration(0)
	for i := 0; i < rounds; i++ {
		s := time.Now()
		loadTest(tenantLoads, fairqueue.New(), cycles)
		d := time.Now().Sub(s)
		fmt.Printf("FairQ, Asymmetric Load Round %d for %d cycles took: %f seconds\n", i, cycles, d.Seconds())
		avgTime += d
	}
	fairQTime := avgTime.Seconds() / float64(rounds)
	fmt.Printf("======> FairQ, Asymmetric Load for %d cycles took: %f seconds\n", cycles, fairQTime)

	avgTime = time.Duration(0)
	for i := 0; i < rounds; i++ {
		s := time.Now()
		loadTest(tenantLoads, workqueue.New(), cycles)
		d := time.Now().Sub(s)
		fmt.Printf("WorkQ, Asymmetric Load Round %d for %d cycles took: %f seconds\n", i, cycles, d.Seconds())
		avgTime += d
	}
	workQTime := avgTime.Seconds() / float64(rounds)
	fmt.Printf("======> WorkQ, Asymmetric Load for %d cycles took: %f seconds\n", cycles, workQTime)

	fmt.Printf("Comparison: FairQ: %f  WorkQ: %f. Speedup [%f]", fairQTime, workQTime, (workQTime-fairQTime)/workQTime*100.0)

}
