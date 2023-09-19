package lichens

import (
	"github.com/lichensio/slichens/pkg/statistics"
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

type SurveyDeltaMap map[SurveyKey]statistics.SurveyDeltaStats

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
	Stat       statistics.SurveyStatsMap
	Min        float64
	Max        float64
}
