<script lang="ts">
  import { onMount } from "svelte";
  import logo from "../assets/logo_ay.jpg";
  import { navigate } from "svelte-routing";
  import axios from "axios";
  import { toast } from "../lib/toastStore";

  let username = "";
  let password = "";
  let captchaVerified = false;
  let loginError = "";
  onMount(() => {
    if (!window.grecaptcha) {
      const script = document.createElement("script");
      script.src = "https://www.google.com/recaptcha/api.js";
      script.async = true;
      script.defer = true;
      document.body.appendChild(script);
    }

    window.onCaptchaSuccess = () => {
      captchaVerified = true;
    };
  });
  async function handleLogin() {
  try {
    await axios.post(
      "http://localhost:8080/api/v1/auth/login",
      { username, password },
      { withCredentials: true }
    );
    toast("Login successful!", 2000);
    navigate("/home");
  } catch (err) {
    const msg = axios.isAxiosError(err)
      ? err.response?.data?.error ?? err.message
      : err instanceof Error
      ? err.message
      : "Unexpected error";
    loginError = "Invalid Credential!";
    toast(`Login failed: ${msg}`, 3000);
  }
}
</script>

<div class="main-container">
  <img src={logo} alt="Logo" class="logo" />
  <div class="form">
    <h1 class="main-title">Log In</h1>

    <div class="form-group">
      <h2 class="label">Username</h2>
      <input
        type="text"
        bind:value={username}
        placeholder="Username"
        class="input"
      />
    </div>

    <div class="form-group">
      <h2 class="label">Password</h2>
      <input
        type="password"
        bind:value={password}
        placeholder="Password"
        class="input"
      />
    </div>

    <div class="form-group captcha-container">
      <div
        class="g-recaptcha"
        data-sitekey="6LddEycrAAAAAC-N3vVTDhSQRVJbo8zfWpDSuFEv"
        data-callback="onCaptchaSuccess"
      ></div>
    </div>

    <button
      type="button"
      class="next-button"
      on:click={handleLogin}
      disabled={!captchaVerified || !username || !password}
    >
      Log In
    </button>
    {#if loginError}
      <div class="login-error">{loginError}</div>
    {/if}
    <div class="links">
      <a href="/">Back to Landing Page</a>
      <a href="/forgot-password">Forgotten Account?</a>
    </div>
  </div>
</div>

<style>
  .main-container {
    max-width: 400px;
    margin: 0 auto;
    padding: 2rem 1rem;
  }

  .logo {
    display: block;
    margin: 0 auto 1rem;
    width: 70px;
  }

  .form {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    background: var(--card-bg-color);
    padding: 2rem;
    border-radius: 0.5rem;
    box-shadow: 0 2px 8px rgba(0,0,0,0.2);
    border: 1px solid var(--border-color);
  }

  .main-title {
    font-size: 1.8rem;
    font-weight: 700;
    text-align: center;
    color: var(--text-color);
  }

  .form-group {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .label {
    font-weight: bold;
    font-size: 1rem;
    color: var(--text-color);
  }

  .input {
    padding: 0.75rem;
    font-size: 1rem;
    border: 1px solid var(--border-color);
    border-radius: 0.5rem;
    background: var(--bg-color);
    color: var(--text-color);
  }

  .captcha-container {
    display: flex;
    justify-content: center;
    margin: 1rem 0;
  }

  .next-button {
    padding: 1rem;
    font-size: 1rem;
    font-weight: bold;
    border: none;
    border-radius: 9999px;
    background-color: var(--primary-color);
    color: #ffffff;
    cursor: pointer;
    transition: background-color 0.2s;
  }

  .next-button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .links {
    display: flex;
    justify-content: space-between;
    margin-top: 1rem;
    font-size: 0.9rem;
  }

  .links a {
    color: var(--primary-color);
    text-decoration: none;
  }

  .links a:hover {
    text-decoration: underline;
  }
  .login-error {
    color: #f44336;
    background: #fff0f0;
    border: 1px solid #f44336;
    border-radius: 0.4rem;
    margin: 0.7rem 0 0.2rem 0;
    padding: 0.6rem 1rem;
    text-align: center;
    font-size: 1rem;
  }
  @media (max-width: 500px) {
  html, body {
    height: 100%;
    margin: 0;
    padding: 0;
  }
  .main-container {
    max-width: 100vw;
    height: 100vh;
    min-height: 100vh;
    padding: 1.5rem 0.5rem;
    box-sizing: border-box;
    overflow-y: auto;
  }
  .form {
    padding: 1.3rem 0.7rem;
    box-sizing: border-box;
    min-width: 0;
  }
}
</style>