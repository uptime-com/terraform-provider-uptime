package provider

func ptr[T any](v T) *T {
	return &v
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
