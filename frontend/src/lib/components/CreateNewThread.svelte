<script lang="ts">
    import { createEventDispatcher, onDestroy, onMount } from "svelte";
    import { Image, Video, BarChart2, Smile, Calendar, MapPin, X } from "lucide-svelte";
    import DonutProgress from "./DonutProgress.svelte";
    import api from "../api";

    export let myAvatar: string = "";
    export let myUserId: string = "";
    export let communities: { id: string; name: string }[] = [];
    const dispatch = createEventDispatcher();
    let content = "";
    let posting = false;
    let postError = false;
    let fileInput: HTMLInputElement;
    let images: { id: string; url: string; type: string }[] = [];
    let selectedCommunity = "";
    let postAsCommunity = false;
    const maxWords = 280;
    $: wordCount = content.length;
    async function loadMyProfile() {
    try {
        const meRes = await api.get("/user/get-me");
        myUserId = meRes.data.id;
        if (meRes.data.profile_picture_id){
            const mediaRes = await api.post<{
            id: string;
            public_url: string;
            }>(
            "/media/get-media",
            { id: meRes.data.profile_picture_id },
            { withCredentials: true }
            );
            myAvatar = mediaRes.data.public_url;
        }
        } catch (err) {
        console.error("Failed to load my profile:", err);
        }
    }
    onMount(async () => {
      await loadMyProfile();
      await getUserCommunities();
    });
    function openFileInput() {
      if (images.length < 4) fileInput.click();
    }
    async function handleFiles(e: Event) {
      const input = e.target as HTMLInputElement;
      if (!input.files?.length) return;
      const file = input.files[input.files.length - 1];
      if (images.length >= 4) {
        input.value = "";
        return;
      }
      try {
        const form = new FormData();
        form.append("file", file);
        const res = await api.post<{ public_url: string; id: string }>(
          "/media/upload",
          form,
          { headers: { "Content-Type": "multipart/form-data" } }
        );
        images = [
          ...images,
          { id: res.data.id, url: res.data.public_url, type: file.type }
        ];
      } catch (err) {
        console.error(err as string);
      } finally {
        input.value = "";
      }
    }
    async function getUserCommunities() {
      try{
        const resData = await api.get(`/users/communities/joined/${myUserId}?page=${1}&size=${1000}`);
        communities = resData.data;
        console.log(communities);
      }
      catch(err){
        console.error(err as string);
      }
    }
    function removeLocalMedia(id: string) {
      images = images.filter(img => img.id !== id);
    }
    function closeModal() {
      dispatch("close");
    }
    async function postThread() {
      if (!content.trim()) return;
      posting = true;
      postError = false;
      try {
        const resp = await api.post<{ id: string }>(
          "/threads",
          {
            user_id: myUserId,
            content,
            community_id: postAsCommunity ? selectedCommunity : null,
            parent_id: null
          },
          { withCredentials: true }
        );
        const threadId = resp.data.id;
        await Promise.all(
          images.map(img =>
            api.post(`/threads/${threadId}/media`, { media_id: img.id })
          )
        );
        content = "";
        images = [];
        closeModal();
      } catch (err) {
        console.error(err as string);
        postError = true;
      } finally {
        posting = false;
      }
    }
    onMount(() => {
        document.body.classList.add("modal-open");
    });
    onDestroy(() => {
        document.body.classList.remove("modal-open");
    });
  </script>
  <div class="thread-modal-backdrop" on:click={closeModal}></div>
  <div class="thread-modal" on:click|stopPropagation>
    <div class="modal-header">
      <button class="close-modal-btn" type="button" on:click={closeModal} aria-label="Close">
        <X size="26" />
      </button>
      <select class="community-dropdown"
        bind:value={selectedCommunity}
        on:change={() => postAsCommunity = !!selectedCommunity}
      >
        <option value="">My Account</option>
        {#each communities as c (c.community_id)}
          <option value={c.community_id}>{c.community_name}</option>
        {/each}
      </select>
    </div>
    <div class="compose-header">
      <img src={myAvatar} alt="You" class="compose-avatar" />
      <textarea
        class="compose-input"
        bind:value={content}
        placeholder="What’s happening?"
        rows={5}
        maxlength={maxWords * 6}
        autofocus
      />
    </div>
    {#if images.length}
      <div class="image-preview images-{images.length}">
        {#each images as img (img.id)}
          <div class="thumb">
            {#if img.type.startsWith("video")}
              <video src={img.url} controls preload="metadata" />
            {:else}
              <img src={img.url} alt="upload" />
            {/if}
            <button type="button"
              class="del-btn"
              aria-label="Remove media"
              on:click={() => removeLocalMedia(img.id)}
            ><X size="18" /></button>
          </div>
        {/each}
      </div>
    {/if}
    <div class="compose-footer">
      <div class="footer-left">
        <DonutProgress value={wordCount} max={maxWords} size={38} color="#1d9bf0" />
        <span class="word-count">{wordCount} / {maxWords}</span>
      </div>
      <div class="icon-row">
        <button on:click={openFileInput} type="button" class="icon-btn media-btn" aria-label="Add image"><Image size="20" /></button>
        <button on:click={openFileInput} type="button" class="icon-btn media-btn" aria-label="Add video"><Video size="20" /></button>
        <button type="button" class="icon-btn" aria-label="Poll"><BarChart2 size="20" /></button>
        <button type="button" class="icon-btn" aria-label="Emoji"><Smile size="20" /></button>
        <button type="button" class="icon-btn" aria-label="Schedule"><Calendar size="20" /></button>
        <button type="button" class="icon-btn" aria-label="Location"><MapPin size="20" /></button>
        <input
          type="file"
          accept="image/*,video/*"
          multiple
          bind:this={fileInput}
          on:change={handleFiles}
          class="hidden"
        />
      </div>
      <button
        class="post-btn"
        type="button"
        disabled={!content.trim() || posting}
        on:click={postThread}
      >
        {#if posting}Posting…{:else}Post{/if}
      </button>
    </div>
    {#if postError}
      <p class="error">Failed to post. Please try again.</p>
    {/if}
  </div>
  <style>
    .thread-modal-backdrop {
      position: fixed;
      inset: 0;
      background: rgba(0,0,0,0.8);
      z-index: 200000;
    }
    .thread-modal {
      position: fixed;
      left: 50%;
      top: 50%;
      transform: translate(-50%,-50%);
      background: #15181c;
      color: #fff;
      min-width: 400px;
      max-width: 550px;
      width: 98vw;
      border-radius: 1.2rem;
      box-shadow: 0 4px 40px rgb(0, 0, 0);
      z-index: 200001;
      padding: 0 0 1.5rem 0;
      animation: modal-in 0.2s cubic-bezier(.4,1.7,.67,.99);
      display: flex;
      flex-direction: column;
      gap: 0;
    }
    @keyframes modal-in {
      from { transform: translate(-50%, -44%) scale(.97);}
      to   { transform: translate(-50%, -50%) scale(1);}
    }
    .modal-header {
      display: flex;
      align-items: center;
      padding: 1.2rem 1.5rem 0.2rem 1.5rem;
      justify-content: space-between;
    }
    .close-modal-btn {
      background: none;
      border: none;
      color: #e3e3e3;
      font-size: 1.25rem;
      cursor: pointer;
      padding: 0.2rem 0.3rem 0.2rem 0.1rem;
      border-radius: 9999px;
      transition: background .12s;
    }
    .close-modal-btn:hover { background: #232327; }
    .community-dropdown {
      border-radius: 1rem;
      border: none;
      background: #181d23;
      color: #5ac6ff;
      font-size: 1.01rem;
      padding: 0.32rem 1rem;
      outline: none;
      font-weight: 500;
      cursor: pointer;
    }
    .compose-header {
      display: flex;
      gap: 1rem;
      padding: 1.1rem 1.5rem 0.3rem 1.5rem;
      align-items: flex-start;
    }
    .compose-avatar {
      width: 48px;
      height: 48px;
      border-radius: 50%;
      object-fit: cover;
      margin-top: 0.1rem;
    }
    .compose-input {
      flex: 1;
      background: transparent;
      border: none;
      resize: none;
      color: #f9f9f9;
      font-size: 1.16rem;
      line-height: 1.6;
      min-height: 6rem;
      height: 8rem;
      outline: none;
      padding: 0.3rem 0;
      font-family: inherit;
      min-height: 4rem;
      height: 8rem;
    }
    .compose-input::placeholder {
      color: #5c6a7b;
      font-size: 1.09rem;
    }
    .image-preview {
      display: flex;
      gap: 0.7rem;
      margin-top: 0.5rem;
      padding-left: 1.5rem;
      padding-right: 1.5rem;
      overflow-x: auto;
      overflow-y: hidden;
      align-items: flex-start;
    }
    .thumb {
      position: relative;
      border-radius: 0.9rem;
      overflow: hidden;
      background: #191b1e;
      display: flex;
      align-items: center;
      justify-content: center;
      min-width: 120px;
      max-width: 140px;
      aspect-ratio: 16/10;
      min-height: 84px;
      box-shadow: 0 2px 6px #0005;
    }
    .thumb img, .thumb video {
      width: 100%;
      height: 100%;
      object-fit: cover;
      background: #202023;
      border-radius: 0.9rem;
      display: block;
    }
    .del-btn {
      position: absolute;
      top: 7px;
      right: 7px;
      background: #000d;
      border: none;
      color: #fff;
      border-radius: 9999px;
      width: 26px;
      height: 26px;
      display: flex;
      align-items: center;
      justify-content: center;
      cursor: pointer;
      font-size: 1.15rem;
      padding: 0;
      z-index: 200001;
      box-shadow: 0 2px 8px #0007;
      transition: background .13s;
    }
    .del-btn:hover { background: #2229; }
    .compose-footer {
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 1.2rem 1.5rem 0 1.5rem;
      margin-top: 0.6rem;
      gap: 1.2rem;
    }
    .footer-left {
      display: flex;
      align-items: center;
      gap: 0.3rem;
      min-width: 120px;
    }
    .word-count {
      color: #aaa;
      font-size: 0.93rem;
      margin-left: 0.3rem;
      min-width: 68px;
      text-align: left;
    }
    .icon-row {
      display: flex;
      gap: 0.7rem;
      align-items: center;
    }
    .icon-btn {
      background: none;
      border: none;
      color: #1d9bf0;
      cursor: pointer;
      padding: 0.3rem;
      border-radius: 9999px;
      font-size: 1.07rem;
      transition: color .16s, background .16s;
    }
    .icon-btn:hover { background: #163751; color: #32bbff; }
    .post-btn {
      background: #1d9bf0;
      border: none;
      color: #fff;
      padding: 0.56rem 1.5rem;
      border-radius: 9999px;
      font-weight: 700;
      font-size: 1.06rem;
      box-shadow: 0 1px 5px #22d2;
      cursor: pointer;
      margin-left: 1rem;
      transition: background .15s;
    }
    .post-btn:disabled { opacity: 0.5; cursor: not-allowed; }
    .post-btn:hover:not(:disabled) { background: #299fff; }
    .hidden { display: none !important; }
    .error {
      color: #ffb1b1;
      background: transparent;
      font-size: 0.95rem;
      padding: 0.6rem 1.5rem 0 1.5rem;
    }
    @media (max-width: 600px) {
      .thread-modal {
        min-width: unset;
        max-width: 90vw;
        width: 90vw;
        border-radius: 0.6rem;
        padding: 0 0 0.5rem 0;
      }
      .modal-header,
      .compose-header
       {
        padding-left: 0.7rem;
        padding-right: 0.7rem;
        max-width: 60vw;
      }
      .icon-btn:not(.media-btn) {
        display: none !important;
      }
      .compose-footer {
        flex-wrap: wrap;
        gap: 0.6rem;
        padding-bottom: 0.7rem;
      }
      .compose-avatar {
        width: 36px;
        height: 36px;
      }
      .compose-input {
        font-size: 1rem;
        min-height: 2.4rem;
        height: 4.5rem;
      }
      .image-preview {
        padding-left: 0.7rem;
        padding-right: 0.7rem;
        gap: 0.4rem;
      }
      .thumb {
        min-width: 70px;
        max-width: 85px;
        min-height: 50px;
        border-radius: 0.45rem;
      }
      .footer-left {
        min-width: 80px;
      }
      .word-count {
        font-size: 0.82rem;
        min-width: 40px;
      }
      .icon-btn {
        font-size: 0.94rem;
        padding: 0.13rem;
      }
      .post-btn {
        padding: 0.37rem 0.7rem;
        font-size: 0.96rem;
        margin-left: 0.3rem;
      }
    }

  </style>