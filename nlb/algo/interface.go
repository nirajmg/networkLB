package algo

type algo interface {
	GetIP([]string) (string, error)
}
