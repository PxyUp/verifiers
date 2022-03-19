# Verifiers

[![codecov](https://codecov.io/gh/PxyUp/verifiers/branch/master/graph/badge.svg)](https://codecov.io/gh/PxyUp/verifiers)

Small GO library for verify async function response.

Provide some basic functionality for conditional check

# Usage

```bash
go get github.com/PxyUp/verifiers
```

**Important**: all function will be finished if condition are matched (it is mean all child routine will be stopped)

### verifiers.All

```go
type Verifier func(ctx context.Context) error

All(fns ...Verifier) error
```

Method verifies is all function finished without error in given context timeout/deadline

Example success:

```go
v := verifiers.New(ctx)
startTime := time.Now()
err := v.All(
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
v := verifiers.New(ctx)
startTime := time.Now()
err := v.All(
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
```

### verifiers.OneOf

```go
type Verifier func(ctx context.Context) error

OneOf(fns ...Verifier) error
```

Method verifies is at least one function finished without error in given context timeout/deadline

```go
v := verifiers.New(ctx)
startTime := time.Now()
err := v.OneOf(
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
assert.True(t, time.Now().Sub(startTime) <= time.Second*3)
```

### verifiers.AtLeast

```go
type Verifier func(ctx context.Context) error

AtLeast(fns ...Verifier) error
```

Method verifies is at least provided amount of functions will be finished without error in given context timeout/deadline

```go
v := verifiers.New(ctx)
startTime := time.Now()
err := v.AtLeast(
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
```