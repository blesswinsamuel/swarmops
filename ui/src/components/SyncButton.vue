<template>
  <button @click="execute()" :class="['btn', 'ml-1', loading && 'loading']">
    <slot>Sync</slot>
  </button>
</template>

<script>
import useApi from "../useApi";

export default {
  name: "SyncButton",
  props: {
    force: Boolean,
  },
  setup(props) {
    const { execute, data, error, loading } = useApi(
      () => `/api/sync` + (props.force ? "?force=true" : ""),
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
