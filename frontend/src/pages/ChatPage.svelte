<script lang="ts">
  import api from "../lib/api";
  import { tick } from "svelte";

  let text = "";
  let result: number | null = null;
  let label = "";
  let loading = false;
  let errorMsg = "";

  async function predict() {
    if (!text.trim()) return;
    loading = true;
    errorMsg = "";
    result = null;
    try {
      const res = await api.post<{ prediction: number }>(
        "/ml/predict",
        { text },
        { withCredentials: true }
      );
      result = res.data.prediction;
      await tick();
    } catch (err) {
      console.error(err);
      errorMsg = "Prediction failed";
    } finally {
      loading = false;
    }
  }

  $: label = result === 1
    ? "Good"
    : result === -1
    ? "Bad"
    : "";
</script>

<style>
  .container {
    max-width: 500px;
    margin: 2rem auto;
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }
  textarea {
    width: 100%;
    height: 6rem;
    padding: 0.5rem;
    font-size: 1rem;
    resize: vertical;
  }
  button {
    align-self: flex-end;
    padding: 0.5rem 1rem;
    font-size: 1rem;
  }
  .result {
    font-weight: bold;
    font-size: 1.25rem;
  }
  .loading {
    opacity: 0.7;
  }
  .error {
    color: #e00;
  }
</style>

<div class="container">
  <textarea
    bind:value={text}
    placeholder="Type your text here…"
    disabled={loading}
  />

  <button on:click={predict} type="button" disabled={loading || !text.trim()}>
    {#if loading}
      Predicting…
    {:else}
      Predict
    {/if}
  </button>

  {#if label}
    <div class="result">
      Sentiment: {label}
    </div>
  {/if}

  {#if errorMsg}
    <div class="error">{errorMsg}</div>
  {/if}
</div>