package actions

import (
	"bytes"
	cryptoRand "crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gobuffalo/envy"

	"encoding/base64"

	"github.com/jinzhu/now"
)

// ClientDetailsArrayToString converts a client details array to a key-pair value
func ClientDetailsArrayToString(client Client) map[string]interface{} {
	keyPair := make(map[string]interface{}, 0)
	keyPair["id"] = client.ID
	keyPair["name"] = client.Name
	keyPair["account_number"] = client.AccountNumber
	keyPair["image"] = client.Image
	lumpedSummary := make([]map[string]interface{}, 0)
	//@TODO: Change quarter date
	if len(client.PVReportSummary) > 0 {
		keyPair["report_date"] = client.PVReportSummary[0].ReportDate.Format("2006-01-02")
	}
	var nominalValueTotal float64
	var valueLCYTotal float64
	headers := LoadPVReportHeadings()
	var receivablesTotalLCY float64
	var receivablesNorminalValue float64
	var receivablesCumulativeCost float64
	filteredSummary := make(map[string]*ClientPVReportSummary)
	for _, data := range client.PVReportSummary {
		nominalValueTotal += data.NominalValue
		valueLCYTotal += data.LCYAmount
		if filteredSummary[data.SecurityType] == nil {
			filteredSummary[data.SecurityType] = &ClientPVReportSummary{
				ID:                data.ID,
				AccountNumber:     data.AccountNumber,
				BPID:              data.BPID,
				ReportDate:        data.ReportDate,
				SecurityType:      data.SecurityType,
				NominalValue:      data.NominalValue,
				CumulativeCost:    data.CumulativeCost,
				LCYAmount:         data.LCYAmount,
				PercentageOfTotal: data.PercentageOfTotal,
			}
		} else {
			filteredSummary[data.SecurityType].NominalValue += data.NominalValue
			filteredSummary[data.SecurityType].CumulativeCost += data.CumulativeCost
			filteredSummary[data.SecurityType].PercentageOfTotal += data.PercentageOfTotal
			filteredSummary[data.SecurityType].LCYAmount += data.LCYAmount
		}
	}
	for _, data := range filteredSummary {
		hasHeader := false
		for _, header := range headers {
			if strings.ToUpper(data.SecurityType) == header.RealName || strings.ToUpper(data.SecurityType) == header.Heading {
				hasHeader = true
				break
			}
		}
		if hasHeader {
			lumpedSummary = append(lumpedSummary, map[string]interface{}{
				"security_type":    data.SecurityType,
				"percentage_total": ToTwoDP((data.LCYAmount / valueLCYTotal) * 100),
				"nominal_value":    FormatWithComma(data.NominalValue, 2),
				"cumulative_cost":  FormatWithComma(data.CumulativeCost, 2),
				"lcy_amount":       FormatWithComma(data.LCYAmount, 2),
				"reports":          getSummaryReportsFromClientBasedOnSecurityType(client, data.SecurityType),
			})
		} else {
			receivablesTotalLCY += data.LCYAmount
			receivablesNorminalValue += data.NominalValue
			receivablesCumulativeCost += data.CumulativeCost
		}
	}
	if receivablesTotalLCY > 0 {
		lumpedSummary = append(lumpedSummary, map[string]interface{}{
			"security_type":    "RECEIVABLES",
			"percentage_total": ToTwoDP((receivablesTotalLCY / valueLCYTotal) * 100),
			"nominal_value":    FormatWithComma(receivablesNorminalValue, 2),
			"cumulative_cost":  FormatWithComma(receivablesCumulativeCost, 2),
			"lcy_amount":       FormatWithComma(receivablesTotalLCY, 2),
		})
	}
	keyPair["nominal_value_total"] = FormatWithComma(nominalValueTotal, 2)
	keyPair["value_lcy_total"] = FormatWithComma(valueLCYTotal, 2)
	keyPair["percentage_total"] = 100
	keyPair["pv_lumped_summary"] = SortLumpedPVSummary(lumpedSummary)
	mDate := time.Now().Format("January 2006")
	keyPair["meeting_date"] = mDate
	return keyPair
}

func getSummaryReportsFromClientBasedOnSecurityType(client Client, security string) []PVReportField {
	reports := make([]PVReportField, 0)
	for _, report := range client.Reports {
		if report.SecurityType == security {
			reports = append(reports, report)
		}
	}
	return reports
}

// ToTwoDP converts a number to 2 dp
func ToTwoDP(number interface{}) string {
	switch num := number.(type) {
	case float32:
		return fmt.Sprintf("%.2f", num)
	case float64:
		return fmt.Sprintf("%.2f", num)
	}
	return fmt.Sprintf("%v", number)
}

// NormalizeWindowsPath normalizes the windows paths by making \ \\
func NormalizeWindowsPath(path string) string {
	return strings.ReplaceAll(path, "/", "\\\\")
}

// reverseBytes reverses a set of given bytes
func reverseBytes(input []byte) []byte {
	if len(input) == 0 {
		return input
	}
	return append(reverseBytes(input[1:]), input[0])
}

// Comma adds comma to a string and formats it as if it was money
func Comma(value string) string {
	values := strings.Split(value, ".")
	var money string
	var isNegative bool
	var isFloat bool
	decimal := "00"
	money = values[0]
	if money[0] == '-' {
		isNegative = true
		money = money[1:]
	}
	if len(values) > 1 {
		decimal = values[1]
		isFloat = true
	}
	output := make([]byte, 0)
	length := len(money) - 1
	for i := 0; i < length+1; i++ {
		if i > 0 && (i%3 == 0) {
			output = append(output, ',')
		}
		output = append(output, money[length-i])
	}
	if isNegative {
		output = append(output, '-')
	}
	if isFloat {
		return fmt.Sprintf("%s.%s", reverseBytes(output), decimal)
	}
	return fmt.Sprintf("%s", reverseBytes(output))
}

// FormatWithComma formats the number passed to the function
func FormatWithComma(number interface{}, dp int) string {
	//TODO: abs for int
	isNegative := false
	switch value := number.(type) {
	case float32:
		val := value
		total := math.RoundToEven(float64(val))
		if total < 0 {
			isNegative = true
		}
		if isNegative {
			return "(" + Comma(fmt.Sprintf("%."+strconv.Itoa(dp)+"f", val))[1:] + ")"
		}
		return Comma(fmt.Sprintf("%."+strconv.Itoa(dp)+"f", val))
	case float64:
		if value < 0 {
			isNegative = true
		}
		if isNegative {
			return "(" + Comma(fmt.Sprintf("%."+strconv.Itoa(dp)+"f", value))[1:] + ")"
		}
		return Comma(fmt.Sprintf("%."+strconv.Itoa(dp)+"f", value))
	case int:
		if value < 0 {
			isNegative = true
		}
		if isNegative {
			return "(" + Comma(fmt.Sprintf("%d", value))[1:] + ")"
		}
		return Comma(fmt.Sprintf("%d", value))
	case int8:
		if value < 0 {
			isNegative = true
		}
		if isNegative {
			return "(" + Comma(fmt.Sprintf("%d", value))[1:] + ")"
		}
		return Comma(fmt.Sprintf("%d", value))
	case int16:
		if value < 0 {
			isNegative = true
		}
		if isNegative {
			return "(" + Comma(fmt.Sprintf("%d", value))[1:] + ")"
		}
		return Comma(fmt.Sprintf("%d", value))
	case int32:
		if value < 0 {
			isNegative = true
		}
		if isNegative {
			return "(" + Comma(fmt.Sprintf("%d", value))[1:] + ")"
		}
		return Comma(fmt.Sprintf("%d", value))
	case int64:
		if value < 0 {
			isNegative = true
		}
		if isNegative {
			return "(" + Comma(fmt.Sprintf("%d", value))[1:] + ")"
		}
		return Comma(fmt.Sprintf("%d", value))
	default:
		return "0.00"
	}
}

// GetQuarterDatesBetween generates the quarters between the years given
func GetQuarterDatesBetween(yearFrom, yearTo int) []QuarterDate {
	dates := make([]QuarterDate, 0)
	for i := yearFrom; i <= yearTo; i++ {
		dates = append(dates, QuarterDate{
			Begin: fmt.Sprintf("%d-01-01", i),
			End:   fmt.Sprintf("%d-03-31", i),
		})
		dates = append(dates, QuarterDate{
			Begin: fmt.Sprintf("%d-04-01", i),
			End:   fmt.Sprintf("%d-06-30", i),
		})
		dates = append(dates, QuarterDate{
			Begin: fmt.Sprintf("%d-07-01", i),
			End:   fmt.Sprintf("%d-09-30", i),
		})
		dates = append(dates, QuarterDate{
			Begin: fmt.Sprintf("%d-10-01", i),
			End:   fmt.Sprintf("%d-12-31", i),
		})
	}
	return dates
}

// GetQuarterDate returns the date for the current quarter @@Changed to use previous quarter as te current
func GetQuarterDate() string {
	today := time.Now()
	month := int(today.Month())
	year := today.Year()
	if month >= 1 && month <= 3 {
		return fmt.Sprintf(`31st December %d`, year-1)
	} else if month >= 4 && month <= 6 {
		return fmt.Sprintf(`31st March %d`, year)
	} else if month >= 7 && month <= 9 {
		return fmt.Sprintf(`30th June %d`, year)
	} else {
		return fmt.Sprintf(`30th September %d`, year)
	}
}

// GetShortQuarterDate returns the date for the current quarter
func GetShortQuarterDate() string {
	today := time.Now()
	month := int(today.Month())
	year := today.Year()
	if month >= 1 && month <= 3 {
		return fmt.Sprintf(`31st Dec %d`, year-1)
	} else if month >= 4 && month <= 6 {
		return fmt.Sprintf(`31st Mar %d`, year)
	} else if month >= 7 && month <= 9 {
		return fmt.Sprintf(`30th Jun %d`, year)
	} else {
		return fmt.Sprintf(`30th Sep %d`, year)
	}
}

// GetShortPreviousQuarterDate returns the date for the previous quarter
func GetShortPreviousQuarterDate() string {
	currentQuarter, _ := time.Parse("2006-01-02", GetQuarterFormalDate())
	date := now.New(currentQuarter.AddDate(0, -3, -5)).EndOfMonth()
	month := int(date.Month())
	year := date.Year()
	if month >= 1 && month <= 3 {
		return fmt.Sprintf(`31st Mar %d`, year)
	} else if month >= 4 && month <= 6 {
		return fmt.Sprintf(`30th Jun %d`, year)
	} else if month >= 7 && month <= 9 {
		return fmt.Sprintf(`30th Sep %d`, year)
	} else {
		return fmt.Sprintf(`31st Dec %d`, year)
	}
}

// GetQuarterFormalDate get a formal date
func GetQuarterFormalDate() string {
	today := time.Now()
	month := int(today.Month())
	year := today.Year()
	if month >= 1 && month <= 3 {
		return fmt.Sprintf(`%d-12-31`, year-1)
	} else if month >= 4 && month <= 6 {
		return fmt.Sprintf(`%d-03-31`, year)
	} else if month >= 7 && month <= 9 {
		return fmt.Sprintf(`%d-06-30`, year)
	} else {
		return fmt.Sprintf(`%d-09-30`, year)
	}
}

// MakeQuarterFormalDate make a formal date
func MakeQuarterFormalDate(quarter, year string) string {
	if quarter == "1" {
		return fmt.Sprintf(`%s-03-31`, year)
	} else if quarter == "2" {
		return fmt.Sprintf(`%s-06-30`, year)
	} else if quarter == "3" {
		return fmt.Sprintf(`%s-09-30`, year)
	} else {
		return fmt.Sprintf(`%s-12-31`, year)
	}
}

// MakeQuarterDate creates a quarter date given the quarter and the year
func MakeQuarterDate(quarter, year string) string {
	if quarter == "1" {
		return fmt.Sprintf(`31st March %s`, year)
	} else if quarter == "2" {
		return fmt.Sprintf(`30th June %s`, year)
	} else if quarter == "3" {
		return fmt.Sprintf(`30th September %s`, year)
	} else {
		return fmt.Sprintf(`31st December %s`, year)
	}
}

func MakeSecLastLicenseRenewalDate(quarter, year string) string {
	qInt, _ := strconv.Atoi(quarter)
	yInt, _ := strconv.Atoi(year)
	if qInt > 2 {
		return fmt.Sprintf("1st July %d", yInt+1)
	}
	return fmt.Sprintf("1st July %d", yInt)
}

// RandomBytes generates a random bytes
func RandomBytes(n int) []byte {
	var letter = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]byte, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return b
}

// MakeLastQuarterFormalDate makes last quarters date from this quarter
func MakeLastQuarterFormalDate(currentQuarter string) string {
	date, err := time.Parse("2006-01-02", currentQuarter)
	if err != nil {
		panic(err)
	}
	return now.New(date.AddDate(0, -3, -5)).EndOfMonth().Format("2006-01-02") //-5 to be certain its in the 3rd prev month
}

// GetLast4QuarterDates returns the quarter dates for the past 4 months
func GetLast4QuarterDates() []time.Time {
	return MakeLast4QuarterDates(GetQuarterFormalDate())
}

// MakeLast4QuarterDates returns the quarter dates for the past 4 months
func MakeLast4QuarterDates(date string) []time.Time {
	dates := make([]time.Time, 4)
	quarter, err := time.Parse("2006-01-02", date)
	if err != nil {
		panic(err)
	}
	dates[0] = now.New(quarter.AddDate(0, -9, -5)).EndOfMonth()
	dates[1] = now.New(quarter.AddDate(0, -6, -5)).EndOfMonth()
	dates[2] = now.New(quarter.AddDate(0, -3, -5)).EndOfMonth()
	dates[3] = now.New(quarter).EndOfMonth()
	return dates
}

// MakeOverviewQuarterDate generates an overview date
func MakeOverviewQuarterDate(date time.Time) string {
	month := int(date.Month())
	year := date.Year()
	if month >= 1 && month <= 3 {
		return fmt.Sprintf(`Q1 - %d`, year)
	} else if month >= 4 && month <= 6 {
		return fmt.Sprintf(`Q2 - %d`, year)
	} else if month >= 7 && month <= 9 {
		return fmt.Sprintf(`Q3 - %d`, year)
	} else {
		return fmt.Sprintf(`Q4 - %d`, year)
	}
}

// FloatArraySum sums the values in the array
func FloatArraySum(array []float64) float64 {
	var total float64
	for _, value := range array {
		total += value
	}
	return total
}

// CalculatePercentageDifference calculates %d in values
func CalculatePercentageDifference(value1, value2 float64) int {
	if value1 > value2 {
		diff := math.Abs(value2 - value1)
		if value1 != 0 {
			return int((diff / value1) * 100)
		}
		return 0
	}
	diff := math.Abs(value1 - value2)
	if value2 != 0 {
		return int((diff / value2) * 100)
	}
	return 0
}

// GetQuarterNumber get the quarter number from the current
func GetQuarterNumber(date time.Time) int {
	month := int(date.Month())
	if month >= 1 && month <= 3 {
		return 4
	} else if month >= 4 && month <= 6 {
		return 1
	} else if month >= 7 && month <= 9 {
		return 2
	} else {
		return 3
	}
}

func GetYearFromQuarter(quarter int) int {
	if quarter == 4 {
		return time.Now().Year() - 1
	}
	return time.Now().Year()
}

// MakeQuarterFancyDate returns the quarter's fancy date
func MakeQuarterFancyDate(date time.Time) string {
	quarterNumber := GetQuarterNumber(date)
	switch quarterNumber {
	case 1:
		return fmt.Sprintf("First Quarter, %d", date.Year())
	case 2:
		return fmt.Sprintf("Second Quarter, %d", date.Year())
	case 3:
		return fmt.Sprintf("Third Quarter, %d", date.Year())
	default:
		return fmt.Sprintf("Fourth Quarter, %d", date.Year()-1)
	}
}

func ExportSecReport(quarterName, lastLicenseRenewalDate string, report GovernanceInfo, schemeDetails []SchemeDetails, info OtherInformation, remarks OfficialReportRemarks, offshoreClients []OffshoreClient) ([]byte, error) {
	type ReportData struct {
		QuarterName            string                `json:"quarter_name"`
		LastLicenseRenewalDate string                `json:"last_license_renewal_date"`
		Report                 GovernanceInfo        `json:"report"`
		SchemeDetails          []SchemeDetails       `json:"scheme_details"`
		Info                   OtherInformation      `json:"info"`
		Remarks                OfficialReportRemarks `json:"remarks"`
		OffshoreClients        []OffshoreClient      `json:"offshore_clients"`
	}
	data, err := json.Marshal(ReportData{
		QuarterName:            quarterName,
		LastLicenseRenewalDate: lastLicenseRenewalDate,
		Report:                 report,
		SchemeDetails:          schemeDetails,
		Info:                   info,
		Remarks:                remarks,
		OffshoreClients:        offshoreClients,
	})
	mailerURL := envy.Get("SEC_REPORT_EXPORT_URL", "http://localhost:9090/sec-report.php")
	response, err := http.Post(mailerURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(response.Body)
		data, err := base64.StdEncoding.DecodeString(string(bodyBytes))
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	return nil, errors.New("error, couldn't convert data to word")
}

func ExportSecCoverLetter() ([]byte, error) {
	data, err := json.Marshal(LoadMailSenderInfo())
	if err != nil {
		log.Println(err.Error())
	}
	mailerURL := envy.Get("SEC_REPORT_EXPORT_URL", "http://localhost:9090/sec-cover-letter.php")
	response, err := http.Post(mailerURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(response.Body)
		data, err := base64.StdEncoding.DecodeString(string(bodyBytes))
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	return nil, errors.New("error, couldn't convert data to word")
}

func IsValidQuarterDayToSendEmails() bool {
	month := time.Now().Month()
	return (month == 3) || (month == 6) || (month == 9) || (month == 12)
}

func ScheduleEmailsToBeSent() {
	if IsValidQuarterDayToSendEmails() {
		day := GetWorkingDayFrom(time.Now()).Add(time.Minute * 3)
		time.AfterFunc(day.Sub(time.Now()), func() {
			SendSecLetterToClients("System")
		})
	}
}

func SendSecLetterToClients(sentBy string) {
	clients := LoadAllLocalClients()
	sender := LoadMailSenderInfo()
	type LetterToClientClient struct {
		Name         string   `json:"name"`
		Emails       []string `json:"emails"`
		AddressLine1 string   `json:"address1"`
		AddressLine2 string   `json:"address2"`
		AddressLine3 string   `json:"address3"`
		AddressLine4 string   `json:"address4"`
		CCList       []string `json:"cc_list"`
	}
	type LetterToClientsData struct {
		Clients  []LetterToClientClient `json:"clients"`
		Sender   MailSenderInfo         `json:"sender"`
		Holidays []string               `json:"holidays"`
	}
	recipients := make([]LetterToClientClient, 0)
	copies := LoadMailService()
	ccEmails := make([]string, 0)
	for _, cc := range copies {
		ccEmails = append(ccEmails, cc.Email)
	}
	for _, client := range clients {
		emails := GetClientEmailsByBPID(client.BPID)
		if len(emails) > 0 {
			clientCCEmails := make([]string, 0)
			for _, email := range emails {
				clientCCEmails = append(clientCCEmails, email.Email)
			}
			recipients = append(recipients, LetterToClientClient{
				Name:         client.Name,
				Emails:       clientCCEmails,
				AddressLine1: client.AddressLine1,
				AddressLine2: client.AddressLine2,
				AddressLine3: client.AddressLine3,
				AddressLine4: client.AddressLine4,
				CCList:       ccEmails,
			})
		}
	}
	values := LetterToClientsData{
		Sender:   sender,
		Clients:  recipients,
		Holidays: GetHolidays(),
	}
	data, err := json.Marshal(values)
	if err != nil {
		log.Println(err.Error())
	}
	mailerURL := "http://localhost:9090/sec-letter-to-client.php"
	_, err = http.Post(mailerURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Println(err.Error())
	}
	LogClientLettersSent(sentBy)
}

func MakeReportableSecuritiesQueryFromHeadings(headings []PVReportHeadings) string {
	var query string
	for _, heading := range headings {
		query += "'" + heading.RealName + "',"
	}
	return "(" + strings.TrimRight(query, ",") + ")"
}

// TODO: make everything standard
func MakeDataForEquityCalculations() string {
	return "('EQUITY SHARE','EQUITY SHARES','GLOBAL EQUITIES','GLOBAL EQUITY','PREFERRED STOCK','PREFERRED STOCKS')"
}

func MakeDataForTotalFixedIncomeInvestmentCalculations() string {
	return "('CORPORATE BOND','CORPORATE BONDS','GOVERNMENT BOND','GOVERNMENT BONDS','FIXED DEPOSIT','FIXED DEPOSITS','TREASURY BILL','TREASURY BILLS','TREASURY NOTE','TREASURY NOTES')"
}

func MakeDataForCapitalInvestmentCalculations() string {
	return "('EQUITY SHARE','EQUITY SHARES','GLOBAL EQUITIES','GLOBAL EQUITY','PREFERRED STOCK','PREFERRED STOCKS','CORPORATE BOND','CORPORATE BONDS','GOVERNMENT BOND','GOVERNMENT BONDS','TREASURY BILL','TREASURY BILLS','TREASURY NOTE','TREASURY NOTES')"
}

func MakeWhereQueryDataUsingClientCodeFromClients(clients []ClientMergedBPAndSCA) string {
	var query string
	for _, client := range clients {
		for _, each := range client.SCA {
			query += "'" + each + "',"
		}
	}
	return "(" + strings.TrimRight(query, ",") + ")"
}

func MakeDBQueryableCurrentQuarterDate() time.Time {
	today := time.Now()
	month := int(today.Month())
	year := today.Year()
	var dateStr string
	if month >= 1 && month <= 3 {
		dateStr = fmt.Sprintf(`%d-12-31`, year-1)
	} else if month >= 4 && month <= 6 {
		dateStr = fmt.Sprintf(`%d-03-31`, year)
	} else if month >= 7 && month <= 9 {
		dateStr = fmt.Sprintf(`%d-06-30`, year)
	} else {
		dateStr = fmt.Sprintf(`%d-09-30`, year)
	}
	date, _ := time.Parse("2006-01-02", dateStr)
	return date
}

func MakeDBQueryableQuarterFirstDate(quarter, year string) time.Time {
	var date time.Time
	if quarter == "1" {
		date, _ = time.Parse("2006-01-02", year+"-01-01")
	} else if quarter == "2" {
		date, _ = time.Parse("2006-01-02", year+"-04-01")
	} else if quarter == "3" {
		date, _ = time.Parse("2006-01-02", year+"-07-01")
	} else {
		date, _ = time.Parse("2006-01-02", year+"-10-01")
	}
	return date
}

func MakeDBQueryableQuarterLastDate(quarter, year string) time.Time {
	var date time.Time
	if quarter == "1" {
		date, _ = time.Parse("2006-01-02", year+"-03-31")
	} else if quarter == "2" {
		date, _ = time.Parse("2006-01-02", year+"-06-30")
	} else if quarter == "3" {
		date, _ = time.Parse("2006-01-02", year+"-09-30")
	} else {
		date, _ = time.Parse("2006-01-02", year+"-12-31")
	}
	return date
}

func CreateFileFromBase64Data(data string, filename string) (string, error) {
	decodedData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	uploadPath := NormalizeWindowsPath(envy.Get("FILE_MANAGER_DIR", "C:\\CRS\\uploads"))
	uploadedFile, err := os.Create(filepath.Join(uploadPath, filename))

	if err != nil {
		return "", err
	}
	defer uploadedFile.Close()
	_, err = uploadedFile.Write(decodedData)
	if err != nil {
		return "", err
	}
	return filename, nil
}

func ConvertFileToBase64(filename string) string {
	uploadPath := NormalizeWindowsPath(envy.Get("FILE_MANAGER_DIR", "C:\\CRS\\uploads"))
	data, _ := ioutil.ReadFile(filepath.Join(uploadPath, filename))
	return base64.StdEncoding.EncodeToString(data)
}

func SendReportActionEmail(emailType string, users []User, title, report, comments string) {
	recipients := make([]map[string]string, 0)
	for _, user := range users {
		recipients = append(recipients, map[string]string{"email": user.Email, "name": user.Fullname})
	}
	values := map[string]interface{}{"recipients": recipients, "type": emailType, "data": map[string]string{"title": title, "report": report, "comment": comments}}

	data, err := json.Marshal(values)
	if err != nil {
		log.Println(err.Error())
	}
	mailerURL := "http://localhost:9090/generic-email.php"
	_, err = http.Post(mailerURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Println(err.Error())
	}
}

func LumpPVSummary(summaryData []ClientPVReportSummary) []ClientPVReportSummary {
	filteredSummary := make(map[string]*ClientPVReportSummary)
	summary := make([]ClientPVReportSummary, 0)
	headers := LoadPVReportHeadings()
	var totalLCYAmount float64
	for _, data := range summaryData {
		if filteredSummary[data.SecurityType] == nil {
			filteredSummary[data.SecurityType] = &ClientPVReportSummary{
				ID:                data.ID,
				AccountNumber:     data.AccountNumber,
				BPID:              data.BPID,
				ReportDate:        data.ReportDate,
				SecurityType:      data.SecurityType,
				NominalValue:      data.NominalValue,
				CumulativeCost:    data.CumulativeCost,
				LCYAmount:         data.LCYAmount,
				PercentageOfTotal: data.PercentageOfTotal,
			}
		} else {
			filteredSummary[data.SecurityType].NominalValue += data.NominalValue
			filteredSummary[data.SecurityType].CumulativeCost += data.CumulativeCost
			filteredSummary[data.SecurityType].PercentageOfTotal += data.PercentageOfTotal
			filteredSummary[data.SecurityType].LCYAmount += data.LCYAmount
		}
		totalLCYAmount += data.LCYAmount
	}
	if len(filteredSummary) > 0 {
		lummpedSummary := ClientPVReportSummary{
			SecurityType: "RECEIVABLES",
		}
		for _, data := range filteredSummary {
			hasHeader := false
			for _, header := range headers {
				if strings.ToUpper(data.SecurityType) == header.RealName || strings.ToUpper(data.SecurityType) == header.Heading {
					hasHeader = true
					break
				}
			}
			if hasHeader {
				if data.LCYAmount == 0 {
					data.PercentageOfTotal = 0
				} else {
					data.PercentageOfTotal = float32((data.LCYAmount / totalLCYAmount) * 100)
				}
				summary = append(summary, *data)
			} else {
				lummpedSummary.PercentageOfTotal += data.PercentageOfTotal
				lummpedSummary.LCYAmount += data.LCYAmount
				lummpedSummary.NominalValue += data.NominalValue
				lummpedSummary.CumulativeCost += data.CumulativeCost
			}
		}
		if lummpedSummary.LCYAmount != 0 && totalLCYAmount != 0 {
			lummpedSummary.PercentageOfTotal = float32((lummpedSummary.LCYAmount / totalLCYAmount) * 100)
			summary = append(summary, lummpedSummary)
		}
	}
	return summary
}

func LumpTrusteeQuarterlyPerformance(performanceData []TrusteePVPerformance) []TrusteePVPerformance {
	filteredPerformance := make(map[string]*TrusteePVPerformance)
	performance := make([]TrusteePVPerformance, 0)
	headers := LoadPVReportHeadings()
	for _, data := range performanceData {
		if filteredPerformance[data.Bond] == nil {
			filteredPerformance[data.Bond] = &TrusteePVPerformance{
				Bond:            data.Bond,
				CurrentQuarter:  data.CurrentQuarter,
				PreviousQuarter: data.PreviousQuarter,
			}
		} else {
			filteredPerformance[data.Bond].CurrentQuarter += data.CurrentQuarter
			filteredPerformance[data.Bond].PreviousQuarter += data.PreviousQuarter
		}
	}
	lummpedPerformance := TrusteePVPerformance{
		Bond: "RECEIVABLES",
	}
	for _, data := range filteredPerformance {
		hasHeader := false
		for _, header := range headers {
			if strings.ToUpper(data.Bond) == header.RealName || strings.ToUpper(data.Bond) == header.Heading {
				hasHeader = true
				break
			}
		}
		if hasHeader {
			performance = append(performance, *data)
		} else {
			lummpedPerformance.CurrentQuarter += data.CurrentQuarter
			lummpedPerformance.PreviousQuarter += data.PreviousQuarter
		}
	}
	performance = append(performance, lummpedPerformance)
	return performance
}

func MakeQuarterMonthAndYear(quarter, year string) string {
	if quarter == "1" {
		return "March " + year

	} else if quarter == "2" {
		return "June " + year
	} else if quarter == "3" {
		return "September " + year
	} else {
		return "December " + year
	}
}

func ExportVarianceToExcel(previousQuarterShortDate, currentQuarterShortDate string, data []VarianceData) (*bytes.Buffer, error) {
	excel := excelize.NewFile()
	index := excel.NewSheet("Sheet1")

	excel.SetCellValue("Sheet1", "A1", "Client")
	excel.SetCellValue("Sheet1", "B1", "Home Country")
	excel.SetCellValue("Sheet1", "C1", fmt.Sprintf("AUA As At %s", previousQuarterShortDate))
	excel.SetCellValue("Sheet1", "D1", fmt.Sprintf("AUA As At %s", currentQuarterShortDate))
	excel.SetCellValue("Sheet1", "E1", "Amount")
	excel.SetCellValue("Sheet1", "F1", "Variance %")
	excel.SetCellValue("Sheet1", "G1", "Remarks")

	for index, datum := range data {
		excel.SetCellValue("Sheet1", fmt.Sprintf("A%d", index+2), datum.Name)
		excel.SetCellValue("Sheet1", fmt.Sprintf("B%d", index+2), datum.Country)
		excel.SetCellValue("Sheet1", fmt.Sprintf("C%d", index+2), FormatWithComma(datum.LastAUA, 2))
		excel.SetCellValue("Sheet1", fmt.Sprintf("D%d", index+2), FormatWithComma(datum.CurrentAUA, 2))
		excel.SetCellValue("Sheet1", fmt.Sprintf("E%d", index+2), replaceNegativeWithBraces(datum.Amount))
		excel.SetCellValue("Sheet1", fmt.Sprintf("F%d", index+2), replaceNegativeWithBraces(datum.Variance))
		excel.SetCellValue("Sheet1", fmt.Sprintf("G%d", index+2), datum.Remarks)
	}
	excel.SetActiveSheet(index)
	return excel.WriteToBuffer()
}

func ExportOutstandingFDsToExcel(certificates []NPRAOutstandingFDCertificate) (*bytes.Buffer, error) {
	excel := excelize.NewFile()
	index := excel.NewSheet("Sheet1")
	excel.SetCellValue("Sheet1", "A1", "Fund Manager")
	excel.SetCellValue("Sheet1", "B1", "Client Name")
	excel.SetCellValue("Sheet1", "C1", "Amount(GHS)")
	excel.SetCellValue("Sheet1", "D1", "Issuer")
	excel.SetCellValue("Sheet1", "E1", "Rate %")
	excel.SetCellValue("Sheet1", "F1", "Tenor")
	excel.SetCellValue("Sheet1", "G1", "Term")
	excel.SetCellValue("Sheet1", "H1", "Effective Date")
	excel.SetCellValue("Sheet1", "I1", "Maturity")
	excel.SetCellValue("Sheet1", "J1", "Status")

	for index, certificate := range certificates {
		isMatured := "Running"
		if certificate.Maturity.Before(time.Now()) {
			isMatured = "Matured"
		}
		excel.SetCellValue("Sheet1", fmt.Sprintf("A%d", index+2), certificate.FundManager)
		excel.SetCellValue("Sheet1", fmt.Sprintf("B%d", index+2), certificate.ClientName)
		excel.SetCellValue("Sheet1", fmt.Sprintf("C%d", index+2), FormatWithComma(certificate.Amount, 2))
		excel.SetCellValue("Sheet1", fmt.Sprintf("D%d", index+2), certificate.Issuer)
		excel.SetCellValue("Sheet1", fmt.Sprintf("E%d", index+2), FormatWithComma(certificate.Rate, 2))
		excel.SetCellValue("Sheet1", fmt.Sprintf("F%d", index+2), certificate.Tenor)
		excel.SetCellValue("Sheet1", fmt.Sprintf("G%d", index+2), certificate.Term)
		excel.SetCellValue("Sheet1", fmt.Sprintf("H%d", index+2), certificate.EffectiveDate.Format("2006-01-02"))
		excel.SetCellValue("Sheet1", fmt.Sprintf("I%d", index+2), certificate.Maturity.Format("2006-01-02"))
		excel.SetCellValue("Sheet1", fmt.Sprintf("J%d", index+2), isMatured)
	}
	excel.SetActiveSheet(index)
	return excel.WriteToBuffer()
}

func replaceNegativeWithBraces(number interface{}) string {
	switch num := number.(type) {
	case float32:
		if !math.Signbit(float64(num)) {
			return fmt.Sprintf("%s", FormatWithComma(num, 2))
		}
		return fmt.Sprintf("(%s)", FormatWithComma(math.Abs(float64(num)), 2))
	case float64:
		if !math.Signbit(num) {
			return fmt.Sprintf("%s", FormatWithComma(num, 2))
		}
		return fmt.Sprintf("(%s)", FormatWithComma(math.Abs(num), 2))
	case int, int8, int16, int32, int64:
		if !math.Signbit(float64(num.(int64))) {
			return fmt.Sprintf("%s", FormatWithComma(num, 2))
		}
		return fmt.Sprintf("(%s)", FormatWithComma(math.Abs(float64(num.(int64))), 2))
	}
	return ""
}

func ExportNPRAUnauthorizedReport() ([]byte, error) {
	data, err := json.Marshal(LoadMailSenderInfo())
	if err != nil {
		log.Println(err.Error())
	}
	mailerURL := envy.Get("SEC_REPORT_EXPORT_URL", "http://localhost:9090/npra-unauthorized-letter.php")
	response, err := http.Post(mailerURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(response.Body)
		data, err := base64.StdEncoding.DecodeString(string(bodyBytes))
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	return nil, errors.New("error, couldn't convert data to word")
}

func ExportNPRAMonthlyReport(declaration NPRADeclaration) ([]byte, error) {
	data, err := json.Marshal(declaration)
	mailerURL := envy.Get("SEC_REPORT_EXPORT_URL", "http://localhost:9090/npra-monthly-report.php")
	response, err := http.Post(mailerURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(response.Body)
		data, err := base64.StdEncoding.DecodeString(string(bodyBytes))
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	return nil, errors.New("error, couldn't convert data to word")
}

func ExportBillingReport(payload BillingReportExportPayload) ([]byte, error) {
	data, err := json.Marshal(payload)
	mailerURL := envy.Get("SEC_REPORT_EXPORT_URL", "http://localhost:9090/invoice.php")
	response, err := http.Post(mailerURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(response.Body)
		data, err := base64.StdEncoding.DecodeString(string(bodyBytes))
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	return nil, errors.New("error, couldn't convert data to word")
}

func ExportDirectorsToWord(payload DirectorsExportPayload) ([]byte, error) {
	data, err := json.Marshal(payload)
	mailerURL := envy.Get("SEC_REPORT_EXPORT_URL", "http://localhost:9090/directors.php")
	response, err := http.Post(mailerURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(response.Body)
		data, err := base64.StdEncoding.DecodeString(string(bodyBytes))
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	return nil, errors.New("error, couldn't convert data to word")
}

func SendPasswordResetEmail(payload UserPasswordResetPayload) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Println(err.Error())
	}
	mailerURL := "http://localhost:9090/password-reset.php"
	_, err = http.Post(mailerURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Println(err.Error())
	}
}

func MakeTrusteeDates(quarter, year string) []time.Time {
	dates := make([]time.Time, 0)
	yearInt, _ := strconv.Atoi(year)
	if quarter == "1" {
		start, _ := time.Parse("2006-01-02", fmt.Sprintf("%d-10-01", yearInt-1))
		for i := 0; i < 6; i++ {
			dates = append(dates, start.AddDate(0, i*1, 0))
		}
	} else {
		quarterInt, _ := strconv.Atoi(quarter)
		months := quarterInt * 3
		startDate, _ := time.Parse("2006-01-02", fmt.Sprintf("%s-01-01", year))
		for i := 0; i < months; i++ {
			dates = append(dates, startDate.AddDate(0, i*1, 0))
		}
	}
	return dates
}

func MakeTrusteeDatesRelativeToCurrentQuarter(quarter, year string) []time.Time {
	dates := make([]time.Time, 0)
	quarterInt, _ := strconv.Atoi(quarter)
	months := quarterInt * 3
	startDate, _ := time.Parse("2006-01-02", fmt.Sprintf("%s-01-01", year))
	for i := 0; i < months; i++ {
		dates = append(dates, startDate.AddDate(0, i*1, 0))
	}
	return dates
}

func GetSecLastLicenseRenewalDate() time.Time {
	year := time.Now().Year()
	expiryDate, _ := time.Parse("01-02", "06-30")
	var renewalDate time.Time
	if time.Now().After(expiryDate) {
		renewalDate, _ = time.Parse("2006-01-02", fmt.Sprintf("%d-07-01", year))
	}
	renewalDate, _ = time.Parse("2006-01-02", fmt.Sprintf("%d-07-01", year-1))
	return renewalDate
}

func lessForSortingMonthlyContributions(first, second string) bool {
	months := map[string]int{"january": 1, "february": 2, "march": 3, "april": 4, "may": 5, "june": 6, "july": 7, "august": 8, "september": 9, "october": 10, "november": 11, "december": 12, "total": 13}
	return months[strings.ToLower(first)] < months[strings.ToLower(second)]
}

func SortContributionsByMonths(contributions []LumpedMonthlyContribution) []LumpedMonthlyContribution {
	sort.Slice(contributions, func(i, j int) bool {
		return lessForSortingMonthlyContributions(strings.Split(contributions[i].Date, " ")[0], strings.Split(contributions[j].Date, " ")[0])
	})
	return contributions
}

func SortGOGMaturitiesByMonths(maturities []GOGSummary) []GOGSummary {
	sort.Slice(maturities, func(i, j int) bool {
		return lessForSortingMonthlyContributions(strings.Split(maturities[i].Month, " ")[0], strings.Split(maturities[j].Month, " ")[0])
	})
	return maturities
}

func lessForSortingPVSummary(first, second string) bool {
	months := map[string]int{"receivables": 999, "cash balance": 1000, "total": 1001}
	if months[strings.ToLower(first)] == 0 && months[strings.ToLower(second)] == 0 {
		return first < second
	}
	return months[strings.ToLower(first)] < months[strings.ToLower(second)]
}

func SortPVSummary(summary []ClientPVReportSummary) []ClientPVReportSummary {
	sort.Slice(summary, func(i, j int) bool {
		return lessForSortingPVSummary(summary[i].SecurityType, summary[j].SecurityType)
	})
	return summary
}

func SortLumpedPVSummary(summary []map[string]interface{}) []map[string]interface{} {
	sort.Slice(summary, func(i, j int) bool {
		return lessForSortingPVSummary(summary[i]["security_type"].(string), summary[j]["security_type"].(string))
	})
	return summary
}

func GroupTradeVolumesByAssetClass(trades []TradeVolumeByAssetClass) []TradeVolumeByAssetClass {
	groupedTrades := make([]TradeVolumeByAssetClass, 0)
	for _, trade := range trades {
		found := false
		for _, data := range groupedTrades {
			if strings.HasPrefix(strings.ToLower(data.Asset), strings.ToLower(trade.Asset[:3])) {
				data.Number += trade.Number
				found = true
				break
			}
		}
		if !found {
			groupedTrades = append(groupedTrades, trade)
		}
	}
	return groupedTrades
}

func RandomNumber() int64 {
	const startRange = 10000000
	const endRange = 99999999
	numberGen, err := cryptoRand.Int(cryptoRand.Reader, big.NewInt(endRange-startRange))
	if err != nil {
		panic(err)
	}
	return numberGen.Int64() + startRange
}

func CountFDs(fds []StrValueQuery) map[string]int {
	results := make(map[string]int, 0)
	for _, fd := range fds {
		var bankName string
		parsedFd := strings.Split(fd.Value, " ")
		parsedLength := len(parsedFd)
		if parsedLength >= 3 {
			bankName = parsedFd[2]
			if parsedLength > 3 {
				bankName += " " + parsedFd[3]
			}
		} else {
			bankName = fd.Value
		}
		results[bankName] += 1
	}
	return results
}

func GetHolidays() []string {
	holidays := LoadHolidays()
	results := make([]string, 0)
	for _, holiday := range holidays {
		results = append(results, holiday.Date)
	}
	return results
}

func DateIsAHoliday(date time.Time) bool {
	parsedDate := date.Format("02-01")
	for _, day := range GetHolidays() {
		if day == parsedDate {
			return true
		}
	}
	return false
}

func GetWorkingDayFrom(date time.Time) time.Time {
	dayOfTheWeek := int(date.Weekday())
	if DateIsAHoliday(date) || (dayOfTheWeek == 6 || dayOfTheWeek == 0) {
		return GetWorkingDayFrom(date.AddDate(0, 0, 1))
	}
	return date
}

func CreateOffshoreClientsTotal(clients []OffshoreClient) []OffshoreClient {
	var totalAUC float64
	for _, client := range clients {
		totalAUC += client.AssetValue
	}
	return append(clients, OffshoreClient{Name: "TOTAL", AssetValue: totalAUC})
}

func CreateVarianceTotal(data []VarianceData) []VarianceData {
	var prevTotal, currTotal float64
	for _, datum := range data {
		prevTotal += datum.LastAUA
		currTotal += datum.CurrentAUA
	}
	return append(data, VarianceData{
		Name:       "TOTAL",
		LastAUA:    prevTotal,
		CurrentAUA: currTotal,
		Amount:     currTotal - prevTotal,
		Variance:   float32(((currTotal - prevTotal) / prevTotal) * 100),
	})
}

// Export301ToExcel EXPORT NPRA 0301 REPORT
func Export301ToExcel(data []NPRA0301) (*bytes.Buffer, error) {
	excel := excelize.NewFile()
	index := excel.NewSheet("sheet1")

	excel.SetCellValue("sheet1", "A1", "Total Portfolio Returns (Net Returns)")
	excel.SetCellValue("sheet1", "B1", "Total Portfolio Returns (Gross Returns)")
	excel.SetCellValue("sheet1", "C1", "Report Code")
	excel.SetCellValue("sheet1", "D1", "Entity ID")
	excel.SetCellValue("sheet1", "E1", "Entity Name")
	excel.SetCellValue("sheet1", "F1", "Reference Period Year")
	excel.SetCellValue("sheet1", "G1", "Reference Period")
	excel.SetCellValue("sheet1", "H1", "Investment Receivables [ghs]")
	excel.SetCellValue("sheet1", "I1", "Total Asset Under Management [ghs]")
	excel.SetCellValue("sheet1", "J1", "Government Securities [ghs]")
	excel.SetCellValue("sheet1", "K1", "Local Government/Satutory Agency Securities [ghs]")
	excel.SetCellValue("sheet1", "L1", "Corporate Debt Securities [ghs]")
	excel.SetCellValue("sheet1", "M1", "Bank Securities [ghs]")
	excel.SetCellValue("sheet1", "N1", "Ordinary/Preference Shares [ghs]")
	excel.SetCellValue("sheet1", "O1", "Collective Investment Scheme [ghs]")
	excel.SetCellValue("sheet1", "P1", "Alternative Investment Scheme [ghs]")
	excel.SetCellValue("sheet1", "Q1", "Bank Balances [ghs]")
	excel.SetRowHeight("sheet1", 1, 25)

	headingColor, _ := excel.NewStyle(`{"fill":{"type":"pattern","color":["#95b3d7"],"pattern":1}}`)
	fmt.Println("########## Level 3 #############", headingColor)
	// if err != nil {
	// 	return err
	// }
	excel.SetCellStyle("sheet1", "A1", "A1", headingColor)
	excel.SetCellStyle("sheet1", "B1", "B1", headingColor)
	excel.SetCellStyle("sheet1", "C1", "C1", headingColor)
	excel.SetCellStyle("sheet1", "D1", "D1", headingColor)
	excel.SetCellStyle("sheet1", "E1", "E1", headingColor)
	excel.SetCellStyle("sheet1", "F1", "F1", headingColor)
	excel.SetCellStyle("sheet1", "G1", "G1", headingColor)
	excel.SetCellStyle("sheet1", "H1", "H1", headingColor)
	excel.SetCellStyle("sheet1", "I1", "I1", headingColor)
	excel.SetCellStyle("sheet1", "J1", "J1", headingColor)
	excel.SetCellStyle("sheet1", "K1", "K1", headingColor)
	excel.SetCellStyle("sheet1", "L1", "L1", headingColor)
	excel.SetCellStyle("sheet1", "M1", "M1", headingColor)
	excel.SetCellStyle("sheet1", "N1", "N1", headingColor)
	excel.SetCellStyle("sheet1", "O1", "O1", headingColor)
	excel.SetCellStyle("sheet1", "P1", "P1", headingColor)
	excel.SetCellStyle("sheet1", "Q1", "Q1", headingColor)

	for index, datum := range data {
		excel.SetCellValue("sheet1", fmt.Sprintf("A%d", index+2), datum.NetReturn)
		excel.SetCellValue("sheet1", fmt.Sprintf("B%d", index+2), datum.GrossReturn)
		excel.SetCellValue("sheet1", fmt.Sprintf("C%d", index+2), datum.ReportCode)
		excel.SetCellValue("sheet1", fmt.Sprintf("D%d", index+2), datum.EntityID)
		excel.SetCellValue("sheet1", fmt.Sprintf("E%d", index+2), datum.EntityName)
		excel.SetCellValue("sheet1", fmt.Sprintf("F%d", index+2), datum.ReferencePeriodYear)
		excel.SetCellValue("sheet1", fmt.Sprintf("G%d", index+2), datum.ReferencePeriod)
		excel.SetCellValue("sheet1", fmt.Sprintf("H%d", index+2), FormatWithComma(datum.InvestmentReceivables, 2))
		excel.SetCellValue("sheet1", fmt.Sprintf("I%d", index+2), FormatWithComma(datum.TotalAssetUnderManagement, 2))
		excel.SetCellValue("sheet1", fmt.Sprintf("J%d", index+2), FormatWithComma(datum.GovernmentSecurities, 2))
		excel.SetCellValue("sheet1", fmt.Sprintf("K%d", index+2), FormatWithComma(datum.LocalGovernmentSecurities, 2))
		excel.SetCellValue("sheet1", fmt.Sprintf("L%d", index+2), FormatWithComma(datum.CorporateDebtSecurities, 2))
		excel.SetCellValue("sheet1", fmt.Sprintf("M%d", index+2), FormatWithComma(datum.BankSecurities, 2))
		excel.SetCellValue("sheet1", fmt.Sprintf("N%d", index+2), FormatWithComma(datum.OrdinaryPreferenceShares, 2))
		excel.SetCellValue("sheet1", fmt.Sprintf("O%d", index+2), FormatWithComma(datum.CollectiveInvestmentScheme, 2))
		excel.SetCellValue("sheet1", fmt.Sprintf("P%d", index+2), FormatWithComma(datum.AlternativeInvestments, 2))
		excel.SetCellValue("sheet1", fmt.Sprintf("Q%d", index+2), FormatWithComma(datum.BankBalances, 2))
	}
	excel.SetActiveSheet(index)
	return excel.WriteToBuffer()
}

// Export302ToExcel EXPORT NPRA 0302 REPORT
func Export302ToExcel(data []NPRA0302) (*bytes.Buffer, error) {
	excel := excelize.NewFile()
	index := excel.NewSheet("sheet1")

	excel.SetCellValue("sheet1", "A1", "Report Code")
	excel.SetCellValue("sheet1", "B1", "Entity ID")
	excel.SetCellValue("sheet1", "C1", "Entity Name")
	excel.SetCellValue("sheet1", "D1", "Reference Period Year")
	excel.SetCellValue("sheet1", "E1", "Reference Period")
	excel.SetCellValue("sheet1", "F1", "Investment ID")
	excel.SetCellValue("sheet1", "G1", "Instrument")
	excel.SetCellValue("sheet1", "H1", "Issuer Name")
	excel.SetCellValue("sheet1", "I1", "Asset Tenure")
	excel.SetCellValue("sheet1", "J1", "Date of Investment")
	excel.SetCellValue("sheet1", "K1", "Reporting Date")
	excel.SetCellValue("sheet1", "L1", "Amount Invested[ghs]")
	excel.SetCellValue("sheet1", "M1", "Accrued Interest/Coupon for the Month")
	excel.SetCellValue("sheet1", "N1", "Coupon Paid")
	excel.SetCellValue("sheet1", "O1", "Accrued Interest Since Purchase/Coupon since payment Date")
	excel.SetCellValue("sheet1", "P1", "Outstanding Interest to Maturity")
	excel.SetCellValue("sheet1", "Q1", "Amount Impaired[ghs]")
	excel.SetCellValue("sheet1", "R1", "Asset Allocation Actual Percent")
	excel.SetCellValue("sheet1", "S1", "Maturity Date")
	excel.SetCellValue("sheet1", "T1", "Type of Investment Charge")
	excel.SetCellValue("sheet1", "U1", "Investment Charge Rate Percent")
	excel.SetCellValue("sheet1", "V1", "Investment Charge Amount[ghs]")
	excel.SetCellValue("sheet1", "W1", "Face Value[ghs]")
	excel.SetCellValue("sheet1", "X1", "Interest Rate Percent")
	excel.SetCellValue("sheet1", "Y1", "Discount Rate Percent")
	excel.SetCellValue("sheet1", "Z1", "Coupon Rate Percent")
	excel.SetCellValue("sheet1", "AA1", "Disposal Proceeds [GHS]")
	excel.SetCellValue("sheet1", "AB1", "Disposal Instructions")
	excel.SetCellValue("sheet1", "AC1", "Yield on Disposal[GHS]")
	excel.SetCellValue("sheet1", "AD1", "Issue Date")
	excel.SetCellValue("sheet1", "AE1", "Price Per Unit/Share at Purchase")
	excel.SetCellValue("sheet1", "AF1", "Price Per Unit/Share at Value Date")
	excel.SetCellValue("sheet1", "AG1", "Capital Gains")
	excel.SetCellValue("sheet1", "AH1", "Dividend Received")
	excel.SetCellValue("sheet1", "AI1", "Number of Units/Shares")
	excel.SetCellValue("sheet1", "AJ1", "Holding Period Return per an investment(Percent)")
	excel.SetCellValue("sheet1", "AK1", "Day Run")
	excel.SetCellValue("sheet1", "AL1", "Currency Conversion Rate")
	excel.SetCellValue("sheet1", "AM1", "Currency")
	excel.SetCellValue("sheet1", "AN1", "Amount Investment Foreign Currency(Eurobond/External Investment)")
	excel.SetCellValue("sheet1", "AO1", "Asset Class")
	excel.SetCellValue("sheet1", "AP1", "Price Per unit/Share at Last value Date")
	excel.SetCellValue("sheet1", "AQ1", "Market Value[ghs]")
	excel.SetCellValue("sheet1", "AR1", "Remaining Days to Maturity")
	excel.SetCellValue("sheet1", "AS1", "Holding Period Return Per an Investment Weighted percent")
	excel.SetRowHeight("sheet1", 1, 25)

	headingColor, _ := excel.NewStyle(`{"fill":{"type":"pattern","color":["#95b3d7"],"pattern":1}}`)
	fmt.Println("########## Level 3 #############", headingColor)
	// if err != nil {
	// 	return err
	// }
	excel.SetCellStyle("sheet1", "A1", "A1", headingColor)
	excel.SetCellStyle("sheet1", "B1", "B1", headingColor)
	excel.SetCellStyle("sheet1", "C1", "C1", headingColor)
	excel.SetCellStyle("sheet1", "D1", "D1", headingColor)
	excel.SetCellStyle("sheet1", "E1", "E1", headingColor)
	excel.SetCellStyle("sheet1", "F1", "F1", headingColor)
	excel.SetCellStyle("sheet1", "G1", "G1", headingColor)
	excel.SetCellStyle("sheet1", "H1", "H1", headingColor)
	excel.SetCellStyle("sheet1", "I1", "I1", headingColor)
	excel.SetCellStyle("sheet1", "J1", "J1", headingColor)
	excel.SetCellStyle("sheet1", "K1", "K1", headingColor)
	excel.SetCellStyle("sheet1", "L1", "L1", headingColor)
	excel.SetCellStyle("sheet1", "M1", "M1", headingColor)
	excel.SetCellStyle("sheet1", "N1", "N1", headingColor)
	excel.SetCellStyle("sheet1", "O1", "O1", headingColor)
	excel.SetCellStyle("sheet1", "P1", "P1", headingColor)
	excel.SetCellStyle("sheet1", "Q1", "Q1", headingColor)
	excel.SetCellStyle("sheet1", "R1", "R1", headingColor)
	excel.SetCellStyle("sheet1", "S1", "S1", headingColor)
	excel.SetCellStyle("sheet1", "T1", "T1", headingColor)
	excel.SetCellStyle("sheet1", "U1", "U1", headingColor)
	excel.SetCellStyle("sheet1", "V1", "V1", headingColor)
	excel.SetCellStyle("sheet1", "W1", "W1", headingColor)
	excel.SetCellStyle("sheet1", "X1", "X1", headingColor)
	excel.SetCellStyle("sheet1", "Y1", "Y1", headingColor)
	excel.SetCellStyle("sheet1", "Z1", "Z1", headingColor)
	excel.SetCellStyle("sheet1", "AA1", "AA1", headingColor)
	excel.SetCellStyle("sheet1", "AB1", "AB1", headingColor)
	excel.SetCellStyle("sheet1", "AC1", "AC1", headingColor)
	excel.SetCellStyle("sheet1", "AD1", "AD1", headingColor)
	excel.SetCellStyle("sheet1", "AE1", "AE1", headingColor)
	excel.SetCellStyle("sheet1", "AF1", "AF1", headingColor)
	excel.SetCellStyle("sheet1", "AG1", "AG1", headingColor)
	excel.SetCellStyle("sheet1", "AH1", "AH1", headingColor)
	excel.SetCellStyle("sheet1", "AI1", "AI1", headingColor)
	excel.SetCellStyle("sheet1", "AJ1", "AJ1", headingColor)
	excel.SetCellStyle("sheet1", "AK1", "AK1", headingColor)
	excel.SetCellStyle("sheet1", "AL1", "AL1", headingColor)
	excel.SetCellStyle("sheet1", "AM1", "AM1", headingColor)
	excel.SetCellStyle("sheet1", "AN1", "AN1", headingColor)
	excel.SetCellStyle("sheet1", "AO1", "AO1", headingColor)
	excel.SetCellStyle("sheet1", "AP1", "AP1", headingColor)
	excel.SetCellStyle("sheet1", "AQ1", "AQ1", headingColor)
	excel.SetCellStyle("sheet1", "AR1", "AR1", headingColor)
	excel.SetCellStyle("sheet1", "AS1", "AS1", headingColor)

	for index, datum := range data {
		excel.SetCellValue("sheet1", fmt.Sprintf("A%d", index+2), datum.ReportCode)
		excel.SetCellValue("sheet1", fmt.Sprintf("B%d", index+2), datum.EntityID)
		excel.SetCellValue("sheet1", fmt.Sprintf("C%d", index+2), datum.EntityName)
		excel.SetCellValue("sheet1", fmt.Sprintf("D%d", index+2), datum.ReferencePeriodYear)
		excel.SetCellValue("sheet1", fmt.Sprintf("E%d", index+2), datum.ReferencePeriod)
		excel.SetCellValue("sheet1", fmt.Sprintf("F%d", index+2), datum.InvestmentID)
		excel.SetCellValue("sheet1", fmt.Sprintf("G%d", index+2), datum.Instrument)
		excel.SetCellValue("sheet1", fmt.Sprintf("H%d", index+2), datum.IssuerName)
		excel.SetCellValue("sheet1", fmt.Sprintf("I%d", index+2), datum.AssetTenure)
		excel.SetCellValue("sheet1", "J%d", "")
		excel.SetCellValue("sheet1", fmt.Sprintf("K%d", index+2), datum.ReportingDate.Format("2006-01-02"))
		excel.SetCellValue("sheet1", "L%d", "")
		excel.SetCellValue("sheet1", "M%d", "")
		excel.SetCellValue("sheet1", "N%d", "")
		excel.SetCellValue("sheet1", "O%d", "")
		excel.SetCellValue("sheet1", "P%d", "")
		excel.SetCellValue("sheet1", "Q%d", "")
		excel.SetCellValue("sheet1", "R%d", "")
		excel.SetCellValue("sheet1", fmt.Sprintf("S%d", index+2), datum.MaturityDate.Format("2006-01-02"))
		excel.SetCellValue("sheet1", "T%d", "")
		excel.SetCellValue("sheet1", "U%d", "")
		excel.SetCellValue("sheet1", "V%d", "")
		excel.SetCellValue("sheet1", fmt.Sprintf("W%d", index+2), FormatWithComma(datum.FaceValue, 2))
		excel.SetCellValue("sheet1", "X%d", "")
		excel.SetCellValue("sheet1", "Y%d", "")
		excel.SetCellValue("sheet1", "Z%d", "")
		excel.SetCellValue("sheet1", "AA%d", "")
		excel.SetCellValue("sheet1", "AB%d", "")
		excel.SetCellValue("sheet1", "AC%d", "")
		excel.SetCellValue("sheet1", fmt.Sprintf("AD%d", index+2), datum.IssueDate.Format("2006-01-02"))
		excel.SetCellValue("sheet1", "AE%d", "")
		excel.SetCellValue("sheet1", "AF%d", "")
		excel.SetCellValue("sheet1", "AG%d", "")
		excel.SetCellValue("sheet1", "AH%d", "")
		excel.SetCellValue("sheet1", "AI%d", "")
		excel.SetCellValue("sheet1", "AJ%d", "")
		excel.SetCellValue("sheet1", "AK%d", "")
		excel.SetCellValue("sheet1", "AL%d", "")
		excel.SetCellValue("sheet1", fmt.Sprintf("AM%d", index+2), datum.Currency)
		excel.SetCellValue("sheet1", "AN%d", "")
		excel.SetCellValue("sheet1", fmt.Sprintf("AO%d", index+2), datum.AssetClass)
		excel.SetCellValue("sheet1", "AP%d", "")
		excel.SetCellValue("sheet1", fmt.Sprintf("AQ%d", index+2), FormatWithComma(datum.MarketValue, 2))
		excel.SetCellValue("sheet1", "AR%d", "")
		excel.SetCellValue("sheet1", "AS%d", "")
	}
	excel.SetActiveSheet(index)
	return excel.WriteToBuffer()
}

// Export303ToExcel EXPORT NPRA 0303 REPORT
func Export303ToExcel(data []NPRA0303) (*bytes.Buffer, error) {
	excel := excelize.NewFile()
	index := excel.NewSheet("sheet1")

	excel.SetCellValue("sheet1", "A1", "Report Code")
	excel.SetCellValue("sheet1", "B1", "Entity ID")
	excel.SetCellValue("sheet1", "C1", "Report Code")
	excel.SetCellValue("sheet1", "D1", "Reference Period Year")
	excel.SetCellValue("sheet1", "E1", "Reference Period")
	excel.SetCellValue("sheet1", "F1", "Unit Price")
	excel.SetCellValue("sheet1", "G1", " Date of Valuation")
	excel.SetCellValue("sheet1", "H1", "Daily NAV")
	excel.SetCellValue("sheet1", "I1", "Unit Number")
	excel.SetCellValue("sheet1", "J1", "NPRA Fees")
	excel.SetCellValue("sheet1", "K1", "Trustee Fees")
	excel.SetCellValue("sheet1", "L1", "Fund Managers Fees")
	excel.SetCellValue("sheet1", "M1", "Fund Custodian Fees")
	excel.SetRowHeight("sheet1", 1, 25)

	headingColor, _ := excel.NewStyle(`{"fill":{"type":"pattern","color":["#95b3d7"],"pattern":1}}`)
	fmt.Println("########## Level 3 #############", headingColor)
	// if err != nil {
	// 	return err
	// }
	excel.SetCellStyle("sheet1", "A1", "A1", headingColor)
	excel.SetCellStyle("sheet1", "B1", "B1", headingColor)
	excel.SetCellStyle("sheet1", "C1", "C1", headingColor)
	excel.SetCellStyle("sheet1", "D1", "D1", headingColor)
	excel.SetCellStyle("sheet1", "E1", "E1", headingColor)
	excel.SetCellStyle("sheet1", "F1", "F1", headingColor)
	excel.SetCellStyle("sheet1", "G1", "G1", headingColor)
	excel.SetCellStyle("sheet1", "H1", "H1", headingColor)
	excel.SetCellStyle("sheet1", "I1", "I1", headingColor)
	excel.SetCellStyle("sheet1", "J1", "J1", headingColor)
	excel.SetCellStyle("sheet1", "K1", "K1", headingColor)
	excel.SetCellStyle("sheet1", "L1", "L1", headingColor)
	excel.SetCellStyle("sheet1", "M1", "M1", headingColor)

	for index, datum := range data {
		excel.SetCellValue("sheet1", fmt.Sprintf("A%d", index+2), datum.ReportCode)
		excel.SetCellValue("sheet1", fmt.Sprintf("B%d", index+2), datum.EntityID)
		excel.SetCellValue("sheet1", fmt.Sprintf("C%d", index+2), datum.EntityName)
		excel.SetCellValue("sheet1", fmt.Sprintf("D%d", index+2), datum.ReferencePeriodYear)
		excel.SetCellValue("sheet1", fmt.Sprintf("E%d", index+2), datum.ReferencePeriod)
		excel.SetCellValue("sheet1", fmt.Sprintf("F%d", index+2), FormatWithComma(datum.UnitPrice, 2))
		excel.SetCellValue("sheet1", fmt.Sprintf("G%d", index+2), datum.DateValuation.Format("2006-01-02"))
		excel.SetCellValue("sheet1", fmt.Sprintf("H%d", index+2), datum.UnitNumber)
		excel.SetCellValue("sheet1", fmt.Sprintf("I%d", index+2), FormatWithComma(datum.DailyNav, 2))
		excel.SetCellValue("sheet1", fmt.Sprintf("J%d", index+2), FormatWithComma(datum.NPRAFees, 2))
		excel.SetCellValue("sheet1", fmt.Sprintf("K%d", index+2), FormatWithComma(datum.TrusteeFees, 2))
		excel.SetCellValue("sheet1", fmt.Sprintf("L%d", index+2), FormatWithComma(datum.FundManagerFees, 2))
		excel.SetCellValue("sheet1", fmt.Sprintf("M%d", index+2), FormatWithComma(datum.FundCustodianFees, 2))
	}
	excel.SetActiveSheet(index)
	return excel.WriteToBuffer()
}
