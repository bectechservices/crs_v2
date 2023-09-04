package actions

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/olekukonko/tablewriter"
)

func parseExcelData(excelData []ParsedPVReport, userID int, PVType string, cashBalance float64) ([]PVClient, []PVUploadError) {
	clientReports := make([]PVClient, 0)
	parseErrors := make([]PVUploadError, 0)
	clientReportsChannel := make(chan PVClient)

	securityNames := LoadMergedSecurityNames()

	var dataParseWaitGroup sync.WaitGroup
	for _, eachParsedPV := range excelData {
		dataParseWaitGroup.Add(1)
		go func(wg *sync.WaitGroup, reportsChannel chan PVClient, data ParsedPVReport, securityNames []MergedSecurityName) {
			defer wg.Done()
			var clientID string
			var clientRawHeader string
			clientReport := PVClient{}
			reportData := make([]ReportData, 0)
			reportErrors := make([]PVUploadError, 0)
			currentReportBlock := ReportData{
				Title:  "",
				Values: make([]ReportField, 0),
			}
			summaryData := make([]ReportSummary, 0)
			hasCapturedClientID := false
			for _, header := range data.Headers {
				if hasReportDate(header) {
					if clientReport.HasNoDate() {
						dateStr := getDatesFromString(header).To //because getDatesFromString returns only To for strings with single dates
						parsedDate := strings.Split(dateStr, "-")
						if len(parsedDate) == 3 {
							clientReport.Date = fmt.Sprintf("%s-%s-%s", parsedDate[2], parsedDate[1], parsedDate[0])
						}
					}
				}
				if !hasCapturedClientID {
					if clientID = getClientIDFromHeader(header); clientID != "" {
						hasCapturedClientID = true
						clientRawHeader = header
					}
				}
			}

			for _, bonds := range data.Bonds {
				currentReportBlock.Title = bonds.Bond
				if !NameIsASecurityName(securityNames, currentReportBlock.Title) {
					reportErrors = append(reportErrors, PVUploadError{
						Type:       "Security Name Error",
						SCA:        clientID,
						ClientInfo: clientRawHeader,
						Date:       clientReport.Date,
						MoreInfo:   currentReportBlock.Title,
					})
				}

				for index, bond := range bonds.Values {
					if index == len(bonds.Values)-1 {
						break
					}
					pvRow, err := createPVReportFieldFromRow(bond)
					if err != nil {
						reportErrors = append(reportErrors, PVUploadError{
							Type:       "Security Parse Error",
							SCA:        clientID,
							ClientInfo: clientRawHeader,
							Date:       clientReport.Date,
							MoreInfo:   currentReportBlock.Title,
						})
					}
					currentReportBlock.Values = append(currentReportBlock.Values, pvRow)
				}
				reportData = append(reportData, currentReportBlock)
				currentReportBlock = ReportData{
					Title:  "",
					Values: make([]ReportField, 0),
				}
			}

			for index, summary := range data.Summary {
				if index == len(data.Summary)-1 {
					break
				}
				summaryRow, err := createSummaryDataFromRow(summary)
				if err != nil {
					reportErrors = append(reportErrors, PVUploadError{
						Type:       "Summary Parse Error",
						SCA:        clientID,
						ClientInfo: clientRawHeader,
						Date:       clientReport.Date,
					})
				} else {
					if !NameIsASecurityName(securityNames, summaryRow.SecurityName) {
						reportErrors = append(reportErrors, PVUploadError{
							Type:       "Security Name Error",
							SCA:        clientID,
							ClientInfo: clientRawHeader,
							Date:       clientReport.Date,
							MoreInfo:   fmt.Sprintf("Error in summary: %s", summaryRow.SecurityName),
						})
					}
				}
				summaryData = append(summaryData, summaryRow)
			}

			hasCapturedClientID = false

			clientReport.UserID = userID
			clientReport.ClientID = clientID
			clientReport.CashBalance = cashBalance
			clientReport.Type = PVType
			clientReport.Report = reportData
			clientReport.Summary = summaryData
			clientReport.RawHeading = clientRawHeader
			clientReport.Error = reportErrors

			clientReportsChannel <- clientReport
		}(&dataParseWaitGroup, clientReportsChannel, eachParsedPV, securityNames)
	}

	go func(wg *sync.WaitGroup, reportsChannel chan PVClient) {
		wg.Wait()
		close(reportsChannel)
	}(&dataParseWaitGroup, clientReportsChannel)

	for report := range clientReportsChannel {
		canIgnoreErrors := PVErrorsAreSecurityNameErrors(report.Error)
		if (len(report.Error) > 0) && !canIgnoreErrors {
			parseErrors = append(parseErrors, report.Error...)
		} else if (len(report.Error) > 0) && canIgnoreErrors {
			parseErrors = append(parseErrors, report.Error...)
			clientReports = append(clientReports, report)
		} else {
			clientReports = append(clientReports, report)
		}
	}

	return clientReports, parseErrors
}

func PVErrorsAreSecurityNameErrors(pvErrors []PVUploadError) bool {
	for _, pvError := range pvErrors {
		if pvError.Type != "Security Name Error" {
			return false
		}
	}
	return true
}

func dumpPVReportToConsole(data PVClient) {
	fmt.Println("Report Date:", data.Date)
	fmt.Println("PV Type:", data.Type)
	fmt.Println("Client ID:", data.ClientID)
	fmt.Println("Cash Balance:", data.CashBalance)
	for i := 0; i < len(data.Report); i++ {
		table := tablewriter.NewWriter(os.Stdout)
		fmt.Println(data.Report[i].Title)
		table.SetHeader([]string{ /*"Security Name",*/ "CDS Code", "ISIN", "SCB Code", "Market Price", "Nominal Value", "Cumulative Cost", "Value", "Percentage of Total"})
		for j := 0; j < len(data.Report[i].Values); j++ {
			table.Append([]string{ /*data.Report[i].Values[j].SecurityName,*/ data.Report[i].Values[j].CDSCode, data.Report[i].Values[j].ISIN, data.Report[i].Values[j].SCBCode, fmt.Sprintf("%.2f", data.Report[i].Values[j].MarketPrice), fmt.Sprintf("%.2f", data.Report[i].Values[j].NominalValue), fmt.Sprintf("%.2f", data.Report[i].Values[j].CumulativeCost), fmt.Sprintf("%.2f", data.Report[i].Values[j].Value), fmt.Sprintf("%.2f", data.Report[i].Values[j].PercentageOfTotal)})
		}
		table.Render()
		fmt.Println("")
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Description", "Nominal Value", "Cumulative Cost", "Value LCY", "Percentage of Total"})
	summaryData := make([][]string, 0)
	for k := 0; k < len(data.Summary); k++ {
		summaryData = append(summaryData, []string{
			data.Summary[k].SecurityName,
			fmt.Sprintf("%.2f", data.Summary[k].NominalValue),
			fmt.Sprintf("%.2f", data.Summary[k].CumulativeCost),
			fmt.Sprintf("%.2f", data.Summary[k].Value),
			fmt.Sprintf("%.2f", data.Summary[k].PercentageOfTotal),
		})
	}
	table.AppendBulk(summaryData)
	fmt.Println("")
	table.Render()
	fmt.Println("")
}
