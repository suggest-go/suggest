// Package metric holds different metrics for sets similarity compare
package metric

// Metric defined here, is not pure mathematics metric definition as distance between each pair of elements of a set.
// Here we can also ask metric to give as minimum intersection between A and B for given alpha,
// min/max candidate cardinality
type Metric interface {
	// MinY returns the minimum ngram cardinality for candidate
	MinY(alpha float64, size int) int
	// MinY returns the maximum ngram cardinality for candidate
	MaxY(alpha float64, size int) int
	// Threshold returns required intersection between A and B for given alpha
	Threshold(alpha float64, sizeA, sizeB int) int
	// Distance calculate distance between 2 sets
	Distance(inter, sizeA, sizeB int) float64
}
