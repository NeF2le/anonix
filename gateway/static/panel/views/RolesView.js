import { ref, onMounted } from '../vue.js';
import { api } from '../api.js';
import { usePagination } from '../composables/usePagination.js';
import AppPagination from '../components/AppPagination.js';

export default {
  setup() {
    const roles   = ref([]);
    const loading = ref(false);
    const error   = ref('');

    const { page, totalPages, pageItems, setPage } = usePagination(roles);

    const loadRoles = async () => {
      loading.value = true;
      error.value   = '';
      try {
        const data  = await api.getRoles();
        roles.value = data?.roles || [];
      } catch (e) {
        error.value = e.message;
      } finally {
        loading.value = false;
      }
    };

    onMounted(loadRoles);

    return { roles, loading, error, page, totalPages, pageItems, setPage };
  },
  components: { AppPagination },
  template: `
    <div class="max-w-2xl mx-auto">
      <h2 class="text-lg font-bold text-slate-900 mb-5">Роли</h2>

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
                <th class="text-left px-4 py-3 font-semibold text-slate-600 text-xs uppercase tracking-wide">ID</th>
                <th class="text-left px-4 py-3 font-semibold text-slate-600 text-xs uppercase tracking-wide">Название</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100">
              <tr v-for="role in pageItems" :key="role.id" class="hover:bg-slate-50 transition-colors">
                <td class="px-4 py-3 text-slate-500">{{ role.id }}</td>
                <td class="px-4 py-3">
                  <span class="inline-block bg-indigo-100 text-indigo-700 text-xs rounded-full px-2.5 py-0.5 font-medium">
                    {{ role.name }}
                  </span>
                </td>
              </tr>
              <tr v-if="roles.length === 0">
                <td colspan="2" class="text-center py-10 text-slate-400 text-sm">Роли не найдены</td>
              </tr>
            </tbody>
          </table>
        </div>
        <AppPagination :page="page" :total-pages="totalPages" @update:page="setPage" />
      </div>
    </div>
  `,
};
