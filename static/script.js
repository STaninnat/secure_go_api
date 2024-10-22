const API_BASE = '/v1';
let currentUserAPIKey = null;
let currentUser = null;

async function createPost() {
    if (!currentUser) {
        alert('Please create a user first');
        return;
    }
    const postContent = document.getElementById('newPostContent').value;
    const response = await fetchWithAlert(`${API_BASE}/posts`,
        {
            method: 'POST', headers: { 'Content-Type': 'application/json', 'Authorization': `ApiKey ${currentUserAPIKey}` },
            body: JSON.stringify({ post: postContent })
        });
    const post = await response.json();
    displayPost(post);
}

async function getUser() {
    const response = await fetchWithAlert(`${API_BASE}/users`, { headers: { 'Authorization': `ApiKey ${currentUserAPIKey}` } });
    return await response.json();
}

async function loadPosts() {
    if (!currentUser) {
        return;
    }
    const response = await fetchWithAlert(`${API_BASE}/posts`, { headers: { 'Authorization': `ApiKey ${currentUserAPIKey}` } });
    const posts = await response.json();
    const postsContainer = document.getElementById('posts');
    postsContainer.innerHTML = '';
    posts.forEach(post => displayPost(post));
}

function displayPost(post) {
    const postElement = document.createElement('div');
    postElement.className = 'post';
    postElement.textContent = post.post;
    document.getElementById('posts').appendChild(postElement);
}

async function createUser() {
    const nameField = document.getElementById('nameField');
    const name = nameField.value.trim();
    if (!name || name === '') {
        alert('Please enter a valid name');
        return;
    }
    const response = await fetchWithAlert(`${API_BASE}/users`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name: name.value })
    });
    const user = await response.json();
    localStorage.setItem('currentUserAPIKey', user.api_key);

    login()

    alert(`User Created: ${user.name}`);

    document.getElementById('userCreationContainer').style.display = 'none';
    document.getElementById('postSection').style.display = 'flex';
}

function logout() {
    localStorage.removeItem('currentUserAPIKey');
    currentUser = null;

    document.getElementById('userCreationContainer').style.display = 'block';
    document.getElementById('postSection').style.display = 'none';

    document.getElementById('userCreationContainer').style.display = 'inline-flex';
    document.getElementById('userCreationContainer').style.justifyContent = 'center';
    document.getElementById('userCreationContainer').style.alignItems = 'center';
}

async function login() {
    currentUserAPIKey = localStorage.getItem('currentUserAPIKey')
    if (!currentUserAPIKey) {
        return;
    }
    
    const user = await getUser();
    currentUser = user;
    currentUserAPIKey = user.api_key;
    await loadPosts();

    document.getElementById('userCreationContainer').style.display = 'none';
    document.getElementById('postSection').style.display = 'flex';
    document.getElementById('greetingMessage').textContent = `Hello ${user.name}!`;
}

async function fetchWithAlert(url, options) {
    const response = await fetch(url, options);
    if (response.status > 299) {
        alert(`Error: ${response.status}`);
        return;
    }
    return response;
}

login();