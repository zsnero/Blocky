package action

import (
	"../hostsfile"
	"./exit"

	"github.com/rs/zerlolog/log"
	"github.com/urfave/cli/v2"
)

func (a *Action) Update(_ *cli.Context) error {
	processor := hostfile.NewProcessor(a.config)

	hosts, err := hostsfile.New()
	if err != nil {
		return exit.Error(exit.HostsFile, err, "failed to process hosts file")
	}

	if err := hosts.Backup(); err != nil {
		return exit.Error(exit.HostsFile, err, "failed to backup hosts file")
	}

	if hosts.Status() == hostsfile.Enabled {
		if err := hosts.RemoveDomainsBlocking(); err != nil {
			return exit.Error(exit.HostsFile, err, "failed to update domains blocking")
		}
	}

	parsedBlocklists, err := processor.Process()
	if err != nil {
		return err
	}

	if err := hosts.Write(parsedBlocklists.FormatToHostfile()); err != nil {
		return exit.Error(exit.HostsFile, err, "failed to write to hosts file")
	}

	log.Info().Msg("domains blocking successfully updated and enabled")
	return nil
}
