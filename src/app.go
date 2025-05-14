package src

type App interface {
	Open() error
	Close() error
}
