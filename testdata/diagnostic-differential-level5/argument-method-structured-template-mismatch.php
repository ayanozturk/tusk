<?php
class Selector {
    /**
     * @template Value of int
     * @param array<Value> $values
     */
    public function accept(array $values): void {}
}
class Collector {
    public array $values;
    /**
     * @template Value of int
     * @param array<Value> $values
     */
    public function __construct(array $values) { $this->values = $values; }
}
/** @param array<string> $values */
function wrong(Selector $selector, array $values): Collector {
    $selector->accept(values: $values);
    return new Collector($values);
}
