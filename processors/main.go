package processors

type Processor interface {
	Process() error
}
