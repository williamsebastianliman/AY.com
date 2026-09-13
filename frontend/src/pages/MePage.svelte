<script lang="ts">
  import { onMount } from "svelte";
  import { toast } from "../lib/toastStore";
  import api from "../lib/api";

  type User = {
    id: string;
    name: string;
    username: string;
    email: string;
    gender: string;
    birth_year: string;
    birth_month: string;
    birth_day: string;
    subscribed_news: boolean;
    profile_picture_id: string;
    banner_media_id: string;
  };

  let user: User | null = null;

  onMount(async () => {
    try {
      const res = await api.get<{
        id: string;
        name: string;
        username: string;
        email: string;
        gender: string;
        birth_year: string;
        birth_month: string;
        birth_day: string;
        subscribed_news: boolean;
        profile_picture_id: string;
        banner_media_id: string;
      }>("/user/get-me");
      user = res.data;
    } catch (err) {
      const msg = err instanceof Error
        ? err.message
        : "Failed to load user info";
      toast(msg, 3000);
    }
  });
</script>

{#if user}
  <div class="container">
    <div class="card">
      <h1 class="heading">Welcome, {user.name}!</h1>
      <ul class="info-list">
        <li><strong>ID:</strong> {user.id}</li>
        <li><strong>Username:</strong> {user.username}</li>
        <li><strong>Email:</strong> {user.email}</li>
        <li><strong>Gender:</strong> {user.gender}</li>
        <li>
          <strong>Birth:</strong>
          {user.birth_day}/{user.birth_month}/{user.birth_year}
        </li>
        <li>
          <strong>Subscribed News:</strong>
          {user.subscribed_news ? "Yes" : "No"}
        </li>
        <li><strong>Profile Picture ID:</strong> {user.profile_picture_id}</li>
        <li><strong>Banner Media ID:</strong> {user.banner_media_id}</li>
      </ul>
    </div>
  </div>
{:else}
  <div class="container">
    <p class="loading">Loading user info...</p>
  </div>
{/if}

<style>
  .container {
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 80vh;
    background-color: var(--bg-color);
    color: var(--text-color);
    padding: 2rem;
  }

  .card {
    background-color: var(--card-bg-color);
    padding: 2rem;
    border-radius: 1rem;
    border: 1px solid var(--border-color);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    width: 100%;
    max-width: 600px;
  }

  .heading {
    font-size: 1.8rem;
    font-weight: bold;
    margin-bottom: 1rem;
    color: var(--text-color);
    text-align: center;
  }

  .info-list {
    list-style: none;
    padding: 0;
    font-size: 1rem;
    line-height: 1.5;
  }

  .info-list li {
    margin-bottom: 0.5rem;
    color: var(--text-color);
  }

  .loading {
    font-size: 1.5rem;
    font-weight: 600;
    text-align: center;
  }
</style>