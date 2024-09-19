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
    name:"dashboard",
    label:"Dashboard",
    icon:"fa-solid fa-dashboard"
  },
  {
    name:"requests",
    label:"Requests",
    icon:"fa-solid fa-tasks"
  },
  {
    name:"templates",
    label:"Templates",
    // icon:"fa-solid fa-code"
    icon:"fa-solid fa-code"
  },
  {
    name:"documents",
    label:"Documents",
    // icon:"fa-solid fa-file"
    icon:"fa-solid fa-file"
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
  <div class="flex mt-2 ml-2 mb-2 bg-white rounded-lg shadow-lg shadow-gray-600">
    <div class="w-full">
      <div :class="isRouteActive(item.name) ? 'menu-active' : 'menu'" v-for="(item, idx) in sideMenu" :key="idx" @click="navigate(item.name)">
        <i class="mx-2 my-auto" :class="item.icon"></i>
        <label class="cursor-pointer">{{item.label}}</label>
      </div>
    </div>
  </div>
</template>

<style scoped>
.menu{
  @apply w-56 ml-2 flex my-auto text-gray-600 px-2 py-3 hover:text-black-900 cursor-pointer text-sm
}

.menu-active{
  @apply w-52 ml-3 flex my-auto font-bold bg-blue-400 rounded-md text-white px-1 py-2  cursor-pointer text-sm
}
</style>