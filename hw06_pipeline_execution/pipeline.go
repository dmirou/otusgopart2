package hw06_pipeline_execution //nolint:golint,stylecheck

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

// ExecutePipeline executes stages as a pipeline.
func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if len(stages) == 0 {
		out := make(Bi)

		go func() {
			defer close(out)

			for {
				select {
				case <-done:
					return
				default:
				}

				select {
				case v, ok := <-in:
					if !ok {
						return
					}
					out <- v
				case <-done:
					return
				}
			}
		}()

		return out
	}

	cur := in

	for _, stage := range stages {
		st := stage

		cur = runStage(cur, done, st)
	}

	return cur
}

// runStage runs stage st with in and done channels
// and write results into out channel. It returns
// out channel.
func runStage(in In, done In, st Stage) (out Bi) {
	out = make(Bi)

	go func() {
		defer close(out)

		ch := st(in)

		for {
			select {
			case <-done:
				return
			default:
			}

			select {
			case v, ok := <-ch:
				if !ok {
					return
				}

				select {
				case out <- v:
					break
				case <-done:
					return
				}

			case <-done:
				return
			}
		}
	}()

	return out
}
