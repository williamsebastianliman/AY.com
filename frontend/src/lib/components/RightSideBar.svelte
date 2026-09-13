<script lang="ts">
  import api from "../api";
  import { onMount } from "svelte";
  import { clickOutside } from "../action/clickOutside";
  import { Search, X } from "lucide-svelte";

  export let myUserId = "";
  let searchQuery = "";

  let trending: { hashtag: string; count: number }[] = [];
  async function loadTrending() {
    try{
      const res = await api.get(`top-hashtag?limit=${5}`);
      trending = res.data;
      console.log(trending.hashtag);
    }catch(err){
      console.error(err as string);
    }
  }


  let showSuggestions = false;
  let recentSearches: string[] = [];
  function loadRecentSearches() {
    const s = localStorage.getItem("recentSearches");
    recentSearches = s ? JSON.parse(s) : [];
  }
  function saveRecentSearches() {
    localStorage.setItem("recentSearches", JSON.stringify(recentSearches.slice(0, 3)));
  }
  function addToRecentSearches(q: string) {
    q = q.trim();
    if (!q) return;
    recentSearches = [q, ...recentSearches.filter(s => s !== q)].slice(0, 3);
    saveRecentSearches();
  }
  function clearRecentSearches() {
    recentSearches = [];
    saveRecentSearches();
  }

  function handleInput(e: Event) {
    searchQuery = (e.target as HTMLInputElement).value;
    showSuggestions = !!searchQuery.trim();
  }
  function handleFocus() {
    showSuggestions = true;
    loadRecentSearches();
  }

  function submitSearch() {
    const q = searchQuery.trim();
    if (!q) return;
    addToRecentSearches(q);
    window.location.href = `/explore?q=${encodeURIComponent(q)}`;
    showSuggestions = false;
  }
  function handleSuggestionClick(q: string) {
    searchQuery = q;
    addToRecentSearches(q);
    window.location.href = `/explore?q=${encodeURIComponent(q)}`;
    showSuggestions = false;
  }

  export interface Profile {
    id: string;
    name: string;
    username: string;
    profile_picture_id: string;
    profile_picture_url: string;
    count: number;
  }

  let profiles: Profile[] = [];
  async function fetchTopProfile(): Promise<Profile[]> {
      try {
        const { data } = await api.get(
          "/user/explore?query=#&page=1&size=3&threshold=100"
        );
        let profiles: Profile[] = await Promise.all(
          data.map(async (p) => {
            let profile_picture_url = "";
            if (p.profile_picture_id) {
              try {
                const resp = await api.post("/media/get-media", { id: p.profile_picture_id });
                profile_picture_url = resp.data.public_url || "";
              } catch (e) {
                console.error(e as string);
                profile_picture_url = "";
              }
            }
            let count = 0;
            try {
              const resp = await api.get(`/user/follower_count/${p.id}`);
              count = resp.data.count;
            } catch (e) {
              console.error(e as string);
              count = 0;
            }
            return {
              ...p,
              profile_picture_url,
              count,
            };
          })
        );
        console.log("pp: ",profiles);
        profiles = profiles.slice(0, 3);
        return profiles;
      } catch (err) {
        console.error(err as string);
        return [];
      }
    }

  function navigateToProfile(id: string){
    window.location.href = `/profile/${id}`;
  }
  onMount(async ()=>{
    loadTrending();
    loadRecentSearches();
    profiles = await fetchTopProfile();
  });
</script>

<div class="sidebar">
  <div class="search-container" use:clickOutside={() => showSuggestions = false}>
    <div class="input-wrapper">
      <input
        class="search-input"
        type="text"
        bind:value={searchQuery}
        placeholder="Search"
        on:input={handleInput}
        on:focus={handleFocus}
        on:keydown={(e) => e.key === "Enter" && submitSearch()}
        autocomplete="off"
      />
      {#if searchQuery}
        <button type="button" class="icon-btn clear-btn" title="Clear" on:click={() => { searchQuery = ""; showSuggestions = false; }}>
          <X size="18" />
        </button>
      {/if}
      <button type="button" class="icon-btn search-btn" title="Search" on:click={submitSearch}>
        <Search size="20" />
      </button>
    </div>
    {#if showSuggestions}
      <div class="recent-dropdown">
        <div class="recent-header">
          <span>Recent Searches</span>
          <button type="button" class="clear-btn" on:click={clearRecentSearches}>Clear All</button>
        </div>
        {#if recentSearches.length === 0}
          <div class="no-recent">No recent searches</div>
        {:else}
          {#each recentSearches as item (item)}
            <div class="recent-item">
              <span on:click={() => handleSuggestionClick(item)}>{item}</span>
            </div>
          {/each}
        {/if}
      </div>
    {/if}
  </div>

  <div class="ad-card">
    <h2>Go Premium</h2>
    <p>Unlock exclusive features!</p>
    <button
      type="button"
      class="ad-btn"
      on:click={() => (window.location.href = "/premium")}
    >
      Verify Now!
    </button>
  </div>

  <div class="section">
    <h2>What's happening</h2>
    {#each trending as t (t.hashtag)}
      <div class="trend-item">
        <span class="trend-tag">#{t.hashtag}</span>
        <span class="trend-count">{t.count}</span>
      </div>
    {/each}
  </div>

  <div class="section">
    <h2>Who to follow</h2>
    {#each profiles as p (p.id)}
      <div class="suggestion-item" on:click={navigateToProfile(p.id)}>
        <img src={p.profile_picture_url} alt="avatar" class="suggestion-avatar" />
        <div class="suggestion-text">
          <div class="suggestion-name">{p.name}</div>
          <div class="suggestion-handle">Follower: {p.count}</div>
        </div>
      </div>
    {/each}
  </div>
</div>

<style>
  :global(*::-webkit-scrollbar) {
    display: none;
  }
  :global(*), :global(*::-webkit-scrollbar-thumb) {
    -ms-overflow-style: none;
    scrollbar-width: none;
  }

  :global(:root) {
    --sidebar-width: 350px;
    --sidebar-bg: #111;
    --border: #333;
    --dropdown-bg: #1a1a1a;
    --hover: rgba(29, 155, 240, 0.1);
    --card-bg: #1a1a1a;
    --primary: #1d9bf0;
    --on-primary: #fff;
  }

  .sidebar {
    width: calc(var(--sidebar-width) * 1.2);
    margin-top: 2vh;

    height: 100vh;
    overflow-y: auto;
    padding: 1rem;
    background: var(--sidebar-bg);
  }

  .search-container {
    position: relative;
    margin-bottom: 1rem;
  }

  .search-input {
    width: 80%;
    max-width: 280px;
    padding: 0.5rem 2.5rem 0.5rem 0.5rem;
    border: 1px solid var(--border);
    border-radius: 0.25rem;
    background: transparent;
    color: #eee;
  }

  .search-container {
    position: relative;
    margin: 2rem auto 1rem;
    max-width: 320px;
  }

  .search-btn {
    position: absolute;
    right: 0.5rem;
    top: 50%;
    transform: translateY(-50%);
    background: none;
    border: none;
    cursor: pointer;
    color: #888;
  }

  .recent-dropdown {
    position: absolute;
    top: 110%;
    left: 0;
    width: 100%;
    background: var(--dropdown-bg);
    border: 1px solid var(--border);
    border-radius: 0.25rem;
    z-index: 10;
  }

  .recent-header {
    display: flex;
    justify-content: space-between;
    padding: 0.5rem;
    font-weight: bold;
    color: #ccc;
  }

  .clear-btn {
    background: none;
    border: none;
    color: var(--primary);
    cursor: pointer;
  }

  .recent-item {
    padding: 0.5rem;
    cursor: pointer;
    color: #eee;
  }

  .recent-item:hover {
    background: var(--hover);
  }

  .ad-card {
    padding: 1rem;
    background: var(--card-bg);
    border-radius: 0.5rem;
    text-align: center;
    margin-bottom: 1rem;
    color: #eee;
  }

  .ad-btn {
    margin-top: 0.5rem;
    padding: 0.5rem 1rem;
    border: none;
    border-radius: 9999px;
    background: var(--primary);
    color: var(--on-primary);
    cursor: pointer;
  }

  .section {
    margin-bottom: 1rem;
    color: #eee;
  }

  .section h2 {
    margin-bottom: 0.5rem;
    font-size: 1.1rem;
    font-weight: bold;
  }

  .trend-item {
    display: flex;
    justify-content: space-between;
    padding: 0.5rem 0;
    cursor: pointer;
    color: #ccc;
  }

  .trend-item:hover {
    background: var(--hover);
  }

  .suggestion-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0.5rem 0;
    cursor: pointer;
    color: #eee;
  }

  .suggestion-item:hover {
    background: var(--hover);
  }

  .suggestion-avatar {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    margin-right: 0.5rem;
  }

  .suggestion-text {
    flex-grow: 1;
  }

  .follow-btn {
    padding: 0.25rem 0.75rem;
    border: none;
    border-radius: 9999px;
    background: var(--primary);
    color: var(--on-primary);
    cursor: pointer;
  }

  .search-container {
  position: relative;
  margin-bottom: 1rem;
  max-width: 320px;
}

.recent-dropdown {
  position: absolute;
  top: 110%;
  left: 0;
  width: 100%;
  background: var(--dropdown-bg);
  border: 1px solid var(--border);
  border-radius: 0.25rem;
  z-index: 100;
}
.recent-header {
  display: flex;
  justify-content: space-between;
  padding: 0.5rem;
  font-weight: bold;
  color: #ccc;
}
.clear-btn {
  background: none;
  border: none;
  color: var(--primary);
  cursor: pointer;
}
.recent-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.45rem 0.55rem;
  cursor: pointer;
  color: #eee;
  border-radius: 0.25rem;
  transition: background 0.14s;
}
.recent-item span { flex: 1; }
.recent-item:hover {
  background: var(--hover);
}
.no-recent {
  color: #888;
  padding: 1.1rem 0;
  text-align: center;
}

.input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.search-input {
  width: 100%;
  padding: 0.6rem 2.8rem 0.6rem 1.6rem;
  border: 1.5px solid #333;
  border-radius: 8px;
  background: #111;
  color: #fff;
  font-size: 1.1rem;
  outline: none;
  transition: border-color .18s;
}

.icon-btn {
  position: absolute;
  background: none;
  border: none;
  cursor: pointer;
  padding: 0;
  top: 50%;
  transform: translateY(-50%);
  display: flex;
  align-items: center;
  z-index: 1;
  border-radius: 50%;
  transition: background .15s;
}
.icon-btn:active, .icon-btn:focus {
  background: #1d9bf010;
}

.clear-btn {
  right: 2.3rem;
}

.search-btn {
  right: 0.6rem;
}


</style>