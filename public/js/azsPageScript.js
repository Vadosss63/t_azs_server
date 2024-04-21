//azsPageScript.js
import AzsService from './AzsService.js';
import { PriceValidator, IntegerValidator } from './Validator.js';
import InputField from './InputField.js';
import ButtonAction from './ButtonAction.js';

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
