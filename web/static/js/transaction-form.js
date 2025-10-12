document.addEventListener('alpine:init', () => {
    Alpine.data('transactionForm', () => ({
        groupCount: 0,
        accountOptions: '',

        init() {
            this.accountOptions = this.getAccountOptions();
            this.addGroup();
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
                        element.value = 'debit';
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
            const accountOptionsElement = document.getElementById('account-options');
            return accountOptionsElement ? accountOptionsElement.innerHTML : '';
        },

        getCurrentDateTime() {
            const now = new Date();
            const year = now.getFullYear();
            const month = String(now.getMonth() + 1).padStart(2, '0');
            const day = String(now.getDate()).padStart(2, '0');
            const hours = String(now.getHours()).padStart(2, '0');
            const minutes = String(now.getMinutes()).padStart(2, '0');
            return `${year}-${month}-${day}T${hours}:${minutes}`;
        }
    }));
});
