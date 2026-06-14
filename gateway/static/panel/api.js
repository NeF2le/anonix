import { translateError } from './errorMessages.js';

const getBase = () => {
  try { return (window.APP_CONFIG && window.APP_CONFIG.API_BASE_URL) || '/api/v1'; }
  catch { return '/api/v1'; }
};

const call = async (method, path, body = null) => {
  const opts = {
    method,
    credentials: 'include',
    headers: { 'Content-Type': 'application/json' },
  };
  if (body !== null) opts.body = JSON.stringify(body);

  const res = await fetch(`${getBase()}${path}`, opts);
  if (!res.ok) {
    let msg = res.statusText;
    try { const d = await res.json(); msg = d.error || d.message || msg; } catch {}
    const err = new Error(translateError(msg) || `HTTP ${res.status}`);
    err.status = res.status;
    throw err;
  }
  const text = await res.text();
  return text ? JSON.parse(text) : null;
};

export const api = {
  login:    (login, password)           => call('POST',   '/auth/signIn', { login, password }),
  register: (login, password, role_id)  => call('POST',   '/auth/signUp', { login, password, role_id }),
  getMe:    ()                          => call('GET',    '/auth/me'),

  getUsers:   ()               => call('GET',    '/user/list'),
  deleteUser: (userId)         => call('DELETE', '/user/delete',     { user_id: userId }),
  assignRole: (userId, roleId) => call('POST',   '/user/assignRole', { user_id: userId, role_id: roleId }),
  removeRole: (userId, roleId) => call('DELETE', '/user/removeRole', { user_id: userId, role_id: roleId }),
  updateClearance: (userId, level) => call('PATCH', '/user/clearance', { user_id: userId, clearance_level: level }),

  getRoles: () => call('GET', '/role/list'),

  getKinds:   ()          => call('GET',    '/kinds/'),
  createKind: (data)      => call('POST',   '/kinds/', data),
  updateKind: (id, data)  => call('PATCH',  `/kinds/${id}`, data),
  deleteKind: (id)        => call('DELETE', `/kinds/${id}`),

  getMappings:    ()           => call('GET',    '/mappings/'),
  deleteMapping:  (id)         => call('DELETE', `/mappings/${id}`),
  updateMapping:  (id, ttlNs)  => call('PATCH',  `/mappings/${id}`, { token_ttl: ttlNs }),

  getAuditLog: () => call('GET', '/audit/'),

  tokenize:    (plaintext, kindId, deterministic, mode, tokenTtlSec, algorithm) => {
    const bytes = new TextEncoder().encode(plaintext);
    const b64   = btoa(String.fromCharCode(...bytes));
    return call('POST', '/tokenizer/tokenize', {
      plaintext:     b64,
      kind_id:       kindId,
      deterministic: deterministic,
      mode:          mode,
      token_ttl:     tokenTtlSec,
      algorithm:     algorithm,
    });
  },
  detokenize:  (token) => call('POST', '/tokenizer/detokenize', { token }),

  rotateMasterKey: () => call('POST', '/admin/keys/rotate-master', {}),
  rotateAllDeks:   () => call('POST', '/admin/keys/rotate-deks', {}),
};
