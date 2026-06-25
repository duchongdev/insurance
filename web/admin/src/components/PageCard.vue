<script setup lang="ts">
defineProps<{
  title?: string
  noPadding?: boolean
}>()
</script>

<template>
  <section class="page-card" :class="{ 'page-card--flat': noPadding }">
    <header v-if="title || $slots.header" class="page-card__header">
      <slot name="header">
        <h3 class="page-card__title">{{ title }}</h3>
      </slot>
      <div v-if="$slots.extra" class="page-card__extra">
        <slot name="extra" />
      </div>
    </header>
    <div class="page-card__body" :class="{ 'page-card__body--flat': noPadding }">
      <slot />
    </div>
  </section>
</template>

<style scoped>
.page-card {
  background: var(--color-bg-card);
  border-radius: var(--radius-card);
  box-shadow: var(--shadow-card);
  margin-bottom: 16px;
}

.page-card__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid var(--color-border);
}

.page-card__title {
  margin: 0;
  font-size: 15px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.page-card__extra {
  display: flex;
  align-items: center;
  gap: 8px;
}

.page-card__body {
  padding: 20px;
}

.page-card__body--flat {
  padding: 0;
}

.page-card--flat .page-card__header {
  border-bottom: none;
  padding-bottom: 0;
}
</style>
