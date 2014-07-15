// Written in 2014 by Petar Maymounkov.
//
// It helps future understanding of past knowledge to save
// this notice, so peers of other times and backgrounds can
// see history clearly, unless you have a better idea.

package faculty

import (
	"sort"

	"github.com/gocircuit/escher/tree"
)

func (attendant *EyeReCognizer) ReCognize(sentence Sentence) {
	ch := make(chan struct{})
	for _, sf_ := range sentence {
		sf := sf_.YieldNil().(SentenceFunctional)
		go func() {
			attendant.recognize[sf.Valve()].ReCognize(sf.Value())
			ch <- struct{}{}
		}()
	}
	for _ = range sentence {
		<-ch
	}
}

func (attendant *EyeReCognizer) cognizeWith(valve tree.Name, value tree.Meaning) {
	attendant.Lock()
	attendant.age++
	attendant.memory.At(valve).Grow("Age", attendant.age).Grow("Value", value).Collapse()
	reply := attendant.formulate()
	attendant.Unlock()
	attendant.cognize(reply)
}

func (attendant *EyeReCognizer) formulate() Sentence {
	var sorting impressionStrength
	for valve, mf := range attendant.memory {
		sorting = append(sorting, MemoryFunctional(mf))
	}
	sort.Sort(sorting)
	return sorting.Verbalize()
}

type impressionStrength []MemoryFunctional

func (x impressionStrength) Verbalize() Sentence {
	s := make(Sentence)
	for i, mf := range x {
		s.Grow(i, tree.Make().Grow("Valve", mf.Valve()).Grow("Value", mf.Value()).Grow("Time", i))
	}
	return s
}

func (x impressionStrength) Len() int {
	return len(x)
}

func (x impressionStrength) Less(i, j int) bool {
	return x[i].Age() > x[j].Age()
}

func (x impressionStrength) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}