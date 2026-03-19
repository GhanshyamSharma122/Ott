// Users management page
const Users = {
    async render() {
        const content = document.getElementById('page-content');
        content.innerHTML = '<div class="loading-spinner"></div>';
        try {
            const data = await API.getUsers();
            const users = data.users || [];
            content.innerHTML = `
                <div class="table-header">
                    <div class="search-box">
                        <span class="material-icons-round">search</span>
                        <input type="text" id="user-search" placeholder="Search users...">
                    </div>
                </div>
                ${users.length === 0
                    ? '<div class="empty-state"><span class="material-icons-round">people</span><p>No users yet.</p></div>'
                    : this.renderTable(users)}
            `;
            const si = document.getElementById('user-search');
            let t;
            si.addEventListener('input', (e) => { clearTimeout(t); t = setTimeout(() => this.searchUsers(e.target.value), 300); });
        } catch (error) {
            content.innerHTML = `<div class="empty-state"><span class="material-icons-round">error</span><p>${error.message}</p></div>`;
        }
    },

    renderTable(users) {
        return `<table class="data-table"><thead><tr>
            <th>Name</th><th>Email</th><th>Subscription</th><th>Role</th><th>Joined</th><th>Actions</th>
        </tr></thead><tbody>${users.map(u => this.renderRow(u)).join('')}</tbody></table>`;
    },

    renderRow(u) {
        const esc = this.escapeHtml;
        return `<tr>
            <td><strong>${esc(u.name)}</strong></td>
            <td style="color:var(--text-secondary)">${esc(u.email)}</td>
            <td><span class="badge ${u.subscription_tier === 'premium' ? 'badge-premium' : 'badge-free'}">${u.subscription_tier}</span></td>
            <td>${u.is_admin ? '<span class="badge badge-admin">Admin</span>' : 'User'}</td>
            <td style="color:var(--text-secondary)">${new Date(u.created_at).toLocaleDateString()}</td>
            <td class="actions-cell">
                <button class="btn btn-sm ${u.subscription_tier === 'premium' ? 'btn-secondary' : 'btn-primary'}" 
                    onclick="Users.togglePremium(${u.id}, '${u.subscription_tier}')">
                    ${u.subscription_tier === 'premium' ? 'Revoke' : 'Grant'} Premium
                </button>
            </td>
        </tr>`;
    },

    async searchUsers(query) {
        try {
            const data = await API.getUsers(query);
            const tbody = document.querySelector('.data-table tbody');
            if (!tbody) return;
            const users = data.users || [];
            tbody.innerHTML = users.length === 0
                ? '<tr><td colspan="6" style="text-align:center;color:var(--text-secondary);padding:40px">No users found</td></tr>'
                : users.map(u => this.renderRow(u)).join('');
        } catch (error) { App.showToast(error.message, 'error'); }
    },

    async togglePremium(userId, currentTier) {
        const newTier = currentTier === 'premium' ? 'free' : 'premium';
        if (!confirm(`${newTier === 'premium' ? 'Grant' : 'Revoke'} premium for this user?`)) return;
        try {
            await API.updateUser(userId, { subscription_tier: newTier });
            App.showToast(`Premium ${newTier === 'premium' ? 'granted' : 'revoked'}!`, 'success');
            this.render();
        } catch (error) { App.showToast(error.message, 'error'); }
    },

    escapeHtml(text) {
        if (!text) return '';
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }
};
