package provider

var App *Application

func init() {
	application, err := InitializeApplication()

	if err != nil {
		panic(err)
	}
	App = application
}
