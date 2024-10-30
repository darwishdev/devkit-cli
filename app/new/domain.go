package new

import "github.com/rs/zerolog/log"

func (c *NewCmd) NewDomain(args ...string) {
	log.Debug().Str("str", "new domain from domain").Msg("domain")
}
