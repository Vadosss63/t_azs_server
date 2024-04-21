// ButtonAction.js
import AzsService from './AzsService.js';
import InputField from './InputField.js';

class ButtonAction {
    constructor(id, action, message, azsService, inputField = null) {
        this.id = id;
        this.action = action;
        this.message = message;
        this.azsService = azsService;
        this.inputField = inputField;
        this.buttonElement = document.getElementById(id);
        this.initButton();
    }

    initButton() {
        if (this.buttonElement) {
            this.buttonElement.addEventListener('click', () => this.handleClick());
        }
    }

    handleClick() {
        const azsId = document.getElementById('azsId').value;
        let value = this.inputField ? this.inputField.getValue() : 0;

        const form = new FormData();
        form.append('value', value);
        form.append('pushedBtn', this.action);
        form.append('id_azs', azsId);

        this.azsService.sendToAzs(form, this.message);
    }
}

export default ButtonAction;
