<script lang="ts">
    import { onMount } from "svelte";
    import { fade } from "svelte/transition";
    import LeftSidebar from "../lib/components/LeftSideBar.svelte";
    import RightSidebar from "../lib/components/RightSideBar.svelte";
    import api from "../lib/api";
    import { clickOutside } from "../lib/action/clickOutside";
    import Thread from "../lib/components/Thread.svelte";

    type Tab = "posts" | "replies" | "likes" | "media";
    let activeTab: Tab = "posts";
    type ThreadMedia = {
      id: string;
      image_url: string;
      extension: string;
    };
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
    interface RepliesGroup {
      parent: Thread;
      replies: Thread[];
    }

    let replyGroups: RepliesGroup[] = [];

    interface GroupedMedia {
      thread: RawThread;
      media_list: ThreadMedia[];
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
    let groupedMedia: GroupedMedia[] = [];

    let currentUserId: string = "";
    let profile: ProfileResponse = {
      name: "",
      username: "",
      bio: "",
      joined_at: "",
      following_count: 0,
      follower_count: 0,
      profile_picture_id: "",
      banner_media_id: ""
    };

    let avatarUrl = "";
    let bannerUrl = "";

    let threads: Thread[] = [];
    // let replies: Thread[] = [];
    let likes: Thread[] = [];

    let showEditModal = false;
    let showPreview = false;
    let previewSrc = "";

    async function fetchProfile() {
      const res = await api.get<ProfileResponse>("/user/get-me");
      profile = res.data;
      currentUserId = res.data.id;
        if (profile.profile_picture_id) {
        const av = await api.post<{ public_url: string }>("/media/get-media", {
          id: profile.profile_picture_id
        });
        avatarUrl = av.data.public_url;
      }

      if (profile.banner_media_id) {
        const bn = await api.post<{ public_url: string }>("/media/get-media", {
          id: profile.banner_media_id
        });
        bannerUrl = bn.data.public_url;
      }
    }

    async function fetchContent() {

      if (activeTab === "posts") {
        const res = await api.get<RawThread[]>(`/users/${currentUserId}/threads`);
        threads = await Promise.all(res.data.map(enrichThread));
      }

      if (activeTab === "replies") {
        const res = await api.get<RepliesGroup[]>(`/replies/user/${currentUserId}`);
        replyGroups = await Promise.all(
          res.data.map(async (group) => ({
            parent: await enrichThread(group.parent_thread),
            replies: await Promise.all(group.replies.map(enrichThread))
          }))
        );
      }

      if (activeTab === "likes") {
        const res = await api.get<RawThread[]>(`/users/${currentUserId}/threads-like`);
        likes = await Promise.all(res.data.map(enrichThread));
      }

      if (activeTab === "media") {
        const res = await api.get<GroupedMedia[]>(`/users/${currentUserId}/media`);
        groupedMedia = await Promise.all(
          res.data.map(async (group) => ({
            thread: await enrichThread(group.thread),
            media_list: group.media_list,
          }))
        );
      }
    }

    function switchTab(tab: Tab) {
      activeTab = tab;
      fetchContent();
    }

    function openPreview(src: string) {
      previewSrc = src;
      showPreview = true;
    }

    function closePreview() {
      showPreview = false;
    }
    //Thread Util and Stuff
    async function fetchThreadMedia(threadId : string) {
      try {
        const res = await api.get<ThreadMedia[]>(`threads/${threadId}/media`);
        return res.data;
      } catch {
        return [];
      }
    }
    async function fetchLikeCount(threadId: string): Promise<number> {
      try {
        const res = await api.get<{ count: number }>(`/threads/${threadId}/likes`);
        return res.data.count;
      } catch (err) {
        console.error("Failed to fetch like count:", err);
        return 0;
      }
    }
    async function fetchLiked(threadId: string): Promise<boolean> {
      try {
        const res = await api.post<{ liked: boolean }>(
          `/threads/${threadId}/liked`,
          { user_id: currentUserId },
          { withCredentials: true }
        );
        return res.data.liked;
      } catch {
        return false;
      }
    }
    async function fetchBookmarked(threadId: string): Promise<boolean> {
      try {
        const res = await api.post<{ bookmarked: boolean }>(
          `/threads/${threadId}/bookmarked`,
          { user_id: currentUserId },
          { withCredentials: true }
        );
        return res.data.bookmarked;
      } catch {
        return false;
      }
    }
    async function fetchUserAndAvatar(userId: string) {
      const userRes = await api.get<{
        id: string;
        name: string;
        username: string;
        profile_picture_id: string;
      }>(`/user/${userId}`);
      const { name, username, profile_picture_id } = userRes.data;

      const mediaRes = await api.post<{
        id: string;
        public_url: string;
        created_at: string;
      }>("/media/get-media", { id: profile_picture_id });

      return {
        author: name,
        username,
        avatar: mediaRes.data.public_url
      };
    }
    //End of Thread Util and Stuff
    onMount(async () => {
      await fetchProfile();
      await fetchContent();
    });
    function isVideo(extension: string) {
      return [".mp4", ".webm", ".mov", ".avi", ".mkv"].includes(extension.toLowerCase());
    }

    // Store video durations so each tile can display it once loaded
    let videoDurations: Record<string, string> = {};

    // Utility to format seconds as M:SS
    function formatDuration(sec: number) {
      const m = Math.floor(sec / 60);
      const s = Math.floor(sec % 60);
      return `${m}:${s.toString().padStart(2, "0")}`;
    }

    // Called by each video thumb to load its duration
    function loadVideoDuration(media: ThreadMedia) {
      if (!isVideo(media.extension)) return;
      // Only load if not already loaded
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

    async function enrichThread(raw: RawThread): Promise<Thread> {
      const { author, username, avatar } = await fetchUserAndAvatar(raw.user_id);

      const { year, month, day, hour, minute, second } = raw.created_at;
      const utcMs = Date.UTC(year, month - 1, day, hour, minute, second);
      const dt    = new Date(utcMs);
      const createdAt = dt.toISOString();
      const media = await fetchThreadMedia(raw.id);
      const like_count = await fetchLikeCount(raw.id);
      const liked_by_me = await fetchLiked(raw.id);
      const bookmarked = await fetchBookmarked(raw.id);
      return {
        id: raw.id,
        author,
        username,
        avatar,
        content: raw.content,
        community_id: raw.community_id,
        like_count: like_count,
        liked_by_me,
        comment_count: raw.comment_count,
        share_count: raw.share_count,
        view_count: raw.view_count,
        created_at: createdAt,
        updated_at: createdAt,
        media,
        bookmarked
      };
    }
  function goToThreadDetail(threadId: string) {
    window.location.href = `/thread/${threadId}`;
  }
</script>

  <div class="profile-page layout">
    <aside class="sidebar">
        {#if currentUserId}
            <LeftSidebar {currentUserId} activePage="profile" />
        {/if}
    </aside>

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
          <small>Joined {new Date(profile.joined_at).toLocaleDateString()}</small>
          {#if profile.bio}
            <p class="bio">{profile.bio}</p>
          {/if}
          <div class="counts">
            <span><strong>{profile.following_count}</strong> Following</span>
            <span><strong>{profile.follower_count}</strong> Followers</span>
          </div>
        </div>
        <button type="button" class="edit" on:click={() => showEditModal = true}>
          Edit Profile
        </button>
      </div>

      <nav class="tabs">
        <button type="button" class:active={activeTab==="posts"} on:click={()=>switchTab("posts")}>Posts</button>
        <button type="button" class:active={activeTab==="replies"} on:click={()=>switchTab("replies")}>Replies</button>
        <button type="button" class:active={activeTab==="likes"} on:click={()=>switchTab("likes")}>Likes</button>
        <button type="button" class:active={activeTab==="media"} on:click={()=>switchTab("media")}>Media</button>
      </nav>

      <section class="content" in:fade>
        {#if activeTab === "posts"}
          {#each threads as thread (thread.id)}
            <Thread {thread} {currentUserId} on:deleted={(e) => {
              threads = threads.filter((t) => t.id !== e.detail.id);
            }} />
          {/each}
        {/if}

        {#if activeTab === "replies"}
          {#if replyGroups.length === 0}
            <p>No replies yet.</p>
          {/if}
          {#each replyGroups as group (group.parent.id)}
            <div class="reply-group-rectangle">
              <div class="parent-thread">
                <Thread thread={group.parent} {currentUserId} />
              </div>
              <div class="replies-list">
                {#each group.replies as reply (reply.id)}
                  <Thread thread={reply} {currentUserId} />
                {/each}
              </div>
            </div>
          {/each}
        {/if}
        {#if activeTab === "likes"}
          {#each likes as thread (thread.id)}
            <Thread {thread} {currentUserId} on:deleted={(e) => {
              threads = threads.filter((t) => t.id !== e.detail.id);
            }} />
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


      </section>

    </main>

    <aside class="sidebar">
      <RightSidebar />
    </aside>

    {#if showPreview}
      <div class="preview-backdrop" on:click={closePreview}>
        <img class="preview-img" src={previewSrc} alt="Preview" />
      </div>
    {/if}

    {#if showEditModal}
      <div class="modal-backdrop" use:clickOutside on:outclick={()=>showEditModal=false}>
        <div class="modal">
          <h2>Edit Profile</h2>
          <button type="button" class="save" on:click={()=>showEditModal=false}>Save</button>
        </div>
      </div>
    {/if}
  </div>

  <style>
    .layout {
      display: grid;
      grid-template-columns: auto 1fr auto;
      height: 100vh;
    }

    .main {
      overflow-y: auto;
      padding: 1rem;
      width: 40vw;
      max-width: 40vw;
      margin: 0 auto;
    }

    .banner img {
      width: 100%;
      height: 200px;
      object-fit: cover;
      cursor: pointer;
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
    .profile-page{
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

    .media-grid {
      display: grid;
      grid-template-columns: repeat(3, 1fr);
      gap: 0.5rem;
    }

    .thumb img {
      width: 100%;
      height: 100px;
      object-fit: cover;
      cursor: pointer;
    }

    .preview-backdrop {
      position: fixed;
      inset: 0;
      background: rgba(0, 0, 0, 0.8);
      display: flex;
      justify-content: center;
      align-items: center;
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
</style>