package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		out := make(chan interface{})
		middle := stage(in)
		go func(middle In, out chan interface{}) {
			defer close(out)
			for i := range middle {
				select {
				case <-done:
					return
				case out <- i:
				}
			}
		}(middle, out)
		in = out
	}
	return in
}
