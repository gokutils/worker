<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# worker

```go
import "github.com/gokutils/worker"
```

## Index

- [type Option](<#Option>)
  - [func OnError\(callback func\(error\)\) Option](<#OnError>)
  - [func WithTimout\(timeout time.Duration\) Option](<#WithTimout>)
- [type Queux](<#Queux>)
  - [func NewQueux\[T any\]\(ctx context.Context, num int, Worker WorkerQeux\[T\], opts ...Option\) \*Queux\[T\]](<#NewQueux>)
  - [func \(impl \*Queux\[T\]\) Push\(v T\) error](<#Queux[T].Push>)
  - [func \(impl \*Queux\[T\]\) Start\(\)](<#Queux[T].Start>)
  - [func \(impl \*Queux\[T\]\) Stop\(\)](<#Queux[T].Stop>)
- [type Tick](<#Tick>)
  - [func NewTick\(ctx context.Context, num int, tickTime time.Duration, Worker WorkerTick, opts ...Option\) \*Tick](<#NewTick>)
  - [func \(impl \*Tick\) Start\(\)](<#Tick.Start>)
  - [func \(impl \*Tick\) Stop\(\)](<#Tick.Stop>)
  - [func \(impl \*Tick\) TickNow\(\) error](<#Tick.TickNow>)
- [type WorkerQeux](<#WorkerQeux>)
- [type WorkerTick](<#WorkerTick>)


<a name="Option"></a>
## type Option



```go
type Option func(*params)
```

<a name="OnError"></a>
### func OnError

```go
func OnError(callback func(error)) Option
```



<a name="WithTimout"></a>
### func WithTimout

```go
func WithTimout(timeout time.Duration) Option
```



<a name="Queux"></a>
## type Queux



```go
type Queux[T any] struct {
    // contains filtered or unexported fields
}
```

<a name="NewQueux"></a>
### func NewQueux

```go
func NewQueux[T any](ctx context.Context, num int, Worker WorkerQeux[T], opts ...Option) *Queux[T]
```



<a name="Queux[T].Push"></a>
### func \(\*Queux\[T\]\) Push

```go
func (impl *Queux[T]) Push(v T) error
```



<a name="Queux[T].Start"></a>
### func \(\*Queux\[T\]\) Start

```go
func (impl *Queux[T]) Start()
```



<a name="Queux[T].Stop"></a>
### func \(\*Queux\[T\]\) Stop

```go
func (impl *Queux[T]) Stop()
```



<a name="Tick"></a>
## type Tick



```go
type Tick struct {
    // contains filtered or unexported fields
}
```

<a name="NewTick"></a>
### func NewTick

```go
func NewTick(ctx context.Context, num int, tickTime time.Duration, Worker WorkerTick, opts ...Option) *Tick
```



<a name="Tick.Start"></a>
### func \(\*Tick\) Start

```go
func (impl *Tick) Start()
```



<a name="Tick.Stop"></a>
### func \(\*Tick\) Stop

```go
func (impl *Tick) Stop()
```



<a name="Tick.TickNow"></a>
### func \(\*Tick\) TickNow

```go
func (impl *Tick) TickNow() error
```



<a name="WorkerQeux"></a>
## type WorkerQeux



```go
type WorkerQeux[T any] interface {
    Do(ctx context.Context, task T) error
}
```

<a name="WorkerTick"></a>
## type WorkerTick



```go
type WorkerTick interface {
    Do(context.Context) error
}
```

Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)
