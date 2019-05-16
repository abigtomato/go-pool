package pool

var DefaultPool = initPool(1000)

func Go(job Job) error {
	return DefaultPool.submit(job)
}

func Close() {
	DefaultPool.close()
}
