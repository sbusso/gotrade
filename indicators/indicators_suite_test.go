package indicators_test

import (
	"encoding/csv"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/thetruetrade/gotrade"
	"github.com/thetruetrade/gotrade/feeds"
	"github.com/thetruetrade/gotrade/indicators"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

var (
	csvFeed          *feeds.CSVFileFeed
	sourceData       []float64        = []float64{5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0, 15.0, 16.0, 17.0, 18.0, 19.0, 20.0}
	sourceDOHLCVData []gotrade.DOHLCV = []gotrade.DOHLCV{gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 5.0, 0.0),
		gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 6.0, 0.0),
		gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 7.0, 0.0),
		gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 8.0, 0.0),
		gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 9.0, 0.0),
		gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 10.0, 0.0),
		gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 11.0, 0.0),
		gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 12.0, 0.0),
		gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 13.0, 0.0),
		gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 14.0, 0.0),
		gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 15.0, 0.0),
		gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 16.0, 0.0),
		gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 17.0, 0.0),
		gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 18.0, 0.0),
		gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 19.0, 0.0),
		gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 20.0, 0.0)}
)

func TestIndicators(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Indicators Suite")
}

var _ = BeforeSuite(func() {
	csvFeed = feeds.NewCSVFileFeedWithDOHLCVFormat("../testdata/JSETOPI.2013.data",
		feeds.DashedYearDayMonthDateParserForLocation(time.Local))
})

var _ = AfterSuite(func() {
	csvFeed = nil
})

// ValidFromBar() int
// Length() int
// MinValue() float64
// MaxValue() float64

type IndicatorSharedSpecInputs struct {
	IndicatorUnderTest indicators.Indicator
}

func ShouldBeAnInitialisedIndicator(inputs *IndicatorSharedSpecInputs) {
	It("the indicator should not be valid from any bar yet", func() {
		Expect(inputs.IndicatorUnderTest.ValidFromBar()).To(Equal(-1))
	})

	It("the indicator stream should have no results", func() {
		Expect(inputs.IndicatorUnderTest.Length()).To(BeZero())
	})

	It("the indicator should have no minimum value set", func() {
		Expect(inputs.IndicatorUnderTest.MinValue()).To(Equal(math.MaxFloat64))
	})

	It("the indicator should have no maximum value set", func() {
		Expect(inputs.IndicatorUnderTest.MaxValue()).To(Equal(math.SmallestNonzeroFloat64))
	})
}

type IndicatorWithLookbackSharedSpecInputs struct {
	IndicatorUnderTest indicators.IndicatorWithLookback
	GetMaximum         GetMaximumFunc
	GetMinimum         GetMinimumFunc
	SourceDataLength   int
}

func ShouldBeAnInitialisedIndicatorWithLookback(inputs *IndicatorWithLookbackSharedSpecInputs) {
	It("the indicator should not be valid from any bar yet", func() {
		Expect(inputs.IndicatorUnderTest.(indicators.Indicator).ValidFromBar()).To(Equal(-1))
	})

	It("the indicator stream should have no results", func() {
		Expect(inputs.IndicatorUnderTest.(indicators.Indicator).Length()).To(BeZero())
	})

	It("the indicator should have no minimum value set", func() {
		Expect(inputs.IndicatorUnderTest.(indicators.Indicator).MinValue()).To(Equal(math.MaxFloat64))
	})

	It("the indicator should have no maximum value set", func() {
		Expect(inputs.IndicatorUnderTest.(indicators.Indicator).MaxValue()).To(Equal(math.SmallestNonzeroFloat64))
	})

	It("the indicator should have a valid lookback period", func() {
		Expect(inputs.IndicatorUnderTest.GetLookbackPeriod()).Should(BeNumerically(">=", indicators.MinimumLookbackPeriod))
		Expect(inputs.IndicatorUnderTest.GetLookbackPeriod()).Should(BeNumerically("<=", indicators.MaximumLookbackPeriod))
	})
}

func ShouldBeAnIndicatorThatHasReceivedFewerTicksThanItsLookbackPeriod(inputs *IndicatorWithLookbackSharedSpecInputs) {

	It("the indicator should not be valid from any bar yet", func() {
		Expect(inputs.IndicatorUnderTest.(indicators.Indicator).ValidFromBar()).To(Equal(-1))
	})

	It("the indicator stream should have no results", func() {
		Expect(inputs.IndicatorUnderTest.(indicators.Indicator).Length()).To(BeZero())
	})

	It("the indicator should have no minimum value set", func() {
		Expect(inputs.IndicatorUnderTest.(indicators.Indicator).MinValue()).To(Equal(math.MaxFloat64))
	})

	It("the indicator should have no maximum value set", func() {
		Expect(inputs.IndicatorUnderTest.(indicators.Indicator).MaxValue()).To(Equal(math.SmallestNonzeroFloat64))
	})
}

func ShouldBeAnIndicatorThatHasReceivedTicksEqualToItsLookbackPeriod(inputs *IndicatorWithLookbackSharedSpecInputs) {
	It("the indicator stream should have a single entry", func() {
		Expect(inputs.IndicatorUnderTest.(indicators.Indicator).Length()).To(Equal(1))
	})

	//It("the indicator min and max should be equal", func() {
	//	Expect(inputs.IndicatorUnderTest.(indicators.Indicator).MaxValue()).To(Equal(inputs.IndicatorUnderTest.(indicators.Indicator).MinValue()))
	//})

	It("the indicator should be valid from the lookback period", func() {
		Expect(inputs.IndicatorUnderTest.(indicators.Indicator).ValidFromBar()).To(Equal(inputs.IndicatorUnderTest.GetLookbackPeriod()))
	})
}

func ShouldBeAnIndicatorThatHasReceivedMoreTicksThanItsLookbackPeriod(inputs *IndicatorWithLookbackSharedSpecInputs) {
	It("the indicator stream should have entries equal to the number of ticks less the lookback period", func() {
		Expect(inputs.IndicatorUnderTest.(indicators.Indicator).Length()).To(Equal(inputs.SourceDataLength - (inputs.IndicatorUnderTest.GetLookbackPeriod() - 1)))
	})

	It("the indicator min should equal the result stream minimum", func() {
		Expect(inputs.IndicatorUnderTest.(indicators.Indicator).MinValue()).To(Equal(inputs.GetMinimum()))
	})

	It("the indicator max should equal the result stream maximum", func() {
		Expect(inputs.IndicatorUnderTest.(indicators.Indicator).MaxValue()).To(Equal(inputs.GetMaximum()))
	})
}

func LoadCSVPriceDataFromFile(fileName string) (results []float64, err error) {
	file, err := os.Open("../testdata/" + fileName)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return nil, err
		}

		priceValue, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			fmt.Println("Error:", err)
			return nil, err
		}
		results = append(results, priceValue)
	}
	return results, nil
}

func LoadCSVBollingerPriceDataFromFile(fileName string) (results []indicators.BollingerBand, err error) {
	file, err := os.Open("../testdata/" + fileName)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return nil, err
		}

		upperBandValue, err := strconv.ParseFloat(strings.TrimSpace(record[0]), 64)
		middleBandValue, err := strconv.ParseFloat(strings.TrimSpace(record[1]), 64)
		lowerBandValue, err := strconv.ParseFloat(strings.TrimSpace(record[2]), 64)

		if err != nil {
			fmt.Println("Error:", err)
			return nil, err
		}
		results = append(results, indicators.NewBollingerBandDataItem(upperBandValue, middleBandValue, lowerBandValue))
	}
	return results, nil
}

func LoadCSVMACDPriceDataFromFile(fileName string) (results []indicators.MACDData, err error) {
	file, err := os.Open("../testdata/" + fileName)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return nil, err
		}

		macd, err := strconv.ParseFloat(strings.TrimSpace(record[0]), 64)
		signal, err := strconv.ParseFloat(strings.TrimSpace(record[1]), 64)
		histogram, err := strconv.ParseFloat(strings.TrimSpace(record[2]), 64)

		if err != nil {
			fmt.Println("Error:", err)
			return nil, err
		}
		results = append(results, indicators.NewMACDDataItem(macd, signal, histogram))
	}
	return results, nil
}

type GetMaximumFunc func() float64

type GetMinimumFunc func() float64

func GetDataMax(dohlcvArray []float64) float64 {
	max := math.SmallestNonzeroFloat64

	for i := range dohlcvArray {
		if max < dohlcvArray[i] {
			max = dohlcvArray[i]
		}
	}

	return max
}

func GetDataMin(dohlcvArray []float64) float64 {
	min := math.MaxFloat64

	for i := range dohlcvArray {
		if min > dohlcvArray[i] {
			min = dohlcvArray[i]
		}
	}

	return min
}

func GetDataMaxDOHLCV(dohlcvArray []gotrade.DOHLCV, selectData gotrade.DataSelectionFunc) float64 {
	max := math.SmallestNonzeroFloat64

	for i := range dohlcvArray {
		var selectedData = selectData(dohlcvArray[i])
		if max < selectedData {
			max = selectedData
		}
	}

	return max
}

func GetDataMinDOHLCV(dohlcvArray []gotrade.DOHLCV, selectData gotrade.DataSelectionFunc) float64 {
	min := math.MaxFloat64

	for i := range dohlcvArray {
		var selectedData = selectData(dohlcvArray[i])
		if min > selectedData {
			min = selectedData
		}
	}

	return min
}

func GetDataMaxBollinger(dohlcvArray []indicators.BollingerBand, selectData indicators.BollingerDataSelectionFunc) float64 {
	max := math.SmallestNonzeroFloat64

	for i := range dohlcvArray {
		var selectedData = selectData(dohlcvArray[i])
		if max < selectedData {
			max = selectedData
		}
	}

	return max
}

func GetDataMinBollinger(dohlcvArray []indicators.BollingerBand, selectData indicators.BollingerDataSelectionFunc) float64 {
	min := math.MaxFloat64

	for i := range dohlcvArray {
		var selectedData = selectData(dohlcvArray[i])
		if min > selectedData {
			min = selectedData
		}
	}

	return min
}

func GetDataMaxMACD(dohlcvArray []indicators.MACDData) float64 {
	max := math.SmallestNonzeroFloat64

	for i := range dohlcvArray {
		macd := dohlcvArray[i].M()
		signal := dohlcvArray[i].S()
		histogram := dohlcvArray[i].H()

		if max < macd {
			max = macd
		}
		if max < signal {
			max = signal
		}
		if max < histogram {
			max = histogram
		}
	}

	return max
}

func GetDataMinMACD(dohlcvArray []indicators.MACDData) float64 {
	min := math.MaxFloat64

	for i := range dohlcvArray {
		macd := dohlcvArray[i].M()
		signal := dohlcvArray[i].S()
		histogram := dohlcvArray[i].H()

		if min > macd {
			min = macd
		}
		if min > signal {
			min = signal
		}
		if min > histogram {
			min = histogram
		}
	}

	return min
}
