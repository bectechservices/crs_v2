import Vue from "vue";
import axios from "axios";
import {format} from "date-fns";
import {
    BPIDSearch,
    IGovernance,
    IOfficialRemarks,
    IOtherSecInformations,
    ISchemeDetails,
    OffshoreClient,
    QuaterYear,
    VarianceData
} from "./lib/types";
import {
    CURRENT_YEAR,
    currentQuarterNumber,
    formatMoney,
    getPreviousQuarterShortDate,
    getShortDate,
    makeCurrentQuarter,
    makeLastSecLicenseRenewalDate,
    makeLastDayOfQuarter,
    makePreviousQuarterShortDate,
    makeQuarterFormalDate,
    makeShortDate
} from "./lib/helpers";
//@ts-ignore
import download from "downloadjs"
import VueSweetalert2 from 'vue-sweetalert2';


Vue.use(VueSweetalert2);

interface MaturedSecurity {
    client: string;
    issuer: string;
    amount_invested: number;
    value: number;
}

interface Data {
    offshoreClients: Array<OffshoreClient>;
    offshoreClientSearchData: BPIDSearch;
    localVariance: Array<VarianceData>;
    foreignVariance: Array<VarianceData>;
    localVarianceSearchData: QuaterYear;
    foreignVarianceSearchData: QuaterYear;
    foreignVarianceLastAUADate: string;
    foreignVarianceCurrentAUADate: string;
    localVarianceLastAUADate: string;
    localVarianceCurrentAUADate: string;
    quarterlyReport: { reportingPeriod: string, quarterDate: string, schemes: Array<ISchemeDetails>, report: IGovernance, info: IOtherSecInformations, remarks: IOfficialRemarks };
    quarterlyReportSearchData: QuaterYear;
    schemes: ISchemeDetails;
    schemeSearchData: {
        bpid: string;
        quarter: string;
        year: string;
        quarterDate: string
    },
    secReportRegistrationRenewalDate: string;
    maturedSecurities: {
        quarter: string;
        year: string;
        formalDate: string;
    },
    maturedSecuritiesData: Array<MaturedSecurity>
}

interface Methods {
    fetchOffshoreClients: (
        bpid: string,
        quarter: string,
        year: string,
        showToast?: boolean
    ) => Promise<void>;
    loadOffshoreClients: () => Promise<void>;
    loadLocalVariance: (quarter: string, year: string, showToast?: boolean) => Promise<void>;
    loadForeignVariance: (quarter: string, year: string, showToast?: boolean) => Promise<void>;
    fetchLocalVariance: () => Promise<void>;
    fetchForeignVariance: () => Promise<void>;
    formatMoney: (amount: number, dp?: number) => string;
    loadQuarterlyReport: (showToast: boolean) => Promise<void>;
    loadClientSchemeDetails: () => Promise<void>;
    exportSecQuarterlyReport: () => Promise<void>;
    exportOffshoreClients: () => Promise<void>;
    exportLocalVariance: () => Promise<void>;
    exportForeignVariance: () => Promise<void>;
    exportMaturedSecurities: () => Promise<void>;
    handleApiResponse: (response: { error: boolean }, successMessage: string) => void;
    formatDate: (date: string, format: string) => void;
    loadMaturedSecurities: (showToast: boolean) => Promise<void>;
}

export default new Vue<Data, Methods>({
    el: ".secReportPage",
    data: {
        offshoreClients: [],
        offshoreClientSearchData: {
            bpid: "",
            quarter: currentQuarterNumber(),
            year: CURRENT_YEAR
        },
        localVariance: [],
        foreignVariance: [],
        localVarianceSearchData: {
            quarter: currentQuarterNumber(),
            year: CURRENT_YEAR
        },
        foreignVarianceSearchData: {
            quarter: currentQuarterNumber(),
            year: CURRENT_YEAR
        },
        foreignVarianceLastAUADate: getPreviousQuarterShortDate(),
        foreignVarianceCurrentAUADate: getShortDate(),
        localVarianceLastAUADate: getPreviousQuarterShortDate(),
        localVarianceCurrentAUADate: getShortDate(),
        quarterlyReport: {
            reportingPeriod: "",
            quarterDate: "",
            schemes: [],
            report: {
                clientName: "",
                reportingOfficer: "",
                reportingDate: "",
                ordinaryShares: [],
                preferenceShares: [],
                schemesUnderCustody: [],
                custodianTransactions: [],
                changeInDirectors: "",
                changeInAgreement: "",
                dealingsApprovedByBoard: "",
                custodianHasUpdatedAssetRegister: "",
                custodianAssetRegistrationDate: "",
                doManagersOfTheSchemeConsultTheLaw: "",
                schemeHadAnyOtherFinancialDealings: "",
                approved: false
            },
            info: {
                areThereAnyClaimOnSchemeAsset: "",
                yesWasCustodianInformedAndApproved: "",
                anyLitigationInvolvingCustodianSccheme: "",
                anySignificantReductionInAssetScheme: "",
                hasMgrsReconciledAssetRegisterCustodian: "",
                significantReductionInSchemeMarketPrice: "",
                howManyTimesDidSchemePublishedPrices: "",
                anyConcernsByInvestors: "",
                anyMattersAttentionSecMgtCustodyOfFund: "",
                hasAccountOfManagersSeparateFromScheme: "",
                companyParentsAffiliateInvolvedInLitigation: "",
                litigationDetails: ""
            },
            remarks: {
                remarks: "",
                reviewingOfficer: "",
                date: "",
                signature: ""
            }
        },
        quarterlyReportSearchData: {
            quarter: currentQuarterNumber(),
            year: CURRENT_YEAR
        },
        schemes: {
            bpid: "",
            nameOfScheme: "",
            numberOfSharesOutstanding: 0,
            numberOfShareholders: 0,
            numberOfRedemptions: 0,
            valueOfRedemptions: 0,
            nameOfManager: "",
            totalValueOfSchemeAssets: 0,
            netAssetValueOfScheme: 0,
            netAssetValuePerShare: 0,
            totalEquityInvestments: 0,
            totalFixedIncomeInvestments: 0,
            netOfMediumTermAssetsHeldByFund: 0,
            capitalMarketsInvestments: 0,
            percentageOfCapitalInvestmentToTotalInvestment: 0,
            areAllCertificatesOfInvestmentWithCustodian: "no",
            totalValueOfUnutilizedFunds: 0,
            valueOfBorrowedFunds: 0,
            reasonsForBorrowing: "",
            wereAllDulyPreparedAccountsDistributed: "no",
            redemptions: 0,
            dividends: 0,
            rights: 0,
            feesOwedCustodian: 0
        },
        schemeSearchData: {
            bpid: "",
            quarter: currentQuarterNumber(),
            year: CURRENT_YEAR,
            quarterDate: ""
        },
        maturedSecurities: {
            quarter: currentQuarterNumber(),
            year: CURRENT_YEAR,
            formalDate: makeQuarterFormalDate(currentQuarterNumber(), CURRENT_YEAR)
        },
        maturedSecuritiesData: [],
        secReportRegistrationRenewalDate: makeLastSecLicenseRenewalDate(currentQuarterNumber(), CURRENT_YEAR)
    },
    beforeMount: async function () {
        try {
            await Promise.all([
                this.fetchOffshoreClients(
                    "",
                    this.offshoreClientSearchData.quarter,
                    this.offshoreClientSearchData.year,
                    false
                ),
                this.loadForeignVariance(
                    this.foreignVarianceSearchData.quarter,
                    this.foreignVarianceSearchData.year,
                    false
                ),
                this.loadLocalVariance(
                    this.localVarianceSearchData.quarter,
                    this.localVarianceSearchData.year,
                    false
                ),
                this.loadMaturedSecurities(false),
                this.loadQuarterlyReport(false)
            ]);
        } catch (e) {
            console.error(e);
        }
    },
    methods: {
        handleApiResponse: function (response: { error: boolean }, successMessage: string) {
            if (response.error) {
                console.log(response);
                this.$swal({
                    toast: true,
                    position: 'top-end',
                    showConfirmButton: false,
                    timer: 3000,
                    type: 'error',
                    title: "Something went wrong. Please try again."
                }).then(console.log)
            } else {
                this.$swal({
                    toast: true,
                    position: 'top-end',
                    showConfirmButton: false,
                    timer: 3000,
                    type: 'success',
                    title: successMessage
                }).then(console.log)
            }
        },
        loadClientSchemeDetails: async function () {
            let response = await axios.post("/fetch-scheme-details", {
                quarterDate: makeLastDayOfQuarter(this.schemeSearchData.quarter, this.schemeSearchData.year),
                bpid: this.schemeSearchData.bpid
            });
            this.handleApiResponse(response.data, "Scheme details loaded");
            this.schemes = response.data.data.schemeDetails;
            this.schemeSearchData.quarterDate = makeCurrentQuarter(this.schemeSearchData.quarter, this.schemeSearchData.year)
        },
        loadQuarterlyReport: async function (showToast: boolean = true) {
            let response = await axios.post("/sec-quarterly-report", {
                year: this.quarterlyReportSearchData.year,
                quarter: this.quarterlyReportSearchData.quarter
            });
            if (showToast) {
                this.handleApiResponse(response.data, "Report loaded");
            }
            this.secReportRegistrationRenewalDate = makeLastSecLicenseRenewalDate(this.quarterlyReportSearchData.quarter, this.quarterlyReportSearchData.year)
            this.quarterlyReport.reportingPeriod = makeLastDayOfQuarter(this.quarterlyReportSearchData.quarter, this.quarterlyReportSearchData.year);
            this.quarterlyReport.quarterDate = makeCurrentQuarter(this.quarterlyReportSearchData.quarter, this.quarterlyReportSearchData.year);
            this.quarterlyReport.schemes = response.data.data.schemeDetails.map((scheme: ISchemeDetails) => {
                return {
                    ...scheme,
                    numberOfSharesOutstanding: this.formatMoney(scheme.numberOfSharesOutstanding),
                    valueOfRedemptions: this.formatMoney(scheme.valueOfRedemptions),
                    netAssetValueOfScheme: this.formatMoney(scheme.netAssetValueOfScheme, 10),
                    netAssetValuePerShare: this.formatMoney(scheme.netAssetValuePerShare, 10),
                    totalEquityInvestments: this.formatMoney(scheme.totalEquityInvestments),
                    totalFixedIncomeInvestments: this.formatMoney(scheme.totalFixedIncomeInvestments),
                    totalValueOfUnutilizedFunds: this.formatMoney(scheme.totalValueOfUnutilizedFunds),
                    valueOfBorrowedFunds: this.formatMoney(scheme.valueOfBorrowedFunds),
                    feesOwedCustodian: this.formatMoney(scheme.feesOwedCustodian),
                    redemptions: this.formatMoney(scheme.redemptions),
                    rights: this.formatMoney(scheme.rights),
                    dividends: this.formatMoney(scheme.dividends)
                }
            });
            this.quarterlyReport.report = {
                ...response.data.data.report,
                ordinaryShares: response.data.data.report.ordinaryShares ? response.data.data.report.ordinaryShares.map((each: any) => {
                    return {
                        name: each.name,
                        shareholding: this.formatMoney(each.shareholdings),
                        percentage: each.percentage
                    }
                }) : [],
                preferenceShares: response.data.data.report.preferenceShares ? response.data.data.report.preferenceShares.map((each: any) => {
                    return {
                        name: each.name,
                        shareholding: this.formatMoney(each.shareholdings),
                        percentage: each.percentage
                    }
                }) : [],
                dateOfReport: format(response.data.data.report.dateOfReport, "Do MMMM YYYY")
            };
            this.quarterlyReport.info = response.data.data.otherInformation;
            this.quarterlyReport.remarks = response.data.data.remarks;
        },
        formatMoney: function (amount: number, dp: number = 2) {
            return formatMoney(amount, dp);
        },
        fetchForeignVariance: async function () {
            await this.loadForeignVariance(this.foreignVarianceSearchData.quarter, this.foreignVarianceSearchData.year);
        },
        fetchLocalVariance: async function () {
            await this.loadLocalVariance(this.localVarianceSearchData.quarter, this.localVarianceSearchData.year)
        },
        loadLocalVariance: async function (quarter: string, year: string, showToast: boolean = true) {
            const response = await axios.post("/load-sec-local-variance", {
                year,
                quarter
            }, {
                headers: {
                    "Content-Type": "application/json"
                }
            });
            if (showToast) {
                this.handleApiResponse(response.data, "Local variance loaded");
            }
            this.localVariance = response.data.data.variance;
            this.localVarianceCurrentAUADate = makeShortDate(quarter, year);
            this.localVarianceLastAUADate = makePreviousQuarterShortDate(quarter, year)
        },
        loadForeignVariance: async function (quarter: string, year: string, showToast: boolean = true) {
            const response = await axios.post("/load-sec-foreign-variance", {
                year,
                quarter
            }, {
                headers: {
                    "Content-Type": "application/json"
                }
            });
            if (showToast) {
                this.handleApiResponse(response.data, "Foreign variance loaded");
            }
            this.foreignVariance = response.data.data.variance;
            this.foreignVarianceCurrentAUADate = makeShortDate(quarter, year);
            this.foreignVarianceLastAUADate = makePreviousQuarterShortDate(quarter, year)
        },
        loadOffshoreClients: async function () {
            await this.fetchOffshoreClients(
                this.offshoreClientSearchData.bpid,
                this.offshoreClientSearchData.quarter,
                this.offshoreClientSearchData.year
            );
        },
        fetchOffshoreClients: async function (
            bpid: string,
            quarter: string,
            year: string,
            showToast: boolean = true
        ) {
            const response = await axios.post("/offshore-clients", {
                bpid,
                quarter,
                year
            });
            if (showToast) {
                this.handleApiResponse(response.data, "Client loaded");
            }
            this.offshoreClients = response.data.data.clients;
        },
        exportSecQuarterlyReport: async function () {
            const response = await axios.post("/export-sec-report", {
                year: this.quarterlyReportSearchData.year,
                quarter: this.quarterlyReportSearchData.quarter
            }, {
                headers: {
                    'Content-Type': 'application/json'
                },
                responseType: 'blob'
            });
            download(response.data, `SEC-Report Q${this.quarterlyReportSearchData.quarter}-${this.quarterlyReportSearchData.year}.docx`);
        },
        exportOffshoreClients: async function () {
            const response = await axios.post("/offshore-clients-excel", {
                year: this.offshoreClientSearchData.year,
                quarter: this.offshoreClientSearchData.quarter
            }, {
                headers: {
                    'Content-Type': 'application/json'
                },
                responseType: 'blob'
            });
            download(response.data, `Offshore Clients Q${this.offshoreClientSearchData.quarter}-${this.offshoreClientSearchData.year}.xlsx`);
        },
        exportLocalVariance: async function () {
            const response = await axios.post("/sec-local-variance-excel", {
                year: this.localVarianceSearchData.year,
                quarter: this.localVarianceSearchData.quarter
            }, {
                headers: {
                    'Content-Type': 'application/json'
                },
                responseType: 'blob'
            });
            download(response.data, `Local Variance Q${this.localVarianceSearchData.quarter}-${this.localVarianceSearchData.year}.xlsx`);
        },
        exportForeignVariance: async function () {
            const response = await axios.post("/sec-foreign-variance-excel", {
                year: this.foreignVarianceSearchData.year,
                quarter: this.foreignVarianceSearchData.quarter
            }, {
                headers: {
                    'Content-Type': 'application/json'
                },
                responseType: 'blob'
            });
            download(response.data, `Foreign Variance Q${this.foreignVarianceSearchData.quarter}-${this.foreignVarianceSearchData.year}.xlsx`);
        },
        exportMaturedSecurities: async function () {
            this.maturedSecurities.formalDate = makeQuarterFormalDate(this.maturedSecurities.quarter, this.maturedSecurities.year)
            const response = await axios.post("/matured-securities-excel", {
                quarter: this.maturedSecurities.formalDate
            }, {
                headers: {
                    'Content-Type': 'application/json'
                },
                responseType: 'blob'
            });
            download(response.data, `Matured Securities Q${this.maturedSecurities.quarter}-${this.maturedSecurities.year}.xlsx`);
        },
        formatDate: function (date: string, dateFormat: string) {
            return format(date, dateFormat)
        },
        loadMaturedSecurities: async function (showToast = true) {
            this.maturedSecurities.formalDate = makeQuarterFormalDate(this.maturedSecurities.quarter, this.maturedSecurities.year)
            const response = await axios.post("/load-matured-securities", {
                quarter: this.maturedSecurities.formalDate
            });
            this.maturedSecuritiesData = response.data.data.maturities;
            showToast && this.handleApiResponse(response.data, 'Matured Securities loaded')
        }
    }
});
