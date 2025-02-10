package anomalydetector

import (
	"math"
)

type RollingWindow struct {
	Prices []float64
	WindowSize int
}

func NewRollingWindow(size int) *RollingWindow {
	return &RollingWindow{Prices: make([]float64, 0), WindowSize: size}
}

// add new price to window and remove oldest if needed
func (rw *RollingWindow) AddPrice(price float64) {
	if len(rw.Prices) >= rw.WindowSize {
		rw.Prices = rw.Prices[1:] // slicing list
	}
	rw.Prices = append(rw.Prices, price)
}

// mean of prices in window
func (rw *RollingWindow) GetMean() float64 {
	if len(rw.Prices) == 0 {
		return 0
	}
	sum := 0.0
	for _, p := range rw.Prices {
		sum += p
	}
	return sum / float64(len(rw.Prices))
}

// std
func (rw *RollingWindow) GetStandardDeviation() float64 {
	if len(rw.Prices) < 2 {
		return 0
	}
	mean := rw.GetMean()
	var variance float64
	for _, p := range rw.Prices {
		variance += math.Pow(p-mean, 2)
	}
	return math.Sqrt(variance / float64(len(rw.Prices)))
}

// check for anomalies based on z-score and bollinger bands
func (rw *RollingWindow) CheckAnomalies(currentPrice float64) (bool, bool) {
	if len(rw.Prices) < 2 {
		return false, false
	}

	mean := rw.GetMean()
	stdDev := rw.GetStandardDeviation()

	// z-score anomaly, so if its outside x std devs
	zScore := (currentPrice - mean) / stdDev
	isZAnomaly := math.Abs(zScore) > 2

	// bollinger bands
	upperBand := mean + (2 * stdDev)
	lowerBand := mean - (2 * stdDev)
	isBollingerAnomaly := currentPrice > upperBand || currentPrice < lowerBand

	return isZAnomaly, isBollingerAnomaly
}
