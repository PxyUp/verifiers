package verifiers

import (
	"context"
	"errors"
)

var (
	// ErrCountMoreThanLength is configuration error, will return if we expect more function than we provide for verifier.AtLeast or verifier.Exact
	ErrCountMoreThanLength = errors.New("cant wait more than exists")
	// ErrMaxAmountOfError wii be returned some function which we not expect return error
	ErrMaxAmountOfError = errors.New("verifier reach max amount of error")
	// ErrMaxAmountOfFinished will be returned if some other function(which we not expect) return success
	ErrMaxAmountOfFinished = errors.New("verifier reach max amount success jobs")
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
	return f.Exact(len(fns), fns...)
}

// AtLeast verifies is at least provided amount of functions will be finished without error in given context timeout/deadline
func (f *verifier) AtLeast(count int, fns ...Verifier) error {
	if count > len(fns) {
		return ErrCountMoreThanLength
	}
	return f.process(len(fns)-count, false, fns...)
}

// OneOf verify at least one function finished without error in given context timeout/deadline
func (f *verifier) OneOf(fns ...Verifier) error {
	return f.AtLeast(1, fns...)
}

// OnlyOne verify exactly one function finished without error in given context timeout/deadline
func (f *verifier) OnlyOne(fns ...Verifier) error {
	return f.Exact(1, fns...)
}

// Exact verify exactly provided amount of functions finished without error in given context timeout/deadline
func (f *verifier) Exact(count int, fns ...Verifier) error {
	if count > len(fns) {
		return ErrCountMoreThanLength
	}
	return f.process(len(fns)-count, true, fns...)
}

// NoOne verifies no one from functions finished without error in given context timeout/deadline
func (f *verifier) NoOne(fns ...Verifier) error {
	return f.Exact(0, fns...)
}

func (f *verifier) process(maxErrorCount int, exact bool, fns ...Verifier) error {
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
			if !exact {
				if doneWithoutError == len(fns)-maxErrorCount {
					return nil
				}
				if doneWithError > maxErrorCount {
					return ErrMaxAmountOfError
				}
				continue
			}

			if doneWithError > maxErrorCount {
				return ErrMaxAmountOfError
			}

			if doneWithoutError > len(fns)-maxErrorCount {
				return ErrMaxAmountOfFinished
			}

			if doneWithError+doneWithoutError == len(fns) {
				return nil
			}
		}
	}
}
