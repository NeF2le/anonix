async function fetchWithAuth(url, options = {}) {
    try {
        const response = await fetch(url, {
            ...options,
            credentials: 'include'
        });

        if (response.status === 401) {
            throw new Error("Сессия истекла. Войдите в аккаунт")
        }

        return response;
    } catch (error) {
        console.error('Fetch error:', error);
        throw error;
    }
}