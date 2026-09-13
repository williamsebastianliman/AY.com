<script lang="ts">
    import { onMount } from "svelte";
    import { Search, X } from "lucide-svelte";
    import api from "../api";
  
    export let myUserId = "";
  
    let open = false;
  
    let trending: { hashtag: string; count: number }[] = [];
    let profiles: {
      id: string;
      name: string;
      username: string;
      profile_picture_id: string;
      profile_picture_url: string;
      count: number;
    }[] = [];
  
    async function loadTrending() {
      try {
        const res = await api.get("top-hashtag?limit=5");
        trending = res.data;
      } catch (err) {
        trending = [];
      }
    }
  
    async function fetchTopProfile() {
      try {
        const { data } = await api.get("/user/explore?query=#&page=1&size=3&threshold=100");
        profiles = await Promise.all(
          data.map(async (p) => {
            let profile_picture_url = "";
            if (p.profile_picture_id) {
              try {
                const resp = await api.post("/media/get-media", { id: p.profile_picture_id });
                profile_picture_url = resp.data.public_url || "";
              } catch {
                profile_picture_url = "";
              }
            }
            let count = 0;
            try {
              const resp = await api.get(`/user/follower_count/${p.id}`);
              count = resp.data.count;
            } catch {
              count = 0;
            }
            return {
              ...p,
              profile_picture_url,
              count,
            };
          })
        );
      } catch {
        profiles = [];
      }
    }
  
    function navigateToProfile(id: string) {
      window.location.href = `/profile/${id}`;
    }
  
    function goPremium() {
      window.location.href = "/premium";
    }
  
    onMount(() => {
      loadTrending();
      fetchTopProfile();
    });
  
    function openSidebar() {
      open = true;
      document.body.style.overflow = "hidden";
    }
    function closeSidebar() {
      open = false;
      document.body.style.overflow = "";
    }
  </script>
  
  <!-- Burger button, only visible below 1200px -->
  <button
    class="burger-btn"
    aria-label="Open right sidebar"
    on:click={openSidebar}
  >
    <svg width="28" height="28" viewBox="0 0 24 24" fill="none">
      <rect y="4" width="24" height="2.5" rx="1" fill="#1d9bf0"/>
      <rect y="10.75" width="24" height="2.5" rx="1" fill="#1d9bf0"/>
      <rect y="17.5" width="24" height="2.5" rx="1" fill="#1d9bf0"/>
    </svg>
  </button>
  
  {#if open}
    <div class="sidebar-backdrop" on:click={closeSidebar}></div>
    <aside class="burger-right-sidebar" in:slide={{ x: 300 }}>
      <button class="close-btn" on:click={closeSidebar} aria-label="Close sidebar">
        <X size="28" />
      </button>
      <div class="burger-content">
        <div class="ad-card">
          <h2>Go Premium</h2>
          <p>Unlock exclusive features!</p>
          <button type="button" class="ad-btn" on:click={goPremium}>
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
            <div class="suggestion-item" on:click={() => navigateToProfile(p.id)}>
              <img src={p.profile_picture_url} alt="avatar" class="suggestion-avatar" />
              <div class="suggestion-text">
                <div class="suggestion-name">{p.name}</div>
                <div class="suggestion-handle">Follower: {p.count}</div>
              </div>
            </div>
          {/each}
        </div>
      </div>
    </aside>
  {/if}
  
  <style>
    .burger-btn {
      display: none;
      position: fixed;
      right: 2vw;
      top: 2vw;
      z-index: 11001;
      background: #222d33;
      border: none;
      padding: 0.45rem 0.6rem;
      border-radius: 10px;
      color: #1d9bf0;
      box-shadow: 0 2px 8px #0003;
      transition: background .14s;
    }
    .burger-btn:hover { background: #1a1a1a; }
    @media (max-width: 1200px) {
      .burger-btn {
        display: block;
      }
    }
    .sidebar-backdrop {
      display: block;
      position: fixed;
      z-index: 11000;
      top: 0; left: 0; right: 0; bottom: 0;
      background: rgba(0,0,0,0.45);
      animation: fadeIn 0.18s;
    }
    @keyframes fadeIn {
      from { opacity: 0; }
      to   { opacity: 1; }
    }
    .burger-right-sidebar {
      position: fixed;
      right: 0; top: 0; bottom: 0;
      width: 350px;
      max-width: 88vw;
      background: #161b22;
      box-shadow: -2px 0 24px #0006;
      z-index: 11002;
      padding: 1.5rem 1rem 1rem 1rem;
      display: flex;
      flex-direction: column;
      overflow-y: auto;
      animation: slideInRight .2s;
    }
    @keyframes slideInRight {
      from { transform: translateX(350px); }
      to { transform: translateX(0); }
    }
    .close-btn {
      position: absolute;
      right: 1rem;
      top: 1rem;
      background: none;
      border: none;
      color: #1d9bf0;
      font-size: 2rem;
      z-index: 11003;
      cursor: pointer;
      padding: 4px;
      border-radius: 50%;
      transition: background .15s;
    }
    .close-btn:hover { background: #202327; }
    .burger-content {
      margin-top: 2.5rem;
      display: flex;
      flex-direction: column;
      gap: 1.5rem;
    }
    .ad-card {
      padding: 1rem;
      background: #1a1a1a;
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
      background: #1d9bf0;
      color: #fff;
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
      background: rgba(29,155,240,0.08);
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
      background: rgba(29,155,240,0.07);
    }
    .suggestion-avatar {
      width: 32px;
      height: 32px;
      border-radius: 50%;
      margin-right: 0.5rem;
      object-fit: cover;
    }
    .suggestion-text {
      flex-grow: 1;
    }
  </style>