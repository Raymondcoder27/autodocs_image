<script setup lang="ts">
import { ref, type Ref } from "vue";
import { useTemplateStore } from "@/domain/templates/stores";
import { useNotificationsStore } from "@/stores/notifications";

type FileAttachment = {
  name: string;
  file: File | null;
};

const payload: Ref<FileAttachment> = ref({
  name: "",
  file: null,
});

const store = useTemplateStore();
const notify = useNotificationsStore();
const loading: Ref<boolean> = ref(false);

function submit() {
  var formData = new FormData();
  formData.set("name", payload.value.name);
  formData.set("template", payload.value.file as Blob);
  loading.value = true;
  store
    .uploadTemplate(formData)
    .then(() => {
      payload.value.file = null;
      loading.value = false;

      notify.success("Template uploaded successfully");

      setTimeout(() => {
        window.location.reload(); // Reload after notification is shown
      }, 1000);
    })
    .catch(() => {
      loading.value = false;
      notify.error("Failed to upload file");
    });
}

function onFileChanged($event: Event) {
  const target = $event.target as HTMLInputElement;
  if (target && target.files) {
    payload.value.file = target.files[0] || null;
  }
}

</script>

<template>
  <div class="flex">
    <div class="w-full">
      <form @submit.prevent="submit" class="flex flex-col items-center">
        <div class="flex flex-col py-1 w-full">
          <label class="font-semibold text-sm">Template Name</label>
          <input class="form-element" type="text" v-model="payload.name" />
        </div>
        <div class="flex flex-col py-1 w-full">
          <label class="font-semibold text-sm">Template File (html)</label>
          <input
            class="relative m-0 block w-full min-w-0 flex-auto cursor-pointer rounded border border-solid border-primary-100 text-white bg-clip-padding px-2 font-normal leading-[2.15] transition duration-300 ease-in-out file:-mx-3 file:-my-[0.32rem] file:cursor-pointer file:overflow-hidden file:rounded-none file:border-0 file:border-solid file:border-inherit file:bg-primary-500 file:px-3 file:py-[0.25rem] file:text-white file:transition file:duration-150 file:ease-in-out file:[border-inline-end-width:1px] file:[margin-inline-end:0.75rem] hover:file:bg-blue-500 focus:border-primary focus:text-blue-400 focus:shadow-text-blue-400 focus:outline-none dark:border-primary-100 dark:text-blue-400 dark:file:bg-blue-400 dark:file:text-white dark:focus:border-blue-400"
            type="file"
            accept="text/html"
            @change="onFileChanged"
          />
        </div>
        <button class="button w-1/5 mt-3" type="submit">
          Save
          <span v-if="loading" class="lds-ring">
            <div></div>
            <div></div>
            <div></div>
            <div></div>
          </span>
        </button>
      </form>
    </div>
  </div>
</template>

<style scoped>
@import "@/assets/styles/buttons.css";
@import "@/assets/styles/inputs.css";
@import "@/assets/styles/loading-ring.css";
</style>