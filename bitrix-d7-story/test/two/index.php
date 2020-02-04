<?
require($_SERVER["DOCUMENT_ROOT"]."/bitrix/header.php");
$APPLICATION->SetTitle("component two");
?>

<?$APPLICATION->IncludeComponent(
    "my:component.two",
    "",
    Array(
        "AJAX_MODE" => "Y",
        "AJAX_OPTION_ADDITIONAL" => "",
        "AJAX_OPTION_HISTORY" => "N",
        "AJAX_OPTION_JUMP" => "N",
        "AJAX_OPTION_STYLE" => "N",
        "CACHE_TIME" => "36000000",
        "CACHE_TYPE" => "A",
        "COUNT" => "10",
        'SHOW_ALWAYS' => 'Y'
    )
);?> 

<?require($_SERVER["DOCUMENT_ROOT"]."/bitrix/footer.php");?>