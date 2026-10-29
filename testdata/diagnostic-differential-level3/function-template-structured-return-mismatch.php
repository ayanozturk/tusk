<?php
/**
 * @template Value
 * @param list<Value> $values
 * @return Value
 */
function firstValue(array $values) { return $values[0]; }

/** @param list<string> $values */
function wrong(array $values): int { return firstValue($values); }
