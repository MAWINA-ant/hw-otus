package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func runStage(in In, done In, stage Stage) Out {
	originalOut := stage(in)
	out := make(Bi)
	go func() {
		defer close(out)
		for i := range originalOut {
			select {
			case <-done:
				return
			default:
				select {
				case <-done:
					return
				case out <- i:
				}
			}
		}
	}()
	return out
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {

	currentChannel := in
	for _, stage := range stages {
		currentChannel = runStage(currentChannel, done, stage)
	}

	return currentChannel
}
