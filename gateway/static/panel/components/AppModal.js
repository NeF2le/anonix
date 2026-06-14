import { useModal } from '../composables/useModal.js';

export default {
  setup() {
    const { modal, availableRolesToAdd, close } = useModal();

    const handleConfirm = () => {
      if (modal.onConfirm) modal.onConfirm();
    };

    const handleAssign = () => {
      if (modal.selectedRole && modal.onAssignRole) {
        modal.onAssignRole(modal.selectedRole);
      }
    };

    const handleRemove = (role) => {
      if (modal.onRemoveRole) modal.onRemoveRole(role);
    };

    return { modal, availableRolesToAdd, close, handleConfirm, handleAssign, handleRemove };
  },
  template: `
    <div
      v-if="modal.show"
      class="fixed inset-0 bg-black bg-opacity-40 flex items-center justify-center z-50 p-4"
      @click.self="close"
    >
      <div class="bg-white rounded-2xl shadow-2xl w-full max-w-md">

        <!-- Header -->
        <div class="px-6 py-4 border-b border-slate-200 flex items-center justify-between">
          <h3 class="font-semibold text-slate-900">{{ modal.title }}</h3>
          <button @click="close" class="text-slate-400 hover:text-slate-700 transition text-2xl leading-none">×</button>
        </div>

        <div class="p-6">

          <!-- Confirm dialog -->
          <template v-if="modal.type === 'confirm'">
            <p class="text-slate-600 text-sm mb-6">{{ modal.message }}</p>
            <div v-if="modal.error" class="mb-4 p-3 bg-red-50 border border-red-200 text-red-700 rounded-lg text-sm">
              {{ modal.error }}
            </div>
            <div class="flex gap-3 justify-end">
              <button @click="close" class="px-4 py-2 text-sm text-slate-600 hover:text-slate-900 transition">
                Отмена
              </button>
              <button @click="handleConfirm" :disabled="modal.loading"
                :class="['px-4 py-2 disabled:opacity-50 text-white text-sm font-medium rounded-lg transition', modal.confirmClass || 'bg-red-600 hover:bg-red-700']">
                {{ modal.loading ? (modal.confirmLoadingLabel || 'Удаление...') : (modal.confirmLabel || 'Удалить') }}
              </button>
            </div>
          </template>

          <!-- Form dialog -->
          <template v-else-if="modal.type === 'form'">
            <form @submit.prevent="handleConfirm">
              <template v-for="field in modal.fields" :key="field.key">
              <div v-if="!field.showIf || field.showIf(modal.values)" class="mb-4">
                <label v-if="field.type !== 'checkbox'" class="block text-sm font-medium text-slate-700 mb-1.5">{{ field.label }}</label>
                <select v-if="field.type === 'select'" v-model="modal.values[field.key]"
                  class="w-full px-3 py-2.5 border border-slate-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500">
                  <option v-for="opt in field.options" :key="opt.value" :value="opt.value">{{ opt.label }}</option>
                </select>
                <label v-else-if="field.type === 'checkbox'" class="flex items-center gap-2 text-sm font-medium text-slate-700 cursor-pointer select-none">
                  <input type="checkbox" v-model="modal.values[field.key]"
                    class="w-4 h-4 rounded border-slate-300 text-indigo-600 focus:ring-indigo-500"/>
                  {{ field.label }}
                </label>
                <input v-else
                  v-model="modal.values[field.key]"
                  :type="field.type"
                  :required="field.required !== false"
                  :placeholder="field.placeholder || ''"
                  class="w-full px-3 py-2.5 border border-slate-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 transition"
                />
              </div>
              </template>
              <div v-if="modal.error" class="mb-4 p-3 bg-red-50 border border-red-200 text-red-700 rounded-lg text-sm">
                {{ modal.error }}
              </div>
              <div class="flex gap-3 justify-end">
                <button type="button" @click="close" class="px-4 py-2 text-sm text-slate-600 hover:text-slate-900 transition">
                  Отмена
                </button>
                <button type="submit" :disabled="modal.loading"
                  class="px-4 py-2 bg-indigo-600 hover:bg-indigo-700 disabled:opacity-50 text-white text-sm font-medium rounded-lg transition">
                  {{ modal.loading ? 'Сохранение...' : 'Сохранить' }}
                </button>
              </div>
            </form>
          </template>

          <!-- Result dialog -->
          <template v-else-if="modal.type === 'result'">
            <p class="text-sm text-slate-500 mb-2">{{ modal.message }}</p>
            <div class="font-mono text-sm bg-slate-50 border border-slate-200 rounded-lg p-3 break-all select-all">
              {{ modal.resultText }}
            </div>
            <div class="flex justify-end mt-5">
              <button @click="close" class="px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white text-sm font-medium rounded-lg transition">
                Закрыть
              </button>
            </div>
          </template>

          <!-- Roles manager -->
          <template v-else-if="modal.type === 'roles'">
            <p class="text-sm text-slate-500 mb-4">
              Пользователь:
              <span class="font-semibold text-slate-800">{{ modal.user && modal.user.login }}</span>
            </p>

            <div class="mb-5">
              <p class="text-xs font-semibold text-slate-500 uppercase tracking-wide mb-2">Текущие роли</p>
              <div class="flex flex-wrap gap-2 min-h-8 p-3 bg-slate-50 rounded-lg border border-slate-200">
                <span
                  v-for="role in modal.userRoles" :key="role.id"
                  class="inline-flex items-center gap-1.5 bg-indigo-100 text-indigo-700 text-xs rounded-full px-3 py-1 font-medium"
                >
                  {{ role.name }}
                  <button @click="handleRemove(role)" :disabled="modal.loading"
                    class="hover:text-red-700 disabled:opacity-40 transition font-bold leading-none">×</button>
                </span>
                <span v-if="modal.userRoles.length === 0" class="text-slate-400 text-xs self-center">
                  Нет ролей
                </span>
              </div>
            </div>

            <div class="mb-4">
              <p class="text-xs font-semibold text-slate-500 uppercase tracking-wide mb-2">Добавить роль</p>
              <div class="flex gap-2">
                <select v-model="modal.selectedRole"
                  class="flex-1 px-3 py-2 border border-slate-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500">
                  <option value="">Выберите роль...</option>
                  <option v-for="role in availableRolesToAdd" :key="role.id" :value="role.id">{{ role.name }}</option>
                </select>
                <button @click="handleAssign" :disabled="!modal.selectedRole || modal.loading"
                  class="px-4 py-2 bg-indigo-600 hover:bg-indigo-700 disabled:opacity-50 text-white text-sm rounded-lg transition font-medium">
                  Добавить
                </button>
              </div>
            </div>

            <div v-if="modal.error" class="mb-4 p-3 bg-red-50 border border-red-200 text-red-700 rounded-lg text-sm">
              {{ modal.error }}
            </div>
            <div class="flex justify-end">
              <button @click="close" class="px-4 py-2 text-sm text-slate-600 hover:text-slate-900 transition">
                Закрыть
              </button>
            </div>
          </template>

        </div>
      </div>
    </div>
  `,
};
