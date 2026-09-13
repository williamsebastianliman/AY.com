<script lang="ts">
    import { onMount } from "svelte";
    import { fade } from "svelte/transition";
    import LeftSidebar from "../lib/components/LeftSideBar.svelte";
    import RightSidebar from "../lib/components/RightSideBar.svelte";
    import BurgerLeftSideBar from "../lib/components/BurgerLeftSideBar.svelte";
    import BurgerRightSideBar from "../lib/components/BurgerRightSideBar.svelte";
    import {
      Users,
      Mail,
      Building2,
      Crown,
      Flag,
      Folder,
      FolderTree,
      CheckCircle,
      XCircle,
      Edit,
      Trash2,
      Plus,
      Search,
      Send
    } from "lucide-svelte";
    import api from "../lib/api";
    type AdminTab = "users" | "newsletter" | "communities" | "premium" | "reports" | "threadCategories" | "communityCategories";
    let activeTab: AdminTab = "users";
    interface User {
      id: string;
      username: string;
      name: string;
      profile_picture_id: string;
    }
    let users: User[] = [];
    let usersLoading = false;
    let userSearchQuery = "";
    let newsletterContent = "";
    let newsletterSubject = "";
    let sendingNewsletter = false;
    let subscribedUsersCount = 0;
    interface CommunityRequest {
      community_id: string;
      community_name: string;
      community_description: string;
      community_rules: string;
      created_at: string;
    }
    let communityRequests: CommunityRequest[] = [];
    let communityRequestsLoading = false;
    interface PremiumRequest {
      card_number: string;
      face_image: string;
      reason: string;
      user_id: string;
      image_url: string;
    }
    let premiumRequests: PremiumRequest[] = [];
    let premiumRequestsLoading = false;
    interface ReportRequest {
      report_id: string;
      reported_user: string;
      reason: string;
    }
    let reportRequests: ReportRequest[] = [];
    let reportRequestsLoading = false;
    interface ThreadCategory {
      category_id: string;
      category_name: string;
      created_at: string;
    }
    let threadCategories: ThreadCategory[] = [];
    let threadCategoriesLoading = false;
    let showThreadCategoryModal = false;
    let editingThreadCategory: ThreadCategory | null = null;
    let threadCategoryForm = {
      name: "",
      description: ""
    };
    interface CommunityCategory {
      id: string;
      name: string;
      description: string;
      created_at: string;
    }
    let communityCategories: CommunityCategory[] = [];
    let communityCategoriesLoading = false;
    let showCommunityCategoryModal = false;
    let editingCommunityCategory: CommunityCategory | null = null;
    let communityCategoryForm = {
      name: "",
      description: ""
    };
    let myUserId = "";
    onMount(async () => {
      await loadMyProfile();
      await loadUsers();
      await loadSubscribedUsersCount();
      await loadCommunityRequests();
      await loadPremiumRequests();
      await loadReportRequests();
      await loadThreadCategories();
      await loadCommunityCategories();
    });
    async function loadMyProfile() {
      try {
        const meRes = await api.get("/user/get-me");
        myUserId = meRes.data.id;
      } catch (err) {
        console.error("Failed to load profile:", err);
      }
    }
    async function loadUsers() {
      usersLoading = true;
      try {
        const response = await api.get<User[]>("/users");
        users = response.data;
      } catch (err) {
        console.error("Failed to load users:", err);
      } finally {
        usersLoading = false;
      }
    }
    async function toggleUserBan(userId: string, currentBanStatus: boolean) {
      try {
        if (currentBanStatus) {
          await api.post(`/admin/users/${userId}/unban`);
        } else {
          await api.post(`/admin/users/${userId}/ban`);
        }
        await loadUsers();
      } catch (err) {
        console.error("Failed to toggle user ban:", err);
      }
    }
    async function loadSubscribedUsersCount() {
      try {
        const response = await api.get<{ count: number }>("/admin/newsletter/subscribers-count");
        subscribedUsersCount = response.data.count;
      } catch (err) {
        console.error("Failed to load subscribers count:", err);
      }
    }
    async function sendNewsletter() {
      if (!newsletterSubject.trim() || !newsletterContent.trim()) return;
      sendingNewsletter = true;
      try {
        await api.post("/admin/newsletter/send", {
          subject: newsletterSubject,
          content: newsletterContent
        });
        newsletterSubject = "";
        newsletterContent = "";
        alert("Newsletter sent successfully!");
      } catch (err) {
        console.error("Failed to send newsletter:", err);
        alert("Failed to send newsletter. Please try again.");
      } finally {
        sendingNewsletter = false;
      }
    }
    async function loadCommunityRequests() {
      communityRequestsLoading = true;
      try {
        const response = await api.get<CommunityRequest[]>("/communities/pending");
        communityRequests = response.data;
      } catch (err) {
        console.error("Failed to load community requests:", err);
      } finally {
        communityRequestsLoading = false;
      }
    }
    async function handleCommunityRequest(requestId: string, action: "accept" | "reject") {
      try {
        await api.post(`/community/${action}/${requestId}`);
        await loadCommunityRequests();
      } catch (err) {
        console.error(`Failed to ${action} community request:`, err);
      }
    }
    async function loadPremiumRequests() {
      premiumRequestsLoading = true;
      try {
        const response = await api.get<PremiumRequest[]>("user/get-premium");
        const withImageUrls = await Promise.all(response.data.requests.map(async (req)=>{
          let image_url = "";
          try{
              const mediaRes = await api.post("/media/get-media", {id: req.face_image});
              image_url = mediaRes.data.public_url;
              return {...req, image_url};
          }
          catch(err)
          {
            console.error(err as string);
          }
        }));
        premiumRequests = withImageUrls;
        console.log(premiumRequests);
      } catch (err) {
        console.error("Failed to load premium requests:", err);
      } finally {
        premiumRequestsLoading = false;
      }
    }
    async function handlePremiumRequest(requestId: string, action: "accept" | "reject") {
      console.log("action: ",action);
      console.log("reqid: ", requestId);

      try {
        await api.post(`/user/premium-request/${action}`, {user_id: requestId});
        premiumRequests = premiumRequests.filter(r => r.user_id !== requestId);
      } catch (err) {
        console.error(`Failed to ${action} premium request:`, err);
      }
    }
    async function loadReportRequests() {
      reportRequestsLoading = true;
      try {
        const response = await api.get<ReportRequest[]>("/user-reports");
        reportRequests = response.data;
      } catch (err) {
        console.error("Failed to load report requests:", err);
      } finally {
        reportRequestsLoading = false;
      }
    }
    async function handleReportRequest(requestId: string, action: "approve" | "reject") {
      try {
        await api.post(`/admin/report-requests/${requestId}/${action}`);
        await loadReportRequests();
      } catch (err) {
        console.error(`Failed to ${action} report request:`, err);
      }
    }
    async function loadThreadCategories() {
      threadCategoriesLoading = true;
      try {
        const response = await api.get<ThreadCategory[]>("/thread-categories");
        threadCategories = response.data;
      } catch (err) {
        console.error("Failed to load thread categories:", err);
      } finally {
        threadCategoriesLoading = false;
      }
    }
    function openThreadCategoryModal(category?: ThreadCategory) {
      editingThreadCategory = category || null;
      threadCategoryForm = {
        name: category?.name || "",
        description: category?.description || ""
      };
      showThreadCategoryModal = true;
    }
    function closeThreadCategoryModal() {
      showThreadCategoryModal = false;
      editingThreadCategory = null;
      threadCategoryForm = { name: "", description: "" };
    }
    async function saveThreadCategory() {
      if (!threadCategoryForm.name.trim()) return;
      try {
        if (editingThreadCategory) {
          await api.put(`/thread-categories/${editingThreadCategory.category_id}`, {category_name: threadCategoryForm.name.trim()});
        } else {
          await api.post("/thread-categories", {category_name: threadCategoryForm.name.trim()});
        }
        await loadThreadCategories();
        closeThreadCategoryModal();
      } catch (err) {
        console.error("Failed to save thread category:", err);
      }
    }
    async function deleteThreadCategory(categoryId: string) {
      try {
        await api.delete(`/thread-categories/${categoryId}`);
        await loadThreadCategories();
      } catch (err) {
        console.error("Failed to delete thread category:", err);
      }
    }
    async function loadCommunityCategories() {
      communityCategoriesLoading = true;
      try {
        const response = await api.get<CommunityCategory[]>("/community-categories");
        communityCategories = response.data;
      } catch (err) {
        console.error("Failed to load community categories:", err);
      } finally {
        communityCategoriesLoading = false;
      }
    }
    function openCommunityCategoryModal(category?: CommunityCategory) {
      editingCommunityCategory = category || null;
      communityCategoryForm = {
        name: category?.name || "",
        description: category?.description || ""
      };
      showCommunityCategoryModal = true;
    }
    function closeCommunityCategoryModal() {
      showCommunityCategoryModal = false;
      editingCommunityCategory = null;
      communityCategoryForm = { name: "", description: "" };
    }
    async function saveCommunityCategory() {
      if (!communityCategoryForm.name.trim()) return;
      try {
        console.log("comz: ", editingCommunityCategory);
        if (editingCommunityCategory) {
          await api.put(`/community-categories/${editingCommunityCategory.category_id}`, {category_name: communityCategoryForm.name.trim()});
        } else {
          await api.post("/community-categories", {category_name: communityCategoryForm.name.trim()});
        }
        await loadCommunityCategories();
        closeCommunityCategoryModal();
      } catch (err) {
        console.error("Failed to save community category:", err);
      }
    }
    async function deleteCommunityCategory(categoryId: string) {
      try {
        await api.delete(`community-categories/${categoryId}`);
        await loadCommunityCategories();
      } catch (err) {
        console.error("Failed to delete community category:", err);
      }
    }
    function formatDate(dateString: string) {
      return new Date(dateString).toLocaleDateString();
    }
    $: filteredUsers = users.filter(user =>
      user.name.toLowerCase().includes(userSearchQuery.toLowerCase()) ||
      user.username.toLowerCase().includes(userSearchQuery.toLowerCase()) ||
      user.email.toLowerCase().includes(userSearchQuery.toLowerCase())
    );
    let windowWidth = window.innerWidth;
    function handleResize() {
      windowWidth = window.innerWidth;
    }
    function navigateToProfile(id: string){
      window.location.href = `/profile/${id}`;
    }
    onMount(() => {
      window.addEventListener("resize", handleResize);
      return () => window.removeEventListener("resize", handleResize);
    });
  </script>
  <div class="admin layout">
    <aside class="sidebar">
      <aside class="sidebar">
        {#if myUserId && windowWidth >= 1250}
          <LeftSidebar currentUserId={myUserId} activePage="" />
        {/if}
      </aside>
      {#if myUserId && windowWidth < 1250}
        <BurgerLeftSideBar currentUserId={myUserId} activePage="" />
      {/if}
    </aside>
    <main class="main">
      <div class="admin-header" in:fade>
        <h1>Admin Dashboard</h1>
        <p>Manage users, communities, and platform content</p>
      </div>
      <div class="tabs">
        <button
          type="button"
          class:active={activeTab === "users"}
          on:click={() => activeTab = "users"}
        >
          <Users size="18" />
          Users
        </button>
        <button
          type="button"
          class:active={activeTab === "newsletter"}
          on:click={() => activeTab = "newsletter"}
        >
          <Mail size="18" />
          Newsletter
        </button>
        <button
          type="button"
          class:active={activeTab === "communities"}
          on:click={() => activeTab = "communities"}
        >
          <Building2 size="18" />
          Communities
        </button>
        <button
          type="button"
          class:active={activeTab === "premium"}
          on:click={() => activeTab = "premium"}
        >
          <Crown size="18" />
          Premium
        </button>
        <button
          type="button"
          class:active={activeTab === "reports"}
          on:click={() => activeTab = "reports"}
        >
          <Flag size="18" />
          Reports
        </button>
        <button
          type="button"
          class:active={activeTab === "threadCategories"}
          on:click={() => activeTab = "threadCategories"}
        >
          <Folder size="18" />
          Thread Categories
        </button>
        <button
          type="button"
          class:active={activeTab === "communityCategories"}
          on:click={() => activeTab = "communityCategories"}
        >
          <FolderTree size="18" />
          Community Categories
        </button>
      </div>
      {#if activeTab === "users"}
        <div class="tab-content" in:fade>
          <div class="section-header">
            <h2>User Management</h2>
            <div class="search-bar">
              <Search size="18" />
              <input
                type="text"
                placeholder="Search users..."
                bind:value={userSearchQuery}
              />
            </div>
          </div>
          {#if usersLoading}
            <div class="loading">Loading users...</div>
          {:else}
            <div class="users-grid">
              {#each filteredUsers as user (user.id)}
                <div class="user-card"  on:click={()=>{ navigateToProfile(user.id)}} in:fade>
                  <div class="user-info">
                    <div class="user-details">
                      <h3>{user.name}</h3>
                      <p>@{user.username}</p>
                    </div>
                  </div>
                </div>
              {/each}
            </div>
          {/if}
        </div>
      {/if}
      {#if activeTab === "newsletter"}
        <div class="tab-content" in:fade>
          <div class="section-header">
            <h2>Newsletter</h2>
            <p class="subscriber-count">{subscribedUsersCount} subscribers</p>
          </div>
          <div class="newsletter-form">
            <div class="form-group">
              <label for="newsletter-subject">Subject</label>
              <input
                id="newsletter-subject"
                type="text"
                placeholder="Newsletter subject..."
                bind:value={newsletterSubject}
              />
            </div>
            <div class="form-group">
              <label for="newsletter-content">Content</label>
              <textarea
                id="newsletter-content"
                placeholder="Write your newsletter content..."
                rows="10"
                bind:value={newsletterContent}
              ></textarea>
            </div>
            <button
              type="button"
              class="send-btn"
              disabled={!newsletterSubject.trim() || !newsletterContent.trim() || sendingNewsletter}
              on:click={sendNewsletter}
            >
              <Send size="18" />
              {sendingNewsletter ? "Sending..." : "Send Newsletter"}
            </button>
          </div>
        </div>
      {/if}
      {#if activeTab === "communities"}
        <div class="tab-content" in:fade>
          <div class="section-header">
            <h2>Community Requests</h2>
          </div>
          {#if communityRequestsLoading}
            <div class="loading">Loading community requests...</div>
          {:else}
            <div class="requests-list">
              {#each communityRequests as request (request.community_id)}
                <div class="request-card" in:fade>
                  <div class="request-info">
                    <h3>{request.community_name}</h3>
                    <p class="description">Description: {request.community_description}</p>
                    <p class="description">Rules: </p>
                    <p class="description">{request.community_rules}</p>
                  </div>
                  {#if request.status === "pending"}
                    <div class="request-actions">
                      <button
                        type="button"
                        class="approve-btn"
                        on:click={() => handleCommunityRequest(request.community_id, "accept")}
                      >
                        <CheckCircle size="16" />
                        Approve
                      </button>
                      <button
                        type="button"
                        class="reject-btn"
                        on:click={() => handleCommunityRequest(request.community_id, "reject")}
                      >
                        <XCircle size="16" />
                        Reject
                      </button>
                    </div>
                  {/if}
                </div>
              {/each}
            </div>
          {/if}
        </div>
      {/if}
      {#if activeTab === "premium"}
        <div class="tab-content" in:fade>
          <div class="section-header">
            <h2>Premium Requests</h2>
          </div>
          {#if premiumRequestsLoading}
            <div class="loading">Loading premium requests...</div>
          {:else}
            <div class="requests-list">
              {#each premiumRequests as request (request.user_id)}
                <div class="request-card" in:fade>
                  <div class="request-header">
                    <div class="request-info">
                      <h3>{request.card_number}</h3>
                      <p class="email">{request.reason}</p>
                    </div>
                    <img class="face-verify-img" src={request.image_url}>
                  </div>
                  <div class="request-actions">
                    <button
                      type="button"
                      class="approve-btn"
                      on:click={() => handlePremiumRequest(request.user_id, "accept")}
                    >
                      <CheckCircle size="16" />
                      Approve
                    </button>
                    <button
                      type="button"
                      class="reject-btn"
                      on:click={() => handlePremiumRequest(request.user_id, "reject")}
                    >
                      <XCircle size="16" />
                      Reject
                    </button>
                  </div>
                </div>
              {/each}
            </div>
          {/if}
        </div>
      {/if}
      {#if activeTab === "reports"}
        <div class="tab-content" in:fade>
          <div class="section-header">
            <h2>Report Requests</h2>
          </div>
          {#if reportRequestsLoading}
            <div class="loading">Loading report requests...</div>
          {:else}
            <div class="requests-list">
              {#each reportRequests as request (request.id)}
                <div class="request-card report-card" in:fade>
                  <div class="request-info">
                    <div class="report-header">
                      <h2>Report against <a href={"/profile/" + request.reported_user} class="user-link">
                        {request.reported_user}
                      </a></h2>
                    </div>
                    <div class="report-header">
                      <h2>Reason: {request.reason}</h2>
                    </div>
                  </div>
                </div>
              {/each}
            </div>
          {/if}
        </div>
      {/if}
      {#if activeTab === "threadCategories"}
        <div class="tab-content" in:fade>
          <div class="section-header">
            <h2>Thread Categories</h2>
            <button
              type="button"
              class="add-btn"
              on:click={() => openThreadCategoryModal()}
            >
              <Plus size="18" />
              Add Category
            </button>
          </div>
          {#if threadCategoriesLoading}
            <div class="loading">Loading thread categories...</div>
          {:else}
            <div class="categories-grid">
              {#each threadCategories as category (category.category_id)}
                <div class="category-card" in:fade>
                  <div class="category-info">
                    <h3>{category.category_name}</h3>
                    <!-- <span class="date">Created: {formatDate(category.created_at)}</span> -->
                  </div>
                  <div class="category-actions">
                    <button
                      type="button"
                      class="edit-btn"
                      on:click={() => openThreadCategoryModal(category)}
                    >
                      <Edit size="16" />
                    </button>
                    <button
                      type="button"
                      class="delete-btn"
                      on:click={() => deleteThreadCategory(category.category_id)}
                    >
                      <Trash2 size="16" />
                    </button>
                  </div>
                </div>
              {/each}
            </div>
          {/if}
        </div>
      {/if}
      {#if activeTab === "communityCategories"}
        <div class="tab-content" in:fade>
          <div class="section-header">
            <h2>Community Categories</h2>
            <button
              type="button"
              class="add-btn"
              on:click={() => openCommunityCategoryModal()}
            >
              <Plus size="18" />
              Add Category
            </button>
          </div>
          {#if communityCategoriesLoading}
            <div class="loading">Loading community categories...</div>
          {:else}
            <div class="categories-grid">
              {#each communityCategories as category (category.category_id)}
                <div class="category-card" in:fade>
                  <div class="category-info">
                    <h3>{category.category_name}</h3>
                    <!-- <span class="date">Created: {formatDate(category.created_at)}</span> -->
                  </div>
                  <div class="category-actions">
                    <button
                      type="button"
                      class="edit-btn"
                      on:click={() => openCommunityCategoryModal(category)}
                    >
                      <Edit size="16" />
                    </button>
                    <button
                      type="button"
                      class="delete-btn"
                      on:click={() => deleteCommunityCategory(category.category_id)}
                    >
                      <Trash2 size="16" />
                    </button>
                  </div>
                </div>
              {/each}
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
  {#if showThreadCategoryModal}
    <div class="modal-overlay" on:click={closeThreadCategoryModal}>
      <div class="modal" on:click|stopPropagation>
        <div class="modal-header">
          <h3>{editingThreadCategory ? "Edit" : "Add"} Thread Category</h3>
          <button type="button" class="modal-close" on:click={closeThreadCategoryModal}>×</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label for="thread-category-name">Name</label>
            <input
              id="thread-category-name"
              type="text"
              placeholder="Category name"
              bind:value={threadCategoryForm.name}
            />
          </div>
        </div>
        <div class="modal-footer">
          <button type="button" class="cancel-btn" on:click={closeThreadCategoryModal}>
            Cancel
          </button>
          <button type="button" class="save-btn" on:click={saveThreadCategory}>
            {editingThreadCategory ? "Update" : "Create"}
          </button>
        </div>
      </div>
    </div>
  {/if}
  {#if showCommunityCategoryModal}
    <div class="modal-overlay" on:click={closeCommunityCategoryModal}>
      <div class="modal" on:click|stopPropagation>
        <div class="modal-header">
          <h3>{editingCommunityCategory ? "Edit" : "Add"} Community Category</h3>
          <button type="button" class="modal-close" on:click={closeCommunityCategoryModal}>×</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label for="community-category-name">Name</label>
            <input
              id="community-category-name"
              type="text"
              placeholder="Category name"
              bind:value={communityCategoryForm.name}
            />
          </div>
        </div>
        <div class="modal-footer">
          <button type="button" class="cancel-btn" on:click={closeCommunityCategoryModal}>
            Cancel
          </button>
          <button type="button" class="save-btn" on:click={saveCommunityCategory}>
            {editingCommunityCategory ? "Update" : "Create"}
          </button>
        </div>
      </div>
    </div>
  {/if}
  <style>
    .admin::-webkit-scrollbar { display: none; }
    .admin, .admin * {
      overflow-y: hidden;
      scrollbar-width: none;
      -ms-overflow-style: none;
    }
    .layout {
      display: grid;
      grid-template-columns: auto 1fr auto;
      height: 100vh;
    }
    .main {
      overflow-y: auto;
      padding: 1rem;
      width: 50vw;
      max-width: 60vw;
      margin: 0 auto;
    }
    .admin-header {
      text-align: center;
      margin-bottom: 2rem;
      padding: 2rem 0;
    }
    .admin-header h1 {
      font-size: 2.5rem;
      font-weight: 700;
      margin-bottom: 0.5rem;
      background: linear-gradient(135deg, #1d9bf0, #00d4ff);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      background-clip: text;
    }
    .admin-header p {
      color: rgba(255,255,255,0.7);
      font-size: 1.1rem;
    }
    .tabs {
      display: flex;
      flex-wrap: wrap;
      gap: 0.5rem;
      margin-bottom: 2rem;
      padding: 1rem;
      background: rgba(30,30,30,0.8);
      border-radius: 1rem;
      border: 1px solid rgba(255,255,255,0.08);
    }
    .tabs button {
      background: none;
      border: none;
      padding: 0.75rem 1rem;
      font-weight: 600;
      font-size: 0.9rem;
      color: rgba(255,255,255,0.6);
      border-radius: 0.5rem;
      cursor: pointer;
      transition: all 0.2s;
      display: flex;
      align-items: center;
      gap: 0.5rem;
    }
    .tabs button:hover {
      color: #fff;
      background: rgba(255,255,255,0.05);
    }
    .tabs button.active {
      color: #fff;
      background: #1d9bf0;
    }
    .tab-content {
      background: rgba(30,30,30,0.6);
      border: 1px solid rgba(255,255,255,0.08);
      border-radius: 1rem;
      padding: 2rem;
      margin-bottom: 2rem;
    }
    .section-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 2rem;
      padding-bottom: 1rem;
      border-bottom: 1px solid rgba(255,255,255,0.1);
    }
    .section-header h2 {
      font-size: 1.5rem;
      font-weight: 600;
    }
    .search-bar {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      background: rgba(255,255,255,0.05);
      border: 1px solid rgba(255,255,255,0.1);
      border-radius: 0.5rem;
      padding: 0.5rem 1rem;
    }
    .search-bar input {
      background: none;
      border: none;
      color: #fff;
      outline: none;
      font-size: 0.9rem;
    }
    .search-bar input::placeholder {
      color: rgba(255,255,255,0.5);
    }
    .subscriber-count {
      color: #1d9bf0;
      font-weight: 600;
    }
    .loading {
      text-align: center;
      padding: 3rem;
      color: rgba(255,255,255,0.7);
      font-size: 1.1rem;
    }
    .users-grid {
      display: grid;
      gap: 1rem;
      grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
    }
    .user-card {
      background: rgba(40,40,40,0.8);
      border: 1px solid rgba(255,255,255,0.08);
      border-radius: 0.75rem;
      padding: 1.5rem;
      transition: transform 0.2s;
    }
    .user-card:hover {
      transform: translateY(-2px);
      background: rgba(45,45,45,0.8);
    }
    .user-info {
      display: flex;
      gap: 1rem;
      margin-bottom: 1rem;
    }
    .user-avatar {
      width: 60px;
      height: 60px;
      border-radius: 50%;
      object-fit: cover;
    }
    .user-details h3 {
      margin: 0 0 0.25rem 0;
      font-size: 1.1rem;
      font-weight: 600;
    }
    .user-details p {
      margin: 0.25rem 0;
      color: rgba(255,255,255,0.7);
      font-size: 0.9rem;
    }
    .user-details .email {
      color: #1d9bf0;
    }
    .user-details .date {
      font-size: 0.8rem;
      color: rgba(255,255,255,0.5);
    }
    .user-actions {
      display: flex;
      justify-content: space-between;
      align-items: center;
    }
    .status {
      padding: 0.25rem 0.75rem;
      border-radius: 9999px;
      font-size: 0.8rem;
      font-weight: 600;
      background: rgba(34, 197, 94, 0.2);
      color: #22c55e;
    }
    .status.banned {
      background: rgba(239, 68, 68, 0.2);
      color: #ef4444;
    }
    .action-btn {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      padding: 0.5rem 1rem;
      border: none;
      border-radius: 0.5rem;
      font-weight: 600;
      cursor: pointer;
      transition: all 0.2s;
    }
    .ban-btn {
      background: rgba(239, 68, 68, 0.2);
      color: #ef4444;
      border: 1px solid rgba(239, 68, 68, 0.3);
    }
    .ban-btn:hover {
      background: rgba(239, 68, 68, 0.3);
    }
    .unban-btn {
      background: rgba(34, 197, 94, 0.2);
      color: #22c55e;
      border: 1px solid rgba(34, 197, 94, 0.3);
    }
    .unban-btn:hover {
      background: rgba(34, 197, 94, 0.3);
    }
    .newsletter-form {
      max-width: 600px;
    }
    .form-group {
      margin-bottom: 1.5rem;
    }
    .form-group label {
      display: block;
      margin-bottom: 0.5rem;
      font-weight: 600;
      color: #fff;
    }
    .form-group input,
    .form-group textarea {
      width: 100%;
      background: rgba(255,255,255,0.05);
      border: 1px solid rgba(255,255,255,0.1);
      border-radius: 0.5rem;
      padding: 0.75rem;
      color: #fff;
      font-size: 0.9rem;
      outline: none;
      transition: border-color 0.2s;
    }
    .form-group input:focus,
    .form-group textarea:focus {
      border-color: #1d9bf0;
    }
    .form-group input::placeholder,
    .form-group textarea::placeholder {
      color: rgba(255,255,255,0.5);
    }
    .send-btn {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      background: #1d9bf0;
      border: none;
      color: #fff;
      padding: 0.75rem 2rem;
      border-radius: 0.5rem;
      font-weight: 600;
      cursor: pointer;
      transition: background 0.2s;
    }
    .send-btn:hover:not(:disabled) {
      background: #299fff;
    }
    .send-btn:disabled {
      opacity: 0.5;
      cursor: not-allowed;
    }
    .add-btn {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      background: #1d9bf0;
      border: none;
      color: #fff;
      padding: 0.5rem 1rem;
      border-radius: 0.5rem;
      font-weight: 600;
      cursor: pointer;
      transition: background 0.2s;
    }
    .add-btn:hover {
      background: #299fff;
    }
    .requests-list {
      display: flex;
      flex-direction: column;
      gap: 1rem;
    }
    .request-card {
      background: rgba(40,40,40,0.8);
      border: 1px solid rgba(255,255,255,0.08);
      border-radius: 0.75rem;
      padding: 1.5rem;
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      gap: 1rem;
    }
    .request-info {
      flex: 1;
    }
    .request-info h3 {
      margin: 0 0 0.5rem 0;
      font-size: 1.1rem;
      font-weight: 600;
    }
    .request-info .description {
      color: rgba(255,255,255,0.8);
      margin-bottom: 1rem;
      line-height: 1.5;
    }
    .request-info .email {
      color: #1d9bf0;
      margin-bottom: 1rem;
    }
    .request-meta {
      display: flex;
      flex-wrap: wrap;
      gap: 1rem;
      font-size: 0.8rem;
      color: rgba(255,255,255,0.6);
    }
    .request-meta strong {
      color: #fff;
    }
    .status-pending {
      background: rgba(234, 179, 8, 0.2);
      color: #eab308;
      padding: 0.25rem 0.5rem;
      border-radius: 0.25rem;
      font-weight: 600;
    }
    .status-approved {
      background: rgba(34, 197, 94, 0.2);
      color: #22c55e;
      padding: 0.25rem 0.5rem;
      border-radius: 0.25rem;
      font-weight: 600;
    }
    .status-rejected {
      background: rgba(239, 68, 68, 0.2);
      color: #ef4444;
      padding: 0.25rem 0.5rem;
      border-radius: 0.25rem;
      font-weight: 600;
    }
    .request-actions {
      display: flex;
      gap: 0.5rem;
      flex-shrink: 0;
    }
    .approve-btn {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      background: rgba(34, 197, 94, 0.2);
      border: 1px solid rgba(34, 197, 94, 0.3);
      color: #22c55e;
      padding: 0.5rem 1rem;
      border-radius: 0.5rem;
      font-weight: 600;
      cursor: pointer;
      transition: all 0.2s;
    }
    .approve-btn:hover {
      background: rgba(34, 197, 94, 0.3);
    }
    .reject-btn {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      background: rgba(239, 68, 68, 0.2);
      border: 1px solid rgba(239, 68, 68, 0.3);
      color: #ef4444;
      padding: 0.5rem 1rem;
      border-radius: 0.5rem;
      font-weight: 600;
      cursor: pointer;
      transition: all 0.2s;
    }
    .reject-btn:hover {
      background: rgba(239, 68, 68, 0.3);
    }
    .report-card .report-header {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      margin-bottom: 0.5rem;
      max-height: fit-content;
    }
    .reason-tag {
      background: rgba(249, 24, 128, 0.2);
      color: #f91880;
      padding: 0.25rem 0.5rem;
      border-radius: 0.25rem;
      font-size: 0.8rem;
      font-weight: 600;
    }
    .categories-grid {
      display: grid;
      gap: 1rem;
      grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    }
    .category-card {
      background: rgba(40,40,40,0.8);
      border: 1px solid rgba(255,255,255,0.08);
      border-radius: 0.75rem;
      padding: 1.5rem;
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      gap: 1rem;
      transition: transform 0.2s;
    }
    .category-card:hover {
      transform: translateY(-2px);
      background: rgba(45,45,45,0.8);
    }
    .category-info {
      flex: 1;
    }
    .category-info h3 {
      margin: 0 0 0.5rem 0;
      font-size: 1.1rem;
      font-weight: 600;
    }
    .category-info p {
      color: rgba(255,255,255,0.8);
      margin-bottom: 1rem;
      line-height: 1.5;
    }
    .category-info .date {
      font-size: 0.8rem;
      color: rgba(255,255,255,0.5);
    }
    .category-actions {
      display: flex;
      gap: 0.5rem;
      flex-shrink: 0;
    }
    .edit-btn {
      background: rgba(59, 130, 246, 0.2);
      border: 1px solid rgba(59, 130, 246, 0.3);
      color: #3b82f6;
      padding: 0.5rem;
      border-radius: 0.5rem;
      cursor: pointer;
      transition: all 0.2s;
    }
    .edit-btn:hover {
      background: rgba(59, 130, 246, 0.3);
    }
    .delete-btn {
      background: rgba(239, 68, 68, 0.2);
      border: 1px solid rgba(239, 68, 68, 0.3);
      color: #ef4444;
      padding: 0.5rem;
      border-radius: 0.5rem;
      cursor: pointer;
      transition: all 0.2s;
    }
    .delete-btn:hover {
      background: rgba(239, 68, 68, 0.3);
    }
    .modal-overlay {
      position: fixed;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background: rgba(0, 0, 0, 0.8);
      display: flex;
      align-items: center;
      justify-content: center;
      z-index: 1000;
      backdrop-filter: blur(4px);
    }
    .modal {
      background: rgba(30, 30, 30, 0.95);
      border: 1px solid rgba(255, 255, 255, 0.1);
      border-radius: 1rem;
      padding: 1vw;
      width: 90%;
      max-width: 500px;
      max-height: 80vh;
      overflow: hidden;
      box-shadow: 0 20px 60px rgba(0, 0, 0, 0.5);
    }
    .modal-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 1.5rem;
      border-bottom: 1px solid rgba(255, 255, 255, 0.1);
    }
    .modal-header h3 {
      margin: 0;
      font-size: 1.2rem;
      font-weight: 600;
    }
    .modal-close {
      background: none;
      border: none;
      color: rgba(255, 255, 255, 0.7);
      font-size: 1.5rem;
      cursor: pointer;
      padding: 0;
      width: 2rem;
      height: 2rem;
      display: flex;
      align-items: center;
      justify-content: center;
      border-radius: 0.25rem;
      transition: all 0.2s;
    }
    .modal-close:hover {
      background: rgba(255, 255, 255, 0.1);
      color: #fff;
    }
    .modal-body {
      padding: 1.5rem;
      max-height: 60vh;
      overflow-y: auto;
    }
    .request-header{
      width: 100%;
      display: flex;
      flex-direction: column;
    }
    .modal-footer {
      display: flex;
      justify-content: flex-end;
      gap: 1rem;
      padding: 1.5rem;
      border-top: 1px solid rgba(255, 255, 255, 0.1);
    }
    .cancel-btn {
      background: none;
      border: 1px solid rgba(255, 255, 255, 0.2);
      color: rgba(255, 255, 255, 0.8);
      padding: 0.5rem 1rem;
      border-radius: 0.5rem;
      cursor: pointer;
      transition: all 0.2s;
    }
    .cancel-btn:hover {
      background: rgba(255, 255, 255, 0.05);
      color: #fff;
    }
    .save-btn {
      background: #1d9bf0;
      border: none;
      color: #fff;
      padding: 0.5rem 1rem;
      border-radius: 0.5rem;
      font-weight: 600;
      cursor: pointer;
      transition: background 0.2s;
    }
    .save-btn:hover {
      background: #299fff;
    }
    .sidebar {
      z-index: 400000;
    }
    @media (max-width: 768px) {
      .main {
        width: 95vw;
        max-width: 95vw;
        padding: 0.5rem;
      }
      .tabs {
        flex-direction: column;
        gap: 0.25rem;
      }
      .tabs button {
        justify-content: center;
        padding: 1rem;
      }
      .section-header {
        flex-direction: column;
        align-items: flex-start;
        gap: 1rem;
      }
      .users-grid {
        grid-template-columns: 1fr;
      }
      .categories-grid {
        grid-template-columns: 1fr;
      }
      .request-card {
        flex-direction: column;
        align-items: stretch;
      }
      .request-actions {
        justify-content: flex-end;
      }
      .category-card {
        flex-direction: column;
        align-items: stretch;
      }
      .category-actions {
        justify-content: flex-end;
      }
      .modal {
        width: 95%;
        margin: 1rem;
      }
    }
    .request-card {
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    background: rgba(40,40,40,0.8);
    border: 1px solid rgba(255,255,255,0.08);
    border-radius: 0.75rem;
    padding: 1.5rem;
    gap: 1rem;
    min-height: 180px;
    position: relative;
  }

  .request-bottom-row {
    display: flex;
    align-items: flex-end;
    gap: 1rem;
    margin-top: auto;
  }

  .face-verify-img {
    max-width: 160px;
    max-height: 160px;
    border-radius: 0.35rem;
    border: 1.5px solid #222e;
    background: #222;
    box-shadow: 0 1px 6px 0 rgba(30,30,30,0.09);
    flex-shrink: 0;
    margin-right: 0.7rem;
  }
</style>