package actions

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/jinzhu/now"
	"golang.org/x/crypto/bcrypt"
)

// GetUserByStaffID retrieves a user from the db based on their staff id
func GetUserByStaffID(staffID int) User {
	user := User{}
	DatabaseConnection.Preload("Roles.Role.Permissions.Permission").Preload("Permissions.Permission").Where("STAFF_ID=?", staffID).Find(&user)
	return user
}

// GetUserByID retrieves a user by ID
func GetUserByID(id int) User {
	user := User{}
	DatabaseConnection.Preload("Roles.Role.Permissions.Permission").Preload("Permissions.Permission").Where("ID=?", id).Find(&user)
	return user
}

// LoadAllClients loads all clients from the database
func LoadAllClients() []Client {
	return LoadNClients(-1)
}

// LoadNClients loads a given number of cients
func LoadNClients(limit int) []Client {
	clients := make([]Client, 0)
	if limit == -1 {
		DatabaseConnection.Order("CLIENTNAME asc").Find(&clients)
	} else {
		DatabaseConnection.Order("CLIENTNAME asc").Limit(limit).Find(&clients)
	}
	return clients
}

// LoadClientDetails loads all client details with their pv summary for trustee
func LoadClientDetails(bpidOrSCA, quarter string) Client {
	client := Client{}
	DatabaseConnection.Where("BP_ID=?", bpidOrSCA).Preload("PVReportSummary", func(db *gorm.DB) *gorm.DB {
		return db.Where("REPORT_DATE = ? AND CLIENT_TYPE=?", quarter, "trustee")
	}).Preload("Reports", func(db *gorm.DB) *gorm.DB {
		return db.Where("REPORT_DATE = ? AND CLIENT_TYPE=?", quarter, "trustee")
	}).First(&client)
	if client.ID == 0 {
		scaClient := GetClientByCode(bpidOrSCA)
		client = GetClientByBPID(scaClient.BPID)
		//to show the sca name as the client name
		client.Name = scaClient.ClientName
		client.BPID = bpidOrSCA
		reports := make([]PVReportField, 0)
		summary := make([]ClientPVReportSummary, 0)
		DatabaseConnection.Where("CLIENT_ID = ? AND REPORT_DATE = ? AND CLIENT_TYPE=?", bpidOrSCA, quarter, "trustee").Find(&reports)
		DatabaseConnection.Where("CLIENT_ID = ? AND REPORT_DATE = ? AND CLIENT_TYPE=?", bpidOrSCA, quarter, "trustee").Find(&summary)
		client.PVReportSummary = summary
		client.Reports = reports
	}
	return client
}

// LoadOverallSecuritiesSummary loads the overall summary for all bonds
func LoadOverallSecuritiesSummary() [][]BondOverallSummary {
	dates := GetQuarterDatesBetween(time.Now().Year()-1, time.Now().Year())
	summary := make([][]BondOverallSummary, len(dates))
	for index, date := range dates {
		DatabaseConnection.Raw(fmt.Sprintf("select security_type as bond,sum(lcy_amt) as value from CRS_PVSUMMARY  where report_date >= ? and report_date <= ? and SECURITY_TYPE IN %s group by SECURITY_TYPE", MakeReportableSecuritiesQueryFromHeadings(LoadPVReportHeadings())), date.Begin, date.End).Scan(&summary[index])
	}
	return summary
}

// InsertOrdinarySharesIntoTheDB inserts the ordinary shares values into the db
func InsertOrdinarySharesIntoTheDB(governanceInfoID int, data []GovernanceShare) {
	DatabaseConnection.Exec("DELETE FROM CRS_SEC_GOV_ORDINARY_SHARES WHERE GOV_INFO_ID=?", governanceInfoID)
	for _, share := range data {
		shareData := OrdinaryShare{
			GovernanceInfoID: governanceInfoID,
			Name:             share.Name,
			Shareholdings:    share.Shareholding,
			Percentage:       share.Percentage,
			CreatedAt:        time.Now(),
		}
		DatabaseConnection.Create(&shareData)
	}
}

// InsertPreferenceSharesIntoTheDB inserts the preference shares values into the db
func InsertPreferenceSharesIntoTheDB(governanceInfoID int, data []GovernanceShare) {
	DatabaseConnection.Exec("DELETE FROM CRS_SEC_GOV_PREFERENCE_SHARES WHERE GOV_INFO_ID=?", governanceInfoID)
	for _, share := range data {
		shareData := PreferenceShare{
			GovernanceInfoID: governanceInfoID,
			Name:             share.Name,
			Shareholdings:    share.Shareholding,
			Percentage:       share.Percentage,
			CreatedAt:        time.Now(),
		}
		DatabaseConnection.Create(&shareData)
	}
}

// InsertTransactionsWithAffiliatesIntoTheDB inserts the transactions with the affiliates into the db
func InsertTransactionsWithAffiliatesIntoTheDB(governanceInfoID int, data []GovernanceCustodianTransaction) {
	DatabaseConnection.Exec("DELETE FROM CRS_SEC_TXN_WITH_CUSTODIAN_AFFILIATES WHERE GOV_INFO_ID=?", governanceInfoID)
	for _, transaction := range data {
		transactionData := AffiliateTransaction{
			GovernanceInfoID:          governanceInfoID,
			NameOfAffiliate:           transaction.NameOfTrustee,
			RelationshipWithCustodian: transaction.RelationshipWIthTrustee,
			TypeOfTransaction:         transaction.TypeOfTransaction,
			Amount:                    transaction.Amount,
			CreatedAt:                 time.Now(),
		}
		DatabaseConnection.Create(&transactionData)
	}

}

// InsertValueVolumeOfShareSchemesIntoTheDB inserts the value volume of shares int the db
func InsertValueVolumeOfShareSchemesIntoTheDB(governanceInfoID int, data []GovernanceSchemesUnderCustody) {
	DatabaseConnection.Exec("DELETE FROM CRS_SEC_GOV_VALUE_VOLUME_OF_SHARE_SCHEMES WHERE GOV_INFO_ID=?", governanceInfoID)
	for _, scheme := range data {
		schemeData := &ValueVolumeOfShareScheme{
			GovernanceInfoID:          governanceInfoID,
			NameOfFirm:                scheme.NameOfFirm,
			NameOfScheme:              scheme.NameOfScheme,
			RelationshipWithCustodian: scheme.RelationshipWIthTrustee,
			Volume:                    scheme.Volume,
			MarkedToMarketValue:       scheme.MarkedToMarketValue,
			CreatedAt:                 time.Now(),
		}
		DatabaseConnection.Create(schemeData)
	}

}

// InsertGovernanceInfoIntoTheDB inserts the governance info data into the db
func InsertGovernanceInfoIntoTheDB(data GovernanceDataUploadRequest) error {
	quarterDate := GetQuarterDate()
	date := time.Now()
	if data.CustodianHasUpdatedAssetRegister == "yes" {
		date, _ = time.Parse("2006-01-02", data.CustodianAssetRegistrationDate)
	}
	governanceInfo := GovernanceInfo{
		ReportRefID:                        quarterDate,
		CustodianName:                      data.ClientName,
		ReportingPeriod:                    quarterDate,
		ReportingOfficer:                   data.ReportingOfficer,
		DateOfReport:                       time.Now(),
		RenewalDate:                        GetSecLastLicenseRenewalDate(),
		ChangeInDirectors:                  data.ChangeInDirectors,
		ChangeInAgreement:                  data.ChangeInAgreement,
		DealingsApprovedByBoard:            data.DealingsApprovedByBoard,
		CustodianHasUpdatedAssetRegister:   data.CustodianHasUpdatedAssetRegister,
		CustodianAssetRegistrationDate:     date,
		DoManagersOfTheSchemeConsultTheLaw: data.DoManagersOfTheSchemeConsultTheLaw,
		SchemeHadAnyOtherFinancialDealings: data.SchemeHadAnyOtherFinancialDealings,
		CreatedAt:                          time.Now(),
		UpdatedAt:                          time.Now(),
	}
	DatabaseConnection.Where(GovernanceInfo{ReportRefID: quarterDate}).Assign(governanceInfo).FirstOrCreate(&governanceInfo)
	if governanceInfo.ID != 0 {
		InsertOrdinarySharesIntoTheDB(governanceInfo.ID, data.OrdinaryShares)
		InsertPreferenceSharesIntoTheDB(governanceInfo.ID, data.PreferenceShares)
		InsertTransactionsWithAffiliatesIntoTheDB(governanceInfo.ID, data.CustodianTransactions)
		InsertValueVolumeOfShareSchemesIntoTheDB(governanceInfo.ID, data.SchemesUnderCustody)
		//send email to checkers
		return nil
	}

	return errors.New("Couldnt create governance data")
}

// LoadGovernanceInfoByQuarterDate retrieves the governance info data id for the given quarter date
func LoadGovernanceInfoByQuarterDate(date string) GovernanceInfo {
	governanceInfo := GovernanceInfo{}
	DatabaseConnection.Where("REPORT_REF_ID=?", date).Preload("OrdinaryShares").Preload("PreferenceShares").Preload("AffiliateTransactions").Preload("ValueVolumeOfShareSchemes").Find(&governanceInfo)
	return governanceInfo
}

// LoadQuarterSchemeSubmissionDetails loads all scheme submission data for the quarter
func LoadQuarterSchemeSubmissionDetails() int64 {
	value := IntValueQuery{}
	quarterDate := GetQuarterDate()
	existingGovernanceInfo := LoadGovernanceInfoByQuarterDate(quarterDate)
	if existingGovernanceInfo.IsEmpty() || !SchemeDetailsHasBeenUploadedForQuarter(existingGovernanceInfo.ID) {
		return 0
	}
	DatabaseConnection.Raw(`SELECT COUNT(*) AS value FROM CRS_SEC_SCHEME_ASSET_HOLDINGS WHERE GOV_INFO_ID= ?`, existingGovernanceInfo.ID).Scan(&value)
	return value.Value
}

// SchemeDetailsHasBeenUploadedForQuarter checks if scheme details data has been uploaded for the quarter
func SchemeDetailsHasBeenUploadedForQuarter(governaceID int) bool {
	count := 0
	DatabaseConnection.Model(&SchemeDetails{}).Where("GOV_INFO_ID = ?", governaceID).Count(&count)
	return count > 0
}

func DeleteSchemeDetails(bpid string) error {
	quarterDate := GetQuarterDate()
	existingGovernanceInfo := LoadGovernanceInfoByQuarterDate(quarterDate)
	if existingGovernanceInfo.IsEmpty() {
		return errors.New("cannot delete scheme details. Governance data for quarter does not exist")
	}
	DatabaseConnection.Exec("DELETE FROM CRS_SEC_SCHEME_ASSET_HOLDINGS WHERE BP_ID=? AND GOV_INFO_ID=?", bpid, existingGovernanceInfo.ID)
	return nil
}

// InsertSchemeDetailsIntoTheDB inserts the scheme details into the db
func InsertSchemeDetailsIntoTheDB(data SchemeDetails, UserID int) error {
	quarterDate := GetQuarterDate()
	existingGovernanceInfo := LoadGovernanceInfoByQuarterDate(quarterDate)
	if existingGovernanceInfo.IsEmpty() {
		return errors.New("cannot upload scheme details. Please upload governance data for quarter first")
	}
	scheme := SchemeDetails{
		GovernanceInfoID:                        data.GovernanceInfoID,
		BPID:                                    data.BPID,
		NameOfScheme:                            data.NameOfScheme,
		NumberOfSharesOutstanding:               data.NumberOfSharesOutstanding,
		NumberOfShareholders:                    data.NumberOfShareholders,
		NumberOfSharesRedeemed:                  data.NumberOfSharesRedeemed,
		ValueOfSharesRedeemed:                   data.ValueOfSharesRedeemed,
		NameOfManager:                           data.NameOfManager,
		TotalValueOfSchemeAssets:                data.TotalValueOfSchemeAssets,
		NetAssetValueOfScheme:                   data.NetAssetValueOfScheme,
		NetAssetValueOfSchemePerUnit:            data.NetAssetValueOfSchemePerUnit,
		TotalEquityInvestment:                   data.TotalEquityInvestment,
		TotalFixedIncomeInvestment:              data.TotalFixedIncomeInvestment,
		NetMediumAssetHeldByFund:                data.NetMediumAssetHeldByFund,
		CapitalMarkerInvestments:                data.CapitalMarkerInvestments,
		PercentageCapitalMarketInvestment:       data.PercentageCapitalMarketInvestment,
		CertificatesOfInvestmentWithCustodian:   data.CertificatesOfInvestmentWithCustodian,
		TotalValueOfUnutilizedFunds:             data.TotalValueOfUnutilizedFunds,
		ValueOfBorrowedFunds:                    data.ValueOfBorrowedFunds,
		ReasonsForBorrowing:                     data.ReasonsForBorrowing,
		AuditedAccountsDistributedToAuthorities: data.AuditedAccountsDistributedToAuthorities,
		Redemptions:                             data.Redemptions,
		Dividends:                               data.Dividends,
		Rights:                                  data.Rights,
		FeesOwedCustodian:                       data.FeesOwedCustodian,
	}
	if strings.TrimSpace(data.AttachedFile) != "" {
		if filename, err := CreateFileFromBase64Data(data.AttachedFile, fmt.Sprintf("%s-%s-SchemeDetails.pdf", existingGovernanceInfo.ReportRefID, data.BPID)); err == nil {
			scheme.AttachedFile = filename
		}
	}
	DatabaseConnection.Where(SchemeDetails{GovernanceInfoID: existingGovernanceInfo.ID, BPID: data.BPID}).Assign(scheme).FirstOrCreate(&scheme)
	if scheme.ID != 0 {
		LogRecentActivity(UserID, data.BPID, UploadScheme)
		return nil
	}
	return errors.New("couldn't store scheme details")
}

// InsertOtherInformationIntoTheDB inserts the scheme details into the db
func InsertOtherInformationIntoTheDB(data OtherInformation) error {
	quarterDate := GetQuarterDate()
	existingGovernanceInfo := LoadGovernanceInfoByQuarterDate(quarterDate)
	if existingGovernanceInfo.IsEmpty() {
		return errors.New("cannot upload other information. Please upload governance data for quarter first")
	}
	info := OtherInformation{
		AreThereAnyClaimOnSchemeAsset:               data.AreThereAnyClaimOnSchemeAsset,
		YesWasCustodianInformedAndApproved:          data.YesWasCustodianInformedAndApproved,
		AnyLitigationInvolvingCustodianSccheme:      data.AnySignificantReductionInAssetScheme,
		AnySignificantReductionInAssetScheme:        data.AnySignificantReductionInAssetScheme,
		SignificantReductionInSchemeMarketPrice:     data.SignificantReductionInSchemeMarketPrice,
		HasMgrsReconciledAssetRegisterCustodian:     data.HasMgrsReconciledAssetRegisterCustodian,
		HowManyTimesDidSchemePublishedPrices:        data.HowManyTimesDidSchemePublishedPrices,
		AnyConcernsByInvestors:                      data.AnyConcernsByInvestors,
		AnyMattersAttentionSecMgtCustodyOfFund:      data.AnyMattersAttentionSecMgtCustodyOfFund,
		HasAccountOfManagersSeparateFromScheme:      data.HasAccountOfManagersSeparateFromScheme,
		CompanyParentsAffiliateInvolvedInLitigation: data.CompanyParentsAffiliateInvolvedInLitigation,
		LitigationDetails:                           data.LitigationDetails,
	}
	DatabaseConnection.Where(OtherInformation{GovernanceInfoID: existingGovernanceInfo.ID}).Assign(info).FirstOrCreate(&info)
	if info.ID != 0 {
		return nil
	}
	return errors.New("Couldnt store other information")
}

// InsertOfficialReportRemarksIntoTheDB inserts the scheme details into the db
func InsertOfficialReportRemarksIntoTheDB(data OfficialReportRemarks) error {
	quarterDate := GetQuarterDate()
	existingGovernanceInfo := LoadGovernanceInfoByQuarterDate(quarterDate)
	if existingGovernanceInfo.IsEmpty() {
		return errors.New("cannot upload official remarks. Please upload governance data for quarter first")
	}
	remarks := OfficialReportRemarks{
		Remarks:          data.Remarks,
		ReviewingOfficer: data.ReviewingOfficer,
		Date:             data.Date,
		Signature:        data.Signature,
	}
	DatabaseConnection.Where(OfficialReportRemarks{GovernanceInfoID: existingGovernanceInfo.ID}).Assign(remarks).FirstOrCreate(&remarks)
	if remarks.ID != 0 {
		return nil
	}
	return errors.New("couldn't store official remarks")
}

// GetSchemeDetailsByGovernanceInfoID finds the scheme details for a governance id
func GetSchemeDetailsByGovernanceInfoID(id int) []SchemeDetails {
	details := make([]SchemeDetails, 0)
	DatabaseConnection.Where("GOV_INFO_ID=?", id).Find(&details)
	return details
}

// GetOtherInformationByGovernanceInfoID finds the other info for a governance id
func GetOtherInformationByGovernanceInfoID(id int) OtherInformation {
	info := OtherInformation{}
	DatabaseConnection.Where("GOV_INFO_ID=?", id).Find(&info)
	return info
}

// GetOfficialRemarksByGovernanceInfoID finds the official remarks for a governance id
func GetOfficialRemarksByGovernanceInfoID(id int) OfficialReportRemarks {
	remarks := OfficialReportRemarks{}
	DatabaseConnection.Where("GOV_INFO_ID=?", id).Find(&remarks)
	return remarks
}

// FetchSchemeDetails loads the scheme details that matches the request
func FetchSchemeDetails(request SchemeDetailsFetchRequest) SchemeDetails {
	scheme := SchemeDetails{}
	governanceInfo := LoadGovernanceInfoByQuarterDate(request.QuarterDate)
	if governanceInfo.ID != 0 {
		DatabaseConnection.Where("GOV_INFO_ID=?", governanceInfo.ID).Where("BP_ID=?", request.BPID).Find(&scheme)
	}
	return scheme
}

// CreateNewUser creates a new user
func CreateNewUser(request UserAddRequest) {
	newPassword := RandomBytes(8)
	if password, err := bcrypt.GenerateFromPassword(newPassword, 10); err == nil {
		user := User{
			StaffID:                request.StaffID,
			Fullname:               request.Fullname,
			Email:                  request.Email,
			Password:               string(password),
			MustResetPassword:      true,
			LastPasswordChangeDate: time.Now(),
			Active:                 request.Active,
		}
		DatabaseConnection.Create(&user)
		if user.ID > 0 {
			for _, role := range request.Roles {
				DatabaseConnection.Create(UserRole{
					UserID: user.ID,
					RoleID: role,
				})
			}
			//TODO: send welcome email
		}
	}
}

// LoadRolesFromDB loads all roles from the DB
func LoadRolesFromDB() []Role {
	roles := make([]Role, 0)
	DatabaseConnection.Find(&roles)
	return roles
}

// LoadUsersFromDB loads the users from the DB
func LoadUsersFromDB() []User {
	users := make([]User, 0)
	DatabaseConnection.Preload("Roles.Role").Where("LOCKED=?", false).Find(&users)
	return users
}

// ResetUserPassword resets the users password
func ResetUserPassword(userID int, email string) {
	newPassword := RandomBytes(8)
	if password, err := bcrypt.GenerateFromPassword(newPassword, 10); err == nil {
		DatabaseConnection.Table("CRS_USERS").Where("id = ?", userID).Updates(map[string]interface{}{"PASSWORD": password, "MUSTRESETPASSWORD": true})
		SendPasswordResetEmail(UserPasswordResetPayload{
			Email:    email,
			Password: string(newPassword),
		})
	}
}

// DeleteUser delete the user
func DeleteUser(request OnlyUserIDRequest) {
	DatabaseConnection.Table("CRS_USERS").Where("id = ?", request.UserID).Update("LOCKED", true)
}

// EditUser edits the users information in the DB
func EditUser(request EditUserRequest) {
	DatabaseConnection.Table("CRS_USERS").Where("id = ?", request.UserID).Updates(map[string]interface{}{"FULLNAME": request.Fullname, "STAFF_ID": request.StaffID, "EMAIL": request.Email, "ACTIVE": request.Active})
	DatabaseConnection.Where("USERID=?", request.UserID).Delete(UserRole{})
	for _, role := range request.Roles {
		DatabaseConnection.Create(UserRole{
			UserID: request.UserID,
			RoleID: role,
		})
	}
}

// LoadClientPVSummary loads the client's pv summary for a given quarter and year
func LoadClientPVSummary(bpidOrSca, date string) []ClientPVReportSummary {
	summary := make([]ClientPVReportSummary, 0)
	DatabaseConnection.Table("CRS_PVSUMMARY").Where("REPORT_DATE=? AND CLIENT_TYPE=?", date, "trustee").Where("BP_ID=? OR CLIENT_ID=?", bpidOrSca, bpidOrSca).Find(&summary)
	return summary
}

// LoadOffshoreClients loads all foreign clients
func LoadOffshoreClients(request OffshoreClientsRequest) []OffshoreClient {
	clients := make([]OffshoreClient, 0)
	date := MakeDBQueryableQuarterLastDate(request.Quarter, request.Year)
	if strings.TrimSpace(request.BPID) == "" {
		DatabaseConnection.Raw("select CLIENTNAME as [name], HOME_COUNTRY as country , (select sum(LCY_AMT) from CRS_PVSUMMARY WHERE BP_ID = CRS_CLIENT.BP_ID AND REPORT_DATE = ?) as assets from CRS_CLIENT WHERE CLIENT_TYPE = 'FOREIGN CLIENT' and( CLOSED is null or CLOSED > ?) order by CLIENTNAME asc ", date, date).Scan(&clients)
	} else {
		DatabaseConnection.Raw("select CLIENTNAME as [name], HOME_COUNTRY as country , (select sum(LCY_AMT) from CRS_PVSUMMARY WHERE BP_ID = CRS_CLIENT.BP_ID AND REPORT_DATE = ?) as assets from CRS_CLIENT WHERE BP_ID=? and( CLOSED is null or CLOSED > ?) order by CLIENTNAME asc", date, request.BPID, date).Scan(&clients)
	}
	return CreateOffshoreClientsTotal(clients)
}

// LoadLocalVariance loads the sec local variance data
func LoadSecLocalVariance(lastQuarter, currentQuarter string) []VarianceData {
	variance := make([]VarianceQueryData, 0)
	DatabaseConnection.Raw("select CLIENTNAME as [name],HOME_COUNTRY as country, (select sum(LCY_AMT) from CRS_PVSUMMARY where BP_ID = CRS_CLIENT.BP_ID and REPORT_DATE=? and CLIENT_TYPE='sec') as l_aua,(select sum(LCY_AMT) from CRS_PVSUMMARY where BP_ID = CRS_CLIENT.BP_ID and REPORT_DATE=? and CLIENT_TYPE='sec') as c_aua , (select TOP(1) remarks from CRS_VARIANCE_REMARKS where BP_ID = CRS_CLIENT.BP_ID and REPORT_DATE=?) as remarks from CRS_CLIENT WHERE CLIENT_TYPE = 'SEC' and( CLOSED is null or CLOSED > ?) ORDER BY CLIENTNAME ASC", lastQuarter, currentQuarter, currentQuarter, currentQuarter).Scan(&variance)
	data := make([]VarianceData, len(variance))
	for index, datum := range variance {
		var varianceNum float32 = 100
		if datum.LastAUA == 0 && datum.CurrentAUA == 0 {
			varianceNum = 0
		}
		if datum.LastAUA != 0 {
			varianceNum = float32(((datum.CurrentAUA - datum.LastAUA) / datum.LastAUA) * 100)
		}
		data[index] = VarianceData{
			Name:       datum.Name,
			Country:    datum.Country,
			LastAUA:    datum.LastAUA,
			CurrentAUA: datum.CurrentAUA,
			Amount:     datum.CurrentAUA - datum.LastAUA,
			Variance:   varianceNum,
			Remarks:    datum.Remarks,
		}
	}
	return CreateVarianceTotal(data)
}

// LoadForeignVariance loads the sec foreign variance data
func LoadSecForeignVariance(lastQuarter, currentQuarter string) []VarianceData {
	variance := make([]VarianceQueryData, 0)
	DatabaseConnection.Raw("select CLIENTNAME as [name],HOME_COUNTRY as country, (select sum(LCY_AMT) from CRS_PVSUMMARY where BP_ID = CRS_CLIENT.BP_ID and REPORT_DATE=? and CLIENT_TYPE='sec') as l_aua,(select sum(LCY_AMT) from CRS_PVSUMMARY where BP_ID = CRS_CLIENT.BP_ID and REPORT_DATE=? and CLIENT_TYPE='sec') as c_aua , (select TOP(1) remarks from CRS_VARIANCE_REMARKS where BP_ID = CRS_CLIENT.BP_ID and REPORT_DATE=?) as remarks from CRS_CLIENT WHERE CLIENT_TYPE = 'FOREIGN CLIENT' and( CLOSED is null or CLOSED > ?) ORDER BY CLIENTNAME ASC", lastQuarter, currentQuarter, currentQuarter, currentQuarter).Scan(&variance)
	data := make([]VarianceData, len(variance))
	for index, datum := range variance {
		var varianceNum float32 = 100
		if datum.LastAUA == 0 && datum.CurrentAUA == 0 {
			varianceNum = 0
		}
		if datum.LastAUA != 0 {
			varianceNum = float32(((datum.CurrentAUA - datum.LastAUA) / datum.LastAUA) * 100)
		}
		data[index] = VarianceData{
			Name:       datum.Name,
			Country:    datum.Country,
			LastAUA:    datum.LastAUA,
			CurrentAUA: datum.CurrentAUA,
			Amount:     datum.CurrentAUA - datum.LastAUA,
			Variance:   varianceNum,
			Remarks:    datum.Remarks,
		}
	}
	return CreateVarianceTotal(data)
}

// LoadNpraLocalVariance loads the npra local variance data
func LoadNpraLocalVariance(lastQuarter, currentQuarter string) []VarianceData {
	variance := make([]VarianceQueryData, 0)
	DatabaseConnection.Raw(`
select 
	distinct(CLIENT_NAME) as [name],
	(select HOME_COUNTRY from CRS_CLIENT where BP_ID = query.BP_ID) as country,
	(select case
		when 
		(select count(*) from CRS_CLIENT_PV_UNIQUE_CODES where CLIENT_NAME=query.CLIENT_NAME) > 1
		then
		(select sum(LCY_AMT) from CRS_PVSUMMARY where BP_ID = query.BP_ID and REPORT_DATE=? and CLIENT_TYPE='npra')
		else
		(select sum(LCY_AMT) from CRS_PVSUMMARY where CLIENT_ID = query.CODE and REPORT_DATE=? and CLIENT_TYPE='npra')
		end
	) as l_aua, 
	(select case
		when 
		(select count(*) from CRS_CLIENT_PV_UNIQUE_CODES where CLIENT_NAME=query.CLIENT_NAME) > 1
		then
		(select sum(LCY_AMT) from CRS_PVSUMMARY where BP_ID = query.BP_ID and REPORT_DATE=? and CLIENT_TYPE='npra')
		else
		(select sum(LCY_AMT) from CRS_PVSUMMARY where CLIENT_ID = query.CODE and REPORT_DATE=? and CLIENT_TYPE='npra')
		end
	) as c_aua ,
	(select case
	when
	(select count(*) from CRS_CLIENT_PV_UNIQUE_CODES where CLIENT_NAME=query.CLIENT_NAME) > 1
	then
	(select TOP(1) remarks from CRS_VARIANCE_REMARKS where BP_ID = query.BP_ID and REPORT_DATE=?)
	else 
	(select TOP(1) remarks from CRS_VARIANCE_REMARKS where BP_ID = query.CODE and REPORT_DATE=?)
	end
	) as remarks
from CRS_CLIENT_PV_UNIQUE_CODES query where CLIENT_TIER in (2,3) and( CLOSED is null or CLOSED > ?)
`, lastQuarter, lastQuarter, currentQuarter, currentQuarter, currentQuarter, currentQuarter, currentQuarter).Scan(&variance)
	data := make([]VarianceData, len(variance))
	for index, datum := range variance {
		var varianceNum float32 = 100
		if datum.LastAUA == 0 && datum.CurrentAUA == 0 {
			varianceNum = 0
		}
		if datum.LastAUA != 0 {
			varianceNum = float32(((datum.CurrentAUA - datum.LastAUA) / datum.LastAUA) * 100)
		}
		data[index] = VarianceData{
			Name:       datum.Name,
			Country:    datum.Country,
			LastAUA:    datum.LastAUA,
			CurrentAUA: datum.CurrentAUA,
			Amount:     datum.CurrentAUA - datum.LastAUA,
			Variance:   varianceNum,
			Remarks:    datum.Remarks,
		}
	}
	return CreateVarianceTotal(data)
}

// LoadNpraForeignVariance loads the npra foreign variance data
func LoadNpraForeignVariance(lastQuarter, currentQuarter string) []VarianceData {
	variance := make([]VarianceQueryData, 0)
	DatabaseConnection.Raw("select CLIENTNAME as [name],HOME_COUNTRY as country, (select sum(LCY_AMT) from CRS_PVSUMMARY where BP_ID = CRS_CLIENT.BP_ID and REPORT_DATE=? and CLIENT_TYPE='npra') as l_aua,(select sum(LCY_AMT) from CRS_PVSUMMARY where BP_ID = CRS_CLIENT.BP_ID and REPORT_DATE=? and CLIENT_TYPE='npra') as c_aua , (select TOP(1) remarks from CRS_VARIANCE_REMARKS where BP_ID = CRS_CLIENT.BP_ID and REPORT_DATE=?) as remarks from CRS_CLIENT WHERE CLIENT_TYPE = 'FOREIGN CLIENT'and( CLOSED is null or CLOSED > ?) ORDER BY CLIENTNAME ASC", lastQuarter, currentQuarter, currentQuarter, currentQuarter).Scan(&variance)
	data := make([]VarianceData, len(variance))
	for index, datum := range variance {
		var varianceNum float32 = 100
		if datum.LastAUA == 0 && datum.CurrentAUA == 0 {
			varianceNum = 0
		}
		if datum.LastAUA != 0 {
			varianceNum = float32(((datum.CurrentAUA - datum.LastAUA) / datum.LastAUA) * 100)
		}
		data[index] = VarianceData{
			Name:       datum.Name,
			Country:    datum.Country,
			LastAUA:    datum.LastAUA,
			CurrentAUA: datum.CurrentAUA,
			Amount:     datum.CurrentAUA - datum.LastAUA,
			Variance:   varianceNum,
			Remarks:    datum.Remarks,
		}
	}
	return CreateVarianceTotal(data)
}

// LoadOverviewData loads the overview data
func LoadOverviewData() []OverviewSummaryData {
	overview := make([]OverviewSummaryData, 4)
	dates := GetLast4QuarterDates()
	lastQuarter := now.New(dates[0].AddDate(0, -3, -5)).EndOfMonth()
	lastQuarterTotal := FloatArraySum(GetOverviewBondData(lastQuarter))
	date1Data := GetOverviewBondData(dates[0])
	date1DataTotal := FloatArraySum(date1Data)
	date2Data := GetOverviewBondData(dates[1])
	date2DataTotal := FloatArraySum(date2Data)
	date3Data := GetOverviewBondData(dates[2])
	date3DataTotal := FloatArraySum(date3Data)
	date4Data := GetOverviewBondData(dates[3])
	date4DataTotal := FloatArraySum(date4Data)

	overview[0] = OverviewSummaryData{
		PrettyDate:            MakeOverviewQuarterDate(dates[0]),
		Status:                true,
		IsMoreThanLastQuarter: date1DataTotal > lastQuarterTotal,
		Percentage:            CalculatePercentageDifference(lastQuarterTotal, date1DataTotal),
		QuarterDate:           dates[0].Format("2006-01-02"),
		Data:                  date1Data,
	}
	overview[1] = OverviewSummaryData{
		PrettyDate:            MakeOverviewQuarterDate(dates[1]),
		Status:                true,
		IsMoreThanLastQuarter: date2DataTotal > date1DataTotal,
		Percentage:            CalculatePercentageDifference(date1DataTotal, date2DataTotal),
		QuarterDate:           dates[1].Format("2006-01-02"),
		Data:                  date2Data,
	}
	overview[2] = OverviewSummaryData{
		PrettyDate:            MakeOverviewQuarterDate(dates[2]),
		Status:                true,
		IsMoreThanLastQuarter: date3DataTotal > date2DataTotal,
		Percentage:            CalculatePercentageDifference(date2DataTotal, date3DataTotal),
		QuarterDate:           dates[2].Format("2006-01-02"),
		Data:                  date3Data,
	}
	overview[3] = OverviewSummaryData{
		PrettyDate:            MakeOverviewQuarterDate(dates[3]),
		Status:                false,
		IsMoreThanLastQuarter: date4DataTotal > date3DataTotal,
		Percentage:            CalculatePercentageDifference(date3DataTotal, date4DataTotal),
		QuarterDate:           dates[3].Format("2006-01-02"),
		Data:                  date4Data,
	}
	return overview
}

// GetOverviewBondData overview bond data
func GetOverviewBondData(date time.Time) []float64 {
	type QAmount struct {
		Amount float64 `gorm:"column:amount"`
	}
	amounts := make([]QAmount, 0)
	DatabaseConnection.Raw(fmt.Sprintf("select SUM(LCY_AMT) as amount from CRS_PVSUMMARY  where REPORT_DATE = ? and SECURITY_TYPE IN %s group by SECURITY_TYPE", MakeReportableSecuritiesQueryFromHeadings(LoadPVReportHeadings())), date.Format("2006-01-02")).Scan(&amounts)
	data := make([]float64, len(amounts))
	for index, amount := range amounts {
		data[index] = amount.Amount
	}
	return data
}

// LoadQuarterOverview loads the quarter overview data
func LoadQuarterOverview(date string) []BondOverallSummary {
	summary := make([]BondOverallSummary, 0)
	DatabaseConnection.Raw(fmt.Sprintf("select SECURITY_TYPE as bond, SUM(LCY_AMT) as value from CRS_PVSUMMARY  where REPORT_DATE = ? and SECURITY_TYPE IN %s group by SECURITY_TYPE", MakeReportableSecuritiesQueryFromHeadings(LoadPVReportHeadings())), date).Scan(&summary)
	return summary
}

// NumberOfClientsAvailable returns the number of clients available
func NumberOfClientsAvailable() int {
	count := 0
	DatabaseConnection.Model(&Client{}).Count(&count)
	return count
}

// LoadUploadedPVData loads the uploaded data
func LoadUploadedPVData() UploadedData {
	clients := NumberOfClientsAvailable()
	count := CountQuery{}
	DatabaseConnection.Raw("select count(distinct BP_ID) as count from CRS_PVREPORTS where REPORT_DATE = ?", GetQuarterFormalDate()).Scan(&count)
	return UploadedData{Uploaded: count.Count, Pending: clients - count.Count}
}

// LoadAlertSettings loads the system alert settings
func LoadAlertSettings() AlertSettings {
	settings := AlertSettings{}
	DatabaseConnection.First(&settings)
	return settings
}

// LoadSecRecentActivities loads the sec recent activities from the database
func LoadSecRecentActivities() []SecActivities {
	activities := make([]SecActivities, 0)
	DatabaseConnection.Where("APPROVED = ?", false).Preload("User").Preload("Client").Order("ID desc").Find(&activities)
	return activities
}

// LoadNPRARecentActivities loads the sec recent activities from the database
func LoadNPRARecentActivities() []NPRAActivities {
	activities := make([]NPRAActivities, 0)
	DatabaseConnection.Where("APPROVED = ?", false).Preload("User").Preload("Client").Order("ID desc").Find(&activities)
	return activities
}

// LoadTrusteeRecentActivities loads the trustee recent activities from the database
func LoadTrusteeRecentActivities() []TrusteeActivities {
	activities := make([]TrusteeActivities, 0)
	DatabaseConnection.Where("APPROVED = ?", false).Preload("User").Preload("Client").Preload("SCAClient").Order("ID desc").Find(&activities)
	return activities
}

// LoadBillingRecentActivities loads the sec recent activities from the database
func LoadBillingRecentActivities() []BillingActivities {
	activities := make([]BillingActivities, 0)
	DatabaseConnection.Where("APPROVED = ?", false).Preload("User").Preload("Client").Order("ID desc").Find(&activities)
	return activities
}

// LogRecentActivity logs sec recent activities
func LogRecentActivity(UserId int, BPID, Activity string) {
	date := MakeDBQueryableCurrentQuarterDate()
	activity := SecActivities{
		BPID:        BPID,
		UserID:      UserId,
		Activity:    Activity,
		Date:        time.Now(),
		QuarterDate: date,
		Approved:    false,
	}
	DatabaseConnection.Where(SecActivities{
		BPID:        BPID,
		Activity:    Activity,
		QuarterDate: date,
	}).Assign(activity).FirstOrCreate(&activity)
}

// LoadSchemeFromActivity loads the details of the scheme from the activities
func LoadSchemeFromActivity(BPID, QuarterDate string) SchemeDetails {
	//TODO: check if the activity has been approved
	scheme := SchemeDetails{}
	parsedDate, err := time.Parse("2006-01-02", QuarterDate)
	if err == nil {
		quarter := GetQuarterNumber(parsedDate)
		year := parsedDate.Year()
		date := MakeQuarterDate(fmt.Sprintf("%d", quarter), fmt.Sprintf("%d", year))
		governance := LoadGovernanceInfoByQuarterDate(date)
		if governance.ID > 0 {
			DatabaseConnection.Where("GOV_INFO_ID=? AND BP_ID=?", governance.ID, BPID).Find(&scheme)
		}
	}
	return scheme
}

// GetClientByBPID loads a client by their BPID
func GetClientByBPID(BPID string) Client {
	client := Client{}
	DatabaseConnection.Where("BP_ID=?", BPID).Find(&client)
	return client
}

func GetClientByCode(code string) ClientPVUniqueCode {
	client := ClientPVUniqueCode{}
	DatabaseConnection.Where("CODE=?", code).Find(&client)
	return client
}

// GetClientSCASByBPID loads a client by their BPID
func GetClientSCASByBPID(BPID string) []ClientPVUniqueCode {
	scas := make([]ClientPVUniqueCode, 0)
	DatabaseConnection.Where("BP_ID=?", BPID).Find(&scas)
	return scas
}

// LoadPVReportDataByQuarter loads the pv report data for a given date by the bpid
func LoadPVReportDataByQuarter(ClientID, Date, pvType string) []PVReportField {
	data := make([]PVReportField, 0)
	//DatabaseConnection.Where("CLIENT_ID=? AND REPORT_DATE=? AND CLIENT_TYPE=?", ClientID, Date, pvType).Order("CREATED_AT ASC").Find(&data)
	DatabaseConnection.Raw("select * from CRS_PVREPORTS where CLIENT_ID=? and CLIENT_TYPE=? and REPORT_DATE=? and BP_ID = (select TOP(1) BP_ID from CRS_PVSUMMARY where CLIENT_ID=? and REPORT_DATE=? and CLIENT_TYPE=?) order by CREATED_AT ASC", ClientID, pvType, Date, ClientID, Date, pvType).Scan(&data)
	return data
}

// LoadPVSummaryDataByQuarter loads the pv summary data for a given date by the bpid
func LoadPVSummaryDataByQuarter(ClientID, Date, pvType string) []ClientPVReportSummary {
	data := make([]ClientPVReportSummary, 0)
	//DatabaseConnection.Where("CLIENT_ID=? AND REPORT_DATE=? AND CLIENT_TYPE=?", ClientID, Date, pvType).Find(&data)
	DatabaseConnection.Raw("select * from CRS_PVSUMMARY where CLIENT_ID=? and CLIENT_TYPE=? and REPORT_DATE=? and BP_ID = (select TOP(1) BP_ID from CRS_PVSUMMARY where CLIENT_ID=? and REPORT_DATE=? and CLIENT_TYPE=?)", ClientID, pvType, Date, ClientID, Date, pvType).Scan(&data)
	return data
}

func LoadTrusteeQuarterlyPerformance(bpidOrSca, quarter, year string) []TrusteePVPerformance {
	data := make([]TrusteePVPerformance, 0)
	currentDate := MakeQuarterFormalDate(quarter, year)
	previousDate := MakeLastQuarterFormalDate(currentDate)
	DatabaseConnection.Raw("select curr.SECURITY_TYPE, SUM(curr.LCY_AMT) as CURRENT_QUARTER,(select SUM(LCY_AMT) from CRS_PVSUMMARY where (BP_ID = ? OR CLIENT_ID=?) AND SECURITY_TYPE = curr.SECURITY_TYPE AND REPORT_DATE=?) as PREVIOUS_QUARTER from CRS_PVSUMMARY curr where (BP_ID=? OR CLIENT_ID=?) AND REPORT_DATE=? and CLIENT_TYPE=? group by SECURITY_TYPE", bpidOrSca, bpidOrSca, previousDate, bpidOrSca, bpidOrSca, currentDate, "trustee").Scan(&data)
	return data
}

func AddClientMonthlyContribution(request ClientMonthlyContributionRequest) ClientMonthlyContribution {
	date, _ := time.Parse("2006-01-02", request.Date)
	contribution := ClientMonthlyContribution{
		Quarter: MakeDBQueryableQuarterLastDate(request.Quarter, request.Year),
		BPID:    request.BPID,
		SCA:     request.SCA,
		Date:    date,
		Amount:  request.Amount,
	}
	DatabaseConnection.Create(&contribution)
	return contribution
}

func EditClientMonthlyContribution(request ClientMonthlyContributionEditRequest) ClientMonthlyContribution {
	date, _ := time.Parse("2006-01-02", request.Date)
	contribution := ClientMonthlyContribution{
		//Quarter: MakeQuarterDate(request.Quarter, request.Year),
		BPID:   request.BPID,
		SCA:    request.SCA,
		Date:   date,
		Amount: request.Amount,
	}
	DatabaseConnection.Where(ClientMonthlyContribution{ID: request.ID}).Assign(contribution).FirstOrCreate(&contribution)
	return contribution
}

func LoadMonthlyContributions(bpidOrSca, quarter, year string) []ClientMonthlyContribution {
	data := make([]ClientMonthlyContribution, 0)
	DatabaseConnection.Where("QUARTER= ?", MakeQuarterFormalDate(quarter, year)).Where("BP_ID=? OR SCA=?", bpidOrSca, bpidOrSca).Find(&data)
	return data
}

func LoadFundManagers() []FundManager {
	data := make([]FundManager, 0)
	DatabaseConnection.Order("NAME asc").Find(&data)
	return data
}

func LoadFundAdministrators(notClosedAsAt time.Time) []FundAdministrators {
	data := make([]FundAdministrators, 0)
	DatabaseConnection.Raw("SELECT DISTINCT(CLIENT_NAME) as NAME,ADMINISTRATOR FROM CRS_CLIENT_PV_UNIQUE_CODES WHERE ADMINISTRATOR NOT IN ('','N/A') and ( CLOSED is null or CLOSED > ?)", notClosedAsAt).Scan(&data)
	return data
}

// AddClientToDB adds a client to the database
func AddClientToDB(request ClientAddRequest) bool {
	if err :=
		DatabaseConnection.Create(&Client{
			Name:               request.ClientName,
			Location:           request.Location,
			Phone:              request.PhoneNumber,
			AccountNumber:      request.SafekeepingAccount,
			Type:               request.ClientType,
			Tier:               request.ClientTier,
			FundManagerID:      request.FundManagerID,
			ContactPerson:      request.ContactPerson,
			ContactPersonPhone: request.ContactPersonNumber,
			HomeCountry:        request.HomeCountry,
			BPID:               request.BPID,
			Email:              request.Email,
			AddressLine1:       request.AddressLine1,
			AddressLine2:       request.AddressLine2,
			AddressLine3:       request.AddressLine3,
			AddressLine4:       request.AddressLine4,
			Image:              request.Image,
			SchemeManager:      request.SchemeManager,
		}); err == nil {
		return true
	}
	return false
}

// EditClient edits a client in the database
func EditClient(request ClientEditRequest) {
	DatabaseConnection.Model(&Client{}).Where("ID=?", request.ID).Updates(map[string]interface{}{
		"CLIENTNAME":       request.ClientName,
		"LOCATION":         request.Location,
		"TELEPHONENO":      request.PhoneNumber,
		"SAFEKEEPINGACCNO": request.SafekeepingAccount,
		"CLIENT_TYPE":      request.ClientType,
		"CLIENT_TIER":      request.ClientTier,
		"FUND_MANAGER_ID":  request.FundManagerID,
		"CONTACTPERSON":    request.ContactPerson,
		"CONTACTPERSONNO":  request.ContactPersonNumber,
		"HOME_COUNTRY":     request.HomeCountry,
		"BP_ID":            request.BPID,
		"EMAIL":            request.Email,
		"ADDRESS_LINE1":    request.AddressLine1,
		"ADDRESS_LINE2":    request.AddressLine2,
		"ADDRESS_LINE3":    request.AddressLine3,
		"ADDRESS_LINE4":    request.AddressLine4,
		"CLIENT_IMAGE":     request.Image,
		"SCHEME_MANAGER":   request.SchemeManager,
	})
}

// AddNPRA0303ToDB adds a client to the database
func AddNPRA0303ToDB(request NPRA0303AddRequest) bool {
	if err :=
		DatabaseConnection.Create(&NPRA0303{
			ReportCode:          request.ReportCode,
			EntityID:            request.EntityID,
			EntityName:          request.EntityName,
			ReferencePeriodYear: request.ReferencePeriodYear,
			ReferencePeriod:     request.ReferencePeriod,
			UnitPrice:           request.UnitPrice,
			DateValuation:       request.DateValuation,
			UnitNumber:          request.UnitNumber,
			DailyNav:            request.DailyNav,
			NPRAFees:            request.NPRAFees,
			TrusteeFees:         request.TrusteeFees,
			FundManagerFees:     request.FundManagerFees,
			FundCustodianFees:   request.FundCustodianFees,
			CreatedBy:           request.CreatedBy,
		}); err == nil {
		return true
	}
	return false
}

// AddAsetClassToDB adds a client to the database
func AddAssetClassToDB(request AssetClassAddRequest) bool {
	if err :=
		DatabaseConnection.Create(&AssetClass{
			SecurityType:   request.SecurityType,
			AssetClassName: request.AssetClassName,
		}); err == nil {
		return true
	}
	return false
}

func LoadNPRA0303Reports() []NPRA0303 {
	data := make([]NPRA0303, 0)
	DatabaseConnection.Order("created_at asc").Find(&data)
	return data
}
func LoadAssetClassReports() []AssetClass {
	data := make([]AssetClass, 0)
	DatabaseConnection.Order("created_at asc").Find(&data)
	return data
}

func LoadNPRA03021Reports(bpOrSca, date string) []NPRA0302 {
	data := make([]NPRA0302, 0)
	DatabaseConnection.Where("bp_id = ? or client_code=? and reporting_date=?", bpOrSca, bpOrSca, date).Order("created_at asc").Find(&data)
	return data
}
func LoadNPRA03011Reports(date string) []NPRA0301 {
	data := make([]NPRA0301, 0)
	DatabaseConnection.Where("reporting_date=?", date).Order("created_at asc").Find(&data)
	return data
}

// // LoadAllnpra303 loads all LoadAllnpra303 from the database
// func LoadAllnpra303() []NPRA0303 {
// 	return LoadNnpra303(-1)
// }

// // LoadNnpra303 loads a given number of clients
// func LoadNnpra303(limit int) []NPRA0303 {
// 	npra303s := make([]Client, 0)
// 	if limit == -1 {
// 		DatabaseConnection.Order("CLIENTNAME asc").Find(&npra303s)
// 	} else {
// 		DatabaseConnection.Order("CLIENTNAME asc").Limit(limit).Find(&npra303s)
// 	}
// 	return npra303s
// }

// EditNPRA0303 edits a client in the database
func EditNPRA0303(request NPRA0303EditRequest) {
	DatabaseConnection.Model(&NPRA0303{}).Where("ID=?", request.ID).Updates(map[string]interface{}{
		"report_code":           request.ReportCode,
		"entity_id":             request.EntityID,
		"entity_name":           request.EntityName,
		"reference_period_year": request.ReferencePeriodYear,
		"reference_period":      request.ReferencePeriod,
		"unit_price":            request.UnitPrice,
		"date_valuation":        request.DateValuation,
		"unit_number":           request.UnitNumber,
		"daily_nav":             request.DailyNav,
		"npra_fees":             request.NPRAFees,
		"trustee_fees":          request.TrusteeFees,
		"fund_manager_fees":     request.FundManagerFees,
		"fund_custodian_fees":   request.FundCustodianFees,
		"created_by":            request.CreatedBy,
	})
}

// AddFundManagerToDB adds a fund manager to the database
func AddFundManagerToDB(request FundManagerAddRequest) {
	DatabaseConnection.Create(&FundManager{
		Name:          request.Name,
		AccountNumber: request.AccountNumber,
		Administrator: request.Administrator,
	})
}

// GetFundManagerByID returns the fund manager by id
func GetFundManagerByID(id int) FundManager {
	manager := FundManager{}
	DatabaseConnection.Where("ID=?", id).First(&manager)
	return manager
}

// EditFundManager returns the fund manager by id
func EditFundManager(request FundManagerEditRequest) {
	DatabaseConnection.Model(&FundManager{}).Where("ID=?", request.ID).Updates(FundManager{
		Name:          request.Name,
		AccountNumber: request.AccountNumber,
		Administrator: request.Administrator,
	})
}

// StoreUnidentifiedPayments uploads unidentified transactions
func StoreUnidentifiedPayments(payments []UnidentifiedPaymentRequest) {
	for _, txn := range payments {
		txnDate, _ := time.Parse("2006-01-02", txn.TransactionDate)
		valueDate, _ := time.Parse("2006-01-02", txn.ValueDate)
		DatabaseConnection.Create(&UnidentifiedPayment{
			ClientBPID:        txn.ClientBPID,
			FundManagerID:     txn.FundManager,
			TransactionDate:   txnDate,
			TransactionType:   txn.TransactionType,
			ValueDate:         valueDate,
			NameOfCompany:     txn.NameOfCompany,
			Amount:            txn.Amount,
			CollectionAccount: txn.CollectionAccountNumber,
			Status:            txn.Status,
		})
	}
}

// LoadUnidentifiedPayments loads all unidentified payments for the client
func LoadUnidentifiedPayments(bpid, year string) []UnidentifiedPayment {
	payments := make([]UnidentifiedPayment, 0)
	DatabaseConnection.Where("CLIENT_BPID=? AND TRANSACTION_DATE BETWEEN ? AND ?", bpid, year+"-01-01", year+"-12-01").Preload("FundManager").Find(&payments)
	return payments
}

// LoadUnidentifiedPaymentSummary loads the client unidentified payment history
func LoadUnidentifiedPaymentSummary(bpid, year string) []UnidentifiedPaymentSummary {
	data := make([]UnidentifiedPaymentSummary, 0)
	DatabaseConnection.Raw(`
	select 
		TRANSACTION_TYPE,
		count(*) as TOTAL,
		(select 
			count(*)
		from 
			CRS_TRUSTEE_UNIDENTIFIED_PAYMENTS
		where 
			CLIENT_BPID=?
		AND
			TRANSACTION_DATE BETWEEN ? AND ?
		AND
			[STATUS] = 'pending'
		AND
			TRANSACTION_TYPE=main.TRANSACTION_TYPE
		) as PENDING,
		(select 
			count(*)
		from
			CRS_TRUSTEE_UNIDENTIFIED_PAYMENTS
		where 
			CLIENT_BPID=?
		AND
			TRANSACTION_DATE BETWEEN ? AND ?
		AND
			[STATUS] = 'done'
		AND 
			TRANSACTION_TYPE=main.TRANSACTION_TYPE
		) as DONE 
	from 
		CRS_TRUSTEE_UNIDENTIFIED_PAYMENTS main
	where 
		CLIENT_BPID=?
	AND
		TRANSACTION_DATE BETWEEN ? AND ?
	group by
		TRANSACTION_TYPE
`,
		bpid, year+"-01-01", year+"-12-01",
		bpid, year+"-01-01", year+"-12-01",
		bpid, year+"-01-01", year+"-12-01",
	).Scan(&data)
	return data
}

// LoadAllLocalClients loads all local clients
func LoadAllLocalClients() []Client {
	clients := make([]Client, 0)
	DatabaseConnection.Model(Client{}).Where("CLIENT_TYPE <> 'FOREIGN CLIENT'").Find(&clients)
	return clients
}

// LoadPVReportHeadings loads all pv report headings from the DB
func LoadPVReportHeadings() []PVReportHeadings {
	headings := make([]PVReportHeadings, 0)
	DatabaseConnection.Find(&headings)
	return headings
}

// CalculatePercentageCapitalMarketValue calculates the market value based on the bpid
func CalculatePercentageCapitalMarketValue(bpid string) (equities, investments, percentage float64) {
	equities = LoadDataForEquityCalculations(bpid)
	investments = LoadDataForTotalFixedIncomeInvestmentCalculations(bpid)
	capitalInvestments := LoadDataForCapitalInvestmentCalculations(bpid)
	aucWithoutCashBalance := ClientAUCWithoutCashBalance(bpid)
	if aucWithoutCashBalance > 0 {
		percentage = (capitalInvestments / aucWithoutCashBalance) * 100
	} else {
		percentage = 0
	}
	return equities, investments, percentage
}

func LoadDataForEquityCalculations(bpid string) float64 {
	type RetVal struct {
		Amount float64 `gorm:"column:amount"`
	}
	ret := RetVal{}
	DatabaseConnection.Raw(fmt.Sprintf("select sum(LCY_AMT) as amount from CRS_PVSUMMARY where SECURITY_TYPE IN %s AND BP_ID= ? AND REPORT_DATE=? AND CLIENT_TYPE='sec'", MakeDataForEquityCalculations()), bpid, GetQuarterFormalDate()).Scan(&ret)
	return ret.Amount
}

func LoadDataForTotalFixedIncomeInvestmentCalculations(bpid string) float64 {
	type RetVal struct {
		Amount float64 `gorm:"column:amount"`
	}
	ret := RetVal{}
	DatabaseConnection.Raw(fmt.Sprintf("select sum(LCY_AMT) as amount from CRS_PVSUMMARY where (SECURITY_TYPE IN %s OR SECURITY_TYPE LIKE '%%RECEI%%') AND BP_ID= ? AND REPORT_DATE=? AND CLIENT_TYPE='sec'", MakeDataForTotalFixedIncomeInvestmentCalculations()), bpid, GetQuarterFormalDate()).Scan(&ret)
	return ret.Amount
}

func LoadDataForCapitalInvestmentCalculations(bpid string) float64 {
	type RetVal struct {
		Amount float64 `gorm:"column:amount"`
	}
	ret := RetVal{}
	DatabaseConnection.Raw(fmt.Sprintf("select sum(LCY_AMT) as amount from CRS_PVSUMMARY where SECURITY_TYPE IN %s AND BP_ID= ? AND REPORT_DATE=? AND CLIENT_TYPE='sec'", MakeDataForCapitalInvestmentCalculations()), bpid, GetQuarterFormalDate()).Scan(&ret)
	return ret.Amount
}

func ClientAUCWithoutCashBalance(bpid string) float64 {
	type RetVal struct {
		Amount float64 `gorm:"column:amount"`
	}
	ret := RetVal{}
	DatabaseConnection.Raw("select sum(LCY_AMT) as amount from CRS_PVSUMMARY where SECURITY_TYPE NOT LIKE '%CASH BAL%' AND BP_ID= ? AND REPORT_DATE=?", bpid, GetQuarterFormalDate()).Scan(&ret)
	return ret.Amount
}

func LoadPensionFund(quarter, year string, shouldAddTotal bool) []NPRAFund {
	fund := make([]NPRAFund, 0)
	clients := LoadTier2Clients(MakeDBQueryableQuarterLastDate(quarter, year))
	var total float64
	for _, client := range clients {
		value := ValueQuery{}
		if client.HasMultipleSCA {
			DatabaseConnection.Raw("SELECT SUM(LCY_AMT) as value FROM CRS_PVSUMMARY WHERE BP_ID= ? AND REPORT_DATE=?", client.BPID, MakeQuarterFormalDate(quarter, year)).Scan(&value)
		} else {
			DatabaseConnection.Raw("SELECT SUM(LCY_AMT) as value FROM CRS_PVSUMMARY WHERE CLIENT_ID= ? AND REPORT_DATE=?", client.SCA[0], MakeQuarterFormalDate(quarter, year)).Scan(&value)
		}
		//ERMM
		code := client.BPID
		if client.HasMultipleSCA {
			code = client.BPID
		} else {
			if len(client.SCA) > 0 {
				code = client.SCA[0]
			}
		}
		fund = append(fund, NPRAFund{Value: value.Value, Name: client.Client, BPID: code})
		total += value.Value
	}
	if shouldAddTotal {
		fund = append(fund, NPRAFund{Name: "TOTAL(GHS)", Value: total})
	}
	return fund
}

func LoadProvidentFund(quarter, year string, shouldAddTotal bool) []NPRAFund {
	fund := make([]NPRAFund, 0)
	clients := LoadTier3Clients(MakeDBQueryableQuarterLastDate(quarter, year))
	var total float64
	for _, client := range clients {
		value := ValueQuery{}
		if client.HasMultipleSCA {
			DatabaseConnection.Raw("SELECT SUM(LCY_AMT) as value FROM CRS_PVSUMMARY WHERE BP_ID= ? AND REPORT_DATE=?", client.BPID, MakeQuarterFormalDate(quarter, year)).Scan(&value)
		} else {
			DatabaseConnection.Raw("SELECT SUM(LCY_AMT) as value FROM CRS_PVSUMMARY WHERE CLIENT_ID= ? AND REPORT_DATE=?", client.SCA[0], MakeQuarterFormalDate(quarter, year)).Scan(&value)
		}
		//ERMMM
		code := client.BPID
		if client.HasMultipleSCA {
			code = client.BPID
		} else {
			if len(client.SCA) > 0 {
				code = client.SCA[0]
			}
		}
		fund = append(fund, NPRAFund{Value: value.Value, Name: client.Client, BPID: code})
		total += value.Value
	}
	if shouldAddTotal {
		fund = append(fund, NPRAFund{Name: "TOTAL(GHS)", Value: total})
	}
	return fund
}

func LoadTier2Clients(notClosedAsAt time.Time) []ClientMergedBPAndSCA {
	clients := make([]ClientPVUniqueCode, 0)
	DatabaseConnection.Raw("select distinct(GROUP_KEY),CLIENT_NAME,BP_ID,CLIENT_TIER,CODE from CRS_CLIENT_PV_UNIQUE_CODES where CLIENT_TIER =? and( CLOSED is null or CLOSED > ?)", 2, notClosedAsAt).Scan(&clients)
	merged := make([]ClientMergedBPAndSCA, 0)
	hackForMerged := make(map[string]*ClientMergedBPAndSCA, 0)
	for _, client := range clients {
		if hackForMerged[client.GroupKey] == nil {
			scas := make([]string, 0)
			scas = append(scas, client.Code)
			hackForMerged[client.GroupKey] = &ClientMergedBPAndSCA{
				Client:         client.GroupKey,
				SCA:            scas,
				BPID:           client.BPID,
				HasMultipleSCA: false,
			}
		} else {
			hackForMerged[client.GroupKey].HasMultipleSCA = true
			hackForMerged[client.GroupKey].SCA = append(hackForMerged[client.GroupKey].SCA, client.Code)
		}
	}
	for _, hack := range hackForMerged {
		merged = append(merged, *hack)
	}
	return merged
}

func LoadTier3Clients(notClosedAsAt time.Time) []ClientMergedBPAndSCA {
	clients := make([]ClientPVUniqueCode, 0)
	DatabaseConnection.Raw("select distinct(GROUP_KEY),CLIENT_NAME,BP_ID,CLIENT_TIER,CODE from CRS_CLIENT_PV_UNIQUE_CODES where CLIENT_TIER =? and( CLOSED is null or CLOSED > ?)", 3, notClosedAsAt).Scan(&clients)
	merged := make([]ClientMergedBPAndSCA, 0)
	hackForMerged := make(map[string]*ClientMergedBPAndSCA, 0)
	for _, client := range clients {
		if hackForMerged[client.GroupKey] == nil {
			scas := make([]string, 0)
			scas = append(scas, client.Code)
			hackForMerged[client.GroupKey] = &ClientMergedBPAndSCA{
				Client:         client.GroupKey,
				SCA:            scas,
				BPID:           client.BPID,
				HasMultipleSCA: false,
			}
		} else {
			hackForMerged[client.GroupKey].HasMultipleSCA = true
			hackForMerged[client.GroupKey].SCA = append(hackForMerged[client.GroupKey].SCA, client.Code)
		}
	}
	for _, hack := range hackForMerged {
		merged = append(merged, *hack)
	}
	return merged
}

func LoadSumTotalOfPensionFundByQuarter(date time.Time) float64 {
	value := ValueQuery{}
	DatabaseConnection.Raw(fmt.Sprintf("SELECT SUM(LCY_AMT) as value FROM CRS_PVSUMMARY WHERE CLIENT_ID IN %s AND REPORT_DATE=?", MakeWhereQueryDataUsingClientCodeFromClients(LoadTier2Clients(date))), date.Format("2006-01-02")).Scan(&value)
	return value.Value
}

func LoadSumTotalOfProvidentFundByQuarter(date time.Time) float64 {
	value := ValueQuery{}
	DatabaseConnection.Raw(fmt.Sprintf("SELECT SUM(LCY_AMT) as value FROM CRS_PVSUMMARY WHERE CLIENT_ID IN %s AND REPORT_DATE=?", MakeWhereQueryDataUsingClientCodeFromClients(LoadTier3Clients(date))), date.Format("2006-01-02")).Scan(&value)
	return value.Value
}

// LoadNPRATrendsForPastFourQuarters loads the trends for NPRA for the past 4 quarters
func LoadNPRATrendsForPastFourQuarters() []NPRATrend {
	dates := GetLast4QuarterDates()
	trends := make([]NPRATrend, 0)
	for _, date := range dates {
		trends = append(trends, NPRATrend{
			Pension:   LoadSumTotalOfPensionFundByQuarter(date),
			Provident: LoadSumTotalOfProvidentFundByQuarter(date),
		})
	}
	return trends
}

func InsertUnauthorizedTransactionsIntoDB(transactions []NPRAUnauthorizedTransactionRequest) {
	for _, transaction := range transactions {
		date, _ := time.Parse("2006-01-02", transaction.Date)
		data := NPRAUnauthorizedTransaction{
			ClientName:         transaction.ClientName,
			TransactionDetails: transaction.TransactionDetails,
			Date:               date,
		}
		DatabaseConnection.Create(&data)
	}
}

func InsertOutstandingFDCertificatesIntoDB(certificates []NPRAOutstandingFDCertificateRequest, userId int) {
	hash := string(RandomBytes(80))
	for _, certificate := range certificates {
		data := NPRAOutstandingFDCertificate{
			FundManager:   certificate.FundManager,
			ClientName:    certificate.ClientName,
			Amount:        certificate.Amount,
			Issuer:        certificate.Issuer,
			Rate:          certificate.Rate,
			Tenor:         certificate.Tenor,
			Term:          certificate.Term,
			EffectiveDate: certificate.EffectiveDate,
			Maturity:      certificate.Maturity,
			Hash:          hash,
		}
		DatabaseConnection.Create(&data)
	}
	DatabaseConnection.Create(&NPRAActivities{
		Hash:        hash,
		UserID:      userId,
		Date:        time.Now(),
		QuarterDate: MakeDBQueryableCurrentQuarterDate(),
		Activity:    UploadOutstandingFDReport,
	})
}

func LoadUnauthorizedTransactions(quarter, year string) []NPRAUnauthorizedTransaction {
	transactions := make([]NPRAUnauthorizedTransaction, 0)
	DatabaseConnection.Where("CREATED_AT BETWEEN ? AND ?", MakeDBQueryableQuarterFirstDate(quarter, year), MakeDBQueryableQuarterLastDate(quarter, year)).Order("CREATED_AT desc").Find(&transactions)
	return transactions
}

func LoadOutstandingFDCertificates(quarter, year string) []NPRAOutstandingFDCertificate {
	certificates := make([]NPRAOutstandingFDCertificate, 0)
	DatabaseConnection.Where("EFFECTIVE_DATE BETWEEN ? AND ?", MakeDBQueryableQuarterFirstDate(quarter, year), MakeDBQueryableQuarterLastDate(quarter, year)).Order("EFFECTIVE_DATE desc").Find(&certificates)
	return certificates
}

func LoadOutstandingFDCertificatesByHash(hash string) []NPRAOutstandingFDCertificate {
	certificates := make([]NPRAOutstandingFDCertificate, 0)
	DatabaseConnection.Where("HASH= ?", hash).Order("EFFECTIVE_DATE desc").Find(&certificates)
	return certificates
}

func LoadClientAUCForPastMonths(bpidOrSca, quarter, year string) []float64 {
	values := make([]float64, 0)
	dates := MakeTrusteeDates(quarter, year)
	for _, date := range dates {
		value := ValueQuery{}
		DatabaseConnection.Raw("select sum(LCY_AMT) as value from CRS_PVSUMMARY WHERE (BP_ID=? OR CLIENT_ID=?)AND REPORT_DATE=? AND CLIENT_TYPE=?", bpidOrSca, bpidOrSca, now.New(date).EndOfMonth().Format("2006-01-02"), "trustee").Scan(&value)
		values = append(values, value.Value)
	}
	return values
}

func LoadClientTransactionVolumesForPastMonths(bpid, quarter, year string) []float64 {
	values := make([]float64, 0)
	dates := MakeTrusteeDatesRelativeToCurrentQuarter(quarter, year)
	quarterDate := MakeDBQueryableQuarterLastDate(quarter, year)
	for _, date := range dates {
		value := ValueQuery{}
		DatabaseConnection.Raw("select COUNT(*) as value from CRS_TRUST_TXN_VOLUMES WHERE BP_ID=? AND Stock_Settled_Date BETWEEN ? AND ? AND QUARTER_DATE=?", bpid, now.New(date).BeginningOfMonth(), now.New(date).EndOfMonth(), quarterDate).Scan(&value)
		values = append(values, value.Value)
	}
	return values
}

func LoadClientTransactionVolumesByAssetClassForPastMonths(bpid, quarter, year string) ([]TradeVolumeByAssetClassSummary, int64) {
	dates := MakeTrusteeDatesRelativeToCurrentQuarter(quarter, year)
	quarterDate := MakeDBQueryableQuarterLastDate(quarter, year)
	assetClasses := make([]StrValueQuery, 0)
	data := make(map[string][]TradeVolumeByAssetClass, 0)
	DatabaseConnection.Raw("select Security_Type as value from CRS_TRUST_TXN_VOLUMES WHERE BP_ID=? AND QUARTER_DATE= ? GROUP BY Security_Type", bpid, quarterDate).Scan(&assetClasses)
	hackFormatsForSorting := make([]string, 0)
	totalNumberOfTransactions := int64(0)
	for _, date := range dates {
		var govtBonds int64
		format := date.Format("Jan")
		hackFormatsForSorting = append(hackFormatsForSorting, format)
		for _, asset := range assetClasses {
			value := IntValueQuery{}
			DatabaseConnection.Raw("select COUNT(*) as value from CRS_TRUST_TXN_VOLUMES WHERE BP_ID=? AND Stock_Settled_Date BETWEEN ? AND ? AND Security_Type=? AND QUARTER_DATE=?", bpid, now.New(date).BeginningOfMonth(), now.New(date).EndOfMonth(), asset.Value, quarterDate).Scan(&value)
			totalNumberOfTransactions += value.Value
			if asset.Value == "Gov Bonds" || asset.Value == "Gov Notes" {
				govtBonds += value.Value
			} else {
				if data[format] == nil {
					data[format] = make([]TradeVolumeByAssetClass, 0)
				}
				data[format] = append(data[format], TradeVolumeByAssetClass{Number: value.Value, Asset: asset.Value})
			}
		}
		data[format] = append(data[format], TradeVolumeByAssetClass{Number: govtBonds, Asset: "Gov Bonds"})
	}

	sortedData := make([]TradeVolumeByAssetClassSummary, 0)
	for _, month := range hackFormatsForSorting {
		sortedData = append(sortedData, TradeVolumeByAssetClassSummary{
			Month: month,
			Data:  GroupTradeVolumesByAssetClass(data[month]),
		})
	}
	return sortedData, totalNumberOfTransactions
}

func MakeTransactionVolumeByAssetClassSummary(data []TradeVolumeByAssetClassSummary) map[string]string {
	var totalCount int64
	hackForAssetData := make(map[string]int64, 0)
	summary := make(map[string]string, 0)
	for _, each := range data {
		for _, asset := range each.Data {
			totalCount += asset.Number
			hackForAssetData[asset.Asset] += asset.Number
		}
	}
	for asset, value := range hackForAssetData {
		result := (float32(value) / float32(totalCount)) * 100
		summary[asset] = fmt.Sprintf("%.2f %%", result)
	}
	return summary
}

func UploadTransactionDetails(transactions []BillingTransactionDetails) {
	for _, transaction := range transactions {
		DatabaseConnection.Create(&transaction)
	}
}

func LoadTransactionDetails(bpOrSca string, date time.Time) []BillingTransactionDetails {
	transactions := make([]BillingTransactionDetails, 0)
	DatabaseConnection.Where("BP_ID= ? OR SCA = ? ", bpOrSca, bpOrSca).Where(" CHARGE_ITEM='TRANSACTION FEE' AND REPORTING_DATE=?", date).Find(&transactions)
	return transactions
}

func UploadCurrencyDetails(details BillingCurrencyDetails) {
	DatabaseConnection.Create(&details)
}

func LoadCurrencyDetails(bporSca string, date time.Time) []BillingCurrencyDetails {
	currencyDetails := make([]BillingCurrencyDetails, 0)
	DatabaseConnection.Where("BP_ID= ? OR SCA=? ", bporSca, bporSca).Where(" REPORTING_DATE= ?", date).Find(&currencyDetails)
	return currencyDetails
}

func LoadClientBillingInfo(bpid string) ClientBillingInfo {
	info := ClientBillingInfo{}
	DatabaseConnection.Where("BP_ID=?", bpid).First(&info)
	return info
}

func LoadClientBasisPoints(bpOrSca string) []ClientBasisPoints {
	points := make([]ClientBasisPoints, 0)
	DatabaseConnection.Where("BP_ID=? OR SCA=?", bpOrSca, bpOrSca).Find(&points)
	return points
}

func UpdateClientBillingInfo(request ClientBillingInfoUpdateRequest) {
	billing := ClientBillingInfo{
		BPID:                 request.BPID,
		MinimumCharge:        request.MinimumCharge,
		ChargePerTransaction: request.ChargePerTransaction,
		ThirdPartyTransfer:   request.ThirdPartyTransfer,
	}
	DatabaseConnection.Where(ClientBillingInfo{BPID: request.BPID}).Assign(billing).FirstOrCreate(&billing)
	//TODO: ... this is terrible but :(
	DatabaseConnection.Where("BP_ID=?", request.BPID).Delete(&ClientBasisPoints{})
	for _, point := range request.BasisPoints {
		DatabaseConnection.Create(&ClientBasisPoints{
			BPID:        request.BPID,
			Maximum:     point.Maximum,
			Minimum:     point.Minimum,
			BasisPoints: point.BasisPoints,
		})
	}
}

func IsSecClient(bpOrSca string) bool {
	/*TODO: use sca for clients with sca*/
	return GetClientByBPID(bpOrSca).Type == "SEC"
}

func CalculateClientInvoiceSummary(bpOrSca string, date time.Time) []ClientInvoiceSummary {
	summary := make([]ClientInvoiceSummary, 3)
	portfolioData := GetClientNAVDetails(bpOrSca, date)
	points := make([]ClientBasisPoints, 0)
	var basisPoint float64
	DatabaseConnection.Where("BP_ID=? OR SCA=?", bpOrSca, bpOrSca).Find(&points)
	for _, bp := range points {
		if bp.Maximum == 0.0 {
			if portfolioData.NAV >= bp.Minimum {
				basisPoint = bp.BasisPoints
				break
			}
		} else {
			if portfolioData.NAV >= bp.Minimum && portfolioData.NAV <= bp.Maximum {
				basisPoint = bp.BasisPoints
				break
			}
		}
	}
	txns := make([]BillingTransactionDetails, 0)
	DatabaseConnection.Where("BP_ID= ? OR SCA = ?", bpOrSca, bpOrSca).Where("CHARGE_ITEM=? AND REPORTING_DATE= ?", "TRANSACTION FEE", date.Format("2006-01-02")).Find(&txns)

	billingInfo := LoadClientBillingInfo(bpOrSca)

	tradesChargeAmount := float64(float32(len(txns)) * billingInfo.ChargePerTransaction)

	portfolioChargeAmount := (portfolioData.NAV * basisPoint) / 120000
	if IsSecClient(bpOrSca) {
		pv5Percent := (portfolioChargeAmount * 5) / 100
		initialGross := pv5Percent + portfolioChargeAmount
		pv12Point5Percent := ((pv5Percent + portfolioChargeAmount) * 12.5) / 100
		tax := pv12Point5Percent + initialGross
		/*Total trades amount is minus the tax...*/
		totalChargeAmount := portfolioChargeAmount + tradesChargeAmount
		total5Percent := (totalChargeAmount * 5) / 100
		totalInitialGross := total5Percent + totalChargeAmount
		total12Point5Percent := (totalInitialGross * 12.5) / 100
		totalTaxAmount := totalInitialGross + total12Point5Percent

		trades5Percent := (tradesChargeAmount * 5) / 100
		tradesInitialGross := trades5Percent + tradesChargeAmount
		trades12Point5Percent := (tradesInitialGross * 12.5) / 100
		tradesChargeAmountWithTax := trades12Point5Percent + tradesInitialGross
		summary[0] = ClientInvoiceSummary{
			ChargeType:             "Portfolio Fee",
			ChargeableQuantity:     portfolioData.NAV,
			ChargeAmount:           portfolioChargeAmount,
			TaxAmount:              pv12Point5Percent + pv5Percent,
			ChargeAmountWithTax:    tax,
			InvoiceAmountWithTax:   tax,
			BasisPoint:             basisPoint,
			FivePercent:            pv5Percent,
			TwelvePointFivePercent: pv12Point5Percent,
			InitialGross:           initialGross,
		}

		summary[1] = ClientInvoiceSummary{
			ChargeType:             "Transaction Fee",
			ChargeableQuantity:     float64(len(txns)),
			ChargeAmount:           tradesChargeAmount,
			TaxAmount:              trades12Point5Percent + trades5Percent,
			ChargeAmountWithTax:    tradesChargeAmountWithTax,
			InvoiceAmountWithTax:   tradesChargeAmountWithTax,
			BasisPoint:             float64(billingInfo.ChargePerTransaction),
			FivePercent:            trades5Percent,
			InitialGross:           tradesInitialGross,
			TwelvePointFivePercent: trades12Point5Percent,
		}

		summary[2] = ClientInvoiceSummary{
			ChargeType:             "Total Amount",
			ChargeableQuantity:     0,
			ChargeAmount:           totalChargeAmount,
			TaxAmount:              total12Point5Percent + total5Percent,
			ChargeAmountWithTax:    totalTaxAmount,
			BasisPoint:             0,
			FivePercent:            total5Percent,
			TwelvePointFivePercent: total12Point5Percent,
			InitialGross:           totalInitialGross,
			InvoiceAmountWithTax:   totalTaxAmount,
		}
	} else {
		summary[0] = ClientInvoiceSummary{
			ChargeType:           "Portfolio Fee",
			ChargeableQuantity:   portfolioData.NAV,
			BasisPoint:           basisPoint,
			InvoiceAmountWithTax: portfolioChargeAmount,
		}

		summary[1] = ClientInvoiceSummary{
			ChargeType:           "Transaction Fee",
			ChargeableQuantity:   float64(len(txns)),
			BasisPoint:           float64(billingInfo.ChargePerTransaction),
			InvoiceAmountWithTax: tradesChargeAmount,
		}

		summary[2] = ClientInvoiceSummary{
			ChargeType:           "Total Amount",
			ChargeableQuantity:   0,
			InvoiceAmountWithTax: tradesChargeAmount + portfolioChargeAmount,
		}
	}
	return summary
}

func LoadUsersByRole(role string) []User {
	users := make([]User, 0)
	DatabaseConnection.Raw("SELECT * FROM CRS_USERS WHERE ID IN (SELECT USERID FROM CRS_USER_ROLE WHERE ROLEID = (SELECT ID FROM CRS_ROLES WHERE [NAME] = ?))", role).Scan(&users)
	return users
}

func ApproveSecCurrentQuarterReport(id int) {
	quarterDate := GetQuarterDate()
	result := DatabaseConnection.Model(GovernanceInfo{}).Where("REPORT_REF_ID=?", quarterDate).Updates(map[string]interface{}{"APPROVED": true, "APPROVED_BY": id})
	DatabaseConnection.Model(SecActivities{}).Where("QUARTER_DATE=?", MakeDBQueryableCurrentQuarterDate()).Updates(map[string]interface{}{"APPROVED": true})
	DatabaseConnection.Model(MaturedSecurity{}).Where("DATE=?", MakeDBQueryableCurrentQuarterDate().Format("2006-01-02")).Updates(map[string]interface{}{"APPROVED": true})
	if result.RowsAffected > 0 {
		SendReportActionEmail("report_approval", LoadUsersByRole(SECMakerRole), "SEC Quarterly Report Approved", "SEC", "")
	}
}

func ApproveNPRAReport(id int) int64 {
	date := MakeDBQueryableCurrentQuarterDate()
	result := DatabaseConnection.Model(NPRADeclaration{}).Where("REPORT_REF_ID=?", date.Format("2006-01-02")).Updates(map[string]interface{}{"APPROVED": true, "APPROVED_BY": id})
	DatabaseConnection.Model(NPRAActivities{}).Where("QUARTER_DATE=?", date).Updates(map[string]interface{}{"APPROVED": true})
	if result.RowsAffected > 0 {
		SendReportActionEmail("report_approval", LoadUsersByRole(NPRAMakerRole), "NPRA Report Approved", "NPRA", "")
	}
	return result.RowsAffected
}

func LogCheckerComments(checkerID int, report, comment string, date string) int64 {
	results := DatabaseConnection.Create(&CheckerDisapprovalComments{
		DisapprovedBy: checkerID,
		ReportRefID:   date,
		Comment:       comment,
		ReportType:    report,
	})
	return results.RowsAffected
}

func DisapproveSecCurrentQuarterReport(checkerID int, comment string) {
	result := LogCheckerComments(checkerID, "sec", comment, GetQuarterDate())
	if result > 0 {
		SendReportActionEmail("report_disapproval", LoadUsersByRole(SECMakerRole), "SEC Quarterly Report Disapproved", "SEC", comment)
	}
}

func DisapproveNPRACurrentQuarterReport(checkerID int, comment string) {
	result := LogCheckerComments(checkerID, "npra", comment, GetQuarterDate())
	if result > 0 {
		SendReportActionEmail("report_disapproval", LoadUsersByRole(NPRAMakerRole), "NPRA Report Disapproved", "NPRA", comment)
	}
}

func LoadSecInternalComplianceForQuarter() SecReportInternalCompliance {
	compliance := SecReportInternalCompliance{}
	compliance.Preparers = LoadUsersByRole(SECMakerRole)
	compliance.Reviewers = LoadUsersByRole(SECCheckerRole)
	user := User{}
	DatabaseConnection.Raw("select * from CRS_USERS where ID = (select APPROVED_BY from CRS_SEC_GOVERNANCE_INFO where REPORT_REF_ID = ?)", GetQuarterDate()).Scan(&user)
	compliance.Authorizer = user
	return compliance
}

func GetClientByName(name string) Client {
	client := Client{}
	DatabaseConnection.Where("CLIENTNAME=?", name).Find(&client)
	return client
}

func GetClientUniquePVCodeByName(name string) []ClientPVUniqueCode {
	clients := make([]ClientPVUniqueCode, 0)
	DatabaseConnection.Where("CLIENT_NAME=?", name).Find(&clients)
	return clients
}

func UploadSECVarianceRemarks(request VarianceRemarksUpdateRequest) {
	//TODO: refactor this
	date := MakeDBQueryableCurrentQuarterDate()
	for index, variance := range request.Client {
		if variance != "TOTAL" {
			bpid := GetClientByName(variance).BPID
			data := VarianceRemarks{ReportDate: date, BPID: bpid, Remarks: request.Remarks[index]}
			DatabaseConnection.Where(VarianceRemarks{BPID: bpid, ReportDate: date}).Assign(data).FirstOrCreate(&data)
		}
	}
}

func UploadNPRAVarianceRemarks(request VarianceRemarksUpdateRequest) {
	//TODO: refactor this
	date := MakeDBQueryableCurrentQuarterDate()
	for index, variance := range request.Client {
		if variance != "TOTAL" {
			clients := GetClientUniquePVCodeByName(variance)
			var bpid string
			if len(clients) > 1 {
				bpid = clients[0].BPID
			} else {
				bpid = clients[0].Code
			}
			data := VarianceRemarks{ReportDate: date, BPID: bpid, Remarks: request.Remarks[index]}
			DatabaseConnection.Where(VarianceRemarks{BPID: bpid, ReportDate: date}).Assign(data).FirstOrCreate(&data)
		}
	}
}

// LoadCurrentNPRADeclaration
func LoadCurrentNPRADeclaration() NPRADeclaration {
	declaration := NPRADeclaration{}
	DatabaseConnection.Where("REPORT_REF_ID=?", MakeDBQueryableCurrentQuarterDate().Format("2006-01-02")).First(&declaration)
	return declaration
}

func UpdateCurrentNPRADeclaration(declaration NPRADeclaration) {
	//TODO: find out why assigned_by becomes 0
	DatabaseConnection.Where(NPRADeclaration{ReportRefID: MakeDBQueryableCurrentQuarterDate().Format("2006-01-02")}).Assign(declaration).FirstOrCreate(&declaration)
}

func UpdateOutstandingFDCertificates(data []OutstandingFDEdit) {
	for _, certificate := range data {
		DatabaseConnection.Exec("UPDATE CRS_NPRA_OUTSTANDING_FD_CERTIFICATES SET RECEIPT_RECEIVED=? WHERE ID=?", certificate.Value, certificate.ID)
	}
}

func GetClientPositionByBPID(bpOrSca string, date time.Time) float64 {
	value := ValueQuery{}
	DatabaseConnection.Raw("SELECT SUM(LCY_AMT) as value FROM CRS_PVSUMMARY WHERE REPORT_DATE=? AND CLIENT_TYPE=? AND (BP_ID=? OR CLIENT_ID=?)", date.Format("2006-01-02"), "billing", bpOrSca, bpOrSca).Scan(&value)
	return value.Value
}

func GenerateRandomInvoiceReference() int64 {
	var number int64
	number = RandomNumber()
	value := ValueQuery{}
	DatabaseConnection.Raw("select count(*) as value from CRS_BILLING_NAV where INVOICE_REFERENCE =?", number).Scan(&value)
	var retries = 0
	for value.Value > 0 {
		number = RandomNumber()
		DatabaseConnection.Raw("select count(*) as value from CRS_BILLING_NAV where INVOICE_REFERENCE =?", number).Scan(&value)
		if retries >= 9000000 {
			panic("Random Invoice reference exhausted")
		}
		retries++
	}
	return number
}

func UpdateClientNAV(request BillingNAVUpdateRequest, authID int) {
	date, _ := time.Parse("2006-01", request.Date)
	parsedDate := now.New(date).EndOfMonth().Format("2006-01-02")
	data := BillingNav{
		RefID:       parsedDate,
		BPID:        request.BPID,
		SCA:         request.SCA,
		Position:    request.Position,
		CashBalance: request.CashBalance,
		Liabilities: request.Liabilities,
		NAV:         request.NAV,
	}
	var clientId string
	if request.SCA == "" {
		clientId = request.BPID
	} else {
		clientId = request.SCA
	}
	billing := BillingQuarterlyReport{Month: parsedDate, ClientID: clientId}
	DatabaseConnection.Where(billing).Assign(billing).FirstOrCreate(&billing)
	existingNAV := BillingNav{}
	DatabaseConnection.Where(BillingNav{RefID: parsedDate, BPID: request.BPID, SCA: request.SCA}).First(&existingNAV)
	if existingNAV.InvoiceReference == 0 {
		data.InvoiceReference = GenerateRandomInvoiceReference()
	} else {
		data.InvoiceReference = existingNAV.InvoiceReference
	}
	DatabaseConnection.Where(BillingNav{RefID: parsedDate, BPID: request.BPID, SCA: request.SCA}).Assign(data).FirstOrCreate(&data)
	DatabaseConnection.Create(&NAVUpdateLog{
		ClientID:       clientId,
		UserID:         authID,
		ReportingMonth: parsedDate,
		Position:       request.Position,
		Liabilities:    request.Liabilities,
		NAV:            request.NAV,
	})
}

func GetClientNAVDetails(bpOrSca string, date time.Time) BillingNav {
	details := BillingNav{}
	DatabaseConnection.Where("BP_ID=? OR SCA=? ", bpOrSca, bpOrSca).Where(" REF_ID=?", date.Format("2006-01-02")).First(&details)
	return details
}

func CalculateInvoiceAmountForClient(bpOrSca string, nav float64, date time.Time) float64 {
	portfolioData := GetClientNAVDetails(bpOrSca, date)
	points := make([]ClientBasisPoints, 0)
	var basisPoint float64
	DatabaseConnection.Where("BP_ID=? OR SCA=?", bpOrSca, bpOrSca).Find(&points)
	for _, bp := range points {
		if bp.Maximum == 0.0 {
			if portfolioData.NAV >= bp.Minimum {
				basisPoint = bp.BasisPoints
				break
			}
		} else {
			if portfolioData.NAV >= bp.Minimum && portfolioData.NAV <= bp.Maximum {
				basisPoint = bp.BasisPoints
				break
			}
		}
	}
	txns := make([]BillingTransactionDetails, 0)
	DatabaseConnection.Where("BP_ID= ? OR SCA = ?", bpOrSca, bpOrSca).Where("CHARGE_ITEM=? AND REPORTING_DATE= ?", "TRANSACTION FEE", date.Format("2006-01-02")).Find(&txns)

	billingInfo := LoadClientBillingInfo(bpOrSca)

	tradesChargeAmount := float64(float32(len(txns)) * billingInfo.ChargePerTransaction)

	portfolioChargeAmount := (portfolioData.NAV * basisPoint) / 120000
	/*Total trades amount is minus the tax...*/
	totalChargeAmount := portfolioChargeAmount + tradesChargeAmount
	if IsSecClient(bpOrSca) {
		total5Percent := (totalChargeAmount * 5) / 100
		totalInitialGross := total5Percent + totalChargeAmount
		total12Point5Percent := ((total5Percent + totalChargeAmount) * 12.5) / 100
		return totalInitialGross + total12Point5Percent
	}
	return tradesChargeAmount + portfolioChargeAmount
}

func LoadClientUnutilizedFunds(bpid string) float64 {
	value := ValueQuery{}
	DatabaseConnection.Raw("select sum(LCY_AMT) as value from CRS_PVSUMMARY WHERE BP_ID = ? AND REPORT_DATE = ? AND SECURITY_TYPE LIKE '%CASH BAL%'", bpid, GetQuarterFormalDate()).Scan(&value)
	return value.Value
}

func LoadClientMonthlyContributions(bpid, quarter, year string) []LumpedMonthlyContribution {
	from := MakeDBQueryableQuarterFirstDate("1", year) //start from year begin
	to := MakeDBQueryableQuarterLastDate(quarter, year)
	contributions := make([]LumpedMonthlyContribution, 0)
	DatabaseConnection.Raw("select FORMAT(DATE, 'MMMM yyyy') as date, SUM(AMOUNT) as contributions from CRS_TRUSTEE_MONTHLY_CONTRIBUTIONS WHERE BP_ID= ? AND QUARTER BETWEEN ? AND ? GROUP BY FORMAT(DATE, 'MMMM yyyy')", bpid, from, to).Scan(&contributions)
	return contributions
}

func LoadMonthlyContributionsForEachSCA(bpid, quarter, year string) []SCAMonthlyContribution {
	from := MakeDBQueryableQuarterFirstDate(quarter, year)
	to := MakeDBQueryableQuarterLastDate(quarter, year)
	clients := GetClientSCASByBPID(bpid)
	data := make([]SCAMonthlyContribution, 0)
	for _, client := range clients {
		contributions := make([]LumpedMonthlyContribution, 0)
		DatabaseConnection.Raw("select FORMAT(DATE, 'MMMM yyyy') as date, SUM(AMOUNT) as contributions from CRS_TRUSTEE_MONTHLY_CONTRIBUTIONS WHERE SCA= ? AND QUARTER BETWEEN ? AND ? GROUP BY FORMAT(DATE, 'MMMM yyyy')", client.Code, from, to).Scan(&contributions)
		if len(contributions) == 0 {
			continue
		}
		data = append(data, SCAMonthlyContribution{
			Client:        client.FundManager,
			HTMLTableID:   rand.Int() % 1000,
			Contributions: SortContributionsByMonths(contributions),
		})
	}
	return data
}

func LoadPVSummaryForEachSCA(bpid, quarter, year string) []ClientIndividualPVSummary {
	clients := GetClientSCASByBPID(bpid)
	data := make([]ClientIndividualPVSummary, 0)
	for _, client := range clients {
		summary := make([]ClientPVReportSummary, 0)
		total := PVSummaryTotal{}
		DatabaseConnection.Raw("SELECT * FROM CRS_PVSUMMARY WHERE CLIENT_ID=? AND REPORT_DATE=? AND CLIENT_TYPE=?", client.Code, MakeQuarterFormalDate(quarter, year), "trustee").Scan(&summary)
		if len(summary) == 0 {
			continue
		}
		DatabaseConnection.Raw("SELECT SUM(NOMINAL_VALUE) as nominal,SUM(LCY_AMT) as lcy,SUM(PERCENTAGE_TOTAL) as percentage FROM CRS_PVSUMMARY WHERE CLIENT_ID=? AND REPORT_DATE=? AND CLIENT_TYPE=?", client.Code, MakeQuarterFormalDate(quarter, year), "trustee").Scan(&total)
		summary = LumpPVSummary(summary)
		summary = append(summary, ClientPVReportSummary{
			SecurityType:      "TOTAL",
			NominalValue:      total.Nominal,
			LCYAmount:         total.LCY,
			PercentageOfTotal: 100,
		})
		data = append(data, ClientIndividualPVSummary{
			Client:      client.FundManager,
			HTMLTableID: rand.Int() % 1000,
			Summary:     SortPVSummary(summary),
		})
	}
	return data
}

func LoadClientGOGMaturities(bpid, quarter, year string) []GOGMaturity {
	maturities := make([]GOGMaturity, 0)
	DatabaseConnection.Where("BP_ID=? AND QUARTER_DATE = ?", bpid, MakeDBQueryableQuarterLastDate(quarter, year)).Find(&maturities)
	return maturities
}

func LoadClientTransactions(bpid, quarter, year string) []CrsTrustTxnVolumes {
	transactions := make([]CrsTrustTxnVolumes, 0)
	DatabaseConnection.Where("BP_ID=? AND QUARTER_DATE = ?", bpid, MakeDBQueryableQuarterLastDate(quarter, year)).Find(&transactions)
	return transactions
}

func StoreGOGMaturitiesInTheDB(request GOGMaturitiesUploadRequest, userId int) {
	hash := string(RandomBytes(80))
	date := MakeDBQueryableQuarterLastDate(request.Quarter, request.Year)
	for _, data := range request.Data {
		data.Hash = hash
		data.QuarterDate = date
		DatabaseConnection.Create(&data)
	}
	if len(request.Data) > 0 {
		DatabaseConnection.Create(&TrusteeActivities{
			Hash:        hash,
			BPID:        request.Data[0].BPID,
			UserID:      userId,
			Date:        time.Now(),
			QuarterDate: date,
			Approved:    false,
			Activity:    UploadGOGMaturities,
		})
	}
}

func StoreTransactionVolumes(request TransactionVolumesUploadRequest, userId int) {
	hash := string(RandomBytes(80))
	date := MakeDBQueryableQuarterLastDate(request.Quarter, request.Year)
	for _, data := range request.Data {
		data.Hash = hash
		data.QuarterDate = date
		DatabaseConnection.Create(&data)
	}
	if len(request.Data) > 0 {
		DatabaseConnection.Create(&TrusteeActivities{
			Hash:        hash,
			BPID:        request.Data[0].BPID,
			UserID:      userId,
			Date:        time.Now(),
			QuarterDate: date,
			Approved:    false,
			Activity:    UploadTransactionVolumes,
		})
	}
}

func LoadClientGOGMaturitiesSummary(bpid, quarter, year string) []GOGSummary {
	summary := make([]GOGSummary, 0)
	DatabaseConnection.Raw("SELECT FORMAT(ENTRY_DATE, 'MMMM yyyy') as month, COUNT(*) as number_of_maturities,SUM(GROSS_AMOUNT) as amount FROM CRS_TRUSTEE_GOG_FD WHERE STATUS='Posted' AND DEPOT_ID = 'GH000047-BANK OF GHANA' AND BP_ID=? AND QUARTER_DATE =? GROUP BY FORMAT(ENTRY_DATE, 'MMMM yyyy')", bpid, MakeDBQueryableQuarterLastDate(quarter, year)).Scan(&summary)
	if len(summary) > 0 {
		total := GOGSummary{Month: "TOTAL"}
		for _, data := range summary {
			total.Amount += data.Amount
			total.NumberOfMaturities += data.NumberOfMaturities
		}
		summary = append(summary, total)
	}
	return summary
}

func LoadClientFDMaturitiesSummary(bpid, quarter, year string) []StrValueQuery {
	data := make([]StrValueQuery, 0)
	DatabaseConnection.Raw("SELECT BASE_SECURITY_ID as value FROM CRS_TRUSTEE_GOG_FD WHERE STATUS='Posted' AND DEPOT_ID = 'GH000046-BANKCASH DEPOSIT' AND BP_ID=? AND QUARTER_DATE= ?", bpid, MakeDBQueryableQuarterLastDate(quarter, year)).Scan(&data)
	return data
}

func LoadClientCorporateActionActivities(bpid, quarter, year string) []CorporateActionActivity {
	actions := make([]CorporateActionActivity, 0)
	DatabaseConnection.Raw(" SELECT EVENT_TYPE as activity,count(*) as number FROM CRS_TRUSTEE_GOG_FD WHERE STATUS='Posted' AND BP_ID=? AND QUARTER_DATE= ? GROUP BY EVENT_TYPE", bpid, MakeDBQueryableQuarterLastDate(quarter, year)).Scan(&actions)
	if len(actions) > 0 {
		total := CorporateActionActivity{Activity: "TOTAL"}
		for _, data := range actions {
			total.Number += data.Number
		}
		actions = append(actions, total)
	}
	return actions
}

func DeleteNPRAPV(request DeleteUploadedPVRequest) {
	DatabaseConnection.Exec("DELETE FROM CRS_NPRA_ACTIVITIES WHERE BP_ID=? AND QUARTER_DATE=?", request.BPID, request.Date)
	DatabaseConnection.Exec("DELETE FROM CRS_PVREPORTS WHERE BP_ID=? AND REPORT_DATE=? AND CLIENT_TYPE=?", request.BPID, request.Date, "npra")
	DatabaseConnection.Exec("DELETE FROM CRS_PVSUMMARY WHERE BP_ID=? AND REPORT_DATE=? AND CLIENT_TYPE=?", request.BPID, request.Date, "npra")
}

func DeleteSECPV(request DeleteUploadedPVRequest) {
	DatabaseConnection.Exec("DELETE FROM CRS_SEC_ACTIVITIES WHERE BP_ID=? AND QUARTER_DATE=?", request.BPID, request.Date)
	DatabaseConnection.Exec("DELETE FROM CRS_PVREPORTS WHERE BP_ID=? AND REPORT_DATE=? AND CLIENT_TYPE=?", request.BPID, request.Date, "sec")
	DatabaseConnection.Exec("DELETE FROM CRS_PVSUMMARY WHERE BP_ID=? AND REPORT_DATE=? AND CLIENT_TYPE=?", request.BPID, request.Date, "sec")
}

func DeleteBillingPV(request DeleteUploadedPVRequest) {
	DatabaseConnection.Exec("DELETE FROM CRS_BILLING_ACTIVITIES WHERE BP_ID=? AND QUARTER_DATE=?", request.BPID, request.Date)
	DatabaseConnection.Exec("DELETE FROM CRS_TRUSTEE_ACTIVITIES WHERE BP_ID=? AND QUARTER_DATE=? AND ACTIVITY=?", request.BPID, request.Date, UploadSecPV)
	DatabaseConnection.Exec("DELETE FROM CRS_PVREPORTS WHERE BP_ID=? AND REPORT_DATE=? AND CLIENT_TYPE=?", request.BPID, request.Date, "billing")
	DatabaseConnection.Exec("DELETE FROM CRS_PVSUMMARY WHERE BP_ID=? AND REPORT_DATE=? AND CLIENT_TYPE=?", request.BPID, request.Date, "billing")
}

func LoadLumpedPVSummaryForPast3Quarters(bpid, quarter, year string) [][]ClientPVReportSummary {
	lumpedSummary := make([][]ClientPVReportSummary, 0)
	dates := MakeLast4QuarterDates(MakeQuarterFormalDate(quarter, year))[1:]
	for _, date := range dates {
		lumpedSummary = append(lumpedSummary, LumpPVSummary(LoadClientPVSummary(bpid, date.Format("2006-01-02"))))
	}
	return lumpedSummary
}

func LoadClientUnidentifiedPayments(bpid, quarter, year string) []UnidentifiedPaymentsSummary {
	payments := make([]UnidentifiedPaymentsSummary, 0)
	firstDay := MakeDBQueryableQuarterFirstDate(quarter, year)
	lastDay := MakeDBQueryableQuarterLastDate(quarter, year)
	DatabaseConnection.Raw("select TRANSACTION_TYPE as type,count(*) as total,(select count(*) from CRS_TRUSTEE_UNIDENTIFIED_PAYMENTS where CLIENT_BPID=? and TRANSACTION_TYPE= base_query.TRANSACTION_TYPE and TRANSACTION_DATE between ? and ? and STATUS='pending') as pending,(select count(*) from CRS_TRUSTEE_UNIDENTIFIED_PAYMENTS where CLIENT_BPID=? and TRANSACTION_TYPE= base_query.TRANSACTION_TYPE and TRANSACTION_DATE between ? and ? and STATUS='done') as done from CRS_TRUSTEE_UNIDENTIFIED_PAYMENTS base_query where CLIENT_BPID=? and TRANSACTION_DATE between ? and ? group by TRANSACTION_TYPE", bpid, firstDay, lastDay, bpid, firstDay, lastDay, bpid, firstDay, lastDay).Scan(&payments)
	return payments
}

func DeleteOutstandingFDCertificates(hash string) {
	DatabaseConnection.Exec("DELETE FROM CRS_NPRA_ACTIVITIES WHERE HASH=?", hash)
	DatabaseConnection.Exec("DELETE FROM CRS_NPRA_OUTSTANDING_FD_CERTIFICATES WHERE HASH=?", hash)
}

func GetClientNAVByBPOrSca(bpOrSca string, date time.Time) BillingNav {
	nav := BillingNav{}
	DatabaseConnection.Where("BP_ID=? OR SCA=?", bpOrSca, bpOrSca).Where("REF_ID=?", date.Format("2006-01-02")).First(&nav)
	return nav
}

func LoadNPRAClients() []Client {
	clients := make([]Client, 0)
	DatabaseConnection.Where("CLIENT_TYPE=?", "NPRA").Find(&clients)
	return clients
}

func LoadClientTrusteeQuarterlyReport(bpid, date string, createIfNotExists bool) TrusteeQuarterlyReport {
	report := TrusteeQuarterlyReport{}
	DatabaseConnection.Where("BP_ID=? AND QUARTER=?", bpid, date).First(&report)
	if report.ID == 0 && createIfNotExists {
		report.Quarter = date
		report.BPID = bpid
		DatabaseConnection.Create(&report)
	}
	return report
}

func ApproveTrusteeReport(id int, bpid string) {
	result := DatabaseConnection.Model(TrusteeQuarterlyReport{}).Where("QUARTER=? AND BP_ID=?", GetQuarterFormalDate(), bpid).Updates(map[string]interface{}{"APPROVED": true, "APPROVED_BY": id})
	date := MakeDBQueryableCurrentQuarterDate()
	DatabaseConnection.Model(TrusteeActivities{}).Where("QUARTER_DATE=? AND BP_ID=?", date, bpid).Updates(map[string]interface{}{"APPROVED": true})
	if result.RowsAffected > 0 {
		//SendReportActionEmail("report_approval", LoadUsersByRole(NPRAMakerRole), "NPRA Report Approved", "NPRA", "")
	}
}

func DisapproveTrusteeReport(checkerID int, comment string) {
	result := LogCheckerComments(checkerID, "trustee", comment, GetQuarterDate())
	if result > 0 {
		//SendReportActionEmail("report_disapproval", LoadUsersByRole(NPRAMakerRole), "NPRA Report Disapproved", "NPRA", comment)
	}
}

func DeleteMonthlyContribution(id int) {
	DatabaseConnection.Exec("DELETE FROM CRS_TRUSTEE_MONTHLY_CONTRIBUTIONS WHERE ID=?", id)
}

func LoadTrusteeActivityByHash(hash string) TrusteeActivities {
	activity := TrusteeActivities{}
	DatabaseConnection.Where("HASH=?", hash).First(&activity)
	return activity
}

func LoadTransactionVolumesByHash(hash string) []CrsTrustTxnVolumes {
	txns := make([]CrsTrustTxnVolumes, 0)
	DatabaseConnection.Where("HASH=?", hash).Find(&txns)
	return txns
}

func LoadGOGMaturitiesByHash(hash string) []GOGMaturity {
	maturities := make([]GOGMaturity, 0)
	DatabaseConnection.Where("HASH=?", hash).Find(&maturities)
	return maturities
}

func DeleteGOGMaturities(hash string) {
	DatabaseConnection.Exec("DELETE FROM CRS_TRUSTEE_GOG_FD WHERE HASH = ?", hash)
	DatabaseConnection.Exec("DELETE FROM CRS_TRUSTEE_ACTIVITIES WHERE HASH = ?", hash)
}

func DeleteTransactionVolumes(hash string) {
	DatabaseConnection.Exec("DELETE FROM CRS_TRUST_TXN_VOLUMES WHERE HASH = ?", hash)
	DatabaseConnection.Exec("DELETE FROM CRS_TRUSTEE_ACTIVITIES WHERE HASH = ?", hash)
}

func GetClientBillingReportInfo(clientID string, month time.Time) BillingQuarterlyReport {
	report := BillingQuarterlyReport{}
	DatabaseConnection.Where("CLIENT_ID=? AND MONTH=?", clientID, month.Format("2006-01-02")).First(&report)
	return report
}

func DisapproveBillingReport(checkerID int, date, comment, clientID string) {
	parsedDate, _ := time.Parse("2006-01-02", date)
	result := LogCheckerComments(checkerID, "billing", comment, now.New(parsedDate).EndOfMonth().Format("02 January 2006"))
	if result > 0 {
		//SendReportActionEmail("report_disapproval", LoadUsersByRole(NPRAMakerRole), "NPRA Report Disapproved", "NPRA", comment)
	}
}

func ApproveBillingReport(checkerID int, date, clientID string) {
	parsedDate, _ := time.Parse("2006-01-02", date)
	endOfMonth := now.New(parsedDate).EndOfMonth()
	result := DatabaseConnection.Model(BillingQuarterlyReport{}).Where("MONTH=? AND CLIENT_ID=?", endOfMonth.Format("2006-01-02"), clientID).Updates(map[string]interface{}{"APPROVED": true, "APPROVED_BY": checkerID})
	//DatabaseConnection.Model(SecActivities{}).Where("QUARTER_DATE=?", quarterDate).Updates(map[string]interface{}{"APPROVED": true})
	if result.RowsAffected > 0 {
		//SendReportActionEmail("report_approval", LoadUsersByRole(SECMakerRole), "SEC Quarterly Report Approved", "SEC", "")
	}
	CreateBillingTransactionJournal(checkerID, clientID, endOfMonth)
}

func AddClientSCA(request SCAManagementRequest) {
	DatabaseConnection.Create(&ClientPVUniqueCode{
		BPID:          request.BPID,
		ClientName:    request.ClientName,
		ClientTier:    request.ClientTier,
		Code:          request.Code,
		FundManager:   request.FundManager,
		Administrator: request.Administrator,
		GroupKey:      request.GroupKey,
	})
}

func EditClientSCA(request SCAManagementRequest) {
	DatabaseConnection.Model(ClientPVUniqueCode{}).Where("ID=?", request.ID).Updates(map[string]interface{}{
		"CLIENT_NAME":   request.ClientName,
		"CLIENT_TIER":   request.ClientTier,
		"CODE":          request.Code,
		"FUND_MANAGER":  request.FundManager,
		"ADMINISTRATOR": request.Administrator,
		"GROUP_KEY":     request.GroupKey,
	})
}

func DeleteClientSCA(id int) {
	DatabaseConnection.Exec("DELETE FROM CRS_CLIENT_PV_UNIQUE_CODES WHERE ID=?", id)
}

func ReverseBillingReportApproval(checkerID int, date, clientID string) {
	parsedDate, _ := time.Parse("2006-01-02", date)
	date = now.New(parsedDate).EndOfMonth().Format("2006-01-02")
	report := BillingQuarterlyReport{}
	DatabaseConnection.Where("MONTH=? AND CLIENT_ID=?", date, clientID).First(&report)

	if report.ID != 0 {
		result := DatabaseConnection.Model(BillingQuarterlyReport{}).Where("MONTH=? AND CLIENT_ID=?", date, clientID).Updates(map[string]interface{}{"APPROVED": false, "APPROVED_BY": nil})
		DatabaseConnection.Create(&ApprovedReportReversal{
			Type:                 "billing",
			ClientID:             clientID,
			ReportDate:           date,
			ReversedBy:           checkerID,
			ReversedOn:           time.Now(),
			PreviouslyApprovedBy: *report.ApprovedBy,
		})
		if result.RowsAffected > 0 {
			//SendReportActionEmail("report_approval", LoadUsersByRole(SECMakerRole), "SEC Quarterly Report Approved", "SEC", "")
		}
	}
}

func CreateBillingTransactionJournal(userID int, clientID string, date time.Time) {
	portfolioData := GetClientNAVDetails(clientID, date)
	points := make([]ClientBasisPoints, 0)
	var basisPoint float64
	DatabaseConnection.Where("BP_ID=? OR SCA=?", clientID, clientID).Find(&points)
	for _, bp := range points {
		if bp.Maximum == 0.0 {
			if portfolioData.NAV >= bp.Minimum {
				basisPoint = bp.BasisPoints
				break
			}
		} else {
			if portfolioData.NAV >= bp.Minimum && portfolioData.NAV <= bp.Maximum {
				basisPoint = bp.BasisPoints
				break
			}
		}
	}
	txns := make([]BillingTransactionDetails, 0)
	DatabaseConnection.Where("BP_ID= ? OR SCA = ?", clientID, clientID).Where("CHARGE_ITEM=? AND REPORTING_DATE= ?", "TRANSACTION FEE", date.Format("2006-01-02")).Find(&txns)

	billingInfo := LoadClientBillingInfo(clientID)

	tradesChargeAmount := float64(float32(len(txns)) * billingInfo.ChargePerTransaction)

	portfolioChargeAmount := (portfolioData.NAV * basisPoint) / 120000
	var invoiceAmount float64
	if IsSecClient(clientID) {
		totalChargeAmount := portfolioChargeAmount + tradesChargeAmount
		total5Percent := (totalChargeAmount * 5) / 100
		totalInitialGross := total5Percent + totalChargeAmount
		total12Point5Percent := (totalInitialGross * 12.5) / 100

		invoiceAmount = totalInitialGross + total12Point5Percent
	} else {
		invoiceAmount = tradesChargeAmount + portfolioChargeAmount
	}

	DatabaseConnection.Create(&BillingTransactionJournal{
		InvoiceDate:                   date,
		InvoiceReference:              portfolioData.InvoiceReference,
		InvoiceAmount:                 invoiceAmount,
		Position:                      portfolioData.Position,
		Liabilities:                   portfolioData.Liabilities,
		NAV:                           portfolioData.NAV,
		ClientID:                      clientID,
		InvoiceDueDate:                date.AddDate(0, 1, 0),
		ApprovedOn:                    time.Now(),
		ApprovedBy:                    userID,
		PortfolioChargeableQuantity:   portfolioData.NAV,
		TransactionChargeableQuantity: len(txns),
		BasisPoint:                    basisPoint,
		ChargePerTransaction:          billingInfo.ChargePerTransaction,
		PortfolioChargeAmount:         portfolioChargeAmount,
		TransactionChargeAmount:       tradesChargeAmount,
	})
}

// N+1 :(
func GetAllBillingDetails() []ClientDetailedBillingDetails {
	data := make([]ClientDetailedBillingDetails, 0)
	billingInfo := make([]ClientBillingInfo, 0)
	DatabaseConnection.Find(&billingInfo)
	for _, info := range billingInfo {
		client := GetClientByBPID(info.BPID)
		if client.Type == "NPRA" {
			data = append(data, ClientDetailedBillingDetails{
				ClientName:           client.Name,
				MinimumCharge:        info.MinimumCharge,
				ChargePerTransaction: info.ChargePerTransaction,
				ThirdPartyTransfer:   info.ThirdPartyTransfer,
				BasisPoints:          LoadClientBasisPoints(info.BPID),
			})
		}
	}
	return data
}

func LoadAuditablePVData(BPID, Date, pvType string) []AuditablePV {
	pvs := make([]AuditablePV, 0)
	codes := make([]PVReportField, 0)
	DatabaseConnection.Select("DISTINCT(CLIENT_ID)").Where("BP_ID=? OR CLIENT_ID=?", BPID, BPID).Where("REPORT_DATE=? AND CLIENT_TYPE=?", Date, pvType).Preload("ClientSCADetails").Find(&codes)
	if len(codes) > 0 {
		for _, code := range codes {
			reports := LoadPVReportDataByQuarter(code.ClientID, Date, pvType)
			summary := SortPVSummary(LoadPVSummaryDataByQuarter(code.ClientID, Date, pvType))
			client := code.ClientSCADetails
			if client.ID == 0 {
				//for some reason this happens
				client = GetClientByCode(code.ClientID)
			}
			pvs = append(pvs, AuditablePV{Client: client, Reports: reports, Summary: summary})
		}
	} else {
		codes := make([]ClientPVReportSummary, 0)
		DatabaseConnection.Select("DISTINCT(CLIENT_ID)").Where("BP_ID=? OR CLIENT_ID=?", BPID, BPID).Where("REPORT_DATE=? AND CLIENT_TYPE=?", Date, pvType).Preload("ClientSCADetails").Find(&codes)
		for _, code := range codes {
			summary := SortPVSummary(LoadPVSummaryDataByQuarter(code.ClientID, Date, pvType))
			client := code.ClientSCADetails
			if client.ID == 0 {
				//for some reason this happens
				client = GetClientByCode(code.ClientID)
			}
			pvs = append(pvs, AuditablePV{Client: client, Summary: summary})
		}
	}
	return pvs
}

/**
PV REPORT INSERT FNCs
*/

// GetClientDetails returns the clients details from the client id
func GetClientDetails(clientID string) []ClientDetailsQuery {
	query := make([]ClientDetailsQuery, 0)
	DatabaseConnection.Raw("select BP_ID,SAFEKEEPINGACCNO,CLIENTNAME from CRS_CLIENT WHERE BP_ID IN (select BP_ID from CRS_CLIENT_PV_UNIQUE_CODES where CODE = ?)", clientID).Scan(&query)
	return query
}

// ReportWithDateExists checks if a report with the given date exists
func ReportWithDateExists(date, bpid, clientID, clientType string) bool {
	report := IntValueQuery{}
	DatabaseConnection.Raw("SELECT COUNT(*) AS value FROM CRS_PVREPORTS WHERE REPORT_DATE = ? AND BP_ID=? AND CLIENT_ID=? AND CLIENT_TYPE=?", date, bpid, clientID, clientType).Scan(&report)
	return report.Value > 0
}

// UploadPVReportData handles the logic to upload the pv report data into the DB
func UploadPVReportData(client PVClient, wg *sync.WaitGroup, errorChannel chan PVUploadError) {
	defer wg.Done()
	allClients := GetClientDetails(client.ClientID)
	if len(allClients) == 0 {
		errorChannel <- PVUploadError{
			Type:       "Unknown Code Error",
			SCA:        client.ClientID,
			ClientInfo: client.RawHeading,
			Date:       client.Date,
			MoreInfo:   fmt.Sprintf("Unknown %s", client.ClientID),
		}
		return
	} else {
		for _, details := range allClients {
			if !ReportWithDateExists(client.Date, details.BPID, client.ClientID, client.Type) {
				if details.BPID != "" {
					transaction := DatabaseConnection.Begin()
					for _, report := range client.Report {
						title := report.Title
						for _, value := range report.Values {
							var dateFrom time.Time
							var dateTo time.Time
							if value.Dates.From != "" {
								if val, err := ParseDate(value.Dates.From); err == nil {
									dateFrom = val
								}
							}
							if value.Dates.To != "" {
								if val, err := ParseDate(value.Dates.To); err == nil {
									dateTo = val
								}
							}

							if err := transaction.Exec("INSERT INTO CRS_PVREPORTS (SEC_TITLE,SAFEKEEPINGACCNO,BP_ID,CLIENT_ID,CLIENT_TYPE,SECURITY_NAME,CDS_CODE,ISIN,SCBCODE,MARKET_PRICE,NOMINAL_VALUE,CUMULATIVE_COST,VALUE_AMOUNT,CREATED_AT,UPDATED_AT,REPORT_DATE,PERCENTAGE_OF_TOTAL,DATE_FROM,DATE_TO) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)", title, details.SafekeepingAccount, details.BPID, client.ClientID, client.Type, value.SecurityName, value.CDSCode, value.ISIN, value.SCBCode, value.MarketPrice, value.NominalValue, value.CumulativeCost, value.Value, time.Now(), time.Now(), client.Date, value.PercentageOfTotal, dateFrom, dateTo).Error; err != nil {
								transaction.Rollback()
								errorChannel <- PVUploadError{
									Type:       "Upload Error",
									SCA:        client.ClientID,
									ClientInfo: client.RawHeading,
									Date:       client.Date,
									MoreInfo:   err.Error(),
								}
								return
							}
						}
					}
					if client.CashBalance != 0 {
						summaryTotal := GetTotalLCYValueFromSummary(client.Summary) + client.CashBalance

						for _, summary := range client.Summary {
							if err := transaction.Exec("INSERT INTO CRS_PVSUMMARY (SAFEKEEPINGACCNO,BP_ID,CLIENT_ID,CLIENT_TYPE,REPORT_DATE,SECURITY_TYPE,NOMINAL_VALUE,CUMULATIVE_COST,LCY_AMT,PERCENTAGE_TOTAL) VALUES (?,?,?,?,?,?,?,?,?,?)", details.SafekeepingAccount, details.BPID, client.ClientID, client.Type, client.Date, summary.SecurityName, summary.NominalValue, summary.CumulativeCost, summary.Value, (summary.Value/summaryTotal)*100).Error; err != nil {
								transaction.Rollback()
								errorChannel <- PVUploadError{
									Type:       "Upload Error",
									SCA:        client.ClientID,
									ClientInfo: client.RawHeading,
									Date:       client.Date,
									MoreInfo:   err.Error(),
								}
								return
							}
						}

						if err := transaction.Exec("INSERT INTO CRS_PVSUMMARY (SAFEKEEPINGACCNO,BP_ID,CLIENT_ID,CLIENT_TYPE,REPORT_DATE,SECURITY_TYPE,NOMINAL_VALUE,CUMULATIVE_COST,LCY_AMT,PERCENTAGE_TOTAL) VALUES (?,?,?,?,?,?,?,?,?,?)", details.SafekeepingAccount, details.BPID, client.ClientID, client.Type, client.Date, "CASH BALANCE", 0, 0, client.CashBalance, (client.CashBalance/summaryTotal)*100).Error; err != nil {
							transaction.Rollback()
							errorChannel <- PVUploadError{
								Type:       "Upload Error",
								SCA:        client.ClientID,
								ClientInfo: client.RawHeading,
								Date:       client.Date,
								MoreInfo:   err.Error(),
							}
							return
						}
					} else {
						for _, summary := range client.Summary {
							if err := transaction.Exec("INSERT INTO CRS_PVSUMMARY (SAFEKEEPINGACCNO,BP_ID,CLIENT_ID,CLIENT_TYPE,REPORT_DATE,SECURITY_TYPE,NOMINAL_VALUE,CUMULATIVE_COST,LCY_AMT,PERCENTAGE_TOTAL) VALUES (?,?,?,?,?,?,?,?,?,?)", details.SafekeepingAccount, details.BPID, client.ClientID, client.Type, client.Date, summary.SecurityName, summary.NominalValue, summary.CumulativeCost, summary.Value, summary.PercentageOfTotal).Error; err != nil {
								transaction.Rollback()
								errorChannel <- PVUploadError{
									Type:       "Upload Error",
									SCA:        client.ClientID,
									ClientInfo: client.RawHeading,
									Date:       client.Date,
									MoreInfo:   err.Error(),
								}
								return
							}
						}
					}
					date, _ := time.Parse("2006-01-02", client.Date)
					//uploader := GetUserByID(client.UserID)
					if client.Type == "sec" {
						if err := transaction.Exec("INSERT INTO CRS_SEC_ACTIVITIES (BP_ID,USER_ID,DATE,ACTIVITY,QUARTER_DATE,APPROVED) VALUES (?,?,?,?,?,?)", details.BPID, client.UserID, time.Now(), UploadSecPV, date, false).Error; err != nil {
							transaction.Rollback()
							errorChannel <- PVUploadError{
								Type:       "Upload Error",
								SCA:        client.ClientID,
								ClientInfo: client.RawHeading,
								Date:       client.Date,
								MoreInfo:   err.Error(),
							}
							return
						}
						checkers := make([]User, 0)
						transaction.Raw("SELECT * FROM CRS_USERS WHERE ID IN (SELECT USERID FROM CRS_USER_ROLE WHERE ROLEID = (SELECT ID FROM CRS_ROLES WHERE [NAME] = 'sec_checker'))").Scan(&checkers)
						//SendPVUploadedToCheckers(uploader, checkers, details.Name, "sec_pv")
					} else if client.Type == "npra" {
						if err := transaction.Exec("INSERT INTO CRS_NPRA_ACTIVITIES (BP_ID,USER_ID,DATE,ACTIVITY,QUARTER_DATE,APPROVED) VALUES (?,?,?,?,?,?)", details.BPID, client.UserID, time.Now(), UploadSecPV, date, false).Error; err != nil {
							transaction.Rollback()
							panic(err.Error())
						}
						checkers := make([]User, 0)
						transaction.Raw("SELECT * FROM CRS_USERS WHERE ID IN (SELECT USERID FROM CRS_USER_ROLE WHERE ROLEID = (SELECT ID FROM CRS_ROLES WHERE [NAME] = 'npra_checker'))").Scan(&checkers)
						//SendPVUploadedToCheckers(uploader, checkers, details.Name, "npra_pv")
					} else if client.Type == "billing" {
						if err := transaction.Exec("INSERT INTO CRS_BILLING_ACTIVITIES (BP_ID,USER_ID,DATE,ACTIVITY,QUARTER_DATE,APPROVED) VALUES (?,?,?,?,?,?)", details.BPID, client.UserID, time.Now(), UploadSecPV, date, false).Error; err != nil {
							transaction.Rollback()
							errorChannel <- PVUploadError{
								Type:       "Upload Error",
								SCA:        client.ClientID,
								ClientInfo: client.RawHeading,
								Date:       client.Date,
								MoreInfo:   err.Error(),
							}
							return
						}
					} else if client.Type == "trustee" {
						if err := transaction.Exec("INSERT INTO CRS_TRUSTEE_ACTIVITIES (BP_ID,USER_ID,DATE,ACTIVITY,QUARTER_DATE,APPROVED) VALUES (?,?,?,?,?,?)", details.BPID, client.UserID, time.Now(), UploadSecPV, date, false).Error; err != nil {
							transaction.Rollback()
							errorChannel <- PVUploadError{
								Type:       "Upload Error",
								SCA:        client.ClientID,
								ClientInfo: client.RawHeading,
								Date:       client.Date,
								MoreInfo:   err.Error(),
							}
							return
						}
					}
					if err := transaction.Commit().Error; err != nil {
						transaction.Rollback()
						errorChannel <- PVUploadError{
							Type:       "Upload Error",
							SCA:        client.ClientID,
							ClientInfo: client.RawHeading,
							Date:       client.Date,
							MoreInfo:   err.Error(),
						}
						return
					}
				} else {
					errorChannel <- PVUploadError{
						Type:       "Unknown Code Error",
						SCA:        client.ClientID,
						ClientInfo: client.RawHeading,
						Date:       client.Date,
						MoreInfo:   fmt.Sprintf("Unknown %s", client.ClientID),
					}
					return
				}
			} else {
				errorChannel <- PVUploadError{
					Type:       "PV Exists Error",
					SCA:        client.ClientID,
					ClientInfo: client.RawHeading,
					Date:       client.Date,
					MoreInfo:   fmt.Sprintf("The PV for %s has already been uploaded", client.ClientID),
				}
				return
			}
		}
	}
}

func AddDirectorToDB(request ManageDirectorRequest) {
	DatabaseConnection.Create(&Director{
		Fullname:           request.Fullname,
		Nationality:        request.Nationality,
		ResidentialAddress: request.ResidentialAddress,
		Position:           request.Position,
	})
}

func EditDirectorInDB(request ManageDirectorRequest) {
	director := Director{}
	DatabaseConnection.Where("ID=?", request.ID).First(&director)
	director.Fullname = request.Fullname
	director.Nationality = request.Nationality
	director.ResidentialAddress = request.ResidentialAddress
	director.Position = request.Position
	DatabaseConnection.Save(&director)
}

func DeleteDirector(id int) {
	DatabaseConnection.Exec("DELETE FROM CRS_DIRECTORS WHERE ID=?", id)
}

func LoadAllDirectors() []Director {
	directors := make([]Director, 0)
	DatabaseConnection.Find(&directors)
	return directors
}

func UploadMaturedSecurities(maturities []MaturedSecurity, userID int) {
	for _, maturity := range maturities {
		DatabaseConnection.Create(&maturity)
	}
	if len(maturities) > 0 {
		DatabaseConnection.Create(&SecActivities{
			UserID:      userID,
			Activity:    UploadedMaturedSecurities,
			Date:        time.Now(),
			QuarterDate: MakeDBQueryableCurrentQuarterDate(),
			Approved:    false,
		})
	}
}

func DeleteMaturedSecurities(date string) {
	DatabaseConnection.Exec("DELETE FROM CRS_MATURED_SECURITIES WHERE DATE=?", date)
	parsedDate, _ := time.Parse("2006-01-02", date)
	DatabaseConnection.Exec("DELETE FROM CRS_SEC_ACTIVITIES WHERE ACTIVITY=? AND QUARTER_DATE=?", UploadedMaturedSecurities, parsedDate)
}

func LoadMaturedSecurities(date string) []MaturedSecurity {
	securities := make([]MaturedSecurity, 0)
	DatabaseConnection.Where("DATE=?", date).Find(&securities)
	return securities
}

func UpdateUserProfile(userID int, request UserProfileUpdateRequest) {
	if password, err := bcrypt.GenerateFromPassword([]byte(request.Password), 10); err == nil {
		DatabaseConnection.Table("CRS_USERS").Where("id = ?", userID).Updates(map[string]interface{}{"PASSWORD": password, "MUSTRESETPASSWORD": false})
	}
}

func CloseClientAccount(bpid string) {
	DatabaseConnection.Model(Client{}).Where("BP_ID=?", bpid).Update("CLOSED", time.Now())
	CloseAllClientSCAs(bpid)
}

func OpenClientAccount(bpid string) {
	DatabaseConnection.Model(Client{}).Where("BP_ID=?", bpid).Update("CLOSED", nil)
	OpenAllClientSCAs(bpid)
}

func CloseAllClientSCAs(bpid string) {
	DatabaseConnection.Model(ClientPVUniqueCode{}).Where("BP_ID=?", bpid).Update("CLOSED", time.Now())
}

func OpenAllClientSCAs(bpid string) {
	DatabaseConnection.Model(ClientPVUniqueCode{}).Where("BP_ID=?", bpid).Update("CLOSED", nil)
}

func CloseSCA(code string) {
	DatabaseConnection.Model(ClientPVUniqueCode{}).Where("CODE=?", code).Update("CLOSED", time.Now())
}

func OpenSCA(code string) {
	DatabaseConnection.Model(ClientPVUniqueCode{}).Where("CODE=?", code).Update("CLOSED", nil)
}

func GetClientEmailsByBPID(bpid string) []ClientEmail {
	emails := make([]ClientEmail, 0)
	DatabaseConnection.Model(ClientEmail{}).Where("BPID=?", bpid).Find(&emails)
	return emails
}

func AddClientEmail(bpid, email string) {
	clientEmail := ClientEmail{BPID: bpid, Email: email}
	DatabaseConnection.Create(&clientEmail)
}

func AddMailServiceEmail(email string) {
	mailService := MailService{Email: email}
	DatabaseConnection.Create(&mailService)
}

func EditClientEmail(id int, email string) {
	DatabaseConnection.Model(ClientEmail{}).Where("ID=?", id).Update("EMAIL", email)
}

func EditMailServiceEmail(id int, email string) {
	DatabaseConnection.Model(MailService{}).Where("ID=?", id).Update("EMAIL", email)
}

func DeleteClientEmail(id int) {
	DatabaseConnection.Exec("DELETE FROM CRS_CLIENT_EMAILS WHERE ID=?", id)
}

func DeleteMailServiceEmail(id int) {
	DatabaseConnection.Exec("DELETE FROM CRS_MAIL_SERVICE WHERE ID=?", id)
}

func LoadMailService() []MailService {
	mails := make([]MailService, 0)
	DatabaseConnection.Find(&mails)
	return mails
}

func LoadBilledClients(date time.Time) []BilledClients {
	clients := make([]BilledClients, 0)
	DatabaseConnection.Raw("select CLIENTNAME as client,CRS_CLIENT.BP_ID as bpid,INVOICE_REFERENCE as invoice_reference from CRS_CLIENT left join CRS_BILLING_NAV on (CRS_BILLING_NAV.BP_ID = CRS_CLIENT.BP_ID) AND CRS_BILLING_NAV.REF_ID=? WHERE CLIENT_TYPE <> 'FOREIGN CLIENT' ORDER BY CLIENTNAME ASC", date.Format("2006-01-02")).Scan(&clients)
	return clients
}

func LoadBilledClientSCAs(bpid string, date time.Time) []BilledClients {
	clients := make([]BilledClients, 0)
	DatabaseConnection.Raw("select CLIENT_NAME as client,CRS_CLIENT_PV_UNIQUE_CODES.CODE as bpid,INVOICE_REFERENCE as invoice_reference from CRS_CLIENT_PV_UNIQUE_CODES left join CRS_BILLING_NAV on (CRS_BILLING_NAV.SCA = CRS_CLIENT_PV_UNIQUE_CODES.CODE) AND CRS_BILLING_NAV.REF_ID=? WHERE CRS_CLIENT_PV_UNIQUE_CODES.BP_ID =?", date.Format("2006-01-02"), bpid).Scan(&clients)
	return clients
}

func LoadMailSenderInfo() MailSenderInfo {
	sender := MailSenderInfo{}
	DatabaseConnection.First(&sender)
	return sender
}

func UpdateMailSenderInfo(request MailSenderInfoUpdateRequest) {
	DatabaseConnection.Model(MailSenderInfo{}).Updates(map[string]interface{}{"NAME": request.Name, "POSITION": request.Position})
}

func LoadHolidays() []Holiday {
	days := make([]Holiday, 0)
	DatabaseConnection.Find(&days)
	return days
}

func AddHoliday(request HolidayMgmtRequest) {
	data := Holiday{Name: request.Name}
	date, _ := time.Parse("2006-01-02", request.Date)
	data.Date = date.Format("02-01")
	DatabaseConnection.Create(&data)
}

func DeleteHoliday(id int) {
	DatabaseConnection.Exec("DELETE FROM CRS_HOLIDAYS WHERE ID=?", id)
}

func AddSecurity(data Securities) {
	data.Name = strings.ToUpper(data.Name)
	data.Plural = strings.ToUpper(data.Plural)
	DatabaseConnection.Create(&data)
}

func EditSecurity(data Securities) {
	data.Name = strings.ToUpper(data.Name)
	data.Plural = strings.ToUpper(data.Plural)
	DatabaseConnection.Where(Securities{ID: data.ID}).Assign(Securities{Name: data.Name, Plural: data.Plural}).FirstOrCreate(&Securities{Name: data.Name, Plural: data.Plural})
}

// AddAssetClass Adding Asset Class
func AddAssetClass(data AssetClass) {
	data.SecurityType = strings.ToUpper(data.SecurityType)
	data.AssetClassName = strings.ToUpper(data.AssetClassName)
	DatabaseConnection.Create(&data)
}

// EditAssetClass Edit Asset Class
func EditAssetClass(data AssetClass) {
	data.SecurityType = strings.ToUpper(data.SecurityType)
	data.AssetClassName = strings.ToUpper(data.AssetClassName)
	DatabaseConnection.Where(AssetClass{ID: data.ID}).Assign(AssetClass{SecurityType: data.SecurityType, AssetClassName: data.AssetClassName}).FirstOrCreate(&AssetClass{SecurityType: data.SecurityType, AssetClassName: data.AssetClassName})
}

// DeleteAssetClass Delete Asset Class
func DeleteAssetClass(id int) {
	DatabaseConnection.Exec("DELETE FROM CRS_ASSET_CLASS WHERE ID=?", id)
}

func LoadSecurities() []Securities {
	securities := make([]Securities, 0)
	DatabaseConnection.Find(&securities)
	return securities
}

func DeleteSecurity(id int) {
	DatabaseConnection.Exec("DELETE FROM CRS_SECURITIES_NAME WHERE ID=?", id)
}

func LoadMergedSecurityNames() []MergedSecurityName {
	securities := make([]MergedSecurityName, 0)
	DatabaseConnection.Raw("select NAME as securities from CRS_SECURITIES_NAME UNION (select PLURAL from CRS_SECURITIES_NAME)").Scan(&securities)
	return securities
}

func LoadClientEmailLog() []QuarterEmailToClients {
	log := make([]QuarterEmailToClients, 0)
	DatabaseConnection.Find(&log)
	return log
}

func LogClientLettersSent(sentBy string) {
	today := time.Now()
	month := int(today.Month())
	year := today.Year()
	var fancyDate string
	if month >= 1 && month <= 3 {
		fancyDate = fmt.Sprintf("Q1 %d", year)
	} else if month >= 4 && month <= 6 {
		fancyDate = fmt.Sprintf("Q2 %d", year)
	} else if month >= 7 && month <= 9 {
		fancyDate = fmt.Sprintf("Q3 %d", year)
	} else {
		fancyDate = fmt.Sprintf("Q4 %d", year)
	}
	DatabaseConnection.Create(&QuarterEmailToClients{
		QuarterYear: fancyDate,
		SentBy:      sentBy,
		SentOn:      time.Now(),
		CreatedAt:   time.Now(),
	})
}

func ClientLettersHaveBeenSent() bool {
	today := time.Now()
	month := int(today.Month())
	year := today.Year()
	var startDateStr, endDateStr string
	if month >= 1 && month <= 3 {
		startDateStr = fmt.Sprintf(`%d-01-01`, year)
		endDateStr = fmt.Sprintf(`%d-03-31`, year)
	} else if month >= 4 && month <= 6 {
		startDateStr = fmt.Sprintf(`%d-04-01`, year)
		endDateStr = fmt.Sprintf(`%d-06-30`, year)
	} else if month >= 7 && month <= 9 {
		startDateStr = fmt.Sprintf(`%d-07-01`, year)
		endDateStr = fmt.Sprintf(`%d-09-30`, year)
	} else {
		startDateStr = fmt.Sprintf(`%d-10-01`, year)
		endDateStr = fmt.Sprintf(`%d-12-31`, year)
	}
	startDate, _ := time.Parse("2006-01-02", startDateStr)
	endDate, _ := time.Parse("2006-01-02", endDateStr)

	log := make([]QuarterEmailToClients, 0)
	DatabaseConnection.Where("SENT_ON BETWEEN ? AND ?", startDate, endDate).Find(&log)
	return len(log) > 0
}

// Load0302DataByMonth loads the pv report data for a given date by the bpid
func Load0302DataByMonth(BPID, Date string) []NPRA0302 {
	data := make([]NPRA0302, 0)
	DatabaseConnection.Raw("select * from CRS_0302_NPRA_REPORT where bp_id=? and reporting_date=? order by created_at asc", BPID, Date).Scan(&data)
	return data
}

// Load0301DataByMonth loads the pv report data for a given date by the bpid
func Load0301DataByMonth(Date string) []NPRA0301 {
	data := make([]NPRA0301, 0)
	DatabaseConnection.Raw("select * from CRS_0301_NPRA_REPORT where reporting_date=? order by created_at asc", Date).Scan(&data)
	return data
}

func Loadnpra0302Data(BPID, Date string) []Auditable0302 {
	np0302 := make([]Auditable0302, 0)
	codes := make([]NPRA0302, 0)
	DatabaseConnection.Select("bp_id").Where("reporting_date=? OR bp_id=? ", Date, BPID).Preload("ClientSCADetails").Find(&codes)
	if len(codes) > 0 {
		for _, code := range codes {
			report0302 := Load0302DataByMonth(code.BP_ID, Date)
			client := code.ClientSCADetails
			if client.ID == 0 {
				//for some reason this happens
				client = GetClientByCode(code.BP_ID)
			}
			np0302 = append(np0302, Auditable0302{Client: client, Report0302: report0302})
		}
	}
	return np0302
}
func Loadnpra0301Data(Date string) []Auditable0301 {
	np0301 := make([]Auditable0301, 0)
	codes := make([]NPRA0301, 0)
	DatabaseConnection.Select("entity_name").Where("reporting_date=? ", Date).Find(&codes)
	if len(codes) > 0 {
		report0301 := Load0301DataByMonth(Date)
		np0301 = append(np0301, Auditable0301{Report0301: report0301})
	}
	return np0301
}

// Load301data loads the sec local variance data
func Load301data(date string) []NPRA0301 {

	npra301 := make([]NPRA0301, 0)
	DatabaseConnection.Raw(" SELECT net_return, gross_return, report_code, entity_id, entity_name, reference_period_year, reference_period, investment_receivables, total_asset_under_management, government_securities, local_government_securities, corporate_debt_securities, bank_securities, ordinary_preference_shares, collective_investment_scheme, alternative_investments, bank_balances, reporting_date FROM CRS_0301_NPRA_REPORT WHERE reporting_date = ? ORDER BY created_at ASC", date).Scan(&npra301)
	data := make([]NPRA0301, len(npra301))
	for index, datum := range npra301 {

		data[index] = NPRA0301{
			NetReturn:                  datum.NetReturn,
			GrossReturn:                datum.GrossReturn,
			ReportCode:                 datum.ReportCode,
			EntityID:                   datum.EntityID,
			EntityName:                 datum.EntityName,
			ReferencePeriodYear:        datum.ReferencePeriodYear,
			ReferencePeriod:            datum.ReferencePeriod,
			InvestmentReceivables:      datum.InvestmentReceivables,
			TotalAssetUnderManagement:  datum.TotalAssetUnderManagement,
			GovernmentSecurities:       datum.GovernmentSecurities,
			LocalGovernmentSecurities:  datum.LocalGovernmentSecurities,
			CorporateDebtSecurities:    datum.CorporateDebtSecurities,
			BankSecurities:             datum.BankSecurities,
			OrdinaryPreferenceShares:   datum.OrdinaryPreferenceShares,
			CollectiveInvestmentScheme: datum.CollectiveInvestmentScheme,
			AlternativeInvestments:     datum.AlternativeInvestments,
			BankBalances:               datum.BankBalances,
		}
	}
	return data
}

// Load302data loads the sec local variance data
func Load302data(date, bpOrSca string) []NPRA0302 {
	fmt.Println("$$$$$$$$$$Load Data $$$$$$$$$$$$$$$$$$$$$", date)
	fmt.Println("$$$$$$$$$$Load Data $$$$$$$$$$$$$$$$$$$$$", bpOrSca)
	npra302 := make([]NPRA0302, 0)
	DatabaseConnection.Raw(" SELECT report_code, entity_id, entity_name, reference_period_year, reference_period, investment_id, instrument, issuer_name, asset_tenure, reporting_date, maturity_date, face_value, issue_date, currency, asset_class, market_value FROM CRS_0302_NPRA_REPORT WHERE reporting_date = ? and bp_id = ? ORDER BY created_at ASC", date, bpOrSca).Scan(&npra302)
	data := make([]NPRA0302, len(npra302))
	for index, datum := range npra302 {

		data[index] = NPRA0302{
			ReportCode:          datum.ReportCode,
			EntityID:            datum.EntityID,
			EntityName:          datum.EntityName,
			ReferencePeriodYear: datum.ReferencePeriodYear,
			ReferencePeriod:     datum.ReferencePeriod,
			InvestmentID:        datum.InvestmentID,
			Instrument:          datum.Instrument,
			IssuerName:          datum.IssuerName,
			AssetTenure:         datum.AssetTenure,
			ReportingDate:       datum.ReportingDate,
			MaturityDate:        datum.MaturityDate,
			FaceValue:           datum.FaceValue,
			IssueDate:           datum.IssueDate,
			Currency:            datum.Currency,
			AssetClass:          datum.AssetClass,
			MarketValue:         datum.MarketValue,
		}
	}
	return data
}
