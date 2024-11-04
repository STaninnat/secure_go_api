const API_BASE = '/v1';

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
            window.location.href = "post.html"; 
        } else {
            alert("login failed. Please check your credentials");
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
            alert("User created successfully. Welcome!");
            window.location.href = "post.html";
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

        if (response.ok) {
            console.log("Token refreshed successfully");
            return true;
        } else {
            alert("failed to refresh token. please log in again");
            window.location.href = "index.html";
            return false;
        }
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
                credentials: 'include',
            },
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
            window.location.href = "index.html";
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

async function fetchWithAlert(url, options = {}) {
    let response = await fetch(url, options);

    if (response.status === 401) {
        const refreshSuccess = await refreshToken();

        if (refreshSuccess) {
            response = await fetch(url, options);
        } else {
            alert("session expired. please log in again");
            window.location.href = "index.html";
            return;
        }
    }
    
    if (response.status > 299) {
        alert(`Error: ${response.status}`);
        return;
    }
    return response;
}

window.onload = function () {
    const path = window.location.pathname;
    if (path.includes('index.html')) {
        initLoginPage();
    } else if (path.includes('create_user.html')) {
        initCreateUserPage();
    } else if (path.includes('post.html')) {
        initPostPage();
    }
};
