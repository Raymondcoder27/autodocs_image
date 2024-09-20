import {defineStore} from "pinia";
import type {AxiosResponse} from "axios";
import type {ApiResponse} from "@/types";
import {ref, type Ref} from "vue";
import type {Doc, GenerationRequest} from "@/domain/documents/types";
import api from "@/config/api";

export const useDocumentStore = defineStore("documents", () => {

    const logs: Ref<Doc[] | undefined> = ref()

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
