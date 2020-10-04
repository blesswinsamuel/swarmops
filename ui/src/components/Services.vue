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
import StacksDropdown from "./StacksDropdown.vue";
import useApi from "../useApi";

export default {
  name: "Stacks",
  components: { StacksDropdown },
  setup() {
    const selectedStack = ref("");
    const { execute, data, error } = useApi(
      () =>
        selectedStack.value
          ? `/api/docker/services?stack=${selectedStack.value}`
          : null,
      (r) => r.json(),
      true
    );
    function setSelected(v) {
      selectedStack.value = v;
      execute();
    }

    return {
      data,
      error,
      selectedStack: selectedStack.value,
      setSelected,
    };
  },
};
</script>
