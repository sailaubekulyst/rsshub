package ports

type WorkerAndIntervalRepository interface {
	SetStarted() error
	SetWorkerNs(numWorkers int) error
	SetInterval(newInterval string) error
	GetInterval() (string, error)
	GetStarted() (string, error)
	GetWorkerNs() (int, error)
	CloseStarted() error
}
