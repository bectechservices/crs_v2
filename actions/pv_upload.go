package actions

import (
	"sync"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
)

// UploadPVReport handles the pv upload to the server
func UploadPVReport(c buffalo.Context) error {
	request := &PVUploadRequest{}
	if err := c.Bind(request); err != nil {
		return err
	}
	parsedPV, parseErrors := parseExcelData(request.Data, AuthID(c), request.PVType, request.CashBalance)
	shouldDebug := envy.Get("DEBUG_UPLOADED_PV", "0")

	var pvUploadWaitGroup sync.WaitGroup
	errorChannel := make(chan PVUploadError)

	for _, data := range parsedPV {
		if shouldDebug == "1" {
			dumpPVReportToConsole(data)
		}
		pvUploadWaitGroup.Add(1)
		go UploadPVReportData(data, &pvUploadWaitGroup, errorChannel)
	}
	go func(wg *sync.WaitGroup, errorChannel chan PVUploadError) {
		wg.Wait()
		close(errorChannel)
	}(&pvUploadWaitGroup, errorChannel)

	for err := range errorChannel {
		parseErrors = append(parseErrors, err)
	}

	allClients := GetClientDetails(parsedPV[0].ClientID)

	//logic for NPRA0301 Report Upload
	if !NPRA0301DataExists(allClients[0].BPID, parsedPV[0].Date) {
		npra301Data := GetNPRA0301ByBPIdAndReportDate(allClients[0].BPID, parsedPV[0].Date)
		InsertNPRA301DataToDB(npra301Data)
	}

	//logic for NPRA0302 Report Upload
	if !NPRA0302DataExists(allClients[0].BPID, parsedPV[0].Date) {
		InsertIntoNPRA0302Report(allClients[0].BPID, parsedPV[0].Date)
	}

	if len(parseErrors) > 0 {
		return c.Render(200, r.JSON(map[string]interface{}{"error": true, "messages": parseErrors}))
	}
	return c.Render(200, r.JSON(map[string]interface{}{"error": false, "message": "File uploaded"}))
}

func ShowUploadedPVErrors(c buffalo.Context) error {
	return c.Render(200, r.HTML("pv-errors.html"))
}
