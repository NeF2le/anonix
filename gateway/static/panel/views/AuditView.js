import { ref, onMounted } from '../vue.js';
import { api } from '../api.js';
import { usePagination } from '../composables/usePagination.js';
import AppPagination from '../components/AppPagination.js';

const formatDate = (str) => {
  if (!str) return '—';
  try {
    return new Date(str).toLocaleString('ru-RU', {
      day: '2-digit', month: '2-digit', year: 'numeric',
      hour: '2-digit', minute: '2-digit', second: '2-digit',
    });
  } catch {
    return str;
  }
};

const ACTION_LABELS = {
  tokenize:   'Токенизация',
  detokenize: 'Детокенизация',
};

export default {
  setup() {
    const entries = ref([]);
    const loading = ref(false);
    const error   = ref('');

    const { page, totalPages, pageItems, setPage } = usePagination(entries);

    const loadEntries = async () => {
      loading.value = true;
      error.value   = '';
      try {
        const data    = await api.getAuditLog();
        entries.value = Array.isArray(data) ? data : [];
      } catch (e) {
        error.value = e.message;
      } finally {
        loading.value = false;
      }
    };

    onMounted(loadEntries);

    return {
      entries, loading, error,
      page, totalPages, pageItems, setPage,
      loadEntries, formatDate,
      actionLabel: (action) => ACTION_LABELS[action] || action,
    };
  },
  components: { AppPagination },
  template: `
    <div class="max-w-6xl mx-auto">
      <div class="flex items-center justify-between mb-5">
        <h2 class="text-lg font-bold text-slate-900">Аудит</h2>
        <button @click="loadEntries"
          class="inline-flex items-center gap-1.5 text-sm text-slate-500 hover:text-slate-800 border border-slate-300 hover:border-slate-400 px-3 py-1.5 rounded-lg transition">
          <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
          </svg>
          Обновить
        </button>
      </div>

      <div v-if="loading" class="flex items-center justify-center py-16 text-slate-400">
        <svg class="animate-spin w-6 h-6 mr-2" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"/>
        </svg>
        Загрузка...
      </div>
      <div v-else-if="error" class="p-4 bg-red-50 border border-red-200 text-red-700 rounded-xl text-sm">{{ error }}</div>
      <div v-else class="bg-white rounded-xl border border-slate-200 overflow-hidden shadow-sm">
        <div class="overflow-y-auto max-h-[60vh]">
          <table class="w-full text-sm">
            <thead class="bg-slate-50 border-b border-slate-200 sticky top-0 z-10">
              <tr>
                <th class="text-left px-4 py-3 font-semibold text-slate-600 text-xs uppercase tracking-wide">Время</th>
                <th class="text-left px-4 py-3 font-semibold text-slate-600 text-xs uppercase tracking-wide">Действие</th>
                <th class="text-left px-4 py-3 font-semibold text-slate-600 text-xs uppercase tracking-wide">Токен</th>
                <th class="text-left px-4 py-3 font-semibold text-slate-600 text-xs uppercase tracking-wide">Вид данных</th>
                <th class="text-left px-4 py-3 font-semibold text-slate-600 text-xs uppercase tracking-wide">Пользователь</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100">
              <tr v-for="entry in pageItems" :key="entry.id" class="hover:bg-slate-50 transition-colors">
                <td class="px-4 py-3 text-slate-600">{{ formatDate(entry.created_at) }}</td>
                <td class="px-4 py-3 text-slate-700">{{ actionLabel(entry.action) }}</td>
                <td class="px-4 py-3 font-mono text-xs text-slate-400">{{ entry.token }}</td>
                <td class="px-4 py-3 text-slate-700">{{ entry.kind ? entry.kind.russian_name : '—' }}</td>
                <td class="px-4 py-3 font-mono text-xs text-slate-400" :title="entry.user_id">{{ entry.user_id.substring(0,8) }}…</td>
              </tr>
              <tr v-if="entries.length === 0">
                <td colspan="5" class="text-center py-10 text-slate-400 text-sm">Записи не найдены</td>
              </tr>
            </tbody>
          </table>
        </div>
        <AppPagination :page="page" :total-pages="totalPages" @update:page="setPage" />
      </div>
    </div>
  `,
};
