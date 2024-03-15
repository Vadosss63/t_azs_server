const message = document.getElementById('message');

function priceValidator(priceInput) {
    const { value, defaultValue } = priceInput;
    if (!value.match(/^(0|[1-9]\d*)(\.[0-9]{1,2})?$/) || parseFloat(value) > 200.99) {
        priceInput.value = defaultValue;
    } else {
        priceInput.defaultValue = value;
    }
}

function floatValidator(priceInput) {
    const { value, defaultValue } = priceInput;
    if (!value.match(/^(0|[1-9]\d*)(\.[0-9]{1,2})?$/)) {
        priceInput.value = defaultValue;
    } else {
        priceInput.defaultValue = value;
    }
}

function convertPriceToInt(price) {
    return Math.round(price * 100);
}

function showMessage(text) {
    message.innerHTML = text;
    message.classList.add('loading');
}

function hideMessage() {
    message.classList.remove('loading');
    window.location.reload();
}

async function sendToAzs(form, msg, value = 0) {
    const confirmed = window.confirm(msg);
    if (!confirmed) return;

    const formData = new FormData(form);
    formData.set("value", value);

    showMessage('Отправка запроса на АЗС...');

    try {
        const response = await fetch(form.action, { method: "POST", body: formData });
        if (!response.ok) throw new Error("Сетевая ошибка");

        let retries = 15;
        const checkStatus = async () => {
            try {
                const statusResponse = await fetch('/azs_button_ready?id_azs=' + formData.get("id_azs"));
                if (statusResponse.ok) {
                    const data = await statusResponse.json();
                    if (data.status === "ready") {
                        alert("Успешно!");
                        hideMessage();
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
                hideMessage();
            }
        };

        checkStatus();
    } catch (error) {
        alert(`Ошибка: ${error.message}`);
        hideMessage();
    }
}

const selectorsArray =  [
    ".resetAzs1",
    ".resetAzs2",
    ".resetDallyCounter",
    ".serviceBtn1",
    ".serviceBtn2",
    ".serviceBtn3",
    ".priceBtn1",
    ".priceBtn2",
    ".priceBtn3",
    ".priceBtn4",
    ".fuelArrivalBtn1",
    ".fuelArrivalBtn2",
    ".lockFuelValueBtn1",
    ".lockFuelValueBtn2"
  ]

document.querySelectorAll(selectorsArray).forEach(form => {
    form.addEventListener("submit", (event) => {
        event.preventDefault();
        const inputMap = {
            'priceBtn1': "price1Input",
            'priceBtn2': "price2Input",
            'priceBtn3': "price1cashlessInput",
            'priceBtn4': "price2cashlessInput",
            'fuelArrivalBtn1': "fuelArrival1Input",
            'fuelArrivalBtn2': "fuelArrival2Input",
            'lockFuelValueBtn1': "lockFuelValue1Input",
            'lockFuelValueBtn2': "lockFuelValue2Input"
        };

        const actionMsgMap = {
            'resetAzs1': "Заблокировать АЗС?",
            'resetAzs2': "Разблокировать АЗС?",
            'resetDallyCounter': "Выполнить Инкассацию?",
            'serviceBtn1': "Выполнить Снятие Z - отчёта?",
            'serviceBtn2': "Выполнить Отключение N?",
            'serviceBtn3': "Выполнить Включение N?",
            'priceBtn1': "Установить цену для 1-й колонки?",
            'priceBtn2': "Установить цену для 2-й колонки?",
            'priceBtn3': "Установить цену для 1-й колонки безналичного расчета?",
            'priceBtn4': "Установить цену для 2-й колонки безналичного расчета?",
            'fuelArrivalBtn1': "Установить приход для 1-й колонки?",
            'fuelArrivalBtn2': "Установить приход для 2-й колонки?",
            'lockFuelValueBtn1': "Установить значение блокировки для 1-й колонки?",
            'lockFuelValueBtn2': "Установить значение блокировки для 2-й колонки?"
        };

        const className = form.classList[0];

        const inputId = inputMap[className];

        let value = 0;
        if (inputId) {
            const inputElement = document.getElementById(inputId);
            if (inputElement) {
                value = convertPriceToInt(inputElement.value);
            }
        }

        const msg = actionMsgMap[className];

        sendToAzs(form, msg, value);
    });
});

["price1Input", "price2Input", "price1cashlessInput", "price2cashlessInput"].forEach(id => {
    const inputElement = document.getElementById(id);
    if (inputElement) {
        inputElement.addEventListener("input", function () {
            priceValidator(this);
        });
    }
});

["fuelArrival1Input", "fuelArrival2Input", "lockFuelValue1Input", "lockFuelValue2Input"].forEach(id => {
    const inputElement = document.getElementById(id);
    if (inputElement) {
        inputElement.addEventListener("input", function () {
            floatValidator(this);
        });
    }
});
