<template>
  <h1>Stacks</h1>
  <table class="table">
    <thead>
      <tr>
        <th>Name</th>
        <th>Namespace</th>
        <th>Orchestrator</th>
        <th>Services</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="stack in data" :key="stack.Name">
        <td>{{ stack.Name }}</td>
        <td>{{ stack.Namespace }}</td>
        <td>{{ stack.Orchestrator }}</td>
        <td>{{ stack.Services }}</td>
      </tr>
    </tbody>
  </table>
</template>

<script>
import useApi from "../useApi";

export default {
  name: "Stacks",
  setup() {
    const { data, error } = useApi(
      () => "/api/docker/stacks",
      (r) => r.json(),
      true
    );

    return {
      data,
      error,
    };
  },
};
</script>
