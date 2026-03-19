// Dashboard page
const Dashboard = {
    async render() {
        const content = document.getElementById('page-content');
        content.innerHTML = '<div class="loading-spinner"></div>';

        try {
            const data = await API.getStats();
            content.innerHTML = `
                <div class="stats-grid">
                    <div class="stat-card">
                        <div class="stat-icon purple">
                            <span class="material-icons-round">people</span>
                        </div>
                        <div class="stat-value">${data.total_users}</div>
                        <div class="stat-label">Total Users</div>
                    </div>
                    <div class="stat-card">
                        <div class="stat-icon green">
                            <span class="material-icons-round">workspace_premium</span>
                        </div>
                        <div class="stat-value">${data.premium_users}</div>
                        <div class="stat-label">Premium Users</div>
                    </div>
                    <div class="stat-card">
                        <div class="stat-icon blue">
                            <span class="material-icons-round">movie</span>
                        </div>
                        <div class="stat-value">${data.total_videos}</div>
                        <div class="stat-label">Total Videos</div>
                    </div>
                    <div class="stat-card">
                        <div class="stat-icon orange">
                            <span class="material-icons-round">visibility</span>
                        </div>
                        <div class="stat-value">${data.total_views.toLocaleString()}</div>
                        <div class="stat-label">Total Views</div>
                    </div>
                </div>
                <div class="stats-grid">
                    <div class="stat-card">
                        <div class="stat-icon blue">
                            <span class="material-icons-round">category</span>
                        </div>
                        <div class="stat-value">${data.total_categories}</div>
                        <div class="stat-label">Categories</div>
                    </div>
                </div>
            `;
        } catch (error) {
            content.innerHTML = `
                <div class="empty-state">
                    <span class="material-icons-round">error</span>
                    <p>${error.message}</p>
                </div>
            `;
        }
    }
};
