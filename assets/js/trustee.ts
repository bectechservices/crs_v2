import Vue from "vue";
import axios from "axios";
import { Chart } from "chart.js";
import {
  currentQuarterNumber,
  currentYear,
  formatMoney,
  makeLastDayOfQuarter,
  makeTrusteeQuarterFancyQuarterDateFromDate,
  parseStringToFloat,
  parseStringToInt,
  sumMemberVariableInObjectArray,
} from "./lib/helpers";
import {
  FundManager,
  IMonthlyContributions,
  IUnidentifiedPayment,
  IUnidentifiedPaymentSummary,
  TrusteeDashboardSearchData,
  TrusteePerformance,
  TrusteeQuarterlyReport,
  TrusteeUploadedPV,
} from "./lib/types";
import { format, parse } from "date-fns";
import VueSweetalert2 from "vue-sweetalert2";
import ExcelFilePreview from "./lib/excel";
import * as moment from "moment";
import { TRUSTEE_CONSTANTS } from "./vuex/constants";
import store from "./vuex";

interface DataIndices {
  valid: boolean;
  indices: Array<number>;
}

Vue.use(VueSweetalert2);

interface Methods {
  loadPageData: () => Promise<void>;
  formatMoney: (amount: number) => string;
  loadPerformanceChart: () => void;
  updatePerformanceCharts: () => void;
  addClientMonthlyContribution: () => Promise<void>;
  formatDate: (date: string, format: string) => string;
  showPPTOptions: () => void;
  handleApiResponse: (
    response: { error: boolean },
    successMessage: string
  ) => void;
  parseGOGMaturities: (event: any) => Promise<void>;
  uploadGOGMaturities: () => Promise<void>;
  parseTransactionVolumes: (event: any) => Promise<void>;
  uploadTransactionVolumes: () => Promise<void>;
  deleteMonthlyContribution: (id: number) => Promise<void>;
  selectMonthlyContributionForEdit: (
    contribution: IMonthlyContributions
  ) => void;
  updateMonthlyContribution: () => Promise<void>;
  sumOf: (data: Array<any>, field: string) => number;
  getGOGDataIndices: (data: Array<string>) => DataIndices;
  getTxnVolsDataIndices: (data: Array<string>) => DataIndices;
  generateMultipleReport: () => void;
  makeTrusteeQuarterFancyQuarterDateFromDate: (date: string) => string;
}

interface ClientMergedBPAndSCA {
  bpid: string;
  client_name: string;
  code: string;
  fund_manager: string;
}

interface Data {
  clientQuarterlyReport: TrusteeQuarterlyReport;
  pv: {
    bpid: string;
    quarter: string;
    year: string;
  };
  chartJS: any;
  pvResponse: {
    pvData: Array<{
      security_type: string;
      account_number: string;
      report_date: string;
      nominal_value: number;
      cumlative_cost: number;
      lcy_amount: number;
      percentage_of_total: number;
    }>;
    quarterlyPerformance: Array<TrusteePerformance>;
    pvAsAt: string;
    monthlyContributions: Array<IMonthlyContributions>;
    clientName: string;
    safekeepingAccount: string;
    unidentifiedPayments: Array<
      IUnidentifiedPayment & { fundManager: FundManager }
    >;
    unidentifiedPaymentsSummary: Array<IUnidentifiedPaymentSummary>;
  };
  monthlyContributions: {
    bpid: string;
    date: string;
    quarter: string;
    year: string;
    amount: string;
    sca: string;
    data: [];
  };
  selectedClientSCA: Array<ClientMergedBPAndSCA>;
  gogMaturities: Array<Array<string>>;
  gogMaturitiesDate: {
    quarter: string;
    year: string;
  };
  transactionVolumes: Array<Array<string>>;
  transactionVolumesDate: {
    quarter: string;
    year: string;
  };
  hasSelectedFile: boolean;
  hasSelectedTxnVolumesFile: boolean;
  selectedMonthlyContributionForEdit: IMonthlyContributions;
  gogDataIndices: Array<number>;
  txnVolsDataIndices: Array<number>;
  multipleReport: {
    tier2: string;
    tier3: string;
  };
}

export default new Vue<Data, Methods>({
  el: ".trusteePage",
  store,
  data: {
    pv: {
      bpid: "",
      quarter: currentQuarterNumber(),
      year: `${currentYear()}`,
    },
    pvResponse: {
      pvData: [],
      quarterlyPerformance: [],
      pvAsAt: "",
      monthlyContributions: [],
      clientName: "",
      safekeepingAccount: "",
      unidentifiedPayments: [],
      unidentifiedPaymentsSummary: [],
    },
    chartJS: null,
    monthlyContributions: {
      bpid: "",
      date: "",
      quarter: currentQuarterNumber(),
      year: `${currentYear()}`,
      amount: "",
      sca: "",
      data: [],
    },
    selectedClientSCA: [],
    gogMaturities: [],
    gogMaturitiesDate: {
      quarter: currentQuarterNumber(),
      year: `${currentYear()}`,
    },
    transactionVolumes: [],
    transactionVolumesDate: {
      quarter: currentQuarterNumber(),
      year: `${currentYear()}`,
    },
    hasSelectedFile: false,
    hasSelectedTxnVolumesFile: false,
    clientQuarterlyReport: {
      id: 0,
      bpid: "",
      approved: false,
      approved_by: 0,
      quarter: "",
    },
    selectedMonthlyContributionForEdit: {
      id: 0,
      bpid: "",
      date: "",
      amount: 0,
      sca: "",
      created_at: "",
    },
    gogDataIndices: [],
    txnVolsDataIndices: [],
    multipleReport: {
      tier2: "",
      tier3: "",
    },
  },
  mounted: async function() {
    this.loadPerformanceChart();
    const data: TrusteeDashboardSearchData = this.$store.state.trustee
      .searchInputData;
    if (data.bpid && data.quarter && data.year) {
      this.pv.bpid = data.bpid;
      // this.pv.quarter = data.quarter;
      // this.pv.year = data.year;
      await this.loadPageData();
    }
  },
  
  methods: {
    showPPTOptions: function() {
      const bpid = this.pv.bpid;
      const quarter = this.pv.quarter;
      const year = this.pv.year;
      if (bpid && quarter && year) {
        window.location.assign(
          `/template-setup?bpid=${bpid}&quarter=${quarter}&year=${year}`
        );
      }
    },
    formatDate: function(date: string, _format: string): string {
      return format(parse(date), _format);
    },
    loadPageData: async function() {
      const searchData: TrusteeDashboardSearchData = {
        bpid: this.pv.bpid,
        quarter: this.pv.quarter,
        year: this.pv.year,
      };
      this.$store.commit(
        TRUSTEE_CONSTANTS.STORE_DASHBOARD_SEARCH_INPUT,
        searchData
      );
      this.hasSelectedFile = false;
      this.hasSelectedTxnVolumesFile = false;
      let response = await axios.post("/trustee-data", searchData);
      this.handleApiResponse(response.data, "Data loaded");
      if (!response.data.error) {
        this.gogMaturities = response.data.data.gog.map((each: any) => {
          let data = [];
          data[0] = each.depot_id;
          data[1] = each.entry_date;
          data[2] = each.event_type;
          data[3] = each.base_security_id;
          data[4] = each.gross_amount;
          data[5] = each.status;
          return data;
        });
        this.transactionVolumes = response.data.data.transactions.map(
          (each: any) => {
            let data = [];
            data[0] = each.stock_settled_date;
            each[1] = each.security_type;
            return data;
          }
        );
        this.selectedClientSCA = response.data.data.scas;
        this.pvResponse.pvData = response.data.data.summary;
        this.pvResponse.quarterlyPerformance = response.data.data.performance;
        this.pvResponse.monthlyContributions = response.data.data.contributions;
        this.pvResponse.clientName = response.data.data.client.name;
        this.pvResponse.safekeepingAccount =
          response.data.data.client.account_number;
        this.pvResponse.unidentifiedPayments =
          response.data.data.unidentifiedPayments;
        this.pvResponse.unidentifiedPaymentsSummary =
          response.data.data.unidentifiedPaymentsSummary;
        this.pvResponse.pvAsAt = makeLastDayOfQuarter(
          this.pv.quarter,
          this.pv.year
        );
        this.monthlyContributions.bpid = this.pv.bpid;
        this.clientQuarterlyReport = response.data.data.report;
        this.$nextTick(() => {
          this.updatePerformanceCharts();
        });
      }
    },
    loadPerformanceChart: function() {
      let setup = (canvas: any) => {
        let context, dpr;
        context = canvas.getContext("2d");
        canvas.style.width = "100%";
        canvas.style.height = "100%";
        canvas.style.marginTop = "5px";
        dpr = window.devicePixelRatio || 1.4;
        canvas.width = canvas.offsetWidth * dpr;
        canvas.height = canvas.offsetHeight * dpr;
        context.scale(dpr, dpr);
        return context;
      };
      let quarterlyPerformance = setup(
        document.querySelector("#quarterlyPerformance")
      );
      this.chartJS = new Chart(quarterlyPerformance, {
        type: "bar",
        data: {
          labels: [],
          datasets: [
            {
              type: "bar",
              backgroundColor: [],
              data: [0, 0, 0, 0, 0, 0, 0, 0, 0],
            },
          ],
        },
        options: {
          title: {
            display: false,
          },
          tooltips: {
            enabled: false,
          },
          hover: {
            animationDuration: 0,
          },
          scales: {
            xAxes: [
              {
                categoryPercentage: 1.0,
                barPercentage: 1.0,
              },
            ],
            yAxes: [
              {
                ticks: {
                  callback: function(value: any) {
                    if (value >= 1000 && value < 1000000) {
                      return `${value / 1000}K`;
                    } else if (value >= 1000000 && value < 1000000000) {
                      return `${value / 1000000}M`;
                    } else if (value >= 1000000000) {
                      return `${value / 1000000000}B`;
                    } else if (value < 2) {
                      return parseFloat(value).toFixed(2);
                    }
                  },
                },
              },
            ],
          },
          legend: { display: false },
          centertext: "",
        },
        layout: {
          padding: {
            left: 35,
            right: 35,
            top: 0,
            bottom: 0,
          },
        },
        animation: {
          onComplete: function() {
            var chartInstance = (this as any).chart;
            var ctx = chartInstance.ctx;
            ctx.textAlign = "center";
            ctx.font = "12px Open Sans";
            ctx.fillStyle = "#fff";
            Chart.helpers.each(
              (this as any).data.datasets.forEach(function(
                dataset: any,
                i: any
              ) {
                var meta = chartInstance.controller.getDatasetMeta(i);
                Chart.helpers.each(
                  meta.data.forEach(function(bar: any, index: any) {
                    let data = dataset.data[index];
                    if (i == 0) {
                      ctx.fillText(data, bar._model.x - 2, bar._model.y + 50);
                    } else {
                      ctx.fillText(data, bar._model.x - 2, bar._model.y + 50);
                    }
                  }),
                  //@ts-ignore
                  this
                );
              }),
              this
            );
          },
        },
        pointLabelFontFamily: "Quadon Extra Bold",
        scaleFontFamily: "Quadon Extra Bold",
      } as any);
    },
    formatMoney: function(amount: number) {
      return formatMoney(amount);
    },
    updatePerformanceCharts: function() {
      const quarterNumber = parseInt(this.pv.quarter);
      const yearNumber = parseInt(this.pv.year);

      const qp = `Q${quarterNumber === 1 ? 4 : quarterNumber - 1}-${
        quarterNumber === 1 ? yearNumber - 1 : yearNumber
      }`;
      const qc = `Q${this.pv.quarter}-${this.pv.year}`;

      this.chartJS.data.labels = this.pvResponse.quarterlyPerformance
        .map((value: TrusteePerformance) => [qp, qc, ""])
        .flat();
      this.chartJS.data.datasets[0].backgroundColor = this.pvResponse.quarterlyPerformance
        .map((value: TrusteePerformance) => ["#039be5", "#0288d1", ""])
        .flat();
      this.chartJS.data.datasets[0].data = this.pvResponse.quarterlyPerformance
        .map((value: TrusteePerformance) => [
          value.previous_quarter,
          value.current_quarter,
          0,
        ])
        .flat();
      (this.chartJS as Chart).update();
    },
    addClientMonthlyContribution: async function() {
      const result = await this.$swal({
        title: "Are you sure?",
        text: "This will update the client's monthly contribution!",
        type: "warning",
        showCancelButton: true,
        confirmButtonColor: "#3085d6",
        cancelButtonColor: "#d33",
        confirmButtonText: "Yes",
      });
      if (result.value) {
        const response = await axios.post("/client-monthly-contribution", {
          bpid: this.monthlyContributions.bpid,
          date: this.monthlyContributions.date,
          quarter: this.monthlyContributions.quarter,
          year: this.monthlyContributions.year,
          amount: parseFloat(this.monthlyContributions.amount),
          sca: this.monthlyContributions.sca,
        });
        this.handleApiResponse(response.data, "Contribution Added");
        if (!response.data.error) {
          this.pvResponse.monthlyContributions = [
            ...this.pvResponse.monthlyContributions,
            response.data.data.contribution,
          ];
          this.monthlyContributions.date = "";
          this.monthlyContributions.amount = "";
          this.monthlyContributions.sca = "";
        }
      }
    },
    handleApiResponse: function(
      response: { error: boolean },
      successMessage: string
    ) {
      if (response.error) {
        console.log(response);
        this.$swal({
          toast: true,
          position: "top-end",
          showConfirmButton: false,
          timer: 3000,
          type: "error",
          title: "Something went wrong. Please try again.",
        }).then(console.log);
      } else {
        this.$swal({
          toast: true,
          position: "top-end",
          showConfirmButton: false,
          timer: 3000,
          type: "success",
          title: successMessage,
        }).then(console.log);
      }
    },
    parseGOGMaturities: async function(event: any) {
      try {
        const file = event.target.files[0];
        this.hasSelectedFile = true;
        const parser = new ExcelFilePreview(file);
        const content = await parser.parseDataToJson();
        //we try to ignore all first lines less than the number of headers we're expecting.
        let headers: Array<any> = [];
        while (headers.length < 6) {
          headers = content.shift() as Array<string>;
        }
        const gogIndices = this.getGOGDataIndices(headers);
        if (!gogIndices.valid) {
          let message = "";
          const indicies = gogIndices.indices;
          if (indicies[0] == -1) {
            message += "Depot ID, ";
          }
          if (indicies[1] == -1) {
            message += "Entry Date, ";
          }
          if (indicies[2] == -1) {
            message += "Event Type, ";
          }
          if (indicies[3] == -1) {
            message += "Base Security ID, ";
          }
          if (indicies[4] == -1) {
            message += "Gross Amount, ";
          }
          if (indicies[5] == -1) {
            message += "Status";
          }
          throw new Error(`Headers: ${message} missing`);
        }
        this.gogDataIndices = gogIndices.indices;
        this.gogMaturities = content.map((each: any) => {
          let parsedContent = [];
          parsedContent[0] = each[this.gogDataIndices[0]] || "";
          parsedContent[1] = parse(each[this.gogDataIndices[1]]);
          parsedContent[2] = each[this.gogDataIndices[2]] || "";
          parsedContent[3] = each[this.gogDataIndices[3]] || "";
          parsedContent[4] = parseStringToFloat(each[this.gogDataIndices[4]]);
          parsedContent[5] = each[this.gogDataIndices[5]] || "";
          return parsedContent;
        });
      } catch (e) {
        this.hasSelectedFile = false;
        this.$swal({
          toast: true,
          position: "top-end",
          showConfirmButton: false,
          timer: 5000,
          type: "error",
          title: e.message,
        }).then(console.log);
      }
    },
    uploadGOGMaturities: async function() {
      const result = await this.$swal({
        title: "Are you sure?",
        text: "This will update the client's GOG maturities!",
        type: "warning",
        showCancelButton: true,
        confirmButtonColor: "#3085d6",
        cancelButtonColor: "#d33",
        confirmButtonText: "Yes",
      });
      if (result.value) {
        const response = await axios.post("/upload-gog-maturities", {
          quarter: this.gogMaturitiesDate.quarter,
          year: this.gogMaturitiesDate.year,
          data: this.gogMaturities.map((each: any) => {
            return {
              bpid: this.monthlyContributions.bpid,
              depot_id: each[0],
              entry_date: each[1],
              event_type: each[2],
              base_security_id: each[3],
              gross_amount: each[4],
              status: each[5],
            };
          }),
        });
        this.handleApiResponse(response.data, "Maturities uploaded");
      }
    },
    parseTransactionVolumes: async function(event: any) {
      try {
        const file = event.target.files[0];
        this.hasSelectedTxnVolumesFile = true;
        const parser = new ExcelFilePreview(file);
        const content = await parser.parseDataToJson();
        const filteredContent = content.filter(
          (each: Array<string>) => each.length >= 10
        ); //should have a reasonable number of fields
        const txnIndices = this.getTxnVolsDataIndices(
          filteredContent.shift() as Array<string>
        );
        if (!txnIndices.valid) {
          let message = "";
          const indices = txnIndices.indices;
          if (indices[0] == -1) {
            message += "Stock Settled Date, ";
          }
          if (indices[1] == -1) {
            message += "Security Type or Asset Class";
          }
          throw new Error(`Headers: ${message} missing`);
        }
        this.txnVolsDataIndices = txnIndices.indices;
        this.transactionVolumes = filteredContent.map((_each: any) => {
          let each = [];
          each[0] = moment(_each[this.txnVolsDataIndices[0]], "M/D/YYYY");
          each[1] = _each[this.txnVolsDataIndices[1]] || "";
          return each;
        });
      } catch (e) {
        this.hasSelectedTxnVolumesFile = false;
        this.$swal({
          toast: true,
          position: "top-end",
          showConfirmButton: false,
          timer: 5000,
          type: "error",
          title: e.message,
        }).then(console.log);
      }
    },
    uploadTransactionVolumes: async function() {
      const result = await this.$swal({
        title: "Are you sure?",
        text: "This will update the client's transaction volumes!",
        type: "warning",
        showCancelButton: true,
        confirmButtonColor: "#3085d6",
        cancelButtonColor: "#d33",
        confirmButtonText: "Yes",
      });
      if (result.value) {
        const response = await axios.post("/upload-txn-volumes", {
          quarter: this.transactionVolumesDate.quarter,
          year: this.transactionVolumesDate.year,
          data: this.transactionVolumes.map((each: any) => {
            return {
              bpid: this.monthlyContributions.bpid,
              stock_settled_date: each[0],
              security_type: each[1],
            };
          }),
        });
        this.handleApiResponse(response.data, "Transaction volumes uploaded");
      }
    },
    deleteMonthlyContribution: async function(id: number) {
      const result = await this.$swal({
        title: "Are you sure?",
        text: "You won't be able to revert this!",
        type: "warning",
        showCancelButton: true,
        confirmButtonColor: "#3085d6",
        cancelButtonColor: "#d33",
        confirmButtonText: "Yes",
      });
      if (result.value) {
        const response = await axios.delete("/trustee-monthly-contribution", {
          data: {
            id,
          },
        });
        if (!response.data.error) {
          this.pvResponse.monthlyContributions = [
            ...this.pvResponse.monthlyContributions.filter(
              (each: IMonthlyContributions) => each.id !== id
            ),
          ];
        }
        this.handleApiResponse(response.data, "Monthly contribution deleted");
      }
    },
    selectMonthlyContributionForEdit: function(
      contribution: IMonthlyContributions
    ) {
      this.selectedMonthlyContributionForEdit = {
        ...contribution,
        date: this.formatDate(contribution.date, "YYYY-MM-DD"),
      };
    },
    updateMonthlyContribution: async function() {
      const result = await this.$swal({
        title: "Are you sure?",
        text: "This will modify the clients monthly contribution!",
        type: "warning",
        showCancelButton: true,
        confirmButtonColor: "#3085d6",
        cancelButtonColor: "#d33",
        confirmButtonText: "Yes",
      });
      if (result.value) {
        const response = await axios.patch("/trustee-monthly-contribution", {
          id: this.selectedMonthlyContributionForEdit.id,
          bpid: this.selectedMonthlyContributionForEdit.bpid,
          date: this.selectedMonthlyContributionForEdit.date,
          amount: parseFloat(this.selectedMonthlyContributionForEdit
            .amount as any),
          sca: this.selectedMonthlyContributionForEdit.sca,
        });
        if (!response.data.error) {
          this.pvResponse.monthlyContributions = [
            ...this.pvResponse.monthlyContributions.map(
              (each: IMonthlyContributions) => {
                if (each.id == this.selectedMonthlyContributionForEdit.id) {
                  return response.data.data.contribution;
                }
                return each;
              }
            ),
          ];
        }
        this.handleApiResponse(response.data, "Monthly contribution edited");
      }
    },
    sumOf: function(data: Array<any>, field: string): number {
      return sumMemberVariableInObjectArray(data, field);
    },
    getGOGDataIndices: function(data: Array<string>): DataIndices {
      const result: DataIndices = { valid: true, indices: [] };
      const headers = [
        "depotid",
        "entrydate|postingdate",
        "eventtype",
        "basesecurityid",
        "grossamount",
        "status",
      ];
      const currentHeaders = data.map((each: string) =>
        each.replace(/\s/g, "").toLowerCase()
      );
      const indices: Array<number> = [];
      headers.forEach((header: string) => {
        let foundIndex = -1;
        header.split("|").forEach((each: string) => {
          const index = currentHeaders.indexOf(each);
          if (index !== -1) {
            foundIndex = index;
          }
        });
        if (foundIndex === -1 && result.valid) {
          result.valid = false;
        }
        indices.push(foundIndex);
      });
      result.indices = indices;
      return result;
    },
    getTxnVolsDataIndices: function(data: Array<string>): DataIndices {
      const result: DataIndices = { valid: true, indices: [] };
      const headers = ["stocksettleddate", "securitytype|assetclass"];
      const currentHeaders = data.map((each: string) =>
        each.replace(/\s/g, "").toLowerCase()
      );
      const indices: Array<number> = [];
      headers.forEach((each: string) => {
        let index;
        if (each.includes("|")) {
          index = each
            .split("|")
            .map((each) => {
              return currentHeaders.indexOf(each);
            })
            .find((each) => each != -1);
        } else {
          index = currentHeaders.indexOf(each);
        }
        if (index === -1 && result.valid) {
          result.valid = false;
        }
        indices.push(index as any);
      });
      result.indices = indices;
      return result;
    },
    generateMultipleReport: function(): void {
      const quarter = this.pv.quarter;
      const year = this.pv.year;
      const tier2 = this.multipleReport.tier2;
      const tier3 = this.multipleReport.tier3;
      if (tier2 && tier3 && quarter && year) {
        window.location.assign(
          `/template-setup?bpid=${tier2},${tier3}&quarter=${quarter}&year=${year}`
        );
      }
    },
    makeTrusteeQuarterFancyQuarterDateFromDate: function(date: string) {
      return makeTrusteeQuarterFancyQuarterDateFromDate(date);
    },
  },
});
