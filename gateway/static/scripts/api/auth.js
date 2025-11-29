async function signUp(payload) {
    if (payload.password !== payload.passwordAgain) {
        throw new Error('Пароли не совпадают');
    }

    const validateErr = validatePassword(payload.password)
    if (validateErr !== "") {
        throw new Error(validateErr);
    }

    const body = {
        login: payload.login,
        password: payload.password,
        role_id: 2
    }

    try {
        const resp = await fetch(`${apiBase}/auth/signUp`, {
            method: 'POST',
            headers: new Headers({ 'Content-Type': 'application/json' }),
            body: JSON.stringify(body)
        })
        const data = await resp.json();

        if (!resp.ok) {
            if (resp.status === 409) {
                throw new Error('Пользователь с таким логином уже существует')
            }
            throw new Error(`Error ${resp.status} ${data.error}`);
        }

        return "Вы зарегистрированы";
    } catch (error) {
        console.error('Error occurred:', error);
        throw new Error(error.message);
    }
}

async function signIn(payload) {
    const body = {
        login: payload.login,
        password: payload.password,
    }

    try {
        const resp = await fetch(`${apiBase}/auth/signIn`, {
            method: 'POST',
            headers: new Headers({ 'Content-Type': 'application/json' }),
            body: JSON.stringify(body)
        })
        const data = await resp.json();

        if (!resp.ok) {
            if (resp.status === 401) {
                throw new Error('Неверный логин или пароль')
            }
            throw new Error(`Error ${resp.status} ${data.error}`);
        }

        const isAdmin = await isAdminCheck(data.user_id);
        if (isAdmin) {
            window.location.href = "admin.html";
            return "signed in as admin"
        }

        return "signed in as default";
    } catch (error) {
        console.error('Error occurred:', error);
        throw new Error(error.message);
    }
}

async function isAdminCheck(userId) {
    const body = {
        user_id: userId
    }

    try {
        const resp = await fetch(`${apiBase}/user/isAdmin`, {
            method: 'POST',
            headers: new Headers({ 'Content-Type': 'application/json' }),
            body: JSON.stringify(body)
        })

        const data = await resp.json()

        if (!resp.ok) {
            throw new Error(`Error ${resp.status} ${data.error}`);
        }

        return data.result;
    }
    catch (error) {
        console.error('Error occurred:', error);
        throw new Error(error.message);
    }
}

