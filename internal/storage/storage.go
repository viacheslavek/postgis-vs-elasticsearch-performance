package storage

type Storage interface {
	AddPoint() error
	DeletePoint() error
	GetPoint() error

	//// и так далее додумать
}

type Point struct {
	// Получше узнать, как представлять точки для добавления
}
