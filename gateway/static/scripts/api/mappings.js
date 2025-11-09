async function getMappings() {
    try {
        const res = await fetchWithAuth(
            `${apiBase}/mappings/`,
            {
                method: "GET",
                headers: {"Content-Type": "application/json"}
            }
        )

        if (!res.ok) {
            console.error('HTTP ${res.status} ${res.statusText} ')
            throw new Error(`Failed to get mappings from server: ${res.statusText}`);
        }

        const data = await res.json();
        if (!Array.isArray(data)) {
            throw new Error(`Data of response is not an array`);
        }

        return data.map(item => {
            return {
                id: item.id,
                cipher_text: item.cipher_text,
                dek_wrapped: item.dek_wrapped,
                token_ttl: ttlToHuman(item.token_ttl),
                created_at: timestampToDate(item.created_at)
            };
        });
    } catch (err) {
        throw err;
    }
}

async function deleteMapping(mappingId) {
    try {
        const res = await fetchWithAuth(
            `${apiBase}/mappings/${mappingId}`,
            {
                method: "DELETE",
                headers: {"Content-Type": "application/json"}
            }
        )

        if (!res.ok) {
            console.error(`HTTP ${res.status} ${res.statusText}`);
            throw new Error(`Failed to delete token ${mappingId}`);
        }

        return;
    } catch (err) {
        throw err;
    }
}

