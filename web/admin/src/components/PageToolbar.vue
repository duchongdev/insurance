<script setup lang="ts">
defineProps<{
  title: string
  subtitle?: string
}>()
</script>

<template>
  <div class="page-toolbar">
    <div class="page-toolbar__info">
      <h2 class="page-toolbar__title">{{ title }}</h2>
      <p v-if="subtitle" class="page-toolbar__subtitle">{{ subtitle }}</p>
    </div>
    <div v-if="$slots.filters" class="page-toolbar__filters">
      <slot name="filters" />
    </div>
    <div v-if="$slots.actions" class="page-toolbar__actions">
      <slot name="actions" />
    </div>
  </div>
</template>

<style scoped>
.page-toolbar {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px 16px;
  padding: 16px 20px;
  background: var(--color-bg-card);
  border-radius: var(--radius-card);
  box-shadow: var(--shadow-card);
  margin-bottom: 16px;
}

.page-toolbar__info {
  flex: 1;
  min-width: 160px;
}

.page-toolbar__title {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.page-toolbar__subtitle {
  margin: 4px 0 0;
  font-size: 13px;
  color: var(--color-text-muted);
}

.page-toolbar__filters {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
}

.page-toolbar__actions {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-left: auto;
}
</style>
