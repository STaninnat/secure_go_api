const API_BASE = '/v1';
let refreshInterval;

async function fetchWithAlert(url, options = {}) {
    const token = sessionStorage.getItem('access_token');

    const response = await fetch(url, {
        ...options,
        credentials: 'include',
        headers: {
            ...options.headers,
        }
    });

    if (response.status === 401) {
        alert("Session expired, trying to refresh token...");
        const refreshResponse = await refreshToken();
        if (refreshResponse && refreshResponse.ok) {
            return fetch(url, {
                ...options,
                credentials: 'include',
                headers: {
                    ...options.headers,
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
    try {
        const response = await fetch(`${API_BASE}/refresh`, {
            method: 'POST',
            credentials: 'include'
        });

        if (!response.ok) {
            console.error(`Error: ${response.status} ${response.statusText}`);
            if (response.status === 500) {
                alert("Server error while refreshing token. Please check your connection and try again.");
            } else {
                alert("Failed to refresh token. Please log in again.");
            }
            window.location.href = "/";
            return { ok: false, message: 'Failed to refresh token' };
        }

        return response;
    } catch (error) {
        console.error('Error during token refresh:', error);
        alert("An error occurred while refreshing the token. Please log in again.");
        window.location.href = "/";
        return { ok: false, message: 'An error occurred while refreshing the token' };
    }
}

function initLoginPage() {
    console.log("Initializing Login Page");

    async function loginUser(event) {
        if (event) event.preventDefault();

        const name = document.getElementById('nameFieldLogin').value.trim();
        const password = document.getElementById('passwordFieldLogin').value.trim();

        if (!name && !password) {
            alert("please enter both username and password");
            return;
        } else if (!name) {
            alert("please enter username");
            return;
        } else if (!password) {
            alert("please enter password");
            return;
        }

        const response = await fetchWithAlert(`${API_BASE}/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name, password })
        });

        if (response && response.ok) {
            alert(`Login successful, welcome ${name}`);
            window.location.href = "/static/posts.html";
        } else {
            const errorData = await response.json();
            console.log("Login failed: ", errorData.error);

            if ((response.status === 400) && (errorData.error === "username not found")) {
                alert("Invalid username. Please try again.");
            } else if ((response.status === 400) && (errorData.error === "incorrect password")) {
                alert("Invalid password. Please try again.");
            } else {
                alert("Login failed. Please check your credentials.");
            }
        }
    }

    window.loginUser = loginUser;
}

function initCreateUserPage() {
    console.log("Initializing Create User Page");

    async function createUser(event) {
        if (event) event.preventDefault();

        const name = document.getElementById('nameFieldCreate').value.trim();
        const password = document.getElementById('passwordFieldCreate').value.trim();
        
        if (!name && !password) {
            alert("please enter both username and password");
            return;
        } else if (!name) {
            alert("please enter username");
            return;
        } else if (!password) {
            alert("please enter password");
            return;
        }
        
        const response = await fetch(`${API_BASE}/users`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name, password })
        });

        if (response.ok) {
            alert("User created successfully. WELCOME!");
            window.location.href = "/static/posts.html";
        } else {
            const errorData = await response.json();
            if ((response.status === 400) && (errorData.error === "username already exists")) {
                alert("Username already exists. Please try again.");
            } else {
                alert("user creation failed");
            }
        }
    }

    window.createUser = createUser;
}

async function initPostPage() {
    console.log("Initializing Post Page");

    async function loadPosts() {
        const response = await fetchWithAlert(`${API_BASE}/posts`, {
            method: 'GET',
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
            clearInterval(refreshInterval);
            window.location.href = "/";
        } else {
            alert("failed to log out");
        }
    }

    window.createPost = createPost;
    window.logout = logout;

    // const refreshed = await refreshToken();
    // if (refreshed && refreshed.ok) {
    //     loadPosts();
    //     console.log("loadPosts()");
    // }
    loadPosts();
    console.log("loadPosts()");

    refreshInterval = setInterval(async () => {
        try {
            const result = await refreshToken();
            if (!result || !result.ok) {
                console.log("Token refresh failed. Redirecting to login.");
                alert("Session expired. Please log in again.");
                window.location.href = "/";
            }
            console.log("refresh token successfully.")
        } catch (error) {
            console.error("Error refreshing token:", error);
            alert("Error refreshing session. Please check your connection.");
        }
    }, 10 * 60 * 1000);
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