<script lang="ts">
    export let value: number = 0;
    export let max: number = 100;
    export let size: number = 48;
    export let strokeWidth: number = 6;
    export let color: string = "#1d9bf0";
    $: radius = (size - strokeWidth) / 2;
    $: circumference = 2 * Math.PI * radius;
    $: progress = Math.min(Math.max(value / max, 0), 1);
    $: offset = circumference * (1 - progress);
  </script>
  <svg width={size} height={size} class="donut">
    <circle
      class="track"
      cx={size / 2}
      cy={size / 2}
      r={radius}
      stroke-width={strokeWidth}
      fill="none"
    />
    <circle
      class="indicator"
      cx={size / 2}
      cy={size / 2}
      r={radius}
      stroke={color}
      stroke-width={strokeWidth}
      fill="none"
      stroke-dasharray={circumference}
      stroke-dashoffset={offset}
      stroke-linecap="round"
    />
  </svg>
  <style>
    .donut {
      transform: rotate(-90deg);
    }
    .track {
      stroke: #e5e7eb;
    }
    .indicator {
      transition: stroke-dashoffset 0.35s;
      transform: rotate(0deg);
      transform-origin: center;
    }
</style>