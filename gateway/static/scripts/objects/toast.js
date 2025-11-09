class Toast {
    constructor(opts = {}) {
        this.position = opts.position || 'top-right';
        this.container = document.getElementById('toastContainer') || this._createContainer();
    }

    _createContainer() {
        const container = document.createElement('div');
        container.id = 'toastContainer';
        container.className = 'toast-container';
        document.body.appendChild(container);
        return container;
    }

    _formatMessage(msg) {
        if (msg === null || msg === undefined) return '';
        if (typeof msg === 'string') return this._escapeHtml(msg);
        try {
            return `<pre style="white-space:pre-wrap;margin:0">${this._escapeHtml(JSON.stringify(msg, null, 2))}</pre>`;
        } catch (e) {
            return this._escapeHtml(String(msg));
        }
    }

    _escapeHtml(str) {
        return String(str)
            .replace(/&/g, '&amp;')
            .replace(/</g, '&lt;')
            .replace(/>/g, '&gt;');
    }

    success(msg) {
        return this.show('success', msg);
    }

    error(msg) {
        return this.show('error', msg);
    }

    show(type, msg = '', timeout = 4000) {
        const text = this._formatMessage(msg);
        const t = document.createElement('div');
        t.className = `toast toast--${type}`;
        t.setAttribute('role', 'status');
        t.innerHTML = `<div class="toast__body">${text}</div>
                   <button class="toast__close">×</button>`;

        const closeBtn = t.querySelector('.toast__close');
        closeBtn.addEventListener('click', () => this._hideToast(t));

        this.container.appendChild(t);
        
        requestAnimationFrame(() => t.classList.add('show'));

        t._hideTimer = setTimeout(() => this._hideToast(t), timeout);

        return t;
    }

    _hideToast(t) {
        if (!t) return;
        clearTimeout(t._hideTimer);
        t.classList.remove('show');
        setTimeout(() => {
            if (t.parentNode) t.parentNode.removeChild(t);
        }, 240);
    }
}