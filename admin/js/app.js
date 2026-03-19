// Main App router & initialization
const App = {
    currentPage: 'dashboard',

    async init() {
        // Check for existing session
        const hasSession = await Auth.checkSession();
        if (!hasSession) {
            this.showLogin();
        }

        // Setup auth
        Auth.init();

        // Logout
        document.getElementById('logout-btn').addEventListener('click', () => Auth.logout());

        // Navigation
        document.querySelectorAll('.nav-item').forEach(item => {
            item.addEventListener('click', (e) => {
                e.preventDefault();
                const page = item.dataset.page;
                this.navigateTo(page);
            });
        });

        // Modal close
        document.getElementById('modal-close').addEventListener('click', () => this.closeModal());
        document.getElementById('modal-overlay').addEventListener('click', (e) => {
            if (e.target === e.currentTarget) this.closeModal();
        });
    },

    showLogin() {
        document.getElementById('login-screen').classList.remove('hidden');
        document.getElementById('app').classList.add('hidden');
    },

    showApp(user) {
        document.getElementById('login-screen').classList.add('hidden');
        document.getElementById('app').classList.remove('hidden');
        document.getElementById('admin-name').textContent = user.name || 'Admin';
        this.navigateTo('dashboard');
    },

    navigateTo(page) {
        this.currentPage = page;

        // Update active nav
        document.querySelectorAll('.nav-item').forEach(item => {
            item.classList.toggle('active', item.dataset.page === page);
        });

        // Update title
        const titles = {
            dashboard: 'Dashboard',
            videos: 'Video Management',
            categories: 'Categories',
            users: 'User Management'
        };
        document.getElementById('page-title').textContent = titles[page] || page;

        // Render page
        switch (page) {
            case 'dashboard': Dashboard.render(); break;
            case 'videos': Videos.render(); break;
            case 'categories': Categories.render(); break;
            case 'users': Users.render(); break;
        }
    },

    showModal(title, bodyHtml) {
        document.getElementById('modal-title').textContent = title;
        document.getElementById('modal-body').innerHTML = bodyHtml;
        document.getElementById('modal-overlay').classList.remove('hidden');
    },

    closeModal() {
        document.getElementById('modal-overlay').classList.add('hidden');
    },

    showToast(message, type = 'success') {
        const toast = document.getElementById('toast');
        const msgEl = document.getElementById('toast-message');
        msgEl.textContent = message;
        toast.className = `toast toast-${type}`;
        setTimeout(() => { toast.classList.add('hidden'); }, 3000);
    }
};

// Initialize on load
document.addEventListener('DOMContentLoaded', () => App.init());
