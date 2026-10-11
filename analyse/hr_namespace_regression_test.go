package analyse

import "testing"

func TestPHPDocArrayShapeAliasesUseTheirOwningNamespace(t *testing.T) {
	issues := analysePHP(t, `<?php
namespace First\Configurator;
final class Placeholder {}

namespace Second\Configurator;
/**
 * @psalm-type Route = array{path: string}
 * @psalm-type Routes = array{"when@dev"?: array<string, Route>, ...<string, Route>}
 */
final class RoutesReference {
    /** @param Routes $config @psalm-return Routes */
    public static function config(array $config): array { return $config; }
}
`)
	for _, issue := range issues {
		if issue.Code == "A.RETURN.TYPE" || issue.Code == level2PHPDocParamTypeCode || issue.Code == level2PHPDocReturnTypeCode {
			t.Fatalf("PHPDoc aliases in a later namespace must keep that namespace's context, got %#v", issues)
		}
	}
}
