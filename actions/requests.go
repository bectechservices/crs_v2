package actions

import (
	"time"

	"github.com/gobuffalo/nulls"
	v "github.com/gobuffalo/validate"
)

// LoginRequest structure of a login request
type LoginRequest struct {
	StaffID  int    `form:"staff_id"`
	Password string `form:"password"`
}

// IsValid validates the login request input
func (loginRequest *LoginRequest) IsValid(errors *v.Errors) {
	if loginRequest.StaffID == 0 {
		errors.Add("staff_id", "Staff ID must not be blank!")
	}
	if loginRequest.Password == "" {
		errors.Add("password", "Password must not be blank!")
	}
}

// PVUploadRequest the structure of a file upload request
type PVUploadRequest struct {
	PVType      string           `json:"report_type"`
	CashBalance float64          `json:"cashbalance"`
	Data        []ParsedPVReport `json:"data"`
}

// GovernanceShare the structure of a governance share
type GovernanceShare struct {
	Name         string  `json:"name"`
	Shareholding int     `json:"shareholding"`
	Percentage   float32 `json:"percentage"`
}

// GovernanceCustodianTransaction structure of a governance txn
type GovernanceCustodianTransaction struct {
	NameOfTrustee           string  `json:"nameOfTrustee"`
	RelationshipWIthTrustee string  `json:"relationshipWithTrustee"`
	TypeOfTransaction       string  `json:"typeOfTransaction"`
	Amount                  float64 `json:"amount"`
}

// GovernanceSchemesUnderCustody structure of schemesundecustody
type GovernanceSchemesUnderCustody struct {
	NameOfFirm              string  `json:"nameOfFirm"`
	NameOfScheme            string  `json:"nameOfScheme"`
	RelationshipWIthTrustee string  `json:"relationshipWithTrustee"`
	Volume                  float64 `json:"volume"`
	MarkedToMarketValue     float32 `json:"markedToMarketValue"`
}

// GovernanceDataUploadRequest the structure of the governance upload request
type GovernanceDataUploadRequest struct {
	ClientName                         string                           `json:"clientName"`
	ReportingOfficer                   string                           `json:"reportingOfficer"`
	ReportingDate                      string                           `json:"reportingDate"`
	ChangeInDirectors                  string                           `json:"changeInDirectors"`
	ChangeInAgreement                  string                           `json:"changeInAgreement"`
	DealingsApprovedByBoard            string                           `json:"dealingsApprovedByBoard"`
	CustodianHasUpdatedAssetRegister   string                           `json:"custodianHasUpdatedAssetRegister"`
	CustodianAssetRegistrationDate     string                           `json:"custodianAssetRegistrationDate"`
	DoManagersOfTheSchemeConsultTheLaw string                           `json:"doManagersOfTheSchemeConsultTheLaw"`
	SchemeHadAnyOtherFinancialDealings string                           `json:"schemeHadAnyOtherFinancialDealings"`
	OrdinaryShares                     []GovernanceShare                `json:"ordinaryShares"`
	PreferenceShares                   []GovernanceShare                `json:"preferenceShares"`
	CustodianTransactions              []GovernanceCustodianTransaction `json:"custodianTransactions"`
	SchemesUnderCustody                []GovernanceSchemesUnderCustody  `json:"schemesUnderCustody"`
}

// SchemeDetailsFetchRequest structure of the request to fetch the scheme details
type SchemeDetailsFetchRequest struct {
	BPID        string `json:"bpid"`
	QuarterDate string `json:"quarterDate"`
}

// Details0301FetchRequest structure of the request to fetch the scheme details
type Details0301FetchRequest struct {
	QuarterDate string `json:"quarterDate"`
}

// UserAddRequest structure of the request to add a new user
type UserAddRequest struct {
	Fullname string `form:"fullname"`
	StaffID  int    `form:"staff_id"`
	Email    string `form:"email"`
	Roles    []int  `form:"roles"`
	Active   bool   `form:"active"`
}

// OnlyUserIDRequest the structure of a request with only the user id
type OnlyUserIDRequest struct {
	UserID int    `form:"user_id"`
	Email  string `form:"email"`
}

// EditUserRequest structure of the request to edit a user
type EditUserRequest struct {
	UserID   int    `form:"user_id"`
	Fullname string `form:"fullname"`
	StaffID  int    `form:"staff_id"`
	Email    string `form:"email"`
	Roles    []int  `form:"roles"`
	Active   bool   `form:"active"`
}

// OffshoreClientsRequest structure of the request to load offshore clients
type OffshoreClientsRequest struct {
	BPID    string `json:"bpid"`
	Quarter string `json:"quarter"`
	Year    string `json:"year"`
}

// QuarterAndYear quarter and year
type QuarterAndYear struct {
	Quarter string `json:"quarter"`
	Year    string `json:"year"`
}

// MonthAndYear quarter and year
type MonthAndYear struct {
	Month string `json:"month"`
	Year  string `json:"year"`
}

// MonthAndYear quarter and year
type BPAndMonthAndYear struct {
	BpOrSca string `json:"bpOrSca"`
	Month   string `json:"month"`
	Year    string `json:"year"`
}

// BPMonthAndYear quarter and year
type BPMonthAndYear struct {
	BPID  string `json:"bpOrSca"`
	Month string `json:"month"`
	Year  string `json:"year"`
}

func (request QuarterAndYear) IsEmpty() bool {
	return request.Quarter == "" && request.Year == ""
}

// PVFetchForEyeBallingRequest Fetch PV only
type PVFetchForEyeBallingRequest struct {
	SchemeDetailsFetchRequest
	PVType string `json:"pv_type"`
}

// TrusteePVDataRequest structure of the request to load the trustee pv data
type TrusteePVDataRequest struct {
	BPID    string `json:"bpid"`
	Quarter string `json:"quarter"`
	Year    string `json:"year"`
}

// ClientMonthlyContributionRequest Client Monthly Contributions
type ClientMonthlyContributionRequest struct {
	BPID    string  `json:"bpid"`
	Date    string  `json:"date"`
	Quarter string  `json:"quarter"`
	Year    string  `json:"year"`
	Amount  float64 `json:"amount"`
	SCA     string  `json:"sca"`
}

// ClientMonthlyContributionEditRequest Client Monthly Contributions
type ClientMonthlyContributionEditRequest struct {
	ID     int     `json:"id"`
	BPID   string  `json:"bpid"`
	Date   string  `json:"date"`
	Amount float64 `json:"amount"`
	SCA    string  `json:"sca"`
}

// ClientAddRequest the request for adding a client
type ClientAddRequest struct {
	ClientName          string `form:"clientName"`
	SafekeepingAccount  string `form:"safekeepingAccount"`
	BPID                string `form:"bpid"`
	HomeCountry         string `form:"homeCountry"`
	Email               string `form:"email"`
	Location            string `form:"location"`
	PhoneNumber         string `form:"phoneNumber"`
	ClientType          string `form:"clientType"`
	ClientTier          string `form:"clientTier"`
	FundManagerID       int    `form:"fundManager"`
	ContactPerson       string `form:"contactPerson"`
	ContactPersonNumber string `form:"contactPersonNumber"`
	AddressLine1        string `form:"address1"`
	AddressLine2        string `form:"address2"`
	AddressLine3        string `form:"address3"`
	AddressLine4        string `form:"address4"`
	Image               string
	SchemeManager       string `form:"scheme_manager"`
}

// ClientEditRequest the request to edit a client
type ClientEditRequest struct {
	ClientAddRequest
	ID int `form:"id"`
}

// FundManagerAddRequest the request to add a fund manager
type FundManagerAddRequest struct {
	Name          string `form:"name"`
	AccountNumber string `form:"accountNumber"`
	Administrator string `form:"administrator"`
}

// FundManagerEditRequest the request to edit a fund manager
type FundManagerEditRequest struct {
	FundManagerAddRequest
	ID int `form:"id"`
}

// UnidentifiedPaymentRequest the request to add an unidentified payment
type UnidentifiedPaymentRequest struct {
	ClientBPID              string  `json:"client_bpid"`
	TransactionDate         string  `json:"txn_date"`
	TransactionType         string  `json:"txn_type"`
	ValueDate               string  `json:"value_date"`
	NameOfCompany           string  `json:"name_of_company"`
	Amount                  float64 `json:"amount"`
	FundManager             int     `json:"fund_manager"`
	CollectionAccountNumber string  `json:"collection_acc_num"`
	Status                  string  `json:"status"`
}

type NPRAUnauthorizedTransactionRequest struct {
	ClientName         string `json:"clientName"`
	TransactionDetails string `json:"txnDetails"`
	Date               string `json:"date"`
}

type NPRAOutstandingFDCertificateRequest struct {
	FundManager   string    `json:"fundManager"`
	ClientName    string    `json:"clientName"`
	Amount        float64   `json:"amount"`
	Issuer        string    `json:"issuer"`
	Rate          float32   `json:"rate"`
	Tenor         int       `json:"tenor"`
	Term          string    `json:"term"`
	EffectiveDate time.Time `json:"effectiveDate"`
	Maturity      time.Time `json:"maturity"`
}

type OnlyBPID struct {
	BPID string `json:"bpid" form:"bpid"`
}

type OnlyID struct {
	ID int `json:"id" form:"id"`
}

type BPOrSCA struct {
	BpOrSca string `json:"bp_or_sca" form:"bpOrSca"`
	Date    string `json:"date" form:"period"`
}

type OnlyCode struct {
	Code string `json:"code"`
}

type BasisPointsRequest struct {
	Minimum     float64 `json:"minimum_amount"`
	Maximum     float64 `json:"maximum_amount"`
	BasisPoints float64 `json:"basis_point"`
}

type ClientBillingInfoUpdateRequest struct {
	BPID                 string               `json:"bpid"`
	SCA                  string               `json:"sca"`
	MinimumCharge        float32              `json:"minimum_charge"`
	ChargePerTransaction float32              `json:"charge_per_txn"`
	ThirdPartyTransfer   float32              `json:"third_party_transfer"`
	BasisPoints          []BasisPointsRequest `json:"basis_points"`
}

type CheckerCommentRequest struct {
	Comment string `form:"comment"`
}

type VarianceRemarksUpdateRequest struct {
	Client  []string `form:"client"`
	Remarks []string `form:"remarks"`
}

type OutstandingFDEdit struct {
	ID    int  `json:"id"`
	Value bool `json:"value"`
}

type BillingNAVUpdateRequest struct {
	Date        string  `json:"date"`
	BPID        string  `json:"bpid"`
	SCA         string  `json:"sca"`
	NAV         float64 `json:"nav"`
	CashBalance float64 `json:"cash_balance"`
	Liabilities float64 `json:"liabilities"`
	Position    float64 `json:"position"`
}

type GOGMaturitiesUploadRequest struct {
	Data    []GOGMaturity `json:"data"`
	Quarter string        `json:"quarter"`
	Year    string        `json:"year"`
}

type TransactionVolumesUploadRequest struct {
	Data    []CrsTrustTxnVolumes `json:"data"`
	Quarter string               `json:"quarter"`
	Year    string               `json:"year"`
}

type DeleteUploadedPVRequest struct {
	BPID string `json:"bpid"`
	Date string `json:"date"`
}

type OnlyHash struct {
	Hash string `form:"hash"`
}

type CheckerCommentRequestWithDateAndClientID struct {
	Comment string `form:"comment"`
	Date    string `form:"date"`
	BPOrSCA string `form:"bpOrSca"`
}

type SCAManagementRequest struct {
	ID            int    `form:"id"`
	BPID          string `form:"bpid"`
	ClientName    string `form:"name"`
	ClientTier    int    `form:"tier"`
	Administrator string `form:"administrator"`
	FundManager   string `form:"fundManager"`
	Code          string `form:"code"`
	GroupKey      string `form:"group_key"`
}

type ManageDirectorRequest struct {
	ID                 int    `form:"id"`
	Fullname           string `form:"fullname"`
	ResidentialAddress string `form:"address"`
	Nationality        string `form:"nationality"`
	Position           string `form:"position"`
}

type MaturedSecuritiesUpload struct {
	Maturities []MaturedSecurity `json:"maturities"`
}

type OnlyQuarter struct {
	Quarter string `form:"quarter" json:"quarter"`
}

type OnlyStaffID struct {
	StaffID int `form:"staffID"`
}

type UserProfileUpdateRequest struct {
	Password string `form:"password"`
}

type SCAAccountMgmt struct {
	BPID string `json:"bpid" form:"bpid"`
	Code string `json:"code" form:"code"`
}

type ClientEmailMgmt struct {
	ID    int    `json:"id" form:"id"`
	BPID  string `json:"bpid" form:"bpid"`
	Email string `json:"email" form:"email"`
}

type MailServiceMgmt struct {
	ID    int    `json:"id" form:"id"`
	Email string `json:"email" form:"email"`
}

type MailSenderInfoUpdateRequest struct {
	Name     string `json:"name" form:"name"`
	Position string `json:"position" form:"position"`
}

type HolidayMgmtRequest struct {
	ID   int    `json:"id" form:"id"`
	Name string `json:"name" form:"name"`
	Date string `json:"date" form:"date"`
}

// AccountSetupRequest request type
type AccountSetupRequest struct {
	//	Password string `form:"password"`
	Password string `json:"password" form:"password"`
}

// ClientAddRequest the request for adding a client
type NPRA0303AddRequest struct {
	ReportCode string `form:"report_code"`
	//	BPID                string       `form:"bp_id"`
	EntityID            string       `form:"entity_id"`
	EntityName          string       `form:"entity_name"`
	ReferencePeriodYear string       `form:"reference_period_year"`
	ReferencePeriod     string       `form:"reference_period"`
	UnitPrice           float64      `form:"unit_price"`
	DateValuation       time.Time    `form:"date_valuation"`
	UnitNumber          string       `form:"unit_number"`
	DailyNav            string       `form:"daily_nav"`
	NPRAFees            float64      `form:"npra_fees"`
	TrusteeFees         float64      `form:"trustee_fees"`
	FundManagerFees     float64      `form:"fund_manager_fees"`
	FundCustodianFees   float64      `form:"fund_custodian_fees"`
	CreatedBy           nulls.String `form:"created_by"`
}

// AssetClassAddRequest the request for adding a client
type AssetClassAddRequest struct {
	SecurityType   string `form:"security_type"`
	AssetClassName string `form:"asset_class_name"`
}

// AssetClassEditRequest the request to edit a client
type AssetClassEditRequest struct {
	AssetClassAddRequest
	ID string `form:"id"`
}

// NPRA0303EditRequest the request to edit a client
type NPRA0303EditRequest struct {
	NPRA0303AddRequest
	ID string `form:"id"`
}

// PVFetchForEyeBallingRequest Fetch PV only
type Fetch0302ForEyeBallingRequest struct {
	SchemeDetailsFetchRequest
}

// PVFetchForEyeBallingRequest Fetch PV only
type Fetch0301ForEyeBallRequest struct {
	Details0301FetchRequest
}

// requests for NPRA Reports
type NPRAReportRequest struct {
	BPID  string `form:"bp_id" json:"bp_id"`
	Type  string `form:"type" json:"type"`
	Month string `form:"month" json:"month"`
	Year  string `form:"year" json:"year"`
}

// requests for NPRA Reports
type NPRA321Request struct {
	BPID string `form:"bpOrSca" json:"bpOrSca"`
	Date string `form:"date" json:"date"`
}
