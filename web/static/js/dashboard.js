document.addEventListener('alpine:init', () => {
    Alpine.data('dashboard', () => ({
        deleteRoute: '',
        projectSlug: '',

        deleteTransaction(transactionId) {
            if (confirm('Are you sure you want to delete this transaction?')) {
                const form = document.createElement('form');
                form.method = 'POST';
                form.action = '/' + this.projectSlug + this.deleteRoute + '?id=' + transactionId;

                document.body.appendChild(form);
                form.submit();
            }
        }
    }));
});
