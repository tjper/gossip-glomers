# gossip-glomers

## TODO

## Takeaways

1. Behavior of receivers is simpler to intuit when point receivers are always used. However, it is tempting to not use a pointer when a method is not mutating the receiver.

2. If an instance method is passed to a function, the receiver is evalualted when the instance method is passed.

This is imporant to remember as the following code will lead to nil-pointer dereference panic. When `v.handleSignal` is passed to `antenna.Handle`, `v` is evaluated and `v.antenna` has not been initialized. Therefore, when a signal is handled in the future, `s.antenna.Broadcast` will result in a nil-pointer dereference error, as `s.antenna` is nil.

```go
func NewVehicle() *Vehicle {
  v := new(Vehicle)

  antenna := antenna.New()
  antenna.Handle("signal", v.handleSignal)

  v.antenna = antenna
  return v
}

func (v Vehicle) handleSignal(_ []byte) error {
	return s.antenna.Broadcast("received")
}
```

There are numerous ways to handle this. I'll break down three approaches and rank order them based in there pros/cons.

### Dependency Injection

In this scenario a Vehicle has an antenna that is used to broadcast information. Instead of having the Vehicle be concerned with the configuraton the antenna, offload this task onto another type, and only accept a configured antenna as an argument to the NewVehicle function. With this solution the vehicle may still may broadcast information, but is no longer responsible for configuring how received data is handled.

(TODO)

```go
// vehicle.go

type Broadcaster interface {
  Broadcast(string) error
}

func NewVehicle(antenna Broadcaster) *Vehicle {
  return &Vehicle{
    antenna: antenna,
  }
}

func (v Vehicle) start(_ []byte) error {
	return s.antenna.Broadcast("Rvvvvvvvv...")
}
```

```go
// antenna.go

func NewAntenna() *Antenna {
  return &Antenna{
    r io.Reader
    w io.Writer
  }
}

func (a *Antenna) Handle(func(

func (a Antenna) Broadcast(msg string) error {
  _, err := fmt.Fprint(a.w, msg)
  return err
}
```

### Pointer Receivers

(TODO)

### Correct Constructor Sequence

(TODO)

3. Ensuring a slice is not leaking memory may require setting an item to nil. This can be done several ways.

(TODO)

4. Passing a loop variable to a closure, or taking a reference of a loop variable.