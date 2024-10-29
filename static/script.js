const API_BASE = '/v1';
let currentUserAPIKey = null;
let currentUserToken = null;

function initLoginPage() {
    console.log("Initializing Login Page");

    async function loginUser() {
        const name = document.getElementById('nameField').value.trim();
        const password = document.getElementById('passwordField').value.trim();
        
        if (!name || !password) {
            alert("Please enter both username and password.");
            return;
        }
        
        const response = await fetch(`${API_BASE}/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name, password })
        });

        if (response.ok) {
            const user = await response.json();
            localStorage.setItem('currentUserToken', user.token);

            alert(`Login successful, welcome ${name}`);
            window.location.href = "post.html"; 
        } else {
            alert("Login failed. Please check your credentials.");
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
            alert("Please enter both username and password.");
            return;
        }
        
        const response = await fetch(`${API_BASE}/users`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name, password })
        });

        if (response.ok) {
            const user = await response.json();

            localStorage.setItem('currentUserAPIKey', user.api_key);
            localStorage.setItem('currentUserToken', user.token);

            alert("User created successfully. Welcome!");

            window.location.href = "post.html";
        } else {
            alert("User creation failed.");
        }
    }

    window.createUser = createUser;
}

function initPostPage() {
    console.log("Initializing Post Page");

    currentUserToken = localStorage.getItem('currentUserToken');

    if (!currentUserToken) {
        alert("Please log in first.");
        window.location.href = "index.html";
        return;
    }

    async function loadPosts() {
        const response = await fetch(`${API_BASE}/posts`, {
            headers: { 'Authorization': `Bearer ${currentUserToken}` }
        });

        if (response.ok) {
            const posts = await response.json();
            const postsContainer = document.getElementById('posts');
            postsContainer.innerHTML = '';
            posts.forEach(post => displayPost(post));
        } else {
        alert('Error loading posts');
    }
    }

    async function createPost() {
        const postContent = document.getElementById('newPostContent').value;

        if (!postContent) {
            alert('Please enter post content');
            return;
        }

        const response = await fetch(`${API_BASE}/posts`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${currentUserToken}`
            },
            body: JSON.stringify({ post: postContent })
        });

        if (response.ok) {
            const post = await response.json();
            displayPost(post);
            document.getElementById('newPostContent').value = '';
        } else {
            alert("Failed to create post.");
        }
    }

    function displayPost(post) {
        const postElement = document.createElement('div');
        postElement.className = 'post';
        postElement.textContent = post.post;
        document.getElementById('posts').appendChild(postElement);
    }

    function logout() {
        localStorage.removeItem('currentUserToken');
        currentUserToken = null;
        window.location.href = "index.html";
    }

    window.createPost = createPost;
    window.logout = logout;
    loadPosts();
}

async function getUser() {
    const response = await fetchWithAlert(`${API_BASE}/users`, {
        headers: {
                'Authorization': `Bearer ${currentUserToken}`
        }
    });

    if (!response) return null; 
    return await response.json();
}

async function fetchWithAlert(url, options) {
    const token = localStorage.getItem('currentUserToken');
    if (token) {
        options.headers = {
            ...options.headers,
            'Authorization': `Bearer ${token}`
        };
    } else {
        alert("Session expired. Please log in again.");
        window.location.href = "index.html";
        return;
    }

    const response = await fetch(url, options);
    if (response.status === 401) {
        alert("Session expired. Please log in again.");
        window.location.href = "index.html";
        return;
    }
    if (response.status > 299) {
        alert(`Error: ${response.status}`);
        return;
    }
    return response;
}

window.onload = function () {
    currentUserToken = localStorage.getItem('currentUserToken');

    const path = window.location.pathname;
    if (path.includes('index.html')) {
        initLoginPage();
    } else if (path.includes('create_user.html')) {
        initCreateUserPage();
    } else if (path.includes('post.html')) {
        initPostPage();
    }
};
