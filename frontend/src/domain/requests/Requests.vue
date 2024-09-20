<script setup lang="ts">
import AppModal from "@/components/AppModal.vue";
import CreateGenerationRequest from "@/domain/documents/CreateGenerationRequest.vue";
import { onMounted, computed, type Ref, ref } from "vue";
import { useDocumentStore } from "@/domain/documents/stores";
import type { AxiosError } from "axios";
import type { ApiErrorResponse } from "@/types";
import { useNotificationsStore } from "@/stores/notifications";
import FileViewer from "@/components/FileViewer.vue";
import { dateTimeFormat } from "../../composables/transformations";
import { useTemplateStore } from "@/domain/templates/stores";
import { Template } from "../templates/types";

const loading: Ref<boolean> = ref(false);
const showCreateRequestModal: Ref<boolean> = ref(false);
const pdfPreview: Ref<boolean> = ref(false);
const jsonPayloadPreview: Ref<boolean> = ref(false);
const currentPage: Ref<number> = ref(1);
const itemsPerPage: number = 10;

const requestLogs: Ref<{ method: string, status: string, refNumber: string, templateId: string, description: string, jsonPayload: string }[]> = ref([]);

onMounted(() => {
    fetch();
});

function fetch() {
    loading.value = true;
    // Simulate fetching documents and templates
    setTimeout(() => {
        loading.value = false;
        requestLogs.value.push({ method: 'GET', status: 'SUCCESS', refNumber: '123', templateId: 'template1', description: 'Sample Document', jsonPayload: '{}' });
    }, 1000);
}

function createDocument(payload) {
    loading.value = true;
    // Simulate creating a document
    setTimeout(() => {
        loading.value = false;
        requestLogs.value.push({ method: 'POST', status: 'SUCCESS', refNumber: '124', templateId: 'template2', description: 'New Document', jsonPayload: JSON.stringify(payload) });
        fetch();
    }, 1000);
}

const paginatedRequests = computed(() => {
    const start = (currentPage.value - 1) * itemsPerPage;
    const end = start + itemsPerPage;
    return requestLogs.value.slice(start, end);
});

function nextPage() {
    if (currentPage.value * itemsPerPage < requestLogs.value.length) {
        currentPage.value++;
    }
}

function prevPage() {
    if (currentPage.value > 1) {
        currentPage.value--;
    }
}

const failureRate = computed(() => {
    const totalRequests = requestLogs.value.length;
    const failedRequests = requestLogs.value.filter(log => log.status === 'FAILURE').length;
    return totalRequests > 0 ? (failedRequests / totalRequests) * 100 : 0;
});
</script>

<template>
    <div class="flex p-2 bg-white shadow-md shadow-black-200 rounded-xl">
        <div class="w-full">
            <div class="flex justify-between">
                <span>
                    <i
                        class="cursor-pointer bg-primary-10 text-blue-400 rounded p-2 fa-solid fa-refresh my-auto"
                        @click="fetch"
                    ></i>
                </span>
                <button class="button" @click="showCreateRequestModal = true">
                    <i class="fa-solid fa-plus"></i> Create Request
                </button>
            </div>

            <div class="grid grid-cols-1 gap-2 py-2">
                <span class="col-span-2">
                    <table class="">
                        <thead class="text-xs">
                            <tr>
                                <th class="header">#</th>
                                <th class="header">Ref Number</th>
                                <th class="header">Template</th>
                                <th class="header">Document</th>
                                <th class="header">Method</th>
                                <th class="header">Status</th>
                                <th class="header">Json Payload</th>
                                <th class="header">Actions</th>
                            </tr>
                        </thead>
                        <thead v-if="loading">
                            <tr>
                                <th colspan="12" style="padding: 0">
                                    <div
                                        class="w-full bg-primary-300 h-1 p-0 m-0 animate-pulse"
                                    ></div>
                                </th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="(log, idx) in paginatedRequests" :key="idx">
                                <td class="text-black">{{ (currentPage - 1)* itemsPerPage + idx + 1 }}</td>
                                <td class="italic text-black-700">{{ log.refNumber }}</td>
                                <td class="text-black-700">
                                    {{ templateStore.templates?.find((t: Template) => t.id == log.templateId)?.templateName || 'Unknown Template' }}
                                </td>
                                <td class="text-black-700">
                                    {{ log.description || 'Unknown Document' }}
                                </td>
                                <td class="text-black-700">
                                    <span
                                    :class="{
                                        'bg-warning-100 border border-warning-400 text-warning-600 font-semibold rounded-sm px-1 py-0.3': log.method === 'POST',
                                        'bg-blue-50 border border-blue-300 text-blue-400 font-semibold rounded-sm px-1 py-0.3': log.method === 'GET',
                                        'bg-red-100 border border-red-500 text-red-600 font-semibold rounded-sm px-1 py-0.3': log.method === 'DELETE',
                                    }">{{ log.method }}</span>
                                </td>
                                <td class="text-black-700"><span class="bg-green-100 text-xs border border-green-300 text-green-500 font-semibold rounded-sm px-1 py-0.3">SUCCESS</span>{{ log.status }}</td>
                                <td class="text-black-700"><button
                                    @click="(selectedDocumentRef = log.refNumber), (jsonPayloadPreview = true)"
                                    class="bg-gray-50 border border-gray-200 text-gray-500 hover:bg-gray-200 hover:text-gray-600 font-semibold rounded-sm px-1 py-0.3 text-xs">
                                    PREVIEW</button>
                                </td>
                                <td>
                                    <div class="flex gap-2">
                                        <button
                                            class=""
                                            @click="
                                                (selectedDocumentRef = log.refNumber),
                                                    (pdfPreview = true)
                                            "
                                        >
                                            <i class="fa-solid fa-eye mx-1 text-xs text-gray-600 bg-gray-100 border border-gray-100 rounded-sm py-0.5 px-2 hover:bg-green-50 hover:text-green-600"></i>
                                        </button>
                                    </div>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </span>

                <div class="flex justify-between mt-4">
          <button
            class="bg-gray-100 border border-gray-200 text-sm px-2 rounded-md text-gray-800 hover:bg-black-900 hover:text-white font-semibold"
            :disabled="currentPage === 1"
            @click="prevPage"
          >
          <i class="fa-solid fa-chevron-left"></i> Previous
          </button>
          <button
            class="bg-gray-100 border border-gray-200 text-sm px-2 rounded-md text-gray-800 hover:bg-black-900 hover:text-white font-semibold"
            :disabled="currentPage * itemsPerPage >= requestLogs.length"
            @click="nextPage"
          >
            Next<i class="fa-solid fa-chevron-right"></i>
          </button>
        </div>
            </div>
        </div>
        
    </div>
    <AppModal v-model="showCreateRequestModal" xl>
        <CreateGenerationRequest @submit="createDocument" />
    </AppModal>
    <AppModal v-model="jsonPayloadPreview" class="flex flex-col py-2" xl>
        <template #title>
            <h2 class="font-semibold text-sm">
                JSON PAYLOAD
            </h2>
        </template>
        <pre class="text-wrap bg-gray-100 text-[10px] p-2 h-auto overflow-auto max-h-[500px] flex-grow">
            {{ requestLogs.find((log) => log.refNumber === selectedDocumentRef)?.jsonPayload 
            ? JSON.stringify(JSON.parse(requestLogs.find((log) => log.refNumber === selectedDocumentRef)?.jsonPayload), null, 2) 
            : 'No JSON Payload available' }}
        </pre>
    </AppModal>
    <AppModal v-model="pdfPreview" width="50%" xl2>
        <template #title>
            <h2 class="font-semibold text-sm">
                {{ selectedPdf?.description.toUpperCase() || "loading..." }}
            </h2>
        </template>

        <div class="flex justify-between items-center">
            <div />
            <button
                @click="downloadPdf"
                class="bg-black-900 rounded-md p-1.5 mb-2 text-sm text-white hover:bg-blue-400"
            >
                <i class="fa-solid fa-download"></i> Download PDF
            </button>
        </div>
        <FileViewer :ref-number="selectedDocumentRef" />
    </AppModal>
</template>

<style scoped>
@import "@/assets/styles/buttons.css";
@import "@/assets/styles/table.css";
</style>