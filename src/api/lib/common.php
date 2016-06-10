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


function uptime_human_readable($seconds){
	$days = floor($seconds / 86400);
	$hours = floor($seconds % 86400 / 3600);
	$mins = floor($seconds % 86400 % 3600 / 60);

	if ($seconds == 0){
		return 0;
	} elseif ($days < 1) {
		return "{$hours}h, {$mins}m";
	} else {
		return "{$days} days";
	}
}



?>