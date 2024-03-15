const message = document.getElementById('message');

function priceValidator(priceInput) {
    const newValue = priceInput.value;
    const oldValue = priceInput.defaultValue;
    if (!newValue.match(/^(0|[1-9]\d*)(\.[0-9]{1,2})?$/)) {
        priceInput.value = oldValue;
        return;
    };

    const priceFloat = parseFloat(newValue);
    if (priceFloat > 200.99) {
        priceInput.value = oldValue;
        return;
    }
    priceInput.defaultValue = newValue;
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

function handleXhrLoad(handlerObject) {
    if (handlerObject.xhr.status === 200) {
        if (handlerObject.noAnswer) {
            alert("АЗС не отвечает!");
        } else {
            alert("Успешно!");
        }
        hideMessage();

    } else if (handlerObject.retryCount < handlerObject.retryTime / 1000) {
        handlerObject.retryTimer = setTimeout(function () {
            handlerObject.xhr.open('GET', handlerObject.isReadyUrl);
            handlerObject.xhr.send();
            handlerObject.retryCount++;
        }, 1000);
    } else if (!handlerObject.noAnswer) {
        handlerObject.noAnswer = true;
        handlerObject.xhr.open('GET', handlerObject.resetUrl);
        handlerObject.xhr.send();
    } else {
        hideMessage();
    }
}

function sendToAzs(form, msg, price = 0) {
    const confirmed = window.confirm(msg);
    if (!confirmed) {
        return;
    }

    const formData = new FormData(form);
    const idAzs = formData.get("id_azs");
    formData.set("price", price);

    showMessage('Отправка запроса на АЗС...');
    fetch(form.action, {
        method: "POST",
        body: formData
    }).then(() => {
        const handlerObject = {
            retryTime: 15000,
            noAnswer: false,
            retryCount: 0,
            idAzs: idAzs,
            retryTimer: null,
            xhr: new XMLHttpRequest(),
            isReadyUrl: '/azs_button_ready?id_azs=' + idAzs,
            resetUrl: '/reset_azs_button?id=' + idAzs,
        };

        handlerObject.xhr.open('GET', handlerObject.isReadyUrl);
        handlerObject.xhr.onload = handleXhrLoad.bind(null, handlerObject);
        handlerObject.xhr.send();
    }).catch((error) => {
        // handle error
        if (error.response.status === 400) {
            alert("Ошибка: неправильный запрос");
        } else {
            alert("Ошибка: " + error);
        }
        hideMessage();
    });
}

const resetAzs1 = document.querySelectorAll(".resetAzs1");
resetAzs1.forEach((form) => {
    form.addEventListener("submit", (event) => {
        event.preventDefault();
        sendToAzs(form, "Заблокировать АЗС?");
    });
});

const resetAzs2 = document.querySelectorAll(".resetAzs2");
resetAzs2.forEach((form) => {
    form.addEventListener("submit", (event) => {
        event.preventDefault();
        sendToAzs(form, "Разблокировать АЗС?");
    });
});

const resetDallyCounter = document.querySelectorAll(".resetDallyCounter");
resetDallyCounter.forEach((form) => {
    form.addEventListener("submit", (event) => {
        event.preventDefault();
        sendToAzs(form, "Выполнить Инкассацию?");
    });
});

const serviceBtn1 = document.querySelectorAll(".serviceBtn1");
serviceBtn1.forEach((form) => {
    form.addEventListener("submit", (event) => {
        event.preventDefault();
        sendToAzs(form, "Выполнить Снятие Z - отчёта?");
    });
});

const serviceBtn2 = document.querySelectorAll(".serviceBtn2");
serviceBtn2.forEach((form) => {
    form.addEventListener("submit", (event) => {
        event.preventDefault();
        sendToAzs(form, "Выполнить Отключение N?");
    });
});

const serviceBtn3 = document.querySelectorAll(".serviceBtn3");
serviceBtn3.forEach((form) => {
    form.addEventListener("submit", (event) => {
        event.preventDefault();
        sendToAzs(form, "Выполнить Включение N?");
    });
});

const price1Input = document.getElementById("price1Input");
const priceBtn1 = document.querySelectorAll(".priceBtn1");
priceBtn1.forEach((form) => {
    form.addEventListener("submit", (event) => {
        event.preventDefault();
        sendToAzs(form, "Установить цену для 1-й колонки?", convertPriceToInt(price1Input.value));
    });
});

price1Input.addEventListener("input", function () {
    priceValidator(this);
});

const priceBtn3 = document.querySelectorAll(".priceBtn3");
const price1cashlessInput = document.getElementById("price1cashlessInput");
if (priceBtn3 !== null && priceBtn3 !== undefined) {
    priceBtn3.forEach((form) => {
        form.addEventListener("submit", (event) => {
            event.preventDefault();
            sendToAzs(form, "Установить цену для 1-й колонки?", convertPriceToInt(price1cashlessInput.value));
        });
    });
}
if (price1cashlessInput !== null && price1cashlessInput !== undefined) {
    price1cashlessInput.addEventListener("input", function () {
        priceValidator(this);
    });
}

const priceBtn2 = document.querySelectorAll(".priceBtn2");
const price2Input = document.getElementById("price2Input");
if (priceBtn2 !== null && priceBtn2 !== undefined) {
    priceBtn2.forEach((form) => {
        form.addEventListener("submit", (event) => {
            event.preventDefault();
            sendToAzs(form, "Установить цену для 2-й колонки?", convertPriceToInt(price2Input.value));
        });
    });
}
if (price2Input !== null && price2Input !== undefined) {
    price2Input.addEventListener("input", function () {
        priceValidator(this);
    });
}

const priceBtn4 = document.querySelectorAll(".priceBtn4");
const price2cashlessInput = document.getElementById("price2cashlessInput");
if (priceBtn4 !== null && priceBtn4 !== undefined) {
    priceBtn4.forEach((form) => {
        form.addEventListener("submit", (event) => {
            event.preventDefault();
            sendToAzs(form, "Установить цену для 2-й колонки?", convertPriceToInt(price2cashlessInput.value));
        });
    });
}
if (price2cashlessInput !== null && price2cashlessInput !== undefined) {
    price2cashlessInput.addEventListener("input", function () {
        priceValidator(this);
    });
}