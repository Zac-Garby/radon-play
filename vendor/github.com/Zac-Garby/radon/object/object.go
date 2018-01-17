package object

// An Object is the interface which every object
// type implements.
type Object interface {
	String() string
	Debug() string
	Equals(Object) bool
	Type() Type
}

// Collection is a child interface of Object,
// which represents an object which can be
// thought of as a list of items
type Collection interface {
	Object
	Elements() []Object
	GetIndex(int) Object
	SetIndex(int, Object)
}

// Container is a child interface of Object,
// which can be accessed by keys - like a map
type Container interface {
	Object
	GetKey(Object) Object
	SetKey(Object, Object)
}

// Methoder is any object which has methods.
// Methods are accesses using the dot operator.
type Methoder interface {
	Object
	GetMethod(name string) (*Builtin, bool)
}
