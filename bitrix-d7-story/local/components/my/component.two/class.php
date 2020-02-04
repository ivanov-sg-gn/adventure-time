<?
if (!defined("B_PROLOG_INCLUDED") || B_PROLOG_INCLUDED !== true) die();

use \Bitrix\Main\UserTable,
    \Bitrix\Main\Context,
    \Bitrix\Main\Localization\Loc,
    \Bitrix\Main\SystemException,
    \Bitrix\Main\Loader,
    \Bitrix\Main\UI,
    \Bitrix\Main\Page\Asset;

Loc::loadMessages(__FILE__);


class twoComponent extends CBitrixComponent
{
    public function onPrepareComponentParams($arParams)
    {
        $arParams["CACHE_TIMES"] = !isset($arParams["CACHE_TIMES"]) ? 3600 : $arParams["CACHE_TIMES"];
        $arParams['AJAX'] = $_REQUEST['AJAX'] == 'Y' ? true : false;
        $arParams['COUNT'] = intval($arParams['COUNT']) > 0 ? intval($arParams['COUNT']) : 10;
        $arParams['SHOW_ALL'] = $arParams['SHOW_ALL'] == 'Y' ? true : false;
        $arParams['SHOW_ALWAYS'] = $arParams['SHOW_ALWAYS'] == 'Y' ? true : false;

        return $arParams;
    }

    public function setDefaultParams()
    {
        CPageOption::SetOptionString('main', 'nav_page_in_session', 'N');

        $this->arResult['OB_REQUEST'] = Context::getCurrent()->getRequest();
    }

    public function executeComponent()
    {
        GLOBAL $USER;

        $this->setDefaultParams();

        $this->getData();

        $this->arResult['NAV_XML'] = $this->getNavHtml();

        $this->includeComponentTemplate();

        if (empty($this->arResult["USERS"])) {
            UserTable::getEntity()->cleanCache();

            throw new SystemException(Loc::getMessage('EMPTY_USER'));
        }
    }

    public function getData()
    {
        $this->arResult["USERS"] = [];
        $this->arResult["USERS_ID"] = [];
        
        $params = [
            'select' => ['ID', 'NAME'],
            'order' => ['NAME' => 'ASC'],
        ];

        if ($this->arParams['SHOW_ALL'] == false) {
            $nav = new UI\PageNavigation('page');
            $nav->allowAllRecords(true)
                ->setPageSize($this->arParams['COUNT'])
                ->initFromUri();

            $params['offset'] = $nav->getOffset();
            $params['limit'] = $nav->getLimit();
            $params['count_total'] = true;
            $params['cache'] = ['ttl' => $this->arParams["CACHE_TIMES"]];

            $this->arResult['NAV_OBJECT'] = $nav;
        }

        $rsUsers = UserTable::getList($params);


        if ($this->arParams['SHOW_ALL'] == false) {
            $this->arResult['NAV_OBJECT']->setRecordCount($rsUsers->getCount());
        }

        //->fetchAll
        while ($arItems = $rsUsers->Fetch()) {
            $this->arResult["USERS"][$arItems['ID']] = $arItems;
            $this->arResult["USERS_ID"][] = $arItems['ID'];
        }
    }

    public function getNavHtml()
    {
        GLOBAL $APPLICATION;

        if ($this->arParams['SHOW_ALL'] == false) {
            ob_start();
            $APPLICATION->IncludeComponent(
                "bitrix:main.pagenavigation",
                "",
                array(
                    "NAV_OBJECT" => $this->arResult['NAV_OBJECT'],
                    'SHOW_ALWAYS' => $this->arParams['SHOW_ALWAYS'],
                ),
                false
            );
            $res = ob_get_clean();

            return $res;
        }

    }

}