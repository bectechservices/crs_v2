package actions

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
)

// ClientsPage loads the clients page
func ClientsPage(c buffalo.Context) error {
	c.Set("clients", LoadAllClients())
	return c.Render(200, r.HTML("clients.html"))
}

// ClientPVReportPage returns the client pv report page
func ClientPVReportPage(c buffalo.Context) error {
	pvData := make([]map[string]interface{}, 0)
	reportApproved := true
	bpids := strings.Split(c.Param("bpid"), ",")
	for _, bpid := range bpids {
		data := map[string]interface{}{}
		client := LoadClientDetails(bpid, MakeQuarterFormalDate(c.Param("quarter"), c.Param("year")))
		data["client"] = ClientDetailsArrayToString(client)
		data["auc_for_past_months"] = LoadClientAUCForPastMonths(bpid, c.Param("quarter"), c.Param("year"))
		data["txn_vols_for_past_months"] = LoadClientTransactionVolumesForPastMonths(bpid, c.Param("quarter"), c.Param("year"))
		txnByAssetClass, totalNumberOfTransactionsVols := LoadClientTransactionVolumesByAssetClassForPastMonths(bpid, c.Param("quarter"), c.Param("year"))
		data["txn_vols_by_asset_class_for_past_months"] = txnByAssetClass
		data["total_txn_vols_by_asset_class"] = totalNumberOfTransactionsVols
		data["txn_by_asset_class_summary"] = MakeTransactionVolumeByAssetClassSummary(txnByAssetClass)
		contributions := LoadClientMonthlyContributions(bpid, c.Param("quarter"), c.Param("year"))
		data["monthly_contributions"] = SortContributionsByMonths(contributions)
		var monthlyContributionsTotal float64
		for _, contribution := range contributions {
			monthlyContributionsTotal += contribution.Contributions
		}
		data["monthly_contributions_total"] = monthlyContributionsTotal
		data["individual_contributions"] = LoadMonthlyContributionsForEachSCA(bpid, c.Param("quarter"), c.Param("year"))
		data["individual_pvs"] = LoadPVSummaryForEachSCA(bpid, c.Param("quarter"), c.Param("year"))
		data["gog_maturities"] = SortGOGMaturitiesByMonths(LoadClientGOGMaturitiesSummary(bpid, c.Param("quarter"), c.Param("year")))
		fdMaturities := LoadClientFDMaturitiesSummary(bpid, c.Param("quarter"), c.Param("year"))
		data["fd_maturities"] = CountFDs(fdMaturities)
		data["fd_maturities_count"] = len(fdMaturities)
		data["corporate_actions"] = LoadClientCorporateActionActivities(bpid, c.Param("quarter"), c.Param("year"))
		data["lumped_pv_summary_for_past_3_quarters"] = LoadLumpedPVSummaryForPast3Quarters(bpid, c.Param("quarter"), c.Param("year"))
		data["unidentified_payments"] = LoadClientUnidentifiedPayments(bpid, c.Param("quarter"), c.Param("year"))
		report := LoadClientTrusteeQuarterlyReport(bpid, MakeQuarterFormalDate(c.Param("quarter"), c.Param("year")), false)
		reportApproved = reportApproved && report.Approved
		pvData = append(pvData, data)
	}
	nameOnReport := ""
	if len(pvData) == 2 {
		nameOnReport = (pvData[0]["client"].(map[string]interface{}))["name"].(string) + " & " + (pvData[1]["client"].(map[string]interface{}))["name"].(string)
	} else {
		nameOnReport = (pvData[0]["client"].(map[string]interface{}))["name"].(string)
	}
	//transactionVolsHeading := ""
	//if c.Param("quarter") == "1" {
	//	yearInt, _ := strconv.Atoi(c.Param("year"))
	//	transactionVolsHeading = fmt.Sprintf("Transactions Volumes for %d and %s", yearInt-1, c.Param("year"))
	//} else {
	transactionVolsHeading := "Transactions Volumes as at Q" + c.Param("quarter") + " " + c.Param("year")
	//}
	c.Set("report_data", pvData)
	c.Set("txn_vols_heading", transactionVolsHeading)
	c.Set("report_has_been_approved", reportApproved)
	c.Set("report_name", nameOnReport)
	c.Set("current_quarter_month_year", MakeQuarterMonthAndYear(c.Param("quarter"), c.Param("year")))
	c.Set("current_quarter", "Q"+c.Param("quarter")+" "+c.Param("year"))
	c.Set("current_quarter_number", c.Param("quarter"))
	return c.Render(200, r.HTML("pv-report.html"))
}

// OffshoreClients loads all offshore clients
func OffshoreClients(c buffalo.Context) error {
	request := OffshoreClientsRequest{}
	if err := c.Bind(&request); err != nil {
		return c.Render(200, r.JSON(map[string]interface{}{"error": true, "message": err.Error()}))
	}
	clients := LoadOffshoreClients(request)
	return c.Render(200, r.JSON(map[string]interface{}{"error": false, "message": "Offshore clients loaded", "data": map[string]interface{}{"clients": clients}}))
}

func ShowClientAdd(c buffalo.Context) error {
	c.Set("fundmanagers", LoadFundManagers())
	return c.Render(200, r.HTML("client-add.html"))
}

func HandleAddClient(c buffalo.Context) error {
	request := &ClientAddRequest{}
	if err := c.Bind(request); err != nil {
		return err
	}

	if file, err := c.File("image"); err == nil {
		uploadPath := NormalizeWindowsPath(envy.Get("FILE_MANAGER_DIR", "C:\\CRS\\uploads"))
		if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
			if err := os.Mkdir(uploadPath, 0755); err != nil {
				return err
			}
		}
		filename := string(RandomBytes(64)) + filepath.Ext(file.Filename)
		if _, err := os.Stat(filepath.Join(uploadPath, filename)); err == nil {
			return err
		}
		uploadedFile, err := os.Create(filepath.Join(uploadPath, filename))
		if err != nil {
			return err
		} else {

		}

		defer uploadedFile.Close()

		if _, err = io.Copy(uploadedFile, file); err != nil {
			return err
		} else {
			request.Image = filename
		}
	}
	// AddClientToDB(*request)
	if err := AddClientToDB(*request); err == true {
		c.Flash().Add("success", "Client created successfully")
		return c.Redirect(302, "/clients")
	}
	c.Flash().Add("Failure", "Client creation failed")
	return c.Redirect(302, "/clients")
}

// ShowClientDetails displays the Client details page
func ShowClientDetails(c buffalo.Context) error {
	c.Set("client", GetClientByBPID(c.Param("bpid")))
	c.Set("billing_info", LoadClientBillingInfo(c.Param("bpid")))
	c.Set("basis_points", LoadClientBasisPoints(c.Param("bpid")))
	return c.Render(200, r.HTML("client-details.html"))
}

// ShowClientEdit loads the client edit page
func ShowClientEdit(c buffalo.Context) error {
	c.Set("client", GetClientByBPID(c.Param("bpid")))
	c.Set("fundmanagers", LoadFundManagers())
	c.Set("billing_info", LoadClientBillingInfo(c.Param("bpid")))
	c.Set("basis_points", LoadClientBasisPoints(c.Param("bpid")))
	return c.Render(200, r.HTML("client-edit.html"))
}

func HandleEditClient(c buffalo.Context) error {
	request := &ClientEditRequest{}
	if err := c.Bind(request); err != nil {
		return err
	}

	if file, err := c.File("image"); err == nil {
		uploadPath := NormalizeWindowsPath(envy.Get("FILE_MANAGER_DIR", "C:\\CRS\\uploads"))
		if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
			if err := os.Mkdir(uploadPath, 0755); err != nil {
				return err
			}
		}
		filename := string(RandomBytes(64)) + filepath.Ext(file.Filename)
		if _, err := os.Stat(filepath.Join(uploadPath, filename)); err == nil {
			return err
		}

		uploadedFile, err := os.Create(filepath.Join(uploadPath, filename))
		if err != nil {
			return err
		} else {

		}

		defer uploadedFile.Close()

		if _, err = io.Copy(uploadedFile, file); err != nil {
			return err
		} else {
			request.Image = filename
		}
	}
	EditClient(*request)
	c.Flash().Add("success", "Client edited successfully")
	return c.Redirect(302, "/client-edit") // ?bpid="+request.BPID)
}

func HandleLoadClientDetails(c buffalo.Context) error {
	request := &OnlyBPID{}
	if err := c.Bind(request); err != nil {
		return c.Render(200, r.JSON(map[string]interface{}{"error": true, "message": err.Error()}))
	}
	return c.Render(200, r.JSON(map[string]interface{}{"error": false, "message": "Client loaded", "data": map[string]interface{}{"client": GetClientByBPID(request.BPID)}}))
}

func HandleLoadClientDetailsWithCode(c buffalo.Context) error {
	request := &OnlyCode{}
	if err := c.Bind(request); err != nil {
		return c.Render(200, r.JSON(map[string]interface{}{"error": true, "message": err.Error()}))
	}
	return c.Render(200, r.JSON(map[string]interface{}{"error": false, "message": "Client loaded", "data": map[string]interface{}{"client": GetClientByCode(request.Code)}}))
}

func HandleUpdateClientBillingInfo(c buffalo.Context) error {
	request := ClientBillingInfoUpdateRequest{}
	if err := c.Bind(&request); err != nil {
		return c.Render(200, r.JSON(map[string]interface{}{"error": true, "message": err.Error()}))
	}
	UpdateClientBillingInfo(request)
	return c.Render(200, r.JSON(map[string]interface{}{"error": false, "message": "Data uploaded"}))
}

func ShowClientSCAsDetails(c buffalo.Context) error {
	c.Set("clients", GetClientSCASByBPID(c.Param("bpid")))
	c.Set("bpid", c.Param("bpid"))
	return c.Render(200, r.HTML("client-scas.html"))
}

func HandleAddClientSCA(c buffalo.Context) error {
	request := SCAManagementRequest{}
	if err := c.Bind(&request); err != nil {
		return err
	}
	AddClientSCA(request)
	c.Flash().Add("success", "Client SCA created successfully")
	return c.Redirect(302, "/client-scas?bpid="+request.BPID)
}

func HandleEditClientSCA(c buffalo.Context) error {
	request := SCAManagementRequest{}
	if err := c.Bind(&request); err != nil {
		return err
	}
	EditClientSCA(request)
	c.Flash().Add("success", "Client SCA edited successfully")
	return c.Redirect(302, "/client-scas?bpid="+request.BPID)
}

func HandleDeleteClientSCA(c buffalo.Context) error {
	request := SCAManagementRequest{}
	if err := c.Bind(&request); err != nil {
		return err
	}
	DeleteClientSCA(request.ID)
	c.Flash().Add("success", "Client SCA deleted successfully")
	return c.Redirect(302, "/client-scas?bpid="+request.BPID)
}

func HandleAccountClose(c buffalo.Context) error {
	request := OnlyBPID{}
	if err := c.Bind(&request); err != nil {
		return err
	}
	CloseClientAccount(request.BPID)
	c.Flash().Add("success", "Client account closed successfully")
	return c.Redirect(302, "/clients")
}

func HandleAccountOpen(c buffalo.Context) error {
	request := OnlyBPID{}
	if err := c.Bind(&request); err != nil {
		return err
	}
	OpenClientAccount(request.BPID)
	c.Flash().Add("success", "Client account opened successfully")
	return c.Redirect(302, "/clients")
}

func HandleSCAClose(c buffalo.Context) error {
	request := SCAAccountMgmt{}
	if err := c.Bind(&request); err != nil {
		return err
	}
	CloseSCA(request.Code)
	c.Flash().Add("success", "SCA closed successfully")
	return c.Redirect(302, "/client-scas?bpid="+request.BPID)
}

func HandleSCAOpen(c buffalo.Context) error {
	request := SCAAccountMgmt{}
	if err := c.Bind(&request); err != nil {
		return err
	}
	OpenSCA(request.Code)
	c.Flash().Add("success", "SCA opened successfully")
	return c.Redirect(302, "/client-scas?bpid="+request.BPID)
}

func HandleLoadClientEmails(c buffalo.Context) error {
	c.Set("emails", GetClientEmailsByBPID(c.Param("bpid")))
	c.Set("bpid", c.Param("bpid"))
	return c.Render(200, r.HTML("client-emails.html"))
}

func HandleAddClientEmail(c buffalo.Context) error {
	request := ClientEmailMgmt{}
	if err := c.Bind(&request); err != nil {
		return err
	}
	AddClientEmail(request.BPID, request.Email)
	c.Flash().Add("success", "Email added successfully")
	return c.Redirect(302, "/client-emails?bpid="+request.BPID)
}

func HandleEditClientEmail(c buffalo.Context) error {
	request := ClientEmailMgmt{}
	if err := c.Bind(&request); err != nil {
		return err
	}
	EditClientEmail(request.ID, request.Email)
	c.Flash().Add("success", "Email edited successfully")
	return c.Redirect(302, "/client-emails?bpid="+request.BPID)
}

func HandleDeleteClientEmail(c buffalo.Context) error {
	request := ClientEmailMgmt{}
	if err := c.Bind(&request); err != nil {
		return err
	}
	DeleteClientEmail(request.ID)
	c.Flash().Add("success", "Email deleted successfully")
	return c.Redirect(302, "/client-emails?bpid="+request.BPID)
}

func HandleLoadMailService(c buffalo.Context) error {
	c.Set("emails", LoadMailService())
	c.Set("senderInfo", LoadMailSenderInfo())
	c.Set("holidays", LoadHolidays())
	return c.Render(200, r.HTML("mail-service.html"))
}

func HandleAddMailService(c buffalo.Context) error {
	request := MailServiceMgmt{}
	if err := c.Bind(&request); err != nil {
		return err
	}
	AddMailServiceEmail(request.Email)
	c.Flash().Add("success", "Email added successfully")
	return c.Redirect(302, "/mail-service")
}

func HandleEditMailService(c buffalo.Context) error {
	request := MailServiceMgmt{}
	if err := c.Bind(&request); err != nil {
		return err
	}
	EditMailServiceEmail(request.ID, request.Email)
	c.Flash().Add("success", "Email edited successfully")
	return c.Redirect(302, "/mail-service")
}

func HandleDeleteMailService(c buffalo.Context) error {
	request := MailServiceMgmt{}
	if err := c.Bind(&request); err != nil {
		return err
	}
	DeleteMailServiceEmail(request.ID)
	c.Flash().Add("success", "Email deleted successfully")
	return c.Redirect(302, "/mail-service")
}

func HandleLoadClientsReadOnly(c buffalo.Context) error {
	c.Set("clients", LoadAllClients())
	return c.Render(200, r.HTML("clients-readonly.html"))
}

func HandleClientSCAsReadOnly(c buffalo.Context) error {
	c.Set("clients", GetClientSCASByBPID(c.Param("bpid")))
	c.Set("bpid", c.Param("bpid"))
	return c.Render(200, r.HTML("client-scas-readonly.html"))
}

func HandleSendClientLetters(c buffalo.Context) error {
	if !ClientLettersHaveBeenSent() {
		SendSecLetterToClients(AuthUser(c).Fullname)
	}
	return c.Redirect(302, "/client-letters")
}

func HandleUpdateMailSenderInfo(c buffalo.Context) error {
	request := MailSenderInfoUpdateRequest{}
	if err := c.Bind(&request); err != nil {
		return err
	}
	UpdateMailSenderInfo(request)
	c.Flash().Add("success", "Sender Info edited successfully")
	return c.Redirect(302, "/mail-service")
}

func HandleLoadMailSenderInfo(c buffalo.Context) error {
	return c.Render(200, r.JSON(map[string]interface{}{"error": false, "message": "Mail sender info loaded", "data": map[string]interface{}{"sender": LoadMailSenderInfo()}}))
}

func HandleAddHoliday(c buffalo.Context) error {
	request := HolidayMgmtRequest{}
	if err := c.Bind(&request); err != nil {
		return err
	}
	AddHoliday(request)
	c.Flash().Add("success", "Holiday added successfully")
	return c.Redirect(302, "/mail-service")
}

func HandleDeleteHoliday(c buffalo.Context) error {
	request := HolidayMgmtRequest{}
	if err := c.Bind(&request); err != nil {
		return err
	}
	DeleteHoliday(request.ID)
	c.Flash().Add("success", "Holiday deleted successfully")
	return c.Redirect(302, "/mail-service")
}
