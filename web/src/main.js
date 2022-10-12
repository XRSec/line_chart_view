import { createApp } from 'vue'
import homePage from './App.vue'
import {createPinia} from "pinia";
const pinia = createPinia()
import axios from 'axios'
import VueAxios from 'vue-axios'

createApp(homePage).use(pinia).use(VueAxios,axios).mount('#app')
