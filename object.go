package dingo

/*
 Types of object, they are bit flag and can be combined.
 These flags are used in:
  - dingo.Use
  - dingo.Listen
*/
var ObjT = struct {
	/*
		when this type used in dingo.App.Use, it means let
		dingo decide which type would be registered to dingo.App.
	*/
	DEFAULT int
	// this object provides dingo.Reporter interface
	REPORTER int
	// this object provides dingo.Store interface
	STORE int
	// this object provides dingo.Producer interface
	PRODUCER int
	// this object provides dingo.Consumer/dingo.NamedConsumer interface
	CONSUMER int
	// this is a dingo.mapper object
	MAPPER int
	// this is a dingo.worker object
	WORKER int
	// this object provides dingo.bridge interface
	BRIDGE int
	// this object provides dingo.NamedConsumer interface
	NAMED_CONSUMER int
	/*
		all object types, when used in dingo.App.Listen, it means
		listen to events from all possible origins.
	*/
	ALL int
}{
	0,
	(1 << 0),
	(1 << 1),
	(1 << 2),
	(1 << 3),
	(1 << 4),
	(1 << 5),
	(1 << 6),
	(1 << 7),
	int(^uint(0) >> 1), // Max int
}

/*
 An object implements this interface means:
  - dingo can have a trigger to release the resource allocated by this object.
  - dingo can aggregate events raised from this object, (those events can be subscribed
    by dingo.App.Listen)

 All objects attached via dingo.App.Use should implement this interface.
*/
type Object interface {
	// what dingo expects from this object
	Expect(types int) error

	// allow dingo to attach event channels used in this object
	Events() ([]<-chan *Event, error)

	// releasing resource
	Close() error
}
