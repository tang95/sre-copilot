package pkg

type Server interface {
	Start() error
	Stop() error
}
