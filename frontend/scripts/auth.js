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
            body: JSON.stringify({ email, password })
        });

        if (!response.ok) {
            throw new Error('Неверный email или пароль');
        }

        const data = await response.json();
        localStorage.setItem('user_id', data.user_id);
        window.location.href = 'main.html';
    } catch (error) {
        alert(error.message);
    }
}); 