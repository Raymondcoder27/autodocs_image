<script setup lang="ts">

import {onMounted, ref, type Ref} from "vue";
import type {GenerationRequest} from "@/domain/documents/types";
import {useDocumentStore} from "@/domain/documents/stores";
import type {AxiosError} from "axios";
import type {ApiErrorResponse} from "@/types";
import {useNotificationsStore} from "@/stores/notifications";
import {useTemplateStore} from "@/domain/templates/stores";

const loading:Ref<boolean> = ref(false)
const json:Ref<string> = ref("")
const description:Ref<string> = ref("")
const notify = useNotificationsStore()

const form:Ref<GenerationRequest> = ref(
    {
      refNumber:"",
      description:"",
      data:""
    }
)

const store = useDocumentStore()
const templateStore = useTemplateStore()

onMounted(() =>{
  fetch()
})

function submit(){
  loading.value = true
  form.value.data = JSON.parse(json.value)
  form.value.description = description.value
  store.sendRequest(form.value)
      .then(() =>{
        loading.value = false
        window.location.reload()
      })
      .catch(() =>{loading.value = false})
}

function fetch(){
  loading.value = true
  templateStore.fetchTemplates()
      .then(() =>{loading.value = false})
      .catch((error:AxiosError<ApiErrorResponse>) =>{
        loading.value = false
        notify.error(error.response?.data.message || "Error fetching templates")
      })
}

</script>

<template>
  <div class="flex">
    <div class="w-full">
      <form @submit.prevent="submit">
        <div class="grid grid-cols-1 gap-2">
          <div class="flex flex-col py-1">
            <label class="font-bold">Reference Number</label>
            <select class="form-element" v-model="form.refNumber">
              <option :value="template.refNumber" v-for="(template, idx) in templateStore.templates" :key="idx">{{template.templateName}}</option>
            </select>
          </div>
          <div class="flex flex-col py-1">
            <label class="font-bold">Document Description</label>
            <input type="text" v-model="description">
          </div>
          <div class="flex flex-col py-1">
            <label class="font-bold">JSON Data</label>
            <textarea class="form-element" v-model="json" rows="4"></textarea>
          </div>
        </div>

        <div class="flex flex-col py-2">
          <label class="font-bold">Data Preview</label>
          <pre class="text-wrap bg-gray-10 text-[10px] p-2 h-full flex-grow">{{json}}</pre>
        </div>

        <div class="flex justify-between py-2">
          <span></span>
          <button class="button" type="submit">
            Send Request
            <span v-if="loading" class="lds-ring">
            <div></div>
            <div></div>
            <div></div>
            <div></div>
          </span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<style scoped>
@import "@/assets/styles/buttons.css";
@import "@/assets/styles/inputs.css";
@import "@/assets/styles/loading-ring.css";
</style>