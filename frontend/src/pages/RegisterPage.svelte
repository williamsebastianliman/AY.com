<script lang="ts" context="module">
  declare global {
    interface Window {
      grecaptcha: unknown;
      onCaptchaSuccess: (token: string) => void;
    }
  }
  export {};
</script>

<script lang="ts">
  import { onMount } from "svelte";
  import logo from "../assets/logo_ay.jpg";
  import StepBar from "../lib/components/StepBar.svelte";
  import axios from "axios";
  import { toast } from "../lib/toastStore";
  import Toaster from "../lib/components/Toaster.svelte";
  import { navigate } from "svelte-routing";

  let step = 1;
  let name = "";
  let username = "";
  let email = "";
  let password = "";
  let month = "";
  let day = "";
  let year = "";
  let gender = "";
  let requestId = "";
  let confirmed_password = "";
  let subscribeNewsletter = false;
  let userId = "";

  let profileFile: File | null = null;
  let bannerFile: File | null = null;
  let profilePreview = "";
  let bannerPreview = "";
  let profileMediaId = "";
  let bannerMediaId = "";
  let captchaVerified = false;

  let nameError = "";
  let usernameError = "";
  let emailError = "";
  let passwordError = "";
  let confirmError = "";
  let dobError = "";
  let genderError = "";
  let profileError = "";
  let bannerError = "";

  let otpDigits: string[] = Array(6).fill("");
  $: otp = otpDigits.join("");

  const securityQuestions = [
    "What was your first pet's name?",
    "What's your mother's maiden name?",
    "What city were you born in?",
    "What was your high school mascot?",
    "What's your favorite food?",
    "What street did you grow up on?"
  ];

  function clearErrors() {
    nameError = usernameError = emailError = passwordError = confirmError = dobError = genderError = profileError = bannerError = "";
  }

  function validateName() {
    if (name.trim().length < 5 || !/^[A-Za-z\s]+$/.test(name)) {
      nameError = "Name must be >= 5 letters and contain only letters/spaces";
    }
  }

  async function validateUsername() {
    if (!username.trim()) {
      usernameError = "Username cannot be blank";
    } else if (!(await checkUsernameUnique(username))) {
      console.log("username alr registered!");
      usernameError = "Username is already taken";
    }
    else{
      console.log("username not regsitered");
    }
  }

  async function validateEmail() {
    const re = /^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.com$/;
    if (!re.test(email)) {
      emailError = "Email must look like user@domain.com";
    } else if (!(await checkEmailUnique(email))) {
      console.log("email alr registered!");
      emailError = "Email is already registered";
    }
    else{
      console.log("email not regsitered");
    }
  }
  function validateProfilePicture() {
    if (!profileFile) {
      profileError = "Profile picture is required";
    }
  }

  function validateBannerPicture() {
    if (!bannerFile) {
      bannerError = "Banner picture is required";
    }
  }
  function validatePassword() {
    const rules: ([RegExp, string])[] = [
      [/.{8,}/, "At least 8 characters"],
      [/[A-Z]/,   "Must include an uppercase letter"],
      [/[a-z]/,   "Must include a lowercase letter"],
      [/\d/,      "Must include a digit"],
      [/\W/,      "Must include a special character"],
    ];
    for (const [pattern, msg] of rules) {
      if (!pattern.test(password)) {
        passwordError = msg;
        break;
      }
    }
  }

  function validateConfirm() {
    if (password !== confirmed_password) {
      confirmError = "Passwords do not match";
    }
  }

  function validateDOB() {
    if (!year || !month || !day) {
      dobError = "Complete your date of birth";
      return;
    }
    const dob = new Date(year, month - 1, day);
    const age = (Date.now() - dob.getTime()) / (365.25 * 24 * 3600 * 1000);
    if (age < 13) {
      dobError = "You must be at least 13 years old";
    }
  }

  function validateGender() {
    if (gender !== "Male" && gender !== "Female") {
      genderError = "Please select Male or Female";
    }
  }
  export async function checkUsernameUnique(u: string): Promise<boolean> {
    try {
      const res = await fetch("http://localhost:8080/api/v1/user/get-by-username", {
        method: "POST",
        headers: { "Content-Type": "application/json"},
        body: JSON.stringify({ username: u }),
      });

      if (res.status === 404) {
        return true;
      }

      if (!res.ok) {
        emailError = "Username is Already Taken!";
        return false;
      }
      emailError = "Username is Already Taken!";
      return false;
    } catch (err) {
      console.error(err);
      return false;
    }
  }
  export async function checkEmailUnique(e: string): Promise<boolean> {
    try {
      const res = await fetch("http://localhost:8080/api/v1/user/get-by-email", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email: e }),
      });

      if (res.status === 404) {
        return true;
      }

      if (!res.ok) {
        emailError = "Username is Already Taken!";
        return false;
      }
      emailError = "Username is Already Taken!";
      return false;
    } catch (err) {
      console.error(err);
      return false;
    }
  }

  type QA = { question: string; answer: string };
  let securityQAs: QA[] = [{ question: "", answer: "" }];

  function addSecurity(): void {
    securityQAs = [...securityQAs, { question: "", answer: "" }];
  }

  function removeSecurity(idx: number): void {
    securityQAs = securityQAs.filter((_, i) => i !== idx);
  }

  function handleQAChange(idx: number, field: keyof QA, value: string): void {
    const updatedQAs = [...securityQAs];
    updatedQAs[idx] = { ...updatedQAs[idx], [field]: value };
    securityQAs = updatedQAs;
  }

  async function uploadSingleMedia(file: File) {
    const form = new FormData();
    form.append("filename", file.name);
    form.append("file", file);

    const { data } = await axios.post(
      "http://localhost:8080/api/v1/media/upload",
      form
    );
    return data as { id: string; public_url: string; created_at: string };
  }

  async function handleProfileChange(e: Event) {
    const files = (e.target as HTMLInputElement).files;
    if (!files?.[0]) return;

    profileFile = files[0];
    profilePreview = URL.createObjectURL(profileFile);

    try {
      const resp = await uploadSingleMedia(profileFile);
      profileMediaId = resp.id;
      toast("Profile image uploaded", 2000);
    } catch (error: unknown) {
      if (axios.isAxiosError(error)) {
        const errBody = error.response?.data as { error?: string };
        const msg = errBody.error ?? error.message;
        toast(`Profile upload failed: ${msg}`, 3000);
      } else if (error instanceof Error) {
        toast(`Profile upload failed: ${error.message}`, 3000);
      } else {
        toast("Profile upload failed", 3000);
      }
    }
  }
    async function submitSecurityAnswers() {
    try {
      await Promise.all(
        securityQAs.map(({ question, answer }) =>
          axios.post(
            "http://localhost:8080/api/v1/security-answer/create",
            { user_id: userId, question, answer }
          )
        )
      );
      toast("Security questions saved!", 2000);
      step=3;
    } catch (err: unknown) {
      const msg = axios.isAxiosError(err)
        ? err.response?.data?.error ?? err.message
        : err instanceof Error
        ? err.message
        : "Unknown error";
      toast(`Failed to save security Q&A: ${msg}`, 3000);
    }
  }
  async function handleBannerChange(e: Event) {
    const files = (e.target as HTMLInputElement).files;
    if (files?.[0]) {
      bannerFile = files[0];
      bannerPreview = URL.createObjectURL(bannerFile);
      try {
        const resp = await uploadSingleMedia(bannerFile);
        bannerMediaId = resp.id;
        toast("Banner image uploaded", 2000);
      } catch {
        toast("Banner upload failed", 3000);
      }
    }
  }

  async function nextStep() {
    clearErrors();
    validateName();
    await validateUsername();
    await validateEmail();
    validatePassword();
    validateConfirm();
    validateDOB();
    validateGender();
    validateBannerPicture();
    validateProfilePicture();

    const hasError = [
      nameError,
      usernameError,
      emailError,
      passwordError,
      confirmError,
      dobError,
      genderError
    ].some((e) => e);
    if(hasError){
      toast("Some Field Have Error!",3000);
      return;
    }
    const payload = {
      name,
      username,
      email,
      password,
      gender,
      birth_year: String(year),
      birth_month: String(month),
      birth_day: String(day),
      profile_picture_id: profileMediaId,
      banner_media_id: bannerMediaId,
      subscribed_news : subscribeNewsletter
    };
    try {
      const response = await axios.post(
        "http://localhost:8080/api/v1/user/register",
        payload
      );
      userId = response.data.id;
      step = 2;
    } catch (err: unknown) {
      if (axios.isAxiosError(err)) {
        alert("Request failed: " + (err.response?.data?.error ?? err.message));
      } else if (err instanceof Error) {
        alert("Unexpected error: " + err.message);
      } else {
        alert("An unknown error occurred");
      }
    }
  }
  async function prevStep() {
    if(step > 1){
      step--;
    }
  }
  async function sendOtp() {
    try{
      const response = await axios.post(
        "http://localhost:8080/api/v1/auth/send-otp",
        {id: userId, email:email}
      );
      requestId = response.data.request_id;
    } catch(err: unknown){
      alert("Request Failed: "+err as string);
    }
  }
  async function submitOtp() {
    try {
      const response = await axios.post(
        "http://localhost:8080/api/v1/auth/verify-otp",
        { request_id: requestId, code: otp }
      );
      const { success } = response.data as { success: boolean };
      toast(success ? "OTP verified!" : "Wrong OTP, please try again.", 3000);
      if (success) {
        navigate("/login");
      }
    } catch (err: unknown) {
      if (axios.isAxiosError(err)) {
        toast(
          "Verification failed: " + (err.response?.data?.error ?? err.message),
          3000
        );
      } else if (err instanceof Error) {
        toast("Error: " + err.message, 3000);
      } else {
        toast("An unknown error occurred.", 3000);
      }
    }
  }

  function handleOtpInput(e: Event, idx: number) {
    const input = e.target as HTMLInputElement;
    const val = input.value.replace(/\D/g, "").slice(-1);
    otpDigits[idx] = val;
    if (val && idx < 5) {
      const next = document.getElementById(`otp-${idx + 1}`) as HTMLInputElement;
      next?.focus();
    }
  }

  function handleOtpKeyDown(e: KeyboardEvent, idx: number) {
    const input = e.target as HTMLInputElement;
    if (e.key === "Backspace" && !input.value && idx > 0) {
      const prev = document.getElementById(`otp-${idx - 1}`) as HTMLInputElement;
      prev?.focus();
    }
  }

  function resendOtp() {
    alert(`Resending code to ${email}`);
    otpDigits = Array(6).fill("");
    setTimeout(() => {
      const first = document.getElementById("otp-0") as HTMLInputElement;
      first?.focus();
    }, 0);
  }

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
  $: if (step === 3){
    sendOtp();
  }
</script>

<div class="main-container">
  <Toaster />
  <img src={logo} alt="Logo" class="logo" />
  <StepBar currentStep={step} />

  {#if step === 1}
    <div class="form">
      <h1 class="main-title">Create your account</h1>

      <div class="form-group">
        <h2 class="label">Name</h2>
        <div class="hinted-field">
          <input
            type="text"
            bind:value={name}
            placeholder="Name"
            maxlength="50"
            class="input"
          />
          <div class="input-counter">{name.length} / 50</div>
        </div>
        {#if nameError}<div class="error">{nameError}</div>{/if}
      </div>

      <div class="form-group">
        <h2 class="label">Username</h2>
        <div class="hinted-field">
          <input
            type="text"
            bind:value={username}
            placeholder="Username"
            maxlength="50"
            class="input"
          />
          <div class="input-counter">{username.length} / 50</div>
        </div>
        {#if usernameError}<div class="error">{usernameError}</div>{/if}
      </div>

      <div class="form-group">
        <h2 class="label">Email</h2>
        <input
          type="email"
          bind:value={email}
          placeholder="Email"
          class="input"
        />
        {#if emailError}<div class="error">{emailError}</div>{/if}
      </div>

      <div class="form-group">
        <h2 class="label">Password</h2>
        <input
          type="password"
          bind:value={password}
          placeholder="Password"
          class="input"
        />
        {#if passwordError}<div class="error">{passwordError}</div>{/if}
      </div>

      <div class="form-group">
        <h2 class="label">Confirm Password</h2>
        <input
          type="password"
          bind:value={confirmed_password}
          placeholder="Password"
          class="input"
        />
        {#if confirmError}<div class="error">{confirmError}</div>{/if}
      </div>

      <div class="form-group">
        <h2 class="label">Date of birth</h2>
        <p class="description">
          This will not be shown publicly. Confirm your age—even if this
          account is for a business, a pet, etc.
        </p>
        <div class="dob-selects">
          <select bind:value={month} class="select">
            <option disabled selected value="">Month</option>
            {#each Array.from({ length: 12 }, (_, i) => i + 1) as m (m)}
              <option value={m}>{m}</option>
            {/each}
          </select>
          <select bind:value={day} class="select">
            <option disabled value="">Day</option>
            {#each Array.from({ length: 31 }, (_, i) => i + 1) as d (d)}
              <option value={d}>{d}</option>
            {/each}
          </select>
          <select bind:value={year} class="select">
            <option disabled value="">Year</option>
            {#each Array.from({ length: 100 }, (_, i) => 2025 - i) as y (y)}
              <option value={y}>{y}</option>
            {/each}
          </select>
        </div>
        {#if dobError}<div class="error">{dobError}</div>{/if}
      </div>

      <div class="form-group">
        <h2 class="label">Gender</h2>
        <select bind:value={gender} class="select">
          <option disabled selected value="">Select gender</option>
          <option>Male</option>
          <option>Female</option>
        </select>
        {#if genderError}<div class="error">{genderError}</div>{/if}
      </div>

      <div class="form-group">
        <h2 class="label">Profile Picture</h2>
        <input
          type="file"
          accept="image/*"
          on:change={handleProfileChange}
          class="input"
        />
        {#if profilePreview}
          <img src={profilePreview} alt="" class="preview" />
        {/if}
        {#if profileError}<div class="error">{profileError}</div>{/if}
      </div>

      <div class="form-group">
        <h2 class="label">Banner Picture</h2>
        <input
          type="file"
          accept="image/*"
          on:change={handleBannerChange}
          class="input"
        />
        {#if bannerPreview}
          <img src={bannerPreview} alt="" class="preview" />
        {/if}
        {#if bannerError}<div class="error">{bannerError}</div>{/if}
      </div>

      <div class="form-group captcha-container">
        <div
          class="g-recaptcha"
          data-sitekey="6LddEycrAAAAAC-N3vVTDhSQRVJbo8zfWpDSuFEv"
          data-theme="dark"
          data-callback="onCaptchaSuccess"
        ></div>
      </div>

      <div class="newsletter-container">
        <label class="checkbox-container">
          <input type="checkbox" bind:checked={subscribeNewsletter} />
          <span class="checkbox-text">Subscribe to our newsletter to get updates</span>
        </label>
      </div>

      <div class="button-container">
        <button
          type="button"
          class="prev-button"
          disabled={true}
          on:click={prevStep}
        >
          Prev
        </button>
        <button
          type="button"
          class="next-button"
          on:click={nextStep}
          disabled={!captchaVerified}
        >
          Next
        </button>
      </div>

    </div>
  {:else if step === 2}
    <div class="security-form">
      <h1 class="main-title">Security Setup</h1>
      <p class="description">
        Set up security questions to help protect your account.
      </p>

      <div class="security-questions-container">
        {#each securityQAs as qa, idx (idx)}
          <div class="security-question-card">
            <div class="form-group">
              <label for={`question-${idx}`} class="label">Security Question {idx + 1}</label>
              <select
                id={`question-${idx}`}
                class="select"
                value={qa.question}
                on:change={(e) => handleQAChange(idx, "question", (e.target as HTMLSelectElement).value)}
              >
                <option value="" disabled selected>Select a question</option>
                {#each securityQuestions as question (question)}
                  <option value={question}>{question}</option>
                {/each}
              </select>
            </div>

            <div class="form-group">
              <label for={`answer-${idx}`} class="label">Your Answer</label>
              <input
                id={`answer-${idx}`}
                type="text"
                class="input-answer"
                placeholder="Answer"
                value={qa.answer}
                on:input={(e) => handleQAChange(idx, "answer", (e.target as HTMLInputElement).value)}
              />
            </div>

            {#if idx > 0}
              <button
                type="button"
                class="remove-question-button"
                on:click={() => removeSecurity(idx)}
                aria-label="Remove security question"
              >
                <span class="x-icon">×</span>
              </button>
            {/if}
          </div>
        {/each}
        {#if securityQAs.length < 3}
          <button type="button" class="add-question-button" on:click={addSecurity}>
            + Add another security question
          </button>
        {/if}
      </div>
      <div class="button-container">
        <button
          type="button"
          class="prev-button"
          disabled={false}
          on:click={prevStep}
        >
          Prev
        </button>
        <button
          type="button"
          class="next-button"
          on:click={submitSecurityAnswers}
          disabled={securityQAs.some(qa => !qa.question || !qa.answer)}
        >
          Next
        </button>
      </div>
    </div>
  {:else if step === 3}
    <div class="otp-form">
      <h1 class="main-title">Enter OTP</h1>
      <p class="description">
        We sent a one-time code to <strong>{email}</strong>. Enter it below:
      </p>
      <div class="otp-inputs">
        {#each otpDigits.map((_, i) => i) as i (i)}
          <input
            id={"otp-" + i}
            type="text"
            inputmode="numeric"
            maxlength="1"
            bind:value={otpDigits[i]}
            on:input={(e) => handleOtpInput(e, i)}
            on:keydown={(e) => handleOtpKeyDown(e, i)}
            class="otp-box"
          />
        {/each}
      </div>
      <div class="button-container">
        <button
          type="button"
          class="prev-button"
          disabled={false}
          on:click={prevStep}
        >
          Prev
        </button>
        <button
          type="button"
          class="next-button"
          on:click={submitOtp}
          disabled={otp.length < 6}
        >
          Verify OTP
        </button>
      </div>

      <p class="resend">
        Didn"t receive it? <a href="#" on:click|preventDefault={resendOtp}>Resend code</a>
      </p>
    </div>
  {/if}
</div>

<style>
  .main-container {
    max-width: 600px;
    margin: 0 auto;
    padding: 2rem 1rem;
    overflow-y: auto;
    height: 90vh;
  }
  .button-container{
    display: flex;
    flex-direction: row;
    gap: 1vw;
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
  }

  .otp-form {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }

  .security-form {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }

  .main-title {
    font-size: 1.8rem;
    font-weight: 700;
    text-align: center;
  }

  .description {
    font-size: 0.95rem;
    color: var(--secondary-color);
    text-align: center;
    margin-bottom: 1rem;
  }

  .form-group .description {
    font-size: 0.85rem;
    color: var(--secondary-color);
    margin-top: -0.5rem;
    margin-bottom: 0.5rem;
    text-align: left;
  }

  .otp-form .description {
    text-align: center;
    color: var(--secondary-color);
  }

  .form-group {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    margin-bottom: 1rem;
  }

  .captcha-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    width: 100%;
    margin: 1rem 0;
  }

  .label {
    font-weight: bold;
    font-size: 1rem;
    width: 80%;
  }

  .hinted-field {
    position: relative;
  }

  .input,
  .otp-box,
  .select {
    width: 100%;
    padding: 1rem;
    font-size: 1rem;
    border: 1px solid var(--border-color);
    border-radius: 0.5rem;
    background-color: #f7f9fa;
    color: black;
  }

  .input-answer{
    width: 93%;
    padding: 1rem;
    font-size: 1rem;
    border: 1px solid var(--border-color);
    border-radius: 0.5rem;
    background-color: #f7f9fa;
  }

  .input-counter {
    position: absolute;
    right: 1rem;
    top: 50%;
    transform: translateY(-50%);
    font-size: 0.75rem;
    color: var(--secondary-color);
  }

  .dob-selects {
    display: flex;
    gap: 0.5rem;
  }

  .otp-inputs {
    display: flex;
    justify-content: center;
    gap: 0.5rem;
    margin: 1rem 0;
  }

  .otp-box {
    width: 3rem;
    height: 3rem;
    text-align: center;
    font-size: 1.5rem;
    font-weight: bold;
  }

  .next-button {
    padding: 1rem;
    font-size: 1rem;
    width: 50%;
    font-weight: bold;
    border: none;
    border-radius: 9999px;
    background-color: var(--primary-color);
    color: var(--bg-color);
    cursor: pointer;
    transition: background-color 0.2s;
  }

  .prev-button {
    padding: 1rem;
    font-size: 1rem;
    width: 50%;
    font-weight: bold;
    border: none;
    border-radius: 9999px;
    background-color: var(--primary-color);
    color: var(--bg-color);
    cursor: pointer;
    transition: background-color 0.2s;
  }

  .next-button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
  .prev-button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .resend {
    text-align: center;
    font-size: 0.9rem;
  }

  .resend a {
    color: var(--primary-color);
    text-decoration: none;
    font-weight: bold;
    cursor: pointer;
  }

  .preview {
    margin-top: 0.5rem;
    width: 100px;
    border-radius: 0.5rem;
  }

  .security-questions-container {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    margin-bottom: 1.5rem;
  }

  .security-question-card {
    position: relative;
    border-radius: 0.5rem;
    padding: 1.25rem;
    border: 1px solid var(--border-color);
    margin-bottom: 0.5rem;
  }

  .remove-question-button {
    position: absolute;
    top: 0.75rem;
    right: 0.75rem;
    width: 1.75rem;
    height: 1.75rem;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
    border: none;
    background-color: #ff3b30;
    color: white;
    font-weight: bold;
    cursor: pointer;
    padding: 0;
    transition: background-color 0.2s;
  }

  .remove-question-button:hover {
    background-color: #e02d26;
  }

  .x-icon {
    font-size: 1.25rem;
    line-height: 1;
  }

  .add-question-button {
    border: 1px dashed var(--border-color);
    background-color: transparent;
    color: var(--primary-color);
    font-weight: 600;
    padding: 0.75rem;
    border-radius: 0.5rem;
    cursor: pointer;
    text-align: center;
    margin-bottom: 0.5rem;
  }

  .newsletter-container {
    color: red;
    border: 3px solid var(--border-color);
    border-radius: 0.5rem;
    padding: 1rem;
    margin-bottom: 1.5rem;
  }

  .checkbox-container {
    display: flex;
    align-items: center;
    cursor: pointer;
  }
  .error{
    color: red;
  }
  .checkbox-container input[type="checkbox"] {
    margin-right: 0.5rem;
    width: 1.1rem;
    height: 1.1rem;
  }

  .checkbox-text {
    font-size: 0.95rem;
    color: var(--text-color);
  }

  @media (max-width: 700px) {
    .main-container {
      max-width: 100vw;
      min-width: 0;
      margin-right: 4vw;
      margin:3vw;
      /* margin-top: -5vh; */
      max-height: 90vh;
      padding: 1rem 0.3rem;
      height: auto;
      overflow-x: hidden;
    }
    .logo {
      width: 48px;
    }
    .main-title {
      font-size: 1.2rem;
      padding: 0 0.4rem;
    }
    .form, .otp-form, .security-form {
      gap: 0.7rem;
      padding: 0 0.2rem;
    }
    .button-container {
      flex-direction: column;
      gap: 0.7rem;
    }
    .prev-button, .next-button {
      width: 100%;
      padding: 0.7rem;
      font-size: 1rem;
    }
    .newsletter-container {
      padding: 0.7rem;
      margin-bottom: 1rem;
      font-size: 0.96rem;
    }
    .dob-selects {
      flex-direction: column;
      gap: 0.3rem;
    }
    .select, .input, .input-answer {
      font-size: 1rem;
      padding: 0.7rem;
    }
    .preview {
      width: 70px;
    }
    .security-question-card {
      padding: 0.9rem;
    }
    .otp-inputs {
      gap: 0.3rem;
    }
    .otp-box {
      width: 2.3rem;
      height: 2.3rem;
      font-size: 1.1rem;
    }
  }

  @media (max-width: 500px) {
    .main-title { font-size: 1rem; }
    .logo { width: 36px; }
    .preview { width: 48px; }
    .newsletter-container { font-size: 0.85rem; }
    .button-container { gap: 0.35rem; }
    .otp-inputs {
      gap: 0.13rem;
    }
    .otp-box {
      width: 3.65rem;
      height: 3.65rem;
      font-size: 1rem;
      padding: 0;
    }
  }

</style>