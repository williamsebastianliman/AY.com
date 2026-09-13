<script lang="ts">
    import { writable } from "svelte/store";
    import { clickOutside } from "../action/clickOutside";
    import { onMount } from "svelte";
    import api from "../api";
    import CreateNewThread from "./CreateNewThread.svelte";
    import { Home, Search, Bell, MessageCircle, Bookmark, Users, Crown, User, Settings, Moon, Sun, Menu } from "lucide-svelte";
  
    export const darkMode = writable(false);
  
    interface User {
      name: string;
      username: string;
      avatar: string;
      role: string;
    }
    const currentUser: User = { name: "", username: "", avatar: "", role: "" };
    export let currentUserId: string;
    export let activePage: string;
  
    let showPostModal = false;
    let showSidebar = false;
    let sectionEl!: HTMLElement;
    let showLogoutMenu = false;
  
    async function fetchUserAndAvatar() {
      try {
        const res = await api.get<{ name: string; username: string; profile_picture_id: string }>(`/user/${currentUserId}`);
        currentUser.name = res.data.name;
        currentUser.username = res.data.username;
        const pp = await api.post<{ public_url: string }>("/media/get-media", { id: res.data.profile_picture_id });
        currentUser.avatar = pp.data.public_url;
      } catch (err) {
        console.error(err);
      }
    }
  
    onMount(fetchUserAndAvatar);
  
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
    if (currentUser.role === "Admin") {
      navItems.push({ label: "Admin", href: "/admin", icon: User, key: "admin" });
    }
  
    function toggleSidebar() {
      showSidebar = !showSidebar;
    }
  </script>
  
  <!-- Burger button -->
  <button class="burger-btn" on:click={toggleSidebar} aria-label="Toggle sidebar">
    <Menu size="24" />
  </button>
  
  {#if showSidebar}
    <div class="sidebar-backdrop" on:click={toggleSidebar} />
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
        {#each navItems as item (item.key)}
          <div class="nav-item" class:active={activePage === item.key} on:click={() => window.location.href = item.href}>
            <svelte:component this={item.icon} size="20" />
            <span>{item.label}</span>
          </div>
        {/each}
      </nav>
      <button class="post-btn" on:click={() => showPostModal = true}>Post</button>
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
            <button class="logout-btn">Log Out</button>
          </div>
        {/if}
      </div>
      {#if showPostModal}
        <div class="modal-backdrop" on:click={() => showPostModal = false}>
          <div class="modal" on:click|stopPropagation>
            <CreateNewThread currentUserId={currentUserId} isAdmin={false} on:close={() => showPostModal = false} />
          </div>
        </div>
      {/if}
    </div>
  {/if}
  <style>
  .burger-btn {
    position: fixed;
    top: 1rem;
    left: 1rem;
    background: none;
    border: none;
    z-index: 1000;
    cursor: pointer;
  }
  .sidebar-backdrop {
    position: fixed;
    inset: 0;
    background: rgba(0,0,0,0.4);
    z-index: 900;
  }
  .sidebar {
    position: fixed;
    top: 0;
    left: 0;
    width: 250px;
    height: 97vh;
    background: var(--sidebar-bg);
    padding: 1rem;
    z-index: 1000;
    display: flex;
    flex-direction: column;
  }
  .header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 1rem; }
  .logo { font-size: 1.5rem; font-weight: bold; }
  .toggle-btn { background: none; border: none; cursor: pointer; }
  .nav { flex-grow: 1; }
  .nav-item { display: flex; align-items: center; gap: 0.8rem; margin: 0.5rem 0; cursor: pointer; padding: 0.6rem 1rem; border-radius: 9999px; }
  .nav-item.active { background: var(--primary); color: #fff; }
  .post-btn { margin: 1rem 0; padding: 0.5rem 1rem; border: none; border-radius: 9999px; background: var(--primary); color: #fff; }
  .user-section { margin-top: auto; position: relative; }
  .user-info { display: flex; align-items: center; cursor: pointer; }
  .avatar { width: 40px; height: 40px; border-radius: 50%; margin-right: 0.5rem; }
  .user-text { display: flex; flex-direction: column; }
  .logout-menu { position: absolute; bottom: 100%; left: 0; background: var(--dropdown-bg); border-radius: 0.5rem; padding: 0.5rem; }
  .logout-btn { background: none; border: none; width: 100%; text-align: left; cursor: pointer; }
  .modal-backdrop { position: fixed; inset: 0; background: rgba(0,0,0,0.5); display: flex; justify-content: center; align-items: center; z-index: 1100; }
  .modal { background: var(--bg); padding: 1rem; border-radius: 0.5rem; width: 90%; max-width: 400px; }
  </style>