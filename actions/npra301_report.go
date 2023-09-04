package actions

import (
	"fmt"
	"time"
)

type PVSummarySumData struct {
	Total      float64 `gorm:"column:total" json:"total"`
	ClientName string  `gorm:"column:CLIENTNAME" json:"client_name"`
}

// Get Investment Receivables
func GetInvestmentReceivablesFromPVSummary(bpId, reportDate string) float64 {
	data := PVSummarySumData{}

	sql := fmt.Sprintf("select sum(LCY_AMT) as total from CRS_PVSUMMARY where BP_ID = '%s' and CLIENT_TYPE = 'npra' ", bpId)
	sql += "and SECURITY_TYPE like '%REC%' and REPORT_DATE = '" + reportDate + "'"

	fmt.Println("SQL: ", sql)
	DatabaseConnection.Raw(sql).Scan(&data)

	return data.Total
}

// Get Local Gov
func GetLocalGovernmentFromPVSummary(bpId, reportDate string) float64 {
	data := PVSummarySumData{}

	sql := fmt.Sprintf("select sum(LCY_AMT) as total from CRS_PVSUMMARY where BP_ID = '%s' and CLIENT_TYPE = 'npra' ", bpId)
	sql += "and SECURITY_TYPE like '%LOCAL %" + "GOV%' and REPORT_DATE = '" + reportDate + "'"

	fmt.Println("SQL: ", sql)

	DatabaseConnection.Raw(sql).Scan(&data)

	return data.Total
}

// Get Gov
func GetGovernmentFromPVSummary(bpId, reportDate string) float64 {
	data := PVSummarySumData{}

	sql := fmt.Sprintf("select sum(LCY_AMT) as total from CRS_PVSUMMARY where BP_ID = '%s' and CLIENT_TYPE = 'npra' ", bpId)
	sql += "and SECURITY_TYPE like 'GOV%' and REPORT_DATE = '" + reportDate + "'"

	fmt.Println("SQL: ", sql)

	DatabaseConnection.Raw(sql).Scan(&data)

	return data.Total
}

// Get COPORATE
func GetCoporateFromPVSummary(bpId, reportDate string) float64 {
	data := PVSummarySumData{}

	sql := fmt.Sprintf("select sum(LCY_AMT) as total from CRS_PVSUMMARY where BP_ID = '%s' and CLIENT_TYPE = 'npra' ", bpId)
	sql += "and SECURITY_TYPE like '%" + "CORP%' and REPORT_DATE = '" + reportDate + "'"

	fmt.Println("SQL: ", sql)
	DatabaseConnection.Raw(sql).Scan(&data)

	return data.Total
}

// Get Bank
func GetBankFromPVSummary(bpId, reportDate string) float64 {
	data := PVSummarySumData{}

	sql := fmt.Sprintf("select sum(LCY_AMT) as total from CRS_PVSUMMARY where BP_ID = '%s' and CLIENT_TYPE = 'npra' ", bpId)
	sql += "and SECURITY_TYPE like '%" + "BANK%' and REPORT_DATE = '" + reportDate + "'"

	fmt.Println("SQL: ", sql)
	DatabaseConnection.Raw(sql).Scan(&data)

	return data.Total
}

// Get Share
func GetShareFromPVSummary(bpId, reportDate string) float64 {
	data := PVSummarySumData{}

	sql := fmt.Sprintf("select sum(LCY_AMT) as total from CRS_PVSUMMARY where BP_ID = '%s' and CLIENT_TYPE = 'npra' ", bpId)
	sql += "and SECURITY_TYPE like '%" + "SHARE%' and REPORT_DATE = '" + reportDate + "'"

	fmt.Println("SQL: ", sql)
	DatabaseConnection.Raw(sql).Scan(&data)

	return data.Total
}

// Get COLLECTIVE
func GetCollectiveFromPVSummary(bpId, reportDate string) float64 {
	data := PVSummarySumData{}

	sql := fmt.Sprintf("select sum(LCY_AMT) as total from CRS_PVSUMMARY where BP_ID = '%s' and CLIENT_TYPE = 'npra' ", bpId)
	sql += "and SECURITY_TYPE like '%" + "COLLECTIVE%' and REPORT_DATE = '%" + reportDate + "%'"

	DatabaseConnection.Raw(sql).Scan(&data)

	return data.Total
}

// Get Bank Balance
func GetBankBalanceFromPVSummary(bpId, reportDate string) float64 {
	data := PVSummarySumData{}

	sql := fmt.Sprintf("select sum(LCY_AMT) as total from CRS_PVSUMMARY where BP_ID = '%s' and CLIENT_TYPE = 'npra' ", bpId)
	sql += "and SECURITY_TYPE like '%" + "BAL%' and REPORT_DATE = '" + reportDate + "'"

	DatabaseConnection.Raw(sql).Scan(&data)

	return data.Total
}

// Get ALTERNATE INV
func GetAltInvestmentFromPVSummary(bpId, reportDate string) float64 {
	data := PVSummarySumData{}

	sql := fmt.Sprintf("select sum(LCY_AMT) as total from CRS_PVSUMMARY where BP_ID = '%s' and CLIENT_TYPE = 'npra' ", bpId)
	sql += "and SECURITY_TYPE NOT like '%" + "REC%'"
	sql += "and SECURITY_TYPE NOT like '%LOCAL %" + "GOV%'"
	sql += "and SECURITY_TYPE NOT like '%" + "GOV%'"
	sql += "and SECURITY_TYPE NOT like '%" + "CORP%'"
	sql += "and SECURITY_TYPE NOT like '%" + "BANK%'"
	sql += "and SECURITY_TYPE NOT like '%" + "SHARE%'"
	sql += "and SECURITY_TYPE NOT like '%" + "BAL%' and REPORT_DATE = '" + reportDate + "'"

	DatabaseConnection.Raw(sql).Scan(&data)

	return data.Total
}

// GetNPRA
func GetNPRA0301ByBPIdAndReportDate(bpId, reportDate string) NPRA0301 {
	npra301 := NPRA0301{}
	data := PVSummarySumData{}
	date, _ := time.Parse("2006-01-02", Explode("T", reportDate)[0])
	fmt.Println("&&&&&&& Report Date", date)

	invReceivables := GetInvestmentReceivablesFromPVSummary(bpId, reportDate)
	localGov := GetLocalGovernmentFromPVSummary(bpId, reportDate)
	gov := GetGovernmentFromPVSummary(bpId, reportDate)
	corporate := GetCoporateFromPVSummary(bpId, reportDate)
	bank := GetBankFromPVSummary(bpId, reportDate)
	bankBal := GetBankBalanceFromPVSummary(bpId, reportDate)
	share := GetShareFromPVSummary(bpId, reportDate)
	collective := GetCollectiveFromPVSummary(bpId, reportDate)
	altInvestment := GetAltInvestmentFromPVSummary(bpId, reportDate)
	totalAssets := invReceivables + localGov + gov + corporate + bank + bankBal + share + collective + altInvestment
	//entity name
	DatabaseConnection.Raw("select CLIENTNAME from CRS_CLIENT where BP_ID = ?", bpId).Scan(&data)
	fmt.Println("############CLIENTNAME: ", data.ClientName)

	npra301.BP_ID = bpId
	fmt.Println("BPID: ", npra301.BP_ID)
	npra301.ReportCode = "INV MONTHLY"
	fmt.Println("ReportCode: ", npra301.ReportCode)
	npra301.EntityID = "Scheme"
	fmt.Println("EnityID: ", npra301.EntityID)
	npra301.EntityName = data.ClientName
	fmt.Println("EntityName: ", npra301.EntityName)
	npra301.ReferencePeriodYear = fmt.Sprintf("%d", date.Year())
	fmt.Println("Ref Period Year: ", npra301.ReferencePeriodYear)
	npra301.ReferencePeriod = fmt.Sprintf("%s", date.Month())
	fmt.Println("Ref Period: ", npra301.ReferencePeriod)
	npra301.InvestmentReceivables = invReceivables
	fmt.Println("InvestmentReceivables: ", npra301.InvestmentReceivables)
	npra301.TotalAssetUnderManagement = totalAssets
	fmt.Println("TotalAssetUnderManagement: ", npra301.TotalAssetUnderManagement)
	npra301.GovernmentSecurities = gov
	fmt.Println("GovernmentSecurities: ", npra301.GovernmentSecurities)
	npra301.LocalGovernmentSecurities = localGov
	fmt.Println("LocalGovernmentSecurities: ", npra301.LocalGovernmentSecurities)
	npra301.CorporateDebtSecurities = corporate
	fmt.Println("CorporateDebtSecurities: ", npra301.CorporateDebtSecurities)
	npra301.BankSecurities = bank
	fmt.Println("BankSecurities: ", npra301.BankSecurities)
	npra301.OrdinaryPreferenceShares = int(share)
	fmt.Println("OrdinaryPreferenceShares: ", npra301.OrdinaryPreferenceShares)
	npra301.CollectiveInvestmentScheme = int(collective)
	fmt.Println("CollectiveInvestmentScheme: ", npra301.CollectiveInvestmentScheme)
	npra301.AlternativeInvestments = int(altInvestment)
	fmt.Println("AlternativeInvestments: ", npra301.AlternativeInvestments)
	npra301.BankBalances = bankBal
	fmt.Println("BankBalances: ", npra301.BankBalances)
	npra301.ReportingDate = date
	fmt.Println("ReportingDate: ", npra301.ReportingDate)
	npra301.CreatedAt = time.Now()
	fmt.Println("CreatedAt: ", npra301.CreatedAt)

	return npra301
}

func InsertNPRA301DataToDB(data NPRA0301) {
	fmt.Println("$$$$$$ NRPA0301: ", data)
	DatabaseConnection.Create(&data)
}
