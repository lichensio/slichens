package slichens

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const headerRows = 14

func ReadMultiCSV(filename string) (SurveyResult, error) {
	var survey SurveyResult

	if !isValidSirettaLFilename(filename) {
		return survey, fmt.Errorf("invalid filename pattern")
	}
	survey.Surveys = make(map[SurveyKey]SurveyDataSlice)

	input, err := os.Open(filename)
	if err != nil {
		return survey, err
	}
	defer input.Close()

	reader := csv.NewReader(input)

	// Parsing header rows
	for i := 0; i < headerRows; i++ {
		record, err := reader.Read()
		if err == io.EOF {
			return survey, errors.New("unexpected EOF when reading headers")
		}
		if err != nil {
			return survey, err
		}

		// Extract relevant information based on the current header row.
		switch i {
		case 6:
			survey.SurveyType = record[1]
			// fmt.Println(record[1])
		case 7:
			// fmt.Println(record[1])
			survey.FileCreated, err = time.Parse("02/01/06 15:04:05", record[1])
			// fmt.Println(survey.FileCreated)
			if err != nil {
				return survey, errors.New("invalid date format in FileCreated")
			}
		case 8:
			survey.IMEINumber = record[1]
		case 9:
			survey.HardwareVersion = record[1]
		case 10:
			survey.ApplicationVersion = record[1]
		case 11:
			survey.FirmwareVersion = record[1]
		case 12:
			survey.Filename = record[1]
		}
		// Allow for variable columns per row.
		reader.FieldsPerRecord = 0
	}
	// fmt.Println(survey)
	// Parsing survey data
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return survey, err
		}

		// Handle special case for "Survey:" row or parse survey data.
		switch record[0] {
		case "Survey:":
			// Handle the "Survey:" case if needed. If there's no specific action, this can be omitted.
		default:
			var key SurveyKey
			var surveyData SurveyData
			if surveyData, key, err = parseSurveyData(record); err != nil {
				return survey, err
			}
			survey.Surveys[key] = append(survey.Surveys[key], surveyData)
		}
	}

	return survey, nil
}

// parseSurveyData extracts survey data from a record row and returns the parsed data.
func parseSurveyData(record []string) (SurveyData, SurveyKey, error) {
	var surveyData SurveyData
	var key SurveyKey

	// Extract and convert data from the record.
	// This is just a simple example, in a real-world scenario,
	// error handling and parsing would be much more exhaustive.
	// You can also consider breaking this into smaller functions.
	surveyData.Survey, _ = strconv.Atoi(record[0])
	surveyData.Timestamp, _ = time.Parse("01/02/06 15:04:05", record[1])
	surveyData.Network = record[2]
	surveyData.Index, _ = strconv.Atoi(record[3])
	surveyData.XRFCN, _ = strconv.Atoi(record[4])
	surveyData.DBM, _ = strconv.ParseFloat(record[5], 64)
	surveyData.Percentage, _ = strconv.ParseFloat(record[6], 64)
	surveyData.RSSI, _ = strconv.ParseFloat(record[7], 64)
	surveyData.MCC, _ = strconv.Atoi(record[8])
	surveyData.MNC, _ = strconv.Atoi(record[9])
	key.CellID, _ = strconv.Atoi(record[10])
	surveyData.LACTAC, _ = strconv.Atoi(record[11])
	surveyData.BandNum, _ = strconv.Atoi(record[12])
	frequencyParts := strings.Split(record[13], " ")
	key.Band, _ = strconv.Atoi(frequencyParts[0])
	surveyData.BSIC = record[14]
	surveyData.SCR = record[15]
	surveyData.ECIO = record[16]
	surveyData.RSCP, _ = strconv.ParseFloat(record[17], 64)
	surveyData.PCI, _ = strconv.Atoi(record[18])
	surveyData.RSRP, _ = strconv.ParseFloat(record[19], 64)
	surveyData.RSRQ, _ = strconv.ParseFloat(record[20], 64)
	surveyData.BW, _ = strconv.Atoi(record[21])
	surveyData.DL, _ = strconv.ParseFloat(record[22], 64)
	surveyData.UL, _ = strconv.ParseFloat(record[23], 64)

	key.NetName = record[24]
	return surveyData, key, nil
}

func isValidSirettaLFilename(filename string) bool {
	// The regex pattern:
	// (?i): case-insensitive match
	// ^L: starts with an 'L'
	// \d{7}: followed by 7 digits
	// \.csv$: ends with .csv (case-insensitive due to (?i))
	pattern := `(?i)^L\d{7}\.csv$`
	matched, err := regexp.MatchString(pattern, filename)
	if err != nil {
		// Handle error (e.g., invalid regex pattern)
		return false
	}
	return matched
}
