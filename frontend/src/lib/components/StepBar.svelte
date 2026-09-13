<script lang="ts">
  export let currentStep = 1;
  const steps = [1, 2, 3];
</script>

<div class="main-container">
  <div class="line"></div>
  <div class="circle-container">
    {#each steps as step, i (i)}
      <div class="circle" class:active={step === currentStep}>
        {step}
      </div>
    {/each}
  </div>
</div>

<style>
  .main-container {
    position: relative;
    width: 30vw;
    margin: 30px auto;
  }

  .line {
    background-color: var(--primary-color);
    width: 100%;
    height: 8px;
    border-radius: 20px;
    position: relative;
  }

  .circle {
    position: absolute;
    top: 50%;
    transform: translate(-50%, -50%);
    width: 50px;
    height: 50px;
    border-radius: 50%;
    background-color: var(--secondary-color);
    border: 3px solid var(--primary-color);
    color: var(--text-color);
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: bold;
    font-size: 1.2rem;
    transition: background-color 0.2s, color 0.2s;
    left: 0%;
  }

  .circle-container .circle:nth-of-type(1) {
    left: 0%;
  }

  .circle-container .circle:nth-of-type(2) {
    left: 50%;
  }

  .circle-container .circle:nth-of-type(3) {
    left: 100%;
  }

  .circle.active {
    background-color: var(--primary-color);
    color: var(--bg-color);
  }
</style>