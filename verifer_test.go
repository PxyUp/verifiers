package verifiers_test

import (
	"context"
	"errors"
	"fmt"
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
	assert.Error(t, verifiers.ErrCountMoreThanLength, v.Exact(2))
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
		startTime := time.Now()
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
		assert.True(t, time.Now().Sub(startTime) <= time.Second*1)
	})
	t.Run("Return: nil at least one without error", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		v := verifiers.New(ctx)
		startTime := time.Now()
		assert.Nil(t, v.OneOf(
			func(ctx context.Context) error {
				return nil
			},
			func(ctx context.Context) error {
				return someError
			},
			func(ctx context.Context) error {
				time.Sleep(time.Second * 3)
				return nil
			},
		))

		assert.True(t, time.Now().Sub(startTime) <= time.Second*1)
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

func TestVerifier_AtLeast(t *testing.T) {
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

func TestVerifier_OnlyOne(t *testing.T) {
	t.Run("Return: nil just one without error", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		v := verifiers.New(ctx)
		startTime := time.Now()
		assert.Nil(t, v.OnlyOne(
			func(ctx context.Context) error {
				time.Sleep(time.Second)
				return nil
			},
			func(ctx context.Context) error {
				time.Sleep(time.Second * 2)
				return errors.New("")
			},
			func(ctx context.Context) error {
				time.Sleep(time.Second * 3)
				return errors.New("")
			},
		))
		assert.True(t, time.Now().Sub(startTime) >= time.Second*3)
	})
	t.Run("Return: err - all with error", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		v := verifiers.New(ctx)
		startTime := time.Now()
		assert.Equal(t, verifiers.ErrMaxAmountOfError, v.OnlyOne(
			func(ctx context.Context) error {
				time.Sleep(time.Second)
				return someError
			},
			func(ctx context.Context) error {
				time.Sleep(time.Second * 2)
				return someError
			},
			func(ctx context.Context) error {
				time.Sleep(time.Second * 3)
				return someError
			},
		))
		assert.True(t, time.Now().Sub(startTime) > time.Second*3)
	})
	t.Run("Return: err - two without error", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		v := verifiers.New(ctx)
		startTime := time.Now()
		assert.Equal(t, verifiers.ErrMaxAmountOfFinished, v.OnlyOne(
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
		assert.True(t, time.Now().Sub(startTime) >= time.Second*2)
		assert.True(t, time.Now().Sub(startTime) < time.Second*3)
	})
	t.Run("Return: err context timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		v := verifiers.New(ctx)
		startTime := time.Now()
		assert.Equal(t, context.DeadlineExceeded, v.OnlyOne(
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
	t.Run("Return: nil - at one return no error + all other finsihed", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		executed := false
		finished := false
		v := verifiers.New(ctx)
		assert.NoError(t, v.OnlyOne(
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
				return someError
			},
		))
		time.Sleep(time.Second * 5)
		assert.True(t, executed)
		assert.True(t, finished)
	})
}

func TestVerifier_NoOne(t *testing.T) {
	t.Run("Return: nil all with error", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		v := verifiers.New(ctx)
		startTime := time.Now()
		assert.Nil(t, v.NoOne(
			func(ctx context.Context) error {
				time.Sleep(time.Second)
				return someError
			},
			func(ctx context.Context) error {
				time.Sleep(time.Second * 2)
				return someError
			},
			func(ctx context.Context) error {
				time.Sleep(time.Second * 3)
				return someError
			},
		))
		assert.True(t, time.Now().Sub(startTime) >= time.Second*3)
	})
	t.Run("Return: err - all without error", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		v := verifiers.New(ctx)
		startTime := time.Now()
		assert.Equal(t, verifiers.ErrMaxAmountOfFinished, v.NoOne(
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
		assert.True(t, time.Now().Sub(startTime) >= time.Second*1)
		assert.True(t, time.Now().Sub(startTime) < time.Second*2)
	})
	t.Run("Return: err - two without error", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		v := verifiers.New(ctx)
		startTime := time.Now()
		assert.Equal(t, verifiers.ErrMaxAmountOfFinished, v.NoOne(
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
		assert.True(t, time.Now().Sub(startTime) >= time.Second*1)
		assert.True(t, time.Now().Sub(startTime) < time.Second*2)
	})
	t.Run("Return: err context timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		v := verifiers.New(ctx)
		startTime := time.Now()
		assert.Equal(t, context.DeadlineExceeded, v.NoOne(
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
	t.Run("Return: nil - at one return no error + all other finsihed", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		executed := false
		finished := false
		v := verifiers.New(ctx)
		assert.Equal(t, verifiers.ErrMaxAmountOfFinished, v.NoOne(
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
				return someError
			},
		))
		time.Sleep(time.Second * 5)
		assert.True(t, executed)
		assert.False(t, finished)
	})
}

func TestVerifier_Exact(t *testing.T) {
	t.Run("Return: nil just one without error", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		v := verifiers.New(ctx)
		startTime := time.Now()
		assert.Nil(t, v.Exact(
			1,
			func(ctx context.Context) error {
				time.Sleep(time.Second)
				return nil
			},
			func(ctx context.Context) error {
				time.Sleep(time.Second * 2)
				return errors.New("")
			},
			func(ctx context.Context) error {
				time.Sleep(time.Second * 3)
				return errors.New("")
			},
		))
		assert.True(t, time.Now().Sub(startTime) >= time.Second*3)
	})
	t.Run("Return: nil - all with error", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		v := verifiers.New(ctx)
		startTime := time.Now()
		assert.NoError(t, v.Exact(
			0,
			func(ctx context.Context) error {
				time.Sleep(time.Second)
				return someError
			},
			func(ctx context.Context) error {
				time.Sleep(time.Second * 2)
				return someError
			},
			func(ctx context.Context) error {
				time.Sleep(time.Second * 3)
				return someError
			},
		))
		assert.True(t, time.Now().Sub(startTime) > time.Second*3)
	})
	t.Run("Return: nil - two without error", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		v := verifiers.New(ctx)
		startTime := time.Now()
		assert.NoError(t, v.Exact(
			2,
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
		assert.True(t, time.Now().Sub(startTime) >= time.Second*3)
	})
	t.Run("Return: err context timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		v := verifiers.New(ctx)
		startTime := time.Now()
		assert.Equal(t, context.DeadlineExceeded, v.Exact(
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
	t.Run("Return: nil - two return no error + all other finsihed", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		executed := false
		finished := false
		v := verifiers.New(ctx)
		assert.NoError(t, v.Exact(
			2,
			func(ctx context.Context) error {
				return nil
			},
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
				return someError
			},
		))
		time.Sleep(time.Second * 5)
		assert.True(t, executed)
		assert.True(t, finished)
	})
}

func TestWithErrorComparator(t *testing.T) {
	errIgnore := errors.New("ignore error")
	t.Run("Return: nil - all without errors because WithErrorComparator options", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		v := verifiers.New(ctx, verifiers.WithErrorComparator(func(err error) bool {
			if err == nil {
				return false
			}

			if errors.Is(err, errIgnore) {
				return false
			}

			return true
		}))
		assert.NoError(t, v.AtLeast(3,
			func(ctx context.Context) error {
				return errIgnore
			},
			func(ctx context.Context) error {
				return fmt.Errorf("new error %w", errIgnore)
			},
			func(ctx context.Context) error {
				return errIgnore
			},
		))
	})
}
