package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func runStage(in In, done In, stage Stage) Out {
	out := make(Bi)
	go func() {
		defer close(out)
		for i := range stage(in) {
			select {
			case <-done:
				return
			case out <- i:
			}
		}
	}()
	return out
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := make(Bi)
	go func() {
		defer close(out)
		currentIn := in
		for _, stage := range stages {
			currentIn = runStage(currentIn, done, stage)
		}
		for i := range currentIn {
			select {
			case <-done:
				return
			case out <- i:
			}
		}
	}()
	return out
}
