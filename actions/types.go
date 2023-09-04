package actions

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/gobuffalo/nulls"
)

// User structure of the user in the DB
type User struct {
	ID                     int              `gorm:"primary_key;column:ID" json:"id"`
	StaffID                int              `gorm:"column:STAFF_ID" json:"staff_id"`
	Fullname               string           `gorm:"column:FULLNAME" json:"fullname"`
	Email                  string           `gorm:"column:EMAIL" json:"email"`
	Password               string           `gorm:"column:PASSWORD" json:"-"`
	Password1              string           `gorm:"column:PASSWORD1" json:"-"`
	Password2              string           `gorm:"column:PASSWORD2" json:"-"`
	Password3              string           `gorm:"column:PASSWORD3" json:"-"`
	Password4              string           `gorm:"column:PASSWORD4" json:"-"`
	Password5              string           `gorm:"column:PASSWORD5" json:"-"`
	Password6              string           `gorm:"column:PASSWORD6" json:"-"`
	Password7              string           `gorm:"column:PASSWORD7" json:"-"`
	Password8              string           `gorm:"column:PASSWORD8" json:"-"`
	Password9              string           `gorm:"column:PASSWORD9" json:"-"`
	LoginTry               int              `gorm:"column:I_TRY" json:"-"`
	LastPasswordChangeDate time.Time        `gorm:"column:LASTPASSCHGDATE" json:"-"`
	MustResetPassword      bool             `gorm:"column:MUSTRESETPASSWORD" json:"-"`
	Active                 bool             `gorm:"column:ACTIVE" json:"active"`
	Locked                 bool             `gorm:"column:LOCKED" json:"locked"`
	Roles                  []UserRole       `gorm:"foreignkey:UserID" json:"roles"`
	Permissions            []UserPermission `gorm:"foreignkey:UserID" json:"permissions"`
}

// TableName sets the users table name for gorm
func (User) TableName() string {
	return "CRS_USERS"
}

// IsEmpty checks if the user is empty
func (user User) IsEmpty() bool {
	return (user.Email == "") && (user.Fullname == "")
}

// HasRole checks if the user has a role
func (user User) HasRole(role string) bool {
	for _, r := range user.Roles {
		if r.Role.Name == role {
			return true
		}
	}
	return false
}

// Can checks if the user has a permission
func (user User) Can(permission string) bool {
	for _, p := range user.Permissions {
		if p.Permission.Name == permission {
			return true
		}
	}
	for _, r := range user.Roles {
		if r.Role.Can(permission) {
			return true
		}
	}
	return false
}

// IsAdmin checks if the user is an admin
func (user User) IsAdmin() bool {
	return user.HasRole(AdminRole)
}

// IsManager checks if the user is a manager
func (user User) IsManager() bool {
	return user.HasRole(ManagerRole)
}

// IsSecMaker checks if the user has the sec maker role
func (user User) IsSecMaker() bool {
	return user.HasRole(SECMakerRole)
}

// IsSecChecker checks if the user has the sec checker role
func (user User) IsSecChecker() bool {
	return user.HasRole(SECCheckerRole)
}

//IsSecApproval checks if the user has the sec Approval role
// func (user User) IsSecApproval() bool {
// 	return user.HasRole(SECApprovalRole)
// }

// IsNPRAMaker checks if the user has the npra maker role
func (user User) IsNPRAMaker() bool {
	return user.HasRole(NPRAMakerRole)
}

// IsNPRAChecker checks if the user has the npra checker role
func (user User) IsNPRAChecker() bool {
	return user.HasRole(NPRACheckerRole)
}

//IsNPRAApprover checks if the user has the npra Approval role
// func (user User) IsNPRAApprover() bool {
// 	return user.HasRole(NPRAApprovalRole)
// }

// IsTrusteeMaker checks if the user has the trustee maker role
func (user User) IsTrusteeMaker() bool {
	return user.HasRole(TrusteeMakerRole)
}

// IsTrusteeChecker checks if the user has the trustee checker role
func (user User) IsTrusteeChecker() bool {
	return user.HasRole(TrusteeCheckerRole)
}

//IsTrusteeApproval checks if the user has the Approval checker role
// func (user User) IsTrusteeApproval() bool {
// 	return user.HasRole(TrusteeApprovalrRole)
// }

// IsNAVMaker checks if the user has the nav maker role
func (user User) IsNAVMaker() bool {
	return user.HasRole(NAVMakerRole)
}

// IsNAVChecker checks if the user has the nav checker role
func (user User) IsNAVChecker() bool {
	return user.HasRole(NAVCheckerRole)
}

//IsNAVApproval checks if the user has the nav Approval role
// func (user User) IsNAVApproval() bool {
// 	return user.HasRole(NAVApprovalRole)
// }

// Role db table
type Role struct {
	ID          int              `gorm:"column:ID" json:"id"`
	Name        string           `gorm:"column:NAME" json:"name"`
	Permissions []RolePermission `gorm:"foreignkey:RoleID" json:"permissions"`
}

// TableName sets the role table name for gorm
func (Role) TableName() string {
	return "CRS_ROLES"
}

// Can checks if a role has a permission
func (role Role) Can(permission string) bool {
	for _, perm := range role.Permissions {
		if perm.Permission.Name == permission {
			return true
		}
	}
	return false
}

// Permission user's permissions
type Permission struct {
	ID   int    `gorm:"column:ID" json:"id"`
	Name string `gorm:"column:NAME" json:"name"`
}

// TableName sets the role table name for gorm
func (Permission) TableName() string {
	return "CRS_PERMISSION"
}

// UserRole db table
type UserRole struct {
	UserID int  `gorm:"column:USERID" json:"user_id"`
	RoleID int  `gorm:"column:ROLEID" json:"role_id"`
	Role   Role `gorm:"foreignkey:RoleID" json:"role"`
}

// TableName sets the role table name for gorm
func (UserRole) TableName() string {
	return "CRS_USER_ROLE"
}

// UserPermission db table
type UserPermission struct {
	UserID       int        `gorm:"column:USERID" json:"user_id"`
	PermissionID int        `gorm:"column:PERMISSIONID" json:"permission_id"`
	Permission   Permission `gorm:"foreignkey:PermissionID" json:"permission"`
}

// TableName sets the role table name for gorm
func (UserPermission) TableName() string {
	return "CRS_USER_PERMISSION"
}

// RolePermission db table
type RolePermission struct {
	RoleID       int        `gorm:"column:ROLEID" json:"role_id"`
	PermissionID int        `gorm:"column:PERMISSIONID" json:"permission_id"`
	Permission   Permission `gorm:"foreignkey:PermissionID" json:"permission"`
}

// TableName sets the role table name for gorm
func (RolePermission) TableName() string {
	return "CRS_ROLE_PERMISSION"
}

// Client db table
type Client struct {
	ID                 int                     `gorm:"column:ID;AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	Name               string                  `gorm:"column:CLIENTNAME" json:"clientname"`
	Location           string                  `gorm:"column:LOCATION" json:"location"`
	Phone              string                  `gorm:"column:TELEPHONENO" json:"phone"`
	AccountNumber      string                  `gorm:"primary_key;column:SAFEKEEPINGACCNO" json:"safekeepingaccno"`
	ClientID           string                  `gorm:"column:CLIENT_ID" json:"client_id"`
	Type               string                  `gorm:"column:CLIENT_TYPE" json:"type"`
	Tier               string                  `gorm:"column:CLIENT_TIER" json:"tier"`
	FundManagerID      int                     `gorm:"column:FUND_MANAGER_ID" json:"fund_manager_id"`
	ContactPerson      string                  `gorm:"column:CONTACTPERSON" json:"contact_person"`
	ContactPersonPhone string                  `gorm:"column:CONTACTPERSONNO" json:"contact_person_phone"`
	HomeCountry        string                  `gorm:"column:HOME_COUNTRY" json:"country"`
	BPID               string                  `gorm:"column:BP_ID" json:"bp_id"`
	Email              string                  `gorm:"column:EMAIL" json:"email"`
	AddressLine1       string                  `gorm:"column:ADDRESS_LINE1" json:"address_line1"`
	AddressLine2       string                  `gorm:"column:ADDRESS_LINE2" json:"address_line2"`
	AddressLine3       string                  `gorm:"column:ADDRESS_LINE3" json:"address_line3"`
	AddressLine4       string                  `gorm:"column:ADDRESS_LINE4" json:"address_line4"`
	Image              string                  `gorm:"column:CLIENT_IMAGE"`
	SchemeManager      string                  `gorm:"column:SCHEME_MANAGER" json:"scheme_manager"`
	PVReportSummary    []ClientPVReportSummary `gorm:"foreignkey:BPID;association_foreignkey:BP_ID" json:"pv_report_summary"`
	Reports            []PVReportField         `gorm:"foreignkey:BPID;association_foreignkey:BP_ID" json:"reports"`
	FundManager        FundManager             `gorm:"foreignkey:FundManagerID" json:"fund_manager"`
	Closed             *time.Time              `gorm:"column:CLOSED" json:"closed"`
}

// TableName sets the role table name for gorm
func (Client) TableName() string {
	return "CRS_CLIENT"
}

func (client Client) ImageBuffer() string {
	if client.Image == "" {
		return ""
	}
	ext := filepath.Ext(client.Image)
	var mime = "image/jpeg"
	if ext == ".png" {
		mime = "image/png"
	}
	buffer := ConvertFileToBase64(client.Image)
	if buffer != "" {
		return "data:" + mime + ";base64," + buffer
	}
	return ""
}

func (client Client) WasClosedAsAt(date time.Time) bool {
	if client.Closed == nil {
		return false
	}
	return client.Closed.Before(date)
}

// ClientPVReportSummary the client report summary
type ClientPVReportSummary struct {
	ID                int                `gorm:"column:ID;AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	AccountNumber     string             `gorm:"column:SAFEKEEPINGACCNO" json:"account_number"`
	BPID              string             `gorm:"column:BP_ID" json:"bpid"`
	ClientID          string             `gorm:"column:CLIENT_ID" json:"client_id"`
	ReportDate        *time.Time         `gorm:"column:REPORT_DATE" json:"report_date"`
	SecurityType      string             `gorm:"column:SECURITY_TYPE" json:"security_type"`
	NominalValue      float64            `gorm:"column:NOMINAL_VALUE" json:"nominal_value"`
	CumulativeCost    float64            `gorm:"column:CUMULATIVE_COST" json:"cumlative_cost"`
	LCYAmount         float64            `gorm:"column:LCY_AMT" json:"lcy_amount"`
	PercentageOfTotal float32            `gorm:"column:PERCENTAGE_TOTAL" json:"percentage_of_total"`
	ClientSCADetails  ClientPVUniqueCode `gorm:"foreignkey:ClientID;association_foreignkey:CODE" json:"client"`
}

// TableName sets the role table name for gorm
func (ClientPVReportSummary) TableName() string {
	return "CRS_PVSUMMARY"
}

// PVReportField the fields in the NAV report
type PVReportField struct {
	ID                int                `gorm:"column:ID;AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	AccountNumber     string             `gorm:"column:SAFEKEEPINGACCNO" json:"account_number"`
	BPID              string             `gorm:"column:BP_ID" json:"bpid"`
	ClientID          string             `gorm:"column:CLIENT_ID" json:"client_id"`
	SecurityName      string             `gorm:"column:SECURITY_NAME" json:"security_name"`
	CDSCode           string             `gorm:"column:CDS_CODE" json:"cds_code"`
	ISIN              string             `gorm:"column:ISIN" json:"isin"`
	SCBCode           string             `gorm:"column:SCBCODE" json:"scb_code"`
	MarketPrice       float32            `gorm:"column:MARKET_PRICE" json:"market_price"`
	NominalValue      float64            `gorm:"column:NOMINAL_VALUE" json:"nominal_value"`
	CumulativeCost    float64            `gorm:"column:CUMULATIVE_COST" json:"cumulative_cost"`
	Value             float64            `gorm:"column:VALUE_AMOUNT" json:"value_amount"`
	SecurityType      string             `gorm:"column:SEC_TITLE" json:"security_type"`
	PercentageOfTotal float32            `gorm:"column:PERCENTAGE_OF_TOTAL" json:"percentage_of_total"`
	DateFrom          string             `gorm:"column:DATE_FROM" json:"date_from"`
	DateTo            string             `gorm:"column:DATE_TO" json:"date_to"`
	ReportDate        time.Time          `gorm:"column:REPORT_DATE" json:"report_date"`
	ClientSCADetails  ClientPVUniqueCode `gorm:"foreignkey:ClientID;association_foreignkey:CODE" json:"client"`
}

// TableName sets the role table name for gorm
func (PVReportField) TableName() string {
	return "CRS_PVREPORTS"
}

// BondOverallSummary loads the overall bond data
type BondOverallSummary struct {
	Bond  string  `gorm:"column:bond" json:"security"`
	Value float64 `gorm:"column:value" json:"value"`
}

// QuarterDate the structure of a quarter date
type QuarterDate struct {
	Begin string
	End   string
}

// CrsSecCusAfiliateTxnRel sets the role table name for gorm
type CrsSecCusAfiliateTxnRel struct {
	// ID                int       `gorm:"column:ID" json:"id"`
	ReportRefID              string    `gorm:"column:REPORT_REF_ID" json:"REPORT_REF_ID"`
	SafeKeepingAccount       string    `gorm:"column:SAFEKEEPINGACCNO" json:"SAFEKEEPINGACCNO"`
	CustodianTrusteeName     string    `gorm:"column:CUSTODIAN_TRUSTEE_NAME" json:"CUSTODIAN_TRUSTEE_NAME"`
	RelationTrusteeCustodian string    `gorm:"column:RELATION_TRUSTEE_CUSTODIAN" json:"RELATION_TRUSTEE_CUSTODIAN"`
	TypeTransaction          string    `gorm:"Column:TYPE_TRANSACTION" json:"TYPE_TRANSACTION"`
	AmountInvolved           float64   `gorm:"Column:AMT_INVLOVED" json:"AMT_INVLOVED"`
	TransactionDate          time.Time `gorm:"Column:TRANSACTION_DATE" json:"TRANSAction_DATE"`
}

// TableName sets the role table name for gorm
func (CrsSecCusAfiliateTxnRel) TableName() string {
	return "CRS_SEC_cus_AFILIATE_TXN_REL"
}

// CrsSecCustodianInfo the fields in the SEC report
type CrsSecCustodianInfo struct {
	// ID                int       `gorm:"column:ID" json:"id"`
	ReportRefID      string    `gorm:"column:REPORT_REF_ID" json:"REPORT_REF_ID"`
	SafeKeepingAccNo string    `gorm:"SAFEKEEPINGACCNO:SAFEKEEPINGACCNO" json:"NO_SHARE_OUTSTANDING"`
	CustodianName    string    `gorm:"CUSTODIAM:SAFEKEEPINGACCNO" json:"SAFEKEEPINGACCNO"`
	ReportingPeriod  string    `gorm:"REPORTING_PERIOD" jason:"REPORTING_PERIOD"`
	ReportingOfficer string    `gorm:"REPORTING_OFFICER" json:"REPORTING_OFFICER"`
	RenewalDate      time.Time `gorm:"RENEWAL_DATE" json:"RENEWAL_DATE"`
	SustantialID     string    `gorm:"SUBSTANTIALID" json:"SUBSTANTIALID"`
	Remark           string    `gorm:"REMARK" json:"REMARK"`
	VerfierID        int32     `gorm:"VERIFIER_ID" json:"VERIFIER_ID"`
	ReviewedDate     time.Time `gorm:"REVIEWED_DATE" json:"REVIEWED_DATE"`
	Signature        string    `gorm:"SIGNATURE" json:"SIGNATURE"`
}

// CrsSecOffshoreClients the fields in the SEC report
type CrsSecOffshoreClients struct {
	// ID                int       `gorm:"column:ID" json:"id"`
	ReportRefID string `gorm:"column:REPORT_REF_ID" json:"REPORT_REF_ID"`
	Client      string `gorm:"column:Client" json:"Client"`
	BPID        string `gorm:"column:bpid" json:"bpid"`
	HomeCountry string `gorm:"column:Home_Country" json:"Home_Country"`
	AssetValue  string `gorm:"column:Asset_Value" json:"Asset_Value"`
}

// TableName sets the role table name for gorm
func (CrsSecOffshoreClients) TableName() string {
	return "CRS_SEC_OFFSHORE_CLIENTS"
}

// CrsSecOtherInformation the fields in the SEC report
type CrsSecOtherInformation struct {
	// ID                int       `gorm:"column:ID" json:"id"`
	ReportRefID               string    `gorm:"column:REPORT_REF_ID" json:"REPORT_REF_ID"`
	ClaimSchemeAsset          string    `gorm:"column:CLAIM_SCHEME_ASSET" json:"CLAIM_SCHEME_ASSET"`
	ClaimInformApproved       bool      `gorm:"column:CLAIM_INFORMED_APPROVED" json:"CLAIM_INFORMED_APPROVED"`
	LitigationCustodianScheme bool      `gorm:"column:LITIGATION_CUSTODIAN_SCHEME" json:"LITIGATION_CUSTODIAN_SCHEME"`
	SigReductionScheme        string    `gorm:"column:SIG_REDUCTION_SCHEME" json:"SIG_REDUCTION_SCHEME"`
	ReconciledAssetCustodian  string    `gorm:"column:RECONCILED_ASSET_CUSTODIAN" json:"RECONCILED_ASSET_CUSTODIAN"`
	TimesSchemePublished      int32     `gorm:"column:TIMES_SCHEME_PUBLISHED" json:"TIMES_SCHEME_PUBLISHED"`
	ConcernsByInvestors       string    `gorm:"column:CONCERNS_BY_INVESTORS" json:"CONCERNS_BY_INVESTORS"`
	ConcernsMgtcustoAfund     string    `gorm:"column:CONCERNS_MGTCUSTO_AFUND" json:"CONCERNS_MGTCUSTO_AFUND"`
	AccountSeparateScheme     bool      `gorm:"column:ACCOUNT_SEPARATE_SCHEME" json:"ACCOUNT_SEPARATE_SCHEME"`
	LitigationByStakeholder   bool      `gorm:"column:LITIGATION_BY_STAKEHOLDER" json:"LITIGATION_BY_STAKEHOLDER"`
	SecInfor                  string    `gorm:"column:SEC_INFOR" json:"SEC_INFOR"`
	TransactionDate           time.Time `gorm:"column:TRANSACTION_DATE" json:"TRANSACTION_DATE"`
}

// TableName sets the role CRS_SEC_OTHER_INFORMATION for gorm
func (CrsSecOtherInformation) TableName() string {
	return "CRS_SEC_OTHER_INFORMATION"
}

// CrsSecSchemeMarkedVolume the fields in the SEC report
type CrsSecSchemeMarkedVolume struct {
	// ID                int       `gorm:"column:ID" json:"id"`
	ReportRefID             string    `gorm:"column:REPORT_REF_ID" json:"REPORT_REF_ID"`
	NameOfFirm              string    `gorm:"column:Name_of_firm" json:"Name_of_firm"`
	NameOfScheme            string    `gorm:"column:Name_of_Scheme" json:"Name_of_Scheme"`
	RelWithCustodianTrustee string    `gorm:"column:Rel_with_Custodian_Trustee" json:"Rel_with_Custodian_Trustee"`
	Volume                  int64     `gorm:"column:Volume" json:"Volume"`
	MarkedToMarketValue     string    `gorm:"column:Marked_to_market_value" json:"Marked_to_market_value"`
	TransactionDate         time.Time `gorm:"column:TRANSACTION_DATE" json:"TRANSACTION_DATE"`
}

// TableName sets the role CRS_SEC_SCHEME_MARKED_VOLUME for gorm
func (CrsSecSchemeMarkedVolume) TableName() string {
	return "CRS_SEC_SCHEME_MARKED_VOLUME"
}

// CrsSustantialShares the fields in the SEC report
type CrsSustantialShares struct {
	// ID                int       `gorm:"column:ID" json:"id"`
	SustantialID  int32   `gorm:"column:SUSTANTIALID" json:"SUSTANTIALID"`
	ShareName     string  `gorm:"column:SHARESNAME" json:"SHARESNAME"`
	ShareType     string  `gorm:"column:SHARETYPE" json:"SHARETYPE"`
	Shareholdings float64 `gorm:"column:SHAREHOLDINGS" json:"SHAREHOLDINGS"`
	Percentage    int64   `gorm:"column:PERCENTAGE" json:"PERCENTAGE"`
}

// TableName sets the role CRS_SEC_SCHEME_MARKED_VOLUME for gorm
func (CrsSustantialShares) TableName() string {
	return "CRS_SUSTANTIAL_SHARES"
}

// CrsTaskTrail the fields in the TaskTrail report
type CrsTaskTrail struct {
	// ID                int       `gorm:"column:ID" json:"id"`
	MakerID          int32     `gorm:"column:MAKERID" json:"MAKERID"`
	CheckerID        int32     `gorm:"column:CHECKERID" json:"CHECKERID"`
	ApprovalID       int32     `gorm:"column:APPROVALID" json:"APPROVALID"`
	FileReportType   string    `gorm:"column:FILE_REPORT_TYPE" json:"FILE_REPORT_TYPE"`
	Filepath         string    `gorm:"column:FILE_PATH" json:"FILE_PATH"`
	MakerFlag        int       `gorm:"column:MAKER_FLAG" json:"MAKER_FLAG"`
	CheckerFlag      int       `gorm:"column:CHECKER_FLAG" json:"CHECKER_FLAG"`
	ApprovalFlag     int       `gorm:"column:APPROVAL_FLAG" json:"APPROVAL_FLAG"`
	MakerUpdateAt    time.Time `gorm:"column:MAKER_UPDATE_AT" json:"MAKER_UPDATE_AT"`
	CheckerUpdateAt  time.Time `gorm:"column:CHECKER_UPDATE_AT" json:"CHECKER_UPDATE_AT"`
	ApprovalUpdateAt time.Time `gorm:"column:APPROVAL_UPDATE_AT" json:"APPROVAL_UPDATE_AT"`
}

// TableName sets the role table for gorm
func (CrsTaskTrail) TableName() string {
	return "CRS_TASKTRAIL"
}

// CrsTiers2Tier3 the fields in the TaskTrail report
type CrsTiers2Tier3 struct {
	// ID                int       `gorm:"column:ID" json:"id"`
	ChargeCode   string `gorm:"column:CHARGECODE" json:"CHARGECODE"`
	ChargeDesc   string `gorm:"column:CHARGE_DESC" json:"CHARGE_DESC"`
	ChargeType   string `gorm:"column:CHARGE_TYPE" json:"CHARGE_TYPE"`
	Basespoint   int32  `gorm:"column:BASESPOINT" json:"BASESPOINT"`
	ChargesTier1 string `gorm:"column:CHARGESTIER1" json:"CHARGESTIER1"`
}

// TableName sets the role table for gorm
func (CrsTiers2Tier3) TableName() string {
	return "CRS_TIER2TIER3_CHARGES"
}

// CrsTrustTxnVolumes the fields in the TaskTrail report
type CrsTrustTxnVolumes struct {
	ID               int       `gorm:"column:ID;AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	BPID             string    `gorm:"column:BP_ID" json:"bpid"`
	StockSettledDate time.Time `gorm:"column:Stock_Settled_Date" json:"stock_settled_date"`
	SecurityType     string    `gorm:"column:Security_Type" json:"security_type"`
	Hash             string    `json:"hash" gorm:"column:HASH"`
	QuarterDate      time.Time `json:"quarter_date" gorm:"column:QUARTER_DATE"`
}

// TableName sets the role table for gorm
func (CrsTrustTxnVolumes) TableName() string {
	return "CRS_TRUST_TXN_VOLUMES"
}

// CrsVolumeShareScheme for the fields in the Table CRS_VOLUME_SHARE_SCHEME
type CrsVolumeShareScheme struct {
	// ID                int       `gorm:"column:ID" json:"id"`
	FirmName          string  `gorm:"column:FIRM_NAME" json:"FIRM_NAME"`
	RelationCustodian string  `gorm:"column:RELATION_CUSTODIAN" json:"RELATION_CUSTODIAN"`
	Volume            int32   `gorm:"column:VOLUME" json:"VOLUME"`
	MarkMarketValue   float64 `gorm:"column:MARK_MARKET_VALUE" json:"MARK_MARKET_VALUE"`
}

// TableName describes the volume of share schemes
func (CrsVolumeShareScheme) TableName() string {
	return "CRS_VOLUME_SHARE_SCHEME"
}

// CrsCashEntitlementList for the cash entitlement Listing
type CrsCashEntitlementList struct {
	// ID                int       `gorm:"column:ID" json:"id"`
	SecurityFundAccountID string  `gorm:"column:Security_Fund_Account_Id" json:"Security_Fund_Account_Id"`
	ExercisedQty          int64   `gorm:"column:Exercised_Qty" json:"Exercised_Qty"`
	GrossAmount           float64 `gorm:"column:Gross_Amount" json:"Gross_Amount"`
	TaxAmount             float64 `gorm:"column:Tax_Amount" json:"Tax_Amount"`
	NetAmount             float64 `gorm:"column:Net_Amount" json:"Net_Amount"`
	BaseSecurityID        string  `gorm:"column:Base_Security_Id" json:"Base_Security_Id"`
	Status                string  `gorm:"column:Status" json:"Status"`
	PositionBucket        string  `gorm:"column:Position_Bucket" json:"Position_Bucket"`
	UpdatedBy             int     `gorm:"column:UPDATED_BY" json:"UPDATED_BY"`
}

// TableName describes the cash entitlement Listing
func (CrsCashEntitlementList) TableName() string {
	return "CRS_CASH_ENTITLEMENT_LIST"
}

// CrsNpraDetail for the NPRA Details
type CrsNpraDetail struct {
	// ID                int       `gorm:"column:ID" json:"id"`
	ContactPersonName        string `gorm:"column:CONTACT_PERSON_NAME" json:"CONTACT_PERSON_NAME"`
	ContactPersonPhone       string `gorm:"column:CONTACT_PERSON_PHONE" json:"CONTACT_PERSON_PHONE"`
	PensionFundCustodian     string `gorm:"column:PENSION_FUND_CUSTODIAN" json:"PENSION_FUND_CUSTODIAN"`
	FundCustodianID          string `gorm:"column:FUND_CUSTODIAN_ID" json:"FUND_CUSTODIAN_ID"`
	CustodianReportSerialNum string `gorm:"column:CUSTODIAN_REPORT_SERIAL_NUM" json:"CUSTODIAN_REPORT_SERIAL_NUM"`
	NpraSummaryID            string `gorm:"column:NPRA_SUMMARY_ID" json:"NPRA_SUMMARY_ID"`
}

// TableName describes the NPRA Details
func (CrsNpraDetail) TableName() string {
	return "CRS_NPRA_DETAIL"
}

// CrsNpraMonthlyStatistics for the NPRA Monthly Statistics
type CrsNpraMonthlyStatistics struct {
	// ID                int       `gorm:"column:ID" json:"id"`
	NameScheme     string    `gorm:"column:NAME_SCHEME" json:"NAME_SCHEME"`
	InvestmentDate time.Time `gorm:"column:INVESTMENT_DATE" json:"INVESTMENT_DATE"`
	IssuerName     string    `gorm:"column:ISSUER_NAME" json:"ISSUER_NAME"`
	AssetClass     string    `gorm:"column:ASSET_CLASS" json:"ASSET_CLASS"`
	AssetTenor     string    `gorm:"column:ASSET_TENOR" json:"ASSET_TENOR"`
	AmountInvolved float64   `gorm:"column:AMOUNT_INVLOVED" json:"AMOUNT_INVLOVED"`
	ExpectedReturn float64   `gorm:"column:EXPECTED_RETURN" json:"EXPECTED_RETURN"`
}

// TableName describes the NPRA Monthly Statistics
func (CrsNpraMonthlyStatistics) TableName() string {
	return "CRS_NPRA_MTHLY_STAT"
}

// OrdinaryShare the structure of an ordinary share
type OrdinaryShare struct {
	ID               int       `gorm:"column:ID;AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	GovernanceInfoID int       `gorm:"column:GOV_INFO_ID" json:"governaceInfoID"`
	Name             string    `gorm:"column:NAME" json:"name"`
	Shareholdings    int       `gorm:"column:SHAREHOLDINGS" json:"shareholdings"`
	Percentage       float32   `gorm:"column:PERCENTAGE" json:"percentage"`
	CreatedAt        time.Time `gorm:"column:CREATED_AT" json:"-"`
}

// TableName for Pension charges
func (OrdinaryShare) TableName() string {
	return "CRS_SEC_GOV_ORDINARY_SHARES"
}

// PreferenceShare the structure of a preference share
type PreferenceShare struct {
	ID               int       `gorm:"column:ID;AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	GovernanceInfoID int       `gorm:"column:GOV_INFO_ID" json:"governaceInfoID"`
	Name             string    `gorm:"column:NAME" json:"name"`
	Shareholdings    int       `gorm:"column:SHAREHOLDINGS" json:"shareholdings"`
	Percentage       float32   `gorm:"column:PERCENTAGE" json:"percentage"`
	CreatedAt        time.Time `gorm:"column:CREATED_AT" json:"-"`
}

// TableName for Pension charges
func (PreferenceShare) TableName() string {
	return "CRS_SEC_GOV_PREFERENCE_SHARES"
}

// AffiliateTransaction describes the structure of an affiliate transaction
type AffiliateTransaction struct {
	ID                        int       `gorm:"column:ID;AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	GovernanceInfoID          int       `gorm:"column:GOV_INFO_ID" json:"governaceInfoID"`
	NameOfAffiliate           string    `gorm:"column:NAME_OF_AFFILIATE" json:"nameOfAffiliate"`
	RelationshipWithCustodian string    `gorm:"column:REL_WITH_CUSTODIAN_TRUSTEE" json:"relationshipWithCustodian"`
	TypeOfTransaction         string    `gorm:"column:TYPE_OF_TXN" json:"typeOfTransaction"`
	Amount                    float64   `gorm:"column:AMOUNT_INVOLVED" json:"amount"`
	CreatedAt                 time.Time `gorm:"column:CREATED_AT" json:"-"`
}

// TableName for Pension charges
func (AffiliateTransaction) TableName() string {
	return "CRS_SEC_TXN_WITH_CUSTODIAN_AFFILIATES"
}

// ValueVolumeOfShareScheme describes the structure of the value volumes of share schemes
type ValueVolumeOfShareScheme struct {
	ID                        int       `gorm:"column:ID;AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	GovernanceInfoID          int       `gorm:"column:GOV_INFO_ID" json:"governaceInfoID"`
	NameOfFirm                string    `gorm:"column:NAME_OF_FIRM" json:"nameOfFirm"`
	NameOfScheme              string    `gorm:"column:NAME_OF_SCHEME" json:"nameOfScheme"`
	RelationshipWithCustodian string    `gorm:"column:REL_WITH_CUSTODIAN_TRUSTEE" json:"relationshipWithCustodian"`
	Volume                    float64   `gorm:"column:VOLUME" json:"volume"`
	MarkedToMarketValue       float32   `gorm:"column:MARKED_TO_MARKET_VALUE" json:"markedToMarketValue"`
	CreatedAt                 time.Time `gorm:"column:CREATED_AT" json:"-"`
}

// TableName for Pension charges
func (ValueVolumeOfShareScheme) TableName() string {
	return "CRS_SEC_GOV_VALUE_VOLUME_OF_SHARE_SCHEMES"
}

// GovernanceInfo the structure of the sec governance info
type GovernanceInfo struct {
	ID                                 int                        `gorm:"column:ID;AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	ReportRefID                        string                     `gorm:"column:REPORT_REF_ID" json:"reportRefID"`
	CustodianName                      string                     `gorm:"column:CUSTODIAN_NAME" json:"custodianName"`
	ReportingPeriod                    string                     `gorm:"column:REPORTING_PERIOD" json:"reportingPeriod"`
	ReportingOfficer                   string                     `gorm:"column:REPORTING_OFFICER" json:"reportingOfficer"`
	DateOfReport                       time.Time                  `gorm:"column:DATE_OF_REPORT" json:"dateOfReport"`
	RenewalDate                        time.Time                  `gorm:"column:RENEWAL_DATE" json:"renewalDate"`
	ChangeInDirectors                  string                     `gorm:"column:CHANGE_IN_DIRECTORS" json:"changeInDirectors"`
	ChangeInAgreement                  string                     `gorm:"column:CHANGE_IN_AGREEMENT" json:"changeInAgreement"`
	DealingsApprovedByBoard            string                     `gorm:"column:DEALINGS_APPROVED_BY_BOARD" json:"dealingsApprovedByBoard"`
	CustodianHasUpdatedAssetRegister   string                     `gorm:"column:CUSTODIAN_HAS_UPDATED_ASSET_REGISTER" json:"custodianHasUpdatedAssetRegister"`
	CustodianAssetRegistrationDate     time.Time                  `gorm:"column:CUSTODIAN_ASSET_REGISTRATION_DATE" json:"custodianAssetRegistrationDate"`
	DoManagersOfTheSchemeConsultTheLaw string                     `gorm:"column:DO_MANAGERS_OF_THE_SCHEME_CONSULT_THE_LAW" json:"doManagersOfSchemeConsultTheLaw"`
	SchemeHadAnyOtherFinancialDealings string                     `gorm:"column:SCHEME_HAD_ANY_OTHER_FINANCIAL_DEALINGS" json:"schemeHadAnyOtherFinancialDealings"`
	Approved                           bool                       `gorm:"column:APPROVED" json:"approved"`
	ApprovedBy                         int                        `gorm:"column:APPROVED_BY" json:"approvedBy"`
	OrdinaryShares                     []OrdinaryShare            `gorm:"foreignkey:GovernanceInfoID" json:"ordinaryShares"`
	PreferenceShares                   []PreferenceShare          `gorm:"foreignkey:GovernanceInfoID" json:"preferenceShares"`
	AffiliateTransactions              []AffiliateTransaction     `gorm:"foreignkey:GovernanceInfoID"  json:"affiliateTransactions"`
	ValueVolumeOfShareSchemes          []ValueVolumeOfShareScheme `gorm:"foreignkey:GovernanceInfoID" json:"valueVolumeOfShares"`
	CreatedAt                          time.Time                  `gorm:"column:CREATED_AT" json:"-"`
	UpdatedAt                          time.Time                  `gorm:"column:UPDATED_AT" json:"-"`
}

// TableName for GovernanceInfo
func (GovernanceInfo) TableName() string {
	return "CRS_SEC_GOVERNANCE_INFO"
}

// IsEmpty checks if the governance is empty
func (governance GovernanceInfo) IsEmpty() bool {
	return (governance.CustodianName == "") && (governance.ReportingOfficer == "") && (governance.ReportingPeriod == "")
}

// SchemeDetails structure of the sec scheme details structure
type SchemeDetails struct {
	ID                                      int       `gorm:"column:ID;AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	GovernanceInfoID                        int       `gorm:"column:GOV_INFO_ID" json:"governaceInfoID"`
	BPID                                    string    `gorm:"column:BP_ID" json:"bpid"`
	NameOfScheme                            string    `gorm:"column:NAME_OF_SCHEME" json:"nameOfScheme"`
	NumberOfSharesOutstanding               float64   `gorm:"column:NO_SHARE_OUTSTANDING" json:"numberOfSharesOutstanding"`
	NumberOfShareholders                    int       `gorm:"column:NO_OF_SHAREHOLDERS" json:"numberOfShareholders"`
	NumberOfSharesRedeemed                  float64   `gorm:"column:NO_SHARES_REDEEMED" json:"numberOfRedemptions"`
	ValueOfSharesRedeemed                   float64   `gorm:"column:VALUE_SHARES_REDEEMED" json:"valueOfRedemptions"`
	NameOfManager                           string    `gorm:"column:NAME_OF_MANAGER" json:"nameOfManager"`
	TotalValueOfSchemeAssets                float64   `gorm:"column:TOTAL_VALUE_OF_SCHEME_ASSETS" json:"totalValueOfSchemeAssets"`
	NetAssetValueOfScheme                   float64   `gorm:"column:NET_ASSET_VALUE_SCHEME" json:"netAssetValueOfScheme"`
	NetAssetValueOfSchemePerUnit            float64   `gorm:"column:NET_ASSET_VALUE_SHARE_PER_UNIT" json:"netAssetValuePerShare"`
	TotalEquityInvestment                   float64   `gorm:"column:TOTAL_EQUITY_INVESTMENT" json:"totalEquityInvestments"`
	TotalFixedIncomeInvestment              float64   `gorm:"column:TOTAL_FIXED_INCOME_INVESTMENT" json:"totalFixedIncomeInvestments"`
	NetMediumAssetHeldByFund                float64   `gorm:"column:NET_MEDIUM_ASSET_HELD_BY_FUND" json:"netOfMediumTermAssetsHeldByFund"`
	CapitalMarkerInvestments                float64   `gorm:"column:CAPITAL_MARKET_INVESTMENTS" json:"capitalMarketsInvestments"`
	PercentageCapitalMarketInvestment       float32   `gorm:"column:PERCENTAGE_CAPITAL_MKT_INVESTMENT_TINVESTMENT" json:"percentageOfCapitalInvestmentToTotalInvestment"`
	CertificatesOfInvestmentWithCustodian   string    `gorm:"column:CERTIFICATES_0F_INVESTMENT_WITH_CUSTODIAN" json:"areAllCertificatesOfInvestmentWithCustodian"`
	TotalValueOfUnutilizedFunds             float64   `gorm:"column:TOTAL_VALUE_OF_UNUTILIZED_FUNDS" json:"totalValueOfUnutilizedFunds"`
	ValueOfBorrowedFunds                    float64   `gorm:"column:VALUE_OF_BORROWED_FUNDS" json:"valueOfBorrowedFunds"`
	ReasonsForBorrowing                     string    `gorm:"column:REASON_FOR_BORROWING" json:"reasonsForBorrowing"`
	AttachedFile                            string    `gorm:"column:ATTACHED_FILE" json:"attachedFile"`
	AuditedAccountsDistributedToAuthorities string    `gorm:"column:AUDITED_ACCOUNTS_DISTRIBUTTED_TO_AUTHORITIES" json:"wereAllDulyPreparedAccountsDistributed"`
	Redemptions                             float64   `gorm:"column:REDEMPTIONS" json:"redemptions"`
	Dividends                               float64   `gorm:"column:DIVIDENDS" json:"dividends"`
	Rights                                  float64   `gorm:"column:RIGHTS" json:"rights"`
	FeesOwedCustodian                       float64   `gorm:"column:FEES_OWED_CUSTODIANS" json:"feesOwedCustodian"`
	CreatedAt                               time.Time `gorm:"column:CREATED_AT" json:"-"`
	UpdatedAt                               time.Time `gorm:"column:UPDATED_AT" json:"-"`
}

// TableName for SchemeDetails
func (SchemeDetails) TableName() string {
	return "CRS_SEC_SCHEME_ASSET_HOLDINGS"
}

// OtherInformation structure of the other sec information
type OtherInformation struct {
	ID                                          int       `gorm:"column:ID;AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	GovernanceInfoID                            int       `gorm:"column:GOV_INFO_ID" json:"governaceInfoID"`
	AreThereAnyClaimOnSchemeAsset               string    `gorm:"column:ANY_CLAIM_ON_SCHEME_ASSET" json:"areThereAnyClaimOnSchemeAsset"`
	YesWasCustodianInformedAndApproved          string    `gorm:"column:YES_WAS_CUSTODIAN_INFORMED_AND_APPROVED" json:"yesWasCustodianInformedAndApproved"`
	AnyLitigationInvolvingCustodianSccheme      string    `gorm:"column:ANY_LITIGATION_INVOLVING_CUSTODIAN_SCHEME" json:"anyLitigationInvolvingCustodianSccheme"`
	AnySignificantReductionInAssetScheme        string    `gorm:"column:ANY_SIGNIFICANT_REDUCTION_IN_ASSET_SCHEME" json:"anySignificantReductionInAssetScheme"`
	SignificantReductionInSchemeMarketPrice     string    `gorm:"column:ANY_SIGNIFICANT_REDUCTION_IN_MARKET_PRICE" json:"significantReductionInSchemeMarketPrice"`
	HasMgrsReconciledAssetRegisterCustodian     string    `gorm:"column:HAS_MGRS_RECONCILED_ASSET_REGISTER_CUSTODIAN" json:"hasMgrsReconciledAssetRegisterCustodian"`
	HowManyTimesDidSchemePublishedPrices        string    `gorm:"column:HOW_MANY_TIMES_DID_SCHEME_PUBLISHED_PRICES" json:"howManyTimesDidSchemePublishedPrices"`
	AnyConcernsByInvestors                      string    `gorm:"column:ANY_CONCERNS_BY_INVESTORS" json:"anyConcernsByInvestors"`
	AnyMattersAttentionSecMgtCustodyOfFund      string    `gorm:"column:ANY_MATTERS_ATTENTION_SEC_MGT_CUSTODY_OF_FUND" json:"anyMattersAttentionSecMgtCustodyOfFund"`
	HasAccountOfManagersSeparateFromScheme      string    `gorm:"column:HAS_ACCOUNT_OF_MGRS_SEPARATE_FROM_SCHEME" json:"hasAccountOfManagersSeparateFromScheme"`
	CompanyParentsAffiliateInvolvedInLitigation string    `gorm:"column:COMPANY_PARENTS_AFFILIATE_INVOLVED_IN_LITIGATION" json:"companyParentsAffiliateInvolvedInLitigation"`
	LitigationDetails                           string    `gorm:"column:LITIGATION_DATA" json:"litigationDetails"`
	CreatedAt                                   time.Time `gorm:"column:CREATED_AT" json:"-"`
}

// TableName for OtherInformation
func (OtherInformation) TableName() string {
	return "CRS_SEC_OTHER_INFORMATION"
}

// OfficialReportRemarks structure of the sec report remarks
type OfficialReportRemarks struct {
	ID               int       `gorm:"column:ID;AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	GovernanceInfoID int       `gorm:"column:GOV_INFO_ID" json:"governaceInfoID"`
	Remarks          string    `gorm:"column:REMARKS" json:"remarks"`
	ReviewingOfficer string    `gorm:"column:NAME_OF_REVIEWING_OFFICER" json:"reviewingOfficer"`
	Date             string    `gorm:"column:DATE" json:"date"`
	Signature        string    `gorm:"column:SIGNATURE" json:"signature"`
	CreatedAt        time.Time `gorm:"column:CREATED_AT" json:"-"`
}

// TableName for OfficialReportRemarks
func (OfficialReportRemarks) TableName() string {
	return "CRS_SEC_REPORT_REMARKS"
}

// KafkaPVReport the structure of the message sent through kafka
type KafkaPVReport struct {
	UserID      int     `json:"user_id"`
	Filename    string  `json:"filename"`
	PVType      string  `json:"pv_type"`
	CashBalance float64 `json:"cashbalance"`
	IsExcel97   bool    `json:"is_excel_97"`
	ExcelData   string  `json:"excel_data"`
}

// OffshoreClient the structure of an offshore client
type OffshoreClient struct {
	Name       string  `gorm:"column:name" json:"name"`
	Country    string  `gorm:"column:country" json:"country"`
	AssetValue float64 `gorm:"column:assets" json:"assetValue"`
}

// VarianceQueryData structure of a variance query data
type VarianceQueryData struct {
	Name       string  `gorm:"column:name"`
	Country    string  `gorm:"column:country"`
	LastAUA    float64 `gorm:"column:l_aua"`
	CurrentAUA float64 `gorm:"column:c_aua"`
	Remarks    string  `gorm:"column:remarks"`
}

// VarianceData the structure of a variance data
type VarianceData struct {
	Name       string  `json:"name"`
	Country    string  `json:"country"`
	LastAUA    float64 `json:"last_aua"`
	CurrentAUA float64 `json:"current_aua"`
	Amount     float64 `json:"amount"`
	Variance   float32 `json:"variance"`
	Remarks    string  `json:"remarks"`
}

// OverviewSummaryData the structure of the summary on the dashboard
type OverviewSummaryData struct {
	PrettyDate            string
	Status                bool
	IsMoreThanLastQuarter bool
	Percentage            int
	QuarterDate           string
	Data                  []float64
}

// UploadedData the structure of an uploaded data
type UploadedData struct {
	Uploaded int
	Pending  int
}

// PrettyPercentage stringifies the percentage data uploaded
func (data UploadedData) PrettyPercentage() string {
	total := data.Pending + data.Uploaded
	diff := total - data.Uploaded
	percentage := int((float32(diff) / float32(total)) * 100)
	if percentage < 0 {
		percentage = 0
	}
	return fmt.Sprintf("%d%%", percentage)
}

// CountQuery holds the resuls of count
type CountQuery struct {
	Count int
}

// AlertSettings the structure of the alert settings stored in the db
type AlertSettings struct {
	ID                int `gorm:"column:ID;AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	DaysToReportPV    int `gorm:"column:DAYS_TO_REPORT_PV"`
	DaysToReportEmail int `gorm:"column:DAYS_TO_REPORT_EMAIL"`
}

// TableName for OtherInformation
func (AlertSettings) TableName() string {
	return "CRS_ALERT_SETTINGS"
}

// Alert the structure of an alert
type Alert struct {
	Type   string `json:"type"`
	Number int    `json:"number"`
	Due    string `json:"due"`
}

// SecActivities structure of the sec activities table
type SecActivities struct {
	ID          int       `gorm:"column:ID;AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	BPID        string    `gorm:"column:BP_ID" json:"bpid"`
	UserID      int       `gorm:"column:USER_ID" json:"user_id"`
	Date        time.Time `gorm:"column:DATE" json:"date"`
	Activity    string    `gorm:"column:ACTIVITY" json:"activity"`
	QuarterDate time.Time `gorm:"column:QUARTER_DATE" json:"quarter_date"`
	Approved    bool      `gorm:"column:APPROVED" json:"approved"`
	User        User      `gorm:"foreignkey:UserID" json:"user"`
	Client      Client    `gorm:"foreignkey:BPID;association_foreignkey:BP_ID" json:"client"`
}

// TableName for OtherInformation
func (SecActivities) TableName() string {
	return "CRS_SEC_ACTIVITIES"
}

// BillingActivities structure of the billing activities table
type BillingActivities struct {
	ID          int       `gorm:"column:ID;AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	BPID        string    `gorm:"column:BP_ID" json:"bpid"`
	UserID      int       `gorm:"column:USER_ID" json:"user_id"`
	Date        time.Time `gorm:"column:DATE" json:"date"`
	Activity    string    `gorm:"column:ACTIVITY" json:"activity"`
	QuarterDate time.Time `gorm:"column:QUARTER_DATE" json:"quarter_date"`
	Approved    bool      `gorm:"column:APPROVED" json:"approved"`
	User        User      `gorm:"foreignkey:UserID" json:"user"`
	Client      Client    `gorm:"foreignkey:BPID;association_foreignkey:BP_ID" json:"client"`
}

// TableName for OtherInformation
func (BillingActivities) TableName() string {
	return "CRS_BILLING_ACTIVITIES"
}

// NPRAActivities structure of the sec activities table
type NPRAActivities struct {
	ID          int       `gorm:"column:ID;AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	BPID        string    `gorm:"column:BP_ID" json:"bpid"`
	Hash        string    `gorm:"column:HASH" json:"hash"`
	UserID      int       `gorm:"column:USER_ID" json:"user_id"`
	Date        time.Time `gorm:"column:DATE" json:"date"`
	Activity    string    `gorm:"column:ACTIVITY" json:"activity"`
	QuarterDate time.Time `gorm:"column:QUARTER_DATE" json:"quarter_date"`
	Approved    bool      `gorm:"column:APPROVED" json:"approved"`
	User        User      `gorm:"foreignkey:UserID" json:"user"`
	Client      Client    `gorm:"foreignkey:BPID;association_foreignkey:BP_ID" json:"client"`
}

// TableName for OtherInformation
func (NPRAActivities) TableName() string {
	return "CRS_NPRA_ACTIVITIES"
}

// TrusteePVPerformance structure of the data returned from the query
type TrusteePVPerformance struct {
	Bond            string  `json:"bond" gorm:"column:SECURITY_TYPE"`
	CurrentQuarter  float64 `json:"current_quarter" gorm:"column:CURRENT_QUARTER"`
	PreviousQuarter float64 `json:"previous_quarter" gorm:"column:PREVIOUS_QUARTER"`
}

// ClientMonthlyContribution structure of table
type ClientMonthlyContribution struct {
	ID        int       `gorm:"column:ID;AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	Quarter   time.Time `gorm:"column:QUARTER" json:"quarter"`
	BPID      string    `json:"bpid" gorm:"column:BP_ID"`
	SCA       string    `json:"sca" gorm:"column:SCA"`
	Date      time.Time `json:"date" gorm:"column:DATE"`
	Amount    float64   `json:"amount" gorm:"column:AMOUNT"`
	CreatedAt time.Time `json:"created_at" gorm:"column:CREATED_AT"`
}

// TableName for ClientMonthlyContribution
func (ClientMonthlyContribution) TableName() string {
	return "CRS_TRUSTEE_MONTHLY_CONTRIBUTIONS"
}

// FundManager structure of a fundmanager
type FundManager struct {
	ID            int       `gorm:"column:ID;AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	Name          string    `json:"name" gorm:"column:NAME"`
	AccountNumber string    `json:"account_number" gorm:"column:ACCOUNT_NUMBER"`
	Administrator string    `json:"administrator" gorm:"column:ADMINISTRATOR"`
	CreatedAt     time.Time `json:"created_at" gorm:"column:CREATED_AT"`
}

// TableName for FundManager
func (FundManager) TableName() string {
	return "CRS_FUNDMANAGERS"
}

// UnidentifiedPayment structure of an unidentified payment
type UnidentifiedPayment struct {
	ID                int         `gorm:"column:ID;AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	ClientBPID        string      `json:"client_bpid" gorm:"column:CLIENT_BPID"`
	TransactionType   string      `json:"txn_type" gorm:"column:TRANSACTION_TYPE"`
	TransactionDate   time.Time   `json:"txn_date" gorm:"column:TRANSACTION_DATE"`
	ValueDate         time.Time   `json:"value_date" gorm:"column:VALUE_DATE"`
	NameOfCompany     string      `json:"name_of_company" gorm:"column:NAME_OF_COMPANY"`
	Amount            float64     `json:"amount" gorm:"column:AMOUNT"`
	FundManagerID     int         `json:"-" gorm:"column:FUND_MANAGER_ID"`
	FundManager       FundManager `json:"fundManager" gorm:"foreignkey:FundManagerID"`
	CollectionAccount string      `json:"collection_acc_num" gorm:"column:COLLECTION_ACCOUNT_NO"`
	Status            string      `json:"status" gorm:"column:STATUS"`
	CreatedAt         time.Time   `json:"-" gorm:"column:CREATED_AT"`
	UpdatedAt         time.Time   `json:"-" gorm:"column:UPDATED_AT"`
}

// TableName for UnidentifiedPayment
func (UnidentifiedPayment) TableName() string {
	return "CRS_TRUSTEE_UNIDENTIFIED_PAYMENTS"
}

// UnidentifiedPaymentSummary structure of the summary
type UnidentifiedPaymentSummary struct {
	TransactionType string `json:"txn_type" gorm:"column:TRANSACTION_TYPE"`
	Total           int    `json:"total" gorm:"column:TOTAL"`
	Pending         int    `json:"pending" gorm:"column:PENDING"`
	Done            int    `json:"done" gorm:"column:DONE"`
}

// PVReportHeadings structure of the pv report heading in the db
type PVReportHeadings struct {
	Heading  string `gorm:"column:HEADING" json:"heading"`
	RealName string `gorm:"column:REAL_NAME" json:"real_name"`
}

// TableName for PVReportHeadings
func (PVReportHeadings) TableName() string {
	return "CRS_PVHEADINGS"
}

// NPRAFund structure of an npra fund
type NPRAFund struct {
	Name  string  `gorm:"column:name" json:"name"`
	Value float64 `gorm:"column:value" json:"value"`
	BPID  string  `gorm:"column:bpid" json:"bpid"`
}

// NPRATrend structure of the npra trend
type NPRATrend struct {
	Pension   float64
	Provident float64
}

type ValueQuery struct {
	Value float64 `gorm:"column:value"`
}

type StrValueQuery struct {
	Value string `gorm:"column:value"`
}

type IntValueQuery struct {
	Value int64 `gorm:"column:value"`
}

// NPRAUnauthorizedTransaction structure of unauthorized transactions
type NPRAUnauthorizedTransaction struct {
	ClientName         string    `gorm:"column:CLIENT_NAME" json:"clientName"`
	TransactionDetails string    `gorm:"column:TRANSACTION_DETAILS" json:"txnDetails"`
	Date               time.Time `gorm:"column:DATE" json:"date"`
}

func (NPRAUnauthorizedTransaction) TableName() string {
	return "CRS_NPRA_UNAUTHORIZED_TRANSACTIONS"
}

// NPRAOutstandingFDCertificate structure in DB
type NPRAOutstandingFDCertificate struct {
	ID            int       `gorm:"column:ID;AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	FundManager   string    `json:"fundManager" gorm:"column:FUND_MANAGER"`
	ClientName    string    `json:"clientName" gorm:"column:CLIENT_NAME"`
	Amount        float64   `json:"amount" gorm:"column:AMOUNT"`
	Issuer        string    `json:"issuer" gorm:"column:ISSUER"`
	Rate          float32   `json:"rate" gorm:"column:RATE"`
	Tenor         int       `json:"tenor" gorm:"column:TENOR"`
	Term          string    `json:"term" gorm:"column:TERM"`
	Hash          string    `json:"hash" gorm:"column:HASH"`
	EffectiveDate time.Time `json:"effectiveDate" gorm:"column:EFFECTIVE_DATE"`
	Maturity      time.Time `json:"maturity" gorm:"column:MATURITY"`
	CreatedAt     time.Time `gorm:"column:CREATED_AT" json:"created_at"`
}

func (NPRAOutstandingFDCertificate) TableName() string {
	return "CRS_NPRA_OUTSTANDING_FD_CERTIFICATES"
}

// BillingTransactionDetails ...
type BillingTransactionDetails struct {
	ID                   int       `gorm:"column:ID;AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	BPID                 string    `json:"bpid" gorm:"column:BP_ID"`
	SCA                  string    `json:"sca" gorm:"column:SCA"`
	ReportingDate        time.Time `json:"reporting_date" gorm:"column:REPORTING_DATE"`
	Date                 time.Time `gorm:"column:DATE" json:"date"`
	Reference            string    `gorm:"column:REFERENCE" json:"reference"`
	SecurityName         string    `gorm:"column:SECURITY_NAME" json:"security_name"`
	SecurityCategory     string    `gorm:"column:SECURITY_CATEGORY" json:"security_category"`
	ChargeType           string    `gorm:"column:CHARGE_TYPE" json:"charge_type"`
	ChargeItem           string    `gorm:"column:CHARGE_ITEM" json:"charge_item"`
	NumberOfUnits        int       `gorm:"column:NUMBER_OF_UNITS" json:"number_of_units"`
	MarketValue          float64   `gorm:"column:MARKET_VALUE" json:"market_value"`
	ChargeAmountWithTax  float64   `gorm:"column:CHARGE_AMOUNT_WITH_TAX" json:"charge_amount_with_tax"`
	InvoiceAmountWithTax float64   `gorm:"column:INVOICE_AMOUNT_WITH_TAX" json:"invoice_amount_with_tax"`
}

func (BillingTransactionDetails) TableName() string {
	return "CRS_BILLING_TRANSACTION_DETAILS"
}

// BillingTransactionDetails ...
type BillingCurrencyDetails struct {
	ID            int       `gorm:"column:ID;AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	BPID          string    `json:"bpid" gorm:"column:BP_ID"`
	SCA           string    `json:"sca" gorm:"column:SCA"`
	Currency      string    `gorm:"column:CURRENCY" json:"currency"`
	Rate          float32   `gorm:"column:RATE" json:"rate"`
	Date          time.Time `gorm:"column:DATE" json:"date"`
	ReportingDate time.Time `gorm:"column:REPORTING_DATE" json:"reporting_date"`
}

func (BillingCurrencyDetails) TableName() string {
	return "CRS_BILLING_CURRENCY_DETAILS"
}

type ClientBillingInfo struct {
	ID                   int     `gorm:"column:ID;AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	BPID                 string  `json:"bpid" gorm:"column:BP_ID"`
	MinimumCharge        float32 `gorm:"column:MINIMUM_CHARGE" json:"minimum_charge"`
	ChargePerTransaction float32 `gorm:"column:CHARGE_PER_TRANSACTION" json:"charge_per_transaction"`
	ThirdPartyTransfer   float32 `gorm:"column:THIRD_PARTY_TRANSFER" json:"third_party_transfer"`
}

func (ClientBillingInfo) TableName() string {
	return "CRS_CLIENT_BILLING_INFO"
}

type ClientBasisPoints struct {
	ID          int     `gorm:"column:ID;AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	BPID        string  `json:"bpid" gorm:"column:BP_ID"`
	SCA         string  `json:"sca" gorm:"column:SCA"`
	Minimum     float64 `gorm:"column:MINIMUM" json:"minimum"`
	Maximum     float64 `gorm:"column:MAXIMUM" json:"maximum"`
	BasisPoints float64 `gorm:"column:BASIS_POINTS" json:"basis_points"`
}

func (ClientBasisPoints) TableName() string {
	return "CRS_CLIENT_BASIS_POINTS"
}

type ClientInvoiceSummary struct {
	ChargeType             string  `json:"charge_type"`
	ChargeableQuantity     float64 `json:"chargeable_quantity"`
	ChargeAmount           float64 `json:"charge_amount"`
	TaxAmount              float64 `json:"tax_amount"`
	ChargeAmountWithTax    float64 `json:"charge_amount_with_tax"`
	InvoiceAmountWithTax   float64 `json:"invoice_amount_with_tax"`
	BasisPoint             float64 `json:"basis_point"`
	FivePercent            float64 `json:"five_percent"`
	TwelvePointFivePercent float64 `json:"twelve_point_five_percent"`
	InitialGross           float64 `json:"initial_gross"`
}

type CheckerDisapprovalComments struct {
	ID            int    `gorm:"column:ID;AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	ReportRefID   string `gorm:"column:REPORT_REF_ID" json:"report_ref_id"`
	DisapprovedBy int    `gorm:"column:DISAPPROVED_BY" json:"disapproved_by"`
	ReportType    string `gorm:"column:REPORT_TYPE" json:"report_type"`
	Comment       string `gorm:"column:COMMENT" json:"comment"`
}

func (CheckerDisapprovalComments) TableName() string {
	return "CRS_CHECKER_DISAPPROVAL_COMMENTS"
}

type SecReportInternalCompliance struct {
	Preparers  []User
	Reviewers  []User
	Authorizer User
}

type VarianceRemarks struct {
	ID         int       `gorm:"column:ID;AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	BPID       string    `json:"bpid" gorm:"column:BP_ID"`
	ReportDate time.Time `json:"report_date" gorm:"column:REPORT_DATE"`
	Remarks    string    `json:"remarks" gorm:"column:REMARKS"`
}

func (VarianceRemarks) TableName() string {
	return "CRS_VARIANCE_REMARKS"
}

type NPRADeclaration struct {
	ID                    int       `gorm:"column:ID;AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	ReportRefID           string    `json:"report_ref_id" gorm:"column:REPORT_REF_ID"`
	NameOfOfficer         string    `json:"nameOfOfficer" gorm:"column:NAME_OF_OFFICER"`
	Designation           string    `json:"designation" gorm:"column:DESIGNATION"`
	HeadOfCustodyServices string    `json:"headOfCustodyServices" gorm:"column:HEAD_OF_CUSTODY_SERVICES"`
	Date                  time.Time `json:"date" gorm:"column:DATE"`
	Approved              bool      `json:"approved" gorm:"column:APPROVED"`
	ApprovedBy            int       `json:"approved_by" gorm:"column:APPROVED_BY"`
}

func (NPRADeclaration) TableName() string {
	return "CRS_NPRA_DECLARATION"
}

type BillingNav struct {
	ID               int     `gorm:"column:ID;AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	RefID            string  `json:"ref_id" gorm:"column:REF_ID"`
	BPID             string  `json:"bpid" gorm:"column:BP_ID"`
	SCA              string  `json:"sca" gorm:"column:SCA"`
	Position         float64 `json:"position" gorm:"column:POSITION"`
	CashBalance      float64 `json:"cash_balance" gorm:"column:CASH_BALANCE"`
	Liabilities      float64 `json:"liabilities" gorm:"column:LIABILITIES"`
	NAV              float64 `json:"nav" gorm:"column:NAV"`
	InvoiceReference int64   `json:"invoice_ref" gorm:"column:INVOICE_REFERENCE"`
}

func (BillingNav) TableName() string {
	return "CRS_BILLING_NAV"
}

type LumpedMonthlyContribution struct {
	Date          string  `gorm:"column:date"`
	Contributions float64 `gorm:"column:contributions"`
}

type ClientPVUniqueCode struct {
	ID            int        `gorm:"column:ID;AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	BPID          string     `json:"bpid" gorm:"column:BP_ID"`
	ClientName    string     `gorm:"column:CLIENT_NAME" json:"client_name"`
	ClientTier    int        `gorm:"column:CLIENT_TIER" json:"client_tier"`
	Code          string     `gorm:"column:CODE" json:"code"`
	FundManager   string     `gorm:"column:FUND_MANAGER" json:"fund_manager"`
	Administrator string     `gorm:"column:ADMINISTRATOR" json:"administrator"`
	GroupKey      string     `gorm:"column:GROUP_KEY" json:"group_key"`
	Closed        *time.Time `gorm:"column:CLOSED" json:"closed"`
}

func (ClientPVUniqueCode) TableName() string {
	return "CRS_CLIENT_PV_UNIQUE_CODES"
}

type ClientMergedBPAndSCA struct {
	Client         string
	SCA            []string
	BPID           string
	HasMultipleSCA bool
}

type ClientIndividualPVSummary struct {
	Client      string
	HTMLTableID int
	Summary     []ClientPVReportSummary
}

type PVSummaryTotal struct {
	Nominal    float64 `gorm:"column:nominal"`
	LCY        float64 `gorm:"column:lcy"`
	Percentage float32 `gorm:"column:percentage"`
}

type SCAMonthlyContribution struct {
	Client        string
	HTMLTableID   int
	Contributions []LumpedMonthlyContribution
}

type GOGMaturity struct {
	ID             int       `gorm:"column:ID;AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	BPID           string    `json:"bpid" gorm:"column:BP_ID"`
	EntryDate      time.Time `json:"entry_date" gorm:"column:ENTRY_DATE"`
	EventType      string    `json:"event_type" gorm:"column:EVENT_TYPE"`
	DepotID        string    `json:"depot_id" gorm:"column:DEPOT_ID"`
	Status         string    `json:"status" gorm:"column:STATUS"`
	GrossAmount    float64   `json:"gross_amount" gorm:"column:GROSS_AMOUNT"`
	BaseSecurityID string    `json:"base_security_id" gorm:"column:BASE_SECURITY_ID"`
	Hash           string    `json:"hash" gorm:"column:HASH"`
	QuarterDate    time.Time `json:"quarter_date" gorm:"column:QUARTER_DATE"`
}

func (GOGMaturity) TableName() string {
	return "CRS_TRUSTEE_GOG_FD"
}

type GOGSummary struct {
	Month              string  `gorm:"column:month"`
	NumberOfMaturities int     `gorm:"column:number_of_maturities"`
	Amount             float64 `gorm:"column:amount"`
}

type CorporateActionActivity struct {
	Activity string `gorm:"column:activity"`
	Number   int    `gorm:"column:number"`
}

type UnidentifiedPaymentsSummary struct {
	Type    string `gorm:"column:type"`
	Total   int    `gorm:"column:total"`
	Done    int    `gorm:"column:done"`
	Pending int    `gorm:"column:pending"`
}

type TradeVolumeByAssetClass struct {
	Asset  string
	Number int64
}

type TradeVolumeByAssetClassSummary struct {
	Month string
	Data  []TradeVolumeByAssetClass
}

type TrusteeUploadedPVs struct {
	Client      string    `json:"client" gorm:"column:CLIENT"`
	BPID        string    `json:"bpid" gorm:"column:BP_ID"`
	UploadedBy  string    `json:"uploaded_by" gorm:"column:FULLNAME"`
	Date        time.Time `gorm:"column:DATE" json:"date"`
	Type        string    `gorm:"column:TYPE" json:"type"`
	QuarterDate time.Time `gorm:"column:QUARTER_DATE" json:"quarter_date"`
}

type TrusteeQuarterlyReport struct {
	ID         int    `gorm:"column:ID;AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	BPID       string `json:"bpid" gorm:"column:BP_ID"`
	Quarter    string `json:"quarter" gorm:"column:QUARTER"`
	Approved   bool   `json:"approved" gorm:"column:APPROVED"`
	ApprovedBy *int   `json:"approved_by" gorm:"column:APPROVED_BY"`
}

func (TrusteeQuarterlyReport) TableName() string {
	return "CRS_TRUSTEE_QUARTERLY_REPORT"
}

type BillingQuarterlyReport struct {
	ID         int    `gorm:"column:ID;AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	ClientID   string `gorm:"column:CLIENT_ID" json:"client_id"`
	Month      string `json:"month" gorm:"column:MONTH"`
	Approved   bool   `json:"approved" gorm:"column:APPROVED"`
	ApprovedBy *int   `json:"approved_by" gorm:"column:APPROVED_BY"`
}

func (BillingQuarterlyReport) TableName() string {
	return "CRS_BILLING_QUARTERLY_REPORT"
}

// TrusteeActivities structure of the trustee activities table
type TrusteeActivities struct {
	ID          int                `gorm:"column:ID;AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	BPID        string             `json:"bpid" gorm:"column:BP_ID"`
	Hash        string             `gorm:"column:HASH" json:"hash"`
	UserID      int                `gorm:"column:USER_ID" json:"user_id"`
	Date        time.Time          `gorm:"column:DATE" json:"date"`
	Activity    string             `gorm:"column:ACTIVITY" json:"activity"`
	QuarterDate time.Time          `gorm:"column:QUARTER_DATE" json:"quarter_date"`
	Approved    bool               `json:"approved" gorm:"column:APPROVED"`
	User        User               `gorm:"foreignkey:UserID" json:"user"`
	Client      Client             `gorm:"foreignkey:BPID;association_foreignkey:BP_ID" json:"client"`
	SCAClient   ClientPVUniqueCode `gorm:"foreignkey:BPID;association_foreignkey:CODE" json:"sca_client"`
}

func (TrusteeActivities) TableName() string {
	return "CRS_TRUSTEE_ACTIVITIES"
}

type ApprovedReportReversal struct {
	ID                   int       `gorm:"column:ID;AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	Type                 string    `json:"type" gorm:"column:TYPE"`
	ClientID             string    `json:"client_id" gorm:"column:CLIENT_ID"`
	ReportDate           string    `json:"report_date" gorm:"column:REPORT_DATE"`
	ReversedBy           int       `json:"reversed_by" gorm:"column:REVERSED_BY"`
	ReversedOn           time.Time `json:"reversed_on" gorm:"column:REVERSED_ON"`
	PreviouslyApprovedBy int       `json:"previously_approved_by" gorm:"column:PREVIOUSLY_APPROVED_BY"`
}

func (ApprovedReportReversal) TableName() string {
	return "CRS_APPROVED_REPORT_REVERSAL"
}

type NAVUpdateLog struct {
	ID             int     `gorm:"column:ID;AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	ClientID       string  `json:"client_id" gorm:"column:CLIENT_ID"`
	UserID         int     `json:"user_id" gorm:"column:USER_ID"`
	ReportingMonth string  `json:"reporting_month" gorm:"column:REPORTING_MONTH"`
	Position       float64 `json:"position" gorm:"column:POSITION"`
	Liabilities    float64 `json:"liabilities" gorm:"column:LIABILITIES"`
	NAV            float64 `json:"nav" gorm:"column:NAV"`
}

func (NAVUpdateLog) TableName() string {
	return "CRS_BILLING_NAV_UPDATES"
}

type BillingTransactionJournal struct {
	ID                            int       `gorm:"column:ID;AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	InvoiceDate                   time.Time `gorm:"column:INVOICE_DATE" json:"invoice_date"`
	InvoiceReference              int64     `gorm:"column:INVOICE_REFERENCE" json:"invoice_reference"`
	InvoiceAmount                 float64   `gorm:"column:INVOICE_AMOUNT" json:"invoice_amount"`
	Position                      float64   `gorm:"column:POSITION" json:"position"`
	Liabilities                   float64   `gorm:"column:LIABILITIES" json:"liabilities"`
	NAV                           float64   `gorm:"column:NAV" json:"nav"`
	ClientID                      string    `gorm:"column:CLIENT_ID" json:"client_id"`
	InvoiceDueDate                time.Time `gorm:"column:INVOICE_DUE_DATE" json:"invoice_due_date"`
	ApprovedOn                    time.Time `gorm:"column:APPROVED_ON" json:"approved_on"`
	ApprovedBy                    int       `gorm:"column:APPROVED_BY" json:"approved_by"`
	PortfolioChargeableQuantity   float64   `gorm:"column:PORTFOLIO_CHARGEABLE_QUANTITY" json:"portfolio_chargeable_quantity"`
	TransactionChargeableQuantity int       `gorm:"column:TRANSACTION_CHARGEABLE_QUANTITY" json:"transaction_chargeable_quantity"`
	BasisPoint                    float64   `gorm:"column:BASIS_POINT" json:"basis_point"`
	ChargePerTransaction          float32   `gorm:"column:CHARGE_PER_TRANSACTION" json:"charge_per_transaction"`
	PortfolioChargeAmount         float64   `gorm:"column:PORTFOLIO_CHARGE_AMOUNT" json:"portfolio_charge_amount"`
	TransactionChargeAmount       float64   `gorm:"column:TRANSACTION_CHARGE_AMOUNT" json:"transaction_charge_amount"`
}

func (BillingTransactionJournal) TableName() string {
	return "CRS_BILLING_TRANSACTION_JOURNAL"
}

type FundAdministrators struct {
	Name          string `json:"name" gorm:"column:NAME"`
	Administrator string `json:"administrator" gorm:"column:ADMINISTRATOR"`
}

type ClientDetailedBillingDetails struct {
	ClientName           string
	MinimumCharge        float32
	ChargePerTransaction float32
	ThirdPartyTransfer   float32
	BasisPoints          []ClientBasisPoints
}

type AuditablePV struct {
	Client  ClientPVUniqueCode      `json:"client"`
	Reports []PVReportField         `json:"reports"`
	Summary []ClientPVReportSummary `json:"summary"`
}

type ParsedPVReport struct {
	Headers []string     `json:"headers"`
	Bonds   []PVBondData `json:"bonds"`
	Summary [][]string   `json:"summary"`
}

type PVBondData struct {
	Bond   string     `json:"bond"`
	Values [][]string `json:"values"`
}

// PVReportData is the structure of a block of NAV report data
type ReportData struct {
	Title  string
	Values []ReportField
}

// PVReportField the fields in the NAV report
type ReportField struct {
	SecurityName      string
	CDSCode           string
	ISIN              string
	SCBCode           string
	MarketPrice       float32
	NominalValue      float64
	CumulativeCost    float64
	Value             float64
	PercentageOfTotal float32
	Dates             ReportDate
}

// PVReportDate the date fields in the PVField
type ReportDate struct {
	From string
	To   string
}

// IsEmpty checks if the date struct is empty
func (date ReportDate) IsEmpty() bool {
	return date.From == "" && date.To == ""
}

// PVClient the client who owns a pv document
type PVClient struct {
	UserID      int
	ClientID    string
	Date        string
	CashBalance float64
	Type        string
	Report      []ReportData
	Summary     []ReportSummary
	Error       []PVUploadError
	RawHeading  string
}

// HasNoDate checks if date is set
func (client PVClient) HasNoDate() bool {
	return client.Date == ""
}

// ReportSummary the report summary
type ReportSummary struct {
	SecurityName      string
	NominalValue      float64
	CumulativeCost    float64
	Value             float64
	PercentageOfTotal float32
}

// ClientDetailsQuery structure of the query to fetch the bpid
type ClientDetailsQuery struct {
	BPID               string `gorm:"column:BP_ID"`
	SafekeepingAccount string `gorm:"column:SAFEKEEPINGACCNO"`
	Name               string `gorm:"column:CLIENTNAME"`
}

type PVUploadError struct {
	Type       string `json:"type"`
	SCA        string `json:"sca"`
	ClientInfo string `json:"client_info"`
	Date       string `json:"date"`
	MoreInfo   string `json:"more_info"`
}

func (err PVUploadError) IsNil() bool {
	return err.Type == "" && err.SCA == "" && err.ClientInfo == "" && err.Date == ""
}

type BillingReportExportPayload struct {
	Client             Client                      `json:"client"`
	Period             string                      `json:"period"`
	InvoiceReference   int64                       `json:"invoice_reference"`
	InvoiceAmount      float64                     `json:"invoice_amount"`
	CurrencyDetails    []BillingCurrencyDetails    `json:"currency_details"`
	Summary            []ClientInvoiceSummary      `json:"summary"`
	TransactionDetails []BillingTransactionDetails `json:"transaction_details"`
}

type Director struct {
	ID                 int    `gorm:"column:ID;AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	Fullname           string `gorm:"column:FULLNAME" json:"fullname"`
	ResidentialAddress string `gorm:"column:RESIDENTIAL_ADDRESS" json:"residential_address"`
	Nationality        string `gorm:"column:NATIONALITY" json:"nationality"`
	Position           string `gorm:"column:POSITION" json:"position"`
}

func (Director) TableName() string {
	return "CRS_DIRECTORS"
}

type DirectorsExportPayload struct {
	Directors []Director `json:"directors"`
}

type UserPasswordResetPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type MaturedSecurity struct {
	ID             int     `gorm:"column:ID;AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	Client         string  `gorm:"column:CLIENT" json:"client"`
	Issuer         string  `gorm:"column:ISSUER" json:"issuer"`
	AmountInvested float64 `gorm:"column:AMOUNT_INVESTED" json:"amount_invested"`
	Value          float64 `gorm:"column:VALUE" json:"value"`
	Date           string  `gorm:"column:DATE" json:"date"`
	Approved       bool    `gorm:"column:APPROVED" json:"approved"`
}

func (MaturedSecurity) TableName() string {
	return "CRS_MATURED_SECURITIES"
}

type ClientEmail struct {
	ID    int    `gorm:"column:ID;AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	BPID  string `gorm:"column:BPID" json:"bpid"`
	Email string `gorm:"column:EMAIL" json:"email"`
}

func (ClientEmail) TableName() string {
	return "CRS_CLIENT_EMAILS"
}

type MailService struct {
	ID    int    `gorm:"column:ID;AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	Email string `gorm:"column:EMAIL" json:"email"`
}

func (MailService) TableName() string {
	return "CRS_MAIL_SERVICE"
}

type BilledClients struct {
	Client           string `gorm:"client"`
	Bpid             string `gorm:"bpid"`
	InvoiceReference int64  `gorm:"invoice_reference"`
}

func (billedClient BilledClients) HasInvoiceRef() bool {
	return billedClient.InvoiceReference > 0
}

type MailSenderInfo struct {
	ID       int    `gorm:"column:ID;AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	Name     string `gorm:"column:NAME" json:"name"`
	Position string `gorm:"column:POSITION" json:"position"`
}

func (MailSenderInfo) TableName() string {
	return "CRS_MAIL_SENDER_INFO"
}

type Holiday struct {
	ID   int    `gorm:"column:ID;AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	Name string `gorm:"column:NAME" json:"name"`
	Date string `gorm:"column:DATE" json:"date"`
}

func (Holiday) TableName() string {
	return "CRS_HOLIDAYS"
}

type Securities struct {
	ID     int    `gorm:"column:ID;AUTO_INCREMENT;PRIMARY_KEY" json:"id" form:"id"`
	Name   string `gorm:"column:NAME" json:"name" form:"name"`
	Plural string `gorm:"column:PLURAL" json:"plural" form:"plural"`
}

func (Securities) TableName() string {
	return "CRS_SECURITIES_NAME"
}

type AssetClass struct {
	ID             int    `gorm:"column:ID;AUTO_INCREMENT;PRIMARY_KEY" json:"id" form:"id"`
	SecurityType   string `gorm:"column:SECURITY_TYPE" json:"security_type" form:"security_type"`
	AssetClassName string `gorm:"column:ASSET_CLASS_NAME" json:"asset_class_name" form:"asset_class_name"`
}

func (AssetClass) TableName() string {
	return "CRS_ASSET_CLASS"
}

type MergedSecurityName struct {
	Security string `gorm:"column:securities"`
}

type QuarterEmailToClients struct {
	ID          int       `gorm:"column:ID;AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	QuarterYear string    `gorm:"column:QUARTER_YEAR" json:"quarter_year"`
	SentBy      string    `gorm:"column:SENT_BY" json:"sent_by"`
	SentOn      time.Time `gorm:"column:SENT_ON" json:"sent_on"`
	CreatedAt   time.Time `gorm:"column:CREATED_AT" json:"created_at"`
}

func (QuarterEmailToClients) TableName() string {
	return "CRS_QUARTER_EMAIL_TO_CLIENTS"
}

// Client db table
type NPRA0303 struct {
	ID         int    `gorm:"column:ID;AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	ReportCode string `gorm:"column:report_code" json:"report_code"`
	//BPID                string       `gorm:"column:bp_id" json:"bp_id"`
	EntityID            string       `gorm:"column:entity_id" json:"entity_id"`
	EntityName          string       `gorm:"column:entity_name" json:"entity_name"`
	ReferencePeriodYear string       `gorm:"column:reference_period_year" json:"reference_period_year"`
	ReferencePeriod     string       `gorm:"column:reference_period" json:"reference_period"`
	DateValuation       time.Time    `gorm:"column:date_valuation" json:"date_valuation"`
	UnitNumber          string       `gorm:"column:unit_number" json:"unit_number"`
	DailyNav            string       `gorm:"column:daily_nav" json:"daily_nav"`
	UnitPrice           float64      `gorm:"column:unit_price" json:"unit_price"`
	NPRAFees            float64      `gorm:"column:npra_fees" json:"npra_fees"`
	TrusteeFees         float64      `gorm:"column:trustee_fees" json:"trustee_fees"`
	FundManagerFees     float64      `gorm:"column:fund_manager_fees" json:"fund_manager_fees"`
	FundCustodianFees   float64      `gorm:"column:fund_custodian_fees" json:"fund_custodian_fees"`
	CreatedBy           nulls.String `gorm:"column:created_by" json:"created_by"`
}

// TableName sets the role table name for gorm
func (NPRA0303) TableName() string {
	return "CRS_0303_NPRA_REPORT"
}

// Client db table
//
//	type NPRA0302 struct {
//		ID                                           int        `gorm:"column:ID;AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
//		BP_ID                                        string     `gorm:"column:bp_id" json:"bp_id"`
//		ReportCode                                   string     `gorm:"column:report_code" json:"report_code"`
//		EntityID                                     string     `gorm:"column:entity_id" json:"entity_id"`
//		EntityName                                   string     `gorm:"column:entity_name" json:"entity_name"`
//		ReferencePeriodYear                          string     `gorm:"column:reference_period_year" json:"reference_period_year"`
//		ReferencePeriod                              string     `gorm:"column:reference_period" json:"reference_period"`
//		InvestmentID                                 string     `gorm:"column:investment_id" json:"investment_id"`
//		Instrument                                   string     `gorm:"column:instrument" json:"instrument"`
//		IssuerName                                   string     `gorm:"column:issuer_name" json:"issuer_name"`
//		AssetTenure                                  string     `gorm:"column:asset_tenure" json:"asset_tenure"`
//		DateInvestment                               nulls.Time `gorm:"column:date_investment" json:"date_investment"`
//		ReportingDate                                time.Time  `gorm:"column:reporting_date" json:"reporting_date"`
//		AmountInvested                               float64    `gorm:"column:amount_invested" json:"amount_invested"`
//		AccruedCouponInterest                        float64    `gorm:"column:accrued_coupon_interest" json:"accrued_coupon_interest"`
//		CouponPaid                                   float64    `gorm:"column:coupon_paid" json:"coupon_paid"`
//		AccruedCouponInterestSinceLastYear           float64    `gorm:"column:accrued_coupon_interest_since_last_year" json:"accrued_coupon_interest_since_last_year"`
//		OutstandingInterestMaturity                  float64    `gorm:"column:outstanding_interest_maturity" json:"outstanding_interest_maturity"`
//		AmountImpaired                               float64    `gorm:"column:amount_impaired" json:"amount_impaired"`
//		AssetAllocationActualPercent                 float64    `gorm:"column:asset_allocation_actual_percent" json:"asset_allocation_actual_percent"`
//		MaturityDate                                 float64    `gorm:"column:maturity_date" json:"maturity_date"`
//		TypeInvestmentCharge                         float64    `gorm:"column:type_investment_charge" json:"type_investment_charge"`
//		InvestmentChargeDatePercent                  float64    `gorm:"column:investment_charge_date_percent" json:"investment_charge_date_percent"`
//		IvestmentChargeAmount                        float64    `gorm:"column:investment_charge_amount" json:"investment_charge_amount"`
//		FaceValue                                    float64    `gorm:"column:face_value" json:"face_value"`
//		InterestRatePercent                          float64    `gorm:"column:interest_rate_percent" json:"interest_rate_percent"`
//		DiscountRatePercent                          float64    `gorm:"column:discount_rate_percent" json:"discount_rate_percent"`
//		CouponRatePercent                            float64    `gorm:"column:coupon_rate_percent" json:"coupon_rate_percent"`
//		DisposalProceeds                             float64    `gorm:"column:disposal_proceeds" json:"disposal_proceeds"`
//		DisposalInstructions                         string     `gorm:"column:disposal_instructions" json:"disposal_instructions"`
//		YieldDisposal                                float64    `gorm:"column:yield_disposal" json:"yield_disposal"`
//		IssueDate                                    time.Time  `gorm:"column:issue_date" json:"issue_date"`
//		PriceShareUnitPurchase                       float64    `gorm:"column:price_share_unit_purchase" json:"price_share_unit_purchase"`
//		PriceShareUnitValueDate                      time.Time  `gorm:"column:price_share_unit_value_date" json:"price_share_unit_value_date"`
//		CapitalGains                                 float64    `gorm:"column:capital_gains" json:"capital_gains"`
//		DividendReceived                             float64    `gorm:"column:dividend_received" json:"dividend_received"`
//		NumberUnitsShares                            int64      `gorm:"column:number_units_shares" json:"number_units_shares"`
//		HoldingPeriodReturnPerInvestment             float64    `gorm:"column:holding_period_return_per_investment" json:"holding_period_return_per_investment"`
//		DaysRun                                      int64      `gorm:"column:days_run" json:"days_run"`
//		CurrencyConversionRate                       float64    `gorm:"column:currency_conversion_rate" json:"currency_conversion_rate"`
//		Currency                                     string     `gorm:"column:currency" json:"currency"`
//		AmountInvestedForeignCurrency                float64    `gorm:"column:amount_invested_foreign_currency" json:"amount_invested_foreign_currency"`
//		AssetClass                                   string     `gorm:"column:asset_class" json:"asset_class"`
//		PriceShareUnitValueDateForeign               time.Time  `gorm:"column:price_share_unit_value_date_foreign" json:"price_share_unit_value_date_foreign"`
//		MarketValue                                  float64    `gorm:"column:market_value" json:"market_value"`
//		RemainingDaysMaturity                        int64      `gorm:"column:remaining_days_maturity" json:"remaining_days_maturity"`
//		HoldingPeriodPerInvestmentWeightedPercentage float64    `gorm:"column:holding_period_per_investment_weighted_percentage" json:"holding_period_per_investment_weighted_percentage"`
//		CreatedBy                                    string     `gorm:"column:created_by" json:"created_by"`
//	}
//
// Client db table
type NPRA0302 struct {
	ID                  int                `gorm:"column:ID;AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	BP_ID               string             `gorm:"column:bp_id" json:"bp_id"`
	ClientCode          string             `gorm:"column:client_code" json:"client_code"`
	ReportCode          string             `gorm:"column:report_code" json:"report_code"`
	EntityID            string             `gorm:"column:entity_id" json:"entity_id"`
	EntityName          string             `gorm:"column:entity_name" json:"entity_name"`
	ReferencePeriodYear string             `gorm:"column:reference_period_year" json:"reference_period_year"`
	ReferencePeriod     string             `gorm:"column:reference_period" json:"reference_period"`
	InvestmentID        string             `gorm:"column:investment_id" json:"investment_id"`
	Instrument          string             `gorm:"column:instrument" json:"instrument"`
	IssuerName          string             `gorm:"column:issuer_name" json:"issuer_name"`
	AssetTenure         string             `gorm:"column:asset_tenure" json:"asset_tenure"`
	ReportingDate       time.Time          `gorm:"column:reporting_date" json:"reporting_date"`
	MaturityDate        time.Time          `gorm:"column:maturity_date" json:"maturity_date"`
	FaceValue           float64            `gorm:"column:face_value" json:"face_value"`
	IssueDate           time.Time          `gorm:"column:issue_date" json:"issue_date"`
	Currency            string             `gorm:"column:currency" json:"currency"`
	AssetClass          string             `gorm:"column:asset_class" json:"asset_class"`
	MarketValue         float64            `gorm:"column:market_value" json:"market_value"`
	ClientSCADetails    ClientPVUniqueCode `gorm:"foreignkey:ClientCode;association_foreignkey:CODE" json:"client"`
	//CreatedBy           time.Time `gorm:"column:created_by" json:"created_by"`
}

// TableName sets the role table name for gorm
func (NPRA0302) TableName() string {
	return "CRS_0302_NPRA_REPORT"
}

type Auditable0302 struct {
	Client     ClientPVUniqueCode `json:"client"`
	Report0302 []NPRA0302         `json:"report0302"`
}
type Auditable0301 struct {
	Report0301 []NPRA0301 `json:"report0301"`
}

// Client db table
type NPRA0301 struct {
	ID                         int     `gorm:"column:ID;AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	BP_ID                      string  `gorm:"column:bp_id" json:"bp_id"`
	ReportCode                 string  `gorm:"column:report_code" json:"report_code"`
	EntityID                   string  `gorm:"column:entity_id" json:"entity_id"`
	EntityName                 string  `gorm:"column:entity_name" json:"entity_name"`
	ReferencePeriodYear        string  `gorm:"column:reference_period_year" json:"reference_period_year"`
	ReferencePeriod            string  `gorm:"column:reference_period" json:"reference_period"`
	NetReturn                  float64 `gorm:"column:net_return" json:"net_return"`
	GrossReturn                float64 `gorm:"column:gross_return" json:"gross_return"`
	InvestmentReceivables      float64 `gorm:"column:investment_receivables" json:"investment_receivables"`
	TotalAssetUnderManagement  float64 `gorm:"column:total_asset_under_management" json:"total_asset_under_management"`
	GovernmentSecurities       float64 `gorm:"column:government_securities" json:"government_securities"`
	LocalGovernmentSecurities  float64 `gorm:"column:local_government_securities" json:"local_government_securities"`
	CorporateDebtSecurities    float64 `gorm:"column:corporate_debt_securities" json:"corporate_debt_securities"`
	BankSecurities             float64 `gorm:"column:bank_securities" json:"bank_securities"`
	OrdinaryPreferenceShares   int     `gorm:"column:ordinary_preference_shares" json:"ordinary_preference_shares"`
	CollectiveInvestmentScheme int     `gorm:"column:collective_investment_scheme" json:"collective_investment_scheme"`
	AlternativeInvestments     int     `gorm:"column:alternative_investments" json:"alternative_investments"`
	BankBalances               float64 `gorm:"column:bank_balances" json:"bank_balances"`
	//CreatedBy                  string    `gorm:"column:created_by" json:"created_by"`
	ReportingDate time.Time `gorm:"column:reporting_date" json:"reporting_date"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
}

// TableName sets the role table name for gorm
func (NPRA0301) TableName() string {
	return "CRS_0301_NPRA_REPORT"
}
