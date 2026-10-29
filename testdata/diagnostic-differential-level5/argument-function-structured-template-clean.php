<?php
/**
 * @template Value of int
 * @param array<Value> $values
 * @return Value
 */
function boundedValue(array $values) { return $values[0]; }

/** @param array<int> $values */
function clean(array $values): int { return boundedValue($values); }
