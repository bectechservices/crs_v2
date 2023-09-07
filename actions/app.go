package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	csrf "github.com/gobuffalo/mw-csrf"
	paramlogger "github.com/gobuffalo/mw-paramlogger"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
//var ENV = envy.Get("GO_ENV", "production")

var ENV = "production"
var app *buffalo.App

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
//
// Routing, middleware, groups, etc... are declared TOP -> DOWN.
// This means if you add a middleware to `app` *after* declaring a
// group, that group will NOT have that new middleware. The same
// is true of resource declarations as well.
//
// It also means that routes are checked in the order they are declared.
// `ServeFiles` is a CATCH-ALL route, so it should always be
// placed last in the route declarations, as it will prevent routes
// declared after it to never be called.
//------------------------------------------------------------

// type myConfig struct {
//     value string
// }

// func (config *myConfig) myMiddlewareFunc() buffalo.MiddlewareFunc {
//     return func(next buffalo.Handler) buffalo.Handler {
//         return func(c buffalo.Context) error {
//             c.Logger().Info("Test ", config.value)

//             return next(c)
//         }
//     }
// }

// -------------------------------------------------------------
func App() *buffalo.App {
	if app == nil {
		store := sessions.NewCookieStore(securecookie.GenerateRandomKey(64))
		store.MaxAge(12 * 60 * 60)
		app = buffalo.New(buffalo.Options{
			Env:          ENV,
			SessionName:  "_crs_session",
			SessionStore: store,
		})

		// Log request parameters (filters apply).
		app.Use(paramlogger.ParameterLogger)

		app.Use(csrf.New)

		app.Use(SetConstantsInContext)

		app.Use(SetCurrentUser)

		app.Use(RequiresAuthorization)

		app.Use(RedirectIfAuthenticated)

		app.Use(MustResetPassword)

		app.Use(MustBeAdmin)

		app.Use(MustBeSec)

		app.Use(MustBeNPRA)

		app.Use(MustBeTrustee)

		app.Use(MustBeNAV)

		app.Use(LoadSystemAlerts)

		app.Middleware.Skip(RequiresAuthorization, HomeHandler, LoginHandler, HandlePasswordReset, Logout, ShowPasswordManager, SavePassword)

		app.Middleware.Skip(MustResetPassword, HomeHandler, LoginHandler, HandlePasswordReset, showProfilePage, HandleUserProfileUpdate, Logout, ShowPasswordManager, SavePassword)

		app.ErrorHandlers[404] = GlobalErrorHandler
		app.ErrorHandlers[500] = GlobalErrorHandler
		app.ErrorHandlers[403] = GlobalErrorHandler

		app.GET("/", HomeHandler)

		app.POST("/login", LoginHandler)

		app.GET("/dashboard", Dashboard)

		app.POST("/portal", UploadPVReport)

		app.GET("/clients", ClientsPage)

		app.GET("/client-pvreport", ClientPVReportPage)

		app.GET("/users", LoadUsersPage)

		app.GET("/quarter-overview", DisplayQuarterOverviewPage)

		app.GET("/custodian", ShowCustodianPage)

		app.GET("/sec", ShowSecDashboard)

		app.GET("/scheme-pv", ShowPVPage)

		app.POST("/upload-governance", HandleGovernanceDataUpload)

		app.POST("/upload-scheme-details", HandleSchemeDetailsDataUpload)

		app.POST("/delete-scheme-details", HandleSchemeDetailsDelete)

		app.POST("/upload-other-information", HandleOtherInformationUpload)

		app.POST("/upload-official-remarks", HandleOfficialReportRemarksUpload)

		app.POST("/sec-quarterly-report", LoadSecQuarterlyReport)

		app.POST("/fetch-scheme-details", HandleFetchSchemeDetails)

		app.POST("/recalculate-scheme-details", HandleRecalculateSchemeDetails)

		app.POST("/users", HandleAddUser)

		app.PATCH("/users", HandleEditUser)

		app.DELETE("/user", HandleUserDelete)

		app.POST("/trigger-password-reset", HandleTriggerPasswordReset)

		app.GET("/secreport", ShowSECReportsPage)

		app.GET("/trustee-report", ShowTrusteeReportPage)

		app.POST("/offshore-clients", OffshoreClients)

		app.POST("/load-sec-local-variance", SecLocalVariance)

		app.POST("/load-sec-foreign-variance", SecForeignVariance)

		app.POST("/load-npra-local-variance", NpraLocalVariance)

		app.POST("/load-npra-foreign-variance", NpraForeignVariance)

		app.GET("/trustee-dashboard", ShowTrusteeDashboardPage)

		app.POST("/trustee-data", LoadTrusteeData)

		app.POST("/client-monthly-contribution", HandleClientMonthlyContribution)

		app.GET("/client-view", ShowClientDetails)

		app.GET("/client-scas", ShowClientSCAsDetails)

		app.POST("/client-scas", HandleAddClientSCA)

		app.PATCH("/client-scas", HandleEditClientSCA)

		app.DELETE("/client-scas", HandleDeleteClientSCA)

		app.GET("/client-add", ShowClientAdd)

		app.POST("/client-add", HandleAddClient)

		app.GET("/client-edit", ShowClientEdit)

		app.PATCH("/client-edit", HandleEditClient)

		app.GET("/fundmanager", ShowFundManager)

		app.GET("/fundmanager-add", ShowAddFundManager)

		app.POST("/fundmanager-add", HandleAddFundManager)

		app.GET("/fundmanager-edit", ShowFundManagerEdit)

		app.PATCH("/fundmanager-edit", HandleEditFundManager)

		app.POST("/load-fund-managers", HandleLoadFundManagers)

		app.POST("/upload-unidentified-payments", HandleUnidentifiedPaymentsUpload)

		app.POST("/export-sec-report", ExportSecReportToWord)

		app.POST("/sec-foreign-variance-excel", ExportSecForeignVarianceToExcel)

		app.POST("/matured-securities-excel", ExportMaturedSecuritiesToExcel)

		app.POST("/sec-local-variance-excel", ExportSecLocalVarianceToExcel)

		app.POST("/directors-word", HandleExportDirectorsToWord)

		app.POST("/offshore-clients-excel", ExportOffshoreClientsToExcel)

		app.POST("/export-cover-letter", ExportCoverLetter)

		app.POST("/export-billing-report", HandleBillingReportExport)

		app.POST("/upload-unauthorized-transactions", HandleUnauthorizedTransactionUpload)

		app.POST("/upload-outstanding-fd-certificates", HandleOutstandingFDCertificateUpload)

		app.POST("/load-unauthorized-transactions", HandleLoadUnauthorizedTransactions)

		app.POST("/load-outstanding-fd-certificates", HandleLoadOutstandingFDCertificates)

		app.POST("/load-npra-quarterly-report", HandleLoadNPRAQuarterlyReport)

		app.GET("/billing-dashboard", ShowNavBillDashboardPage)

		app.POST("/load-client-details", HandleLoadClientDetails)

		app.POST("/load-client-details-with-code", HandleLoadClientDetailsWithCode)

		app.POST("/upload-billing-transaction-details", HandleTransactionDetailsUpload)

		app.POST("/upload-billing-currency-details", HandleCurrencyDetailsUpload)

		app.POST("/update-client-billing-info", HandleUpdateClientBillingInfo)

		app.POST("/approve-sec-current-quarter-pv", HandleSecPVApproval)

		app.POST("/approve-npra-current-report", HandleNPRAReportApproval)

		app.POST("/disapprove-sec-current-quarter-report", HandleSecCurrentQuarterReportDisapproval)

		app.POST("/disapprove-npra-current-quarter-report", HandleNPRACurrentQuarterReportDisapproval)

		app.POST("/sec-send-progress-to-checker", HandleSecSendProgressToChecker)

		app.POST("/npra-send-progress-to-checker", HandleNPRASendProgressToChecker)

		app.POST("/upload-variance-remarks", HandleUploadVarianceRemarks)

		app.POST("/upload-npra-variance-remarks", HandleUploadNPRAVarianceRemarks)

		app.POST("/load-current-npra-declaration", HandleLoadCurrentNPRADeclaration)

		app.POST("/upload-npra-declaration", HandleUploadNPRADeclaration)

		app.POST("/update-outstanding-fds", HandleOutstandingFDCertificateEdit)

		app.POST("/load-client-position", HandleLoadClientPosition)

		app.POST("/load-client-nav", HandleLoadClientNAV)

		app.POST("/update-clients-nav", HandleUpdateClientsNAV)

		app.POST("/upload-gog-maturities", HandleGOGMaturitiesUpload)

		app.DELETE("/delete-sec-uploaded-pv", HandleDeleteSECUploadedPV)

		app.DELETE("/delete-billing-uploaded-pv", HandleDeleteBillingUploadedPV)

		app.DELETE("/delete-npra-uploaded-pv", HandleDeleteNPRAUploadedPV)

		app.POST("/upload-txn-volumes", HandleUploadTransactionVolumes)

		app.GET("/outstanding-fd", ShowUploadedOutstandingFDCertificates)

		app.DELETE("/outstanding-fd", HandleOutstandingFDCertificatesDelete)

		app.POST("/approve-trustee-report", HandleTrusteeReportApproval)

		app.POST("/disapprove-trustee-report", HandleTrusteeReportDisapproval)

		app.DELETE("/trustee-monthly-contribution", HandleMonthlyContributionDelete)

		app.PATCH("/trustee-monthly-contribution", HandleMonthlyContributionEdit)

		app.GET("/gog-maturities", ShowUploadedGOGMaturities)

		app.DELETE("/gog-maturities", HandleGOGMaturitiesDelete)

		app.GET("/transaction-volumes", ShowUploadedTransactionVolumes)

		app.DELETE("/transaction-volumes", HandleTransactionVolumesDelete)

		app.POST("/disapprove-nav-current-month-report", HandleNAVMonthlyReportDisapproval)

		app.POST("/approve-nav-current-month-report", HandleNAVMonthlyReportApproval)

		app.POST("/reverse-nav-current-month-report-approval", HandleNAVMonthlyReportReverseApproval)

		app.POST("/npra-local-variance-excel", ExportNPRALocalVarianceToExcel)

		app.POST("/npra-outstanding-fd-excel", ExportNPRAOutstandingFDToExcel)

		app.POST("/npra-monthly-report-excel", ExportNPRAMonthlyReportToExcel)

		app.POST("/npra-unauthorized-report-word", ExportNPRAUnauthorizedReportToWord)

		app.POST("/npra-monthly-report-word", ExportNPRAMonthlyReportToWord)

		app.GET("/pv-upload-errors", ShowUploadedPVErrors)

		app.GET("/pv-history", LoadPVHistoryPage)

		app.GET("/directors", HandleShowDirectorsPage)

		app.POST("/directors", HandleAddDirector)

		app.PATCH("/directors", HandleEditDirector)

		app.DELETE("/directors", HandleDirectorDelete)

		app.GET("/matured-securities", HandleLoadMaturedSecurities)

		app.POST("/matured-securities", HandleMaturedSecuritiesUpload)

		app.POST("/load-matured-securities", HandleLoadMaturedSecuritiesJSON)

		app.DELETE("/matured-securities", HandleMaturedSecuritiesDelete)

		app.POST("/reset-password", HandlePasswordReset)

		app.POST("/close-account", HandleAccountClose)

		app.POST("/open-account", HandleAccountOpen)

		app.POST("/close-client-sca", HandleSCAClose)

		app.POST("/open-client-sca", HandleSCAOpen)

		app.POST("/mail-sender-info", HandleUpdateMailSenderInfo)

		app.POST("/get-mail-sender-info", HandleLoadMailSenderInfo)

		app.POST("/logout", Logout)

		// TODO:NEW ADDONS

		app.GET("/sec-report-preview", ShowSecReportPreview)

		app.GET("/reports", ShowReportsPage)

		app.GET("/npra-dashboard", showNpraDashboard)

		app.GET("/npra", generateNpraReport)

		// New NPRA Reports
		app.GET("/npra030", ShowNPRA030ReportPage) // To be Removed later

		app.GET("/npra0301", LoadNPRA0301ReportPage)
		//app.GET("/npra0301/get-report", HandleNPRA0301ReportData)

		app.GET("/npra0302", LoadNPRA0302ReportPage)

		app.GET("/npra03021", LoadNPRA03021ReportPage)
		app.GET("/npra03011", LoadNPRA03011ReportPage)
		//app.GET("/npra0302/get-report", HandleNPRA0302ReportData)

		app.GET("/npra0303", ShowNPRA0303Report)

		app.GET("/npra0303-add", ShowNPRA0303ReportAdd)

		app.POST("/npra0303-add", HandleNPRA0303ReportAdd)

		app.GET("/npra0303-edit", ShowNPRA0303Edit)

		app.PATCH("/npra0303-edit", HandleNPRA0303Edit)

		app.GET("/npra0303-view", ShowNPRA0303Details)

		app.POST("/pv-details-for-npra0302", Loadnpra030PVForEyeBalling)

		app.POST("/details-0301-for-eyeballing", Load301ForEyeBalling)

		//----------------------------------------
		app.GET("/npra-301-report-excel", Export301ReportToExcel)

		app.GET("/npra-302-report-excel", Export302ReportToExcel)

		app.GET("/npra-303-report-excel", Export303ReportToExcel)
		//---------------------------------------
		app.GET("/assetclass", HandleLoadAssetClass)

		app.POST("/assetclass", HandleAddAssetClass)

		app.DELETE("/assetclass", HandleDeleteAssetClass)

		app.PATCH("/assetclass", HandleEditAssetClass)
		// --------------------------------------

		app.GET("/npra-report", ShowNPRAReportPage)

		app.GET("/unidentified-report", ShowUnidentifiedReport)

		app.GET("/template-setup", ShowTemplateSetup)

		app.GET("/billing", ShowNavBillPage)

		app.GET("/profile", showProfilePage)

		app.POST("/profile", HandleUserProfileUpdate)

		app.GET("/audit-pv", EyeBallPvPage)

		app.GET("/npra-preview", ShowNpraPreview)

		app.POST("/pv-details-for-eyeballing", LoadUploadedPVForEyeBalling)

		app.GET("/billing-report", ShowBillingReportPage)

		app.GET("/client-emails", HandleLoadClientEmails)

		app.POST("/client-emails", HandleAddClientEmail)

		app.PATCH("/client-emails", HandleEditClientEmail)

		app.DELETE("/client-emails", HandleDeleteClientEmail)

		app.GET("/mail-service", HandleLoadMailService)

		app.POST("/mail-service", HandleAddMailService)

		app.POST("/holiday", HandleAddHoliday)

		app.DELETE("/holiday", HandleDeleteHoliday)

		app.PATCH("/mail-service", HandleEditMailService)

		app.DELETE("/mail-service", HandleDeleteMailService)

		app.GET("/client-readonly", HandleLoadClientsReadOnly)

		app.GET("/client-scas-readonly", HandleClientSCAsReadOnly)

		app.GET("/billed-clients", HandleLoadBilledClients)

		app.POST("/send-client-letters", HandleSendClientLetters)

		app.POST("/export-unauthorized-transactions", HandleExportUnauthorizedTransactions)

		app.GET("/securities", HandleLoadSecurities)

		app.POST("/securities", HandleAddSecurity)

		app.DELETE("/securities", HandleDeleteSecurity)

		app.PATCH("/securities", HandleEditSecurity)

		app.GET("/client-letters", HandleLetterToClients)

		app.GET("/PasswordManager", ShowPasswordManager)

		app.POST("/Password-Manager", SavePassword)

		//TODO: stop sending files over the wire as base64 and send the link using this
		app.ServeFiles("/file-manager", http.Dir(envy.Get("FILE_MANAGER_DIR", ""))) //make files in the file manager accessible

		app.ServeFiles("/", assetsBox) // serve files from the public directory
	}

	return app
}
