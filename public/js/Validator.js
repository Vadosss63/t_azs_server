// Validator.js
class Validator {
    setInputElement(inputElement) {
        this.inputElement = inputElement;
    }

    validate() {
        const { value, defaultValue } = this.inputElement;
        const { isValid, validatedValue } = this.isValid(value);

        if (!isValid) {
            this.inputElement.value = defaultValue;
        } else {
            this.inputElement.value = validatedValue;
            this.inputElement.defaultValue = validatedValue;
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
        value = String(value).replace(',', '.');
        value = value.replace(/^0+(?!$)/, '');   // remove leading zeros

        const originalValue = value;

        if (value.trim() === '') {
            return {
                isValid: true,
                validatedValue: ''
            };
        }

        const isTested = /^(\d+)(\.[0-9]{0,2})?$/.test(value);

        if (!isTested) {
            return { isValid: false, validatedValue: null };
        }

        const floatVal = parseFloat(value);
        const isValid = floatVal <= 200.99;

        return {
            isValid,
            validatedValue: isValid ? originalValue : null
        };
    }

    getValue() {
        return this.convertPriceToInt(this.inputElement.value);
    }

    convertPriceToInt(price) {
        if (price.trim() === '') {
            return 0;
        }
        price = String(price).replace(',', '.');
        return Math.round(parseFloat(price) * 100);
    }
}

class IntegerValidator extends Validator {

    isValid(value) {

        value = value.replace(/^0+(?!$)/, '');   // remove leading zeros
        const originalValue = value;

        if (value.trim() === '-') {
            return {
                isValid: true,
                validatedValue: '-'
            };
        }


        if (value.trim() === '') {
            return {
                isValid: true,
                validatedValue: ''
            };
        }

        if (typeof value !== 'string' || !/^(-?\d+)$/.test(value)) {
            return { isValid: false, validatedValue: null };
        }

        const number = Number(value);
        const isValid = number >= -100000 && number <= 100000;

        return { isValid, validatedValue: originalValue };
    }

    getValue() {
        const value = this.inputElement.value;
        if (value.trim() === '' || value.trim() === '-') {
            return 0;
        }
        return parseInt(value, 10);
    }
}

export { Validator, PriceValidator, IntegerValidator };
