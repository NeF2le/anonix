import { ref, onMounted, onUnmounted } from '../vue.js';
import { api } from '../api.js';
import { useToast } from '../composables/useToast.js';
import { useModal } from '../composables/useModal.js';
import { usePagination } from '../composables/usePagination.js';
import AppPagination from '../components/AppPagination.js';

const formatDate = (str) => {
  if (!str) return '—';
  try {
    return new Date(str).toLocaleDateString('ru-RU', { day: '2-digit', month: '2-digit', year: 'numeric' });
  } catch {
    return str;
  }
};

// Parses a Go time.Duration string (e.g. "24h0m0s") into milliseconds.
const parseGoDuration = (str) => {
  if (!str) return null;
  const re = /(\d+(?:\.\d+)?)(ns|us|µs|ms|s|m|h)/g;
  let totalMs = 0;
  let found = false;
  let match;
  while ((match = re.exec(str)) !== null) {
    found = true;
    const value = parseFloat(match[1]);
    switch (match[2]) {
      case 'h':  totalMs += value * 3600000; break;
      case 'm':  totalMs += value * 60000; break;
      case 's':  totalMs += value * 1000; break;
      case 'ms': totalMs += value; break;
      case 'us':
      case 'µs': totalMs += value / 1000; break;
      case 'ns': totalMs += value / 1e6; break;
    }
  }
  return found ? totalMs : null;
};

const formatRemainingMs = (ms) => {
  if (ms < 60000) return '< 1 мин';
  const totalMinutes = Math.floor(ms / 60000);
  const days    = Math.floor(totalMinutes / 1440);
  const hours   = Math.floor((totalMinutes % 1440) / 60);
  const minutes = totalMinutes % 60;
  const parts = [];
  if (days)    parts.push(`${days}д`);
  if (hours)   parts.push(`${hours}ч`);
  if (minutes || parts.length === 0) parts.push(`${minutes}м`);
  return parts.join(' ');
};

const b64ToText = (b64) => {
  try {
    const bytes = Uint8Array.from(atob(b64), c => c.charCodeAt(0));
    return new TextDecoder().decode(bytes);
  } catch {
    return b64;
  }
};

const OTHER_KIND_OPTION = { value: 0, label: 'Другое' };

const ALGO_LABELS = {
  'aes-256-siv':              'AES-SIV',
  'aes-256-siv-random':       'AES-SIV',
  'gost-kuznechik-mgm':       'ГОСТ Кузнечик',
  'gost-kuznechik-mgm-random':'ГОСТ Кузнечик',
  '':                         'AES-SIV',
};

export default {
  setup() {
    const { show: toast }        = useToast();
    const { modal, open, close } = useModal();

    const tokens  = ref([]);
    const kinds   = ref([]);
    const loading = ref(false);
    const error   = ref('');

    // Ticks once a minute so the "remaining TTL" column stays up to date.
    const now = ref(Date.now());
    const nowTimer = setInterval(() => { now.value = Date.now(); }, 60000);
    onUnmounted(() => clearInterval(nowTimer));

    const formatTtl = (token) => {
      const ttlMs = parseGoDuration(token.token_ttl);
      if (ttlMs === null || ttlMs === 0) return '∞';
      const createdAt = new Date(token.created_at).getTime();
      if (isNaN(createdAt)) return '—';
      const remainingMs = createdAt + ttlMs - now.value;
      return remainingMs <= 0 ? 'Истёк' : formatRemainingMs(remainingMs);
    };

    const { page, totalPages, pageItems, setPage } = usePagination(tokens);

    const loadTokens = async () => {
      loading.value = true;
      error.value   = '';
      try {
        const data   = await api.getMappings();
        tokens.value = Array.isArray(data) ? data : [];
      } catch (e) {
        error.value = e.message;
      } finally {
        loading.value = false;
      }
    };

    const loadKinds = async () => {
      try {
        const data  = await api.getKinds();
        kinds.value = Array.isArray(data) ? data : (data?.kinds ?? []);
      } catch {}
    };

    // ── Tokenize ─────────────────────────────────────────────────────────────
    const openTokenizeModal = () => {
      const kindOptions = [OTHER_KIND_OPTION, ...kinds.value.map(k => ({ value: k.id, label: k.russian_name || k.name }))];
      const modeOptions = [
        { value: 'pseudonymize', label: 'Псевдонимизация' },
        { value: 'anonymize',    label: 'Анонимизация' },
      ];
      const algorithmOptions = [
        { value: 'aes-siv',       label: 'AES-SIV' },
        { value: 'gost-kuznechik', label: 'ГОСТ Кузнечик (MGM)' },
      ];

      open({
        type:   'form',
        title:  'Зашифровать данные',
        fields: [
          { key: 'plaintext',     label: 'Исходные данные',            type: 'text',     required: true,  placeholder: 'Введите текст' },
          { key: 'mode',          label: 'Вид шифрования',             type: 'select',   required: true,  options: modeOptions },
          { key: 'kind_id',       label: 'Вид данных',                 type: 'select',   required: true,  options: kindOptions },
          { key: 'deterministic', label: 'Детерминированный',          type: 'checkbox' },
          { key: 'algorithm',     label: 'Алгоритм шифрования',        type: 'select',   required: false, options: algorithmOptions, showIf: v => v.mode === 'pseudonymize' },
          { key: 'token_ttl',     label: 'Срок жизни (часов, 0 = ∞)', type: 'number',   required: false, placeholder: '0', showIf: v => v.mode === 'pseudonymize' },
        ],
        values: { plaintext: '', mode: 'pseudonymize', kind_id: OTHER_KIND_OPTION.value, deterministic: true, algorithm: 'aes-siv', token_ttl: '0' },
        onConfirm: async () => {
          modal.loading = true;
          modal.error   = '';
          try {
            const ttlSec = parseInt(modal.values.token_ttl || '0') * 3600;
            const result = await api.tokenize(
              modal.values.plaintext,
              modal.values.kind_id,
              !!modal.values.deterministic,
              modal.values.mode,
              ttlSec,
              modal.values.algorithm,
            );
            close();
            if (result && result.id) {
              toast(`Токен создан: ${result.token}`);
              await loadTokens();
            } else {
              open({
                type:       'result',
                title:      'Результат шифрования',
                message:    'Токен (необратимый — сохранить негде, дешифровка невозможна):',
                resultText: result.token,
              });
            }
          } catch (e) {
            modal.error = e.message;
          } finally {
            modal.loading = false;
          }
        },
      });
    };

    // ── Detokenize ────────────────────────────────────────────────────────────
    const openDetokenizeModal = (prefillToken = '') => {
      open({
        type:   'form',
        title:  'Дешифровать токен',
        fields: [
          { key: 'token', label: 'Токен', type: 'text', required: true, placeholder: 'Например, fio_7f82a1c3' },
        ],
        values: { token: prefillToken },
        onConfirm: async () => {
          modal.loading = true;
          modal.error   = '';
          try {
            const resp = await api.detokenize(modal.values.token);
            const text = b64ToText(resp.plaintext);
            close();
            open({
              type:       'result',
              title:      'Результат дешифрования',
              message:    'Исходные данные:',
              resultText: text,
            });
          } catch (e) {
            modal.error = e.status === 404 ? 'Токен не найден или истёк срок жизни' : e.message;
          } finally {
            modal.loading = false;
          }
        },
      });
    };

    // ── Edit TTL ──────────────────────────────────────────────────────────────
    const openEditTtlModal = (token) => {
      open({
        type:   'form',
        title:  'Изменить TTL токена',
        fields: [{ key: 'ttl_hours', label: 'Срок жизни (часов)', type: 'number', required: true, placeholder: '24' }],
        values: { ttl_hours: '24' },
        onConfirm: async () => {
          modal.loading = true;
          modal.error   = '';
          try {
            const ttlNs = parseInt(modal.values.ttl_hours) * 3600 * 1_000_000_000;
            await api.updateMapping(token.id, ttlNs);
            close();
            toast('TTL обновлён');
            await loadTokens();
          } catch (e) {
            modal.error = e.message;
          } finally {
            modal.loading = false;
          }
        },
      });
    };

    // ── Delete ────────────────────────────────────────────────────────────────
    const openDeleteModal = (token) => {
      open({
        type:    'confirm',
        title:   'Удалить токен',
        message: `Удалить токен ${token.token}? Действие необратимо.`,
        onConfirm: async () => {
          modal.loading = true;
          try {
            await api.deleteMapping(token.id);
            close();
            toast('Токен удалён');
            await loadTokens();
          } catch (e) {
            modal.error = e.message;
          } finally {
            modal.loading = false;
          }
        },
      });
    };

    onMounted(() => { loadTokens(); loadKinds(); });

    return {
      tokens, loading, error,
      page, totalPages, pageItems, setPage,
      loadTokens, formatDate, formatTtl,
      openTokenizeModal, openDetokenizeModal,
      openEditTtlModal, openDeleteModal,
      algoLabel: (algoName) => ALGO_LABELS[algoName] || algoName || 'AES-SIV',
    };
  },
  components: { AppPagination },
  template: `
    <div class="max-w-6xl mx-auto">
      <div class="flex items-center justify-between mb-5">
        <h2 class="text-lg font-bold text-slate-900">Токены</h2>
        <div class="flex items-center gap-2">
          <button @click="openDetokenizeModal()"
            class="inline-flex items-center gap-1.5 text-sm text-slate-700 hover:text-slate-900 border border-slate-300 hover:border-slate-400 bg-white px-3 py-1.5 rounded-lg transition font-medium">
            🔓 Дешифровать
          </button>
          <button @click="openTokenizeModal"
            class="inline-flex items-center gap-1.5 text-sm text-white bg-indigo-600 hover:bg-indigo-700 px-3 py-1.5 rounded-lg transition font-medium">
            🔒 Зашифровать
          </button>
          <button @click="loadTokens"
            class="inline-flex items-center gap-1.5 text-sm text-slate-500 hover:text-slate-800 border border-slate-300 hover:border-slate-400 px-3 py-1.5 rounded-lg transition">
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
            </svg>
            Обновить
          </button>
        </div>
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
                <th class="text-left px-4 py-3 font-semibold text-slate-600 text-xs uppercase tracking-wide">Токен</th>
                <th class="text-left px-4 py-3 font-semibold text-slate-600 text-xs uppercase tracking-wide">Вид данных</th>
                <th class="text-left px-4 py-3 font-semibold text-slate-600 text-xs uppercase tracking-wide">Осталось</th>
                <th class="text-left px-4 py-3 font-semibold text-slate-600 text-xs uppercase tracking-wide">Создан</th>
                <th class="text-left px-4 py-3 font-semibold text-slate-600 text-xs uppercase tracking-wide">Тип</th>
                <th class="text-left px-4 py-3 font-semibold text-slate-600 text-xs uppercase tracking-wide">Алгоритм</th>
                <th class="text-right px-4 py-3 font-semibold text-slate-600 text-xs uppercase tracking-wide">Действия</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100">
              <tr v-for="token in pageItems" :key="token.id" class="hover:bg-slate-50 transition-colors">
                <td class="px-4 py-3 font-mono text-xs text-slate-400" :title="token.id">{{ token.token }}</td>
                <td class="px-4 py-3 text-slate-700">{{ token.kind ? token.kind.russian_name : 'Другое' }}</td>
                <td class="px-4 py-3 text-slate-600 font-mono text-xs">{{ formatTtl(token) }}</td>
                <td class="px-4 py-3 text-slate-600">{{ formatDate(token.created_at) }}</td>
                <td class="px-4 py-3">
                  <span
                    :class="token.deterministic ? 'bg-emerald-100 text-emerald-700' : 'bg-blue-100 text-blue-700'"
                    class="text-xs rounded-full px-2.5 py-0.5 font-medium"
                  >
                    {{ token.deterministic ? 'Детерм.' : 'Случ.' }}
                  </span>
                </td>
                <td class="px-4 py-3 text-slate-600 text-xs">{{ algoLabel(token.algo_name) }}</td>
                <td class="px-4 py-3 text-right">
                  <button @click="openDetokenizeModal(token.token)"
                    class="text-emerald-600 hover:text-emerald-800 text-xs font-medium transition mr-3">Дешифр.</button>
                  <button @click="openEditTtlModal(token)"
                    class="text-indigo-600 hover:text-indigo-800 text-xs font-medium transition mr-3">Изм. TTL</button>
                  <button @click="openDeleteModal(token)"
                    class="text-red-500 hover:text-red-700 text-xs font-medium transition">Удалить</button>
                </td>
              </tr>
              <tr v-if="tokens.length === 0">
                <td colspan="7" class="text-center py-10 text-slate-400 text-sm">Токены не найдены</td>
              </tr>
            </tbody>
          </table>
        </div>
        <AppPagination :page="page" :total-pages="totalPages" @update:page="setPage" />
      </div>
    </div>
  `,
};
