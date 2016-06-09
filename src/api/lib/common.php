<?php

// common functions


$units = explode(' ', 'B KB MB GB TB PB');
function format_size($size) {
    global $units;
    $mod = 1024;
    for ($i = 0; $size > $mod; $i++) {
        $size /= $mod;
    }
    $endIndex = strpos($size, ".")+3;
    return substr( $size, 0, $endIndex).' '.$units[$i];
}



?>