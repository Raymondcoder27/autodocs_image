<script setup lang="ts">
import { computed } from "vue"
import { cva } from "class-variance-authority"
import { XMarkIcon } from "@heroicons/vue/24/solid"
import { TransitionRoot, TransitionChild, Dialog, DialogPanel, DialogTitle } from "@headlessui/vue"
const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  static: Boolean,
  xs: Boolean,
  sm: Boolean,
  md: Boolean,
  lg: Boolean,
  xl: Boolean,
  xl2: Boolean,
  xl4: Boolean,
  xl6: Boolean,
})
const emit = defineEmits(["update:modelValue"])
const modalClass = computed(() => {
  return cva("w-full transform rounded-md z-[25] bg-white text-left align-middle shadow-xl transition-all", {
    variants: {
      xs: { true: "max-w-xs" },
      sm: { true: "max-w-sm" },
      md: { true: "max-w-md" },
      lg: { true: "max-w-lg" },
      xl: { true: "max-w-xl" },
      xl2: { true: "max-w-2xl" },
      xl4: { true: "max-w-4xl" },
      xl6: { true: "max-w-4xl" },
    }
  })({
    xs: props.xs,
    sm: props.sm,
    md: props.md,
    lg: props.lg,
    xl: props.xl,
    xl2: props.xl2,
    xl4: props.xl4,
    xl6: props.xl6,
  })
})
const closeModal: () => void = () => {
  emit("update:modelValue", props.static)
}
</script>
<template>
  <TransitionRoot appear :show="modelValue" as="template">
    <Dialog as="div" @close="closeModal" class="relative">
      <!-- Overlay Transition -->
      <TransitionChild
        as="template"
        enter="transition-opacity duration-300 ease-out"
        enter-from="opacity-0"
        enter-to="opacity-100"
        leave="transition-opacity duration-200 ease-in"
        leave-from="opacity-100"
        leave-to="opacity-0"
      >
        <div class="fixed inset-0 bg-black bg-opacity-25 z-[25]" />
      </TransitionChild>

      <div id="bgwrapper" class="fixed inset-0 bg-black/25 z-[23]" aria-hidden="true" />

      <div class="fixed inset-0 z-[25]">
        <div class="flex min-h-full items-center justify-center text-center">
          <TransitionChild
            as="template"
            enter="transition-transform transition-opacity duration-300 ease-out"
            enter-from="opacity-0 scale-95"
            enter-to="opacity-100 scale-100"
            leave="transition-transform transition-opacity duration-200 ease-in"
            leave-from="opacity-100 scale-100"
            leave-to="opacity-0 scale-95"
          >
            <DialogPanel :class="modalClass" class="max-h-[100vh]">
              <!-- Header Section -->
              <div class="w-full p-6 py-3 border-b border-neutral-300 flex flex-row items-center">
                <div class="flex-grow">
                  <slot name="title" />
                </div>
                <div @click="emit('update:modelValue', false)" class="cursor-pointer flex items-center">
                  <XMarkIcon class="w-5 h-5 text-neutral-900" aria-hidden="true" />
                </div>
              </div>

              <!-- Body Section -->
              <div class="p-6 pt-3 pb-5 h-auto">
                <slot />
              </div>
            </DialogPanel>
          </TransitionChild>
        </div>
      </div>
    </Dialog>
  </TransitionRoot>
</template>

<style scoped>
  #bgwrapper {
    position: fixed;
    top: 0;
    right: 0;
    left: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
  }
</style>



