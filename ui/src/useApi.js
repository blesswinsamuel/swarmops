import { ref } from "vue";

export default function useApi(factory, handleResponse) {
  const loading = ref(false);
  const result = ref(null);
  const error = ref(null);
  const execute = async (...args) => {
    const request = factory(...args);

    loading.value = true;
    error.value = null;
    try {
      const response = await fetch(request);
      const valueResponse = await handleResponse(response);

      result.value = valueResponse;
      return valueResponse;
    } catch (e) {
      error.value = e;
      result.value = null;
    } finally {
      loading.value = false;
    }
  };

  return {
    loading,
    result,
    error,
    execute,
  };
}
