package fft

import (
    "gonum.org/v1/gonum/dsp/fourier"
)

func CalculateFFT(data []float64) []complex128 {
    fft := fourier.NewFFT(len(data))
    return fft.Coefficients(nil, data)
}

func DifferentialAnalysis(data1 []float64, data2 []float64) []complex128 {
    if (len(data1) != len(data2)) { return nil }

    diff := make([]float64, len(data1))
    for i := range data1 {
        diff[i] = data1[i] - data2[i]
    }

    return CalculateFFT(diff)
}