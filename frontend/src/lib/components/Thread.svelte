<script lang="ts">
  import { fly } from "svelte/transition";
  import { MessageCircle, Repeat, Heart, Eye, Share2, Bookmark,Pin} from "lucide-svelte";
  import api from "../api";
  import { MoreVertical } from "lucide-svelte";
  import { clickOutside } from "../action/clickOutside";
  import { createEventDispatcher } from "svelte";
  import { BadgeCheck } from "lucide-svelte";
  import MediaOverlay from "../components/MediaOverlay.svelte";
  import RichText from "./RichText.svelte";

  let showMenu = false;
  const dispatch = createEventDispatcher();
  let menuEl: HTMLElement;

  export let thread: {
    id: string;
    user_id: string;
    avatar: string;
    author: string;
    is_premium: boolean;
    username: string;
    created_at: string;
    community_id?: string;
    content: string;
    media: ThreadMedia[];
    comment_count: number;
    share_count: number;
    like_count: number;
    liked_by_me: boolean;
    view_count: number;
    bookmarked: boolean;
    repost_count: number;
    is_reposted: boolean;
    is_pinned: boolean;
  };
  let showOverlay = false;
  let overlayIndex = 0;

  function openMediaOverlay(idx: number) {
    overlayIndex = idx;
    showOverlay = true;
  }
  function closeMediaOverlay() {
    showOverlay = false;
  }

  export let is_profile_thread: boolean = false;

  type ThreadMedia = {
    id: string;
    image_url: string;
    extension: string;
  };
  export let currentUserId: string;
  console.log("cur-state: ", thread);
  console.log("cur uid: ", currentUserId);

  function goToUserProfile() {
    window.location.href = `/profile/${thread.user_id}`;
  }
  function fmtDate(iso: string) {
    const d = new Date(iso);
    if (isNaN(d.getTime())) return "";
    const diff = Date.now() - d.getTime();
    const mins = Math.floor(diff / 60000);
    if (mins < 1) return "just now";
    if (mins < 60) return `${mins}m ago`;
    const hrs = Math.floor(mins / 60);
    if (hrs < 24) return `${hrs}h ago`;
    const days = Math.floor(hrs / 24);
    return days > 7 ? d.toLocaleDateString() : `${days}d ago`;
  }

  async function toggleLike(): Promise<void> {
    const url = `/threads/${thread.id}/like`;
    try {
      if (thread.liked_by_me) {
        await api.delete(url, { data: { user_id: currentUserId } });
        thread.like_count -= 1;
        thread.liked_by_me = false;
      } else {
        await api.post(url, { user_id: currentUserId });
        thread.like_count += 1;
        thread.liked_by_me = true;
      }
      thread = { ...thread };
    } catch (err) {
      console.error("Failed to toggle like:", err);
    }
  }

  async function toggleBookmark(): Promise<void> {
    const url = `/threads/${thread.id}/bookmark`;
    try {
      if (thread.bookmarked) {
        await api.delete(url, { data: { user_id: currentUserId } });
        thread.bookmarked = false;
      } else {
        await api.post(url, { user_id: currentUserId });
        thread.bookmarked = true;
      }
      thread = { ...thread };
    } catch (err) {
      console.error("Failed to toggle bookmark:", err);
    }
  }

  async function handleDelete(id: string) {
    const url = `threads/${id}`;
    try {
      await api.delete(url);
      dispatch("deleted", { id });
    } catch (err) {
      console.log(err as string);
    }
    showMenu = false;
  }

  function isVideo(extension: string) {
    return [".mp4", ".webm", ".ogg"].includes(extension.toLowerCase());
  }
  function isImage(extension: string) {
    return [".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp"].includes(extension.toLowerCase());
  }

  function handleCardClick(event: MouseEvent) {
    const path = event.composedPath();
    for (const node of path) {
      if (
        node instanceof HTMLElement &&
        (
          node.tagName === "BUTTON" ||
          node.tagName === "A" ||
          node.tagName === "INPUT" ||
          node.tagName === "TEXTAREA" ||
          node.classList?.contains("actions") ||
          node.classList?.contains("header-actions") ||
          node.classList?.contains("dropdown-menu")
        )
      ) {
        return;
      }
    }
    window.location.href = `/thread/${thread.id}`;
  }
  async function handleRepost(threadId: string) {
    try {
      let currentId = threadId;

      while (true) {
        const res = await api.get<{ repost_id: string }>(`/threads/${currentId}`);
        const { repost_id } = res.data;
        if (!repost_id || repost_id === "") break;
        currentId = repost_id;
      }
      if (!thread.is_reposted){
          const payload = {
          user_id: currentUserId,
          content: "Reposted",
          parent_id: "",
          community_id: "",
          repost_id: currentId,
        };
        await api.post("/threads/repost", payload);
        thread.is_reposted = true;
      }

    } catch (err) {
      console.error("Failed to repost:", err);
    }
  }
  async function handlePin(threadID: string) {
    try {
      await api.patch(`/threads/pin/${threadID}`);
      dispatch("pinned");
    } catch (error) {
      console.error("Failed to pin thread:", error);
    }
    showMenu = false;
  }

</script>

<div class="thread-container" in:fly={{ y: 20, duration: 300 }}>
  <div class="thread" on:click={handleCardClick} tabindex="0">
    <img src={thread.avatar} alt={thread.author} class="avatar" on:click|stopPropagation={goToUserProfile} />
    <div class="body">
      <div class="header-row">
        <div class="header">
          <div class="header-info">
            <strong class="profile-link" on:click|stopPropagation={goToUserProfile}>
              {thread.author}
              {#if thread.is_premium}
                <BadgeCheck size="16" class="verified-badge" />
              {/if}
            </strong>
            <span class="handle">@{thread.username}</span>
            <small class="timestamp">{fmtDate(thread.created_at)}</small>
            {#if thread.is_pinned && is_profile_thread}
              <Pin size="16" class="pin-badge" title="Pinned" style="margin-left: 0.4em; color: #f4d35e;" />
            {/if}
          </div>
        </div>
        {#if currentUserId == thread.user_id}
          <div class="header-actions" use:clickOutside on:outclick={() => showMenu = false} bind:this={menuEl}>
            <button type="button" aria-label="More" class="three-dot" on:click={() => showMenu = !showMenu}>
              <MoreVertical size="20" />
            </button>
            {#if showMenu}
              <div class="dropdown-menu">
                <button type="button" class="dropdown-item" on:click={handleDelete(thread.id)}>Delete</button>
                {#if is_profile_thread}
                  <button type="button" class="dropdown-item" on:click={handlePin(thread.id)}>Pin</button>
                {/if}
              </div>
            {/if}
          </div>
        {/if}
      </div>
      {#if thread.community_id}
        <div class="community">{thread.community_id}</div>
      {/if}

      {#if thread.repost_thread}
          <div
            class="repost-banner clickable"
            on:click|stopPropagation={() => window.location.href = `/thread/${thread.repost_thread.id}`}
            title="Go to original thread"
            tabindex="0"
            role="button"
          >
            <span class="repost-label">Reposted</span>
            <div class="repost-content">
              <div class="repost-header">
                <img src={thread.repost_thread.avatar} alt={thread.repost_thread.author} class="avatar-small" />
                <span><strong>{thread.repost_thread.author}</strong> <span class="repost-username">@{thread.repost_thread.username}</span></span>
              </div>
              <RichText content={thread.repost_thread.content} />
              {#if thread.repost_thread.media && thread.repost_thread.media.length}
                <div class="image-preview images-{thread.repost_thread.media.length}">
                  {#each thread.repost_thread.media as item (item.image_url)}
                    <div class="thumb">
                      {#if isVideo(item.extension)}
                        <video src={item.image_url} controls preload="metadata"></video>
                      {:else if isImage(item.extension)}
                        <img src={item.image_url} alt="thread media" />
                      {:else}
                        <span>Unsupported format: {item.extension}</span>
                      {/if}
                    </div>
                  {/each}
                </div>
              {/if}
            </div>
          </div>
      {:else}
        <RichText content={thread.content} />
        {#if thread.media && thread.media.length}
          <div class="image-preview images-{thread.media.length}">
            {#each thread.media as item, idx (item.image_url)}
            <div class="thumb" on:click|stopPropagation={() => openMediaOverlay(idx)}>
              {#if isVideo(item.extension)}
                <video src={item.image_url} controls preload="metadata" />
                <div class="media-play-overlay">
                </div>
              {:else if isImage(item.extension)}
                <img src={item.image_url} alt="thread media" />
              {:else}
                <span>Unsupported format: {item.extension}</span>
              {/if}
            </div>
            {/each}
          </div>
        {/if}
      {/if}

      <div class="actions">
        <button type="button" aria-label="Comment">
          <MessageCircle size="18" /><span>{thread.comment_count}</span>
        </button>
        <button
          type="button"
          aria-label="Repost"
          on:click|stopPropagation={() => handleRepost(thread.id)}
          class:reposted={thread.is_reposted}
        >
          <Repeat size="18" /><span>{thread.repost_count}</span>
        </button>
        <button
          type="button"
          aria-label="Like"
          on:click={() => toggleLike(thread)}
          class:liked={thread.liked_by_me}
        >
          <Heart size="18" /><span>{thread.like_count}</span>
        </button>
        <button type="button" aria-label="View">
          <Eye size="18" /><span>{thread.view_count}</span>
        </button>
        <button type="button" aria-label="Share">
          <Share2 size="18" />
        </button>
        <button type="button" aria-label="Bookmark" on:click={() => toggleBookmark(thread)} class:bookmarked={thread.bookmarked}>
          <Bookmark size="18" />
        </button>
      </div>
    </div>
  </div>
  <MediaOverlay
    media={thread.media}
    startIndex={overlayIndex}
    open={showOverlay}
    on:close={closeMediaOverlay}
  />
</div>

<style>
.thread-container { margin-bottom: 1rem; }
.thread {
  display: flex;
  gap: 1rem;
  padding: 1.25rem;
  background: rgba(30,30,30,0.6);
  border: 1px solid rgba(255,255,255,0.08);
  border-radius: 1rem;
  position: relative;
  transition: background 0.2s, transform 0.2s;
  cursor: pointer;
}
.thread:focus { outline: 2px solid #1d9bf0; }
.thread:hover {
  background: rgba(35,35,35,0.8);
  transform: translateY(-2px);
}
.avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  object-fit: cover;
  box-shadow: 0 2px 8px rgba(0,0,0,0.2);
}
.body { flex: 1; }
.header-row {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: 0.2rem;
}
.header {
  display: flex;
  align-items: center;
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
.image-preview {
  display: grid;
  gap: 0.5rem;
  margin-top: 0.5rem;
}
.image-preview.images-1 { grid-template-columns: 1fr; }
.image-preview.images-2 { grid-template-columns: 1fr 1fr; }
.image-preview.images-3 {
  grid-auto-flow: column;
  grid-auto-columns: calc((100% - 0.5rem * 2) / 3);
  overflow-x: auto;
}
.image-preview.images-4 {
  grid-auto-flow: column;
  grid-auto-columns: calc((100% - 0.5rem * 3) / 4);
  overflow-x: auto;
}
.thumb {
  position: relative;
  border-radius: 0.5rem;
  overflow: hidden;
  cursor:"pointer"
}
.thumb img {
  width: 100%;
  height: auto;
  display: block;
}
.actions {
  display: flex;
  gap: 1rem;
  color: rgba(255,255,255,0.6);
  margin-top: 1rem;
}
.actions button {
  background: none;
  border: none;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  cursor: pointer;
  transition: color 0.15s;
  z-index: 3;
  position: relative;
}
.actions button:hover { color: #fff; }
.actions button.liked { color: #f91880; }
.actions button.reposted { color: #199610; }
.actions button.bookmarked { color: #1d9bf0; }
.actions button:hover:nth-child(1) { color: #1d9bf0; }
.actions button:hover:nth-child(3) { color: #f91880; }

.header-actions {
  position: relative;
  display: flex;
  align-items: center;
  margin-left: 1rem;
  z-index: 3;
}
.three-dot {
  background: none;
  border: none;
  cursor: pointer;
  padding: 0.2rem;
  border-radius: 9999px;
  transition: background 0.13s;
}
.three-dot:hover {
  background: rgba(255,255,255,0.09);
}
.dropdown-menu {
  position: absolute;
  right: 0;
  top: 2rem;
  min-width: 120px;
  background: #181c24;
  color: #fff;
  border-radius: 0.5rem;
  box-shadow: 0 2px 10px rgba(0,0,0,0.12);
  padding: 0.25rem 0;
  z-index: 20;
  border: 1px solid rgba(255,255,255,0.08);
  display: flex;
  flex-direction: column;
}
.dropdown-item {
  background: none;
  border: none;
  width: 100%;
  color: #fff;
  padding: 0.55rem 1.1rem;
  font-size: 1rem;
  text-align: left;
  cursor: pointer;
  border-radius: 0.3rem;
  transition: background 0.13s;
}
.dropdown-item:hover {
  background: rgba(240,57,77,0.16);
  color: #f14e67;
}
video {
  width: 100%;
  height: auto;
}
.repost-banner {
  border-left: 4px solid #1d9bf0;
  background: rgba(29,155,240,0.04);
  padding: 0.7rem 1rem;
  margin-bottom: 1rem;
  border-radius: 0.75rem;
}
.repost-content {
  margin-top: 0.3rem;
}
.repost-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 0.2rem;
}
.avatar-small {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  object-fit: cover;
  box-shadow: 0 1px 3px rgba(0,0,0,0.09);
}
.repost-username {
  color: #aaa;
  font-size: 0.96em;
  margin-left: 0.25em;
}
.repost-label {
  color: #1d9bf0;
  font-size: 0.92em;
  font-weight: 500;
  margin-bottom: 0.3rem;
  display: block;
}
.media-play-overlay {
  position: absolute;
  top: 0; left: 0; right: 0; bottom: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 3rem;
  opacity: 0.85;
  pointer-events: none;
}

.image-preview {
  display: grid;
  gap: 0.35rem;
  margin-top: 0.7rem;
  width: 100%;
  max-width: 520px;
  max-height: 520px;
  border-radius: 18px;
  overflow: hidden;
}

.image-preview.images-1 {
  grid-template-columns: 1fr;
  grid-template-rows: 1fr;
  aspect-ratio: 1/1;
  height: 320px;
}

.image-preview.images-2 {
  grid-template-columns: 1fr 1fr;
  grid-template-rows: 1fr;
  aspect-ratio: 2/1;
  height: 260px;
}

.image-preview.images-3,
.image-preview.images-4 {
  grid-template-columns: 1fr 1fr;
  grid-template-rows: 1fr 1fr;
  aspect-ratio: 1/1;
  height: 320px;
}

.thumb {
  position: relative;
  width: 100%;
  height: 100%;
  aspect-ratio: 1/1;
  overflow: hidden;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #191b1f;
  box-sizing: border-box;
}

.thumb img,
.thumb video {
  width: 100%;
  height: 100%;
  object-fit: cover;
  object-position: center;
  display: block;
  border-radius: 14px;
}

.media-play-overlay {
  display: none !important;
}

.thumb video::-webkit-media-controls {
  opacity: 0.92;
}
.thumb video {
  background: #131518;
  pointer-events: none;
  filter: brightness(0.93);
  pointer-events:none;
}

@media (max-width: 650px) {
  .image-preview,
  .image-preview.images-1,
  .image-preview.images-2,
  .image-preview.images-3,
  .image-preview.images-4 {
    max-width: 92vw;
    max-height: 92vw;
    height: auto;
    aspect-ratio: 1/1;
  }
}

.image-preview:last-child {
  margin-bottom: 0;
}

@media (max-width: 1200px) {
  .thread {
    flex-direction: column;
    padding: 1rem;
    gap: 0.7rem;
  }
  .avatar {
    width: 40px;
    height: 40px;
    margin-bottom: 0.5rem;
  }
  .header-row {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.3rem;
  }
  .header-actions {
    margin-left: 0;
    margin-top: 0.5rem;
    align-self: flex-end;
  }
  .body {
    width: 100%;
  }
  .actions {
    gap: 0.7rem;
    flex-wrap: wrap;
  }
  .image-preview,
  .image-preview.images-1,
  .image-preview.images-2,
  .image-preview.images-3,
  .image-preview.images-4 {
    max-width: 100vw;
    max-height: 70vw;
    height: auto;
    border-radius: 12px;
  }
  .actions {
    justify-content: center;
    margin: 1rem 0;
  }

  .actions button[aria-label="View"],
  .actions button[aria-label="Share"] {
    display: none;
  }

  .header-info strong,
  .header-info .handle {
    font-size: 0.9rem;
  }

  .header-info {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
  }
}

@media (max-width: 600px) {
  .thread {
    padding: 0.7rem;
    border-radius: 0.7rem;
    font-size: 0.95rem;
  }
  .avatar {
    width: 36px;
    height: 36px;
  }
  .header strong, .header .handle, .community, .dot, small {
    font-size: 0.96em;
  }
  .actions button span {
    font-size: 0.98em;
  }
  .dropdown-menu {
    min-width: 95px;
    font-size: 0.92em;
  }
  .image-preview,
  .image-preview.images-1,
  .image-preview.images-2,
  .image-preview.images-3,
  .image-preview.images-4 {
    max-width: 98vw;
    max-height: 55vw;
    border-radius: 8px;
  }
}


</style>