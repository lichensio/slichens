package slichens

import "time"

const MinimumSampleCount = 2
const MinimumSignalLevel = -129.99

type SurveyKey struct {
	Band    int    // 0 all
	CellID  int    // 0 all
	NetName string // all all
}

// Front structure and information about the survey
// and back data of the survey in Survey

type SurveyResult struct {
	SurveyType         string
	FileCreated        time.Time
	IMEINumber         string
	HardwareVersion    string
	ApplicationVersion string
	FirmwareVersion    string
	Filename           string
	Timestamp          int
	Surveys            SurveyMap
}

type KeySlice []SurveyKey

type SurveyData struct {
	Survey     int
	Timestamp  time.Time
	Network    string
	Index      int
	XRFCN      int
	DBM        float64
	Percentage float64
	RSSI       float64
	MCC        int
	MNC        int
	CellID     int
	LACTAC     int
	BandNum    int
	Band       int
	BSIC       string
	SCR        string
	ECIO       string
	RSCP       float64
	PCI        int
	RSRP       float64
	RSRQ       float64
	BW         int
	DL         float64
	UL         float64
	NetName    string
	Signal     string
}

type SurveyDataSlice []SurveyData

type SurveyMap map[SurveyKey]SurveyDataSlice

type SurveyStat struct {
	Number                uint
	RSRPMean              float64
	RSRPMedian            float64
	RSRPMode              float64
	RSRPRange             float64
	RSRPQuartiles         float64
	RSRPMin               float64
	RSRPMax               float64
	RSRPVariance          float64
	RSRPSkewness          float64
	RSRPKurtosis          float64
	RSRPStandardDeviation float64
}

type SurveyStatMap map[SurveyKey]SurveyStat

type SurveySummary struct {
	SurveyType string
	Stat       SurveyStatMap
}

type SurveyTwoSamples struct {
	Number                   uint
	RSRPavOut                float64
	RSRQavOut                float64
	RSRPmaxOut               float64
	RSRQmaxOut               float64
	RSRPminOut               float64
	RSRQminOut               float64
	RSRPStandardDeviationOut float64
	RSRPavIn                 float64
	RSRQavIn                 float64
	RSRPmaxIn                float64
	RSRQmaxIn                float64
	RSRPminIn                float64
	RSRQminIn                float64
	RSRPStandardDeviationIn  float64
	DeltaRSRP                float64
	DeltaRSRQ                float64
	T                        float64
	P                        float64
	Df                       float64
}

type SurveyTwoSamplesMap map[SurveyKey]SurveyTwoSamples

type SurveyTwoSamplesSummary struct {
	SurveyType string
	Data       SurveyTwoSamplesMap
}
