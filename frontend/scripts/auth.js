document.getElementById('loginForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;

    try {
        const response = await fetch('http://localhost:8080/auth/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ email, password }),
            credentials: 'include'
        });

        if (!response.ok) {
            throw new Error('Неверный email или пароль');
        }

        const data = await response.json();
        sessionStorage.setItem('access_token', data.access_token);
        // Получаем user_id через защищённый эндпоинт
        const userResp = await fetch('http://localhost:8080/users', {
            headers: {
                'Authorization': 'Bearer ' + data.access_token
            },
            credentials: 'include'
        });
        if (!userResp.ok) {
            throw new Error('Ошибка при получении профиля пользователя');
        }
        const userData = await userResp.json();
        localStorage.setItem('user_id', userData.user_id);
        window.location.href = 'main.html';
    } catch (error) {
        alert(error.message);
    }
}); 