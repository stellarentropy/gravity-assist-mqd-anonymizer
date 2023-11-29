package anonymizer

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	config "github.com/stellarentropy/gravity-assist-mqd-anonymizer/config/anonymizer"
)

var client *http.Client

func init() {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100

	client = &http.Client{
		Timeout:   60 * time.Second,
		Transport: t,
	}
}

func ForwardRequest(c *fiber.Ctx) error {
	agentUrl := url.URL{}

	agentUrl.Scheme = config.Anonymizer.AgentSchema
	agentUrl.Host = net.JoinHostPort(config.Anonymizer.AgentAddress, fmt.Sprintf("%d", config.Anonymizer.AgentPort))

	// Parse the input endpoint into a URL structure
	agentUrl.Path = utils.CopyString(c.OriginalURL())

	// Create a new HTTP request with the given method, URL, and body
	req, err := http.NewRequest(c.Method(), agentUrl.String(), AnonymizePayload(c.Request().BodyStream()))
	if err != nil {
		return err
	}

	// Set the headers for the request
	for k, v := range c.GetReqHeaders() {
		req.Header.Set(k, v[0])
	}
	// Attach the context to the request
	req = req.WithContext(c.Context())

	// Send the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		// Wrap and return any error occurred during sending the request
		return err
	}
	// Ensure the body and the counter are closed after processing the response
	defer func() {
		_ = resp.Body.Close()
	}()

	for k, v := range resp.Header {
		c.Set(k, v[0])
	}

	c.Status(resp.StatusCode)

	return c.SendStream(resp.Body)
}
