package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func runStage(done In, in In, stage Stage) Out {
	// create channel for next stage
	nextStageChan := make(Bi)

	go func() {
		// close unused channel in end
		defer close(nextStageChan)

		for {
			select {
			case <-done:
				// break stage
				return
			case v, ok := <-in:
				if !ok {
					// previous channel was closed, close current
					return
				}
				// throw value for next stage
				nextStageChan <- v
			}
		}
	}()

	// run current stage
	return stage(nextStageChan)
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in
	for _, s := range stages {
		stage := s

		// run current stage and generate channel for next stage
		out = runStage(done, out, stage)
	}

	return out
}
