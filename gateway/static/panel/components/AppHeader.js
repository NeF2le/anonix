export default {
  props: {
    screen:      { type: String, required: true },
    screenTitle: { type: String, default: '' },
  },
  emits: ['menu', 'logout'],
  template: `
    <header class="bg-white border-b border-slate-200 px-6 py-3.5 flex items-center justify-between flex-shrink-0">
      <div class="flex items-center gap-3">
        <!-- Logo mark -->
        <div class="w-7 h-7 bg-indigo-600 rounded-lg flex items-center justify-center flex-shrink-0">
          <svg class="w-4 h-4 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/>
          </svg>
        </div>
        <span class="font-bold text-slate-900">Anonix</span>

        <!-- Breadcrumb -->
        <template v-if="screen !== 'menu'">
          <span class="text-slate-300">/</span>
          <button @click="$emit('menu')" class="text-slate-500 hover:text-indigo-600 text-sm transition">
            Меню
          </button>
          <span class="text-slate-300">/</span>
          <span class="text-slate-700 text-sm font-medium">{{ screenTitle }}</span>
        </template>
      </div>

      <!-- Logout -->
      <button @click="$emit('logout')"
        class="flex items-center gap-1.5 text-sm text-slate-500 hover:text-red-600 transition">
        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"/>
        </svg>
        Выйти
      </button>
    </header>
  `,
};
