import {defineStore} from "pinia";
import type {AxiosResponse} from "axios";
import type {ApiResponse} from "@/types";
import {ref, type Ref} from "vue";
import type {Doc, GenerationRequest} from "@/domain/documents/types";
import api from "@/config/api";

export const useDocumentStore = defineStore("documents", () => {

    const documents: Ref<Doc[] | undefined> = ref()
    const fileBase64: Ref<string | undefined> = ref()
    const generationResponse: Ref<object | undefined> = ref()

    const fetchDocuments = async () => {
        return api.get("/documents")
            .then((response: AxiosResponse<ApiResponse<Doc[]>>) => {
                documents.value = response.data.data
            })
    }

    const fetchDocumentFile = async (ref:string) => {
        return api.get("/documents/preview/"+ref)
            .then((response: AxiosResponse<ApiResponse<string>>) => {
                fileBase64.value = response.data.data
            })
    }

    const sendRequest = async (payload: GenerationRequest) => {
        return api.post("/generate", payload)
            .then((response: AxiosResponse<ApiResponse<any>>) => {
                generationResponse.value = response.data
            })
    }

    const deleteDocument = async (ref:string) => {
        return api.delete("/documents/"+ref)
            .then((response: AxiosResponse<ApiResponse<any>>) => {
                generationResponse.value = response.data
            })
    }


    return {
        documents,
        generationResponse,
        fileBase64,
        fetchDocuments,
        fetchDocumentFile,
        deleteDocument,
        sendRequest
    }
})