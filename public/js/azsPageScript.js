class AzsService {
    constructor() {
        this.messageElement = document.getElementById('message');
    }

    showMessage(text) {
        this.messageElement.innerHTML = text;
        this.messageElement.classList.add('loading');
    }

    hideMessage() {
        this.messageElement.classList.remove('loading');
        window.location.reload();
    }

    async sendToAzs(formData, msg) {
        const confirmed = window.confirm(msg);
        if (!confirmed) return;

        this.showMessage('Отправка запроса на АЗС...');

        try {
            const response = await fetch('/push_azs_button', { method: "POST", body: formData });
            if (!response.ok) throw new Error("Сетевая ошибка");

            let retries = 15;
            const checkStatus = async () => {
                try {
                    const statusResponse = await fetch('/azs_button_ready?id_azs=' + formData.get("id_azs"));
                    if (statusResponse.ok) {
                        const data = await statusResponse.json();
                        if (data.status === "ready") {
                            alert("Успешно!");
                            this.hideMessage();
                        } else if (--retries > 0) {
                            setTimeout(checkStatus, 1000);
                        } else {
                            throw new Error("АЗС не отвечает!");
                        }
                    } else {
                        throw new Error("Ошибка сервера при проверке статуса");
                    }
                } catch (error) {
                    alert(`Ошибка: ${error.message}`);
                    this.hideMessage();
                }
            };

            checkStatus();
        } catch (error) {
            alert(`Ошибка: ${error.message}`);
            this.hideMessage();
        }
    }
}

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
        return value.match(/^(0|[1-9]\d*)$/);
    }
}

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

const azsService = new AzsService();

new ButtonAction("service1Btn", "serviceBtn1", "Выполнить Снятие Z - отчёта?", azsService);
new ButtonAction("service2Btn", "serviceBtn2", "Выполнить Отключение N?", azsService);
new ButtonAction("service3Btn", "serviceBtn3", "Выполнить Включение N?", azsService);
new ButtonAction("resetCountersBtn", "resetCounters", "Выполнить Инкассацию?", azsService);
new ButtonAction("blockAzsNodeBtn", "blockAzsNode", "Заблокировать АЗС?", azsService);
new ButtonAction("unblockAzsNodeBtn", "unblockAzsNode", "Разблокировать АЗС?", azsService);
new ButtonAction("priceCash1Btn", "setPriceCash1", "Установить цену для 1-й колонки?", azsService, new InputField("priceCash1Input", new PriceValidator()));
new ButtonAction("priceCashless1Btn", "setPriceCashless1", "Установить цену для 1-й колонки безналичного расчета?", azsService, new InputField("priceCashless1Input", new PriceValidator()));
new ButtonAction("priceCash2Btn", "setPriceCash2", "Установить цену для 2-й колонки?", azsService, new InputField("priceCash2Input", new PriceValidator()));
new ButtonAction("priceCashless2Btn", "setPriceCashless2", "Установить цену для 2-й колонки безналичного расчета?", azsService, new InputField("priceCashless2Input", new PriceValidator()));
new ButtonAction("fuelArrival1Btn", "setFuelArrival1", "Установить приход для 1-й колонки?", azsService, new InputField("fuelArrival1Input", new IntegerValidator()));
new ButtonAction("lockFuelValue1Btn", "setLockFuelValue1", "Установить значение блокировки для 1-й колонки?", azsService, new InputField("lockFuelValue1Input", new IntegerValidator()));
new ButtonAction("fuelArrival2Btn", "setFuelArrival2", "Установить приход для 2-й колонки?", azsService, new InputField("fuelArrival2Input", new IntegerValidator()));
new ButtonAction("lockFuelValue2Btn", "setLockFuelValue2", "Установить значение блокировки для 2-й колонки?", azsService, new InputField("lockFuelValue2Input", new IntegerValidator()));
