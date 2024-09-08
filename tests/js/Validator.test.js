//How to run:
// replace in public/js/Validator.js:
// #export { Validator, PriceValidator, IntegerValidator };
// to
// #module.exports = { Validator, PriceValidator, IntegerValidator };
// run: node tests/js/Validator.test.js

const assert = require('assert');


const { PriceValidator, IntegerValidator } = require('../../public/js/Validator');


function testValidator(validator, input, expectedOutput, expectedInputValue, message) {
    validator.setInputElement(input);
    validator.validate();
    assert.strictEqual(validator.getValue(), expectedOutput, message);
    assert.strictEqual(input.value, expectedInputValue, `Test failed: input value should be "${expectedInputValue}".`);
}

const priceValidator = new PriceValidator();

testValidator(
    priceValidator,
    { value: '01', defaultValue: '01' },
    100,
    '1',
    'Test failed: "01" should be converted to 100 (cents).'
);

testValidator(
    priceValidator,
    { value: '01.', defaultValue: '01' },
    100,
    '1.',
    'Test failed: "01" should be converted to 100 (cents).'
);

testValidator(
    priceValidator,
    { value: '2,', defaultValue: '2' },
    200,
    '2.',
    'Test failed: "2," should be converted to 200 (cents).'
);


testValidator(
    priceValidator,
    { value: '2,0', defaultValue: '2,0' },
    200,
    '2.0',
    'Test failed: "2,0" should be converted to 200 (cents).'
);

testValidator(
    priceValidator,
    { value: '1.00', defaultValue: '1.00' },
    100,
    '1.00',
    'Test failed: "1.00" should be converted to 100 (cents).'
);

testValidator(
    priceValidator,
    { value: '1.00.', defaultValue: '1.00' },
    100,
    '1.00',
    'Test failed: "1.00" should be converted to 100 (cents).'
);


testValidator(
    priceValidator,
    { value: 's', defaultValue: '1.00' },
    100,
    '1.00',
    'Test failed: "1.00" should be converted to 100 (cents).'
);

testValidator(
    priceValidator,
    { value: '2,', defaultValue: '2' },
    200,
    '2.',
    'Test failed: "2," should be converted to 200 (cents).'
);

testValidator(
    priceValidator,
    { value: '', defaultValue: '2' },
    0,
    '',
    'Test failed: "" should be converted to 0 (cents).'
);

testValidator(
    priceValidator,
    { value: '-', defaultValue: '2' },
    200,
    '2',
    'Test failed: "" should be converted to 0 (cents).'
);


testValidator(
    priceValidator,
    { value: '01', defaultValue: '1' },
    100,
    '1',
    'Test failed: "01" should be converted to 100 (cents).'
);

testValidator(
    priceValidator,
    { value: '0', defaultValue: '' },
    0,
    '0',
    'Test failed: "0" should be converted to 0.'
);

const integerValidator = new IntegerValidator();

testValidator(
    integerValidator,
    { value: '', defaultValue: '-1000' },
    0,
    '',
    'Test failed: "-1000" should be valid.'
);

testValidator(
    integerValidator,
    { value: '-', defaultValue: '-1' },
    0,
    '-',
    'Test failed: "-1000" should be valid.'
);


testValidator(
    integerValidator,
    { value: '01', defaultValue: '1' },
    1,
    '1',
    'Test failed: "-1000" should be valid.'
);

testValidator(
    integerValidator,
    { value: '-1000', defaultValue: '-1000' },
    -1000,
    '-1000',
    'Test failed: "-1000" should be valid.'
);

testValidator(
    integerValidator,
    { value: '123.', defaultValue: '123' },
    123,
    '123',
    'Test failed: "123" should be invalid.'
);

testValidator(
    integerValidator,
    { value: '123.0', defaultValue: '123' },
    123,
    '123',
    'Test failed: "123" should be invalid.'
);

testValidator(
    integerValidator,
    { value: '123,0', defaultValue: '123' },
    123,
    '123',
    'Test failed: "123" should be invalid.'
);

testValidator(
    integerValidator,
    { value: '-100001', defaultValue: '-10000' },
    -10000,
    '-10000',
    'Test failed: "-100001" should be invalid.'
);

testValidator(
    integerValidator,
    { value: '100001', defaultValue: '10000' },
    10000,
    '10000',
    'Test failed: "100001" should be invalid.'
);


console.log('All tests passed!');
