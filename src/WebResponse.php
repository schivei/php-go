<?php

namespace Schivei\PhpGo;

class WebResponse
{
    public $status;
    public $headers;
    public $body;

    public function __construct(string $responseJson)
    {
        $response = \json_decode($responseJson, true);

        $this->status = $response['status'];
        $this->headers = $response['headers'];
        $this->body = $response['body'];
    }

    /**
     * Write HTTP response to the client
     * @return void
     */
    public function write()
    {
        \http_response_code($this->status);

        foreach ($this->headers as $header) {
            \header($header);
        }

        if (\file_exists($this->body)) {
            \readfile($this->body);

            \unlink($this->body);
        } else {
            echo $this->body;
        }
    }

    public static function deserialize(string $responseJson): WebResponse
    {
        return new WebResponse($responseJson);
    }
}
