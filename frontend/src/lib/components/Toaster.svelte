<script lang="ts">
  import Toast from "./Toast.svelte";
  import { toasts } from "../toastStore";

  function handleRemove(event: CustomEvent<{ id: string }>): void {
    const { id } = event.detail;
    toasts.update((all) => all.filter((t) => t.id !== id));
  }
</script>

<div class="toaster">
  {#each $toasts as t (t.id)}
    <Toast
      id={t.id}
      message={t.message}
      duration={t.duration}
      on:remove={handleRemove}
    />
  {/each}
</div>

<style>
  .toaster {
    position: fixed;
    top: 1rem;
    right: 1rem;
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    z-index: 9999;
  }
</style>