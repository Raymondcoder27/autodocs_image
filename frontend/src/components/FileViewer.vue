<script setup lang="ts">
import {onMounted, ref, type Ref, watch} from "vue";
import {useDocumentStore} from "@/domain/documents/stores";

let filePath:Ref<string> = ref("x")
const loading: Ref<boolean> = ref(false)
const store = useDocumentStore()
const props = defineProps({
  refNumber:String
})

onMounted(() =>{
  filePath.value = import.meta.env.VITE_APP_BASE_URL +"/documents/preview/" + props.refNumber
  fetch()
})

function fetch(){
  loading.value = true
  store.fetchDocumentFile(props.refNumber || "")
      .then(() =>{ loading.value = false})
      .catch(() =>{ loading.value = false})
}

watch(() => props.refNumber,
    () => {
      filePath.value = import.meta.env.VITE_APP_BASE_URL + "/documents/preview/" + props.refNumber
      fetch()
    })

</script>

<template>
  <div class="flex">
    <div class="w-full">
      <div class="grid grid-cols-1 shadow" v-if="props.refNumber">
        <div class="relative w-full aspect-[1/1.1] overflow-hidden border border-gray-50 rounded-lg bg-white" >
          <iframe v-if="store.fileBase64"
                  allowtransparency="true"
                  loading="lazy"
                  class="pdf"
                  :src="'data:application/pdf;base64,'+store.fileBase64+ '#toolbar=0&navpanes=0&scrollbar=0&transparent=0'"></iframe>
        </div>
      </div>
      <div class="flex shadow" v-else style="height: 80vh">
        <div class="bg-gray-10 p-2 m-2 rounded flex w-full text-gray-500 justify-center items-center text-center">
          Previewer
        </div>
      </div>
    </div>
  </div>

  
</template>

<style scoped>
@import "@/assets/styles/loading-ring.css";
@import "@/assets/styles/buttons.css";
@import "@/assets/styles/table.css";

.pdf {
  @apply w-full absolute h-full object-center object-none rounded-md;
}


</style>