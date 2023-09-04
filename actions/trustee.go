package actions

import (
	"github.com/gobuffalo/buffalo"
)

// ShowTrusteeDashboardPage display trustee dashboard
func ShowTrusteeDashboardPage(c buffalo.Context) error {
	date := MakeDBQueryableCurrentQuarterDate()
	c.Set("recent_activities", LoadTrusteeRecentActivities())
	c.Set("t2_clients", LoadTier2Clients(date))
	c.Set("t3_clients", LoadTier3Clients(date))
	return c.Render(200, r.HTML("trustee-dashboard.html"))
}

// ShowUnidentifiedReport display unidentified report page
func ShowUnidentifiedReport(c buffalo.Context) error {
	c.Set("clients", LoadAllClients())
	return c.Render(200, r.HTML("gen-unidentified-report.html"))
}

// ShowTemplateSetup display Trustee Presentation setup page
func ShowTemplateSetup(c buffalo.Context) error {
	return c.Render(200, r.HTML("template-setup.html"))
}

//LoadTrusteeData loads the data for the trustee
func LoadTrusteeData(c buffalo.Context) error {
	request := &TrusteePVDataRequest{}
	if err := c.Bind(request); err != nil {
		return c.Render(200, r.JSON(map[string]interface{}{"error": true, "message": err.Error()}))
	}
	pvSummary := LoadClientPVSummary(request.BPID, MakeQuarterFormalDate(request.Quarter, request.Year))
	performance := LoadTrusteeQuarterlyPerformance(request.BPID, request.Quarter, request.Year)
	contributions := LoadMonthlyContributions(request.BPID, request.Quarter, request.Year)
	client := GetClientByBPID(request.BPID)
	if client.ID == 0 {
		scaClient := GetClientByCode(request.BPID)
		client = GetClientByBPID(scaClient.BPID)
		client.BPID = request.BPID
	}
	scas := GetClientSCASByBPID(request.BPID)
	unidentifiedPayments := LoadUnidentifiedPayments(request.BPID, request.Year)
	unidentifiedPaymentsSummary := LoadUnidentifiedPaymentSummary(request.BPID, request.Year)
	gog := LoadClientGOGMaturities(request.BPID, request.Quarter, request.Year)
	transactions := LoadClientTransactions(request.BPID, request.Quarter, request.Year)
	return c.Render(200, r.JSON(map[string]interface{}{
		"error":   false,
		"message": "PV data loaded successfully",
		"data": map[string]interface{}{
			"summary":                     SortPVSummary(LumpPVSummary(pvSummary)),
			"performance":                 LumpTrusteeQuarterlyPerformance(performance),
			"contributions":               contributions,
			"client":                      client,
			"unidentifiedPayments":        unidentifiedPayments,
			"unidentifiedPaymentsSummary": unidentifiedPaymentsSummary,
			"scas":                        scas,
			"gog":                         gog,
			"transactions":                transactions,
			"report":                      LoadClientTrusteeQuarterlyReport(request.BPID, MakeQuarterFormalDate(request.Quarter, request.Year), client.ID > 0),
		}}))
}

//HandleClientMonthlyContribution Client Monthly Contributions
func HandleClientMonthlyContribution(c buffalo.Context) error {
	request := ClientMonthlyContributionRequest{}
	if err := c.Bind(&request); err != nil {
		return c.Render(200, r.JSON(map[string]interface{}{"error": true, "message": err.Error()}))
	}
	return c.Render(201, r.JSON(map[string]interface{}{"error": false, "message": "Contribution Added", "data": map[string]interface{}{"contribution": AddClientMonthlyContribution(request)}}))
}

//HandleLoadFundManagers handles the request to load all fund managers
func HandleLoadFundManagers(c buffalo.Context) error {
	managers := LoadFundManagers()
	return c.Render(201, r.JSON(map[string]interface{}{"error": false, "message": "Fund Managers loaded", "data": map[string]interface{}{
		"managers": managers,
	}}))
}

//HandleUnidentifiedPaymentsUpload handles the upload of unidentified payments
func HandleUnidentifiedPaymentsUpload(c buffalo.Context) error {
	request := make([]UnidentifiedPaymentRequest, 0)
	if err := c.Bind(&request); err != nil {
		return c.Render(200, r.JSON(map[string]interface{}{"error": true, "message": err.Error()}))
	}
	StoreUnidentifiedPayments(request)
	return c.Render(201, r.JSON(map[string]interface{}{"error": false, "message": "Payments Added"}))
}

func HandleGOGMaturitiesUpload(c buffalo.Context) error {
	request := GOGMaturitiesUploadRequest{}
	if err := c.Bind(&request); err != nil {
		return c.Render(200, r.JSON(map[string]interface{}{"error": true, "message": err.Error()}))
	}
	StoreGOGMaturitiesInTheDB(request, AuthID(c))
	return c.Render(201, r.JSON(map[string]interface{}{"error": false, "message": "Maturities Uploaded"}))
}

func HandleUploadTransactionVolumes(c buffalo.Context) error {
	request := TransactionVolumesUploadRequest{}
	if err := c.Bind(&request); err != nil {
		return c.Render(200, r.JSON(map[string]interface{}{"error": true, "message": err.Error()}))
	}
	StoreTransactionVolumes(request, AuthID(c))
	return c.Render(201, r.JSON(map[string]interface{}{"error": false, "message": "Transaction volumes Uploaded"}))
}

func HandleTrusteeReportApproval(c buffalo.Context) error {
	request := &OnlyBPID{}
	if err := c.Bind(request); err != nil {
		return err
	}
	ApproveTrusteeReport(AuthID(c), request.BPID)
	return c.Redirect(302, "/trustee-dashboard")
}

func HandleTrusteeReportDisapproval(c buffalo.Context) error {
	request := &CheckerCommentRequest{}
	if err := c.Bind(request); err != nil {
		return err
	}
	DisapproveTrusteeReport(AuthID(c), request.Comment)
	return c.Redirect(302, "/trustee-dashboard")
}

func HandleMonthlyContributionDelete(c buffalo.Context) error {
	request := &OnlyID{}
	if err := c.Bind(request); err != nil {
		return err
	}
	DeleteMonthlyContribution(request.ID)
	return c.Render(201, r.JSON(map[string]interface{}{"error": false, "message": "Contribution Deleted"}))
}

func HandleMonthlyContributionEdit(c buffalo.Context) error {
	request := ClientMonthlyContributionEditRequest{}
	if err := c.Bind(&request); err != nil {
		return err
	}
	return c.Render(201, r.JSON(map[string]interface{}{"error": false, "message": "Contribution Edited", "data": map[string]interface{}{"contribution": EditClientMonthlyContribution(request)}}))
}

func ShowUploadedGOGMaturities(c buffalo.Context) error {
	c.Set("activity", LoadTrusteeActivityByHash(c.Param("hash")))
	c.Set("maturities", LoadGOGMaturitiesByHash(c.Param("hash")))
	return c.Render(200, r.HTML("gog-maturities.html"))
}

func ShowUploadedTransactionVolumes(c buffalo.Context) error {
	c.Set("activity", LoadTrusteeActivityByHash(c.Param("hash")))
	c.Set("transactions", LoadTransactionVolumesByHash(c.Param("hash")))
	return c.Render(200, r.HTML("transaction-volumes.html"))
}

func HandleGOGMaturitiesDelete(c buffalo.Context) error {
	request := OnlyHash{}
	if err := c.Bind(&request); err != nil {
		return err
	}
	DeleteGOGMaturities(request.Hash)
	return c.Redirect(302, "/trustee-dashboard")
}

func HandleTransactionVolumesDelete(c buffalo.Context) error {
	request := OnlyHash{}
	if err := c.Bind(&request); err != nil {
		return err
	}
	DeleteTransactionVolumes(request.Hash)
	return c.Redirect(302, "/trustee-dashboard")
}
