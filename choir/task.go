package choir

type Task interface {
	Job(run, generation int)
}
