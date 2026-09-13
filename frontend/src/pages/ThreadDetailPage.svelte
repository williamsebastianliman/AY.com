<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import LeftSidebar from "../lib/components/LeftSideBar.svelte";
  import RightSidebar from "../lib/components/RightSideBar.svelte";
  import Thread from "../lib/components/Thread.svelte";
  import { Loader2, Video, Image } from "lucide-svelte";
  import BurgerLeftSideBar from "../lib/components/BurgerLeftSideBar.svelte";
  import BurgerRightSideBar from "../lib/components/BurgerRightSideBar.svelte";
  import api from "../lib/api";

  export let params: { id: string };
  let myUserId = "";
  let myAvatar = "no-image";
  let mainThread = null;
  let loadingMain = true;
  let mainError = false;
  let replies: Thread[] = [];
  let loadingReplies = false;
  let repliesError = false;
  let repliesPage = 1;
  let hasMoreReplies = true;
  let observerTarget: HTMLDivElement;
  let observer: IntersectionObserver;
  import {
    enrichThread,
  } from "../lib/threadUtil";
  type ThreadMedia = {
    id: string;
    image_url: string;
    extension: string;
  };
  interface Thread {
    id: string;
    author: string;
    username: string;
    avatar: string;
    content: string;
    community_id?: string;
    parent_id?: string;
    like_count: number;
    comment_count: number;
    share_count: number;
    view_count: number;
    created_at: string;
    updated_at: string;
    liked_by_me?: boolean;
    bookmarked: boolean;
    media: ThreadMedia[];
  }
  let replyContent = "";
  let replyImages: UploadedMedia[] = [];
  let replyPosting = false;
  let replyError = false;
  let replyFileInput: HTMLInputElement | undefined;
  interface UploadedMedia {
    id: string;
    url: string;
  }
  async function loadMyProfile() {
    try {
      const meRes = await api.get("/user/get-me");
      myUserId = meRes.data.id;
      if (meRes.data.profile_picture_id){
        const mediaRes = await api.post("/media/get-media", { id: meRes.data.profile_picture_id }, { withCredentials: true });
        myAvatar = mediaRes.data.public_url;
      }
    } catch (err) {
      // Silent fail
      console.error(err as string);
    }
  }
  async function loadMainThread() {
    loadingMain = true;
    mainError = false;
    try {
      const res = await api.get(`/threads/${params.id}`);
      console.log("midd: ", myUserId);
      mainThread = await enrichThread(res.data, myUserId);
    } catch (e) {
      console.error(e);
      mainError = true;
    } finally {
      loadingMain = false;
    }
  }
  async function loadReplies() {
    if (loadingReplies || !hasMoreReplies) return;
    loadingReplies = true;
    repliesError = false;
    try {
      const res = await api.get(`/replies?page=${repliesPage}&size=12&parent_id=${params.id}`);
      const raws = res.data;
      if (!raws || !raws.length || raws.length < 12) hasMoreReplies = false;
      if (raws && raws.length) {
        const enriched = await Promise.all(raws.map((r) => enrichThread(r, myUserId)));
        replies = [...replies, ...enriched];
        repliesPage += 1;
      }
    } catch (e) {
      console.error(e);
      repliesError = true;
    } finally {
      loadingReplies = false;
    }
  }
  function setupIntersectionObserver() {
    observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting && !loadingReplies && hasMoreReplies) {
          loadReplies();
        }
      },
      { rootMargin: "200px" }
    );
    if (observerTarget) observer.observe(observerTarget);
  }
  async function handleReplyFiles(e: Event) {
    const input = e.target as HTMLInputElement;
    if (!input.files?.length) return;
    const file = input.files[input.files.length - 1];
    if (replyImages.length >= 4) {
      input.value = "";
      return;
    }
    try {
      const form = new FormData();
      form.append("file", file);
      const res = await api.post("/media/upload", form, {
        headers: { "Content-Type": "multipart/form-data" }
      });
      replyImages = [...replyImages, { id: res.data.id, url: res.data.public_url }];
    } catch (err) {
      console.error(err as string);
    } finally {
      input.value = "";
    }
  }
  function removeReplyImage(id: string) {
    replyImages = replyImages.filter(img => img.id !== id);
  }
  async function postReply() {
    if (!replyContent.trim()) return;
    replyPosting = true;
    replyError = false;
    try {
      const resp = await api.post("/threads", {
        user_id: myUserId,
        content: replyContent,
        community_id: null,
        parent_id: params.id
      }, { withCredentials: true });
      const newReplyId = resp.data.id;
      await Promise.all(replyImages.map(img =>
        api.post(`/threads/${newReplyId}/media`, { media_id: img.id })
      ));
      const raw = (await api.get(`/threads/${newReplyId}`)).data;
      const enriched = await enrichThread(raw, myUserId);
      replies = [enriched, ...replies];
      replyContent = "";
      replyImages = [];
    } catch (err) {
      console.error(err as string);
      replyError = true;
    } finally {
      replyPosting = false;
    }
  }
  onMount(async () => {
    await loadMyProfile();
    await loadMainThread();
    await loadReplies();
    setupIntersectionObserver();
  });
  onDestroy(() => {
    if (observer && observerTarget) observer.unobserve(observerTarget);
  });
  let windowWidth = window.innerWidth;
  function handleResize() {
    windowWidth = window.innerWidth;
  }
  onMount(() => {
    window.addEventListener("resize", handleResize);
    return () => window.removeEventListener("resize", handleResize);
  });
</script>

<div class="thread-detail-layout">
  <aside class="sidebar">
    {#if myUserId && windowWidth >= 1250}
      <LeftSidebar currentUserId={myUserId} activePage="premium" />
    {/if}
  </aside>
  {#if myUserId && windowWidth < 1250}
    <BurgerLeftSideBar currentUserId={myUserId} activePage="premium" />
  {/if}
<main class="main">
  {#if loadingMain}
    <div class="skeleton skeleton-thread"></div>
  {:else if mainError}
    <div class="error">Thread not found or failed to load.</div>
  {:else if mainThread}
    <Thread thread={mainThread} currentUserId={myUserId} isDetail={true} />
  {/if}

  <section class="compose reply-compose">
    <div class="compose-header">
      <img src={myAvatar} alt="You" class="compose-avatar" />
      <textarea
        class="compose-input"
        bind:value={replyContent}
        placeholder="Write your reply…"
      />
      <input
        type="file"
        accept="image/*,video/*"
        multiple
        bind:this={replyFileInput}
        on:change={handleReplyFiles}
        class="hidden"
      />
    </div>
    {#if replyImages.length}
      <div class="image-preview images-{replyImages.length}">
        {#each replyImages as img (img.id)}
          <div class="thumb">
            {#if img.url.endsWith(".mp4") || img.url.endsWith(".webm") || img.url.endsWith(".mov") || img.url.endsWith(".mkv")}
              <video src={img.url} controls preload="metadata" />
            {:else}
              <img src={img.url} alt="upload" />
            {/if}
            <button type="button" on:click={() => removeReplyImage(img.id)} class="del-btn">×</button>
          </div>
        {/each}
      </div>
    {/if}
    <div class="compose-footer">
      <div class="icon-row">
        <button type="button" class="icon-btn" aria-label="Add image" on:click={() => replyFileInput.click()}>
          <Image size="20" />
        </button>
        <button type="button" class="icon-btn" aria-label="Add video" on:click={() => replyFileInput.click()}>
          <Video size="20" />
        </button>
      </div>
      <button
        class="post-btn"
        type="button"
        disabled={!replyContent.trim() || replyPosting}
        on:click={postReply}
      >
        {#if replyPosting}Replying…{:else}Reply{/if}
      </button>
      {#if replyError}
        <p class="error">Failed to reply. Please try again.</p>
      {/if}
    </div>
  </section>

  <div class="replies-list">
    {#if replies.length === 0 && loadingReplies}
      {#each Array.from({ length: 3 }, (_, i) => i) as i (i)}
        <div class="skeleton skeleton-thread"></div>
      {/each}
    {:else}
      {#each replies as thread (thread.id)}
        <Thread {thread} currentUserId={myUserId} />
      {/each}
    {/if}
    {#if repliesError}
      <div class="error">Failed to load replies.<button type="button" on:click={loadReplies}>Retry</button></div>
    {/if}
    {#if loadingReplies && replies.length > 0}
      <div class="loading"><Loader2 size="20" class="loading-icon" /> Loading more…</div>
    {:else if !hasMoreReplies && replies.length > 0}
      <div class="end">No more replies.</div>
    {/if}
    <div bind:this={observerTarget} class="observer"></div>
  </div>
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
  .thread-detail-layout {
 display: grid;
 grid-template-columns: auto 1fr auto;
 height: 100vh;
}

.main {
 overflow-y: auto;
 padding: 1.5rem 1.5rem 0 1.5rem;
 width: 40vw;
 max-width: 40vw;
 margin: 0 auto;
 background: transparent;
}

.reply-box {
 margin-bottom: 2.5rem;
 background: rgba(32, 36, 48, 0.85);
 border: 1px solid rgba(255, 255, 255, 0.08);
 border-radius: 1rem;
 padding: 1rem;
 box-shadow: 0 2px 8px rgba(0,0,0,0.10);
}

.compose-header {
 display: flex;
 gap: 1rem;
 align-items: flex-start;
}

.compose-avatar {
 width: 48px;
 height: 48px;
 border-radius: 50%;
 object-fit: cover;
 box-shadow: 0 2px 8px rgba(0,0,0,0.14);
}

.compose-input {
 flex: 1;
 background: transparent;
 border: none;
 resize: none;
 color: #f9f9f9;
 font-size: 1.1rem;
 line-height: 1.4;
 min-height: 2.2rem;
 outline: none;
 padding: 0.4rem 0.8rem;
 border-radius: 0.5rem;
 transition: background 0.15s;
}

.compose-input::placeholder {
 color: #888;
}

.compose-footer {
 display: flex;
 justify-content: space-between;
 align-items: center;
 margin-top: 0.7rem;
}

.icon-row {
 display: flex;
 gap: 1.1rem;
}

.icon-btn {
 background: none;
 border: none;
 color: #1d9bf0;
 cursor: pointer;
 padding: 0.4rem;
 border-radius: 9999px;
 transition: background 0.18s, color 0.18s;
}

.icon-btn:hover {
 background: rgba(29,155,240,0.08);
 color: #299fff;
}

.post-btn {
 background: #1d9bf0;
 border: none;
 color: #fff;
 padding: 0.5rem 1.3rem;
 border-radius: 9999px;
 font-weight: 600;
 font-size: 1.02rem;
 transition: background 0.15s;
 box-shadow: 0 1px 4px rgba(29,155,240,0.08);
}
.post-btn:disabled {
 opacity: 0.55;
 cursor: not-allowed;
}

.post-btn:hover:not(:disabled) {
 background: #299fff;
}

.image-preview {
 display: flex;
 gap: 0.7rem;
 margin-top: 0.7rem;
 flex-wrap: wrap;
}

.thumb {
 border-radius: 0.7rem;
 overflow: hidden;
 position: relative;
 width: 110px;
 height: 110px;
 background: #191b20;
 display: flex;
 align-items: center;
 justify-content: center;
}

.thumb img,
.thumb video {
 width: 100%;
 height: 100%;
 object-fit: cover;
 background: #181a1f;
 border-radius: 0.7rem;
}

.del-btn {
 position: absolute;
 top: 6px;
 right: 6px;
 background: rgba(0,0,0,0.7);
 color: #fff;
 border: none;
 border-radius: 9999px;
 width: 1.5rem;
 height: 1.5rem;
 display: flex;
 align-items: center;
 justify-content: center;
 cursor: pointer;
 font-size: 1.15rem;
}

.replies-list {
 margin-bottom: 2.5rem;
}

.skeleton-thread {
 height: 86px;
 background: #1a1c22;
 border-radius: 1rem;
 margin-bottom: 1.1rem;
 animation: shimmer 1.5s infinite linear;
 background: linear-gradient(90deg, #23263b 25%, #222 50%, #23263b 75%);
 background-size: 200% 100%;
}

@keyframes shimmer {
 0% { background-position: 200% 0; }
 100% { background-position: -200% 0; }
}

.error {
 background: rgba(220,38,38,0.13);
 color: #fff;
 border-radius: 0.4rem;
 padding: 0.55rem 0.8rem;
 margin: 1rem 0;
}

.loading {
 text-align: center;
 color: #aaa;
 margin: 1.2rem 0 0.2rem 0;
}

.end {
 text-align: center;
 color: #5b5d66;
 margin: 1.5rem 0 0.4rem 0;
}

.observer {
 height: 10px;
 margin-bottom: 100px;
}

.sidebar {
 min-width: 60px;
 background: transparent;
}
.compose.reply-compose {
  background: rgba(30,30,30,0.8);
  border: 1px solid rgba(255,255,255,0.08);
  border-radius: 1rem;
  padding: 1rem;
  margin-bottom: 1.5rem;
}
.compose-header {
  display: flex;
  gap: 1rem;
  align-items: flex-start;
}
.compose-avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  object-fit: cover;
  box-shadow: 0 2px 8px rgba(0,0,0,0.14);
}
.compose-input {
  flex: 1;
  background: transparent;
  border: none;
  resize: none;
  color: #f9f9f9;
  font-size: 1.1rem;
  line-height: 1.4;
  min-height: 2.2rem;
  outline: none;
  padding: 0.4rem 0.8rem;
  border-radius: 0.5rem;
  transition: background 0.15s;
}
.compose-input::placeholder { color: #888; }
.compose-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 0.7rem;
}
.icon-row {
  display: flex;
  gap: 1.1rem;
}
.icon-btn {
  background: none;
  border: none;
  color: #1d9bf0;
  cursor: pointer;
  padding: 0.4rem;
  border-radius: 9999px;
  transition: background 0.18s, color 0.18s;
}
.icon-btn:hover {
  background: rgba(29,155,240,0.08);
  color: #299fff;
}
.post-btn {
  background: #1d9bf0;
  border: none;
  color: #fff;
  padding: 0.5rem 1.3rem;
  border-radius: 9999px;
  font-weight: 600;
  font-size: 1.02rem;
  transition: background 0.15s;
  box-shadow: 0 1px 4px rgba(29,155,240,0.08);
}
.post-btn:disabled {
  opacity: 0.55;
  cursor: not-allowed;
}
.post-btn:hover:not(:disabled) {
  background: #299fff;
}
.hidden {
  display: none !important;
}
.sidebar {
  z-index: 400000;
}
</style>