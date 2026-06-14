import { ref, computed, watch } from '../vue.js';

const DEFAULT_PAGE_SIZE = 10;

// Paginates a reactive array. Resets to the last valid page if the
// underlying list shrinks (e.g. after a delete or reload).
export function usePagination(itemsRef, pageSize = DEFAULT_PAGE_SIZE) {
  const page = ref(1);

  const totalPages = computed(() => Math.max(1, Math.ceil((itemsRef.value?.length || 0) / pageSize)));

  const pageItems = computed(() => {
    const start = (page.value - 1) * pageSize;
    return (itemsRef.value || []).slice(start, start + pageSize);
  });

  watch(itemsRef, () => {
    if (page.value > totalPages.value) page.value = totalPages.value;
  });

  const setPage = (p) => {
    page.value = Math.min(Math.max(1, p), totalPages.value);
  };

  return { page, totalPages, pageItems, setPage };
}
