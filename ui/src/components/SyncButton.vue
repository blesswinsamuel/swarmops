<template>
  <button @click="execute()">Sync</button>
</template>

<script>
import useApi from "../useApi";

export default {
  name: "SyncButton",
  props: {
    force: Boolean,
  },
  setup() {
    const { execute, data, error, loading } = useApi(
      () => ({ url: `/api/sync` }),
      (r) => r.json()
    );

    return {
      execute,
      data,
      error,
      loading,
    };
  },
};
</script>
