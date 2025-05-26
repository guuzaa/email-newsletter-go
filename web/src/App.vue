<template>
  <div id="app">
    <router-view />
  </div>
</template>

<script setup lang="ts">
// App component handles routing
</script>

<style>
@import './styles/global.css';

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  line-height: 1.6;
  color: #333;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  min-height: 100vh;
}

#app {
  width: 100%;
  min-height: 100vh;
}
</style>
