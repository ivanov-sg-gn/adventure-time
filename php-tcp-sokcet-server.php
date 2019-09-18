<?php
$thisFile = basename(__FILE__);

# Check starting this script
exec("/bin/ps awx | grep '$thisFile'", $arExec);
foreach ($arExec as $str) {
    if (strpos($str, '/bin/sh -c') !== false)
        continue;
    if (strpos($str, "/$thisFile") !== false)
        exit();
}

# If start from console, DOCUMENT_ROOT not exist
if (empty($_SERVER['DOCUMENT_ROOT']))
    $_SERVER["DOCUMENT_ROOT"] = dirname(__DIR__);


# Launch socket with use ssl
$stream_context = stream_context_create();
stream_context_set_option($stream_context, 'ssl', 'local_cert', '/var/www/ssl/certificate.crt');
stream_context_set_option($stream_context, 'ssl', 'local_pk', '/var/www/ssl/certificate.key');
stream_context_set_option($stream_context, 'ssl', 'allow_self_signed', true);
stream_context_set_option($stream_context, 'ssl', 'verify_peer', false);

$socketServer = stream_socket_server('ssl://0.0.0.0:33377', $errno, $errstr, STREAM_SERVER_BIND | STREAM_SERVER_LISTEN, $stream_context);

# Check status
if (!$socketServer) {
    print_r([
        'errno' => $errno,
        'errstr' => $errstr
    ]);
    die("\nsocket failure\n");
} else {
    echo "socket started\n";
}

# Array connections ( yes, there may be several (= )
$connections = [];

while (true) {
    $read = $connections;

    # Add new connections
    $read[] = $socketServer;


    $write = $except = null;
    # Waiting for flow status (without timeout)
    if (!stream_select($read, $write, $except, null))
        break;

    # Handshake with Web
    if (in_array($socketServer, $read)) {
        # Accept new connections
        if (($connection = stream_socket_accept($socketServer, -1)) && ($info = mysocket::handshake($connection))) {
            # Add to the list of processed
            $connections[] = $connection;

            mysocket::send($socketServer, 'Hello man =)');
        }

        unset($read[array_search($socketServer, $read)]);
    }

    # Process each connection ( because php is not asynchronous (= )
    foreach ($read as $connection) {
        $data = fread($connection, 100000);

        # Close connection
        if ($data === flse) {
            fclose($connection);
            continue;
        }

        # Get messages
        $data = mysocket::decode($data);
        if (empty($data['payload'])) {
            continue;
        }
        $data = json_decode($data['payload'], 1);

        # Check keep alive status from web socket
        if ($data["method"] != "keep_alive") continue;


        if (is_array($data)) {
            // code ...
        }
    }
}


class mysocket
{
    static function send($connection, string $messages = ''): array
    {
        if (!is_resource($connection)) return false;

        $messages = mysocket::encode(json_encode($messages));

        return fwrite($connection, $messages) === false;
    }

    static function encode($payload, $type = 'text', $masked = false)
    {
        $frameHead = array();
        $payloadLength = strlen($payload);

        switch ($type) {
            case 'text':
                # first byte indicates FIN, Text-Frame (10000001):
                $frameHead[0] = 129;
                break;
            case 'close':
                # first byte indicates FIN, Close Frame(10001000):
                $frameHead[0] = 136;
                break;
            case 'ping':
                # first byte indicates FIN, Ping frame (10001001):
                $frameHead[0] = 137;
                break;
            case 'pong':
                # first byte indicates FIN, Pong frame (10001010):
                $frameHead[0] = 138;
                break;
        }

        # set mask and payload length (using 1, 3 or 9 bytes)
        if ($payloadLength > 65535) {
            $payloadLengthBin = str_split(sprintf('%064b', $payloadLength), 8);
            $frameHead[1] = ($masked === true) ? 255 : 127;

            for ($i = 0; $i < 8; $i++) {
                $frameHead[$i + 2] = bindec($payloadLengthBin[$i]);
            }

            # most significant bit MUST be 0
            if ($frameHead[2] > 127) {
                return array('type' => '', 'payload' => '', 'error' => 'frame too large (1004)');
            }
        } elseif ($payloadLength > 125) {
            $payloadLengthBin = str_split(sprintf('%016b', $payloadLength), 8);

            $frameHead[1] = ($masked === true) ? 254 : 126;
            $frameHead[2] = bindec($payloadLengthBin[0]);
            $frameHead[3] = bindec($payloadLengthBin[1]);
        } else {
            $frameHead[1] = ($masked === true) ? $payloadLength + 128 : $payloadLength;
        }

        # convert frame-head to string:
        foreach (array_keys($frameHead) as $i) {
            $frameHead[$i] = chr($frameHead[$i]);
        }

        if ($masked === true) {
            # generate a random mask:
            $mask = array();

            for ($i = 0; $i < 4; $i++) {
                $mask[$i] = chr(rand(0, 255));
            }

            $frameHead = array_merge($frameHead, $mask);
        }

        $frame = implode('', $frameHead);

        # append payload to frame:
        for ($i = 0; $i < $payloadLength; $i++) {
            $frame .= ($masked === true) ? substr($payload, $i, 1) ^ $mask[$i % 4] : substr($payload, $i, 1);
        }

        return $frame;
    }

    static function decode($data)
    {
        $unmaskedPayload = '';
        $decodedData = array();

        # estimate frame type:
        $firstByteBinary = sprintf('%08b', ord($data[0]));
        $secondByteBinary = sprintf('%08b', ord($data[1]));
        $opcode = bindec(substr($firstByteBinary, 4, 4));
        $isMasked = ($secondByteBinary[0] == '1') ? true : false;
        $payloadLength = ord($data[1]) & 127;

        # unmasked frame is received:
        if (!$isMasked) {
            return array('type' => '', 'payload' => '', 'error' => 'protocol error (1002)');
        }

        switch ($opcode) {
            # text frame:
            case 1:
                $decodedData['type'] = 'text';
                break;

            case 2:
                $decodedData['type'] = 'binary';
                break;

            # connection close frame:
            case 8:
                $decodedData['type'] = 'close';
                break;

            # ping frame:
            case 9:
                $decodedData['type'] = 'ping';
                break;

            # pong frame:
            case 10:
                $decodedData['type'] = 'pong';
                break;

            default:
                return array('type' => '', 'payload' => '', 'error' => 'unknown opcode (1003)');
        }

        if ($payloadLength === 126) {
            $mask = substr($data, 4, 4);
            $payloadOffset = 8;
            $dataLength = bindec(sprintf('%08b', ord($data[2])) . sprintf('%08b', ord($data[3]))) + $payloadOffset;
        } elseif ($payloadLength === 127) {
            $mask = substr($data, 10, 4);
            $payloadOffset = 14;
            $tmp = '';
            for ($i = 0; $i < 8; $i++) {
                $tmp .= sprintf('%08b', ord($data[$i + 2]));
            }
            $dataLength = bindec($tmp) + $payloadOffset;
            unset($tmp);
        } else {
            $mask = substr($data, 2, 4);
            $payloadOffset = 6;
            $dataLength = $payloadLength + $payloadOffset;
        }

        /**
         * We have to check for large frames here. socket_recv cuts at 1024 bytes
         * so if websocket-frame is > 1024 bytes we have to wait until whole
         * data is transferd.
         */
        if (strlen($data) < $dataLength) {
            return false;
        }

        if ($isMasked) {
            for ($i = $payloadOffset; $i < $dataLength; $i++) {
                $j = $i - $payloadOffset;
                if (isset($data[$i])) {
                    $unmaskedPayload .= $data[$i] ^ $mask[$j % 4];
                }
            }
            $decodedData['payload'] = $unmaskedPayload;
        } else {
            $payloadOffset = $payloadOffset - 4;
            $decodedData['payload'] = substr($data, $payloadOffset);
        }

        return $decodedData;
    }

    static function handshake($connections)
    {
        $info = array();

        $line = fgets($connections);
        $header = explode(' ', $line);
        $info['method'] = $header[0];
        $info['uri'] = $header[1];

        # Get header fom connection
        while ($line = rtrim(fgets($connections))) {
            if (preg_match('/\A(\S+): (.*)\z/', $line, $matches)) {
                $info[$matches[1]] = $matches[2];
            } else {
                break;
            }
        }

        $address = explode(':', stream_socket_get_name($connections, true));
        $info['ip'] = $address[0];
        $info['port'] = $address[1];

        if (empty($info['Sec-WebSocket-Key'])) {
            return false;
        }

        # Send the header according to the websocket Protocol
        $SecWebSocketAccept = base64_encode(pack('H*', sha1($info['Sec-WebSocket-Key'] . '258EAFA5-E914-47DA-95CA-C5AB0DC85B11')));
        $upgrade = "HTTP/1.1 101 Web Socket Protocol Handshake\r\n" . "Upgrade: websocket\r\n" . "Connection: Upgrade\r\n" . "Sec-WebSocket-Accept:$SecWebSocketAccept\r\n\r\n";
        fwrite($connections, $upgrade);

        return $info;
    }
}
