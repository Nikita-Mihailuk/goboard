// Проверка авторизации
const userId = localStorage.getItem('user_id');
if (!userId) {
    window.location.href = 'index.html';
}

// DOM элементы
const mainContent = document.getElementById('mainContent');
const profileContent = document.getElementById('profileContent');
const mainTab = document.getElementById('mainTab');
const profileTab = document.getElementById('profileTab');
const logoutBtn = document.getElementById('logoutBtn');
const searchInput = document.getElementById('searchInput');
const sortSelect = document.getElementById('sortSelect');
const articlesList = document.getElementById('articlesList');
const userArticlesList = document.getElementById('userArticlesList');
const userPhoto = document.getElementById('userPhoto');
const userName = document.getElementById('userName');
const userEmail = document.getElementById('userEmail');
const editProfileBtn = document.getElementById('editProfileBtn');
const createArticleBtn = document.getElementById('createArticleBtn');
const editProfileModal = document.getElementById('editProfileModal');
const editProfileForm = document.getElementById('editProfileForm');
const closeEditProfileModal = document.getElementById('closeEditProfileModal');
const prevPageBtn = document.getElementById('prevPage');
const nextPageBtn = document.getElementById('nextPage');
const pageInfo = document.getElementById('pageInfo');

// Состояние приложения
let currentPage = 1;
const articlesPerPage = 9;
let totalArticles = 0;
let currentTab = 'main';
let allArticles = [];

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
        // Сохраняем user_id в localStorage
        if (data.user_id) {
            localStorage.setItem('user_id', data.user_id);
        }
        const user = data.user || data;
        userName.textContent = user.name;
        userEmail.textContent = user.email;
        userPhoto.src = user.photo_url ? `${API_URL}/${user.photo_url}` : 'images/default-avatar.jpg';
    } catch (error) {
        console.error('Error loading user info:', error);
        alert('Ошибка при загрузке информации о пользователе');
    }
}

// Создание карточки статьи
function createArticleCard(article) {
    const card = document.createElement('div');
    card.className = 'article-card';
    card.onclick = () => {
        window.location.href = `article.html?id=${article.id}&returnTab=${currentTab}`;
    };

    const title = document.createElement('h3');
    title.textContent = article.title;

    const meta = document.createElement('div');
    meta.className = 'article-meta';
    meta.textContent = `${article.author_name} • ${new Date(article.created_at).toLocaleDateString()}`;

    card.appendChild(title);
    card.appendChild(meta);

    return card;
}

// Загрузка статей
async function loadArticles() {
    try {
        if (allArticles.length === 0) {
            const response = await fetchWithAuth(`${API_URL}/articles`);
            if (!response.ok) {
                throw new Error('Ошибка при загрузке статей');
            }
            allArticles = await response.json();
        }

        let filteredArticles = allArticles;

        if (searchInput.value) {
            const searchQuery = searchInput.value.toLowerCase();
            filteredArticles = filteredArticles.filter(article => 
                article.title.toLowerCase().includes(searchQuery) || 
                article.author_name.toLowerCase().includes(searchQuery)
            );
        }

        filteredArticles.sort((a, b) => {
            if (sortSelect.value === 'newest') {
                return new Date(b.created_at) - new Date(a.created_at);
            } else {
                return new Date(a.created_at) - new Date(b.created_at);
            }
        });

        totalArticles = filteredArticles.length;
        const startIndex = (currentPage - 1) * articlesPerPage;
        const endIndex = startIndex + articlesPerPage;
        const pageArticles = filteredArticles.slice(startIndex, endIndex);

        articlesList.innerHTML = '';
        pageArticles.forEach(article => {
            articlesList.appendChild(createArticleCard(article));
        });

        const isLastPage = endIndex >= filteredArticles.length;
        updatePagination(isLastPage);
    } catch (error) {
        console.error('Error loading articles:', error);
        alert('Ошибка при загрузке статей');
    }
}

// Загрузка статей пользователя
async function loadUserArticles() {
    try {
        const response = await fetchWithAuth(`${API_URL}/articles/author`);
        if (!response.ok) {
            console.log(response)
            throw new Error('Ошибка при загрузке статей пользователя');
        }
        const data = await response.json();
        
        userArticlesList.innerHTML = '';
        data.forEach(article => {
            userArticlesList.appendChild(createArticleCard(article));
        });
    } catch (error) {
        console.error('Error loading user articles:', error);
        alert('Ошибка при загрузке статей пользователя');
    }
}

// Обновление пагинации
function updatePagination(isLastPage) {
    const totalPages = Math.ceil(totalArticles / articlesPerPage);
    pageInfo.textContent = `Страница ${currentPage} из ${totalPages}`;
    prevPageBtn.disabled = currentPage === 1;
    nextPageBtn.disabled = isLastPage;
}

// Обработчики событий
mainTab.onclick = (e) => {
    e.preventDefault();
    currentTab = 'main';
    mainContent.classList.remove('hidden');
    profileContent.classList.add('hidden');
    mainTab.classList.add('active');
    profileTab.classList.remove('active');
    loadArticles();
};

profileTab.onclick = (e) => {
    e.preventDefault();
    currentTab = 'profile';
    mainContent.classList.add('hidden');
    profileContent.classList.remove('hidden');
    mainTab.classList.remove('active');
    profileTab.classList.add('active');
    loadUserArticles();
};

logoutBtn.onclick = async (e) => {
    e.preventDefault();
    try {
        await fetchWithAuth(`${API_URL}/auth/logout`, { method: 'DELETE' });
    } catch {}
    localStorage.removeItem('user_id');
    sessionStorage.removeItem('access_token');
    window.location.href = 'index.html';
};

searchInput.oninput = () => {
    currentPage = 1;
    loadArticles();
};

sortSelect.onchange = () => {
    currentPage = 1;
    loadArticles();
};

prevPageBtn.onclick = () => {
    if (currentPage > 1) {
        currentPage--;
        loadArticles();
    }
};

nextPageBtn.onclick = () => {
    currentPage++;
    loadArticles();
};

editProfileBtn.onclick = () => {
    editProfileModal.classList.add('active');
};

closeEditProfileModal.onclick = () => {
    editProfileModal.classList.remove('active');
};

createArticleBtn.onclick = () => {
    window.location.href = 'article-edit.html?new=true';
};

editProfileForm.onsubmit = async (e) => {
    e.preventDefault();
    
    const formData = new FormData(editProfileForm);
    
    try {
        const response = await fetchWithAuth(`${API_URL}/users`, {
            method: 'PATCH',
            body: formData
        });

        if (!response.ok) {
            throw new Error('Неверный текущий пароль');
        }

        editProfileModal.classList.remove('active');
        loadUserInfo();
    } catch (error) {
        console.error('Error updating profile:', error);
        alert(error.message);
    }
};

// Проверка хэша URL для определения начальной вкладки
if (window.location.hash === '#profile') {
    profileTab.click();
} else {
    mainTab.click();
}

// Инициализация
loadUserInfo();