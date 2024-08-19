<script setup lang="ts">

import AppModal from "@/components/AppModal.vue";
import {onMounted, ref, type Ref} from "vue";
import UploadTemplate from "@/domain/templates/UploadTemplate.vue";
import {useTemplateStore} from "@/domain/templates/stores";
import {useNotificationsStore} from "@/stores/notifications";
import type {AxiosError} from "axios";
import type {ApiErrorResponse} from "@/types";
import {dateTimeFormat} from "@/composables/transformations";
import TemplateViewer from "@/components/TemplateViewer.vue";

const loading:Ref<boolean> = ref(false)
const showTemplateModal:Ref<boolean> = ref(false)
const showDeleteModal:Ref<boolean> = ref(false)
const selectedTemplateRef:Ref<string> = ref("")
const store = useTemplateStore()
const notify = useNotificationsStore()

onMounted(() =>{
 fetch()
})

function fetch(){
  loading.value = true
  store.fetchTemplates()
      .then(() =>{loading.value = false})
      .catch((error:AxiosError<ApiErrorResponse>) =>{
        loading.value = false
        notify.error(error.response?.data.message || "Error fetching templates")
      })
}

function deleteTemplate(){
  loading.value = true
  store.deleteTemplate(selectedTemplateRef.value)
      .then(() =>{
        loading.value = false
        showDeleteModal.value = false
        fetch()
      })
      .catch((error:AxiosError<ApiErrorResponse>) =>{
        loading.value = false
        notify.error(error.response?.data.message || "Error deleting the template")
      })
}
</script>

<template>
  <div class="flex p-2 bg-white">
    <div class="w-full">
      <div class="flex justify-between">
        <span>
          <i class="cursor-pointer bg-primary-10 text-primary rounded p-2 fa-solid fa-refresh my-auto" @click="fetch"></i>
        </span>
        <button class="button" @click="showTemplateModal = true"><i class="fa-solid fa-plus"></i> Create Template</button>
      </div>
      <div class="grid grid-cols-3 gap-2 py-2">
        <span class="col-span-2">
          <table>
          <thead>
          <tr>
            <th class="header">#</th>
            <th class="header">Template Name</th>
            <th class="header">Reference</th>
            <th class="header">Date</th>
            <th class="header">Actions</th>
          </tr>
          </thead>
          <thead v-if="loading">
          <tr>
            <th colspan="12" style="padding: 0">
              <div class="w-full bg-primary-300 h-1 p-0 m-0 animate-pulse"></div>
            </th>
          </tr>
          </thead>
          <tbody>
          <tr v-for="(template, idx) in store.templates" :key="idx">
            <td>{{idx + 1}}</td>
            <td>
              <span class="font-bold">{{template.templateName}}</span>
            </td>
            <td>
              <label class="italic">{{template.refNumber}}</label>
            </td>
            <td>{{dateTimeFormat(template.created_at)}}</td>
            <td>
              <div class="flex gap-2">
                <button class="action-btn" @click="selectedTemplateRef = template.refNumber"><i class="fa-solid fa-eye mx-1"></i></button>
                <button class="danger-action-btn"  @click="showDeleteModal = true; selectedTemplateRef = template.refNumber;"><i class="fa-solid fa-trash mx-1"></i></button>
              </div>
            </td>
          </tr>
          </tbody>
        </table>
        </span>
        <TemplateViewer :ref-number="selectedTemplateRef"/>
      </div>

    </div>
  </div>
  <AppModal v-model="showTemplateModal" xl>
    <UploadTemplate/>
  </AppModal>

  <AppModal v-model="showDeleteModal" xl>
    <div class="flex">
      <div class="w-full">
        <div class="flex">
          <span class="mx-auto text-center justify-center">
            <i class="mx-auto fa-solid fa-exclamation-circle text-3xl text-warning-600"></i>
          </span>
        </div>
        <p class="py-5">
          Are you sure you want to delete this template and lose it completely?
        </p>
        <div class="grid grid-cols-2 gap-2">
          <button class="button" @click="showDeleteModal = false">
            <i class="fa-solid fa-times-circle mx-1"></i> Cancel
          </button>

          <button class="bg-danger text-white p-2 rounded" @click="deleteTemplate">
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