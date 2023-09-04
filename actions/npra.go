package actions

import (
	"fmt"
	"io"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/jinzhu/now"
)

//@TODO:NPRA DASHBOARD

// ShowNpraDashboard displays the users page
func showNpraDashboard(c buffalo.Context) error {
	c.Set("recent_activities", LoadNPRARecentActivities())
	c.Set("trends", LoadNPRATrendsForPastFourQuarters())
	return c.Render(200, r.HTML("npra-dashboard.html"))
}

// Generate NPRA Report displays the users page
func generateNpraReport(c buffalo.Context) error {
	c.Set("npra_report_approved", LoadCurrentNPRADeclaration().Approved)
	return c.Render(200, r.HTML("npra.html"))
}

// ShowNpraPreview displays the NPRA Quarterly Previw Page
func ShowNpraPreview(c buffalo.Context) error {
	c.Set("npra_report_approved", LoadCurrentNPRADeclaration().Approved)
	currentQuarter := GetQuarterFormalDate()
	c.Set("local_variance", LoadNpraLocalVariance(MakeLastQuarterFormalDate(currentQuarter), currentQuarter))
	return c.Render(200, r.HTML("npra-preview.html"))
}

// HandleUnauthorizedTransactionUpload handles the request to upload unauthorized txns
func HandleUnauthorizedTransactionUpload(c buffalo.Context) error {
	type RequestData struct {
		Transactions []NPRAUnauthorizedTransactionRequest `json:"transactions"`
	}
	request := RequestData{}
	if err := c.Bind(&request); err != nil {
		return c.Render(200, r.JSON(map[string]interface{}{"error": true, "message": err.Error()}))
	}
	InsertUnauthorizedTransactionsIntoDB(request.Transactions)
	return c.Render(201, r.JSON(map[string]interface{}{"error": false, "message": "Data uploaded"}))
}

func HandleOutstandingFDCertificateUpload(c buffalo.Context) error {
	type RequestData struct {
		Certificates []NPRAOutstandingFDCertificateRequest `json:"certificates"`
	}
	request := RequestData{}
	if err := c.Bind(&request); err != nil {
		return c.Render(200, r.JSON(map[string]interface{}{"error": true, "message": err.Error()}))
	}
	InsertOutstandingFDCertificatesIntoDB(request.Certificates, AuthID(c))
	return c.Render(201, r.JSON(map[string]interface{}{"error": false, "message": "Data uploaded"}))
}

func HandleLoadUnauthorizedTransactions(c buffalo.Context) error {
	request := &QuarterAndYear{}
	if err := c.Bind(request); err != nil {
		return c.Render(200, r.JSON(map[string]interface{}{"error": true, "message": err.Error()}))
	}
	transactions := LoadUnauthorizedTransactions(request.Quarter, request.Year)
	return c.Render(200, r.JSON(map[string]interface{}{"error": false, "message": "Transactions loaded", "data": map[string]interface{}{
		"transactions": transactions,
	}}))
}

func HandleLoadOutstandingFDCertificates(c buffalo.Context) error {
	request := &QuarterAndYear{}
	if err := c.Bind(request); err != nil {
		return c.Render(200, r.JSON(map[string]interface{}{"error": true, "message": err.Error()}))
	}
	certificates := LoadOutstandingFDCertificates(request.Quarter, request.Year)
	return c.Render(200, r.JSON(map[string]interface{}{"error": false, "message": "Certificates loaded", "data": map[string]interface{}{
		"certificates": certificates,
	}}))
}

func HandleLoadNPRAQuarterlyReport(c buffalo.Context) error {
	request := &QuarterAndYear{}
	if err := c.Bind(request); err != nil {
		return c.Render(200, r.JSON(map[string]interface{}{"error": true, "message": err.Error()}))
	}
	pensions := LoadPensionFund(request.Quarter, request.Year, true)
	provident := LoadProvidentFund(request.Quarter, request.Year, true)
	administrators := LoadFundAdministrators(MakeDBQueryableQuarterLastDate(request.Quarter, request.Year))
	return c.Render(200, r.JSON(map[string]interface{}{"error": false, "message": "NPRA quarterly report loaded", "data": map[string]interface{}{
		"pensions":       pensions,
		"provident":      provident,
		"administrators": administrators,
	}}))
}

func HandleLoadCurrentNPRADeclaration(c buffalo.Context) error {
	return c.Render(200, r.JSON(map[string]interface{}{"error": false, "message": "Declaration loaded", "data": map[string]interface{}{"declaration": LoadCurrentNPRADeclaration()}}))
}

func HandleUploadNPRADeclaration(c buffalo.Context) error {
	request := NPRADeclaration{}
	if err := c.Bind(&request); err != nil {
		return c.Render(200, r.JSON(map[string]interface{}{"error": true, "message": err.Error()}))
	}
	UpdateCurrentNPRADeclaration(request)
	return c.Render(201, r.JSON(map[string]interface{}{"error": false, "message": "Declaration uploaded"}))
}

func HandleNPRASendProgressToChecker(c buffalo.Context) error {
	SendReportActionEmail("report_prepared", LoadUsersByRole(NPRACheckerRole), "NPRA Report", "NPRA", "")
	return c.Render(201, r.JSON(map[string]interface{}{"error": false, "message": "Data uploaded"}))
}

func HandleNPRACurrentQuarterReportDisapproval(c buffalo.Context) error {
	request := &CheckerCommentRequest{}
	if err := c.Bind(request); err != nil {
		return err
	}
	DisapproveNPRACurrentQuarterReport(AuthID(c), request.Comment)
	return c.Redirect(302, "/npra-preview")
}

func HandleNPRAReportApproval(c buffalo.Context) error {
	approved := ApproveNPRAReport(AuthID(c))
	if approved == 0 {
		c.Flash().Add("error", "All input from MAKER is required.")
	}
	return c.Redirect(302, "/npra-preview")
}

func HandleOutstandingFDCertificateEdit(c buffalo.Context) error {
	type Request struct {
		Data []OutstandingFDEdit `json:"data"`
	}
	request := Request{}
	if err := c.Bind(&request); err != nil {
		return c.Render(200, r.JSON(map[string]interface{}{"error": true, "message": err.Error()}))
	}
	UpdateOutstandingFDCertificates(request.Data)
	return c.Render(201, r.JSON(map[string]interface{}{"error": false, "message": "Outstanding receipt data edited"}))
}

func HandleDeleteNPRAUploadedPV(c buffalo.Context) error {
	request := DeleteUploadedPVRequest{}
	if err := c.Bind(&request); err != nil {
		return c.Render(200, r.JSON(map[string]interface{}{"error": true, "message": err.Error()}))
	}
	DeleteNPRAPV(request)
	return c.Render(200, r.JSON(map[string]interface{}{"error": false, "message": "PV deleted"}))
}

func ShowUploadedOutstandingFDCertificates(c buffalo.Context) error {
	c.Set("outstanding_fd_certs", LoadOutstandingFDCertificatesByHash(c.Param("hash")))
	return c.Render(200, r.HTML("fd-certs.html"))
}

func HandleOutstandingFDCertificatesDelete(c buffalo.Context) error {
	request := OnlyHash{}
	if err := c.Bind(&request); err != nil {
		return err
	}
	DeleteOutstandingFDCertificates(request.Hash)
	return c.Redirect(302, "/npra-dashboard")
}

// ExportNPRALocalVarianceToExcel exports the foreign variance data to excel
func ExportNPRALocalVarianceToExcel(c buffalo.Context) error {
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

		data := LoadNpraLocalVariance(MakeLastQuarterFormalDate(currentQuarter), currentQuarter)

		bytes, err := ExportVarianceToExcel(previousQuarterShortDate, currentQuarterShortDate, data)
		if err != nil {
			return err
		}
		_, writeError := w.Write(bytes.Bytes())
		return writeError
	}))
}

// ExportNPRAOutstandingFDToExcel exports the outstanding fd data to excel
func ExportNPRAOutstandingFDToExcel(c buffalo.Context) error {
	request := &QuarterAndYear{}
	if err := c.Bind(request); err != nil {
		return err
	}
	c.Response().Header().Set("Content-Disposition", "attachment; filename=local-variance.xlsx")
	c.Response().Header().Set("Content-Transfer-Encoding", "binary")
	return c.Render(200, r.Func("application/octet-stream", func(w io.Writer, d render.Data) error {
		bytes, err := ExportOutstandingFDsToExcel(LoadOutstandingFDCertificates(request.Quarter, request.Year))
		if err != nil {
			return err
		}
		_, writeError := w.Write(bytes.Bytes())
		return writeError
	}))
}

// ExportNPRAMonthlyReportToExcel exports the monthly report data to excel
func ExportNPRAMonthlyReportToExcel(c buffalo.Context) error {
	request := &QuarterAndYear{}
	if err := c.Bind(request); err != nil {
		return err
	}
	c.Response().Header().Set("Content-Disposition", "attachment; filename=local-variance.xlsx")
	c.Response().Header().Set("Content-Transfer-Encoding", "binary")
	return c.Render(200, r.Func("application/octet-stream", func(w io.Writer, d render.Data) error {
		pensions := LoadPensionFund(request.Quarter, request.Year, false)
		provident := LoadProvidentFund(request.Quarter, request.Year, false)
		administrators := LoadFundAdministrators(MakeDBQueryableQuarterLastDate(request.Quarter, request.Year))
		date := MakeQuarterDate(request.Quarter, request.Year)
		excel := excelize.NewFile()
		index := excel.NewSheet("Summary")

		length := len(pensions)
		pensionsLength := len(pensions)
		providentLength := len(provident)
		if providentLength > pensionsLength {
			length = providentLength
		}
		excelRowCount := 1 //to skip the headings

		excel.SetCellValue("Summary", "A1", "")
		excel.SetCellValue("Summary", "B1", "")
		excel.SetCellValue("Summary", "C1", "")
		excel.SetCellValue("Summary", "D1", "")
		excel.SetCellValue("Summary", "A2", "Pension Fund(Tier 2)")
		excel.SetCellValue("Summary", "B2", fmt.Sprintf("AUA As At %s", date))
		excel.SetCellValue("Summary", "C2", "Provident Fund(Tier 3)")
		excel.SetCellValue("Summary", "D2", fmt.Sprintf("AUA As At %s", date))
		excel.SetRowHeight("Summary", 2, 25)

		headingColor, err := excel.NewStyle(`{"fill":{"type":"pattern","color":["#95b3d7"],"pattern":1}}`)
		if err != nil {
			return err
		}
		excel.SetCellStyle("Summary", "A2", "A2", headingColor)
		excel.SetCellStyle("Summary", "B2", "B2", headingColor)
		excel.SetCellStyle("Summary", "C2", "C2", headingColor)
		excel.SetCellStyle("Summary", "D2", "D2", headingColor)

		for i := 0; i < length; i++ {
			if i < pensionsLength {
				excel.SetCellValue("Summary", fmt.Sprintf("A%d", excelRowCount+2), pensions[i].Name)
				excel.SetCellValue("Summary", fmt.Sprintf("B%d", excelRowCount+2), FormatWithComma(pensions[i].Value, 2))
			} else {
				excel.SetCellValue("Summary", fmt.Sprintf("A%d", excelRowCount+2), "")
				excel.SetCellValue("Summary", fmt.Sprintf("B%d", excelRowCount+2), "")
			}
			if i < providentLength {
				excel.SetCellValue("Summary", fmt.Sprintf("C%d", excelRowCount+2), provident[i].Name)
				excel.SetCellValue("Summary", fmt.Sprintf("D%d", excelRowCount+2), FormatWithComma(provident[i].Value, 2))
			} else {
				excel.SetCellValue("Summary", fmt.Sprintf("C%d", excelRowCount+2), "")
				excel.SetCellValue("Summary", fmt.Sprintf("D%d", excelRowCount+2), "")
			}
			if i == length {
				excel.SetCellValue("Summary", fmt.Sprintf("A%d", excelRowCount+2), "")
				excel.SetCellValue("Summary", fmt.Sprintf("B%d", excelRowCount+2), "")
				excel.SetCellValue("Summary", fmt.Sprintf("C%d", excelRowCount+2), "")
				excel.SetCellValue("Summary", fmt.Sprintf("D%d", excelRowCount+2), "")

				excel.SetCellValue("Summary", fmt.Sprintf("A%d", excelRowCount+3), "")
				excel.SetCellValue("Summary", fmt.Sprintf("B%d", excelRowCount+3), "")
				excel.SetCellValue("Summary", fmt.Sprintf("C%d", excelRowCount+3), "")
				excel.SetCellValue("Summary", fmt.Sprintf("D%d", excelRowCount+3), "")

				excelRowCount += 2
			} else {
				excelRowCount += 1
			}
		}
		for i := 0; i < 2; i++ {
			excel.SetCellValue("Summary", fmt.Sprintf("A%d", excelRowCount+2), "")
			excel.SetCellValue("Summary", fmt.Sprintf("B%d", excelRowCount+2), "")
			excel.SetCellValue("Summary", fmt.Sprintf("C%d", excelRowCount+2), "")
			excel.SetCellValue("Summary", fmt.Sprintf("D%d", excelRowCount+2), "")
			excelRowCount += 1
		}
		excel.SetCellValue("Summary", fmt.Sprintf("A%d", excelRowCount+2), "Total number of  Tier 2  clients")
		excel.SetCellValue("Summary", fmt.Sprintf("B%d", excelRowCount+2), len(pensions))
		excelRowCount += 1
		excel.SetCellValue("Summary", fmt.Sprintf("A%d", excelRowCount+2), "Total number of  Tier 3  clients")
		excel.SetCellValue("Summary", fmt.Sprintf("B%d", excelRowCount+2), len(provident))
		excelRowCount += 1

		for i := 0; i < 4; i++ {
			excel.SetCellValue("Summary", fmt.Sprintf("A%d", excelRowCount+2), "")
			excel.SetCellValue("Summary", fmt.Sprintf("B%d", excelRowCount+2), "")
			excel.SetCellValue("Summary", fmt.Sprintf("C%d", excelRowCount+2), "")
			excel.SetCellValue("Summary", fmt.Sprintf("D%d", excelRowCount+2), "")
			excelRowCount += 1
		}
		excel.SetCellValue("Summary", fmt.Sprintf("A%d", excelRowCount+2), "Fund Name")
		excel.SetCellValue("Summary", fmt.Sprintf("B%d", excelRowCount+2), "Administrators")
		excelRowCount += 1
		for _, admin := range administrators {
			excel.SetCellValue("Summary", fmt.Sprintf("A%d", excelRowCount+2), admin.Name)
			excel.SetCellValue("Summary", fmt.Sprintf("B%d", excelRowCount+2), admin.Administrator)
			excelRowCount += 1
		}

		excel.NewSheet("Charges")

		//excel.MergeCell("Charges", "B1", "D1")
		excel.SetCellValue("Charges", "A1", "Clients")
		excel.SetCellValue("Charges", "B1", "Charges")
		chargesRowCount := 0
		billingInfo := GetAllBillingDetails()
		clientNameColor, err := excel.NewStyle(`{"fill":{"type":"pattern","color":["#ffbf00"],"pattern":1}}`)
		if err != nil {
			return err
		}
		for _, info := range billingInfo {
			excel.SetCellValue("Charges", fmt.Sprintf("A%d", chargesRowCount+2), info.ClientName)
			excel.SetCellStyle("Charges", fmt.Sprintf("A%d", chargesRowCount+2), fmt.Sprintf("A%d", chargesRowCount+2), clientNameColor)
			chargesRowCount += 1

			excel.SetCellValue("Charges", fmt.Sprintf("A%d", chargesRowCount+2), "")
			excel.SetCellValue("Charges", fmt.Sprintf("B%d", chargesRowCount+2), "Transaction")
			excel.SetCellValue("Charges", fmt.Sprintf("C%d", chargesRowCount+2), fmt.Sprintf("GHS %s", FormatWithComma(info.ChargePerTransaction, 2)))
			excel.SetCellValue("Charges", fmt.Sprintf("D%d", chargesRowCount+2), "Per Transaction")
			chargesRowCount += 1

			if len(info.BasisPoints) == 1 {
				min := int64(info.BasisPoints[0].Minimum)
				max := int64(info.BasisPoints[0].Maximum)
				if min == 0 && max == 0 {
					excel.SetCellValue("Charges", fmt.Sprintf("A%d", chargesRowCount+2), "")
					excel.SetCellValue("Charges", fmt.Sprintf("B%d", chargesRowCount+2), "Portfolio Fee")
					excel.SetCellValue("Charges", fmt.Sprintf("C%d", chargesRowCount+2), fmt.Sprintf("%.f bpts", info.BasisPoints[0].BasisPoints))
					excel.SetCellValue("Charges", fmt.Sprintf("D%d", chargesRowCount+2), "")
					chargesRowCount += 1
				}
			} else {
				for _, bpInfo := range info.BasisPoints {
					minimum := int64(bpInfo.Minimum)
					maximum := int64(bpInfo.Maximum)
					var minVal int64
					var maxVal int64
					var bpStr string
					if minimum > 0 {
						minVal = minimum / 1000000
					}
					if maximum > 0 {
						maxVal = maximum / 1000000
					}
					if minVal > 0 && maxVal > 0 {
						bpStr = fmt.Sprintf("%d-%d mil", minVal, maxVal)
					} else {
						bpStr = fmt.Sprintf("%d mil+", minVal)
					}
					excel.SetCellValue("Charges", fmt.Sprintf("A%d", chargesRowCount+2), "")
					excel.SetCellValue("Charges", fmt.Sprintf("B%d", chargesRowCount+2), "Portfolio Fee")
					excel.SetCellValue("Charges", fmt.Sprintf("C%d", chargesRowCount+2), bpStr)
					excel.SetCellValue("Charges", fmt.Sprintf("D%d", chargesRowCount+2), fmt.Sprintf("%.f bpts", bpInfo.BasisPoints))
					chargesRowCount += 1
				}
			}

			excel.SetCellValue("Charges", fmt.Sprintf("A%d", chargesRowCount+2), "")
			excel.SetCellValue("Charges", fmt.Sprintf("B%d", chargesRowCount+2), "Third Party Transfer")
			excel.SetCellValue("Charges", fmt.Sprintf("C%d", chargesRowCount+2), fmt.Sprintf("GHS %s", FormatWithComma(info.ThirdPartyTransfer, 2)))
			excel.SetCellValue("Charges", fmt.Sprintf("D%d", chargesRowCount+2), "Flat")
			chargesRowCount += 1

			excel.SetCellValue("Charges", fmt.Sprintf("A%d", chargesRowCount+2), "")
			excel.SetCellValue("Charges", fmt.Sprintf("B%d", chargesRowCount+2), "Minimum Fee")
			excel.SetCellValue("Charges", fmt.Sprintf("C%d", chargesRowCount+2), fmt.Sprintf("GHS %s", FormatWithComma(info.MinimumCharge, 2)))
			excel.SetCellValue("Charges", fmt.Sprintf("D%d", chargesRowCount+2), "Flat")
			chargesRowCount += 1
		}
		excel.DeleteSheet("Sheet1")
		excel.SetActiveSheet(index)

		bytes, err := excel.WriteToBuffer()
		if err != nil {
			return err
		}
		_, writeError := w.Write(bytes.Bytes())
		return writeError
	}))
}

func ExportNPRAMonthlyReportToWord(c buffalo.Context) error {
	c.Response().Header().Set("Content-Disposition", "attachment; filename=npra-monthly.docx")
	return c.Render(200, r.Func("application/vnd.openxmlformats-officedocument.wordprocessingml.document", func(w io.Writer, d render.Data) error {
		data, err := ExportNPRAMonthlyReport(LoadCurrentNPRADeclaration())
		if err != nil {
			return err
		}
		_, writeError := w.Write(data)
		return writeError
	}))
}

func ExportNPRAUnauthorizedReportToWord(c buffalo.Context) error {
	c.Response().Header().Set("Content-Disposition", "attachment; filename=Unauthorized-Report.doc")
	return c.Render(200, r.Func("application/vnd.openxmlformats-officedocument.wordprocessingml.document", func(w io.Writer, d render.Data) error {
		data, err := ExportNPRAUnauthorizedReport()
		if err != nil {
			return err
		}
		_, writeError := w.Write(data)
		return writeError
	}))
}

func HandleUploadNPRAVarianceRemarks(c buffalo.Context) error {
	request := VarianceRemarksUpdateRequest{}
	if err := c.Bind(&request); err != nil {
		return err
	}
	UploadNPRAVarianceRemarks(request)
	return c.Redirect(302, "/npra-preview")
}

func HandleExportUnauthorizedTransactions(c buffalo.Context) error {
	request := &QuarterAndYear{}
	if err := c.Bind(request); err != nil {
		return err
	}
	c.Response().Header().Set("Content-Disposition", "attachment; filename=local-variance.xlsx")
	c.Response().Header().Set("Content-Transfer-Encoding", "binary")
	transactions := LoadUnauthorizedTransactions(request.Quarter, request.Year)
	return c.Render(200, r.Func("application/octet-stream", func(w io.Writer, d render.Data) error {
		excel := excelize.NewFile()

		excel.SetCellValue("Sheet1", "A1", "Client Name")
		excel.SetCellValue("Sheet1", "B1", "Transaction Details")
		excel.SetCellValue("Sheet1", "C1", "Date Of Occurrence")

		if len(transactions) > 0 {
			for index, datum := range transactions {
				excel.SetCellValue("Sheet1", fmt.Sprintf("A%d", index+2), datum.ClientName)
				excel.SetCellValue("Sheet1", fmt.Sprintf("B%d", index+2), datum.TransactionDetails)
				excel.SetCellValue("Sheet1", fmt.Sprintf("C%d", index+2), datum.Date.Format("2006-01-02"))
			}
		} else {
			excel.SetCellValue("Sheet1", fmt.Sprintf("A%d", 2), "Nil")
			excel.SetCellValue("Sheet1", fmt.Sprintf("B%d", 2), "Nil")
			excel.SetCellValue("Sheet1", fmt.Sprintf("C%d", 2), "Nil")
		}

		bytes, err := excel.WriteToBuffer()
		if err != nil {
			return err
		}
		_, writeError := w.Write(bytes.Bytes())
		return writeError
	}))
}

// ShowNPRA030ReportPage NPRA Report displays the 030 Reports
func ShowNPRA030ReportPage(c buffalo.Context) error {
	c.Set("npra030_report_approved", LoadCurrentNPRADeclaration())
	return c.Render(200, r.HTML("npra030.html"))
}

// LoadNPRA0301ReportPage NPRA Report displays the 030 Reports
func LoadNPRA0301ReportPage(c buffalo.Context) error {
	return c.Render(200, r.HTML("npra0301.html"))
}

// // LoadNPRA0301ReportPage NPRA Report displays the 030 Reports
// func HandleNPRA0301ReportData(c buffalo.Context) error {
// 	request := NPRAReportRequest{}
// 	request.BPID = c.Param("bp_id")
// 	request.Type = c.Param("type")
// 	request.Month = c.Param("month")
// 	request.Year = c.Param("year")

// 	npra301Data := GetNPRA0301Data(request)

// 	return c.Render(200, r.JSON(npra301Data))
// }

// LoadNPRA0302ReportPage NPRA Report displays the 030 Reports
func LoadNPRA0302ReportPage(c buffalo.Context) error {
	return c.Render(200, r.HTML("npra0302.html"))
}

func LoadNPRA03021ReportPage(c buffalo.Context) error {
	bpOrSca := c.Param("bpOrSca")
	month := c.Param("month")
	year := c.Param("year")
	var sdate time.Time
	var npra0302s []NPRA0302
	parsedDate, _ := time.Parse("2006-01", year+"-"+month)
	sdate = now.New(parsedDate).EndOfMonth()
	date := sdate.Format("2006-01-02")
	npra0302s = LoadNPRA03021Reports(bpOrSca, date)
	c.Set("npra302s", npra0302s)
	c.Set("bpOrSca", bpOrSca)
	c.Set("month", month)
	c.Set("year", year)
	return c.Render(200, r.HTML("npra03021.html"))
}

func LoadNPRA03011ReportPage(c buffalo.Context) error {
	month := c.Param("month")
	year := c.Param("year")
	var sdate time.Time
	var npra0301s []NPRA0301
	parsedDate, _ := time.Parse("2006-01", year+"-"+month)
	sdate = now.New(parsedDate).EndOfMonth()
	date := sdate.Format("2006-01-02")
	npra0301s = LoadNPRA03011Reports(date)
	c.Set("npra301s", npra0301s)
	c.Set("month", month)
	c.Set("year", year)
	return c.Render(200, r.HTML("npra03011.html"))
}

// LoadNPRA0301ReportPage NPRA Report displays the 030 Reports
func HandleNPRA0302ReportData(c buffalo.Context) error {
	request := NPRA321Request{}
	request.BPID = c.Param("bp_id")
	request.Date = c.Param("date")
	npra302Data := LoadNPRA03021Reports(request.BPID, request.Date)
	return c.Render(200, r.JSON(npra302Data))
}

// ShowNPRA0303Report NPRA Report displays the 030 Reports
func ShowNPRA0303Report(c buffalo.Context) error {
	c.Set("npra303s", LoadNPRA0303Reports())
	c.Set("clients", LoadAllClients())
	return c.Render(200, r.HTML("npra0303.html"))
}

func ShowNPRA0303ReportAdd(c buffalo.Context) error {
	c.Set("clients", LoadAllClients())
	return c.Render(200, r.HTML("npra0303-add.html"))
}

// HandleNPRA0303ReportAdd NPRA Report displays the 030 Reports
func HandleNPRA0303ReportAdd(c buffalo.Context) error {
	request := &NPRA0303AddRequest{}
	if err := c.Bind(request); err != nil {
		return err
	}
	// AddNPRA0303ToDB(*request)
	if err := AddNPRA0303ToDB(*request); err == true {
		c.Flash().Add("success", "NPRA 0303 Client Data Added successfully")
		return c.Redirect(302, "/npra0303")
	}
	c.Flash().Add("Failure", "NPRA 0303 Client Data failed")
	return c.Redirect(302, "/npra0303")
}

// ShowClientNPRA0303 loads the client edit page
func ShowNPRA0303Edit(c buffalo.Context) error {
	c.Set("npra303", GetNPRA0303ByID(c.Param("id")))
	c.Set("npra303s", LoadNPRA0303Reports())
	return c.Render(200, r.HTML("npra0303-edit.html"))
}

func HandleNPRA0303Edit(c buffalo.Context) error {
	request := &NPRA0303EditRequest{}
	if err := c.Bind(request); err != nil {
		return err
	}

	EditNPRA0303(*request)
	c.Flash().Add("success", "npra 0303 edited successfully")
	return c.Redirect(302, "/npra0303")
}

// ShowClientDetails displays the Client details page
func ShowNPRA0303Details(c buffalo.Context) error {
	c.Set("npra303", GetNPRA0303ByID(c.Param("id")))
	c.Set("npra303s", LoadNPRA0303Reports())
	return c.Render(200, r.HTML("npra0303-view.html"))
}

// GetClientByID loads a NPRA0303 by their BPID
func GetNPRA0303ByID(ID string) NPRA0303 {
	npra0303 := NPRA0303{}
	DatabaseConnection.Where("ID=?", ID).Find(&npra0303)
	return npra0303
}

func NPRA0302DataExists(bpId, reportDate string) bool {
	npra302Data := []NPRA0302{}

	DatabaseConnection.Raw("select * from CRS_0302_NPRA_REPORT where bp_id = ? and reporting_date = ?", bpId, reportDate).Scan(&npra302Data)

	return len(npra302Data) > 1
}

func NPRA0301DataExists(bpId, reportDate string) bool {
	npra301Data := []NPRA0301{}

	DatabaseConnection.Raw("select * from CRS_0301_NPRA_REPORT where bp_id = ? and reporting_date = ?",
		bpId, reportDate).Scan(&npra301Data)

	return len(npra301Data) > 1
}

// LoadUploadedPVForEyeBalling loads the pv data for eyeballing
func Loadnpra030PVForEyeBalling(c buffalo.Context) error {
	request := &Fetch0302ForEyeBallingRequest{}
	if err := c.Bind(request); err != nil {
		return c.Render(200, r.JSON(map[string]interface{}{"error": true, "message": err.Error()}))
	}
	data := Loadnpra0302Data(request.BPID, request.QuarterDate)
	return c.Render(200, r.JSON(map[string]interface{}{"error": false, "message": "PV details loaded", "data": map[string]interface{}{"reportData": data}}))
}

// LoadUploadedPVForEyeBalling loads the pv data for eyeballing
func Load301ForEyeBalling(c buffalo.Context) error {
	request := &Fetch0301ForEyeBallRequest{}
	if err := c.Bind(request); err != nil {
		return c.Render(200, r.JSON(map[string]interface{}{"error": true, "message": err.Error()}))
	}
	data := Loadnpra0301Data(request.QuarterDate)
	return c.Render(200, r.JSON(map[string]interface{}{"error": false, "message": "PV details loaded", "data": map[string]interface{}{"reportData": data}}))
}

// Export301ReportToExcel exports the NPRA 0301 data to excel
func Export301ReportToExcel(c buffalo.Context) error {
	request := &MonthAndYear{}

	month := c.Param("month")
	year := c.Param("year")
	c.Set("month", month)
	c.Set("year", year)
	var sdate time.Time

	parsedDate, _ := time.Parse("2006-01", request.Year+"-"+request.Month)
	sdate = now.New(parsedDate).EndOfMonth()

	if err := c.Bind(request); err != nil {
		return err
	}
	c.Response().Header().Set("Content-Disposition", "attachment; filename=local-variance.xlsx")
	c.Response().Header().Set("Content-Transfer-Encoding", "binary")
	return c.Render(200, r.Func("application/octet-stream", func(w io.Writer, d render.Data) error {

		date := sdate.Format("2006-01-02")
		data := Load301data(date)

		bytes, err := Export301ToExcel(data)
		if err != nil {
			return err
		}
		_, writeError := w.Write(bytes.Bytes())
		return writeError
	}))
}

// Export302ReportToExcel exports the NPRA 0302 data to excel
func Export302ReportToExcel(c buffalo.Context) error {
	request := &BPAndMonthAndYear{}
	month := c.Param("month")
	bpOrSca := c.Param("bpOrSca")
	year := c.Param("year")
	c.Set("month", month)
	c.Set("bpOrSca", bpOrSca)
	c.Set("year", year)
	var sdate time.Time
	fmt.Println("############## Level 1 date @@@@@@@@@@@@@@@@@@@@@", month)
	fmt.Println("############## Level 2 DATE @@@@@@@@@@@@@@@@@@@@@", year)
	//---------------------------
	parsedDate, _ := time.Parse("2006-01", request.Year+"-"+request.Month)
	sdate = now.New(parsedDate).EndOfMonth()
	fmt.Println("############## Level 3 date @@@@@@@@@@@@@@@@@@@@@", parsedDate)
	fmt.Println("############## Level 4 DATE @@@@@@@@@@@@@@@@@@@@@", sdate)
	if err := c.Bind(request); err != nil {
		return err
	}
	c.Response().Header().Set("Content-Disposition", "attachment; filename=local-variance.xlsx")
	c.Response().Header().Set("Content-Transfer-Encoding", "binary")
	return c.Render(200, r.Func("application/octet-stream", func(w io.Writer, d render.Data) error {

		date := sdate.Format("2006-01-02")
		fmt.Println("@@@@@@@@@@@@@@LOOKING for DATE @@@@@@@@@@@@@@@@@@@@@", date)
		data := Load302data(date, bpOrSca)
		fmt.Println("@@@@@@@@@@@@@@LOOKING for DATA for Date @@@@@@@@@@@@@@@@@@@@@", data)
		bytes, err := Export302ToExcel(data)
		if err != nil {
			return err
		}
		_, writeError := w.Write(bytes.Bytes())
		return writeError
	}))
}
