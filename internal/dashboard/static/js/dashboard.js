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
    // Set up SVG container
    const width = document.getElementById('forest-graph').clientWidth;
    const height = document.getElementById('forest-graph').clientHeight || 600;
    
    const svg = d3.select('#forest-graph')
        .append('svg')
        .attr('width', width)
        .attr('height', height);
        
    // Create zoom behavior
    const zoom = d3.zoom()
        .scaleExtent([0.1, 4])
        .on('zoom', (event) => {
            container.attr('transform', event.transform);
        });
        
    svg.call(zoom);
    
    // Create container for zoomable content
    const container = svg.append('g');
    
    // Create force simulation
    const simulation = d3.forceSimulation()
        .force('link', d3.forceLink().id(d => d.id).distance(100))
        .force('charge', d3.forceManyBody().strength(-300))
        .force('center', d3.forceCenter(width / 2, height / 2));
        
    // Process data into nodes and links
    const nodes = [];
    const links = [];
    
    function processNode(node, parentId = null) {
        const nodeData = {
            id: node.path,
            name: node.name || node.path.split('/').pop(),
            type: node.type,
            status: node.status
        };
        nodes.push(nodeData);
        
        if (parentId) {
            links.push({
                source: parentId,
                target: nodeData.id
            });
        }
        
        // Check if Children exists and is an object
        if (node.Children && typeof node.Children === 'object') {
            // Handle both array and object cases
            const children = Array.isArray(node.Children) ? 
                node.Children : 
                Object.values(node.Children);
                
            children.forEach(child => processNode(child, nodeData.id));
        }
    }
    
    // Handle forest data which might be an array or single node
    if (Array.isArray(data)) {
        data.forEach(node => processNode(node));
    } else {
        processNode(data);
    }

    // Create links
    const link = container.append('g')
        .selectAll('line')
        .data(links)
        .join('line')
        .attr('stroke', '#999')
        .attr('stroke-opacity', 0.6)
        .attr('stroke-width', 2);
        
    // Create nodes
    const node = container.append('g')
        .selectAll('g')
        .data(nodes)
        .join('g')
        .call(d3.drag()
            .on('start', dragStarted)
            .on('drag', dragged)
            .on('end', dragEnded));
            
    // Add circles for nodes
    node.append('circle')
        .attr('r', 10)
        .attr('fill', d => d.type === 'leaf' ? '#69b3a2' : '#3498db')
        .attr('stroke', '#fff')
        .attr('stroke-width', 2);
        
    // Add labels
    node.append('text')
        .text(d => d.name)
        .attr('x', 15)
        .attr('y', 5)
        .attr('font-size', '12px');
        
    // Add click handler
    node.on('click', (event, d) => {
        event.stopPropagation();
        loadViewData('tree', d.id);
    });
    
    // Update force simulation
    simulation
        .nodes(nodes)
        .on('tick', () => {
            link
                .attr('x1', d => d.source.x)
                .attr('y1', d => d.source.y)
                .attr('x2', d => d.target.x)
                .attr('y2', d => d.target.y);
                
            node
                .attr('transform', d => `translate(${d.x},${d.y})`);
        });
        
    simulation.force('link').links(links);
    
    // Drag functions
    function dragStarted(event) {
        if (!event.active) simulation.alphaTarget(0.3).restart();
        event.subject.fx = event.subject.x;
        event.subject.fy = event.subject.y;
    }
    
    function dragged(event) {
        event.subject.fx = event.x;
        event.subject.fy = event.y;
    }
    
    function dragEnded(event) {
        if (!event.active) simulation.alphaTarget(0);
        event.subject.fx = null;
        event.subject.fy = null;
    }
    
    // Add zoom controls
    document.getElementById('zoom-in').onclick = () => {
        zoom.scaleBy(svg.transition().duration(750), 1.2);
    };
    
    document.getElementById('zoom-out').onclick = () => {
        zoom.scaleBy(svg.transition().duration(750), 0.8);
    };
    
    document.getElementById('reset-view').onclick = () => {
        svg.transition().duration(750).call(
            zoom.transform,
            d3.zoomIdentity
        );
    };
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

let currentPage = 1;
let hasMore = false;
let lastEventId = '';
let isLoadingLogs = false;

async function loadLogs(reset = false) {
    if (isLoadingLogs) return;
    isLoadingLogs = true;

    try {
        if (reset) {
            currentPage = 1;
            lastEventId = '';
            document.getElementById('logs-list').innerHTML = '';
        }

        const level = document.getElementById('log-level').value;
        const url = new URL('/api/logs', window.location.origin);
        url.searchParams.set('page', currentPage);
        if (level !== 'all') url.searchParams.set('level', level);
        if (lastEventId) url.searchParams.set('last_event_id', lastEventId);

        const response = await fetch(url);
        const data = await response.json();

        // Only proceed if we have new logs
        if (data.logs.length > 0) {
            const lastLog = data.logs[data.logs.length - 1];
            lastEventId = lastLog.timestamp.toString();
            renderLogsPage(data.logs, reset);
        }

        hasMore = data.has_more;
        if (hasMore) {
            currentPage = data.next_page;
        }
    } catch (error) {
        console.error('Error loading logs:', error);
    } finally {
        isLoadingLogs = false;
    }
}

function renderLogsPage(logs, reset) {
    const logsList = document.getElementById('logs-list');
    
    // Create a unique identifier for each log entry
    const existingLogs = new Set(Array.from(logsList.querySelectorAll('.log-entry, .log-group'))
        .map(el => el.getAttribute('data-timestamp')));
    
    // Filter out logs we already have
    const newLogs = logs.filter(log => !existingLogs.has(log.timestamp.toString()));
    
    if (newLogs.length === 0) return;

    const processedLogs = processLogsIntoGroups(newLogs);
    const newLogsHtml = renderLogGroups(processedLogs);

    if (reset) {
        logsList.innerHTML = newLogsHtml;
    } else {
        logsList.insertAdjacentHTML('beforeend', newLogsHtml);
    }

    // Add click handlers for collapsible groups
    logsList.querySelectorAll('.log-group-header:not([data-handler])').forEach(header => {
        header.setAttribute('data-handler', 'true');
        header.addEventListener('click', (e) => {
            e.stopPropagation();
            const group = header.closest('.log-group');
            if (group) {
                group.classList.toggle('collapsed');
            }
        });
    });
}

function processLogsIntoGroups(logs) {
    const groups = [];
    const groupStack = [];
    
    logs.forEach(log => {
        if (log.type === 'begin') {
            // Create new group with proper indentation
            const group = {
                header: log,
                logs: [],
                indent: groupStack.length, // Use stack length for proper nesting
                subgroups: []
            };
            
            // Add to parent group or top level
            if (groupStack.length > 0) {
                groupStack[groupStack.length - 1].subgroups.push(group);
            } else {
                groups.push(group);
            }
            
            groupStack.push(group);
        } else if (log.type === 'end') {
            // Pop the last group from stack
            if (groupStack.length > 0) {
                groupStack.pop();
            }
        } else {
            // Add to current group or top level with proper indentation
            const currentIndent = groupStack.length;
            const logEntry = { ...log, indent: currentIndent };
            
            if (groupStack.length > 0) {
                groupStack[groupStack.length - 1].logs.push(logEntry);
            } else {
                groups.push({ logs: [logEntry] });
            }
        }
    });

    return groups;
}

function renderLogEntry(log) {
    return `
        <div class="log-entry ${log.level.toLowerCase()}" style="padding-left: ${(log.indent + 1) * 20}px">
            <span class="log-timestamp">${new Date(log.timestamp).toLocaleString()}</span>
            <span class="log-level">${log.level}</span>
            <span class="log-message">${escapeHtml(log.message)}</span>
            ${log.trace ? `<pre class="log-trace">${escapeHtml(log.trace)}</pre>` : ''}
        </div>
    `;
}

function escapeHtml(str) {
    const div = document.createElement('div');
    div.textContent = str;
    return div.innerHTML;
}

function renderLogGroups(groups) {
    return groups.map(group => {
        if (!group.header) {
            // Single log entry
            return group.logs.map(log => `
                <div class="log-entry" data-timestamp="${log.timestamp}">
                    ${renderLogEntry(log)}
                </div>
            `).join('');
        }

        // Group with header
        const groupClass = group.header.level.toLowerCase();
        const hasContent = group.logs.length > 0 || group.subgroups.length > 0;
        
        return `
            <div class="log-group" data-timestamp="${group.header.timestamp}">
                <div class="log-group-header ${groupClass}" style="padding-left: ${group.indent * 20}px">
                    <span class="toggle-icon">${hasContent ? '▼' : '▹'}</span>
                    <span class="log-timestamp">${new Date(group.header.timestamp).toLocaleString()}</span>
                    <span class="log-level">${group.header.level}</span>
                    <span class="log-message">${escapeHtml(group.header.message)}</span>
                </div>
                ${hasContent ? `
                    <div class="log-group-content">
                        ${group.logs.map(log => renderLogEntry(log)).join('')}
                        ${group.subgroups.map(subgroup => renderLogGroups([subgroup])).join('')}
                    </div>
                ` : ''}
            </div>
        `;
    }).join('');
}

// Add infinite scroll
const logsView = document.getElementById('logs-view');
logsView?.addEventListener('scroll', () => {
    if (hasMore && !isLoadingLogs && 
        logsView.scrollHeight - logsView.scrollTop <= logsView.clientHeight + 100) {
        loadLogs();
    }
});

// Add real-time updates
setInterval(async () => {
    if (document.getElementById('logs-view').style.display === 'block') {
        const oldScrollHeight = logsView.scrollHeight;
        await loadLogs(false);
        
        // Maintain scroll position if user was not at bottom
        if (logsView.scrollTop < oldScrollHeight - logsView.clientHeight) {
            logsView.scrollTop = logsView.scrollHeight - oldScrollHeight + logsView.scrollTop;
        }
    }
}, 5000); // Check for updates every 5 seconds

// Add log filtering functionality
document.getElementById('log-level')?.addEventListener('change', async (e) => {
    const level = e.target.value;
    await loadLogs(true); // Use true to reset the view
});

document.getElementById('log-search')?.addEventListener('input', (e) => {
    const searchTerm = e.target.value.toLowerCase();
    const logGroups = document.querySelectorAll('.log-group');
    const singleLogs = document.querySelectorAll('.log-entry:not(.log-group .log-entry)');
    
    // Search in single logs
    singleLogs.forEach(entry => {
        const text = entry.textContent.toLowerCase();
        entry.style.display = text.includes(searchTerm) ? 'block' : 'none';
    });

    // Search in groups
    logGroups.forEach(group => {
        const groupContent = group.textContent.toLowerCase();
        const hasMatch = groupContent.includes(searchTerm);
        group.style.display = hasMatch ? 'block' : 'none';
        
        if (hasMatch) {
            // If there's a match, expand the group
            group.classList.remove('collapsed');
        }
    });
});

// Add click handler for the logs tab
document.querySelector('a[data-view="logs"]')?.addEventListener('click', async () => {
    document.querySelectorAll('.view').forEach(view => view.style.display = 'none');
    document.getElementById('logs-view').style.display = 'block';
    await loadLogs(true); // Reset logs when switching to logs tab
});
