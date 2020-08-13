package kernel

type serve func()

var containers []serve

func Inject(handle serve) {
	containers = append(containers, handle)
}

func Load() {
	if len(containers) > 0 {
		for _, handle := range containers {
			handle()
		}
	}
}
