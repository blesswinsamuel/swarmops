<template>
  <Layout>
    <Routes />
  </Layout>
</template>

<script>
import Layout from "./components/Layout.vue";
import Routes from "./Routes.vue";
import "spectre.css/dist/spectre.css";

export default {
  name: "App",
  components: { Layout, Routes },
};
</script>
