import { createApp } from 'vue'
import './style.css'
import "@fortawesome/fontawesome-free/css/all.min.css"
import App from './App.vue'
import { createPinia } from "pinia"
import router from "./router";

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.mount('#app')
