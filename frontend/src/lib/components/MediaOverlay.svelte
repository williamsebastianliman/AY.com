<script lang="ts">
    import { createEventDispatcher, onMount } from "svelte";
    import { X, ChevronLeft, ChevronRight } from "lucide-svelte";
    export let media: ThreadMedia[] = [];
    export let startIndex = 0;
    export let open = false;
    const dispatch = createEventDispatcher();
    let current = startIndex;
    type ThreadMedia = {
        id: string;
        image_url: string;
        extension: string;
    };
    function close() {
      dispatch("close");
    }
    function prev() {
      if (media.length > 1) current = (current - 1 + media.length) % media.length;
    }
    function next() {
      if (media.length > 1) current = (current + 1) % media.length;
    }
    function onBackdrop(e: MouseEvent) {
      if ((e.target as HTMLElement).classList.contains("media-overlay-backdrop")) close();
    }
    function isVideo(extension: string) {
        return [".mp4", ".webm", ".ogg"].includes(extension.toLowerCase());
    }
    // function isImage(extension: string) {
    //     return [".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp"].includes(extension.toLowerCase());
    // }
    onMount(() => {
      function esc(e: KeyboardEvent) {
        if (e.key === "Escape") close();
      }
      window.addEventListener("keydown", esc);
      return () => window.removeEventListener("keydown", esc);
    });
  </script>
  {#if open}
    <div class="media-overlay-backdrop" on:click={onBackdrop}>
      <div class="media-overlay-content">
        <button type="button" class="close-btn" aria-label="Close" on:click={close}><X size={28} /></button>
        {#if media.length > 1}
          <button type="button" class="nav-btn left" on:click={prev} aria-label="Previous media"><ChevronLeft size={36} /></button>
        {/if}
        <div class="media-display">
          {#if isVideo(media[current].extension)}
            <video src={media[current].image_url} controls autoplay class="media-el" />
          {:else}
            <img src={media[current].image_url} alt="Media" class="media-el" />
          {/if}
        </div>
        {#if media.length > 1}
          <button type="button" class="nav-btn right" on:click={next} aria-label="Next media"><ChevronRight size={36} /></button>
        {/if}
        <slot></slot>
      </div>
    </div>
  {/if}
  <style>
  .media-overlay-backdrop {
    z-index: 500000;
    position: fixed;
    inset: 0;
    background: rgba(10, 11, 15, 0.96);
    display: flex;
    align-items: center;
    justify-content: center;
  }
  .media-overlay-content {
    position: relative;
    background: transparent;
    max-width: 96vw;
    max-height: 96vh;
    display: flex;
    flex-direction: column;
    align-items: center;
  }
  .close-btn {
    position: absolute;
    top: 18px;
    left: 20px;
    z-index: 11;
    background: rgba(0,0,0,0.68);
    color: #fff;
    border: none;
    border-radius: 999px;
    padding: 6px;
    cursor: pointer;
    box-shadow: 0 2px 8px #0002;
    transition: background 0.16s;
  }
  .close-btn:hover { background: #222; }
  .nav-btn {
    position: absolute;
    top: 50%;
    transform: translateY(-50%);
    z-index: 10;
    background: rgba(0,0,0,0.65);
    color: #fff;
    border: none;
    border-radius: 50%;
    width: 48px;
    height: 48px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: background 0.16s;
  }
  .nav-btn.left  { left: 16px; }
  .nav-btn.right { right: 16px; }
  .nav-btn:hover { background: #222; }
  .media-display {
    max-width: 90vw;
    max-height: 85vh;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  .media-el {
    max-width: 70vw;
    max-height: 80vh;
    border-radius: 16px;
    box-shadow: 0 8px 44px #0008;
    background: #15151a;
  }
  </style>