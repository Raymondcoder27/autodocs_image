import {defineStore} from "pinia";
import type {AxiosResponse} from "axios";
import type {ApiResponse} from "@/types";
import {ref, type Ref} from "vue";
import api from "@/config/api";
import type {Template} from "@/domain/templates/types";

export const useTemplateStore = defineStore("templates", () => {

    const templates: Ref<Template[] | undefined> = ref()
    const fileBase64: Ref<string | undefined> = ref()
    const uploadResponse: Ref<object | undefined> = ref()

    const fetchTemplates = async () => {
        return api.get("/templates")
            .then((response: AxiosResponse<ApiResponse<Template[]>>) => {
                templates.value = response.data.data
            })
    }

    const fetchTemplateFile = async (ref:string) => {
        return api.get("/templates/preview/"+ref)
            .then((response: AxiosResponse<ApiResponse<string>>) => {
                fileBase64.value = response.data.data
            })
    }

    const uploadTemplate = async (payload:FormData) => {
        return api.post("/upload-template", payload)
            .then((response: AxiosResponse<ApiResponse<any>>) => {
                uploadResponse.value = response.data.data
            })
    }

    const deleteTemplate = async (ref:string) => {
        return api.delete("/templates/"+ref)
            .then((response: AxiosResponse<ApiResponse<any>>) => {
                uploadResponse.value = response.data.data
            })
    }

    return {
        templates,
        uploadResponse,
        fileBase64,
        fetchTemplateFile,
        fetchTemplates,
        deleteTemplate,
        uploadTemplate,
    }
})