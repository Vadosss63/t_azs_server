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

async function sendToAzs(formData, msg) {
    const confirmed = window.confirm(msg);
    if (!confirmed) return;

    showMessage('Отправка запроса на АЗС...');

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

const buttonIdArray = [
    "service1Btn",
    "service2Btn",
    "service3Btn",
    "resetCountersBtn",
    "blockAzsNodeBtn",
    "unblockAzsNodeBtn",
    "priceCash1Btn",
    "priceCashless1Btn",
    "priceCash2Btn",
    "priceCashless2Btn",
    "fuelArrival1Btn",
    "lockFuelValue1Btn",
    "fuelArrival2Btn",
    "lockFuelValue2Btn"
]

const actionMap = {
    "service1Btn": "serviceBtn1",
    "service2Btn": "serviceBtn2",
    "service3Btn": "serviceBtn3",
    "resetCountersBtn": "resetCounters",
    "blockAzsNodeBtn": "blockAzsNode",
    "unblockAzsNodeBtn": "unblockAzsNode",
    "priceCash1Btn": "setPriceCash1",
    "priceCashless1Btn": "setPriceCashless1",
    "priceCash2Btn": "setPriceCash2",
    "priceCashless2Btn": "setPriceCashless2",
    "fuelArrival1Btn": "setFuelArrival1",
    "lockFuelValue1Btn": "setLockFuelValue1",
    "fuelArrival2Btn": "setFuelArrival2",
    "lockFuelValue2Btn": "setLockFuelValue2"
};

const inputMap = {
    "priceCash1Btn": "priceCash1Input",
    "priceCashless1Btn": "priceCashless1Input",
    "priceCash2Btn": "priceCash2Input",
    "priceCashless2Btn": "priceCashless2Input",
    "fuelArrival1Btn": "fuelArrival1Input",
    "lockFuelValue1Btn": "lockFuelValue1Input",
    "fuelArrival2Btn": "fuelArrival2Input",
    "lockFuelValue2Btn": "lockFuelValue2Input"
};

const actionMsgMap = {
    "service1Btn": "Выполнить Снятие Z - отчёта?",
    "service2Btn": "Выполнить Отключение N?",
    "service3Btn": "Выполнить Включение N?",
    "resetCountersBtn": "Выполнить Инкассацию?",
    "blockAzsNodeBtn": "Заблокировать АЗС?",
    "unblockAzsNodeBtn": "Разблокировать АЗС?",
    "priceCash1Btn": "Установить цену для 1-й колонки?",
    "priceCashless1Btn": "Установить цену для 1-й колонки безналичного расчета?",
    "priceCash2Btn": "Установить цену для 2-й колонки",
    "priceCashless2Btn": "Установить цену для 2-й колонки безналичного расчета?",
    "fuelArrival1Btn": "Установить приход для 1-й колонки?",
    "lockFuelValue1Btn": "Установить значение блокировки для 1-й колонки?",
    "fuelArrival2Btn": "Установить приход для 2-й колонки?",
    "lockFuelValue2Btn": "Установить значение блокировки для 2-й колонки?"
};

buttonIdArray.forEach(id => {
    const button = document.getElementById(id);
    if (button) {
        button.addEventListener('click', () => {
            const azsId = document.getElementById('azsId').value;
            const action = actionMap[id];
            const inputId = inputMap[id];

            let value = 0;
            if (inputId) {
                const inputElement = document.getElementById(inputId);
                if (inputElement) {
                    value = convertPriceToInt(inputElement.value);
                }
            }

            const form = new FormData();
            form.append('value', value);
            form.append('pushedBtn', action);
            form.append('id_azs', azsId);

            const msg = actionMsgMap[id];
            sendToAzs(form, msg);
        });
    }
});

["priceCash1Input", "priceCash2Input", "priceCashless1Input", "priceCashless2Input"].forEach(id => {
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
