package main

import "fmt"

type summary struct {
	mp3Bitrate      int
	numResults      int
	successRate     float64
	medianNumSwaps  float64
	meanElapsedTime float64
	significant     bool
}

func NewSummary(bitrateResults []*result) (*summary, error) {
	n := len(bitrateResults)
	if n == 0 {
		return nil, fmt.Errorf("Error: Cannot summarize an empty slice of results.")
	}

	numSuccess := 0
	numSwaps := []int{}
	elapsedTime := []float64{}

	for _, result := range bitrateResults {
		if result.Success() {
			numSuccess++
		}
		numSwaps = append(numSwaps, result.NumSwaps())
		elapsedTime = append(elapsedTime, result.ElapsedTime())
	}

	medianNumSwaps, _ := median(numSwaps)
	meanElapsedTime, _ := mean(elapsedTime)
	significant := rightTailedBinomialTest(n, numSuccess, 0.5) < 0.05

	return &summary{
		mp3Bitrate:      bitrateResults[0].Mp3Bitrate(),
		numResults:      n,
		successRate:     float64(numSuccess) / float64(n),
		medianNumSwaps:  medianNumSwaps,
		meanElapsedTime: meanElapsedTime,
		significant:     significant,
	}, nil
}

func summarize(results []*result) (map[int]*summary, error) {
	summaries := make(map[int]*summary)
	bitrateResults := make(map[int][]*result)

	for _, res := range results {
		if value, ok := bitrateResults[res.Mp3Bitrate()]; !ok {
			bitrateResults[res.Mp3Bitrate()] = []*result{res}
		} else {
			bitrateResults[res.Mp3Bitrate()] = append(value, res)
		}
	}

	for bitrate, value := range bitrateResults {
		summary, err := NewSummary(value)
		if err != nil {
			return nil, err
		}
		summaries[bitrate] = summary
	}

	return summaries, nil
}

func (s summary) Mp3Bitrate() int {
	return s.mp3Bitrate
}

func (s summary) NumResults() int {
	return s.numResults
}

func (s summary) SuccessRate() float64 {
	return s.successRate
}

func (s summary) MedianNumSwaps() float64 {
	return s.medianNumSwaps
}

func (s summary) MeanElapsedTime() float64 {
	return s.meanElapsedTime
}

func (s summary) Significant() bool {
	return s.significant
}

func SummaryFields() []string {
	return []string{"MP3 kbps", "Success rate", "Swaps (median)", "Seconds passed (mean)", "Significant? (p < 0.05)"}
}

func (s summary) SummaryValuesTSV() string {
	return fmt.Sprintf("%s\t%-*s of %s\t%.0f\t%.2f\t%s", fmtBitrate(s.mp3Bitrate), 4, fmt.Sprintf("%.0f%%", s.successRate*100), fmtNumResults(s.numResults), s.medianNumSwaps, s.meanElapsedTime, fmtBool(s.significant, "Yes", "No"))
}
