import Vue from "vue";
import store from './vuex';
import JabPPT from "./lib/JabPPT";

interface Data {
    client: any
    misc: any
}

interface Methods {
    exportToPPT: () => void;
}

export default new Vue<Data, Methods>({
    el: ".pvReportPage",
    store,
    data: {
        client: {},
        misc: {}
    },
    mounted() {
        let clients: any = [];
        let misc: any = [];
        for (let i = 0; i < (window as any).NumberOfRecords; i++) {
            clients.push((window as any)[`CRS_CLIENT_${i}`]);
            misc.push((window as any)[`CRS_PPT_MISC_${i}`]);
        }
        this.client = clients;
        this.misc = misc;
    },
    methods: {
        exportToPPT: function () {
            const jab = new JabPPT(this.client, this.misc, store.getters.pptTemplateOptions);
            jab.generatePPT();
        }
    },
    computed: {
        templateOptions: function () {
            return store.getters.pptTemplateOptions;
        }
    }
});
