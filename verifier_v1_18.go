//go:build go1.18
// +build go1.18

package verifiers

import "context"

// Only for v1.18 +
// FromArray generate Verifier from static generic array
func FromArray[T any](arr []T, cmp func(context.Context, T) error) []Verifier {
	fns := make([]Verifier, len(arr))

	for index, _ := range arr {
		value := arr[index]
		fns[index] = func(ctx context.Context) error {
			return cmp(ctx, value)
		}
	}

	return fns
}
