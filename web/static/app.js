function loadFragment(fragmentName) {
    const xhr = new XMLHttpRequest();
    xhr.open('GET', `/static/fragments/${fragmentName}.html`, true);
    xhr.onload = function () {
        if (xhr.status >= 200 && xhr.status < 300) {
            document.getElementById('content').innerHTML = xhr.responseText;
        } else {
            document.getElementById('content').innerHTML = '<h1>404 page not found</h1>';
        }
    };
    xhr.send();
}

// Load default fragment (home)
document.addEventListener('DOMContentLoaded', function() {
    loadFragment('home');
});
