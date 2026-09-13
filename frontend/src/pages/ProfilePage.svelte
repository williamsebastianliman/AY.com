<script lang="ts">
  import { onMount } from "svelte";
  import { fade } from "svelte/transition";
  import LeftSidebar from "../lib/components/LeftSideBar.svelte";
  import RightSidebar from "../lib/components/RightSideBar.svelte";
  import BurgerLeftSideBar from "../lib/components/BurgerLeftSideBar.svelte";
  import BurgerRightSideBar from "../lib/components/BurgerRightSideBar.svelte";
  import Thread from "../lib/components/Thread.svelte";
  import api from "../lib/api";
  import { clickOutside } from "../lib/action/clickOutside";

  export let params: { id: string };
  const profileId = params.id;

  type Tab = "posts" | "replies" | "media" | "likes";
  let activeTab: Tab = "posts";

  import {
    enrichThread,
    type RawThread,
  } from "../lib/threadUtil";

  interface ProfileResponse {
    id: string;
    name: string;
    username: string;
    bio?: string;
    joined_at: string;
    following_count: number;
    follower_count: number;
    profile_picture_id: string;
    banner_media_id: string;
    is_private?: boolean;
    is_followed_by_me?: boolean;
    is_blocked_by_me?: boolean;
  }
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

  type ThreadMedia = {
    id: string;
    image_url: string;
    extension: string;
  };

  interface RepliesGroup {
    parent: Thread;
    replies: Thread[];
  }

  interface GroupedMedia {
      thread: RawThread;
      media_list: ThreadMedia[];
  }

  let currentUserId: string = "";
  let isSelf = false;
  let profile: ProfileResponse = {
    id: "",
    name: "",
    username: "",
    bio: "",
    joined_at: "",
    following_count: 0,
    follower_count: 0,
    profile_picture_id: "",
    banner_media_id: "",
    is_private: false,
    is_followed_by_me: false,
    is_blocked_by_me: false
  };

  let avatarUrl = "";
  let bannerUrl = "";
  let threads: Thread[] = [];
  let replyGroups: RepliesGroup[] = [];
  let likes: Thread[] = [];
  let showReportModal = false;
  let reportReason = "";
  let showEditModal = false;
  let showPreview = false;
  let previewSrc = "";

  let videoDurations: Record<string, string> = {};

  async function followUserAPI(targetId: string) {
    await api.post("/user/follow", {
      user_id: targetId,
      follower_id: currentUserId
    });
  }

  async function unfollowUserAPI(targetId: string) {
    await api.post("/user/unfollow", {
      user_id: targetId,
      follower_id: currentUserId
    });
  }

  async function fetchFollowerCount(userId: string): Promise<number> {
    const res = await api.get<{ count: number }>(`/user/follower_count/${userId}`);
    return res.data.count;
  }

  async function fetchFollowingCount(userId: string): Promise<number> {
    const res = await api.get<{ count: number }>(`/user/following_count/${userId}`);
    return res.data.count;
  }

  function isVideo(extension: string) {
    return [".mp4", ".webm", ".mov", ".avi", ".mkv"].includes(extension.toLowerCase());
  }
  function formatDuration(sec: number) {
    const m = Math.floor(sec / 60);
    const s = Math.floor(sec % 60);
    return `${m}:${s.toString().padStart(2, "0")}`;
  }
  function loadVideoDuration(media: ThreadMedia) {
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

  let groupedMedia: GroupedMedia[] = [];

  async function fetchMe() {
    const res = await api.get<ProfileResponse>("/user/get-me");
    currentUserId = res.data.id;
    isSelf = profileId === currentUserId;
  }
  async function fetchProfile() {
    console.log("res: ",isSelf);
    const res = isSelf
      ? await api.get<ProfileResponse>("/user/get-me")
      : await api.get<ProfileResponse>(`/user/${profileId}`);
    profile = res.data;
    profile.follower_count = await fetchFollowerCount(profile.id);
    profile.following_count = await fetchFollowingCount(profile.id);
    console.log("profile: ",profile);
    if (!isSelf && profile.id && currentUserId) {
      profile.is_followed_by_me = await fetchIsFollowed(profile.id, currentUserId);
    }
    if (profile.profile_picture_id) {
      const av = await api.post<{ public_url: string }>("/media/get-media", {
        id: profile.profile_picture_id
      });
      avatarUrl = av.data.public_url;
    } else {
      avatarUrl = "";
    }
    if (profile.banner_media_id) {
      const bn = await api.post<{ public_url: string }>("/media/get-media", {
        id: profile.banner_media_id
      });
      bannerUrl = bn.data.public_url;
    } else {
      bannerUrl = "";
    }
  }

  async function fetchContent() {
    const userId = profile.id;
    if (activeTab === "posts") {
      const res = await api.get<RawThread[]>(`/users/${userId}/threads`);
      threads = await Promise.all(res.data.map((r) => enrichThread(r, currentUserId)));
    }
    if (activeTab === "replies") {
      const userId = profile.id;
      const res = await api.get<RepliesGroup[]>(`/replies/user/${userId}`);
        replyGroups = await Promise.all(
          res.data.map(async (group) => ({
            parent: await enrichThread(group.parent_thread, currentUserId),
            replies: await Promise.all(group.replies.map((r) => enrichThread(r, currentUserId)))
          }))
        );
    }
    if (activeTab === "media") {
      const userId = profile.id;
      const res = await api.get<GroupedMedia[]>(`/users/${userId}/media`);
        groupedMedia = await Promise.all(
          res.data.map(async (group) => ({
            thread: await enrichThread(group.thread, currentUserId),
            media_list: group.media_list,
          }))
        );
    }
    if (activeTab === "likes" && isSelf) {
      const userId = profile.id;
      const res = await api.get<RawThread[]>(`/users/${userId}/threads-like`);
      likes = await Promise.all(res.data.map((r) => enrichThread(r, currentUserId)));
    }
  }
  let editProfile = {
  name: "",
  username: "",
  email: "",
  password: "",
  gender: "Male",
  date_of_birth: "",
  profile_picture_id: "",
  banner_media_id: ""
};
let updateError = "";

function openEditModal() {
  editProfile = {
    name: profile.name,
    username: profile.username,
    email: profile.email,
    password: "",
    gender: profile.gender,
    date_of_birth: formatDate(profile),
    profile_picture_id: profile.profile_picture_id || "",
    banner_media_id: profile.banner_media_id || "",
  };
  updateError = "";
  showEditModal = true;
}

function formatDate(p: ProfileResponse) {
  if (!p.joined_at) return "";
  const dob = p.date_of_birth || "";
  if (typeof dob === "string") return dob.slice(0, 10);
  if (dob instanceof Date) return dob.toISOString().slice(0, 10);
  return "";
}


let profileMediaId = profile.profile_picture_id || "";
let bannerMediaId = profile.banner_media_id || "";
let profilePreview = avatarUrl || "";
let bannerPreview = bannerUrl || "";
let profileFile: File | null = null;
let bannerFile: File | null = null;
let profileError = "";
let bannerError = "";

async function handleProfileChange(e: Event) {
  const files = (e.target as HTMLInputElement).files;
  if (!files?.[0]) return;
  profileFile = files[0];
  profilePreview = URL.createObjectURL(profileFile);
  try {
    const resp = await uploadSingleMedia(profileFile);
    profileMediaId = resp.id;
  } catch {
    profileError = "Failed to upload profile image";
  }
}

async function handleBannerChange(e: Event) {
  const files = (e.target as HTMLInputElement).files;
  if (!files?.[0]) return;
  bannerFile = files[0];
  bannerPreview = URL.createObjectURL(bannerFile);
  try {
    const resp = await uploadSingleMedia(bannerFile);
    bannerMediaId = resp.id;
  } catch {
    bannerError = "Failed to upload banner image";
  }
}

async function uploadSingleMedia(file: File) {
  const form = new FormData();
  form.append("filename", file.name);
  form.append("file", file);

  const { data } = await api.post("/media/upload", form);
  return data as { id: string; public_url: string; created_at: string };
}


async function submitProfileUpdate() {
  updateError = "";

  const dob = new Date(editProfile.date_of_birth);
  const now = new Date();
  if (now.getFullYear() - dob.getFullYear() < 13) {
    updateError = "You must be at least 13 years old.";
    return;
  }

  if (profileFile && !profileMediaId) {
    try {
      const resp = await uploadSingleMedia(profileFile);
      profileMediaId = resp.id;
    } catch {
      updateError = "Failed to upload profile image";
      return;
    }
  }
  if (bannerFile && !bannerMediaId) {
    try {
      const resp = await uploadSingleMedia(bannerFile);
      bannerMediaId = resp.id;
    } catch {
      updateError = "Failed to upload banner image";
      return;
    }
  }

  const birthYear = String(dob.getFullYear());
  const birthMonth = String(dob.getMonth() + 1);
  const birthDay = String(dob.getDate());

  const payload = {
    user_id: profile.id,
    name: editProfile.name,
    username: editProfile.username,
    email: editProfile.email,
    password: editProfile.password || "",
    gender: editProfile.gender,
    birth_year: birthYear,
    birth_month: birthMonth,
    birth_day: birthDay,
    profile_picture_id: profileMediaId || profile.profile_picture_id || "",
    banner_media_id: bannerMediaId || profile.banner_media_id || "",
  };

  try {
    await api.patch("/user/update-profile", payload);
    await fetchProfile();
    showEditModal = false;
    profileFile = null;
    bannerFile = null;
  } catch (e) {
    updateError = e?.response?.data?.error || "Failed to update profile";
  }
}

async function blockUser() {
  try{
    api.post("/user/block", {user_id: currentUserId, target_id: params.id});
    profile.is_blocked_by_me = true;
  }catch(err){
    console.error(err as string);
  }
}






  onMount(async () => {
    await fetchMe();
    await fetchProfile();
    await isBlocked();
  });
  $: if (profile && profile.id && activeTab) {
    fetchContent();
  }

  function switchTab(tab: Tab) {
    activeTab = tab;
  }
  function openPreview(src: string) {
    previewSrc = src;
    showPreview = true;
  }
  function closePreview() {
    showPreview = false;
  }
  async function followUser() {
    if (profile.is_followed_by_me) {
      await unfollowUserAPI(profile.id);
      profile.is_followed_by_me = false;
    } else {
      await followUserAPI(profile.id);
      profile.is_followed_by_me = true;
    }
    profile.follower_count = await fetchFollowerCount(profile.id);
    profile.following_count = await fetchFollowingCount(profile.id);
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
  function openReportModal() {
    showReportModal = true;
  }
  function submitReport(reported_user: string, reason: string) {
    try{
      api.post("/user-reports", {reported_user: reported_user, reason: reason});
    }
    catch(err){
      console.error(err as string);
    }
    showReportModal = false;
    reportReason = "";
  }
  function goToThreadDetail(threadId: string) {
    window.location.href = `/thread/${threadId}`;
  }
  async function isBlocked() {
    try{
      const res = await api.post("/user/is-blocked", {user_id: currentUserId, target_id: params.id});
      profile.is_blocked_by_me = res.data.is_blocked;
    }catch(err){
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

<div class="profile-page layout">
  <aside class="sidebar">
    {#if currentUserId && windowWidth >= 1250}
      <LeftSidebar {currentUserId} activePage="profile" />
    {/if}
  </aside>
  {#if currentUserId && windowWidth < 1250}
    <BurgerLeftSideBar {currentUserId} activePage="profile" />
  {/if}
  <main class="main">
    <div class="banner" on:click={() => openPreview(bannerUrl)}>
      <img src={bannerUrl} alt="Profile Banner" />
    </div>
    <div class="avatar" on:click={() => openPreview(avatarUrl)}>
      <img src={avatarUrl} alt="Profile Avatar" />
    </div>
    <div class="header">
      <div class="info">
        <h1>{profile.name}</h1>
        <p class="handle">@{profile.username}</p>
        <small>Joined {profile.joined_at && new Date(profile.joined_at).toLocaleDateString()}</small>
        {#if profile.bio}
          <p class="bio">{profile.bio}</p>
        {/if}
        <div class="counts">
          <span><strong>{profile.following_count}</strong> Following</span>
          <span><strong>{profile.follower_count}</strong> Followers</span>
        </div>
      </div>
      <div class="actions">
        {#if isSelf}
          <button type="button" class="edit" on:click={openEditModal}>Edit Profile</button>
        {:else if !profile.is_blocked_by_me}
          <button type="button" class="follow-btn" on:click={followUser}>
            {profile.is_followed_by_me ? "Unfollow" : "Follow"}
          </button>
          <button type="button" class="report-btn" on:click={openReportModal}>Report</button>
          <button type="button" class="block-btn" on:click={blockUser}>Block</button>
        {:else}
          <span class="blocked-text">You have blocked this user</span>
        {/if}
      </div>
    </div>
    <nav class="tabs">
      <button type="button" class:active={activeTab === "posts"} on:click={() => switchTab("posts")}>Posts</button>
      <button type="button" class:active={activeTab === "replies"} on:click={() => switchTab("replies")}>Replies</button>
      <button type="button" class:active={activeTab === "media"} on:click={() => switchTab("media")}>Media</button>
      {#if isSelf}
        <button type="button" class:active={activeTab === "likes"} on:click={() => switchTab("likes")}>Likes</button>
      {/if}
    </nav>
    <section class="content" in:fade>
      {#if profile.is_blocked_by_me && !isSelf}
        <div class="blocked-banner">You have blocked this user. Their content is hidden.</div>
      {:else if profile.is_private && !profile.is_followed_by_me && !isSelf}
        <div class="private-banner">This user's profile is private. Follow to see their posts.</div>
      {:else}
        {#if activeTab === "posts"}
          {#each threads as thread (thread.id)}
            <Thread is_profile_thread={true} {thread} {currentUserId} on:deleted={(e) => {
          threads = threads.filter((t) => t.id !== e.detail.id);}} on:pinned={fetchContent} />
          {/each}
        {/if}
        {#if activeTab === "replies"}
          {#each replyGroups as group (group.parent.id)}
            <div class="reply-group-rectangle">
              <div class="parent-thread">
                <Thread is_profile_thread={true} thread={group.parent} {currentUserId} />
              </div>
              <div class="replies-list">
                {#each group.replies as reply (reply.id)}
                  <Thread is_profile_thread={true} thread={reply} {currentUserId} />
                {/each}
              </div>
            </div>
          {/each}
        {/if}
        {#if activeTab === "media"}
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
          </div>
        {/if}
        {#if activeTab === "likes" && isSelf}
          {#each likes as thread (thread.id)}
            <Thread {thread} {currentUserId} />
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
  {#if currentUserId && windowWidth < 1460}
    <BurgerRightSideBar {currentUserId} activePage="home" />
  {/if}

  {#if showPreview}
    <div class="preview-backdrop" on:click={closePreview}>
      <img class="preview-img" src={previewSrc} alt="Preview" />
    </div>
  {/if}
  {#if showEditModal}
    <div class="modal-backdrop" use:clickOutside on:outclick={() => showEditModal = false}>
      <div class="modal">
        <button type="button" class="close-btn" on:click={() => showEditModal = false}>✕</button>
        <h2>Edit Profile</h2>
        <form on:submit|preventDefault={submitProfileUpdate}>
          <label>
            Name
            <input type="text" bind:value={editProfile.name} required minlength="5" pattern="[A-Za-z\s]+" />
          </label>
          <label>
            Username
            <input type="text" bind:value={editProfile.username} required minlength="4" />
          </label>
          <label>
            Email
            <input type="email" bind:value={editProfile.email} required pattern="^[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.com$" />
          </label>
          <label>
            Password
            <input type="password" bind:value={editProfile.password} placeholder="Leave blank to keep current" minlength="8" />
          </label>
          <label>
            Gender
            <select bind:value={editProfile.gender} required>
              <option value="Male">Male</option>
              <option value="Female">Female</option>
            </select>
          </label>
          <label>
            Date of Birth
            <input type="date" bind:value={editProfile.date_of_birth} required />
          </label>
          <label>
            Profile Picture
            <input type="file" accept="image/*" on:change={handleProfileChange} />
            {#if profilePreview}
              <img src={profilePreview} alt="" class="preview" />
            {/if}
            {#if profileError}<div class="error">{profileError}</div>{/if}
          </label>
          <label>
            Banner Picture
            <input type="file" accept="image/*" on:change={handleBannerChange} />
            {#if bannerPreview}
              <img src={bannerPreview} alt="" class="preview" />
            {/if}
            {#if bannerError}<div class="error">{bannerError}</div>{/if}
          </label>
          <div class="edit-actions">
            <button type="submit" class="save">Save</button>
          </div>
          {#if updateError}
            <div class="update-error">{updateError}</div>
          {/if}
        </form>
      </div>
    </div>
  {/if}
  {#if showReportModal}
    <div class="modal-backdrop" use:clickOutside on:outclick={() => showReportModal = false}>
      <div class="modal">
        <h2>Report User</h2>
        <textarea
          class="report-textarea"
          placeholder="Enter your reason for reporting this account"
          bind:value={reportReason}
        ></textarea>
        <button type="button" class="submit-btn" on:click={submitReport(params.id,reportReason)}>Submit Report</button>
      </div>
    </div>
  {/if}
</div>

<style>
  .actions {
    display: flex;
    gap: 0.75rem;
    align-items: flex-start;
  }
  .follow-btn, .report-btn, .block-btn {
    border: none;
    border-radius: 9999px;
    padding: 0.5rem 1.25rem;
    font-weight: 600;
    cursor: pointer;
  }
  .follow-btn:hover {
    background: #000;
    color: #fff;
    outline: 2.5px solid #1d9bf0;
    outline-offset: 2px;
    transition: background 0.15s, color 0.15s, outline 0.1s;
  }
  .report-btn:hover {
    background: #000;
    color: #fff;
    outline: 2.5px solid #1d9bf0;
    outline-offset: 2px;
    transition: background 0.15s, color 0.15s, outline 0.1s;
  }
  .block-btn:hover {
    background: #000;
    color: #fff;
    outline: 2.5px solid #1d9bf0;
    outline-offset: 2px;
    transition: background 0.15s, color 0.15s, outline 0.1s;
  }
  .layout {
  display: grid;
  grid-template-columns: auto 1fr auto;
  height: 100vh;
  z-index: 1;
}
.main {
  overflow-y: auto;
  padding: 1rem;
  width: 45vw;
  max-width: 45vw;
  margin: 0 auto;
  z-index: 2;
}
.banner img {
  width: 100%;
  height: 200px;
  object-fit: cover;
  cursor: pointer;
  border-radius: 1.1rem 1.1rem 0 0;
}
.avatar {
  position: relative;
  top: -40px;
  margin-left: 1rem;
  cursor: pointer;
}
.avatar img {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  border: 3px solid white;
}
.profile-page {
  z-index: 1;
  overflow-y: hidden;
}
.header {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  margin-top: -1.5rem;
}
.info h1 {
  margin: 0;
}
.handle {
  color: #888;
}
.bio {
  margin: 0.5rem 0;
}
.counts span {
  margin-right: 1rem;
}
.edit {
  padding: 0.5rem 1rem;
  border-radius: 9999px;
  background: #1d9bf0;
  color: white;
  cursor: pointer;
}
.tabs {
  display: flex;
  border-bottom: 1px solid #444;
  margin: 1rem 0;
}
.tabs button {
  flex: 1;
  background: none;
  border: none;
  padding: 0.75rem;
  color: #888;
  cursor: pointer;
}
.tabs button.active {
  color: white;
  border-bottom: 2px solid #1d9bf0;
}
.content {
  min-height: 60vh;
  padding-bottom: 10vh;
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
.blocked-banner, .private-banner {
  text-align: center;
  margin: 2rem 0;
  padding: 2rem;
  background: rgba(220,38,38,0.08);
  color: #f14e67;
  border-radius: 0.5rem;
  font-size: 1.1rem;
}
.preview-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.8);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 10;
}
.preview-img {
  max-width: 90%;
  max-height: 90%;
}
.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 10;
}
.modal {
  background: #111;
  padding: 1.5rem;
  border-radius: 0.5rem;
  width: 90%;
  max-width: 400px;
}
.save {
  padding: 0.5rem 1rem;
  border-radius: 0.5rem;
  background: #1d9bf0;
  color: white;
  cursor: pointer;
}
.report-textarea {
  width: 100%;
  min-height: 4rem;
  margin-bottom: 1rem;
  padding: 0.5rem;
  border-radius: 0.3rem;
  border: 1px solid #555;
}
.submit-btn {
  background: #1d9bf0;
  color: white;
  border-radius: 9999px;
  padding: 0.5rem 1.25rem;
  border: none;
  font-weight: 600;
  cursor: pointer;
}
.blocked-text { color: #f14e67; font-weight: bold; }

@media (max-width: 900px) {
  .main { width: 98vw; max-width: 98vw; }
  .media-thread-grid { grid-template-columns: repeat(2, 1fr);}
  .banner img { height: 120px;}
  .avatar img { width: 56px; height: 56px;}
}
@media (max-width: 900px) {
      .reply-group-rectangle { padding: 0.5rem; }
      .parent-thread, .replies-list { margin-left: 0; }
    }
    .reply-group-rectangle {
    background: rgb(18, 18, 18);
    border-radius: 1.1rem;
    box-shadow: 0 2px 16px rgba(0,0,0,0.10);
    margin-bottom: 2rem;
    padding: 1.1rem 1rem 0.5rem 1rem;
    border: 1.5px solid rgba(29,155,240,0.09);
    transition: box-shadow 0.14s;
  }

  .reply-group-rectangle:hover {
    box-shadow: 0 4px 20px rgba(29,155,240,0.10), 0 2px 8px rgba(0,0,0,0.17);
  }

  .parent-thread {
    margin-bottom: 0.6rem;
  }

  .replies-list {
    margin-top: 0.4rem;
    padding-left: 1.7rem;
    border-left: 2.5px solid rgb(29, 156,   240);
    display: flex;
    flex-direction: column;
    gap: 1.1rem;
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
  .sidebar{
    z-index: 400000;
  }

  .modal .close-btn {
    background: none;
    border: none;
    font-size: 1.5rem;
    position: absolute;
    top: 10px;
    right: 15px;
    color: #fff;
    cursor: pointer;
    z-index: 2;
  }
  .modal form {
    display: flex;
    flex-direction: column;
    gap: 0.8rem;
  }
  .modal label {
    display: flex;
    flex-direction: column;
    font-weight: 600;
    color: #dedede;
    margin-bottom: 0.4rem;
  }
  .modal input, .modal select, .modal textarea {
    padding: 0.4rem;
    margin-top: 0.12rem;
    background: #222;
    color: #fff;
    border-radius: 0.35rem;
    border: 1px solid #444;
  }
  .edit-actions {
    display: flex;
    justify-content: flex-end;
    gap: 0.7rem;
    margin-top: 1rem;
  }
  .update-error {
    color: #f14e67;
    font-weight: bold;
    margin-top: 0.8rem;
    text-align: center;
  }
  .preview {
    margin-top: 0.5rem;
    width: 90px;
    border-radius: 0.4rem;
  }
  /* @media (max-width: 1200px){
    *{
      overflow-x: hidden;
    }
  } */

</style>