import Vue from "vue";
import {
    convertFileToBase64,
    CURRENT_YEAR,
    currentQuarter,
    currentQuarterNumber,
    dataURLToBlob,
    isNumeric,
    isPDFFile,
    lastDayOfQuarter,
    lastSecLicenseRenewalDate,
    parseStringToFloat
} from "./lib/helpers";
import {
    ApiResponse,
    ICanDuplicate,
    ICustodianTransaction,
    IFile,
    IGovernance,
    IGovernanceError,
    IKeyableCustodianTransaction,
    IKeyableSchemesUnderCustody,
    IKeyableShare,
    IOfficialRemarks,
    IOfficialRemarksError,
    IOtherSecInformations,
    IOtherSecInformationsError,
    ISchemeDetails,
    ISchemeDetailsError,
    ISchemesUnderCustody,
    IShare,
    SchemeDetailsRequest
} from "./lib/types";
import store from "./vuex";
import {CUSTODIAN_ACTION_CONSTANTS, CUSTODIAN_CONSTANTS} from "./vuex/constants";
import axios from "axios";
import {format, parse} from "date-fns";
import collect from "collect.js";
import VueSweetalert2 from 'vue-sweetalert2';
import ExcelFilePreview from "./lib/excel";


Vue.use(VueSweetalert2);

interface Data {
    files: Array<Array<IFile>>;
    filesUploaded: boolean;
    containsUnknownFileType: boolean;
    activeFile: IFile;
    lastDayOfTheQuarter: string;
    lastSecLicenseRenewalDate: string;
    currentQuarter: string;
    ordinarySharesComponents: ICanDuplicate[];
    preferenceSharesComponents: ICanDuplicate[];
    custodianTransactionComponents: ICanDuplicate[];
    schemesUnderCustodyComponents: ICanDuplicate[];
    governanceData: IGovernance;
    governanceDataError: IGovernanceError;
    schemeDetails: ISchemeDetails;
    otherSecInformations: IOtherSecInformations;
    officialRemarks: IOfficialRemarks;
    canEditData: boolean;
    schemeDetailsError: ISchemeDetailsError;
    otherSecInformationsError: IOtherSecInformationsError;
    officialRemarksError: IOfficialRemarksError;
    schemeSearch: string;
    schemeSearchError: boolean;
    uploadedSchemeDetailsTemplate: Array<Array<string>>
}

interface Methods {
    onFileChange: (event: any) => void;
    onSchemeTemplateUpload: (event: any) => Promise<void>;
    loadPDFFileByName: (name: string) => void;
    addOrdinaryShareInput: () => void;
    deleteOrdinaryShareInput: (item: ICanDuplicate) => void;
    addPreferenceShareInput: () => void;
    deletePreferenceShareInput: (item: ICanDuplicate) => void;
    addCustodianTransactionInput: () => void;
    deleteCustodianTransactionInput: (item: ICanDuplicate) => void;
    addSchemesUnderCustodyInput: () => void;
    deleteSchemesUnderCustodyInput: (item: ICanDuplicate) => void;
    submitGovernanceData: () => void;
    validateGovernanceData: () => boolean;
    validateSchemeDetailsData: () => boolean;
    saveSchemeDetails: () => void;
    emptySchemeDetailsError: () => ISchemeDetailsError;
    validateOtherSecInfoData: () => boolean;
    saveOtherSecInfo: () => void;
    validateOfficialRemarks: () => boolean;
    saveOfficialRemarks: () => void;
    addMoreSchemes: () => void;
    loadSchemeByBPID: () => void;
    recalculateFormulas: () => void;
    previewReport: () => void;
    submitChanges: () => Promise<void>;
    handleApiResponse: (response: ApiResponse, successMessage: string) => void;
    deleteScheme: () => Promise<void>;
}

interface ISharesComponentData {
    name: string;
    shareholding: string | number;
    percentage: string | number;
}

interface ISharesComponentMethods {
    commitChanges: () => void;
    onShareNameChange: ({target}: { target: HTMLInputElement }) => void;
    onShareHoldingChange: ({target}: { target: HTMLInputElement }) => void;
    onSharePercentageChange: ({target}: { target: HTMLInputElement }) => void;
}

Vue.component<ISharesComponentData, ISharesComponentMethods, {}, {}>(
    "ordinary-shares",
    {
        template: `
  <tr>
    <td><input type="text" class="assetFormInput h-100 w-100 uk-input" autocomplete="new-field" @input="onShareNameChange" v-bind:style="{border: errors ? (errors.name ? '1px solid red' : ''):''}" :value="name"/></td>
    <td><input type="number" class="assetFormInput h-100 w-100 uk-input" autocomplete="new-field" @input="onShareHoldingChange" v-bind:style="{border: errors ? (errors.shareholding ? '1px solid red' : ''):''}" :value="shareholding"/></td>
    <td><input type="number" min="0" max="100" step="0.01" class="assetFormInput h-100 w-100 uk-input" autocomplete="new-field" @input="onSharePercentageChange" v-bind:style="{border: errors ? (errors.percentage ? '1px solid red' : ''):''}" :value="percentage"/></td>
    <td>
      <button class="w-100" v-if="canAdd" @click="$emit('add-ordinary-shares-clicked')" type="button">Add</button>
      <button class="w-100" v-if="canDelete" @click="$emit('delete-ordinary-shares-clicked')" type="button">Delete</button>
    </td>
  </tr>
  `,
        props: ["canAdd", "canDelete", "share"],
        data: function () {
            const share: IShare = this.$store.getters.ordinaryShareInputData(
                this.$props.share.key
            );
            return {
                name: share ? share.name : "",
                shareholding: share ? share.shareholding : "",
                percentage: share ? share.percentage : ""
            };
        },
        methods: {
            commitChanges: function () {
                this.$nextTick(() => {
                    this.$store.commit(
                        CUSTODIAN_CONSTANTS.MODIFY_ORDINARY_SHARE_INPUT_DATA,
                        {
                            key: this.$props.share.key,
                            name: this.name,
                            shareholding: this.shareholding,
                            percentage: this.percentage
                        }
                    );
                });
            },
            onShareNameChange: function ({target}: { target: HTMLInputElement }) {
                this.name = target.value;
                this.commitChanges();
            },
            onShareHoldingChange: function ({target}: { target: HTMLInputElement }) {
                this.shareholding = parseFloat(target.value);
                this.commitChanges();
            },
            onSharePercentageChange: function ({
                                                   target
                                               }: {
                target: HTMLInputElement;
            }) {
                this.percentage = parseFloat(target.value);
                this.commitChanges();
            }
        },
        computed: {
            errors: function () {
                return this.$store.getters.ordinaryShareInputHasError(
                    this.$props.share.key
                );
            }
        }
    }
);

Vue.component<ISharesComponentData, ISharesComponentMethods, {}, {}>(
    "preference-shares",
    {
        template: `
  <tr>
    <td><input type="text" class="assetFormInput h-100 w-100 uk-input" autocomplete="new-field" @input="onShareNameChange" v-bind:style="{border: errors ? (errors.name ? '1px solid red' : ''):''}" :value="name"/></td>
    <td><input type="number" step="0.01" class="assetFormInput h-100 w-100 uk-input" autocomplete="new-field" @input="onShareHoldingChange" v-bind:style="{border: errors ? (errors.shareholding ? '1px solid red' : ''):''}" :value="shareholding"/></td>
    <td><input type="number" min="0" max="100" step="0.01" class="assetFormInput h-100 w-100 uk-input" autocomplete="new-field" @input="onSharePercentageChange" v-bind:style="{border: errors ? (errors.percentage ? '1px solid red' : ''):''}" :value="percentage"/></td>
    <td>
      <button style="width:100%" v-if="canAdd" @click="$emit('add-preference-shares-clicked')" type="button">Add</button>
      <button style="width:100%" v-if="canDelete" @click="$emit('delete-preference-shares-clicked')" type="button">Delete</button>
    </td>
  </tr>
  `,
        props: ["canAdd", "canDelete", "share"],
        data: function () {
            const share: IShare = this.$store.getters.preferenceShareInputData(
                this.$props.share.key
            );
            return {
                name: share ? share.name : "",
                shareholding: share ? share.shareholding : "",
                percentage: share ? share.percentage : ""
            };
        },
        methods: {
            commitChanges: function () {
                this.$nextTick(() => {
                    this.$store.commit(
                        CUSTODIAN_CONSTANTS.MODIFY_PREFERENCE_SHARE_INPUT_DATA,
                        {
                            key: this.$props.share.key,
                            name: this.name,
                            shareholding: this.shareholding,
                            percentage: this.percentage
                        }
                    );
                });
            },
            onShareNameChange: function ({target}: { target: HTMLInputElement }) {
                this.name = target.value;
                this.commitChanges();
            },
            onShareHoldingChange: function ({target}: { target: HTMLInputElement }) {
                this.shareholding = parseFloat(target.value);
                this.commitChanges();
            },
            onSharePercentageChange: function ({
                                                   target
                                               }: {
                target: HTMLInputElement;
            }) {
                this.percentage = parseFloat(target.value);
                this.commitChanges();
            }
        },
        computed: {
            errors: function () {
                return this.$store.getters.preferenceShareInputHasError(
                    this.$props.share.key
                );
            }
        }
    }
);

interface ICustodianTransactionComponentData {
    nameOfTrustee: string;
    relationshipWithTrustee: string;
    typeOfTransaction: string;
    amount: number;
}

interface ICustodianTransactionComponentMethods {
    commitChanges: () => void;
    onNameOfTrusteeChange: ({target}: { target: HTMLInputElement }) => void;
    onRelationsipWithTrusteeChange: ({
                                         target
                                     }: {
        target: HTMLInputElement;
    }) => void;
    onTypeOfTransactionChange: ({target}: { target: HTMLInputElement }) => void;
    onAmountChange: ({target}: { target: HTMLInputElement }) => void;
}

Vue.component<ICustodianTransactionComponentData,
    ICustodianTransactionComponentMethods,
    {},
    {}>("custodian-transactions", {
    template: `
  <tr>
    <td><input type="text" class="assetFormInput h-100 w-100 uk-input" autocomplete="new-field" @input="onNameOfTrusteeChange" v-bind:style="{border: errors ? (errors.nameOfTrustee ? '1px solid red' : ''):''}" :value="nameOfTrustee"/></td>
    <td><input type="text" class="assetFormInput h-100 w-100 uk-input" autocomplete="new-field" @input="onRelationsipWithTrusteeChange" v-bind:style="{border: errors ? (errors.relationshipWithTrustee ? '1px solid red' : ''):''}" :value="relationshipWithTrustee"/></td>
    <td><input type="text" class="assetFormInput h-100 w-100 uk-input" autocomplete="new-field" @input="onTypeOfTransactionChange"  v-bind:style="{border: errors ? (errors.typeOfTransaction ? '1px solid red' : ''):''}" :value="typeOfTransaction"/></td>
    <td><input type="number" class="assetFormInput h-100 w-100 uk-input" autocomplete="new-field" @input="onAmountChange" v-bind:style="{border: errors ? (errors.amount ? '1px solid red' : ''):''}" :value="amount" step="0.01"/></td>
    <td>
      <button style="width:100%" v-if="canAdd" @click="$emit('add-custodian-transaction-clicked')" type="button">Add</button>
      <button style="width:100%" v-if="canDelete" @click="$emit('delete-custodian-transaction-clicked')" type="button">Delete</button>
    </td>
  </tr>
  `,
    props: ["canAdd", "canDelete", "transaction"],
    data: function () {
        const transaction: ICustodianTransaction = this.$store.getters.custodianTransactionsInputData(
            this.$props.transaction.key
        );
        return {
            nameOfTrustee: transaction ? transaction.nameOfTrustee : "",
            relationshipWithTrustee: transaction
                ? transaction.relationshipWithTrustee
                : "",
            typeOfTransaction: transaction ? transaction.typeOfTransaction : "",
            amount: transaction ? transaction.amount : 0
        };
    },
    methods: {
        commitChanges: function () {
            this.$nextTick(() => {
                this.$store.commit(
                    CUSTODIAN_CONSTANTS.MODIFY_CUSTODIAN_TRANSACTION_INPUT_DATA,
                    {
                        key: this.$props.transaction.key,
                        nameOfTrustee: this.nameOfTrustee,
                        relationshipWithTrustee: this.relationshipWithTrustee,
                        typeOfTransaction: this.typeOfTransaction,
                        amount: this.amount
                    }
                );
            });
        },
        onNameOfTrusteeChange: function ({target}: { target: HTMLInputElement }) {
            this.nameOfTrustee = target.value;
            this.commitChanges();
        },
        onRelationsipWithTrusteeChange: function ({
                                                      target
                                                  }: {
            target: HTMLInputElement;
        }) {
            this.relationshipWithTrustee = target.value;
            this.commitChanges();
        },
        onTypeOfTransactionChange: function ({
                                                 target
                                             }: {
            target: HTMLInputElement;
        }) {
            this.typeOfTransaction = target.value;
            this.commitChanges();
        },
        onAmountChange: function ({target}: { target: HTMLInputElement }) {
            this.amount = parseFloat(target.value);
            this.commitChanges();
        }
    },
    computed: {
        errors: function () {
            return this.$store.getters.custodianTransactionsInputHasError(
                this.$props.transaction.key
            );
        }
    }
});

interface ISchemesUnderCustodyComponentData {
    nameOfFirm: string;
    nameOfScheme: string;
    relationshipWithTrustee: string;
    volume: number;
    markedToMarketValue: number;
}

interface ISchemesUnderCustodyComponentMethods {
    commitChanges: () => void;
    onNameOfFirmChange: ({target}: { target: HTMLInputElement }) => void;
    onNameOfSchemeChange: ({target}: { target: HTMLInputElement }) => void;
    onRelationshipWithTrusteeChange: ({
                                          target
                                      }: {
        target: HTMLInputElement;
    }) => void;
    onVolumeChange: ({target}: { target: HTMLInputElement }) => void;
    onMarkedToMarketValueChange: ({
                                      target
                                  }: {
        target: HTMLInputElement;
    }) => void;
}

Vue.component<ISchemesUnderCustodyComponentData,
    ISchemesUnderCustodyComponentMethods,
    {},
    {}>("schemes-under-custody", {
    template: `
  <tr>
    <td><input type="text" class="assetFormInput h-100 w-100 uk-input" autocomplete="new-field" @input="onNameOfFirmChange" v-bind:style="{border: errors ? (errors.nameOfFirm ? '1px solid red' : ''):''}" :value="nameOfFirm"/></td>
    <td><input type="text" class="assetFormInput h-100 w-100 uk-input" autocomplete="new-field" @input="onNameOfSchemeChange" v-bind:style="{border: errors ? (errors.nameOfScheme ? '1px solid red' : ''):''}" :value="nameOfScheme"/></td>
    <td><input type="text" class="assetFormInput h-100 w-100 uk-input" autocomplete="new-field" @input="onRelationshipWithTrusteeChange" v-bind:style="{border: errors ? (errors.relationshipWithTrustee ? '1px solid red' : ''):''}" :value="relationshipWithTrustee"/></td>
    <td><input type="number" class="assetFormInput h-100 w-100 uk-input" autocomplete="new-field" @input="onVolumeChange" v-bind:style="{border: errors ? (errors.volume ? '1px solid red' : ''):''}" :value="volume" step="0.01"/></td>
    <td><input type="number" class="assetFormInput h-100 w-100 uk-input" autocomplete="new-field" @input="onMarkedToMarketValueChange" v-bind:style="{border: errors ? (errors.markedToMarketValue ? '1px solid red' : ''):''}" :value="markedToMarketValue" step="0.01"/></td>
    <td>
      <button style="width:100%" v-if="canAdd" @click="$emit('add-scheme-under-custody-clicked')" type="button">Add</button>
      <button style="width:100%" v-if="canDelete" @click="$emit('delete-scheme-under-custody-clicked')" type="button">Delete</button>
    </td>
  </tr>
  `,
    props: ["canAdd", "canDelete", "scheme"],
    data: function () {
        const scheme: ISchemesUnderCustody = this.$store.getters.schemesUnderCustodyInputData(
            this.$props.scheme.key
        );
        return {
            nameOfFirm: scheme ? scheme.nameOfFirm : "",
            nameOfScheme: scheme ? scheme.nameOfScheme : "",
            relationshipWithTrustee: scheme ? scheme.relationshipWithTrustee : "",
            volume: scheme ? scheme.volume : 0,
            markedToMarketValue: scheme ? scheme.markedToMarketValue : 0
        };
    },
    methods: {
        commitChanges: function () {
            this.$nextTick(() => {
                this.$store.commit(
                    CUSTODIAN_CONSTANTS.MODIFY_SCHEME_UNDER_CUSTODY_INPUT_DATA,
                    {
                        key: this.$props.scheme.key,
                        nameOfFirm: this.nameOfFirm,
                        nameOfScheme: this.nameOfScheme,
                        relationshipWithTrustee: this.relationshipWithTrustee,
                        volume: this.volume,
                        markedToMarketValue: this.markedToMarketValue
                    }
                );
            });
        },
        onNameOfFirmChange: function ({target}: { target: HTMLInputElement }) {
            this.nameOfFirm = target.value;
            this.commitChanges();
        },
        onNameOfSchemeChange: function ({target}: { target: HTMLInputElement }) {
            this.nameOfScheme = target.value;
            this.commitChanges();
        },
        onRelationshipWithTrusteeChange: function ({
                                                       target
                                                   }: {
            target: HTMLInputElement;
        }) {
            this.relationshipWithTrustee = target.value;
            this.commitChanges();
        },
        onVolumeChange: function ({target}: { target: HTMLInputElement }) {
            this.volume = parseFloat(target.value);
            this.commitChanges();
        },
        onMarkedToMarketValueChange: function ({
                                                   target
                                               }: {
            target: HTMLInputElement;
        }) {
            this.markedToMarketValue = parseFloat(target.value);
            this.commitChanges();
        }
    },
    computed: {
        errors: function () {
            return this.$store.getters.schemesUnderCustodyInputHasError(
                this.$props.scheme.key
            );
        }
    }
});

export default new Vue<Data, Methods>({
    el: ".custodianPage",
    store,
    data: {
        files: [],
        filesUploaded: false,
        containsUnknownFileType: false,
        activeFile: {
            name: "",
            data: "",
            rawBase64Data: ""
        },
        lastDayOfTheQuarter: lastDayOfQuarter(),
        lastSecLicenseRenewalDate: lastSecLicenseRenewalDate
        (),
        currentQuarter: currentQuarter(),
        ordinarySharesComponents: [
            {
                key: Math.floor(Math.random() * 10000)
            }
        ],
        preferenceSharesComponents: [
            {
                key: Math.floor(Math.random() * 10000)
            }
        ],
        custodianTransactionComponents: [
            {
                key: Math.floor(Math.random() * 10000)
            }
        ],
        schemesUnderCustodyComponents: [
            {
                key: Math.floor(Math.random() * 10000)
            }
        ],
        governanceData: {
            clientName: "Standard Chartered Bank",
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
        governanceDataError: {
            clientName: false,
            reportingOfficer: false,
            reportingDate: false,
            custodianAssetRegistrationDate: false
        },
        schemeDetails: {
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
        otherSecInformations: {
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
        officialRemarks: {
            remarks: "",
            reviewingOfficer: "",
            date: "",
            signature: ""
        },
        canEditData: true,// Hack due to changes
        schemeDetailsError: {
            bpid: false,
            nameOfScheme: false,
            numberOfSharesOutstanding: false,
            numberOfShareholders: false,
            numberOfRedemptions: false,
            valueOfRedemptions: false,
            nameOfManager: false,
            totalValueOfSchemeAssets: false,
            netAssetValueOfScheme: false,
            netAssetValuePerShare: false,
            totalEquityInvestments: false,
            totalFixedIncomeInvestments: false,
            netOfMediumTermAssetsHeldByFund: false,
            capitalMarketsInvestments: false,
            percentageOfCapitalInvestmentToTotalInvestment: false,
            totalValueOfUnutilizedFunds: false,
            valueOfBorrowedFunds: false,
            reasonsForBorrowing: false,
            redemptions: false,
            dividends: false,
            rights: false,
            feesOwedCustodian: false
        },
        otherSecInformationsError: {
            litigationDetails: false
        },
        officialRemarksError: {
            remarks: false,
            reviewingOfficer: false,
            date: false
        },
        schemeSearch: "",
        schemeSearchError: false,
        uploadedSchemeDetailsTemplate: []
    },
    beforeMount: async function () {
        store.commit(CUSTODIAN_CONSTANTS.CLEAR_ALL_DATA); //reset state

        let response = await axios.post("/sec-quarterly-report", {
            year: CURRENT_YEAR,
            quarter: currentQuarterNumber()
        });
        const report = response.data.data.report;
        this.schemeDetails =
            response.data.data.schemeDetails.length &&
            response.data.data.schemeDetails[0];
        if (response.data.data.schemeDetails.length) {
            if (response.data.data.schemeDetails[0].attachedFile) {
                this.activeFile = {
                    name: "Attached File.pdf",
                    data: URL.createObjectURL(
                        await dataURLToBlob(`data:application/pdf;base64,${response.data.data.schemeDetails[0].attachedFile}`)
                    ),
                    rawBase64Data: response.data.data.schemeDetails[0].attachedFile
                };
                this.filesUploaded = true;
            }
        }
        this.otherSecInformations = response.data.data.otherInformation;
        this.officialRemarks = response.data.data.remarks;
        if (report && report.id !== 0) {
            if (report.ordinaryShares.length) {
                this.ordinarySharesComponents = report.ordinaryShares.map(
                    (share: any) => {
                        const key = Math.floor(Math.random() * 10000);
                        store.commit(CUSTODIAN_CONSTANTS.ADD_ORDINARY_SHARE_INPUT, key);
                        store.commit(CUSTODIAN_CONSTANTS.MODIFY_ORDINARY_SHARE_INPUT_DATA, {
                            key,
                            name: share.name,
                            shareholding: share.shareholdings,
                            percentage: share.percentage
                        });
                        return {
                            key
                        };
                    }
                ) as Array<IKeyableShare>;
            } else {
                store.commit(
                    CUSTODIAN_CONSTANTS.ADD_ORDINARY_SHARE_INPUT,
                    this.ordinarySharesComponents[0].key
                );
            }

            if (report.preferenceShares.length) {
                this.preferenceSharesComponents = report.preferenceShares.map(
                    (share: any) => {
                        const key = Math.floor(Math.random() * 10000);
                        store.commit(CUSTODIAN_CONSTANTS.ADD_PREFERENCE_SHARE_INPUT, key);
                        store.commit(
                            CUSTODIAN_CONSTANTS.MODIFY_PREFERENCE_SHARE_INPUT_DATA,
                            {
                                key,
                                name: share.name,
                                shareholding: share.shareholdings,
                                percentage: share.percentage
                            }
                        );
                        return {
                            key
                        };
                    }
                ) as Array<IKeyableShare>;
            } else {
                store.commit(
                    CUSTODIAN_CONSTANTS.ADD_PREFERENCE_SHARE_INPUT,
                    this.preferenceSharesComponents[0].key
                );
            }

            if (report.affiliateTransactions.length) {
                this.custodianTransactionComponents = report.affiliateTransactions.map(
                    (transaction: any) => {
                        const key = Math.floor(Math.random() * 10000);
                        store.commit(
                            CUSTODIAN_CONSTANTS.ADD_CUSTODIAN_TRANSACTION_INPUT,
                            key
                        );
                        store.commit(
                            CUSTODIAN_CONSTANTS.MODIFY_CUSTODIAN_TRANSACTION_INPUT_DATA,
                            {
                                key,
                                nameOfTrustee: transaction.nameOfAffiliate,
                                relationshipWithTrustee: transaction.relationshipWithCustodian,
                                typeOfTransaction: transaction.typeOfTransaction,
                                amount: transaction.amount
                            } as IKeyableCustodianTransaction
                        );
                        return {
                            key
                        };
                    }
                ) as Array<IKeyableCustodianTransaction>;
            } else {
                store.commit(
                    CUSTODIAN_CONSTANTS.ADD_CUSTODIAN_TRANSACTION_INPUT,
                    this.custodianTransactionComponents[0].key
                );
            }

            if (report.valueVolumeOfShares.length) {
                this.schemesUnderCustodyComponents = report.valueVolumeOfShares.map(
                    (value: any) => {
                        const key = Math.floor(Math.random() * 10000);
                        store.commit(
                            CUSTODIAN_CONSTANTS.ADD_SCHEME_UNDER_CUSTODY_INPUT,
                            key
                        );
                        store.commit(
                            CUSTODIAN_CONSTANTS.MODIFY_SCHEME_UNDER_CUSTODY_INPUT_DATA,
                            {
                                key,
                                nameOfFirm: value.nameOfFirm,
                                nameOfScheme: value.nameOfScheme,
                                relationshipWithTrustee: value.relationshipWithCustodian,
                                volume: value.volume,
                                markedToMarketValue: value.markedToMarketValue
                            } as IKeyableSchemesUnderCustody
                        );
                        return {
                            key
                        };
                    }
                ) as Array<IKeyableSchemesUnderCustody>;
            } else {
                store.commit(
                    CUSTODIAN_CONSTANTS.ADD_SCHEME_UNDER_CUSTODY_INPUT,
                    this.schemesUnderCustodyComponents[0].key
                );
            }

            const governanceData = {
                clientName: report.custodianName,
                reportingOfficer: report.reportingOfficer,
                reportingDate: format(parse(report.dateOfReport), "YYYY-MM-DD"),
                ordinaryShares: [],
                preferenceShares: [],
                schemesUnderCustody: [],
                custodianTransactions: [],
                changeInDirectors: report.changeInDirectors,
                changeInAgreement: report.changeInAgreement,
                dealingsApprovedByBoard: report.dealingsApprovedByBoard,
                custodianHasUpdatedAssetRegister:
                report.custodianHasUpdatedAssetRegister,
                custodianAssetRegistrationDate: format(
                    parse(report.custodianAssetRegistrationDate),
                    "YYYY-MM-DD"
                ),
                doManagersOfTheSchemeConsultTheLaw:
                report.doManagersOfSchemeConsultTheLaw,
                schemeHadAnyOtherFinancialDealings:
                report.schemeHadAnyOtherFinancialDealings,
                approved: report.approved
            };
            this.governanceData = governanceData;
        } else {
            this.governanceData.approved = true;
            store.commit(
                CUSTODIAN_CONSTANTS.ADD_ORDINARY_SHARE_INPUT,
                this.ordinarySharesComponents[0].key
            );
            store.commit(
                CUSTODIAN_CONSTANTS.ADD_PREFERENCE_SHARE_INPUT,
                this.preferenceSharesComponents[0].key
            );
            store.commit(
                CUSTODIAN_CONSTANTS.ADD_CUSTODIAN_TRANSACTION_INPUT,
                this.custodianTransactionComponents[0].key
            );
            store.commit(
                CUSTODIAN_CONSTANTS.ADD_SCHEME_UNDER_CUSTODY_INPUT,
                this.schemesUnderCustodyComponents[0].key
            );
        }
    },
    methods: {
        onFileChange: function (event: any) {
            const uploadedFiles = event.target.files;
            const length = (uploadedFiles as FileList).length;
            const files: Array<IFile> = [];
            for (let i = 0; i < length; i++) {
                if (!isPDFFile((uploadedFiles as FileList)[i])) {
                    this.containsUnknownFileType = true;
                }
            }
            this.$nextTick(async () => {
                if (!this.containsUnknownFileType && uploadedFiles.length) {
                    this.filesUploaded = true;
                    for (let i = 0; i < length; i++) {
                        const file = (uploadedFiles as FileList)[i];
                        const data = await convertFileToBase64(file);
                        files.push({
                            name: file.name,
                            data: URL.createObjectURL(
                                await dataURLToBlob(`data:application/pdf;base64,${data}`)
                            ),
                            rawBase64Data: data
                        });
                        this.files = collect(files)
                            .chunk(3)
                            .toArray();
                    }
                    this.activeFile = files[0];
                }
            });
        },
        loadPDFFileByName: function (name: string) {
            this.activeFile = collect(this.files)
                .collapse()
                .toArray()
                .find((file: any) => file.name === name) as IFile;
        },
        addOrdinaryShareInput: function () {
            const key = Math.floor(Math.random() * 10000);
            this.ordinarySharesComponents.push({
                key
            });
            store.commit(CUSTODIAN_CONSTANTS.ADD_ORDINARY_SHARE_INPUT, key);
        },
        deleteOrdinaryShareInput: function (item: ICanDuplicate) {
            this.ordinarySharesComponents = this.ordinarySharesComponents.filter(
                (component: ICanDuplicate) => component.key !== item.key
            );
            store.commit(CUSTODIAN_CONSTANTS.DELETE_ORDINARY_SHARE_INPUT, item.key);
        },
        addPreferenceShareInput: function () {
            const key = Math.floor(Math.random() * 10000);
            this.preferenceSharesComponents.push({
                key
            });
            store.commit(CUSTODIAN_CONSTANTS.ADD_PREFERENCE_SHARE_INPUT, key);
        },
        deletePreferenceShareInput: function (item: ICanDuplicate) {
            this.preferenceSharesComponents = this.preferenceSharesComponents.filter(
                (component: ICanDuplicate) => component.key !== item.key
            );
            store.commit(CUSTODIAN_CONSTANTS.DELETE_PREFERENCE_SHARE_INPUT, item.key);
        },
        addCustodianTransactionInput: function () {
            const key = Math.floor(Math.random() * 10000);
            this.custodianTransactionComponents.push({
                key
            });
            store.commit(CUSTODIAN_CONSTANTS.ADD_CUSTODIAN_TRANSACTION_INPUT, key);
        },
        deleteCustodianTransactionInput: function (item: ICanDuplicate) {
            this.custodianTransactionComponents = this.custodianTransactionComponents.filter(
                (component: ICanDuplicate) => component.key !== item.key
            );
            store.commit(
                CUSTODIAN_CONSTANTS.DELETE_CUSTODIAN_TRANSACTION_INPUT,
                item.key
            );
        },
        addSchemesUnderCustodyInput: function () {
            const key = Math.floor(Math.random() * 10000);
            this.schemesUnderCustodyComponents.push({
                key
            });
            store.commit(CUSTODIAN_CONSTANTS.ADD_SCHEME_UNDER_CUSTODY_INPUT, key);
        },
        deleteSchemesUnderCustodyInput: function (item: ICanDuplicate) {
            this.schemesUnderCustodyComponents = this.schemesUnderCustodyComponents.filter(
                (component: ICanDuplicate) => component.key !== item.key
            );
            store.commit(
                CUSTODIAN_CONSTANTS.DELETE_SCHEME_UNDER_CUSTODY_INPUT,
                item.key
            );
        },
        submitGovernanceData: function () {
            store.commit(
                CUSTODIAN_CONSTANTS.VALIDATE_GOVERNANCE_MULTIPLE_FIELDS_DATA
            );
            this.$nextTick(async () => {
                const multipleFieldsHasError =
                    store.getters.custodianMultipleInputFieldHasError;
                let governanceDataIsValid = this.validateGovernanceData();
                if (!multipleFieldsHasError && !governanceDataIsValid) {
                    this.$swal({
                        type: 'error',
                        title: "All inputs are required.",
                        animation: false
                    }).then(console.log);
                } else if (!multipleFieldsHasError && governanceDataIsValid) {
                    const governanceData: IGovernance = {...this.governanceData};
                    const state = store.getters.multipleFieldsData;
                    governanceData.ordinaryShares = state.ordinarySharesInputData.map(
                        (share: IKeyableShare) => {
                            const {key: _, ...data} = share;
                            return data;
                        }
                    );
                    governanceData.preferenceShares = state.preferenceSharesInputData.map(
                        (share: IKeyableShare) => {
                            const {key: _, ...data} = share;
                            return data;
                        }
                    );
                    governanceData.custodianTransactions = state.custodianTransactionsInputData.map(
                        (transaction: IKeyableCustodianTransaction) => {
                            const {key: _, ...data} = transaction;
                            return data;
                        }
                    );
                    governanceData.schemesUnderCustody = state.schemesUnderCustodyInputData.map(
                        (scheme: IKeyableSchemesUnderCustody) => {
                            const {key: _, ...data} = scheme;
                            return data;
                        }
                    );
                    if (this.canEditData) {
                        const response = await store.dispatch(
                            CUSTODIAN_ACTION_CONSTANTS.UPLOAD_GOVERNANCE_DATA,
                            governanceData
                        );
                        this.handleApiResponse(response, "Governance info uploaded")
                    }
                }
            });
        },
        validateGovernanceData: function (): boolean {
            this.governanceDataError.clientName = false;
            this.governanceDataError.reportingOfficer = false;
            this.governanceDataError.reportingDate = false;
            this.governanceDataError.custodianAssetRegistrationDate = false;
            let governanceDataIsValid = true;
            if (!this.governanceData.clientName) {
                this.governanceDataError.clientName = true;
                governanceDataIsValid = false;
            }
            if (!this.governanceData.reportingOfficer) {
                this.governanceDataError.reportingOfficer = true;
                governanceDataIsValid = false;
            }
            if (!this.governanceData.reportingDate) {
                this.governanceDataError.reportingDate = true;
                governanceDataIsValid = false;
            }
            if (
                !this.governanceData.custodianAssetRegistrationDate &&
                this.governanceData.custodianHasUpdatedAssetRegister == "yes"
            ) {
                this.governanceDataError.custodianAssetRegistrationDate = true;
                governanceDataIsValid = false;
            }
            return (
                governanceDataIsValid &&
                (Boolean(this.governanceData.changeInDirectors) &&
                    Boolean(this.governanceData.changeInAgreement) &&
                    Boolean(this.governanceData.dealingsApprovedByBoard) &&
                    Boolean(this.governanceData.custodianHasUpdatedAssetRegister) &&
                    Boolean(this.governanceData.doManagersOfTheSchemeConsultTheLaw) &&
                    Boolean(this.governanceData.schemeHadAnyOtherFinancialDealings))
            );
        },
        validateSchemeDetailsData: function (): boolean {
            this.schemeDetailsError = this.emptySchemeDetailsError();
            let schemeDetailsIsValid = true;
            let schemeDetails: ISchemeDetails = this.schemeDetails;
            if (!schemeDetails.bpid) {
                this.schemeDetailsError.bpid = true;
                schemeDetailsIsValid = false;
            }
            if (!schemeDetails.nameOfScheme) {
                this.schemeDetailsError.nameOfScheme = true;
                schemeDetailsIsValid = false;
            }
            if (!isNumeric(schemeDetails.numberOfSharesOutstanding)) {
                this.schemeDetailsError.numberOfSharesOutstanding = true;
                schemeDetailsIsValid = false;
            }
            if (!isNumeric(schemeDetails.numberOfShareholders)) {
                this.schemeDetailsError.numberOfShareholders = true;
                schemeDetailsIsValid = false;
            }
            if (!isNumeric(schemeDetails.numberOfRedemptions)) {
                this.schemeDetailsError.numberOfRedemptions = true;
                schemeDetailsIsValid = false;
            }
            if (!isNumeric(schemeDetails.valueOfRedemptions)) {
                this.schemeDetailsError.valueOfRedemptions = true;
                schemeDetailsIsValid = false;
            }
            if (!schemeDetails.nameOfManager) {
                this.schemeDetailsError.nameOfManager = true;
                schemeDetailsIsValid = false;
            }
            if (!isNumeric(schemeDetails.totalValueOfSchemeAssets)) {
                this.schemeDetailsError.totalValueOfSchemeAssets = true;
                schemeDetailsIsValid = false;
            }
            if (!isNumeric(schemeDetails.netAssetValueOfScheme)) {
                this.schemeDetailsError.netAssetValueOfScheme = true;
                schemeDetailsIsValid = false;
            }
            if (!isNumeric(schemeDetails.netAssetValuePerShare)) {
                this.schemeDetailsError.netAssetValuePerShare = true;
                schemeDetailsIsValid = false;
            }
            if (!isNumeric(schemeDetails.totalEquityInvestments)) {
                this.schemeDetailsError.totalEquityInvestments = true;
                schemeDetailsIsValid = false;
            }
            if (!isNumeric(schemeDetails.totalFixedIncomeInvestments)) {
                this.schemeDetailsError.totalFixedIncomeInvestments = true;
                schemeDetailsIsValid = false;
            }
            if (!isNumeric(schemeDetails.netOfMediumTermAssetsHeldByFund)) {
                this.schemeDetailsError.netOfMediumTermAssetsHeldByFund = true;
                schemeDetailsIsValid = false;
            }
            if (!isNumeric(schemeDetails.capitalMarketsInvestments)) {
                this.schemeDetailsError.capitalMarketsInvestments = true;
                schemeDetailsIsValid = false;
            }
            if (
                !isNumeric(schemeDetails.percentageOfCapitalInvestmentToTotalInvestment)
            ) {
                this.schemeDetailsError.percentageOfCapitalInvestmentToTotalInvestment = true;
                schemeDetailsIsValid = false;
            }
            if (!isNumeric(schemeDetails.totalValueOfUnutilizedFunds)) {
                this.schemeDetailsError.totalValueOfUnutilizedFunds = true;
                schemeDetailsIsValid = false;
            }
            if (!isNumeric(schemeDetails.valueOfBorrowedFunds)) {
                this.schemeDetailsError.valueOfBorrowedFunds = true;
                schemeDetailsIsValid = false;
            } else {
                if (
                    schemeDetails.valueOfBorrowedFunds > 0 &&
                    !schemeDetails.reasonsForBorrowing
                ) {
                    this.schemeDetailsError.reasonsForBorrowing = true;
                    schemeDetailsIsValid = false;
                }
            }
            if (!isNumeric(schemeDetails.redemptions)) {
                this.schemeDetailsError.redemptions = true;
                schemeDetailsIsValid = false;
            }
            if (!isNumeric(schemeDetails.dividends)) {
                this.schemeDetailsError.dividends = true;
                schemeDetailsIsValid = false;
            }
            if (!isNumeric(schemeDetails.rights)) {
                this.schemeDetailsError.rights = true;
                schemeDetailsIsValid = false;
            }
            if (!isNumeric(schemeDetails.feesOwedCustodian)) {
                this.schemeDetailsError.feesOwedCustodian = true;
                schemeDetailsIsValid = false;
            }
            return schemeDetailsIsValid;
        },
        emptySchemeDetailsError: function (): ISchemeDetailsError {
            return {
                bpid: false,
                nameOfScheme: false,
                numberOfSharesOutstanding: false,
                numberOfShareholders: false,
                numberOfRedemptions: false,
                valueOfRedemptions: false,
                nameOfManager: false,
                totalValueOfSchemeAssets: false,
                netAssetValueOfScheme: false,
                netAssetValuePerShare: false,
                totalEquityInvestments: false,
                totalFixedIncomeInvestments: false,
                netOfMediumTermAssetsHeldByFund: false,
                capitalMarketsInvestments: false,
                percentageOfCapitalInvestmentToTotalInvestment: false,
                totalValueOfUnutilizedFunds: false,
                valueOfBorrowedFunds: false,
                reasonsForBorrowing: false,
                redemptions: false,
                dividends: false,
                rights: false,
                feesOwedCustodian: false
            };
        },
        saveSchemeDetails: async function () {
            if (this.validateSchemeDetailsData()) {
                if (this.canEditData) {
                    const details: ISchemeDetails = this.schemeDetails;
                    details.numberOfSharesOutstanding = parseFloat(
                        `${details.numberOfSharesOutstanding}`
                    );
                    details.numberOfShareholders = parseFloat(
                        `${details.numberOfShareholders}`
                    );
                    details.numberOfRedemptions = parseFloat(
                        `${details.numberOfRedemptions}`
                    );
                    details.valueOfRedemptions = parseFloat(
                        `${details.valueOfRedemptions}`
                    );
                    details.totalValueOfSchemeAssets = parseFloat(
                        `${details.totalValueOfSchemeAssets}`
                    );
                    details.netAssetValueOfScheme = parseFloat(
                        `${details.netAssetValueOfScheme}`
                    );
                    details.netAssetValuePerShare = parseFloat(
                        `${details.netAssetValuePerShare}`
                    );
                    details.totalEquityInvestments = parseFloat(
                        `${details.totalEquityInvestments}`
                    );
                    details.totalFixedIncomeInvestments = parseFloat(
                        `${details.totalFixedIncomeInvestments}`
                    );
                    details.netOfMediumTermAssetsHeldByFund = parseFloat(
                        `${details.netOfMediumTermAssetsHeldByFund}`
                    );
                    details.capitalMarketsInvestments = parseFloat(
                        `${details.capitalMarketsInvestments}`
                    );
                    details.percentageOfCapitalInvestmentToTotalInvestment = parseFloat(
                        `${details.percentageOfCapitalInvestmentToTotalInvestment}`
                    );
                    details.totalValueOfUnutilizedFunds = parseFloat(
                        `${details.totalValueOfUnutilizedFunds}`
                    );
                    details.valueOfBorrowedFunds = parseFloat(
                        `${details.valueOfBorrowedFunds}`
                    );
                    details.redemptions = parseFloat(`${details.redemptions}`);
                    details.dividends = parseFloat(`${details.dividends}`);
                    details.rights = parseFloat(`${details.rights}`);
                    details.feesOwedCustodian = parseFloat(
                        `${details.feesOwedCustodian}`
                    );
                    details.attachedFile = this.activeFile.rawBase64Data;
                    const response = await store.dispatch(
                        CUSTODIAN_ACTION_CONSTANTS.UPLOAD_SCHEME_DETAILS,
                        details
                    );
                    this.handleApiResponse(response, "Scheme details uploaded");
                    console.log(response);
                    this.uploadedSchemeDetailsTemplate = [];
                }
            }
        },
        validateOfficialRemarks: function (): boolean {
            let isValid = true;
            this.officialRemarksError = {
                remarks: false,
                reviewingOfficer: false,
                date: false
            } as IOfficialRemarksError;
            if (!this.officialRemarks.remarks) {
                this.officialRemarksError.remarks = true;
                isValid = false;
            }
            if (!this.officialRemarks.reviewingOfficer) {
                this.officialRemarksError.reviewingOfficer = true;
                isValid = false;
            }
            if (!this.officialRemarks.date) {
                this.officialRemarksError.date = true;
                isValid = false;
            }
            return isValid;
        },
        validateOtherSecInfoData: function (): boolean {
            let hasError = false;
            for (let info in this.otherSecInformations) {
                if ((this.otherSecInformations as any)[info] == "") {
                    if (["litigationDetails", "id", "governaceInfoID"].includes(info)) {
                        continue;
                    }
                    this.$swal({
                        type: 'error',
                        title: "All inputs are required.",
                        animation: false
                    }).then(console.log);
                    hasError = true;
                    break;
                }
            }
            if (!hasError) {
                this.otherSecInformationsError.litigationDetails = false;
                if (
                    this.otherSecInformations
                        .companyParentsAffiliateInvolvedInLitigation === "yes"
                ) {
                    if (this.otherSecInformations.litigationDetails) {
                        return true;
                    } else {
                        this.otherSecInformationsError.litigationDetails = true;
                        return false;
                    }
                }
                return true;
            }
            return !hasError;
        },
        saveOfficialRemarks: async function () {
            if (this.validateOfficialRemarks()) {
                if (this.canEditData) {
                    const response = await store.dispatch(
                        CUSTODIAN_ACTION_CONSTANTS.UPLOAD_OFFICIAL_REPORT_REMARKS,
                        this.officialRemarks
                    );
                    this.handleApiResponse(response, "Official remarks uploaded");
                    console.log(response);
                }
            }
        },
        saveOtherSecInfo: async function () {
            if (this.validateOtherSecInfoData()) {
                if (this.canEditData) {
                    const response = await store.dispatch(
                        CUSTODIAN_ACTION_CONSTANTS.UPLOAD_OTHER_INFORMATION,
                        this.otherSecInformations
                    );
                    this.handleApiResponse(response, "Other info uploaded");
                    console.log(response);
                }
            }
        },
        addMoreSchemes: function () {
            this.schemeDetails = {
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
            };
            this.activeFile = {
                name: "",
                data: "",
                rawBase64Data: ""
            };
            this.filesUploaded = false;
            this.uploadedSchemeDetailsTemplate = [];
        },
        loadSchemeByBPID: async function () {
            this.schemeSearchError = false;
            this.uploadedSchemeDetailsTemplate = [];
            if (this.schemeSearch) {
                const response = await store.dispatch(
                    CUSTODIAN_ACTION_CONSTANTS.FETCH_SCHEME_DETAILS,
                    {
                        bpid: this.schemeSearch,
                        quarterDate: lastDayOfQuarter()
                    } as SchemeDetailsRequest
                );
                this.handleApiResponse(response, "Scheme details loaded");
                this.schemeDetails = response.data.schemeDetails;
                if (response.data.schemeDetails.attachedFile) {
                    this.activeFile = {
                        name: "Attached File.pdf",
                        data: URL.createObjectURL(
                            await dataURLToBlob(`data:application/pdf;base64,${response.data.schemeDetails.attachedFile}`)
                        ),
                        rawBase64Data: response.data.schemeDetails.attachedFile
                    };
                    this.filesUploaded = true;
                } else {
                    this.activeFile = {
                        name: "",
                        data: "",
                        rawBase64Data: ""
                    };
                    this.filesUploaded = false;
                }
            } else {
                this.schemeSearchError = true;
            }
        },
        previewReport: function () {
            window.location.assign(
                `/sec-report-preview`
            );
        },
        submitChanges: async function () {
            const result = await this.$swal({
                title: 'Are you sure?',
                text: "You won't be able to revert this!",
                type: 'warning',
                showCancelButton: true,
                confirmButtonColor: '#3085d6',
                cancelButtonColor: '#d33',
                confirmButtonText: 'Yes'
            });
            if (result.value) {
                const response = await store.dispatch(CUSTODIAN_ACTION_CONSTANTS.SEND_CHANGES_TO_CHECKER);
                this.handleApiResponse(response, "Changes Submitted");
                console.log(response)
            }
        },
        handleApiResponse: function (response: ApiResponse, successMessage: string) {
            if (response.error) {
                console.log(response);
                this.$swal({
                    toast: true,
                    position: 'top-end',
                    showConfirmButton: false,
                    timer: 3000,
                    type: 'error',
                    title: response.message
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
        onSchemeTemplateUpload: async function (event: any) {
            this.addMoreSchemes();
            const parser = new ExcelFilePreview(event.target.files[0]);
            const data = await parser.parseDataToJson();
            this.uploadedSchemeDetailsTemplate = [];
            const validData = data.filter((datum: Array<string>) => datum.length >= 7);
            this.schemeSearch = validData[1][6];
            await this.loadSchemeByBPID();
            this.uploadedSchemeDetailsTemplate = validData.map((datum: Array<string>) => {
                const data = [];
                data.push(datum[0]);
                data.push(datum[datum.length - 1]);
                return data;
            });
            this.schemeDetails.numberOfSharesOutstanding = parseStringToFloat(validData[2][6]);
            this.schemeDetails.numberOfShareholders = parseStringToFloat(validData[3][6]);
            this.schemeDetails.numberOfRedemptions = parseStringToFloat(validData[4][7]);
            this.schemeDetails.valueOfRedemptions = parseStringToFloat(validData[5][7]);
            this.schemeDetails.nameOfManager = validData[6][6];
            this.schemeDetails.netAssetValueOfScheme = parseStringToFloat(validData[7][6]);
            this.schemeDetails.netAssetValuePerShare = parseStringToFloat(validData[8][6]);
            this.schemeDetails.netOfMediumTermAssetsHeldByFund = parseStringToFloat(validData[9][6]);
            this.schemeDetails.areAllCertificatesOfInvestmentWithCustodian = validData[10][6].toLowerCase() == "yes" ? "yes" : "no";
            this.schemeDetails.totalValueOfUnutilizedFunds = parseStringToFloat(validData[11][6]);
            this.schemeDetails.valueOfBorrowedFunds = parseStringToFloat(validData[12][6]);
            this.schemeDetails.reasonsForBorrowing = validData[13][6];
            this.schemeDetails.wereAllDulyPreparedAccountsDistributed = validData[14][6].toLowerCase() == "yes" ? "yes" : "no";
            this.schemeDetails.redemptions = parseStringToFloat(validData[15][6]);
            this.schemeDetails.dividends = parseStringToFloat(validData[16][6]);
            this.schemeDetails.rights = parseStringToFloat(validData[17][6]);
            this.schemeDetails.feesOwedCustodian = parseStringToFloat(validData[18][6]);
        },
        deleteScheme: async function () {
            const result = await this.$swal({
                title: 'Are you sure?',
                text: "You won't be able to revert this!",
                type: 'warning',
                showCancelButton: true,
                confirmButtonColor: '#3085d6',
                cancelButtonColor: '#d33',
                confirmButtonText: 'Yes'
            });
            if (result.value) {
                const response = await store.dispatch(
                    CUSTODIAN_ACTION_CONSTANTS.DELETE_SCHEME_DETAILS,
                    {bpid: this.schemeDetails.bpid}
                );
                if (!response.error) {
                    this.addMoreSchemes();
                }
                this.handleApiResponse(response, "Scheme deleted");
            }
        },
        recalculateFormulas: async function () {
            const result = await this.$swal({
                title: 'Are you sure?',
                text: "You won't be able to revert this!",
                type: 'warning',
                showCancelButton: true,
                confirmButtonColor: '#3085d6',
                cancelButtonColor: '#d33',
                confirmButtonText: 'Yes'
            });
            if(result.value){
                const response = await store.dispatch(
                    CUSTODIAN_ACTION_CONSTANTS.RECALCULATE_FORMULAS,
                    {
                        bpid: this.schemeSearch,
                        quarterDate: lastDayOfQuarter()
                    } as SchemeDetailsRequest
                );
                this.handleApiResponse(response, "Scheme details recalculated");
                this.schemeDetails = response.data.schemeDetails;
                if (response.data.schemeDetails.attachedFile) {
                    this.activeFile = {
                        name: "Attached File.pdf",
                        data: URL.createObjectURL(
                            await dataURLToBlob(`data:application/pdf;base64,${response.data.schemeDetails.attachedFile}`)
                        ),
                        rawBase64Data: response.data.schemeDetails.attachedFile
                    };
                    this.filesUploaded = true;
                } else {
                    this.activeFile = {
                        name: "",
                        data: "",
                        rawBase64Data: ""
                    };
                    this.filesUploaded = false;
                }
            }
        }
    }
});
