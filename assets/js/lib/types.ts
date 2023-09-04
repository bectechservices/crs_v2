export type MaybeNull<T> = T | null;

export interface ParsedPVReport {
    headers: Array<string>;
    bonds: Array<BondData>;
    summary: Array<Array<string>>;
}

export interface BondData {
    bond: string;
    values: Array<Array<string>>;
}

export interface IKeyable {
    key: number;
}

export interface IFile {
    name: string;
    data: string;
    rawBase64Data: string;
}

export type ICanDuplicate = IKeyable;

export interface IShare {
    name: string;
    shareholding: number;
    percentage: number;
}

export type IKeyableShare = IKeyable & IShare;

export type YesNo = "yes" | "no" | "";
export type YesNo_NA = YesNo | "n/a" | "";

export interface IGovernance {
    clientName: string;
    reportingOfficer: string;
    reportingDate: string;
    ordinaryShares: Array<IShare>;
    preferenceShares: Array<IShare>;
    custodianTransactions: Array<ICustodianTransaction>;
    schemesUnderCustody: Array<ISchemesUnderCustody>;
    changeInDirectors: YesNo;
    changeInAgreement: YesNo;
    dealingsApprovedByBoard: YesNo_NA;
    custodianHasUpdatedAssetRegister: YesNo;
    custodianAssetRegistrationDate: string;
    doManagersOfTheSchemeConsultTheLaw: YesNo;
    schemeHadAnyOtherFinancialDealings: YesNo;
    approved: boolean;
}

export interface RootState {
    loading: boolean;
    pvUploadErrors: Array<string>
}

export interface ICustodianTransaction {
    nameOfTrustee: string;
    relationshipWithTrustee: string;
    typeOfTransaction: string;
    amount: number;
}

export type IKeyableCustodianTransaction = ICustodianTransaction & IKeyable;

export interface ISchemesUnderCustody {
    nameOfFirm: string;
    nameOfScheme: string;
    relationshipWithTrustee: string;
    volume: number;
    markedToMarketValue: number;
}

export type IKeyableSchemesUnderCustody = ISchemesUnderCustody & IKeyable;

export interface IError {
    hasError: boolean;
}

export interface IKeyableShareError extends IKeyable {
    name: boolean;
    shareholding: boolean;
    percentage: boolean;
}

export interface IKeyableCustodianTransactionError extends IKeyable {
    nameOfTrustee: boolean;
    relationshipWithTrustee: boolean;
    typeOfTransaction: boolean;
    amount: boolean;
}

export interface IKeyableSchemesUnderCustodyError extends IKeyable {
    nameOfFirm: boolean;
    nameOfScheme: boolean;
    relationshipWithTrustee: boolean;
    volume: boolean;
    markedToMarketValue: boolean;
}

export interface CustodianState {
    ordinarySharesInputData: Array<IKeyableShare>;
    ordinarySharesInputError: Array<IKeyableShareError>;
    preferenceSharesInputData: Array<IKeyableShare>;
    preferenceSharesInputError: Array<IKeyableShareError>;
    custodianTransactionsInputData: Array<IKeyableCustodianTransaction>;
    custodianTransactionsInputError: Array<IKeyableCustodianTransactionError>;
    schemesUnderCustodyInputData: Array<IKeyableSchemesUnderCustody>;
    schemesUnderCustodyInputError: Array<IKeyableSchemesUnderCustodyError>;
    fieldHasError: boolean;
}

export interface PPTState {
    templateOptions: PPTTemplate;
    bpid: string;
    quarter: number;
    year: number;
}

export interface IGovernanceError {
    clientName: boolean;
    reportingOfficer: boolean;
    reportingDate: boolean;
    custodianAssetRegistrationDate: boolean;
}

export interface IMultipleGovernanceFieldsData {
    ordinarySharesInputData: Array<IKeyableShare>;
    preferenceSharesInputData: Array<IKeyableShare>;
    custodianTransactionsInputData: Array<IKeyableCustodianTransaction>;
    schemesUnderCustodyInputData: Array<IKeyableSchemesUnderCustody>;
}

export interface ApiResponse {
    error: boolean;
    message: string;
    data?: Array<any>;
}

export interface ISchemeDetails {
    bpid: string;
    nameOfScheme: string;
    numberOfSharesOutstanding: number;
    numberOfShareholders: number;
    numberOfRedemptions: number;
    valueOfRedemptions: number;
    nameOfManager: string;
    totalValueOfSchemeAssets: number;
    netAssetValueOfScheme: number;
    netAssetValuePerShare: number;
    totalEquityInvestments: number;
    totalFixedIncomeInvestments: number;
    netOfMediumTermAssetsHeldByFund: number;
    capitalMarketsInvestments: number;
    percentageOfCapitalInvestmentToTotalInvestment: number;
    areAllCertificatesOfInvestmentWithCustodian: YesNo;
    totalValueOfUnutilizedFunds: number;
    valueOfBorrowedFunds: number;
    reasonsForBorrowing: string;
    wereAllDulyPreparedAccountsDistributed: YesNo;
    redemptions: number;
    dividends: number;
    rights: number;
    feesOwedCustodian: number;
    attachedFile?: string;
}

export interface IOtherSecInformations {
    areThereAnyClaimOnSchemeAsset: YesNo_NA;
    yesWasCustodianInformedAndApproved: YesNo;
    anyLitigationInvolvingCustodianSccheme: YesNo_NA;
    anySignificantReductionInAssetScheme: YesNo;
    hasMgrsReconciledAssetRegisterCustodian: YesNo;
    significantReductionInSchemeMarketPrice: YesNo;
    howManyTimesDidSchemePublishedPrices: YesNo_NA;
    anyConcernsByInvestors: YesNo;
    anyMattersAttentionSecMgtCustodyOfFund: YesNo_NA;
    hasAccountOfManagersSeparateFromScheme: YesNo;
    companyParentsAffiliateInvolvedInLitigation: YesNo_NA;
    litigationDetails: string;
}

export interface IOfficialRemarks {
    remarks: string;
    reviewingOfficer: string;
    date: string;
    signature?: string;
}

export interface ISchemeDetailsError {
    bpid: boolean;
    nameOfScheme: boolean;
    numberOfSharesOutstanding: boolean;
    numberOfShareholders: boolean;
    numberOfRedemptions: boolean;
    valueOfRedemptions: boolean;
    nameOfManager: boolean;
    totalValueOfSchemeAssets: boolean;
    netAssetValueOfScheme: boolean;
    netAssetValuePerShare: boolean;
    totalEquityInvestments: boolean;
    totalFixedIncomeInvestments: boolean;
    netOfMediumTermAssetsHeldByFund: boolean;
    capitalMarketsInvestments: boolean;
    percentageOfCapitalInvestmentToTotalInvestment: boolean;
    totalValueOfUnutilizedFunds: boolean;
    valueOfBorrowedFunds: boolean;
    reasonsForBorrowing: boolean;
    redemptions: boolean;
    dividends: boolean;
    rights: boolean;
    feesOwedCustodian: boolean;
}

export interface IOtherSecInformationsError {
    litigationDetails: boolean;
}

export interface IOfficialRemarksError {
    remarks: boolean;
    reviewingOfficer: boolean;
    date: boolean;
}

export interface SchemeDetailsRequest {
    bpid: string;
    quarterDate: string;
}

export interface SecPerformance {
    security: string;
    value: number;
}

export interface PPTTemplate {
    total_summary_of_auc: boolean;
    auc_trend: boolean;
    trade_volumes: boolean;
    pv_report: boolean;
    total_contribution: boolean;
    corporate_action: boolean;
    gog_and_fd_maturities: boolean;
    appendix_i: boolean;
    appendix_ii: boolean;
    unidentified_payments: boolean;
}

export interface OffshoreClient {
    name: string;
    country: string;
    assetValue: number;
}

export interface QuaterYear {
    quarter: string;
    year: string;
}

export interface BPIDSearch extends QuaterYear {
    bpid: string;
}

export interface VarianceData {
    name: string;
    country: string;
    last_aua: string;
    current_aua: string;
    amount: string;
    variance: string;
}

export interface PVReport {
    account_number: string;
    security_name: string
    cds_code: string
    isin: string
    scb_code: string
    market_price: number
    nominal_value: number
    cumulative_cost: number
    value_amount: number
    security_type: string
    percentage_of_total: number
    date_from: string
    date_to: string
    report_date: string
}

export interface ClientDetails {
    name: string;
    bpid: string;
    safekeepingAccount: string;
    client_id: string;
}

export interface PVSummary {
    description: string;
    norminal_value: number;
    cummulative_cost: number;
    value: number;
    percentage_of_total: number;
}

export interface IMaturable {
    is_matured: boolean;
}

export interface TrusteePerformance {
    bond: string;
    current_quarter: number;
    previous_quarter: number;
}

export interface IMonthlyContributions {
    id: number;
    bpid: string;
    date: string;
    amount: number;
    sca: string;
    created_at: string;
}

export interface FundManager {
    id: number;
    name: string;
    account_number: string;
    administrator?: string;
}

export interface FundAdministrator {
    name: string;
    administrator: string;
}

export interface IUnidentifiedPayment {
    key: number;
    client_bpid: string;
    txn_date: string;
    txn_type: string;
    value_date: string;
    name_of_company: string;
    amount: number;
    fund_manager: number;
    collection_acc_num: string;
    status: "pending" | "done"
}

export interface TrusteeState {
    unidentifiedPayments: Array<IUnidentifiedPayment>;
    searchInputData: TrusteeDashboardSearchData;
}
export interface npra321State {
    search321InputData: npra321SearchData;
}

export interface npra321SearchData {
    bpid: string;
    month: string;
    year: string;
}
export interface TrusteeDashboardSearchData {
    bpid: string;
    quarter: string;
    year: string;
}

export interface IUnidentifiedPaymentSummary {
    txn_type: string;
    total: number;
    pending: number;
    done: number;
}

export interface IUnauthorizedTransaction {
    key: number;
    clientName: string;
    txnDetails: string;
    date: string;
}


export interface NPRAState {
    unauthorizedTransactions: Array<IUnauthorizedTransaction>;
}

export interface IOutstandingFDCertificate {
    fundManager: string;
    clientName: string;
    amount: number;
    issuer: string;
    rate: number;
    tenor: number;
    term: string;
    effectiveDate: string;
    maturity: string;
}

export interface IFund {
    name: string;
    value: number;
    bpid: string;
}

export interface GroupedOutstandingFDCertificate {
    month: string;
    certificates: Array<IOutstandingFDCertificate>
}

export interface TransactionDetails {
    bpid: string;
    reporting_date: string;
    sca: string;
    date: string;
    reference: string;
    security_name: string;
    security_category: string;
    charge_type: string;
    charge_item: string;
    number_of_units: number;
    market_value: number;
    charge_amount_with_tax: number;
    invoice_amount_with_tax: number;
}

export interface RangedBasisPoint extends IKeyable {
    minimum_amount: number;
    maximum_amount: number;
    basis_point: number;
}

export interface ClientState {
    rangedBasisPoint: Array<RangedBasisPoint>
}

export interface INPRADeclaration {
    nameOfOfficer: string;
    designation: string;
    headOfCustodyServices: string;
    date: string
}

export interface TrusteeUploadedPV {
    client: string;
    bpid: string;
    uploaded_by: string;
    date: string;
    type: string;
    quarter_date: string;
}

export interface TrusteeQuarterlyReport {
    id: number;
    bpid: string;
    quarter: string;
    approved: boolean;
    approved_by: number
}

export interface BillingQuarterlyReport {
    id: number;
    client_id: string;
    month: string;
    approved: boolean;
    approved_by: number
}

export interface BillingState {
    searchData: BillingSearchData
}

export interface BillingSearchData {
    bpOrSca: string;
    date: string;
}

export interface FormattedPVData {
    bond: string;
    data: Array<PVReport & IMaturable>
}

export interface AuditablePV {
    client: { client_name: string; code: string; };
    reports: Array<FormattedPVData>;
}


export interface Npra0301 {
    report_code: string;
    entity_id: string
    entity_name: string
    reference_period_year: string
    reference_period: string
    net_return: number
    gross_return: number
    investment_receivables: number
    total_asset_under_management: number
    government_securities: number
    local_government_securities: number
    corporate_debt_securities: number
    bank_securities: number
    ordinary_preference_shares: number
    collective_investment_scheme: number
    alternative_investments: number
    bank_balances: number
    reporting_date: string
}


// export interface 0302 {
//     report_code: string;
//     entity_id: string
//     entity_name: string
//     reference_period_year: string
//     reference_period: string
//     investment_id: string
//     instrument: string
//     issuer_name: string
//     asset_tenure: string
//     reporting_date: string
//     amount_invested: number
//     accrued_coupon_interest: number
//     coupon_paid: number
//     accrued_coupon_interest_since_last_year: number
//     outstanding_interest_maturity: number
//     amount_impaired: number
//     asset_allocation_actual_percent: number
//     type_investment_charge: number
//     investment_charge_date_percent: number
//     investment_charge_mount: number
//     face_value: number
//     interest_rate_percent: number
//     discount_rate_percent: number
//     disposal_proceeds: number
//     disposal_instructions: string
//     yield_disposal: number
//     issue_date: string
//     price_share_unit_purchase: number
//     price_share_unit_value_date: string
//     capital_gains: number
//     dividend_received: number
//     number_units_shares: number
//     holding_period_return_per_investment: number
//     days_run: number
//     currency_conversion_rate: number
//     currency: string
//     amount_invested_foreign_currency: number
//     asset_class: string
//     price_share_unit_value_date_foreign: number
//     market_value: number
//     remaining_days_maturity: number
//     holding_period_per_investment_weighted_percentage: number
// }
export interface Npra0302 {
    bp_id: string;
    report_code: string
    entity_id: string
    entity_name: string
    reference_period_year: string
    reference_period: string
    investment_id: string
    instrument: string
    issuer_name: string
    asset_tenure: string
    reporting_date: string
    face_value: number
    currency: string
    asset_class: string
    maturity_date: number
    market_value: number
}


export interface Formatted301Data {
    data: Array<Npra0301 & IMaturable>
}

export interface Auditable0301 {
    report0301: Array<Formatted301Data>;
}

export interface Formatted302Data {
    bpid: string;
    data: Array<Npra0302 & IMaturable>
}

export interface Auditable0302 {
    client: { client_name: string; code: string; };
    report0302: Array<Formatted302Data>;
}

export interface Npra0303 {
    report_code: string;
    entity_id: string
    entity_name: string
    reference_period_year: string
    reference_period: string
    unit_price: number
    date_valuation: string
    unit_number: number
    daily_nav: number
    npra_fees: number
    trustee_fees: number
    fund_manager_fees: number
    fund_custodian_fees: number
}