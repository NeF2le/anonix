import { ref, reactive, computed } from '../vue.js';
import { api } from '../api.js';

export default {
  emits: ['login'],
  setup(_, { emit }) {
    const mode = ref('login'); // 'login' | 'register'

    const form    = reactive({ login: '', password: '', passwordAgain: '' });
    const loading = ref(false);
    const error   = ref('');
    const success = ref('');

    const isRegister = computed(() => mode.value === 'register');

    const reset = (newMode) => {
      mode.value          = newMode;
      error.value         = '';
      success.value       = '';
      form.login          = '';
      form.password       = '';
      form.passwordAgain  = '';
    };

    const login = async () => {
      loading.value = true;
      error.value   = '';
      try {
        const resp = await api.login(form.login, form.password);
        emit('login', resp?.roles || []);
      } catch (e) {
        error.value = e.status === 401 ? 'Неверный логин или пароль' : (e.message || 'Ошибка входа');
      } finally {
        loading.value = false;
      }
    };

    const register = async () => {
      if (form.password !== form.passwordAgain) {
        error.value = 'Пароли не совпадают';
        return;
      }
      loading.value = true;
      error.value   = '';
      success.value = '';
      try {
        await api.register(form.login, form.password, 2);
        success.value = 'Аккаунт создан — войдите в систему';
        setTimeout(() => reset('login'), 1500);
      } catch (e) {
        error.value = e.status === 409 ? 'Пользователь с таким логином уже существует' : (e.message || 'Ошибка регистрации');
      } finally {
        loading.value = false;
      }
    };

    return { mode, form, loading, error, success, isRegister, reset, login, register };
  },
  template: `
    <div class="min-h-screen flex items-center justify-center p-4">
      <div class="bg-white rounded-2xl shadow-lg p-8 w-full max-w-sm">

        <!-- Logo -->
        <div class="text-center mb-8">
          <div class="w-12 h-12 bg-indigo-600 rounded-xl flex items-center justify-center mx-auto mb-3">
            <svg class="w-6 h-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/>
            </svg>
          </div>
          <h1 class="text-2xl font-bold text-slate-900">Anonix</h1>
          <p class="text-slate-500 text-sm mt-1">{{ isRegister ? 'Регистрация' : 'Панель управления' }}</p>
        </div>

        <!-- Login form -->
        <form v-if="!isRegister" @submit.prevent="login">
          <div class="mb-4">
            <label class="block text-sm font-medium text-slate-700 mb-1.5">Логин</label>
            <input v-model="form.login" type="text" required autocomplete="off"
              class="w-full px-3 py-2.5 border border-slate-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition"
              placeholder="Введите логин"/>
          </div>
          <div class="mb-5">
            <label class="block text-sm font-medium text-slate-700 mb-1.5">Пароль</label>
            <input v-model="form.password" type="password" required autocomplete="off"
              class="w-full px-3 py-2.5 border border-slate-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition"
              placeholder="Введите пароль"/>
          </div>
          <div v-if="error" class="mb-4 p-3 bg-red-50 border border-red-200 text-red-700 rounded-lg text-sm">{{ error }}</div>
          <button type="submit" :disabled="loading"
            class="w-full bg-indigo-600 hover:bg-indigo-700 disabled:opacity-50 text-white font-semibold py-2.5 rounded-lg text-sm transition">
            <span v-if="loading" class="inline-flex items-center justify-center gap-2">
              <svg class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"/>
              </svg>
              Вход...
            </span>
            <span v-else>Войти</span>
          </button>
          <p class="text-center text-sm text-slate-500 mt-5">
            Нет аккаунта?
            <button type="button" @click="reset('register')" class="text-indigo-600 hover:text-indigo-800 font-medium transition">
              Зарегистрироваться
            </button>
          </p>
        </form>

        <!-- Register form -->
        <form v-else @submit.prevent="register">
          <div class="mb-4">
            <label class="block text-sm font-medium text-slate-700 mb-1.5">Логин</label>
            <input v-model="form.login" type="text" required autocomplete="off"
              class="w-full px-3 py-2.5 border border-slate-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition"
              placeholder="Придумайте логин"/>
          </div>
          <div class="mb-4">
            <label class="block text-sm font-medium text-slate-700 mb-1.5">Пароль</label>
            <input v-model="form.password" type="password" required autocomplete="off"
              class="w-full px-3 py-2.5 border border-slate-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition"
              placeholder="Придумайте пароль"/>
          </div>
          <div class="mb-5">
            <label class="block text-sm font-medium text-slate-700 mb-1.5">Повторите пароль</label>
            <input v-model="form.passwordAgain" type="password" required autocomplete="off"
              class="w-full px-3 py-2.5 border border-slate-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition"
              placeholder="Повторите пароль"/>
          </div>
          <div v-if="error"   class="mb-4 p-3 bg-red-50     border border-red-200   text-red-700   rounded-lg text-sm">{{ error }}</div>
          <div v-if="success" class="mb-4 p-3 bg-emerald-50 border border-emerald-200 text-emerald-700 rounded-lg text-sm">{{ success }}</div>
          <button type="submit" :disabled="loading"
            class="w-full bg-indigo-600 hover:bg-indigo-700 disabled:opacity-50 text-white font-semibold py-2.5 rounded-lg text-sm transition">
            <span v-if="loading" class="inline-flex items-center justify-center gap-2">
              <svg class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"/>
              </svg>
              Регистрация...
            </span>
            <span v-else>Создать аккаунт</span>
          </button>
          <p class="text-center text-sm text-slate-500 mt-5">
            Уже есть аккаунт?
            <button type="button" @click="reset('login')" class="text-indigo-600 hover:text-indigo-800 font-medium transition">
              Войти
            </button>
          </p>
        </form>

      </div>
    </div>
  `,
};
