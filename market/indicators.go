package market

import "math"

// =============================================================================
// 技术指标计算函数
// 这些函数是纯计算函数，不依赖任何外部状态
// =============================================================================

// calculateEMA 计算EMA
func calculateEMA(klines []Kline, period int) float64 {
	if len(klines) < period {
		return 0
	}

	// 计算SMA作为初始EMA
	sum := 0.0
	for i := 0; i < period; i++ {
		sum += klines[i].Close
	}
	ema := sum / float64(period)

	// 计算EMA
	multiplier := 2.0 / float64(period+1)
	for i := period; i < len(klines); i++ {
		ema = (klines[i].Close-ema)*multiplier + ema
	}

	return ema
}

// calculateMACD 计算MACD
func calculateMACD(klines []Kline) float64 {
	if len(klines) < 26 {
		return 0
	}

	// 计算12期和26期EMA
	ema12 := calculateEMA(klines, 12)
	ema26 := calculateEMA(klines, 26)

	// MACD = EMA12 - EMA26
	return ema12 - ema26
}

// calculateRSI 计算RSI
func calculateRSI(klines []Kline, period int) float64 {
	if len(klines) <= period {
		return 0
	}

	gains := 0.0
	losses := 0.0

	// 计算初始平均涨跌幅
	for i := 1; i <= period; i++ {
		change := klines[i].Close - klines[i-1].Close
		if change > 0 {
			gains += change
		} else {
			losses += -change
		}
	}

	avgGain := gains / float64(period)
	avgLoss := losses / float64(period)

	// 使用Wilder平滑方法计算后续RSI
	for i := period + 1; i < len(klines); i++ {
		change := klines[i].Close - klines[i-1].Close
		if change > 0 {
			avgGain = (avgGain*float64(period-1) + change) / float64(period)
			avgLoss = (avgLoss * float64(period-1)) / float64(period)
		} else {
			avgGain = (avgGain * float64(period-1)) / float64(period)
			avgLoss = (avgLoss*float64(period-1) + (-change)) / float64(period)
		}
	}

	if avgLoss == 0 {
		return 100
	}

	rs := avgGain / avgLoss
	rsi := 100 - (100 / (1 + rs))

	return rsi
}

// calculateATR 计算ATR
func calculateATR(klines []Kline, period int) float64 {
	if len(klines) <= period {
		return 0
	}

	trs := make([]float64, len(klines))
	for i := 1; i < len(klines); i++ {
		high := klines[i].High
		low := klines[i].Low
		prevClose := klines[i-1].Close

		tr1 := high - low
		tr2 := math.Abs(high - prevClose)
		tr3 := math.Abs(low - prevClose)

		trs[i] = math.Max(tr1, math.Max(tr2, tr3))
	}

	// 计算初始ATR
	sum := 0.0
	for i := 1; i <= period; i++ {
		sum += trs[i]
	}
	atr := sum / float64(period)

	// Wilder平滑
	for i := period + 1; i < len(klines); i++ {
		atr = (atr*float64(period-1) + trs[i]) / float64(period)
	}

	return atr
}

// calculateATRSeries 计算ATR序列，返回最近10个点的ATR值
func calculateATRSeries(klines []Kline, period int) []float64 {
	if len(klines) <= period {
		return []float64{}
	}

	// 计算所有True Range
	trs := make([]float64, len(klines))
	for i := 1; i < len(klines); i++ {
		high := klines[i].High
		low := klines[i].Low
		prevClose := klines[i-1].Close

		tr1 := high - low
		tr2 := math.Abs(high - prevClose)
		tr3 := math.Abs(low - prevClose)

		trs[i] = math.Max(tr1, math.Max(tr2, tr3))
	}

	// 计算初始ATR (第period个点)
	sum := 0.0
	for i := 1; i <= period; i++ {
		sum += trs[i]
	}
	atr := sum / float64(period)

	// 收集所有ATR值
	allATRs := make([]float64, 0, len(klines)-period)
	allATRs = append(allATRs, atr)

	// Wilder平滑计算后续ATR
	for i := period + 1; i < len(klines); i++ {
		atr = (atr*float64(period-1) + trs[i]) / float64(period)
		allATRs = append(allATRs, atr)
	}

	// 返回最近10个点
	if len(allATRs) > 10 {
		return allATRs[len(allATRs)-10:]
	}
	return allATRs
}

// =============================================================================
// Efficiency Ratio (ER) - Kaufman 效率系数
// =============================================================================

// calculateEfficiencyRatio 计算 Kaufman 效率系数
// ER = |Direction| / Volatility
// Direction = Close[n] - Close[0] (净价格变动)
// Volatility = Sum of |Close[i] - Close[i-1]| (每日变动之和)
// 返回值范围 0.0-1.0，值越大趋势越强
// 数据不足时返回 NaN（使用 math.IsNaN() 检测）
func calculateEfficiencyRatio(klines []Kline, period int) float64 {
	if len(klines) <= period {
		return math.NaN()
	}

	// 使用最近 period+1 个数据点
	startIdx := len(klines) - period - 1
	endIdx := len(klines) - 1

	// 计算 Direction (净变动)
	direction := math.Abs(klines[endIdx].Close - klines[startIdx].Close)

	// 计算 Volatility (每日变动之和)
	volatility := 0.0
	for i := startIdx + 1; i <= endIdx; i++ {
		volatility += math.Abs(klines[i].Close - klines[i-1].Close)
	}

	// 避免除零
	if volatility == 0 {
		return 0
	}

	er := direction / volatility

	// ER 理论上应该在 0-1 之间，但由于浮点精度可能略微超出
	if er > 1 {
		er = 1
	}

	return er
}

// calculateERSeries 计算 ER 序列，返回最近 10 个点的 ER 值
func calculateERSeries(klines []Kline, period int) []float64 {
	if len(klines) <= period {
		return []float64{}
	}

	// 计算所有可能的 ER 值
	allERs := make([]float64, 0, len(klines)-period)

	for endIdx := period; endIdx < len(klines); endIdx++ {
		startIdx := endIdx - period

		// 计算 Direction (净变动)
		direction := math.Abs(klines[endIdx].Close - klines[startIdx].Close)

		// 计算 Volatility (每日变动之和)
		volatility := 0.0
		for i := startIdx + 1; i <= endIdx; i++ {
			volatility += math.Abs(klines[i].Close - klines[i-1].Close)
		}

		var er float64
		if volatility == 0 {
			er = 0
		} else {
			er = direction / volatility
			if er > 1 {
				er = 1
			}
		}
		allERs = append(allERs, er)
	}

	// 返回最近 10 个点
	if len(allERs) > 10 {
		return allERs[len(allERs)-10:]
	}
	return allERs
}

// =============================================================================
// Bollinger Bands 布林带
// =============================================================================

// calculateBollingerBands 计算布林带指标
// 返回 %B (价格在带中的位置) 和 Bandwidth (带宽，波动率指标)
// %B = (Close - Lower) / (Upper - Lower)
// Bandwidth = (Upper - Lower) / Middle
// 数据不足时返回 NaN, NaN（使用 math.IsNaN() 检测）
func calculateBollingerBands(klines []Kline, period int, stdDevMultiplier float64) (percentB, bandwidth float64) {
	if len(klines) < period {
		return math.NaN(), math.NaN()
	}

	// 取最近 period 个收盘价
	startIdx := len(klines) - period
	closes := make([]float64, period)
	for i := 0; i < period; i++ {
		closes[i] = klines[startIdx+i].Close
	}

	// 计算 SMA (中轨)
	sum := 0.0
	for _, c := range closes {
		sum += c
	}
	middle := sum / float64(period)

	// 计算标准差
	variance := 0.0
	for _, c := range closes {
		diff := c - middle
		variance += diff * diff
	}
	stdDev := math.Sqrt(variance / float64(period))

	// 计算上下轨
	upper := middle + stdDevMultiplier*stdDev
	lower := middle - stdDevMultiplier*stdDev

	// 当前价格
	currentPrice := klines[len(klines)-1].Close

	// 计算 %B
	bandWidth := upper - lower
	if bandWidth == 0 {
		// 价格恒定，没有波动
		percentB = 0.5
		bandwidth = 0
		return
	}
	percentB = (currentPrice - lower) / bandWidth

	// 计算 Bandwidth (相对带宽)
	bandwidth = bandWidth / middle

	return
}

// calculateBollingerSeries 计算 Bollinger Bands 序列，返回最近 10 个点的 %B 和 Bandwidth 值
func calculateBollingerSeries(klines []Kline, period int, stdDevMultiplier float64) (percentBs, bandwidths []float64) {
	if len(klines) < period {
		return []float64{}, []float64{}
	}

	// 计算所有可能的 Bollinger 值
	allPercentBs := make([]float64, 0, len(klines)-period+1)
	allBandwidths := make([]float64, 0, len(klines)-period+1)

	for endIdx := period; endIdx <= len(klines); endIdx++ {
		// 取 [endIdx-period, endIdx) 范围的收盘价
		startIdx := endIdx - period
		closes := make([]float64, period)
		for i := 0; i < period; i++ {
			closes[i] = klines[startIdx+i].Close
		}

		// 计算 SMA (中轨)
		sum := 0.0
		for _, c := range closes {
			sum += c
		}
		middle := sum / float64(period)

		// 计算标准差
		variance := 0.0
		for _, c := range closes {
			diff := c - middle
			variance += diff * diff
		}
		stdDev := math.Sqrt(variance / float64(period))

		// 计算上下轨
		upper := middle + stdDevMultiplier*stdDev
		lower := middle - stdDevMultiplier*stdDev

		// 当前价格（窗口最后一个）
		currentPrice := klines[endIdx-1].Close

		// 计算 %B 和 Bandwidth
		bandWidth := upper - lower
		var percentB, bw float64
		if bandWidth == 0 {
			percentB = 0.5
			bw = 0
		} else {
			percentB = (currentPrice - lower) / bandWidth
			bw = bandWidth / middle
		}

		allPercentBs = append(allPercentBs, percentB)
		allBandwidths = append(allBandwidths, bw)
	}

	// 返回最近 10 个点
	if len(allPercentBs) > 10 {
		return allPercentBs[len(allPercentBs)-10:], allBandwidths[len(allBandwidths)-10:]
	}
	return allPercentBs, allBandwidths
}

    // calculateSMA 计算SMA (简单移动平均线)
    func calculateSMA(klines []Kline, period int) float64 {
	        if len(klines) < period {
		            return 0
	        }

	        sum := 0.0
	        // 取最后 period 个点
	        start := len(klines) - period
	        for i := start; i < len(klines); i++ {
		    sum += klines[i].Close
	        }

	        return sum / float64(period)
}

		// =============================================================================
		// 缠论专用 MACD (ChanLun MACD)
		// 参数: Fast=34, Slow=89, Signal=13 (注意: 信号线使用 SMA 而非 EMA)
		// =============================================================================
		
		// CalculateChanLunMACDState 计算缠论MACD的状态（金叉/死叉）
		// 返回值: 
		// - macdLine: DIF (快线 - 慢线)
		// - signalLine: DEA (DIF的SMA)
		// - histogram: 柱状图
		// - crossType: 0=无交叉, 1=金叉(Bullish), 2=死叉(Bearish)
		func CalculateChanLunMACDState(klines []Kline) (macdLine, signalLine, histogram float64, crossType int) {
			// 确保有足够的数据计算
			// 需要足够的数据来预热 EMA (通常建议 3-4 倍周期长度)
			if len(klines) < 100 { 
				return 0, 0, 0, 0
			}
		
			fastPeriod := 34
			slowPeriod := 89
			signalPeriod := 13
		
			// 1. 计算快线 EMA 34 序列
			emaFast := calculateEMASeries(klines, fastPeriod)
			// 2. 计算慢线 EMA 89 序列
			emaSlow := calculateEMASeries(klines, slowPeriod)
		
			// 3. 计算 DIF (MACD Line) 序列
			// 我们至少需要 signalPeriod 个点来计算最后的 SMA
			dataLen := len(klines)
			if len(emaFast) != dataLen || len(emaSlow) != dataLen {
				return 0, 0, 0, 0
			}
		
			difSeries := make([]float64, dataLen)
			for i := 0; i < dataLen; i++ {
				difSeries[i] = emaFast[i] - emaSlow[i]
			}
		
			// 4. 计算 DEA (Signal Line) - 使用 SMA 算法 (这是该脚本的特殊之处)
			// 我们只需要计算最后两个点的 DEA 来判断交叉
			prevDea := calculateSMAFromSeries(difSeries[:dataLen-1], signalPeriod)
			currDea := calculateSMAFromSeries(difSeries, signalPeriod)
		
			// 获取最后两个点的 DIF
			prevDif := difSeries[dataLen-2]
			currDif := difSeries[dataLen-1]
		
			// 5. 计算当前的柱状图
			currHist := currDif - currDea
		
			// 6. 判断交叉
			// 金叉: 上一刻 DIF < DEA 且 当前 DIF > DEA
			if prevDif < prevDea && currDif > currDea {
				crossType = 1 // 金叉
			} else if prevDif > prevDea && currDif < currDea {
				crossType = 2 // 死叉
			} else {
				crossType = 0 // 无新交叉
			}
		
			return currDif, currDea, currHist, crossType
		}
		
		// 辅助函数：计算 EMA 序列 (返回完整数组)
		func calculateEMASeries(klines []Kline, period int) []float64 {
			length := len(klines)
			if length < period {
				return make([]float64, length)
			}
			result := make([]float64, length)
		
			// 初始值使用 SMA
			sum := 0.0
			for i := 0; i < period; i++ {
				sum += klines[i].Close
			}
			result[period-1] = sum / float64(period)
		
			multiplier := 2.0 / float64(period+1)
			for i := period; i < length; i++ {
				result[i] = (klines[i].Close-result[i-1])*multiplier + result[i-1]
			}
			return result
		}
		
		// 辅助函数：根据 float64 数组计算最后一个点的 SMA
		func calculateSMAFromSeries(data []float64, period int) float64 {
			length := len(data)
			if length < period {
				return 0
			}
			sum := 0.0
			for i := length - period; i < length; i++ {
				sum += data[i]
			}
			return sum / float64(period)
		}
