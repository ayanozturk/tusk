<?php
/**
 * @template Value of int
 * @param array<Value> $values
 * @return Value
 */
function boundedValue(array $values) { return $values[0]; }

/** @param array<string> $values */
function wrong(array $values): string { return boundedValue($values); }
