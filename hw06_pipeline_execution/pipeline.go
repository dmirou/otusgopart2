package hw06_pipeline_execution //nolint:golint,stylecheck

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if stages == nil || len(stages) == 0 {
		out := make(Bi)

		close(out)

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
// and write results into out channel.
func runStage(in In, done In, st Stage) Out {
	out := make(Bi)

	go func() {
		defer close(out)

		ch := st(in)

		if done == nil {
			for {
				v, ok := <-ch
				if !ok {
					return
				}

				out <- v
			}
		}

		for {
			select {
			case <-done:
				return
			default:
			}

			select {
			case v := <-ch:
				out <- v
			default:
			}

			select {
			case <-done:
				return
			default:
			}
		}
	}()

	return out
}
