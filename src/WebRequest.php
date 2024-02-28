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
        if (!extension_loaded('phpgo')) {
            throw new \Exception("phpgo extension not loaded");
        }

        if (!function_exists('phpgo_load')) {
            throw new \Exception("phpgo_load function not found");
        }

        $this->module = \phpgo_load($module, $name);
        if (!$this->module) {
            throw new \Exception("Module [$name::$module] not found");
        }

        $this->method = $_SERVER['REQUEST_METHOD'];
        $this->url = $_SERVER['REQUEST_URI'];
        $this->schema = $protocol = (!empty($_SERVER['HTTPS']) && $_SERVER['HTTPS'] !== 'off' || $_SERVER['SERVER_PORT'] == 443) ? "https://" : "http://";
        $this->headers = getallheaders();
        $this->headers['REMOTE_ADDR'] = $_SERVER['REMOTE_ADDR'];
        $this->body = file_get_contents('php://input');
        $this->form = $_POST;
        $this->files = $_FILES;

        if (count($this->files) === 0) {
            $this->files = null;
        }

        if (count($this->form) === 0) {
            $this->form = null;
        }

        if (count($this->headers) === 0) {
            $this->headers = null;
        }
    }

    public function serialize() {
        return json_encode($this, JSON_UNESCAPED_SLASHES | JSON_UNESCAPED_UNICODE);
    }

    /**
     * Run the module with the serialized request
     * @return \Schivei\PhpGo\WebResponse
     */
    public function run() {
        $jsonResponse = $this->module->run($this->serialize());

        return \Schivei\PhpGo\WebResponse::deserialize($jsonResponse);
    }
}
