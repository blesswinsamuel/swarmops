<template>
  <select
    class="form-control"
    :selected="selected"
    :required="true"
    @input="onChange($event.target.value)"
  >
    <option value="">Choose Stack</option>
    <option v-for="option in data" :key="option.Name" :value="option.Name">{{
      option.Name
    }}</option>
  </select>
</template>

<script>
import useSWRV from "swrv";
import fetcher from "../fetcher";

export default {
  name: "StacksDropdown",
  props: {
    selected: String,
    onChange: Function,
  },
  setup() {
    const { data, error } = useSWRV("/api/docker/stacks", fetcher);

    return {
      data,
      error,
    };
  },
};
</script>
