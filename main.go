package main

import (
	"fmt"
	"sort"

	Dictionary "github.com/go-composites/dictionary/src"
	Result "github.com/go-composites/result/src"
)

func main() {
	d := Dictionary.New()
	d.Set("one", 1).Set("two", 2).Set("three", 3)

	fmt.Printf("Len = %d\n", d.Len())
	fmt.Printf("Has(\"two\") = %t\n", d.Has("two").ToGoBool())

	got := d.Get("two")
	fmt.Printf("Get(\"two\") payload = %v, hasError = %t\n",
		got.Payload(), got.HasError())

	miss := d.Get("missing")
	fmt.Printf("Get(\"missing\") hasError = %t, message = %q\n",
		miss.HasError(), miss.Error().Message())

	// Keys are unordered; sort the string keys for a stable demo print.
	keys := []string{}
	d.Keys().Each(func(_ int, k interface{}) Result.Interface {
		keys = append(keys, k.(string))
		return Result.New()
	})
	sort.Strings(keys)
	fmt.Printf("Keys (sorted) = %v\n", keys)

	d.Delete("two")
	fmt.Printf("after Delete(\"two\"), Len = %d\n", d.Len())
}
