import { useToast } from '../composables/useToast.js';

// Roles required to access each section (matches backend RBAC)
const REQUIRED_ROLES = {
  users:    ['admin'],
  roles:    ['admin'],
  kinds:    ['admin'],
  tokens:   ['admin', 'specialist'],
  audit:    ['admin', 'auditor'],
  security: ['admin'],
};

export default {
  props: {
    roles: { type: Array, default: () => [] },
  },
  emits: ['navigate'],
  setup(props, { emit }) {
    const { show } = useToast();

    const items = [
      { id: 'users',  label: 'Пользователи',   icon: '👤', desc: 'Учётные записи и роли'       },
      { id: 'roles',  label: 'Роли',           icon: '🔑', desc: 'Список системных ролей'      },
      { id: 'kinds',  label: 'Виды токенов',   icon: '🏷️', desc: 'Типы данных для токенизации' },
      { id: 'tokens', label: 'Токены',         icon: '🔐', desc: 'Активные маппинги'            },
      { id: 'audit',  label: 'Аудит',          icon: '📋', desc: 'Журнал операций с ПДн'        },
      { id: 'security', label: 'Безопасность', icon: '🛡️', desc: 'Ротация ключей шифрования'    },
    ];

    const hasAccess = (id) => {
      const required = REQUIRED_ROLES[id];
      if (!required) return true;
      const names = props.roles.map(r => r.name);
      return required.some(r => names.includes(r));
    };

    const handleClick = (item) => {
      if (!hasAccess(item.id)) {
        show(`Нет доступа к разделу «${item.label}»`, 'error');
        return;
      }
      emit('navigate', item.id);
    };

    return { items, hasAccess, handleClick };
  },
  template: `
    <div class="max-w-2xl mx-auto py-8">
      <h2 class="text-xl font-bold text-slate-900 mb-6 text-center">Выберите раздел</h2>
      <div class="grid grid-cols-2 gap-4">
        <button
          v-for="item in items" :key="item.id"
          @click="handleClick(item)"
          :class="[
            'rounded-2xl border p-6 text-left transition-all group',
            hasAccess(item.id)
              ? 'bg-white border-slate-200 hover:border-indigo-300 hover:shadow-md'
              : 'bg-slate-50 border-slate-200 opacity-50 cursor-not-allowed',
          ]"
        >
          <div class="text-3xl mb-3">{{ item.icon }}</div>
          <div :class="['font-semibold transition-colors', hasAccess(item.id) ? 'text-slate-800 group-hover:text-indigo-700' : 'text-slate-500']">
            {{ item.label }}
          </div>
          <div class="text-xs text-slate-500 mt-1">{{ item.desc }}</div>
          <div v-if="!hasAccess(item.id)" class="text-xs text-red-400 mt-2">Недостаточно прав</div>
        </button>
      </div>
    </div>
  `,
};
