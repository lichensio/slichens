package student

import "math"

// WelchTest performs Welch's t-test on two groups of observations.
// It compares the means of the two groups, assuming unequal variances.
// The function returns the t-value and p-value.
func WelchTest(groupA, groupB []float64) (float64, float64) {
	nA := len(groupA)
	nB := len(groupB)

	// Calculate the means
	meanA := Mean(groupA)
	meanB := Mean(groupB)

	// Calculate the variances
	varA := Variance(groupA)
	varB := Variance(groupB)

	// Calculate the degrees of freedom
	nAFloat := float64(nA)
	nBFloat := float64(nB)
	dfNum := math.Pow((varA/nAFloat + varB/nBFloat), 2)
	dfDen := (math.Pow(varA/nAFloat, 2) / (nAFloat - 1)) + (math.Pow(varB/nBFloat, 2) / (nBFloat - 1))
	df := dfNum / dfDen

	// Calculate the t-value
	t := (meanA - meanB) / math.Sqrt((varA/nAFloat)+(varB/nBFloat))

	// Calculate the p-value
	p := 2 * (1 - StudentCDF(t, df))

	return t, p
}

// Mean calculates the mean of a slice of float64 values.
func Mean(data []float64) float64 {
	sum := 0.0
	for _, value := range data {
		sum += value
	}
	return sum / float64(len(data))
}

// Variance calculates the variance of a slice of float64 values.
func Variance(data []float64) float64 {
	mean := Mean(data)
	sum := 0.0
	for _, value := range data {
		sum += math.Pow(value-mean, 2)
	}
	return sum / float64(len(data)-1)
}

// StudentCDF calculates the cumulative distribution function of the Student's t-distribution.
// It returns the probability that a value is less than or equal to the given t-value and degrees of freedom.
func StudentCDF(t, df float64) float64 {
	x := (t + math.Sqrt(t*t+df)) / (2 * math.Sqrt(t*t+df))
	return StudentBetaInc(x, df/2, 0.5) / 2
}

// StudentBetaInc calculates the incomplete beta function for the Student's t-distribution.
// It is used to calculate the cumulative distribution function (CDF).
func StudentBetaInc(x, a, b float64) float64 {
	// Implementation of the incomplete beta function using continued fractions
	epsilon := 1e-15
	maxIterations := 1000

	// Initialize variables
	f := 1.0
	c := 1.0
	d := 0.0
	h := 1.0
	n := 1.0

	// Perform continued fraction expansion
	for i := 1; i <= maxIterations; i++ {
		aPlusB := a + b
		an := float64(i) * (b - float64(i)) * x / ((aPlusB - 2*float64(i) + 1) * (aPlusB - 2*float64(i) + 2))
		bn := -x * float64(i) * (a + float64(i)) * (aPlusB - float64(i)) / ((aPlusB - 2*float64(i) + 2) * (aPlusB - 2*float64(i) + 3))
		d = 1 + an*bn
		if math.Abs(d) < epsilon {
			d = epsilon
		}
		c = 1 + an/c
		if math.Abs(c) < epsilon {
			c = epsilon
		}
		d = 1 / d
		delta := c * d
		h *= delta
		if math.Abs(delta-1) < epsilon {
			break
		}
		n *= -1
	}

	return f * h * math.Pow(x, a) * math.Pow(1-x, b) / a
}
