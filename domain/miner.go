package domain

type MainerController interface {
	Stop() bool
	Start() bool
	Status() bool
}
