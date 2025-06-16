const API_URL = 'http://localhost:8080';

document.getElementById('loginForm').addEventListener('submit', async (e) => {
    e.preventDefault();

    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;

    try {
        const response = await fetch(`${API_URL}/auth/login`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                email: email,
                password: password
            }),
            credentials: 'include'
        });

        if (!response.ok) {
            throw new Error('Неверный email или пароль');
        }

        const data = await response.json();
        sessionStorage.setItem('access_token', data.access_token);

        // Получаем данные пользователя сразу после успешной авторизации
        await getUserData();
        window.location.href = 'main.html';
    } catch (error) {
        alert(error.message);
    }
});

async function getUserData() {
    try {
        const response = await fetch(`${API_URL}/users`, {
            headers: {
                'Authorization': 'Bearer ' + sessionStorage.getItem('access_token')
            },
            credentials: 'include'
        });

        if (!response.ok) {
            throw new Error('Failed to fetch user data');
        }

        const data = await response.json();
        // Сохраняем данные пользователя
        localStorage.setItem('user_id', data.user_id);
        localStorage.setItem('user_name', data.user.name);
        localStorage.setItem('user_photo_url', data.user.photo_url || '');

    } catch (error) {
        alert('Error fetching user data');
    }
}