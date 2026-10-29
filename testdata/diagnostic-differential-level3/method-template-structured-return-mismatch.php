<?php
class Selector {
    /**
     * @template Value
     * @param list<Value> $values
     * @return Value
     */
    public function first(array $values) { return $values[0]; }
    /**
     * @template Value
     * @param list<Value> $values
     * @return Value
     */
    public static function select(array $values) { return $values[0]; }
}
class ChildSelector extends Selector {}
interface Selection {
    /**
     * @template Value
     * @param list<Value> $values
     * @return Value
     */
    public function first(array $values);
}
/** @param list<string> $values */
function inherited(ChildSelector $selector, array $values): int { return $selector->first(values: $values); }
/** @param list<string> $values */
function fromInterface(Selection $selector, array $values): int { return $selector->first($values); }
/** @param list<string> $values */
function fromStatic(array $values): int { return Selector::select($values); }
