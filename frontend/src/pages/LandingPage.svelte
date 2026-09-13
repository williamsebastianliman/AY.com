<script>
  import logo from "../assets/logo_ay.jpg";
  import { navigate } from "svelte-routing";
  function goToRegister(){
    navigate("/register");
  }
  function goToLogin(){
    navigate("/login");
  }
</script>

<main>
  <div class="main-layout">
    <div class="left-section">
      <img src={logo} alt="Logo" />
    </div>

    <div class="right-section">
      <div class="title">Happening now</div>
      <div class="tagline">Join today.</div>

      <div class="signup-section">
        <div class="signup-upper">
          <button type="button" class="temp-button">Sign in using Google</button>
          <button type="button" class="temp-button">Sign up with Apple</button>
        </div>

        <div class="signup-middle">
          <div class="line"></div>
          <div class="or-text">OR</div>
          <div class="line"></div>
        </div>

        <div class="signup-bottom">
          <button type="button" class="create-account-button" onclick={goToRegister}>
            Create account
          </button>
        </div>

        <div class="login-section">
          <div class="sub-text">Already have an account?</div>
          <button type="button" class="secondary" onclick={goToLogin}>Sign in</button>
        </div>
      </div>
    </div>
  </div>
</main>

<style>
  .main-layout {
    font-size: 10px;
    display: flex;
    flex-direction: row;
    background-color: black;
    align-items: center;
    /* height: 100vh; */
    min-height: 100vh;
  }

  .left-section {
    margin-left: 5vw;
  }

  img {
    width: 700px;
    height: 700px;
  }

  .right-section {
    margin-left: 10vw;
    display: flex;
    flex-direction: column;
    gap: 2vh;
  }

  .signup-section {
    display: flex;
    flex-direction: column;
    gap: 20px;
  }

  .signup-upper {
    display: flex;
    flex-direction: column;
    gap: 20px;
  }

  .temp-button {
    width: 50%;
  }

  .title {
    font-size: 80px;
    font-weight: bold;
  }

  .tagline {
    font-size: 40px;
    font-weight: bold;
  }

  .line {
    width: 6vw;
    height: 1px;
    background-color: gray;
  }

  .signup-middle {
    display: flex;
    flex-direction: row;
    align-items: center;
    gap: 2vw;
  }

  .create-account-button {
    width: 50%;
  }

  .or-text {
    font-size: 20px;
  }

  .sub-text {
    font-size: 15px;
    font-weight: bold;
  }

  .login-section {
    display: flex;
    flex-direction: column;
    margin-top: 20px;
    gap: 20px;
    width: 50%;
  }

  @media (max-width: 1400px) {
  .main-layout {
    flex-direction: column;
    align-items: center;
    justify-content: flex-start;
    height: auto;
    max-height: 10vh;
    /* padding: 32px 0 48px 0; */
    overflow-y: auto;
  }

  .left-section {
    margin-left: 0;
    margin-bottom: 24px;
    display: flex;
    justify-content: center;
    width: 100%;
  }

  img {
    width: 60vw;
    max-width: 320px;
    height: auto;
    display: block;
    margin: 0 auto;
  }

  .right-section {
    margin-left: 0;
    align-items: center;
    text-align: center;
    width: 100%;
  }

  .title {
    font-size: 42px;
  }

  .tagline {
    font-size: 26px;
  }

  .signup-section,
  .signup-upper,
  .login-section {
    width: 100%;
    align-items: center;
  }

  .temp-button,
  .create-account-button {
    width: 90%;
    max-width: 400px;
    font-size: 1rem;
  }

  .login-section {
    width: 90%;
    margin: 24px auto 0 auto;
    gap: 12px;
  }

  .signup-middle {
    gap: 12px;
  }

  .line {
    width: 15vw;
    max-width: 60px;
  }
  .or-text {
    font-size: 18px;
  }
}

@media (max-width: 500px) {
  .title {
    font-size: 2.2rem;
  }
  .tagline {
    font-size: 1.2rem;
  }
  .line {
    width: 18vw;
    max-width: 36px;
  }
}

</style>