<?php

enum Month: int
{
    case January = 1;
}

function requireInteger(int $value): void {}

function normalizeMonth(Month|int $month): void
{
    if ($month instanceof Month) {
        $month = $month->value;
    }

    requireInteger($month);
}
