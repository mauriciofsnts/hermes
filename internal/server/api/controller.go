package api

type Controller interface {
	Route(r Router)
}
