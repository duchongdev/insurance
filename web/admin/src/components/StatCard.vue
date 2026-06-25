<script setup lang="ts">
defineProps<{
  title: string
  value: number | string
  suffix?: string
  trend?: 'up' | 'down' | 'neutral'
  trendText?: string
}>()
</script>

<template>
  <div class="stat-card">
    <div class="stat-card__label">{{ title }}</div>
    <div class="stat-card__value-row">
      <span class="stat-card__value">{{ value }}</span>
      <span v-if="suffix" class="stat-card__suffix">{{ suffix }}</span>
    </div>
    <div v-if="trendText" class="stat-card__trend" :class="`stat-card__trend--${trend || 'neutral'}`">
      {{ trendText }}
    </div>
  </div>
</template>

<style scoped>
.stat-card {
  background: var(--color-bg-card);
  border-radius: var(--radius-card);
  box-shadow: var(--shadow-card);
  padding: 20px;
  height: 100%;
}

.stat-card__label {
  font-size: 14px;
  color: var(--color-text-secondary);
  margin-bottom: 8px;
}

.stat-card__value-row {
  display: flex;
  align-items: baseline;
  gap: 4px;
}

.stat-card__value {
  font-size: 30px;
  font-weight: 600;
  line-height: 1.2;
  color: var(--color-text-primary);
  font-variant-numeric: tabular-nums;
}

.stat-card__suffix {
  font-size: 14px;
  color: var(--color-text-secondary);
}

.stat-card__trend {
  margin-top: 8px;
  font-size: 12px;
}

.stat-card__trend--up {
  color: var(--color-success);
}

.stat-card__trend--down {
  color: var(--color-warning);
}

.stat-card__trend--neutral {
  color: var(--color-text-muted);
}
</style>
