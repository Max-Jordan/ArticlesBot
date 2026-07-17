package bot

type Requestable interface {
	params() (Params, error)
	method() string
}
