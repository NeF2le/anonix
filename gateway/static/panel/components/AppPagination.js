export default {
  props: {
    page:       { type: Number, required: true },
    totalPages: { type: Number, required: true },
  },
  emits: ['update:page'],
  template: `
    <div v-if="totalPages > 1" class="flex items-center justify-between px-4 py-3 border-t border-slate-200 bg-slate-50">
      <button @click="$emit('update:page', page - 1)" :disabled="page <= 1"
        class="px-3 py-1.5 rounded-lg border border-slate-300 text-slate-600 text-xs font-medium hover:bg-white disabled:opacity-40 disabled:cursor-not-allowed transition">
        ← Назад
      </button>
      <span class="text-xs text-slate-500">Страница {{ page }} из {{ totalPages }}</span>
      <button @click="$emit('update:page', page + 1)" :disabled="page >= totalPages"
        class="px-3 py-1.5 rounded-lg border border-slate-300 text-slate-600 text-xs font-medium hover:bg-white disabled:opacity-40 disabled:cursor-not-allowed transition">
        Вперёд →
      </button>
    </div>
  `,
};
