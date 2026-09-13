// src/global.d.ts
declare global {
  namespace svelte.JSX {
    interface HTMLAttributes {
      "on:outclick"?: () => void;
    }
  }
}

export {};
