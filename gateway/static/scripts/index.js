const toast = new Toast();
const apiBase = window.APP_CONFIG.API_BASE_URL;

const signInModalInstance = new Modal(
    document.getElementById('signInOverlay'),
    {
        onSubmit: signIn,
        toast
    }
);

const signUpModalInstance = new Modal(
    document.getElementById('signUpOverlay'),
    {
        onSubmit: signUp,
        toast
    }
);

signUpBtn = document.getElementById('signUpBtn');
signInBtn = document.getElementById('signInBtn');

signInBtn.addEventListener('click', () => signInModalInstance.open());
signUpBtn.addEventListener('click', () => signUpModalInstance.open());