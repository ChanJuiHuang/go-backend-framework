package job

import "fmt"

type ExampleJob struct {
	Frequency string
	Name      string
}

func NewExampleJob() *ExampleJob {
	return &ExampleJob{
		Frequency: "* * * * * *",
		Name:      "example job",
	}
}

func (ej *ExampleJob) GetFrequency() string {
	return ej.Frequency
}

func (ej *ExampleJob) Execute() {
	fmt.Printf("The %s job is executed\n", ej.Name)
}
