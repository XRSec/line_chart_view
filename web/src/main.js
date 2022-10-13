import {createApp} from 'vue'
import homePage from './App.vue'
import axios from 'axios'
import VueAxios from 'vue-axios'

createApp(homePage)
    .use(VueAxios, axios)
    .mount('#app')
