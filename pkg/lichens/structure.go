package lichens

import (
	"time"
)

const MinimumSampleCount = 2
const MinimumSignalLevel = -129.99

type GSMAType int64

const (
	GSM GSMAType = iota
	ThreeG
	FourG
	FiveG
)

type DeltaType string

const (
	IndoorOutdoor DeltaType = "IndoorOutdoor"
	IndoorBooster DeltaType = "IndoorBooster"
)

type SurveyDeltaStatsSummary struct {
	SurveyType string
	DeltaStats SurveyDeltaMap
	DeltaType  DeltaType
	Min        float64
	Max        float64
}

type SurveyDeltaMap map[SurveyKey]SurveyDeltaStats

type SurveyKey struct {
	Band        int    // 0 all
	CellID      int    // 0 all
	NetName     string //
	NetworkType string
}

// Front structure and information about the survey
// and back data of the survey in Survey

type SurveyInfo struct {
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

type SurveyMap map[SurveyKey]SurveyDataSlice

type SurveyDataSlice []SurveyData

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

type SurveySummary struct {
	SurveyType string
	Stat       SurveyStatsMap
	Min        float64
	Max        float64
}

type Stats struct {
	Number            uint
	Mean              float64
	Median            float64
	Mode              float64
	Range             float64
	Quartiles         float64
	Min               float64
	Max               float64
	Variance          float64
	Skewness          float64
	Kurtosis          float64
	StandardDeviation float64
}
type SurveyStats map[string]Stats

type SurveyStatsMap map[SurveyKey]SurveyStats

// Delta statistics

type DeltaStats struct {
	Number1                uint
	Number2                uint
	CorrelationCoefficient float64
	TTestValue             float64
	PValue                 float64
	AreSignificantlyDiff   bool
	Alpha                  float64
	Delta                  float64
}

type SurveyDeltaStats map[string]DeltaStats

type SurveyDeltaStatsMap map[SurveyKey]SurveyDeltaStats
