<?php


$body = '[
{
	"union": {
		"input1": {"selection": {
			"input": { "relation": { "name": "staff" } },
				"attr": "age",  "selector": ">=", "arg": 31
		}},
			"input2": {"selection": {
				"input": { "relation": { "name": "staff" } },
				"attr": "name", "selector": "==", "arg": "山田"
			}}
	}
},
{
	"selection": {
		"input": { "relation": { "name": "rank" } },
		"attr": "rank",  "selector": ">=", "arg": 1
	}
}
]';

$context = stream_context_create( array('http' =>
			  array(
					'method'  => 'POST',
					'header'  => "Content-Type: application/json\r\nContent-Length: ".strlen($body)."\r\n",
					'content' => $body,
					'timeout' => 60
				   )
			 ));

$url = 'http://localhost:3050/php/dir2';
$result = file_get_contents($url, false, $context);

var_export(unserialize($result));