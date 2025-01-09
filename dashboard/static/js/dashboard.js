async function loadUsers() {
    try {
        const response = await fetch('/api/users', {
            headers: {
                'Authorization': `Bearer ${getCookie('session_token')}`
            }
        });
        
        const data = await response.json();
        
        if (!response.ok) {
            throw new Error(data.error || 'Failed to load users');
        }
        
        renderUsers(data);
    } catch (error) {
        if (error.message.includes('No authentication token')) {
            window.location.href = '/';
        }
        console.error('Error loading users:', error.message);
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
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${getCookie('session_token')}`
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
document.addEventListener('DOMContentLoaded', async () => {
    const token = getCookie('session_token');
    if (!token) {
        window.location.href = '/';
        return;
    }

    // Setup profile menu functionality
    const userProfile = document.getElementById('user-profile');
    const profileMenu = document.getElementById('profile-menu');
    const serverSettings = document.getElementById('server-settings');
    const logoutButton = document.getElementById('logout');

    userProfile.addEventListener('click', (e) => {
        e.stopPropagation();
        profileMenu.classList.toggle('active');
    });

    document.addEventListener('click', () => {
        profileMenu.classList.remove('active');
    });

    serverSettings.addEventListener('click', async (e) => {
        e.preventDefault();
        showServerSettingsModal();
    });

    logoutButton.addEventListener('click', async (e) => {
        e.preventDefault();
        try {
            await fetch('/api/logout', {
                method: 'POST',
                headers: {
                    'X-User-ID': getCookie('user_id'),
                    'Authorization': `Bearer ${getCookie('session_token')}`
                }
            });
        } catch (error) {
            console.error('Error during logout:', error);
        } finally {
            window.location.href = '/';
        }
    });

    try {
        const profileResponse = await fetch('/api/user/profile', {
            headers: {
                'Authorization': `Bearer ${token}`,
                'X-User-ID': getCookie('user_id')
            }
        });
        
        if (!profileResponse.ok) {
            throw new Error('Failed to load user profile');
        }
        
        const userData = await profileResponse.json();
        updateUserProfile(userData.username);
        
        // Then load view data
        await loadViewData('forest');
        await loadUsers();
    } catch (error) {
        window.location.href = '/';
        console.error('Error during initialization:', error);
    }
});

async function loadViewData(viewName) {
    try {
        const token = getCookie('session_token');
        const endpoints = {
            forest: '/api/forest',
            tree: '/api/tree',
            logs: '/api/logs',
            users: '/api/users',
            settings: '/api/settings'
        };

        const response = await fetch(endpoints[viewName], {
            headers: {
                'X-User-ID': getCookie('user_id'),
                'Authorization': `Bearer ${token}`,
            }
        });
        
        if (response.status === 401) {
            window.location.href = '/';
            console.error('Unauthorized access');
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
    const logsList = document.getElementById('logs-list');
    
    if (!logsList) {
        console.error('Logs list element not found');
        return;
    }

    // Handle empty or null data
    const logs = data || [];
    
    logsList.innerHTML = logs.length ? logs.map(log => `
        <div class="log-entry ${log.level.toLowerCase()}">
            <span class="log-timestamp">${new Date(log.timestamp).toLocaleString()}</span>
            <span class="log-level">${log.level}</span>
            <span class="log-message">${log.message}</span>
            ${log.trace ? `<pre class="log-trace">${log.trace}</pre>` : ''}
        </div>
    `).join('') : '<div class="no-logs">No logs found</div>';
}

// Add Permission enum to match Go constants
const Permission = {
    0: 'Read',
    1: 'Write',
    2: 'Admin'
};

function getCookie(name) {
    // console.log('Getting cookie:', name);
    // console.log('All cookies:', document.cookie);
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) {
        const cookieValue = parts.pop().split(';').shift();
        // console.log('Found cookie value:', cookieValue);
        return cookieValue;
    }
    // console.log('Cookie not found');
    return null;
}

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

        const response = await fetch('/api/settings/update', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${getCookie('session_token')}`
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
        const response = await fetch('/api/users/profile', {
            headers: {
                'Authorization': `Bearer ${getCookie('session_token')}`,
                'X-User-ID': getCookie('user_id')
            }
        });
        
        if (!response.ok) {
            throw new Error('Failed to load user profile');
        }
        
        const user = await response.json();
        updateUserProfile(user.username);
        
    } catch (error) {
        if (error.message.includes('No authentication token')) {
            window.location.href = '/';
        }
        console.error('Error loading user profile:', error);
    }
}

function updateUserProfile(username) {
    const organizationDisplay = document.getElementById('organization-display');
    const usernameDisplay = document.getElementById('username-display');
    
    if (username) {
        const initials = "org"
        organizationDisplay.textContent = initials;
        usernameDisplay.textContent = username;
    }
}

async function refreshToken() {
    try {
        const response = await fetch('/api/refresh', {
            method: 'POST',
            credentials: 'same-origin'
        });
        
        if (!response.ok) {
            throw new Error('Failed to refresh token');
        }
        
        return true;
    } catch (error) {
        window.location.href = '/';
        console.error('Error refreshing token:', error);
        return false;
    }
}

// Update the fetch wrapper to include credentials
async function fetchWithAuth(url, options = {}) {
    try {
        const response = await fetch(url, {
            ...options,
            credentials: 'include',
            headers: {
                ...options.headers,
                'Authorization': `Bearer ${getCookie('session_token')}`
            }
        });
        
        if (response.status === 401) {
            const refreshed = await refreshToken();
            if (refreshed) {
                return fetch(url, {
                    ...options,
                    credentials: 'include',
                    headers: {
                        ...options.headers,
                        'Authorization': `Bearer ${getCookie('session_token')}`
                    }
                });
            }
        }
        
        return response;
    } catch (error) {
        console.error('Fetch error:', error);
        throw error;
    }
}

// Add token refresh interval
setInterval(async () => {
    const sessionExpiry = getCookie('session_expiry');
    if (sessionExpiry && Date.now() > parseInt(sessionExpiry) - (5 * 60 * 1000)) {
        await refreshToken();
    }
}, 300000); 

async function loadForest() {
    try {
        const response = await fetch('/api/forest', {
            headers: {
                'Authorization': `Bearer ${getCookie('session_token')}`
            }
        });
        
        if (!response.ok) {
            throw new Error('Failed to load forest data');
        }
        
        const data = await response.json();
        renderForest(data);
    } catch (error) {
        console.error('Error loading forest:', error);
    }
}

async function loadLogs() {
    try {
        const response = await fetchWithAuth('/api/logs', {
            headers: {
                'Authorization': `Bearer ${getCookie('session_token')}`
            }
        });
        
        if (!response.ok) {
            throw new Error('Failed to load logs');
        }

        const data = await response.json();
        console.log('Logs data:', data);
        renderLogs(data);
    } catch (error) {
        console.error('Error loading logs:', error);
        renderLogs([]); // Render empty state on error
    }
}

// Add log filtering functionality
document.getElementById('log-level')?.addEventListener('change', async (e) => {
    const level = e.target.value;
    const url = level === 'all' ? '/api/logs' : `/api/logs?level=${level}`;
    
    try {
        const response = await fetch(url, {
            headers: {
                'Authorization': `Bearer ${getCookie('session_token')}`
            }
        });
        
        if (!response.ok) {
            throw new Error('Failed to filter logs');
        }
        
        const data = await response.json();
        renderLogs(data);
    } catch (error) {
        console.error('Error filtering logs:', error);
    }
});

document.getElementById('log-search')?.addEventListener('input', (e) => {
    const searchTerm = e.target.value.toLowerCase();
    const logEntries = document.querySelectorAll('.log-entry');
    
    logEntries.forEach(entry => {
        const text = entry.textContent.toLowerCase();
        entry.style.display = text.includes(searchTerm) ? 'block' : 'none';
    });
});

// Add click handler for the logs tab
document.querySelector('a[data-view="logs"]')?.addEventListener('click', async () => {
    document.querySelectorAll('.view').forEach(view => view.style.display = 'none');
    document.getElementById('logs-view').style.display = 'block';
    await loadLogs();
});
