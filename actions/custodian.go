package actions

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"io"
	"strings"
	"time"
)

//ShowCustodianPage displays the custodian page
func ShowCustodianPage(c buffalo.Context) error {
	quarter := GetQuarterNumber(time.Now())
	currentQuarterNumber := fmt.Sprintf("%d", quarter)
	currentYearNumber := fmt.Sprintf("%d", GetYearFromQuarter(quarter))
	c.Set("offshore_clients", LoadOffshoreClients(OffshoreClientsRequest{
		Quarter: currentQuarterNumber,
		Year:    currentYearNumber,
	}))
	return c.Render(200, r.HTML("custodian.html"))
}

//ShowPVPage displays the custodian page
func ShowPVPage(c buffalo.Context) error {
	c.Set("scheme", LoadSchemeFromActivity(c.Param("bpid"), c.Param("quarter")))
	return c.Render(200, r.HTML("viewpv.html"))
}

//HandleGovernanceDataUpload handles the request to upload the governance data
func HandleGovernanceDataUpload(c buffalo.Context) error {
	request := &GovernanceDataUploadRequest{}
	if err := c.Bind(request); err != nil {
		return c.Render(200, r.JSON(map[string]interface{}{"error": true, "message": err.Error()}))
	}
	if err := InsertGovernanceInfoIntoTheDB(*request); err != nil {
		return c.Render(200, r.JSON(map[string]interface{}{"error": true, "message": err.Error()}))
	}
	return c.Render(201, r.JSON(map[string]interface{}{"error": false, "message": "Data uploaded"}))
}

//LoadSecQuarterlyReport retrieves the sec quarterly report
func LoadSecQuarterlyReport(c buffalo.Context) error {
	request := &QuarterAndYear{}
	if err := c.Bind(request); err != nil {
		return c.Render(200, r.JSON(map[string]interface{}{"error": true, "message": err.Error()}))
	}
	report, schemeDetails, info, remarks := getQuarterlyReportByDate(MakeQuarterDate(request.Quarter, request.Year))

	if report.DateOfReport.IsZero() {
		report.DateOfReport = time.Now()
	}

	if len(schemeDetails) > 0 {
		schemeDetails[0].AttachedFile = ConvertFileToBase64(schemeDetails[0].AttachedFile)
	}

	return c.Render(200, r.JSON(map[string]interface{}{"error": false, "message": "Report loaded", "data": map[string]interface{}{"report": report, "schemeDetails": schemeDetails, "otherInformation": info, "remarks": remarks}}))
}

//ShowSecDashboard displays the sec dashboard
func ShowSecDashboard(c buffalo.Context) error {
	c.Set("shares", LoadGovernanceInfoByQuarterDate(GetQuarterDate()))
	c.Set("scheme_submission_details", LoadQuarterSchemeSubmissionDetails())
	c.Set("recent_activities", LoadSecRecentActivities())
	return c.Render(200, r.HTML("sec.html"))
}

//HandleSchemeDetailsDataUpload handles the request to upload the scheme details
func HandleSchemeDetailsDataUpload(c buffalo.Context) error {
	request := &SchemeDetails{}
	if err := c.Bind(request); err != nil {
		return c.Render(200, r.JSON(map[string]interface{}{"error": true, "message": err.Error()}))
	}
	if err := InsertSchemeDetailsIntoTheDB(*request, AuthID(c)); err != nil {
		return c.Render(200, r.JSON(map[string]interface{}{"error": true, "message": err.Error()}))
	}
	return c.Render(201, r.JSON(map[string]interface{}{"error": false, "message": "Data uploaded"}))
}

func HandleSchemeDetailsDelete(c buffalo.Context) error {
	request := OnlyBPID{}
	if err := c.Bind(&request); err != nil {
		return c.Render(200, r.JSON(map[string]interface{}{"error": true, "message": err.Error()}))
	}
	if err := DeleteSchemeDetails(request.BPID); err != nil {
		return c.Render(200, r.JSON(map[string]interface{}{"error": true, "message": err.Error()}))
	}
	return c.Render(201, r.JSON(map[string]interface{}{"error": false, "message": "Scheme deleted"}))
}

//HandleOtherInformationUpload handles the request to upload other information
func HandleOtherInformationUpload(c buffalo.Context) error {
	request := &OtherInformation{}
	if err := c.Bind(request); err != nil {
		return c.Render(200, r.JSON(map[string]interface{}{"error": true, "message": err.Error()}))
	}
	if err := InsertOtherInformationIntoTheDB(*request); err != nil {
		return c.Render(200, r.JSON(map[string]interface{}{"error": true, "message": err.Error()}))
	}
	return c.Render(201, r.JSON(map[string]interface{}{"error": false, "message": "Data uploaded"}))
}

//HandleOfficialReportRemarksUpload handles the request to upload the official report remarks
func HandleOfficialReportRemarksUpload(c buffalo.Context) error {
	request := &OfficialReportRemarks{}
	if err := c.Bind(request); err != nil {
		return c.Render(200, r.JSON(map[string]interface{}{"error": true, "message": err.Error()}))
	}
	if err := InsertOfficialReportRemarksIntoTheDB(*request); err != nil {
		return c.Render(200, r.JSON(map[string]interface{}{"error": true, "message": err.Error()}))
	}
	return c.Render(201, r.JSON(map[string]interface{}{"error": false, "message": "Data uploaded"}))
}

//HandleFetchSchemeDetails handles the request to upload the fetch the scheme details
func HandleFetchSchemeDetails(c buffalo.Context) error {
	request := &SchemeDetailsFetchRequest{}
	scheme := SchemeDetails{}
	if err := c.Bind(&request); err != nil {
		return c.Render(400, r.JSON(map[string]interface{}{"error": true, "message": err.Error(), "data": map[string]interface{}{"schemeDetails": scheme}}))
	}
	scheme = FetchSchemeDetails(*request)
	scheme.BPID = request.BPID
	if scheme.AttachedFile != "" {
		scheme.AttachedFile = ConvertFileToBase64(scheme.AttachedFile)
	}
	if strings.TrimSpace(scheme.NameOfScheme) == "" {
		client := GetClientByBPID(request.BPID)
		//quarter, _ := time.Parse("2006-01-02", request.QuarterDate)
		//if client.WasClosedAsAt(quarter) {
		//	return c.Render(400, r.JSON(map[string]interface{}{"error": true, "message": "Client Closed", "data": map[string]interface{}{"schemeDetails": SchemeDetails{}}}))
		//}
		scheme.NameOfScheme = client.Name
		scheme.NameOfManager = client.SchemeManager
		equities, investments, percentage := CalculatePercentageCapitalMarketValue(request.BPID)
		scheme.PercentageCapitalMarketInvestment = float32(percentage)
		scheme.TotalEquityInvestment = equities
		scheme.TotalFixedIncomeInvestment = investments
		scheme.TotalValueOfUnutilizedFunds = LoadClientUnutilizedFunds(request.BPID)
	}
	return c.Render(200, r.JSON(map[string]interface{}{"error": false, "message": "Scheme details loaded", "data": map[string]interface{}{"schemeDetails": scheme}}))
}

//HandleRecalculateSchemeDetails handles the request to upload the fetch the scheme details
func HandleRecalculateSchemeDetails(c buffalo.Context) error {
	request := &SchemeDetailsFetchRequest{}
	scheme := SchemeDetails{}
	if err := c.Bind(&request); err != nil {
		return c.Render(400, r.JSON(map[string]interface{}{"error": true, "message": err.Error(), "data": map[string]interface{}{"schemeDetails": scheme}}))
	}
	scheme = FetchSchemeDetails(*request)
	scheme.BPID = request.BPID
	if scheme.AttachedFile != "" {
		scheme.AttachedFile = ConvertFileToBase64(scheme.AttachedFile)
	}
	equities, investments, percentage := CalculatePercentageCapitalMarketValue(request.BPID)
	scheme.PercentageCapitalMarketInvestment = float32(percentage)
	scheme.TotalEquityInvestment = equities
	scheme.TotalFixedIncomeInvestment = investments
	scheme.TotalValueOfUnutilizedFunds = LoadClientUnutilizedFunds(request.BPID)

	return c.Render(200, r.JSON(map[string]interface{}{"error": false, "message": "Scheme details loaded", "data": map[string]interface{}{"schemeDetails": scheme}}))
}

//ShowSecReportPreview loads the sec report preview page
func ShowSecReportPreview(c buffalo.Context) error {
	date := GetQuarterDate()
	report, schemeDetails, info, remarks := getQuarterlyReportByDate(date)
	if report.DateOfReport.IsZero() {
		report.DateOfReport = time.Now()
	}
	c.Set("compliance", LoadSecInternalComplianceForQuarter())
	c.Set("report", report)
	c.Set("schemeDetails", schemeDetails)
	c.Set("otherInformation", info)
	c.Set("remarks", remarks)
	quarterNum := GetQuarterNumber(time.Now())
	currentQuarterNumber := fmt.Sprintf("%d", quarterNum)
	currentYearNumber := fmt.Sprintf("%d", GetYearFromQuarter(quarterNum))
	c.Set("offshore_clients", LoadOffshoreClients(OffshoreClientsRequest{
		Quarter: currentQuarterNumber,
		Year:    currentYearNumber,
	}))
	currentQuarter := MakeQuarterFormalDate(currentQuarterNumber, currentYearNumber)
	c.Set("local_variance", LoadSecLocalVariance(MakeLastQuarterFormalDate(currentQuarter), currentQuarter))
	c.Set("foreign_variance", LoadSecForeignVariance(MakeLastQuarterFormalDate(currentQuarter), currentQuarter))
	c.Set("directors", LoadAllDirectors())
	c.Set("maturities", LoadMaturedSecurities(currentQuarter))
	c.Set("quarter", currentQuarter)
	c.Set("mailSenderInfo", LoadMailSenderInfo())
	return c.Render(200, r.HTML("sec-report-preview.html"))
}

//ShowSECReportsPage displays the sec report page
func ShowSECReportsPage(c buffalo.Context) error {
	return c.Render(200, r.HTML("secreport.html"))
}

// EyeBallPvPage displays the Pv for the Checker to audit data before upload
func EyeBallPvPage(c buffalo.Context) error {
	return c.Render(200, r.HTML("audit-pv.html"))
}

//LoadUploadedPVForEyeBalling loads the pv data for eyeballing
func LoadUploadedPVForEyeBalling(c buffalo.Context) error {
	request := &PVFetchForEyeBallingRequest{}
	if err := c.Bind(request); err != nil {
		return c.Render(200, r.JSON(map[string]interface{}{"error": true, "message": err.Error()}))
	}
	data := LoadAuditablePVData(request.BPID, request.QuarterDate, request.PVType)
	return c.Render(200, r.JSON(map[string]interface{}{"error": false, "message": "PV details loaded", "data": map[string]interface{}{"reportData": data,}}))
}

func getQuarterlyReportByDate(date string) (report GovernanceInfo, schemeDetails []SchemeDetails, info OtherInformation, remarks OfficialReportRemarks) {
	report = LoadGovernanceInfoByQuarterDate(date)
	schemeDetails = make([]SchemeDetails, 0)
	info = OtherInformation{}
	remarks = OfficialReportRemarks{}
	if report.ID != 0 {
		schemeDetails = GetSchemeDetailsByGovernanceInfoID(report.ID)
		info = GetOtherInformationByGovernanceInfoID(report.ID)
		remarks = GetOfficialRemarksByGovernanceInfoID(report.ID)
	}
	return report, schemeDetails, info, remarks
}

//ExportSecReportToWord exports the sec report to word
func ExportSecReportToWord(c buffalo.Context) error {
	request := &QuarterAndYear{}
	if err := c.Bind(request); err != nil {
		return err
	}
	var quarterFormalDate, date, quarterName, lastLicenseRenewalDate string
	if request.IsEmpty() {
		quarterFormalDate = GetQuarterFormalDate()
		date = GetQuarterDate()
		quarter := GetQuarterNumber(time.Now())
		year := GetYearFromQuarter(quarter)
		request.Quarter = fmt.Sprintf("%d", quarter)
		request.Year = fmt.Sprintf("%d", year)
	} else {
		quarterFormalDate = MakeQuarterFormalDate(request.Quarter, request.Year)
		date = MakeQuarterDate(request.Quarter, request.Year)
	}

	switch request.Quarter {
	case "1":
		quarterName = "First Quarter, " + request.Year
	case "2":
		quarterName = "Second Quarter, " + request.Year
	case "3":
		quarterName = "Third Quarter, " + request.Year
	default:
		quarterName = "Fourth Quarter, " + request.Year
	}

	lastLicenseRenewalDate = MakeSecLastLicenseRenewalDate(request.Quarter, request.Year)

	filename := fmt.Sprintf("sec-%s-report.docx", quarterFormalDate)
	c.Response().Header().Set("Content-Disposition", "attachment; filename="+filename)
	return c.Render(200, r.Func("application/vnd.openxmlformats-officedocument.wordprocessingml.document", func(w io.Writer, d render.Data) error {

		report, schemeDetails, info, remarks := getQuarterlyReportByDate(date)
		clients := LoadOffshoreClients(OffshoreClientsRequest{Year: request.Year, Quarter: request.Quarter})

		data, err := ExportSecReport(quarterName, lastLicenseRenewalDate, report, schemeDetails, info, remarks, clients)
		if err != nil {
			return err
		}
		_, writeError := w.Write(data)
		return writeError
	}))
}

//ExportSecForeignVarianceToExcel exports the foreign variance data to excel
func ExportSecForeignVarianceToExcel(c buffalo.Context) error {
	request := &QuarterAndYear{}
	if err := c.Bind(request); err != nil {
		return err
	}
	c.Response().Header().Set("Content-Disposition", "attachment; filename=foreign-variance.xlsx")
	c.Response().Header().Set("Content-Transfer-Encoding", "binary")
	return c.Render(200, r.Func("application/octet-stream", func(w io.Writer, d render.Data) error {

		currentQuarter := GetQuarterFormalDate()
		currentQuarterShortDate := GetShortQuarterDate()
		previousQuarterShortDate := GetShortPreviousQuarterDate()
		if (request.Year != "") && (request.Quarter != "") {
			currentQuarter = MakeQuarterFormalDate(request.Quarter, request.Year)
			//TODO: show aua as at for arbitrary dates
		}

		data := LoadSecForeignVariance(MakeLastQuarterFormalDate(currentQuarter), currentQuarter)

		bytes, err := ExportVarianceToExcel(previousQuarterShortDate, currentQuarterShortDate, data)
		if err != nil {
			return err
		}
		_, writeError := w.Write(bytes.Bytes())
		return writeError
	}))
}

//ExportSecLocalVarianceToExcel exports the foreign variance data to excel
func ExportSecLocalVarianceToExcel(c buffalo.Context) error {
	request := &QuarterAndYear{}
	if err := c.Bind(request); err != nil {
		return err
	}
	c.Response().Header().Set("Content-Disposition", "attachment; filename=local-variance.xlsx")
	c.Response().Header().Set("Content-Transfer-Encoding", "binary")
	return c.Render(200, r.Func("application/octet-stream", func(w io.Writer, d render.Data) error {

		currentQuarter := GetQuarterFormalDate()
		currentQuarterShortDate := GetShortQuarterDate()
		previousQuarterShortDate := GetShortPreviousQuarterDate()
		if (request.Year != "") && (request.Quarter != "") {
			currentQuarter = MakeQuarterFormalDate(request.Quarter, request.Year)
			//TODO: show aua as at for arbitrary dates
		}

		data := LoadSecLocalVariance(MakeLastQuarterFormalDate(currentQuarter), currentQuarter)

		bytes, err := ExportVarianceToExcel(previousQuarterShortDate, currentQuarterShortDate, data)
		if err != nil {
			return err
		}
		_, writeError := w.Write(bytes.Bytes())
		return writeError
	}))
}

//ExportOffshoreClientsToExcel exports the foreign variance data to excel
func ExportOffshoreClientsToExcel(c buffalo.Context) error {
	request := &QuarterAndYear{}
	if err := c.Bind(request); err != nil {
		return err
	}
	var currentQuarterNumber, currentYearNumber string
	if request.IsEmpty() {
		quarter := GetQuarterNumber(time.Now())
		currentQuarterNumber = fmt.Sprintf("%d", quarter)
		currentYearNumber = fmt.Sprintf("%d", GetYearFromQuarter(quarter))
	} else {
		currentQuarterNumber = request.Quarter
		currentYearNumber = request.Year
	}
	c.Response().Header().Set("Content-Disposition", "attachment; filename=offshore-clients.xlsx")
	c.Response().Header().Set("Content-Transfer-Encoding", "binary")
	return c.Render(200, r.Func("application/octet-stream", func(w io.Writer, d render.Data) error {
		data := LoadOffshoreClients(OffshoreClientsRequest{
			Quarter: currentQuarterNumber,
			Year:    currentYearNumber,
		})

		excel := excelize.NewFile()
		index := excel.NewSheet("Sheet1")

		excel.SetCellValue("Sheet1", "A1", "Client")
		excel.SetCellValue("Sheet1", "B1", "Home Country")
		excel.SetCellValue("Sheet1", "C1", "Asset Value (GHS)")

		for index, datum := range data {
			excel.SetCellValue("Sheet1", fmt.Sprintf("A%d", index+2), datum.Name)
			excel.SetCellValue("Sheet1", fmt.Sprintf("B%d", index+2), datum.Country)
			excel.SetCellValue("Sheet1", fmt.Sprintf("C%d", index+2), FormatWithComma(datum.AssetValue, 2))
		}
		excel.SetActiveSheet(index)
		bytes, err := excel.WriteToBuffer()
		if err != nil {
			return err
		}
		_, writeError := w.Write(bytes.Bytes())
		return writeError
	}))
}

//ExportCoverLetter exports the cover letter
func ExportCoverLetter(c buffalo.Context) error {
	filename := fmt.Sprintf("sec-%s-cover-letter.docx", GetQuarterFormalDate())
	c.Response().Header().Set("Content-Disposition", "attachment; filename="+filename)
	return c.Render(200, r.Func("application/vnd.openxmlformats-officedocument.wordprocessingml.document", func(w io.Writer, d render.Data) error {
		data, err := ExportSecCoverLetter()
		if err != nil {
			return err
		}
		_, writeError := w.Write(data)
		return writeError
	}))
}

func HandleSecPVApproval(c buffalo.Context) error {
	ApproveSecCurrentQuarterReport(AuthID(c))
	return c.Redirect(302, "/sec-report-preview")
}

func HandleSecCurrentQuarterReportDisapproval(c buffalo.Context) error {
	request := &CheckerCommentRequest{}
	if err := c.Bind(request); err != nil {
		return err
	}
	DisapproveSecCurrentQuarterReport(AuthID(c), request.Comment)
	return c.Redirect(302, "/sec-report-preview")
}

func HandleSecSendProgressToChecker(c buffalo.Context) error {
	SendReportActionEmail("report_prepared", LoadUsersByRole(SECCheckerRole), "SEC Quarterly Report", "SEC", "")
	return c.Render(201, r.JSON(map[string]interface{}{"error": false, "message": "Data uploaded"}))
}

func HandleUploadVarianceRemarks(c buffalo.Context) error {
	request := VarianceRemarksUpdateRequest{}
	if err := c.Bind(&request); err != nil {
		return err
	}
	UploadSECVarianceRemarks(request)
	return c.Redirect(302, "/sec-report-preview")
}

func HandleDeleteSECUploadedPV(c buffalo.Context) error {
	request := DeleteUploadedPVRequest{}
	if err := c.Bind(&request); err != nil {
		return c.Render(200, r.JSON(map[string]interface{}{"error": true, "message": err.Error()}))
	}
	DeleteSECPV(request)
	return c.Render(200, r.JSON(map[string]interface{}{"error": false, "message": "PV deleted"}))
}

func HandleLoadMaturedSecurities(c buffalo.Context) error {
	maturities := LoadMaturedSecurities(c.Param("quarter"))
	isApproved := false
	if len(maturities) > 0 {
		isApproved = maturities[0].Approved
	}
	c.Set("maturities", maturities)
	c.Set("approved", isApproved)
	c.Set("quarter", c.Param("quarter"))
	return c.Render(200, r.HTML("matured-securities.html"))
}

func HandleMaturedSecuritiesUpload(c buffalo.Context) error {
	request := MaturedSecuritiesUpload{}
	if err := c.Bind(&request); err != nil {
		return c.Render(200, r.JSON(map[string]interface{}{"error": true, "message": err.Error()}))
	}
	UploadMaturedSecurities(request.Maturities, AuthID(c))
	return c.Render(201, r.JSON(map[string]interface{}{"error": false, "message": "Matured Securities Uploaded"}))
}

func HandleMaturedSecuritiesDelete(c buffalo.Context) error {
	request := OnlyQuarter{}
	if err := c.Bind(&request); err != nil {
		return err
	}
	DeleteMaturedSecurities(request.Quarter)
	return c.Redirect(302, "/sec")
}

func HandleLoadMaturedSecuritiesJSON(c buffalo.Context) error {
	request := OnlyQuarter{}
	if err := c.Bind(&request); err != nil {
		return c.Render(200, r.JSON(map[string]interface{}{"error": true, "message": err.Error()}))
	}
	return c.Render(200, r.JSON(map[string]interface{}{"error": false, "message": "matured securities loaded", "data": map[string]interface{}{"maturities": LoadMaturedSecurities(request.Quarter)}}))
}

func ExportMaturedSecuritiesToExcel(c buffalo.Context) error {
	request := OnlyQuarter{}
	if err := c.Bind(&request); err != nil {
		return err
	}
	c.Response().Header().Set("Content-Disposition", "attachment; filename=matured-securities.xlsx")
	c.Response().Header().Set("Content-Transfer-Encoding", "binary")
	return c.Render(200, r.Func("application/octet-stream", func(w io.Writer, d render.Data) error {
		date := request.Quarter
		if date == "" {
			date = GetQuarterFormalDate()
		}
		data := LoadMaturedSecurities(date)

		excel := excelize.NewFile()
		index := excel.NewSheet("Sheet1")

		excel.SetCellValue("Sheet1", "A1", "CLIENT")
		excel.SetCellValue("Sheet1", "B1", "ISSUER")
		excel.SetCellValue("Sheet1", "C1", "AMOUNT INVESTED")
		excel.SetCellValue("Sheet1", "D1", "VALUE AS AT "+request.Quarter+" (GHS)")

		for index, datum := range data {
			excel.SetCellValue("Sheet1", fmt.Sprintf("A%d", index+2), datum.Client)
			excel.SetCellValue("Sheet1", fmt.Sprintf("B%d", index+2), datum.Issuer)
			excel.SetCellValue("Sheet1", fmt.Sprintf("C%d", index+2), FormatWithComma(datum.AmountInvested, 2))
			excel.SetCellValue("Sheet1", fmt.Sprintf("D%d", index+2), FormatWithComma(datum.Value, 2))
		}
		excel.SetActiveSheet(index)
		bytes, err := excel.WriteToBuffer()
		if err != nil {
			return err
		}
		_, writeError := w.Write(bytes.Bytes())
		return writeError
	}))
}

func HandleLetterToClients(c buffalo.Context) error {
	c.Set("this_quarter_has_been_sent", ClientLettersHaveBeenSent())
	c.Set("logs", LoadClientEmailLog())
	return c.Render(200, r.HTML("client-letters.html"))
}
