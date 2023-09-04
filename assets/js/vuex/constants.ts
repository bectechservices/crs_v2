export enum CUSTODIAN_CONSTANTS {
    ADD_ORDINARY_SHARE_INPUT = "addOrdinaryShareInput",
    DELETE_ORDINARY_SHARE_INPUT = "deleteOrdinaryShareInput",
    MODIFY_ORDINARY_SHARE_INPUT_DATA = "modifyOrdinaryShareInputData",
    ADD_PREFERENCE_SHARE_INPUT = "addPreferenceShareInput",
    DELETE_PREFERENCE_SHARE_INPUT = "deletePreferenceShareInput",
    MODIFY_PREFERENCE_SHARE_INPUT_DATA = "modifyPreferenceShareInputData",
    ADD_CUSTODIAN_TRANSACTION_INPUT = "addCustodianTransactionInput",
    DELETE_CUSTODIAN_TRANSACTION_INPUT = "deleteCustodianTransactionInput",
    MODIFY_CUSTODIAN_TRANSACTION_INPUT_DATA = "modifyCustodianTransactionInputData",
    ADD_SCHEME_UNDER_CUSTODY_INPUT = "addSchemeUnderCustodyInput",
    DELETE_SCHEME_UNDER_CUSTODY_INPUT = "deleteSchemeUnderCustodyInput",
    MODIFY_SCHEME_UNDER_CUSTODY_INPUT_DATA = "modifySchemeUnderCustodyInputData",
    VALIDATE_GOVERNANCE_MULTIPLE_FIELDS_DATA = "validateGovernanceMultipleFieldData",
    CLEAR_ALL_DATA = "clearAllData"
}

export enum CUSTODIAN_ACTION_CONSTANTS {
    UPLOAD_GOVERNANCE_DATA = "uploadGovernanceData",
    UPLOAD_SCHEME_DETAILS = "uploadSchemeDetails",
    UPLOAD_OTHER_INFORMATION = "uploadOtherInformation",
    UPLOAD_OFFICIAL_REPORT_REMARKS = "uploadOfficialReportRemarks",
    FETCH_SCHEME_DETAILS = "fetchSchemeDetails",
    RECALCULATE_FORMULAS = "recalculateFormulas",
    SEND_CHANGES_TO_CHECKER = "sendChangesToChecker",
    DELETE_SCHEME_DETAILS = "deleteSchemeDetails"
}

export enum ROOT_STATE_CONSTANTS {
    SET_LOADING_STATUS = "setLoadingStatus",
    STORE_PV_UPLOAD_ERRORS = "storePvUploadErrors"
}

export enum PPT_OPTIONS_CONSTANTS {
    STORE_PPT_OPTIONS = "storePPTOptions"
}

export enum TRUSTEE_CONSTANTS {
    ADD_UNIDENTIFIED_PAYMENT_INPUT = "addUnidentifiedPaymentInput",
    DELETE_UNIDENTIFIED_PAYMENT_INPUT = "deleteUnidentifiedPaymentInput",
    MODIFY_UNIDENTIFIED_PAYMENT_DATA = "modifyUnidentifiedPaymentData",
    CLEAR_ALL_UNIDENTIFIED_PAYMENTS = "clearAllUnidentifiedPayments",
    STORE_DASHBOARD_SEARCH_INPUT = "storeDashboardSearchInput"
}

export enum NPRA321_CONSTANTS {
    STORE_321_SEARCH_INPUT = "store321SearchInput"
}

export enum TRUSTEE_ACTION_CONSTANTS {
    UPLOAD_UNIDENTIFIED_PAYMENTS = "uploadUnidentifiedPayments"
}

export enum NPRA_CONSTANTS {
    ADD_UNAUTHORIZED_TRANSACTION_INPUT = "addUnauthorizedTransactionInput",
    DELETE_UNAUTHORIZED_TRANSACTION_INPUT = "deleteUnauthorizedTransactionInput",
    MODIFY_UNAUTHORIZED_TRANSACTION_DATA = "modifyUnauthorizedTransactionData",
    CLEAR_ALL_UNAUTHORIZED_TRANSACTION = "clearAllUnauthorizedTransactions"
}

export enum NPRA_ACTION_CONSTANTS {
    UPLOAD_UNAUTHORIZED_TRANSACTIONS = "uploadUnauthorizedTransactions",
    UPLOAD_OUTSTANDING_FD_CERTIFICATES = "uploadOutstandingFDCertificates",
    SEND_CHANGES_TO_CHECKERS = "npraSendChangesToChecker"
}

export enum CLIENT_CONSTANTS {
    ADD_RANGED_BASIS_POINT_INPUT = "addRangedBasisPointInput",
    DELETE_RANGED_BASIS_POINT_INPUT = "deleteRangedBasisPointInput",
    MODIFY_RANGED_BASIS_POINT_DATA = "modifyRangedBasisPointData",
    CLEAR_ALL_RANGED_BASIS_POINT_INPUTS = "clearAllRangedBasisPointInputs"
}

export enum CLIENT_ACTIONS_CONSTANTS {
    UPLOAD_CLIENT_BILLING_INFORMATION = "uploadClientBillingInfo"
}

export enum BILLING_CONSTANTS {
    STORE_BILLING_SEARCH_DATA = "storeBillingSearchData"
}