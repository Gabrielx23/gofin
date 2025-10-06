document.addEventListener('DOMContentLoaded', function () {
    const digitInputs = document.querySelectorAll('input[type="text"]');

    digitInputs.forEach((input, index, inputs) => {
        input.addEventListener('input', function (e) {
            if (e.target.value.length === 1 && index < inputs.length - 1) {
                inputs[index + 1].focus();
            }
        });

        input.addEventListener('keydown', function (e) {
            if (e.key === 'Backspace' && e.target.value === '' && index > 0) {
                inputs[index - 1].focus();
            }
        });
    });
});
