<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  status: 'success' | 'warning' | 'danger' | 'default'
  label: string
}>()

const tagType = computed(() => {
  const map = {
    success: 'success',
    warning: 'warning',
    danger: 'danger',
    default: 'info',
  } as const
  return map[props.status]
})
</script>

<template>
  <el-tag :type="tagType" size="small" effect="light">{{ label }}</el-tag>
</template>
