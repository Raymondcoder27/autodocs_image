<script setup lang="ts">
import AppModal from "@/components/AppModal.vue";
import { onMounted, computed, ref, type Ref } from "vue";
import UploadTemplate from "@/domain/templates/UploadTemplate.vue";
import { useTemplateStore } from "@/domain/templates/stores";
import { useNotificationsStore } from "@/stores/notifications";
import type { AxiosError } from "axios";
import type { ApiErrorResponse } from "@/types";
import { dateTimeFormat } from "@/composables/transformations";
import TemplateViewer from "@/components/TemplateViewer.vue";

const loading: Ref<boolean> = ref(false);
const showTemplateModal: Ref<boolean> = ref(false);
const templatePreview: Ref<boolean> = ref(false);
const showDeleteModal: Ref<boolean> = ref(false);
const selectedTemplateRef: Ref<string> = ref("");
// const selectedTemplateRef2:Ref<string> = ref("")

const store = useTemplateStore();
const notify = useNotificationsStore();

onMounted(() => {
  fetch();
});

function fetch() {
  loading.value = true;
  store
    .fetchTemplates()
    .then(() => {
      loading.value = false;
    })
    .catch((error: AxiosError<ApiErrorResponse>) => {
      loading.value = false;
      notify.error(error.response?.data.message || "Error fetching templates");
    });
}

function deleteTemplate() {
  loading.value = true;
  store
    .deleteTemplate(selectedTemplateRef.value)
    .then(() => {
      loading.value = false;
      showDeleteModal.value = false;
      fetch();
    })
    .catch((error: AxiosError<ApiErrorResponse>) => {
      loading.value = false;
      notify.error(
        error.response?.data.message || "Error deleting the template"
      );
    });
}

// Computed property to find the selected template
const selectedTemplate = computed(() => {
  return store.templates.find(
    (template) => template.refNumber === selectedTemplateRef.value
  );
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
        <button class="button" @click="showTemplateModal = true">
          <i class="fa-solid fa-plus"></i> Create Template
        </button>
      </div>
      <div class="grid grid-cols-1 gap-2 py-2 w-full">
        <span class="col-span-2">
          <table>
            <thead>
              <tr>
                <th class="header">#</th>
                <!-- <th class="header">Template Name</th>
                <th class="header">Reference</th>
                <th class="header">Date</th>
                <th class="header">Actions</th> -->
                <th class="header">TEMPLATE NAME</th>
                <th class="header">REFERENCE</th>
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
              <tr v-for="(template, idx) in store.templates" :key="idx">
                <td class="text-black-700">{{ idx + 1 }}</td>
                <td class="text-black-700">
                  <span class="font-bold">{{ template.templateName }}</span>
                </td>
                <td class="text-black-700">
                  <label class="italic">{{ template.refNumber }}</label>
                </td>
                <td class="text-black-700">{{ dateTimeFormat(template.created_at) }}</td>
                <td>
                  <div class="flex gap-2">
                    <button
                      class=""
                      @click="
                        (selectedTemplateRef = template.refNumber),
                          (templatePreview = true)
                      "
                    >
                      <i class="fa-solid fa-eye mx-1 text-xs text-gray-500 bg-gray-200 rounded-sm p-1 hover:bg-gray-300 hover:text-black-400"></i>

                    </button>
                    <button
                      class=""
                      @click="
                        showDeleteModal = true;
                        selectedTemplateRef = template.refNumber;
                      "
                    >
                      <i class="fa-solid fa-trash mx-1 text-xs text-gray-500 bg-gray-200 rounded-sm p-1 hover:bg-red-100 hover:text-red-500"></i>

                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </span>
        <!-- <TemplateViewer :ref-number="selectedTemplateRef"/> -->
      </div>
    </div>
  </div>
  <AppModal v-model="showTemplateModal" xl>
    <UploadTemplate />
  </AppModal>

  <AppModal v-model="templatePreview" xl2>
    <template #title>
    <h2 class="font-semibold text-sm">
      {{ selectedTemplate?.templateName.toUpperCase() || "Loading..." }}
    </h2>
  </template>
  <TemplateViewer :ref-number="selectedTemplateRef" />

      <!-- <div>
        <h2 class="text-center font-semibold text-sm">
      {{ selectedTemplate?.templateName.toUpperCase() || "Loading..." }}
    </h2>

    <TemplateViewer :ref-number="selectedTemplateRef" />
      </div> -->
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
        <div class="flex justify-center mx-auto w-1/2 gap-2">
          <button class="bg-blue-400 hover:bg-blue-500 rounded w-1/2 text-white" @click="showDeleteModal = false">
            <i class="fa-solid fa-times-circle mx-1"></i> Cancel
          </button>

          <button
            class="bg-danger text-white p-1 w-1/2 rounded hover:bg-red-800"
            @click="deleteTemplate"
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