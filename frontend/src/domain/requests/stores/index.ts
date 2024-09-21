import {defineStore} from "pinia";
import type {AxiosResponse} from "axios";
import type {ApiResponse} from "@/types";
import type {Log} from "@/domain/requests/types";
import {ref, type Ref} from "vue";
import api from "@/config/api";

export const useLogStore = defineStore("logs", () => {

    const logs: Ref<Log[] | undefined> = ref()

    const fetchLogs = async () => {
        return api.get("/logs")
            .then((response: AxiosResponse<ApiResponse<Doc[]>>) => {
                logs.value = response.data.data
            })
    }
    return {
        fetchLogs,
    }
})
