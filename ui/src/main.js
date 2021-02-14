import Vue from 'vue'
import App from './App.vue'
import vuetify from './plugins/vuetify'
import router from './router'
import store from './store'
import axios from 'axios'
axios.defaults.withCredentials = false
axios.defaults.baseURL = 'http://127.0.0.1:9000';
axios.defaults.headers.post['Content-Type'] = 'application/json';
Vue.prototype.$http = axios;

Vue.config.productionTip = false

new Vue({
  router,
  store,
  vuetify,
  render: h => h(App)
}).$mount('#app')
