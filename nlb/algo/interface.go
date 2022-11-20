package algo

import "nlb/k8s"

type Algorithm interface {
	GetIP(*[]*k8s.PodDetails, string) (string, error)
}
