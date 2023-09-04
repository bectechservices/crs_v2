import Vue from "vue";
import {PPTTemplate} from "./lib/types";
import store from './vuex';
import {PPT_OPTIONS_CONSTANTS} from "./vuex/constants";

interface Data {
    template: PPTTemplate;
    bpid: string;
    quarter: number;
    year: number;
    all_template_options: boolean;
}

interface Methods {
    showPPTDemo: () => void;
}

export default new Vue<Data, Methods>({
    el: ".pptTemplatePage",
    store,
    data: {
        template: {
            total_summary_of_auc: false,
            auc_trend: false,
            trade_volumes: false,
            pv_report: false,
            total_contribution: false,
            corporate_action: false,
            gog_and_fd_maturities: false,
            appendix_i: false,
            appendix_ii: false,
            unidentified_payments: false
        },
        bpid: "",
        quarter: 0,
        year: 0,
        all_template_options: false
    },
    mounted: function () {
        const url = new URL(window.location.href);
        this.bpid = url.searchParams.get("bpid") || "";
        this.quarter = parseInt(url.searchParams.get("quarter") as string) || 0;
        this.year = parseInt(url.searchParams.get("year") as string) || 0;
        // this.template = {...this.$store.state.ppt.templateOptions};
    },
    methods: {
        showPPTDemo: function () {
            if (this.bpid) {
                store.commit(PPT_OPTIONS_CONSTANTS.STORE_PPT_OPTIONS, {
                    templateOptions: this.template,
                    bpid: this.bpid,
                    year: this.year,
                    quarter: this.quarter
                });
                window.location.assign(`/client-pvreport?bpid=${this.bpid}&quarter=${this.quarter}&year=${this.year}`);
            }
        }
    },
    watch: {
        all_template_options: function (value) {
            this.template.total_summary_of_auc = value;
            this.template.auc_trend = value;
            this.template.trade_volumes = value;
            this.template.pv_report = value;
            this.template.total_contribution = value;
            this.template.corporate_action = value;
            this.template.gog_and_fd_maturities = value;
            this.template.appendix_i = value;
            this.template.appendix_ii = value;
            this.template.unidentified_payments = value;
        }
    }
});
