<script lang="ts">
  import { writable } from "svelte/store";
  import { clickOutside } from "../action/clickOutside";
  import { onMount } from "svelte";
  import api from "../api";
  import CreateNewThread from "./CreateNewThread.svelte";
  import {
    Home,
    Search,
    Bell,
    MessageCircle,
    Bookmark,
    Users,
    Crown,
    User,
    Settings,
    Moon,
    Sun,
  } from "lucide-svelte";

  async function logout() {
    try {
      await api.post("/logout");
      window.location.href = "/login";
    } catch (err) {
      alert("Failed to logout. Try again!");
      console.error(err);
    }
  }
  export const darkMode = writable(false);
  interface User {
    name: string,
    username: string,
    avatar: string,
    role: string,
  }
  const currentUser: User = { name: "", username: "", avatar: "", role:"" };
  export let currentUserId: string;
  export let activePage: string;
  let showPostModal = false;
  async function fetchUserAndAvatar() {
    try {
      console.log("uid: ",currentUserId);
      const res = await api.get<{ name: string, username: string, profile_picture_id: string }>(`/user/${currentUserId}`);
      currentUser.name = res.data.name;
      currentUser.username = res.data.username;
      try {
        const pp = await api.post<{ public_url: string }>("/media/get-media", { id: res.data.profile_picture_id });
        currentUser.avatar = pp.data.public_url;
      } catch (err) {
        console.error(err as string);
      }
    } catch (err) {
      console.error(err as string);
    }
  }

  let sectionEl!: HTMLElement;
  let showLogoutMenu = false;

  onMount(() => { fetchUserAndAvatar(); });

  const navItems = [
    { label: "Home", href: "/home", icon: Home, key: "home" },
    { label: "Explore", href: "/explore", icon: Search, key: "explore" },
    { label: "Notifications", href: "/notifications", icon: Bell, key: "notifications" },
    { label: "Messages", href: "/messages", icon: MessageCircle, key: "messages" },
    { label: "Bookmarks", href: "/bookmark", icon: Bookmark, key: "bookmarks" },
    { label: "Communities", href: "/community", icon: Users, key: "communities" },
    { label: "Premium", href: "/premium", icon: Crown, key: "premium" },
    { label: "Profile", href: `/profile/${currentUserId}`, icon: User, key: "profile" },
    { label: "Settings", href: "/settings", icon: Settings, key: "settings" },
  ];
  if (currentUser.role == "Admin"){
    navItems.push({ label: "Admin", href: "/admin", icon: User, key: "admin" });
  }
</script>

<div class="sidebar">
  <div class="header">
    <h1 class="logo">AY</h1>
    <button type="button" class="toggle-btn" on:click={() => darkMode.update(v => !v)} aria-label="Toggle dark mode">
      {#if $darkMode}
        <Moon size="22" />
      {:else}
        <Sun size="22" />
      {/if}
    </button>
  </div>
  <nav class="nav">
    {#each navItems as item (item.href)}
      <div
        class="nav-item"
        class:active={activePage === item.key}
        on:click={() => window.location.href = item.href}
        aria-current={activePage === item.key ? "page" : undefined}
      >
        <svelte:component this={item.icon} size="20" class="nav-icon" />
        <span>{item.label}</span>
      </div>
    {/each}
  </nav>
  <button type="button" class="post-btn" on:click={() => showPostModal = true}>Post</button>
  <div class="user-section" use:clickOutside bind:this={sectionEl}>
    <div class="user-info" on:click={() => showLogoutMenu = !showLogoutMenu}>
      <img src={currentUser.avatar} alt="avatar" class="avatar" />
      <div class="user-text">
        <div class="user-name">{currentUser.name}</div>
        <div class="user-handle">{currentUser.username}</div>
      </div>
    </div>
    {#if showLogoutMenu}
      <div class="logout-menu">
        <button type="button" class="logout-btn" on:click={logout}>Log Out</button>
      </div>
    {/if}
  </div>
  {#if showPostModal}
    <div class="modal-backdrop" on:click={() => showPostModal = false}>
      <div on:click|stopPropagation>
        <CreateNewThread
          currentUserId
          isAdmin={false}
          on:close={() => showPostModal = false}
        />
      </div>
    </div>
  {/if}

</div>

<style>
  /* @media (max-width: 1050px) {
    .sidebar {
      display: none
    }
  } */
  .sidebar {
    display: flex;
    flex-direction: column;
    width: 250px;
    height: 95vh;
    overflow-y: auto;
    padding: 1rem;
    background: var(--sidebar-bg);
  }
  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
    z-index: 300000;
  }
  .logo {
    font-size: 1.5rem;
    font-weight: bold;
  }
  .toggle-btn {
    background: none;
    border: none;
    cursor: pointer;
    font-size: 1.25rem;
    display: flex;
    align-items: center;
    padding: 0.2rem;
  }
  .nav {
    flex-grow: 1;
  }
  .nav-item {
    display: flex;
    align-items: center;
    gap: 0.8rem;
    margin: 0.5rem 0;
    font-size: 1.1rem;
    cursor: pointer;
    color: var(--text-secondary);
    border-radius: 9999px;
    padding: 0.6rem 1rem;
    transition: background 0.12s, color 0.12s;
  }
  .nav-item.active {
    background: var(--primary, #1d9bf0);
    color: #fff;
  }
  .nav-icon {
    transition: color 0.12s;
  }
  .nav-item.active .nav-icon {
    color: #fff;
  }
  .post-btn {
    padding: 0.5rem 1rem;
    margin: 1rem 0;
    border: none;
    border-radius: 9999px;
    background-color: var(--primary);
    color: var(--on-primary);
    cursor: pointer;
  }
  .user-section {
    margin-top: auto;
    position: relative;
  }
  .user-info {
    display: flex;
    align-items: center;
    cursor: pointer;
  }
  .avatar {
    width: 40px;
    height: 40px;
    border-radius: 9999px;
    margin-right: 0.5rem;
  }
  .user-text {
    display: flex;
    flex-direction: column;
  }
  .user-name {
    font-weight: bold;
  }
  .user-handle {
    font-size: 0.85rem;
    color: var(--text-secondary);
  }
  .logout-menu {
    position: absolute;
    bottom: 100%;
    left: 0;
    background: var(--dropdown-bg);
    border: 1px solid var(--border);
    border-radius: 0.5rem;
    padding: 0.5rem;
  }
  .logout-btn {
    background: none;
    border: none;
    padding: 0.25rem 0.5rem;
    width: 100%;
    text-align: left;
    cursor: pointer;
  }
  .modal-backdrop {
    position: fixed;
    inset: 0;
    background: rgba(0,0,0,0.5);
    display: flex;
    justify-content: center;
    align-items: center;
  }
  .modal {
    background: var(--bg);
    padding: 1rem;
    border-radius: 0.5rem;
    width: 90%;
    max-width: 400px;
  }
  .modal-textarea {
    width: 100%;
    height: 6rem;
    margin: 0.5rem 0;
    padding: 0.5rem;
    border: 1px solid var(--border);
    border-radius: 0.25rem;
  }
  .submit-btn {
    padding: 0.5rem 1rem;
    border: none;
    border-radius: 9999px;
    background-color: var(--primary);
    color: var(--on-primary);
    cursor: pointer;
  }
</style>