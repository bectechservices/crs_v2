import Vue from "vue";
import store from "./vuex"

export default new Vue({
    el: '.pvErrorsPage',
    store,
    computed: {
        uploadErrors: function () {
            return this.$store.getters.uploadedPvErrors;
        }
    }
})