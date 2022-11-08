package algo

type Algorithm interface {
	GetIP(*[]string) (string, error)
}
