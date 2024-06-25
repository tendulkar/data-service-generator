package parser

import (
	"encoding/json"
	"os"

	"stellarsky.ai/platform/codegen/data-service-generator/base"
)

func ReadJsonTo[T any](filename string) (*T, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		base.LOG.Error("Error reading json file", "error", err, "filename", filename)
		return nil, err
	}
	var t T
	err = json.Unmarshal(data, &t)
	if err != nil {
		base.LOG.Error("Error parsing json", "error", err, "filename", filename)
		return &t, err
	}
	return &t, nil
}

func MustReadJsonTo[T any](filename string) *T {
	t, err := ReadJsonTo[T](filename)
	if err != nil {
		base.LOG.Error("MustReadJsonTo FAILED PANICING", "error", err, "filename", filename)
		panic(err)
	}
	return t
}

func ReadJsonToSlice[T any](filename string) ([]T, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		base.LOG.Error("Error reading json file", "error", err, "filename", filename)
		return nil, err
	}
	var t []T
	err = json.Unmarshal(data, &t)
	if err != nil {
		base.LOG.Error("Error parsing json", "error", err, "filename", filename)
		return t, err
	}
	return t, nil
}

func MustReadJsonToSlice[T any](filename string) []T {
	t, err := ReadJsonToSlice[T](filename)
	if err != nil {
		base.LOG.Error("MustReadJsonToSlice FAILED PANICING", "error", err, "filename", filename)
		panic(err)
	}
	return t
}
