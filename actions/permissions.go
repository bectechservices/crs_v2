package actions

const (
	//Maker Permissions

	//UploadPV permission to upload PVs
	UploadPV = "upload_pv"
	//InputClientInformation permission to input client information
	InputClientInformation = "input_client_information"
	//InputSECInformation permission to input sec information
	InputSECInformation = "input_sec_information"
	//InputClientBankDetails permission to input client bank details
	InputClientBankDetails = "input_client_bank_details"

	//Checker Permissions

	//ViewPVDetails permission to view pv details
	ViewPVDetails = "view_pv_details"
	//ViewClientInformation permission to view client information
	ViewClientInformation = "view_client_information"
	//ApproveSECMakerRequest permission to approve maker requests
	ApproveSECMakerRequest = "approve_sec_maker_request"
	//ApproveNPRAMakerRequest permission to approve maker requests
	ApproveNPRAMakerRequest = "approve_npra_maker_request"
	//ApproveTrusteeMakerRequest permission to approve maker requests
	ApproveTrusteeMakerRequest = "approve_trustee_maker_request"
	//ApproveNAVMakerRequest permission to approve maker requests
	ApproveNAVMakerRequest = "approve_nav_maker_request"
	//AccessClientBasedAnalytics permission to access client based analytics
	AccessClientBasedAnalytics = "access_client_based_analytics"
	//AccessOverallBasedAnalytics permission to access overall analytics
	AccessOverallBasedAnalytics = "access_overall_analytics"

	//Alert Permissions

	//PVReportOutstanding permission to view pv outstanding reports alerts
	PVReportOutstanding = "pv_report_outstanding"
	//SECRequests permission to view sec request alerts
	SECRequests = "sec_requests"
	//NPRARequests permission to view npra request alerts
	NPRARequests = "npra_requests"
	//ClientLetters permission to aview client letter alerts
	ClientLetters = "client_letters"
	//MissedLiabilities permission to view missed liabilities
	MissedLiabilities = "missed_liabilities"
	//NPRAReportDelays permission to view npra report delay alerts
	NPRAReportDelays = "npra_report_delays"
	//SECReportDelays permission to view sec report delay alerts
	SECReportDelays = "sec_report_delays"
	//ClientReportDelays permission to view client report delay alerts
	ClientReportDelays = "client_report_delay"
	//InvoiceDelay permission to view invoice delay alerts
	InvoiceDelay = "invoice_delay"

	
	//InputNPRADetails permission to input npra details
	InputNPRADetails = "input_npra_details"
	//InputTrusteeDetails permission to input trustee details
	InputTrusteeDetails = "input_trustee_details"
	//InputNAVDetails permission to input nav details
	InputNAVDetails = "input_nav_details"
	//AccessMailPage permission to access mail page
	AccessMailPage = "access_mail_page"
	//ViewSentMails permission to view sent mails
	ViewSentMails = "view_sent_mails"
	//ModifyMailTemplate permission to modify mail template
	ModifyMailTemplate = "modify_mail_template"
)
