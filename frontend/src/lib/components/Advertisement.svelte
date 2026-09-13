<script lang="ts">
    import { fly } from "svelte/transition";
    import { MessageCircle, Repeat, Heart, Eye, Share2, Bookmark } from "lucide-svelte";
    import api from "../api";
    import { MoreVertical } from "lucide-svelte";
    import { clickOutside } from "../action/clickOutside";
    import { createEventDispatcher } from "svelte";
    import { BadgeCheck } from "lucide-svelte";
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
    };
    type ThreadMedia = {
      id: string;
      image_url: string;
      extension: string;
    };
    export let currentUserId: string;
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
  </script>
  <div class="thread-container" in:fly={{ y: 20, duration: 300 }}>
    <div class="thread" on:click={handleCardClick} tabindex="0">
      <img src={thread.avatar} alt={thread.author} class="avatar" on:click|stopPropagation={goToUserProfile} />
      <div class="body">
        <div class="header-row">
          <div class="header">
            <strong class="profile-link" on:click|stopPropagation={goToUserProfile}>
              {thread.author}
              {#if thread.is_premium}
                <BadgeCheck size="16" class="verified-badge" />
              {/if}
            </strong>
            <p
              class="profile-link handle"
              on:click|stopPropagation={goToUserProfile}
            >
              @{thread.username}
            </p>
            <span class="dot">·</span>
            <small>{fmtDate(thread.created_at)}</small>
          </div>
          <div class="header-actions" use:clickOutside on:outclick={() => showMenu = false} bind:this={menuEl}>
            <button type="button" aria-label="More" class="three-dot" on:click={() => showMenu = !showMenu}>
              <MoreVertical size="20" />
            </button>
            {#if showMenu}
              <div class="dropdown-menu">
                <button type="button" class="dropdown-item" on:click={handleDelete(thread.id)}>Delete</button>
              </div>
            {/if}
          </div>
        </div>
        {#if thread.community_id}
          <div class="community">{thread.community_id}</div>
        {/if}
        <p class="text">{thread.content}</p>
        {#if thread.media && thread.media.length}
          <div class="image-preview images-{thread.media.length}">
            {#each thread.media as item (item.image_url)}
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
        <div class="actions">
          <button type="button" aria-label="Comment">
            <MessageCircle size="18" /><span>{thread.comment_count}</span>
          </button>
          <button type="button" aria-label="Repost">
            <Repeat size="18" /><span>{thread.share_count}</span>
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
  </style>