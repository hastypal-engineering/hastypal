package types

type Repository[T any] interface {
	Find(criteria Criteria) ([]T, error)
	FindOne(criteria Criteria) (T, error)
	Save(entity T) error
	Update(entity T) error
}
