// Проверка авторизации
const userId = localStorage.getItem('user_id');
if (!userId) {
    window.location.href = 'index.html';
}

// Получение параметров из URL
const urlParams = new URLSearchParams(window.location.search);
const articleId = urlParams.get('id');
const isNewArticle = urlParams.get('new') === 'true';

// DOM элементы
const articleForm = document.getElementById('articleForm');
const titleInput = document.getElementById('title');
const contentInput = document.getElementById('content');
const deleteBtn = document.getElementById('deleteBtn');
const previewContent = document.getElementById('previewContent');
const tabButtons = document.querySelectorAll('.tab-button');
const tabPanes = document.querySelectorAll('.tab-pane');

// Информация о текущем пользователе
let currentUser = null;

// Настройка marked.js
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

// Загрузка информации о пользователе
async function loadUserInfo() {
    try {
        const response = await fetchWithAuth(`${API_URL}/users`);
        if (!response.ok) {
            throw new Error('Ошибка при загрузке информации о пользователе');
        }
        const data = await response.json();
        currentUser = data.user || data;
        if (data.user_id) {
            localStorage.setItem('user_id', data.user_id);
        }
    } catch (error) {
        console.error('Error loading user info:', error);
        alert('Ошибка при загрузке информации о пользователе');
        window.location.href = 'main.html';
    }
}

// Скрыть кнопку удаления для новой статьи
if (isNewArticle) {
    deleteBtn.style.display = 'none';
    document.title = 'GoBoard - Создание статьи';
    loadUserInfo();
} else {
    loadArticle();
}

// Загрузка существующей статьи
async function loadArticle() {
    try {
        const response = await fetchWithAuth(`${API_URL}/articles/${articleId}`);
        
        if (!response.ok) {
            throw new Error('Статья не найдена');
        }

        const article = await response.json();
        
        // Проверка прав на редактирование
        if (article.author_id !== parseInt(userId)) {
            alert('У вас нет прав на редактирование этой статьи');
            window.location.href = 'main.html';
            return;
        }

        titleInput.value = article.title;
        contentInput.value = article.content;
        updatePreview();
    } catch (error) {
        console.error('Error loading article:', error);
        alert('Ошибка при загрузке статьи');
        window.location.href = 'main.html';
    }
}

// Обработка отправки формы
articleForm.addEventListener('submit', async (e) => {
    e.preventDefault();

    const articleData = {
        title: titleInput.value,
        content: contentInput.value,
        author_id: parseInt(localStorage.getItem('user_id'))
    };

    try {
        let response;
        
        if (isNewArticle) {
            if (!currentUser) {
                throw new Error('Информация о пользователе не загружена');
            }

            articleData.author_name = currentUser.name;
            if (currentUser.photo_url) {
                articleData.author_photo_url = currentUser.photo_url;
            }
            
            response = await fetchWithAuth(`${API_URL}/articles/`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(articleData)
            });
        } else {
            response = await fetchWithAuth(`${API_URL}/articles/${articleId}`, {
                method: 'PATCH',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    title: articleData.title,
                    content: articleData.content
                })
            });
        }

        if (!response.ok) {
            throw new Error('Ошибка при сохранении статьи');
        }

        window.location.href = 'main.html';
    } catch (error) {
        console.error('Error saving article:', error);
        alert(error.message);
    }
});

// Обработка удаления статьи
deleteBtn.addEventListener('click', async () => {
    if (!confirm('Вы уверены, что хотите удалить эту статью?')) {
        return;
    }

    try {
        const response = await fetchWithAuth(`${API_URL}/articles/${articleId}`, {
            method: 'DELETE'
        });

        if (!response.ok) {
            throw new Error('Ошибка при удалении статьи');
        }

        window.location.href = 'main.html';
    } catch (error) {
        console.error('Error deleting article:', error);
        alert(error.message);
    }
});

// Обработка переключения вкладок
tabButtons.forEach(button => {
    button.addEventListener('click', () => {
        const tab = button.dataset.tab;

        tabButtons.forEach(btn => btn.classList.remove('active'));
        button.classList.add('active');

        tabPanes.forEach(pane => pane.classList.remove('active'));
        document.getElementById(`${tab}Tab`).classList.add('active');

        if (tab === 'preview') {
            updatePreview();
        }
    });
});

// Обновление предпросмотра при вводе
contentInput.addEventListener('input', updatePreview);

// Функция обновления предпросмотра
function updatePreview() {
    const markdown = contentInput.value;
    previewContent.innerHTML = marked.parse(markdown);
} 