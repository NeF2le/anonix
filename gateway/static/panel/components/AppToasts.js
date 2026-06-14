import { useToast } from '../composables/useToast.js';

export default {
  setup() {
    const { toasts } = useToast();
    return { toasts };
  },
  template: `
    <div class="fixed bottom-4 right-4 z-50 flex flex-col gap-2 pointer-events-none">
      <div
        v-for="toast in toasts" :key="toast.id"
        :class="toast.type === 'error' ? 'bg-red-600' : 'bg-emerald-600'"
        class="text-white px-4 py-3 rounded-xl shadow-lg text-sm max-w-xs flex items-start gap-2 pointer-events-auto"
      >
        <svg v-if="toast.type === 'error'" class="w-4 h-4 flex-shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
        </svg>
        <svg v-else class="w-4 h-4 flex-shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
        </svg>
        {{ toast.message }}
      </div>
    </div>
  `,
};
