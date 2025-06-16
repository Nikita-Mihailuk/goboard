const urlParams = new URLSearchParams(window.location.search);
const articleId = urlParams.get('id');
const returnTab = urlParams.get('returnTab') || 'main';
const userId = localStorage.getItem('user_id');
const userName = localStorage.getItem('user_name');
const userPhotoUrl = localStorage.getItem('user_photo_url');

if (!articleId) {
    window.location.href = 'main.html';
}

const articleTitle = document.getElementById('articleTitle');
const articleContent = document.getElementById('articleContent');
const authorName = document.getElementById('authorName');
const authorPhoto = document.getElementById('authorPhoto');
const articleDate = document.getElementById('articleDate');
const returnButton = document.querySelector('nav button');
const nav = document.querySelector('nav');

const commentsList = document.getElementById('commentsList');
const commentTemplate = document.getElementById('commentTemplate');
const commentText = document.getElementById('commentText');
const submitButton = document.getElementById('submitComment');

returnButton.onclick = () => {
    window.location.href = `main.html${returnTab === 'profile' ? '#profile' : ''}`;
};

marked.setOptions({
    breaks: true,
    gfm: true,
    headerIds: true,
    mangle: false,
    sanitize: false
});

const API_URL = 'http://localhost:8080';

async function fetchWithAuth(url, options = {}) {
    let accessToken = sessionStorage.getItem('access_token');
    let headers = options.headers ? { ...options.headers } : {};
    if (accessToken) {
        headers['Authorization'] = 'Bearer ' + accessToken;
    }
    headers['Cache-Control'] = 'no-cache, no-store, must-revalidate';
    options.headers = headers;
    options.credentials = 'include';

    let response = await fetch(url, options);
    if (response.status === 401) {
        // Пробуем обновить токен
        const refreshResponse = await fetch(`${API_URL}/auth/refresh`, { method: 'POST', credentials: 'include' });
        if (refreshResponse.ok) {
            const refreshData = await refreshResponse.json();
            sessionStorage.setItem('access_token', refreshData.access_token);
            headers['Authorization'] = 'Bearer ' + refreshData.access_token;
            response = await fetch(url, { ...options, headers });
        }
    }
    return response;
}

// Загрузка статьи
async function loadArticle() {
    try {
        const response = await fetchWithAuth(`${API_URL}/articles/${articleId}`);
        
        if (!response.ok) {
            throw new Error('Статья не найдена');
        }

        const article = await response.json();

        document.title = `GoBoard - ${article.title}`;
        articleTitle.textContent = article.title;
        articleContent.innerHTML = marked.parse(article.content);
        authorName.textContent = article.author_name;
        authorPhoto.src = article.author_photo_url ? `${API_URL}/${article.author_photo_url}` : 'images/default-avatar.jpg';
        articleDate.textContent = new Date(article.created_at).toLocaleDateString();

        // Проверяем, является ли текущий пользователь автором статьи
        const isAuthor = article.author_id === parseInt(localStorage.getItem('user_id'));

        const buttons = nav.querySelectorAll('button:not(:first-child)');
        buttons.forEach(button => button.remove());

        // Если пользователь является автором, показываем кнопки редактирования и удаления
        if (isAuthor) {
            const editButton = document.createElement('button');
            editButton.className = 'btn-primary';
            editButton.textContent = 'Редактировать';
            editButton.onclick = () => {
                window.location.href = `article-edit.html?id=${articleId}`;
            };

            nav.appendChild(editButton);
        }
    } catch (error) {
        console.error('Error loading article:', error);
        alert('Ошибка при загрузке статьи');
        window.location.href = 'main.html';
    }
}

async function loadComments() {
    try {
        const response = await fetchWithAuth(`${API_URL}/comments/article/${articleId}`);
        if (!response.ok) throw new Error('Failed to load comments');
        
        const comments = await response.json();
        renderComments(comments);
    } catch (error) {
        console.error('Error loading comments:', error);
    }
}

function renderComments(comments) {
    commentsList.innerHTML = '';

    const sortOldestFirst = (a, b) => new Date(a.created_at) - new Date(b.created_at);
    
    const sortedComments = comments.sort(sortOldestFirst);

    const topLevelComments = sortedComments.filter(c => !c.parent_id);
    const replies = sortedComments.filter(c => c.parent_id);

    topLevelComments.forEach(comment => {
        const commentElement = createCommentElement(comment);
        if (!commentElement) return;

        const commentReplies = replies
            .filter(reply => reply.parent_id === comment.id)
            .sort(sortOldestFirst);
        
        const repliesContainer = commentElement.querySelector('.replies');
        commentReplies.forEach(reply => {
            const replyElement = createCommentElement(reply);
            if (replyElement) {
                repliesContainer.appendChild(replyElement);
            }
        });
        
        commentsList.appendChild(commentElement);
    });
}

function createCommentElement(comment) {
    if (!comment) return null;

    const template = commentTemplate.content.cloneNode(true);
    const commentElement = template.querySelector('.comment');
    
    commentElement.dataset.commentId = comment.id;
    commentElement.dataset.rootId = comment.parent_id || comment.id;

    const authorPhoto = commentElement.querySelector('.author-photo');
    if (!authorPhoto) return null;

    const photoUrl = comment.author_photo_url ? `${API_URL}/${comment.author_photo_url}` : 'images/default-avatar.jpg';
    authorPhoto.src = photoUrl;

    const authorName = commentElement.querySelector('.author-name');
    if (!authorName) return null;
    authorName.textContent = comment.author_name || 'Неизвестный пользователь';

    const contentElement = commentElement.querySelector('.comment-content');
    if (!contentElement) return null;
    contentElement.textContent = comment.content;

    const dateElement = commentElement.querySelector('.comment-date');
    if (!dateElement) return null;
    const date = new Date(comment.created_at);
    dateElement.textContent = date.toLocaleString();

    if (comment.updated_at && comment.updated_at !== comment.created_at) {
        const editedMarkElement = commentElement.querySelector('.edited-mark');
        if (editedMarkElement) {
            const editedDate = new Date(comment.updated_at);
            editedMarkElement.textContent = `Изменено ${editedDate.toLocaleString()}`;
        }
    }

    const menuContainer = commentElement.querySelector('.comment-menu');
    if (menuContainer) {
        if (comment.author_id === parseInt(userId)) {
            menuContainer.style.display = 'block';
        } else {
            menuContainer.style.display = 'none';
        }
    }

    const replyForm = commentElement.querySelector('.reply-form');
    if (!replyForm) return null;

    const repliesContainer = commentElement.querySelector('.replies');
    if (!repliesContainer) return null;
    
    return commentElement;
}

function toggleReplyForm(commentElement) {
    const replyForm = commentElement.querySelector('.reply-form');
    const isVisible = replyForm.style.display === 'block';

    document.querySelectorAll('.reply-form').forEach(form => {
        form.style.display = 'none';
    });

    replyForm.style.display = isVisible ? 'none' : 'block';
    
    if (!isVisible) {
        replyForm.querySelector('textarea').focus();
    }
}

async function handleNewComment(parentId = null, rootId = null) {
    const textArea = parentId ? 
        document.querySelector(`[data-comment-id="${parentId}"] .reply-form textarea`) :
        commentText;
    
    if (!textArea) {
        console.error('Could not find textarea');
        return;
    }

    const content = textArea.value.trim();
    if (!content) return;

    try {
        const response = await fetchWithAuth(`${API_URL}/comments`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                article_id: articleId,
                parent_id: rootId || parentId,
                content: content,
                author_id: parseInt(userId),
                author_name: userName,
                author_photo_url: userPhotoUrl || null
            })
        });

        if (!response.ok) throw new Error('Failed to create comment');

        textArea.value = '';
        if (parentId) {
            const commentElement = document.querySelector(`[data-comment-id="${parentId}"]`);
            if (commentElement) {
                toggleReplyForm(commentElement);
            }
        }
        await loadComments();

    } catch (error) {
        console.error('Error creating comment:', error);
        alert('Ошибка при создании комментария');
    }
}

async function handleEditComment(commentElement) {
    const contentElement = commentElement.querySelector('.comment-content');
    const currentContent = contentElement.textContent;

    contentElement.innerHTML = `
        <textarea class="edit-textarea">${currentContent}</textarea>
        <div class="edit-actions">
            <button class="save-edit btn-primary">Сохранить</button>
            <button class="cancel-edit">Отменить</button>
        </div>
    `;

    const textarea = contentElement.querySelector('textarea');
    const saveBtn = contentElement.querySelector('.save-edit');
    const cancelBtn = contentElement.querySelector('.cancel-edit');

    saveBtn.addEventListener('click', async () => {
        const newContent = textarea.value.trim();
        if (!newContent) return;

        try {
            const response = await fetchWithAuth(`${API_URL}/comments/${commentElement.dataset.commentId}`, {
                method: 'PATCH',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    content: newContent
                })
            });

            if (!response.ok) throw new Error('Failed to update comment');

            await loadComments();
        } catch (error) {
            console.error('Error updating comment:', error);
            alert('Ошибка при обновлении комментария');
        }
    });

    cancelBtn.addEventListener('click', () => {
        contentElement.textContent = currentContent;
    });
}

async function handleDeleteComment(commentElement) {
    if (!confirm('Вы уверены, что хотите удалить этот комментарий?')) return;

    try {
        const response = await fetchWithAuth(`${API_URL}/comments/${commentElement.dataset.commentId}`, {
            method: 'DELETE'
        });

        if (!response.ok) throw new Error('Failed to delete comment');

        await loadComments();
    } catch (error) {
        console.error('Error deleting comment:', error);
        alert('Ошибка при удалении комментария');
    }
}

function toggleCommentMenu(menuToggle) {
    const dropdown = menuToggle.nextElementSibling;
    dropdown.classList.toggle('active');
}

function setupCommentEventListeners() {
    submitButton.addEventListener('click', () => handleNewComment());

    commentsList.addEventListener('click', (e) => {
        const target = e.target;
        const comment = target.closest('.comment');
        
        if (!comment) return;

        if (target.classList.contains('menu-toggle')) {
            toggleCommentMenu(target);
        } else if (target.classList.contains('edit-comment')) {
            handleEditComment(comment);
        } else if (target.classList.contains('delete-comment')) {
            handleDeleteComment(comment);
        } else if (target.classList.contains('reply-button')) {
            toggleReplyForm(comment);
        } else if (target.classList.contains('submit-reply')) {
            const rootId = comment.dataset.rootId;
            handleNewComment(comment.dataset.commentId, rootId);
        } else if (target.classList.contains('cancel-reply')) {
            toggleReplyForm(comment);
        }
    });

    document.addEventListener('click', (e) => {
        if (!e.target.classList.contains('menu-toggle')) {
            document.querySelectorAll('.menu-dropdown.active').forEach(menu => {
                menu.classList.remove('active');
            });
        }
    });
}

document.addEventListener('DOMContentLoaded', () => {
    loadArticle();
    loadComments();
    setupCommentEventListeners();
}); 