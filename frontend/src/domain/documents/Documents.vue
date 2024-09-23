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
const showDeleteModal: Ref<boolean> = ref(false);
const selectedDocumentRef: Ref<string> = ref("");
const store = useDocumentStore();
const templateStore = useTemplateStore();
const notify = useNotificationsStore();
const pdfPreview: Ref<boolean> = ref(false);
const currentPage: Ref<number> = ref(1);
const itemsPerPage: number = 10;

const requestLogs: Ref<{ method: string, status: string }[]> = ref([]);

onMounted(() => {
  fetch();
});

function fetch() {
  loading.value = true;
  store
    .fetchDocuments()
    .then(() => {
      loading.value = false;
      requestLogs.value.push({ method: 'GET', status: 'SUCCESS' });
    })
    .catch((error: AxiosError<ApiErrorResponse>) => {
      loading.value = false;
      requestLogs.value.push({ method: 'GET', status: 'FAILURE' });
      notify.error(error.response?.data.message || "Error fetching documents");
    });

  templateStore
    .fetchTemplates()
    .then(() => {
      loading.value = false;
      requestLogs.value.push({ method: 'GET', status: 'SUCCESS' });
    })
    .catch((error: AxiosError<ApiErrorResponse>) => {
      loading.value = false;
      requestLogs.value.push({ method: 'GET', status: 'FAILURE' });
      notify.error(error.response?.data.message || "Error fetching templates");
    });
}

function deleteDocument() {
  loading.value = true;
  store
    .deleteDocument(selectedDocumentRef.value)
    .then(() => {
      loading.value = false;
      showDeleteModal.value = false;
      requestLogs.value.push({ method: 'DELETE', status: 'SUCCESS' });
      fetch();
    })
    .catch((error: AxiosError<ApiErrorResponse>) => {
      loading.value = false;
      requestLogs.value.push({ method: 'DELETE', status: 'FAILURE' });
      notify.error(
        error.response?.data.message || "Error deleting the document"
      );
    });
}

const selectedPdf = computed(() => {
  return store.documents.find(
    (document) => document.refNumber === selectedDocumentRef.value
  );
});

function downloadPdf() {
  const base64String = store.fileBase64;
  const byteCharacters = atob(base64String);
  const byteNumbers = new Array(byteCharacters.length);
  for (let i = 0; i < byteCharacters.length; i++) {
    byteNumbers[i] = byteCharacters.charCodeAt(i);
  }
  const byteArray = new Uint8Array(byteNumbers);
  const blob = new Blob([byteArray], { type: "application/pdf" });

  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;

  const selectedDoc = store.documents.find(
    (doc) => doc.refNumber === selectedDocumentRef.value
  );
  const fileName = selectedDoc ? `${selectedDoc.refNumber}.pdf` : "document.pdf";

  a.download = fileName;

  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);

  URL.revokeObjectURL(url);
}

// const paginatedDocuments = computed(() => {
//   const start = (currentPage.value - 1) * itemsPerPage;
//   const end = start + itemsPerPage;
//   return store.documents.slice(start, end);
// });

const paginatedDocuments = computed(() => {
  const documents = store.documents || []; // Ensure `documents` is an array
  const start = (currentPage.value - 1) * itemsPerPage;
  const end = start + itemsPerPage;
  return documents.slice(start, end);
});

function nextPage() {
  if (currentPage.value * itemsPerPage < store.documents.length) {
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
                <th class="header">DESCRIPTION</th>
                <th class="header">REFERENCE</th>
                <th class="header">TEMPLATE</th>
                <th class="header">DATE</th>
                <th class="header">ACTIONS</th>
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
              <tr v-for="(document, idx) in paginatedDocuments" :key="idx">
                <td class="text-black">{{ (currentPage - 1) * itemsPerPage + idx + 1 }}</td>
                <td class="font-bold text-black-700">{{ document.description }}</td>
                <td class="italic text-black-700">{{ document.refNumber }}</td>
                <td class="text-black-700">
                  {{ templateStore.templates?.find((t: Template) => t.id == document.templateId)?.templateName || 'Unknown Template' }}
                </td>
                <td class="text-black-700">{{ dateTimeFormat(document.created_at) }}</td>
                <td>
                  <div class="flex gap-2">
                    <button
                      class=""
                      @click="
                        (selectedDocumentRef = document.refNumber),
                        (pdfPreview = true)
                      "
                    >
                      <i class="fa-solid fa-eye mx-1 text-xs text-gray-600 bg-gray-100 border border-gray-100 rounded-sm py-0.5 px-2 hover:bg-green-50 hover:text-green-600"></i>
                    </button>
                    <button
                      class=""
                      @click="
                        showDeleteModal = true;
                        selectedDocumentRef = document.refNumber;
                      "
                    >
                      <i class="fa-solid fa-trash mx-1 text-xs text-gray-600 bg-gray-100  border border-gray-100 rounded-sm py-0.5 px-2 hover:bg-red-50 hover:text-red-500"></i>
                    </button>
                  </div>
                </td>
              </tr>

             <!-- displaying from logs table -->
              <!-- <tr v-for="(log, idx) in paginatedLogs" :key="idx">
                <td class="text-black">{{ (currentPage - 1) * itemsPerPage + idx + 1 }}</td>
                <td class="font-bold text-black-700">{{ log.description }}</td>
                <td class="italic text-black-700">{{ log.status }}</td>
              </tr> -->



            </tbody>
          </table>
        </span>
        <div class="flex justify-between mt-4">
          <button
            class="bg-gray-100 border border-gray-200 text-sm px-1 rounded-md text-gray-800 hover:bg-black-900 hover:text-white font-semibold"
            :disabled="currentPage === 1"
            @click="prevPage"
          >
          <i class="fa-solid fa-chevron-left"></i> Previous
          </button>
          <button
            class="bg-gray-100 border border-gray-200 text-sm px-1 rounded-md text-gray-800 hover:bg-black-900 hover:text-white font-semibold"
            :disabled="currentPage * itemsPerPage >= store.documents.length"
            @click="nextPage"
          >
            Next<i class="fa-solid fa-chevron-right"></i>
          </button>
        </div>
      </div>
    </div>
  </div>
  <AppModal v-model="showCreateRequestModal" xl>
    <CreateGenerationRequest />
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
          <button class="bg-blue-400 hover:bg-blue-500 w-1/2 rounded text-white" @click="showDeleteModal = false">
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