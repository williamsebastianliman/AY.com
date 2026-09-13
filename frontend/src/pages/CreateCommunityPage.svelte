<script lang="ts">
    import { onMount } from "svelte";
    import LeftSidebar from "../lib/components/LeftSideBar.svelte";
    import RightSidebar from "../lib/components/RightSideBar.svelte";
    import BurgerLeftSideBar from "../lib/components/BurgerLeftSideBar.svelte";
    import BurgerRightSideBar from "../lib/components/BurgerRightSideBar.svelte";
    import api from "../lib/api";
    import { fade } from "svelte/transition";
    import { Upload, X, Plus, AlertCircle } from "lucide-svelte";

    interface MeResponse {
      id: string;
      name: string;
      username: string;
      email?: string;
      profile_picture_id: string;
    }

    let myUserId = "";
    let error = "";
    let success = false;
    let submitting = false;

    let communityName = "";
    let description = "";
    let selectedCategories: string[] = [];
    let rules = "";
    let iconFile: File | null = null;
    let bannerFile: File | null = null;
    let iconPreview = "";
    let bannerPreview = "";

    let availableCategories: string[] = [];
    let showCategoryDropdown = false;

    let nameError = "";
    let descriptionError = "";
    let categoryError = "";
    let rulesError = "";
    let iconError = "";
    let bannerError = "";

    async function loadMyProfile() {
      try {
        const meRes = await api.get<MeResponse>("/user/get-me");
        myUserId = meRes.data.id;
      } catch (err) {
        console.error("Failed to load profile:", err);
        window.location.href = "/login";
      }
    }

    async function fetchCategories() {
      try {
        const res = await api.get<{ categories: string[] }>("/community-categories");
        availableCategories = res.data;
        console.log("catego: ",availableCategories);
      } catch (err) {
        console.error("Failed to load categories:", err);
        availableCategories = ["Technology", "Gaming", "Art", "Music", "Sports", "Education", "Business", "Health", "Science"];
      }
    }
    function handleIconUpload(event: Event) {
      const target = event.target as HTMLInputElement;
      const file = target.files?.[0];
      if (file) {
        if (file.size > 5 * 1024 * 1024) { // 5MB limit
          iconError = "Icon must be smaller than 5MB";
          return;
        }
        if (!file.type.startsWith("image/")) {
          iconError = "Please upload an image file";
          return;
        }
        iconFile = file;
        iconError = "";
        const reader = new FileReader();
        reader.onload = (e) => {
          iconPreview = e.target?.result as string;
        };
        reader.readAsDataURL(file);
      }
    }

    function handleBannerUpload(event: Event) {
      const target = event.target as HTMLInputElement;
      const file = target.files?.[0];

      if (file) {
        if (file.size > 10 * 1024 * 1024) {
          bannerError = "Banner must be smaller than 10MB";
          return;
        }

        if (!file.type.startsWith("image/")) {
          bannerError = "Please upload an image file";
          return;
        }

        bannerFile = file;
        bannerError = "";
        const reader = new FileReader();
        reader.onload = (e) => {
          bannerPreview = e.target?.result as string;
        };
        reader.readAsDataURL(file);
      }
    }

    function toggleCategory(category: string) {
      if (selectedCategories.includes(category)) {
        selectedCategories = selectedCategories.filter(c => c !== category);
      } else {
        selectedCategories = [...selectedCategories, category];
      }
      categoryError = "";
    }

    function removeCategory(category: string) {
      selectedCategories = selectedCategories.filter(c => c !== category);
    }

    function validateForm(): boolean {
      let isValid = true;
      nameError = descriptionError = categoryError = rulesError = iconError = bannerError = "";
      if (!communityName.trim()) {
        nameError = "Community name is required";
        isValid = false;
      } else if (communityName.length < 3) {
        nameError = "Community name must be at least 3 characters";
        isValid = false;
      } else if (communityName.length > 50) {
        nameError = "Community name must be less than 50 characters";
        isValid = false;
      }
      if (!description.trim()) {
        descriptionError = "Description is required";
        isValid = false;
      } else if (description.length < 10) {
        descriptionError = "Description must be at least 10 characters";
        isValid = false;
      } else if (description.length > 500) {
        descriptionError = "Description must be less than 500 characters";
        isValid = false;
      }
      if (selectedCategories.length === 0) {
        categoryError = "At least one category is required";
        isValid = false;
      } else if (selectedCategories.length > 5) {
        categoryError = "Maximum 5 categories allowed";
        isValid = false;
      }
      if (!rules.trim()) {
        rulesError = "Community rules are required";
        isValid = false;
      } else if (rules.length < 20) {
        rulesError = "Rules must be at least 20 characters";
        isValid = false;
      } else if (rules.length > 2000) {
        rulesError = "Rules must be less than 2000 characters";
        isValid = false;
      }
      if (!iconFile) {
        iconError = "Community icon is required";
        isValid = false;
      }
      if (!bannerFile) {
        bannerError = "Community banner is required";
        isValid = false;
      }
      return isValid;
    }
    async function handleSubmit(event: Event) {
      event.preventDefault();
      if (!validateForm()) {
        return;
      }
      submitting = true;
      error = "";
      try {
        let iconId = "";
        if (iconFile) {
          const iconFormData = new FormData();
          iconFormData.append("file", iconFile);
          const iconResponse = await api.post<{ id: string }>("/media/upload", iconFormData);
          iconId = iconResponse.data.id;
        }
        let bannerId = "";
        if (bannerFile) {
          const bannerFormData = new FormData();
          bannerFormData.append("file", bannerFile);
          const bannerResponse = await api.post<{ id: string }>("/media/upload", bannerFormData);
          bannerId = bannerResponse.data.id;
        }
        const categoryIds = selectedCategories.map(cat => cat.category_id);
        const communityData = {
          community_name: communityName.trim(),
          community_description: description.trim(),
          community_rules: rules.trim(),
          status: "pending",
          community_logo: iconId,
          community_banner: bannerId,
          category_ids: categoryIds,
          owner_id: myUserId
        };
        await api.post("/communities", communityData);
        success = true;
        communityName = "";
        description = "";
        selectedCategories = [];
        rules = "";
        iconFile = null;
        bannerFile = null;
        iconPreview = "";
        bannerPreview = "";
      } catch (err) {
        console.error("Failed to create community:", err);
        error = err.response?.data?.message || "Failed to create community request. Please try again.";
      }
      submitting = false;
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
      await fetchCategories();
    });
  </script>
  <div class="create-community layout">
    <aside class="sidebar">
      {#if myUserId && windowWidth >= 1250}
        <LeftSidebar currentUserId={myUserId} activePage="communities" />
      {/if}
    </aside>
    {#if myUserId && windowWidth < 1250}
      <BurgerLeftSideBar currentUserId={myUserId} activePage="communities" />
    {/if}
    <main class="main">
      <div class="header">
        <h1>Create Community</h1>
        <p class="subtitle">Start your own community and bring people together</p>
      </div>
      {#if success}
        <div class="success-message" in:fade>
          <AlertCircle size={24} />
          <div>
            <h3>Community Request Submitted!</h3>
            <p>Your community creation request has been submitted and is pending admin approval. You'll be notified once it's reviewed.</p>
          </div>
        </div>
      {/if}
      {#if error}
        <div class="error-message" in:fade>
          <AlertCircle size={20} />
          <span>{error}</span>
        </div>
      {/if}
      <form class="create-form" on:submit={handleSubmit} in:fade>
        <div class="form-group">
          <label for="name">Community Name *</label>
          <input
            id="name"
            type="text"
            bind:value={communityName}
            placeholder="Enter community name"
            class:error={nameError}
            maxlength="50"
          />
          {#if nameError}
            <span class="error-text">{nameError}</span>
          {/if}
          <div class="char-count">{communityName.length}/50</div>
        </div>
        <div class="form-group">
          <label for="description">Description *</label>
          <textarea
            id="description"
            bind:value={description}
            placeholder="Describe what your community is about..."
            class:error={descriptionError}
            maxlength="500"
            rows="4"
          ></textarea>
          {#if descriptionError}
            <span class="error-text">{descriptionError}</span>
          {/if}
          <div class="char-count">{description.length}/500</div>
        </div>
        <div class="form-group">
          <label>Categories * (Select up to 5)</label>
          <div class="selected-categories">
            {#each selectedCategories as category (category.category_id)}
              <span class="selected-category">
                {category.category_name}
                <button type="button" on:click={() => removeCategory(category)}>
                  <X size={16} />
                </button>
              </span>
            {/each}
          </div>
          <div class="category-selector">
            <button
              type="button"
              class="category-toggle"
              on:click={() => showCategoryDropdown = !showCategoryDropdown}
            >
              <Plus size={16} />
              Add Category
            </button>
            {#if showCategoryDropdown}
              <div class="category-dropdown" in:fade>
                {#each availableCategories as category (category.category_id)}
                  <button
                    type="button"
                    class="category-option"
                    class:selected={selectedCategories.includes(category.category_name)}
                    disabled={selectedCategories.includes(category) || selectedCategories.length >= 5}
                    on:click={() => toggleCategory(category)}
                  >
                    {category.category_name}
                  </button>
                {/each}
              </div>
            {/if}
          </div>
          {#if categoryError}
            <span class="error-text">{categoryError}</span>
          {/if}
        </div>
        <div class="form-group">
          <label>Community Icon *</label>
          <div class="upload-section">
            <input
              type="file"
              accept="image/*"
              on:change={handleIconUpload}
              id="icon-upload"
              class="file-input"
            />
            <label for="icon-upload" class="upload-button">
              <Upload size={20} />
              Choose Icon
            </label>
            {#if iconPreview}
              <div class="image-preview">
                <img src={iconPreview} alt="Icon preview" class="icon-preview" />
              </div>
            {/if}
          </div>
          {#if iconError}
            <span class="error-text">{iconError}</span>
          {/if}
          <div class="upload-hint">Recommended: Square image, max 5MB</div>
        </div>
        <div class="form-group">
          <label>Community Banner *</label>
          <div class="upload-section">
            <input
              type="file"
              accept="image/*"
              on:change={handleBannerUpload}
              id="banner-upload"
              class="file-input"
            />
            <label for="banner-upload" class="upload-button">
              <Upload size={20} />
              Choose Banner
            </label>
            {#if bannerPreview}
              <div class="image-preview">
                <img src={bannerPreview} alt="Banner preview" class="banner-preview" />
              </div>
            {/if}
          </div>
          {#if bannerError}
            <span class="error-text">{bannerError}</span>
          {/if}
          <div class="upload-hint">Recommended: 16:9 aspect ratio, max 10MB</div>
        </div>
        <div class="form-group">
          <label for="rules">Community Rules *</label>
          <textarea
            id="rules"
            bind:value={rules}
            placeholder="Define the rules and guidelines for your community..."
            class:error={rulesError}
            maxlength="2000"
            rows="6"
          ></textarea>
          {#if rulesError}
            <span class="error-text">{rulesError}</span>
          {/if}
          <div class="char-count">{rules.length}/2000</div>
        </div>
        <div class="form-actions">
          <button
            type="button"
            class="cancel-btn"
            on:click={() => window.location.href = "/communities"}
          >
            Cancel
          </button>
          <button
            type="submit"
            class="submit-btn"
            disabled={submitting}
          >
            {submitting ? "Submitting..." : "Submit for Review"}
          </button>
        </div>
      </form>
      <div class="info-box">
        <AlertCircle size={20} />
        <div>
          <h4>Review Process</h4>
          <p>Your community request will be reviewed by our administrators. This process typically takes 1-3 business days. You'll receive a notification once your community is approved or if any changes are needed.</p>
        </div>
      </div>
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
    .create-community.layout {
      display: grid;
      grid-template-columns: auto 1fr auto;
      min-height: 100vh;
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
      height: 90vh;
      /* margin-bottom: 5vh; */
    }
    .header {
      margin-bottom: 2rem;
    }
    .header h1 {
      font-size: 2rem;
      font-weight: 700;
      color: #e0e8ff;
      margin-bottom: 0.5rem;
    }
    .subtitle {
      color: #b2c2d8;
      font-size: 1.1rem;
    }
    .success-message {
      display: flex;
      align-items: flex-start;
      gap: 1rem;
      background: #1a4d3a;
      border: 1px solid #29e57b;
      border-radius: 0.8rem;
      padding: 1.2rem;
      margin-bottom: 1.5rem;
      color: #29e57b;
    }
    .success-message h3 {
      margin: 0 0 0.5rem 0;
      font-weight: 600;
    }
    .success-message p {
      margin: 0;
      color: #b8f5d1;
    }
    .error-message {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      background: #4d1a1a;
      border: 1px solid #ff6b6b;
      border-radius: 0.8rem;
      padding: 1rem;
      margin-bottom: 1.5rem;
      color: #ff6b6b;
    }
    .create-form {
      display: flex;
      flex-direction: column;
      gap: 1.5rem;
      margin-bottom: 2rem;
    }
    .form-group {
      display: flex;
      flex-direction: column;
      gap: 0.5rem;
    }
    .form-group label {
      font-weight: 600;
      color: #e0e8ff;
      font-size: 1rem;
    }
    .form-group input,
    .form-group textarea {
      padding: 0.75rem;
      border-radius: 0.8rem;
      border: 1px solid #285a97;
      background: #161a1d;
      color: #f7f7f7;
      font-size: 1rem;
      transition: border-color 0.2s;
    }
    .form-group input:focus,
    .form-group textarea:focus {
      outline: none;
      border-color: #1d9bf0;
    }
    .form-group input.error,
    .form-group textarea.error {
      border-color: #ff6b6b;
    }
    .error-text {
      color: #ff6b6b;
      font-size: 0.9rem;
    }
    .char-count {
      color: #abb6c4;
      font-size: 0.85rem;
      text-align: right;
    }
    .selected-categories {
      display: flex;
      flex-wrap: wrap;
      gap: 0.5rem;
      margin-bottom: 0.5rem;
    }
    .selected-category {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      background: #232a39;
      color: #5bc0f8;
      padding: 0.4rem 0.8rem;
      border-radius: 0.5rem;
      font-size: 0.9rem;
      font-weight: 500;
    }
    .selected-category button {
      background: none;
      border: none;
      color: #5bc0f8;
      cursor: pointer;
      display: flex;
      align-items: center;
      padding: 0;
    }
    .category-selector {
      position: relative;
    }
    .category-toggle {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      background: #1d9bf0;
      color: #fff;
      border: none;
      border-radius: 0.6rem;
      padding: 0.5rem 1rem;
      font-weight: 500;
      cursor: pointer;
      font-size: 0.9rem;
    }
    .category-dropdown {
      position: absolute;
      top: 100%;
      left: 0;
      right: 0;
      background: #161a1d;
      border: 1px solid #285a97;
      border-radius: 0.8rem;
      padding: 0.5rem;
      margin-top: 0.5rem;
      display: flex;
      flex-wrap: wrap;
      gap: 0.5rem;
      z-index: 10;
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
    }
    .category-option {
      background: #232a39;
      color: #b2c2d8;
      border: none;
      border-radius: 0.5rem;
      padding: 0.4rem 0.8rem;
      cursor: pointer;
      font-size: 0.9rem;
      transition: all 0.2s;
    }
    .category-option:hover:not(:disabled) {
      background: #1d9bf0;
      color: #fff;
    }
    .category-option.selected {
      background: #5bc0f8;
      color: #000;
    }
    .category-option:disabled {
      opacity: 0.5;
      cursor: not-allowed;
    }
    .upload-section {
      display: flex;
      align-items: center;
      gap: 1rem;
    }
    .file-input {
      display: none;
    }
    .upload-button {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      background: #1d9bf0;
      color: #fff;
      border: none;
      border-radius: 0.6rem;
      padding: 0.6rem 1.2rem;
      font-weight: 500;
      cursor: pointer;
      font-size: 0.9rem;
    }
    .image-preview {
      display: flex;
      align-items: center;
    }
    .icon-preview {
      width: 60px;
      height: 60px;
      border-radius: 0.8rem;
      object-fit: cover;
      border: 2px solid #1d9bf0;
    }
    .banner-preview {
      width: 120px;
      height: 68px;
      border-radius: 0.6rem;
      object-fit: cover;
      border: 2px solid #1d9bf0;
    }
    .upload-hint {
      color: #abb6c4;
      font-size: 0.85rem;
    }
    .form-actions {
      display: flex;
      gap: 1rem;
      justify-content: flex-end;
      margin-top: 1rem;
    }
    .cancel-btn {
      background: #232a39;
      color: #b2c2d8;
      border: none;
      border-radius: 0.6rem;
      padding: 0.75rem 1.5rem;
      font-weight: 500;
      cursor: pointer;
      font-size: 1rem;
    }
    .submit-btn {
      background: #1d9bf0;
      color: #fff;
      border: none;
      border-radius: 0.6rem;
      padding: 0.75rem 1.5rem;
      font-weight: 600;
      cursor: pointer;
      font-size: 1rem;
      transition: background 0.2s;
    }
    .submit-btn:hover:not(:disabled) {
      background: #1a8cd8;
    }
    .submit-btn:disabled {
      background: #285a97;
      cursor: not-allowed;
    }
    .info-box {
      display: flex;
      align-items: flex-start;
      gap: 1rem;
      background: #161a1e;
      border: 1px solid #285a97;
      border-radius: 0.8rem;
      padding: 1.2rem;
      color: #b2c2d8;
    }
    .info-box h4 {
      margin: 0 0 0.5rem 0;
      color: #e0e8ff;
      font-weight: 600;
    }
    .info-box p {
      margin: 0;
      line-height: 1.5;
    }
    @media (max-width: 900px) {
      .main {
        width: 98vw;
        max-width: 98vw;
        padding: 1rem 0.5rem;
      }
      .form-actions {
        flex-direction: column;
      }
      .upload-section {
        flex-direction: column;
        align-items: flex-start;
      }
    }
    .sidebar{
      z-index: 400000;
    }
  </style>