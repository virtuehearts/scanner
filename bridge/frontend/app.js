let apiKey = localStorage.getItem('fortune_api_key') || '';
document.getElementById('api-key-input').value = apiKey;

function saveApiKey() {
    apiKey = document.getElementById('api-key-input').value;
    localStorage.setItem('fortune_api_key', apiKey);
    alert('API Key saved and applied');
    loadFiles();
    updateStatus();
}

async function apiCall(endpoint, method = 'GET', body = null) {
    const options = {
        method,
        headers: {
            'x-api-key': apiKey,
            'Content-Type': 'application/json'
        }
    };
    if (body) options.body = JSON.stringify(body);

    try {
        const response = await fetch(endpoint, options);
        if (response.status === 403) {
            console.warn('Invalid API Key');
            return null;
        }
        return await response.json();
    } catch (error) {
        console.error('API Call failed:', error);
        return null;
    }
}

async function loadFiles() {
    const data = await apiCall('/api/files');
    if (data && data.files) {
        const container = document.getElementById('file-list');
        container.innerHTML = '';
        data.files.forEach(file => {
            const div = document.createElement('div');
            div.className = 'flex items-center space-x-2 p-1 hover:bg-slate-800 rounded';
            div.innerHTML = `
                <input type="checkbox" value="${file}" class="addr-file-checkbox w-4 h-4 rounded bg-slate-700 border-slate-600 text-blue-500">
                <span class="text-xs text-slate-300 truncate">${file}</span>
            `;
            container.appendChild(div);
        });
    }
}

async function updateStatus() {
    const data = await apiCall('/api/status');
    if (!data) return;

    const dot = document.getElementById('status-dot');
    const text = document.getElementById('status-text');
    const startBtn = document.getElementById('start-btn');
    const stopBtn = document.getElementById('stop-btn');

    if (data.running) {
        dot.className = 'w-3 h-3 rounded-full bg-emerald-500 mr-2 animate-pulse';
        text.innerText = 'RUNNING';
        text.className = 'text-sm font-medium uppercase tracking-wider text-emerald-500';
        startBtn.classList.add('hidden');
        stopBtn.classList.remove('hidden');
    } else {
        dot.className = 'w-3 h-3 rounded-full bg-slate-500 mr-2';
        text.innerText = 'OFFLINE';
        text.className = 'text-sm font-medium uppercase tracking-wider text-slate-400';
        startBtn.classList.remove('hidden');
        stopBtn.classList.add('hidden');
    }

    document.getElementById('stat-iops').innerText = (data.iops || 0).toLocaleString();
    document.getElementById('stat-tried').innerText = (data.total_tried || 0).toLocaleString();

    if (data.found_keys && data.found_keys.length > 0) {
        renderFoundKeys(data.found_keys);
    }
}

function renderFoundKeys(keys) {
    const container = document.getElementById('found-keys-container');
    container.innerHTML = '';
    keys.forEach(item => {
        const div = document.createElement('div');
        div.className = 'bg-emerald-900/20 border border-emerald-500/30 rounded-lg p-3';
        div.innerHTML = `
            <div class="text-[10px] text-emerald-500 font-bold mb-1 uppercase">${item.timestamp}</div>
            <div class="mono text-xs text-emerald-100 break-all">${item.key_info}</div>
        `;
        container.prepend(div);
    });
}

async function loadLogs() {
    const data = await apiCall('/api/logs');
    if (data && data.logs) {
        const container = document.getElementById('logs-container');
        // Only update if we have new logs to avoid flickering
        const newLogs = data.logs.join('');
        if (container.dataset.lastLog !== newLogs) {
            container.innerHTML = data.logs.map(line => `<div>${line}</div>`).join('');
            container.scrollTop = container.scrollHeight;
            container.dataset.lastLog = newLogs;
        }
    }
}

async function startScanner() {
    const selectedFiles = Array.from(document.querySelectorAll('.addr-file-checkbox:checked')).map(cb => cb.value);
    const config = {
        command: document.getElementById('config-command').value,
        workers: parseInt(document.getElementById('config-workers').value),
        files: selectedFiles,
        night: document.getElementById('config-night').checked,
        telegram_token: document.getElementById('config-tg-token').value || null,
        telegram_channel: document.getElementById('config-tg-channel').value || null
    };

    const result = await apiCall('/api/start', 'POST', config);
    if (result) {
        updateStatus();
    }
}

async function stopScanner() {
    const result = await apiCall('/api/stop', 'POST');
    if (result) {
        updateStatus();
    }
}

function clearLogs() {
    document.getElementById('logs-container').innerHTML = '';
}

// Initial loads
loadFiles();
updateStatus();
loadLogs();

// Intervals
setInterval(updateStatus, 1000);
setInterval(loadLogs, 2000);
