--TEST--
Error conditions
--FILE--
<?php
$module = phpgo_load(__DIR__ . "/fixtures/golang/test.so", "test");
$module->boolAnd(true);
$module->boolAnd(true, true, true);

$module = phpgo_load(__DIR__ . "/fixtures/golang/nop.so", "test");
$module = phpgo_load(__DIR__ . "/fixtures/golang/test.so", "nop");
--EXPECTF--
Warning: PHPGo\Module\test_%s::boolAnd() expects exactly 2 parameters, 1 given in %s on line %d

Warning: PHPGo\Module\test_%s::boolAnd() expects exactly 2 parameters, 3 given in %s on line %d

Warning: Failed loading %s/nop.so (test): %s in %s on line %d

Warning: Failed loading %s/test.so (nop): Not exports found for `nop` in %s on line %d
