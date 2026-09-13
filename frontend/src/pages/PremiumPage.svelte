<script lang="ts">
    import { onMount } from "svelte";
    import LeftSidebar from "../lib/components/LeftSideBar.svelte";
    import RightSidebar from "../lib/components/RightSideBar.svelte";
    import BurgerLeftSideBar from "../lib/components/BurgerLeftSideBar.svelte";
    import BurgerRightSideBar from "../lib/components/BurgerRightSideBar.svelte";
    import { BadgeCheck } from "lucide-svelte";
    import api from "../lib/api";
    let facePhoto: File | null = null;
    let facePreview: string | null = null;
    let facePhotoId: string;
    async function handlePhotoChange(e: Event) {
      const input = e.target as HTMLInputElement;
      facePhoto = input.files && input.files.length ? input.files[0] : null;
      facePreview = null;
      if (facePhoto) {
          const form = new FormData();
          form.append("file", facePhoto);

          try {
              const uploadRes = await api.post<{ id: string }>(
                  "/media/upload",
                  form,
                  { headers: { "Content-Type": "multipart/form-data" } }
              );

              const mediaId = uploadRes.data.id;
              const getMediaRes = await api.post<{ id: string; public_url: string; extension: string}>(
                  "/media/get-media",
                  { id: mediaId }
              );
              facePreview = getMediaRes.data.public_url;

              facePhotoId = mediaId;
          } catch (err) {
              console.error("Failed to upload or fetch image:", err);
              facePhoto = null;
              facePreview = null;
          }
      }
    }

    function removePhoto() {
        facePhoto = null;
        facePreview = null;
    }

    let me: { id: string; name: string; username: string; is_premium: boolean; premium_pending: boolean } | null = null;
    let idNumber = "";
    let reason = "";
    let submitting = false;
    let submitError = "";
    let submitSuccess = false;
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
    onMount(async () => {
      try {
        const resp = await api.get("/user/get-me");
        me = resp.data;
        if (!me?.id) {
          window.location.href = "/login";
        }
      } catch {
        window.location.href = "/login";
      }
    });
    async function handleSubmit() {
      submitError = "";
      submitSuccess = false;

      if (!idNumber.trim() || !reason.trim() || !facePhoto) {
        submitError = "All fields are required.";
        return;
      }

      submitting = true;
      try {

        await api.post("/user/request-premium", {
          user_id: me.id,
          card_number: idNumber.trim(),
          reason: reason.trim(),
          face_image: facePhotoId
        });

        submitSuccess = true;
        me.premium_pending = true;
      } catch (err) {
        console.error(err as string);
        submitError = "Failed to submit verification request. You already request it!";
      }
      submitting = false;
    }

    let myUserId = "";

    async function loadMyProfile() {
    try {
        const meRes = await api.get<MeResponse>("/user/get-me");
        const me = meRes.data;

        myUserId = me.id;
    } catch (err) {
        console.error("Failed to load my profile:", err);
    }
    }
    onMount(async () => {
        await loadMyProfile();
    });
    let windowWidth = window.innerWidth;
    function handleResize() {
      windowWidth = window.innerWidth;
    }
    onMount(() => {
      window.addEventListener("resize", handleResize);
      return () => window.removeEventListener("resize", handleResize);
    });
  </script>
  <div class="premium-layout">
    <aside class="sidebar">
      {#if myUserId && windowWidth >= 1250}
        <LeftSidebar currentUserId={myUserId} activePage="premium" />
      {/if}
    </aside>
    {#if myUserId && windowWidth < 1250}
      <BurgerLeftSideBar currentUserId={myUserId} activePage="premium" />
    {/if}
    <main class="premium-main">
      <h1>Premium Verification</h1>
      <div class="benefit-card">
        <BadgeCheck size={24} class="text-sky-500" />
        <span>
          <b>Premium users get a blue checkmark</b>
          <BadgeCheck size={18} class="inline text-sky-500 mx-1" />
          next to their name for increased recognition!
        </span>
      </div>
      {#if me?.is_premium}
        <div class="success">You are already a Premium member! <span class="blue-check">✔</span></div>
      {:else if me?.premium_pending || submitSuccess}
        <div class="pending">Your verification request is pending. Please wait for admin approval.</div>
      {:else}
        <form class="verify-form" on:submit|preventDefault={handleSubmit}>
          <label>
            National Identity Card Number
            <input type="text" bind:value={idNumber} maxlength="32" required />
          </label>
          <label>
            Why do you want to be verified?
            <textarea bind:value={reason} required rows="3"></textarea>
          </label>
          <label>
            Upload a recent photo of your face
            {#if facePreview}
              <div class="preview-wrapper">
                <div class="preview-container">
                  <img src={facePreview} alt="Preview" class="preview-image" />
                  <button type="button" class="remove-btn" on:click={removePhoto}>×</button>
                </div>
              </div>
            {:else}
              <input type="file" accept="image/*" on:change={handlePhotoChange} required />
            {/if}
          </label>
          {#if submitError}
            <div class="error">{submitError}</div>
          {/if}
          <button type="submit" class="submit-btn" disabled={submitting}>
            {submitting ? "Submitting..." : "Submit Verification"}
          </button>
        </form>
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
    :global(html, body) {
    height: 100%;
    margin: 0;
    overflow: hidden;
    }
    .premium-layout {
        display: grid;
        grid-template-columns: auto 1fr auto;
        overflow-y: auto;
        scrollbar-width: none;
        -ms-overflow-style: none;
    }
    .premium-main {
        max-width: 630px;
        margin: 4vh auto 0 auto;
        padding: 2rem;
        background: rgba(30,40,55,0.85);
        border-radius: 1.25rem;
        box-shadow: 0 2px 24px #0002;
        max-height: 90vh;
    }
    .benefit-card {
        display: flex;
        align-items: center;
        background: #121f36;
        color: #88c1ff;
        padding: 1rem 1.3rem;
        border-radius: 0.7rem;
        gap: 1rem;
        margin-bottom: 2rem;
        font-size: 1.1rem;
    }
    .checkmark, .blue-check {
        color: #1da1f2;
        font-size: 1.6em;
        margin-right: 0.4em;
    }
    .verify-form {
        display: flex;
        flex-direction: column;
        gap: 1.4rem;
    }
    .verify-form label {
        font-weight: 500;
        color: #eee;
        display: flex;
        flex-direction: column;
        gap: 0.4rem;
    }
    .verify-form input[type="text"], .verify-form textarea {
        padding: 0.55rem 1rem;
        border-radius: 0.6rem;
        border: 1px solid #245490;
        background: #181e29;
        color: #f7f7f7;
    }
    .verify-form input[type="file"] {
        background: none;
        color: #f7f7f7;
        margin-top: 0.2rem;
    }
    .submit-btn {
        background: #1da1f2;
        color: #fff;
        padding: 0.7rem 0;
        border: none;
        border-radius: 0.7rem;
        font-weight: bold;
        font-size: 1.12rem;
        margin-top: 0.3rem;
        transition: background 0.18s;
    }
    .submit-btn:disabled {
        background: #999;
        cursor: not-allowed;
    }
    .error {
        color: #ff5462;
        font-weight: 500;
    }
    .success, .pending {
        padding: 1.4rem 1rem;
        border-radius: 0.7rem;
        background: #193254;
        color: #36e2a0;
        font-weight: 600;
        font-size: 1.1rem;
        text-align: center;
    }
    .pending { color: #ffd36e; }
    .preview-wrapper {
    display: flex;
    justify-content: center;
    margin-top: 1rem;
    }

    .preview-container {
    position: relative;
    display: inline-block;
    }

    .preview-image {
    width: 180px;
    height: 180px;
    object-fit: cover;
    border-radius: 1rem;
    border: 3px solid #1da1f2;
    box-shadow: 0 0 6px #1da1f288;
    }

    .remove-btn {
        position: absolute;
        top: -14px;
        right: -14px;
        width: 32px;
        height: 32px;
        background: #ff5462;
        color: white;
        font-size: 1.2rem;
        border: none;
        border-radius: 9999px;
        cursor: pointer;
        display: flex;
        align-items: center;
        justify-content: center;
        box-shadow: 0 2px 6px #0003;
        transition: background 0.2s ease;
    }

    .remove-btn:hover {
    background: #e63c4b;
    }

    @media (max-width: 900px) {
  .premium-layout {
    grid-template-columns: 1fr;
    padding: 0 !important;
    min-height: 100vh;
    height: auto;
  }
  .premium-main {
    max-width: 90vw;
    min-width: 0;
    padding: 1.2rem 0.7rem;
    margin: 2vh auto 0 auto;
    border-radius: 0.7rem;
    box-shadow: 0 1px 8px #0002;
    max-height: unset;
    overflow-y: visible;
  }
  .benefit-card {
    font-size: 1rem;
    gap: 0.6rem;
    padding: 0.6rem 0.6rem;
    margin-bottom: 1.2rem;
  }
}

@media (max-width: 600px) {
  .preview-image {
    width: 90px;
    height: 90px;
  }
  .submit-btn {
    font-size: 1rem;
    padding: 0.45rem 0;
  }
}

.premium-layout {
  height: 100vh;
  min-height: 100vh;
  overflow-y: auto;
}

</style>