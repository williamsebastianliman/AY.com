<script lang="ts">
  import { createEventDispatcher, onMount } from "svelte";

  export let id: string;
  export let message: string;
  export let duration: number = 3000;

  const dispatch = createEventDispatcher();
  let visible = false;

  onMount(() => {
    visible = true;
    const hideTimer = setTimeout(() => {
      visible = false;
      setTimeout(() => dispatch("remove", { id }), 300);
    }, duration);

    return () => clearTimeout(hideTimer);
  });
</script>

<div class="toast {visible ? "enter" : "leave"}">
  {message}
</div>

<style>
  .toast {
    background: rgba(0, 0, 0, 0.85);
    color: white;
    padding: 0.75rem 1rem;
    margin-top: 0.5rem;
    border-radius: 0.25rem;
    max-width: 300px;
    transform: translateX(100%);
    opacity: 0;
    transition: transform 0.3s ease, opacity 0.3s ease;
  }
  .toast.enter {
    transform: translateX(0);
    opacity: 1;
  }
  .toast.leave {
    transform: translateX(100%);
    opacity: 0;
  }
</style>