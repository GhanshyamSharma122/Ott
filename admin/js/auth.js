// Auth module
const Auth = {
    init() {
        const loginForm = document.getElementById('login-form');
        loginForm.addEventListener('submit', (e) => {
            e.preventDefault();
            this.handleLogin();
        });
    },

    async handleLogin() {
        const email = document.getElementById('login-email').value;
        const password = document.getElementById('login-password').value;
        const errorEl = document.getElementById('login-error');
        const btn = document.getElementById('login-btn');

        errorEl.textContent = '';
        btn.textContent = 'Signing in...';
        btn.disabled = true;

        try {
            const data = await API.login(email, password);

            if (!data.user.is_admin) {
                throw new Error('Access denied. Admin privileges required.');
            }

            API.setToken(data.token);
            App.showApp(data.user);
        } catch (error) {
            errorEl.textContent = error.message;
        } finally {
            btn.textContent = 'Sign In';
            btn.disabled = false;
        }
    },

    async checkSession() {
        const token = API.getToken();
        if (!token) return false;

        try {
            const data = await API.getMe();
            if (data.user && data.user.is_admin) {
                App.showApp(data.user);
                return true;
            }
        } catch {
            API.clearToken();
        }
        return false;
    },

    logout() {
        API.clearToken();
        App.showLogin();
    }
};
