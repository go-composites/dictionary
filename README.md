<p align="center"><img src="https://raw.githubusercontent.com/go-composites/brand/main/social/go-composites.png" alt="go-composites/dictionary" width="720"></p>

# dictionary

A key→value composite for Composition-Oriented Programming — the sibling of
[`array`](https://github.com/go-composites/array). A `Dictionary` wraps a
`map[interface{}]interface{}` (arbitrary comparable keys, arbitrary values)
behind a small `Interface` that follows the go-composites grammar:

- **Never nil / Null-Object**: every constructor and method returns a real
  object; `Null()` provides an inert variant and `IsNull()` distinguishes it.
- **Result-based errors**: fallible lookups return a
  [`Result`](https://github.com/go-composites/result) — a payload on a hit, or a
  `Result` carrying `Error.New("key not found")` on a miss. No `(value, ok)`, no
  panics, no bare nils.
- **Composite returns**: presence tests return a
  [`Boolean`](https://github.com/go-composites/boolean); `Keys()` and `Values()`
  return an [`Array`](https://github.com/go-composites/array).

## Install

```sh
go get github.com/go-composites/dictionary@main
```

## Usage

```go
package main

import (
	"fmt"

	Dictionary "github.com/go-composites/dictionary/src"
	Result "github.com/go-composites/result/src"
)

func main() {
	d := Dictionary.New()
	d.Set("one", 1).Set("two", 2) // Set returns the Dictionary, so calls chain.

	fmt.Println(d.Len())                 // 2
	fmt.Println(d.Has("two").ToGoBool()) // true

	hit := d.Get("one")
	fmt.Println(hit.HasError(), hit.Payload()) // false 1

	miss := d.Get("nope")
	fmt.Println(miss.HasError(), miss.Error().Message()) // true key not found

	// Each short-circuits on the first Result whose HasError() is true.
	d.Each(func(key, value interface{}) Result.Interface {
		fmt.Printf("%v=%v\n", key, value)
		return Result.New()
	})

	// Keys() and Values() return an Array (order is unspecified — Go map order).
	_ = d.Keys()
	_ = d.Values()

	d.Delete("two")
}
```

### API

| Method | Returns | Notes |
| --- | --- | --- |
| `New(opts...)` | `Dictionary.Interface` | empty by default; `WithPairs` seeds it |
| `Null()` | `Dictionary.Interface` | inert Null-Object; `IsNull()` is `true` |
| `Set(key, value)` | `Dictionary.Interface` | returns the receiver (chainable) |
| `Get(key)` | `Result.Interface` | payload on hit; error `"key not found"` on miss |
| `Has(key)` | `Boolean.Interface` | presence test |
| `Delete(key)` | `Dictionary.Interface` | no-op when absent; chainable |
| `Len()` | `int` | number of pairs |
| `Keys()` | `Array.Interface` | keys, unspecified order |
| `Values()` | `Array.Interface` | values, unspecified order |
| `Each(fn)` | `Result.Interface` | iterate; short-circuit on `HasError()` |
| `IsNull()` | `bool` | `false` for a real Dictionary |

## License

BSD-3-Clause — see [LICENSE](LICENSE).
