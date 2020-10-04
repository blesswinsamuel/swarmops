<template>
  <h1>Services</h1>
  <StacksDropdown :selected="selectedStack" :onChange="setSelected" />
  <table class="table">
    <thead>
      <tr>
        <th>ID</th>
        <th>Image</th>
        <th>Mode</th>
        <th>Name</th>
        <th>Ports</th>
        <th>Replicas</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="stack in data" :key="stack.Name">
        <td>{{ stack.ID }}</td>
        <td>{{ stack.Image }}</td>
        <td>{{ stack.Mode }}</td>
        <td>{{ stack.Name }}</td>
        <td>{{ stack.Ports }}</td>
        <td>{{ stack.Replicas }}</td>
      </tr>
    </tbody>
  </table>
</template>

<script>
import { ref } from "vue";
import useSWRV from "swrv";
import StacksDropdown from "./StacksDropdown.vue";
import fetcher from "../fetcher";

export default {
  name: "Stacks",
  components: { StacksDropdown },
  setup() {
    const selectedStack = ref("");
    function setSelected(v) {
      console.log(v);
      selectedStack.value = v;
    }
    const { data, error } = useSWRV(
      () =>
        selectedStack.value
          ? `/api/docker/services?stack=${selectedStack.value}`
          : null,
      fetcher
    );

    return {
      data,
      error,
      selectedStack: selectedStack.value,
      setSelected,
    };
  },
};
</script>
