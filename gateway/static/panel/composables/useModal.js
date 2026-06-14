import { reactive, computed } from '../vue.js';

const modal = reactive({
  show:    false,
  type:    '',
  title:   '',
  message: '',

  fields:  [],
  values:  {},

  loading: false,
  error:   '',

  onConfirm:    null,
  onAssignRole: null,
  onRemoveRole: null,

  user:         null,
  userRoles:    [],
  allRoles:     [],
  selectedRole: '',
});

const availableRolesToAdd = computed(() => {
  const ids = new Set((modal.userRoles || []).map(r => r.id));
  return (modal.allRoles || []).filter(r => !ids.has(r.id));
});

export function useModal() {
  const open = (config) => {
    Object.assign(modal, { show: true, loading: false, error: '', ...config });
  };

  const close = () => {
    modal.show        = false;
    modal.error       = '';
    modal.loading     = false;
    modal.onConfirm   = null;
    modal.onAssignRole = null;
    modal.onRemoveRole = null;
  };

  return { modal, availableRolesToAdd, open, close };
}
