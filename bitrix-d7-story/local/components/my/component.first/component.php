<?
if (!defined("B_PROLOG_INCLUDED") || B_PROLOG_INCLUDED !== true) die();

CPageOption::SetOptionString("main", "nav_page_in_session", "N");

use Bitrix\Main\Loader,
    Bitrix\Iblock;


if (!isset($arParams["CACHE_TIME"])) $arParams["CACHE_TIME"] = 36000000;
$arParams["IBLOCK_TYPE"] = trim($arParams["IBLOCK_TYPE"]);
$arParams["IBLOCK_ID"] = trim($arParams["IBLOCK_ID"]);
$arParams["PREFIX"] = htmlspecialchars(trim($arParams["PREFIX"]));
$arParams['COUNT'] = intval($arParams['COUNT']) > 0 ? intval($arParams['COUNT']) : 10;

$navParams = [
    'nPageSize' => $arParams['COUNT']
];

$navigation = \CDBResult::GetNavParams($this->navParams);

if ($this->startResultCache(false, [$navigation, $USER->GetUserGroupArray()])) {
    if (!Loader::includeModule("iblock")) {
        $this->abortResultCache();
        ShowError(GetMessage("IBLOCK_MODULE_NOT_INSTALLED"));
        return;
    }

    $arFilter = [
        'ACTIVE' => 'Y',
        'ACTIVE_DATE' => 'Y',
    ];


    if (!empty($arParams["IBLOCK_TYPE"])) {
        $arFilter['IBLOCK_TYPE'] = $arParams["IBLOCK_TYPE"];
    }
    if (!empty($arParams["IBLOCK_ID"])) {
        $arFilter['IBLOCK_ID'] = $arParams["IBLOCK_ID"];
    }


    $arResult["ITEMS"] = [];
    $arResult["ELEMENTS"] = [];

    $rsElement = CIBlockElement::GetList(
        [
            'IBLOCK_SECTION_ID' => 'asc',
            'sort' => 'asc',
            'name' => 'asc',
        ],
        $arFilter,
        false,
        $navParams,
        ['IBLOCK_ID', 'ID', 'NAME']
    );

    while ($arItem = $rsElement->fetch()) {
        # buttons
        $arButtons = CIBlock::GetPanelButtons(
            $arItem["IBLOCK_ID"],
            $arItem["ID"],
            $arItem["SECTION_ID"] ?? 0
        );
        $arItem["EDIT_LINK"] = $arButtons["edit"]["edit_element"]["ACTION_URL"];
        $arItem["DELETE_LINK"] = $arButtons["edit"]["delete_element"]["ACTION_URL"];

        # add prefix
        if (!empty($arParams["PREFIX"])) {
            $arItem['NAME'] = $arParams["PREFIX"] . ' | ' . $arItem['NAME'];
        }

        $arResult["ITEMS"][$arItem["ID"]] = $arItem;
        $arResult["ELEMENTS"][] = $arItem["ID"];
    }


    # get sections
    if (!empty($arResult["ITEMS"])) {

        $rsSection = CIBlockElement::GetElementGroups($arResult["ELEMENTS"], true, ['IBLOCK_ELEMENT_ID', 'ID', 'NAME']);

        while ($arItem = $rsSection->fetch()) {
            $arResult["ITEMS"][$arItem["IBLOCK_ELEMENT_ID"]]['SECTIONS'][$arItem['ID']] = [
                'ID' => $arItem['ID'],
                'NAME' => $arItem['NAME']
            ];
        }

    }


    $arResult["NAV_STR"] = $rsElement->GetPageNavString(GetMessage("NAV_MESSAGE"));


    if (empty($arResult["ELEMENTS"])) {
        $this->abortResultCache();
        Iblock\Component\Tools::process404(
            GetMessage("404_MESSAGE")
            , true
            , true
            , true
        );
        return;
    }


    $this->setResultCacheKeys(["ELEMENTS"]);
    $this->includeComponentTemplate();
}