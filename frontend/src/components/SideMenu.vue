<script setup lang="ts">
import {ref, type Ref} from "vue";
import {useRoute} from "vue-router";
import router from "@/router";

const route = useRoute()

type SideMenuLink ={
  name:string
  label:string
  icon:string
}

const sideMenu:Ref<Array<SideMenuLink>> = ref([
  {
    name:"documents",
    label:"Documents",
    icon:"fa-solid fa-file"
  },
  {
    name:"templates",
    label:"Templates",
    icon:"fa-solid fa-code"
  }
])

function isRouteActive(routeName:string){
  return route.name === routeName
}

function navigate(routeName:string){
  router.push({name:routeName})
}
</script>

<template>
  <div class="flex m-2 bg-white rounded-lg">
    <div class="w-full">
      <div :class="isRouteActive(item.name) ? 'menu-active' : 'menu'" v-for="(item, idx) in sideMenu" :key="idx" @click="navigate(item.name)">
        <i class="mx-2 my-auto" :class="item.icon"></i>
        <label>{{item.label}}</label>
      </div>
    </div>
  </div>
</template>

<style scoped>
.menu{
  @apply w-56 flex my-auto text-primary px-2 py-3 hover:bg-primary-10 hover:text-primary-700
}

.menu-active{
  @apply w-56 flex my-auto font-bold bg-primary-10 text-primary-700 px-2 py-3 hover:bg-gray-10
}
</style>