document.addEventListener('DOMContentLoaded', () => {
    function loadFragment(fragment) {
        fetch(`/fragments/${fragment}.html`)
            .then(response => response.text())
            .then(html => {
                document.getElementById('main-content').innerHTML = html;
            })
            .catch(error => console.error('Error loading fragment:', error));
    }

    function navigate(event) {
        if (event.target.tagName === 'A' && event.target.href) {
            event.preventDefault();
            const fragment = event.target.getAttribute('href').substring(1);
            loadFragment(fragment);
            history.pushState(null, '', `#${fragment}`);
        }
    }

    function handleLogin(event) {
        event.preventDefault();
        const email = document.getElementById('email').value;
        const password = document.getElementById('password').value;

        fetch('/api/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ email, password })
        })
        .then(response => response.json())
        .then(data => {
            if (data.error) {
                console.error('Login error:', data.error);
            } else {
                console.log('Login successful');
                loadFragment('home');
            }
        })
        .catch(error => console.error('Error during login:', error));
    }

    window.addEventListener('popstate', () => {
        const fragment = location.hash.substring(1) || 'home';
        loadFragment(fragment);
    });

    document.querySelectorAll('a.nav-link').forEach(link => {
        link.addEventListener('click', navigate);
    });

    document.getElementById('login-form')?.addEventListener('submit', handleLogin);

    const initialFragment = location.hash.substring(1) || 'home';
    loadFragment(initialFragment);
});
