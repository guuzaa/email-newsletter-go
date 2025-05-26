<template>
  <div class="feature">
    <div class="feature-icon">{{ icon }}</div>
    <h3>{{ title }}</h3>
    <p>{{ description }}</p>
  </div>
</template>

<script setup lang="ts">
interface Props {
  icon: string
  title: string
  description: string
}

defineProps<Props>()
</script>

<style scoped>
.feature {
  background: rgba(255, 255, 255, 0.1);
  padding: 30px;
  border-radius: 15px;
  text-align: center;
  backdrop-filter: blur(10px);
}

.feature-icon {
  font-size: 3rem;
  margin-bottom: 15px;
}

.feature h3 {
  margin-bottom: 10px;
  font-size: 1.3rem;
}

.feature p {
  opacity: 0.8;
}
</style>
