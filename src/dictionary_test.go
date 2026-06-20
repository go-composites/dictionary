package Dictionary_test

import (
	"sort"

	Dictionary "github.com/go-composites/dictionary/src"
	Error "github.com/go-composites/error/src"
	Result "github.com/go-composites/result/src"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

// errResult builds a Result that reports HasError() == true, so Each
// short-circuits on it (HasError() is !error.IsNull()).
func errResult() Result.Interface {
	return Result.New(Result.WithError(Error.New("sentinel")))
}

// stringKeys collects an Array of (string-typed) keys into a sorted Go slice so
// assertions are independent of Go's unspecified map-iteration order.
func stringKeys(arr interface {
	Each(func(int, interface{}) Result.Interface) Result.Interface
}) []string {
	out := []string{}
	arr.Each(func(_ int, item interface{}) Result.Interface {
		out = append(out, item.(string))
		return Result.New()
	})
	sort.Strings(out)
	return out
}

// intValues collects an Array of (int-typed) values into a sorted Go slice.
func intValues(arr interface {
	Each(func(int, interface{}) Result.Interface) Result.Interface
}) []int {
	out := []int{}
	arr.Each(func(_ int, item interface{}) Result.Interface {
		out = append(out, item.(int))
		return Result.New()
	})
	sort.Ints(out)
	return out
}

var _ = ginkgo.Describe("Dictionary", func() {
	ginkgo.Describe("New", func() {
		ginkgo.It("returns a non-nil, non-null, empty Dictionary", func() {
			d := Dictionary.New()
			gomega.Expect(d).NotTo(gomega.BeNil())
			gomega.Expect(d.IsNull()).To(gomega.BeFalse())
			gomega.Expect(d.Len()).To(gomega.Equal(0))
		})

		ginkgo.It("seeds pairs via WithPairs", func() {
			d := Dictionary.New(Dictionary.WithPairs(map[interface{}]interface{}{
				"a": 1,
				"b": 2,
			}))
			gomega.Expect(d.Len()).To(gomega.Equal(2))
			gomega.Expect(d.Get("a").Payload()).To(gomega.Equal(1))
			gomega.Expect(d.Get("b").Payload()).To(gomega.Equal(2))
		})
	})

	ginkgo.Describe("Set", func() {
		ginkgo.It("stores values and returns the receiver for chaining", func() {
			d := Dictionary.New()
			ret := d.Set("one", 1).Set("two", 2)
			gomega.Expect(ret).To(gomega.BeIdenticalTo(d))
			gomega.Expect(d.Len()).To(gomega.Equal(2))
		})

		ginkgo.It("overwrites an existing key", func() {
			d := Dictionary.New()
			d.Set("k", 1).Set("k", 9)
			gomega.Expect(d.Len()).To(gomega.Equal(1))
			gomega.Expect(d.Get("k").Payload()).To(gomega.Equal(9))
		})
	})

	ginkgo.Describe("Get", func() {
		ginkgo.It("returns a payload-bearing Result on a hit", func() {
			d := Dictionary.New().Set("k", 42)
			r := d.Get("k")
			gomega.Expect(r).NotTo(gomega.BeNil())
			gomega.Expect(r.HasError()).To(gomega.BeFalse())
			gomega.Expect(r.Payload()).To(gomega.Equal(42))
		})

		ginkgo.It("returns a Result with 'key not found' on a miss", func() {
			d := Dictionary.New()
			r := d.Get("absent")
			gomega.Expect(r).NotTo(gomega.BeNil())
			gomega.Expect(r.HasError()).To(gomega.BeTrue())
			gomega.Expect(r.Error().Message()).To(gomega.Equal("key not found"))
		})
	})

	ginkgo.Describe("Has", func() {
		ginkgo.It("is true for a present key", func() {
			d := Dictionary.New().Set("k", 1)
			gomega.Expect(d.Has("k").ToGoBool()).To(gomega.BeTrue())
		})

		ginkgo.It("is false for an absent key", func() {
			d := Dictionary.New()
			gomega.Expect(d.Has("k").ToGoBool()).To(gomega.BeFalse())
		})
	})

	ginkgo.Describe("Delete", func() {
		ginkgo.It("removes a key and returns the receiver", func() {
			d := Dictionary.New().Set("a", 1).Set("b", 2)
			ret := d.Delete("a")
			gomega.Expect(ret).To(gomega.BeIdenticalTo(d))
			gomega.Expect(d.Has("a").ToGoBool()).To(gomega.BeFalse())
			gomega.Expect(d.Len()).To(gomega.Equal(1))
		})

		ginkgo.It("is a no-op for an absent key", func() {
			d := Dictionary.New().Set("a", 1)
			d.Delete("missing")
			gomega.Expect(d.Len()).To(gomega.Equal(1))
		})
	})

	ginkgo.Describe("Keys and Values", func() {
		ginkgo.It("returns all keys and values, order-independent", func() {
			d := Dictionary.New().Set("a", 1).Set("b", 2).Set("c", 3)
			gomega.Expect(stringKeys(d.Keys())).To(
				gomega.Equal([]string{"a", "b", "c"}))
			gomega.Expect(intValues(d.Values())).To(
				gomega.Equal([]int{1, 2, 3}))
		})

		ginkgo.It("returns empty Arrays for an empty Dictionary", func() {
			d := Dictionary.New()
			gomega.Expect(stringKeys(d.Keys())).To(gomega.BeEmpty())
			gomega.Expect(intValues(d.Values())).To(gomega.BeEmpty())
		})
	})

	ginkgo.Describe("Each", func() {
		ginkgo.It("visits every pair and returns a clean Result", func() {
			d := Dictionary.New().Set("a", 1).Set("b", 2).Set("c", 3)
			count := 0
			res := d.Each(func(key, value interface{}) Result.Interface {
				count++
				return Result.New()
			})
			gomega.Expect(count).To(gomega.Equal(3))
			gomega.Expect(res).NotTo(gomega.BeNil())
			gomega.Expect(res.HasError()).To(gomega.BeFalse())
		})

		ginkgo.It("short-circuits on the first error Result", func() {
			d := Dictionary.New().Set("a", 1).Set("b", 2).Set("c", 3)
			count := 0
			res := d.Each(func(key, value interface{}) Result.Interface {
				count++
				return errResult()
			})
			gomega.Expect(count).To(gomega.Equal(1))
			gomega.Expect(res.HasError()).To(gomega.BeTrue())
		})
	})

	ginkgo.Describe("Null", func() {
		ginkgo.It("is a Null-Object: IsNull true and inert", func() {
			n := Dictionary.Null()
			gomega.Expect(n).NotTo(gomega.BeNil())
			gomega.Expect(n.IsNull()).To(gomega.BeTrue())
			gomega.Expect(n.Len()).To(gomega.Equal(0))

			// Mutators are no-ops that return the receiver.
			gomega.Expect(n.Set("k", 1)).To(gomega.BeIdenticalTo(n))
			gomega.Expect(n.Delete("k")).To(gomega.BeIdenticalTo(n))
			gomega.Expect(n.Len()).To(gomega.Equal(0))

			// Lookups always miss.
			gomega.Expect(n.Has("k").ToGoBool()).To(gomega.BeFalse())
			get := n.Get("k")
			gomega.Expect(get.HasError()).To(gomega.BeTrue())
			gomega.Expect(get.Error().Message()).To(gomega.Equal("key not found"))

			// Collections are empty.
			gomega.Expect(stringKeys(n.Keys())).To(gomega.BeEmpty())
			gomega.Expect(intValues(n.Values())).To(gomega.BeEmpty())

			// Each returns a clean Result without invoking fn.
			called := false
			res := n.Each(func(key, value interface{}) Result.Interface {
				called = true
				return errResult()
			})
			gomega.Expect(called).To(gomega.BeFalse())
			gomega.Expect(res.HasError()).To(gomega.BeFalse())
		})
	})
})
