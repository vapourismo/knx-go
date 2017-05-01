package knx

type reaper <-chan struct{}

type scythe func()

func makeReaper() (reaper, scythe) {
	r := make(chan struct{}, 1)

	s := func () {
		select {
			case r <- struct{}{}:
			default:
		}
	}

	return r, s
}
