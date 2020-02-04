<?php
// если не нужно собирать статистику, выполнять агенты
//define("NO_KEEP_STATISTIC", true);
//define("NO_AGENT_CHECK", true);

// убираем мусор типа отладки
define('PUBLIC_AJAX_MODE', true);
require($_SERVER["DOCUMENT_ROOT"] . "/bitrix/modules/main/include/prolog_before.php");


header('Content-Type: text/csv; charset=utf-8');
header('Content-Disposition: attachment; filename=users.csv');


CBitrixComponent::includeComponentClass("my:component.two");


$twoComponent = new twoComponent();
$twoComponent->arParams = $twoComponent->onPrepareComponentParams(['SHOW_ALL' => 'Y']);
$twoComponent->getData();


if (!empty($twoComponent->arResult['USERS'])) {
    $out = fopen('php://output', 'w');

    fputcsv($out, ['id', 'name'], ";");

    foreach ($twoComponent->arResult['USERS'] as $item) {

        fputcsv($out, [$item['ID'], $item['NAME']], ';');
    }

    fclose($out);
}


require($_SERVER["DOCUMENT_ROOT"] . "/bitrix/modules/main/include/epilog_after.php");
?>
