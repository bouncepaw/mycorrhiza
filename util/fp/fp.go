// Package fp provides classic functional-programming-style functions.
package fp

// TODO: Use these fancy shining functions across old code.

func If[Branch any](cond bool, then Branch, els Branch) Branch {
	if cond {
		return then
	}
	return els
}

func FirstOrEmpty[T any](arr []T) T {
	if len(arr) == 0 {
		var t T
		return t
	}
	return arr[0]
}

func CoalesceOrEmpty[T comparable](ts ...T) T {
	var zero T
	for _, t := range ts {
		if t != zero {
			return t
		}
	}
	return zero
}
