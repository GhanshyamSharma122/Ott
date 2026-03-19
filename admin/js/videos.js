// Videos management page
const Videos = {
    categories: [],

    async render() {
        const content = document.getElementById('page-content');
        content.innerHTML = '<div class="loading-spinner"></div>';

        try {
            // Load categories for the dropdown
            const catData = await API.getCategories();
            this.categories = catData.categories || [];

            const data = await API.getVideos();
            const videos = data.videos || [];

            content.innerHTML = `
                <div class="table-header">
                    <div class="search-box">
                        <span class="material-icons-round">search</span>
                        <input type="text" id="video-search" placeholder="Search videos...">
                    </div>
                    <button class="btn btn-primary" id="add-video-btn">
                        <span class="material-icons-round">add</span>
                        Add Video
                    </button>
                </div>
                ${videos.length === 0
                    ? `<div class="empty-state">
                        <span class="material-icons-round">movie</span>
                        <p>No videos yet. Add your first video!</p>
                       </div>`
                    : `<table class="data-table">
                        <thead>
                            <tr>
                                <th>Thumbnail</th>
                                <th>Title</th>
                                <th>Category</th>
                                <th>Type</th>
                                <th>Status</th>
                                <th>Views</th>
                                <th>Actions</th>
                            </tr>
                        </thead>
                        <tbody>
                            ${videos.map(v => `
                                <tr>
                                    <td><img class="table-thumb" src="${v.custom_thumbnail_url || v.thumbnail_url || `https://img.youtube.com/vi/${v.youtube_video_id}/mqdefault.jpg`}" alt="" onerror="this.src='https://via.placeholder.com/64x36?text=No+Img'"></td>
                                    <td>${this.escapeHtml(v.title)}</td>
                                    <td>${v.category ? this.escapeHtml(v.category.name) : '—'}</td>
                                    <td><span class="badge ${v.is_premium ? 'badge-premium' : 'badge-free'}">${v.is_premium ? 'Premium' : 'Free'}</span></td>
                                    <td><span class="badge ${v.is_published ? 'badge-published' : 'badge-draft'}">${v.is_published ? 'Published' : 'Draft'}</span></td>
                                    <td>${v.view_count.toLocaleString()}</td>
                                    <td class="actions-cell">
                                        <button class="btn-icon" onclick="Videos.editVideo(${v.id})" title="Edit">
                                            <span class="material-icons-round">edit</span>
                                        </button>
                                        <button class="btn-icon" onclick="Videos.deleteVideo(${v.id}, '${this.escapeHtml(v.title)}')" title="Delete">
                                            <span class="material-icons-round">delete</span>
                                        </button>
                                    </td>
                                </tr>
                            `).join('')}
                        </tbody>
                    </table>`
                }
            `;

            // Event listeners
            document.getElementById('add-video-btn').addEventListener('click', () => this.showAddModal());

            const searchInput = document.getElementById('video-search');
            let searchTimeout;
            searchInput.addEventListener('input', (e) => {
                clearTimeout(searchTimeout);
                searchTimeout = setTimeout(() => this.searchVideos(e.target.value), 300);
            });

        } catch (error) {
            content.innerHTML = `<div class="empty-state"><span class="material-icons-round">error</span><p>${error.message}</p></div>`;
        }
    },

    async searchVideos(query) {
        try {
            const data = await API.getVideos(query);
            const videos = data.videos || [];
            const tbody = document.querySelector('.data-table tbody');
            if (!tbody) return;

            if (videos.length === 0) {
                tbody.innerHTML = `<tr><td colspan="7" style="text-align:center;color:var(--text-secondary);padding:40px">No videos found</td></tr>`;
                return;
            }

            tbody.innerHTML = videos.map(v => `
                <tr>
                    <td><img class="table-thumb" src="${v.custom_thumbnail_url || v.thumbnail_url || `https://img.youtube.com/vi/${v.youtube_video_id}/mqdefault.jpg`}" alt="" onerror="this.src='https://via.placeholder.com/64x36?text=No+Img'"></td>
                    <td>${this.escapeHtml(v.title)}</td>
                    <td>${v.category ? this.escapeHtml(v.category.name) : '—'}</td>
                    <td><span class="badge ${v.is_premium ? 'badge-premium' : 'badge-free'}">${v.is_premium ? 'Premium' : 'Free'}</span></td>
                    <td><span class="badge ${v.is_published ? 'badge-published' : 'badge-draft'}">${v.is_published ? 'Published' : 'Draft'}</span></td>
                    <td>${v.view_count.toLocaleString()}</td>
                    <td class="actions-cell">
                        <button class="btn-icon" onclick="Videos.editVideo(${v.id})" title="Edit"><span class="material-icons-round">edit</span></button>
                        <button class="btn-icon" onclick="Videos.deleteVideo(${v.id}, '${this.escapeHtml(v.title)}')" title="Delete"><span class="material-icons-round">delete</span></button>
                    </td>
                </tr>
            `).join('');
        } catch (error) {
            App.showToast(error.message, 'error');
        }
    },

    showAddModal() {
        const categoryOptions = this.categories.map(c => `<option value="${c.id}">${this.escapeHtml(c.name)}</option>`).join('');

        App.showModal('Add New Video', `
            <form id="video-form">
                <div class="form-group">
                    <label>Title *</label>
                    <input type="text" id="video-title" required>
                </div>
                <div class="form-group">
                    <label>Description</label>
                    <textarea id="video-description"></textarea>
                </div>
                <div class="form-group">
                    <label>YouTube URL *</label>
                    <input type="text" id="video-youtube-url" placeholder="https://www.youtube.com/watch?v=..." required>
                </div>
                <div class="form-group">
                    <label>Custom Thumbnail URL (optional)</label>
                    <input type="text" id="video-custom-thumb" placeholder="Leave empty to use YouTube thumbnail">
                </div>
                <div class="form-group">
                    <label>Category</label>
                    <select id="video-category">
                        <option value="">— Select Category —</option>
                        ${categoryOptions}
                    </select>
                </div>
                <div class="form-row">
                    <div class="form-group">
                        <label>Duration</label>
                        <input type="text" id="video-duration" placeholder="e.g. 2h 30m">
                    </div>
                    <div class="form-group">
                        <label>Release Date</label>
                        <input type="date" id="video-release-date">
                    </div>
                </div>
                <div class="form-row">
                    <div class="checkbox-group">
                        <input type="checkbox" id="video-is-premium">
                        <label for="video-is-premium">Premium Content</label>
                    </div>
                    <div class="checkbox-group">
                        <input type="checkbox" id="video-is-published" checked>
                        <label for="video-is-published">Published</label>
                    </div>
                </div>
                <div class="modal-actions">
                    <button type="button" class="btn btn-secondary" onclick="App.closeModal()">Cancel</button>
                    <button type="submit" class="btn btn-primary">Add Video</button>
                </div>
            </form>
        `);

        document.getElementById('video-form').addEventListener('submit', async (e) => {
            e.preventDefault();
            await this.handleAddVideo();
        });
    },

    async handleAddVideo() {
        const data = {
            title: document.getElementById('video-title').value,
            description: document.getElementById('video-description').value,
            youtube_url: document.getElementById('video-youtube-url').value,
            custom_thumbnail_url: document.getElementById('video-custom-thumb').value,
            category_id: parseInt(document.getElementById('video-category').value) || 0,
            duration: document.getElementById('video-duration').value,
            release_date: document.getElementById('video-release-date').value,
            is_premium: document.getElementById('video-is-premium').checked,
            is_published: document.getElementById('video-is-published').checked,
        };

        try {
            await API.createVideo(data);
            App.closeModal();
            App.showToast('Video added successfully!', 'success');
            this.render();
        } catch (error) {
            App.showToast(error.message, 'error');
        }
    },

    async editVideo(id) {
        try {
            const data = await API.request('GET', `/videos/${id}`);
            const video = data.video;
            const categoryOptions = this.categories.map(c =>
                `<option value="${c.id}" ${c.id === video.category_id ? 'selected' : ''}>${this.escapeHtml(c.name)}</option>`
            ).join('');

            App.showModal('Edit Video', `
                <form id="edit-video-form">
                    <div class="form-group">
                        <label>Title *</label>
                        <input type="text" id="edit-video-title" value="${this.escapeHtml(video.title)}" required>
                    </div>
                    <div class="form-group">
                        <label>Description</label>
                        <textarea id="edit-video-description">${this.escapeHtml(video.description || '')}</textarea>
                    </div>
                    <div class="form-group">
                        <label>YouTube URL</label>
                        <input type="text" id="edit-video-youtube-url" value="https://www.youtube.com/watch?v=${video.youtube_video_id}">
                    </div>
                    <div class="form-group">
                        <label>Custom Thumbnail URL</label>
                        <input type="text" id="edit-video-custom-thumb" value="${this.escapeHtml(video.custom_thumbnail_url || '')}">
                    </div>
                    <div class="form-group">
                        <label>Category</label>
                        <select id="edit-video-category">
                            <option value="">— Select Category —</option>
                            ${categoryOptions}
                        </select>
                    </div>
                    <div class="form-row">
                        <div class="form-group">
                            <label>Duration</label>
                            <input type="text" id="edit-video-duration" value="${video.duration || ''}">
                        </div>
                        <div class="form-group">
                            <label>Release Date</label>
                            <input type="date" id="edit-video-release-date" value="${video.release_date || ''}">
                        </div>
                    </div>
                    <div class="form-row">
                        <div class="checkbox-group">
                            <input type="checkbox" id="edit-video-is-premium" ${video.is_premium ? 'checked' : ''}>
                            <label for="edit-video-is-premium">Premium Content</label>
                        </div>
                        <div class="checkbox-group">
                            <input type="checkbox" id="edit-video-is-published" ${video.is_published ? 'checked' : ''}>
                            <label for="edit-video-is-published">Published</label>
                        </div>
                    </div>
                    <div class="modal-actions">
                        <button type="button" class="btn btn-secondary" onclick="App.closeModal()">Cancel</button>
                        <button type="submit" class="btn btn-primary">Save Changes</button>
                    </div>
                </form>
            `);

            document.getElementById('edit-video-form').addEventListener('submit', async (e) => {
                e.preventDefault();
                const updateData = {
                    title: document.getElementById('edit-video-title').value,
                    description: document.getElementById('edit-video-description').value,
                    youtube_url: document.getElementById('edit-video-youtube-url').value,
                    custom_thumbnail_url: document.getElementById('edit-video-custom-thumb').value,
                    category_id: parseInt(document.getElementById('edit-video-category').value) || 0,
                    duration: document.getElementById('edit-video-duration').value,
                    release_date: document.getElementById('edit-video-release-date').value,
                    is_premium: document.getElementById('edit-video-is-premium').checked,
                    is_published: document.getElementById('edit-video-is-published').checked,
                };

                try {
                    await API.updateVideo(id, updateData);
                    App.closeModal();
                    App.showToast('Video updated successfully!', 'success');
                    this.render();
                } catch (error) {
                    App.showToast(error.message, 'error');
                }
            });
        } catch (error) {
            App.showToast(error.message, 'error');
        }
    },

    async deleteVideo(id, title) {
        if (!confirm(`Are you sure you want to delete "${title}"?`)) return;

        try {
            await API.deleteVideo(id);
            App.showToast('Video deleted successfully!', 'success');
            this.render();
        } catch (error) {
            App.showToast(error.message, 'error');
        }
    },

    escapeHtml(text) {
        if (!text) return '';
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }
};
