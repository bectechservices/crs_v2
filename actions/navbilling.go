package actions

import (
	"fmt"
	"io"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/jinzhu/now"
)

// ShowNavBillDashboardPage display Nav Billing dashboard
func ShowNavBillPage(c buffalo.Context) error {
	date, err := time.Parse("2006-01-02", c.Param("period"))
	if err != nil {
		return err
	}
	c.Set("report_date", date)
	c.Set("transaction_details", LoadTransactionDetails(c.Param("bpOrSca"), date))
	c.Set("currency_details", LoadCurrencyDetails(c.Param("bpOrSca"), date))
	c.Set("invoice_summary", CalculateClientInvoiceSummary(c.Param("bpOrSca"), date))
	var client Client
	client = GetClientByBPID(c.Param("bpOrSca"))
	if client.ID == 0 {
		scaClient := GetClientByCode(c.Param("bpOrSca"))
		client = GetClientByBPID(scaClient.BPID)
		client.BPID = c.Param("bpOrSca")
	}
	c.Set("client", client)
	details := GetClientNAVDetails(c.Param("bpOrSca"), date)
	c.Set("nav_details", details)
	c.Set("invoice_amount", CalculateInvoiceAmountForClient(c.Param("bpOrSca"), details.NAV, date))
	c.Set("invoice_approved", GetClientBillingReportInfo(c.Param("bpOrSca"), date).Approved)
	return c.Render(200, r.HTML("billing.html"))
}

func ShowNavBillDashboardPage(c buffalo.Context) error {
	c.Set("recent_activities", LoadBillingRecentActivities())
	return c.Render(200, r.HTML("billing-dashboard.html"))
}

func HandleTransactionDetailsUpload(c buffalo.Context) error {
	type Request struct {
		Transactions []BillingTransactionDetails `json:"transactions"`
	}
	request := Request{}
	if err := c.Bind(&request); err != nil {
		return c.Render(200, r.JSON(map[string]interface{}{"error": true, "message": err.Error()}))
	}
	UploadTransactionDetails(request.Transactions)
	return c.Render(201, r.JSON(map[string]interface{}{"error": false, "message": "Data uploaded"}))
}

func HandleCurrencyDetailsUpload(c buffalo.Context) error {
	request := BillingCurrencyDetails{}
	if err := c.Bind(&request); err != nil {
		return c.Render(200, r.JSON(map[string]interface{}{"error": true, "message": err.Error()}))
	}
	UploadCurrencyDetails(request)
	return c.Render(201, r.JSON(map[string]interface{}{"error": false, "message": "Data uploaded"}))
}

func HandleLoadClientPosition(c buffalo.Context) error {
	request := &BPOrSCA{}
	if err := c.Bind(request); err != nil {
		return c.Render(200, r.JSON(map[string]interface{}{"error": true, "message": err.Error()}))
	}
	parsedDate, _ := time.Parse("2006-01", request.Date)
	date := now.New(parsedDate).EndOfMonth()
	return c.Render(200, r.JSON(map[string]interface{}{"error": false, "message": "Client position loaded", "data": map[string]interface{}{"position": GetClientPositionByBPID(request.BpOrSca, date)}}))
}

func HandleUpdateClientsNAV(c buffalo.Context) error {
	request := BillingNAVUpdateRequest{}
	if err := c.Bind(&request); err != nil {
		return c.Render(200, r.JSON(map[string]interface{}{"error": true, "message": err.Error()}))
	}
	UpdateClientNAV(request, AuthID(c))
	return c.Render(200, r.JSON(map[string]interface{}{"error": false, "message": "Data uploaded"}))
}

func HandleDeleteBillingUploadedPV(c buffalo.Context) error {
	request := DeleteUploadedPVRequest{}
	if err := c.Bind(&request); err != nil {
		return c.Render(200, r.JSON(map[string]interface{}{"error": true, "message": err.Error()}))
	}
	DeleteBillingPV(request)
	return c.Render(200, r.JSON(map[string]interface{}{"error": false, "message": "PV deleted"}))
}

func ShowBillingReportPage(c buffalo.Context) error {
	parsedDate, err := time.Parse("2006-01", c.Param("period"))
	if err != nil {
		return err
	}
	hack := now.New(parsedDate).EndOfMonth().Day()
	date, _ := time.Parse("2006-01-02", fmt.Sprintf("%s-%d", c.Param("period"), hack))
	c.Set("report_date", date)
	c.Set("transaction_details", LoadTransactionDetails(c.Param("bpOrSca"), date))
	c.Set("currency_details", LoadCurrencyDetails(c.Param("bpOrSca"), date))
	c.Set("invoice_summary", CalculateClientInvoiceSummary(c.Param("bpOrSca"), date))
	var client Client
	client = GetClientByBPID(c.Param("bpOrSca"))
	if client.ID == 0 {
		scaClient := GetClientByCode(c.Param("bpOrSca"))
		client = GetClientByBPID(scaClient.BPID)
		client.BPID = c.Param("bpOrSca")
	}
	c.Set("client", client)
	details := GetClientNAVDetails(c.Param("bpOrSca"), date)
	c.Set("nav_details", details)
	c.Set("invoice_amount", CalculateInvoiceAmountForClient(c.Param("bpOrSca"), details.NAV, date))
	return c.Render(200, r.HTML("billing-report.html"))
}

func HandleLoadClientNAV(c buffalo.Context) error {
	request := &BPOrSCA{}
	if err := c.Bind(request); err != nil {
		return c.Render(200, r.JSON(map[string]interface{}{"error": true, "message": err.Error()}))
	}
	parsedDate, _ := time.Parse("2006-01", request.Date)
	date := now.New(parsedDate).EndOfMonth()
	return c.Render(200, r.JSON(map[string]interface{}{"error": false, "message": "Client position loaded", "data": map[string]interface{}{"nav": GetClientNAVByBPOrSca(request.BpOrSca, date), "report": GetClientBillingReportInfo(request.BpOrSca, date)}}))
}

func HandleNAVMonthlyReportDisapproval(c buffalo.Context) error {
	request := &CheckerCommentRequestWithDateAndClientID{}
	if err := c.Bind(request); err != nil {
		return err
	}
	DisapproveBillingReport(AuthID(c), request.Date, request.Comment, request.BPOrSCA)
	return c.Redirect(302, "/billing?bpOrSca="+request.BPOrSCA+"&period="+request.Date)
}

func HandleNAVMonthlyReportApproval(c buffalo.Context) error {
	request := &CheckerCommentRequestWithDateAndClientID{}
	if err := c.Bind(request); err != nil {
		return err
	}
	ApproveBillingReport(AuthID(c), request.Date, request.BPOrSCA)
	return c.Redirect(302, "/billing?bpOrSca="+request.BPOrSCA+"&period="+request.Date)
}

func HandleNAVMonthlyReportReverseApproval(c buffalo.Context) error {
	request := &CheckerCommentRequestWithDateAndClientID{}
	if err := c.Bind(request); err != nil {
		return err
	}
	ReverseBillingReportApproval(AuthID(c), request.Date, request.BPOrSCA)
	return c.Redirect(302, "/billing?bpOrSca="+request.BPOrSCA+"&period="+request.Date)
}

func HandleBillingReportExport(c buffalo.Context) error {
	request := BPOrSCA{}
	if err := c.Bind(&request); err != nil {
		return err
	}
	filename := fmt.Sprintf("%s-%s-invoice.docx", request.BpOrSca, request.Date)
	c.Response().Header().Set("Content-Disposition", "attachment; filename="+filename)
	return c.Render(200, r.Func("application/vnd.openxmlformats-officedocument.wordprocessingml.document", func(w io.Writer, d render.Data) error {
		payload := BillingReportExportPayload{}
		var client Client
		client = GetClientByBPID(request.BpOrSca)
		if client.ID == 0 {
			scaClient := GetClientByCode(request.BpOrSca)
			client = GetClientByBPID(scaClient.BPID)
			client.BPID = c.Param(request.BpOrSca)
		}
		payload.Client = client
		date, err := time.Parse("2006-01-02", request.Date)
		if err != nil {
			return err
		}
		payload.Period = date.Format("02 January 2006")
		details := GetClientNAVDetails(request.BpOrSca, date)
		payload.InvoiceReference = details.InvoiceReference
		payload.InvoiceAmount = CalculateInvoiceAmountForClient(request.BpOrSca, details.NAV, date)
		payload.CurrencyDetails = LoadCurrencyDetails(request.BpOrSca, date)
		payload.Summary = CalculateClientInvoiceSummary(request.BpOrSca, date)
		payload.TransactionDetails = LoadTransactionDetails(request.BpOrSca, date)

		data, err := ExportBillingReport(payload)
		if err != nil {
			return err
		}
		_, writeError := w.Write(data)
		return writeError
	}))
}

func HandleLoadBilledClients(c buffalo.Context) error {
	bpid := c.Param("bpid")
	month := c.Param("month")
	year := c.Param("year")
	var date time.Time
	var clients []BilledClients
	if bpid == "" && month == "" && year == "" {
		date = now.New(time.Now().AddDate(0, -1, -5)).EndOfMonth()
		clients = LoadBilledClients(date)
		c.Set("is_sca_report", false)
	} else if bpid != "" {
		parsedDate, _ := time.Parse("2006-01", year+"-"+month)
		date = now.New(parsedDate).EndOfMonth()
		clients = LoadBilledClientSCAs(bpid, date)
		c.Set("is_sca_report", true)
	} else {
		parsedDate, _ := time.Parse("2006-01", year+"-"+month)
		date = now.New(parsedDate).EndOfMonth()
		clients = LoadBilledClients(date)
		c.Set("is_sca_report", false)
	}
	c.Set("clients", clients)
	c.Set("bpid", c.Param("bpid"))
	c.Set("month", fmt.Sprintf("%02d", int(date.Month())))
	c.Set("year", date.Year())
	return c.Render(200, r.HTML("billed-clients.html"))
}
