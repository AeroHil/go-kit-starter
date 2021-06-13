package common

// endpoint enum
type endpointEnum struct {
	Health   string
	Greeting string
}

// EndpointNames - variables from endpoint enums
var EndpointNames = endpointEnum{
	Health:   "Health",
	Greeting: "Greeting",
}

// EndpointNamesMap - iterable endpoint names
var EndpointNamesMap = map[string]bool{
	EndpointNames.Health:   true,
	EndpointNames.Greeting: true,
}
