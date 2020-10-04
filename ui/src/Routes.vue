<script>
import { h } from "vue";
import Stacks from "./components/Stacks.vue";
import Services from "./components/Services.vue";
import "spectre.css/dist/spectre.css";

const NotFound = { template: "<p>Page not found</p>" };

const routes = {
  "/": Stacks,
  "/stacks": Stacks,
  "/services": Services,
};

export default {
  name: "Routes",
  data: () => ({
    currentRoute: window.location.pathname,
  }),
  mounted: function() {
    this.onRouteChange();
    window.addEventListener("hashchange", this.onRouteChange);
  },
  beforeUnmount: function() {
    window.removeEventListener("hashchange", this.onRouteChange);
  },
  methods: {
    onRouteChange: function() {
      this.currentRoute = window.location.hash.substring(1);
    },
  },
  computed: {
    ViewComponent() {
      return routes[this.currentRoute] || NotFound;
    },
  },
  render() {
    return h(this.ViewComponent);
  },
};
</script>
