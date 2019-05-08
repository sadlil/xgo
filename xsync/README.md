### xsync
Package `xsync` contains some helper functions that are similar or incremental features
of Golang's `sync` package.


#### MustOnce
MustOnce is provides an API to perform exactly one action until the action is succeed.
All failed run for the action will have the opportunity to run the action again.

```go
var mo xsync.MustOnce

mo.Do(func() error {
	// Do your task here and return success or error
})
```

Do calls the function f if and only if f is not being called successfully
for this instance of MustOnce. In other words, given
`var mo MustOnce`
if `mo.Do(f)` is called multiple times, only until the first invocation to f that is
succeeded will be executed. After the first successful invocation of f
Do will not invoke f even if f has a different value in each invocation.
A new instance of MustOnce is required for different function to execute successfully.

Do is intended for initialization that must be run exactly once successfully. f must
return an error in case of any failed invocation. it may be necessary to use
a function literal to capture the arguments to a function to be invoked by Do.

Because no call to Do returns until the one call to f returns, if f causes
Do to be called, it will deadlock.

If f panics, Do considers it to have returned successfully; future calls of Do return
without calling f.

This implementation is an adaption from https://golang.org/pkg/sync/#Once. It only varies
by the successful invocation of f.