<script lang="ts">
  import logo from "../../assets/logo_ay.jpg";
  export let onClose: () => void;
  let name = "";
  let email = "";
  let month = "";
  let day = "";
  let year = "";
</script>

<div class="register-container">
  <button type="button" class="close-button" on:click={onClose}>
    ×
  </button>

  <img src={logo} alt="Logo" class="logo" />

  <h1 class="title">Create your account</h1>

  <div class="form-group">
    <input
      type="text"
      bind:value={name}
      placeholder="Name"
      maxlength={50}
      class="input"
    />
    <div class="input-counter">{name.length} / 50</div>
  </div>

  <div class="form-group">
    <input
      type="email"
      bind:value={email}
      placeholder="Email"
      class="input"
    />
  </div>

  <div class="dob-section">
    <h2 class="dob-title">Date of birth</h2>
    <p class="dob-description">
      This will not be shown publicly. Confirm your own age, even if this
      account is for a business, a pet, or something else.
    </p>

    <div class="dob-selects">
      <select bind:value={month} class="select">
        <option disabled value="">Month</option>
        <option value="January">January</option>
        <option value="February">February</option>
        <option value="March">March</option>
        <option value="April">April</option>
        <option value="May">May</option>
        <option value="June">June</option>
        <option value="July">July</option>
        <option value="August">August</option>
        <option value="September">September</option>
        <option value="October">October</option>
        <option value="November">November</option>
        <option value="December">December</option>
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
  </div>

  <button type="button" class="next-button">Next</button>
</div>

<style>
.register-container {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 2rem 1rem;
  max-width: 400px;
  margin: 0 auto;
}

.close-button {
  position: absolute;
  top: 1rem;
  left: 1rem;
  background: transparent;
  border: none;
  font-size: 1.5rem;
  color: var(--text-color);
  cursor: pointer;
}

.logo {
  width: 48px;
  height: 48px;
  margin-bottom: 1rem;
}

.title {
  font-size: 1.8rem;
  font-weight: 700;
  margin-bottom: 2rem;
}

.form-group {
  width: 100%;
  position: relative;
  margin-bottom: 1rem;
}

.input {
  width: 100%;
  padding: 1rem;
  background-color: var(--bg-color);
  border: 1px solid var(--border-color);
  border-radius: 0.5rem;
  color: var(--text-color);
  font-size: 1rem;
}

.input-counter {
  position: absolute;
  top: 50%;
  right: 1rem;
  transform: translateY(-50%);
  font-size: 0.8rem;
  color: var(--secondary-color);
}

.dob-section {
  width: 100%;
  margin: 2rem 0;
}

.dob-title {
  font-weight: 700;
  margin-bottom: 0.5rem;
}

.dob-description {
  font-size: 0.9rem;
  color: var(--secondary-color);
  margin-bottom: 1rem;
}

.dob-selects {
  display: flex;
  gap: 0.5rem;
}

.select {
  flex: 1;
  padding: 1rem;
  background-color: var(--bg-color);
  border: 1px solid var(--border-color);
  border-radius: 0.5rem;
  color: var(--text-color);
  font-size: 1rem;
}

.next-button {
  width: 100%;
  background-color: var(--primary-color);
  color: var(--bg-color);
  border: none;
  border-radius: 9999px;
  padding: 0.8rem;
  font-weight: bold;
  font-size: 1rem;
  cursor: pointer;
  margin-top: 2rem;
}
</style>