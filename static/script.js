const API_BASE = '/v1';
let refreshInterval;

async function fetchWithAlert(url, options = {}) {
    const response = await fetch(url, {
        ...options,
        credentials: 'include',
        headers: {
            ...options.headers,
        }
    });

    if (response.status === 401) {
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
            window.location.href = "/";
            return;
        }
    }
    
    if (response.status > 299) {
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
            window.location.href = "/";
            return response;
        }

        return response;
    } catch (error) {
        window.location.href = "/";
        return { ok: false, message: 'An error occurred while refreshing the token' };
    }
}

function displayMessage(elementId, message, isError = false, duration = 0) {
    const element = document.getElementById(elementId);
    element.textContent = message;
    element.style.color = isError ? 'red' : 'green';

    if (duration > 0) {
        setTimeout(() => {
            element.textContent = '';
        }, duration);
    }
}

function initLoginPage() {
    async function loginUser(event) {
        if (event) event.preventDefault();

        const name = document.getElementById('nameFieldLogin').value.trim();
        const password = document.getElementById('passwordFieldLogin').value.trim();
        const alertElement = 'alertLogin';

        if (!name && !password) {
            displayMessage(alertElement, "Please enter both username and password.", true);
            return;
        } else if (!name) {
            displayMessage(alertElement, "Please enter username.", true);
            return;
        } else if (!password) {
            displayMessage(alertElement, "Please enter password.", true);
            return;
        }

        const response = await fetchWithAlert(`${API_BASE}/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name, password })
        });

        if (response && response.ok) {
            displayMessage(alertElement, `Login successful...WELCOME BACK ${name}`);
            sessionStorage.setItem('username', name);
            setTimeout(() => {
                window.location.href = "/static/posts.html";
            }, 2000);
        } else {
            const errorData = await response.json();

            if ((response.status === 400) && (errorData.error === "username not found")) {
                displayMessage(alertElement, "Invalid username. Please try again.", true);
            } else if ((response.status === 400) && (errorData.error === "incorrect password")) {
                displayMessage(alertElement, "Invalid password. Please try again.", true);
            } else {
                displayMessage(alertElement, "Login failed. Please check your credentials.", true);
            }
        }
    }

    window.loginUser = loginUser;
}

function initCreateUserPage() {
    async function createUser(event) {
        if (event) event.preventDefault();

        const name = document.getElementById('nameFieldCreate').value.trim();
        const password = document.getElementById('passwordFieldCreate').value.trim();
        const alertElement = 'alertCreate';
        
        if (!name && !password) {
            displayMessage(alertElement, "Please enter both username and password", true);
            return;
        } else if (!name) {
            displayMessage(alertElement, "Please enter username", true);
            return;
        } else if (!password) {
            displayMessage(alertElement, "Please enter password", true);
            return;
        }
        
        const response = await fetch(`${API_BASE}/users`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name, password })
        });

        if (response && response.ok) {
            displayMessage(alertElement, "User created successfully...WELCOME!");
            setTimeout(() => {
                window.location.href = "/static/posts.html";
            }, 2000);
        } else {
            const errorData = await response.json();
            if ((response.status === 400) && (errorData.error === "username already exists")) {
                displayMessage(alertElement, "Username already exists. Please try again.", true);
            } else {
                console.error(`User creation failed ${errorData}`);
                displayMessage(alertElement, "User creation failed", true);
            }
        }
    }

    window.createUser = createUser;
}

async function initPostPage() {
    const username = sessionStorage.getItem('username');
    const alertElement = 'alertPost';
    
    if (username) {
        document.getElementById('greetingMessage').textContent = `Hello, ${username}! What's on your mind today?`;
    } else {
        document.getElementById('greetingMessage').textContent = `Hello! What's on your mind today?`;
    }

    async function loadPosts() {
        const response = await fetchWithAlert(`${API_BASE}/posts`, {
            method: 'GET',
            credentials: 'include'
        });

        const postsContainer = document.getElementById('posts');
        postsContainer.innerHTML = '';

        if (response.ok) {
            const posts = await response.json();
            if (posts.length === 0) {
                const noPostsMessage = document.createElement('p');
                noPostsMessage.className = 'no-posts-message';
                noPostsMessage.id = 'noPostsMessage';
                noPostsMessage.textContent = 'No posts yet. Be the first to post something!';
                postsContainer.appendChild(noPostsMessage);
            } else {
                posts.forEach(post => displayPost(post));
            }
        } else {
            console.error(`Error: ${response.status} ${response.statusText}`);
        }
    }

    async function createPost() {
        const postContent = document.getElementById('newPostContent').value;

        if (!postContent) {
            displayMessage(alertElement, 'Please enter post content', true, 5000);
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
            const noPostsMessage = document.getElementById('noPostsMessage');
            if (noPostsMessage) {
                noPostsMessage.remove();
            }

            const post = await response.json();
            displayPost(post);
            document.getElementById('newPostContent').value = '';
            displayMessage(alertElement, 'Post created successfully!', false, 5000);
        } else {
            displayMessage(alertElement, 'Failed to create post', true);
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
            displayMessage(alertElement, 'Logged out successfully');

            setTimeout(() => {
                sessionStorage.removeItem('username');
                clearInterval(refreshInterval);
                window.location.href = "/";
            }, 2000);
        } else {
            displayMessage(alertElement, 'Failed to log out', true);
        }
    }

    window.createPost = createPost;
    window.logout = logout;

    // const refreshed = await refreshToken();
    // if (refreshed && refreshed.ok) {
    //     loadPosts();
    // }
    loadPosts();

    refreshInterval = setInterval(async () => {
        try {
            const result = await refreshToken();
            if (!result || !result.ok) {
                alert("Session expired. Please log in again.");
                window.location.href = "/";
            }
        } catch (error) {
            alert("Error refreshing session. Please check your connection.");
        }
    }, 10 * 60 * 1000);
}

function initContainer() {
    const container = document.getElementById('container');
    const registerBtn = document.getElementById('register');
    const loginBtn = document.getElementById('login');
    const alertCreate = document.getElementById('alertCreate');
    const alertLogin = document.getElementById('alertLogin');

    registerBtn.addEventListener('click', () => {
        container.classList.add("active");
        if (alertLogin) {
            alertLogin.textContent = '';
        }
    });

    loginBtn.addEventListener('click', () => {
        container.classList.remove("active");
        if (alertCreate) {
            alertCreate.textContent = '';
        }
    });
}

window.onload = function () {
    const path = window.location.pathname;
    if (path === '/' || path.endsWith('/index.html')) {
        initContainer();
        initLoginPage();
        initCreateUserPage();
    } else if (path.endsWith('/static/posts.html')) {
        initPostPage();
    }
};