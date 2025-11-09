function validatePassword(pw) {
    if (pw.length > 72) {
        return "Слишком длинный пароль";
    }
    if (pw.length < 8) {
        return "Пароль должен быть от 8 символов";
    }

    let hasLetter = false;
    let hasDigit = false;

    for (let i = 0; i < pw.length; i++) {
        const char = pw[i];
        const code = char.charCodeAt(0);

        if (code > 127) {
            return "Недопустимые символы в пароле";
        }

        if ((code >= 65 && code <= 90) || (code >= 97 && code <= 122)) {
            hasLetter = true;
        } else if (code >= 48 && code <= 57) {
            hasDigit = true;
        }
    }

    if (!hasLetter || !hasDigit) {
        return "Пароль должен содержать латинские буквы и цифры";
    }

    return "";
}
