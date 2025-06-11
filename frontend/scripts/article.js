const urlParams = new URLSearchParams(window.location.search);
const articleId = urlParams.get('id');
const returnTab = urlParams.get('returnTab') || 'main';
const userId = localStorage.getItem('user_id');

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

// Инициализация
loadArticle(); 