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
            .then((response: AxiosResponse<ApiResponse<Log[]>>) => {
                logs.value = response.data.data
            })
    }

    const clearLogs = async (ref:string) => {
        return api.delete("/clear-logs")
            .then((response: AxiosResponse<ApiResponse<any>>) => {
                uploadResponse.value = response.data.data
            })
    }

    return {
        logs,
        fetchLogs,
    }
})
