import Vue from 'vue'
import App from './App.vue'
import moment from 'moment'
import axios from 'axios'

Vue.config.productionTip = false
Vue.prototype.$http = axios

Vue.filter('formatDate', function(value) {
  if (value) {
      return moment(String(value)).format('MM/DD/YYYY')
  }
});

new Vue({
  render: h => h(App),
}).$mount('#app')
