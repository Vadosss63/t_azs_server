// Validator.js
class Validator {
    setInputElement(inputElement) {
        this.inputElement = inputElement;
    }

    validate() {
        const { value, defaultValue } = this.inputElement;
        if (!this.isValid(value)) {
            this.inputElement.value = defaultValue;
        } else {
            this.inputElement.defaultValue = value;
        }
    }

    isValid(value) {
        throw new Error("Method 'isValid(value)' must be implemented.");
    }

    getValue() {
        return this.inputElement.value;
    }
}

class PriceValidator extends Validator {
    isValid(value) {
        return value.match(/^(0|[1-9]\d*)(\.[0-9]{1,2})?$/) && parseFloat(value) <= 200.99;
    }

    getValue() {
        return this.convertPriceToInt(this.inputElement.value);
    }

    convertPriceToInt(price) {
        return Math.round(price * 100);
    }
}

class IntegerValidator extends Validator {
    isValid(value) {
        if (!value.match(/^(-?(0|[1-9]\d*))$/)) {
            return false;
        }
        const number = parseInt(value, 10);
        return number >= -100000 && number <= 100000;
    }
}

export { Validator, PriceValidator, IntegerValidator };
