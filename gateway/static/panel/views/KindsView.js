import { ref, onMounted } from '../vue.js';
import { api } from '../api.js';
import { useToast } from '../composables/useToast.js';
import { useModal } from '../composables/useModal.js';
import { usePagination } from '../composables/usePagination.js';
import AppPagination from '../components/AppPagination.js';

const KIND_FIELDS = [
  { key: 'name',         label: 'Название (EN)',   type: 'text',   required: true, placeholder: 'passport' },
  { key: 'russian_name', label: 'Название (RU)',   type: 'text',   required: true, placeholder: 'Паспорт'  },
  { key: 'short_name',   label: 'Короткий префикс', type: 'text',  required: false, placeholder: 'psp' },
  { key: 'access_level', label: 'Уровень доступа', type: 'number', required: true, placeholder: '3'        },
  { key: 'mask',         label: 'Маска (regex)',   type: 'text',   required: false, placeholder: '^\\d{4} \\d{6}$' },
];

export default {
  setup() {
    const { show: toast }        = useToast();
    const { modal, open, close } = useModal();

    const kinds   = ref([]);
    const loading = ref(false);
    const error   = ref('');

    const { page, totalPages, pageItems, setPage } = usePagination(kinds);

    const loadKinds = async () => {
      loading.value = true;
      error.value   = '';
      try {
        const data  = await api.getKinds();
        kinds.value = Array.isArray(data) ? data : [];
      } catch (e) {
        error.value = e.message;
      } finally {
        loading.value = false;
      }
    };

    // Reads current form values and builds the API payload.
    const payload = () => ({
      name:         modal.values.name,
      russian_name: modal.values.russian_name,
      short_name:   modal.values.short_name,
      access_level: parseInt(modal.values.access_level),
      mask:         modal.values.mask,
    });

    const openCreateModal = () => {
      open({
        type:   'form',
        title:  'Добавить вид токена',
        fields: KIND_FIELDS,
        values: { name: '', russian_name: '', short_name: '', access_level: '', mask: '' },
        onConfirm: async () => {
          modal.loading = true;
          modal.error   = '';
          try {
            await api.createKind(payload());
            close();
            toast('Вид токена создан');
            await loadKinds();
          } catch (e) {
            modal.error = e.message;
          } finally {
            modal.loading = false;
          }
        },
      });
    };

    const openEditModal = (kind) => {
      open({
        type:   'form',
        title:  'Изменить вид токена',
        fields: KIND_FIELDS,
        values: { name: kind.name, russian_name: kind.russian_name, short_name: kind.short_name || '', access_level: kind.access_level, mask: kind.mask || '' },
        onConfirm: async () => {
          modal.loading = true;
          modal.error   = '';
          try {
            await api.updateKind(kind.id, payload());
            close();
            toast('Вид токена обновлён');
            await loadKinds();
          } catch (e) {
            modal.error = e.message;
          } finally {
            modal.loading = false;
          }
        },
      });
    };

    const openDeleteModal = (kind) => {
      open({
        type:    'confirm',
        title:   'Удалить вид токена',
        message: `Удалить «${kind.russian_name}» (${kind.name})?`,
        onConfirm: async () => {
          modal.loading = true;
          try {
            await api.deleteKind(kind.id);
            close();
            toast('Вид токена удалён');
            await loadKinds();
          } catch (e) {
            modal.error = e.message;
          } finally {
            modal.loading = false;
          }
        },
      });
    };

    onMounted(loadKinds);

    return {
      kinds, loading, error,
      page, totalPages, pageItems, setPage,
      openCreateModal, openEditModal, openDeleteModal,
    };
  },
  components: { AppPagination },
  template: `
    <div class="max-w-5xl mx-auto">
      <div class="flex items-center justify-between mb-5">
        <h2 class="text-lg font-bold text-slate-900">Виды токенов</h2>
        <button @click="openCreateModal"
          class="inline-flex items-center gap-1.5 bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded-lg text-sm font-medium transition">
          <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
          </svg>
          Добавить
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
                <th class="text-left px-4 py-3 font-semibold text-slate-600 text-xs uppercase tracking-wide">ID</th>
                <th class="text-left px-4 py-3 font-semibold text-slate-600 text-xs uppercase tracking-wide">Название (EN)</th>
                <th class="text-left px-4 py-3 font-semibold text-slate-600 text-xs uppercase tracking-wide">Название (RU)</th>
                <th class="text-left px-4 py-3 font-semibold text-slate-600 text-xs uppercase tracking-wide">Префикс</th>
                <th class="text-left px-4 py-3 font-semibold text-slate-600 text-xs uppercase tracking-wide">Уровень доступа</th>
                <th class="text-left px-4 py-3 font-semibold text-slate-600 text-xs uppercase tracking-wide">Маска</th>
                <th class="text-right px-4 py-3 font-semibold text-slate-600 text-xs uppercase tracking-wide">Действия</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100">
              <tr v-for="kind in pageItems" :key="kind.id" class="hover:bg-slate-50 transition-colors">
                <td class="px-4 py-3 text-slate-500">{{ kind.id }}</td>
                <td class="px-4 py-3 font-mono text-xs text-slate-700">{{ kind.name }}</td>
                <td class="px-4 py-3 font-medium text-slate-900">{{ kind.russian_name }}</td>
                <td class="px-4 py-3 font-mono text-xs text-slate-500">{{ kind.short_name || '—' }}</td>
                <td class="px-4 py-3">
                  <span class="inline-block bg-amber-100 text-amber-700 text-xs rounded-full px-2.5 py-0.5 font-medium">
                    {{ kind.access_level }}
                  </span>
                </td>
                <td class="px-4 py-3 font-mono text-xs text-slate-500">{{ kind.mask || '—' }}</td>
                <td class="px-4 py-3 text-right">
                  <button @click="openEditModal(kind)"
                    class="text-indigo-600 hover:text-indigo-800 text-xs font-medium transition mr-3">Изменить</button>
                  <button @click="openDeleteModal(kind)"
                    class="text-red-500 hover:text-red-700 text-xs font-medium transition">Удалить</button>
                </td>
              </tr>
              <tr v-if="kinds.length === 0">
                <td colspan="7" class="text-center py-10 text-slate-400 text-sm">Виды токенов не найдены</td>
              </tr>
            </tbody>
          </table>
        </div>
        <AppPagination :page="page" :total-pages="totalPages" @update:page="setPage" />
      </div>
    </div>
  `,
};
