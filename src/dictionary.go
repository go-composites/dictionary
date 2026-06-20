package Dictionary

import (
	Array "github.com/go-composites/array/src"
	Boolean "github.com/go-composites/boolean/src"
	Error "github.com/go-composites/error/src"
	Result "github.com/go-composites/result/src"
)

// Interface is the public contract of a key→value composite — the sibling of
// Array. Keys are arbitrary comparable values; values are arbitrary. Fallible
// lookups return a Result rather than (value, ok) or a bare nil, and presence
// tests return a Boolean — both in keeping with the go-composites grammar.
type Interface interface {
	Set(key, value interface{}) Interface
	Get(key interface{}) Result.Interface
	Has(key interface{}) Boolean.Interface
	Delete(key interface{}) Interface
	Len() int
	Keys() Array.Interface
	Values() Array.Interface
	Each(fn func(key, value interface{}) Result.Interface) Result.Interface
	IsNull() bool
}

type data struct {
	value map[interface{}]interface{}
}

// Option is a functional option for New.
type Option func(*data)

// WithPairs seeds the Dictionary with an initial set of key→value pairs.
func WithPairs(pairs map[interface{}]interface{}) Option {
	return func(d *data) {
		for k, v := range pairs {
			d.value[k] = v
		}
	}
}

// New creates an empty Dictionary (unless seeded via options).
func New(options ...Option) Interface {
	d := &data{
		value: make(map[interface{}]interface{}),
	}
	for _, opt := range options {
		opt(d)
	}
	return d
}

// Set stores value under key and returns the receiver so calls chain.
func (d *data) Set(key, value interface{}) Interface {
	d.value[key] = value
	return d
}

// Get returns a Result carrying the value on a hit, or a Result whose error is
// Error.New("key not found") on a miss. It never returns nil and never panics.
func (d *data) Get(key interface{}) Result.Interface {
	if value, ok := d.value[key]; ok {
		return Result.New(
			Result.WithPayload(value),
		)
	}
	return Result.New(
		Result.WithError(Error.New("key not found")),
	)
}

// Has reports whether key is present, as a Boolean.Interface.
func (d *data) Has(key interface{}) Boolean.Interface {
	_, ok := d.value[key]
	return Boolean.New(ok)
}

// Delete removes key (a no-op when absent) and returns the receiver so calls
// chain.
func (d *data) Delete(key interface{}) Interface {
	delete(d.value, key)
	return d
}

// Len returns the number of stored pairs.
func (d *data) Len() int {
	return len(d.value)
}

// Keys returns an Array of the Dictionary's keys. Iteration order over a Go map
// is unspecified, so the order of the returned keys is unspecified too.
func (d *data) Keys() Array.Interface {
	keys := Array.New()
	for k := range d.value {
		keys.Push(k)
	}
	return keys
}

// Values returns an Array of the Dictionary's values, in the same unspecified
// order as Keys.
func (d *data) Values() Array.Interface {
	values := Array.New()
	for _, v := range d.value {
		values.Push(v)
	}
	return values
}

// Each iterates over the pairs, invoking fn for each. It short-circuits and
// returns the first Result for which HasError() is true; on a full pass it
// returns a fresh Result.New(). Iteration order is unspecified (Go map order).
func (d *data) Each(
	fn func(key, value interface{}) Result.Interface,
) Result.Interface {
	for key, value := range d.value {
		if result := fn(key, value); result.HasError() {
			return result
		}
	}
	return Result.New()
}

// IsNull reports that this is a real (non-null) Dictionary.
func (d *data) IsNull() bool {
	return false
}

// null is the Null-Object variant of a Dictionary: an empty, immutable
// placeholder that honours the full Interface without ever being nil. Mutating
// methods are no-ops that return the receiver; lookups always miss.
type null struct{}

// Null returns the Null-Object Dictionary.
func Null() Interface {
	return &null{}
}

func (n *null) Set(key, value interface{}) Interface { return n }

func (n *null) Get(key interface{}) Result.Interface {
	return Result.New(
		Result.WithError(Error.New("key not found")),
	)
}

func (n *null) Has(key interface{}) Boolean.Interface { return Boolean.False() }

func (n *null) Delete(key interface{}) Interface { return n }

func (n *null) Len() int { return 0 }

func (n *null) Keys() Array.Interface { return Array.New() }

func (n *null) Values() Array.Interface { return Array.New() }

func (n *null) Each(
	fn func(key, value interface{}) Result.Interface,
) Result.Interface {
	return Result.New()
}

// IsNull reports that this is the null Dictionary.
func (n *null) IsNull() bool { return true }
