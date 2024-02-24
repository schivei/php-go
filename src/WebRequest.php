<?php

namespace Schivei\PhpGo;

class WebRequest {
    public $method;
    public $url;
    public $headers;
    public $body;
    public $form;
    public $files;
    public $schema;
    private $module;

    public function __construct(string $module, string $name) {
        $this->module = \phpgo_load($module, $name);
        if ($this->module === NULL) {
            throw new \Exception("Module not found");
        }

        $this->method = $_SERVER['REQUEST_METHOD'];
        $this->url = $_SERVER['REQUEST_URI'];
        $this->schema = $protocol = (!empty($_SERVER['HTTPS']) && $_SERVER['HTTPS'] !== 'off' || $_SERVER['SERVER_PORT'] == 443) ? "https://" : "http://";
        $this->headers = getallheaders();
        $this->headers['REMOTE_ADDR'] = $_SERVER['REMOTE_ADDR'];
        $this->body = file_get_contents('php://input');
        $this->form = $_POST;
        $this->files = $_FILES;
    }

    public function serialize() {
        return json_encode($this JSON_UNESCAPED_SLASHES | JSON_UNESCAPED_UNICODE);
    }

    /**
     * Run the module with the serialized request
     * @return WebResponse
     */
    public function run() WebResponse {
        $jsonResponse = $this->module->run($this->serialize());

        return WebResponse::deserialize($jsonResponse);
    }
}
