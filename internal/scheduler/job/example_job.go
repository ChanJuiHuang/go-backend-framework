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

func (ej *ExampleJob) GetFrequency() string {
	return ej.frequency
}

func (ej *ExampleJob) Execute() {
	fmt.Printf("The %s job is executed\n", ej.name)
}
