package actions

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//SendPVUploadedToCheckers sends emails to all checkers when a pv is uploaded
func SendPVUploadedToCheckers(uploadedBy User, recipients []User, client, pvType string) {
	to := make([]map[string]string, 0)
	for _, recipient := range recipients {
		to = append(to, map[string]string{"email": recipient.Email, "name": recipient.Fullname})
	}
	sendEmail(to, uploadedBy.Fullname, client, pvType)
}

//sendEmail sends the email by sending a request to the php mail server
func sendEmail(recipients []map[string]string, sender, client, pvType string) {
	values := map[string]interface{}{"recipients": recipients, "type": pvType, "data": map[string]string{"uploader": sender, "client": client}}
	data, err := json.Marshal(values)
	if err != nil {
		log.Println(err.Error())
	}
	mailerURL := "http://localhost:9090/generic-email.php"
	_, err = http.Post(mailerURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Println(err.Error())
	}
	/**
	defer response.Body.Close()

	type APIResponse struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}
	if response.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Println(err.Error())
		}
		var responseJSON APIResponse
		json.Unmarshal(bodyBytes, &responseJSON)
		if responseJSON.Error {
			log.Println(responseJSON.Message)
		}
	}
	**/
}

func stringToFloat64(value string) float64 {
	if value != "" {
		regex := regexp.MustCompile(`(?m)[A-Da-d]|\s|,|[F-Zf-z]`)
		if amount, err := strconv.ParseFloat(regex.ReplaceAllString(value, ""), 64); err == nil {
			return toFixed(amount, 2)
		}
	}
	return 0
}

func stringToFloat32(value string) float32 {
	if value != "" {
		regex := regexp.MustCompile(`(?m)[A-Da-d]|\s|,|[F-Zf-z]`)
		if amount, err := strconv.ParseFloat(regex.ReplaceAllString(value, ""), 32); err == nil {
			return float32(toFixed(amount, 2))
		}
	}
	return 0
}

//getDatesFromString uses regex to get dates from the string
func getDatesFromString(value string) ReportDate {
	date := ReportDate{}
	var regex = regexp.MustCompile(`(?i)\d{1,2}[.,\/]\d{1,2}[.,\/]\d{2,4}`)
	match := regex.FindAllString(value, -1)
	switch len(match) {
	case 0:
		dateWithMonths := getDateWithMonthNameFromString(value)
		switch len(dateWithMonths) {
		case 1:
			date.To = dateWithMonths[0]
			break
		case 2:
			date.From = dateWithMonths[0]
			date.To = dateWithMonths[1]
			break
		}
		break
	case 1:
		dateWithMonths := getDateWithMonthNameFromString(value)
		if len(dateWithMonths) > 0 {
			if dateComesBefore(value, match[0], dateWithMonths[0]) {
				date.From = normalizeDate(match[0])
				date.To = dateWithMonths[0]
			} else {
				date.From = dateWithMonths[0]
				date.To = normalizeDate(match[0])
			}
		} else {
			date.To = normalizeDate(match[0])
		}
		break
	case 2:
		date.From = normalizeDate(match[0])
		date.To = normalizeDate(match[1])
		break
	}
	return date
}

//getDateWithMonthNameFromString returns the date with month names in it
func getDateWithMonthNameFromString(date string) []string {
	var regex = regexp.MustCompile(`(?i)\d{1,2}\s\w{3,}\s\d{2,4}`)
	match := regex.FindAllString(date, -1)
	if len(match) > 0 {
		return match
	}
	return []string{}
}

//normalizeDate creates a unified date format for dates
func normalizeDate(date string) string {
	return strings.ReplaceAll(strings.ReplaceAll(date, ".", "-"), "/", "-")
}

//createPVReportFieldFromRow loads data from the row
func createPVReportFieldFromRow(row []string) (ReportField, error) {
	length := len(row)
	if length < 8 {
		return ReportField{}, errors.New("cannot parse row")
	}
	if length == 8 {
		return ReportField{
			SecurityName:      row[0],
			CDSCode:           row[1],
			ISIN:              row[2],
			SCBCode:           row[3],
			MarketPrice:       stringToFloat32(row[4]),
			NominalValue:      stringToFloat64(row[5]),
			CumulativeCost:    stringToFloat64(row[6]),
			Value:             stringToFloat64(row[7]),
			PercentageOfTotal: 0,
			Dates:             getDatesFromString(row[0]),
		}, nil
	}
	return ReportField{
		SecurityName:      row[0],
		CDSCode:           row[1],
		ISIN:              row[2],
		SCBCode:           row[3],
		MarketPrice:       stringToFloat32(row[4]),
		NominalValue:      stringToFloat64(row[5]),
		CumulativeCost:    stringToFloat64(row[6]),
		Value:             stringToFloat64(row[7]),
		PercentageOfTotal: stringToFloat32(row[8]),
		Dates:             getDatesFromString(row[0]),
	}, nil
}

//dateComesBefore compares the position of two dates
func dateComesBefore(value, date1, date2 string) bool {
	return strings.Index(value, date1) < strings.Index(value, date2)
}

//hasReportDate checks if the string has the report date
func hasReportDate(value string) bool {
	return strings.Contains(strings.ToLower(value), "portfolio valuation report as at")
}

//createSummaryDataFromRow creates a summary data struct from the given row
func createSummaryDataFromRow(row []string) (ReportSummary, error) {
	length := len(row)
	//if length < 5 || length > 6 {
	if length < 5 {
		return ReportSummary{}, errors.New("cannot parse summary")
	}
	if length == 5 { //idk why
		return ReportSummary{
			SecurityName:      row[0],
			NominalValue:      stringToFloat64(row[1]),
			CumulativeCost:    stringToFloat64(row[2]),
			Value:             stringToFloat64(row[3]),
			PercentageOfTotal: stringToFloat32(row[4]),
		}, nil
	}
	cumlativeCostIndex := 2
	if (strings.TrimSpace(row[2]) == "") && strings.TrimSpace(row[3]) != "" {
		cumlativeCostIndex = 3
	}
	return ReportSummary{
		SecurityName:      row[0],
		NominalValue:      stringToFloat64(row[1]),
		CumulativeCost:    stringToFloat64(row[cumlativeCostIndex]),
		Value:             stringToFloat64(row[4]),
		PercentageOfTotal: stringToFloat32(row[5]),
	}, nil
}

//GetTotalLCYValueFromSummary returns the sum total of the values in the pv summary
func GetTotalLCYValueFromSummary(summary []ReportSummary) float64 {
	var total float64
	for _, data := range summary {
		total += data.Value
	}
	return total
}

//ParseDate returns the parsed date from the string passed
func ParseDate(value string) (time.Time, error) {
	var date time.Time
	var err error
	if date, err = time.Parse("02-01-2006", value); err != nil {
		if date, err = time.Parse("02-01-06", value); err != nil {
			if date, err = time.Parse("02 January 2006", value); err != nil {
				if date, err = time.Parse("02 January 06", value); err != nil {
					if date, err = time.Parse("02 Jan 2006", value); err != nil {
						if date, err = time.Parse("02 Jan 06", value); err != nil {
							return time.Now(), err
						}
					}
				}
			}
		}
	}
	return date, nil
}

//getClientIDFromHeader checks if the header value has a client id and returns is
func getClientIDFromHeader(value string) string {
	var regex = regexp.MustCompile(`(?m):[^-]+`)
	result := strings.TrimSpace(regex.FindString(value))
	if len(result) > 1 {
		return strings.TrimSpace(result[1:]) //to remove the matched :
	}
	regex = regexp.MustCompile(`(?m)Client:\s-.+`)
	result = regex.FindString(value)
	if len(result) > 1 {
		return strings.TrimSpace(result[9:]) //cut out the Client: -
	}
	regex = regexp.MustCompile(`(?m)Safekeeping\sAccount:\s-.+`)
	result = regex.FindString(value)
	if len(result) > 1 {
		return strings.TrimSpace(result[22:]) //cut out the Safekeeping Account: -
	}
	return ""
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func NameIsASecurityName(securityNames []MergedSecurityName, name string) bool {
	for _, secName := range securityNames {
		if (secName.Security == name) || strings.Contains(strings.ToLower(name), "recei") {
			return true
		}
	}
	return false
}
