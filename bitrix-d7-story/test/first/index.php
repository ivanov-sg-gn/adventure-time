<?
require($_SERVER["DOCUMENT_ROOT"]."/bitrix/header.php");
$APPLICATION->SetTitle("component first");
?>

<?$APPLICATION->IncludeComponent(
    "my:component.first",
    "",
    Array(
        "CACHE_TIME" => "36000000",
        "CACHE_TYPE" => "A",
        "IBLOCK_ID" => "5",
        "IBLOCK_TYPE" => "new_type",
        "PREFIX" => "{$_GET['name']}"
    )
);?>

<?require($_SERVER["DOCUMENT_ROOT"]."/bitrix/footer.php");?>