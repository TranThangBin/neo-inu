package pkg

type App interface {
	Init()
	Open() error
	Close() error
}
