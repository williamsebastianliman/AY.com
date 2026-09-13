<script lang="ts">
    import api from "../lib/api";
    import { onMount } from "svelte";
    import LeftSidebar from "../lib/components/LeftSideBar.svelte";
    import RightSidebar from "../lib/components/RightSideBar.svelte";
    import BurgerLeftSideBar from "../lib/components/BurgerLeftSideBar.svelte";
    import BurgerRightSideBar from "../lib/components/BurgerRightSideBar.svelte";
    import Thread from "../lib/components/Thread.svelte";
    import {Search} from "lucide-svelte";
    import { debounce } from "lodash";
    import {
      enrichThread,
      type RawThread,
    } from "../lib/threadUtil";
    interface RawThread {
      id: string;
      user_id: string;
      community_id?: string;
      content: string;
      like_count: number;
      comment_count: number;
      share_count: number;
      view_count: number;
      created_at: {
        year: number;
        month: number;
        day: number;
        hour: number;
        minute: number;
        second: number;
        timezone: string;
      };
      updated_at: {
        year: number;
        month: number;
        day: number;
        hour: number;
        minute: number;
        second: number;
        timezone: string;
      };
    }
    interface Thread {
      id: string;
      author: string;
      username: string;
      avatar: string;
      content: string;
      community_id?: string;
      like_count: number;
      comment_count: number;
      share_count: number;
      view_count: number;
      created_at: string;
      updated_at: string;
      media: string[];
      liked_by_me?: boolean;
    }
    interface MeResponse {
      id: string;
      name: string;
      username: string;
      profile_picture_id: string;
    }
    let threads: Thread[] = [];
    let loading = false;
    let loadError = false;
    let myUserId = "";
    let search = "";
    const debouncedFetch = debounce(() => loadBookmarks(true), 300);

    $: if (myUserId) {
      debouncedFetch();
    }
    async function onSearchInput() {
      debouncedFetch();
    }
    async function loadMyProfile() {
      try {
        const meRes = await api.get<MeResponse>("/user/get-me");
        myUserId = meRes.data.id;
      } catch (err) {
        console.error("Failed to load my profile:", err);
      }
    }
    async function loadBookmarks() {
      loading = true;
      loadError = false;
      try {
        const url = `/users/${myUserId}/threads-bookmark?query=${search}`;
        const res = await api.get<RawThread[]>(url);
        const raws = res.data;
        console.log("love uid: ", myUserId);
        const enriched = await Promise.all(raws.map((r) => enrichThread(r, myUserId)));
        threads = enriched;
      } catch (e) {
        loadError = true;
        console.error("Failed to load bookmarks:", e);
      } finally {
        loading = false;
      }
    }
    function handleSearch(e: Event) {
      e.preventDefault();

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
      await loadBookmarks();
    });
  </script>
  <div class="bookmarks layout">
    <aside class="sidebar">
      {#if myUserId && windowWidth >= 1250}
        <LeftSidebar currentUserId={myUserId} activePage="bookmarks" />
      {/if}
    </aside>
    {#if myUserId && windowWidth < 1250}
      <BurgerLeftSideBar currentUserId={myUserId} activePage="bookmarks" />
    {/if}
    <main class="main">
      <div class="header-bookmarks">
        <h2>Bookmarks</h2>
        <form class="search-bar" on:submit={handleSearch}>
          <Search size="18" class="search-icon" />
          <input
            type="text"
            placeholder="Search Bookmarks"
            bind:value={search}
            on:input={onSearchInput}
            class="search-input"
          />
        </form>
      </div>
      <section>
        {#if loading}
          {#each Array.from({ length: 3 }, (_, i) => i) as i (i)}
            <div class="skeleton skeleton-thread" />
          {/each}
        {:else}
          {#if threads.length === 0}
            <div class="no-bookmarks"><p>No bookmarks yet.</p></div>
          {:else if loadError}
          <div class="error"><p>Failed to load bookmarks.</p>
            <button type="button" class="retry-button" on:click={loadBookmarks}>Retry</button>
          </div>
          {:else}
            {#each threads as thread (thread.id)}
                <Thread {thread} currentUserId={myUserId} on:deleted={(e) => {
                  threads = threads.filter((t) => t.id !== e.detail.id);
                }} />
            {/each}
          {/if}
        {/if}
      </section>
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
    .layout {
      display: grid;
      grid-template-columns: auto 1fr auto;
      height: 100vh;
      overflow-y: hidden;
    }
    .main {
      overflow-y: auto;
      padding: 1rem;
      width: 40vw;
      max-width: 40vw;
      margin: 0 auto;
      border-left: 1px solid rgba(255,255,255,0.06);
      border-right: 1px solid rgba(255,255,255,0.06);
    }
    .header-bookmarks {
      display: flex;
      flex-direction: column;
      gap: 0.75rem;
      margin-bottom: 1.2rem;
      margin-top: 1rem;
    }
    .header-bookmarks h2 {
      font-size: 1.7rem;
      font-weight: 700;
      margin: 0;
      color: #fff;
    }
    .search-bar {
      display: flex;
      align-items: center;
      background: #171717;
      border-radius: 9999px;
      padding: 0.3rem 1rem;
      border: 1px solid rgba(255,255,255,0.07);
      margin-top: 0.4rem;
    }
    .search-input {
      border: none;
      outline: none;
      background: transparent;
      color: #f9f9f9;
      padding: 0.4rem;
      flex: 1;
      font-size: 1.02rem;
    }
    .search-icon {
      margin-right: 0.6rem;
      color: #888;
    }
    .skeleton {
      border-radius: 1rem;
      margin-bottom: 1rem;
      height: 120px;
      background: linear-gradient(90deg,rgba(40,40,40,0.4) 25%,rgba(60,60,60,0.4) 50%,rgba(40,40,40,0.4) 75%);
      background-size: 200% 100%;
      animation: shimmer 1.5s infinite;
    }

    @keyframes shimmer {
      0%   { background-position: 200% 0 }
      100% { background-position: -200% 0 }
    }
    .loading, .end, .no-bookmarks {
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

    @media (max-width: 900px) {
    .layout {
      min-height: 100vh;
      height: auto;
      width: 90vw;
      padding: 0;
      margin: 0 auto;
      display: grid;
      justify-content: center;
    }
    .main {
      width: 90vw;
      max-width: 90vw;
      margin: 0 auto;
      padding: 0.7rem 0.2rem 1.5rem 0.2rem;
      border-left: none;
      border-right: none;
      border-radius: 0;
      display: flex;
      flex-direction: column;
      align-items: center;
    }
    .header-bookmarks,
    section {
      width: 100%;
      max-width: 600px;
      margin: 0 auto;
    }
    .header-bookmarks h2 {
      font-size: 1.2rem;
    }
    .search-bar {
      padding: 0.2rem 0.7rem;
      font-size: 0.98rem;
    }
  }
</style>