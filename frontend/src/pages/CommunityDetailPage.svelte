<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import LeftSidebar from "../lib/components/LeftSideBar.svelte";
  import RightSidebar from "../lib/components/RightSideBar.svelte";
  import BurgerLeftSideBar from "../lib/components/BurgerLeftSideBar.svelte";
  import BurgerRightSideBar from "../lib/components/BurgerRightSideBar.svelte";
  import Thread from "../lib/components/Thread.svelte";
  import { ArrowLeft, Loader2, Users} from "lucide-svelte";
  import api from "../lib/api";
  import {
    enrichThread,
    type RawThread,
  } from "../lib/threadUtil";
  export let params: { id: string };

  let myUserId = "";

  let currentUserRole = "none";
  let community: Community | null = null;
  let loadingCommunity = true;
  let communityError = false;
  let activeTab: "top" | "latest" | "media" | "about" | "manage";
  let threads: Thread[] = [];
  const loadingThreads = false;

  let videoDurations: Record<string, string> = {};
  function formatDuration(seconds: number): string {
    const m = Math.floor(seconds / 60);
    const s = Math.floor(seconds % 60);
    return `${m}:${s.toString().padStart(2, "0")}`;
  }

  function loadVideoDuration(media: ThreadMedia) {
    if (![".mp4", ".webm", ".mov", ".avi", ".mkv"].includes(media.extension.toLowerCase())) return;
    if (videoDurations[media.id]) return;

    const video = document.createElement("video");
    video.src = media.image_url;
    video.preload = "metadata";
    video.muted = true;
    video.addEventListener("loadedmetadata", () => {
      videoDurations = { ...videoDurations, [media.id]: formatDuration(video.duration) };
    });
  }
  interface Category {
    category_id: string;
    category_name: string;
  }

  interface Community {
    community_id: string;
    community_name: string;
    community_description: string;
    logo_url?: string;
    community_logo: string;
    banner_url?: string;
    categories: Category[];
    created_at: string;
    community_rules?: string;
    member_count: number;
  }

  interface Thread {
    id: string;
    author: string;
    username: string;
    avatar: string;
    content: string;
    community_id: string;
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

  interface ThreadMedia {
    id: string;
    image_url: string;
    extension: string;
  }

  interface GroupedMedia {
      thread: RawThread;
      media_list: ThreadMedia[];
  }

  interface Member {
    id: string;
    name: string;
    username: string;
    avatar: string;
    follower_count: number;
    profile_picture_id: string;
    pp_url: string;
    is_moderator?: boolean;
  }


  async function loadMyProfile() {
    try {
      const meRes = await api.get("/user/get-me");
      myUserId = meRes.data.id;
    } catch (err) {
      console.error(err as string);
    }
  }

  async function loadCommunity() {
    loadingCommunity = true;
    communityError = false;
    try {
      const res = await api.get(`/communities/${params.id}`);
      console.log("comss: ", res.data);
      const rawCommunity = res.data;
      let logo_url = "";
      let banner_url = "";
      if (rawCommunity.community_logo) {
        const logoRes = await api.post("/media/get-media", { id: rawCommunity.community_logo });
        logo_url = logoRes.data.public_url;
      }
      if (rawCommunity.community_banner) {
        const bannerRes = await api.post("/media/get-media", { id: rawCommunity.community_banner });
        banner_url = bannerRes.data.public_url;
      }
      let lenMem = 0;
      try{
        const res = await api.get(`community/count/${params.id}`);
        lenMem = res.data.count;
      }catch(err){
        console.log(err as string);
      }
      let categories: Category[];
      try{
        const categoryRes = await api.get(`/communities/categories/${params.id}`);
        categories = categoryRes.data;
      } catch(err){
        console.error(err as string);
      }
      console.log("categ: ", categories);
      community = {
        community_id: rawCommunity.community_id,
        community_name: rawCommunity.community_name,
        community_description: rawCommunity.community_description,
        logo_url,
        banner_url,
        categories: categories || [],
        created_at: rawCommunity.created_at,
        community_rules: rawCommunity.community_rules,
        member_count: lenMem
      };
    } catch (e) {
      console.error(e);
      communityError = true;
    } finally {
      loadingCommunity = false;
    }
  }

  let topMembers: Member[] = [];
  async function fetchTopUser(){
    try{
      console.log("dwhuwd");
      const userRes = await api.get(`/community/top-user/${params.id}`);
      const members = await Promise.all(userRes.data.member_ids.map(async (id)=>{
        const userMod = await api.get(`user/${id}`);
        return userMod.data;
      }));
      console.log("res: ",topMembers);
      const members_1 = await Promise.all(members.map(async (member)=>{
        const userPp = await api.post("media/get-media", {id: member.profile_picture_id});
        const pp = userPp.data;
        return {...member, public_url: pp.public_url};
      }));
      const members_2 = await Promise.all(members_1.map(async (member)=>{
        const userCount = await api.get(`/user/follower_count/${member.id}`);
        const c = userCount.data;
        return {...member, count: c.count};
      }));
      topMembers = members_2;
      console.log("topss: ",topMembers);
    }catch(err){
      console.error(err as string);
    }
  }
  let topThread: RawThread[];
  async function fetchTopThreads(page = 1) {
    const size = 2;
    try{
      const tData = await api.get(`/threads/community/${params.id}/by-like?page=${page}&size=${size}`);
      const enriched = await Promise.all(tData.data.map((r) => enrichThread(r, myUserId)));
      topThread = enriched;
    }catch(err){
      console.log(err as string);
    }
  }
  let latestThread: RawThread[];
  async function fetchLatestThreads() {
    try{
      const tData =  await api.get(`/threads/community/${params.id}/by-time?page=${1}&size=${100}`);
      const enriched = await Promise.all(tData.data.map((r) => enrichThread(r, myUserId)));
      latestThread = enriched;
    }
    catch(err){
      console.error(err as string);
    }
  }

  let groupedMedia: GroupedMedia[] = [];
  let mediaPage = 1;
  const mediaSize = 3;
  let hasMoreMedia = true;
  let loadingMedia = false;
  let mediaObserverTarget: HTMLDivElement;
  let mediaObserver: IntersectionObserver;

  async function fetchCommunityMedia(page = 1) {
  if (loadingMedia || !hasMoreMedia) return;
  loadingMedia = true;
  try {
    const res = await api.get<GroupedMedia[]>(`/threads/community/${params.id}/media?page=${page}&size=${mediaSize}`);
    let data = res.data;
    data = await Promise.all(
      data.map(async (group) => ({
        thread: await enrichThread(group.thread, myUserId),
        media_list: group.media_list,
      }))
    );
    if (page === 1) {
      groupedMedia = data;
    } else {
      groupedMedia = [...groupedMedia, ...data];
    }
    if (data.length < mediaSize) hasMoreMedia = false;
    else mediaPage = page + 1;
  } catch (err) {
    console.error(err);
    hasMoreMedia = false;
  } finally {
    loadingMedia = false;
  }
}

  function setupMediaObserver() {
    if (mediaObserver) mediaObserver.disconnect();
    mediaObserver = new IntersectionObserver(
      (entries) => {
        if (
          entries[0].isIntersecting &&
          !loadingMedia &&
          hasMoreMedia &&
          activeTab === "media"
        ) {
          fetchCommunityMedia(mediaPage);
        }
      },
      { rootMargin: "200px" }
    );
    if (mediaObserverTarget) mediaObserver.observe(mediaObserverTarget);
  }

  let owner: Member;
  async function getOwnerProfile()
  {
    try{
      const ownerRes = await api.get(`communities/members/${params.id}?role=owner`);
      const members_0 = await Promise.all(ownerRes.data.map(async (member)=>{
        const userPp = await api.get(`user/${member.member_id}`);
        const pp = userPp.data;
        return {...member, profile_picture_id: pp.profile_picture_id, username: pp.username, name: pp.name};
      }));
      const members_1 = await Promise.all(members_0.map(async (member)=>{
        const userPp = await api.post("media/get-media", {id: member.profile_picture_id});
        const pp = userPp.data;
        return {...member, public_url: pp.public_url};
      }));
      const members_2 = await Promise.all(members_1.map(async (member)=>{
        const userCount = await api.get(`/user/follower_count/${member.member_id}`);
        const c = userCount.data;
        return {...member, count: c.count};
      }));
      owner = members_2;
      console.log("owner: ", owner);

    }catch(err)
    {
      console.error(err as string);
    }
  }
  let moderator: Member[];
  async function getAllModeratorProfile()
  {
    try{
      const moderatorRes = await api.get(`communities/members/${params.id}?role=moderator`);
      const members_0 = await Promise.all(moderatorRes.data.map(async (member)=>{
        const userPp = await api.get(`user/${member.member_id}`);
        const pp = userPp.data;
        return {...member, profile_picture_id: pp.profile_picture_id, username: pp.username, name: pp.name};
      }));
      const members_1 = await Promise.all(members_0.map(async (member)=>{
        const userPp = await api.post("media/get-media", {id: member.profile_picture_id});
        const pp = userPp.data;
        return {...member, public_url: pp.public_url};
      }));
      const members_2 = await Promise.all(members_1.map(async (member)=>{
        const userCount = await api.get(`/user/follower_count/${member.member_id}`);
        const c = userCount.data;
        return {...member, count: c.count};
      }));
      moderator = members_2;
      console.log("moderator: ", moderator);

    }catch(err)
    {
      console.error(err as string);
    }
  }
  let showMembersModal = false;
  let membersTab: "members" | "moderators" = "members";
  let memberSearch = "";
  let memberPage = 1;
  let memberPageSize = 25;
  let totalMembers = 0;
  let totalModerators = 0;
  let membersList: Member[] = [];
  let moderatorsList: Member[] = [];
  let pendingList: Member[] = [];
  let loadingMemberList = false;

  let searchTimeout: number | undefined;
  function onSearchInput(e: Event) {
    memberSearch = (e.target as HTMLInputElement).value;
    clearTimeout(searchTimeout);
    searchTimeout = setTimeout(fetchMembersTab, 350);
  }
  async function fetchMembersTab() {
    loadingMemberList = true;
    if (membersTab === "members") {
      const res = await api.get(
        `/communities/members/${params.id}?role=accepted&page=${memberPage}&size=${memberPageSize}&search=${encodeURIComponent(memberSearch)}`
      );
      membersList = res.data;
      const members_0 = await Promise.all(membersList.map(async (member)=>{
        const userPp = await api.get(`user/${member.member_id}`);
        const pp = userPp.data;
        return {...member, profile_picture_id: pp.profile_picture_id, username: pp.username, name: pp.name};
      }));
      const members_1 = await Promise.all(members_0.map(async (member)=>{
        const userPp = await api.post("media/get-media", {id: member.profile_picture_id});
        const pp = userPp.data;
        return {...member, public_url: pp.public_url};
      }));
      const members_2 = await Promise.all(members_1.map(async (member)=>{
        const userCount = await api.get(`/user/follower_count/${member.member_id}`);
        const c = userCount.data;
        return {...member, count: c.count};
      }));
      membersList = members_2;
      console.log("mem Lists: ", membersList);
      totalMembers = res.data.length;
      console.log("tot member: ", totalMembers);
    } else {
      const res = await api.get(
        `/communities/members/${params.id}?role=moderator&page=${memberPage}&size=${memberPageSize}&search=${encodeURIComponent(memberSearch)}`
      );
      const members_0 = await Promise.all(res.data.map(async (member)=>{
        const userPp = await api.get(`user/${member.member_id}`);
        const pp = userPp.data;
        return {...member, profile_picture_id: pp.profile_picture_id, username: pp.username, name: pp.name};
      }));
      const members_1 = await Promise.all(members_0.map(async (member)=>{
        const userPp = await api.post("media/get-media", {id: member.profile_picture_id});
        const pp = userPp.data;
        return {...member, public_url: pp.public_url};
      }));
      const members_2 = await Promise.all(members_1.map(async (member)=>{
        const userCount = await api.get(`/user/follower_count/${member.member_id}`);
        const c = userCount.data;
        return {...member, count: c.count};
      }));
      moderatorsList = members_2;
      console.log("mod list: ", moderatorsList);
      totalModerators = res.data.length;
      console.log("tot mod: ", totalModerators);
    }
    loadingMemberList = false;
  }
  function nextMemberPage() {
    memberPage++;
    fetchMembersTab();
  }
  function prevMemberPage() {
    if (memberPage > 1) {
      memberPage--;
      fetchMembersTab();
    }
  }
  async function openMembersModal(tab: "members" | "moderators" = "members") {
    showMembersModal = true;
    membersTab = tab;
    memberPage = 1;
    memberPageSize = 25;
    memberSearch = "";
    await fetchMembersTab();
  }

  async function promoteToModerator(userId: string) {
    try {
      await api.post(`/communities/members/promote/${params.id}/${userId}`);
      await fetchMembersTab();
    } catch (err) {
      console.error("Failed to promote user:", err);
    }
  }
  async function demoteModerator(userId: string) {
    try {
      await api.post(`/communities/members/demote/${params.id}/${userId}`);
      await fetchMembersTab();
    } catch (err) {
      console.error("Failed to demote moderator:", err);
    }
  }
  async function fetchUserRole(communityId: string, userId: string): Promise<string | null> {
    try {
      const res = await api.get(`/user/community/role/${communityId}/${userId}`);
      currentUserRole = res.data.role;
      return currentUserRole;
    } catch (err) {
      console.error("Failed to fetch user role:", err);
      return null;
    }
  }


  async function getPendingMembers() {
    try{
      const res = await api.get(
        `/communities/members/${params.id}?role=pending&page=${memberPage}&size=${memberPageSize}&search=${encodeURIComponent(memberSearch)}`
      );
      const members_0 = await Promise.all(res.data.map(async (member)=>{
        const userPp = await api.get(`user/${member.member_id}`);
        const pp = userPp.data;
        return {...member, profile_picture_id: pp.profile_picture_id, username: pp.username, name: pp.name};
      }));
      const members_1 = await Promise.all(members_0.map(async (member)=>{
        const userPp = await api.post("media/get-media", {id: member.profile_picture_id});
        const pp = userPp.data;
        return {...member, public_url: pp.public_url};
      }));
      const members_2 = await Promise.all(members_1.map(async (member)=>{
        const userCount = await api.get(`/user/follower_count/${member.member_id}`);
        const c = userCount.data;
        return {...member, count: c.count};
      }));
      pendingList = members_2;
    }
    catch(err){
      console.error(err as string);
    }
  }

  function changeTab(tab: typeof activeTab) {
    activeTab = tab;
    if (tab === "top") {
      threads = [];
      fetchTopUser();
      fetchTopThreads();
    } else if (tab === "latest") {
      fetchLatestThreads();
    } else if (tab === "media") {
      groupedMedia = [];
      mediaPage = 1;
      hasMoreMedia = true;
      fetchCommunityMedia(1);
      setTimeout(setupMediaObserver, 100);
    } else if (tab === "about") {
      getOwnerProfile();
      getAllModeratorProfile();
    }
    else if (tab === "manage") {
      getPendingMembers();
    }
  }

  function navigateToProfile(userId: string) {
    window.location.href = `/profile/${userId}`;
  }
  function changeMemberPageSize(newSize: number) {
    memberPageSize = newSize;
    memberPage = 1;
    fetchMembersTab();
  }
  onMount(async () => {
    await loadMyProfile();
    await loadCommunity();
    await fetchUserRole(params.id, myUserId);
    changeTab("top");
  });
  function goToThreadDetail(threadId: string) {
    window.location.href = `/thread/${threadId}`;
  }
  onDestroy(() => {
    onDestroy(() => {
      if (mediaObserver && mediaObserverTarget) mediaObserver.unobserve(mediaObserverTarget);
    });
  });
  function isVideo(extension: string) {
    return [".mp4", ".webm", ".mov", ".avi", ".mkv"].includes(extension.toLowerCase());
  }
  function goBack() {
    if (window.history.length > 1) {
      window.history.back();
    } else {
      window.location.href = "/home";
    }
  }
  async function approveMember(userId: string) {
    try {
      await api.post(`/communities/requests/accept/${params.id}/${userId}`);
      pendingList = pendingList.filter(member => member.member_id !== userId);
    } catch (err) {
      console.error("Failed to approve member:", err);
    }
  }

  async function rejectMember(userId: string) {
    try {
      await api.post(`/communities/requests/reject/${params.id}/${userId}`);
      pendingList = pendingList.filter(member => member.member_id !== userId);
    } catch (err) {
      console.error("Failed to reject member:", err);
    }
  }

  async function requestJoin() {
    try{
      await api.post(`/communities/join/${params.id}`, {user_id: myUserId});
      currentUserRole = "pending";
    } catch(err){
      console.error(err as string);
    }
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

<div class="twitter-layout">
  <aside class="sidebar">
    {#if myUserId && windowWidth >= 1250}
      <LeftSidebar currentUserId={myUserId} activePage="communities" />
    {/if}
  </aside>
  {#if myUserId && windowWidth < 1250}
    <BurgerLeftSideBar currentUserId={myUserId} activePage="communities" />
  {/if}

  <main class="main-content">
    <div class="header-bar">
      <button type="button" class="back-btn" on:click={goBack}>
        <ArrowLeft size="20" />
      </button>
      <div class="header-info">
        <h1>Community</h1>
        {#if community}
        <span class="member-count" on:click={() => openMembersModal("members")}>
          <Users size="16" />
          {community.member_count} members
        </span>
        {/if}
      </div>
    </div>

    {#if loadingCommunity}
      <div class="loading-state">
        <div class="skeleton-header"></div>
        <div class="skeleton-banner"></div>
        <div class="skeleton-info"></div>
      </div>
    {:else if communityError}
      <div class="error-state">
        <p>This community doesn't exist or has been deleted.</p>
      </div>
    {:else if community}
      {#if community.banner_url}
        <div class="banner-container">
          <img src={community.banner_url} alt="Community Banner" class="banner" />
        </div>
      {:else}
        <div class="banner-placeholder"></div>
      {/if}

      <div class="community-info">
        <div class="community-header">
          <div class="avatar-container">
            {#if community.logo_url}
              <img src={community.logo_url} alt="Community Logo" class="community-avatar" />
            {:else}
              <div class="community-avatar-placeholder">
                <Users size="40" />
              </div>
            {/if}
          </div>
        </div>

        <div class="community-details">
          <h2 class="community-name">{community.community_name}</h2>
          <p class="community-description">{community.community_description}</p>
          <div class="community-meta">
            <span class="member-count" on:click={() => openMembersModal("members")}>
              <Users size="16" />
              {community.member_count} members
            </span>
            <span class="join-date">
              Created {community.created_at.year}-{community.created_at.month}-{community.created_at.day}
            </span>
          </div>

          {#if community.categories && community.categories.length > 0}
            <div class="categories">
              {#each community.categories as category (category.category_id)}
                <span class="category-tag">{category.category_name}</span>
              {/each}
            </div>
          {/if}
        </div>
        {#if !currentUserRole || currentUserRole === "none"}
          <div class="join-community-bar">
            <button type="button" class="join-community-btn" on:click={requestJoin}>Join Community</button>
          </div>
        {/if}

      </div>

      <div class="nav-tabs" data-active-tab={["top","latest","media","about", "manage"].indexOf(activeTab)}>
        <button
          type="button"
          class="nav-tab"
          class:active={activeTab === "top"}
          on:click={() => changeTab("top")}
        >Top</button>
        <button
          type="button"
          class="nav-tab"
          class:active={activeTab === "latest"}
          on:click={() => changeTab("latest")}
        >Latest</button>
        <button
          type="button"
          class="nav-tab"
          class:active={activeTab === "media"}
          on:click={() => changeTab("media")}
        >Media</button>
        <button
          type="button"
          class="nav-tab"
          class:active={activeTab === "about"}
          on:click={() => changeTab("about")}
        >About</button>
        {#if currentUserRole == "moderator" || currentUserRole == "owner"}
          <button
            type="button"
            class="nav-tab"
            class:active={activeTab === "manage"}
            on:click={() => changeTab("manage")}
          >Manage</button>
        {/if}
        <!-- <div class="tab-indicator"></div> -->
      </div>

      <div class="tab-content">
        {#if activeTab === "top"}
          <div class="posts-feed">
            <div class="user-section">
              <h1 class="top-user-title">Top Users</h1>
              <div class="top-user-list">
                {#each topMembers as member (member.id)}
                  <div class="top-user-card">
                    <img
                      class="user-avatar"
                      src={member.public_url}
                      alt={member.name}
                    />
                    <div class="user-info">
                      <div class="user-name">{member.name}</div>
                      <div class="user-username">@{member.username}</div>
                      <div class="user-name">Follower: {member.count}</div>
                      {#if member.bio}
                        <div class="user-bio">{member.bio}</div>
                      {/if}
                    </div>
                  </div>
                {/each}
              </div>
            </div>
            <div class="thread-section">
              <h1>Thread Section</h1>
              {#each topThread as thread (thread.id)}
                <Thread {thread} currentUserId={myUserId} />
              {/each}
            </div>
          </div>
        {:else if activeTab === "latest"}
          <div class="posts-feed">
            <div class="thread-section">
              {#each latestThread as thread (thread.id)}
                <Thread {thread} currentUserId={myUserId} />
              {/each}
            </div>
          </div>
        {:else if activeTab === "media"}
        <div class="media-thread-grid">
          {#each groupedMedia as group (group.thread.id)}
            {#if group.media_list.length > 0}
              <div
                class="media-thread-block"
                on:click={() => goToThreadDetail(group.thread.id)}
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
          <div class="observer" bind:this={mediaObserverTarget}></div>
        </div>
        {:else if activeTab === "about"}
          <div class="community-overview-card">
            <h2 class="overview-title">Community Overview</h2>
            <div class="overview-row">
              <span class="overview-label">Created by:</span>
              {#if owner && owner.length > 0}
                <div class="owner-info" on:click={() => navigateToProfile(owner[0].member_id)}>
                  <img class="owner-avatar" src={owner[0].public_url} alt="Owner Avatar" />
                  <span class="owner-name">{owner[0].name}</span>
                  <span class="owner-username">@{owner[0].username}</span>
                </div>
              {:else}
                <span>Unknown</span>
              {/if}
            </div>
            <div class="overview-row">
              <span class="overview-label">Created at:</span>
              <span>
                {#if community}
                  {community.created_at.year}-{community.created_at.month}-{community.created_at.day}
                {/if}
              </span>
            </div>
            <div class="overview-row">
              <span class="overview-label">Community Rules:</span>
              <span class="rules-text">{community.community_rules}</span>
            </div>
            <div class="overview-row mods">
              <span class="overview-label">Moderators:</span>
              <div class="moderators-list">
                {#if moderator && moderator.length > 0}
                  {#each moderator as mod (mod.member_id)}
                    <div class="moderator-card" on:click={() => navigateToProfile(mod.member_id)}>
                      <img class="mod-avatar" src={mod.public_url} alt="Moderator Avatar" />
                      <span class="mod-name">{mod.name}</span>
                      <span class="mod-username">@{mod.username}</span>
                    </div>
                  {/each}
                {:else}
                  <span class="no-moderator">No moderators yet</span>
                {/if}
              </div>
            </div>
          </div>

          {:else if activeTab === "manage"}
          <div class="manage-section">
            <h3 class="manage-title">Pending Member Approvals</h3>
            {#if pendingList.length === 0}
              <div class="empty-modal-list">No pending member requests.</div>
            {:else}
              <div class="pending-member-list">
                {#each pendingList as member (member.member_id)}
                  <div class="pending-member-card">
                    <img class="pending-avatar" src={member.public_url || member.avatar} alt={member.name} />
                    <div class="pending-info">
                      <div class="pending-name">{member.name}</div>
                      <div class="pending-username">@{member.username}</div>
                    </div>
                    <div class="pending-actions">
                      <button type="button" class="approve-btn" on:click={() => approveMember(member.member_id)}>Approve</button>
                      <button type="button" class="reject-btn" on:click={() => rejectMember(member.member_id)}>Reject</button>
                    </div>
                  </div>
                {/each}
              </div>
            {/if}
          </div>
        {/if}

      </div>

      {#if loadingThreads && threads.length > 0}
        <div class="loading-more">
          <Loader2 size="16" class="spin" /> Loading more posts...
        </div>
      {/if}

    {/if}
    {#if showMembersModal}
      <div class="members-modal-backdrop" on:click={() => showMembersModal = false}></div>
      <div class="members-modal" on:click|stopPropagation>
        <div class="members-modal-header">
          <h2>Community Members</h2>
          <button type="button" class="close-btn" on:click={() => showMembersModal = false}>&times;</button>
        </div>
        <div class="members-modal-tabs">
          <button
            type="button"
            class:active={membersTab === "members"}
            on:click={() => { membersTab = "members"; memberPage = 1; fetchMembersTab(); }}
          >
            Members
          </button>
          <button
            type="button"
            class:active={membersTab === "moderators"}
            on:click={() => { membersTab = "moderators"; memberPage = 1; fetchMembersTab(); }}
          >
            Moderators
          </button>
        </div>
        <div class="members-modal-search">
          <input
            type="text"
            placeholder="Search users…"
            value={memberSearch}
            on:input={onSearchInput}
          />
          <select bind:value={memberPageSize} on:change={(e) => changeMemberPageSize(Number(e.target.value))}>
            <option value={25}>25 / page</option>
            <option value={30}>30 / page</option>
            <option value={35}>35 / page</option>
          </select>
        </div>
        <div class="members-modal-list">
          {#if loadingMemberList}
            <div class="loading-modal-members">Loading...</div>
          {:else if membersTab === "members" && membersList.length === 0}
            <div class="empty-modal-list">No members found.</div>
          {:else if membersTab === "moderators" && moderatorsList.length === 0}
            <div class="empty-modal-list">No moderators found.</div>
          {:else}
            {#if membersTab === "members"}
              {#each membersList as member (member.member_id)}
                <div class="members-modal-item" on:click={() => navigateToProfile(member.member_id)}>
                  <img class="modal-user-avatar" src={member.public_url || member.avatar} alt={member.name} />
                  <div class="modal-user-info">
                    <div class="modal-user-name">{member.name}</div>
                    <div class="modal-user-username">@{member.username}</div>
                  </div>
                  {#if currentUserRole == "owner"}
                    <button
                      type="button"
                      class="promote-btn"
                      title="Promote to Moderator"
                      on:click|stopPropagation={() => promoteToModerator(member.member_id)}
                    >Promote</button>
                  {/if}
                </div>
              {/each}
            {:else}
              {#each moderatorsList as moderator (moderator.member_id)}
                <div class="members-modal-item" on:click={() => navigateToProfile(moderator.member_id)}>
                  <img class="modal-user-avatar" src={moderator.public_url || moderator.avatar} alt={moderator.name} />
                  <div class="modal-user-info">
                    <div class="modal-user-name">{moderator.name}</div>
                    <div class="modal-user-username">@{moderator.username}</div>
                  </div>
                  {#if currentUserRole == "owner"}
                    <button
                      type="button"
                      class="demote-btn"
                      title="Demote to Member"
                      on:click|stopPropagation={() => demoteModerator(moderator.member_id)}
                    >Demote</button>
                  {/if}
                </div>
              {/each}
            {/if}
          {/if}
        </div>
        <div class="members-modal-pagination">
          <button type="button" disabled={memberPage === 1} on:click={prevMemberPage}>&lt; Prev</button>
          <span>Page {memberPage}</span>
          <button
            type="button"
            disabled={(membersTab === "members" && memberPage * memberPageSize >= totalMembers) ||
              (membersTab === "moderators" && memberPage * memberPageSize >= totalModerators)}
            on:click={nextMemberPage}
          >Next &gt;</button>
        </div>
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
  .observer{
    height: 10px;
  }
  .tab-content{
    flex: 1 1 auto;
    min-height: 0;
    overflow-y: auto;
  }
  .twitter-layout {
    display: grid;
    grid-template-columns: auto 1fr auto;
    margin: 0 auto;
    width: 100vw;
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

  .header-bar {
    position: sticky;
    top: 0;
    z-index: 10;
    background: rgba(0, 0, 0, 0.85);
    backdrop-filter: blur(12px);
    padding: 2px 4px;
    border-bottom: 1px solid #2f3336;
    display: flex;
    align-items: center;
    gap: 20px;
  }
  .categories{
    display: flex;
    flex-direction: row;
    gap: 15px;
  }
  .back-btn {
    background: none;
    border: none;
    color: #fff;
    cursor: pointer;
    padding: 8px;
    border-radius: 50%;
    transition: background-color 0.2s;
  }

  .back-btn:hover {
    background: rgba(255, 255, 255, 0.1);
  }

  .header-info h1 {
    font-size: 20px;
    font-weight: 800;
    margin: 0;
    line-height: 1.2;
  }

  .header-info .member-count {
    font-size: 13px;
    color: #71767b;
    font-weight: 400;
    cursor: pointer;
  }

  .banner-container {
    width: 100%;
    height: 150px;
    min-height: 150px;
    max-height: 150px;
    overflow: hidden;
    background: #16181c;
  }

  .banner {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .banner-placeholder {
    width: 100%;
    height: 200px;
    background: #16181c;
  }

  .community-info {
    padding: 12px 6px 0;
  }

  .community-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 16px;
  }

  .avatar-container {
    margin-top: -74px;
  }

  .community-avatar {
    width: 134px;
    height: 134px;
    border-radius: 50%;
    border: 4px solid #000;
    background: #16181c;
    object-fit: cover;
  }

  .community-avatar-placeholder {
    width: 134px;
    height: 134px;
    border-radius: 50%;
    border: 4px solid #000;
    background: #16181c;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #71767b;
  }

  .action-buttons {
    display: flex;
    gap: 8px;
    margin-top: 12px;
  }

  .icon-btn {
    width: 36px;
    height: 36px;
    border: 1px solid #536471;
    background: none;
    color: #fff;
    border-radius: 50%;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: background-color 0.2s;
  }

  .icon-btn:hover {
    background: rgba(255, 255, 255, 0.1);
  }

  .join-button {
    background: #1d9bf0;
    color: #fff;
    border: none;
    padding: 8px 16px;
    border-radius: 20px;
    font-weight: 700;
    cursor: pointer;
    transition: background-color 0.2s;
  }

  .join-button:hover {
    background: #1a8cd8;
  }

  .joined-button {
    background: #000;
    color: #fff;
    border: 1px solid #536471;
    padding: 8px 16px;
    border-radius: 20px;
    font-weight: 700;
    cursor: pointer;
    transition: all 0.2s;
  }

  .community-details {
    margin-bottom: 16px;
  }

  .community-name {
    font-size: 22px;
    font-weight: 800;
    margin: 0 0 8px 0;
  }

  .community-description {
    font-size: 15px;
    color: #fff;
    margin: 0 0 12px 0;
    line-height: 1.4;
  }

  .community-meta {
    display: flex;
    gap: 16px;
    margin-bottom: 12px;
    color: #71767b;
    font-size: 15px;
  }

  .community-meta span {
    display: flex;
    align-items: center;
  }
  .nav-tabs {
    width: 100%;
    background: none;
    border-bottom: 2px solid #222a;
    margin: 0 0 14px 0;
    position: relative;
    overflow-x: auto;
    user-select: none;
    min-height: 50px;
    max-height: 50px;
  }

  .nav-tabs button {
    appearance: none;
    background: none;
    border: none;
    color: #888;
    font-size: 1rem;
    font-weight: 600;
    padding: 14px 26px 10px 26px;
    cursor: pointer;
    border-radius: 0;
    transition: color 0.17s cubic-bezier(.4,0,.2,1);
    outline: none;
    letter-spacing: 0.02em;
    z-index: 2;

  }

  @media (max-width: 768px) {
    .nav-tabs button { padding: 12px 12px 8px 12px; font-size: 0.95rem; }
    .tab-indicator { width: 60px; }
  }

  .nav-tabs {
  display: flex;
  position: relative;
  background: rgba(20, 20, 20, 0.85);
  border-radius: 999px;
  padding: 0.2rem;
  gap: 0.25rem;
  margin-bottom: 2rem;
  height: 90px;
}

.nav-tab {
  flex: 1;
  background: none;
  border: none;
  color: #aaa;
  font-weight: 600;
  padding: 1rem 0.2rem;
  border-radius: 999px;
  font-size: 1rem;
  cursor: pointer;
  position: relative;
  z-index: 1;
  transition: color 0.2s;
}

.nav-tab.active {
  color: #fff;
}

.tab-indicator {
  position: absolute;
  left: 0;
  bottom: 0;
  height: 4px;
  width: 20%;
  background: #1d9bf0;
  border-radius: 999px 999px 0 0;
  transition: transform 0.25s cubic-bezier(.4,1,.6,1);
  z-index: 0;
  transform: translateX(calc(var(--tab-index, 0) * 100%));
}

.nav-tabs[data-active-tab="0"] .tab-indicator { --tab-index: 0; }
.nav-tabs[data-active-tab="1"] .tab-indicator { --tab-index: 1; }
.nav-tabs[data-active-tab="2"] .tab-indicator { --tab-index: 2; }
.nav-tabs[data-active-tab="3"] .tab-indicator { --tab-index: 3; }
.nav-tabs[data-active-tab="4"] .tab-indicator { --tab-index: 4; }
.category-tag{
  background-color: #1d9bf0;
  padding-left: 10px;
  padding-right: 10px;
  padding-top: 5px;
  padding-bottom: 5px;
  border-radius: 10px;
}
.sidebar{
  z-index: 400000;
}

.user-section {
  margin-bottom: 2.5rem;
}
.top-user-title {
  font-weight: 700;
  margin-bottom: 1rem;
  color: #fff;
}
.top-user-list {
  display: flex;
  flex-direction: column;
  gap: 1.1rem;
}
.top-user-card {
  display: flex;
  align-items: center;
  background: #181818;
  border-radius: 1rem;
  padding: 1rem 1.3rem;
  box-shadow: 0 2px 8px 0 rgba(0,0,0,0.10);
  transition: box-shadow 0.15s;
}
.top-user-card:hover {
  box-shadow: 0 4px 20px 0 rgba(30,50,90,0.19);
}
.user-avatar {
  width: 50px;
  height: 50px;
  border-radius: 9999px;
  object-fit: cover;
  margin-right: 1.2rem;
  border: 2px solid #0ea5e9;
  background: #2a2a2a;
}
.user-info {
  display: flex;
  flex-direction: column;
}
.user-name {
  font-weight: 600;
  color: #fff;
  font-size: 1.1rem;
}
.user-username {
  font-size: 0.97rem;
  color: #a5b4fc;
  margin-bottom: 0.15rem;
}
.user-bio {
  color: #bababa;
  font-size: 0.97rem;
  margin-top: 0.25rem;
}

.media-thread-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1.2rem;
  margin-top: 1.5rem;
  width: 100%;
}

.media-thread-block {
  position: relative;
  cursor: pointer;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0,0,0,0.11);
  background: #181c20;
  aspect-ratio: 1 / 1;
  min-width: 0;
  min-height: 0;
  display: flex;
  align-items: stretch;
  justify-content: stretch;
  border-radius: 1.2rem;
  transition: box-shadow 0.13s;
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

.stacked-badge {
  position: absolute;
  top: 8px;
  right: 10px;
  z-index: 2;
  background: transparent;
  pointer-events: none;
}

.observer {
  height: 10px;
}

@media (max-width: 900px) {
  .media-thread-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 600px) {
  .media-thread-grid {
    grid-template-columns: 1fr;
    gap: 0.7rem;
  }
}

.community-overview-card {
  background: #16181c;
  border-radius: 1.1rem;
  box-shadow: 0 2px 16px rgba(0,0,0,0.10);
  margin-bottom: 2.2rem;
  padding: 1.4rem 1.6rem 1.2rem 1.6rem;
  border: 1.5px solid rgba(29,155,240,0.09);
}
.overview-title {
  font-size: 1.35rem;
  font-weight: 800;
  margin-bottom: 1.2rem;
  color: #fff;
  letter-spacing: 0.01em;
}
.overview-row {
  display: flex;
  align-items: center;
  gap: 1.2rem;
  margin-bottom: 1rem;
  flex-wrap: wrap;
}
.overview-label {
  font-weight: 600;
  color: #1d9bf0;
  min-width: 110px;
  font-size: 1rem;
}
.owner-info {
  display: flex;
  align-items: center;
  cursor: pointer;
  gap: 0.5rem;
}
.owner-avatar {
  width: 36px;
  height: 36px;
  border-radius: 999px;
  border: 2px solid #1d9bf0;
  margin-right: 0.5rem;
}
.owner-name, .mod-name {
  font-weight: 600;
  color: #fff;
  font-size: 1.05rem;
}
.owner-username, .mod-username {
  font-size: 0.98rem;
  color: #a5b4fc;
  margin-left: 0.15rem;
}
.rules-text {
  color: #e3e3e3;
  font-size: 1.04rem;
  max-width: 500px;
  word-break: break-word;
}
.moderators-list {
  display: flex;
  flex-wrap: wrap;
  gap: 1.15rem;
  margin-left: -0.8rem;
}
.moderator-card {
  display: flex;
  align-items: center;
  background: #222c;
  padding: 0.4rem 1rem;
  border-radius: 1rem;
  cursor: pointer;
  transition: box-shadow 0.13s;
  border: 1.5px solid rgba(29,155,240,0.06);
}
.moderator-card:hover {
  box-shadow: 0 4px 16px rgba(29,155,240,0.13);
  background: #233344;
}
.mod-avatar {
  width: 30px;
  height: 30px;
  border-radius: 999px;
  border: 1.5px solid #1d9bf0;
  margin-right: 0.5rem;
}
.no-moderator {
  color: #aaa;
  font-style: italic;
  margin-left: 0.6rem;
}

.member-count {
  cursor: pointer;
  transition: background 0.17s;
  border-radius: 8px;
  padding: 4px 10px;
  display: inline-flex;
  align-items: center;
  user-select: none;
}
.member-count:hover {
  background: #232d40;
  color: #fff;
}

.members-modal-backdrop {
  position: fixed;
  z-index: 10000;
  inset: 0;
  background: rgba(0,0,0,0.64);
  transition: background 0.25s;
}

.members-modal {
  position: fixed;
  z-index: 11000;
  left: 50%;
  top: 50%;
  transform: translate(-50%, -50%);
  min-width: 350px;
  width: 95vw;
  max-width: 430px;
  max-height: 85vh;
  background: #181c20;
  border-radius: 18px;
  box-shadow: 0 8px 40px rgba(0,0,0,0.36), 0 0 0 2px #1d9bf044;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  animation: pop-modal 0.23s cubic-bezier(.33,1.2,.64,1) both;
}

@keyframes pop-modal {
  0% { transform: translate(-50%, -40%) scale(0.95);}
  100% { transform: translate(-50%, -50%) scale(1);}
}

.members-modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1.2rem 1.6rem 0.8rem 1.6rem;
  border-bottom: 1.5px solid #22334a;
  background: #181c20;
}
.members-modal-header h2 {
  margin: 0;
  font-size: 1.28rem;
  font-weight: 700;
  color: #fff;
}
.close-btn {
  background: none;
  border: none;
  color: #fff;
  font-size: 2rem;
  line-height: 1;
  cursor: pointer;
  border-radius: 50%;
  transition: background 0.18s;
  padding: 0 0.5rem;
}
.close-btn:hover {
  background: #233344;
}

.members-modal-tabs {
  display: flex;
  padding: 0.7rem 1.6rem 0 1.6rem;
  background: #181c20;
  gap: 1rem;
}
.members-modal-tabs button {
  background: none;
  border: none;
  font-size: 1.05rem;
  font-weight: 700;
  color: #a0b6ce;
  cursor: pointer;
  padding: 0.38rem 1.3rem;
  border-radius: 12px 12px 0 0;
  transition: color 0.19s, background 0.16s;
}
.members-modal-tabs button.active {
  color: #fff;
  background: #1d9bf0;
}

.members-modal-search {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 0.6rem 1.6rem;
  background: #181c20;
}
.members-modal-search input {
  flex: 1;
  padding: 0.6rem 1rem;
  font-size: 1rem;
  border-radius: 0.7rem;
  border: 1.5px solid #22334a;
  background: #232d40;
  color: #fff;
  outline: none;
}
.members-modal-search select {
  padding: 0.35rem 0.7rem;
  border-radius: 0.6rem;
  background: #181c20;
  color: #fff;
  border: 1.5px solid #22334a;
  font-size: 1rem;
}

.members-modal-list {
  flex: 1 1 0;
  min-height: 120px;
  overflow-y: auto;
  background: #181c20;
  padding: 0.6rem 1.6rem 0.6rem 1.6rem;
}
.members-modal-item {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 0.75rem 0.2rem;
  cursor: pointer;
  border-radius: 8px;
  transition: background 0.16s;
}
.members-modal-item:hover {
  background: #222d34;
}
.modal-user-avatar {
  width: 38px;
  height: 38px;
  border-radius: 999px;
  object-fit: cover;
  background: #344;
  border: 2px solid #1d9bf0;
}
.modal-user-info {
  flex: 1 1 0;
}
.modal-user-name {
  font-weight: 600;
  color: #fff;
}
.modal-user-username {
  color: #98b4d8;
  font-size: 0.96rem;
}
.promote-btn, .demote-btn {
  background: #233344;
  color: #1d9bf0;
  border: none;
  border-radius: 0.5rem;
  font-weight: 600;
  font-size: 0.93rem;
  padding: 0.45rem 1.05rem;
  cursor: pointer;
  transition: background 0.16s, color 0.13s;
}
.promote-btn:hover, .demote-btn:hover {
  background: #1d9bf0;
  color: #fff;
}

.members-modal-pagination {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1.2rem;
  padding: 1rem 1.6rem;
  border-top: 1.5px solid #22334a;
  background: #181c20;
}

.loading-modal-members, .empty-modal-list {
  text-align: center;
  padding: 2.5rem 0;
  color: #aaa;
  font-size: 1.06rem;
}

.manage-section {
  background: #16181c;
  border-radius: 1.1rem;
  box-shadow: 0 2px 16px rgba(0,0,0,0.10);
  padding: 1.4rem 1.6rem 1.2rem 1.6rem;
  border: 1.5px solid rgba(29,155,240,0.09);
  margin-bottom: 2.2rem;
}

.manage-title {
  font-size: 1.3rem;
  font-weight: 700;
  color: #fff;
  margin-bottom: 1.1rem;
}

.pending-member-list {
  display: flex;
  flex-direction: column;
  gap: 1.2rem;
}

.pending-member-card {
  display: flex;
  align-items: center;
  background: #181818;
  border-radius: 1rem;
  padding: 1rem 1.3rem;
  box-shadow: 0 2px 8px 0 rgba(0,0,0,0.10);
  transition: box-shadow 0.15s;
}

.pending-avatar {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  object-fit: cover;
  background: #2a2a2a;
  margin-right: 1.1rem;
  border: 2px solid #1d9bf0;
}

.pending-info {
  flex: 1;
}

.pending-name {
  font-weight: 600;
  color: #fff;
  font-size: 1.05rem;
}

.pending-username {
  font-size: 0.97rem;
  color: #a5b4fc;
}

.pending-actions {
  display: flex;
  gap: 0.6rem;
}

.approve-btn,
.reject-btn {
  border: none;
  padding: 0.5rem 1.2rem;
  border-radius: 8px;
  font-weight: 600;
  font-size: 0.95rem;
  cursor: pointer;
  transition: background 0.14s, color 0.13s;
}

.approve-btn {
  background: #1d9bf0;
  color: #fff;
}
.approve-btn:hover {
  background: #177dc0;
}

.reject-btn {
  background: #222c;
  color: #f87171;
}
.reject-btn:hover {
  background: #f87171;
  color: #fff;
}

.join-community-bar {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  margin: 16px 0 12px 0;
}

.join-community-btn {
  background: #1d9bf0;
  color: #fff;
  font-weight: 700;
  border: none;
  border-radius: 20px;
  padding: 10px 26px;
  font-size: 1.04rem;
  cursor: pointer;
  box-shadow: 0 2px 8px rgba(29, 155, 240, 0.09);
  transition: background 0.16s, box-shadow 0.16s;
}

.join-community-btn:hover {
  background: #1873c5;
  box-shadow: 0 4px 18px rgba(29,155,240,0.17);
}

@media (max-width: 1250px) {
  .main-content {
    width: 98vw;
    max-width: 98vw;
    padding-left: 1vw;
    padding-right: 1vw;
    border-left: none;
    border-right: none;
  }
  .twitter-layout {
    grid-template-columns: 1fr;
    width: 100vw;
  }
  .community-info,
  .community-details {
    padding-left: 0;
    padding-right: 0;
  }
  .banner-container,
  .banner-placeholder {
    height: 100px;
    min-height: 100px;
    max-height: 100px;
  }
  .community-avatar,
  .community-avatar-placeholder {
    width: 82px;
    height: 82px;
    min-width: 82px;
    min-height: 82px;
    margin-top: -40px;
  }
  .nav-tabs {
    flex-direction: row;
    gap: 0;
    padding: 0.1rem;
    margin-bottom: 1rem;
    min-height: 44px;
    max-height: 44px;
    overflow-x: auto;
    border-radius: 0;
    background: #16181c;
  }
  .nav-tab {
    font-size: 0.98rem;
    padding: 0.7rem 0.1rem;
  }
  .community-overview-card,
  .manage-section {
    padding: 1rem 0.4rem 0.7rem 0.4rem;
  }
  .top-user-card,
  .pending-member-card {
    flex-direction: column;
    align-items: flex-start;
    padding: 0.8rem 0.8rem;
    gap: 0.7rem;
  }
  .user-avatar,
  .pending-avatar {
    margin-right: 0;
    margin-bottom: 0.4rem;
    width: 38px;
    height: 38px;
  }
  .media-thread-grid {
    grid-template-columns: 1fr;
    gap: 0.8rem;
  }
}

/* @media (max-width: 600px) {
  .main-content {
    padding: 0;
  }
  .community-name {
    font-size: 1.1rem;
  }
  .community-description {
    font-size: 0.97rem;
  }
  .community-overview-card,
  .manage-section {
    padding: 0.8rem 0.2rem 0.7rem 0.2rem;
  }
} */

</style>