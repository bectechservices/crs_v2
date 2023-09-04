import Vue from "vue";
import axios from "axios";
import {isExcel97File, isExcelFile, makeClientInfoForPVList, worker_script} from "./lib/helpers";
import {MaybeNull, ParsedPVReport} from "./lib/types";
import ExcelFilePreview from "./lib/excel";
import PVReportExcelParser from "./lib/pv-excel-parser";
import ShouldReadExcelFile from "./lib/IExcel";
import Excel97FilePreview from "./lib/excel97";
import VueSweetalert2 from 'vue-sweetalert2';
import store from "./vuex";
import {ROOT_STATE_CONSTANTS} from "./vuex/constants";
import WorkerPool from "./lib/WorkerPool";

Vue.use(VueSweetalert2);

interface Data {
    file: MaybeNull<File>;
    loading: boolean;
    hasLoadedExcelFile: boolean;
    excelFileData: Array<ParsedPVReport>;
    cashBalance: string;
    loadingFiles: boolean;
    excelFileList: Array<string>,
    hasLoadedExcelAsList: boolean,
    workerPool: MaybeNull<WorkerPool>;
}

interface Methods {
    uploadDocument: () => Promise<void>;
    onFileChange: (event: any) => Promise<void>;
    clearFileUploadContent: () => void;
    handleApiResponse: (response: { error: boolean, messages: Array<string> }) => void;
}

Vue.component("excel-data-preview", {
    template: `
      <span>
                <span v-if="reports.length">
                    <span v-for="(data,key) in reports" :key="key">
                    <table class="uk-table excelData">
                      <thead>
                        <tr>
                          <th colspan="9" rowspan="3" style="text-align: center">
                            <div v-for="(header,key) in data.headers" :key="key">{{header}}</div>
                          </th>
                        </tr>
                      </thead>
                      <tbody v-for="(bond,key) in data.bonds" :key="key">
                        <tr class="securityHead">
                          <th colspan="9">{{bond.bond}}</th>
                        </tr>
                        <tr class="securityCols">
                          <th>Security name</th>
                          <th>TENOR/CODE</th>
                          <th>isin</th>
                          <th>scb code</th>
                          <th>mkt price</th>
                          <th>nominal value</th>
                          <th>cumulative cost</th>
                          <th>value[ghs]</th>
                          <th>&percnt; of total</th>
                        </tr>
                          <tr class="securityColsData" v-for="(value,tKey) in bond.values" :key="tKey">
                            <td>{{value[0]}}</td>
                            <td>{{value[1]}}</td>
                            <td>{{value[2]}}</td>
                            <td>{{value[3]}}</td>
                            <td>{{value[4]}}</td>
                            <td>{{value[5]}}</td>
                            <td>{{value[6]}}</td>
                            <td>{{value[7]}}</td>
                            <td>{{value[8]}}</td>
                          </tr>
                        </tbody>
                      </table>
                    <table class="uk-table excelData excelDataSummary">
                      <thead>
                        <tr>
                          <th colspan="9" rowspan="3" style="text-align: center">
                            <div v-for="(header,key) in data.headers" :key="key">{{header}}</div>
                          </th>
                        </tr>
                      </thead>
                      <tr class="securityHead">
                        <th colspan="9">Summary</th>
                      </tr>
                      <tr class="securityCols">
                        <th>description</th>
                        <th>nominal value</th>
                        <th>cummulative cost</th>
                        <th>value[ghs]</th>
                        <th>&percnt; of total</th>
                      </tr>
                      <tbody>
                        <tr v-if="data.summary.length && data.summary[0].length == 6" class="securityColsData"
                            v-for="(sum,nKey) in data.summary" :key="nKey">
                          <td>{{sum[0]}}</td>
                          <td>{{sum[1]}}</td>
                          <td>{{sum[2] ? sum[2] : sum[3]}}</td>
                          <td>{{sum[4]}}</td>
                          <td>{{sum[5]}}</td>
                        </tr>
                        <tr v-if="data.summary.length && data.summary[0].length == 5" class="securityColsData"
                            v-for="(sum,nKey) in data.summary" :key="nKey">
                          <td>{{sum[0]}}</td>
                          <td>{{sum[1]}}</td>
                          <td>{{sum[2]}}</td>
                          <td>{{sum[3]}}</td>
                          <td>{{sum[4]}}</td>
                        </tr>
                      </tbody>
                    </table>
                    <div class="userSelectNone" style="font-size: 13px;letter-spacing: 0.25px;font-weight: 400;">THE
                      PRICES QUOTED ARE INTENDED FOR INTERNAL FOR ADMINISTRATIVE AND VALUATION PROCESS.SCB ACCEPTS
                      NO
                      LIABILITY FOR ITS ACCURACY AND COMPLETENESS</div>
                  </span>
                </span>
                <span v-else>
                    <div style="    display: flex;justify-content: center;padding: 40px">
                        <p style="font-weight: 500">No data to report for PV</p>
                    </div>
                </span>
            </span>
    `,
    props: ["reports"],
    data: function () {
        return {};
    }
});
Vue.component("excel-data-list", {
    template: `
      <span>
                <ol v-if="reports.length > 0">
                    <li v-for="(report,key) in reports" :key="key">{{report}}</li>
                </ol>
                <div style="display: flex;justify-content: center;padding: 40px" v-else>
                        <p style="font-weight: 500">No data to report for PVs</p>
                    </div>
            </span>
    `,
    props: ["reports"],
    data: function () {
        return {};
    }
});

export default new Vue<Data, Methods>({
    el: ".uploadablePVPage",
    store,
    data: {
        file: null,
        loading: false,
        hasLoadedExcelFile: false,
        excelFileData: [],
        cashBalance: "0",
        loadingFiles: false,
        hasLoadedExcelAsList: false,
        excelFileList: [],
        workerPool: null
    },
    mounted: function () {
        if (WorkerPool.hasWebWorkerSupport()) {
            this.workerPool = new WorkerPool(worker_script(), 1);
            this.workerPool.startWorkers();
        }
        
    },
    beforeDestroy: function () {
        this.workerPool && this.workerPool.stopWorkers();
    },
    methods: {
        uploadDocument: async function () {
            
            const pvType = (this.$refs.pvTypeInput as HTMLInputElement).value;
            if (this.hasLoadedExcelFile && pvType) {

                this.loading = true;
                try {
                    const response = await axios.post("/portal", {
                        report_type: pvType,
                        cashbalance: parseFloat(this.cashBalance || "0"),
                        data: this.excelFileData
                    }, {
                        headers: {
                            "Content-Type": "application/json"
                        }
                    });
                    this.handleApiResponse(response.data);
                } catch (error) {
                    console.error(error.response.data);
                } finally {
                    this.loading = false;
                }
            }
        }
        ,
        onFileChange: async function (event: any) {
            this.hasLoadedExcelFile = false;
            this.loadingFiles = true;
            this.excelFileData = [];
            this.excelFileList = [];
            const length = event.target.files.length;
            if (WorkerPool.hasWebWorkerSupport() && length > 1) {
                (this.workerPool as WorkerPool).run<FileList, ParsedPVReport>([...event.target.files], (data: Array<ParsedPVReport>) => {
                    this.excelFileData = [...this.excelFileData, ...data];
                    this.excelFileList = [...this.excelFileList, ...data.map(each => makeClientInfoForPVList(each.headers))]
                }, () => {
                    this.$nextTick(() => {
                        this.loadingFiles = false;
                        this.hasLoadedExcelFile = true;
                        this.hasLoadedExcelAsList = true;
                        const button = document.querySelector('.has-loader') as any;
                        if (button) {
                            button.addEventListener('click', function () {
                                //@ts-ignore
                                this.classList.toggle('active');
                            });
                        }

                    });
                })
            } else {
                for (let i = 0; i < length; i++) {
                    this.file = event.target.files[i];
                    let preview: ShouldReadExcelFile;
                    if (isExcelFile(this.file as File)) {
                        if (isExcel97File(this.file as File)) {
                            preview = new Excel97FilePreview(this.file as File);
                        } else {
                            preview = new ExcelFilePreview(this.file as File);
                        }
                        let results: string[][] = [];
                        try {
                            results = await preview.parseDataToJson(true);
                            console.log("results 1: ", results);
                        } catch (e) {
                            if (isExcel97File(this.file as File)) {
                                preview = new ExcelFilePreview(this.file as File);
                                results = await preview.parseDataToJson(true);
                                console.log("results 2: ", results);
                            }
                        }
                        if (results.length && (results[0][0] === '___PATCHED_PV__')) {
                            results.shift();
                            results.pop();
                            results = Excel97FilePreview.parseConvertedPVToAKnownFormat(results);
                            console.log("results 2: ", results);
                        }
                        this.excelFileData = [...this.excelFileData, ...PVReportExcelParser.parsePVReport(
                            results
                        )];
                    }
                }
                this.$nextTick(() => {
                    this.loadingFiles = false;
                    this.hasLoadedExcelFile = true;
                    this.hasLoadedExcelAsList = false;
                    const button = document.querySelector('.has-loader') as any;
                    if (button) {
                        button.addEventListener('click', function () {
                            //@ts-ignore
                            this.classList.toggle('active');
                        });
                    }

                });
            }
        }
        ,
        clearFileUploadContent: function () {
            this.hasLoadedExcelFile = false;
            this.loading = false;
            const button = document.querySelector('.has-loader') as any;
            if (button) {
                button.classList.toggle('active');
            }
        }
        ,
        handleApiResponse: function (response: { error: boolean, messages: Array<string> }) {
            let swalObj = {
                type: 'success',
                title: 'PV Report Uploaded'
            };
            if (response.error) {
                swalObj = {
                    ...swalObj,
                    showCancelButton: true,
                    cancelButtonText: 'Ok',
                    confirmButtonText: 'Show errors'
                } as any
            }
            this.$swal(swalObj as any).then((result) => {
                if (result.value && response.error) {
                    this.$store.commit(ROOT_STATE_CONSTANTS.STORE_PV_UPLOAD_ERRORS, response.messages);
                    window.open('/pv-upload-errors', '_blank')
                }
                window.location.reload();
            })
        }
        ,
    }
});
