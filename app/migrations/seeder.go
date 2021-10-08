package migrations

// DataSeeds data to seeds
var DataSeeds []interface{} = []interface{}{}

func strptr(s string) *string {
	return &s
}

func float64ptr(f float64) *float64 {
	return &f
}
