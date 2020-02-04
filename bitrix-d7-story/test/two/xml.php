<?php
// если не нужно собирать статистику, выполнять агенты
//define("NO_KEEP_STATISTIC", true);
//define("NO_AGENT_CHECK", true);

// убираем мусор типа отладки
define('PUBLIC_AJAX_MODE', true);
require($_SERVER["DOCUMENT_ROOT"] . "/bitrix/modules/main/include/prolog_before.php");


header('Content-Type: text/xml; charset-utf-8');
header('Content-Disposition: attachment; filename=users.xml');


CBitrixComponent::includeComponentClass("my:component.two");


$twoComponent = new twoComponent();
$twoComponent->arParams = $twoComponent->onPrepareComponentParams(['SHOW_ALL' => 'Y']);
$twoComponent->getData();

$xml = new SimpleXMLElement('<users/>');

foreach ($twoComponent->arResult['USERS'] as $item) {
    $arXMLUser = $xml->addChild('user');

    $arXMLUser->addChild('ID', $item['ID']);
    $arXMLUser->addChild('NAME', $item['NAME']);
}

echo $xml->asXML();

require($_SERVER["DOCUMENT_ROOT"] . "/bitrix/modules/main/include/epilog_after.php");
?>