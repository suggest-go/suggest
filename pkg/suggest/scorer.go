package suggest

import (
	"github.com/suggest-go/suggest/pkg/merger"
	"github.com/suggest-go/suggest/pkg/metric"
)

// Scorer is responsible for scoring an index position
type Scorer interface {
	// Score returns the score of the given candidate
	Score(position merger.MergeCandidate) float64
}

type metricScorer struct {
	sizeA, sizeB int
	metric metric.Metric
}

// NewMetricScorer creates a new scorer that uses metric as a score value
func NewMetricScorer(metric metric.Metric, sizeA, sizeB int) Scorer {
	return &metricScorer{
		metric: metric,
		sizeA: sizeA,
		sizeB: sizeB,
	}
}

// Score returns the score of the given candidate
func (s *metricScorer) Score(candidate merger.MergeCandidate) float64 {
	return 1 - s.metric.Distance(candidate.Overlap(), s.sizeA, s.sizeB)
}
