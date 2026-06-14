import { createApp, ref, computed, onMounted, onUnmounted } from './vue.js';

import { api } from './api.js';
import AppHeader  from './components/AppHeader.js';
import AppModal   from './components/AppModal.js';
import AppToasts  from './components/AppToasts.js';
import LoginView  from './views/LoginView.js';
import MenuView   from './views/MenuView.js';
import UsersView  from './views/UsersView.js';
import RolesView  from './views/RolesView.js';
import KindsView  from './views/KindsView.js';
import TokensView from './views/TokensView.js';
import AuditView  from './views/AuditView.js';
import SecurityView from './views/SecurityView.js';

const MENU_ITEMS = [
  { id: 'users',    label: 'Пользователи'  },
  { id: 'roles',    label: 'Роли'          },
  { id: 'kinds',    label: 'Виды токенов'  },
  { id: 'tokens',   label: 'Токены'        },
  { id: 'audit',    label: 'Аудит'         },
  { id: 'security', label: 'Безопасность'  },
];

const App = {
  components: {
    AppHeader, AppModal, AppToasts,
    LoginView, MenuView,
    UsersView, RolesView, KindsView, TokensView, AuditView, SecurityView,
  },

  setup() {
    const validScreens = new Set(['login', 'menu', ...MENU_ITEMS.map(m => m.id)]);

    const initialScreen = () => {
      const hash = location.hash.slice(1);
      return validScreens.has(hash) ? hash : 'login';
    };

    const screen = ref(initialScreen());
    const userRoles = ref([]);

    const screenTitle = computed(() =>
      MENU_ITEMS.find(m => m.id === screen.value)?.label || ''
    );

    const navigate = (id) => {
      screen.value = id;
      history.pushState({ screen: id }, '', '#' + id);
    };

    const onLogin = (roles) => {
      userRoles.value = roles || [];
      navigate('menu');
    };

    const logout = () => {
      userRoles.value = [];
      screen.value = 'login';
      history.pushState({ screen: 'login' }, '', '#login');
    };

    const onPopState = (e) => {
      const id = e.state?.screen ?? initialScreen();
      screen.value = validScreens.has(id) ? id : 'login';
    };

    onMounted(async () => {
      if (!location.hash) history.replaceState({ screen: screen.value }, '', '#' + screen.value);
      window.addEventListener('popstate', onPopState);

      if (screen.value !== 'login') {
        try {
          const resp = await api.getMe();
          userRoles.value = resp?.roles || [];
        } catch {
          userRoles.value = [];
          screen.value = 'login';
          history.replaceState({ screen: 'login' }, '', '#login');
        }
      }
    });

    onUnmounted(() => {
      window.removeEventListener('popstate', onPopState);
    });

    return { screen, screenTitle, userRoles, navigate, onLogin, logout };
  },

  template: `
    <div>
      <LoginView v-if="screen === 'login'" @login="onLogin" />

      <div v-else class="min-h-screen flex flex-col">
        <AppHeader
          :screen="screen"
          :screenTitle="screenTitle"
          @menu="navigate('menu')"
          @logout="logout"
        />
        <main class="flex-1 p-6">
          <MenuView   v-if="screen === 'menu'"            :roles="userRoles" @navigate="navigate" />
          <UsersView  v-else-if="screen === 'users'"  />
          <RolesView  v-else-if="screen === 'roles'"  />
          <KindsView  v-else-if="screen === 'kinds'"  />
          <TokensView v-else-if="screen === 'tokens'" />
          <AuditView  v-else-if="screen === 'audit'"  />
          <SecurityView v-else-if="screen === 'security'" />
        </main>
      </div>

      <!-- Global overlays -->
      <AppModal />
      <AppToasts />
    </div>
  `,
};

createApp(App).mount('#app');
