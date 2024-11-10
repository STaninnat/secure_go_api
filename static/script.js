const API_BASE = '/v1';

async function fetchWithAlert(url, options = {}) {
    const token = sessionStorage.getItem('access_token');

    const response = await fetch(url, {
        ...options,
        credentials: 'include',
        headers: {
            ...options.headers,
            'Authorization': token ? `Bearer ${token}` : ''
        }
    });

    if (response.status === 401) {
        const refreshResponse = await refreshToken();
        if (refreshResponse && refreshResponse.ok) {
            const newToken = sessionStorage.getItem('access_token');
            return fetch(url, {
                ...options,
                credentials: 'include',
                headers: {
                    ...options.headers,
                    'Authorization': newToken ? `Bearer ${newToken}` : ''
                }
            });
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

async function refreshToken() {
    const response = await fetchWithAlert(`${API_BASE}/refresh`, {
        method: 'POST',
        credentials: 'include'
    });

    if (response && response.ok) {
        const data = await response.json();
        if (data.access_token) {
            sessionStorage.setItem('access_token', data.access_token);
            console.log("Token refreshed successfully");
        } else {
            alert("failed to refresh token. please log in again");
            window.location.href = "/";
        }
    } else {
        alert("failed to refresh token. please log in again");
        window.location.href = "/";
    }
    return response;
}

function initLoginPage() {
    console.log("Initializing Login Page");

    async function loginUser(event) {
        event.preventDefault();

        const name = document.getElementById('nameFieldLogin').value.trim();
        const password = document.getElementById('passwordFieldLogin').value.trim();

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
            const data = await response.json();
            sessionStorage.setItem('access_token', data.access_token);
            alert(`Login successful, welcome ${name}`);
            window.location.href = "/static/posts.html";
        } else {
            if (response.status === 401) {
                alert("Invalid username or password. Please try again.");
            } else {
                alert("Login failed. Please check your credentials.");
            }
        }
    }

    window.loginUser = loginUser;
}

function initCreateUserPage() {
    console.log("Initializing Create User Page");

    async function createUser() {
        const name = document.getElementById('nameFieldCreate').value.trim();
        const password = document.getElementById('passwordFieldCreate').value.trim();
        
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
            const data = await response.json();
            alert("User created successfully. WELCOME!");
            sessionStorage.setItem('access_token', data.access_token);
            console.log("Access Token Stored:", data.access_token);
            window.location.href = "/static/posts.html"; 
            // if (data.access_token) {
            //     alert("User created successfully. WELCOME!");
            //     sessionStorage.setItem('access_token', data.access_token);
            //     console.log("Access Token Stored:", data.access_token);
            //     window.location.href = "/static/posts.html"; 
            // } else {
            //     alert("User creation successful, but no token received. Please log in.");
            //     window.location.href = "/";
            // }
        } else {
            alert("user creation failed");
        }
    }

    window.createUser = createUser;
}

async function initPostPage() {
    console.log("Initializing Post Page");

    if (!sessionStorage.getItem('access_token')) {
        alert("No access token found. Please log in.");
        window.location.href = "/";
        return;
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
            sessionStorage.removeItem('access_token');
            window.location.href = "/";
        } else {
            alert("failed to log out");
        }
    }

    window.createPost = createPost;
    window.logout = logout;

    const refreshed = await refreshToken();
    if (refreshed && refreshed.ok) {
        loadPosts();
    }

    setInterval(async () => {
        const result = await refreshToken();
        if (!result || !result.ok) {
            console.log("Token refresh failed. Redirecting to login.");
            alert("Session expired. Please log in again.");
            window.location.href = "/";
        }
    }, 25 * 60 * 1000);
}

function initContainer() {
    console.log("Initializing Container");

    const container = document.getElementById('container');
    const registerBtn = document.getElementById('register');
    const loginBtn = document.getElementById('login');

    registerBtn.addEventListener('click', () => {
        container.classList.add("active");
    });

    loginBtn.addEventListener('click', () => {
        container.classList.remove("active");
    });
}

window.onload = function () {
    const path = window.location.pathname;
    console.log("Current path:", path);
    if (path === '/' || path.endsWith('/index.html')) {
        initContainer();
        initLoginPage();
        initCreateUserPage();
    } else if (path.endsWith('/static/posts.html')) {
        initPostPage();
    }
};