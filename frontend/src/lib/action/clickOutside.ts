export function clickOutside(node: HTMLElement) {
  const handle = (event: MouseEvent) => {
    if (!node.contains(event.target as Node)) {
      node.dispatchEvent(new CustomEvent("outclick"));
    }
  };
  document.addEventListener("click", handle, true);
  return {
    destroy() {
      document.removeEventListener("click", handle, true);
    },
  };
}
