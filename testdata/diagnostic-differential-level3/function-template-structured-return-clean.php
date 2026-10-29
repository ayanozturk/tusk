<?php
/**
 * @template Value
 * @param list<Value> $values
 * @return Value
 */
function firstValue(array $values) { return $values[0]; }

/** @param list<string> $values */
function clean(array $values): string { return firstValue($values); }
