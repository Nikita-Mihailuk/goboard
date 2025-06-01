document.getElementById('registerForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const name = document.getElementById('name').value;
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;
    const confirmPassword = document.getElementById('confirmPassword').value;
    const errorElement = document.getElementById('passwordError');

    // Проверка совпадения паролей
    if (password !== confirmPassword) {
        errorElement.textContent = 'Пароли не совпадают';
        return;
    }
    errorElement.textContent = '';

    try {
        const response = await fetch('http://localhost:8080/auth/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ name, email, password })
        });

        if (!response.ok) {
            throw new Error('Ошибка при регистрации');
        }

        alert('Регистрация успешна! Теперь вы можете войти.');
        window.location.href = 'index.html';
    } catch (error) {
        alert(error.message);
    }
}); 