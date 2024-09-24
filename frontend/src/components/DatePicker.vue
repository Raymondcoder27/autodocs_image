<template>
    <div>
      <label :for="id">{{ label }}</label>
      <input :id="id" type="date" v-model="dateValue" @input="updateDate" class="ring-2 ring-gray-100 bg-gray-50 rounded-sm" />
    </div>
  </template>
  
  <script setup lang="ts">
  import { defineProps, defineEmits, ref, watch } from 'vue';
  
  const props = defineProps({
    modelValue: {
      type: String,
      required: true
    },
    label: {
      type: String,
      required: true
    },
    id: {
      type: String,
      required: true
    }
  });
  
  const emits = defineEmits(['update:modelValue']);
  
  const dateValue = ref(props.modelValue);
  
  watch(dateValue, (newValue) => {
    emits('update:modelValue', newValue);
  });
  
  function updateDate(event: Event) {
    const target = event.target as HTMLInputElement;
    dateValue.value = target.value;
  }
  </script>
  
  <style scoped>
  /* Add any styles you need for the date picker */
  </style>