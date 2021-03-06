package stream_test

import (
	op "github.com/MercuryThePlanet/optional"
	"github.com/MercuryThePlanet/stream"
	"testing"
)

type Array struct {
	ts  op.Ts
	idx int
}

func NewArray(ts op.Ts) *Array {
	return &Array{ts, 0}
}

func (a *Array) TryAdvance(c stream.Consumer) bool {
	if a.idx >= len(a.ts) {
		return false
	}
	c(a.ts[a.idx])
	a.idx++
	return true
}

func (a *Array) ForEachRemaining(c stream.Consumer) {
	for a.TryAdvance(c) {
	}
}

func (a *Array) TrySplit() (s stream.Spliterator) {
	remain := len(a.ts) - a.idx
	if remain > 1 {
		var split int = (remain) / 2
		s = &Array{ts: a.ts[split:], idx: 0}
		a.ts = a.ts[:split]
		return
	}
	return nil
}

func TestStream(t *testing.T) {
	gen := func() *Array {
		return NewArray(makeRange(1, 10))
	}

	stream.Of(gen()).Map(func(t op.T) op.T {
		return t.(int) * 2
	}).Filter(func(t op.T) bool {
		return t.(int) > 1
	}).FindAny()

	matches := stream.Of(gen()).AnyMatch(func(t op.T) bool {
		return t.(int) > 1
	})
	println(matches)

	stream.Of(gen()).Limit(10).ForEach(func(t op.T) {
		println(t.(int))
	})

	i := 0
	supplier := func() op.T {
		i++
		return i
	}
	arr := stream.Generate(supplier).Skip(10).Limit(10).ToSlice()
	for _, v := range arr {
		println(v.(int))
	}

	stream.Iterate(1, func(t op.T) op.T {
		return t.(int) * 2
	}).Limit(63).ForEach(func(t op.T) {
		println(t.(int))
	})

	stream.Of(gen()).Peek(func(t op.T) {
		print("Found one: ")
		println(t.(int))
	}).Filter(func(t op.T) bool {
		return t.(int) > 5
	}).Map(func(t op.T) op.T {
		return t.(int) * 5
	}).Peek(func(t op.T) {
		print("Values: ")
		println(t.(int))
	}).Reduce(func(t1, t2 op.T) op.T {
		return t1.(int) + t2.(int)
	}).IfPresent(func(t op.T) {
		print("Reduced value: ")
		println(t.(int))
	})
	stream.Of(NewArray(op.Ts{"Hello", "World", "!"})).Reduce(func(t1, t2 op.T) op.T {
		return t1.(string) + " " + t2.(string)
	}).IfPresent(func(t op.T) {
		println(t.(string))
	})

	println("TESTING COLLECTOR")
	collector := stream.NewGeneralCollector()
	stream.Of(gen()).Collect(collector)
	for _, t := range collector.Get() {
		println(t.(int))
	}
}

func makeRange(min, max int) []interface{} {
	a := make([]interface{}, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}
