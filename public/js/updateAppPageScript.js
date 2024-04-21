// updateAppPageScript.js

import AzsService from './AzsService.js';
import InputField from './InputField.js';
import ButtonAction from './ButtonAction.js';

const pushUrl = "/app_update_button";
const statusUrl = "/app_update_button_ready";
const resetUrl = "/app_update_button_reset";

const azsService = new AzsService(pushUrl, statusUrl, 200, resetUrl);

new ButtonAction("installAppBtn", "install", "Установить обновление?", azsService, new InputField("availableFiles"));
new ButtonAction("deleteAppFileBtn", "delete", "Удалить с сервера обновление?", azsService, new InputField("availableFiles"));
new ButtonAction("downloadAppBtn", "download", "Скачать на сервер обновление?", azsService, new InputField("availableTags"));
