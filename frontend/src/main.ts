import { createApp } from 'vue'
import './style.css'
import "@fortawesome/fontawesome-free/css/all.min.css"
import App from './App.vue'
import { createPinia } from "pinia"
import router from "./router";
import CanvasJSChart from '@canvasjs/vue-chartjs'

const app = createApp(App)
app.use(createPinia())
app.use(CanvasJSChart)
app.use(router)
app.mount('#app')
