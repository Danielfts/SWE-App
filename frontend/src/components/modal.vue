<script setup lang="ts">
const props = defineProps({
  modelValue: { type: Boolean, required: true },
  stock: {type: Object, required: true}
})

const emit = defineEmits(['update:modelValue'])

const close = () => emit('update:modelValue', false)
</script>

<template>
    <div v-if="modelValue" @click.self="close" class="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50 p-4">
        <div class="bg-white rounded-lg shadow-2xl max-w-2xl w-full max-h-[90vh] overflow-y-auto p-6">
            <h1 style="font-size: xx-large; font-weight: bold;">Recommended Stock of the day</h1>
            <br>
            <div v-for="k in Object.keys(stock)" style="display: grid; grid-template-columns: 50% 50%;">
                <strong>{{ k }}</strong>
                <span>{{ stock[k] }}</span>
            </div>
           
            <button @click="close" class="mt-6 px-6 py-2 bg-[#3B1CEA] text-white font-semibold rounded-lg hover:bg-[#2D15B8] transition-colors shadow-md">Close</button>
        </div>
    </div>
</template>