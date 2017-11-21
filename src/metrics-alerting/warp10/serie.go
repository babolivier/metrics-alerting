package warp10

type FloatTimeSerie struct {
	// Class name of the serie
	Class string `json:"c"`
	// Labels of the serie
	Labels map[string]string `json:"l"`
	// Attributes of the serie
	Attributes map[string]string `json:"a"`
	// Datapoints: each element of the slice is a point, represented as a slice
	// with the timestamp as its first element and the value of the point as its
	// second element
	Datapoints [][]float64 `json:"v"`
}
