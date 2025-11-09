class Modal {
    constructor(overlay, options = {}) {
        if (!(overlay instanceof HTMLElement)) throw new Error('overlay should be a HTMLElement');

        this.overlay = overlay;
        this.modal = overlay.querySelector('.modal');
        this.form = overlay.querySelector('form');
        this.message = overlay.querySelector('.modal__message');
        this.lastFocused = null;

        this._onKeyDown = this._onKeyDown.bind(this);
        this._onOverlayClick = this._onOverlayClick.bind(this);
        this._onModalClick = (e) => e.stopPropagation();

        this._onSubmit = this._onSubmit.bind(this);
        this._boundCloseBtnsHandler = this._boundCloseBtnsHandler?.bind(this);

        this.firstFocusSelector = options.firstFocusSelector || 'input,button,[tabindex]:not([tabindex="-1"])';
        this.onSubmit = (typeof options.onSubmit === 'function') ? options.onSubmit : null;

        this._modalMessage = this.message || this._createModalMessage();
        this.toast = options.toast instanceof Toast ? options.toast : new Toast();

        this._submitBtn = this.form ? this.form.querySelector('[type="submit"]') : null;
        this._closeBtns = Array.from(this.overlay.querySelectorAll('[data-modal-close], .btn--close, [data-action="close"]'));

        this._handleCloseBtnClick = this._handleCloseBtnClick.bind(this);
        this._closeBtns.forEach(b => b.addEventListener('click', this._handleCloseBtnClick));
    }

    open() {
        if (this.form) this.form.reset();

        if (this.message) {
            this.message.textContent = '';
            this.message.style.display = 'none';
        }

        this.lastFocused = document.activeElement;
        this.overlay.classList.add('open');
        document.documentElement.style.overflow = 'hidden';

        requestAnimationFrame(() => {
            const first = this.modal.querySelector(this.firstFocusSelector);
            if (first) first.focus();
            this.modal.focus?.();
        });

        document.addEventListener('keydown', this._onKeyDown);
        this.overlay.addEventListener('click', this._onOverlayClick);
        this.modal.addEventListener('click', this._onModalClick);

        if (this.form) this.form.addEventListener('submit', this._onSubmit);
    }

    close() {
        this.overlay.classList.remove('open');
        document.documentElement.style.overflow = '';

        document.removeEventListener('keydown', this._onKeyDown);
        this.overlay.removeEventListener('click', this._onOverlayClick);
        this.modal.removeEventListener('click', this._onModalClick);

        if (this.form) this.form.removeEventListener('submit', this._onSubmit);

        if (this.lastFocused && typeof this.lastFocused.focus === 'function') {
            this.lastFocused.focus();
        }
    }

    _handleCloseBtnClick() {
        this.close();
    }

    showInlineMessage(type = 'success', text, timeout = 30000) {
        const el = this._modalMessage || this.message;
        if (!el) return;

        clearTimeout(el._hideTimer);

        el.className = 'modal__message';
        el.style.display = 'block';
        el.hidden = false;

        if (type === 'error') {
            el.classList.add('modal__message--error');
        } else {
            el.classList.add('modal__message--success');
        }

        el.textContent = String(text);

        if (timeout && timeout > 0) {
            el._hideTimer = setTimeout(() => {
                this.hideInlineMessage();
            }, timeout);
        }
    }

    hideInlineMessage() {
        const el = this._modalMessage || this.message;
        if (!el) return;
        clearTimeout(el._hideTimer);
        el.textContent = '';
        el.hidden = true;
        el.style.display = 'none';
        el.className = 'modal-message';
    }

    _onKeyDown(e) {
        if (e.key === 'Escape') {
            this.close();
            return;
        }
        if (e.key === 'Tab') {
            this._maintainFocus(e);
        }
    }

    _onOverlayClick(e) {
        if (e.target === this.overlay) {
            this.close();
        }
    }

    _getFocusable() {
        const selectors = 'a[href], area[href], input:not([disabled]):not([type=hidden]), select:not([disabled]), textarea:not([disabled]), button:not([disabled]), [tabindex]:not([tabindex="-1"])';
        return Array.from(this.modal.querySelectorAll(selectors)).filter(el => el.offsetParent !== null);
    }

    _maintainFocus(e) {
        const focusable = this._getFocusable();
        if (focusable.length === 0) { e.preventDefault(); return; }
        const first = focusable[0], last = focusable[focusable.length - 1];
        if (e.shiftKey && document.activeElement === first) {
            e.preventDefault(); last.focus();
        } else if (!e.shiftKey && document.activeElement === last) {
            e.preventDefault(); first.focus();
        }
    }

    async _onSubmit(e) {
        e.preventDefault();
        this.hideInlineMessage();

        if (!this.form) return;

        const formData = new FormData(this.form);
        const payload = {};
        for (const [k, v] of formData.entries()) payload[k] = v;

        const firstInput = this.modal.querySelector('input[type="text"], input:not([type])');
        if (firstInput && !firstInput.value.trim()) {
            this.showInlineMessage('error', 'Пожалуйста, введите текст.');
            firstInput.focus();
            return;
        }

        if (this._submitBtn) {
            this._submitBtn.disabled = true;
            this._submitBtn.setAttribute('aria-busy', 'true');
        }

        try {
            const result = await this.onSubmit(payload);

            if (!result) {
                throw new Error('Error: empty response');
            }

            this.toast.show('success', result);
            this.showInlineMessage('success', result);

            return result;
        } catch (err) {
            const message = err && err.message ? err.message : 'Не удалось выполнить операцию';
            this.showInlineMessage('error', message);
            this.toast.show('error', message);
            throw err;
        } finally {
            if (this._submitBtn) {
                this._submitBtn.disabled = false;
                this._submitBtn.removeAttribute('aria-busy');
            }
        }
    }

    _createModalMessage() {
        const el = document.createElement('div');
        el.id = 'modalMessage';
        el.className = 'modal__message';
        el.hidden = true;

        if (this.form && this.form.parentNode) {
            this.form.parentNode.insertBefore(el, this.form.nextSibling);
        } else {
            this.modal.appendChild(el);
        }
        return el;
    }
}
