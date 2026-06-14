import { api } from '../api.js';
import { useToast } from '../composables/useToast.js';
import { useModal } from '../composables/useModal.js';

export default {
  setup() {
    const { show: toast }        = useToast();
    const { modal, open, close } = useModal();

    const showResult = (title, result) => {
      open({
        type:    'result',
        title,
        message: 'Результат операции:',
        resultText: `Обновлено: ${result.updated_count}, ошибок: ${result.failed_count}`,
      });
    };

    const openRotateMasterKeyModal = () => {
      open({
        type:    'confirm',
        title:   'Ротация мастер-ключа',
        message: 'Будет создана новая версия мастер-ключа Vault, и все обёртки ключей шифрования данных (DEK) будут перешифрованы новой версией. Сами данные не изменятся. Продолжить?',
        confirmLabel: 'Запустить ротацию',
        confirmLoadingLabel: 'Выполняется...',
        confirmClass: 'bg-indigo-600 hover:bg-indigo-700',
        onConfirm: async () => {
          modal.loading = true;
          try {
            const result = await api.rotateMasterKey();
            close();
            toast('Ротация мастер-ключа завершена');
            showResult('Ротация мастер-ключа', result);
          } catch (e) {
            modal.error = e.message;
          } finally {
            modal.loading = false;
          }
        },
      });
    };

    const openRotateDeksModal = () => {
      open({
        type:    'confirm',
        title:   'Ротация ключей шифрования данных',
        message: 'Все маппинги будут полностью перешифрованы новыми ключами шифрования данных (DEK). Значения токенов, видимые пользователям, не изменятся. Это тяжёлая операция — она может занять продолжительное время. Продолжить?',
        confirmLabel: 'Запустить ротацию',
        confirmLoadingLabel: 'Выполняется...',
        confirmClass: 'bg-indigo-600 hover:bg-indigo-700',
        onConfirm: async () => {
          modal.loading = true;
          try {
            const result = await api.rotateAllDeks();
            close();
            toast('Ротация ключей шифрования данных завершена');
            showResult('Ротация ключей шифрования данных', result);
          } catch (e) {
            modal.error = e.message;
          } finally {
            modal.loading = false;
          }
        },
      });
    };

    return { openRotateMasterKeyModal, openRotateDeksModal };
  },
  template: `
    <div class="max-w-3xl mx-auto">
      <h2 class="text-lg font-bold text-slate-900 mb-5">Безопасность</h2>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div class="bg-white rounded-xl border border-slate-200 p-6 shadow-sm">
          <h3 class="font-semibold text-slate-800 mb-2">Ротация мастер-ключа</h3>
          <p class="text-sm text-slate-500 mb-4">
            Создаёт новую версию мастер-ключа (KEK) в Vault и перешифровывает обёртки
            ключей шифрования данных (DEK) всех маппингов последней версией ключа.
            Сами данные не изменяются.
          </p>
          <button @click="openRotateMasterKeyModal"
            class="px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white text-sm font-medium rounded-lg transition">
            Ротировать мастер-ключ
          </button>
        </div>

        <div class="bg-white rounded-xl border border-slate-200 p-6 shadow-sm">
          <h3 class="font-semibold text-slate-800 mb-2">Ротация ключей шифрования данных</h3>
          <p class="text-sm text-slate-500 mb-4">
            Полностью перешифровывает данные всех маппингов новыми ключами шифрования
            данных (DEK). Значения токенов, видимые пользователям, не изменятся.
            Тяжёлая операция.
          </p>
          <button @click="openRotateDeksModal"
            class="px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white text-sm font-medium rounded-lg transition">
            Ротировать ключи данных
          </button>
        </div>
      </div>
    </div>
  `,
};
