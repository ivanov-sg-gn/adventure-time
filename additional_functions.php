<?
/**
 * Сортировка массива по убыванию значения поля.
 * @param array  $array Массив для сортировки.
 * @param string $field Поле, по которому осуществляется сортировка.
 * @return array Отсортированный $array.
 */

if ( !function_exists( 'sortDown' ) ) {
    function sortDown ( $array, $field ) {
        usort( $array, function ( $a, $b ) use ( $field ) {
            if ( $a[ $field ] == $b[ $field ] )
                return 0;
            return ( $a[ $field ] > $b[ $field ] ) ? - 1 : 1;
        } );
        return $array;
    }
}

/**
 * Сортировка массива по возрастанию значения поля.
 * @param array  $array Массив для сортировки.
 * @param string $field Поле, по которому осуществляется сортировка.
 * @return array Отсортированный $array.
 */
if ( !function_exists( 'sortUp' ) ) {
    function sortUp ( $array, $field ) {
        usort( $array, function ( $a, $b ) use ( $field ) {
            if ( $a[ $field ] == $b[ $field ] ) {
                return 0;
            }
            return ( $a[ $field ] < $b[ $field ] ) ? - 1 : 1;
        } );
        return $array;
    }
}

/**
 * Сортировка массива по значения поля.
 * @param array  $array     Массив для сортировки.
 * @param string $direction Порядок сортировки (asc - по возрастанию,
 *                          desc - по убыванию).
 * @param string $field     Поле, по которому осуществляется сортировка.
 * @return array Отсортированный $array.
 */

if ( !function_exists( 'imSort' ) ) {
    function imSort ( $array, $direction, $field ) {
        if ( $direction == 'asc' )
            return sortUp( $array, $field );
        return sortDown( $array, $field );
    }
}


/**
 * Форматирование вывода времени
 * @param      $unixTime
 * @param bool $onlyHours
 * @return float|string
 */
if ( !function_exists( 'formatUnix' ) ) {
    function formatUnix ( $unixTime, $onlyHours = false ) {
        $hours = floor( $unixTime / 3600 );
        $mins = floor( $unixTime / 60 % 60 );
        $secs = floor( $unixTime % 60 );

        if ( $onlyHours ) {
            return $hours;
        }
        else {
            return sprintf( '%02d:%02d:%02d', $hours, $mins, $secs );
        }
    }
}


# Получить только цифры из текста
if ( !function_exists( 'getOnlyNumbers' ) ) {
    function getOnlyNumbers ( $text ) {
        return preg_replace( '/[^\d]/', '', $text );
    }
}
# Убрать только цифры из текста
if ( !function_exists( 'getOnlyText' ) ) {
    function getOnlyText ( $text ) {
        return preg_replace( '/[0-9]/', '', $text );
    }
}

# форматирование телефона в нужном формате !! Используется в импортЕ !!
if ( !function_exists( 'formatPhone' ) ) {
    function formatPhone ( $phone ) {
        return preg_replace( '/^(\+7)|8/isU', '7', getOnlyNumbers( $phone ) );
    }
}

# Формирование цены
if ( !function_exists( 'formatPrice' ) ) {
    function formatPrice ( $price ) {
        return number_format( floatval( $price ), 2, '.', ' ' );
    }
}

# узнать сколько в секундах часов/мин/сек
if ( !function_exists( 'parse_time' ) ) {
    function parse_time ( $time ) {
        $h = intval( $time / 3600 );
        $i = intval( ( $time - ( $h * 3600 ) ) / 60 );
        $s = intval( $time - ( ( $h * 3600 ) + ( $i * 60 ) ) );

        return [ $h, $i, $s ];
    }
}

// время
if ( !function_exists( 'formatPrice' ) ) {
    function microtime_float () {
        list( $usec, $sec ) = explode( " ", microtime() );
        return ( (float) $usec + (float) $sec );
    }
}

/**
 * Создание архива и запихивание в него файл
 * @param       $filename - имя файла с путём БЕЗ DOCUMENT_ROOT !!
 * @param array $path - массив путей с файлами, которые будут туда помещены
 *                    url - путь БЕЗ DOCUMENT_ROOT !!
 *                    name - имя файла (не обязательно)
 * @return array|bool
 */
function create_zip($filename, $path=[]){
    $zip = new ZipArchive();

    $error = [];

    if(file_exists($_SERVER['DOCUMENT_ROOT'].$filename)){
        return true;
    }

    $res = $zip->open($_SERVER['DOCUMENT_ROOT'].$filename, ZipArchive::CREATE);
    if($res === true){
        foreach($path as $item){
            if(empty($item['name'])){
                $item['name'] = str_replace('/', '_', $item['url']);
            }

            if(!$zip->addFile($_SERVER['DOCUMENT_ROOT'].$item['url'], $item['name'])){
                $error[] = "Файл <$item> не добавлен";
            }
        }
    }
    else{
        switch($res){
            case ZipArchive::ER_EXISTS:
                $error[] = "File already exists.";
                break;

            case ZipArchive::ER_INCONS:
                $error[] = "Zip archive inconsistent.";
                break;

            case ZipArchive::ER_MEMORY:
                $error[] = "Malloc failure.";
                break;

            case ZipArchive::ER_NOENT:
                $error[] = "No such file.";
                break;

            case ZipArchive::ER_NOZIP:
                $error[] = "Not a zip archive.";
                break;

            case ZipArchive::ER_OPEN:
                $error[] = "Can't open file.";
                break;

            case ZipArchive::ER_READ:
                $error[] = "Read error.";
                break;

            case ZipArchive::ER_SEEK:
                $error[] = "Seek error.";
                break;

            default:
                $error[] = "Unknow (Code $arOpen)";
                break;
        }
    }

    @$zip->close();

    return empty($error) ? true : $error;
}
