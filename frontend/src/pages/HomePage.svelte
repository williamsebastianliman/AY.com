<script lang="ts">
  import api from "../lib/api";
  import { onMount, onDestroy } from "svelte";
  import { fade } from "svelte/transition";
  import LeftSidebar from "../lib/components/LeftSideBar.svelte";
  import BurgerLeftSideBar from "../lib/components/BurgerLeftSideBar.svelte";
  import BurgerRightSideBar from "../lib/components/BurgerRightSideBar.svelte";
  import RightSidebar from "../lib/components/RightSideBar.svelte";
  import Thread from "../lib/components/Thread.svelte";
  import DonutProgress from "../lib/components/DonutProgress.svelte";
  import {
    Loader2,
    Image,
    Video,
    BarChart2,
    Smile,
    Calendar,
    MapPin
  } from "lucide-svelte";
  import {
    enrichThread,
    type RawThread,
  } from "../lib/threadUtil";

  let fileInput: HTMLInputElement;
  type Tab = "forYou" | "following";
  let activeTab: Tab = "forYou";
  const savedTab = localStorage.getItem("activeHomeTab");
  if (savedTab === "forYou" || savedTab === "following") {
    activeTab = savedTab as Tab;
  }

  interface MeResponse {
    id: string;
    name: string;
    username: string;
    email?: string;
    gender?: string;
    birth_year?: string;
    birth_month?: string;
    birth_day?: string;
    subscribed_news?: boolean;
    profile_picture_id: string;
    banner_media_id?: string;
  }
  interface UploadedMedia {
    id: string;
    url: string;
  }

  let newContent = "";

  let threads: Thread[] = [];
  let loading = false;
  let initialLoad = true;
  let loadError = false;
  let page = 1;
  const size = 5;
  let hasMore = true;
  let posting = false;
  let postError = false;

  let observerTarget: HTMLDivElement;
  let observer: IntersectionObserver;
  let myAvatar = "no-image";
  let myUserId = "";

  let images: UploadedMedia[] = [];
  let communities: { id: string; name: string }[] = [];
  let selectedCommunity = "";
  let postAsCommunity = false;
  const maxWords = 280;
  $: wordCount = newContent.length;

  function removeLocalMedia(id: string) {
    images = images.filter(img => img.id !== id);
  }

  async function handleFiles(e: Event) {
    const input = e.target as HTMLInputElement;
    if (!input.files?.length) return;

    const file = input.files[input.files.length - 1];

    if (images.length >= 4) {
      input.value = "";
      return;
    }

    try {
      const form = new FormData();
      form.append("file", file);

      const res = await api.post<{ public_url: string; id: string }>(
        "/media/upload",
        form,
        { headers: { "Content-Type": "multipart/form-data" } }
      );

      images = [
        ...images,
        { id: res.data.id, url: res.data.public_url }
      ];
    } catch (err) {
      console.error("Upload failed:", err);
    } finally {
      input.value = "";
    }
  }

  async function postThread() {
    if (!newContent.trim()) return;

    posting = true;
    postError = false;

    try {
      const resp = await api.post<{ id: string }>(
        "/threads",
        {
          user_id: myUserId,
          content: newContent,
          community_id: null
        },
        { withCredentials: true }
      );

      const raw = (await api.get<RawThread>(
        `/threads/${resp.data.id}`
      )).data;
      const threadId = resp.data.id;
      await Promise.all(
        images.map(img =>
          api.post(`/threads/${threadId}/media`, { media_id: img.id })
        )
      );

      const enriched = await enrichThread(raw, myUserId);
      threads = [enriched, ...threads];
      newContent = "";
      images = [];
    } catch (err) {
      console.error("Post failed:", err);
      postError = true;
    } finally {
      posting = false;
    }
  }

  async function loadMyProfile() {
    try {
      const meRes = await api.get<MeResponse>("/user/get-me");
      myUserId = meRes.data.id;
      if (meRes.data.profile_picture_id){
        const mediaRes = await api.post<{
          id: string;
          public_url: string;
        }>(
          "/media/get-media",
          { id: meRes.data.profile_picture_id },
          { withCredentials: true }
        );
        myAvatar = mediaRes.data.public_url;
      }
    } catch (err) {
      console.error("Failed to load my profile:", err);
    }
  }

  async function loadMore() {
    if (loading || !hasMore) return;

    loading = true;
    loadError = false;

    let url = "";
    if (activeTab === "forYou") {
      url = `/threads?page=${page}&size=${size}`;
    } else if (activeTab === "following") {
      if (!myUserId) {
        loading = false;
        return;
      }
      url = `/threads/following?user_id=${myUserId}&page=${page}&size=${size}`;
    }

    try {
      const rawRes = await api.get<RawThread[]>(url);
      const raws = rawRes.data;
      if (raws.length < size) hasMore = false;

      const enriched = await Promise.all(raws.map((r) => enrichThread(r, myUserId)));
      threads = [...threads, ...enriched];
      page += 1;
      initialLoad = false;
    } catch (e) {
      console.error("Failed to load threads:", e);
      loadError = true;
    } finally {
      loading = false;
    }
  }


  function setupIntersectionObserver() {
    observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting && !loading && hasMore) {
          loadMore();
        }
      },
      { rootMargin: "5px" }
    );
    if (observerTarget) observer.observe(observerTarget);
  }

  function switchTab(tab: Tab) {
    if (tab === activeTab) return;
    activeTab = tab;
    localStorage.setItem("activeHomeTab", tab);
    threads = [];
    page = 1;
    hasMore = true;
    initialLoad = true;
    loadError = false;
    loadMore();
  }
  let windowWidth = window.innerWidth;
  function handleResize() {
    windowWidth = window.innerWidth;
  }
  onMount(() => {
    window.addEventListener("resize", handleResize);
    return () => window.removeEventListener("resize", handleResize);
  });
  onMount(async () => {
    await loadMyProfile();
    await loadMore();
    setupIntersectionObserver();
  });
  onDestroy(() => {
    if (observer && observerTarget) observer.unobserve(observerTarget);
  });
</script>

<div class="home layout">
  <aside class="sidebar">
    {#if myUserId && windowWidth >= 1250}
      <LeftSidebar currentUserId={myUserId} activePage="home" />
    {/if}
  </aside>
  {#if myUserId && windowWidth < 1250}
    <BurgerLeftSideBar currentUserId={myUserId} activePage="home" />
  {/if}

  <main class="main">
    <div class="tabs">
      <button
        type="button"
        class:active={activeTab === "forYou"}
        on:click={()=> switchTab("forYou")}
      >For you</button>
      <button
        type="button"
        class:active={activeTab === "following"}
        on:click={()=> switchTab("following")}
      >Following</button>
    </div>

    <section class="compose" in:fade>
      <div class="compose-header">
        <img src={myAvatar} alt="You" class="compose-avatar" />
        <textarea
          class="compose-input"
          bind:value={newContent}
          placeholder="What’s happening?"
          maxlength={maxWords * 6}
        />
        <input
          type="file"
          accept="image/*,video/*"
          multiple
          bind:this={fileInput}
          on:change={handleFiles}
          class="hidden"
        />
      </div>
      {#if images.length}
        <div class="image-preview images-{images.length}">
          {#each images as img (img.id)}
            <div class="thumb">
              {#if img.url.endsWith(".mp4") || img.url.endsWith(".webm") || img.url.endsWith(".mov") || img.url.endsWith(".mkv")}
                <video src={img.url} controls preload="metadata" />
              {:else}
                <img src={img.url} alt="upload" />
              {/if}
              <button type="button"
                on:click={() => removeLocalMedia(img.id)}
                class="del-btn">×</button>
            </div>
          {/each}
        </div>
      {/if}
      <div class="compose-footer">
        <div class="footer-left">
          <DonutProgress value={wordCount} max={maxWords} size={38} color="#1d9bf0" />
        </div>
        <select class="community-dropdown"
          bind:value={selectedCommunity}
          on:change={() => postAsCommunity = !!selectedCommunity}
        >
          <option value="">My Account</option>
          {#each communities as c (c.id)}
            <option value={c.id}>{c.name}</option>
          {/each}
        </select>
        <div class="icon-row">
          <button on:click={()=> fileInput.click()} type="button" class="icon-btn" aria-label="Add image"><Image size="20" /></button>
          <button type="button" class="icon-btn" aria-label="Add video" on:click={()=>fileInput.click()}><Video size="20" /></button>
          <button type="button" class="icon-btn" aria-label="Poll"><BarChart2 size="20" /></button>
          <button type="button" class="icon-btn" aria-label="Emoji"><Smile size="20" /></button>
          <button type="button" class="icon-btn" aria-label="Schedule"><Calendar size="20" /></button>
          <button type="button" class="icon-btn" aria-label="Location"><MapPin size="20" /></button>
        </div>
        <button
          class="post-btn"
          type="button"
          disabled={!newContent.trim() || posting}
          on:click={postThread}
        >
          {#if posting}
            Posting…
          {:else}
            Post
          {/if}
        </button>
        {#if postError}
          <p class="error">Failed to post. Please try again.</p>
        {/if}
      </div>
    </section>

    {#if initialLoad}
      {#each Array.from({ length: 3 }, (_, i) => i) as i (i)}
        <div class="skeleton skeleton-thread" in:fade={{ delay: i * 100 }} />
      {/each}
    {:else}
      {#each threads as thread (thread.id)}
        <Thread {thread} currentUserId={myUserId}
        on:deleted={(e) => {
          threads = threads.filter((t) => t.id !== e.detail.id);
        }} />
      {/each}

      {#if loadError}
        <div class="error" in:fade>
          <p>Failed to load threads.</p>
          <button type="button" class="retry-button" on:click={loadMore}>Retry</button>
        </div>
      {/if}
    {/if}

    {#if loading && !initialLoad}
      <div class="loading" in:fade>
        <Loader2 size="20" class="loading-icon" /> Loading more…
      </div>
    {:else if !hasMore && threads.length > 0}
      <div class="end" in:fade>You’ve reached the end!</div>
    {/if}

    <div bind:this={observerTarget} class="observer" />
  </main>

  <aside class="sidebar">
    {#if windowWidth >= 1460}
      <RightSidebar />
    {/if}
  </aside>
  {#if myUserId && windowWidth < 1460}
    <BurgerRightSideBar currentUserId={myUserId} activePage="home" />
  {/if}
</div>

<style>
  .home::-webkit-scrollbar { display: none; }
  .home, .home * {
    overflow-y: hidden;
    scrollbar-width: none;
    -ms-overflow-style: none;
  }

  /* :global(body) {
    margin: 0;
    padding: 0;
    overflow-x: auto;
    overflow-y: hidden;
    font-family: "Inter", sans-serif;
    background: #0a0a0a;
    color: #f9f9f9;
  } */

  .hidden {
    display: none !important;
  }

  .layout {
    display: grid;
    grid-template-columns: auto 1fr auto;
    height: 100vh;
  }

  .main {
      overflow-y: auto;
      padding: 1rem;
      width: 40vw;
      max-width: 40vw;
      margin: 0 auto;
  }

  .compose {
    background: rgba(30,30,30,0.8);
    border: 1px solid rgba(255,255,255,0.08);
    border-radius: 1rem;
    padding: 1rem;
    margin-bottom: 1.5rem;
  }
  .compose-header {
    display: flex;
    gap: 1rem;
  }
  .compose-avatar {
    width: 48px; height: 48px; border-radius: 50%; object-fit: cover;
  }
  .compose-input {
    flex: 1;
    background: transparent;
    border: none;
    resize: none;
    color: #f9f9f9;
    font-size: 1.1rem;
    line-height: 1.4;
    min-height: 3rem;
    outline: none;
  }
  .compose-input::placeholder { color: #666; }
  .compose-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: 2.75rem;
  }
  .icon-row { display: flex; gap: 1rem; }
  .icon-btn {
    background: none; border: none;
    color: #1d9bf0; cursor: pointer;
    padding: 0.3rem; transition: color 0.2s;
  }
  .icon-btn:hover { color: #299fff; }
  .post-btn {
    background: #1d9bf0; border: none;
    color: #fff; padding: 0.5rem 1.25rem;
    border-radius: 9999px; font-weight: 600;
  }
  .post-btn:disabled { opacity: 0.5; cursor: not-allowed; }
  .post-btn:hover:not(:disabled) { background: #299fff; }

  .tabs{
    margin-top: 5vh;
    display: flex;
    margin-bottom: 1.5rem;
    border-bottom: 1px solid rgba(255,255,255,0.1);
    background: rgba(10,10,10,0.8);
    backdrop-filter: blur(8px);
    z-index: 10;
    padding: 0.5rem 0;
  }

  .tabs button {
    background: none; border: none;
    padding: 0.75rem 1.5rem;
    font-weight: 600; font-size: 0.95rem;
    color: rgba(255,255,255,0.6);
    border-bottom: 2px solid transparent;
    cursor: pointer; transition: color 0.2s;
  }
  .tabs button:hover { color: #fff; }
  .tabs button.active {
    color: #fff; border-color: #1d9bf0;
  }

  .thread-container { margin-bottom: 1rem; }
  .thread {
    display: flex; gap: 1rem;
    padding: 1.25rem;
    background: rgba(30,30,30,0.6);
    border: 1px solid rgba(255,255,255,0.08);
    border-radius: 1rem;
  }
  .thread:hover {
    background: rgba(35,35,35,0.8);
    transform: translateY(-2px);
  }
  .avatar {
    width: 48px; height: 48px; border-radius: 50%;
    object-fit: cover; box-shadow: 0 2px 8px rgba(0,0,0,0.2);
  }
  .body { flex: 1; }
  .header {
    display: flex; align-items: center;
    color: rgba(255,255,255,0.7);
    margin-bottom: 0.25rem;
  }
  .header strong { color: #fff; margin-right: 0.5rem; }
  .dot { margin: 0 0.5rem; opacity: 0.5; }
  .community {
    display: inline-block;
    background: rgba(29,155,240,0.1);
    color: #1d9bf0;
    padding: 0.2rem 0.5rem;
    border-radius: 4px;
    font-size: 0.85rem;
    margin-bottom: 0.75rem;
  }
  .text {
    margin: 0.75rem 0;
    line-height: 1.5;
    white-space: pre-wrap;
  }
  .actions {
    display: flex; gap: 2rem;
    color: rgba(255,255,255,0.6);
    margin-top: 1rem;
  }
  .actions button.liked {
    color: #f91880;
  }
  .actions button {
    background: none; border: none;
    display: flex; align-items: center;
    gap: 0.5rem; cursor: pointer;
    transition: color 0.15s;
  }
  .actions button:hover { color: #fff; }
  .actions button:hover:nth-child(1) { color: #1d9bf0; }
  .actions button:hover:nth-child(3) { color: #f91880; }

  .loading, .end {
    text-align: center;
    padding: 2rem 0;
    color: rgba(255,255,255,0.7);
  }
  .loading-icon {
    animation: spin 1s linear infinite;
    margin-right: 0.5rem;
  }
  @keyframes spin {
    from { transform: rotate(0deg); }
    to { transform: rotate(360deg); }
  }
  .error {
    background: rgba(220,38,38,0.1);
    border: 1px solid rgba(220,38,38,0.3);
    color: #fff;
    padding: 1rem;
    border-radius: 0.5rem;
    text-align: center;
    margin: 1rem 0;
  }
  .sidebar{
    z-index: 400000;
  }
  .skeleton {
    border-radius: 1rem;
    margin-bottom: 1rem;
    background: linear-gradient(
      90deg,
      rgba(40,40,40,0.4) 25%,
      rgba(60,60,60,0.4) 50%,
      rgba(40,40,40,0.4) 75%
    );
    background-size: 200% 100%;
    animation: shimmer 1.5s infinite;
  }
  @keyframes shimmer {
    0%   { background-position: 200% 0 }
    100% { background-position: -200% 0 }
  }
  .skeleton-thread { height: 120px; }

  .observer{
    height: 10px;
    margin-bottom: 200px;
  }
  textarea{
    font-family: sans-serif;
  }
  .image-preview {
    display: grid;
    gap: 0.7rem;
    margin-top: 0.5rem;
    max-width: 100%;
    max-height: 40vh;
    overflow: auto;
    grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
    align-items: center;
  }

  .image-preview.images-1 {
    grid-template-columns: 1fr;
    max-width: 480px;
    margin-left: auto;
    margin-right: auto;
  }
  .image-preview.images-1 .thumb {
    aspect-ratio: 16/9;
    min-height: 180px;
  }

  .image-preview.images-2 {
    grid-template-columns: 1fr 1fr;
  }
  .image-preview.images-2 .thumb {
    aspect-ratio: 4/3;
  }

  .image-preview.images-3,
  .image-preview.images-4 {
    grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
  }

  .image-preview.images-5,
  .image-preview.images-6,
  .image-preview.images-7,
  .image-preview.images-8 {
    grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
    max-height: 40vh;
    overflow-y: auto;
  }

  .thumb {
    position: relative;
    border-radius: 0.75rem;
    overflow: hidden;
    background: #111;
    display: flex;
    align-items: center;
    justify-content: center;
    aspect-ratio: 16/9;
    min-width: 0;
  }

  .thumb img,
  .thumb video {
    width: 100%;
    height: 100%;
    object-fit: cover;
    display: block;
    border-radius: 0.75rem;
    background: #181818;
    max-height: 40vh;
  }

  .thumb video {
    background: #16161a;
  }

  .del-btn {
    position: absolute;
    top: 4px; right: 4px;
    background: rgba(0,0,0,0.6);
    color: #fff;
    border: none;
    border-radius: 9999px;
    width: 1.5rem; height: 1.5rem;
    display: flex;
    align-items: center; justify-content: center;
    cursor: pointer;
  }
  .community-row {
  display: flex;
  justify-content: flex-end;
  margin-top: -2.3rem;
  margin-bottom: 3.6rem;
  margin-right: 1.4rem;
  z-index: 2;
  position: relative;

}

.community-dropdown {
  border-radius: 999px;
  border: none;
  background: #222d33;
  color: #5ac6ff;
  font-size: 1rem;
  font-weight: 600;
  padding: 0.3rem 1.3rem;
  outline: none;
  min-width: 100px;
  cursor: pointer;
  box-shadow: 0 1px 7px #0001;
  appearance: none;
  transition: background .12s;
}
.community-dropdown:focus, .community-dropdown:hover {
  background: #232a2f;
}

@media (max-width: 1250px) {
  .main {
    width: 98vw;
    max-width: 98vw;
    padding: 0.7rem 0.5rem;
    overflow-y: auto;
  }
  .compose,
  .thread {
    padding: 0.7rem;
    border-radius: 0.8rem;
    width: 100%;
  }
  .compose-header {
    gap: 0.6rem;
  }
  .compose-avatar {
    width: 40px;
    height: 40px;
  }
  .compose-input {
    font-size: 1rem;
    min-height: 2.2rem;
  }
  .tabs button {
    font-size: 0.97rem;
    padding: 0.6rem 1rem;
  }
  .compose-footer {
    flex-direction: column;
    align-items: stretch;
    gap: 0.7rem;
    margin-top: 1.6rem;
  }
  .icon-row {
    gap: 0.7rem;
  }
  .post-btn {
    width: 100%;
    margin-top: 0.5rem;
  }
}

@media (max-width: 1460px) {
  .main {
    width: 60vw;
    max-width: 60vw;
    padding: 0.7rem 0.5rem;
    overflow-y: auto;
  }
  .compose,
  .thread {
    min-width: 50vw;
    padding: 0.7rem;
    border-radius: 0.8rem;
  }
  .compose-header {
    gap: 0.6rem;
  }
  .compose-avatar {
    width: 40px;
    height: 40px;
  }
  .compose-input {
    font-size: 1rem;
    min-height: 2.2rem;
  }
  .tabs button {
    font-size: 0.97rem;
    padding: 0.6rem 1rem;
  }
  .compose-footer {
    flex-direction: column;
    align-items: stretch;
    gap: 0.7rem;
    margin-top: 1.6rem;
  }
  .icon-row {
    gap: 0.7rem;
  }
  .post-btn {
    width: 100%;
    margin-top: 0.5rem;
  }
}

@media (max-width: 500px) {
  .layout {
    grid-template-columns: none;
    min-height: 99vh;
    width: 100vw;
    margin: 0;
    padding: 0;
    background: #000;
    overflow-x: hidden;
  }
  .main {
    width: 100vw;
    max-width: 100vw;
    padding: 0.5rem 0.02rem 1.2rem 0.02rem;
    margin: 0;
    min-height: 100vh;
    border-radius: 0;
    box-shadow: none;
    overflow-x: hidden;
  }
}

</style>