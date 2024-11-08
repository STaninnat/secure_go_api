const API_BASE = '/v1';

async function fetchWithAlert(url, options = {}) {
    const response = await fetch(url, {
        ...options,
        credentials: 'include'
    });

    if (response.status === 401) {
        const refreshResponse = await refreshToken();
        if (refreshResponse && refreshResponse.ok) {
            return fetch(url, { ...options, credentials: 'include' });
        } else {
            alert("session expired. please log in again");
            window.location.href = "/";
            return;
        }
    }
    
    if (response.status > 299) {
        alert(`Error: ${response.status}`);
        return response;
    }
    return response;
}

function initLoginPage() {
    console.log("Initializing Login Page");

    async function loginUser() {
        const name = document.getElementById('nameField').value.trim();
        const password = document.getElementById('passwordField').value.trim();
        
        if (!name || !password) {
            alert("please enter both username and password");
            return;
        }
        
        const response = await fetchWithAlert(`${API_BASE}/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name, password })
        });

        if (response.ok) {
            alert(`Login successful, welcome ${name}`);
            window.location.href = "/static/posts.html"; 
        } else {
            alert("login failed. please check your credentials");
        }
    }

    window.loginUser = loginUser;
}

function initCreateUserPage() {
    console.log("Initializing Create User Page");

    async function createUser() {
        const name = document.getElementById('nameField').value.trim();
        const password = document.getElementById('passwordField').value.trim();
        
        if (!name || !password) {
            alert("please enter both username and password");
            return;
        }
        
        const response = await fetchWithAlert(`${API_BASE}/users`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name, password })
        });

        if (response.ok) {
            alert("User created successfully. WELCOME!");
            window.location.href = "posts.html";
        } else {
            alert("user creation failed");
        }
    }

    window.createUser = createUser;
}

function initPostPage() {
    console.log("Initializing Post Page");

    async function refreshToken() {
        const response = await fetchWithAlert(`${API_BASE}/refresh`, {
            method: 'POST',
            credentials: 'include'
        });

        if (response && response.ok) {
            console.log("Token refreshed successfully");
        } else {
            alert("failed to refresh token. please log in again");
            window.location.href = "/";
        }
        
        return response;
    }

    async function loadPosts() {
        const response = await fetchWithAlert(`${API_BASE}/posts`, {
            credentials: 'include'
        });

        if (response.ok) {
            const posts = await response.json();
            const postsContainer = document.getElementById('posts');
            postsContainer.innerHTML = '';
            posts.forEach(post => displayPost(post));
        } else {
            alert('error loading posts');
        }
    }

    async function createPost() {
        const postContent = document.getElementById('newPostContent').value;

        if (!postContent) {
            alert('please enter post content');
            return;
        }

        const response = await fetchWithAlert(`${API_BASE}/posts`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: 'include',
            body: JSON.stringify({ post: postContent })
        });

        if (response.ok) {
            const post = await response.json();
            displayPost(post);
            document.getElementById('newPostContent').value = '';
        } else {
            alert("failed to create post");
        }
    }

    function displayPost(post) {
        const postElement = document.createElement('div');
        postElement.className = 'post';
        postElement.textContent = post.post;
        document.getElementById('posts').appendChild(postElement);
    }

    async function logout() {
        const response = await fetchWithAlert(`${API_BASE}/logout`, {
            method: 'POST',
            credentials: 'include'
        });

        if (response.ok) {
            alert("Logged out successfully");
            window.location.href = "/";
        } else {
            alert("failed to log out");
        }
    }

    window.createPost = createPost;
    window.logout = logout;
    
    refreshToken().then((refreshed) => {
        if (refreshed) {
            loadPosts();
        }
    });

    setInterval(refreshToken, 25 * 60 * 1000);
}

window.onload = function () {
    const path = window.location.pathname;
    console.log("Current path:", path);
    if (path === '/' || path.endsWith('/index.html')) {
        initLoginPage();
    } else if (path.endsWith('/static/create_user.html')) {
        initCreateUserPage();
    } else if (path.endsWith('/static/posts.html')) {
        initPostPage();
    }
};
