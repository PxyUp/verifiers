//go:build go1.18
// +build go1.18

package verifiers_test

import (
	"context"
	"errors"
	"github.com/PxyUp/verifiers"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFromArray(t *testing.T) {
	t.Run("All:", func(t *testing.T) {
		v := verifiers.New(context.Background())
		err := v.All(verifiers.FromArray([]int{1, 2, 3}, func(ctx context.Context, value int) error {
			if value > 2 {
				return nil
			}
			return errors.New("value less than or equal 2")
		})...)
		assert.Equal(t, verifiers.ErrMaxAmountOfError, err)
		err = v.All(verifiers.FromArray([]int{1, 2, 3}, func(ctx context.Context, value int) error {
			if value > 0 {
				return nil
			}
			return errors.New("value less than or equal 2")
		})...)
		assert.NoError(t, err)
	})
	t.Run("NoOne:", func(t *testing.T) {
		v := verifiers.New(context.Background())
		err := v.NoOne(verifiers.FromArray([]int{1, 2, 3}, func(ctx context.Context, value int) error {
			if value > 2 {
				return nil
			}
			return errors.New("value less than or equal 2")
		})...)
		assert.Equal(t, verifiers.ErrMaxAmountOfFinished, err)
		err = v.NoOne(verifiers.FromArray([]int{1, 2, 3}, func(ctx context.Context, value int) error {
			if value > 5 {
				return nil
			}
			return errors.New("value less than or equal 2")
		})...)
		assert.NoError(t, err)
	})
	t.Run("OnyOne:", func(t *testing.T) {
		v := verifiers.New(context.Background())
		err := v.OnlyOne(verifiers.FromArray([]int{1, 2, 3}, func(ctx context.Context, value int) error {
			if value > 1 {
				return nil
			}
			return errors.New("value less than or equal 2")
		})...)
		assert.Equal(t, verifiers.ErrMaxAmountOfFinished, err)
		err = v.OnlyOne(verifiers.FromArray([]int{1, 2, 3}, func(ctx context.Context, value int) error {
			if value > 2 {
				return nil
			}
			return errors.New("value less than or equal 2")
		})...)
		assert.NoError(t, err)
	})
	t.Run("Exact:", func(t *testing.T) {
		v := verifiers.New(context.Background())
		err := v.Exact(1, verifiers.FromArray([]int{1, 2, 3}, func(ctx context.Context, value int) error {
			if value > 1 {
				return nil
			}
			return errors.New("value less than or equal 2")
		})...)
		assert.Equal(t, verifiers.ErrMaxAmountOfFinished, err)
		err = v.Exact(1, verifiers.FromArray([]int{1, 2, 3}, func(ctx context.Context, value int) error {
			if value > 2 {
				return nil
			}
			return errors.New("value less than or equal 2")
		})...)
		assert.NoError(t, err)
	})
	t.Run("AtLeast:", func(t *testing.T) {
		v := verifiers.New(context.Background())
		err := v.AtLeast(1, verifiers.FromArray([]int{1, 2, 3}, func(ctx context.Context, value int) error {
			if value > 5 {
				return nil
			}
			return errors.New("value less than or equal 2")
		})...)
		assert.Equal(t, verifiers.ErrMaxAmountOfError, err)
		err = v.AtLeast(1, verifiers.FromArray([]int{1, 2, 3}, func(ctx context.Context, value int) error {
			if value > 2 {
				return nil
			}
			return errors.New("value less than or equal 2")
		})...)
		assert.NoError(t, err)
	})
	t.Run("Custom array test:", func(t *testing.T) {
		type People struct {
			Name string
			Age  int
		}

		people := []*People{
			{
				Age:  20,
				Name: "First",
			},
			{
				Age:  25,
				Name: "Second",
			},
			{
				Age:  30,
				Name: "Third",
			},
		}

		fns := verifiers.FromArray(people, func(ctx context.Context, p *People) error {
			if p.Age >= 25 {
				return nil
			}
			return errors.New("to old")
		})

		v := verifiers.New(context.Background())

		assert.Equal(t, verifiers.ErrMaxAmountOfError, v.All(fns...))
		assert.Equal(t, verifiers.ErrMaxAmountOfFinished, v.NoOne(fns...))
		assert.Equal(t, nil, v.OneOf(fns...))
		assert.Equal(t, nil, v.Exact(2, fns...))
		assert.Equal(t, verifiers.ErrMaxAmountOfFinished, v.OnlyOne(fns...))
	})
}
