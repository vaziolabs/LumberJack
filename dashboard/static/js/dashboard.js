async function loadUsers() {
    try {
        const response = await fetch('/api/users', {
            headers: {
                'Authorization': `Bearer ${getCookie('auth_token')}`
            }
        });
        
        const data = await response.json();
        
        if (!response.ok) {
            throw new Error(data.error || 'Failed to load users');
        }
        
        renderUsers(data);
    } catch (error) {
        console.error('Error loading users:', error.message);
        if (error.message.includes('No authentication token')) {
            window.location.href = '/';  // Redirect to login
        }
    }
}

function renderUsers(data) {
    const usersView = document.getElementById('users-view');
    const usersList = usersView.querySelector('.users-list');
    
    if (!usersList) return;
    
    // Handle empty or null data
    const users = data || [];
    
    usersList.innerHTML = users.length ? users.map(user => `
        <div class="user-item">
            <div class="user-info">
                <span class="username">${user.username}</span>
                <span class="email">${user.email}</span>
            </div>
            <div class="permissions">
                ${user.permissions.map(perm => `
                    <span class="permission-tag">${Permission[perm]}</span>
                `).join('')}
            </div>
        </div>
    `).join('') : '<div class="no-users">No users found</div>';
}

function showModal() {
    document.getElementById('user-modal').style.display = 'block';
}

function closeModal() {
    document.getElementById('user-modal').style.display = 'none';
}

document.getElementById('add-user-btn').addEventListener('click', showModal);

document.getElementById('user-form').addEventListener('submit', async (e) => {
    e.preventDefault();
    const formData = {
        username: document.getElementById('username').value,
        email: document.getElementById('email').value,
        permissions: Array.from(document.querySelectorAll('.permissions-group input:checked'))
            .map(input => parseInt(input.value))
    };

    try {
        const response = await fetch('/api/users', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(formData)
        });

        if (response.ok) {
            closeModal();
            loadUsers();
        } else {
            console.error('Error creating user:', await response.text());
        }
    } catch (error) {
        console.error('Error creating user:', error);
    }
});

function showView(viewName) {
    const views = document.querySelectorAll('.view');
    views.forEach(view => view.style.display = 'none');
    document.getElementById(`${viewName}-view`).style.display = 'block';
}

// Add click handlers for navigation
document.querySelectorAll('.sidebar a').forEach(link => {
    link.addEventListener('click', (e) => {
        e.preventDefault();
        const view = e.target.getAttribute('data-view');
        showView(view);
    });
});

// Modify the initial load section
async function checkAuthAndLoadData() {
    const token = getCookie('auth_token');
    if (!token) {
        window.location.href = '/';  // Redirect to login if no token
        return;
    }
    
    await loadViewData('tree');
    await loadUsers();
}

async function loadViewData(viewName) {
    try {
        const token = getCookie('auth_token');
        const response = await fetch(`/api/${viewName}`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        
        if (response.status === 401) {
            window.location.href = '/';  // Redirect to login if unauthorized
            return;
        }
        
        const data = await response.json();
        switch(viewName) {
            case 'forest':
                renderForest(data);
                break;
            case 'tree':
                renderEvents(data);
                break;
            case 'logs':
                renderLogs(data);
                break;
            case 'users':
                renderUsers(data);
                break;
        }
    } catch (error) {
        console.error(`Error loading ${viewName} data:`, error);
    }
}

function renderForest(data) {
    const forestView = document.getElementById('forest-view');
    forestView.innerHTML = `
        <div class="forest-container">
            <div class="forest-controls">
                <button id="zoom-in" class="btn">+</button>
                <button id="zoom-out" class="btn">-</button>
                <button id="reset-view" class="btn">Reset</button>
            </div>
            <div id="forest-graph" class="forest-graph"></div>
        </div>
    `;

    // Initialize forest visualization
    initForestGraph(data);
}

function initForestGraph(data) {
    // TODO: Implement D3.js or similar visualization library
    // to create interactive forest graph visualization
    // Example structure:
    // - Nodes represented as circles/rectangles
    // - Lines connecting parent-child relationships
    // - Zoom/pan capabilities
    // - Click to navigate into nodes
}

function renderEvents(data) {
    const eventsView = document.getElementById('events-view');
    eventsView.innerHTML = `
        <div class="events-container">
            <h2>Events</h2>
            <div class="events-list">
                ${data.map(event => `
                    <div class="event-item">
                        <div class="event-header">${event.content}</div>
                        <div class="event-time">${new Date(event.timestamp).toLocaleString()}</div>
                    </div>
                `).join('')}
            </div>
        </div>
    `;
}

function renderLogs(data) {
    const logsView = document.getElementById('logs-view');
    logsView.innerHTML = `
        <div class="logs-container">
            <h2>System Logs</h2>
            <div class="logs-filter">
                <select id="log-level">
                    <option value="all">All Levels</option>
                    <option value="debug">Debug</option>
                    <option value="info">Info</option>
                    <option value="notice">Notice</option>
                    <option value="warn">Warning</option>
                    <option value="error">Error</option>
                    <option value="critical">Critical</option>
                    <option value="alert">Alert</option>
                    <option value="emergency">Emergency</option>
                    <option value="success">Success</option>
                    <option value="failure">Failure</option>
                    <option value="enter">Enter</option>
                    <option value="exit">Exit</option>
                    <option value="other">Uncategorized</option>
                </select>
                <input type="text" id="log-search" placeholder="Search logs...">
            </div>
            <div class="logs-list" id="logs-list"></div>
        </div>
    `;

    updateLogs(data);
}

function updateLogs(logs) {
    const logsList = document.getElementById('logs-list');
    logsList.innerHTML = logs.map(log => `
        <div class="log-entry ${log.level.toLowerCase()}">
            <span class="log-timestamp">${new Date(log.timestamp).toLocaleString()}</span>
            <span class="log-level">${log.level}</span>
            <span class="log-message">${log.message}</span>
        </div>
    `).join('');
}

// Add Permission enum to match Go constants
const Permission = {
    0: 'Read',
    1: 'Write',
    2: 'Admin'
};

function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
}

// Initialize on page load
document.addEventListener('DOMContentLoaded', checkAuthAndLoadData);

// Profile menu functionality
document.addEventListener('DOMContentLoaded', () => {
    const userProfile = document.getElementById('user-profile');
    const profileMenu = document.getElementById('profile-menu');
    const serverSettings = document.getElementById('server-settings');
    const logoutButton = document.getElementById('logout');

    // Toggle menu on profile click
    userProfile.addEventListener('click', (e) => {
        e.stopPropagation();
        profileMenu.classList.toggle('active');
    });

    // Close menu when clicking outside
    document.addEventListener('click', () => {
        profileMenu.classList.remove('active');
    });

    // Server Settings handler
    serverSettings.addEventListener('click', async (e) => {
        e.preventDefault();
        // TODO: Implement server settings modal
        showServerSettingsModal();
    });

    // Logout handler
    logoutButton.addEventListener('click', async (e) => {
        e.preventDefault();
        try {
            const response = await fetch('/api/logout', {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${getCookie('auth_token')}`
                }
            });
            
            // Clear auth token regardless of response
            document.cookie = 'auth_token=; Path=/; Expires=Thu, 01 Jan 1970 00:00:01 GMT;';
            window.location.href = '/';
            
        } catch (error) {
            console.error('Error during logout:', error);
            // Force redirect to login page even if logout fails
            window.location.href = '/';
        }
    });

    // Load user profile on page load
    loadUserProfile();
});

// Add Server Settings Modal functionality
function showServerSettingsModal() {
    const modalHtml = `
        <div id="server-settings-modal" class="modal">
            <div class="modal-content">
                <h3>Server Settings</h3>
                <form id="server-settings-form">
                    <div class="form-group">
                        <label for="api-endpoint">API Endpoint</label>
                        <input type="text" id="api-endpoint" value="${window.location.origin}/api" readonly>
                    </div>
                    <div class="form-group">
                        <label for="log-level">Default Log Level</label>
                        <select id="log-level">
                            <option value="debug">Debug</option>
                            <option value="info">Info</option>
                            <option value="warn">Warning</option>
                            <option value="error">Error</option>
                        </select>
                    </div>
                    <div class="form-actions">
                        <button type="submit" class="btn">Save</button>
                        <button type="button" class="btn btn-secondary" onclick="closeServerSettingsModal()">Cancel</button>
                    </div>
                </form>
            </div>
        </div>
    `;

    document.body.insertAdjacentHTML('beforeend', modalHtml);
    const modal = document.getElementById('server-settings-modal');
    modal.style.display = 'block';

    // Add form submit handler
    document.getElementById('server-settings-form').addEventListener('submit', handleServerSettingsSave);
}

function closeServerSettingsModal() {
    const modal = document.getElementById('server-settings-modal');
    modal.remove();
}

async function handleServerSettingsSave(e) {
    e.preventDefault();
    try {
        const settings = {
            logLevel: document.getElementById('log-level').value
        };

        const response = await fetch('/api/settings', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${getCookie('auth_token')}`
            },
            body: JSON.stringify(settings)
        });

        if (response.ok) {
            closeServerSettingsModal();
        } else {
            console.error('Failed to save settings:', await response.text());
        }
    } catch (error) {
        console.error('Error saving settings:', error);
    }
}

async function loadUserProfile() {
    try {
        const response = await fetch('/api/user/profile', {
            headers: {
                'Authorization': `Bearer ${getCookie('auth_token')}`
            }
        });
        
        if (!response.ok) {
            throw new Error('Failed to load user profile');
        }
        
        const user = await response.json();
        updateUserProfile(user.username);
        
    } catch (error) {
        console.error('Error loading user profile:', error);
        if (error.message.includes('No authentication token')) {
            window.location.href = '/';
        }
    }
}

function updateUserProfile(username) {
    const userInitials = document.getElementById('user-initials');
    const usernameDisplay = document.getElementById('username-display');
    
    if (username) {
        const initials = username
            .split(' ')
            .map(word => word[0])
            .join('')
            .toUpperCase();
        userInitials.textContent = initials;
        usernameDisplay.textContent = username;
    }
}
