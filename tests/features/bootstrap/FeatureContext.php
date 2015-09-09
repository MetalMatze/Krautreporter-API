<?php

use Behat\Behat\Context\Context;
use Behat\Behat\Context\SnippetAcceptingContext;
use Behat\Behat\Tester\Exception\PendingException;
use Behat\Gherkin\Node\PyStringNode;
use Behat\Gherkin\Node\TableNode;
use Behat\MinkExtension\Context\MinkContext;
use GuzzleHttp\Client;
use GuzzleHttp\Exception\BadResponseException;
use GuzzleHttp\Message\Response;

/**
 * Defines application features from the specific context.
 */
class FeatureContext extends MinkContext implements Context, SnippetAcceptingContext
{
    /**
     * @var Client
     */
    protected $client;

    /**
     * @var string
     */
    private $resource;

    /**
     * @var
     */
    private $response;

    /**
     * @var
     */
    private $responsePayload;

    /**
     * @var
     */
    private $requestPayload;

    /**
     * @var
     */
    private $scope;

    /**
     * Initializes context.
     *
     * Every scenario gets its own context instance.
     * You can also pass arbitrary arguments to the
     * context constructor through behat.yml.
     */
    public function __construct($url)
    {
        $config['base_url'] = $url;
        $this->client = new Client($config);
    }

    /**
     * @When /^I request "(GET|PUT|POST|DELETE)" "([^"]*)"$/
     */
    public function iRequest($method, $resource)
    {
        $method = strtolower($method);
        $options = [];

        $this->resource = $resource;

        if($this->requestPayload != null)
        {
            $options = array_merge($options, ['body' => $this->requestPayload]);
        }

        try
        {
            $this->response = $this->client->$method($resource, $options);
        }
        catch(BadResponseException $e)
        {
            $response = $e->getResponse();

            if($response == null)
            {
                throw $e;
            }

            $this->response = $response;
        }
    }

    /**
     * @Then /^I get a "(\d+)" response$/
     */
    public function iGetAResponse($statusCode)
    {
        $contentType = $this->response->getHeader('Content-Type');
        if($contentType === 'application/json')
        {
            $body = $this->response->getBody();
        }
        else
        {
            $body = 'Output is '.$contentType.', which is not JSON and is therefore scary. Run the request manually.';
        }

        PHPUnit_Framework_Assert::assertSame((int) $statusCode, (int) $this->response->getStatusCode(), $body);
    }

    /**
     * @Then /^I get a "(\d+)" error response$/
     */
    public function iGetAErrorResponse($status_code, $message = null)
    {
        switch($status_code)
        {
            case 400:
                $this->iGetAResponse(400);
                $this->thePropertyIsAIntegerEqualling('status_code', 400);
                $this->thePropertyEquals('message', is_null($message) ? 'Bad Request' : $message);
                break;
            case 401:
                $this->iGetAResponse(401);
                $this->thePropertyIsAIntegerEqualling('status_code', 401);
                $this->thePropertyEquals('message', is_null($message) ? 'Failed to authenticate because of bad credentials or an invalid authorization header.' : $message);
                break;
            case 403:
                $this->iGetAResponse(403);
                $this->thePropertyIsAIntegerEqualling('status_code', 403);
                $this->thePropertyEquals('message', is_null($message) ? 'Forbidden' : $message);
                break;
            case 404:
                $this->iGetAResponse(404);
                $this->thePropertyIsAIntegerEqualling('status_code', 404);
                $this->thePropertyEquals('message', is_null($message) ? '404 Not Found' : $message);
                break;
            case 422:
                $this->iGetAResponse(422);
                $this->scopeIntoTheProperty('error');
                $this->thePropertyIsAStringEqualling('type', 'ValidationException');
                $this->thePropertyIsAString('message');
                $this->thePropertyIsAnObject('errors');
                break;
            default:
                throw new PendingException;
        }
    }

    /**
     * @Then /^I get a "(\d+)" error response with message "([^"]*)"$/
     */
    public function iGetAErrorResponseWithMessage($status_code, $message)
    {
        $this->iGetAErrorResponse($status_code, $message);
    }

    /**
     * @Given /^the property "([^"]*)" equals "([^"]*)"$/
     */
    public function thePropertyEquals($property, $expectedValue)
    {
        $payload = $this->getScopePayload();
        $actualValue = $this->arrayGet($payload, $property);
        PHPUnit_Framework_Assert::assertEquals(
            $actualValue,
            $expectedValue,
            "Asserting the [$property] property in current scope equals [$expectedValue]: ".json_encode($payload)
        );
    }

    /**
     * @Given /^the property "([^"]*)" exists$/
     */
    public function thePropertyExists($property)
    {
        $payload = $this->getScopePayload();
        $message = sprintf(
            'Asserting the [%s] property exists in the scope [%s]: %s',
            $property,
            $this->scope,
            json_encode($payload)
        );
        if (is_object($payload)) {
            PHPUnit_Framework_Assert::assertTrue(array_key_exists($property, get_object_vars($payload)), $message);
        } else {
            PHPUnit_Framework_Assert::assertTrue(array_key_exists($property, $payload), $message);
        }
    }
    /**
     * @Given /^the property "([^"]*)" is absent$/
     */
    public function thePropertyIsAbsent($property)
    {
        $payload = $this->getScopePayload();
        $message = sprintf(
            'Asserting the [%s] property is absent in the scope [%s]: %s',
            $property,
            $this->scope,
            json_encode($payload)
        );
        if (is_object($payload)) {
            PHPUnit_Framework_Assert::assertFalse(array_key_exists($property, get_object_vars($payload)), $message);
        } else {
            PHPUnit_Framework_Assert::assertFalse(array_key_exists($property, $payload), $message);
        }
    }
    /**
     * @Given /^the property "([^"]*)" is an array$/
     */
    public function thePropertyIsAnArray($property)
    {
        $payload = $this->getScopePayload();
        $actualValue = $this->arrayGet($payload, $property);
        PHPUnit_Framework_Assert::assertTrue(
            is_array($actualValue),
            "Asserting the [$property] property in current scope [{$this->scope}] is an array: ".json_encode($payload)
        );
    }
    /**
     * @Given /^the property "([^"]*)" is an object$/
     */
    public function thePropertyIsAnObject($property)
    {
        $payload = $this->getScopePayload();
        $actualValue = $this->arrayGet($payload, $property);
        PHPUnit_Framework_Assert::assertTrue(
            is_object($actualValue),
            "Asserting the [$property] property in current scope [{$this->scope}] is an object: ".json_encode($payload)
        );
    }
    /**
     * @Given /^the property "([^"]*)" is an empty array$/
     */
    public function thePropertyIsAnEmptyArray($property)
    {
        $payload = $this->getScopePayload();
        $scopePayload = $this->arrayGet($payload, $property);
        PHPUnit_Framework_Assert::assertTrue(
            is_array($scopePayload) and $scopePayload === [],
            "Asserting the [$property] property in current scope [{$this->scope}] is an empty array: ".json_encode($payload)
        );
    }
    /**
     * @Given /^the property "([^"]*)" contains "(\d+)" items$/
     */
    public function thePropertyContainsItems($property, $count)
    {
        $payload = $this->getScopePayload();
        PHPUnit_Framework_Assert::assertCount(
            (int) $count,
            $this->arrayGet($payload, $property),
            "Asserting the [$property] property contains [$count] items: ".json_encode($payload)
        );
    }
    /**
     * @Given /^the property "([^"]*)" is an integer$/
     */
    public function thePropertyIsAnInteger($property)
    {
        $payload = $this->getScopePayload();
        PHPUnit_Framework_Assert::isType(
            'int',
            $this->arrayGet($payload, $property),
            "Asserting the [$property] property in current scope [{$this->scope}] is an integer: ".json_encode($payload)
        );
    }
    /**
     * @Given /^the property "([^"]*)" is a string$/
     */
    public function thePropertyIsAString($property)
    {
        $payload = $this->getScopePayload();
        PHPUnit_Framework_Assert::isType(
            'string',
            $this->arrayGet($payload, $property),
            "Asserting the [$property] property in current scope [{$this->scope}] is a string: ".json_encode($payload)
        );
    }
    /**
     * @Given /^the property "([^"]*)" is a string equalling "([^"]*)"$/
     */
    public function thePropertyIsAStringEqualling($property, $expectedValue)
    {
        $payload = $this->getScopePayload();
        $this->thePropertyIsAString($property);
        $actualValue = $this->arrayGet($payload, $property);
        PHPUnit_Framework_Assert::assertSame(
            $actualValue,
            $expectedValue,
            "Asserting the [$property] property in current scope [{$this->scope}] is a string equalling [$expectedValue]."
        );
    }
    /**
     * @Given /^the property "([^"]*)" is a boolean$/
     */
    public function thePropertyIsABoolean($property)
    {
        $payload = $this->getScopePayload();
        PHPUnit_Framework_Assert::assertTrue(
            gettype($this->arrayGet($payload, $property)) == 'boolean',
            "Asserting the [$property] property in current scope [{$this->scope}] is a boolean."
        );
    }
    /**
     * @Given /^the property "([^"]*)" is a boolean equalling "([^"]*)"$/
     */
    public function thePropertyIsABooleanEqualling($property, $expectedValue)
    {
        $payload = $this->getScopePayload();
        $actualValue = $this->arrayGet($payload, $property);
        if (! in_array($expectedValue, ['true', 'false'])) {
            throw new \InvalidArgumentException("Testing for booleans must be represented by [true] or [false].");
        }
        $this->thePropertyIsABoolean($property);
        PHPUnit_Framework_Assert::assertSame(
            $actualValue,
            $expectedValue == 'true',
            "Asserting the [$property] property in current scope [{$this->scope}] is a boolean equalling [$expectedValue]."
        );
    }
    /**
     * @Given /^the property "([^"]*)" is a integer equalling "([^"]*)"$/
     */
    public function thePropertyIsAIntegerEqualling($property, $expectedValue)
    {
        $payload = $this->getScopePayload();
        $actualValue = $this->arrayGet($payload, $property);
        $this->thePropertyIsAnInteger($property);
        PHPUnit_Framework_Assert::assertSame(
            $actualValue,
            (int) $expectedValue,
            "Asserting the [$property] property in current scope [{$this->scope}] is an integer equalling [$expectedValue]."
        );
    }
    /**
     * @Given /^the property "([^"]*)" is either:$/
     */
    public function thePropertyIsEither($property, PyStringNode $options)
    {
        $payload = $this->getScopePayload();
        $actualValue = $this->arrayGet($payload, $property);
        $valid = explode("\n", (string) $options);
        PHPUnit_Framework_Assert::assertTrue(
            in_array($actualValue, $valid),
            sprintf(
                "Asserting the [%s] property in current scope [{$this->scope}] is in array of valid options [%s].",
                $property,
                implode(', ', $valid)
            )
        );
    }
    /**
     * @Given /^the property "([^"]*)" has items where "([^"]*)" is sorted like "([^"]*)"$/
     */
    public function thePropertyHasItemsWhereIsSortedLike($property, $key, $list)
    {
        $payload = $this->getScopePayload();
        $actualArray = $this->arrayGet($payload, $property);
        $actualList = explode(',', $list);
        foreach($actualArray as $index => $item) {
            PHPUnit_Framework_Assert::assertEquals(
                $item->$key,
                $actualList[$index]
            );
        }
    }

    /**
     * @Given /^I scope into the first property "([^"]*)"$/
     */
    public function scopeIntoTheFirstProperty($scope)
    {
        $this->scope = "{$scope}.0";
    }
    /**
     * @Given /^I scope into the property "([^"]*)"$/
     */
    public function scopeIntoTheProperty($scope)
    {
        $this->scope = $scope;
    }

    /**
     * @Given /^the properties exist:$/
     */
    public function thePropertiesExist(PyStringNode $propertiesString)
    {
        foreach (explode("\n", (string) $propertiesString) as $property) {
            $this->thePropertyExists($property);
        }
    }

    /**
     * @Given /^reset scope$/
     */
    public function resetScope()
    {
        $this->scope = null;
    }

    /**
     * Checks the response exists and returns it.
     *
     * @return  Response
     */
    protected function getResponse()
    {
        if (! $this->response) {
            throw new Exception("You must first make a request to check a response.");
        }
        return $this->response;
    }

    /**
     * Return the response payload from the current response.
     *
     * @return  mixed
     */
    protected function getResponsePayload()
    {
        if (! $this->responsePayload) {
            $json = json_decode($this->getResponse()->getBody(true));
            if (json_last_error() !== JSON_ERROR_NONE) {
                $message = 'Failed to decode JSON body ';
                switch (json_last_error()) {
                    case JSON_ERROR_DEPTH:
                        $message .= '(Maximum stack depth exceeded).';
                        break;
                    case JSON_ERROR_STATE_MISMATCH:
                        $message .= '(Underflow or the modes mismatch).';
                        break;
                    case JSON_ERROR_CTRL_CHAR:
                        $message .= '(Unexpected control character found).';
                        break;
                    case JSON_ERROR_SYNTAX:
                        $message .= '(Syntax error, malformed JSON).';
                        break;
                    case JSON_ERROR_UTF8:
                        $message .= '(Malformed UTF-8 characters, possibly incorrectly encoded).';
                        break;
                    default:
                        $message .= '(Unknown error).';
                        break;
                }
                throw new Exception($message);
            }
            $this->responsePayload = $json;
        }
        return $this->responsePayload;
    }

    /**
     * Returns the payload from the current scope within
     * the response.
     *
     * @return mixed
     */
    protected function getScopePayload()
    {
        $payload = $this->getResponsePayload();
        if (! $this->scope) {
            return $payload;
        }
        return $this->arrayGet($payload, $this->scope);
    }

    /**
     * Get an item from an array using "dot" notation.
     *
     * @copyright   Taylor Otwell
     * @link        http://laravel.com/docs/helpers
     * @param       array   $array
     * @param       string  $key
     * @param       mixed   $default
     * @return      mixed
     */
    protected function arrayGet($array, $key)
    {
        if (is_null($key)) {
            return $array;
        }
        // if (isset($array[$key])) {
        //     return $array[$key];
        // }
        foreach (explode('.', $key) as $segment) {
            if (is_object($array)) {
                if (! isset($array->{$segment})) {
                    return;
                }
                $array = $array->{$segment};
            } elseif (is_array($array)) {
                if (! array_key_exists($segment, $array)) {
                    return;
                }
                $array = $array[$segment];
            }
        }
        return $array;
    }

}
