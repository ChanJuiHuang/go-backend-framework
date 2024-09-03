package job

import "fmt"

type ExampleJob struct {
	frequency string
	name      string
}

func NewExampleJob() *ExampleJob {
	return &ExampleJob{
		frequency: "* * * * * *",
		name:      "example job",
	}
}

func (job *ExampleJob) GetFrequency() string {
	return job.frequency
}

func (job *ExampleJob) Execute() {
	fmt.Printf("The %s is finished\n", job.name)
}
