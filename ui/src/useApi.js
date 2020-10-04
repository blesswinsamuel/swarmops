import { ref, onMounted } from "vue";

export default function useApi(factory, handleResponse, loadOnMount = false) {
  const loading = ref(false);
  const data = ref(null);
  const error = ref(null);
  const execute = async (...args) => {
    const request = factory(...args);
    if (!request) {
      data.value = null;
      error.value = null;
      loading.value = false;
      return;
    }

    loading.value = true;
    error.value = null;
    try {
      const response = await fetch(request);
      const valueResponse = await handleResponse(response);

      data.value = valueResponse;
      return valueResponse;
    } catch (e) {
      error.value = e;
      data.value = null;
    } finally {
      loading.value = false;
    }
  };

  onMounted(() => {
    if (loadOnMount) {
      execute();
    }
  });

  return {
    loading,
    data,
    error,
    execute,
  };
}
