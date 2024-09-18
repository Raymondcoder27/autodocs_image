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

onMounted(() => {
  fetch();
});

function fetch() {
  loading.value = true;
  store
    .fetchDocuments()
    .then(() => {
      loading.value = false;
    })
    .catch((error: AxiosError<ApiErrorResponse>) => {
      loading.value = false;
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
}

function deleteDocument() {
  loading.value = true;
  store
    .deleteDocument(selectedDocumentRef.value)
    .then(() => {
      loading.value = false;
      showDeleteModal.value = false;
      fetch();
    })
    .catch((error: AxiosError<ApiErrorResponse>) => {
      loading.value = false;
      notify.error(
        error.response?.data.message || "Error deleting the document"
      );
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
  const fileName = selectedDoc ? `${selectedDoc.refNumber}.pdf` : "document.pdf";

  a.download = fileName

  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);

  // Revoke the object URL to free up memory
  URL.revokeObjectURL(url);
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
                <!-- <th class="header">Name</th> -->
                <!-- <th class="header">Description</th>
                <th class="header">Reference</th>
                <th class="header">Template</th>
                <th class="header">Date</th>
                <th class="header">Actions</th> -->
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
              <tr v-for="(document, idx) in store.documents" :key="idx">
                <td class="text-black">{{ idx + 1 }}</td>
                <!-- <td>
                  <span
                    class="text-primary italic border border-primary-50 rounded text-sm px-2 py-[2px]"
                    >{{ document.documentName }}
                    <i class="fa-solid fa-copy mx-2"></i
                  ></span>
                </td> -->
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
                      <!-- <i class="fa-solid fa-eye mx-1 text-gray-500 hover:text-blue-400"></i> -->
                      <i class="fa-solid fa-eye mx-1 text-xs text-gray-500 bg-gray-200 rounded-sm p-1 hover:bg-gray-300 hover:text-black-400"></i>

                    </button>
                    <button
                      class=""
                      @click="
                        showDeleteModal = true;
                        selectedDocumentRef = document.refNumber;
                      "
                    >
                      <!-- <i class="fa-solid fa-trash mx-1 text-red-300 hover:text-red-500"></i> -->
                      <i class="fa-solid fa-trash mx-1 text-xs text-gray-500 bg-gray-200 rounded-sm p-1 hover:bg-red-100 hover:text-red-500"></i>

                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </span>
        <!-- <FileViewer :ref-number="selectedDocumentRef"/> -->
      </div>
    </div>
  </div>
  <AppModal v-model="showCreateRequestModal" xl>
    <CreateGenerationRequest />
  </AppModal>
  <!-- <AppModal v-model="pdfPreview" width="50%" xl2>
  <template #title>
    <div class="flex items-center justify-between">
      <h2 class="font-semibold">
        {{ selectedPdf?.description || "loading..." }}
      </h2>
      <button
        @click="downloadPdf"
        class="bg-black-900 rounded-md p-2 text-sm text-white hover:bg-blue-400"
      >
        <i class="fa-solid fa-download"></i> Download PDF
      </button>
    </div>
  </template>

  <FileViewer :ref-number="selectedDocumentRef" />
</AppModal> -->
<AppModal v-model="pdfPreview" width="50%" xl2>
  <template #title>
    <h2 class="font-semibold text-sm">
      {{ selectedPdf?.description.toUpperCase() || "loading..." }}
    </h2>
  </template>

  <div class="flex justify-between items-center">
    <!-- Keep the FileViewer component here -->
    <!-- <FileViewer :ref-number="selectedDocumentRef" /> -->
     <div />

    <!-- Place the download button in the body section -->
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
