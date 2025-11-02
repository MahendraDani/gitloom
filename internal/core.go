package internal

type Object interface {
	Type() string
	Serialize() []byte
}
