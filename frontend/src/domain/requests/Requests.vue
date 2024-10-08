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
import { Doc } from "@/domain/documents/types";
import { Log } from "@/domain/requests/types";
import { useLogStore } from "@/domain/requests/stores";

const loading: Ref<boolean> = ref(false);
const showCreateRequestModal: Ref<boolean> = ref(false);
const showDeleteModal: Ref<boolean> = ref(false);
const showDeleteLogsRequest: Ref<boolean> = ref(false);
const selectedDocumentRef: Ref<string> = ref("");
const store = useDocumentStore();
const logStore = useLogStore();
// const logStore = useLogStore();
logStore.logs = logStore.logs || [];

const templateStore = useTemplateStore();
const documentStore = useDocumentStore();
const notify = useNotificationsStore();
const pdfPreview: Ref<boolean> = ref(false);
const jsonPayloadPreview: Ref<boolean> = ref(false);
const currentPage: Ref<number> = ref(1);
const itemsPerPage: number = 10;

onMounted(() => {
  fetch();
});

function fetch() {
  loading.value = true;
  store
    .fetchDocuments()
    .then(() => {
      loading.value = false;
      // requestLogs.value.push({ method: 'GET', status: 'SUCCESS' });
    })
    .catch((error: AxiosError<ApiErrorResponse>) => {
      loading.value = false;
      // requestLogs.value.push({ method: 'GET', status: 'FAILURE' });
      notify.error(error.response?.data.message || "Error fetching documents");
    });

  templateStore
    .fetchTemplates()
    .then(() => {
      loading.value = false;
    })
    .catch((error: AxiosError<ApiErrorResponse>) => {
      loading.value = false;
      notify.error(error.response?.data.message || "Error fetching templates");
    });

logStore
  .fetchLogs()
  .then(() => {
    loading.value = false;
    if (!logStore.logs) logStore.logs = []; // Fallback to empty array if undefined
  })
  .catch((error: AxiosError<ApiErrorResponse>) => {
    loading.value = false;
    notify.error(error.response?.data.message || "Error fetching logs");
    logStore.logs = []; // Fallback to empty array on error
  });

}


//clear logs
function clearLogs() {
  loading.value = true;
  logStore
    .clearLogs()
    .then(() => {
      loading.value = false;
      showDeleteLogsRequest.value = false;
      fetch();
    })
    .catch((error: AxiosError<ApiErrorResponse>) => {
      loading.value = false;
      notify.error(error.response?.data.message || "Error deleting the logs");
    });
}

//computed property to find selected pdf:
const selectedPdf = computed(() => {
  return store.documents.find(
    (document) => document.refNumber === selectedDocumentRef.value
  );
});

function downloadPdf() {
  // Convert Base64 to a Blob
  const base64String = store.fileBase64;
  const byteCharacters = atob(base64String);
  const byteNumbers = new Array(byteCharacters.length);
  for (let i = 0; i < byteCharacters.length; i++) {
    byteNumbers[i] = byteCharacters.charCodeAt(i);
  }
  const byteArray = new Uint8Array(byteNumbers);
  const blob = new Blob([byteArray], { type: "application/pdf" });

  // Create a download link
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;

  // Find the selected document using selectedDocumentRef
  const selectedDoc = store.documents.find(
    (doc) => doc.refNumber === selectedDocumentRef.value
  );
  const fileName = selectedDoc
    ? `${selectedDoc.refNumber}.pdf`
    : "document.pdf";

  a.download = fileName;

  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);

  // Revoke the object URL to free up memory
  URL.revokeObjectURL(url);
}

const paginatedLogs = computed(() => {
  if (!logStore.logs || !logStore.logs.length) {
    return [];
  }
  const start = (currentPage.value - 1) * itemsPerPage;
  const end = start + itemsPerPage;
  return logStore.logs.slice(start, end);
});

function nextPage() {
  if (
    logStore.logs &&
    currentPage.value * itemsPerPage < logStore.logs.length
  ) {
    currentPage.value++;
  }
}

function prevPage() {
  if (currentPage.value > 1) {
    currentPage.value--;
  }
}
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
        <button
          class="bg-red-500 text-white text-xs font-semibold rounded-md px-1 py-0"
          @click="showDeleteLogsRequest = true"
        >
          <i class="fa-solid fa-trash"></i> Clear logs
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
            <tbody v-if="paginatedLogs && paginatedLogs.length">
              <!-- <tr v-for="(document, idx) in paginatedRequests" :key="idx"> -->
              <tr v-for="(log, idx) in paginatedLogs" :key="idx">
                <td class="text-black">
                  {{ (currentPage - 1) * itemsPerPage + idx + 1 }}
                </td>
                <!-- <td class="italic text-black-700">{{ document.refNumber }}</td> -->
                <td class="italic text-black-700">
                  {{ log.refNumber || "--" }}
                </td>

                <td class="text-black-700">
                  <!-- {{ logStore.logs?.find((l: Log) => l.templateId == log.templateId)?.templateName || '--' }} -->
                  {{ templateStore.templates?.find((t: Template) => t.id == log.templateId)?.templateName || '--' }}
                </td>
                <!-- <td class="text-black-700">
                  {{ logStore.logs?.find((l: Log) => l.refNumber == log.refNumber)?.templateName || 'All Templates' }}
                </td> -->
                <td class="text-black-700">
                  <!-- {{ logStore.logs?.find((l: Log) => l.refNumber == log.refNumber)?.description || 'All Documents' }} -->
                  {{ logStore.logs?.find((l: Log) => l.refNumber == log.refNumber)?.description || '--' }}
                </td>
                <td class="text-black-700">
                  <span
                    :class="{
                      'bg-warning-100 border border-warning-400 text-warning-600 font-semibold rounded-sm px-1 py-0.3':
                        log.requestMethod === 'POST',
                      'bg-blue-50 border border-blue-300 text-blue-400 font-semibold rounded-sm px-1 py-0.3':
                        log.requestMethod === 'GET',
                      'bg-red-100 border border-red-300 text-red-500 font-semibold rounded-sm px-1 py-0.3':
                        log.requestMethod === 'DELETE',
                    }"
                    >{{ log.requestMethod }}</span
                  >
                </td>
                <!-- <td class="text-black-700"><span class="bg-green-100 text-xs border border-green-300 text-green-500 font-semibold rounded-sm px-1 py-0.3">SUCCESS</span>{{ document.status }}</td> -->
                <!-- <td class="text-black-700">
                  <span
                    class="bg-green-100 text-xs border border-green-300 text-green-500 font-semibold rounded-sm px-1 py-0.3"
                    >SUCCESS</span
                  >{{ log.status }}
                </td> -->
                <!-- Making the status dynamic -->
                <td class="text-black">
                    <span
                    :class="{
                        'bg-green-100 text-xs border border-green-300 text-green-500 font-semibold rounded-sm px-1 py-0.3':
                        log.requestStatus === 'SUCCESS',
                        'bg-red-100 text-xs border border-red-300 text-red-500 font-semibold rounded-sm px-1 py-0.3':
                        log.requestStatus === 'FAILED',
                    }"  
                >
                    {{ log.requestStatus }}
                </span>
                </td>

                <td class="text-black-700">
                  <button
                    @click="
                      (selectedDocumentRef = log.refNumber),
                        (jsonPayloadPreview = true)
                    "
                    class="bg-gray-50 border border-gray-200 text-gray-500 hover:bg-gray-200 hover:text-gray-600 font-semibold rounded-sm px-1 py-0.3 text-xs"
                  >
                    PREVIEW
                  </button>
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
                      <i
                        class="fa-solid fa-eye mx-1 text-xs text-gray-600 bg-gray-100 border border-gray-100 rounded-sm py-0.5 px-2 hover:bg-green-50 hover:text-green-600"
                      ></i>
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </span>
        <div
          class="flex justify-between mt-4"
          v-if="logStore.logs.length > itemsPerPage"
        >
          <button
            :disabled="currentPage === 1"
            @click="prevPage"
            class="bg-gray-100 border border-gray-200 text-sm px-1 rounded-md text-gray-800 hover:bg-black-900 hover:text-white font-semibold"
          >
            <i class="fa-solid fa-chevron-left"></i> Previous
          </button>
          <button
            :disabled="currentPage * itemsPerPage >= logStore.logs.length"
            @click="nextPage"
            class="bg-gray-100 border border-gray-200 text-sm px-1 rounded-md text-gray-800 hover:bg-black-900 hover:text-white font-semibold"
          >
            Next<i class="fa-solid fa-chevron-right"></i>
          </button>
        </div>
      </div>
    </div>
  </div>
  <AppModal v-model="showDeleteLogsRequest" xl>
    <!-- <CreateGenerationRequest @submit="createDocument" /> -->
    <div class="flex">
      <div class="w-full">
        <div class="flex">
          <span class="mx-auto text-center justify-center">
            <i
              class="mx-auto fa-solid fa-exclamation-circle text-3xl text-danger"
            ></i>
          </span>
        </div>
        <p class="py-5">
          Are you sure you want to delete the logs and lose them completely?
        </p>
        <div class="flex w-1/2 gap-2 justify-center mx-auto">
          <button
            class="bg-blue-400 hover:bg-blue-500 w-1/2 rounded text-white"
            @click="showDeleteModal = false"
          >
            <i class="fa-solid fa-times-circle mx-1"></i> Cancel
          </button>

          <button
            class="bg-danger text-white p-1 w-1/2 rounded hover:bg-red-800"
            @click="clearLogs"
          >
            <i class="fa-solid fa-check-circle mx-1"></i> Confirm
          </button>
        </div>
      </div>
    </div>
  </AppModal>
  <AppModal v-model="jsonPayloadPreview" class="flex flex-col py-2" xl>
    <template #title>
      <h2 class="font-semibold text-sm">JSON PAYLOAD</h2>
    </template>
    <pre
      class="text-wrap bg-gray-100 text-[10px] p-2 h-auto overflow-auto max-h-[500px] flex-grow"
    >
            {{
        store.documents.find((doc) => doc.refNumber === selectedDocumentRef)
          ?.jsonPayload
          ? JSON.stringify(
              JSON.parse(
                store.documents.find(
                  (doc) => doc.refNumber === selectedDocumentRef
                )?.jsonPayload
              ),
              null,
              2
            )
          : "No JSON Payload available"
      }}
        </pre
    >
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

  <AppModal v-model="showDeleteModal" xl>
    <div class="flex">
      <div class="w-full">
        <div class="flex">
          <span class="mx-auto text-center justify-center">
            <i
              class="mx-auto fa-solid fa-exclamation-circle text-3xl text-danger"
            ></i>
          </span>
        </div>
        <p class="py-5">
          Are you sure you want to delete this template and lose it completely?
        </p>
        <div class="flex w-1/2 gap-2 justify-center mx-auto">
          <button
            class="bg-blue-400 hover:bg-blue-500 w-1/2 rounded text-white"
            @click="showDeleteModal = false"
          >
            <i class="fa-solid fa-times-circle mx-1"></i> Cancel
          </button>

          <button
            class="bg-danger text-white p-1 w-1/2 rounded hover:bg-red-800"
            @click="deleteDocument"
          >
            <i class="fa-solid fa-check-circle mx-1"></i> Confirm
          </button>
        </div>
      </div>
    </div>
  </AppModal>
</template>

<style scoped>
@import "@/assets/styles/buttons.css";
@import "@/assets/styles/table.css";
</style>