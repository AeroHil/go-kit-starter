package common

// endpoint enum
type endpointEnum struct {
	Health   string
	Greeting string
}
// variables from endpoint enums
var EndpointNames = endpointEnum{
	Health:   "Health",
	Greeting: "Greeting",
}
// iterable endpoint names
var EndpointNamesMap = map[string]bool{
	EndpointNames.Health:   true,
	EndpointNames.Greeting: true,
}
