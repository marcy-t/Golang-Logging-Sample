package domain

type SampleResponse struct {
	Name string
}

type ApplicationResponse struct {
	Type       string
	StatusCode int
	AppCode    string
	Message    interface{}
}
