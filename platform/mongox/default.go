package mongox

var defaultClient *Client

func SetDefault(client *Client) { defaultClient = client }
func DB() *Client               { return defaultClient }
