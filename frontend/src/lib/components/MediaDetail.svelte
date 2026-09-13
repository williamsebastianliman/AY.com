    <script lang="ts">
      import { createEventDispatcher, onMount, onDestroy } from "svelte";
      import { fade, fly } from "svelte/transition";

      export let threadId: string;
      export let mediaUrls: string[] = [];

      const dispatch = createEventDispatcher();
      let current = 0;

      function close() {
        dispatch("close");
      }
      function next() {
        current = (current + 1) % mediaUrls.length;
      }
      function prev() {
        current = (current - 1 + mediaUrls.length) % mediaUrls.length;
      }

      function onKey(e: KeyboardEvent) {
        if (e.key === "Escape") close();
        if (e.key === "ArrowRight") next();
        if (e.key === "ArrowLeft") prev();
      }

      onMount(() => window.addEventListener("keydown", onKey));
      onDestroy(() => window.removeEventListener("keydown", onKey));
    </script>


  <div
    class="overlay"
    role="dialog"
    aria-modal="true"
    tabindex="-1"
    on:click|self={close}
    on:keydown={(e) => e.key === "Escape" && close()}
    transition:fade
  >
    <div class="carousel" transition:fly={{ y: 20, duration: 200 }}>
      <button
        class="btn prev"
        on:click|stopPropagation={prev}
        aria-label="Previous slide"
        type="button"
      >
        ←
      </button>

      <img
        src={mediaUrls[current]}
        alt={"Media " + (current + 1)}
      />

      <button
        class="btn next"
        on:click|stopPropagation={next}
        aria-label="Next slide"
        type="button"
      >
        →
      </button>

      <button
        class="close"
        on:click|stopPropagation={close}
        aria-label="Close gallery"
        type="button"
      >
        ×
      </button>

      <div class="dots">
        {#each mediaUrls as url, i (url)}
          <button
            class="dot"
            class:active={i === current}
            aria-label={"Go to slide " + (i + 1)}
            on:click|stopPropagation={() => current = i}
            type="button"
          >
          </button>
        {/each}
      </div>
    </div>
  </div>


    <style>
      .overlay {
        position: fixed;
        top: 0; left: 0; right: 0; bottom: 0;
        background: rgba(0,0,0,0.8);
        display: flex;
        justify-content: center;
        align-items: center;
        z-index: 1000;
      }
      .carousel {
        position: relative;
        max-width: 80vw;
        max-height: 80vh;
        display: flex;
        align-items: center;
        overflow: hidden;
      }
      img {
        max-width: 100%;
        max-height: 100%;
        user-select: none;
      }
      .btn {
        position: absolute;
        background: rgba(255,255,255,0.3);
        border: none;
        padding: 0.5rem;
        border-radius: 50%;
        cursor: pointer;
        transition: background 0.2s;
      }
      .btn:hover { background: rgba(255,255,255,0.6); }
      .prev { left: 1rem; }
      .next { right: 1rem; }
      .close {
        position: absolute;
        top: 1rem; right: 1rem;
        background: none;
        color: white;
        font-size: 2rem;
        border: none;
        cursor: pointer;
      }
      .dots {
        position: absolute;
        bottom: 1rem;
        display: flex;
        gap: 0.5rem;
      }
      .dot {
        width: 0.75rem; height: 0.75rem;
        background: rgba(255,255,255,0.5);
        border-radius: 50%;
        cursor: pointer;
        transition: background 0.2s;
      }
      .dot.active { background: white; }
    </style>