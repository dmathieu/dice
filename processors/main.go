package processors

import "context"

type Processor interface {
	Process(context.Context) error
}
