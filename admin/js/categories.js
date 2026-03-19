// Categories management page
const Categories = {
    async render() {
        const content = document.getElementById('page-content');
        content.innerHTML = '<div class="loading-spinner"></div>';

        try {
            const data = await API.getCategories();
            const categories = data.categories || [];

            content.innerHTML = `
                <div class="table-header">
                    <div></div>
                    <button class="btn btn-primary" id="add-category-btn">
                        <span class="material-icons-round">add</span>
                        Add Category
                    </button>
                </div>
                ${categories.length === 0
                    ? `<div class="empty-state">
                        <span class="material-icons-round">category</span>
                        <p>No categories yet. Add your first category!</p>
                       </div>`
                    : `<table class="data-table">
                        <thead>
                            <tr>
                                <th>Name</th>
                                <th>Description</th>
                                <th>Sort Order</th>
                                <th>Actions</th>
                            </tr>
                        </thead>
                        <tbody>
                            ${categories.map(c => `
                                <tr>
                                    <td><strong>${this.escapeHtml(c.name)}</strong></td>
                                    <td style="color:var(--text-secondary)">${this.escapeHtml(c.description || '—')}</td>
                                    <td>${c.sort_order}</td>
                                    <td class="actions-cell">
                                        <button class="btn-icon" onclick="Categories.editCategory(${c.id}, '${this.escapeHtml(c.name)}', '${this.escapeHtml(c.description || '')}', '${this.escapeHtml(c.image_url || '')}', ${c.sort_order})" title="Edit">
                                            <span class="material-icons-round">edit</span>
                                        </button>
                                        <button class="btn-icon" onclick="Categories.deleteCategory(${c.id}, '${this.escapeHtml(c.name)}')" title="Delete">
                                            <span class="material-icons-round">delete</span>
                                        </button>
                                    </td>
                                </tr>
                            `).join('')}
                        </tbody>
                    </table>`
                }
            `;

            document.getElementById('add-category-btn').addEventListener('click', () => this.showAddModal());
        } catch (error) {
            content.innerHTML = `<div class="empty-state"><span class="material-icons-round">error</span><p>${error.message}</p></div>`;
        }
    },

    showAddModal() {
        App.showModal('Add New Category', `
            <form id="category-form">
                <div class="form-group">
                    <label>Name *</label>
                    <input type="text" id="cat-name" required>
                </div>
                <div class="form-group">
                    <label>Description</label>
                    <textarea id="cat-description"></textarea>
                </div>
                <div class="form-group">
                    <label>Image URL</label>
                    <input type="text" id="cat-image-url" placeholder="Optional cover image URL">
                </div>
                <div class="form-group">
                    <label>Sort Order</label>
                    <input type="number" id="cat-sort-order" value="0">
                </div>
                <div class="modal-actions">
                    <button type="button" class="btn btn-secondary" onclick="App.closeModal()">Cancel</button>
                    <button type="submit" class="btn btn-primary">Add Category</button>
                </div>
            </form>
        `);

        document.getElementById('category-form').addEventListener('submit', async (e) => {
            e.preventDefault();
            const data = {
                name: document.getElementById('cat-name').value,
                description: document.getElementById('cat-description').value,
                image_url: document.getElementById('cat-image-url').value,
                sort_order: parseInt(document.getElementById('cat-sort-order').value) || 0,
            };

            try {
                await API.createCategory(data);
                App.closeModal();
                App.showToast('Category added!', 'success');
                this.render();
            } catch (error) {
                App.showToast(error.message, 'error');
            }
        });
    },

    editCategory(id, name, description, imageUrl, sortOrder) {
        App.showModal('Edit Category', `
            <form id="edit-category-form">
                <div class="form-group">
                    <label>Name *</label>
                    <input type="text" id="edit-cat-name" value="${this.escapeHtml(name)}" required>
                </div>
                <div class="form-group">
                    <label>Description</label>
                    <textarea id="edit-cat-description">${this.escapeHtml(description)}</textarea>
                </div>
                <div class="form-group">
                    <label>Image URL</label>
                    <input type="text" id="edit-cat-image-url" value="${this.escapeHtml(imageUrl)}">
                </div>
                <div class="form-group">
                    <label>Sort Order</label>
                    <input type="number" id="edit-cat-sort-order" value="${sortOrder}">
                </div>
                <div class="modal-actions">
                    <button type="button" class="btn btn-secondary" onclick="App.closeModal()">Cancel</button>
                    <button type="submit" class="btn btn-primary">Save Changes</button>
                </div>
            </form>
        `);

        document.getElementById('edit-category-form').addEventListener('submit', async (e) => {
            e.preventDefault();
            const data = {
                name: document.getElementById('edit-cat-name').value,
                description: document.getElementById('edit-cat-description').value,
                image_url: document.getElementById('edit-cat-image-url').value,
                sort_order: parseInt(document.getElementById('edit-cat-sort-order').value) || 0,
            };

            try {
                await API.updateCategory(id, data);
                App.closeModal();
                App.showToast('Category updated!', 'success');
                this.render();
            } catch (error) {
                App.showToast(error.message, 'error');
            }
        });
    },

    async deleteCategory(id, name) {
        if (!confirm(`Delete category "${name}"? Videos in this category won't be deleted.`)) return;

        try {
            await API.deleteCategory(id);
            App.showToast('Category deleted!', 'success');
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
