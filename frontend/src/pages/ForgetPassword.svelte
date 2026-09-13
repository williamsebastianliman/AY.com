<script lang="ts">

    interface SecurityQuestions{
      id : string,
      question: string,
    }
    import logo from "../assets/logo_ay.jpg";
    import api from "../lib/api";
    let email = "";
    let emailError = "";
    let step = 1;
    let securityQuestions : SecurityQuestions = [];
    let selectedQuestion = "";
    let answer = "";
    let newPassword = "";
    let confirmPassword = "";
    let passwordError = "";
    let resetSuccess = false;
    let myUserId = "";
    async function fetchUser(email: string) {
      try{
        const data = await api.post("/user/get-by-email",{email: email});
        myUserId = data.data.id;
        fetchQuestions();
        step = 2;
      }catch(err){
        console.error(err as string);
        emailError = err;
      }
    }
    async function fetchQuestions() {
      try{
        const res = await api.post("/security-answer", {user_id: myUserId});
        securityQuestions = res.data.answers;
      }catch(err){
        console.error(err as string);
      }
    }
    async function handleEmailNext() {
      emailError = "";
      if (!email.includes("@")) {
        emailError = "Enter a valid email.";
        return;
      }
      await fetchUser(email);
    }
    function validatePassword(pw) {
      const rules = [
        { regex: /.{8,}/, message: "Password must be at least 8 characters long" },
        { regex: /[A-Z]/, message: "Password must contain at least one uppercase letter" },
        { regex: /[a-z]/, message: "Password must contain at least one lowercase letter" },
        { regex: /\d/, message: "Password must contain at least one digit" },
        { regex: /\W/, message: "Password must contain at least one special character" }
      ];
      for (const rule of rules) {
        if (!rule.regex.test(pw)) return rule.message;
      }
      return "";
    }
    async function resetPassword(){
      if (newPassword != confirmPassword){
        passwordError = "Password and Confirm Password Must Be The Same Value!";
        return;
      }
      const pwValidationMsg = validatePassword(newPassword);
      if (pwValidationMsg) {
        passwordError = pwValidationMsg;
        return;
      }
      try{
        const data = await api.post("/auth/check-security-answer", {user_id: myUserId, question: selectedQuestion, answer: answer, new_password: newPassword});
        if (data.data.success == false){
          if(data.data.error_message=="failed to change password: rpc error: code = Internal desc = failed to change password: new password cannot be the same as the old password")
          {
           passwordError = "New Password Cannot Be The Same As The New One!";
          }
          else{
            passwordError = data.data.error_message;
          }
        }
        else{
          passwordError = "";
          resetSuccess = true;
        }
      }catch(err){
        passwordError = err;
      }
    }
  </script>
  <div class="main-container">
    <img src={logo} alt="Logo" class="logo" />
    <div class="form">
      <h1 class="main-title">Forgotten Account</h1>
      {#if step === 1}
        <div class="form-group">
          <h2 class="label">Registered Email</h2>
          <input
            type="email"
            class="input"
            placeholder="Enter your email"
            bind:value={email}
          />
        </div>
        {#if emailError}
          <div class="login-error">{emailError}</div>
        {/if}
        <button
          type="button"
          class="next-button"
          on:click={handleEmailNext}
          disabled={!email}
        >
          Next
        </button>
      {:else if step === 2}
        <div class="form-group">
          <h2 class="label">Security Question</h2>
          <select class="input" bind:value={selectedQuestion}>
            <option value="" disabled selected>Select a question</option>
            {#each securityQuestions as q (q.id)}
              <option value={q.question}>{q.question}</option>
            {/each}
          </select>
        </div>
        <div class="form-group">
          <h2 class="label">Your Answer</h2>
          <input
            type="text"
            class="input"
            placeholder="Answer"
            bind:value={answer}
          />
        </div>
        <div class="form-group">
          <h2 class="label">New Password</h2>
          <input
            type="password"
            class="input"
            placeholder="New Password"
            bind:value={newPassword}
          />
        </div>
        <div class="form-group">
          <h2 class="label">Confirm New Password</h2>
          <input
            type="password"
            class="input"
            placeholder="Confirm Password"
            bind:value={confirmPassword}
          />
        </div>
        {#if passwordError}
          <div class="login-error">{passwordError}</div>
        {/if}
        <button
          type="button"
          class="next-button"
          on:click={resetPassword}
          disabled={!selectedQuestion || !answer || !newPassword || !confirmPassword}
        >
          Reset Password
        </button>
        {#if resetSuccess}
          <div class="reset-success">Your password has been reset. <a href="/">Go to Landing Page</a></div>
        {/if}
      {/if}
      <div class="links">
        <a href="/">Back to Landing Page</a>
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
    select.input {
      appearance: none;
      -webkit-appearance: none;
      background: var(--bg-color);
      color: var(--text-color);
      border-radius: 0.5rem;
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
    .reset-success {
      color: #43a047;
      background: #f0fff0;
      border: 1px solid #43a047;
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
        padding: 1rem 0.3rem;
        box-sizing: border-box;
        overflow-y: auto;
      }
      .form {
        padding: 1.2rem 0.6rem;
        box-sizing: border-box;
        min-width: 0;
      }
    }

  </style>