import { writable, type Writable } from "svelte/store";

export interface ToastItem {
  id: string;
  message: string;
  duration?: number;
}

// single source of truth for all toasts
export const toasts: Writable<ToastItem[]> = writable([]);

/**
 * Show a toast message.
 * @param message the text to display
 * @param duration how long (ms) before auto-dismissal
 */
export function toast(message: string, duration = 3000): void {
  const id = crypto.randomUUID();
  toasts.update((all) => [...all, { id, message, duration }]);
}
