#!/usr/bin/env php
<?php

$finder = PhpCsFixer\Finder::create()->in("src");

$config = new PhpCsFixer\Config();
return $config->setRules([
    '@PSR2' => true,
    'array_syntax' => ['syntax' => 'short'],
    'binary_operator_spaces' => ['align_double_arrow' => true],
    'blank_line_before_statement' => ['statements' => ['return']],
    'global_namespace_import' => true,
    'native_function_casing' => true,
    'no_extra_blank_lines' => true,
    'no_unused_imports' => true,
    'ordered_class_elements' => ['sortAlgorithm' => 'alpha', 'order' => ['use_trait', 'constant_public', 'constant_protected', 'constant_private', 'property_public', 'property_protected', 'property_private', 'construct', 'destruct', 'magic', 'phpunit', 'method_public', 'method_protected', 'method_private']],
    'ordered_imports' => true,
    'phpdoc_add_missing_param_annotation' => ['only_untyped' => true],
    'phpdoc_indent' => true,
    'phpdoc_line_span' => ['const' => 'single', 'method' => 'multi', 'property' => 'single'],
    'phpdoc_no_package' => true,
    'phpdoc_order' => true,
    'phpdoc_scalar' => true,
    'phpdoc_separation' => true,
    'phpdoc_trim_consecutive_blank_line_separation' => true,
    'phpdoc_trim' => true,
    'single_quote' => true,
    'trailing_comma_in_multiline_array' => true,
])->setFinder($finder);
