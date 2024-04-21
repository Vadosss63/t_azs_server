class AzsService {
    constructor(pushUrl, statusUrl, maxRetries = 15, resetUrl = null) {
        this.pushUrl = pushUrl;
        this.statusUrl = statusUrl;
        this.maxRetries = maxRetries;
        this.resetUrl = resetUrl;
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
            const response = await fetch(this.pushUrl, { method: "POST", body: formData });
            if (!response.ok) throw new Error("Сетевая ошибка");
            await this.checkStatus(formData);
        } catch (error) {
            alert(`Ошибка: ${error.message}`);
            this.hideMessage();
        }
    }

    async checkStatus(formData) {
        let retries = this.maxRetries;
        const check = async () => {
            try {
                const statusResponse = await fetch(`${this.statusUrl}?id_azs=${formData.get("id_azs")}`);
                if (!statusResponse.ok) throw new Error("Ошибка сервера при проверке статуса");

                const data = await statusResponse.json();
                if (data.status === "ready") {
                    alert("Успешно!");
                    this.hideMessage();
                } else if (--retries > 0) {
                    setTimeout(check, 1000);
                } else {
                    if(this.resetUrl)
                    {
                        const status = await fetch(`${this.resetUrl}?id=${formData.get("id_azs")}`);
                        if (!status.ok) throw new Error("Ошибка сервера при проверке статуса");
                    }
                    throw new Error("АЗС не отвечает!");
                }
            } catch (error) {
                alert(`Ошибка: ${error.message}`);
                this.hideMessage();
            }
        };
        check();
    }
}

export default AzsService;
