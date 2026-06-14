import { ref } from '../vue.js';

// Module-level state — shared across all components that call useToast().
const toasts = ref([]);

export function useToast() {
  const show = (message, type = 'success') => {
    const id = Date.now() + Math.random();
    toasts.value.push({ id, message, type });
    setTimeout(() => {
      toasts.value = toasts.value.filter(t => t.id !== id);
    }, 3500);
  };

  return { toasts, show };
}
