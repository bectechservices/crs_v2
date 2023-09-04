package actions

import "github.com/gobuffalo/buffalo"

//ShowReportsPage displays the report page
func ShowReportsPage(c buffalo.Context) error {
	return c.Render(200, r.HTML("reports.html"))
}

//ShowTrusteeReportPage displays the Trustee report page
func ShowTrusteeReportPage(c buffalo.Context) error {
	date := MakeDBQueryableCurrentQuarterDate()
	c.Set("t2_clients", LoadTier2Clients(date))
	c.Set("t3_clients", LoadTier3Clients(date))
	return c.Render(200, r.HTML("trustee-report.html"))
}

// TODO:NPRA REPORT PAGE

//ShowNPRAReportPage displays the NPRA report page
func ShowNPRAReportPage(c buffalo.Context) error {
	return c.Render(200, r.HTML("npra-report.html"))
}
