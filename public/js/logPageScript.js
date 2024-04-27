// updateAppPageScript.js

import AzsService from './AzsService.js';
import ButtonAction from './ButtonAction.js';

const pushUrl = "/log_button";
const statusUrl = "/log_button_ready";
const resetUrl = "/log_button_reset";

const azsService = new AzsService(pushUrl, statusUrl, 25, resetUrl);

new ButtonAction("downloadLogsBtn", "download", "Получить логи с АЗС?", azsService);
new ButtonAction("deleteLogsBtn", "delete", "Удалить все логи?", azsService);
