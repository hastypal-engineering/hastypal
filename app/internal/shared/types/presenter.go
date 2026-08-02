package types

type Presenter interface {
	Format(data any) (ServerResponse, error)
}
