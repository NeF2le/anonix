import { ref, onMounted } from '../vue.js';
import { api } from '../api.js';
import { useToast } from '../composables/useToast.js';
import { useModal } from '../composables/useModal.js';
import { usePagination } from '../composables/usePagination.js';
import AppPagination from '../components/AppPagination.js';

export default {
  setup() {
    const { show: toast }         = useToast();
    const { modal, open, close }  = useModal();

    const users   = ref([]);
    const loading = ref(false);
    const error   = ref('');
    const roles   = ref([]);  // cached for role pickers

    const { page, totalPages, pageItems, setPage } = usePagination(users);

    // ── Data loaders ──────────────────────────────────────────────────────────

    const loadUsers = async () => {
      loading.value = true;
      error.value   = '';
      try {
        const data  = await api.getUsers();
        users.value = data?.users || [];
      } catch (e) {
        error.value = e.message;
      } finally {
        loading.value = false;
      }
    };

    const loadRoles = async () => {
      if (roles.value.length > 0) return;
      try {
        const data  = await api.getRoles();
        roles.value = data?.roles || [];
      } catch {}
    };

    // ── Modals ────────────────────────────────────────────────────────────────

    const openCreateModal = async () => {
      await loadRoles();
      open({
        type:   'form',
        title:  'Добавить пользователя',
        fields: [
          { key: 'login',    label: 'Логин',  type: 'text',     required: true,  placeholder: 'user123'  },
          { key: 'password', label: 'Пароль', type: 'password', required: true,  placeholder: '••••••••' },
          { key: 'role_id',  label: 'Роль',   type: 'select',
            options: roles.value.map(r => ({ value: r.id, label: r.name })) },
        ],
        values: { login: '', password: '', role_id: roles.value[0]?.id ?? 2 },
        onConfirm: async () => {
          modal.loading = true;
          modal.error   = '';
          try {
            await api.register(modal.values.login, modal.values.password, parseInt(modal.values.role_id) || 2);
            close();
            toast('Пользователь создан');
            await loadUsers();
          } catch (e) {
            modal.error = e.message;
          } finally {
            modal.loading = false;
          }
        },
      });
    };

    const openDeleteModal = (user) => {
      open({
        type:    'confirm',
        title:   'Удалить пользователя',
        message: `Удалить пользователя «${user.login}»? Действие необратимо.`,
        onConfirm: async () => {
          modal.loading = true;
          try {
            await api.deleteUser(user.id);
            close();
            toast('Пользователь удалён');
            await loadUsers();
          } catch (e) {
            modal.error = e.message;
          } finally {
            modal.loading = false;
          }
        },
      });
    };

    const openRolesModal = async (user) => {
      await loadRoles();
      open({
        type:         'roles',
        title:        'Управление ролями',
        user,
        userRoles:    [...(user.roles || [])],
        allRoles:     roles.value,
        selectedRole: '',
        onAssignRole: async (roleId) => {
          modal.loading = true;
          modal.error   = '';
          try {
            await api.assignRole(user.id, parseInt(roleId));
            const role = roles.value.find(r => r.id === parseInt(roleId));
            if (role) modal.userRoles.push(role);
            modal.selectedRole = '';
            toast('Роль добавлена');
            await loadUsers();
          } catch (e) {
            modal.error = e.message;
          } finally {
            modal.loading = false;
          }
        },
        onRemoveRole: async (role) => {
          modal.loading = true;
          modal.error   = '';
          try {
            await api.removeRole(user.id, role.id);
            modal.userRoles = modal.userRoles.filter(r => r.id !== role.id);
            toast('Роль удалена');
            await loadUsers();
          } catch (e) {
            modal.error = e.message;
          } finally {
            modal.loading = false;
          }
        },
      });
    };

    const openClearanceModal = (user) => {
      open({
        type:   'form',
        title:  'Уровень допуска',
        fields: [
          { key: 'clearance_level', label: 'Допуск', type: 'select',
            options: [1, 2, 3, 4].map(l => ({ value: l, label: String(l) })) },
        ],
        values: { clearance_level: user.clearance_level || 1 },
        onConfirm: async () => {
          modal.loading = true;
          modal.error   = '';
          try {
            await api.updateClearance(user.id, parseInt(modal.values.clearance_level));
            close();
            toast('Уровень допуска обновлён');
            await loadUsers();
          } catch (e) {
            modal.error = e.message;
          } finally {
            modal.loading = false;
          }
        },
      });
    };

    onMounted(loadUsers);

    return {
      users, loading, error,
      page, totalPages, pageItems, setPage,
      openCreateModal, openDeleteModal, openRolesModal, openClearanceModal,
    };
  },
  components: { AppPagination },
  template: `
    <div class="max-w-5xl mx-auto">
      <div class="flex items-center justify-between mb-5">
        <h2 class="text-lg font-bold text-slate-900">Пользователи</h2>
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
                <th class="text-left px-4 py-3 font-semibold text-slate-600 text-xs uppercase tracking-wide">Логин</th>
                <th class="text-left px-4 py-3 font-semibold text-slate-600 text-xs uppercase tracking-wide">Роли</th>
                <th class="text-left px-4 py-3 font-semibold text-slate-600 text-xs uppercase tracking-wide">Допуск</th>
                <th class="text-right px-4 py-3 font-semibold text-slate-600 text-xs uppercase tracking-wide">Действия</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100">
              <tr v-for="user in pageItems" :key="user.id" class="hover:bg-slate-50 transition-colors">
                <td class="px-4 py-3 font-mono text-xs text-slate-400" :title="user.id">{{ user.id.substring(0,8) }}…</td>
                <td class="px-4 py-3 font-medium text-slate-900">{{ user.login }}</td>
                <td class="px-4 py-3">
                  <span v-for="role in (user.roles || [])" :key="role.id"
                    class="inline-block bg-indigo-100 text-indigo-700 text-xs rounded-full px-2.5 py-0.5 mr-1 font-medium">
                    {{ role.name }}
                  </span>
                  <span v-if="!user.roles || user.roles.length === 0" class="text-slate-400 text-xs">—</span>
                </td>
                <td class="px-4 py-3">
                  <button @click="openClearanceModal(user)"
                    class="inline-block bg-slate-100 hover:bg-slate-200 text-slate-700 text-xs rounded-full px-2.5 py-0.5 font-medium transition">
                    {{ user.clearance_level || 1 }}
                  </button>
                </td>
                <td class="px-4 py-3 text-right">
                  <button @click="openRolesModal(user)"
                    class="text-indigo-600 hover:text-indigo-800 text-xs font-medium transition mr-3">Роли</button>
                  <button @click="openDeleteModal(user)"
                    class="text-red-500 hover:text-red-700 text-xs font-medium transition">Удалить</button>
                </td>
              </tr>
              <tr v-if="users.length === 0">
                <td colspan="5" class="text-center py-10 text-slate-400 text-sm">Пользователи не найдены</td>
              </tr>
            </tbody>
          </table>
        </div>
        <AppPagination :page="page" :total-pages="totalPages" @update:page="setPage" />
      </div>
    </div>
  `,
};
