package main

import (
	"net/http"
	"testing"
)

func Test_avarageResponseTime(t *testing.T) {

	var tests = []struct {
		n        []float64 // input
		expected float64   // expected result
	}{
		{[]float64{3.5, 3.5}, 3.5},
		{[]float64{1, 2, 3, 4, 5}, 3.0},
		{[]float64{0, 5, 1, 2}, 2},
		{[]float64{6.5, 6.5}, 6.5},
		{[]float64{3, 8, 5, 1}, 4.25},
		{[]float64{1.024, 0.654, 0.654, 5.545, 4.545}, 2.4844},
	}
	for _, tt := range tests {
		actual := avarageResponseTime(tt.n)
		if actual != tt.expected {
			t.Errorf("avarageResponseTime(%d): expected %d, actual %d", tt.n, tt.expected, actual)
		}
	}
}

func Test_showStats(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			showStats(tt.args.w, tt.args.r)
		})
	}
}

func Test_getAvarageResults(t *testing.T) {
	type args struct {
		data float64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getAvarageResults(tt.args.data); got != tt.want {
				t.Errorf("getAvarageResults() = %v, want %v", got, tt.want)
			}
		})
	}
}
