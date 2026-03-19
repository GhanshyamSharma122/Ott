// API Client for Admin Panel
const API = {
    BASE_URL: 'http://localhost:8080/api',

    getToken() {
        return localStorage.getItem('admin_token');
    },

    setToken(token) {
        localStorage.setItem('admin_token', token);
    },

    clearToken() {
        localStorage.removeItem('admin_token');
    },

    async request(method, path, body = null) {
        const headers = {
            'Content-Type': 'application/json',
        };

        const token = this.getToken();
        if (token) {
            headers['Authorization'] = `Bearer ${token}`;
        }

        const options = { method, headers };
        if (body) {
            options.body = JSON.stringify(body);
        }

        try {
            const response = await fetch(`${this.BASE_URL}${path}`, options);
            const data = await response.json();

            if (!response.ok) {
                throw new Error(data.error || 'Request failed');
            }

            return data;
        } catch (error) {
            if (error.message === 'Failed to fetch') {
                throw new Error('Cannot connect to server. Make sure the backend is running.');
            }
            throw error;
        }
    },

    // Auth
    login(email, password) {
        return this.request('POST', '/auth/login', { email, password });
    },

    getMe() {
        return this.request('GET', '/auth/me');
    },

    // Admin Stats
    getStats() {
        return this.request('GET', '/admin/stats');
    },

    // Videos
    getVideos(search = '') {
        const query = search ? `?search=${encodeURIComponent(search)}` : '';
        return this.request('GET', `/admin/videos${query}`);
    },

    createVideo(data) {
        return this.request('POST', '/admin/videos', data);
    },

    updateVideo(id, data) {
        return this.request('PUT', `/admin/videos/${id}`, data);
    },

    deleteVideo(id) {
        return this.request('DELETE', `/admin/videos/${id}`);
    },

    // Categories
    getCategories() {
        return this.request('GET', '/categories');
    },

    createCategory(data) {
        return this.request('POST', '/admin/categories', data);
    },

    updateCategory(id, data) {
        return this.request('PUT', `/admin/categories/${id}`, data);
    },

    deleteCategory(id) {
        return this.request('DELETE', `/admin/categories/${id}`);
    },

    // Users
    getUsers(search = '') {
        const query = search ? `?search=${encodeURIComponent(search)}` : '';
        return this.request('GET', `/admin/users${query}`);
    },

    updateUser(id, data) {
        return this.request('PUT', `/admin/users/${id}`, data);
    },
};
