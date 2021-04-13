#!/usr/bin/env php
<?php

$finder = PhpCsFixer\Finder::create()->in("src");

$config = new PhpCsFixer\Config();
return $config->setRules([
    '@PSR2' => true,
    'concat_space' => ['spacing' => 'one'],
    'ordered_class_elements' => ['sortAlgorithm' => 'alpha', 'order' => ['use_trait', 'constant_public', 'constant_protected', 'constant_private', 'property_public', 'property_protected', 'property_private', 'construct', 'destruct', 'magic', 'phpunit', 'method_public', 'method_protected', 'method_private']],
    'php_unit_internal_class' => false,
    'php_unit_test_class_requires_covers' => false,
    'phpdoc_line_span' => ['const' => 'single', 'method' => 'multi', 'property' => 'single'],
    'return_assignment' => false,
    'yoda_style' => false,
])->setFinder($finder);
