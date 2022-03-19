# Verifiers

[![codecov](https://codecov.io/gh/PxyUp/verifiers/branch/master/graph/badge.svg)](https://codecov.io/gh/PxyUp/verifiers)

Small GO library for verify async function response.

Provide some basic functionality for conditional check

# Usage

```bash
go get github.com/PxyUp/verifiers
```

**Important**: all function will be finished if condition are matched (it is mean all child routine will be stopped)

# Methods

- [verifier.All(...Verifier)](#verifierall) - is equal verifier.Exact(len(fns), fns ...Verifier)
- [verifier.OneOf(...Verifier)](#verifieroneof) - is equal verifier.AtLeast(1, ...Verifier)
- [verifier.AtLeast(int, ...Verifier)](#verifieratleast) 
- [verifier.Exact(int, ...Verifier)](#verifierexact) 
- [verifier.OnlyOne(...Verifier)](#verifieronlyone) - is equal verifier.Exact(1, ...Verifier)
- [verifier.NoOne(...Verifier)](#verifiernoone) - is equal verifier.Exact(0, ...Verifier)

# List of errors
```go
// ErrCountMoreThanLength is configuration error.
// Will return if we expect more function than we provide for verifier.AtLeast or verifier.Exact
verifiers.ErrCountMoreThanLength = errors.New("cant wait more than exists")
// ErrMaxAmountOfError wii be returned some function which we not expect return error
verifiers.ErrMaxAmountOfError = errors.New("verifier reach max amount of error")
// ErrMaxAmountOfFinished will be returned if some other function(which we not expect) return success
verifiers.ErrMaxAmountOfFinished = errors.New("verifier reach max amount success jobs")
```

### verifier.All

```go
type Verifier func(ctx context.Context) error

All(fns ...Verifier) error
```

Method verifies is all function finished without error in given context timeout/deadline

Example success:

```go
verifier := verifiers.New(ctx)
startTime := time.Now()
err := verifier.All(
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
)
// Because we should wait latest function
assert.True(t, time.Now().Sub(startTime) >= time.Second*3)
```

Example error:

```go
verifier := verifiers.New(ctx)
startTime := time.Now()
err := verifier.All(
    func(ctx context.Context) error {
        time.Sleep(time.Second)
        return errors.New("")
    },
    func(ctx context.Context) error {
        time.Sleep(time.Second * 2)
        return nil
    },
    func(ctx context.Context) error {
        time.Sleep(time.Second * 3)
        return nil
    },
)
// Because we will throw error after first return it
assert.True(t, time.Now().Sub(startTime) < time.Second*2)
assert.Err(t, verifiers.ErrMaxAmountOfError)
```

### verifier.OneOf

```go
type Verifier func(ctx context.Context) error

OneOf(fns ...Verifier) error
```

Method verifies is at least one function finished without error in given context timeout/deadline

```go
verifier := verifiers.New(ctx)
startTime := time.Now()
err := verifier.OneOf(
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
        return nil
    },
)
// Because we should wait second function
assert.True(t, time.Now().Sub(startTime) >= time.Second*1)
assert.True(t, time.Now().Sub(startTime) < time.Second*2)
assert.Nil(t, err)
```

### verifier.AtLeast

```go
type Verifier func(ctx context.Context) error

AtLeast(count int, fns ...Verifier) error
```

Method verifies is at least provided amount of functions will be finished without error in given context timeout/deadline

```go
verifier := verifiers.New(ctx)
startTime := time.Now()
err := verifier.AtLeast(
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
        return errors.New("")
    },
)
// Because we should wait second function
assert.True(t, time.Now().Sub(startTime) > time.Second*1)
assert.True(t, time.Now().Sub(startTime) <= time.Second*3)
assert.Nil(t, err)
```

### verifier.Exact

```go
type Verifier func(ctx context.Context) error

Exact(count int, fns ...Verifier) error
```

Method verify exactly provided amount of functions finished without error in given context timeout/deadline

```go
verifier := verifiers.New(ctx)
startTime := time.Now()
err := verifier.Exact(
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
        return errors.New("")
    },
)
// Because we should wait three function(all functions should be finished)
assert.True(t, time.Now().Sub(startTime) >= time.Second*3)
assert.Nil(t, err)
```

### verifier.OnlyOne

```go
type Verifier func(ctx context.Context) error

OnyOne(fns ...Verifier) error
```

Method verify exactly one function finished without error in given context timeout/deadline

```go
verifier := verifiers.New(ctx)
startTime := time.Now()
err := verifier.OnlyOne(
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
)
// Because we should wait three function(all functions should be finished)
assert.True(t, time.Now().Sub(startTime) >= time.Second*3)
assert.Nil(t, err)
```

### verifier.NoOne

```go
type Verifier func(ctx context.Context) error

NoOne(fns ...Verifier) error
```

Method verifies no one from functions finished without error in given context timeout/deadline

```go
verifier := verifiers.New(ctx)
startTime := time.Now()
err := verifier.NoOne(
    func(ctx context.Context) error {
        time.Sleep(time.Second)
        return errors.New("")
    },
    func(ctx context.Context) error {
        time.Sleep(time.Second * 2)
        return errors.New("")
    },
    func(ctx context.Context) error {
        time.Sleep(time.Second * 3)
        return errors.New("")
    },
)
// Because we should wait three function(all functions should be finished)
assert.True(t, time.Now().Sub(startTime) >= time.Second*3)
assert.Nil(t, err)
```