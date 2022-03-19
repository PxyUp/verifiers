package verifiers

import (
	"context"
	"errors"
)

var (
	ErrCountMoreThanLength = errors.New("cant wait more than exists")
	ErrMaxAmountOfError    = errors.New("verifier reach max amount of error")
)

type Verifier func(ctx context.Context) error

type verifier struct {
	ctx context.Context
}

// New return new verifier with provided context
// If context is nil will be used context.Background()
func New(ctx context.Context) *verifier {
	if ctx == nil {
		ctx = context.Background()
	}
	return &verifier{
		ctx: ctx,
	}
}

// All verify all function finished without error in given context timeout/deadline
func (f *verifier) All(fns ...Verifier) error {
	return f.process(0, fns...)
}

// AtLeast verifies is at least provided amount of functions will be finished without error in given context timeout/deadline
func (f *verifier) AtLeast(count int, fns ...Verifier) error {
	if count > len(fns) {
		return ErrCountMoreThanLength
	}
	return f.process(len(fns)-count, fns...)
}

// OneOf verify at least one function finished without error in given context timeout/deadline
func (f *verifier) OneOf(fns ...Verifier) error {
	return f.process(len(fns)-1, fns...)
}

func (f *verifier) process(maxErrorCount int, fns ...Verifier) error {
	if len(fns) == 0 {
		return nil
	}
	var doneWithoutError, doneWithError int
	childrenCtx, cancel := context.WithCancel(f.ctx)
	defer cancel()
	resp := make(chan error)
	for _, fn := range fns {
		go func(verifier Verifier) {
			resp <- verifier(childrenCtx)
		}(fn)
	}
	for {
		select {
		case <-f.ctx.Done():
			return f.ctx.Err()
		case errInner, ok := <-resp:
			if !ok {
				return context.Canceled
			}
			if errInner == nil {
				doneWithoutError += 1
			} else {
				doneWithError += 1
			}
			if doneWithoutError == len(fns)-maxErrorCount {
				return nil
			}
			if doneWithError > maxErrorCount {
				return ErrMaxAmountOfError
			}
		}
	}
}
