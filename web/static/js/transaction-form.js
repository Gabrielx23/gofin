document.addEventListener('alpine:init', () => {
    Alpine.data('transactionForm', () => ({
        groupCount: 0,
        accountOptions: '',
        transactionTypeOptions: '',
        currencyOptions: '',

        init() {
            this.accountOptions = this.getAccountOptions();
            this.transactionTypeOptions = this.getTransactionTypeOptions();
            this.currencyOptions = this.getCurrencyOptions();
            this.addGroup();
            this.setupAccountCreationHandlers();
        },

        addGroup() {
            this.groupCount++;
            const groupsContainer = document.getElementById('transactionGroups');
            const template = groupsContainer.querySelector('[data-template="true"]');
            const newGroup = template.cloneNode(true);

            newGroup.removeAttribute('data-template');
            newGroup.style.display = 'block';

            this.updateGroupFields(newGroup, this.groupCount);
            groupsContainer.appendChild(newGroup);
            this.updateRemoveButtons();
            this.updateSubmitButton();
        },

        removeGroup(event) {
            event.target.closest('.transaction-group').remove();
            this.updateRemoveButtons();
            this.updateSubmitButton();
        },

        updateRemoveButtons() {
            const groups = document.querySelectorAll('.transaction-group:not([data-template="true"])');
            const removeButtons = document.querySelectorAll('.remove-group-btn');

            removeButtons.forEach((btn) => {
                if (groups.length > 1) {
                    btn.style.display = 'block';
                } else {
                    btn.style.display = 'none';
                }
            });
        },

        updateSubmitButton() {
            const groups = document.querySelectorAll('.transaction-group:not([data-template="true"])');
            const submitBtn = document.querySelector('.create-transaction-button.primary');
            if (submitBtn) {
                if (groups.length === 0) {
                    submitBtn.style.display = 'none';
                } else {
                    submitBtn.style.display = 'inline-block';
                }
            }
        },

        updateGroupFields(groupElement, groupNumber) {
            const index = groupNumber - 1;

            const fields = [
                { selector: 'input[name*="name"]', name: `groups[${index}].name`, id: `name_${index}` },
                { selector: 'input[name*="value"]', name: `groups[${index}].value`, id: `value_${index}` },
                { selector: 'select[name*="type"]', name: `groups[${index}].type`, id: `type_${index}` },
                { selector: 'select[name*="account_id"]', name: `groups[${index}].account_id`, id: `account_${index}` },
                { selector: 'input[name*="date"]', name: `groups[${index}].date`, id: `date_${index}` }
            ];

            fields.forEach(field => {
                const element = groupElement.querySelector(field.selector);
                if (element) {
                    element.name = field.name;
                    element.id = field.id;
                    element.value = '';

                    if (field.selector.includes('name') || field.selector.includes('value') ||
                        field.selector.includes('type') || field.selector.includes('account_id')) {
                        element.required = true;
                    }

                    if (field.selector.includes('date')) {
                        element.value = this.getCurrentDateTime();
                    }
                    if (field.selector.includes('type')) {
                        element.innerHTML = this.transactionTypeOptions;
                        element.value = 'debit';
                    }
                    if (field.selector.includes('account_id')) {
                        element.innerHTML = this.getAccountOptions();
                    }
                }
            });

            const labels = groupElement.querySelectorAll('label');
            labels.forEach(label => {
                const forAttr = label.getAttribute('for');
                if (forAttr && forAttr.includes('_template')) {
                    label.setAttribute('for', forAttr.replace('_template', `_${index}`));
                }
            });

            const removeBtn = groupElement.querySelector('.remove-group-btn');
            if (removeBtn) {
                removeBtn.addEventListener('click', (e) => {
                    this.removeGroup(e);
                });
            }
        },

        getAccountOptions() {
            const template = document.querySelector('[data-template="true"]');
            const accountSelect = template.querySelector('.account-select');
            return accountSelect ? accountSelect.innerHTML : '';
        },

        getTransactionTypeOptions() {
            const template = document.querySelector('[data-template="true"]');
            const typeSelect = template.querySelector('select[name*="type"]');
            return typeSelect ? typeSelect.innerHTML : '';
        },

        getCurrentDateTime() {
            const now = new Date();
            const year = now.getFullYear();
            const month = String(now.getMonth() + 1).padStart(2, '0');
            const day = String(now.getDate()).padStart(2, '0');
            const hours = String(now.getHours()).padStart(2, '0');
            const minutes = String(now.getMinutes()).padStart(2, '0');
            return `${year}-${month}-${day}T${hours}:${minutes}`;
        },

        getCurrencyOptions() {
            const template = document.querySelector('[data-template="true"]');
            const currencySelect = template.querySelector('.new-account-currency');
            return currencySelect ? currencySelect.innerHTML : '';
        },

        setupAccountCreationHandlers() {
            document.addEventListener('click', (e) => {
                if (e.target.classList.contains('add-account-btn')) {
                    this.showCreateAccountForm(e.target);
                } else if (e.target.classList.contains('create-account-btn')) {
                    this.createAccount(e.target);
                } else if (e.target.classList.contains('cancel-create-account-btn') ||
                    e.target.classList.contains('create-account-backdrop')) {
                    this.hideCreateAccountForm(e.target);
                }
            });

            document.addEventListener('keydown', (e) => {
                if (e.key === 'Escape') {
                    const visibleForm = document.querySelector('.create-account-section[style*="block"]');
                    if (visibleForm) {
                        this.hideCreateAccountForm(visibleForm);
                    }
                }
            });
        },

        showCreateAccountForm(button) {
            const container = button.closest('.account-select-container');
            const createSection = container.querySelector('.create-account-section');
            const backdrop = container.querySelector('.create-account-backdrop');
            const select = container.querySelector('.account-select');

            backdrop.style.display = 'block';
            createSection.style.display = 'block';
            select.style.display = 'none';
            button.style.display = 'none';

            document.body.style.overflow = 'hidden';
        },

        hideCreateAccountForm(button) {
            const container = button.closest('.account-select-container');
            const createSection = container.querySelector('.create-account-section');
            const backdrop = container.querySelector('.create-account-backdrop');
            const select = container.querySelector('.account-select');
            const addBtn = container.querySelector('.add-account-btn');

            backdrop.style.display = 'none';
            createSection.style.display = 'none';
            select.style.display = 'block';
            addBtn.style.display = 'block';

            const nameInput = container.querySelector('.new-account-name');
            nameInput.value = '';

            document.body.style.overflow = '';

            this.hideAccountError(container);
        },

        async createAccount(button) {
            const container = button.closest('.account-select-container');
            const nameInput = container.querySelector('.new-account-name');
            const currencySelect = container.querySelector('.new-account-currency');
            const select = container.querySelector('.account-select');

            const name = nameInput.value.trim();
            const currency = currencySelect.value;

            if (!name) {
                this.showAccountError(container, 'Please enter an account name');
                return;
            }

            try {
                const projectSlug = window.location.pathname.split('/')[1];
                const response = await fetch(`/${projectSlug}/accounts/create`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        name: name,
                        currency: currency
                    })
                });

                const result = await response.json();

                if (!response.ok || result.error) {
                    this.showAccountError(container, result.error || 'Failed to create account');
                    return;
                }

                this.addAccountToAllSelects(result);
                this.hideCreateAccountForm(button);

            } catch (error) {
                this.showAccountError(container, 'Network error. Please try again.');
            }
        },

        addAccountToAllSelects(account) {
            this.updateAccountOptionsTemplate(account);

            const option = document.createElement('option');
            option.value = account.id;
            option.textContent = `${account.name} (${account.currency})`;
            option.selected = true;

            const existingSelects = document.querySelectorAll('.account-select:not([data-template="true"] .account-select)');
            existingSelects.forEach(select => {
                const newOption = option.cloneNode(true);
                select.appendChild(newOption);
            });
        },

        updateAccountOptionsTemplate(account) {
            const template = document.querySelector('[data-template="true"]');
            const templateSelect = template.querySelector('.account-select');
            if (templateSelect) {
                const option = document.createElement('option');
                option.value = account.id;
                option.textContent = `${account.name} (${account.currency})`;
                templateSelect.appendChild(option);
            }
        },

        showAccountError(container, message) {
            this.hideAccountError(container);

            const errorDiv = document.createElement('div');
            errorDiv.className = 'account-error-message';
            errorDiv.textContent = message;
            errorDiv.style.cssText = `
                color: #dc3545;
                font-size: 0.9rem;
                margin-top: 0.5rem;
                padding: 0.5rem;
                background: #f8d7da;
                border: 1px solid #f5c6cb;
                border-radius: 4px;
            `;

            container.appendChild(errorDiv);
        },

        hideAccountError(container) {
            const existingError = container.querySelector('.account-error-message');
            if (existingError) {
                existingError.remove();
            }
        }
    }));
});
