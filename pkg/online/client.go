package online

type Client interface {
	GetStatus(id []int) ([]Status, error)
}
