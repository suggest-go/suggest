package suggest

import "github.com/alldroll/suggest/metric"

type Metric = metric.Metric

// CosineMetric returns a Metric that represents Cosine Metric
func CosineMetric() Metric {
	return metric.CosineMetric()
}

// DiceMetric returns a Metric that represents Dice Metric
func DiceMetric() Metric {
	return metric.DiceMetric()
}

// ExactMetric returns a Metric that represents exact matching between 2 ngram sets
func ExactMetric() Metric {
	return metric.ExactMetric()
}

// JaccardMetric returns a Metric that represents Jaccard Metric
func JaccardMetric() Metric {
	return metric.JaccardMetric()
}

// OverlapMetric returns a Metric that represents Overlap metric
func OverlapMetric() Metric {
	return metric.OverlapMetric()
}
