package verifiers_test

import (
	"context"
	"errors"
	"github.com/PxyUp/verifiers"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	someError = errors.New("some error")
)

func TestNew(t *testing.T) {
	v := verifiers.New(nil)
	assert.NoError(t, v.All())
	assert.Error(t, verifiers.ErrCountMoreThanLength, v.AtLeast(2))
}

func TestVerifier_All(t *testing.T) {
	t.Run("Return: nil", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		v := verifiers.New(ctx)
		startTime := time.Now()
		assert.Nil(t, v.All(
			func(ctx context.Context) error {
				time.Sleep(time.Second)
				return nil
			},
			func(ctx context.Context) error {
				time.Sleep(time.Second * 2)
				return nil
			},
			func(ctx context.Context) error {
				time.Sleep(time.Second * 3)
				return nil
			},
		))
		assert.True(t, time.Now().Sub(startTime) >= time.Second*3)
	})
	t.Run("Return: err context timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		v := verifiers.New(ctx)
		startTime := time.Now()
		assert.Equal(t, context.DeadlineExceeded, v.All(
			func(ctx context.Context) error {
				time.Sleep(time.Second * 4)
				return nil
			},
			func(ctx context.Context) error {
				time.Sleep(time.Second * 4)
				return nil
			},
			func(ctx context.Context) error {
				time.Sleep(time.Second * 4)
				return nil
			},
		))
		assert.True(t, time.Now().Sub(startTime) <= time.Second*2)
	})
	t.Run("Return: err - at least one return error", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		v := verifiers.New(ctx)
		assert.Equal(t, verifiers.ErrMaxAmountOfError, v.All(
			func(ctx context.Context) error {
				return someError
			},
			func(ctx context.Context) error {
				return nil
			},
			func(ctx context.Context) error {
				return nil
			},
		))
	})
	t.Run("Return: err - at least one return error all other stopped", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		executed := false
		finished := false
		v := verifiers.New(ctx)
		assert.Equal(t, verifiers.ErrMaxAmountOfError, v.All(
			func(ctx context.Context) error {
				return someError
			},
			func(ctx context.Context) error {
				executed = true
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(time.Second * 3):
					finished = true
				}
				return nil
			},
		))
		time.Sleep(time.Second * 5)
		assert.True(t, executed)
		assert.False(t, finished)
	})
}

func TestVerifier_OneOf(t *testing.T) {
	t.Run("Return: nil all without error", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		v := verifiers.New(ctx)
		startTime := time.Now()
		assert.Nil(t, v.OneOf(
			func(ctx context.Context) error {
				time.Sleep(time.Second)
				return nil
			},
			func(ctx context.Context) error {
				time.Sleep(time.Second * 2)
				return nil
			},
			func(ctx context.Context) error {
				time.Sleep(time.Second * 3)
				return nil
			},
		))
		assert.True(t, time.Now().Sub(startTime) < time.Second*2)
	})
	t.Run("Return: nil all without error", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		v := verifiers.New(ctx)
		startTime := time.Now()
		assert.Nil(t, v.OneOf(
			func(ctx context.Context) error {
				time.Sleep(time.Second)
				return someError
			},
			func(ctx context.Context) error {
				time.Sleep(time.Second * 2)
				return nil
			},
			func(ctx context.Context) error {
				time.Sleep(time.Second * 3)
				return someError
			},
		))
		assert.True(t, time.Now().Sub(startTime) >= time.Second*1)
		assert.True(t, time.Now().Sub(startTime) <= time.Second*3)
	})
	t.Run("Return: nil at least one without error", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		v := verifiers.New(ctx)
		assert.Nil(t, v.OneOf(
			func(ctx context.Context) error {
				return nil
			},
			func(ctx context.Context) error {
				return someError
			},
			func(ctx context.Context) error {
				return someError
			},
		))
	})
	t.Run("Return: err context timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		v := verifiers.New(ctx)
		startTime := time.Now()
		assert.Equal(t, context.DeadlineExceeded, v.OneOf(
			func(ctx context.Context) error {
				time.Sleep(time.Second * 4)
				return nil
			},
			func(ctx context.Context) error {
				time.Sleep(time.Second * 4)
				return nil
			},
			func(ctx context.Context) error {
				time.Sleep(time.Second * 4)
				return nil
			},
		))
		assert.True(t, time.Now().Sub(startTime) <= time.Second*2)
	})
	t.Run("Return: err - all with errors", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		v := verifiers.New(ctx)
		assert.Equal(t, verifiers.ErrMaxAmountOfError, v.OneOf(
			func(ctx context.Context) error {
				return someError
			},
			func(ctx context.Context) error {
				return someError
			},
			func(ctx context.Context) error {
				return someError
			},
		))
	})
	t.Run("Return: nil - at least one return no error + all other stopped", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		executed := false
		finished := false
		v := verifiers.New(ctx)
		assert.NoError(t, v.OneOf(
			func(ctx context.Context) error {
				return nil
			},
			func(ctx context.Context) error {
				executed = true
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(time.Second * 3):
					finished = true
				}
				return nil
			},
		))
		time.Sleep(time.Second * 5)
		assert.True(t, executed)
		assert.False(t, finished)
	})
}

func TestVerifier_Count(t *testing.T) {
	t.Run("Return: nil all without error", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		v := verifiers.New(ctx)
		startTime := time.Now()
		assert.Nil(t, v.AtLeast(3,
			func(ctx context.Context) error {
				time.Sleep(time.Second)
				return nil
			},
			func(ctx context.Context) error {
				time.Sleep(time.Second * 2)
				return nil
			},
			func(ctx context.Context) error {
				time.Sleep(time.Second * 3)
				return nil
			},
		))
		assert.True(t, time.Now().Sub(startTime) >= time.Second*3)
	})
	t.Run("Return: nil some without error", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		v := verifiers.New(ctx)
		startTime := time.Now()
		assert.Nil(t, v.AtLeast(2,
			func(ctx context.Context) error {
				time.Sleep(time.Second)
				return nil
			},
			func(ctx context.Context) error {
				time.Sleep(time.Second * 2)
				return nil
			},
			func(ctx context.Context) error {
				time.Sleep(time.Second * 3)
				return someError
			},
		))
		assert.True(t, time.Now().Sub(startTime) > time.Second*1)
		assert.True(t, time.Now().Sub(startTime) < time.Second*3)
	})
	t.Run("Return: err some without error", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		v := verifiers.New(ctx)
		assert.Error(t, v.AtLeast(2,
			func(ctx context.Context) error {
				return nil
			},
			func(ctx context.Context) error {
				return someError
			},
			func(ctx context.Context) error {
				return someError
			},
		))
	})
	t.Run("Return: nil some without error", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		v := verifiers.New(ctx)
		assert.NoError(t, v.AtLeast(2,
			func(ctx context.Context) error {
				return nil
			},
			func(ctx context.Context) error {
				return nil
			},
			func(ctx context.Context) error {
				return someError
			},
			func(ctx context.Context) error {
				return someError
			},
		))
	})
	t.Run("Return: err context timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		v := verifiers.New(ctx)
		startTime := time.Now()
		assert.Equal(t, context.DeadlineExceeded, v.AtLeast(
			2,
			func(ctx context.Context) error {
				time.Sleep(time.Second * 4)
				return nil
			},
			func(ctx context.Context) error {
				time.Sleep(time.Second * 4)
				return nil
			},
			func(ctx context.Context) error {
				time.Sleep(time.Second * 4)
				return nil
			},
		))
		assert.True(t, time.Now().Sub(startTime) <= time.Second*2)
	})
	t.Run("Return: nil - all without errors", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		v := verifiers.New(ctx)
		assert.NoError(t, v.AtLeast(3,
			func(ctx context.Context) error {
				return nil
			},
			func(ctx context.Context) error {
				return nil
			},
			func(ctx context.Context) error {
				return nil
			},
		))
	})
	t.Run("Return: err - all with errors but count is 0", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		v := verifiers.New(ctx)
		assert.NoError(t, v.AtLeast(0,
			func(ctx context.Context) error {
				return someError
			},
			func(ctx context.Context) error {
				return someError
			},
			func(ctx context.Context) error {
				return someError
			},
		))
	})
	t.Run("Return: nil - at least one return no error + all other stopped", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		executed := false
		finished := false
		v := verifiers.New(ctx)
		assert.NoError(t, v.AtLeast(
			1,
			func(ctx context.Context) error {
				return nil
			},
			func(ctx context.Context) error {
				executed = true
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(time.Second * 3):
					finished = true
				}
				return nil
			},
		))
		time.Sleep(time.Second * 5)
		assert.True(t, executed)
		assert.False(t, finished)
	})
}
