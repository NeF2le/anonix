function encodeBase64Unicode(str) {
    return btoa(encodeURIComponent(str).replace(/%([0-9A-F]{2})/g,
        function(match, p1) {
            return String.fromCharCode('0x' + p1);
        }));
}

function decodeBase64Unicode(str) {
    return decodeURIComponent(Array.prototype.map.call(atob(str),
        function(c) {
            return '%' + c.charCodeAt(0).toString(16).padStart(2, '0');
        }).join(''));
}

async function tokenize(payload) {
    let ttl = null;
    if ('ttl_select' in payload) {
        if (payload.ttl_select === 'custom') {
            const n = parseInt((payload.ttl_custom || ' ').trim(), 10);
            ttl = Number.isFinite(n) && n > 0 ? n : null;
        } else {
            const n = parseInt(payload.ttl_select, 10);
            ttl = Number.isFinite(n) ? n : null;
        }
    }

    const body = {
        plaintext: encodeBase64Unicode(payload.plaintext),
        token_ttl: ttl === null ? undefined : ttl,
        deterministic: true,
        reversible: true,
    }

    try {
        const resp = await fetchWithAuth(
            `${apiBase}/tokenizer/tokenize`,
            {
                method: 'POST',
                headers: new Headers({ 'Content-Type': 'application/json' }),
                body: JSON.stringify(body)
            }
        )
        const data = await resp.json();

        if (!resp.ok) {
            throw new Error(`Error ${resp.status} ${data.error}`);
        }

        return data.id;
    } catch (error) {
        console.error('Error occurred:', error);
        throw new Error(error.message);
    }
}

async function detokenize(payload) {
    const body = {token: payload.token}

    try {
        const resp = await fetchWithAuth(
            `${apiBase}/tokenizer/detokenize`,
            {
                method: 'POST',
                headers: new Headers({ 'Content-Type': 'application/json' }),
                body: JSON.stringify(body)
            }
        )
        const data = await resp.json();

        if (!resp.ok) {
            if (resp.status === 400) {
                throw new Error('Некорректный токен')
            }
            throw new Error(`Error ${resp.status} ${data.error}`);
        }

        return decodeBase64Unicode(data.plaintext);
    } catch (error) {
        console.error('Error occurred:', error);
        throw new Error(error.message);
    }
}