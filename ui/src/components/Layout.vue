<template>
  <div class="container grid-lg">
    <header class="navbar">
      <section class="navbar-section">
        <a href="#/" class="navbar-brand text-bold mr-2">Docker Swarm</a>
        <a href="#/stacks" class="btn btn-link">Stacks</a>
        <a href="#/services" class="btn btn-link">Services</a>
      </section>
      <section class="navbar-section">
        <SyncButton>Sync</SyncButton>
        <SyncButton :force="true">Force Sync</SyncButton>
      </section>
    </header>

    <slot></slot>
  </div>
</template>

<script>
import SyncButton from "./SyncButton.vue";

export default {
  name: "Layout",
  components: { SyncButton },
};
</script>

<style scoped>
.navbar {
  margin: 12px 0;
}
</style>
