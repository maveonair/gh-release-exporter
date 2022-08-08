package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/maveonair/gh-release-exporter/internal/config"
	"github.com/maveonair/gh-release-exporter/internal/github"
	"github.com/maveonair/gh-release-exporter/internal/metrics"
	"github.com/maveonair/gh-release-exporter/internal/release"

	"github.com/Masterminds/semver"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	log.SetOutput(os.Stdout)
}

func main() {
	configFilePathPtr := flag.String("config", "", "config file")
	flag.Parse()

	if *configFilePathPtr == "" {
		log.Fatal("argument -config is not set")
	}

	config, err := config.NewConfig(*configFilePathPtr)
	if err != nil {
		log.WithError(err).Fatal()
	}

	go metrics.Setup(config.ListeningAddr)

	for {
		checkReleasesForUpdate(config.Releases)

		log.WithField("interval", config.Interval).Info("Sleep until next update")
		time.Sleep(config.Interval)
	}
}

func checkReleasesForUpdate(releases map[string]release.Release) {
	for key, release := range releases {
		log.WithFields(log.Fields{
			"name":               key,
			"last_known_version": release.LastKnownVersion,
		}).Info("Check latest release")

		c, err := semver.NewConstraint(fmt.Sprintf("> %s", release.LastKnownVersion))
		if err != nil {
			metrics.IncreaseErrors()

			log.WithFields(log.Fields{
				"name":               key,
				"last_known_version": release.LastKnownVersion,
			}).Error(err)

			continue
		}

		tags, err := github.GetLatestTags(release.GitHubRepo)
		if err != nil {
			metrics.IncreaseErrors()

			log.WithFields(log.Fields{
				"name":               key,
				"last_known_version": release.LastKnownVersion,
			}).Error(err)

			continue
		}

		newVersion := ""
		for _, version := range tags {
			v, err := semver.NewVersion(version)
			if err != nil {
				log.WithFields(log.Fields{
					"name":               key,
					"last_known_version": release.LastKnownVersion,
					"version":            version,
				}).Debug(err)

				continue
			}

			if c.Check(v) {
				newVersion = version
			}
		}

		if newVersion != "" {
			metrics.SetReleaseSuccessProbe(key, newVersion, 0)
		} else {
			metrics.SetReleaseSuccessProbe(key, release.LastKnownVersion, 1)
		}
	}
}
