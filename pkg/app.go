package pkg

type App interface {
	Init() error
	Open() error
	Close() error
}
