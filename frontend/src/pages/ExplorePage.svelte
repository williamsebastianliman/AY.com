<script lang="ts">
    import { onMount } from "svelte";
    import LeftSidebar from "../lib/components/LeftSideBar.svelte";
    import BurgerLeftSideBar from "../lib/components/BurgerLeftSideBar.svelte";
    import RightSidebar from "../lib/components/RightSideBar.svelte";
    import BurgerRightSideBar from "../lib/components/BurgerRightSideBar.svelte";
    import Thread from "../lib/components/Thread.svelte";
    import { Search, X, Hash, Loader2, ArrowLeft } from "lucide-svelte";
    import { clickOutside } from "../lib/action/clickOutside";
    import api from "../lib/api";
    let myUserId = "";
    let inputQuery = "";
    let searchQuery = "";
    let searchLoading = false;
    let searchTab: "top"|"latest"|"people"|"media"|"communities" = "top";
    let showSuggestions = false;
    let recentSearches: string[] = [];
    let peopleFilter: "everyone"|"following"|"verified" = "everyone";
    let selectedCategory = "";
    let categories: string[] = [];
    let trendingHashtags: {hashtag: string; count: number}[] = [];
    let topHashtags: {tag: string; count: number}[] = [];
    let pagePeople = 1;
    let perPagePeople = 25;
    let totalPeople = 0;
    let pageCommunities = 1;
    let perPageCommunities = 25;
    let totalCommunities = 0;
    let recommendedProfiles: any[] = [];

    let curtain: boolean = true;
    export interface Profile {
      id: string;
      name: string;
      username: string;
      profile_picture_id: string;
      profile_picture_url: string;
      count: number;
      is_followed: boolean;
    }

      import {
      enrichThread,
      type RawThread,
    } from "../lib/threadUtil";

    let trendingThreads: Thread[] = [];
    let results: any = {
      people: [],
      threads: [],
      latest: [],
      media: [],
      communities: [],
    };
    let trendingTab: boolean = false;
    let debounceTimeout: ReturnType<typeof setTimeout>;
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

    async function fetchProfileRecommendation(): Promise<Profile[]> {
      try {
        const { data } = await api.get(`/user/explore?query=${encodeURIComponent(searchQuery)}&page=1&size=3&threshold=3`);
        const profiles: Profile[] = await Promise.all(
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
            return {
              ...p,
              profile_picture_url,
            };
          })
        );
        return profiles;
      } catch (err) {
        console.error(err as string);
        return [];
      }
    }

    let recommendedProfile: Profile[] = [];
    function handleInput(e: Event) {
      inputQuery = (e.target as HTMLInputElement).value;
      showSuggestions = !!inputQuery.trim();
      clearTimeout(debounceTimeout);
      debounceTimeout = setTimeout(async () => {
        if (inputQuery.trim()) {
          searchQuery = inputQuery;
          recommendedProfile = await fetchProfileRecommendation();
        }
      }, 320);
    }
    function submitSearch() {
        console.log("sq: ",searchQuery);
        console.log("iq: ", inputQuery);
        searchQuery = inputQuery;
        if(searchQuery == "")
        {
          curtain = true;
        }
        else{
          curtain=false;
        }
        addToRecentSearches(inputQuery);
        showSuggestions = false;
        doSearch();
    }
    function handleSuggestionClick(q: string) {
      inputQuery = q;
      searchQuery = q;
      addToRecentSearches(q);
      doSearch();
      curtain=false;
      showSuggestions = false;
    }
    function setTab(t: typeof searchTab) { searchTab = t; doSearch(); }
    function setPeopleFilter(f: typeof peopleFilter) { peopleFilter = f; doSearch();}
    function setCategory(cat: string) { selectedCategory = cat; doSearch();}
    async function handleHashtagClick(tag: string) {
      addToRecentSearches(tag);
      doSearch();
      showSuggestions = false;
      await loadTrendingThreads(tag);
      trendingTab = true;
    }
    let topThreeProfile: Profile[] = [];
    let peopleProfile: Profile[] = [];
    let relatedMedia: SearchMediaGroupedMediaJSON[] = [];
    let relatedCommunities: Community[] = [];
    async function doSearch() {
        if(curtain == false){
          searchLoading = true;
          searchLoading = false;
          if(searchTab == "top"){
            topThreeProfile =  await fetchTopProfile();
            await fetchTopThreads();
          }
          else if(searchTab == "latest")
          {
            await fetchLatestThreads();
            console.log("datazz: ", latestThreads);
          }
          else if(searchTab == "people")
          {
            peopleProfile = await topProfilePaginate(pagePeople,perPagePeople);
          }
          else if(searchTab == "media")
          {
            const data = await fetchRelatedMedia(searchQuery, 3, 1, 12);
            if(data)
            {
              relatedMedia = data.groups;
            }
          }
          else if(searchTab == "communities")
          {
            const data = await fetchCommunityPage(searchQuery,3, pageCommunities, perPageCommunities);
            relatedCommunities = data;
          }
        }

    }
    async function fetchTrendingHashtags() {
      const data = (await api.get("/top-hashtag")).data;
      if (data!= null){
        topHashtags = data;
      }
    }
    async function loadMyProfile() {
      try {
        const meRes = await api.get("/user/get-me");
        myUserId = meRes.data.id;
      } catch (err) {
        console.error("Failed to load my profile:", err);
      }
    }
    async function loadTrendingThreads(hashtag:string) {
      try{
        trendingThreads = [];
        const res = await api.get(`/threads/hashtag?hashtag=${hashtag}&size=${1000}`);
        const enriched = await Promise.all(res.data.map((r) => enrichThread(r, myUserId)));
        trendingThreads = [...trendingThreads, ...enriched];
      }catch(err){
        console.error(err as string);
      }
    }
    async function fetchTopProfile(): Promise<Profile[]> {
      try {
        const { data } = await api.get(`/user/explore?query=${encodeURIComponent(searchQuery)}&page=1&size=3&threshold=3`);
        const profiles: Profile[] = await Promise.all(
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
            return {
              ...p,
              profile_picture_url,
            };
          })
        );
        const profiles_0: Profile[] = await Promise.all(
          profiles.map(async (p) => {
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
              count,
            };
          })
        );
        return profiles_0;
      } catch (err) {
        console.error(err as string);
        return [];
      }
    }
    let topThreads: Thread[] = [];

    async function fetchTopThreads() {
      topThreads = [];
      try {
        const resp = await api.get<RawThread[]>(`/threads/explore?keyword=${encodeURIComponent(searchQuery)}&page=1&size=1000&threshold=3`);
        const rawThreads = resp.data.threads;
        const enriched = await Promise.all(rawThreads.map((r) => enrichThread(r, myUserId)));

        topThreads = enriched;
        console.log("tt: ",topThreads);
      } catch (err) {
        console.error("Fetch top threads failed:", err);
      }
    }
    let latestThreads: Thread[] = [];
    async function fetchLatestThreads() {
      latestThreads = [];
      try {
        const resp = await api.get<RawThread[]>(`/threads/explore?keyword=${encodeURIComponent(searchQuery)}&page=1&size=1000&threshold=3`);
        const rawThreads = resp.data.threads;
        const enriched = await Promise.all(rawThreads.map((r) => enrichThread(r, myUserId)));

        latestThreads = enriched;
        console.log("tt: ",topThreads);
      } catch (err) {
        console.error("Fetch top threads failed:", err);
      }
    }
    async function topProfilePaginate(page: number, size: number): Promise<Profile[]> {
      try {
        const { data } = await api.get(
          `/user/explore?query=${encodeURIComponent(searchQuery)}&page=${page}&size=${size}&threshold=3`
        );
        const profiles: Profile[] = await Promise.all(
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
            let is_followed = false;
            try {
              is_followed = await fetchIsFollowed(p.id, myUserId);
            } catch (e) {
              console.error(e as string);
              is_followed = false;
            }
            return {
              ...p,
              profile_picture_url,
              count,
              is_followed
            };
          })
        );
        console.log("pp: ",profiles);
        return profiles;
      } catch (err) {
        console.error(err as string);
        return [];
      }
    }


    function navigateToProfile(id: string){
      window.location.href = `/profile/${id}`;
    }
    function navigateToCommunityDetail(id: string){
      window.location.href = `/community/${id}`;
    }
    async function followUserAPI(targetId: string) {
      await api.post("/user/follow", {
        user_id: targetId,
        follower_id: myUserId
      });
    }

    async function unfollowUserAPI(targetId: string) {
        await api.post("/user/unfollow", {
          user_id: targetId,
          follower_id: myUserId
        });
    }
    async function fetchIsFollowed(targetUserId: string, currentUserId: string): Promise<boolean> {
      if (!targetUserId || !currentUserId || targetUserId === currentUserId) return false;
      try {
        const res = await api.post<{ is_followed: boolean }>("/user/is-followed", {
          user_id: targetUserId,
          follower_id: currentUserId,
        });
        return res.data.is_followed;
      } catch {
        return false;
      }
    }

    export interface Community {
      community_id: string;
      community_name: string;
      community_description: string;
      community_rules: string;
      status: string;
      created_at: {
        year: number;
        month: number;
        day: number;
        hour: number;
        minute: number;
        second: number;
        timezone: string;
      };
      community_logo: string;
      logo_url: string;
      community_banner: string;
    }

    async function fetchCommunityPage(query: string, threshold: number, page: number, size: number): Promise<Community[]> {
      try {
        const res = await api.get(
          `/community/explore?query=${encodeURIComponent(query)}&threshold=${threshold}&page=${page}&size=${size}`
        );
        const communities: Community[] = await Promise.all(
          (res.data as Community[]).map(async (c) => {
            let logo_url = "";
            if (c.community_logo) {
              try {
                const logoResp = await api.post("/media/get-media", { id: c.community_logo });
                logo_url = logoResp.data.public_url || "";
              } catch {
                logo_url = "";
              }
            }
            return { ...c, logo_url };
          })
        );
        return communities;
      } catch (err) {
        console.error("Failed to fetch community page:", err);
        return [];
      }
    }


    export interface TimeJSON {
      year: number;
      month: number;
      day: number;
      hour: number;
      minute: number;
      second: number;
      timezone: string;
    }

    export interface ThreadMediaJSON {
      thread_id: string;
      image_url: string;
      extension: string;
      created_at: TimeJSON;
    }

    export interface ThreadJSON {
      id: string;
      user_id: string;
      community_id: string;
      content: string;
      like_count: number;
      comment_count: number;
      share_count: number;
      view_count: number;
      created_at: TimeJSON;
      updated_at: TimeJSON;
      repost_id: string;
      is_pinned: boolean;
    }

    export interface SearchMediaGroupedMediaJSON {
      thread: ThreadJSON;
      media_list: ThreadMediaJSON[];
    }

    export interface SearchMediaGroupedResponse {
      groups: SearchMediaGroupedMediaJSON[];
      total: number;
    }

    function isVideo(extension: string) {
      return [".mp4", ".webm", ".mov", ".avi", ".mkv"].includes((extension || "").toLowerCase());
    }
    let videoDurations: Record<string, string> = {};
    function formatDuration(sec: number) {
      const m = Math.floor(sec / 60);
      const s = Math.floor(sec % 60);
      return `${m}:${s.toString().padStart(2, "0")}`;
    }
    function loadVideoDuration(media: { id: string; image_url: string; extension: string }) {
      if (!isVideo(media.extension)) return;
      if (videoDurations[media.id]) return;
      const video = document.createElement("video");
      video.src = media.image_url;
      video.preload = "metadata";
      video.muted = true;
      video.addEventListener("loadedmetadata", () => {
        videoDurations = {
          ...videoDurations,
          [media.id]: formatDuration(video.duration)
        };
      });
    }


    async function fetchRelatedMedia(
      keyword: string,
      threshold = 3,
      page = 1,
      size = 12
    ): Promise<SearchMediaGroupedResponse | null> {
      try {
        const params = new URLSearchParams({
          keyword,
          threshold: threshold.toString(),
          page: page.toString(),
          size: size.toString()
        }).toString();
        const resp = await api.get<SearchMediaGroupedResponse>(`/threads-media/explore?${params}`);
        return resp.data;
      } catch (err) {
        console.error("Fetch related media failed:", err);
        return null;
      }
    }

    $: if (pagePeople || perPagePeople) {
      doSearch();
    }

    $: if (pageCommunities || perPageCommunities) {
      doSearch();
    }


    onMount(async () => {
      inputQuery="";
      await loadMyProfile();
      await fetchTrendingHashtags();
      loadRecentSearches();
      const params = new URLSearchParams(window.location.search);
      const q = params.get("q")?.trim() || "";
      inputQuery = q;
      searchQuery = q;
      if (q) {
        curtain = false;
        await doSearch(q);
      }
    });

    let windowWidth = window.innerWidth;
    function handleResize() {
      windowWidth = window.innerWidth;
    }
    onMount(() => {
      window.addEventListener("resize", handleResize);
      return () => window.removeEventListener("resize", handleResize);
    });
    // onMount(() => {
    //   const input = document.querySelector("input[name='search']") as HTMLInputElement;
    //   if (input && input.value) {
    //     inputQuery = input.value;
    //   }
    // });
  </script>
  <div class="twitter-layout">
    <aside class="sidebar">
      {#if myUserId && windowWidth >= 1250}
        <LeftSidebar currentUserId={myUserId} activePage="explore" />
      {/if}
    </aside>
    {#if myUserId && windowWidth < 1250}
      <BurgerLeftSideBar currentUserId={myUserId} activePage="explore" />
    {/if}
    <main class="main-content">
      <div class="explore-header">
        <h1>Explore</h1>
      </div>
      <div class="search-bar-container" use:clickOutside on:outclick={() => showSuggestions = false}>
        <div class="search-bar">
          <input
            name="search"
            placeholder="Search…"
            bind:value={inputQuery}
            on:input={handleInput}
            on:focus={() => showSuggestions = true}
            on:keydown={(e) => e.key === "Enter" && submitSearch()}
          />
          <button type="button" class="search-btn" on:click={submitSearch}>
            <Search size="20" />
          </button>
          {#if inputQuery}
            <button type="button" class="clear-btn" on:click={() => { inputQuery = ""; showSuggestions = false; curtain=true;}}>
              <X size="18" />
            </button>
          {/if}
        </div>
        {#if showSuggestions}
          <div class="recent-search-dropdown">
            <div class="recent-header">
              <span>Recent Searches</span>
              <button type="button" class="clear-recent" on:click={clearRecentSearches}>Clear All</button>
            </div>
            {#if recentSearches.length === 0}
              <div class="no-recent">No recent searches</div>
            {:else}
              {#each recentSearches as item (item)}
                <div class="recent-item" on:click={() => handleSuggestionClick(item)}>{item}</div>
              {/each}
              <h3>Recommended Profile</h3>
              {#each recommendedProfile as profile (profile.id)}
                <div class="profile-suggestion" on:click={() => navigateToProfile(profile.id)}>
                  <img class="profile-pic" src={profile.profile_picture_url || "/assets/default-avatar.png"} alt="Profile Picture" />
                  <div class="profile-info">
                    <span class="profile-name">{profile.name}</span>
                    <span class="profile-username">@{profile.username}</span>
                  </div>
                </div>
              {/each}
            {/if}
          </div>
        {/if}
      </div>
      <div class="filters-row">
        <label>People:</label>
        <select bind:value={peopleFilter} on:change={(e) => setPeopleFilter(e.target.value)}>
          <option value="everyone">Everyone</option>
          <option value="following">People you follow</option>
          <option value="verified">Verified accounts only</option>
        </select>
        <label>Category:</label>
        <select bind:value={selectedCategory} on:change={(e) => setCategory(e.target.value)}>
          <option value="">All</option>
          {#each categories as cat (cat)}
            <option value={cat}>{cat}</option>
          {/each}
        </select>
      </div>
      {#if curtain}
        <div class="trending-section">
          {#if !trendingTab}
          <h2>Trending Hashtags</h2>
          <div class="trending-list">
            {#each topHashtags as tagObj (tagObj.hashtag)}
              <div class="trending-item" on:click={() => handleHashtagClick(tagObj.hashtag)}>
                <span class="hashtag"><Hash size="16" />{tagObj.hashtag}</span>
                <span class="count">{tagObj.count}</span>
              </div>
            {/each}
          </div>
          {:else}
            <div class="trending-threads-header">
              <button type="button" class="clear-btn trending-back-btn" on:click={() => {trendingTab = false;}}>
                <ArrowLeft size="20" />
              </button>
              <h2>Trending Threads</h2>
              <div class="trending-divider"></div>
            </div>
            <div class="trending-threads-scroll">
              {#if trendingThreads.length === 0}
                <div class="no-threads">No trending threads for this hashtag.</div>
              {:else}
                {#each trendingThreads as t (t.id)}
                  <Thread is_profile_thread={false} thread={t} currentUserId={myUserId}
                    on:deleted={(e) => {
                      t = t.filter((t) => t.id !== e.detail.id);
                    }} />
                {/each}
              {/if}
            </div>
          {/if}
        </div>
      {:else}
        <div class="nav-tabs">
          <button type="button" class:active={searchTab === "top"} on:click={() => setTab("top")}>Top</button>
          <button type="button" class:active={searchTab === "latest"} on:click={() => setTab("latest")}>Latest</button>
          <button type="button" class:active={searchTab === "people"} on:click={() => setTab("people")}>People</button>
          <button type="button" class:active={searchTab === "media"} on:click={() => setTab("media")}>Media</button>
          <button type="button" class:active={searchTab === "communities"} on:click={() => setTab("communities")}>Communities</button>
        </div>
        <div class="tab-content">
          {#if searchLoading}
            <div class="loading-center"><Loader2 size="24" class="spin" /></div>
          {:else if searchTab === "top"}
            <div class="search-top">
              <h3>Top People</h3>
              <div class="top-people-list">
                {#each topThreeProfile as profile (profile.id)}
                  <div class="user-card" on:click={() => navigateToProfile(profile.id)}>
                    <img class="avatar" src={profile.profile_picture_url} alt={profile.name} />
                    <div class="info">
                      <span class="name">{profile.name}</span>
                      <span class="username">@{profile.username}</span>
                      <span class="followers">{profile.count} followers</span>
                    </div>
                  </div>
                {/each}
                <button type="button" class="view-all-btn" on:click={() => setTab("people")}>View All</button>
              </div>
              <h3>Threads</h3>
              <div class="top-people-list">
                {#each topThreads as t (t.id)}
                  <Thread thread={t} currentUserId={myUserId} />
                {/each}
              </div>
            </div>
          {:else if searchTab === "latest"}
            <div class="thread-list">
              {#each latestThreads as thread (thread.id)}
                <Thread {thread} currentUserId={myUserId} />
              {/each}
            </div>
          {:else if searchTab === "people"}
            <div class="people-list">
              {#each peopleProfile as profile (profile.id)}
                {#if myUserId !=profile.id}
                  <div class="user-card" on:click={() => navigateToProfile(profile.id)}>
                    <img class="avatar" src={profile.profile_picture_url} alt={profile.name} />
                    <div class="info">
                      <span class="name">{profile.name}</span>
                      <span class="username">@{profile.username}</span>
                      <span class="followers">{profile.count} followers</span>
                    </div>

                    {#if profile.is_followed == false}
                      <button type="button" class="follow-btn" on:click|stopPropagation={followUserAPI(profile.id)}>Follow</button>
                    {:else}
                      <button type="button" class="follow-btn" on:click|stopPropagation={unfollowUserAPI(profile.id)}>Unfollow</button>
                    {/if}
                  </div>
                {/if}
              {/each}
              <div class="pagination">
                <button type="button" on:click={() => pagePeople--} disabled={pagePeople === 1}>Prev</button>
                <span>Page {pagePeople}</span>
                <button type="button" on:click={() => pagePeople++} disabled={peopleProfile.length < perPagePeople}>Next</button>
                <select bind:value={perPagePeople}>
                  <option value={25}>25</option>
                  <option value={30}>30</option>
                  <option value={35}>35</option>
                </select>
              </div>
            </div>
            {:else if searchTab === "media"}
            <div class="media-thread-grid">
              {#each relatedMedia as group (group.thread.id)}
                {#if group.media_list.length > 0}
                  <div
                    class="media-thread-block"
                    on:click={() => window.location.href = `/thread/${group.thread.id}`}
                    tabindex="0"
                  >
                    <div class="media-thumb-wrapper">
                      {#if isVideo(group.media_list[0].extension)}
                        <video
                          src={group.media_list[0].image_url}
                          class="media-thumb-rect"
                          muted
                          preload="metadata"
                          playsinline
                          on:loadedmetadata={() => loadVideoDuration(group.media_list[0])}
                          tabindex="-1"
                        />
                        {#if videoDurations[group.media_list[0].id]}
                          <span class="video-duration-badge">{videoDurations[group.media_list[0].id]}</span>
                        {/if}
                      {:else}
                        <img
                          class="media-thumb-rect"
                          src={group.media_list[0].image_url}
                          alt="media"
                          draggable="false"
                        />
                      {/if}
                      {#if group.media_list.length > 1}
                        <div class="stacked-badge">
                          <svg width="30" height="30" viewBox="0 0 24 24" fill="none">
                            <rect x="2" y="7" width="17" height="11" rx="2.5" fill="#F4F4F4" opacity="0.7" />
                            <rect x="5" y="3" width="17" height="11" rx="2.5" fill="#fff" />
                          </svg>
                        </div>
                      {/if}
                    </div>
                  </div>
                {/if}
              {/each}
            </div>
          {:else if searchTab === "communities"}
            <div class="community-list">
              {#each relatedCommunities as c (c.community_id)}
                <div class="community-card" on:click={()=>{navigateToCommunityDetail(c.community_id)}}>
                  <img class="logo" src={c.logo_url} alt={c.community_name} />
                  <div class="c-info">
                    <span class="name">{c.community_name}</span>
                    <span class="desc">{c.community_description}</span>
                  </div>
                </div>
              {/each}
              <div class="pagination">
                <button type="button" on:click={() => pageCommunities--} disabled={pageCommunities === 1}>Prev</button>
                <span>Page {pageCommunities}</span>
                <button type="button" on:click={() => pageCommunities++} disabled={relatedCommunities.length < perPageCommunities}>Next</button>
                <select bind:value={perPageCommunities}>
                  <option value={25}>25</option>
                  <option value={30}>30</option>
                  <option value={35}>35</option>
                </select>
              </div>
            </div>
          {/if}
        </div>
      {/if}
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
  .twitter-layout {
    display: grid;
    grid-template-columns: auto 1fr auto;
    background: #000;
    color: #fff;
  }
  .main-content {
    margin: 0 auto;
    border-left: 1px solid #2f3336;
    border-right: 1px solid #2f3336;
    width: 45vw;
    padding-left: 5vw;
    padding-right: 5vw;
    background: #000;
    z-index: 2;
    display: flex;
    flex-direction: column;
    height: 100vh;
  }
  .explore-header {
    padding: 1.2rem 0 0.8rem 0;
    font-size: 1.5rem;
    font-weight: 800;
    border-bottom: 1px solid #222a;
    margin-bottom: 0.5rem;
  }
  .search-bar-container {
    position: relative;
    width: 100%;
    margin-bottom: 1.2rem;
  }
  .search-bar {
    display: flex;
    align-items: center;
    background: #181c20;
    border-radius: 2rem;
    border: 1.5px solid #222a;
    padding: 0.25rem 0.6rem;
    position: relative;
  }
  .search-bar input {
    flex: 1;
    background: none;
    border: none;
    color: #fff;
    font-size: 1.1rem;
    padding: 0.6rem 0.8rem;
    outline: none;
  }
  .search-btn, .clear-btn {
    background: none;
    border: none;
    color: #bbb;
    cursor: pointer;
    margin-right: 0.15rem;
    padding: 0.1rem 0.3rem;
  }
  .clear-btn { margin-left: 0.2rem;}
  .recent-search-dropdown {
    position: absolute;
    left: 0; right: 0; top: 100%; z-index: 100;
    background: #181c20;
    border: 1px solid #222a;
    border-radius: 0 0 14px 14px;
    padding: 0.65rem 1.1rem;
    box-shadow: 0 2px 24px #000a;
    margin-top: 2px;
  }
  .recent-header {
    display: flex; justify-content: space-between; align-items: center;
    font-weight: 600; color: #1d9bf0; font-size: 1.01rem; margin-bottom: 0.65rem;
  }
  .clear-recent {
    background: none; border: none; color: #bbb; cursor: pointer; font-size: 0.95rem;
  }
  .no-recent {
    color: #999; padding: 1rem 0; text-align: center;
  }
  .recent-item {
    cursor: pointer;
    padding: 0.33rem 0.3rem;
    border-radius: 7px;
    font-size: 1.07rem;
    color: #e2e8f0;
    transition: background 0.14s;
  }
  .recent-item:hover {
    background: #1d9bf0; color: #fff;
  }
  .filters-row {
    display: flex; align-items: center; gap: 1.3rem;
    margin-bottom: 1.2rem; font-size: 1rem;
  }
  .filters-row select {
    padding: 0.22rem 0.7rem; background: #16181c; color: #fff; border: 1.5px solid #222a; border-radius: 8px; outline: none;
  }
  .trending-section {
    margin-top: 1.7rem;
  }
  .trending-list {
    display: flex; flex-direction: column; gap: 0.8rem; margin-top: 1.2rem;
  }
  .trending-item {
    display: flex; align-items: center; gap: 1.1rem;
    background: #181c20;
    padding: 0.7rem 1.1rem;
    border-radius: 12px;
    font-size: 1.13rem;
    cursor: pointer;
    transition: background 0.16s;
  }
  .trending-item:hover { background: #2b2b2b; color: #fff; }
  .hashtag { color: #1d9bf0; display: flex; align-items: center; gap: 6px;}
  .count { margin-left: auto; color: #bbb; font-size: 0.98rem;}
  .nav-tabs {
    width: 100%; border-bottom: 2px solid #222a; margin-bottom: 1.3rem; margin-top: 1.3rem;
    display: flex; gap: 1rem;
  }
  .nav-tabs button {
    appearance: none; background: none; border: none;
    color: #aaa; font-size: 1rem; font-weight: 600;
    padding: 13px 22px 9px 22px; border-radius: 999px;
    cursor: pointer; transition: color 0.16s;
  }
  .nav-tabs button.active { color: #fff; background: #1d9bf0; }
  .tab-content { flex: 1 1 auto; min-height: 0; max-height: 60vh; overflow-y: auto;}
  .loading-center { text-align: center; margin: 2.4rem 0; }
  .spin { animation: spin 1s linear infinite;}
  @keyframes spin { 100% {transform: rotate(360deg);} }
  .search-top h3 { margin: 1.1rem 0 0.65rem 0; overflow-y: auto;}
  .top-people-list { display: flex; flex-direction: column; gap: 1.1rem;}
  .user-card {
    display: flex; align-items: center;
    background: #181818; border-radius: 1rem;
    padding: 1rem 1.3rem; box-shadow: 0 2px 8px 0 #0002;
  }
  .user-card .avatar {
    width: 50px; height: 50px; border-radius: 9999px;
    object-fit: cover; margin-right: 1.2rem;
    border: 2px solid #0ea5e9; background: #2a2a2a;
  }
  .user-card .info {
    display: flex; flex-direction: column;
  }
  .user-card .name { font-weight: 600; color: #fff; font-size: 1.1rem;}
  .user-card .username { font-size: 0.97rem; color: #a5b4fc; margin-bottom: 0.15rem;}
  .user-card .bio { color: #bababa; font-size: 0.97rem; margin-top: 0.25rem;}
  .user-card .followers { font-size: 0.95rem; color: #1d9bf0; }
  .follow-btn { margin-left: auto; background: #1d9bf0; color: #fff; border-radius: 999px; border: none; padding: 0.5rem 1.1rem; font-weight: 700; cursor: pointer;}
  .view-all-btn { margin: 0.5rem 0 0 0; background: #1d9bf0; color: #fff; border: none; border-radius: 9px; padding: 0.4rem 1.3rem; font-weight: 700; cursor: pointer;}
  .media-grid { display: grid; grid-template-columns: repeat(3,1fr); gap: 1.15rem; margin-top: 1.5rem; width: 100%;}
  .media-item { position: relative; cursor: pointer; overflow: hidden; background: #181c20; aspect-ratio: 1 / 1; border-radius: 1.2rem;}
  .media-item img, .media-item video { width: 100%; height: 100%; object-fit: cover; border-radius: 1.2rem;}
  .community-list { display: flex; flex-direction: column; gap: 1.2rem;}
  .community-card { display: flex; align-items: center; background: #181818; border-radius: 1rem; padding: 1rem 1.3rem;}
  .community-card .logo { width: 48px; height: 48px; border-radius: 14px; object-fit: cover; margin-right: 1.2rem; }
  .community-card .c-info { flex: 1 1 0; }
  .community-card .name { font-weight: 600; color: #fff; font-size: 1.1rem;}
  .community-card .desc { font-size: 0.96rem; color: #bababa;}
  .join-btn { background: #1d9bf0; color: #fff; border: none; border-radius: 9px; padding: 0.5rem 1.2rem; font-weight: 700; cursor: pointer; margin-left: 1.1rem;}
  .pagination { display: flex; align-items: center; gap: 1.2rem; margin-top: 1.6rem;}
  .pagination button, .pagination select { background: #181c20; border: 1.5px solid #222a; color: #fff; border-radius: 7px; padding: 0.3rem 0.9rem; font-size: 1.01rem;}
  .pagination select { padding: 0.3rem 0.9rem;}
  @media (max-width: 900px) { .media-grid { grid-template-columns: repeat(2, 1fr); } }
  @media (max-width: 600px) { .media-grid { grid-template-columns: 1fr; gap: 0.7rem; } }

  .trending-threads-header {
  display: flex;
  align-items: center;
  gap: 1.1rem;
  padding-bottom: 0.5rem;
  margin-bottom: 0.7rem;
  background: transparent;
}

.trending-back-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-weight: 500;
  font-size: 1.08rem;
  background: #16181c;
  color: #fff;
  border-radius: 999px;
  border: none;
  padding: 0.38rem 1rem;
  cursor: pointer;
  transition: background 0.14s;
}
.trending-back-btn:hover {
  background: #23272e;
  color: #1d9bf0;
}
.trending-threads-header h2 {
  margin: 0;
  font-size: 1.23rem;
  font-weight: 700;
  color: #1d9bf0;
  letter-spacing: -0.5px;
}
.trending-divider {
  flex: 1;
  height: 1.5px;
  background: linear-gradient(90deg,#222a,#1d9bf0 50%,#222a);
  border-radius: 1px;
  margin-left: 1.1rem;
}

.trending-threads-scroll {
  max-height: 55vh;
  min-height: 320px;
  overflow-y: auto;
  background: #16181c;
  border-radius: 1.2rem;
  box-shadow: 0 3px 24px 0 #0004;
  padding: 1.2rem 1rem;
  display: flex;
  flex-direction: column;
  gap: 1.2rem;
}

.trending-threads-scroll::-webkit-scrollbar {
  width: 8px;
  background: #222a;
  border-radius: 8px;
}
.trending-threads-scroll::-webkit-scrollbar-thumb {
  background: #222a;
  border-radius: 8px;
}

.no-threads {
  text-align: center;
  color: #bbb;
  padding: 2.5rem 0 2rem 0;
  font-size: 1.11rem;
  font-style: italic;
  opacity: 0.82;
}

.profile-suggestion {
  display: flex;
  align-items: center;
  padding: 0.65rem 1.25rem;
  cursor: pointer;
  transition: background 0.18s;
  gap: 0.9rem;
}

.profile-suggestion:hover {
  background: #1a92d2;
}

.profile-pic {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid #e8eaf0;
  background: #eee;
  flex-shrink: 0;
}

.profile-info {
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 2px;
}

.profile-name {
  font-weight: 600;
  font-size: 1.08rem;
  color: #ffffff;
  margin-bottom: 0;
}

.profile-username {
  font-size: 0.93rem;
  color: #ffffff;
}

.no-recent {
  color: #b8b8b8;
  font-size: 1rem;
  padding: 0.6rem 1.25rem;
}

.recent-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.2rem 1.25rem 0.7rem 1.25rem;
  font-weight: 500;
  color: #697386;
  font-size: 1rem;
  border-bottom: 1px solid #f1f1f1;
}

.media-thread-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1.2rem;
  margin-top: 1.5rem;
  width: 100%;
}

@media (max-width: 900px) {
  .media-thread-grid { grid-template-columns: repeat(2, 1fr); }
}
@media (max-width: 600px) {
  .media-thread-grid { grid-template-columns: 1fr; }
}

.media-thread-block {
  position: relative;
  cursor: pointer;
  overflow: hidden;
  background: #181c20;
  aspect-ratio: 1 / 1;
  min-width: 0;
  min-height: 0;
  display: flex;
  align-items: stretch;
  justify-content: stretch;
  border-radius: 1.2rem;
  transition: box-shadow 0.13s;
  box-shadow: 0 2px 8px rgba(0,0,0,0.11);
}
.media-thread-block:hover {
  box-shadow: 0 8px 20px rgba(29,155,240,0.17), 0 2px 8px rgba(0,0,0,0.17);
}

.media-thumb-wrapper {
  position: relative;
  width: 100%;
  height: 100%;
  aspect-ratio: 1 / 1;
  display: flex;
  align-items: stretch;
  justify-content: stretch;
  background: #222;
  overflow: hidden;
  border-radius: 1.2rem;
}

.media-thumb-rect {
  width: 100%;
  height: 100%;
  object-fit: cover;
  object-position: center;
  display: block;
  background: #222;
  border-radius: 0;
  pointer-events: none;
}

.stacked-badge {
  position: absolute;
  top: 8px;
  right: 10px;
  z-index: 2;
  background: transparent;
  pointer-events: none;
}

.video-duration-badge {
  position: absolute;
  left: 10px;
  bottom: 10px;
  background: rgba(0,0,0,0.72);
  color: #fff;
  font-size: 1.08rem;
  font-weight: 600;
  border-radius: 0.4rem;
  padding: 0.1em 0.7em 0.12em 0.5em;
  pointer-events: none;
  user-select: none;
  letter-spacing: 0.01em;
}

@media (max-width: 1249px) {
  .twitter-layout {
    grid-template-columns: 1fr auto;
  }
  .main-content {
    width: 90vw;

    min-width: 0;
    max-width: 91vw;
    padding-left: 0;
    padding-right: 5vw;
    border-left: none;
    border-right: none;
  }
}
@media (max-width: 600px) {

  .trending-section{
    width: 100%;
  }
  .nav-tabs {
    overflow-x: auto;
    flex-wrap: nowrap;
    white-space: nowrap;
    display: flex;
    gap: 0.3rem;
    scrollbar-width: none;
    -ms-overflow-style: none;
    background: #000;
    position: sticky;
    top: 0;
    z-index: 20;
    padding-bottom: 30px;
    margin-bottom: 12px;
    overflow-y: hidden;
  }
  .nav-tabs button.active { color: #fff; background: rgba(255, 255, 255, 0); }
  .nav-tabs::-webkit-scrollbar {
    display: none;
  }
  .filters-row {
    display: none;
  }
  .filters-row label, .filters-row select {
    width: 100%;
    font-size: 0.97rem;
  }
}

</style>