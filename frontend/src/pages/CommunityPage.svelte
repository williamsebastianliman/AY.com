<script lang="ts">
    import { onMount } from "svelte";
    import LeftSidebar from "../lib/components/LeftSideBar.svelte";
    import RightSidebar from "../lib/components/RightSideBar.svelte";
    import BurgerLeftSideBar from "../lib/components/BurgerLeftSideBar.svelte";
    import BurgerRightSideBar from "../lib/components/BurgerRightSideBar.svelte";
    import api from "../lib/api";
    import { fade } from "svelte/transition";
    import { Plus } from "lucide-svelte";
    type Tab = "joined" | "requested" | "discover";
    interface Category {
      category_id: string;
      category_name: string;
    }
    interface Community {
      community_id: string;
      community_name: string;
      community_description: string;
      community_logo?: string;
      logo_url?: string;
      status: string;
      categories: Category[];
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
    let activeTab: Tab = "joined";
    let communities: Community[] = [];
    let currentPage = 1;
    let perPage = 25;
    const perPageOptions = [25, 30, 35];
    let categories: string[] = [];
    let selectedCategory = "";
    let searchQuery = "";
    let loading = false;
    let error = "";
    const joinRequesting: Record<string, boolean> = {};
    let myUserId = "";
    async function fetchCategories() {
      try {
        const res = await api.get<{ categories: string[] }>("/community-categories");
        categories = res.data;
        console.log("cat: ",categories);
      } catch {
        categories = [];
      }
    }
    async function fetchCommunities() {
      loading = true;
      error = "";
      let url = "";
      const params: Record<string, string | number> = {
        page: currentPage,
        size: perPage,
      };
      if (activeTab === "joined") {
        url = `/users/communities/joined/${myUserId}`;
      } else if (activeTab === "requested") {
        url = `/users/communities/requests/${myUserId}`;
      } else {
        url = "/communities";
        params.user_id = myUserId;
        if (searchQuery) params.query = searchQuery;
        if (selectedCategory) params.category = selectedCategory;
      }
      try {
        const res = await api.get<{ communities: Community[]; total?: number }>(url, { params });
        communities = await Promise.all(
          res.data.map(async (c) => {
            let logo_url = "";
            if (c.community_logo) {
              try {
                const m = await api.post<{ public_url: string }>("/media/get-media", { id: c.community_logo });
                logo_url = m.data.public_url;
              } catch (error) {
                console.error(error as string);
              }
            }
            let categories: Category[] = [];
            try {
              const catsRes = await api.get<{ categories: Category[] }>(`/communities/categories/${c.community_id}`);
              categories = catsRes.data;
            } catch (error) {
              console.error("Failed to fetch categories for community:", c.community_id, error);
            }
            return { ...c, logo_url, categories };
          })
        );
        console.log("comz: ", communities);
      } catch (e) {
        error = "Failed to load communities";
        console.error(e as string);
      }
      loading = false;
    }

    function switchTab(tab: Tab) {
      if (tab !== activeTab) {
        activeTab = tab;
        currentPage = 1;
        fetchCommunities();
      }
    }
    function changePerPage(n: number) {
      perPage = n;
      currentPage = 1;
      fetchCommunities();
    }
    function goToCommunity(id: string) {
      window.location.href = `/community/${id}`;
    }
    async function loadMyProfile()
    {
      try {
        const meRes = await api.get<MeResponse>("/user/get-me");
        myUserId = meRes.data.id;
        console.log("user idss: ", myUserId);
      } catch (err) {
        console.error("Failed to load my profile:", err);
      }
    }
    async function handleJoinRequest(id: string) {
      joinRequesting[id] = true;
      try {
        await api.post(`/communities/join/${id}`,{user_id: myUserId});
        fetchCommunities();
      } catch(error) {
        console.error(error as string);
      }
      joinRequesting[id] = false;
    }
    function handleSearch(e: Event) {
      e.preventDefault();
      currentPage = 1;
      fetchCommunities();
    }
    onMount(async () => {
      await fetchCategories();
      await loadMyProfile();
      fetchCommunities();
    });
    $: if (activeTab === "discover" && (selectedCategory || searchQuery)) {
      currentPage = 1;
      fetchCommunities();
    }
    let windowWidth = window.innerWidth;
    function handleResize() {
      windowWidth = window.innerWidth;
    }
    onMount(() => {
      window.addEventListener("resize", handleResize);
      return () => window.removeEventListener("resize", handleResize);
    });
  </script>
  <div class="communities layout">
    <aside class="sidebar">
      {#if myUserId && windowWidth >= 1250}
        <LeftSidebar currentUserId={myUserId} activePage="communities" />
      {/if}
    </aside>
    {#if myUserId && windowWidth < 1250}
      <BurgerLeftSideBar currentUserId={myUserId} activePage="communities" />
    {/if}
    <main class="main">
      <div class="communities-header">
        <h1>Communities</h1>
        <button type="button" class="create-btn" on:click={() => window.location.href = "/create-community"}>
          <Plus size={20} /> Create Community
        </button>
      </div>
      <nav class="tabs">
        <button type="button" class:active={activeTab === "joined"} on:click={() => switchTab("joined")}>Joined</button>
        <button type="button" class:active={activeTab === "requested"} on:click={() => switchTab("requested")}>Requested</button>
        <button type="button" class:active={activeTab === "discover"} on:click={() => switchTab("discover")}>Discover</button>
      </nav>
      <section class="toolbar" in:fade>
        {#if activeTab === "discover"}
          <form class="search-bar" on:submit={handleSearch}>
            <input
              type="search"
              placeholder="Search communities…"
              bind:value={searchQuery}
              class="search-input"
            />
            <button type="submit" class="search-btn">Search</button>
          </form>
          <select bind:value={selectedCategory} class="category-filter" on:change={fetchCommunities}>
            <option value="">All Categories</option>
            {#each categories as c (c.category_id)}
              <option value={c.category_id}>{c.category_name}</option>
            {/each}
          </select>
          <div class="pagination-controls">
            <button
              type="button"
              on:click={() => { if (currentPage > 1) { currentPage--; fetchCommunities(); }}}
              disabled={currentPage <= 1}
            >Prev</button>
            <span>Page {currentPage}</span>
            <button
              type="button"
              on:click={() => { if (communities.length === perPage) { currentPage++; fetchCommunities(); }}}
              disabled={communities.length < perPage}
            >Next</button>
          </div>
        {:else}
          <div class="pagination-controls">
            <button
              type="button"
              on:click={() => { if (currentPage > 1) { currentPage--; fetchCommunities(); }}}
              disabled={currentPage <= 1}
            >Prev</button>
            <span>Page {currentPage}</span>
            <button
              type="button"
              on:click={() => { if (communities.length === perPage) { currentPage++; fetchCommunities(); }}}
              disabled={communities.length < perPage}
            >Next</button>
          </div>
        {/if}
        <div class="pagination">
            <span>Show</span>
            {#each perPageOptions as n (n)}
              <button
                type="button"
                class:active={perPage === n}
                on:click={() => changePerPage(n)}
              >{n}</button>
            {/each}
            <span>per page</span>
        </div>
      </section>
      <section class="community-list" in:fade>
        {#if loading}
          <div class="loading">Loading…</div>
        {:else if error}
          <div class="error">{error}</div>
        {:else if communities.length === 0}
          <div class="empty">No communities found.</div>
        {:else}
          {#each communities as c (c.community_id)}
            <div class="community-entry" on:click={() => goToCommunity(c.community_id)} tabindex="0">
              <img class="community-logo" src={c.logo_url || "/default-community.png"} alt="logo" />
              <div class="community-info">
                <div class="community-title">{c.community_name}</div>
                <div class="community-desc">{c.community_description}</div>
                <div class="community-categories">
                  {#each c.categories as cat (cat.category_id)}
                    <span class="category">{cat.category_name}</span>
                  {/each}
                </div>
              </div>
              {#if activeTab === "discover"}
                <button
                  class="join-btn"
                  type="button"
                  on:click|stopPropagation={() => handleJoinRequest(c.community_id)}
                  disabled={joinRequesting[c.id]}
                >
                  {joinRequesting[c.id] ? "Requesting…" : "Request to Join"}
                </button>
              {/if}
              {#if activeTab === "requested"}
                <span class="pending-label">Pending</span>
              {/if}
              {#if activeTab === "joined"}
                <span class="joined-label">Joined</span>
              {/if}
            </div>
          {/each}
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
    .communities.layout {
      display: grid;
      grid-template-columns: auto 1fr auto;
      /* min-height: 100vh; */
      overflow-y: hidden;
    }

.main {
  overflow-y: auto;
  padding: 1.5rem 1rem 2rem 1rem;
  width: 45vw;
  max-width: 45vw;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  min-height: 0;
  height: 100vh;
  background: none;
  -webkit-overflow-scrolling: touch;
}

.communities-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.75rem;
  flex-wrap: wrap;
  gap: 1rem;
}

.create-btn {
  background: #1d9bf0;
  color: #fff;
  border: none;
  border-radius: 9999px;
  padding: 0.5rem 1.4rem;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 0.7rem;
  cursor: pointer;
  font-size: 1.01rem;
  box-shadow: 0 2px 7px #1d9bf031;
}

.tabs {
  display: flex;
  border-bottom: 1.3px solid #313950;
  margin-bottom: 1.2rem;
  flex-wrap: wrap;
}
.tabs button {
  background: none;
  border: none;
  flex: 1;
  padding: 0.7rem 0.2rem;
  font-size: 1rem;
  font-weight: 500;
  color: #abb6c4;
  cursor: pointer;
  border-bottom: 2.3px solid transparent;
  transition: color 0.12s, border-color 0.12s;
}
.tabs button.active {
  color: #fff;
  border-bottom: 2.5px solid #1d9bf0;
}

.toolbar {
  display: flex;
  align-items: center;
  gap: 1.5rem;
  margin-bottom: 1.5rem;
  flex-wrap: wrap;
}
.search-bar {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}
.search-input {
  padding: 0.45rem 0.95rem;
  border-radius: 0.7rem;
  border: 1px solid #285a97;
  background: #161a1d;
  color: #f7f7f7;
  font-size: 1rem;
  min-width: 140px;
}
.search-btn {
  background: #1d9bf0;
  color: white;
  border: none;
  border-radius: 0.7rem;
  padding: 0.45rem 1.2rem;
  font-weight: 500;
  cursor: pointer;
}
.category-filter {
  padding: 0.45rem 0.9rem;
  border-radius: 0.7rem;
  border: 1px solid #285a97;
  background: #161a1d;
  color: #f7f7f7;
  font-size: 1rem;
}

.pagination, .pagination-controls {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  flex-wrap: wrap;
}
.pagination {
  margin-left: auto;
}
.pagination button {
  background: none;
  border: none;
  color: #1d9bf0;
  font-weight: 600;
  cursor: pointer;
  border-radius: 0.5rem;
  padding: 0.35rem 0.9rem;
}
.pagination button.active {
  background: #1d9bf0;
  color: #fff;
}
.pagination-controls button {
  background: #191f27;
  color: #1d9bf0;
  border: none;
  border-radius: 0.6rem;
  padding: 0.45rem 1.3rem;
  font-weight: 600;
  font-size: 1.01rem;
  cursor: pointer;
}
.pagination-controls button:disabled {
  color: #555;
  background: #13181f;
  cursor: not-allowed;
}

.community-list {
  display: flex;
  flex-direction: column;
  gap: 1.2rem;
  overflow-y: auto;
  flex: 1;
  min-height: 0;
  padding: 2vh 0;
}
.community-entry {
  display: flex;
  align-items: flex-start;
  gap: 1.4rem;
  padding: 1.1rem 1.6rem;
  background: #161a1e;
  border-radius: 1.2rem;
  cursor: pointer;
  box-shadow: 0 1.5px 8px #1d9bf012;
  transition: background 0.11s;
  border: 1.5px solid rgba(29,155,240,0.07);
  position: relative;
}
.community-entry:hover {
  background: #191f27;
}
.community-logo {
  width: 56px;
  height: 56px;
  border-radius: 2rem;
  object-fit: cover;
  background: #222;
  border: 2px solid #1d9bf0;
}
.community-info {
  flex: 1;
}
.community-title {
  font-size: 1.19rem;
  font-weight: 700;
  color: #e0e8ff;
  margin-bottom: 0.2rem;
}
.community-desc {
  color: #b2c2d8;
  font-size: 1.02rem;
  margin-bottom: 0.6rem;
}
.community-categories {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}
.category {
  background: #232a39;
  color: #5bc0f8;
  padding: 0.22rem 0.75rem;
  border-radius: 0.5rem;
  font-size: 0.91rem;
  font-weight: 500;
}
.join-btn {
  margin-top: 1vh;
  background: #1d9bf0;
  color: #fff;
  border: none;
  border-radius: 0.8rem;
  padding: 0.5rem 1.2rem;
  font-weight: 600;
  cursor: pointer;
  margin-left: 1rem;
}
.joined-label, .pending-label {
  background: #29e57b;
  color: #fff;
  border-radius: 0.8rem;
  padding: 0.5rem 1.2rem;
  font-weight: 700;
  margin-left: 1.3rem;
  margin-top: 1vh;
}
.pending-label {
  background: #ffd070;
  color: #444;
}
.empty, .loading, .error {
  text-align: center;
  color: #fff;
  padding: 2rem 0;
}

@media (max-width: 1200px) {
  .main {
    width: 98vw;
    max-width: 98vw;
    padding: 1.2rem 0.4rem 2.5rem 0.4rem;
  }
  .community-entry {
    padding: 0.7rem 0.7rem;
    gap: 1rem;
  }
  .communities.layout {
    grid-template-columns: none;
    max-height: 90vh;
    width: 100vw;
    margin: 0;
    padding: 0;
    background: #101522;
    overflow-x: hidden;
    overflow-y: auto;
  }
  .main {
    width: 100vw;
    max-width: 100vw;
    padding: 1rem 0.1rem 2.2rem 0.1rem;
    margin: 0;
    border-radius: 0;
    box-shadow: none;
    display: flex;
    flex-direction: column;
    align-items: center;
    max-height: 90vh;
    overflow-y: auto;
    height: 90vh;
  }
  .communities-header,
  .tabs,
  .toolbar {
    flex-direction: column;
    align-items: stretch;
    gap: 0.5rem;
  }
  .community-entry {
    padding: 0.6rem 0.4rem;
    gap: 0.5rem;
  }
  .community-logo {
    width: 36px;
    height: 36px;
  }
  .community-title {
    font-size: 1rem;
  }
  .category {
    font-size: 0.85rem;
    padding: 0.18rem 0.55rem;
  }
}

@media (max-width: 500px) {
  .main {
    width: 100vw;
    max-width: 100vw;
    padding: 0.5rem 0.01rem 1.3rem 0.01rem;
  }
  .community-entry {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.5rem;
    padding: 0.5rem 0.2rem;
    border-radius: 0.9rem;
  }
  .community-logo {
    width: 30px;
    height: 30px;
  }
  .community-title {
    font-size: 0.98rem;
  }
  .category {
    font-size: 0.82rem;
    padding: 0.11rem 0.38rem;
  }
  .pagination-controls,
  .pagination {
    gap: 0.4rem;
    font-size: 0.91rem;
  }
}


</style>