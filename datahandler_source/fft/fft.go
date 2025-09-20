package fft

import (
    "gonum.org/v1/gonum/dsp/fourier"
)

// Perform FFT on float64 array
func CalculateFFT(data []float64) []complex128 {
    fft := fourier.NewFFT(len(data))
    return fft.Coefficients(nil, data)
}
