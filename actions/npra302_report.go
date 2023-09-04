package actions

import (
	"fmt"
	"strings"
	"time"
)

func GetInvestmentID(request NPRAReportRequest) {}

func GetClientNameById(clientId string) string {
	data := ClientPVUniqueCode{}
	DatabaseConnection.Raw("select CLIENT_NAME from CRS_CLIENT_PV_UNIQUE_CODES where CODE = ?", clientId).First(&data)
	return data.ClientName
}

// func ExtractMaturityDate(value string) string {
// 	value = strings.TrimSpace(value)
// 	expValue := Explode(" ", value)

// 	//fmt.Println("@@@@@@@@@@@@@@", expValue)
// 	sort.Sort(sort.Reverse(sort.StringSlice(expValue)))

// 	for i := 0; i < len(expValue); i++ {
// 		//fmt.Println("!!!!!!!!!!!!!", expValue[i])
// 		if strings.ContainsAny(expValue[i], "/") {
// 			maturityDate := Explode("-", expValue[i])
// 			return maturityDate[len(maturityDate)-1]
// 		}
// 	}

// 	return ""
// }

func Explode(delimiter, text string) []string {
	if len(delimiter) > len(text) {
		return strings.Split(delimiter, text)
	} else {
		return strings.Split(text, delimiter)
	}
}

func GetIssuerName(securityName, secTitle string) string {
	expValue := Explode("-", securityName)

	if expValue[0] == "REPUBLIC OF GHANA" || strings.Contains(secTitle, "GOVERNMENT") ||
		strings.Contains(secTitle, "TREASURY") {
		return "GOG"
	}
	return strings.TrimSpace(strings.ToUpper(expValue[0]))

}

func GetAssetClass(securityType string) string {
	record := AssetClass{}

	sql := "select * from CRS_ASSET_CLASS where SECURITY_TYPE LIKE '%" + securityType + "%'"
	DatabaseConnection.Raw(sql).First(&record)

	return record.AssetClassName

}

func GetPVReportsByBPIdAndReportDate(bpId string, reportDate string) []PVReportField {
	reports := []PVReportField{}
	sql := "select * from CRS_PVREPORTS where BP_ID = '" + bpId + "' and REPORT_DATE like '%" + reportDate + "%'"
	fmt.Println(sql)
	DatabaseConnection.Raw(sql).Scan(&reports)

	return reports
}

func GetPVReportsByReportDate(bpId, reportDate string) []PVReportField {
	reports := []PVReportField{}
	sql := "select * from CRS_PVREPORTS where BP_ID = '" + bpId + "' and REPORT_DATE like '%" + reportDate + "%'"
	fmt.Println(sql)
	DatabaseConnection.Raw(sql).Scan(&reports)
	fmt.Println("@@@@@@@@@@@@@@PVT Rep Count 1: ", len(reports))

	return reports
}

func GetClientCode(cdsCode, isin string) string {
	if len(cdsCode) > len(isin) {
		return cdsCode
	}
	return isin
}

func InsertIntoNPRA0302Report(bpId, reportDate string) {
	fmt.Println("@@@@@@@@@@@@@@@@@2Report Date: ", reportDate)
	//get pvc reports
	pvcReports := GetPVReportsByReportDate(bpId, reportDate)
	fmt.Println("@@@@@@@@@@@@@@PVT Rep Count 2: ", len(pvcReports))
	//date,_:=time.Parse("2006-01-02",reportDate)

	for _, report := range pvcReports {

		maturityDate, _ := time.Parse("2006-01-02", Explode("T", report.DateTo)[0])
		//fmt.Println("@@@@@@@@ Actual Maturity Date: ", report.DateTo)

		_npra302 := NPRA0302{
			BP_ID:               report.BPID,
			ClientCode:          GetClientCode(report.CDSCode, report.ISIN),
			ReportCode:          "INV_MONTHLY",
			EntityID:            "Scheme",
			EntityName:          GetClientNameById(report.ClientID),
			ReferencePeriodYear: Explode("-", reportDate)[0],
			ReferencePeriod:     Explode("-", reportDate)[1],
			InvestmentID:        strings.TrimSpace(report.CDSCode),
			Instrument:          report.SecurityType + " -%",
			AssetTenure:         report.SecurityName,
			ReportingDate:       report.ReportDate,
			FaceValue:           report.NominalValue,
			Currency:            "GHS",
			MarketValue:         report.Value,
			MaturityDate:        maturityDate,
			IssuerName:          GetIssuerName(report.SecurityName, report.SecurityType),
		}

		DatabaseConnection.Create(&_npra302)
	}
}
