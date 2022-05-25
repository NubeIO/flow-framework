package urls

import "fmt"

const FlowNetworkUrl string = "/api/flow_network"
const FlowNetworkCloneUrl string = "/api/flow_network_clones"
const StreamUrl string = "/api/streams"
const StreamCloneUrl string = "/api/stream_clones"
const ProducerUrl string = "/api/producers"
const ConsumerUrl string = "/api/consumers"
const WriterCloneUrl string = "/api/producers/writer_clones"
const WriterUrl string = "/api/consumers/writers"

func SingularUrl(url, uuid string) string {
	return fmt.Sprintf("%s/%s", url, uuid)
}

func SingularUrlByArg(url, name, value string) string {
	return fmt.Sprintf("%s/one/args?%s=%s", url, name, value)
}

func PluralUrlByArg(url, name, value string) string {
	return fmt.Sprintf("%s?%s=%s", url, name, value)
}
