// InputField.js
import { Validator } from './Validator.js';

class InputField {
    constructor(id, validator = null) {
        this.id = id;
        this.validator = validator;
        this.inputElement = document.getElementById(id);
        this.initValidation();
    }

    initValidation() {
        if (this.inputElement && this.validator) {
            this.validator.setInputElement(this.inputElement);
            this.inputElement.addEventListener("input", () => this.validator.validate());
        }
    }

    getValue() {
        if (!this.inputElement) return null;
        return this.validator ? this.validator.getValue() : this.inputElement.value;
    }
}

export default InputField;
